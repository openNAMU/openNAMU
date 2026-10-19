package test

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

// Setup_query_api_db builds a full schema and seeds documents that reference
// shared external hostnames (plus distinct subdomain and sibling hosts) so the
// V5 host-based query APIs can be verified against real rows.
func Setup_query_api_db(t *testing.T) *sql.DB {
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
	db := tool.DB_connect()

	exec := func(query string, args ...any) {
		if _, err := db.Exec(query, args...); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}

	// Seed default permission groups (mirrors Always_init in the real app).
	// Without these, Get_auth_info returns only "nothing"/"group:user" and
	// the user->ip->document->edit->view->site_view chain never forms, so
	// Check_acl("render") always returns false for an anonymous IP.
	exec("insert into alist (name, acl) values ('owner', 'owner')")
	exec("insert into alist (name, acl) values ('admin', 'admin')")
	exec("insert into alist (name, acl) values ('user', 'user')")
	exec("insert into alist (name, acl) values ('ip', 'ip')")

	// Documents: Alpha, Beta, page, Other, Secret.
	exec("insert into data (title, data) values ('Alpha','content Alpha')")
	exec("insert into data (title, data) values ('Beta','content Beta')")
	exec("insert into data (title, data) values ('page','content page')")
	exec("insert into data (title, data) values ('Other','content Other')")

	// External backlinks (V5: title = normalized hostname, link = source doc).
	// example.com  <- Alpha, Beta, Other
	// www.example.com <- page
	// sub.example.com <- Other
	exec("insert into back (title, link, type, data) values ('example.com','Alpha','external','')")
	exec("insert into back (title, link, type, data) values ('example.com','Beta','external','')")
	exec("insert into back (title, link, type, data) values ('example.com','Other','external','')")
	exec("insert into back (title, link, type, data) values ('www.example.com','page','external','')")
	exec("insert into back (title, link, type, data) values ('sub.example.com','Other','external','')")

	return db
}

// Test_Query_api_docs_for_host verifies route.Api_external_docs_for_host returns only
// the renderable docs that link to the exact normalized host.
func Test_Query_api_docs_for_host(t *testing.T) {
	db := Setup_query_api_db(t)
	defer tool.DB_close(db)

	config := tool.Config{}

	// Exact match: Alpha, Beta, Other link to example.com (not www, not sub).
	got := route.Api_external_docs_for_host(config, "1", "https://example.com/Page")
	if got["response"] != "ok" {
		t.Fatalf("response = %v", got["response"])
	}
	data := got["data"].([][]string)
	if len(data) != 3 {
		t.Fatalf("docs_for_host count = %d, want 3 (Alpha, Beta, Other)", len(data))
	}
	found := map[string]bool{}
	for _, row := range data {
		found[row[0]] = true
		if row[1] != "example.com" {
			t.Errorf("docs_for_host host = %q, want example.com", row[1])
		}
	}
	if !found["Alpha"] || !found["Beta"] || !found["Other"] {
		t.Fatalf("docs_for_host missing docs, got %v", found)
	}
}

// Test_Query_api_docs_for_host_case verifies mixed-case / scheme input is
// normalized to the same host before querying.
func Test_Query_api_docs_for_host_case(t *testing.T) {
	db := Setup_query_api_db(t)
	defer tool.DB_close(db)

	config := tool.Config{}

	got := route.Api_external_docs_for_host(config, "1", "HTTP://EXAMPLE.COM/a")
	if got["response"] != "ok" {
		t.Fatalf("response = %v", got["response"])
	}
	data := got["data"].([][]string)
	if len(data) != 3 {
		t.Fatalf("docs_for_host case count = %d, want 3", len(data))
	}
	for _, row := range data {
		if row[1] != "example.com" {
			t.Errorf("docs_for_host case host = %q, want example.com", row[1])
		}
	}
}

