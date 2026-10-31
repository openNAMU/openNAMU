package test

// trusted_proxy_ip_security_test.go
//
// 25-item regression matrix for the trusted-proxy IP identity hardening.
//
// Proves that the *identity used for authorization* (ban, permission, document
// ACL) is the resolved TCP peer when the peer is not in trusted_proxies, not a
// forged X-Real-IP / CF-Connecting-IP header. Covers:
//
//   - untrusted peer -> header ignored (default / X-Real-IP / CF-Connecting-IP)
//   - trusted peer -> header honored
//   - REMOTE_ADDR mode always ignores headers
//   - comma-separated / CIDR (IPv4, IPv6, IPv4-mapped-IPv6) trust
//   - empty & catch-all (0.0.0.0/0, ::/0) fail-safe (peer only)
//   - invalid RemoteAddr -> fail-closed (no ::1 fallback)
//   - settings API rejects invalid values, persists valid ones
//   - Is_https honors X-Forwarded-Proto only from trusted peers
//   - authenticated session identity is preserved
//
// All tests run sequentially against isolated SQLite DBs (no t.Parallel).

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
	"opennamu/route"
	"opennamu/route/tool"
)

// ---------- gin / session helpers ----------

// fakeSession implements sessions.Session for the id field read by Get_IP.
type fakeSession struct {
	id   string
	data map[string]interface{}
}

func (f *fakeSession) ID() string { return "" }
func (f *fakeSession) Get(key interface{}) interface{} {
	k, ok := key.(string)
	if !ok {
		return nil
	}
	if f.data != nil {
		if v, p := f.data[k]; p {
			return v
		}
	}
	if k == "id" {
		return f.id
	}
	return nil
}
func (f *fakeSession) Set(key interface{}, val interface{}) {
	k, ok := key.(string)
	if !ok {
		return
	}
	if f.data == nil {
		f.data = map[string]interface{}{}
	}
	f.data[k] = val
}
func (f *fakeSession) Delete(key interface{})                 {}
func (f *fakeSession) Clear()                                 { f.data = nil }
func (f *fakeSession) AddFlash(v interface{}, vars ...string) {}
func (f *fakeSession) Flashes(vars ...string) []interface{}   { return nil }
func (f *fakeSession) Options(o sessions.Options)             {}
func (f *fakeSession) Save() error                            { return nil }

// ipContext builds a gin.Context with the given RemoteAddr and headers.
// authID is the session id ("" means no authenticated identity).
func ipContext(remote string, headers map[string]string, authID string) *gin.Context {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = &http.Request{
		Method:     "GET",
		URL:        &url.URL{Path: "/"},
		RemoteAddr: remote,
		Header:     headerFromMap(headers),
	}
	// Get_IP -> Get_session -> sessions.Default(c) panics if this key is unset,
	// so every context carries a (possibly empty) session.
	c.Set(sessions.DefaultKey, &fakeSession{id: authID})
	return c
}

// headerFromMap converts a map[string]string into an http.Header.
func headerFromMap(m map[string]string) http.Header {
	h := http.Header{}
	for k, v := range m {
		h.Add(k, v)
	}
	return h
}

// ---------- isolated DB helpers ----------

func new_ip_db(t *testing.T) *sql.DB {
	t.Helper()
	db := new_test_db(t, filepath.Join(t.TempDir(), "ipsecurity.db"))
	// Check_acl -> Check_document_filter_acl reads this table; keep it present
	// with the columns the SELECT expects (kind, plus, plus_t).
	if _, err := db.Exec(`create table if not exists html_filter (kind, plus, plus_t, coverage)`); err != nil {
		t.Fatalf("create html_filter: %v", err)
	}
	return db
}

func loadSetting(t *testing.T, db *sql.DB, name, value string) {
	t.Helper()
	if _, err := db.Exec(`delete from other where name = ?`, name); err != nil {
		t.Fatalf("delete other %s: %v", name, err)
	}
	if _, err := db.Exec(`insert into other (name,data,coverage) values (?,?,?)`, name, value, ""); err != nil {
		t.Fatalf("insert other %s: %v", name, err)
	}
}

// seedGroups inserts the base auth groups referenced by Get_auth_group_info /
// Check_auth, plus any extra groups the test needs.
func seedGroups(t *testing.T, db *sql.DB, groups ...string) {
	t.Helper()
	for _, g := range []string{"owner", "admin", "user", "ip", "ban", "ban_without_site", "ban_without_login"} {
		if _, err := db.Exec(`delete from alist where name = ?`, g); err != nil {
			t.Fatalf("delete alist %s: %v", g, err)
		}
	}
	if _, err := db.Exec(`insert into alist (name,acl) values ('ban_without_site','nothing')`); err != nil {
		t.Fatalf("insert alist ban_without_site: %v", err)
	}
	for _, g := range []string{"owner", "admin", "user", "ip", "ban", "ban_without_login"} {
		if _, err := db.Exec(`insert into alist (name,acl) values (?,?)`, g, g); err != nil {
			t.Fatalf("insert alist %s: %v", g, err)
		}
	}
	for _, g := range groups {
		if _, err := db.Exec(`delete from alist where name = ?`, g); err != nil {
			t.Fatalf("delete alist %s: %v", g, err)
		}
		if _, err := db.Exec(`insert into alist (name,acl) values (?,?)`, g, g); err != nil {
			t.Fatalf("insert alist %s: %v", g, err)
		}
	}
}

