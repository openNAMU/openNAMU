package test

import (
	"strings"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

// Test_view_main_search_post_search_type verifies the search-type redirect
// logic reused by the content-search POST routes: search_type="data" (passed by
// /search_data and /search_data_page POST requests) redirects to a data-search
// page, while an empty search_type (title search / default) redirects to
// title search. This is the regression guard for the content-search POST
// routes staying on content search.
func Test_view_main_search_post_search_type(t *testing.T) {
	db := Setup_rss_db(t, false)
	defer tool.DB_close(db)

	config := tool.Config{IP: "192.168.0.1", Session: &emptySession{}}

	keyword := "검색어"
	encoded := tool.Url_parser(keyword)

	// content search: the search_type the /search_data* POST routes now pass.
	out := route.View_main_search_post(config, "data", false, keyword)
	if !strings.Contains(out, "/search_data_page/1/"+encoded) {
		t.Fatalf("data search should redirect to a data page, got %q", out)
	}
	if strings.Contains(out, "/search/1/"+encoded) {
		t.Fatalf("data search must not fall back to title search; got %q", out)
	}

	// title search: empty search_type, the existing title-search behavior.
	out = route.View_main_search_post(config, "", false, keyword)
	if !strings.Contains(out, "/search/"+encoded) {
		t.Fatalf("title search should redirect to /search/, got %q", out)
	}
}
