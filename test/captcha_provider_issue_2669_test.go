package test

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"

	altcha "github.com/altcha-org/altcha-lib-go/v2"
	"github.com/gin-gonic/gin"

	"opennamu/route"
	"opennamu/route/tool"
)

var captcha_issue_2669_setting_names = []string{"recaptcha", "sec_re", "altcha_sec_re", "recaptcha_ver"}

func Captcha_issue_2669_test_backends() []string {
	backends := []string{"sqlite", "turso"}
	if os.Getenv("MYSQL_AUDIT_DSN") != "" {
		backends = append(backends, "mysql")
	}
	return backends
}

func Captcha_issue_2669_test_setup(t *testing.T, backend string) *sql.DB {
	t.Helper()
	if backend == "mysql" {
		db, _ := mysql_audit_setup(t)
		return db
	}
	return audit_setup(t, backend)
}

func Captcha_issue_2669_test_set_settings(t *testing.T, db *sql.DB, values map[string]string) {
	t.Helper()
	for _, name := range captcha_issue_2669_setting_names {
		if _, err := db.Exec(`delete from other where name = ? and coverage = ''`, name); err != nil {
			t.Fatalf("clear CAPTCHA setting %s: %v", name, err)
		}
		value, exists := values[name]
		if !exists {
			continue
		}
		if _, err := db.Exec(`insert into other (name, data, coverage) values (?, ?, '')`, name, value); err != nil {
			t.Fatalf("insert CAPTCHA setting %s: %v", name, err)
		}
	}
}

func Captcha_issue_2669_test_setting(t *testing.T, db *sql.DB, name string) string {
	t.Helper()
	value, err := tool.Captcha_setting_value_unique(db, name)
	if err != nil {
		t.Fatalf("read CAPTCHA setting %s: %v", name, err)
	}
	return value
}

func Test_Captcha_issue_2669_provider_resolution_and_frontend_assets(t *testing.T) {
	provider_cases := []struct {
		name             string
		provider         string
		frontend_marker  string
		forbidden_marker string
	}{
		{name: "google_v2", provider: "v2", frontend_marker: "www.google.com/recaptcha/api.js", forbidden_marker: "js.hcaptcha.com"},
		{name: "google_v3", provider: "v3", frontend_marker: "grecaptcha.execute", forbidden_marker: "js.hcaptcha.com"},
		{name: "hcaptcha", provider: "h", frontend_marker: "js.hcaptcha.com/1/api.js", forbidden_marker: "www.google.com/recaptcha/api.js"},
		{name: "turnstile", provider: "cf", frontend_marker: "challenges.cloudflare.com/turnstile/v0/api.js", forbidden_marker: "js.hcaptcha.com"},
		{name: "altcha", provider: "altcha_low", frontend_marker: "<altcha-widget", forbidden_marker: "js.hcaptcha.com"},
	}

	for _, backend := range Captcha_issue_2669_test_backends() {
		t.Run(backend, func(t *testing.T) {
			db := Captcha_issue_2669_test_setup(t, backend)
			for _, provider_case := range provider_cases {
				t.Run(provider_case.name, func(t *testing.T) {
					Captcha_issue_2669_test_set_settings(t, db, map[string]string{
						"recaptcha":     "captcha-test-public-key",
						"sec_re":        "captcha-test-google-secret",
						"altcha_sec_re": "captcha-test-altcha-secret",
						"recaptcha_ver": provider_case.provider,
					})

					resolved_config, err := tool.Captcha_config_resolve(db)
					if err != nil {
						t.Fatalf("resolve %s configuration: %v", provider_case.provider, err)
					}
					if resolved_config.Provider != provider_case.provider {
						t.Fatalf("resolved provider = %q, want %q", resolved_config.Provider, provider_case.provider)
					}

					frontend := tool.Get_captcha_ui(db, tool.Config{IP: "192.0.2.99"})
					if !strings.Contains(frontend, provider_case.frontend_marker) {
						t.Fatalf("%s frontend omitted %q: %s", provider_case.provider, provider_case.frontend_marker, frontend)
					}
					if strings.Contains(frontend, provider_case.forbidden_marker) {
						t.Fatalf("%s frontend unexpectedly contains %q: %s", provider_case.provider, provider_case.forbidden_marker, frontend)
					}
					if strings.Contains(frontend, "captcha-test-google-secret") || strings.Contains(frontend, "captcha-test-altcha-secret") {
						t.Fatalf("%s frontend exposed a CAPTCHA secret", provider_case.provider)
					}
				})
			}
		})
	}
}