// grant assigns an auth group to an identity via the acl user_set row.
func grant(t *testing.T, db *sql.DB, id, group string) {
	t.Helper()
	if _, err := db.Exec(`delete from user_set where id = ? and name = 'acl'`, id); err != nil {
		t.Fatalf("delete user_set %s: %v", id, err)
	}
	if _, err := db.Exec(`insert into user_set (name,id,data) values ('acl',?,?)`, id, group); err != nil {
		t.Fatalf("insert user_set %s: %v", id, err)
	}
}

// ---------- matrix constants ----------

const (
	forgedIP = "198.51.100.10"
	actualIP = "203.0.113.20"
	peerIPv6 = "2001:db8:85a3::8a2e:370:7334"
	trustV6  = "2001:db8:85a3:0:0:0:0:0/48"
)

// resolveIP refreshes the global trust state from db and returns the identity
// Get_IP resolves for a request with the given peer + headers.
func resolveIP(t *testing.T, db *sql.DB, mode, trusted string, remote string, headers map[string]string, authID string) string {
	t.Helper()
	tool.Get_ip_select(db)
	tool.Get_trusted_proxies(db)
	return tool.Get_IP(ipContext(remote, headers, authID))
}

// authDB seeds groups + settings and returns (db, resolvedIdentity).
func authDB(t *testing.T, mode, trusted string, remote string, headers map[string]string, groups ...string) (*sql.DB, string) {
	t.Helper()
	db := new_ip_db(t)
	seedGroups(t, db, groups...)
	loadSetting(t, db, "load_ip_select", mode)
	loadSetting(t, db, "trusted_proxies", trusted)
	tool.Get_ip_select(db)
	tool.Get_trusted_proxies(db)
	return db, tool.Get_IP(ipContext(remote, headers, ""))
}

// ---------- 1-11: basic header trust ----------

func Test_TRP_01_default_untrusted_header_ignored(t *testing.T) {
	db := new_ip_db(t)
	loadSetting(t, db, "load_ip_select", "default")
	loadSetting(t, db, "trusted_proxies", "")
	ip := resolveIP(t, db, "default", "", actualIP, map[string]string{"X-Real-IP": forgedIP}, "")
	if ip != actualIP {
		t.Fatalf("default untrusted: got %q, want %q", ip, actualIP)
	}
}

func Test_TRP_02_untrusted_cf_header_ignored(t *testing.T) {
	db := new_ip_db(t)
	loadSetting(t, db, "load_ip_select", "default")
	loadSetting(t, db, "trusted_proxies", "")
	ip := resolveIP(t, db, "default", "", actualIP, map[string]string{"CF-Connecting-IP": forgedIP}, "")
	if ip != actualIP {
		t.Fatalf("default cf untrusted: got %q, want %q", ip, actualIP)
	}
}

func Test_TRP_03_remote_addr_ignores_headers(t *testing.T) {
	db := new_ip_db(t)
	loadSetting(t, db, "load_ip_select", "REMOTE_ADDR")
	loadSetting(t, db, "trusted_proxies", actualIP)
	ip := resolveIP(t, db, "REMOTE_ADDR", actualIP, actualIP, map[string]string{"X-Real-IP": forgedIP}, "")
	if ip != actualIP {
		t.Fatalf("REMOTE_ADDR: got %q, want %q", ip, actualIP)
	}
}

func Test_TRP_04_default_no_header(t *testing.T) {
	db := new_ip_db(t)
	loadSetting(t, db, "load_ip_select", "default")
	loadSetting(t, db, "trusted_proxies", "")
	ip := resolveIP(t, db, "default", "", actualIP, nil, "")
	if ip != actualIP {
		t.Fatalf("default no header: got %q, want %q", ip, actualIP)
	}
}

func Test_TRP_05_real_ip_trusted_honored(t *testing.T) {
	db := new_ip_db(t)
	loadSetting(t, db, "load_ip_select", "HTTP_X_REAL_IP")
	loadSetting(t, db, "trusted_proxies", actualIP)
	ip := resolveIP(t, db, "HTTP_X_REAL_IP", actualIP, actualIP, map[string]string{"X-Real-IP": forgedIP}, "")
	if ip != forgedIP {
		t.Fatalf("X-Real-IP trusted: got %q, want %q", ip, forgedIP)
	}
}

