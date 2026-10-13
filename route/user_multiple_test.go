package route

import (
	"database/sql"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"opennamu/route/tool"
)

// Setup_multiple_db builds the full schema and seeds IP groups that mix
// hideip and non-hideip members, for multi-account regression tests.
func Setup_multiple_db(t *testing.T) *sql.DB {
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

	exec(`delete from alist where name in ('inspector_group','victim_group')`)
	exec(`delete from user_set where name = 'acl' and id in ('inspector','victim')`)
	exec(`delete from ua_d where name in ('victim','open','dave')`)

	// inspector has check; victim holds hideip; open and dave hold none.
	exec(`insert into alist (name, acl) values ('inspector_group','check')`)
	exec(`insert into alist (name, acl) values ('victim_group','hideip')`)
	exec(`insert into user_set (id, name, data) values ('inspector','acl','inspector_group')`)
	exec(`insert into user_set (id, name, data) values ('victim','acl','victim_group')`)

	// 192.168.0.10: hideip (victim) + non-hideip (open) -> excluded group.
	// 192.168.0.30: two non-hideip users (open, dave) -> shown group.
	exec(`insert into ua_d (name, ip, ua, today, sub) values ('victim','192.168.0.10','Mozilla/5.0',?,'0')`, "20261008")
	exec(`insert into ua_d (name, ip, ua, today, sub) values ('open','192.168.0.10','Mozilla/5.1',?,'0')`, "20261008")
	exec(`insert into ua_d (name, ip, ua, today, sub) values ('open','192.168.0.30','Mozilla/5.2',?,'0')`, "20261008")
	exec(`insert into ua_d (name, ip, ua, today, sub) values ('dave','192.168.0.30','Mozilla/5.3',?,'0')`, "20261008")

	return db
}

// Test_Get_ua_multiple_rows_excludes_hideip_groups verifies that an IP group
// containing any hideip member is dropped, while a regular-only group is kept.
func Test_Get_ua_multiple_rows_excludes_hideip_groups(t *testing.T) {
	db := Setup_multiple_db(t)
	defer tool.DB_close(db)

	rows := tool.Get_ua_multiple_rows(db, "", "recent", 0)

	ips := map[string]bool{}
	for _, r := range rows {
		ips[r.IP] = true
	}
	if ips["192.168.0.10"] {
		t.Fatalf("hideip co-tenancy group 192.168.0.10 was exposed; rows=%v", rows)
	}
	if !ips["192.168.0.30"] {
		t.Fatalf("regular group 192.168.0.30 was missing; rows=%v", rows)
	}
	for _, r := range rows {
		for _, name := range r.Names {
			if name == "victim" {
				t.Fatalf("hideip member 'victim' leaked into group %s", r.IP)
			}
		}
	}
}

// Test_Get_ua_multiple_rows_search_by_hideip_name verifies that searching by a
// hideip member's name still does not surface the group.
func Test_Get_ua_multiple_rows_search_by_hideip_name(t *testing.T) {
	db := Setup_multiple_db(t)
	defer tool.DB_close(db)

	rows := tool.Get_ua_multiple_rows(db, "victim", "recent", 0)
	for _, r := range rows {
		if r.IP == "192.168.0.10" {
			t.Fatalf("search by hideip name 'victim' exposed group 192.168.0.10")
		}
	}
}

// Test_View_user_multiple_excludes_hideip_groups verifies the rendered multi
// account HTML omits hideip co-tenancy groups.
func Test_View_user_multiple_excludes_hideip_groups(t *testing.T) {
	db := Setup_multiple_db(t)
	defer tool.DB_close(db)

	config := tool.Config{IP: "inspector", Session: &emptySession{}}
	out := View_user_multiple(config, "1", "recent", "")

	if strings.Contains(out, "victim") || strings.Contains(out, "192.168.0.10") {
		t.Fatalf("hideip co-tenancy leaked in multi-account view")
	}
	if !strings.Contains(out, "dave") {
		t.Fatalf("regular group member 'dave' missing from multi-account view")
	}
}

