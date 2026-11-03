// BBS board + post slow-mode regression tests (record-based implementation).
//
// Two suites:
//   1. Core logic (explicit SQLite DBs): validates the enforcement helper,
//      validation, cooldowns, exemption, state-on-deny, and concurrency.
//   2. API integration (full schema via audit_setup): verifies the write paths
//      return the slow_mode error contract and that a denied write creates no
//      content.
//
// Cooldown state is derived from existing BBS write records in bbs_data:
//   - top-level posts/topics contribute via set_name='date' and 'user_id';
//   - comments/nested replies/discussion replies contribute via
//     set_name='comment_date' and 'comment_user_id'.
// No dedicated bss_slow_mode table is used or created.
//
// All tests run sequentially against isolated SQLite DBs (no t.Parallel).

package test

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"

	_ "modernc.org/sqlite"
	"opennamu/route"
	"opennamu/route/tool"
)

// ---------------------------------------------------------------------------
// Core-logic DB helpers
// ---------------------------------------------------------------------------

// Set_slow_test_db_config selects the SQL dialect for a manually opened test
// database and restores the environment-derived configuration during cleanup.
// DB_boot is required because DB helpers inspect the process-global backend,
// not only the environment variables.
func Set_slow_test_db_config(t *testing.T, db_name, db_type string) {
	t.Helper()
	old_db, had_db := os.LookupEnv("NAMU_DB")
	old_db_type, had_db_type := os.LookupEnv("NAMU_DB_TYPE")
	t.Cleanup(func() {
		if had_db {
			_ = os.Setenv("NAMU_DB", old_db)
		} else {
			_ = os.Unsetenv("NAMU_DB")
		}
		if had_db_type {
			_ = os.Setenv("NAMU_DB_TYPE", old_db_type)
		} else {
			_ = os.Unsetenv("NAMU_DB_TYPE")
		}
		tool.DB_boot()
	})
	if err := os.Setenv("NAMU_DB", db_name); err != nil {
		t.Fatalf("set test DB path: %v", err)
	}
	if err := os.Setenv("NAMU_DB_TYPE", db_type); err != nil {
		t.Fatalf("set test DB type: %v", err)
	}
	tool.DB_boot()
}

// new_slow_db opens an isolated SQLite DB with the tables the slow-mode code
// paths need: bbs_set, bbs_data, plus the auth tables (alist, user_set) used
// by permission resolution. The "ip" group is seeded so a valid IP with no
// user_set row falls back to it. No dedicated slow-mode table is created.
func new_slow_db(t *testing.T) *sql.DB {
	t.Helper()
	db_name := filepath.Join(t.TempDir(), "bbs_slow")
	Set_slow_test_db_config(t, db_name, "sqlite")
	db := new_test_db(t, db_name+".db")
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("core seed %s: %v", q, err)
		}
	}
	exec(`create table if not exists bbs_set (set_name, set_code, set_id, set_data)`)
	exec(`create table if not exists bbs_data (set_name, set_code, set_id, set_data)`)
	exec(`create table if not exists re_admin (who, what, time)`)
	exec(`insert into alist (name, acl) values ('ip', 'ip')`)
	return db
}

// nowWall returns a current local wall-clock string, matching Get_time.
func nowWall() string {
	return tool.Get_time()
}

// seedPostWrite records a top-level post (or discussion topic) write for
// `actor` in `board`, post code `postCode`, at the given wall time. This is
// what a successful new-post/topic write inserts (date + user_id).
func seedPostWrite(db *sql.DB, actor, board, postCode, walltime string) {
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`,
		"date", postCode, board, walltime); err != nil {
		panic("seedPostWrite date: " + err.Error())
	}
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`,
		"user_id", postCode, board, actor); err != nil {
		panic("seedPostWrite user_id: " + err.Error())
	}
}

// seedCommentWrite records a comment/nested-reply/discussion-reply write for
// `actor` inside post `postCode` of `board`, comment code `commentCode`, at
// the given wall time. This is what a successful comment write inserts
// (comment_date + comment_user_id, with set_id = board-post).
func seedCommentWrite(db *sql.DB, actor, board, postCode, commentCode, walltime string) {
	commentSetID := board + "-" + postCode
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`,
		"comment_date", commentCode, commentSetID, walltime); err != nil {
		panic("seedCommentWrite comment_date: " + err.Error())
	}
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`,
		"comment_user_id", commentCode, commentSetID, actor); err != nil {
		panic("seedCommentWrite comment_user_id: " + err.Error())
	}
}

// enforceOnce begins a transaction on db, runs Bbs_slow_enforce, and rolls back
// (the check is read-only). It mirrors how the write paths call it inside
// tool.DB_transaction. Returns the remaining cooldown (seconds); 0 means allowed.
func enforceOnce(db *sql.DB, actor, board, post string, boardInt, postInt int) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	remaining, err2 := tool.Bbs_slow_enforce(db, tx, actor, board, post, boardInt, postInt)
	_ = tx.Rollback()
	return remaining, err2
}

// setSlowBoard writes the board-level interval in bbs_set.
func setSlowBoard(db *sql.DB, t *testing.T, setID string, seconds int) {
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed %s: %v", q, err)
		}
	}
	exec(`delete from bbs_set where set_name = 'bbs_slow_mode_seconds' and set_id = ?`, setID)
	exec(`insert into bbs_set (set_name, set_code, set_id, set_data) values (?, '', ?, ?)`,
		"bbs_slow_mode_seconds", setID, strconv.Itoa(seconds))
}