func Test_TRP_06_real_ip_trusted_no_header(t *testing.T) {
	db := new_ip_db(t)
	loadSetting(t, db, "load_ip_select", "HTTP_X_REAL_IP")
	loadSetting(t, db, "trusted_proxies", actualIP)
	ip := resolveIP(t, db, "HTTP_X_REAL_IP", actualIP, actualIP, nil, "")
	if ip != actualIP {
		t.Fatalf("X-Real-IP trusted no header: got %q, want %q", ip, actualIP)
	}
}

func Test_TRP_07_real_ip_untrusted(t *testing.T) {
	db := new_ip_db(t)
	loadSetting(t, db, "load_ip_select", "HTTP_X_REAL_IP")
	loadSetting(t, db, "trusted_proxies", "")
	ip := resolveIP(t, db, "HTTP_X_REAL_IP", "", actualIP, map[string]string{"X-Real-IP": forgedIP}, "")
	if ip != actualIP {
		t.Fatalf("X-Real-IP untrusted: got %q, want %q", ip, actualIP)
	}
}

func Test_TRP_08_cf_trusted_honored(t *testing.T) {
	db := new_ip_db(t)
	loadSetting(t, db, "load_ip_select", "HTTP_CF_CONNECTING_IP")
	loadSetting(t, db, "trusted_proxies", actualIP)
	ip := resolveIP(t, db, "HTTP_CF_CONNECTING_IP", actualIP, actualIP, map[string]string{"CF-Connecting-IP": forgedIP}, "")
	if ip != forgedIP {
		t.Fatalf("CF-Connecting-IP trusted: got %q, want %q", ip, forgedIP)
	}
}

func Test_TRP_09_cf_untrusted(t *testing.T) {
	db := new_ip_db(t)
	loadSetting(t, db, "load_ip_select", "HTTP_CF_CONNECTING_IP")
	loadSetting(t, db, "trusted_proxies", "")
	ip := resolveIP(t, db, "HTTP_CF_CONNECTING_IP", "", actualIP, map[string]string{"CF-Connecting-IP": forgedIP}, "")
	if ip != actualIP {
		t.Fatalf("CF-Connecting-IP untrusted: got %q, want %q", ip, actualIP)
	}
}

func Test_TRP_10_default_trusted_honored(t *testing.T) {
	db := new_ip_db(t)
	loadSetting(t, db, "load_ip_select", "default")
	loadSetting(t, db, "trusted_proxies", actualIP)
	ip := resolveIP(t, db, "default", actualIP, actualIP, map[string]string{"X-Real-IP": forgedIP}, "")
	if ip != forgedIP {
		t.Fatalf("default trusted: got %q, want %q", ip, forgedIP)
	}
}

func Test_TRP_11_default_trusted_no_header(t *testing.T) {
	db := new_ip_db(t)
	loadSetting(t, db, "load_ip_select", "default")
	loadSetting(t, db, "trusted_proxies", actualIP)
	ip := resolveIP(t, db, "default", actualIP, actualIP, nil, "")
	if ip != actualIP {
		t.Fatalf("default trusted no header: got %q, want %q", ip, actualIP)
	}
}

// ---------- 12-13: lists and fail-closed ----------

func Test_TRP_12_comma_separated(t *testing.T) {
	db := new_ip_db(t)
	loadSetting(t, db, "load_ip_select", "default")
	loadSetting(t, db, "trusted_proxies", actualIP+",203.0.113.21")
	ip := resolveIP(t, db, "default", actualIP+",203.0.113.21", actualIP, map[string]string{"X-Real-IP": forgedIP}, "")
	if ip != forgedIP {
		t.Fatalf("comma trusted: got %q, want %q", ip, forgedIP)
	}
	ip = resolveIP(t, db, "default", actualIP+",203.0.113.21", "203.0.113.99", map[string]string{"X-Real-IP": forgedIP}, "")
	if ip != "203.0.113.99" {
		t.Fatalf("comma untrusted: got %q, want %q", ip, "203.0.113.99")
	}
}

func Test_TRP_13_invalid_remote_addr_fail_closed(t *testing.T) {
	db := new_ip_db(t)
	loadSetting(t, db, "load_ip_select", "default")
	loadSetting(t, db, "trusted_proxies", "")
	ip := resolveIP(t, db, "default", "", "not-an-ip", nil, "")
	if ip != "invalid_remote_address" {
		t.Fatalf("invalid remote addr: got %q, want invalid_remote_address", ip)
	}
}

// ---------- 14-17: CIDR / catch-all ----------

func Test_TRP_14_ipv4_cidr_trust(t *testing.T) {
	db := new_ip_db(t)
	loadSetting(t, db, "load_ip_select", "default")
	loadSetting(t, db, "trusted_proxies", "203.0.113.0/24")
	ip := resolveIP(t, db, "default", "203.0.113.0/24", "203.0.113.20", map[string]string{"X-Real-IP": forgedIP}, "")
	if ip != forgedIP {
		t.Fatalf("ipv4 cidr trusted: got %q, want %q", ip, forgedIP)
	}
	ip = resolveIP(t, db, "default", "203.0.113.0/24", "203.0.114.5", map[string]string{"X-Real-IP": forgedIP}, "")
	if ip != "203.0.114.5" {
		t.Fatalf("ipv4 cidr untrusted: got %q, want %q", ip, "203.0.114.5")
	}
}

