package test

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"opennamu/route/tool"
)

type wsgi_identity_case struct {
	Remote_addr string            `json:"remote_addr"`
	Headers     map[string]string `json:"headers"`
}

type wsgi_probe_result struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

type wsgi_backend_identity struct {
	Identity          string              `json:"identity"`
	Owner_permission  bool                `json:"owner_permission"`
	Admin_permission  bool                `json:"admin_permission"`
	Forwarded_headers map[string][]string `json:"forwarded_headers"`
}

func Run_python_wsgi_cases(t *testing.T, python_path string, server_host string, server_port int, cases []wsgi_identity_case) []wsgi_probe_result {
	t.Helper()
	_, current_file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate WSGI integration test source")
	}
	helper_path := filepath.Join(filepath.Dir(current_file), "wsgi_ip_bridge_probe.py")
	input_data, err := json.Marshal(cases)
	if err != nil {
		t.Fatalf("encode WSGI probe cases: %v", err)
	}
	command := exec.Command(python_path, helper_path, server_host, strconv.Itoa(server_port))
	command.Stdin = bytes.NewReader(input_data)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("run isolated WSGI bridge probe: %v\n%s", err, output)
	}
	results := []wsgi_probe_result{}
	if err := json.Unmarshal(output, &results); err != nil {
		t.Fatalf("decode WSGI probe output %q: %v", output, err)
	}
	if len(results) != len(cases) {
		t.Fatalf("WSGI result count = %d, want %d", len(results), len(cases))
	}
	return results
}

