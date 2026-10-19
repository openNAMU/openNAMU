package test

import (
	"opennamu/route/tool/markup/backlink"
	"strings"
	"testing"
)

// Test_Normalize_external_host covers the V5 hostname rules: http/https only,
// hostname extracted (scheme/path/query/fragment/port dropped), lowercased,
// trailing DNS dot stripped, userinfo/invalid/IPv4/IPv6/localhost rejected.
func Test_Normalize_external_host(t *testing.T) {
	total253 := strings.Join([]string{
		strings.Repeat("a", 63),
		strings.Repeat("a", 63),
		strings.Repeat("a", 63),
		strings.Repeat("a", 57),
		"com",
	}, ".")
	total254 := strings.Join([]string{
		strings.Repeat("a", 63),
		strings.Repeat("a", 63),
		strings.Repeat("a", 63),
		strings.Repeat("a", 58),
		"com",
	}, ".")
	cases := []struct {
		name string
		in   string
		want string
	}{
		// V5 spec table (hostname, not full URL).
		{"canonical hostname", "https://EXAMPLE.COM/Article", "example.com"},
		{"http query fragment path", "http://EXAMPLE.com/path?query=1#frag", "example.com"},
		{"path case ignored", "https://Example.Com/path", "example.com"},
		{"www is a distinct domain", "https://www.example.com/Page", "www.example.com"},
		{"multi-level TLD", "https://example.co.uk/Page", "example.co.uk"},
		{"subdomain is a distinct domain", "https://sub.domain.example.com/a/b", "sub.domain.example.com"},
		{"port not stored", "https://example.com:443/a", "example.com"},
		{"valid port not stored", "https://example.com:8443/a", "example.com"},
		{"trailing dot stripped", "https://example.com./a", "example.com"},
		{"trailing dot multi label", "https://sub.example.co.uk./a", "sub.example.co.uk"},
		{"userinfo excluded", "https://user:pass@example.com/a", ""},
		{"userinfo single name excluded", "https://user@sub.example.com/a", ""},
		{"bad port (99999)", "http://example.com:99999/...", ""},
		{"bad port (0)", "http://example.com:0/a", ""},
		{"bad port (65536)", "http://example.com:65536/a", ""},
		// DNS-only: IP addresses and localhost are not domains.
		{"ipv4 excluded", "https://192.168.1.1/a", ""},
		{"ipv4 private excluded", "https://10.0.0.1/a", ""},
		{"ipv6 excluded", "https://[::1]/a", ""},
		{"ipv6 excluded 2", "https://[2001:db8::1]/a", ""},
		{"localhost excluded", "https://localhost/a", ""},
		{"valid two-label domain", "http://example.net/a", "example.net"},
		// Invalid DNS shape.
		{"leading hyphen excluded", "https://-example.com/a", ""},
		{"trailing hyphen label excluded", "https://example-.com/a", ""},
		{"consecutive dots excluded", "https://example..com/a", ""},
		{"empty host excluded", "http://", ""},
		{"no host https excluded", "https://", ""},
		{"protocol relative excluded", "//example.com/a", ""},
		{"no scheme excluded", "example.com/a", ""},
		{"non http scheme ftp excluded", "ftp://example.com/a", ""},
		{"non http scheme mailto excluded", "mailto:user@example.com", ""},
		{"empty string excluded", "", ""},
		{"percent encoded host invalid", "https://exam%2Ecom/a", ""},
		// IDN is normalized to Punycode (golang.org/x/net/idna.ToASCII),
		// not rejected. Mixed case is lower-cased.
		{"idn normalized to punycode", "https://예.커/a", "xn--2j5b.xn--ve7b"},
		{"idn multi-label normalized to punycode", "https://한국어.도메인.kr/a", "xn--3e0bk47br7k.xn--hq1bm8jm9l.kr"},
		{"idn mixed case lowercased", "https://예.커.KR/a", "xn--2j5b.xn--ve7b.kr"},
		// Multi-level public suffixes are stored as the full hostname (no
		// eTLD+1 reduction; publicsuffix is not used for storage).
		{"co.kr stored as full hostname", "https://example.co.kr/a", "example.co.kr"},
		{"co.uk stored as full hostname", "https://example.co.uk/a", "example.co.uk"},
		{"subdomain of co.uk stored as full hostname", "https://sub.example.co.uk/a", "sub.example.co.uk"},
		// DNS length limits (checked on the ASCII/Punycode host): total 253,
		// per-label 63.
		{"label 64 too long", "https://" + strings.Repeat("a", 64) + ".com/a", ""},
		{"label 63 ok", "https://" + strings.Repeat("a", 63) + ".com/a", strings.Repeat("a", 63) + ".com"},
		{"total 253 ok", "https://" + total253 + "/a", total253},
		{"total 254 too long", "https://" + total254 + "/a", ""},
	}
	for _, c := range cases {
		if got := backlink.Normalize_external_host(c.in); got != c.want {
			t.Errorf("backlink.Normalize_external_host(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// Test_Normalize_external_host_www_distinct verifies that a host with and
// without "www" are distinct (no public-suffix stripping).
func Test_Normalize_external_host_www_distinct(t *testing.T) {
	if got := backlink.Normalize_external_host("https://www.example.com/"); got != "www.example.com" {
		t.Errorf("www example = %q, want %q", got, "www.example.com")
	}
	if got := backlink.Normalize_external_host("https://example.com/"); got != "example.com" {
		t.Errorf("plain example = %q, want %q", got, "example.com")
	}
}

// Test_Get_backlink_markdown covers markdown bracketed and autolink external
// links (now hostname-based) plus code-block exclusion and the V4 "s"
// autolink regex bug fix.
func Test_Get_backlink_markdown(t *testing.T) {
	cases := []struct {
		name  string
		data  string
		want  map[string][]string
		count int
	}{
		{
			name:  "inline external link",
			data:  `[a](https://Example.Com/Page#frag)`,
			want:  map[string][]string{"example.com": {"external"}},
			count: 0,
		},
		{
			name:  "inline link no fragment",
			data:  `[a](https://example.com/Page)`,
			want:  map[string][]string{"example.com": {"external"}},
			count: 0,
		},
		{
			name:  "autolink",
			data:  `<https://example.com/article?b=1&a=2>`,
			want:  map[string][]string{"example.com": {"external"}},
			count: 0,
		},
		{
			name:  "autolink with s in path is matched (V4 regex bug fix)",
			data:  `<https://example.com/with-s>`,
			want:  map[string][]string{"example.com": {"external"}},
			count: 0,
		},
		{
			name:  "autolink with s in host is matched (V4 regex bug fix)",
			data:  `<https://has-s.example.com/and-1>`,
			want:  map[string][]string{"has-s.example.com": {"external"}},
			count: 0,
		},
		{
			name:  "both inline and autolink same host dedupe",
			data:  `See [a](https://example.com/Page) or <https://example.com/other>.`,
			want:  map[string][]string{"example.com": {"external"}},
			count: 0,
		},
		{
			name:  "fenced code block excluded",
			data:  "```\nhttps://example.com/hidden\n```",
			want:  map[string][]string{},
			count: 0,
		},
		{
			name:  "tilde code block excluded",
			data:  "~~~\nhttps://example.com/hidden\n~~~",
			want:  map[string][]string{},
			count: 0,
		},
		{
			name:  "inline code excluded",
			data:  "`https://example.com/inline` is hidden",
			want:  map[string][]string{},
			count: 0,
		},
		{
			name:  "image excluded (V1)",
			data:  `![alt](https://example.com/image)`,
			want:  map[string][]string{},
			count: 0,
		},
		{
			name: "external plus internal count only internal",
			data: `[ext](https://example.com/Page) and [in](category:Main)`,
			want: map[string][]string{
				"example.com":   {"external"},
				"category:Main": {"cat"},
			},
			count: 1,
		},
		{
			name:  "non-http bracketed link falls to internal (V4 behavior preserved)",
			data:  `[ext](ftp://example.com/Page)`,
			want:  map[string][]string{"ftp://example.com/Page": {""}},
			count: 1,
		},
		{
			name:  "empty data",
			data:  "",
			want:  map[string][]string{},
			count: 0,
		},
		{
			name:  "relative link stays internal",
			data:  "[rel](./foo)",
			want:  map[string][]string{},
			count: 0,
		},
		{
			name:  "non-external bracketed link becomes internal",
			data:  "[rel](internal/doc)",
			want:  map[string][]string{"internal/doc": {""}},
			count: 1,
		},
	}
	for _, c := range cases {
		got, count, supported := backlink.Get_backlink(c.data, "markdown")
		if !supported {
			t.Errorf("%s: backlink.Get_backlink supported = false", c.name)
			continue
		}
		if len(got) != len(c.want) {
			t.Errorf("%s: got %d keys, want %d keys (got keys %v, want keys %v)", c.name, len(got), len(c.want), keys(got), keys(c.want))
			continue
		}
		for k, wantTypes := range c.want {
			gotTypes, ok := got[k]
			if !ok {
				t.Errorf("%s: missing key %q in got", c.name, k)
				continue
			}
			if !Equal_strings(gotTypes, wantTypes) {
				t.Errorf("%s: got[%q] = %v, want %v", c.name, k, gotTypes, wantTypes)
			}
		}
		if count != c.count {
			t.Errorf("%s: link_count = %d, want %d", c.name, count, c.count)
		}
	}
}

func Test_Get_backlink_macromark(t *testing.T) {
	cases := []struct {
		name  string
		data  string
		want  map[string][]string
		count int
	}{
		{
			name:  "ex macro external",
			data:  `[ex(https://Example.Com/Page)]`,
			want:  map[string][]string{"example.com": {"external"}},
			count: 0,
		},
		{
			name:  "ex macro label ignored",
			data:  `[ex(https://example.com/Page,Label)]`,
			want:  map[string][]string{"example.com": {"external"}},
			count: 0,
		},
		{
			name:  "ex macro http",
			data:  `[ex(http://example.com/Page)]`,
			want:  map[string][]string{"example.com": {"external"}},
			count: 0,
		},
		{
			name:  "ex macro same host multiple URLs dedupe",
			data:  `[ex(https://example.com/a,1)] [ex(https://example.com/b,2)]`,
			want:  map[string][]string{"example.com": {"external"}},
			count: 0,
		},
		{
			name:  "a macro internal counts",
			data:  `[a(category:Main)]`,
			want:  map[string][]string{"category:Main": {"cat"}},
			count: 1,
		},
		{
			name:  "include macro internal counts",
			data:  `[include(category:Main)]`,
			want:  map[string][]string{"category:Main": {"include"}},
			count: 1,
		},
		{
			name: "ex plus a",
			data: `[ex(https://example.com/Page)] [a(category:Main)]`,
			want: map[string][]string{
				"example.com":   {"external"},
				"category:Main": {"cat"},
			},
			count: 1,
		},
		{
			name:  "a macro external URL not extracted",
			data:  "[a(https://example.com/Page)]",
			want:  map[string][]string{},
			count: 0,
		},
		{
			name:  "include macro external URL not extracted",
			data:  "[include(https://example.com/Page)]",
			want:  map[string][]string{},
			count: 0,
		},
		{
			name:  "ex macro invalid target not internal (pollution fix)",
			data:  "[ex(abc,Link)]",
			want:  map[string][]string{},
			count: 0,
		},
		{
			name:  "ex macro non-http target not internal (pollution fix)",
			data:  "[ex(ftp://example.com/a)]",
			want:  map[string][]string{},
			count: 0,
		},
		{
			name:  "ex macro localhost target not internal (pollution fix)",
			data:  "[ex(localhost)]",
			want:  map[string][]string{},
			count: 0,
		},
		{
			name:  "empty data",
			data:  "",
			want:  map[string][]string{},
			count: 0,
		},
	}
	for _, c := range cases {
		got, count, supported := backlink.Get_backlink(c.data, "macromark")
		if !supported {
			t.Errorf("%s: backlink.Get_backlink supported = false", c.name)
			continue
		}
		if len(got) != len(c.want) {
			t.Errorf("%s: got %d keys, want %d keys (got keys %v, want keys %v)", c.name, len(got), len(c.want), keys(got), keys(c.want))
			continue
		}
		for k, wantTypes := range c.want {
			gotTypes, ok := got[k]
			if !ok {
				t.Errorf("%s: missing key %q in got", c.name, k)
				continue
			}
			if !Equal_strings(gotTypes, wantTypes) {
				t.Errorf("%s: got[%q] = %v, want %v", c.name, k, gotTypes, wantTypes)
			}
		}
		if count != c.count {
			t.Errorf("%s: link_count = %d, want %d", c.name, count, c.count)
		}
	}
}

func Test_Normalize_link_internal(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"plain internal", "RealTarget", "RealTarget"},
		{"category", "category:Main", "category:Main"},
		{"category korean", "분류:Main", "category:Main"},
		// backlink.Normalize_link does not strip relative paths; that is
		// backlink.Normalize_markdown_link's job.
		{"relative kept by backlink.Normalize_link", "./foo", "./foo"},
		{"internal path", "sub/dir", "sub/dir"},
		{"external http not internal", "https://example.com/a", ""},
		{"internal hash anchor dropped", "doc#anchor", "doc"},
	}
	for _, c := range cases {
		if got := backlink.Normalize_link(c.in); got != c.want {
			t.Errorf("backlink.Normalize_link(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// Test_Normalize_markdown_link verifies markdown relative-path stripping,
// which is separate from backlink.Normalize_link.
func Test_Normalize_markdown_link(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"relative ./ stripped", "./foo", ""},
		{"relative ../ stripped", "../foo", ""},
		{"absolute kept", "category:Main", "category:Main"},
		{"internal path kept", "sub/dir", "sub/dir"},
		{"external not internal", "https://example.com/a", ""},
	}
	for _, c := range cases {
		if got := backlink.Normalize_markdown_link(c.in); got != c.want {
			t.Errorf("backlink.Normalize_markdown_link(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func keys(m map[string][]string) []string {
	var out []string
	for k := range m {
		out = append(out, k)
	}
	return out
}

func Equal_strings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