func Test_TRP_15_ipv6_cidr_trust(t *testing.T) {
	db := new_ip_db(t)
	loadSetting(t, db, "load_ip_select", "default")
	loadSetting(t, db, "trusted_proxies", trustV6)
	ip := resolveIP(t, db, "default", trustV6, peerIPv6, map[string]string{"X-Real-IP": forgedIP}, "")
	if ip != forgedIP {
		t.Fatalf("ipv6 cidr trusted: got %q, want %q", ip, forgedIP)
	}
}

func Test_TRP_16_ipv4_mapped_v6_vs_v4_prefix(t *testing.T) {
	db := new_ip_db(t)
	loadSetting(t, db, "load_ip_select", "default")
	loadSetting(t, db, "trusted_proxies", "203.0.113.0/24")
	mapped := "::ffff:203.0.113.20"
	ip := resolveIP(t, db, "default", "203.0.113.0/24", mapped, map[string]string{"X-Real-IP": forgedIP}, "")
	if ip != forgedIP {
		t.Fatalf("v4-mapped v6 vs v4 prefix: got %q, want %q", ip, forgedIP)
	}
}

func Test_TRP_17_empty_trusted(t *testing.T) {
	db := new_ip_db(t)
	loadSetting(t, db, "load_ip_select", "default")
	loadSetting(t, db, "trusted_proxies", "")
	ip := resolveIP(t, db, "default", "", actualIP, map[string]string{"X-Real-IP": forgedIP}, "")
	if ip != actualIP {
		t.Fatalf("empty trusted: got %q, want %q (caught)", ip, actualIP)
	}
}

func Test_TRP_18_reject_catch_all_0_0_0_0(t *testing.T) {
	db := new_ip_db(t)
	loadSetting(t, db, "load_ip_select", "default")
	loadSetting(t, db, "trusted_proxies", "0.0.0.0/0")
	ip := resolveIP(t, db, "default", "0.0.0.0/0", actualIP, map[string]string{"X-Real-IP": forgedIP}, "")
	if ip != actualIP {
		t.Fatalf("catch-all 0.0.0.0/0: got %q, want %q (caught)", ip, actualIP)
	}
}

func Test_TRP_19_reject_catch_all_double_colon(t *testing.T) {
	db := new_ip_db(t)
	loadSetting(t, db, "load_ip_select", "default")
	loadSetting(t, db, "trusted_proxies", "::/0")
	ip := resolveIP(t, db, "default", "::/0", actualIP, map[string]string{"X-Real-IP": forgedIP}, "")
	if ip != actualIP {
		t.Fatalf("catch-all ::/0: got %q, want %q (caught)", ip, actualIP)
	}
}

func Test_TRP_20_reject_mixed_values(t *testing.T) {
	db := new_ip_db(t)
	loadSetting(t, db, "load_ip_select", "default")
	loadSetting(t, db, "trusted_proxies", actualIP+",not-an-ip")
	ip := resolveIP(t, db, "default", actualIP+",not-an-ip", actualIP, map[string]string{"X-Real-IP": forgedIP}, "")
	if ip != actualIP {
		t.Fatalf("mixed valid+invalid: got %q, want %q (caught)", ip, actualIP)
	}
	if err := tool.Validate_trusted_proxies("203.0.113.0/24,not-an-ip"); err == nil {
		t.Fatalf("Validate_trusted_proxies accepted mixed value")
	}
}

// ---------- 21: authenticated session ----------

func Test_TRP_21_authenticated_session_bypasses(t *testing.T) {
	db := new_ip_db(t)
	loadSetting(t, db, "load_ip_select", "default")
	loadSetting(t, db, "trusted_proxies", "")
	ip := resolveIP(t, db, "default", "", actualIP, map[string]string{"X-Real-IP": forgedIP}, "alice")
	if ip != "alice" {
		t.Fatalf("authenticated: got %q, want alice", ip)
	}
}

// ---------- 22-24: authorization decisions ----------

// 22. ban: actual peer banned via 'ban_without_site'; forges an unbanned IP.
func Test_TRP_22_ban_identity_decision(t *testing.T) {
	db := new_ip_db(t)
	seedGroups(t, db)
	loadSetting(t, db, "load_ip_select", "default")
	loadSetting(t, db, "trusted_proxies", "")
	grant(t, db, actualIP, "ban_without_site")
	tool.Get_ip_select(db)
	tool.Get_trusted_proxies(db)
	ip := tool.Get_IP(ipContext(actualIP, map[string]string{"X-Real-IP": forgedIP}, ""))
	if ip != actualIP {
		t.Fatalf("ban: resolved identity got %q, want %q", ip, actualIP)
	}
	if tool.Check_permission(db, "site_view", ip) {
		t.Fatalf("banned actual peer must be denied")
	}
	if !tool.Check_permission(db, "site_view", forgedIP) {
		t.Fatalf("forged unbanned IP should be allowed")
	}
}