// Test_Query_api_docs_for_host_www verifies a subdomain host is a distinct key.
func Test_Query_api_docs_for_host_www(t *testing.T) {
	db := Setup_query_api_db(t)
	defer tool.DB_close(db)

	config := tool.Config{}

	got := route.Api_external_docs_for_host(config, "1", "https://www.example.com/a")
	if got["response"] != "ok" {
		t.Fatalf("response = %v", got["response"])
	}
	data := got["data"].([][]string)
	if len(data) != 1 {
		t.Fatalf("docs_for_host www count = %d, want 1 (page)", len(data))
	}
	if data[0][0] != "page" || data[0][1] != "www.example.com" {
		t.Fatalf("docs_for_host www row = %v, want [page www.example.com]", data[0])
	}
}

// Test_Query_api_hosts_of_doc verifies route.Api_external_hosts_of_doc returns the
// document's distinct external hosts, sorted.
func Test_Query_api_hosts_of_doc(t *testing.T) {
	db := Setup_query_api_db(t)
	defer tool.DB_close(db)

	config := tool.Config{}

	got := route.Api_external_hosts_of_doc(config, "Other")
	if got["response"] != "ok" {
		t.Fatalf("response = %v", got["response"])
	}
	hosts, ok := got["data"].([]string)
	if !ok {
		t.Fatalf("data not []string")
	}
	if !reflect.DeepEqual(hosts, []string{"example.com", "sub.example.com"}) {
		t.Fatalf("hosts_of_doc = %v, want [example.com sub.example.com]", hosts)
	}

	gotSingle := route.Api_external_hosts_of_doc(config, "Alpha")
	hostsSingle, _ := gotSingle["data"].([]string)
	if !reflect.DeepEqual(hostsSingle, []string{"example.com"}) {
		t.Fatalf("hosts_of_doc Alpha = %v, want [example.com]", hostsSingle)
	}
}

// Test_Query_api_count_for_host verifies route.Api_external_count_for_host returns the
// number of distinct docs linking to the host.
func Test_Query_api_count_for_host(t *testing.T) {
	db := Setup_query_api_db(t)
	defer tool.DB_close(db)

	config := tool.Config{}

	got := route.Api_external_count_for_host(config, "https://example.com/Page")
	if got["response"] != "ok" {
		t.Fatalf("response = %v", got["response"])
	}
	count, ok := got["count"].(int)
	if !ok {
		t.Fatalf("count not int")
	}
	if count != 3 {
		t.Fatalf("count_for_host = %d, want 3", count)
	}

	gotwww := route.Api_external_count_for_host(config, "https://www.example.com/a")
	countwww, okwww := gotwww["count"].(int)
	if !okwww {
		t.Fatalf("count www not int")
	}
	if countwww != 1 {
		t.Fatalf("count_for_host www = %d, want 1", countwww)
	}
}

// Test_Query_api_count_for_host_empty verifies empty/invalid input returns 0.
func Test_Query_api_count_for_host_empty(t *testing.T) {
	db := Setup_query_api_db(t)
	defer tool.DB_close(db)

	config := tool.Config{}

	for _, bad := range []string{"", "not a host", "mailto:x", "//rel", "#anchor", "ftp://example.com", "localhost"} {
		got := route.Api_external_count_for_host(config, bad)
		if got["response"] != "ok" {
			t.Fatalf("bad host %q: response = %v", bad, got["response"])
		}
		count, ok := got["count"].(int)
		if !ok || count != 0 {
			t.Fatalf("bad host %q: count = %v, want 0", bad, count)
		}
	}
}

