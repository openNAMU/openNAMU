package test

import (
	"strings"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

// Test_view_main_search_search_type_highlight verifies the currently selected
// search-type link label is wrapped in <strong>: title search highlights the
// title-search link, data search highlights the data-search link. Both links
// remain (clickable) in every case. The default /search/*keyword path passes
// search_type="title".
func Test_view_main_search_search_type_highlight(t *testing.T) {
	db := Setup_rss_db(t, false)
	defer tool.DB_close(db)

	config := tool.Config{IP: "192.168.0.1", Session: &emptySession{}}

	keyword := "검색어"
	encoded := tool.Url_parser(keyword)
	title_href := "/search_page/1/" + encoded
	data_href := "/search_data_page/1/" + encoded
	// ko-KR is seeded by Setup_rss_db, so these are the ko-KR translations.
	title_label := "(문서명 검색)"
	data_label := "(문서 내용 검색)"

	title_strong := `<a href="` + title_href + `"><strong>` + title_label + `</strong></a>`
	title_plain := `<a href="` + title_href + `">` + title_label + `</a>`
	data_strong := `<a href="` + data_href + `"><strong>` + data_label + `</strong></a>`
	data_plain := `<a href="` + data_href + `">` + data_label + `</a>`

	// title search (default /search/*keyword): title link emphasized.
	out := route.View_main_search(config, keyword, "1", "title")
	if !strings.Contains(out, title_strong) {
		t.Fatalf("title search: title link not emphasized")
	}
	if !strings.Contains(out, data_plain) {
		t.Fatalf("title search: data link should stay plain (not emphasized)")
	}
	// both links remain present (clickable)
	if !strings.Contains(out, `href="`+title_href+`"`) ||
		!strings.Contains(out, `href="`+data_href+`"`) {
		t.Fatalf("title search: a search-type link href is missing")
	}

	// data search (/search_data*): data link emphasized.
	out = route.View_main_search(config, keyword, "1", "data")
	if !strings.Contains(out, data_strong) {
		t.Fatalf("data search: data link not emphasized")
	}
	if !strings.Contains(out, title_plain) {
		t.Fatalf("data search: title link should stay plain (not emphasized)")
	}
	if !strings.Contains(out, `href="`+title_href+`"`) ||
		!strings.Contains(out, `href="`+data_href+`"`) {
		t.Fatalf("data search: a search-type link href is missing")
	}
}