func Test_Captcha_issue_2669_initialization_migrations_are_safe_and_idempotent(t *testing.T) {
	for _, backend := range Captcha_issue_2669_test_backends() {
		t.Run(backend, func(t *testing.T) {
			db := Captcha_issue_2669_test_setup(t, backend)

			t.Run("legacy_google_keys_become_v2", func(t *testing.T) {
				Captcha_issue_2669_test_set_settings(t, db, map[string]string{
					"recaptcha": "legacy-test-public",
					"sec_re":    "legacy-test-secret",
				})
				tool.Init_captcha(db)
				if got := Captcha_issue_2669_test_setting(t, db, "recaptcha_ver"); got != "v2" {
					t.Fatalf("migrated provider = %q, want v2", got)
				}
				if got := Captcha_issue_2669_test_setting(t, db, "recaptcha"); got != "legacy-test-public" {
					t.Fatalf("legacy public key changed to %q", got)
				}
				if got := Captcha_issue_2669_test_setting(t, db, "sec_re"); got != "legacy-test-secret" {
					t.Fatalf("legacy secret changed during provider migration")
				}
				tool.Init_captcha(db)
				if got := Captcha_issue_2669_test_setting(t, db, "recaptcha_ver"); got != "v2" {
					t.Fatalf("second initialization changed provider to %q", got)
				}
			})

			t.Run("empty_legacy_provider_becomes_v2", func(t *testing.T) {
				Captcha_issue_2669_test_set_settings(t, db, map[string]string{
					"recaptcha":     "legacy-test-public",
					"sec_re":        "legacy-test-secret",
					"recaptcha_ver": "",
				})
				tool.Init_captcha(db)
				if got := Captcha_issue_2669_test_setting(t, db, "recaptcha_ver"); got != "v2" {
					t.Fatalf("empty legacy provider migrated to %q, want v2", got)
				}
			})

			t.Run("legacy_altcha_secret_uses_verified_migration_shape", func(t *testing.T) {
				Captcha_issue_2669_test_set_settings(t, db, map[string]string{
					"sec_re":        "legacy-altcha-secret",
					"recaptcha_ver": "altcha_high",
				})
				tool.Init_captcha(db)
				if got := Captcha_issue_2669_test_setting(t, db, "sec_re"); got != "" {
					t.Fatal("legacy ALTCHA migration did not clear the old shared secret field")
				}
				if got := Captcha_issue_2669_test_setting(t, db, "altcha_sec_re"); got != "legacy-altcha-secret" {
					t.Fatal("verified legacy ALTCHA secret was not preserved in its dedicated field")
				}
				resolved_config, err := tool.Captcha_config_resolve(db)
				if err != nil || resolved_config.Provider != "altcha_high" {
					t.Fatalf("legacy ALTCHA configuration = %q, %v", resolved_config.Provider, err)
				}
			})

			t.Run("explicit_altcha_does_not_move_google_credentials", func(t *testing.T) {
				Captcha_issue_2669_test_set_settings(t, db, map[string]string{
					"recaptcha":     "legacy-google-public",
					"sec_re":        "legacy-google-secret",
					"recaptcha_ver": "altcha_high",
				})
				tool.Init_captcha(db)
				if got := Captcha_issue_2669_test_setting(t, db, "sec_re"); got != "legacy-google-secret" {
					t.Fatal("ALTCHA initialization changed the reCAPTCHA secret")
				}
				if got := Captcha_issue_2669_test_setting(t, db, "altcha_sec_re"); strings.TrimSpace(got) == "" || got == "legacy-google-secret" {
					t.Fatal("ALTCHA provider did not receive an independent secret key")
				}
				resolved_config, err := tool.Captcha_config_resolve(db)
				if err != nil || resolved_config.Provider != "altcha_high" {
					t.Fatalf("ALTCHA configuration = %q, %v", resolved_config.Provider, err)
				}
			})

			t.Run("new_install_defaults_to_altcha_and_is_idempotent", func(t *testing.T) {
				Captcha_issue_2669_test_set_settings(t, db, nil)
				if _, err := tool.Captcha_config_resolve(db); err == nil {
					t.Fatal("missing provider and credentials unexpectedly resolved as an active configuration")
				}
				tool.Init_captcha(db)
				first_secret := Captcha_issue_2669_test_setting(t, db, "altcha_sec_re")
				if first_secret == "" || Captcha_issue_2669_test_setting(t, db, "recaptcha_ver") != "altcha_high" {
					t.Fatal("new installation did not initialize ALTCHA high with a secret")
				}
				tool.Init_captcha(db)
				if second_secret := Captcha_issue_2669_test_setting(t, db, "altcha_sec_re"); second_secret != first_secret {
					t.Fatal("repeated initialization replaced the generated ALTCHA secret")
				}
			})

			t.Run("unsupported_provider_is_not_rewritten", func(t *testing.T) {
				Captcha_issue_2669_test_set_settings(t, db, map[string]string{
					"recaptcha":     "test-public",
					"sec_re":        "test-secret",
					"altcha_sec_re": "existing-altcha-secret",
					"recaptcha_ver": "unrecognized-provider",
				})
				tool.Init_captcha(db)
				if got := Captcha_issue_2669_test_setting(t, db, "recaptcha_ver"); got != "unrecognized-provider" {
					t.Fatalf("unsupported provider was rewritten to %q", got)
				}
				if got := Captcha_issue_2669_test_setting(t, db, "sec_re"); got != "test-secret" {
					t.Fatal("unsupported provider handling changed the reCAPTCHA secret")
				}
				if _, err := tool.Captcha_config_resolve(db); err == nil {
					t.Fatal("unsupported provider unexpectedly resolved")
				}
			})

			t.Run("non_global_provider_rows_do_not_override_default", func(t *testing.T) {
				Captcha_issue_2669_test_set_settings(t, db, nil)
				for _, setting := range [][2]string{{"recaptcha_ver", "h"}, {"recaptcha", "custom-public"}, {"sec_re", "custom-secret"}} {
					if _, err := db.Exec(`insert into other (name, data, coverage) values (?, ?, 'custom')`, setting[0], setting[1]); err != nil {
						t.Fatalf("insert covered setting: %v", err)
					}
				}
				tool.Init_captcha(db)
				if got := Captcha_issue_2669_test_setting(t, db, "recaptcha_ver"); got != "altcha_high" {
					t.Fatalf("non-global settings selected provider %q, want global ALTCHA default", got)
				}
				if got := Captcha_issue_2669_test_setting(t, db, "recaptcha"); got != "" {
					t.Fatalf("non-global public key leaked into global setting: %q", got)
				}
			})

			t.Run("conflicting_duplicate_provider_rows_fail_closed", func(t *testing.T) {
				Captcha_issue_2669_test_set_settings(t, db, map[string]string{
					"recaptcha":     "test-public",
					"sec_re":        "test-secret",
					"recaptcha_ver": "v2",
				})
				if _, err := db.Exec(`insert into other (name, data, coverage) values ('recaptcha_ver', 'h', '')`); err != nil {
					t.Fatalf("insert conflicting provider: %v", err)
				}
				tool.Init_captcha(db)
				if _, err := tool.Captcha_config_resolve(db); err == nil {
					t.Fatal("conflicting provider rows unexpectedly resolved")
				}
				if got := Captcha_issue_2669_test_setting(t, db, "sec_re"); got != "test-secret" {
					t.Fatal("initialization changed a secret while provider rows conflicted")
				}
			})

			t.Run("conflicting_duplicate_credentials_fail_closed", func(t *testing.T) {
				Captcha_issue_2669_test_set_settings(t, db, map[string]string{
					"recaptcha":     "test-public",
					"sec_re":        "test-secret-one",
					"recaptcha_ver": "v2",
				})
				if _, err := db.Exec(`insert into other (name, data, coverage) values ('sec_re', 'test-secret-two', '')`); err != nil {
					t.Fatalf("insert conflicting secret: %v", err)
				}
				if _, err := tool.Captcha_config_resolve(db); err == nil {
					t.Fatal("conflicting provider secret rows unexpectedly resolved")
				}
				if tool.Captcha_check(db, &emptySession{}, "192.0.2.99", "test-token") {
					t.Fatal("conflicting provider secrets were accepted")
				}
			})

			t.Run("same_duplicate_provider_rows_are_stable_and_whitespace_is_normalized", func(t *testing.T) {
				Captcha_issue_2669_test_set_settings(t, db, map[string]string{
					"recaptcha":     "test-public",
					"sec_re":        "test-secret",
					"recaptcha_ver": " v2 ",
				})
				if _, err := db.Exec(`insert into other (name, data, coverage) values ('recaptcha_ver', 'v2', '')`); err != nil {
					t.Fatalf("insert equivalent provider: %v", err)
				}
				resolved_config, err := tool.Captcha_config_resolve(db)
				if err != nil || resolved_config.Provider != "v2" {
					t.Fatalf("normalized provider = %q, %v", resolved_config.Provider, err)
				}
			})
		})
	}
}