func Test_SEC04_wsgi_to_go_authorization_identity_boundary(t *testing.T) {
	python_path, err := exec.LookPath("python3")
	if err != nil {
		python_path, err = exec.LookPath("python")
	}
	if err != nil {
		t.Skip("Python is unavailable; run test/wsgi_ip_bridge_probe.py with the WSGI integration harness")
	}

	for _, backend := range []string{"sqlite", "turso"} {
		t.Run(backend, func(t *testing.T) {
			db := audit_setup(t, backend)
			audit_seed_auth(db, t)

			const (
				owner_identity = "203.0.113.42"
				admin_identity = "203.0.113.43"
			)
			grant(t, db, owner_identity, "owner")
			grant(t, db, admin_identity, "admin")
			loadSetting(t, db, "load_ip_select", "default")
			loadSetting(t, db, "trusted_proxies", "127.0.0.1")
			tool.Get_ip_select(db)
			tool.Get_trusted_proxies(db)

			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set(sessions.DefaultKey, &emptySession{})
				c.Next()
			})
			var backend_requests atomic.Int64
			router.GET("/_test/wsgi_identity", func(c *gin.Context) {
				backend_requests.Add(1)
				identity := tool.Get_IP(c)
				c.JSON(http.StatusOK, wsgi_backend_identity{
					Identity:         identity,
					Owner_permission: tool.Check_permission(db, "owner", identity),
					Admin_permission: tool.Check_permission(db, "admin", identity),
					Forwarded_headers: map[string][]string{
						"X-Real-IP":        c.Request.Header.Values("X-Real-IP"),
						"CF-Connecting-IP": c.Request.Header.Values("CF-Connecting-IP"),
						"X-Forwarded-For":  c.Request.Header.Values("X-Forwarded-For"),
						"Forwarded":        c.Request.Header.Values("Forwarded"),
					},
				})
			})
			server := httptest.NewServer(router)
			defer server.Close()
			server_host, server_port_string, err := net.SplitHostPort(strings.TrimPrefix(server.URL, "http://"))
			if err != nil {
				t.Fatalf("parse local Go server address %q: %v", server.URL, err)
			}
			server_port, err := strconv.Atoi(server_port_string)
			if err != nil {
				t.Fatalf("parse local Go server port %q: %v", server_port_string, err)
			}

			// The integration client runs app.py's real WSGI application. The Go test
			// backend sees an actual loopback TCP peer and trusts exactly that peer.
			// A reverse-proxy case represents WSGI REMOTE_ADDR after the ingress has
			// validated and canonicalized the client address.
			default_cases := []wsgi_identity_case{
				{
					Remote_addr: "198.51.100.77",
					Headers: map[string]string{
						"X-Real-IP":        owner_identity,
						"CF-Connecting-IP": admin_identity,
						"X-Forwarded-For":  admin_identity,
						"Forwarded":        "for=" + admin_identity,
					},
				},
				{
					Remote_addr: "192.0.2.9",
					Headers:     map[string]string{"X-Real-IP": admin_identity},
				},
				{
					Remote_addr: owner_identity,
					Headers: map[string]string{
						"X-Real-IP":        admin_identity,
						"CF-Connecting-IP": admin_identity,
						"X-Forwarded-For":  admin_identity,
					},
				},
				{
					Remote_addr: admin_identity,
					Headers:     map[string]string{"X-Real-IP": owner_identity},
				},
				{
					Remote_addr: "2001:db8::77",
					Headers:     map[string]string{"X-Real-IP": owner_identity},
				},
			}
			results := Run_python_wsgi_cases(t, python_path, server_host, server_port, default_cases)
			want_identities := []string{"198.51.100.77", "192.0.2.9", owner_identity, admin_identity, "2001:db8::77"}
			for index, result := range results {
				if result.Status != "200 OK" {
					t.Fatalf("default case %d WSGI status = %q, want 200 OK: %s", index, result.Status, result.Body)
				}
				backend_result := wsgi_backend_identity{}
				if err := json.Unmarshal([]byte(result.Body), &backend_result); err != nil {
					t.Fatalf("decode default case %d backend result %q: %v", index, result.Body, err)
				}
				if backend_result.Identity != want_identities[index] {
					t.Fatalf("default case %d resolved identity = %q, want %q", index, backend_result.Identity, want_identities[index])
				}
				want_owner := want_identities[index] == owner_identity
				want_admin := want_owner || want_identities[index] == admin_identity
				if backend_result.Owner_permission != want_owner || backend_result.Admin_permission != want_admin {
					t.Fatalf("default case %d permissions = owner:%v admin:%v, want owner:%v admin:%v", index, backend_result.Owner_permission, backend_result.Admin_permission, want_owner, want_admin)
				}
				for _, header_name := range []string{"X-Real-IP", "CF-Connecting-IP", "X-Forwarded-For"} {
					values := backend_result.Forwarded_headers[header_name]
					if len(values) != 1 || values[0] != want_identities[index] {
						t.Fatalf("default case %d backend %s values = %v, want one sanitized value %q", index, header_name, values, want_identities[index])
					}
				}
				if values := backend_result.Forwarded_headers["Forwarded"]; len(values) != 0 {
					t.Fatalf("default case %d untrusted Forwarded header reached backend: %v", index, values)
				}
			}

			// A missing, malformed, or host:port WSGI REMOTE_ADDR is rejected before
			// app.py reads/forwards the request body or contacts the Go backend.
			backend_count := backend_requests.Load()
			invalid_cases := []wsgi_identity_case{
				{Remote_addr: "", Headers: map[string]string{"X-Real-IP": owner_identity}},
				{Remote_addr: "not-an-ip", Headers: map[string]string{"X-Real-IP": owner_identity}},
				{Remote_addr: "198.51.100.7:8080", Headers: map[string]string{"X-Real-IP": owner_identity}},
			}
			invalid_results := Run_python_wsgi_cases(t, python_path, server_host, server_port, invalid_cases)
			for index, result := range invalid_results {
				if result.Status != "400 Bad Request" {
					t.Fatalf("invalid REMOTE_ADDR case %d status = %q, want 400 Bad Request", index, result.Status)
				}
			}
			if got := backend_requests.Load(); got != backend_count {
				t.Fatalf("invalid WSGI addresses reached backend: requests before=%d after=%d", backend_count, got)
			}

			// Verify supported non-default Go header modes also receive only the
			// WSGI server's address, never the client's raw XFF/CF header.
			for _, mode_case := range []struct {
				mode      string
				header    string
				client_ip string
			}{
				{mode: "HTTP_X_FORWARDED_FOR", header: "X-Forwarded-For", client_ip: "198.51.100.88"},
				{mode: "HTTP_CF_CONNECTING_IP", header: "CF-Connecting-IP", client_ip: "198.51.100.89"},
			} {
				loadSetting(t, db, "load_ip_select", mode_case.mode)
				tool.Get_ip_select(db)
				mode_results := Run_python_wsgi_cases(t, python_path, server_host, server_port, []wsgi_identity_case{{
					Remote_addr: mode_case.client_ip,
					Headers:     map[string]string{mode_case.header: owner_identity},
				}})
				if len(mode_results) != 1 || mode_results[0].Status != "200 OK" {
					t.Fatalf("mode %s WSGI result = %+v, want HTTP 200", mode_case.mode, mode_results)
				}
				mode_backend_result := wsgi_backend_identity{}
				if err := json.Unmarshal([]byte(mode_results[0].Body), &mode_backend_result); err != nil {
					t.Fatalf("decode mode %s result: %v", mode_case.mode, err)
				}
				if mode_backend_result.Identity != mode_case.client_ip || mode_backend_result.Owner_permission || mode_backend_result.Admin_permission {
					t.Fatalf("mode %s resolved %q owner:%v admin:%v; want unprivileged %s", mode_case.mode, mode_backend_result.Identity, mode_backend_result.Owner_permission, mode_backend_result.Admin_permission, mode_case.client_ip)
				}
			}

			// Missing and malformed trust lists keep Go's established peer-only
			// fallback. The WSGI-supplied forged owner header cannot override it.
			for _, trusted_proxy_setting := range []string{"", "127.0.0.1,not-an-ip"} {
				loadSetting(t, db, "load_ip_select", "default")
				loadSetting(t, db, "trusted_proxies", trusted_proxy_setting)
				tool.Get_ip_select(db)
				tool.Get_trusted_proxies(db)
				fallback_results := Run_python_wsgi_cases(t, python_path, server_host, server_port, []wsgi_identity_case{{
					Remote_addr: "198.51.100.90",
					Headers:     map[string]string{"X-Real-IP": owner_identity},
				}})
				fallback_result := wsgi_backend_identity{}
				if len(fallback_results) != 1 || fallback_results[0].Status != "200 OK" || json.Unmarshal([]byte(fallback_results[0].Body), &fallback_result) != nil {
					t.Fatalf("trust setting %q WSGI fallback response = %+v", trusted_proxy_setting, fallback_results)
				}
				if fallback_result.Identity != "127.0.0.1" || fallback_result.Owner_permission || fallback_result.Admin_permission {
					t.Fatalf("trust setting %q accepted forged identity: %+v", trusted_proxy_setting, fallback_result)
				}
			}
		})
	}
}

