package test

import (
	"database/sql"
	"opennamu/route/tool/markup"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"opennamu/route/tool"
)

// Setup_namumark_db builds a full schema and returns a connected sqlite DB.
func Setup_namumark_db(t *testing.T) *sql.DB {
	t.Helper()

	db_dir := filepath.Join(t.TempDir(), "db")
	if err := os.MkdirAll(db_dir, 0755); err != nil {
		t.Fatalf("mkdir db dir: %v", err)
	}
	os.Setenv("NAMU_DB", db_dir)
	os.Setenv("NAMU_DB_TYPE", "sqlite")
	defer os.Unsetenv("NAMU_DB")
	defer os.Unsetenv("NAMU_DB_TYPE")

	tool.DB_init()
	return tool.DB_connect()
}

func Count_rows(db *sql.DB, t *testing.T, query string, args ...any) int {
	var c int
	if err := db.QueryRow(query, args...).Scan(&c); err != nil {
		t.Errorf("count query failed: %v", err)
	}
	return c
}

// Test_Namumark_external_backlink_collect verifies that a namumark document
// containing an inline external link produces an external backlink row with the
// normalized host title, and that link_count is NOT incremented for it.
func Test_Namumark_external_backlink_collect(t *testing.T) {
	db := Setup_namumark_db(t)
	defer tool.DB_close(db)

	doc_name := "ExtLinkDoc"
	content := "intro line\n[[https://Example.Com/Page|label text]]\nfinal"

	if _, err := db.Exec("insert into data (title, data) values (?, ?)", doc_name, content); err != nil {
		t.Fatalf("seed: %v", err)
	}

	markup.Render_namumark_compat_backlink(db, doc_name, content)

	if Count_rows(db, t, "select count(*) from back where title=? and link=? and type='external'", "example.com", doc_name) != 1 {
		t.Fatalf("external backlink not persisted for %s", doc_name)
	}
	if c := Count_rows(db, t, "select count(*) from back where link=? and type='external'", doc_name); c != 1 {
		t.Fatalf("expected exactly 1 external backlink, got %d", c)
	}

	var lc int
	if err := db.QueryRow("select set_data from data_set where doc_name=? and set_name='link_count'", doc_name).Scan(&lc); err == nil && lc != 0 {
		t.Fatalf("link_count should be 0 for an external-only document, got %d", lc)
	}
}

// Test_Namumark_external_link_with_anchor_collect verifies that a namumark
// external link with query and fragment normalizes to just the host (the
// host is the only thing stored in V5).
func Test_Namumark_external_link_with_anchor_collect(t *testing.T) {
	db := Setup_namumark_db(t)
	defer tool.DB_close(db)

	doc_name := "AnchorDoc"
	content := "[[https://example.com/a?b=1#c]]\n"

	if _, err := db.Exec("insert into data (title, data) values (?, ?)", doc_name, content); err != nil {
		t.Fatalf("seed: %v", err)
	}

	markup.Render_namumark_compat_backlink(db, doc_name, content)

	if Count_rows(db, t, "select count(*) from back where title=? and link=? and type='external'", "example.com", doc_name) != 1 {
		t.Fatalf("external backlink with host-normalized title not persisted for %s", doc_name)
	}
}

// Test_Namumark_external_redirect_collect verifies that a namumark document
// whose first line redirects to an external URL produces an external backlink
// row (the redirect branch calls Add_backlink before the collect_only
// early-return). The host is normalized (case lowercased) and stored.
func Test_namumark_redirect_first_line_differential_corpus(t *testing.T) {
	db := Setup_namumark_db(t)
	defer tool.DB_close(db)

	for _, test_case := range []struct {
		name          string
		input         string
		render_type   string
		collect       bool
		want_data     string
		want_redirect bool
	}{
		{name: "external redirect and trailing lines", input: "#redirect https://Example.com/path#section\ntail", render_type: "view", want_data: `<a href="https://Example.com/path#section">(GO)</a><br>tail`, want_redirect: true},
		{name: "case-insensitive piped internal redirect", input: "  #ReDiReCt [[Target|label]]  \ntail", render_type: "view", want_data: `<a href="/w_from/Target">(GO)</a><br>tail`, want_redirect: true},
		{name: "Korean redirect and anchor", input: "#넘겨주기 Target#anchor\ntail", render_type: "view", want_data: `<a href="/w_from/Target#anchor">(GO)</a><br>tail`, want_redirect: true},
		{name: "later line is not a redirect", input: "body\n#redirect Target", render_type: "view", want_data: "body<br>#redirect Target"},
		{name: "limited rendering leaves redirect text", input: "#redirect Target\ntail", render_type: "thread", want_data: "#redirect Target<br>tail"},
		{name: "collect-only records redirect without output", input: "#redirect https://Example.com/path#section\ntail", render_type: "backlink", collect: true, want_redirect: true},
	} {
		t.Run(test_case.name, func(t *testing.T) {
			result := markup.Render_namumark_compat(db, "RedirectCorpus", test_case.input, test_case.render_type, nil, 0, test_case.collect)
			got_data, _ := result["data"].(string)
			got_redirect, _ := result["redirect"].(bool)
			if got_data != test_case.want_data {
				t.Fatalf("rendered data = %q, want %q", got_data, test_case.want_data)
			}
			if got_redirect != test_case.want_redirect {
				t.Fatalf("redirect flag = %v, want %v", got_redirect, test_case.want_redirect)
			}
		})
	}
}