func Test_Captcha_issue_2669_external_settings_validate_before_saving(t *testing.T) {
	for _, backend := range Captcha_issue_2669_test_backends() {
		t.Run(backend, func(t *testing.T) {
			db := Captcha_issue_2669_test_setup(t, backend)
			audit_seed_auth(db, t)
			Captcha_issue_2669_test_set_settings(t, db, map[string]string{
				"recaptcha":     "stored-test-public",
				"sec_re":        "stored-test-google-secret",
				"altcha_sec_re": "stored-test-altcha-secret",
				"recaptcha_ver": "v2",
			})
			if _, err := db.Exec(`insert into other (name, data, coverage) values ('smtp_server', 'original-test-smtp', '')`); err != nil {
				t.Fatalf("seed SMTP setting: %v", err)
			}

			invalid_form := map[string]string{
				"recaptcha_ver": "unrecognized-provider",
				"smtp_server":   "must-not-be-saved",
			}
			invalid_result := route.Api_setting_external_post(tool.Config{IP: "admin_ip"}, invalid_form)
			if invalid_result["response"] != "error" || invalid_result["data"] != "captcha configuration" {
				t.Fatalf("unsupported provider response = %#v", invalid_result)
			}
			if got := tool.Get_setting_value_exact(db, "smtp_server", "", ""); got != "original-test-smtp" {
				t.Fatalf("invalid provider partially saved SMTP setting %q", got)
			}
			if got := Captcha_issue_2669_test_setting(t, db, "recaptcha_ver"); got != "v2" {
				t.Fatalf("invalid provider changed the saved selection to %q", got)
			}

			missing_provider_form := map[string]string{"smtp_server": "also-must-not-be-saved"}
			missing_provider_result := route.Api_setting_external_post(tool.Config{IP: "admin_ip"}, missing_provider_form)
			if missing_provider_result["response"] != "error" {
				t.Fatalf("missing provider field response = %#v", missing_provider_result)
			}
			if got := tool.Get_setting_value_exact(db, "smtp_server", "", ""); got != "original-test-smtp" {
				t.Fatalf("missing provider partially saved SMTP setting %q", got)
			}

			valid_form := map[string]string{
				"recaptcha":     "stored-test-public",
				"sec_re":        "",
				"altcha_sec_re": "",
				"recaptcha_ver": "v2",
				"smtp_server":   "saved-test-smtp",
				"ai_provider":   "ollama",
			}
			valid_result := route.Api_setting_external_post(tool.Config{IP: "admin_ip"}, valid_form)
			if valid_result["response"] != "ok" {
				t.Fatalf("valid unchanged settings response = %#v", valid_result)
			}
			if got := Captcha_issue_2669_test_setting(t, db, "sec_re"); got != "stored-test-google-secret" {
				t.Fatal("blank password input erased the stored reCAPTCHA secret")
			}
			if got := Captcha_issue_2669_test_setting(t, db, "altcha_sec_re"); got != "stored-test-altcha-secret" {
				t.Fatal("blank password input erased the stored ALTCHA secret")
			}
			if got := tool.Get_setting_value_exact(db, "smtp_server", "", ""); got != "saved-test-smtp" {
				t.Fatalf("valid settings did not save SMTP value: %q", got)
			}

			legacy_provider_form := map[string]string{
				"recaptcha":     "stored-test-public",
				"sec_re":        "",
				"altcha_sec_re": "",
				"recaptcha_ver": "",
				"smtp_server":   "saved-test-smtp",
				"ai_provider":   "ollama",
			}
			legacy_provider_result := route.Api_setting_external_post(tool.Config{IP: "admin_ip"}, legacy_provider_form)
			if legacy_provider_result["response"] != "ok" || Captcha_issue_2669_test_setting(t, db, "recaptcha_ver") != "v2" {
				t.Fatalf("legacy empty provider form was not saved as v2: %#v", legacy_provider_result)
			}

			altcha_form := map[string]string{
				"recaptcha":     "stored-test-public",
				"sec_re":        "",
				"altcha_sec_re": "",
				"recaptcha_ver": "altcha_low",
				"smtp_server":   "saved-test-smtp",
				"ai_provider":   "ollama",
			}
			altcha_result := route.Api_setting_external_post(tool.Config{IP: "admin_ip"}, altcha_form)
			if altcha_result["response"] != "ok" || Captcha_issue_2669_test_setting(t, db, "recaptcha_ver") != "altcha_low" {
				t.Fatalf("valid ALTCHA selection was not saved: %#v", altcha_result)
			}
			if got := Captcha_issue_2669_test_setting(t, db, "sec_re"); got != "stored-test-google-secret" {
				t.Fatal("switching to ALTCHA modified the existing reCAPTCHA secret")
			}
			if got := Captcha_issue_2669_test_setting(t, db, "altcha_sec_re"); got != "stored-test-altcha-secret" {
				t.Fatal("switching to ALTCHA did not preserve its independent secret")
			}

			Captcha_issue_2669_test_set_settings(t, db, map[string]string{"recaptcha_ver": "v2"})
			incomplete_form := map[string]string{
				"recaptcha":     "",
				"sec_re":        "",
				"altcha_sec_re": "",
				"recaptcha_ver": "v2",
				"smtp_server":   "must-still-not-be-saved",
			}
			incomplete_result := route.Api_setting_external_post(tool.Config{IP: "admin_ip"}, incomplete_form)
			if incomplete_result["response"] != "error" {
				t.Fatalf("incomplete provider response = %#v", incomplete_result)
			}
			if got := tool.Get_setting_value_exact(db, "smtp_server", "", ""); got != "saved-test-smtp" {
				t.Fatalf("incomplete provider partially saved SMTP setting %q", got)
			}
		})
	}
}

