package test

import (
	"strings"
	"testing"

	"opennamu/route/tool"
)

// Test_get_template_og_meta verifies that Get_template emits og:title and
// og:description meta tags (reusing the existing page title and description)
// in the rendered head, while preserving the existing og:image tag and
// HTML-escaping values.
func Test_get_template_og_meta(t *testing.T) {
	db := Setup_rss_db(t, false)
	defer tool.DB_close(db)

	config := tool.Config{IP: "192.168.0.1", Session: &emptySession{}}

	title := "문서 <test> & \"quote\""
	description := "설명 & 특수문자"
	og_image := "https://example.com/img.jpg"

	out := tool.Get_template(
		db,
		config,
		title,
		"body",
		[]any{"", "", "", description, "0"},
		[][]any{},
		map[string]string{"og_image": og_image},
	)
	if !strings.Contains(out, `<meta property="og:title" content="문서 &lt;test&gt; &amp; &#34;quote&#34;">`) {
		t.Fatalf("og:title missing or not escaped correctly")
	}
	if !strings.Contains(out, `<meta property="og:description" content="설명 &amp; 특수문자">`) {
		t.Fatalf("og:description missing or not escaped correctly")
	}
	if !strings.Contains(out, `<meta property="og:image" content="https://example.com/img.jpg">`) {
		t.Fatalf("existing og:image tag lost")
	}
}

// Test_get_template_og_meta_no_description verifies that when there is no
// description, no og:description tag is emitted (empty description is skipped),
// while og:title is still emitted.
func Test_get_template_og_meta_no_description(t *testing.T) {
	db := Setup_rss_db(t, false)
	defer tool.DB_close(db)

	config := tool.Config{IP: "192.168.0.1", Session: &emptySession{}}

	title := "문서"

	out := tool.Get_template(
		db,
		config,
		title,
		"body",
		[]any{"", "", "", "", "0"},
		[][]any{},
		map[string]string{},
	)
	if !strings.Contains(out, `<meta property="og:title" content="문서">`) {
		t.Fatalf("og:title missing")
	}
	if strings.Contains(out, `property="og:description"`) {
		t.Fatalf("og:description emitted even though description is empty")
	}
}