// Test_Query_api_hosts_of_doc_missing verifies a non-existent doc is handled
// gracefully: the API returns "ok" with an empty host list (Check_acl does not
// verify doc existence; it only enforces view permission, which a regular
// user has via the seeded permission chain), and never panics.
func Test_Query_api_hosts_of_doc_missing(t *testing.T) {
	db := Setup_query_api_db(t)
	defer tool.DB_close(db)

	config := tool.Config{}

	got := route.Api_external_hosts_of_doc(config, "NonExistentDoc")
	if got["response"] != "ok" {
		t.Fatalf("hosts_of_doc for missing doc: response = %v, want ok", got["response"])
	}
	if got["data"] == nil {
		t.Fatal("hosts_of_doc for missing doc: data should be empty list, got nil")
	}
}

// Test_Query_api_require_auth verifies that a doc with a restricted view ACL
// returns "require auth" for a user that lacks permission.
func Test_Query_api_require_auth(t *testing.T) {
	db := Setup_query_api_db(t)
	defer tool.DB_close(db)

	config := tool.Config{}

	exec := func(query string, args ...any) {
		if _, err := db.Exec(query, args...); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
	exec("insert into data (title, data) values ('Secret','hidden content')")
	exec("insert into acl (title, data, type) values ('Secret','owner','view')")
	exec("insert into back (title, link, type, data) values ('example.com','Secret','external','')")

	// A regular user does not have owner permission -> require auth.
	got := route.Api_external_hosts_of_doc(config, "Secret")
	if got["response"] != "require auth" {
		t.Fatalf("hosts_of_doc for owner-restricted doc: response = %v, want require auth", got)
	}
}

// Test_Query_api_hosts_of_doc_case_sibling verifies that case-sibling docs
// (Page and page) keep their external hosts separate when queried by exact
// name, and that a public doc's query never leaks a non-public sibling's host.
// On SQLite the SQL is case-sensitive; this test guards against a regression
// that would merge case-siblings via a case-insensitive collation.
func Test_Query_api_hosts_of_doc_case_sibling(t *testing.T) {
	db := Setup_query_api_db(t)
	defer tool.DB_close(db)

	exec := func(query string, args ...any) {
		if _, err := db.Exec(query, args...); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
	// Seed a true case-sibling pair: Page (uppercase, non-public) and page
	// (lowercase, public) on distinct hosts. The setup already has `page`, so
	// clear its backlink and re-seed it with the distinct host; add `Page`.
	exec("insert into data (title, data) values ('Page','content Page')")
	exec("delete from back where link = 'page'")
	exec("insert into back (title, link, type, data) values ('example.com','Page','external','')")
	exec("insert into back (title, link, type, data) values ('other.example.com','page','external','')")
	// The uppercase sibling is non-public (owner view); the lowercase is public.
	exec("insert into acl (title, data, type) values ('Page','owner','view')")
	// An admin (owner) can read both docs for the mutual case-isolation check.
	exec("insert into user_set (name, id, data) values ('acl','admin_ip','owner')")

	admin := tool.Config{IP: "admin_ip"}
	anon := tool.Config{}

	// Mutual case isolation: each doc returns only its own distinct host.
	pageOut := route.Api_external_hosts_of_doc(admin, "Page")
	pageHosts, _ := pageOut["data"].([]string)
	if !reflect.DeepEqual(pageHosts, []string{"example.com"}) {
		t.Fatalf("Page hosts = %v, want [example.com]", pageHosts)
	}
	lcPageOut := route.Api_external_hosts_of_doc(admin, "page")
	lcPageHosts, _ := lcPageOut["data"].([]string)
	if !reflect.DeepEqual(lcPageHosts, []string{"other.example.com"}) {
		t.Fatalf("page hosts = %v, want [other.example.com]", lcPageHosts)
	}

	// Security regression: the PUBLIC doc's query must not leak the NON-PUBLIC
	// case-sibling's host.
	lcPublicOut := route.Api_external_hosts_of_doc(anon, "page")
	hosts, ok := lcPublicOut["data"].([]string)
	if !ok {
		t.Fatalf("page (public) hosts: data not []string")
	}
	if !reflect.DeepEqual(hosts, []string{"other.example.com"}) {
		t.Fatalf("page (public) hosts = %v, want [other.example.com]", hosts)
	}
	for _, h := range hosts {
		if h == "example.com" {
			t.Fatalf("public doc page leaked non-public case-sibling Page host: %v", hosts)
		}
	}
}

// Test_Query_api_docs_for_host_excludes_deleted_doc verifies a deleted document's
// stale external backlink row is not reported as a reference doc, and that the
// count agrees with the (now empty) list.
func Test_Query_api_docs_for_host_excludes_deleted_doc(t *testing.T) {
	db := Setup_query_api_db(t)
	defer tool.DB_close(db)

	config := tool.Config{}
	exec := func(query string, args ...any) {
		if _, err := db.Exec(query, args...); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
	exec("insert into data (title, data) values ('Orphan','content Orphan')")
	exec("insert into back (title, link, type, data) values ('orphan.example.com','Orphan','external','')")

	got := route.Api_external_docs_for_host(config, "1", "orphan.example.com")
	if rows := got["data"].([][]string); len(rows) != 1 || rows[0][0] != "Orphan" {
		t.Fatalf("docs_for_host for Orphan = %v, want 1 row [Orphan]", got["data"])
	}

	exec("delete from data where title = 'Orphan'")
	got2 := route.Api_external_docs_for_host(config, "1", "orphan.example.com")
	if rows := got2["data"].([][]string); len(rows) != 0 {
		t.Fatalf("deleted doc Orphan still reported after delete: %v", got2["data"])
	}

	count, _ := route.Api_external_count_for_host(config, "orphan.example.com")["count"].(int)
	if count != 0 {
		t.Fatalf("count_for_host for deleted Orphan = %d, want 0", count)
	}
}

// Test_Query_api_count_matches_list verifies the reference-doc count equals the
// number of distinct docs returned on page 1 (no paging truncation at this size).
func Test_Query_api_count_matches_list(t *testing.T) {
	db := Setup_query_api_db(t)
	defer tool.DB_close(db)

	config := tool.Config{}
	got := route.Api_external_docs_for_host(config, "1", "https://example.com/Page")
	data := got["data"].([][]string)
	count, _ := route.Api_external_count_for_host(config, "https://example.com/Page")["count"].(int)
	if count != len(data) {
		t.Fatalf("count=%d but list length=%d (want equal)", count, len(data))
	}
}

// Test_Query_api_all_hosts verifies the admin all-hosts list returns per-host
// distinct-document counts in stable alphabetical order, and requires the
// setting_backlink permission.
func Test_Query_api_all_hosts(t *testing.T) {
	db := Setup_query_api_db(t)
	defer tool.DB_close(db)

	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed %v: %v", q, err)
		}
	}
	// Bind admin_ip to the owner group (which grants setting_backlink).
	exec("insert into user_set (name, id, data) values ('acl','admin_ip','owner')")
	admin := tool.Config{IP: "admin_ip"}

	// Non-admin (anonymous) is denied.
	got := route.Api_external_all_hosts(tool.Config{}, "1")
	if got["response"] != "require auth" {
		t.Fatalf("non-admin all_hosts response = %v, want require auth", got["response"])
	}

	// Admin gets per-host distinct doc counts, alphabetical order, page 1.
	got = route.Api_external_all_hosts(admin, "1")
	if got["response"] != "ok" {
		t.Fatalf("admin all_hosts response = %v, want ok", got["response"])
	}
	data, ok := got["data"].([][]string)
	if !ok {
		t.Fatalf("admin all_hosts data not [][]string: %T", got["data"])
	}
	want := [][]string{
		{"example.com", "3"},
		{"sub.example.com", "1"},
		{"www.example.com", "1"},
	}
	if !reflect.DeepEqual(data, want) {
		t.Fatalf("all_hosts page 1 = %v, want %v", data, want)
	}
}

// Test_Query_api_all_hosts_paging verifies the all-hosts list is paginated at
// 50 per page (page 1 is full, later pages carry the remainder).
func Test_Query_api_all_hosts_paging(t *testing.T) {
	db := Setup_query_api_db(t)
	defer tool.DB_close(db)

	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed %v: %v", q, err)
		}
	}
	exec("insert into user_set (name, id, data) values ('acl','admin_ip','owner')")
	admin := tool.Config{IP: "admin_ip"}

	// Seed extra distinct hosts so the total exceeds one page (50).
	const extra = 51
	for i := 0; i < extra; i++ {
		host := fmt.Sprintf("host%03d", i)
		exec("insert into back (title, link, type, data) values (?, 'doc', 'external', '')", host)
	}
	total := 3 + extra

	page1 := route.Api_external_all_hosts(admin, "1")
	data1, _ := page1["data"].([][]string)
	if len(data1) != 50 {
		t.Fatalf("all_hosts page 1 count = %d, want 50 (total=%d)", len(data1), total)
	}
	page2 := route.Api_external_all_hosts(admin, "2")
	data2, _ := page2["data"].([][]string)
	want2 := total - 50
	if len(data2) != want2 {
		t.Fatalf("all_hosts page 2 count = %d, want %d", len(data2), want2)
	}
	// A page far beyond the data is empty.
	pageLast := route.Api_external_all_hosts(admin, "10")
	dataLast, _ := pageLast["data"].([][]string)
	if len(dataLast) != 0 {
		t.Fatalf("all_hosts page 10 not empty (count=%d), want 0", len(dataLast))
	}
}

