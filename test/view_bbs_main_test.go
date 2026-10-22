package test

import (
	"strings"
	"testing"

	"opennamu/route"
)

// Test_Search_highlight verifies the search highlight never marks inside HTML
// entities (e.g. the `amp` of `&amp;`), and that all non-matching and matching
// regions are escaped while only matches are wrapped in <mark>. Covers the
// title/body-snippet/tag search inputs which all call Search_highlight.
func Test_Search_highlight(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		key     string
		want    string
	}{
		{
			name: "entity amp not marked",
			data: "Amp & Co",
			key:  "amp",
			want: `<mark>Amp</mark> &amp; Co`,
		},
		{
			name: "quotes and angle brackets escaped, match only keyword",
			data: `It's "a" <b>`,
			key:  "a",
			want: `It&#39;s &#34;<mark>a</mark>&#34; &lt;b&gt;`,
		},
		{
			name: "korean keyword",
			data: "한국어 검색",
			key:  "검색",
			want: "한국어 <mark>검색</mark>",
		},
		{
			name: "case-insensitive multiple matches",
			data: "Hello hello HELLO",
			key:  "hello",
			want: `<mark>Hello</mark> <mark>hello</mark> <mark>HELLO</mark>`,
		},
		{
			name: "regex dot treated literally",
			data: "a.b c.b",
			key:  "a.b",
			want: `<mark>a.b</mark> c.b`,
		},
		{
			name: "multiple keywords with special chars",
			data: "a.b and a_b",
			key:  "a.b a_b",
			want: `<mark>a.b</mark> and <mark>a_b</mark>`,
		},
		{
			name: "no match keeps escaped original",
			data: "A & Co",
			key:  "xyz",
			want: "A &amp; Co",
		},
		{
			name: "no match korean",
			data: "한국어",
			key:  "abc",
			want: "한국어",
		},
		{
			name: "ampersand keyword marked as entity",
			data: "X & Y",
			key:  "&",
			want: `X <mark>&amp;</mark> Y`,
		},
		{
			name: "empty keyword",
			data: "abc",
			key:  "",
			want: "abc",
		},
		{
			name: "whitespace-only keyword",
			data: "abc",
			key:  "   ",
			want: "abc",
		},
		{
			name: "empty data",
			data: "",
			key:  "abc",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := route.Search_highlight(tt.data, tt.key)
			if got != tt.want {
				t.Fatalf("Search_highlight(%q, %q) = %q; want %q", tt.data, tt.key, got, tt.want)
			}
		})
	}

	// HTML tag-like input must be escaped, not become a real <b> tag.
	t.Run("html tags not interpreted", func(t *testing.T) {
		want := "&lt;<mark>b</mark>&gt;<mark>b</mark>old&lt;/<mark>b</mark>&gt;"
		got := route.Search_highlight("<b>bold</b>", "b")
		if got != want {
			t.Fatalf("html tag input exact = %q; want %q", got, want)
		}
		if !strings.Contains(got, "&lt;") || !strings.Contains(got, "&gt;") {
			t.Fatalf("angle brackets not escaped: %q", got)
		}
		if !strings.Contains(got, "<mark>b</mark>") {
			t.Fatalf("expected <mark> around b: %q", got)
		}
		if strings.Contains(got, "<b>") || strings.Contains(got, "</b>") {
			t.Fatalf("raw <b> tag produced: %q", got)
		}
	})

	// Title-search preview: raw title from DB, with an entity.
	t.Run("title search preview", func(t *testing.T) {
		got := route.Search_highlight("Amp & Co", "amp")
		if !strings.Contains(got, `<mark>Amp</mark>`) {
			t.Fatalf("title not highlighted: %q", got)
		}
		if strings.Contains(got, `&<mark>`) {
			t.Fatalf("entity broken by mark: %q", got)
		}
	})

	// Body-snippet preview: raw body text with special chars.
	t.Run("snippet body preview", func(t *testing.T) {
		got := route.Search_highlight(`He said "hi" & <x>y</x>`, "y")
		if !strings.Contains(got, `<mark>y</mark>`) {
			t.Fatalf("snippet y not highlighted: %q", got)
		}
		if strings.Contains(got, `<x>`) || strings.Contains(got, `</x>`) {
			t.Fatalf("raw tags in snippet: %q", got)
		}
		if !strings.Contains(got, "&amp;") {
			t.Fatalf("entity not preserved in snippet: %q", got)
		}
	})

	// Tag highlight: short tags with special chars.
	t.Run("tag highlight", func(t *testing.T) {
		got := route.Search_highlight("a<b>", "a")
		if got != `<mark>a</mark>&lt;b&gt;` {
			t.Fatalf("tag highlight = %q; want <mark>a</mark>&lt;b&gt;", got)
		}
	})
}