// setSlowPost writes the post-level interval in bbs_data.
func setSlowPost(db *sql.DB, t *testing.T, setID, setCode string, seconds int) {
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed %s: %v", q, err)
		}
	}
	exec(`delete from bbs_data where set_name = 'post_slow_mode_seconds' and set_id = ? and set_code = ?`, setID, setCode)
	exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`,
		"post_slow_mode_seconds", setCode, setID, strconv.Itoa(seconds))
}

// grantExemptUser seeds a group carrying only slow_edit_pass and a user in it.
func grantExemptUser(db *sql.DB, t *testing.T, ip, group string) {
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed %s: %v", q, err)
		}
	}
	exec(`delete from alist where name = ?`, group)
	exec(`insert into alist (name, acl) values (?, 'slow_edit_pass')`, group)
	exec(`delete from user_set where id = ? and name = 'acl'`, ip)
	exec(`insert into user_set (name, id, data) values ('acl', ?, ?)`, ip, group)
}

// ---------------------------------------------------------------------------
// Core-logic tests
// ---------------------------------------------------------------------------

func Test_Bbs_slow_sqlite_fixture_after_mysql_backend(t *testing.T) {
	old_db, had_db := os.LookupEnv("NAMU_DB")
	old_db_type, had_db_type := os.LookupEnv("NAMU_DB_TYPE")
	t.Cleanup(func() {
		if had_db {
			_ = os.Setenv("NAMU_DB", old_db)
		} else {
			_ = os.Unsetenv("NAMU_DB")
		}
		if had_db_type {
			_ = os.Setenv("NAMU_DB_TYPE", old_db_type)
		} else {
			_ = os.Unsetenv("NAMU_DB_TYPE")
		}
		tool.DB_boot()
	})

	// Simulate a preceding MySQL audit fixture without opening a MySQL
	// connection. The local slow-mode fixture must not inherit its SQL dialect.
	t.Setenv("NAMU_DB", filepath.Join(t.TempDir(), "mysql_audit"))
	t.Setenv("NAMU_DB_TYPE", "mysql")
	tool.DB_boot()

	db := new_slow_db(t)
	defer tool.DB_close(db)
	setSlowBoard(db, t, "board", 300)
	seedPostWrite(db, "192.0.2.21", "board", "p1", nowWall())

	remaining, err := enforceOnce(db, "192.0.2.21", "board", "", 300, 0)
	if err != nil || remaining <= 0 {
		t.Fatalf("SQLite fixture lost actor match after MySQL setup: remaining=%d err=%v backend=%q equality=%q", remaining, err, tool.Get_DB_type(), tool.Exact_binary_eq("u.set_data"))
	}
}

func Test_BbsSlowMode_valid_and_seconds(t *testing.T) {
	type tc struct {
		raw   string
		valid bool
		want  int
	}
	cases := []tc{
		{"", true, 0},
		{"0", true, 0},
		{"1", true, 1},
		{"300", true, 300},
		{"86400", true, 86400},
		{"86401", false, 0},
		{"999999999999999999999", false, 0},
		{"abc", false, 0},
		{"12.5", false, 0},
		{"1.0", false, 0},
		{"-1", false, 0},
		{"-86400", false, 0},
		// Leading/trailing whitespace is input normalization and is accepted.
		{" 12", true, 12},
		{"12 ", true, 12},
		{"1a", false, 0},
		{"a1", false, 0},
		{"+", false, 0},
		{"1e2", false, 0},
	}
	for _, c := range cases {
		if got := tool.Bbs_slow_mode_valid(c.raw); got != c.valid {
			t.Errorf("Bbs_slow_mode_valid(%q) = %v, want %v", c.raw, got, c.valid)
		}
		if c.valid {
			if got := tool.Bbs_slow_mode_seconds(c.raw); got != c.want {
				t.Errorf("Bbs_slow_mode_seconds(%q) = %d, want %d", c.raw, got, c.want)
			}
		}
	}
}

func Test_Bbs_board_slow_interval(t *testing.T) {
	db := new_slow_db(t)
	defer tool.DB_close(db)
	setSlowBoard(db, t, "board_a", 120)
	if got := tool.Bbs_board_slow_interval(db, "board_a"); got != 120 {
		t.Fatalf("board interval board_a = %d, want 120", got)
	}
	if got := tool.Bbs_board_slow_interval(db, "missing"); got != 0 {
		t.Fatalf("board interval missing = %d, want 0", got)
	}
}

func Test_Bbs_post_slow_interval(t *testing.T) {
	db := new_slow_db(t)
	defer tool.DB_close(db)
	setSlowPost(db, t, "board_b", "p1", 60)
	if got := tool.Bbs_post_slow_interval(db, "board_b", "p1"); got != 60 {
		t.Fatalf("post interval board_b/p1 = %d, want 60", got)
	}
	if got := tool.Bbs_post_slow_interval(db, "board_b", "missing"); got != 0 {
		t.Fatalf("post interval missing = %d, want 0", got)
	}
}

func Test_Bbs_slow_mode_post_setting_api_auth_and_validation(t *testing.T) {
	db := apiSlowDB(t)
	defer tool.DB_close(db)
	const board, post, owner = "22", "1", "userA"
	bbs_comment_seed_board(db, t, board, post, "setting test", owner)
	seedBbsNonExempt(db, t, owner)
	admin := tool.Config{IP: "admin_ip", Session: &emptySession{}}
	user := tool.Config{IP: owner, Session: &emptySession{}}

	if response := route.Api_bbs_w_post_slow_mode(user, board, post); response["response"] != "require auth" {
		t.Fatalf("unauthorized setting read response=%v, want require auth", response["response"])
	}
	if response := route.Api_bbs_w_post_slow_mode_put(user, board, post, "120"); response["response"] != "require auth" {
		t.Fatalf("unauthorized setting write response=%v, want require auth", response["response"])
	}
	if response := route.Api_bbs_w_post_slow_mode_put(admin, board, post, "120"); response["response"] != "ok" {
		t.Fatalf("authorized setting write response=%v, want ok", response["response"])
	}
	if response := route.Api_bbs_w_post_slow_mode(admin, board, post); response["response"] != "ok" || response["data"] != "120" {
		t.Fatalf("authorized setting read=%v, want 120", response)
	}
	for _, invalid := range []string{"-1", "86401", "1.5", "abc"} {
		if response := route.Api_bbs_w_post_slow_mode_put(admin, board, post, invalid); response["response"] != "error" {
			t.Fatalf("invalid interval %q response=%v, want error", invalid, response["response"])
		}
		if got := tool.Bbs_post_slow_interval(db, board, post); got != 120 {
			t.Fatalf("invalid interval %q changed stored value to %d", invalid, got)
		}
	}

	if response := route.Api_bbs_w_set_put(user, board, "bbs_slow_mode_seconds", "30", ""); response["response"] != "require auth" {
		t.Fatalf("unauthorized board setting write response=%v, want require auth", response["response"])
	}
	if response := route.Api_bbs_w_set_put(admin, board, "bbs_slow_mode_seconds", " 30 ", ""); response["response"] != "ok" {
		t.Fatalf("normalized board setting write response=%v, want ok", response["response"])
	}
	if got := tool.Bbs_board_slow_interval(db, board); got != 30 {
		t.Fatalf("normalized board interval = %d, want 30", got)
	}
	if response := route.Api_bbs_w_set_put(admin, board, "bbs_slow_mode_seconds", "86401", ""); response["response"] != "error" {
		t.Fatalf("invalid board setting response=%v, want error", response["response"])
	}
	if got := tool.Bbs_board_slow_interval(db, board); got != 30 {
		t.Fatalf("invalid board setting changed stored value to %d", got)
	}
}

func Test_Bbs_slow_enforce_allow_no_cooldown(t *testing.T) {
	db := new_slow_db(t)
	defer tool.DB_close(db)
	// No intervals configured (0 disables the scope) and no prior write; both
	// scopes active but empty, so the write is allowed with no retry.
	if remaining, err := enforceOnce(db, "198.51.100.7", "board", "", 0, 0); err != nil || remaining != 0 {
		t.Fatalf("expect allow; got remaining=%d err=%v", remaining, err)
	}
	if remaining, err := enforceOnce(db, "", "board", "", 0, 0); err != nil || remaining != 0 {
		t.Fatalf("disabled slow mode must preserve identityless writes; remaining=%d err=%v", remaining, err)
	}
	if remaining, err := enforceOnce(db, "", "board", "", 60, 0); err != nil || remaining != 1 {
		t.Fatalf("active slow mode must fail closed with a one-second retry without an identity; remaining=%d err=%v", remaining, err)
	}
}

func Test_Bbs_slow_enforce_board_cooldown(t *testing.T) {
	db := new_slow_db(t)
	defer tool.DB_close(db)
	setSlowBoard(db, t, "board", 300)

	// First write: no prior board write for this actor -> allowed.
	if remaining, err := enforceOnce(db, "192.0.2.10", "board", "", 300, 0); err != nil || remaining != 0 {
		t.Fatalf("first write should be allowed; got remaining=%d err=%v", remaining, err)
	}
	// Simulate a successful new post in the board.
	seedPostWrite(db, "192.0.2.10", "board", "p1", nowWall())

	// Second write: the board cooldown is now active.
	remaining, err := enforceOnce(db, "192.0.2.10", "board", "", 300, 0)
	if remaining <= 0 {
		t.Fatalf("expected active cooldown; remaining=%d err=%v", remaining, err)
	}
	if remaining <= 0 || remaining > 300 {
		t.Fatalf("retry_after = %d, want in (0,300]", remaining)
	}

	// A different actor is not bound by this actor's cooldown.
	remaining2, err2 := enforceOnce(db, "203.0.113.5", "board", "", 300, 0)
	if !errors.Is(err2, nil) || remaining2 != 0 {
		t.Fatalf("actor 203.0.113.5 should be allowed; got remaining=%d err=%v", remaining2, err2)
	}
}

func Test_Bbs_slow_enforce_post_cooldown(t *testing.T) {
	db := new_slow_db(t)
	defer tool.DB_close(db)
	setSlowPost(db, t, "board", "p1", 300)

	// First write in post p1 -> allowed.
	if remaining, err := enforceOnce(db, "192.0.2.20", "board", "p1", 0, 300); err != nil || remaining != 0 {
		t.Fatalf("first write should be allowed; got remaining=%d err=%v", remaining, err)
	}
	// Simulate a successful comment in post p1.
	seedCommentWrite(db, "192.0.2.20", "board", "p1", "c1", nowWall())

	// Second write in post p1 -> post cooldown active.
	remaining, err := enforceOnce(db, "192.0.2.20", "board", "p1", 0, 300)
	if remaining <= 0 {
		t.Fatalf("expected active cooldown; remaining=%d err=%v", remaining, err)
	}
	if remaining <= 0 || remaining > 300 {
		t.Fatalf("retry_after = %d, want in (0,300]", remaining)
	}

	// A different post must not be affected.
	remaining2, err2 := enforceOnce(db, "192.0.2.20", "board", "p2", 0, 300)
	if !errors.Is(err2, nil) || remaining2 != 0 {
		t.Fatalf("post p2 should be allowed; got remaining=%d err=%v", remaining2, err2)
	}
	// A different actor must not be affected.
	remaining3, err3 := enforceOnce(db, "203.0.113.5", "board", "p1", 0, 300)
	if !errors.Is(err3, nil) || remaining3 != 0 {
		t.Fatalf("actor 203.0.113.5 should be allowed; got remaining=%d err=%v", remaining3, err3)
	}
}

func Test_Bbs_slow_enforce_combined_max(t *testing.T) {
	db := new_slow_db(t)
	defer tool.DB_close(db)
	setSlowBoard(db, t, "board", 120)
	setSlowPost(db, t, "board", "p1", 300)

	// First write (board+post) -> allowed.
	if remaining, err := enforceOnce(db, "192.0.2.30", "board", "p1", 120, 300); err != nil || remaining != 0 {
		t.Fatalf("first write should be allowed; got remaining=%d err=%v", remaining, err)
	}
	seedPostWrite(db, "192.0.2.30", "board", "p1", nowWall())

	// Combined write: the binding constraint is the 300s post interval, so the
	// remaining should reflect ~300, not the smaller 120s board interval.
	remaining, err := enforceOnce(db, "192.0.2.30", "board", "p1", 120, 300)
	if remaining <= 0 {
		t.Fatalf("expected active cooldown; remaining=%d err=%v", remaining, err)
	}
	if remaining < 100 || remaining > 400 {
		t.Fatalf("retry_after should reflect the 300s post interval, got %d", remaining)
	}
}

func Test_Bbs_slow_enforce_zero_no_bypass(t *testing.T) {
	// Zero post interval must NOT bypass an active board interval.
	db := new_slow_db(t)
	defer tool.DB_close(db)
	setSlowBoard(db, t, "board", 300)

	if remaining, err := enforceOnce(db, "192.0.2.40", "board", "p1", 300, 0); err != nil || remaining != 0 {
		t.Fatalf("first write should be allowed; got remaining=%d err=%v", remaining, err)
	}
	seedCommentWrite(db, "192.0.2.40", "board", "p1", "c1", nowWall())

	remaining, err := enforceOnce(db, "192.0.2.40", "board", "p1", 300, 0)
	if remaining <= 0 {
		t.Fatalf("zero post bypassed board; remaining=%d err=%v", remaining, err)
	}
	if remaining <= 0 {
		t.Fatalf("zero post should not bypass board; retry_after = %d", remaining)
	}

	// Zero board interval must NOT bypass an active post interval.
	db2 := new_slow_db(t)
	defer tool.DB_close(db2)
	setSlowPost(db2, t, "board", "p1", 300)
	if remaining, err := enforceOnce(db2, "192.0.2.50", "board", "p1", 0, 300); err != nil || remaining != 0 {
		t.Fatalf("first write should be allowed; got remaining=%d err=%v", remaining, err)
	}
	seedCommentWrite(db2, "192.0.2.50", "board", "p1", "c1", nowWall())

	remaining, err = enforceOnce(db2, "192.0.2.50", "board", "p1", 0, 300)
	if remaining <= 0 {
		t.Fatalf("zero board bypassed post; remaining=%d err=%v", remaining, err)
	}
	if remaining <= 0 {
		t.Fatalf("zero board should not bypass post; retry_after = %d", remaining)
	}
}

func Test_Bbs_slow_enforce_exemption(t *testing.T) {
	db := new_slow_db(t)
	defer tool.DB_close(db)
	setSlowBoard(db, t, "board", 30)
	grantExemptUser(db, t, "192.0.2.60", "pass_group")

	// Exempt user may write repeatedly without cooldown (even with prior writes).
	seedPostWrite(db, "192.0.2.60", "board", "p1", nowWall())
	for i := 0; i < 3; i++ {
		remaining, err := enforceOnce(db, "192.0.2.60", "board", "", 30, 0)
		if err != nil {
			t.Fatalf("exempt write %d error: %v", i, err)
		}
		if remaining != 0 {
			t.Fatalf("exempt write %d should be allowed; got remaining=%d", i, remaining)
		}
	}
	// Non-exempt user is bound by the cooldown.
	if remaining, err := enforceOnce(db, "192.0.2.61", "board", "p1", 30, 0); err != nil || remaining != 0 {
		t.Fatalf("plain first write should be allowed; got remaining=%d err=%v", remaining, err)
	}
	seedCommentWrite(db, "192.0.2.61", "board", "p1", "c1", nowWall())
	remaining, err := enforceOnce(db, "192.0.2.61", "board", "p1", 30, 0)
	if remaining <= 0 {
		t.Fatalf("plain second write should be denied, got %v", err)
	}
	if remaining <= 0 {
		t.Fatalf("plain second write retry_after = %d", remaining)
	}
}

func Test_Bbs_slow_expired(t *testing.T) {
	db := new_slow_db(t)
	defer tool.DB_close(db)
	setSlowBoard(db, t, "board", 60)
	// A write far in the past is already expired.
	seedPostWrite(db, "192.0.2.70", "board", "p1", "2020-01-01 00:00:00")
	remaining, err := enforceOnce(db, "192.0.2.70", "board", "", 60, 0)
	if !errors.Is(err, nil) || remaining != 0 {
		t.Fatalf("expired cooldown should be allowed; got remaining=%d err=%v", remaining, err)
	}
}

func Test_Bbs_slow_comment_edit_no_reset(t *testing.T) {
	// A comment edit (updating the comment text but NOT comment_date) must not
	// reset the comment's post-level cooldown.
	db := new_slow_db(t)
	defer tool.DB_close(db)
	setSlowPost(db, t, "board", "p1", 60)
	// Seed a comment with a known (now) timestamp.
	seedCommentWrite(db, "192.0.2.80", "board", "p1", "c1", nowWall())
	remaining, err := enforceOnce(db, "192.0.2.80", "board", "p1", 0, 60)
	if remaining <= 0 {
		t.Fatalf("expected post cooldown after a fresh comment; remaining=%d err=%v", remaining, err)
	}
	// Simulate a comment edit: update the comment text only, leave comment_date
	// unchanged. The post-level cooldown must still be active.
	if _, err := db.Exec(`update bbs_data set set_data = 'edited comment' where set_name = 'comment' and set_id = ? and set_code = ?`,
		"board-p1", "c1"); err != nil {
		_ = err
	}
	remaining2, err2 := enforceOnce(db, "192.0.2.80", "board", "p1", 0, 60)
	if remaining2 <= 0 {
		t.Fatalf("comment edit must not reset post cooldown; got remaining=%d err=%v", remaining2, err2)
	}
}

// new_wal_slow_db opens an isolated SQLite DB in WAL mode with a 5s busy
// timeout, matching the runtime connection settings used by openNAMU. WAL plus
// busy timeout let multiple concurrent writer transactions serialize on the
// write lock instead of failing with SQLITE_BUSY.
func new_wal_slow_db(t *testing.T) *sql.DB {
	db_name := filepath.Join(t.TempDir(), "wal_slow")
	Set_slow_test_db_config(t, db_name, "sqlite")
	dbPath := db_name + ".db"
	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout=5000")
	if err != nil {
		t.Fatalf("open wal db: %v", err)
	}
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("wal seed %s: %v", q, err)
		}
	}
	exec(`create table if not exists rb (block, end, today, blocker, why, band, ongoing, login)`)
	exec(`create table if not exists bbs_set (set_name, set_code, set_id, set_data)`)
	exec(`create table if not exists bbs_data (set_name, set_code, set_id, set_data)`)
	exec(`create table if not exists alist (name, acl)`)
	exec(`create table if not exists user_set (name, id, data)`)
	exec(`create table if not exists re_admin (who, what, time)`)
	exec(`insert into alist (name, acl) values ('ip', 'ip')`)
	return db
}

func Test_Bbs_slow_concurrent(t *testing.T) {
	db := new_wal_slow_db(t)
	defer tool.DB_close(db)
	setSlowBoard(db, t, "board", 300)

	// Simulate: one writer "wins" by doing the first write (seedPostWrite),
	// the rest see the active cooldown. Because enforce is read-only, exactly
	// one simulated successful write precedes the others in this deterministic
	// single-threaded sequence.
	seedPostWrite(db, "192.0.2.95", "board", "p1", nowWall())
	if remaining, err := enforceOnce(db, "192.0.2.95", "board", "", 300, 0); err != nil || remaining <= 0 {
		t.Fatalf("post-write should be denied; got remaining=%d err=%v", remaining, err)
	}
	// A brand-new actor is still allowed (independent identity).
	remaining, err := enforceOnce(db, "192.0.2.96", "board", "", 300, 0)
	if !errors.Is(err, nil) || remaining != 0 {
		t.Fatalf("new actor should be allowed; got remaining=%d err=%v", remaining, err)
	}
}

// ---------------------------------------------------------------------------
// API integration tests (full schema)
// ---------------------------------------------------------------------------

// seedBbsNonExempt seeds a user in the "ip" group (BBS-capable but not
// slow_edit_pass-exempt). The group is already present via audit_seed_auth.
func seedBbsNonExempt(db *sql.DB, t *testing.T, user string) {
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', ?, 'ip')`, user); err != nil {
		t.Fatalf("seed user %s: %v", user, err)
	}
}