// 23. permission: forged IP granted admin; actual peer must not inherit it.
func Test_TRP_23_permission_identity_decision(t *testing.T) {
	db := new_ip_db(t)
	seedGroups(t, db)
	loadSetting(t, db, "load_ip_select", "default")
	loadSetting(t, db, "trusted_proxies", "")
	grant(t, db, forgedIP, "admin")
	tool.Get_ip_select(db)
	tool.Get_trusted_proxies(db)
	ip := tool.Get_IP(ipContext(actualIP, map[string]string{"X-Real-IP": forgedIP}, ""))
	if ip != actualIP {
		t.Fatalf("permission: resolved identity got %q, want %q", ip, actualIP)
	}
	if tool.Check_permission(db, "admin", ip) {
		t.Fatalf("actual peer must not inherit forged admin")
	}
	if !tool.Check_permission(db, "admin", forgedIP) {
		t.Fatalf("forged admin IP should be allowed")
	}
}

// 24. document ACL: document restricted to 'restricted_group'; actual peer denied,
// forged (in-group) IP would be allowed.
func Test_TRP_24_document_acl_identity_decision(t *testing.T) {
	doc := "secret_doc"
	db := new_ip_db(t)
	seedGroups(t, db, "restricted_group")
	// restricted_group must also grant site_view (the ip base) so an
	// in-group identity is even considered for the document ACL.
	if _, err := db.Exec(`delete from alist where name = 'restricted_group'`); err != nil {
		t.Fatalf("delete restricted_group: %v", err)
	}
	if _, err := db.Exec(`insert into alist (name,acl) values ('restricted_group','ip')`); err != nil {
		t.Fatalf("insert restricted_group: %v", err)
	}
	loadSetting(t, db, "load_ip_select", "default")
	loadSetting(t, db, "trusted_proxies", "")
	if _, err := db.Exec(`delete from acl`); err != nil {
		t.Fatalf("delete acl: %v", err)
	}
	if _, err := db.Exec(`insert into acl (title,data,type) values (?,?,?)`, doc, "restricted_group", "view"); err != nil {
		t.Fatalf("insert acl doc: %v", err)
	}
	grant(t, db, forgedIP, "restricted_group")
	tool.Get_ip_select(db)
	tool.Get_trusted_proxies(db)
	ip := tool.Get_IP(ipContext(actualIP, map[string]string{"X-Real-IP": forgedIP}, ""))
	if ip != actualIP {
		t.Fatalf("doc acl: resolved identity got %q, want %q", ip, actualIP)
	}
	if tool.Check_acl(db, doc, "", "render", ip) {
		t.Fatalf("actual peer must be denied doc ACL")
	}
	if !tool.Check_acl(db, doc, "", "render", forgedIP) {
		t.Fatalf("forged IP in group should be allowed")
	}
}

// ---------- 25: settings API ----------

// routeDB points the global db_set at an isolated temp DB via the production
// DB_boot/DB_init path. Not parallel.
func routeDB(t *testing.T) (dir string) {
	t.Helper()
	dir = filepath.Join(t.TempDir(), "db")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	os.Setenv("NAMU_DB", dir)
	os.Setenv("NAMU_DB_TYPE", "sqlite")
	defer os.Unsetenv("NAMU_DB")
	defer os.Unsetenv("NAMU_DB_TYPE")
	tool.DB_init()
	return
}

// firstSetting reads a setting row from a route DB.
func firstSetting(db *sql.DB, name string) string {
	var v string
	row := db.QueryRow(`select data from other where name = ?`, name)
	if row.Scan(&v) == nil {
		return v
	}
	return ""
}

func Test_TRP_25_settings_api_validation(t *testing.T) {
	_ = routeDB(t)
	db := tool.DB_connect()
	defer tool.DB_close(db)

	// Seed the owner group so the test identity can edit main settings.
	if _, err := db.Exec(`insert into alist (name,acl) values ('owner','owner')`); err != nil {
		t.Fatalf("insert owner: %v", err)
	}
	if _, err := db.Exec(`insert into user_set (name,id,data) values ('acl','admin-9','owner')`); err != nil {
		t.Fatalf("grant owner: %v", err)
	}

	cfg := tool.Config{IP: "admin-9"}

	// 25a. valid value persists.
	if r := route.Api_setting_main_post(cfg, map[string]string{"trusted_proxies": "203.0.113.0/24"}); r["response"] != "ok" {
		t.Fatalf("valid trusted_proxies rejected: %v", r)
	}
	if got := firstSetting(db, "trusted_proxies"); got != "203.0.113.0/24" {
		t.Fatalf("valid value not persisted: %q", got)
	}

	// 25b. catch-all rejected; previous value preserved.
	if r := route.Api_setting_main_post(cfg, map[string]string{"trusted_proxies": "0.0.0.0/0"}); r["response"] != "error" {
		t.Fatalf("catch-all accepted: %v", r)
	}
	if got := firstSetting(db, "trusted_proxies"); got != "203.0.113.0/24" {
		t.Fatalf("previous value not preserved after catch-all: %q", got)
	}

	// 25c. malformed rejected; previous value preserved.
	if r := route.Api_setting_main_post(cfg, map[string]string{"trusted_proxies": "not-an-ip"}); r["response"] != "error" {
		t.Fatalf("malformed accepted: %v", r)
	}
	if got := firstSetting(db, "trusted_proxies"); got != "203.0.113.0/24" {
		t.Fatalf("previous value not preserved after malformed: %q", got)
	}
}