// Setup_large_multiple_db seeds 1000 IP groups: even index is a visible group
// (2 normal users), odd index is a hideip group (1 normal + 1 hideip user).
// Ordering is deterministic (newest first) so visible groups are 0,2,4,...,998
// -> 500 visible groups -> 10 pages of 50.
func Setup_large_multiple_db(t *testing.T) *sql.DB {
	t.Helper()
	db := Setup_multiple_db(t)
	exec := func(q string, a ...any) {
		if _, err := db.Exec(q, a...); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
	// Drop the base seed's groups so only the 1000 groups seeded below exist,
	// keeping the visible count exact (500 of them).
	exec(`delete from ua_d`)
	exec(`insert into alist (name, acl) values ('big_hideip_group','hideip')`)
	for i := 0; i < 1000; i++ {
		ip := fmt.Sprintf("9.9.%d.%d", i/10, i%10)
		day := fmt.Sprintf("%03d", 999-i)
		u := fmt.Sprintf("big_u_%d", i)
		v := fmt.Sprintf("big_v_%d", i)
		exec(`insert into ua_d (name, ip, ua, today, sub) values (?, ?, 'Mozilla/5.0', ?, '0')`, u, ip, day)
		exec(`insert into ua_d (name, ip, ua, today, sub) values (?, ?, 'Mozilla/5.0', ?, '0')`, v, ip, day)
		if i%2 == 1 {
			exec(`insert into user_set (name, id, data) values ('acl',?, 'big_hideip_group')`, v)
		}
	}
	return db
}

func Test_Get_ua_multiple_rows_large_scale_pagination(t *testing.T) {
	db := Setup_large_multiple_db(t)
	defer tool.DB_close(db)

	page := func(offset int) []tool.UA_multiple_data {
		return tool.Get_ua_multiple_rows(db, "", "recent", offset)
	}
	page1 := page(0)
	page2 := page(50)
	page10 := page(450)
	page11 := page(500)

	if len(page1) != 50 {
		t.Fatalf("page1 len=%d want 50", len(page1))
	}
	if len(page2) != 50 {
		t.Fatalf("page2 len=%d want 50", len(page2))
	}
	if len(page10) != 50 {
		t.Fatalf("page10 len=%d want 50", len(page10))
	}
	if len(page11) != 0 {
		t.Fatalf("page11 len=%d want 0", len(page11))
	}

	// Deterministic boundaries (visible groups are even index: 0,2,...,998).
	if page1[0].IP != "9.9.0.0" {
		t.Fatalf("page1[0].IP=%q want 9.9.0.0", page1[0].IP)
	}
	if page1[49].IP != "9.9.9.8" {
		t.Fatalf("page1[49].IP=%q want 9.9.9.8", page1[49].IP)
	}
	if page10[0].IP != "9.9.90.0" {
		t.Fatalf("page10[0].IP=%q want 9.9.90.0", page10[0].IP)
	}

	// No hideip member may be leaked in any returned group.
	check := func(lists ...[]tool.UA_multiple_data) {
		for _, list := range lists {
			for _, g := range list {
				for _, n := range g.Names {
					if tool.Check_permission(db, "hideip", n) {
						t.Fatalf("hideip member %q leaked in group %s", n, g.IP)
					}
				}
			}
		}
	}
	check(page1, page2, page10, page11)
}

// Test_Get_ua_multiple_rows_performance times the bounded batch walk over a
// 1000-group dataset. The full result is never materialised; memory stays within
// one 256-row batch plus 50 collected groups.
func Test_Get_ua_multiple_rows_performance(t *testing.T) {
	db := Setup_large_multiple_db(t)
	defer tool.DB_close(db)

	start := time.Now()
	rows := tool.Get_ua_multiple_rows(db, "", "recent", 0)
	elapsed := time.Since(start)
	if len(rows) != 50 {
		t.Fatalf("page1 len=%d want 50", len(rows))
	}
	if elapsed >= 5*time.Second {
		t.Fatalf("first page too slow: %v (want <5s)", elapsed)
	}

	// The last page walks the whole visible range from the start; time it too.
	start = time.Now()
	rows = tool.Get_ua_multiple_rows(db, "", "recent", 450)
	elapsed = time.Since(start)
	if len(rows) != 50 {
		t.Fatalf("page10 len=%d want 50", len(rows))
	}
	if elapsed >= 5*time.Second {
		t.Fatalf("last page too slow: %v (want <5s)", elapsed)
	}
}

// Setup_null_today_multiple_db seeds the NULL-today scenario for Section 1:
//   - 192.0.2.10: alice (hideip, today=NULL) + bob (normal)  -> excluded.
//   - 198.51.100.20: carol + dave (normal)                    -> shown.
//   - 203.0.113.9: erin + frank (normal, both today=NULL)     -> shown, Date empty.
func Setup_null_today_multiple_db(t *testing.T) *sql.DB {
	t.Helper()
	db := Setup_multiple_db(t)
	exec := func(q string, a ...any) {
		if _, err := db.Exec(q, a...); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
	exec(`delete from ua_d`)
	exec(`insert into alist (name, acl) values ('null_hideip_group','hideip')`)
	exec(`insert into user_set (name, id, data) values ('acl','alice','null_hideip_group')`)
	// alice holds hideip and her only ua_d record has today=NULL.
	exec(`insert into ua_d (name, ip, ua, today, sub) values (?, ?, 'Mozilla/5.0', NULL, '0')`, "alice", "192.0.2.10")
	// bob is the normal co-tenant on the same IP (valid today, so the group's
	// max(today) is non-empty and the group survives the group query).
	exec(`insert into ua_d (name, ip, ua, today, sub) values (?, ?, 'Mozilla/5.1', ?, '0')`, "bob", "192.0.2.10", "20261010")
	// regular group: two non-hideip users with valid today -> shown.
	exec(`insert into ua_d (name, ip, ua, today, sub) values (?, ?, 'Mozilla/5.2', ?, '0')`, "carol", "198.51.100.20", "20261010")
	exec(`insert into ua_d (name, ip, ua, today, sub) values (?, ?, 'Mozilla/5.3', ?, '0')`, "dave", "198.51.100.20", "20261010")
	// regular group whose max(today) is NULL (both members NULL today) -> shown
	// with Date empty.
	exec(`insert into ua_d (name, ip, ua, today, sub) values (?, ?, 'Mozilla/5.4', NULL, '0')`, "erin", "203.0.113.9")
	exec(`insert into ua_d (name, ip, ua, today, sub) values (?, ?, 'Mozilla/5.5', NULL, '0')`, "frank", "203.0.113.9")
	return db
}

// Test_Get_ua_multiple_rows_null_today_hideip_excluded reproduces the Section 1
// case: a hideip member whose only ua_d record has today NULL must still be
// read (so the group is excluded). Without the fail-closed scan handling,
// Attach_names would drop the member and the hideip group would be exposed as
// safe.
func Test_Get_ua_multiple_rows_null_today_hideip_excluded(t *testing.T) {
	db := Setup_null_today_multiple_db(t)
	defer tool.DB_close(db)

	data := tool.Get_ua_multiple_rows(db, "", "recent", 0)
	by_ip := make(map[string]int)
	for _, g := range data {
		by_ip[g.IP]++
	}
	if _, exposed := by_ip["192.0.2.10"]; exposed {
		t.Fatalf("hideip group 192.0.2.10 (member alice has today=NULL) was exposed; want it excluded")
	}
	if _, shown := by_ip["198.51.100.20"]; !shown {
		t.Fatalf("regular group 198.51.100.20 missing; want it shown")
	}
}

// Test_Get_ua_multiple_rows_null_today_regular_shown verifies a regular (no
// hideip) group whose max(today) is NULL is still shown (it is safe to display),
// with an empty Date rather than being dropped or erroring.
func Test_Get_ua_multiple_rows_null_today_regular_shown(t *testing.T) {
	db := Setup_null_today_multiple_db(t)
	defer tool.DB_close(db)

	data := tool.Get_ua_multiple_rows(db, "", "recent", 0)
	var found bool
	for _, g := range data {
		if g.IP != "203.0.113.9" {
			continue
		}
		found = true
		if g.Count != 2 {
			t.Errorf("NULL-today regular group Count=%d want 2", g.Count)
		}
		if g.Date != "" {
			t.Errorf("NULL-today regular group Date=%q; want empty string", g.Date)
		}
		break
	}
	if !found {
		t.Fatalf("NULL-today regular group 203.0.113.9 not shown; want it present with empty Date")
	}
}

// Test_Get_ua_multiple_rows_past_end_early_exit verifies that a page past the
// end (offset == math.MaxInt, the value Page_offset returns for an overflowed
// page number) returns an empty page without walking the groups, i.e. without a
// full-table scan.
func Test_Get_ua_multiple_rows_past_end_early_exit(t *testing.T) {
	db := Setup_large_multiple_db(t)
	defer tool.DB_close(db)

	start := time.Now()
	past := tool.Get_ua_multiple_rows(db, "", "recent", math.MaxInt)
	elapsed := time.Since(start)
	if len(past) != 0 {
		t.Fatalf("past-end page returned %d groups; want 0", len(past))
	}
	if elapsed > 50*time.Millisecond {
		t.Errorf("past-end page took %v; want < 50ms (no full walk)", elapsed)
	}
}
