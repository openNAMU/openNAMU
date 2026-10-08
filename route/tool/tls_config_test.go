package tool

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_TLS_setting_priority(t *testing.T) {
	t.Setenv("NAMU_TLS_FOO", "") // empty env => fall through to default
	if got := TLS_setting("foo", "def"); got != "def" {
		t.Fatalf("empty env should use default, got %q", got)
	}

	t.Setenv("NAMU_TLS_FOO", "from-env")
	if got := TLS_setting("foo", "def"); got != "from-env" {
		t.Fatalf("env value not used, got %q", got)
	}
}

func Test_TLS_mode_normalization(t *testing.T) {
	cases := map[string]string{
		"":      "off",  // unset -> default off
		"AUTO":  "auto", // case-insensitive
		" Off ": "off",  // whitespace trimmed
	}

	for input, want := range cases {
		t.Setenv("NAMU_TLS_MODE", input)
		if got := TLS_mode(); got != want {
			t.Fatalf("TLS_mode(%q) = %q, want %q", input, got, want)
		}
	}
}

func Test_IsValid_tls_mode(t *testing.T) {
	valids, invalids := []string{"off", "file", "auto"}, []string{"weird", "HTTP"}
	for _, v := range valids {
		if !IsValid_tls_mode(v) {
			t.Fatalf("IsValid_tls_mode(%q) = false, want true", v)
		}
	}
	for _, v := range invalids {
		if IsValid_tls_mode(v) {
			t.Fatalf("IsValid_tls_mode(%q) = true, want false", v)
		}
	}
}

// Test_TLS_bind_addr covers off/file default ports and CLI positional override.
// Auto uses fixed ports 80/443 (run_auto_tls); TLS_bind_addr is not used for auto (#4).
func Test_TLS_bind_addr(t *testing.T) {
	cases := []struct {
		mode   string
		cli    []string
		expect string
	}{
		{"off", nil, "0.0.0.0:3000"},                     // default port
		{"file", nil, "0.0.0.0:443"},                     // default port
		{"off", []string{"8080"}, "0.0.0.0:8080"},        // CLI override allowed
		{"file", []string{"8443"}, "0.0.0.0:8443"},       // CLI override allowed
		{"off", []string{"--localhost"}, "0.0.0.0:3000"}, // flag skipped, default kept
	}

	for _, c := range cases {
		if got := TLS_bind_addr(c.mode, "0.0.0.0", c.cli); got != c.expect {
			t.Fatalf("TLS_bind_addr(%q,%v) = %q, want %q", c.mode, c.cli, got, c.expect)
		}
	}
}

func Test_TLS_has_positional_override(t *testing.T) {
	cases := map[string]bool{
		"--localhost": false, // host flag => not an override (skipped in TLS_bind_addr / auto warn)
		"dev":         false, // dev mode flag => not an override
		"":            true,  // empty positional => counts as override (off/file would use it)
		"8080":        true,  // explicit port => override
	}

	for arg, want := range cases {
		if got := TLS_has_positional_override(arg); got != want {
			t.Fatalf("TLS_has_positional_override(%q) = %v, want %v", arg, got, want)
		}
	}
}

func Test_Split_tls_domain(t *testing.T) {
	got := Split_tls_domain("a.com, b.com ,,c.com")
	want := []string{"a.com", "b.com", "c.com"}
	if len(got) != len(want) {
		t.Fatalf("Split_tls_domain len = %d, want %d (got %v)", len(got), len(want), got)
	}

	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("Split_tls_domain[%d] = %q, want %q", i, got[i], want[i])
		}
	}

	if got := Split_tls_domain(""); len(got) != 0 {
		t.Fatalf("empty domain list should be empty, got %v", got)
	}

	// Explicit subdomain entries are kept as-is; allow-list enforcement (exact match)
	// is delegated to autocert.HostWhitelist(domains...) in run_auto_tls.
	got = Split_tls_domain("example.com,www.example.com")
	want = []string{"example.com", "www.example.com"}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("Split_tls_domain subdomain entries not preserved: %v", got)
	}
}