// ---------- Is_https trust boundary ----------

func Test_TRP_Is_https_trusted_proxy(t *testing.T) {
	db := new_ip_db(t)
	seedGroups(t, db)
	loadSetting(t, db, "trusted_proxies", actualIP)
	tool.Get_trusted_proxies(db)

	// trusted peer + https header -> true.
	c := ipContext(actualIP+":443", map[string]string{"X-Forwarded-Proto": "https"}, "")
	if !tool.Is_https(c) {
		t.Fatalf("trusted peer + https header: expected true")
	}

	// untrusted peer + https header -> false (no forge of Secure).
	loadSetting(t, db, "trusted_proxies", "10.0.0.1")
	tool.Get_trusted_proxies(db)
	c2 := ipContext(actualIP+":443", map[string]string{"X-Forwarded-Proto": "https"}, "")
	if tool.Is_https(c2) {
		t.Fatalf("untrusted peer + https header: expected false")
	}

	// trusted peer, no https header -> false.
	loadSetting(t, db, "trusted_proxies", actualIP)
	tool.Get_trusted_proxies(db)
	c3 := ipContext(actualIP+":8080", nil, "")
	if tool.Is_https(c3) {
		t.Fatalf("trusted peer without https header: expected false")
	}
}

// ---------- Validate_trusted_proxies ----------

func Test_TRP_validate_trusted_proxies(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want bool
	}{
		{"empty", "", true},
		{"single v4", "203.0.113.20", true},
		{"v4 prefix", "203.0.113.0/24", true},
		{"v6", "2001:db8::1", true},
		{"v6 prefix", trustV6, true},
		{"comma list", actualIP + ",203.0.113.21", true},
		{"catch-all v4", "0.0.0.0/0", false},
		{"catch-all v6", "::/0", false},
		{"malformed", "not-an-ip", false},
		{"mixed", "203.0.113.0/24,not-an-ip", false},
	}
	for _, tc := range cases {
		err := tool.Validate_trusted_proxies(tc.raw)
		if (err == nil) != tc.want {
			t.Fatalf("Validate_trusted_proxies(%s): valid=%v want valid=%v", tc.name, err == nil, tc.want)
		}
	}
}

// ---------- Trusted_prefix_count ----------

func Test_TRP_trusted_prefix_count(t *testing.T) {
	db := new_ip_db(t)
	loadSetting(t, db, "trusted_proxies", "")
	tool.Get_trusted_proxies(db)
	if n := tool.Trusted_prefix_count(); n != 0 {
		t.Fatalf("empty: count=%d want 0", n)
	}
	loadSetting(t, db, "trusted_proxies", actualIP+",203.0.113.21")
	tool.Get_trusted_proxies(db)
	if n := tool.Trusted_prefix_count(); n != 2 {
		t.Fatalf("two prefixes: count=%d want 2", n)
	}
	loadSetting(t, db, "trusted_proxies", "0.0.0.0/0")
	tool.Get_trusted_proxies(db)
	if n := tool.Trusted_prefix_count(); n != 0 {
		t.Fatalf("catch-all: count=%d want 0", n)
	}
}

// ---------- netip parsing contract ----------

func Test_TRP_netip_parse_contract(t *testing.T) {
	// v4-mapped v6: ParseAddr succeeds and the address is unmappable to a plain
	// v4 address for prefix matching.
	a, err := netip.ParseAddr("::ffff:203.0.113.20")
	if err != nil || !a.Is4In6() {
		t.Fatalf("ParseAddr(::ffff:203.0.113.20): err=%v is4in6=%v", err, a.Is4In6())
	}
	px, perr := netip.ParsePrefix("203.0.113.0/24")
	if perr != nil {
		t.Fatalf("ParsePrefix: %v", perr)
	}
	if !px.Contains(a.Unmap()) {
		t.Fatalf("v4 prefix must contain 203.0.113.20 after Unmap")
	}
	// catch-all prefixes are detectable (Bits == 0) so Validate can reject them.
	if ca, cerr := netip.ParsePrefix("0.0.0.0/0"); cerr != nil || ca.Bits() != 0 {
		t.Fatalf("catch-all v4: err=%v bits=%d", cerr, ca.Bits())
	}
	if ca6, cerr := netip.ParsePrefix("::/0"); cerr != nil || ca6.Bits() != 0 {
		t.Fatalf("catch-all v6: err=%v bits=%d", cerr, ca6.Bits())
	}
	// ParseAddr rejects host:port forms (that form is handled by SplitHostPort).
	if _, err := netip.ParseAddr("203.0.113.20:80"); err == nil {
		t.Fatalf("ParseAddr should reject 203.0.113.20:80")
	}
}

