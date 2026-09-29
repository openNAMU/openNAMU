package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/acme/autocert"
	"opennamu/route/tool"
)

const (
	tls_auto_challenge = "80"  // ACME HTTP-01 challenge / redirect port (fixed for auto mode)
	tls_auto_https     = "443" // HTTPS listener port (fixed for auto mode)
)

// Start_servers dispatches TLS modes and returns any startup error.
// It never calls log.Fatalf; the caller (main) performs process exit (#1, #2).
func Start_servers(handler http.Handler, host string, cliArgs []string) error {
	mode := tool.TLS_mode()
	if !tool.IsValid_tls_mode(mode) {
		return fmt.Errorf("invalid TLS mode: %q (allowed: off, file, auto)", mode)
	}

	switch mode {
	case "file":
		return run_file_tls(handler, host, cliArgs)
	case "auto":
		return run_auto_tls(handler, host, cliArgs)
	default: // off (and empty normalized to off by IsValid_tls_mode); keep explicit
		return run_off_http(handler, host, cliArgs)
	}
}

func run_off_http(handler http.Handler, host string, cliArgs []string) error {
	addr := tool.TLS_bind_addr("off", host, cliArgs) // 3000 default; CLI positional override allowed
	log.Printf("Run in http://%s", addr)

	if err := serveOff(addr, handler); err != nil {
		return fmt.Errorf("server failed: %w", err) // port-in-use / permission-denied propagated to Start_servers → main (#12)
	}

	return nil
}

func run_file_tls(handler http.Handler, host string, cliArgs []string) error {
	certFile := tool.TLS_setting("cert_file", "")
	keyFile := tool.TLS_setting("key_file", "")

	switch {
	case certFile == "": // env unset/empty and no set.json tls_cert_file → require explicit path (#11)
		return fmt.Errorf("TLS_MODE=file requires NAMU_TLS_CERT_FILE (or data/set.json tls_cert_file)")
	case keyFile == "":
		return fmt.Errorf("TLS_MODE=file requires NAMU_TLS_KEY_FILE (or data/set.json tls_key_file)")
	}

	for _, f := range []string{certFile, keyFile} {
		if !tool.File_exist_check(f) { // preserve which path is wrong; cert vs key distinct messages (#11)
			log.Printf("TLS certificate file not found: %s", f)
			return fmt.Errorf("TLS certificate file not found: %s", f)
		}
	}

	addr := tool.TLS_bind_addr("file", host, cliArgs) // 443 default; CLI positional override allowed
	log.Printf("Run in https://%s (manual PEM)", addr)

	srv := &http.Server{Addr: addr, Handler: handler}
	if err := serveTLS(srv, certFile, keyFile); err != nil {
		return fmt.Errorf("server failed: %w", err) // ListenAndServeTLS failure propagated to Start_servers → main (#12)
	}

	return nil
}

func run_auto_tls(handler http.Handler, host string, cliArgs []string) error {
	domains := tool.Split_tls_domain(tool.TLS_setting("domain", ""))
	if len(domains) == 0 {
		return fmt.Errorf("TLS_MODE=auto requires NAMU_TLS_DOMAIN (comma-separated list of allowed domains)")
	}

	for _, d := range domains { // reject obviously-malformed entries so a bad config fails fast with a clear cause instead of being silently misused by HostWhitelist (#5); keeps the existing error path, no new abstraction
		if strings.ContainsAny(d, "/:") || strings.HasPrefix(d, "http://") || strings.HasPrefix(d, "https://") {
			return fmt.Errorf("invalid TLS domain: %q; expected a hostname such as example.com", d)
		}
	}

	manager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache(tool.TLS_cache_dir()),
		HostPolicy: autocert.HostWhitelist(domains...), // exact allow-list (no custom subdomain logic)
		Email:      tool.TLS_setting("email", ""),
	}

	override := false
	for _, arg := range cliArgs {
		if tool.TLS_has_positional_override(arg) {
			override = true
			break
		}
	}
	if override {
		log.Printf("WARNING: TLS auto mode uses fixed ports %s and %s; CLI port override is ignored", tls_auto_challenge, tls_auto_https)
	}

	if host == "127.0.0.1" { // bound to localhost only (from --localhost); public ACME validation normally unreachable (#11)
		log.Printf("WARNING: TLS auto mode is bound to localhost.")
		log.Printf("Ports 80 (HTTP-01 and HTTP-to-HTTPS redirects) and 443 (TLS-ALPN-01 and HTTPS) are used for automatic issuance; they should normally be reachable from the Internet")
	}

	err_ch := make(chan error, 2) // collect errors from BOTH servers (fail-fast if either fails) (#6)

	go func() {
		var err error
		srv := &http.Server{Addr: host + ":" + tls_auto_https, Handler: handler, TLSConfig: manager.TLSConfig()}
		if err = serveTLS(srv, "", ""); !errors.Is(err, http.ErrServerClosed) {
			log.Printf("error: https server failed: %v", err)
		}

		err_ch <- err
	}()

	go func() {
		var err error
		challengeAddr := host + ":" + tls_auto_challenge // ACME HTTP-01 challenge / redirect to https (#7)
		if err = serveHTTPChallenge(challengeAddr, manager); !errors.Is(err, http.ErrServerClosed) {
			log.Printf("error: http server failed: %v", err)
		}

		err_ch <- err
	}()

	select {
	case err := <-err_ch: // block until either server exits (or both close normally)
		if errors.Is(err, http.ErrServerClosed) { // normal shutdown of a listener (#2)
			return nil
		}
		log.Printf("auto TLS server failed: %v", err)
		return fmt.Errorf("auto TLS server failed: %w", err)
	}
	// Function ends here; the buffered err_ch (buffer 2, one send per goroutine that always runs) is
	// never exhausted before select returns. No trailing return needed (#1/#2).
}

func serveOff(addr string, handler http.Handler) error {
	srv := &http.Server{Addr: addr, Handler: handler}
	return srv.ListenAndServe()
}

func serveHTTPChallenge(addr string, manager *autocert.Manager) error {
	srv := &http.Server{Addr: addr, Handler: manager.HTTPHandler(nil)}
	return srv.ListenAndServe()
}

func serveTLS(srv *http.Server, certFile, keyFile string) error {
	return srv.ListenAndServeTLS(certFile, keyFile)
}