// #6 empty env / set.json fallback semantics. data/set.json has no tls_mode key, so unset/empty env
// falls back to the default ("off"); a non-empty env always wins over the default (#6). The env is the
// explicit value; only an empty env triggers the set.json/default path (TLS_setting: env non-empty →
// returned directly; empty → Tls_read_set_value then default). This test uses t.Setenv only, so it does
// not touch the user's data/set.json. Actual set.json-present behavior is documented as runtime note (#6).
func Test_TLS_mode_env_priority_and_default(t *testing.T) {
	// unset/empty env → no tls_mode in data/set.json → default "off" (#6 row 1 & 2)
	t.Setenv("NAMU_TLS_MODE", "")
	if got := TLS_mode(); got != "off" {
		t.Fatalf("unset/empty NAMU_TLS_MODE should fall back to set.json/default (off), got %q", got)
	}

	// non-empty env always wins over the default, regardless of case (#6: explicit env value)
	for _, v := range []string{"file", "auto", "OFF"} {
		t.Setenv("NAMU_TLS_MODE", v)
		if got := TLS_mode(); strings.ToLower(got) != strings.ToLower(v) {
			t.Fatalf("explicit NAMU_TLS_MODE=%q should win, got %q (default)", v, got)
		}
	}

	// IsValid_tls_mode gate used by Start_servers to decide whether to open a listener (#3/#6):
	if !IsValid_tls_mode("") || !IsValid_tls_mode("off") || !IsValid_tls_mode("file") || !IsValid_tls_mode("auto") {
		t.Fatalf("empty/off/file/auto must be valid modes")
	}
	if IsValid_tls_mode("atuo") {
		t.Fatalf("'atuo' must be invalid so Start_servers returns an error instead of opening a listener")
	}
}

// set.json fallback + empty-env semantics, verified against the REAL data-set.json read path (TLS_setting -> Tls_read_set_value reads "data/set.json" relative to CWD).
// Uses os.Chdir into t.TempDir() so the user's real data/set.json is never touched or overwritten. Go runs tests sequentially (no Parallel), so restoring CWD via defer is safe.
// Locks every row of the openNAMU env > set.json > default convention for TLS_MODE, including the key empty-env case (row 4) that differs from a naive "presence-wins" LookupEnv loader like Get_DB_set. See TLS_IMPLEMENTATION_PROGRESS.md #2.
func Test_TLS_setting_fallback_semantics(t *testing.T) {
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd: %v", err)
	}
	dir := t.TempDir()
	dataDir := filepath.Join(dir, "data")
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		t.Fatalf("MkdirAll data dir: %v", err)
	}
	setPath := filepath.Join(dataDir, "set.json")

	// chdir into the temp dir so TLS_setting reads its data/set.json via the relative path ("data/set.json");
	// the user's real data/set.json is never touched. Go runs tests sequentially (no Parallel), so restoring CWD at end is safe.
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("os.Chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(origWd) })

	writeSet := func(content string) {
		if content == "" {
			_ = os.Remove(setPath)
		} else {
			_ = os.WriteFile(setPath, []byte(content), 0o644)
		}
	}

	run := func(setContent, envValue string) string {
		writeSet(setContent)
		if envValue != "<unset>" {
			os.Setenv("NAMU_TLS_MODE", envValue)
		} else {
			os.Unsetenv("NAMU_TLS_MODE")
		}
		return TLS_setting("mode", "off") // default_value is off, matching the real convention (TLS_mode fallback)
	}

	cases := []struct {
		name string
		set  string // set.json content ("" = no file at all)
		env  string // "<unset>" => env not present; "" or other => explicit value
		want string
	}{
		{"set.json no key + env unset -> off", `{}`, "<unset>", "off"},                                                           // default
		{"set.json no key + env empty -> off (empty falls back to set.json/default)", ``, "", "off"},                             // empty env triggers fallback
		{"set.json mode=auto + env unset -> auto", `{"mode":"auto"}`, "<unset>", "auto"},                                         // env absent => use set.json
		{"set.json mode=auto + env empty -> auto (empty falls back to set.json)", `{"mode":"auto"}`, "", "auto"},                 // KEY row: empty env beats default, uses set.json value
		{"set.json mode=auto + env off -> off (non-empty env beats both set.json and default)", `{"mode":"auto"}`, "off", "off"}, // explicit env
		{"set.json mode=off + env unset -> off", `{"mode":"off"}`, "<unset>", "off"},                                             // set.json value used
	}

	for _, c := range cases {
		if got := run(c.set, c.env); got != c.want {
			t.Fatalf("case %q: TLS_setting(mode) = %q, want %q (set=%q env=%q)", c.name, got, c.want, c.set, c.env)
		}
	}

	// sanity: a non-TLS key still falls through to its default when unset/empty and is overridden by an explicit value (#2/#6)
	if got := run("", "<unset>"); got != "off" {
		t.Fatalf("non-present mode should be default off, got %q", got)
	}
	if got := run(`{"mode":"auto"}`, "file"); got != "file" {
		t.Fatalf("explicit env file must beat set.json auto, got %q", got)
	}
}