// Test_Query_api_all_hosts_count_excludes_deleted_and_case_distinct verifies that
// the admin all-hosts count counts only docs that actually exist in `data`,
// case-sensitively (PageX and pagex are distinct), and agrees with the per-host
// detailed view.
func Test_Query_api_all_hosts_count_excludes_deleted_and_case_distinct(t *testing.T) {
	db := Setup_query_api_db(t)
	defer tool.DB_close(db)

	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed exec %v: %v", q, err)
		}
	}
	// tld.example.com: two existing docs (DocA, DocB) plus a stale backlink to a
	// deleted doc (Ghost, absent from data) that must be excluded from the count.
	exec("insert into data (title, data) values ('DocA','c')")
	exec("insert into data (title, data) values ('DocB','c')")
	exec("insert into back (title, link, type, data) values ('tld.example.com','DocA','external','')")
	exec("insert into back (title, link, type, data) values ('tld.example.com','DocB','external','')")
	exec("insert into back (title, link, type, data) values ('tld.example.com','Ghost','external','')")
	// case.example.com: two case-sibling docs (PageX, pagex), both existing; the
	// count must treat them as distinct (not merged by a case-insensitive collation).
	exec("insert into data (title, data) values ('PageX','c')")
	exec("insert into data (title, data) values ('pagex','c')")
	exec("insert into back (title, link, type, data) values ('case.example.com','PageX','external','')")
	exec("insert into back (title, link, type, data) values ('case.example.com','pagex','external','')")

	exec("insert into user_set (name, id, data) values ('acl','admin_ip','owner')")
	admin := tool.Config{IP: "admin_ip"}
	got := route.Api_external_all_hosts(admin, "1")
	if got["response"] != "ok" {
		t.Fatalf("all_hosts response = %v, want ok", got["response"])
	}
	data, ok := got["data"].([][]string)
	if !ok {
		t.Fatalf("all_hosts data not [][]string: %T", got["data"])
	}
	countFor := func(host string) string {
		for _, row := range data {
			if row[0] == host {
				return row[1]
			}
		}
		return ""
	}
	if n := countFor("tld.example.com"); n != "2" {
		t.Errorf("tld.example.com count=%s, want 2 (Ghost deleted, excluded; DocA+DocB)", n)
	}
	if n := countFor("case.example.com"); n != "2" {
		t.Errorf("case.example.com count=%s, want 2 (PageX and pagex distinct)", n)
	}

	// The per-host count must agree with the all-hosts list (admin renders all docs).
	if c, okc := route.Api_external_count_for_host(admin, "tld.example.com")["count"].(int); !okc || c != 2 {
		t.Errorf("count_for_host tld.example.com = %v, want 2", route.Api_external_count_for_host(admin, "tld.example.com")["count"])
	}
}

