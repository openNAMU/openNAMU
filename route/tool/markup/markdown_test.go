package markup

import (
	"database/sql"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

func markdownTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	if _, err = db.Exec("create table other (name text, data text)"); err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec("insert into other values ('language', 'ko-KR')"); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})
	return db
}

func TestMarkdownRendersThroughMacromark(t *testing.T) {
	result := Markdown(markdownTestDB(t), map[string]string{
		"doc_name": "test",
		"data":     "# Title (test)\n\n**bold** [link](https://example.com/a_(b)) [path](/help) [anchor](#top) [mail](mailto:test@example.com)",
	})
	rendered := result["data"].(string)

	for _, expected := range []string{
		"<h1><a href=\"#toc\">1. </a>Title (test)</h1>",
		"<b>bold</b>",
		"<a class=\"opennamu_link_out\" href=\"https://example.com/a_(b)\">link</a>",
		"<a href=\"/help\">path</a>",
		"<a href=\"#top\">anchor</a>",
		"mail",
	} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("rendered Markdown does not contain %q: %s", expected, rendered)
		}
	}

	if strings.Contains(rendered, "mailto:") {
		t.Fatalf("mailto link was rendered: %s", rendered)
	}
}

func TestMarkdownRendersBlocksAndEscapesCode(t *testing.T) {
	result := Markdown(markdownTestDB(t), map[string]string{
		"doc_name": "test",
		"data":     "~~~go\n<hello>\n~~~\n\n- one\n- two\n\n> quoted\n\n---",
	})
	rendered := result["data"].(string)

	for _, expected := range []string{
		"<pre><code class=\"language-go\">&lt;hello&gt;",
		"<ul><li>one</li><li>two</li></ul>",
		"<q>quoted</q>",
		"<hr>",
	} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("rendered Markdown does not contain %q: %s", expected, rendered)
		}
	}
}

func TestMarkdownEscapesRawHTML(t *testing.T) {
	result := Markdown(nil, map[string]string{
		"doc_name": "test",
		"data":     "<script>alert(1)</script>",
	})
	rendered := result["data"].(string)

	if strings.Contains(rendered, "<script>") || !strings.Contains(rendered, "&lt;script&gt;") {
		t.Fatalf("raw HTML was not escaped: %s", rendered)
	}
}

func TestMacromarkRejectsMailto(t *testing.T) {
	renderer := Macromark_new(nil, map[string]string{
		"data": "[ex(mailto:test@example.com,mail)]",
	}, "html")
	rendered := renderer.main()["data"].(string)

	if strings.Contains(rendered, "mailto:") || strings.Contains(rendered, "href=") {
		t.Fatalf("mailto link was rendered: %s", rendered)
	}
	if rendered != "mail" {
		t.Fatalf("unexpected mailto fallback: %s", rendered)
	}
}
