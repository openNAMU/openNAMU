package test

import (
	"database/sql"
	"strings"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

// seed_diff_doc seeds a document with three history revisions:
//   - rev 1 and rev 2 have identical content (no changes)
//   - rev 3 differs from rev 1 (has changes)
func seed_diff_doc(db *sql.DB, t *testing.T) {
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed %v: %v", q, err)
		}
	}
	exec(`insert into data (title, data) values ('DiffDoc', 'Base content')`)
	exec(`insert into history (id, title, data, date, ip, send, leng, hide, type)
		values ('1', 'DiffDoc', 'Base content', '2026-01-01 00:00:00', 'admin_ip', 'edit', '+10', '', '')`)
	exec(`insert into history (id, title, data, date, ip, send, leng, hide, type)
		values ('2', 'DiffDoc', 'Base content', '2026-01-02 00:00:00', 'admin_ip', 'edit', '+10', '', '')`)
	exec(`insert into history (id, title, data, date, ip, send, leng, hide, type)
		values ('3', 'DiffDoc', 'Base content
Modified', '2026-01-03 00:00:00', 'admin_ip', 'edit', '+10', '', '')`)
}

// Test_view_w_diff_no_changes verifies that identical revisions (and the same
// revision compared to itself) show the "no changes" message instead of an
// empty page, while revisions with differences keep the existing diff table.
func Test_view_w_diff_no_changes(t *testing.T) {
	db := Setup_rebuild_db(t)
	defer tool.DB_close(db)

	seed_diff_doc(db, t)

	config := tool.Config{IP: "admin_ip", Session: &emptySession{}}

	// identical revisions (1 vs 2) -> no changes message, no diff table
	out := route.View_w_diff(config, "DiffDoc", "1", "2")
	if !strings.Contains(out, "변경 사항이 없습니다.") {
		t.Fatalf("identical revisions should show no-changes message, got: %q", out[:300])
	}
	if strings.Contains(out, "<table") {
		t.Fatalf("identical revisions should not show diff table, got: %q", out[:300])
	}

	// same revision (1 vs 1) -> no changes message, no diff table
	out = route.View_w_diff(config, "DiffDoc", "1", "1")
	if !strings.Contains(out, "변경 사항이 없습니다.") {
		t.Fatalf("same revision should show no-changes message, got: %q", out[:300])
	}
	if strings.Contains(out, "<table") {
		t.Fatalf("same revision should not show diff table, got: %q", out[:300])
	}

	// different revisions (1 vs 3) -> existing diff table, no no-changes message
	out = route.View_w_diff(config, "DiffDoc", "1", "3")
	if strings.Contains(out, "변경 사항이 없습니다.") {
		t.Fatalf("different revisions should not show no-changes message, got: %q", out[:300])
	}
	if !strings.Contains(out, "<table") {
		t.Fatalf("different revisions should show diff table, got: %q", out[:300])
	}
	if !strings.Contains(out, "r1") {
		t.Fatalf("diff table should reference revision 1, got: %q", out[:300])
	}
}

// Test_view_w_diff_error_handling preserves the existing behavior for missing
// revisions and missing permissions.
func Test_view_w_diff_error_handling(t *testing.T) {
	db := Setup_rebuild_db(t)
	defer tool.DB_close(db)

	seed_diff_doc(db, t)

	// nonexistent revision -> redirect to /history/DiffDoc
	out := route.View_w_diff(tool.Config{IP: "admin_ip", Session: &emptySession{}}, "DiffDoc", "1", "99")
	if !strings.Contains(out, "/history/DiffDoc") {
		t.Fatalf("nonexistent revision should redirect to history page, got: %q", out[:300])
	}

	// missing permission -> auth error page, not the diff
	out = route.View_w_diff(tool.Config{IP: "stranger", Session: &emptySession{}}, "DiffDoc", "1", "2")
	if strings.Contains(out, "변경 사항이 없습니다.") {
		t.Fatalf("missing permission should not render the no-changes message, got: %q", out[:300])
	}
}