// ---------- fail-closed invalid-remote-address regression (post-fix) ----------
//
// These tests assert the CORRECTED behavior: the Invalid_remote_address
// sentinel is denied by the shared authorization boundary and never gains
// user-level privileges, while valid anonymous, authenticated, and
// trusted-proxy identities are unchanged from before. They replace the
// removed Test_TRP_invalid_remote_address_pre_fix_repro reproduction test.

// test_invalid_peer_db builds an isolated DB with production-style groups,
// a public document, a few users, and an explicitly banned IP.
func test_invalid_peer_db(t *testing.T) *sql.DB {
	t.Helper()
	db := new_ip_db(t)
	seedGroups(t, db)
	// A public document (no per-document view ACL).
	if _, err := db.Exec(`insert into data (title,data,type) values (?,?,?)`,
		"sentinel_doc", "", "normal"); err != nil {
		t.Fatalf("insert document: %v", err)
	}
	// users for the "behaves as before" cases
	grant(t, db, "alice", "user")
	grant(t, db, "bob", "admin")
	grant(t, db, "carol", "owner")
	// an explicitly banned IP
	grant(t, db, "198.51.100.50", "ban")
	// default load mode, no trusted proxy: the TCP peer is the identity
	loadSetting(t, db, "load_ip_select", "REMOTE_ADDR")
	loadSetting(t, db, "trusted_proxies", "")
	tool.Get_ip_select(db)
	tool.Get_trusted_proxies(db)
	return db
}

// Test_TRP_26_malformed_remote_addr_denied_by_site_view_middleware
//
// Builds a real gin request with an unparseable RemoteAddr and runs the EXACT
// decision of Site_view_middleware (routes_access.go):
//
//	user_id := tool.Get_IP(c)
//	if !tool.Get_auth_info(db, user_id)["site_view"] { abort 403 }
//
// and asserts the sentinel takes the 403 path.
func Test_TRP_26_malformed_remote_addr_denied_by_site_view_middleware(t *testing.T) {
	db := test_invalid_peer_db(t)
	c := ipContext("garbage:8080", nil, "")

	user_id := tool.Get_IP(c)
	if user_id != tool.Invalid_remote_address {
		t.Fatalf("invalid peer resolved to %q; want sentinel", user_id)
	}
	if tool.Get_auth_info(db, user_id)["site_view"] {
		t.Fatalf("the site_view_middleware must deny the invalid peer (got site_view=true)")
	}
}

// Test_TRP_27_invalid_identity_has_no_effective_user_permissions
//
// The sentinel must not grant any user-level privilege at the permission
// boundary, whether checked through Get_auth_info, Check_permission, or
// Check_permission.
func Test_TRP_27_invalid_identity_has_no_effective_user_permissions(t *testing.T) {
	db := test_invalid_peer_db(t)
	info := tool.Get_auth_info(db, tool.Invalid_remote_address)
	for _, perm := range []string{"site_view", "view", "edit", "bbs_use", "bbs_edit",
		"bbs_comment", "document_edit", "document_delete", "image_edit", "admin",
		"owner"} {
		if info[perm] {
			t.Fatalf("invalid identity must not grant %q, got %v", perm, info[perm])
		}
	}
	// The permission boundary (Check_permission) must also deny every listed
	// privilege for the sentinel.
	for _, perm := range []string{"site_view", "view", "edit"} {
		if tool.Check_permission(db, perm, tool.Invalid_remote_address) {
			t.Fatalf("Check_permission(%q, sentinel) must be false, got true", perm)
		}
	}
	// The site-view gate (used by Site_view_middleware) must be false too.
	if info["site_view"] {
		t.Fatalf("site_view must be false for the sentinel")
	}
}

// Test_TRP_28_invalid_identity_denied_from_document_and_protected_api
//
// The sentinel must be denied at document access (render and edit ACL) and at
// any permission-gated API. An explicitly banned IP is denied the same way.
func Test_TRP_28_invalid_identity_denied_from_document_and_protected_api(t *testing.T) {
	db := test_invalid_peer_db(t)
	sentinel := tool.Invalid_remote_address
	if tool.Check_acl(db, "sentinel_doc", "", "render", sentinel) {
		t.Fatalf("invalid identity must not view a document (Check_acl render=true)")
	}
	if tool.Check_acl(db, "sentinel_doc", "", "document_edit", sentinel) {
		t.Fatalf("invalid identity must not edit a document (Check_acl document_edit=true)")
	}
	// the explicitly banned IP is also denied site access
	if tool.Check_permission(db, "site_view", "198.51.100.50") {
		t.Fatalf("banned IP 198.51.100.50 must not have site_view")
	}
}

