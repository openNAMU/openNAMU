package route

import (
	"database/sql"
	"math"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

// Create_user_check_db builds a minimal SQLite schema for testing user-check
// masking. It only contains the tables touched by the code under test.
func Create_user_check_db(t *testing.T) *sql.DB {
	t.Helper()
	db_path := filepath.Join(t.TempDir(), "user_check.db")
	db, err := sql.Open("sqlite", db_path)
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}
	if _, err := db.Exec(`CREATE TABLE alist (name TEXT, acl TEXT)`); err != nil {
		t.Fatalf("create alist table: %v", err)
	}
	if _, err := db.Exec(`CREATE TABLE user_set (name TEXT, id TEXT, data TEXT)`); err != nil {
		t.Fatalf("create user_set table: %v", err)
	}
	if _, err := db.Exec(`CREATE TABLE ua_d (name TEXT, ip TEXT, ua TEXT, today TEXT, sub TEXT)`); err != nil {
		t.Fatalf("create ua_d table: %v", err)
	}
	if _, err := db.Exec(`CREATE TABLE rb (block TEXT, end TEXT, today TEXT, blocker TEXT, why TEXT, band TEXT, login TEXT, ongoing TEXT)`); err != nil {
		t.Fatalf("create rb table: %v", err)
	}
	if _, err := db.Exec(`CREATE TABLE other (name TEXT, data TEXT, coverage TEXT)`); err != nil {
		t.Fatalf("create other table: %v", err)
	}

	// auth groups: inspector can check, victim can hide its own records, open has none.
	if _, err := db.Exec(`INSERT INTO alist (name, acl) VALUES ('inspector_group', 'check')`); err != nil {
		t.Fatalf("insert inspector acl: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO alist (name, acl) VALUES ('victim_group', 'hideip')`); err != nil {
		t.Fatalf("insert victim acl: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO alist (name, acl) VALUES ('open_group', '')`); err != nil {
		t.Fatalf("insert open acl: %v", err)
	}

	// per-user ACL assignments.
	if _, err := db.Exec(`INSERT INTO user_set (name, id, data) VALUES ('acl', 'inspector', 'inspector_group')`); err != nil {
		t.Fatalf("insert inspector acl entry: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO user_set (name, id, data) VALUES ('acl', 'victim', 'victim_group')`); err != nil {
		t.Fatalf("insert victim acl entry: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO user_set (name, id, data) VALUES ('acl', 'open_subject', 'open_group')`); err != nil {
		t.Fatalf("insert open acl entry: %v", err)
	}

	// recorded ua data.
	if _, err := db.Exec(`INSERT INTO ua_d (name, ip, ua, today, sub) VALUES ('victim', '192.168.0.10', 'Mozilla/5.0 (Victim UA)', '2026-01-01 00:00:00', '')`); err != nil {
		t.Fatalf("insert victim ua: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO ua_d (name, ip, ua, today, sub) VALUES ('open_subject', '192.168.0.20', 'Mozilla/5.0 (Open UA)', '2026-01-01 00:00:00', '')`); err != nil {
		t.Fatalf("insert open ua: %v", err)
	}

	// approval data for the victim, to verify QA suppression under hideip.
	if _, err := db.Exec(`INSERT INTO user_set (name, id, data) VALUES ('approval_question', 'victim', 'What is your favorite color?')`); err != nil {
		t.Fatalf("insert victim approval question: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO user_set (name, id, data) VALUES ('approval_question_answer', 'victim', 'Blue')`); err != nil {
		t.Fatalf("insert victim approval answer: %v", err)
	}

	return db
}

func Test_user_check_require_auth(t *testing.T) {
	db := Create_user_check_db(t)
	defer db.Close()

	// "no_auth" has no acl entry and resolves to the default "user" group,
	// which was never granted "check".
	result := Api_user_check_rows(db, "no_auth", "victim", "normal", 1, "")
	if response, ok := result["response"].(string); !ok || response != "require auth" {
		t.Fatalf("expected require auth, got response=%q rows=%v", response, result)
	}
}

func Test_user_check_masked_subject(t *testing.T) {
	db := Create_user_check_db(t)
	defer db.Close()

	result := Api_user_check_rows(db, "inspector", "victim", "normal", 1, "")
	if response, ok := result["response"].(string); !ok || response != "success" {
		t.Fatalf("expected success, got response=%q", response)
	}
	rows := result["rows"].([]map[string]string)
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	row := rows[0]
	if row["name"] != "victim" {
		t.Errorf("row name = %q, expected victim", row["name"])
	}
	if row["masked"] != "true" {
		t.Errorf("masked = %q, expected true for a subject with hideip", row["masked"])
	}
	if row["ip"] != "" {
		t.Errorf("masked row ip should be empty, got %q", row["ip"])
	}
	if row["ua"] != "" {
		t.Errorf("masked row ua should be empty, got %q", row["ua"])
	}
	show_qa, ok := result["show_qa"].(bool)
	if ok && show_qa {
		t.Errorf("show_qa should be false for a subject with hideip")
	}
}

func Test_user_check_unmasked_subject(t *testing.T) {
	db := Create_user_check_db(t)
	defer db.Close()

	result := Api_user_check_rows(db, "inspector", "open_subject", "normal", 1, "")
	if response, ok := result["response"].(string); !ok || response != "success" {
		t.Fatalf("expected success, got response=%q", response)
	}
	rows := result["rows"].([]map[string]string)
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	row := rows[0]
	if row["name"] != "open_subject" {
		t.Errorf("row name = %q, expected open_subject", row["name"])
	}
	if row["masked"] != "false" {
		t.Errorf("masked = %q, expected false for a subject without hideip", row["masked"])
	}
	if row["ip"] != "192.168.0.20" {
		t.Errorf("ip = %q, expected 192.168.0.20", row["ip"])
	}
	if row["ua"] != "Mozilla/5.0 (Open UA)" {
		t.Errorf("ua = %q, expected Mozilla/5.0 (Open UA)", row["ua"])
	}
}

func Test_user_check_simple_masked(t *testing.T) {
	db := Create_user_check_db(t)
	defer db.Close()

	result := Api_user_check_rows(db, "inspector", "victim", "simple", 1, "")
	if response, ok := result["response"].(string); !ok || response != "success" {
		t.Fatalf("expected success, got response=%q", response)
	}
	simple, ok := result["simple"].(bool)
	if !ok || !simple {
		t.Fatalf("expected simple mode to be true")
	}
	rows := result["rows"].([]map[string]string)
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	row := rows[0]
	if row["masked"] != "true" {
		t.Errorf("masked = %q, expected true", row["masked"])
	}
	if row["other_name"] != "" {
		t.Errorf("masked simple row other_name should be empty, got %q", row["other_name"])
	}
}

func Test_user_check_simple_ip_excludes_hideip(t *testing.T) {
	db := Create_user_check_db(t)
	defer db.Close()

	// IP reverse-lookup excludes hideip accounts entirely. "victim" holds
	// hideip, so no row may be returned for an IP used only by victim.
	result := Api_user_check_rows(db, "inspector", "192.168.0.10", "simple", 1, "")
	if response, ok := result["response"].(string); !ok || response != "success" {
		t.Fatalf("expected success, got response=%q", response)
	}
	rows := result["rows"].([]map[string]string)
	if len(rows) != 0 {
		t.Fatalf("expected 0 rows (hideip account excluded), got %d: %v", len(rows), rows)
	}
}

// Create_user_check_db_shared seeds the base user-check db and adds an IP
// (203.0.113.5) shared by the hideip user "victim" and the non-hideip user
// "open_subject", for IP reverse-lookup exclusion tests.
func Create_user_check_db_shared(t *testing.T) *sql.DB {
	t.Helper()
	db := Create_user_check_db(t)
	if _, err := db.Exec(
		`INSERT INTO ua_d (name, ip, ua, today, sub) VALUES (?, ?, ?, ?, ?)`,
		"open_subject", "203.0.113.5", "Mozilla/5.0 (Open shared)", "2026-01-01 00:00:00", ""); err != nil {
		t.Fatalf("insert open shared ua: %v", err)
	}
	if _, err := db.Exec(
		`INSERT INTO ua_d (name, ip, ua, today, sub) VALUES (?, ?, ?, ?, ?)`,
		"victim", "203.0.113.5", "Mozilla/5.0 (Victim shared)", "2026-01-01 00:00:00", ""); err != nil {
		t.Fatalf("insert victim shared ua: %v", err)
	}
	return db
}

func Test_user_check_simple_ip_keeps_non_hideip(t *testing.T) {
	db := Create_user_check_db_shared(t)
	defer db.Close()

	// 203.0.113.5 is used by both a hideip user (victim) and a non-hideip
	// user (open_subject). Only the non-hideip user may appear.
	result := Api_user_check_rows(db, "inspector", "203.0.113.5", "simple", 1, "")
	if response, ok := result["response"].(string); !ok || response != "success" {
		t.Fatalf("expected success, got response=%q", response)
	}
	rows := result["rows"].([]map[string]string)
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	row := rows[0]
	if row["masked"] != "false" {
		t.Errorf("masked = %q, expected false for a non-hideip user", row["masked"])
	}
	if row["other_name"] != "open_subject" {
		t.Errorf("other_name = %q, expected open_subject", row["other_name"])
	}
}

func Test_user_check_normal_ip_excludes_hideip(t *testing.T) {
	db := Create_user_check_db_shared(t)
	defer db.Close()

	result := Api_user_check_rows(db, "inspector", "203.0.113.5", "normal", 1, "")
	if response, ok := result["response"].(string); !ok || response != "success" {
		t.Fatalf("expected success, got response=%q", response)
	}
	rows := result["rows"].([]map[string]string)
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	row := rows[0]
	if row["name"] != "open_subject" {
		t.Errorf("name = %q, expected open_subject", row["name"])
	}
	if row["masked"] != "false" {
		t.Errorf("masked = %q, expected false for a non-hideip user", row["masked"])
	}
	if row["ip"] != "203.0.113.5" {
		t.Errorf("ip = %q, expected 203.0.113.5", row["ip"])
	}
}

func Test_user_check_compare_per_row_mask(t *testing.T) {
	db := Create_user_check_db(t)
	defer db.Close()

	result := Api_user_check_rows(db, "inspector", "victim", "normal", 1, "open_subject")
	if response, ok := result["response"].(string); !ok || response != "success" {
		t.Fatalf("expected success, got response=%q", response)
	}
	rows := result["rows"].([]map[string]string)
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	masked := map[string]int{}
	for _, row := range rows {
		masked[row["masked"]]++
	}
	if masked["true"] != 1 {
		t.Errorf("expected 1 masked row, got %d", masked["true"])
	}
		if masked["false"] != 1 {
			t.Errorf("expected 1 unmasked row, got %d", masked["false"])
		}
	}

// Create_compare_db extends the base user-check db with a mix of hideip and
// non-hideip accounts on distinct IPs so the compare per-row policy can be
// exercised for every name/plus_name combination.
func Create_compare_db(t *testing.T) *sql.DB {
	t.Helper()
	db := Create_user_check_db(t)
	exec := func(q string, a ...any) {
		if _, err := db.Exec(q, a...); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
	// alice and eve hold hideip (via the existing victim_group); bob/dave hold none.
	exec(`insert into user_set (name, id, data) values ('acl', 'alice', 'victim_group')`)
	exec(`insert into user_set (name, id, data) values ('acl', 'eve', 'victim_group')`)
	exec(`insert into ua_d (name, ip, ua, today, sub) values ('bob', '10.0.0.1', 'Mozilla/5.0 (Bob)', '2026-01-02', '')`)
	exec(`insert into ua_d (name, ip, ua, today, sub) values ('dave', '10.0.0.2', 'Mozilla/5.0 (Dave)', '2026-01-02', '')`)
	exec(`insert into ua_d (name, ip, ua, today, sub) values ('alice', '20.0.0.10', 'Mozilla/5.0 (Alice)', '2026-01-02', '')`)
	exec(`insert into ua_d (name, ip, ua, today, sub) values ('eve', '30.0.0.50', 'Mozilla/5.0 (Eve)', '2026-01-02', '')`)
	return db
}

// Collect_normal runs a normal-mode compare and returns {name: masked} for each
// returned row.
func Collect_normal(db *sql.DB, inspector, name, plus_name string) map[string]string {
	result := Api_user_check_rows(db, inspector, name, "normal", 1, plus_name)
	got := map[string]string{}
	for _, row := range result["rows"].([]map[string]string) {
		got[row["name"]] = row["masked"]
	}
	return got
}

// Same_set reports whether two {name: masked} maps are equal.
func Same_set(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		bv, ok := b[k]
		if !ok || bv != v {
			return false
		}
	}
	return true
}

// Test_normal_compare_per_row_policy verifies the per-row policy for every
// combination of name/plus_name type (account/account, account/IP, IP/account,
// IP/IP): hideip rows reached only through an IP are excluded, while hideip
// rows explicitly named on either side are kept but masked.
func Test_normal_compare_per_row_policy(t *testing.T) {
	db := Create_compare_db(t)
	defer db.Close()

	cases := []struct {
		name      string
		plus_name string
		want      map[string]string
	}{
		// A: both accounts, no hideip discovered -> both shown unmasked.
		{"bob", "dave", map[string]string{"bob": "false", "dave": "false"}},
		// B: subject account is hideip -> masked; other account unmasked.
		{"victim", "dave", map[string]string{"victim": "true", "dave": "false"}},
		// D (bug fix): account + IP must exclude hideip rows found via the IP;
		// only bob (the explicitly named account) is returned.
		{"bob", "20.0.0.10", map[string]string{"bob": "false"}},
		// E: IP + explicitly-named hideip account -> that account is masked;
		// the other IP co-tenants are excluded.
		{"20.0.0.10", "alice", map[string]string{"alice": "true"}},
		// F: both IPs, all co-tenants hideip, none explicit -> empty.
		{"20.0.0.10", "30.0.0.50", map[string]string{}},
		// C: IP only, hideip co-tenant -> empty.
		{"20.0.0.10", "", map[string]string{}},
	}
	for _, c := range cases {
		got := Collect_normal(db, "inspector", c.name, c.plus_name)
		if !Same_set(got, c.want) {
			t.Errorf("name=%q plus_name=%q: want %v, got %v", c.name, c.plus_name, c.want, got)
		}
	}
}

// Test_effective_permission_drive_exclude verifies exclusion is driven by the
// actual hideip permission (here a distinct group), not by a hardcoded group.
func Test_effective_permission_drive_exclude(t *testing.T) {
	db := Create_compare_db(t)
	defer db.Close()
	exec := func(q string, a ...any) {
		if _, err := db.Exec(q, a...); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
	exec(`insert into alist (name, acl) values ('alt_hideip_group','hideip')`)
	exec(`insert into user_set (name, id, data) values ('acl','frank','alt_hideip_group')`)
	exec(`insert into ua_d (name, ip, ua, today, sub) values ('frank','40.0.0.1','Mozilla/5.0','2026-01-02','')`)
	exec(`insert into ua_d (name, ip, ua, today, sub) values ('grace','40.0.0.1','Mozilla/5.0','2026-01-02','')`)

	got := Collect_normal(db, "inspector", "40.0.0.1", "")
	if !Same_set(got, map[string]string{"grace": "false"}) {
		t.Fatalf("IP reverse lookup: want %v, got %v", map[string]string{"grace": "false"}, got)
	}
}

// Test_user_check_page_number_guard verifies Page_offset's clamping and the
// early-exit guards: page numbers below 1 are clamped to page 1; past-the-end
// and overflowed page numbers return an empty "success" page (no panic, no full
// batch walk); valid pages return the expected rows.
func Test_user_check_page_number_guard(t *testing.T) {
	db := Create_user_check_db(t)
	defer db.Close()

	cases := []struct {
		name      string
		mode      string
		page_num  int
		want_rows int
	}{
		{"open_subject", "normal", 1, 1},
		{"open_subject", "normal", 2, 0},
		{"open_subject", "normal", 0, 1},
		{"open_subject", "normal", -1, 1},
		{"open_subject", "normal", 999999, 0},
		{"192.168.0.20", "simple", 1, 1},
		{"192.168.0.20", "simple", 2, 0},
	}
	for i, c := range cases {
		result := Api_user_check_rows(db, "inspector", c.name, c.mode, c.page_num, "")
		resp, ok := result["response"].(string)
		if !ok || resp != "success" {
			t.Fatalf("case %d (name=%q mode=%q page=%d): response=%q want success", i, c.name, c.mode, c.page_num, resp)
		}
		rows, _ := result["rows"].([]map[string]string)
		if len(rows) != c.want_rows {
			t.Errorf("case %d (name=%q mode=%q page=%d): row_count=%d want %d", i, c.name, c.mode, c.page_num, len(rows), c.want_rows)
		}
	}

	// Overflow boundary: Page_offset returns math.MaxInt for these pages; the
	// guard must yield an empty page without panicking or walking.
	for _, page := range []int{math.MaxInt / 50 + 2, math.MaxInt} {
		for _, mode := range []string{"normal", "simple"} {
			result := Api_user_check_rows(db, "inspector", "open_subject", mode, page, "")
			resp, ok := result["response"].(string)
			if !ok || resp != "success" {
				t.Fatalf("page=%d mode=%q: response=%q want success (no panic, no walk)", page, mode, resp)
			}
			rows, _ := result["rows"].([]map[string]string)
			if len(rows) != 0 {
				t.Errorf("page=%d mode=%q: row_count=%d want 0", page, mode, len(rows))
			}
		}
	}
}

// Test_user_check_past_end_empty verifies that an explicit math.MaxInt offset
// (the value Page_offset hands up for an overflowed page number) yields an
// empty page in both modes without a batch walk.
func Test_user_check_past_end_empty(t *testing.T) {
	db := Create_user_check_db_shared(t)
	defer db.Close()

	for _, mode := range []string{"normal", "simple"} {
		result := Api_user_check_rows(db, "inspector", "open_subject", mode, math.MaxInt, "")
		resp, ok := result["response"].(string)
		if !ok || resp != "success" {
			t.Fatalf("mode=%q past-end: response=%q want success", mode, resp)
		}
		rows, _ := result["rows"].([]map[string]string)
		if len(rows) != 0 {
			t.Fatalf("mode=%q past-end: row_count=%d want 0", mode, len(rows))
		}
	}
}
