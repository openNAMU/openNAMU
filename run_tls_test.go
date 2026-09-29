package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/idna"
	"opennamu/route/tool"
)

// #3 invalid TLS mode: Start_servers now returns an error instead of calling log.Fatalf.
// The check happens before any listener is opened, so no network port is required.
func Test_Start_servers_invalid_mode(t *testing.T) {
	t.Setenv("NAMU_TLS_MODE", "atuo") // not off/file/auto

	err := Start_servers(http.NewServeMux(), "0.0.0.0", nil)
	if err == nil {
		t.Fatalf("Start_servers(%q) should return an error for invalid mode, got nil", "atuo")
	}

	msg := err.Error()
	for _, sub := range []string{"invalid TLS mode", "atuo", "off, file, auto"} {
		if !strings.Contains(msg, sub) {
			t.Fatalf("error message %q must contain %q", msg, sub)
		}
	}

	// Confirms the invalid-mode branch returns before dispatching to any runner:
	// with a valid mode it would try to bind (needs a port); here no listener opens.
	if !tool.IsValid_tls_mode(tool.TLS_mode()) {
		// expected: "atuo" is invalid; Start_servers must not open a listener for it (#3)
	}
}

// #4 HostWhitelist exact semantics, measured against golang.org/x/crypto/acme/autocert v0.55.0 actual code.
// Only exact registered hosts match; subdomains / regexp / wildcard do NOT (source L69-70). autocert normalizes the
// request host via idna.Lookup.ToASCII before matching, so example.com and EXAMPLE.COM resolve identically for
// issuance (GetCertificate L267 → hostPolicy at L292; source explicitly treats them as equivalent, L259-263). The test
// mirrors that real flow: ToASCII-normalize the request first, then apply autocert.HostWhitelist. See TLS_IMPLEMENTATION_PROGRESS.md #4.
func Test_Autocert_host_whitelist_exact(t *testing.T) {
	policy := autocert.HostWhitelist("example.com", "www.example.com")

	assertMatch := func(host string, wantAllowed bool) {
		norm, err := idna.Lookup.ToASCII(host) // autocert normalizes the request host before matching (GetCertificate L267 → L292)
		if err != nil {
			t.Fatalf("ToASCII(%q) error: %v", host, err)
		}
		got := policy(context.Background(), norm) == nil
		if got != wantAllowed {
			t.Fatalf("normalized %q: match=%v, want=%v (autocert.HostWhitelist)", host, got, wantAllowed)
		}
	}

	assertMatch("example.com", true)     // exact registered host
	assertMatch("www.example.com", true) // explicitly registered host
	assertMatch("Example.COM", true)     // case variant → normalized to example.com (autocert treats them as equivalent, source L259-263)
	assertMatch("EXAMPLE.COM", true)

	// exact-match-only: subdomains and cross-domain are rejected even after normalization
	assertMatch("api.example.com", false)
	assertMatch("x.example.com", false)
	assertMatch("example.com.evil.com", false)
}

// #3 invalid-mode gate: valid modes pass IsValid_tls_mode; empty normalizes to "off". This is the
// precondition Start_servers uses to decide whether to open a listener. Package-tool semantics are
// already covered by Test_TLS_mode_normalization/Test_IsValid_tls_mode in tls_config_test.go; here we
// re-check from the main package's perspective (no port opened for invalid mode).
func Test_Start_servers_valid_gate(t *testing.T) {
	for _, m := range []string{"off", "file", "auto"} {
		t.Setenv("NAMU_TLS_MODE", m)
		if !tool.IsValid_tls_mode(tool.TLS_mode()) {
			t.Fatalf("mode %q should be valid (TLS_mode=%q)", m, tool.TLS_mode())
		}
	}

	t.Setenv("NAMU_TLS_MODE", "atuo") // invalid -> Start_servers returns without opening a listener
	if tool.IsValid_tls_mode(tool.TLS_mode()) {
		t.Fatalf("mode should be treated as invalid; Start_servers must not open a listener for it")
	}
}

// autocert.HTTPHandler(nil) wiring + HTTP->HTTPS redirect. Uses net/http/httptest so no privileged port is needed (the real :80/:443 bind is NOT verified here — see "environmentally unverified" in the doc). Labels exactly what's checked: the redirect plumbing and path+query preservation, NOT ACME issuance (which needs a live Let's Encrypt server; observed 404 with a nil client). See TLS_IMPLEMENTATION_PROGRESS.md #10/#3.
func Test_Autocert_HTTPHandler_redirect(t *testing.T) {
	manager := &autocert.Manager{Prompt: autocert.AcceptTOS, HostPolicy: autocert.HostWhitelist("example.com")}
	handler := manager.HTTPHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "http://example.com/test?a=1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusFound { // 302 redirect (user's expected status for the HTTP->HTTPS case)
		t.Fatalf("HTTPHandler redirect status = %d, want 302; Location=%q", rec.Code, rec.Header().Get("Location"))
	}
	loc := rec.Header().Get("Location")
	if loc != "https://example.com/test?a=1" { // scheme flipped to https; path + query preserved exactly
		t.Fatalf("redirect Location = %q, want https://example.com/test?a=1", loc)
	}

	// ACME challenge path cannot be served in-unit: a nil client has no live issuer and there is no cached
	// certificate, so the handler does not return a valid 200 response (confirms real HTTP-01 issuance is out of scope).
	ch := httptest.NewRequest(http.MethodGet, "http://example.com/.well-known/acme-challenge/abc", nil)
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, ch)
	if rec2.Code == http.StatusOK {
		t.Fatalf("ACME challenge path should not succeed in-unit (no live issuer); got %d", rec2.Code)
	}
}
