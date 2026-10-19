package test

import (
	"testing"

	"opennamu/route/tool/markup/backlink"
)

// Test_Normalize_external_host_bare_hostname is a regression test for the V5
// bug where host-based query inputs are bare host names (no http/https scheme)
// but backlink.Normalize_external_host required a scheme, returning "" and making the
// host->docs and host->count APIs report zero results.
func Test_Normalize_external_host_bare_hostname(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"bare full host", "example.com", "example.com"},
		{"bare subdomain", "sub.example.com", "sub.example.com"},
		{"bare with scheme", "https://example.com", "example.com"},
		{"bare uppercase", "EXAMPLE.COM", "example.com"},
		{"bare idn", "예.커", "xn--2j5b.xn--ve7b"},
		{"ip rejected", "1.2.3.4", ""},
		{"ipv4 localhost rejected", "localhost", ""},
		{"trailing dot stripped", "example.com.", "example.com"},
	}
	for _, c := range cases {
		got := backlink.Normalize_external_host(c.input)
		if got != c.want {
			t.Errorf("backlink.Normalize_external_host(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}