func Test_SEC04_direct_go_client_cannot_spoof_owner_or_admin(t *testing.T) {
	db := audit_setup(t, "sqlite")
	audit_seed_auth(db, t)
	grant(t, db, "203.0.113.42", "owner")
	grant(t, db, "203.0.113.43", "admin")
	loadSetting(t, db, "trusted_proxies", "127.0.0.1")

	for _, mode_case := range []struct {
		mode   string
		header string
	}{
		{mode: "default", header: "X-Real-IP"},
		{mode: "HTTP_X_FORWARDED_FOR", header: "X-Forwarded-For"},
		{mode: "HTTP_CF_CONNECTING_IP", header: "CF-Connecting-IP"},
	} {
		loadSetting(t, db, "load_ip_select", mode_case.mode)
		tool.Get_ip_select(db)
		tool.Get_trusted_proxies(db)
		for _, spoofed_identity := range []string{"203.0.113.42", "203.0.113.43"} {
			request_context := ipContext("198.51.100.77:443", map[string]string{mode_case.header: spoofed_identity}, "")
			identity := tool.Get_IP(request_context)
			if identity != "198.51.100.77" {
				t.Fatalf("direct untrusted peer in mode %s resolved to %q, want peer IP", mode_case.mode, identity)
			}
			if tool.Check_permission(db, "owner", identity) || tool.Check_permission(db, "admin", identity) {
				t.Fatalf("direct untrusted peer acquired owner/admin by spoofing %s in mode %s", spoofed_identity, mode_case.mode)
			}
		}
	}
}
