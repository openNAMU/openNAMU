package tool

import (
	"os"
	"path/filepath"
	"strings"
)

// TLS 설정은 환경 변수를 우선으로 읽으며 data/set.json을 폴백으로 사용합니다.
// 파일 경로만 저장하고 private key 내용을 절대 저장하지 않습니다.

func tls_set_path() string {
	return filepath.Join("data", "set.json")
}

func tls_read_set_value(key string) (string, bool) {
	path := tls_set_path()
	raw, err := os.ReadFile(path)
	if err != nil || len(raw) == 0 {
		return "", false
	}

	tmp := map[string]string{}
	if json.Unmarshal(raw, &tmp) != nil {
		return "", false
	}

	v, ok := tmp[key]
	return strings.TrimSpace(v), ok
}

// TLS_setting reads NAMU_TLS_<KEY>. The actual environment variable name is case-sensitive and upper-case,
// e.g. NAMU_TLS_CERT_FILE; only the internal key passed here is normalized case-insensitively (so both
// "cert_file" and "CERT_FILE" map to the same lookup). It falls back to data/set.json then default_value.
func TLS_setting(key string, default_value string) string {
	env_v := os.Getenv(strings.ToUpper("NAMU_TLS_" + strings.ToLower(key)))
	if strings.TrimSpace(env_v) == "" {
		set_v, set_ok := tls_read_set_value(key)
		if set_ok && set_v != "" {
			return Choose(set_v, default_value)
		}
	}

	return Choose(env_v, default_value)
}

func TLS_mode() string {
	m := TLS_setting("mode", "off")
	return strings.ToLower(strings.TrimSpace(m))
}

func IsValid_tls_mode(mode string) bool {
	switch mode {
	case "", "off", "file", "auto":
		return true
	default:
		return false
	}
}

// TLS_has_positional_override reports whether an arg is a positional port override.
// Flags --localhost and dev are skipped; any other arg counts as an explicit port override.
func TLS_has_positional_override(arg string) bool {
	return arg != "--localhost" && arg != "dev"
}

// TLS_bind_addr resolves the server bind address for off/file modes from CLI args.
// It is the single place that decides a port; auto uses fixed ports 80/443 (see run_tls).
func TLS_bind_addr(mode string, host string, cliArgs []string) string {
	var port string
	switch mode {
	case "off":
		port = "3000"
	case "file":
		port = "443"
	}

	for _, arg := range cliArgs {
		if TLS_has_positional_override(arg) {
			return host + ":" + arg // positional override (off/file only)
		}
	}

	return host + ":" + port
}

// Split_tls_domain parses a comma-separated domain allow-list and drops empties.
func Split_tls_domain(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))

	for _, p := range parts {
		if d := strings.TrimSpace(p); d != "" {
			out = append(out, d)
		}
	}

	return out
}

func TLS_cache_dir() string {
	return filepath.Join("data", "tls")
}