// apiSlowDB returns a fully-initialized DB (schema + keeper + base auth).
func apiSlowDB(t *testing.T) *sql.DB {
	db := audit_setup(t, "sqlite")
	audit_seed_auth(db, t)
	var version_rows int
	if err := db.QueryRow(`select count(*) from other where name = 'ver'`).Scan(&version_rows); err != nil {
		t.Fatalf("count version anchor: %v", err)
	}
	if version_rows == 0 {
		if _, err := db.Exec(`insert into other (name, data, coverage) values ('ver', 'audit', '')`); err != nil {
			t.Fatalf("seed version anchor: %v", err)
		}
	}
	return db
}

func bbsSlowBoardSeed(db *sql.DB, t *testing.T, setID string, seconds int) {
	if setID == "-1" {
		var board_rows int
		if err := db.QueryRow(`select count(*) from bbs_set where set_name = 'bbs_name' and set_id = ?`, setID).Scan(&board_rows); err != nil {
			t.Fatalf("count thread board anchor: %v", err)
		}
		if board_rows == 0 {
			for _, value := range [][]string{{"bbs_name", "thread"}, {"bbs_type", "thread"}} {
				if _, err := db.Exec(`insert into bbs_set (set_name, set_code, set_id, set_data) values (?, '', ?, ?)`, value[0], setID, value[1]); err != nil {
					t.Fatalf("seed thread board anchor %s: %v", value[0], err)
				}
			}
		}
	}
	if _, err := db.Exec(`delete from bbs_set where set_name = 'bbs_slow_mode_seconds' and set_id = ?`, setID); err != nil {
		t.Fatalf("seed delete %s: %v", setID, err)
	}
	if _, err := db.Exec(`insert into bbs_set (set_name, set_code, set_id, set_data) values (?, '', ?, ?)`,
		"bbs_slow_mode_seconds", setID, strconv.Itoa(seconds)); err != nil {
		t.Fatalf("seed insert %s: %v", setID, err)
	}
}