// Test_Query_api_batch_in_large verifies that when a single host is referenced by
// more docs than one IN batch can hold (external_doc_in_batch), the existing-doc
// IN query is batched and the count/list still reflect every existing doc.
func Test_Query_api_batch_in_large(t *testing.T) {
	db := Setup_query_api_db(t)
	defer tool.DB_close(db)

	const total = 600 // greater than external_doc_in_batch (500) to force multiple batches
	exec := func(query string, args ...any) {
		if _, err := db.Exec(query, args...); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
	for i := 0; i < total; i++ {
		doc_name := fmt.Sprintf("BatchDoc%04d", i)
		exec("insert into data (title, data) values (?, ?)", doc_name, "c")
		exec("insert into back (title, link, type, data) values ('large.example.com', ?, 'external', '')", doc_name)
	}

	// An admin (owner) can render all docs for the per-doc view check.
	exec("insert into user_set (name, id, data) values ('acl','admin_ip','owner')")
	admin := tool.Config{IP: "admin_ip"}

	count, _ := route.Api_external_count_for_host(admin, "large.example.com")["count"].(int)
	if count != total {
		t.Fatalf("count_for_host large.example.com = %d, want %d (IN query must handle >1 batch)", count, total)
	}

	all := route.Api_external_all_hosts(admin, "1")
	if all["response"] != "ok" {
		t.Fatalf("all_hosts response = %v, want ok", all["response"])
	}
	allData, _ := all["data"].([][]string)
	found := false
	for _, row := range allData {
		if row[0] != "large.example.com" {
			continue
		}
		found = true
		if row[1] != fmt.Sprintf("%d", total) {
			t.Fatalf("all_hosts large.example.com count=%s, want %d", row[1], total)
		}
	}
	if !found {
		t.Fatalf("all_hosts page 1 missing large.example.com; only %d rows returned", len(allData))
	}

	docRows, _ := route.Api_external_docs_for_host(admin, "1", "large.example.com")["data"].([][]string)
	if len(docRows) != 50 {
		t.Fatalf("docs_for_host large.example.com page 1 = %d rows, want 50", len(docRows))
	}
}

// Test_View_external_menu_single_slash verifies the rendered page menu href is a
// single-slash path (/ext_backlink/<host>, /w/<doc>) and never a double slash
// (//ext_backlink, //w/). The ringo template prepends a leading "/" to each menu
// path, so View menu paths must be supplied without one.
func Test_View_external_menu_single_slash(t *testing.T) {
	db := Setup_query_api_db(t)
	defer tool.DB_close(db)

	config := tool.Config{IP: "anonymous", Session: &emptySession{}}

	out := route.View_external_docs_for_host(config, "example.com", "1")
	if strings.Contains(out, "//ext_backlink") {
		t.Fatalf("rendered menu/content contains double-slash ext_backlink path")
	}
	if !strings.Contains(out, "/ext_backlink/example.com") {
		t.Fatalf("rendered menu/content missing single-slash ext_backlink path")
	}

	out2 := route.View_external_hosts_of_doc(config, "Alpha")
	if strings.Contains(out2, "//ext_backlink") || strings.Contains(out2, "//w/Alpha") {
		t.Fatalf("rendered menu/content contains double-slash path")
	}
	if !strings.Contains(out2, "/w/Alpha") {
		t.Fatalf("rendered menu/content missing single-slash w/ path")
	}
}