func Test_Captcha_issue_2669_external_settings_redact_secrets_and_flag_invalid_provider(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		t.Run(backend, func(t *testing.T) {
			db := Captcha_issue_2669_test_setup(t, backend)
			audit_seed_auth(db, t)
			Captcha_issue_2669_test_set_settings(t, db, map[string]string{
				"recaptcha":     "visible-test-public-key",
				"sec_re":        "hidden-test-google-secret",
				"altcha_sec_re": "hidden-test-altcha-secret",
				"recaptcha_ver": "unsupported-test-provider",
			})

			page := route.View_setting_external(tool.Config{IP: "admin_ip", Session: &emptySession{}})
			if strings.Contains(page, "hidden-test-google-secret") || strings.Contains(page, "hidden-test-altcha-secret") {
				t.Fatal("external settings page exposed a CAPTCHA secret")
			}
			if !strings.Contains(page, tool.Get_language(db, "captcha_config_error", true)) {
				t.Fatal("external settings page omitted invalid-configuration recovery guidance")
			}
			if !strings.Contains(page, tool.Get_language(db, "captcha_provider_invalid", true)) {
				t.Fatal("external settings page did not identify the unsupported saved provider")
			}
			if strings.Contains(page, `<option value="v2" selected="selected">`) {
				t.Fatal("unsupported provider was silently presented as a saved v2 selection")
			}
		})
	}
}