func bbsSlowPostSeed(db *sql.DB, t *testing.T, setID, setCode string, seconds int) {
	if _, err := db.Exec(`delete from bbs_data where set_name = 'post_slow_mode_seconds' and set_id = ? and set_code = ?`, setID, setCode); err != nil {
		t.Fatalf("seed delete %s: %v", setID, err)
	}
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`,
		"post_slow_mode_seconds", setCode, setID, strconv.Itoa(seconds)); err != nil {
		t.Fatalf("seed insert %s: %v", setID, err)
	}
}

func bbsCommentPostCount(db *sql.DB, t *testing.T, setID, setCode string) int {
	parentID := setID + "-" + setCode
	var n int64
	query := `select count(*) from bbs_data where set_name = 'comment_user_id' and set_id = ?`
	if err := db.QueryRow(query, parentID).Scan(&n); err != nil {
		return 0
	}
	return int(n)
}

func Test_BbsSlow_api_comment_post_contract(t *testing.T) {
	db := apiSlowDB(t)
	defer tool.DB_close(db)
	const setID, setCode, user = "7", "1", "test_user"
	bbs_comment_seed_board(db, t, setID, setCode, "board post", "post_user")
	seedBbsNonExempt(db, t, user)
	bbsSlowPostSeed(db, t, setID, setCode, 30)

	cfg := func() tool.Config {
		return tool.Config{IP: user, Session: &emptySession{}}
	}
	r1 := route.Api_bbs_w_comment_post(cfg(), setID, setCode, "0", "comment one")
	if r1["response"] != "ok" {
		t.Fatalf("first comment response = %v, want ok", r1)
	}
	r2 := route.Api_bbs_w_comment_post(cfg(), setID, setCode, "0", "comment two")
	if r2["response"] != "error" || r2["data"] != "slow_mode" {
		t.Fatalf("second comment response=%q data=%q, want error/slow_mode", r2["response"], r2["data"])
	}
	if ra, ok := r2["retry_after_seconds"].(int); !ok || ra <= 0 {
		t.Fatalf("second comment retry_after_seconds missing or <=0: %v", r2["retry_after_seconds"])
	}
	// A denied comment must not add a comment row.
	if n := bbsCommentPostCount(db, t, setID, setCode); n != 1 {
		t.Fatalf("denied comment created content; comment count = %d, want 1", n)
	}
}

func Test_Bbs_slow_single_comment_deletion_erases_derived_cooldown_evidence(t *testing.T) {
	db := apiSlowDB(t)
	defer tool.DB_close(db)
	const board, post, user = "23", "1", "delete_comment_user"
	bbs_comment_seed_board(db, t, board, post, "board post", "post_owner")
	seedBbsNonExempt(db, t, user)
	bbsSlowBoardSeed(db, t, board, 60)
	cfg := func(ip string) tool.Config { return tool.Config{IP: ip, Session: &emptySession{}} }

	created := route.Api_bbs_w_comment_post(cfg(user), board, post, "0", "comment to delete")
	if created["response"] != "ok" {
		t.Fatalf("initial comment response=%v, want ok", created)
	}
	comment_num, ok := created["data"].(string)
	if !ok || comment_num == "" {
		t.Fatalf("comment code=%v, want non-empty string", created["data"])
	}
	blocked := route.Api_bbs_w_comment_post(cfg(user), board, post, "0", "blocked before delete")
	if blocked["response"] != "error" || blocked["data"] != "slow_mode" {
		t.Fatalf("cooldown was not active before comment deletion: response=%v data=%v", blocked["response"], blocked["data"])
	}
	deleted := route.Api_bbs_w_comment_one_delete(cfg("admin_ip"), board, post+"-"+comment_num)
	if deleted["response"] != "ok" {
		t.Fatalf("comment deletion response=%v, want ok", deleted)
	}
	second := route.Api_bbs_w_comment_post(cfg(user), board, post, "0", "write after deletion")
	if second["response"] != "ok" {
		t.Fatalf("comment deletion did not erase derived cooldown evidence as documented: response=%v data=%v", second["response"], second["data"])
	}
}

func Test_Bbs_slow_post_deletion_erases_derived_cooldown_evidence(t *testing.T) {
	// Hard deletion removes the content rows that currently serve as cooldown
	// history. With the no-ledger/no-new-schema constraint, this evasion cannot
	// be prevented without changing deletion semantics or storing persistent
	// history; this test records the existing, unavoidable behavior.
	db := apiSlowDB(t)
	defer tool.DB_close(db)
	const board, user = "24", "delete_post_user"
	bbs_comment_seed_board(db, t, board, "1", "first post", "post_owner")
	seedBbsNonExempt(db, t, user)
	bbsSlowBoardSeed(db, t, board, 60)
	cfg := func(ip string) tool.Config { return tool.Config{IP: ip, Session: &emptySession{}} }

	first := route.Api_bbs_w_comment_post(cfg(user), board, "1", "0", "comment before delete")
	if first["response"] != "ok" {
		t.Fatalf("initial comment response=%v, want ok", first)
	}
	blocked := route.Api_bbs_w_comment_post(cfg(user), board, "1", "0", "blocked before delete")
	if blocked["response"] != "error" || blocked["data"] != "slow_mode" {
		t.Fatalf("cooldown was not active before deletion: response=%v data=%v", blocked["response"], blocked["data"])
	}
	deleted := route.Api_bbs_w_delete(cfg("admin_ip"), board, "1")
	if deleted["response"] != "ok" {
		t.Fatalf("post deletion response=%v, want ok", deleted)
	}
	bbs_comment_seed_board(db, t, board, "2", "second post", "post_owner")
	allowed := route.Api_bbs_w_comment_post(cfg(user), board, "2", "0", "write after deletion")
	if allowed["response"] != "ok" {
		t.Fatalf("post deletion did not erase derived cooldown evidence as documented: response=%v data=%v", allowed["response"], allowed["data"])
	}
}

func Test_Bbs_slow_sqlite_simultaneous_api_writes(t *testing.T) {
	db := apiSlowDB(t)
	defer tool.DB_close(db)
	const board, post, user = "25", "1", "sqlite_race_user"
	bbs_comment_seed_board(db, t, board, post, "concurrent board post", "post_owner")
	for _, value := range [][]string{{"comment_count", "0"}, {"last_activity", tool.Get_time()}} {
		if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, value[0], post, board, value[1]); err != nil {
			t.Fatalf("seed post write-lock metadata %s: %v", value[0], err)
		}
	}
	seedBbsNonExempt(db, t, user)
	bbsSlowBoardSeed(db, t, board, 30)
	config := tool.Config{IP: user, Session: &emptySession{}}

	const workers = 8
	var allowed int
	var firstErr error
	var mutex sync.Mutex
	var wait_group sync.WaitGroup
	for worker := 0; worker < workers; worker++ {
		wait_group.Add(1)
		go func() {
			defer wait_group.Done()
			response := route.Api_bbs_w_comment_post(config, board, post, "0", "simultaneous comment")
			mutex.Lock()
			defer mutex.Unlock()
			if response["response"] == "ok" {
				allowed++
				return
			}
			if response["response"] != "error" || response["data"] != "slow_mode" {
				if firstErr == nil {
					firstErr = errors.New("unexpected concurrent write response")
				}
			}
		}()
	}
	wait_group.Wait()
	if firstErr != nil {
		t.Fatal(firstErr)
	}
	if allowed != 1 {
		t.Fatalf("serialized same-board writes: allowed=%d, want exactly 1", allowed)
	}
	var stored_comments int
	if err := db.QueryRow(`select count(*) from bbs_data where set_name = 'comment_user_id' and set_id = ? and set_data = ?`, board+"-"+post, user).Scan(&stored_comments); err != nil {
		t.Fatalf("count stored concurrent comments: %v", err)
	}
	if stored_comments != allowed {
		t.Fatalf("stored comment rows=%d, successful API writes=%d", stored_comments, allowed)
	}
}

func Test_BbsSlow_turso_comment_slow_mode(t *testing.T) {
	if !tool.Turso_platform_supported() {
		t.Skip("Turso is unsupported on this platform")
	}
	db := audit_setup(t, "turso")
	defer tool.DB_close(db)
	audit_seed_auth(db, t)
	const board, post = "31", "1"
	bbs_comment_seed_board(db, t, board, post, "Turso board post", "post_owner")
	for _, value := range [][]string{{"comment_count", "0"}, {"last_activity", tool.Get_time()}} {
		if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, value[0], post, board, value[1]); err != nil {
			t.Fatalf("seed post write-lock metadata %s: %v", value[0], err)
		}
	}
	bbsSlowBoardSeed(db, t, board, 30)
	config := func(user string) tool.Config { return tool.Config{IP: user, Session: &emptySession{}} }

	first := route.Api_bbs_w_comment_post(config("TursoUser"), board, post, "0", "first Turso comment")
	if first["response"] != "ok" {
		t.Fatalf("first Turso comment response=%v, want ok", first)
	}
	case_variant := route.Api_bbs_w_comment_post(config("tursouser"), board, post, "0", "case-variant Turso comment")
	if case_variant["response"] != "ok" {
		t.Fatalf("Turso case variant shared cooldown: response=%v data=%v", case_variant["response"], case_variant["data"])
	}
	repeat := route.Api_bbs_w_comment_post(config("TursoUser"), board, post, "0", "repeated Turso comment")
	if repeat["response"] != "error" || repeat["data"] != "slow_mode" {
		t.Fatalf("same Turso identity response=%v data=%v, want error/slow_mode", repeat["response"], repeat["data"])
	}
}

func Test_BbsSlow_api_thread_post_contract(t *testing.T) {
	db := apiSlowDB(t)
	defer tool.DB_close(db)
	const user = "test_user"
	seedBbsNonExempt(db, t, user)
	bbsSlowBoardSeed(db, t, "-1", 30) // thread board

	cfg := func() tool.Config {
		return tool.Config{IP: user, Session: &emptySession{}}
	}
	// New topic: board-level slow mode only.
	r1 := route.Api_thread_bbs_post(cfg(), "0", "", "content", "New Topic", "")
	if r1["response"] != "ok" {
		t.Fatalf("first thread post response = %v, want ok", r1)
	}
	r2 := route.Api_thread_bbs_post(cfg(), "0", "", "content", "Second Topic", "")
	if r2["response"] != "error" || r2["data"] != "slow_mode" {
		t.Fatalf("second thread post response=%q data=%q, want error/slow_mode", r2["response"], r2["data"])
	}
	if ra, ok := r2["retry_after_seconds"].(int); !ok || ra <= 0 {
		t.Fatalf("second thread post retry_after_seconds missing or <=0: %v", r2["retry_after_seconds"])
	}
}

// countThreadPosts counts existing thread topics so we can assert a denied
// write created no content.
func countThreadPosts(db *sql.DB, t *testing.T) int {
	var n int64
	if err := db.QueryRow(`select count(*) from bbs_data where set_name = 'title' and set_id = '-1'`).Scan(&n); err != nil {
		return 0
	}
	return int(n)
}

func Test_BbsSlow_api_thread_post_no_content_on_deny(t *testing.T) {
	db := apiSlowDB(t)
	defer tool.DB_close(db)
	const user = "test_user"
	seedBbsNonExempt(db, t, user)
	bbsSlowBoardSeed(db, t, "-1", 30)
	cfg := func() tool.Config {
		return tool.Config{IP: user, Session: &emptySession{}}
	}
	before := countThreadPosts(db, t)
	r1 := route.Api_thread_bbs_post(cfg(), "0", "", "content", "T1", "")
	if r1["response"] != "ok" {
		t.Fatalf("first thread post response = %v, want ok", r1)
	}
	r2 := route.Api_thread_bbs_post(cfg(), "0", "", "content", "T2", "")
	if r2["response"] != "error" || r2["data"] != "slow_mode" {
		t.Fatalf("second thread post response=%q data=%q, want error/slow_mode", r2["response"], r2["data"])
	}
	after := countThreadPosts(db, t)
	if after != before+1 {
		t.Fatalf("denied thread write created content; before=%d after=%d want +1", before, after)
	}
} // Test_Bbs_slow_timezone_roundtrip verifies that the wall-clock generated by
// tool.Get_time() round-trips through the same parse used by
// bbs_parse_last_write (time.ParseInLocation(layout, s, time.Local)) to an
// absolute Unix time equal to time.Now().Unix() within 2 seconds. This is the
// invariant that makes bbs_remaining an accurate elapsed time. It must hold
// under UTC and under DST timezones (run the suite with TZ= to exercise them).
func Test_Bbs_slow_timezone_roundtrip(t *testing.T) {
	const layout = "2006-01-02 15:04:05"
	nowUnix := time.Now().Unix()
	wall := tool.Get_time()
	parsed, err := time.ParseInLocation(layout, wall, time.Local)
	if err != nil {
		t.Fatalf("parse wall clock %q: %v", wall, err)
	}
	diff := nowUnix - parsed.Unix()
	if diff < 0 {
		diff = -diff
	}
	if diff > 2 {
		t.Errorf("Get_time/parse round-trip disagreement: now=%d parsed=%d diff=%ds", nowUnix, parsed.Unix(), diff)
	}
}

// Test_Bbs_slow_timezone_past_elapsed verifies that a wall clock N seconds in
// the past, parsed with the same rules, yields a remaining value equal to
// interval - N (i.e. elapsed time is computed in absolute Unix time, not in a
// way that depends on the current local offset).
func Test_Bbs_slow_timezone_past_elapsed(t *testing.T) {
	db := new_slow_db(t)
	defer tool.DB_close(db)
	seedBbsNonExempt(db, t, "userA")
	setSlowBoard(db, t, "board", 60*60) // one hour

	elapsed := 10 * time.Minute
	seedPostWrite(db, "userA", "board", "1",
		walltimeAt(time.Now().Add(-elapsed)))
	rem, err := enforceOnce(db, "userA", "board", "", 60*60, 0)
	if err != nil {
		t.Fatalf("enforceOnce: %v", err)
	}
	want := 60*60 - int(elapsed.Seconds())
	if rem < want || rem > want+5 {
		t.Errorf("elapsed-time mismatch: expected rem in [%d, %d], got %d", want, want+5, rem)
	}
}

func Test_Bbs_slow_remaining_at_boundaries(t *testing.T) {
	now := time.Date(2026, time.November, 3, 12, 0, 0, 0, time.Local)
	wall := func(at time.Time) string {
		return at.Format("2006-01-02 15:04:05")
	}
	cases := []struct {
		name     string
		last     string
		interval int
		want     int
	}{
		{"fresh_60_seconds", wall(now), 60, 60},
		{"one_second_before_60_boundary", wall(now.Add(-59 * time.Second)), 60, 1},
		{"exact_60_boundary", wall(now.Add(-60 * time.Second)), 60, 0},
		{"one_second_before_day_boundary", wall(now.Add(-86399 * time.Second)), 86400, 1},
		{"exact_day_boundary", wall(now.Add(-86400 * time.Second)), 86400, 0},
		{"future_timestamp_capped", wall(now.Add(time.Hour)), 60, 60},
		{"missing_timestamp", "", 60, 0},
		{"malformed_timestamp", "not-a-date", 60, 0},
		{"disabled_interval", wall(now), 0, 0},
	}
	for _, test_case := range cases {
		t.Run(test_case.name, func(t *testing.T) {
			if got := tool.Bbs_slow_remaining_at(test_case.last, test_case.interval, now); got != test_case.want {
				t.Fatalf("remaining(%q, %d, %s) = %d, want %d", test_case.last, test_case.interval, now, got, test_case.want)
			}
		})
	}
}

func Test_Bbs_slow_remaining_at_local_timezone(t *testing.T) {
	now := time.Date(2026, time.April, 5, 12, 0, 0, 0, time.Local)
	last := now.Add(-10 * time.Minute).Format("2006-01-02 15:04:05")
	got := tool.Bbs_slow_remaining_at(last, 60*60, now)
	if want := 50 * 60; got != want {
		t.Fatalf("local-zone elapsed time: got %d, want %d in %s", got, want, time.Local)
	}
}

func Test_Bbs_slow_remaining_at_new_york_spring_dst(t *testing.T) {
	if time.Local.String() != "America/New_York" {
		t.Skip("run with TZ=America/New_York to test local DST transition")
	}
	// 01:30 EST to 03:30 EDT is one elapsed hour on the 2026 spring-forward day.
	now := time.Date(2026, time.March, 8, 3, 30, 0, 0, time.Local)
	last := "2026-03-08 01:30:00"
	if got, want := tool.Bbs_slow_remaining_at(last, 2*60*60, now), 60*60; got != want {
		t.Fatalf("spring DST elapsed time: got %d, want %d in %s", got, want, time.Local)
	}
}