func Test_Namumark_external_redirect_collect(t *testing.T) {
	db := Setup_namumark_db(t)
	defer tool.DB_close(db)

	doc_name := "RedirectDoc"
	content := "#redirect https://Example.Com/Article#frag\ntail"

	if _, err := db.Exec("insert into data (title, data) values (?, ?)", doc_name, content); err != nil {
		t.Fatalf("seed: %v", err)
	}

	markup.Render_namumark_compat_backlink(db, doc_name, content)

	if Count_rows(db, t, "select count(*) from back where title=? and link=? and type='external'", "example.com", doc_name) != 1 {
		t.Fatalf("external redirect backlink not persisted for %s", doc_name)
	}
}

// Test_Namumark_external_and_internal_counts verifies that a document with both
// internal and external links has both backlinks collected, and link_count
// reflects only the internal link.
func Test_Namumark_external_and_internal_counts(t *testing.T) {
	db := Setup_namumark_db(t)
	defer tool.DB_close(db)

	if _, err := db.Exec("insert into data (title, data) values ('RealTarget','target content')"); err != nil {
		t.Fatalf("seed: %v", err)
	}
	doc_name := "MixedDoc"
	content := "intro\n[[RealTarget]]\n[[https://example.com/Other]]\n"
	if _, err := db.Exec("insert into data (title, data) values (?, ?)", doc_name, content); err != nil {
		t.Fatalf("seed: %v", err)
	}

	markup.Render_namumark_compat_backlink(db, doc_name, content)

	if Count_rows(db, t, "select count(*) from back where title=? and link=? and type=''", "RealTarget", doc_name) != 1 {
		t.Fatalf("internal backlink to RealTarget not persisted")
	}
	if Count_rows(db, t, "select count(*) from back where title=? and link=? and type='external'", "example.com", doc_name) != 1 {
		t.Fatalf("external backlink to example.com not persisted")
	}

	var lc int
	if err := db.QueryRow("select set_data from data_set where doc_name=? and set_name='link_count'", doc_name).Scan(&lc); err != nil {
		t.Fatalf("link_count query failed: %v", err)
	}
	if lc != 1 {
		t.Fatalf("link_count should be 1 (only the internal link), got %d", lc)
	}
}

// Test_Merge_child_skips_external verifies that Merge_child merges the child's
// internal backlink into the parent but skips the child's external backlink.
// (Merge_child is only invoked from full-render include/inter paths, so this is
// tested directly rather than through the backlink-collection path.)
func Test_Namumark_middle_preserves_known_slash_closed_block(t *testing.T) {
	renderer := markup.New_namumark_compat_renderer(nil, "MiddleTest", "", "view", nil, 0, false)
	close_token := renderer.Reserve_slash("}")
	input := "{{{#!if\nparameter==value\nbody" + close_token + "}}"

	got := renderer.Restore(renderer.Process_middle(input))
	want := "{{{#!if\nparameter==value\nbody}}}"
	if got != want {
		t.Fatalf("Process_middle() = %q, want %q", got, want)
	}
}

func Test_Namumark_middle_size_and_color_share_rendering_contract(t *testing.T) {
	for _, test_case := range []struct {
		name         string
		middle       string
		want         string
		collect_want string
	}{
		{name: "size", middle: "+2 text", want: `<span style="font-size:140%">text</span>`, collect_want: "text"},
		{name: "color", middle: "#f00 text", want: `<span style="color:#f00">text</span>`, collect_want: "text"},
	} {
		t.Run(test_case.name, func(t *testing.T) {
			renderer := markup.New_namumark_compat_renderer(nil, "MiddleStyleTest", "", "view", nil, 0, false)
			got := renderer.Restore(renderer.Process_middle_block(test_case.middle))
			if got != test_case.want {
				t.Fatalf("rendered middle = %q, want %q", got, test_case.want)
			}

			collect_renderer := markup.New_namumark_compat_renderer(nil, "MiddleStyleTest", "", "backlink", nil, 0, true)
			if got := collect_renderer.Process_middle_block(test_case.middle); got != test_case.collect_want {
				t.Fatalf("collected middle = %q, want %q", got, test_case.collect_want)
			}
		})
	}
}

