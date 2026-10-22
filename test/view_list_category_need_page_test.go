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

// setupNeedCategoryDb builds a fresh DB with exactly n "need" categories.
// Rank r (1-indexed, 1 = first in sort order) gets category `cat-r` with
// distinct link count n-r+1, so the query's `count(distinct link) desc`
// ordering is fully deterministic (cat-1, cat-2, ...).
func setupNeedCategoryDb(t *testing.T, n int) *sql.DB {
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
	// admin user.
	exec("insert into alist (name, acl) values ('owner','owner')")
	exec("insert into alist (name, acl) values ('admin','admin')")
	exec("insert into user_set (name, id, data) values ('acl','admin_ip','owner')")
	// need categories: title 'category:cat-N', type 'cat', no matching data row.
	for r := 1; r <= n; r++ {
		count := n - r + 1
		for l := 1; l <= count; l++ {
			exec(`insert into back (title, link, type, data) values (?, ?, 'cat', '')`,
				fmt.Sprintf("category:cat-%d", r), fmt.Sprintf("c%d_l%d", r, l))
		}
	}
	return db
}

// renderNeedCategory renders View_list_category_need_page for n categories on the page.
func renderNeedCategory(t *testing.T, n int, page int) string {
	t.Helper()
	db := setupNeedCategoryDb(t, n)
	defer tool.DB_close(db)
	config := tool.Config{IP: "admin_ip", Session: &emptySession{}}
	return route.View_list_category_need_page(config, strconv.Itoa(page))
}

// categoryLinkCount counts the category rows shown (each is `<a href="/w/category:...`).
func categoryLinkCount(s string) int {
	return strings.Count(s, `href="/w/category`)
}

// hasPageLink reports whether a page-control link to the given page exists.
func hasPageLink(s string, page int) bool {
	return strings.Contains(s, fmt.Sprintf(`href="/list/category/need/%d"`, page))
}

// Test_view_list_category_need_page verifies the next-page link appears only when a
// next page exists (categories in an exact multiple of 50 no longer show a bogus page),
// while per-page counts, sorting, link-count display and category links are preserved.
func Test_view_list_category_need_page(t *testing.T) {
	cases := []struct {
		name     string
		page     int
		n        int
		wantCnt  int  // categories listed on this page (<=50)
		wantNext bool
		wantPrev bool
	}{
		{"no cats: no next/prev", 1, 0, 0, false, false},
		{"49 cats: no next page", 1, 49, 49, false, false},
		{"50 cats: no next page (exact multiple, previously buggy)", 1, 50, 50, false, false},
		{"51 cats: page 2 exists", 1, 51, 50, true, false},
		{"100 cats page1: 50 shown, page 2 exists", 1, 100, 50, true, false},
		{"100 cats page2: 50 shown, no page 3, prev exists", 2, 100, 50, false, true},
		{"101 cats page1: 50 shown, page 2 exists", 1, 101, 50, true, false},
		{"101 cats page3: 1 shown, prev exists", 3, 101, 1, false, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out := renderNeedCategory(t, c.n, c.page)
			if got := categoryLinkCount(out); got != c.wantCnt {
				t.Fatalf("want %d cats on page %d, got %d", c.wantCnt, c.page, got)
			}
			if hasPageLink(out, c.page+1) != c.wantNext {
				t.Fatalf("want next page %d, got %v (have=%v)", c.page+1, c.wantNext, hasPageLink(out, c.page+1))
			}
			if hasPageLink(out, c.page-1) != c.wantPrev {
				t.Fatalf("want prev page %d, got %v (have=%v)", c.page-1, c.wantPrev, hasPageLink(out, c.page-1))
			}
		})
	}

	// n=51: verify sorting, link-count display, category-link format, and page2.
	p1 := renderNeedCategory(t, 51, 1)
	p2 := renderNeedCategory(t, 51, 2)
	if !strings.Contains(p1, `href="/w/category:cat-1">`) {
		t.Fatalf("page1 missing cat-1 link")
	}
	if !strings.Contains(p1, `href="/w/category:cat-50">`) {
		t.Fatalf("page1 missing cat-50 link")
	}
	if strings.Contains(p1, `href="/w/category:cat-51">`) {
		t.Fatalf("page1 must not include cat-51")
	}
	// cat-1 (count 51) must appear before cat-50 (count 2) -> count-desc sort kept.
	if strings.Index(p1, "cat-1") > strings.Index(p1, "cat-50") {
		t.Fatalf("cat-1 (count 51) should come before cat-50 (count 2)")
	}
	// link-count display for cat-1 (count 51) and cat-50 (count 2) on page1.
	if !strings.Contains(p1, ": 51") {
		t.Fatalf("page1 should show cat-1 link count 51")
	}
	if !strings.Contains(p1, ": 2") {
		t.Fatalf("page1 should show cat-50 link count 2")
	}
	// page2 shows only cat-51 (count 1).
	if !strings.Contains(p2, `href="/w/category:cat-51">`) {
		t.Fatalf("page2 of 51 should show cat-51")
	}
	if !strings.Contains(p2, ": 1") {
		t.Fatalf("page2 of 51 should show cat-51 link count 1")
	}
}