// Test_TRP_29_valid_anonymous_ipv4_ipv6_unchanged
//
// Valid anonymous IPv4 and IPv6 peers still resolve to themselves and keep
// the existing ip-level permissions.
func Test_TRP_29_valid_anonymous_ipv4_ipv6_unchanged(t *testing.T) {
	db := test_invalid_peer_db(t)

	c4 := ipContext("203.0.113.20:8080", nil, "")
	id4 := tool.Get_IP(c4)
	if id4 != "203.0.113.20" {
		t.Fatalf("valid IPv4 Get_IP: got %q, want 203.0.113.20", id4)
	}
	info4 := tool.Get_auth_info(db, id4)
	if !info4["site_view"] || !info4["view"] || !info4["bbs_use"] {
		t.Fatalf("valid IPv4 must keep ip-level perms (site_view=%v view=%v bbs_use=%v)",
			info4["site_view"], info4["view"], info4["bbs_use"])
	}

	c6 := ipContext("2001:db8:85a3::8a2e:370:7334", nil, "")
	id6 := tool.Get_IP(c6)
	if id6 != "2001:db8:85a3::8a2e:370:7334" {
		t.Fatalf("valid IPv6 Get_IP: got %q, want 2001:db8:85a3::8a2e:370:7334", id6)
	}
	info6 := tool.Get_auth_info(db, id6)
	if !info6["site_view"] || !info6["view"] || !info6["bbs_use"] {
		t.Fatalf("valid IPv6 must keep ip-level perms (site_view=%v view=%v bbs_use=%v)",
			info6["site_view"], info6["view"], info6["bbs_use"])
	}
}

// Test_TRP_30_valid_authenticated_session_unchanged
//
// A logged-in session still resolves to the username (not the peer) and keeps
// the user-group privileges.
func Test_TRP_30_valid_authenticated_session_unchanged(t *testing.T) {
	db := test_invalid_peer_db(t)
	c := ipContext("203.0.113.20:8080", nil, "alice")
	id := tool.Get_IP(c)
	if id != "alice" {
		t.Fatalf("authenticated Get_IP: got %q, want alice", id)
	}
	info := tool.Get_auth_info(db, id)
	if !info["site_view"] || !info["edit"] {
		t.Fatalf("authenticated user must keep site_view+edit (site_view=%v edit=%v)",
			info["site_view"], info["edit"])
	}
}

// Test_TRP_31_trusted_proxy_header_handling_unaffected
//
// A configured trusted peer still honors the forwarded header, and a malformed
// peer still fails closed to the sentinel. The fix did not change trusted-proxy
// resolution.
func Test_TRP_31_trusted_proxy_header_handling_unaffected(t *testing.T) {
	db := test_invalid_peer_db(t)
	loadSetting(t, db, "load_ip_select", "X-Real-IP")
	loadSetting(t, db, "trusted_proxies", "203.0.113.0/24")
	tool.Get_ip_select(db)
	tool.Get_trusted_proxies(db)

	// trusted peer: the header is honored (existing behavior)
	c := ipContext("203.0.113.20:8080", map[string]string{"X-Real-IP": "198.51.100.10"}, "")
	id := tool.Get_IP(c)
	if id != "198.51.100.10" {
		t.Fatalf("trusted peer must honor the header, got %q", id)
	}
	// malformed peer still fails closed, independent of the trust config
	c2 := ipContext("garbage:8080", nil, "")
	id2 := tool.Get_IP(c2)
	if id2 != tool.Invalid_remote_address {
		t.Fatalf("malformed peer must still fail closed, got %q", id2)
	}
}

// Test_TRP_32_untrusted_forwarding_headers_remain_ignored
//
// In a forwarding mode with no matching trusted proxy, forwarded headers are
// ignored and the TCP peer is used. This covers both X-Real-IP and
// CF-Connecting-IP.
func Test_TRP_32_untrusted_forwarding_headers_remain_ignored(t *testing.T) {
	db := test_invalid_peer_db(t)
	loadSetting(t, db, "load_ip_select", "X-Real-IP")
	loadSetting(t, db, "trusted_proxies", "")
	tool.Get_ip_select(db)
	tool.Get_trusted_proxies(db)

	// no trusted proxy: X-Real-IP is ignored, peer is used
	c := ipContext("203.0.113.20:8080", map[string]string{"X-Real-IP": "198.51.100.10"}, "")
	id := tool.Get_IP(c)
	if id != "203.0.113.20" {
		t.Fatalf("untrusted X-Real-IP must be ignored, got %q", id)
	}
	// CF-Connecting-IP ignored on an untrusted v6 peer
	c2 := ipContext("2001:db8:85a3::8a2e:370:7334",
		map[string]string{"CF-Connecting-IP": "198.51.100.10"}, "")
	id2 := tool.Get_IP(c2)
	if id2 != "2001:db8:85a3::8a2e:370:7334" {
		t.Fatalf("untrusted CF-Connecting-IP must be ignored, got %q", id2)
	}
}