func Test_namumark_middle_differential_corpus(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		parameter map[string]any
		want      string
	}{
		{name: "empty input", input: "", want: ""},
		{name: "nested unknown literal", input: "{{{outer {{{inner}}} tail}}}", want: "outer {{{inner}}} tail"},
		{name: "nested style blocks", input: "{{{+2 outer {{{#f00 inner}}} tail}}}", want: `<span style="font-size:140%">outer <span style="color:#f00">inner</span> tail</span>`},
		{name: "conditional", input: "{{{#!if name==alice\\nvisible}}}", parameter: map[string]any{"name": "alice"}, want: "visible"},
		{name: "nested conditional", input: "{{{#!if name==alice\\n{{{#f00 visible}}}\\n}}}", parameter: map[string]any{"name": "alice"}, want: "<span style=\"color:#f00\">visible</span>\\n"},
		{name: "HTML entity", input: "{{{#!html\\n<b>A &amp; B</b>}}}", want: "<b>A &amp; B</b>"},
		{name: "HTML with nested literal", input: "{{{#!html\\n<b>{{{+2 bold}}}</b>\\n}}}", want: "<b>{{{+2 bold}}}</b>"},
		{name: "HTML with multiple nested literals", input: "{{{#!html\\n<b>{{{+2 bold}}}</b> and {{{#f00 red}}}\\n}}}", want: "<b>{{{+2 bold}}}</b> and {{{#f00 red}}}"},
		{name: "HTML nested child with deep braces", input: "{{{#!html\\n<div>{{{#f00 outer {{{+2 inner}}} tail}}}</div>\\n}}}", want: "<div>{{{#f00 outer {{{+2 inner}}} tail}}}</div>"},
		{name: "six brace legacy block", input: "{{{{{{#!html\\n<b>legacy</b>\\n}}}}}}", want: "{{{#!html\\n<b>legacy</b>\\n}}}"},
		{name: "incomplete triple token", input: "before {{{#!if\\nname==alice\\nunfinished", want: "before {{{#!if\\nname==alice\\nunfinished"},
		{name: "incomplete six brace token", input: "before {{{{{{#!html\\nunfinished", want: "before {{{{{{#!html\\nunfinished"},
		{name: "repeated middle tokens", input: "{{{+1 x}}} {{{+1 y}}} {{{#0f0 z}}}", want: `<span style="font-size:120%">x</span> <span style="font-size:120%">y</span> <span style="color:#0f0">z</span>`},
		{name: "deep nesting", input: strings.Repeat("{{{", 8) + "deep" + strings.Repeat("}}}", 8), want: strings.Repeat("{", 21) + "deep" + strings.Repeat("}", 21)},
		{name: "mixed nested syntax", input: "{{{+1 outer {{{#!if name==alice\\ninner}}} tail}}}", parameter: map[string]any{"name": "alice"}, want: `<span style="font-size:120%">outer inner tail</span>`},
	}
	for _, test_case := range cases {
		t.Run(test_case.name, func(t *testing.T) {
			renderer := markup.New_namumark_compat_renderer(nil, "MiddleCorpus", "", "view", test_case.parameter, 0, false)
			input := strings.ReplaceAll(test_case.input, `\n`, "\n")
			want := strings.ReplaceAll(test_case.want, `\n`, "\n")
			got := renderer.Restore(renderer.Process_middle(input))
			if got != want {
				t.Fatalf("Process_middle() = %q, want %q", got, want)
			}
		})
	}

	t.Run("escaped closing delimiter", func(t *testing.T) {
		renderer := markup.New_namumark_compat_renderer(nil, "MiddleCorpus", "", "view", nil, 0, false)
		close_token := renderer.Reserve_slash("}")
		input := strings.ReplaceAll("{{{#!if name==alice\\nvisible"+close_token+"}}", `\n`, "\n")
		got := renderer.Restore(renderer.Process_middle(input))
		if want := "{{{#!if name==alice\nvisible}}}"; got != want {
			t.Fatalf("Process_middle() = %q, want %q", got, want)
		}
	})

	t.Run("escaped nested HTML delimiter", func(t *testing.T) {
		renderer := markup.New_namumark_compat_renderer(nil, "MiddleCorpus", "", "view", nil, 0, false)
		close_token := renderer.Reserve_slash("}")
		input := strings.ReplaceAll("{{{#!html\\n<b>{{{+2 bold"+close_token+"}}</b>\\n}}}", `\n`, "\n")
		got := renderer.Restore(renderer.Process_middle(input))
		if want := "<b>{{{+2 bold}}}</b>"; got != want {
			t.Fatalf("Process_middle() = %q, want %q", got, want)
		}
	})

	t.Run("literal restoration", func(t *testing.T) {
		renderer := markup.New_namumark_compat_renderer(nil, "MiddleCorpus", "", "view", nil, 0, false)
		literal := renderer.Reserve_literal("{{{raw}}}")
		got := renderer.Restore(renderer.Process_middle("before " + literal + " after"))
		if want := "before {{{raw}}} after"; got != want {
			t.Fatalf("restored literal = %q, want %q", got, want)
		}
	})

	t.Run("collect-only", func(t *testing.T) {
		renderer := markup.New_namumark_compat_renderer(nil, "MiddleCorpus", "", "backlink", nil, 0, true)
		got := renderer.Restore(renderer.Process_middle("{{{+2 body}}}"))
		if want := "body"; got != want {
			t.Fatalf("collect-only middle = %q, want %q", got, want)
		}
	})
}

