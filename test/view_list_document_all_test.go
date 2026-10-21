package test

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

// setupDocListDb builds a fresh DB containing only n document rows (doc001..)
// plus an admin user, so pagination counts are not affected by seeded wiki docs.
func setupDocListDb(t *testing.T, n int) *sql.DB {
	t.Helper()
	db_dir := filepath.Join(t.TempDir(), "db")
	os.MkdirAll(db_dir, 0755)
	os.Setenv("NAMU_DB", db_dir)
	os.Setenv("NAMU_DB_TYPE", "sqlite")
	defer os.Unsetenv("NAMU_DB")
	defer os.Unsetenv("NAMU_DB_TYPE")

	tool.DB_init()
	db := tool.DB_connect()

	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed %v: %v", q, err)
		}
	}
	// admin user bound to the owner group.
	exec("insert into alist (name, acl) values ('owner','owner')")
	exec("insert into alist (name, acl) values ('admin','admin')")
	exec("insert into user_set (name, id, data) values ('acl','admin_ip','owner')")

	// zero-padded so doc001..docNNN sort in numeric order.
	for i := 1; i <= n; i++ {
		exec(`insert into data (title, data) values (?, '')`, fmt.Sprintf("doc%03d", i))
	}
	return db
}

// renderDocList renders View_list_document_all for n docs on the given page.
func renderDocList(t *testing.T, n int, page int) string {
	t.Helper()
	db := setupDocListDb(t, n)
	defer tool.DB_close(db)
	config := tool.Config{IP: "admin_ip", Session: &emptySession{}}
	return route.View_list_document_all(config, strconv.Itoa(page))
}

// docLinkCount counts document rows (each is `<li>N. <a href="/w/docNNN">`).
func docLinkCount(s string) int {
	return strings.Count(s, `href="/w/doc`)
}

// Test_view_list_document_all_pagination verifies the next-page link appears only
// when a next page exists (50-row pages no longer show a bogus page 2), and that
// per-page counts, numbering and ordering stay correct after fetching 51 rows.
func Test_view_list_document_all_pagination(t *testing.T) {
	cases := []struct {
		name     string
		page     int
		n        int
		wantDocs int // docs listed on this page (<=50)
		wantNext bool
		wantPrev bool
	}{
		{"no docs: no next/prev", 1, 0, 0, false, false},
		{"49 docs: no next page", 1, 49, 49, false, false},
		{"50 docs: no next page (exact multiple, previously buggy)", 1, 50, 50, false, false},
		{"51 docs: page 2 exists", 1, 51, 50, true, false},
		{"100 docs page1: 50 shown, page 2 exists", 1, 100, 50, true, false},
		{"100 docs page2: 50 shown, no page 3, prev exists", 2, 100, 50, false, true},
		{"101 docs page1: 50 shown, page 2 exists", 1, 101, 50, true, false},
		{"101 docs page3: 1 shown, prev exists", 3, 101, 1, false, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out := renderDocList(t, c.n, c.page)
			if got := docLinkCount(out); got != c.wantDocs {
				t.Fatalf("want %d docs on page %d, got %d", c.wantDocs, c.page, got)
			}
			if next := strings.Contains(out, fmt.Sprintf(`href="/list/document/all/%d"`, c.page+1)); next != c.wantNext {
				t.Fatalf("want next page %v, got %v", c.wantNext, next)
			}
			if prev := strings.Contains(out, fmt.Sprintf(`href="/list/document/all/%d"`, c.page-1)); prev != c.wantPrev {
				t.Fatalf("want prev page %v, got %v", c.wantPrev, prev)
			}
		})
	}

	// ordering + numbering: page 2 of 100 docs must show 51..100 in order.
	out := renderDocList(t, 100, 2)
	for i := 51; i <= 100; i++ {
		if !strings.Contains(out, fmt.Sprintf(`<li>%d. <a href="/w/doc%03d"`, i, i)) {
			t.Fatalf("page 2 missing doc %d in expected number", i)
		}
	}
	// page 3 of 101 shows only doc 101.
	out3 := renderDocList(t, 101, 3)
	if !strings.Contains(out3, `<li>101. <a href="/w/doc101">`) {
		t.Fatalf("page 3 of 101 should show doc 101")
	}
	if strings.Contains(out3, `<li>102. <a href="/w/doc102">`) {
		t.Fatalf("page 3 must not include doc 102")
	}
	// 50-doc page 1 ends at 50 (no 51st), 0-doc page shows the empty notice.
	out50 := renderDocList(t, 50, 1)
	if strings.Contains(out50, `href="/w/doc051">`) {
		t.Fatalf("50-doc page 1 must not include doc 51")
	}
	out0 := renderDocList(t, 0, 1)
	if !strings.Contains(out0, "정보 없음") {
		t.Fatalf("0-doc page should show the data-missing notice")
	}
}