type captcha_issue_2669_verification_request struct {
	path   string
	fields url.Values
}

type captcha_issue_2669_transport struct {
	base   http.RoundTripper
	target *url.URL
	mu     sync.Mutex
	hosts  []string
}

func (transport *captcha_issue_2669_transport) RoundTrip(request *http.Request) (*http.Response, error) {
	transport.mu.Lock()
	transport.hosts = append(transport.hosts, request.URL.Host)
	transport.mu.Unlock()

	request_copy := request.Clone(request.Context())
	url_copy := *transport.target
	url_copy.Path = request.URL.Path
	url_copy.RawPath = request.URL.RawPath
	url_copy.RawQuery = request.URL.RawQuery
	request_copy.URL = &url_copy
	request_copy.Host = ""
	return transport.base.RoundTrip(request_copy)
}

func (transport *captcha_issue_2669_transport) reset() {
	transport.mu.Lock()
	transport.hosts = nil
	transport.mu.Unlock()
}

func (transport *captcha_issue_2669_transport) snapshot() []string {
	transport.mu.Lock()
	defer transport.mu.Unlock()
	return append([]string(nil), transport.hosts...)
}

func Captcha_issue_2669_test_post_login(t *testing.T, client *http.Client, target string, fields url.Values) string {
	t.Helper()
	request, err := http.NewRequest(http.MethodPost, target+"/login", strings.NewReader(fields.Encode()))
	if err != nil {
		t.Fatalf("make login request: %v", err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("post login: %v", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("read login response: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("login status = %d, want 200", response.StatusCode)
	}
	return string(body)
}

func Test_Captcha_issue_2669_login_http_flow_uses_the_selected_provider(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		t.Run(backend, func(t *testing.T) {
			db := Captcha_issue_2669_test_setup(t, backend)
			audit_seed_auth(db, t)
			for _, group := range []string{"ip", "user"} {
				if _, err := db.Exec(`insert into alist (name, acl) values (?, 'login_available')`, group); err != nil {
					t.Fatalf("allow test login for %s: %v", group, err)
				}
			}
			create_result := route.Api_add_user(tool.Config{}, "captcha_test_user", "captcha-test-password", "", "")
			if create_result["response"] != "ok" {
				t.Fatalf("create test login account: %#v", create_result)
			}

			verify_requests := make([]captcha_issue_2669_verification_request, 0)
			var verify_requests_mu sync.Mutex
			verify_server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				if err := request.ParseForm(); err != nil {
					writer.WriteHeader(http.StatusBadRequest)
					return
				}
				verify_requests_mu.Lock()
				verify_requests = append(verify_requests, captcha_issue_2669_verification_request{path: request.URL.Path, fields: request.Form})
				verify_requests_mu.Unlock()
				writer.Header().Set("Content-Type", "application/json")
				_, _ = io.WriteString(writer, `{"success":true}`)
			}))
			t.Cleanup(verify_server.Close)
			target_url, err := url.Parse(verify_server.URL)
			if err != nil {
				t.Fatalf("parse local verifier URL: %v", err)
			}
			previous_transport := http.DefaultTransport
			verification_transport := &captcha_issue_2669_transport{base: previous_transport, target: target_url}
			http.DefaultTransport = verification_transport
			t.Cleanup(func() { http.DefaultTransport = previous_transport })

			previous_gin_mode := gin.Mode()
			gin.SetMode(gin.TestMode)
			t.Cleanup(func() { gin.SetMode(previous_gin_mode) })
			config := tool.Config{IP: "192.0.2.99", Session: &emptySession{}}
			router := gin.New()
			router.GET("/login", func(context *gin.Context) {
				context.String(http.StatusOK, route.View_login_login(config))
			})
			router.POST("/login", func(context *gin.Context) {
				_ = context.Request.ParseForm()
				captcha_response := tool.Captcha_response_internal(
					context.PostForm("g-recaptcha"),
					context.PostForm("g-recaptcha-response"),
					context.PostForm("h-captcha-response"),
					context.PostForm("cf-turnstile-response"),
					context.PostForm("altcha"),
				)
				context.String(http.StatusOK, route.View_login_login_post_full(config, context.PostForm("id"), context.PostForm("password"), captcha_response))
			})
			router_server := httptest.NewServer(router)
			t.Cleanup(router_server.Close)
			router_client := router_server.Client()

			provider_cases := []struct {
				provider string
				field    string
				host     string
				path     string
				asset    string
			}{
				{provider: "v2", field: "g-recaptcha-response", host: "www.google.com", path: "/recaptcha/api/siteverify", asset: "www.google.com/recaptcha/api.js"},
				{provider: "v3", field: "g-recaptcha", host: "www.google.com", path: "/recaptcha/api/siteverify", asset: "grecaptcha.execute"},
				{provider: "h", field: "h-captcha-response", host: "hcaptcha.com", path: "/siteverify", asset: "js.hcaptcha.com/1/api.js"},
				{provider: "cf", field: "cf-turnstile-response", host: "challenges.cloudflare.com", path: "/turnstile/v0/siteverify", asset: "challenges.cloudflare.com/turnstile/v0/api.js"},
			}
			for _, provider_case := range provider_cases {
				t.Run(provider_case.provider, func(t *testing.T) {
					Captcha_issue_2669_test_set_settings(t, db, map[string]string{
						"recaptcha":     "captcha-test-public-key",
						"sec_re":        "captcha-test-secret",
						"recaptcha_ver": provider_case.provider,
					})
					get_response, err := router_client.Get(router_server.URL + "/login")
					if err != nil {
						t.Fatalf("get login form: %v", err)
					}
					form_html, read_err := io.ReadAll(get_response.Body)
					_ = get_response.Body.Close()
					if read_err != nil {
						t.Fatalf("read login form: %v", read_err)
					}
					if get_response.StatusCode != http.StatusOK || !strings.Contains(string(form_html), provider_case.asset) {
						t.Fatalf("%s login form omitted provider asset %q", provider_case.provider, provider_case.asset)
					}

					verification_transport.reset()
					form := url.Values{
						"id":                {"captcha_test_user"},
						"password":          {"captcha-test-password"},
						provider_case.field: {"captcha-test-response"},
					}
					login_result := Captcha_issue_2669_test_post_login(t, router_client, router_server.URL, form)
					if !strings.Contains(login_result, `<!--opennamu_redirect:"/user"-->`) {
						t.Fatalf("%s CAPTCHA-verified login did not reach the normal login redirect", provider_case.provider)
					}
					if hosts := verification_transport.snapshot(); len(hosts) != 1 || hosts[0] != provider_case.host {
						t.Fatalf("%s verifier hosts = %v, want only %q", provider_case.provider, hosts, provider_case.host)
					}

					verify_requests_mu.Lock()
					last_request := verify_requests[len(verify_requests)-1]
					verify_requests_mu.Unlock()
					if last_request.path != provider_case.path {
						t.Fatalf("%s verifier path = %q, want %q", provider_case.provider, last_request.path, provider_case.path)
					}
					if last_request.fields.Get("secret") != "captcha-test-secret" || last_request.fields.Get("response") != "captcha-test-response" {
						t.Fatalf("%s verifier request did not carry the expected test-only token fields", provider_case.provider)
					}
				})
			}

			t.Run("unknown_provider_fails_closed_without_network_fallback", func(t *testing.T) {
				Captcha_issue_2669_test_set_settings(t, db, map[string]string{
					"recaptcha":     "captcha-test-public-key",
					"sec_re":        "captcha-test-secret",
					"recaptcha_ver": "unknown-test-provider",
				})
				verification_transport.reset()
				form := url.Values{
					"id":                   {"captcha_test_user"},
					"password":             {"captcha-test-password"},
					"g-recaptcha-response": {"captcha-test-response"},
				}
				login_result := Captcha_issue_2669_test_post_login(t, router_client, router_server.URL, form)
				if !strings.Contains(login_result, tool.Get_language(db, "captcha_config_error", true)) {
					t.Fatal("unknown provider did not show configuration recovery guidance")
				}
				if hosts := verification_transport.snapshot(); len(hosts) != 0 {
					t.Fatalf("unknown provider triggered a verifier request: %v", hosts)
				}
			})

			t.Run("active_provider_missing_keys_fails_closed", func(t *testing.T) {
				Captcha_issue_2669_test_set_settings(t, db, map[string]string{"recaptcha_ver": "v2"})
				verification_transport.reset()
				form := url.Values{
					"id":                   {"captcha_test_user"},
					"password":             {"captcha-test-password"},
					"g-recaptcha-response": {"captcha-test-response"},
				}
				login_result := Captcha_issue_2669_test_post_login(t, router_client, router_server.URL, form)
				if !strings.Contains(login_result, tool.Get_language(db, "captcha_config_error", true)) {
					t.Fatal("incomplete provider did not show configuration recovery guidance")
				}
				if hosts := verification_transport.snapshot(); len(hosts) != 0 {
					t.Fatalf("incomplete provider triggered a verifier request: %v", hosts)
				}
			})

			t.Run("altcha_challenge_and_login_stay_local", func(t *testing.T) {
				Captcha_issue_2669_test_set_settings(t, db, map[string]string{
					"recaptcha":     "legacy-public-key",
					"sec_re":        "legacy-google-secret",
					"altcha_sec_re": "captcha-test-altcha-secret",
					"recaptcha_ver": "altcha_low",
				})
				form_response, err := router_client.Get(router_server.URL + "/login")
				if err != nil {
					t.Fatalf("get ALTCHA login form: %v", err)
				}
				form_html, read_err := io.ReadAll(form_response.Body)
				_ = form_response.Body.Close()
				if read_err != nil || form_response.StatusCode != http.StatusOK || !strings.Contains(string(form_html), "<altcha-widget") {
					t.Fatalf("ALTCHA login form status/body invalid: %d", form_response.StatusCode)
				}
				challenge, err := route.Api_captcha_challenge()
				if err != nil {
					t.Fatalf("create ALTCHA challenge through API: %v", err)
				}
				solution, err := altcha.SolveChallenge(altcha.SolveChallengeOptions{
					Challenge: challenge,
					DeriveKey: altcha.DeriveKeyPBKDF2(),
				})
				if err != nil || solution == nil {
					t.Fatalf("solve local ALTCHA test challenge: %v", err)
				}
				payload, err := json.Marshal(altcha.Payload{Challenge: challenge, Solution: *solution})
				if err != nil {
					t.Fatalf("encode ALTCHA test solution: %v", err)
				}
				verification_transport.reset()
				login_result := Captcha_issue_2669_test_post_login(t, router_client, router_server.URL, url.Values{
					"id":       {"captcha_test_user"},
					"password": {"captcha-test-password"},
					"altcha":   {base64.StdEncoding.EncodeToString(payload)},
				})
				if !strings.Contains(login_result, `<!--opennamu_redirect:"/user"-->`) {
					t.Fatal("valid local ALTCHA solution did not reach the normal login redirect")
				}
				if hosts := verification_transport.snapshot(); len(hosts) != 0 {
					t.Fatalf("ALTCHA unexpectedly contacted an external verifier: %v", hosts)
				}
				if got := Captcha_issue_2669_test_setting(t, db, "sec_re"); got != "legacy-google-secret" {
					t.Fatal("ALTCHA login changed the stored reCAPTCHA secret")
				}
			})
		})
	}
}