func Test_Namumark_table_parameter_prefix_rendering(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "raw parameters and cell styles",
			input: "||<bgcolor=#123456><color=white>Head||\n||<width=120><keepall>Value||",
			want:  `<div class="table_safe" style=""><table class="" style=""><tr style=""><td colspan="1" rowspan="" style="background:#123456;color:white;text-align: left;">Head</td></tr><tr style=""><td colspan="1" rowspan="" style="width:120px;word-break:keep-all !important;text-align: left;">Value</td></tr></table></div>`,
		},
		{
			name:  "escaped parameters remain escaped",
			input: "||&lt;bgcolor=#123456&gt;&lt;color=white&gt;Head||\n||&lt;width=120&gt;&lt;keepall&gt;Value||",
			want:  `<div class="table_safe" style=""><table class="" style=""><tr style=""><td colspan="1" rowspan="" style="text-align: left;">&amp;lt;bgcolor=#123456&amp;gt;&amp;lt;color=white&amp;gt;Head</td></tr><tr style=""><td colspan="1" rowspan="" style="text-align: left;">&amp;lt;width=120&amp;gt;&amp;lt;keepall&amp;gt;Value</td></tr></table></div>`,
		},
		{
			name:  "unknown and incomplete parameters",
			input: "||<color=red>&lt;width=50&gt;Mixed||\n||<bogus>literal||\n||<span>safe tag||\n||<width=80>Valid||<width=20 Broken||",
			want:  `<div class="table_safe" style=""><table class="" style=""><tr style=""><td colspan="1" rowspan="" style="color:red;text-align: left;">&amp;lt;width=50&amp;gt;Mixed</td></tr><tr style=""><td colspan="1" rowspan="" style="background:bogus;text-align: left;">literal</td></tr><tr style=""><td colspan="1" rowspan="" style="text-align: left;">&lt;span&gt;safe tag</td></tr><tr style=""><td colspan="1" rowspan="" style="width:80px;text-align: left;">Valid</td><td colspan="1" rowspan="" style="text-align: left;">&lt;width=20 Broken</td></tr></table></div>`,
		},
		{
			name:  "colspan rowspan and row styles",
			input: "||<-2><rowbgcolor=#eeeeee><bgcolor=blue>Span||<^|2><color=red>Top||\n||Next||",
			want:  `<div class="table_safe" style=""><table class="" style=""><tr style="background:#eeeeee;"><td colspan="2" rowspan="" style="background:blue;text-align: left;">Span</td><td colspan="1" rowspan="2" style="vertical-align:top;color:red;text-align: left;">Top</td></tr><tr style=""><td colspan="1" rowspan="" style="text-align: left;">Next</td></tr></table></div>`,
		},
	}

	for _, test_case := range cases {
		t.Run(test_case.name, func(t *testing.T) {
			result := markup.Render_namumark_compat(nil, "TableTest", test_case.input, "view", nil, 0, false)
			got, ok := result["data"].(string)
			if !ok || got != test_case.want {
				t.Fatalf("rendered table = %q (type %T), want %q", result["data"], result["data"], test_case.want)
			}
		})
	}
}

func Test_Merge_child_skips_external(t *testing.T) {
	db := Setup_namumark_db(t)
	defer tool.DB_close(db)

	parent := markup.New_namumark_compat_renderer(db, "Parent", "", "view", nil, 0, false)
	child := markup.New_namumark_compat_renderer(db, "Child", "", "include", nil, 0, false)

	child.Add_backlink("Internal", "", "")
	child.Add_backlink("https://child.example.com/e", "external", "")

	parent.Merge_child(child)

	if !parent.Backlink_exists("Internal", "") {
		t.Fatalf("child's internal backlink not merged into Parent")
	}
	if parent.Backlink_exists("https://child.example.com/e", "external") {
		t.Fatalf("child's external backlink was merged into Parent")
	}
	if parent.Has_external_backlink() {
		t.Fatalf("an external backlink was merged into Parent")
	}
}
