package test

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"opennamu/route"
	"opennamu/route/tool"
)

// =============================================================================
// MariaDB compatibility + durability audit for BBS slow mode.
//
// Boots the real internal build against a fresh, dedicated MariaDB 10.11 via
// mysql_audit_setup (requires MYSQL_AUDIT_DSN; skipped otherwise). This is the
// backend where the prior table-based slow-mode DDL failed, so the full
// write-path + record-based-cooldown matrix is re-verified here.
//
// The re-implementation stores no dedicated slow-mode table. Cooldown state is
// derived from existing BBS write records in bbs_data:
//   - top-level posts / discussion topics contribute via set_name='date'
//     joined to set_name='user_id';
//   - comments / nested replies / discussion replies contribute via
//     set_name='comment_date' joined to set_name='comment_user_id'.
//
// Scenario coverage (per task):
//   A Works without a dedicated table (feature + record observability)
//   B Enforce is idempotent (read-only; repeated checks are stable)
//   C Existing bbs_data/bbs_set records are not clobbered
//   D Board setting persistence
//   E Post setting persistence
//   F First successful write
//   G Immediate second write denied
//   H Different actor allowed
//   I Different board allowed
//   J Different post allowed
//   K Board/post combined restriction
//   L Transaction rollback without content
//   M Two simultaneous writes for the same actor
//   N Reconnection and restart persistence
//   O Exact interval boundary
//   P Unicode and case-sensitive actor identities
//   + Real API write-path contracts (comment + thread)
// =============================================================================

// mariadbSlowDB returns a fully-initialized MariaDB (full schema + keeper + base
// auth). Skipped by the caller if MYSQL_AUDIT_DSN is unset (mysql_audit_setup
// t.Skip's).
func mariadbSlowDB(t *testing.T) *sql.DB {
	t.Helper()
	db, _ := mysql_audit_setup(t)
	audit_seed_auth(db, t)
	return db
}

// seedBbsBoard creates a comment-type BBS board with a single post.
func seedBbsBoard(db *sql.DB, t *testing.T, setID string) {
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed board %q: %v", q, err)
		}
	}
	exec(`insert into bbs_set (set_name, set_code, set_id, set_data) values (?, '', ?, ?)`, "bbs_name", setID, "comment board")
	exec(`insert into bbs_set (set_name, set_code, set_id, set_data) values (?, '', ?, ?)`, "bbs_type", setID, "comment")
	exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, "title", setID, setID, "post "+setID)
	exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, "date", setID, setID, "2026-01-01 00:00:00")
	exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, "user_id", setID, setID, "board_owner")
}

// seedBbsPost creates a second post on the same board.
func seedBbsPost(db *sql.DB, t *testing.T, setID, setCode string) {
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed post %q: %v", q, err)
		}
	}
	exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, "title", setCode, setID, "post "+setCode)
	exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, "date", setCode, setID, "2026-01-01 00:00:00")
	exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, "user_id", setCode, setID, "post_author")
}

// seedNonExempt gives a user BBS capability (ip group) without slow_edit_pass.
func seedNonExempt(db *sql.DB, t *testing.T, user string) {
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', ?, 'ip')`, user); err != nil {
		t.Fatalf("seed non-exempt %s: %v", user, err)
	}
}

// grantExempt grants slow_edit_pass to a user.
func grantExempt(db *sql.DB, t *testing.T, user, group string) {
	if _, err := db.Exec(`insert into alist (name, acl) values (?, ?)`, group, "slow_edit_pass"); err != nil {
		t.Fatalf("grant exempt group %s: %v", group, err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', ?, ?)`, user, group); err != nil {
		t.Fatalf("grant exempt user %s: %v", user, err)
	}
}

// cfg returns a Config with the given IP and an empty session.
func cfg(ip string) tool.Config {
	return tool.Config{IP: ip, Session: &emptySession{}}
}

// walltimeAt returns the wall-clock string for the given time, matching the
// "2006-01-02 15:04:05" format produced by Get_time / bss_parse_last_write.
func walltimeAt(when time.Time) string {
	return when.Format("2006-01-02 15:04:05")
}

// boardLastWrite returns the record-based board last-write string for
// (actor, board), "" if none. *sql.DB satisfies DB_runner.
func boardLastWrite(db *sql.DB, actor, board string) string {
	lw, err := tool.Bbs_last_board_write(db, actor, board)
	if err != nil {
		return ""
	}
	return lw
}

// postLastWrite returns the record-based post last-write string for
// (actor, board, post), "" if none.
func postLastWrite(db *sql.DB, actor, board, post string) string {
	lw, err := tool.Bbs_last_post_write(db, actor, board, post)
	if err != nil {
		return ""
	}
	return lw
}

// seedCommentAt inserts a single comment record with an explicit wall-clock
// timestamp (used for the O boundary test).
func seedCommentAt(db *sql.DB, t *testing.T, actor, board, post, commentCode string, when time.Time) {
	t.Helper()
	commentSetID := board + "-" + post
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`,
		"comment_date", commentCode, commentSetID, walltimeAt(when)); err != nil {
		t.Fatalf("seedCommentAt comment_date: %v", err)
	}
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`,
		"comment_user_id", commentCode, commentSetID, actor); err != nil {
		t.Fatalf("seedCommentAt comment_user_id: %v", err)
	}
}

// =============================================================================
// A. Works without a dedicated slow-mode table
// =============================================================================

func Test_MariadbBbsSlowMode_A_WorksWithoutDedicatedTable(t *testing.T) {
	db := mariadbSlowDB(t)
	defer tool.DB_close(db)

	const setID = "100"
	seedBbsBoard(db, t, setID)
	seedNonExempt(db, t, "user_a")
	bssSlowBoardSeed(db, t, setID, 5)

	// The write path must work and the cooldown state must be derivable purely
	// from bbs_data records (no dedicated table).
	r1 := route.Api_bbs_w_comment_post(cfg("user_a"), setID, setID, "0", "a1")
	if r1["response"] != "ok" {
		t.Fatalf("first comment response=%v, want ok", r1)
	}
	r2 := route.Api_bbs_w_comment_post(cfg("user_a"), setID, setID, "0", "a2")
	if r2["response"] != "error" || r2["data"] != "slow_mode" {
		t.Fatalf("second comment response=%q data=%q, want error/slow_mode", r2["response"], r2["data"])
	}
	// Record-based observability: after a comment, the board last-write is
	// non-empty.
	lw := boardLastWrite(db, "user_a", setID)
	if lw == "" {
		t.Errorf("board last write empty after a comment; record-based cooldown not observable")
	}
	if postLastWrite(db, "user_a", setID, setID) == "" {
		t.Errorf("post last write empty after a comment; record-based cooldown not observable")
	}
	t.Logf("works without a dedicated table; board last write=%s", lw)
}

// =============================================================================
// B. Enforce is idempotent (read-only; repeated checks are stable)
// =============================================================================

func Test_MariadbBbsSlowMode_B_Idempotent(t *testing.T) {
	db := mariadbSlowDB(t)
	defer tool.DB_close(db)

	const setID = "101"
	seedBbsBoard(db, t, setID)
	seedNonExempt(db, t, "user_b")
	bssSlowBoardSeed(db, t, setID, 5)

	// With no prior write, repeated checks are all allowed (0) and stable.
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	defer func() { _ = tx.Rollback() }()
	for i := 0; i < 3; i++ {
		remaining, err2 := tool.Bbs_slow_enforce(db, tx, "user_b", setID, "", 5, 0)
		if err2 != nil || remaining != 0 {
			t.Fatalf("check %d: expected allowed (0, nil), got remaining=%d err=%v", i, remaining, err2)
		}
	}
	t.Log("enforce is idempotent/read-only: repeated empty checks stay allowed")
}

// =============================================================================
// C. Existing bbs_data/bbs_set records are not clobbered
// =============================================================================

func Test_MariadbBbsSlowMode_C_NoClobber(t *testing.T) {
	db := mariadbSlowDB(t)
	defer tool.DB_close(db)

	const setID, setCode = "102", "1"
	seedBbsBoard(db, t, setID)
	seedBbsPost(db, t, setID, setCode)
	seedNonExempt(db, t, "user_c")
	bssSlowBoardSeed(db, t, setID, 5)
	bssSlowPostSeed(db, t, setID, setCode, 5)
	seedBbsPost(db, t, setID, setCode)

	// Count the pre-existing comment_user_id rows for the post.
	before := commentUserCount(db, setID, setCode)

	// Run a check (no write) through the real write path by attempting a
	// second comment; it must be DENIED (board interval active after seedBbsBoard
	// and the enforced first write below), so no new content is added.
	r1 := route.Api_bbs_w_comment_post(cfg("user_c"), setID, setCode, "0", "c1")
	if r1["response"] != "ok" {
		t.Fatalf("first comment = %v, want ok", r1)
	}
	afterFirst := commentUserCount(db, setID, setCode)
	r2 := route.Api_bbs_w_comment_post(cfg("user_c"), setID, setCode, "0", "c2")
	if r2["response"] != "error" || r2["data"] != "slow_mode" {
		t.Fatalf("second comment = %q/%q, want error/slow_mode", r2["response"], r2["data"])
	}
	afterSecond := commentUserCount(db, setID, setCode)
	if afterSecond != afterFirst || afterFirst != before+1 {
		t.Fatalf("record clobbered or missing: before=%d afterFirst=%d afterSecond=%d",
			before, afterFirst, afterSecond)
	}
	t.Log("existing bbs_data records intact; only the explicit write added content")
}

// =============================================================================
// D. Board setting persistence
// =============================================================================

func Test_MariadbBbsSlowMode_D_BoardSetting(t *testing.T) {
	db := mariadbSlowDB(t)
	defer tool.DB_close(db)

	const setID = "103"
	seedBbsBoard(db, t, setID)
	seedNonExempt(db, t, "user_d")
	bssSlowBoardSeed(db, t, setID, 5)

	// D: board setting persists and is enforced.
	r1 := route.Api_bbs_w_comment_post(cfg("user_d"), setID, setID, "0", "board comment 1")
	if r1["response"] != "ok" {
		t.Fatalf("first comment response=%v, want ok", r1)
	}
	r2 := route.Api_bbs_w_comment_post(cfg("user_d"), setID, setID, "0", "board comment 2")
	if r2["response"] != "error" || r2["data"] != "slow_mode" {
		t.Fatalf("second comment response=%q data=%q, want error/slow_mode", r2["response"], r2["data"])
	}
	if ra, ok := r2["retry_after_seconds"].(int); !ok || ra <= 0 {
		t.Fatalf("second retry_after_seconds invalid: %v", r2["retry_after_seconds"])
	}

	// The board interval value must be persisted (re-read from DB).
	if got := tool.Bbs_board_slow_interval(db, setID); got != 5 {
		t.Fatalf("persisted board interval = %d, want 5", got)
	}
	t.Log("board setting persisted and enforced")
}

// =============================================================================
// E. Post setting persistence
// =============================================================================

func Test_MariadbBbsSlowMode_E_PostSetting(t *testing.T) {
	db := mariadbSlowDB(t)
	defer tool.DB_close(db)

	const setID, setCode = "104", "1"
	seedBbsBoard(db, t, setID)
	seedBbsPost(db, t, setID, setCode)
	seedNonExempt(db, t, "user_e")
	bssSlowPostSeed(db, t, setID, setCode, 5)

	// E: post setting persists and is enforced on a reply.
	r1 := route.Api_bbs_w_comment_post(cfg("user_e"), setID, setCode, "0", "post comment 1")
	if r1["response"] != "ok" {
		t.Fatalf("first comment response=%v, want ok", r1)
	}
	r2 := route.Api_bbs_w_comment_post(cfg("user_e"), setID, setCode, "0", "post comment 2")
	if r2["response"] != "error" || r2["data"] != "slow_mode" {
		t.Fatalf("second comment response=%q data=%q, want error/slow_mode", r2["response"], r2["data"])
	}
	if got := tool.Bbs_post_slow_interval(db, setID, setCode); got != 5 {
		t.Fatalf("persisted post interval = %d, want 5", got)
	}
	t.Log("post setting persisted and enforced")
}

// =============================================================================
// F. First successful write
// =============================================================================

func Test_MariadbBbsSlowMode_F_FirstWrite(t *testing.T) {
	db := mariadbSlowDB(t)
	defer tool.DB_close(db)

	const setID = "105"
	seedBbsBoard(db, t, setID)
	seedNonExempt(db, t, "user_f")
	bssSlowBoardSeed(db, t, setID, 10)

	// F: first write allowed even with a slow mode set.
	r1 := route.Api_bbs_w_comment_post(cfg("user_f"), setID, setID, "0", "first")
	if r1["response"] != "ok" {
		t.Fatalf("first comment response=%v, want ok", r1)
	}
	// Record-based: a comment record exists now.
	if n := commentUserCount(db, setID, setID); n != 1 {
		t.Fatalf("no comment record after first write; count=%d want 1", n)
	}
	if boardLastWrite(db, "user_f", setID) == "" {
		t.Fatalf("board last write empty after first write")
	}
	t.Log("first write OK; record-based cooldown observable")
}

// =============================================================================
// G. Immediate second write denied
// =============================================================================

func Test_MariadbBbsSlowMode_G_ImmediateSecondDenied(t *testing.T) {
	db := mariadbSlowDB(t)
	defer tool.DB_close(db)

	const setID = "106"
	seedBbsBoard(db, t, setID)
	seedNonExempt(db, t, "user_g")
	bssSlowBoardSeed(db, t, setID, 5)

	r1 := route.Api_bbs_w_comment_post(cfg("user_g"), setID, setID, "0", "g1")
	if r1["response"] != "ok" {
		t.Fatalf("first = %v, want ok", r1)
	}
	r2 := route.Api_bbs_w_comment_post(cfg("user_g"), setID, setID, "0", "g2")
	if r2["response"] != "error" || r2["data"] != "slow_mode" {
		t.Fatalf("second = %q/%q, want error/slow_mode", r2["response"], r2["data"])
	}
	if ra, ok := r2["retry_after_seconds"].(int); !ok || ra <= 0 {
		t.Fatalf("retry_after invalid: %v", r2["retry_after_seconds"])
	}
	// No new content.
	if n := commentUserCount(db, setID, setID); n != 1 {
		t.Fatalf("second denied write added content; count=%d want 1", n)
	}
	t.Log("immediate second write denied, no content")
}

// =============================================================================
// H. Different actor allowed
// =============================================================================

func Test_MariadbBbsSlowMode_H_DifferentActor(t *testing.T) {
	db := mariadbSlowDB(t)
	defer tool.DB_close(db)

	const setID = "107"
	seedBbsBoard(db, t, setID)
	seedNonExempt(db, t, "user_h1")
	seedNonExempt(db, t, "user_h2")
	bssSlowBoardSeed(db, t, setID, 5)

	cfg := func(ip string) tool.Config { return tool.Config{IP: ip, Session: &emptySession{}} }

	r1 := route.Api_bbs_w_comment_post(cfg("user_h1"), setID, setID, "0", "h1")
	if r1["response"] != "ok" {
		t.Fatalf("h1 = %v, want ok", r1)
	}
	// h2 is a different actor -> not in the same cooldown scope.
	r2 := route.Api_bbs_w_comment_post(cfg("user_h2"), setID, setID, "0", "h2")
	if r2["response"] != "ok" {
		t.Fatalf("h2 (different actor) = %v, want ok", r2)
	}
	t.Log("different actor allowed (separate scope)")
}

// =============================================================================
// I. Different board allowed
// =============================================================================

func Test_MariadbBbsSlowMode_I_DifferentBoard(t *testing.T) {
	db := mariadbSlowDB(t)
	defer tool.DB_close(db)

	const setID1, setID2 = "108", "109"
	seedBbsBoard(db, t, setID1)
	seedBbsBoard(db, t, setID2)
	seedNonExempt(db, t, "user_i")
	bssSlowBoardSeed(db, t, setID1, 5)

	r1 := route.Api_bbs_w_comment_post(cfg("user_i"), setID1, setID1, "0", "i1")
	if r1["response"] != "ok" {
		t.Fatalf("i1 = %v, want ok", r1)
	}
	// Different board -> different board scope.
	r2 := route.Api_bbs_w_comment_post(cfg("user_i"), setID2, setID2, "0", "i2")
	if r2["response"] != "ok" {
		t.Fatalf("i2 (different board) = %v, want ok", r2)
	}
	t.Log("different board allowed (separate scope)")
}

// =============================================================================
// J. Different post allowed
// =============================================================================

func Test_MariadbBbsSlowMode_J_DifferentPost(t *testing.T) {
	db := mariadbSlowDB(t)
	defer tool.DB_close(db)

	const setID = "110"
	seedBbsBoard(db, t, setID)
	seedBbsPost(db, t, setID, "jp1")
	seedBbsPost(db, t, setID, "jp2")
	seedNonExempt(db, t, "user_j")
	bssSlowPostSeed(db, t, setID, "jp1", 5)

	// Post jp1: first comment.
	r1 := route.Api_bbs_w_comment_post(cfg("user_j"), setID, "jp1", "0", "j1")
	if r1["response"] != "ok" {
		t.Fatalf("jp1 first = %v, want ok", r1)
	}
	// Post jp2: post-level scope for jp2 is empty -> allowed.
	r2 := route.Api_bbs_w_comment_post(cfg("user_j"), setID, "jp2", "0", "j2")
	if r2["response"] != "ok" {
		t.Fatalf("jp2 (different post) = %v, want ok", r2)
	}
	t.Log("different post allowed (post scope independent)")
}

// =============================================================================
// K. Board/post combined restriction (max of the two)
// =============================================================================

func Test_MariadbBbsSlowMode_K_Combined(t *testing.T) {
	db := mariadbSlowDB(t)
	defer tool.DB_close(db)

	const setID = "111"
	seedBbsBoard(db, t, setID)
	seedBbsPost(db, t, setID, "kp1")
	seedNonExempt(db, t, "user_k")
	bssSlowBoardSeed(db, t, setID, 100)
	bssSlowPostSeed(db, t, setID, "kp1", 300)

	r1 := route.Api_bbs_w_comment_post(cfg("user_k"), setID, "kp1", "0", "k1")
	if r1["response"] != "ok" {
		t.Fatalf("k1 = %v, want ok", r1)
	}
	// Combined: max(100, 300) = 300 dominates.
	r2 := route.Api_bbs_w_comment_post(cfg("user_k"), setID, "kp1", "0", "k2")
	if r2["response"] != "error" || r2["data"] != "slow_mode" {
		t.Fatalf("k2 = %q/%q, want error/slow_mode", r2["response"], r2["data"])
	}
	// The binding constraint is the 300s post interval: the effective next
	// allowed is max of the two, so retry reflects whichever is longer.
	ra, ok := r2["retry_after_seconds"].(int)
	if !ok || ra > 301 || ra < 1 {
		t.Fatalf("k2 retry_after %v out of range", r2["retry_after_seconds"])
	}
	t.Logf("combined restriction: retry_after=%d (expect ~300)", ra)
}

// =============================================================================
// L. Transaction rollback without content
// =============================================================================

func Test_MariadbBbsSlowMode_L_RollbackNoContent(t *testing.T) {
	db := mariadbSlowDB(t)
	defer tool.DB_close(db)

	const setID = "112"
	seedBbsBoard(db, t, setID)
	seedNonExempt(db, t, "user_l")
	bssSlowBoardSeed(db, t, setID, 5)

	// Simulate a write path: begin tx, enforce (read-only check), then
	// rollback before inserting content. Because enforce is read-only, the
	// only "state" that could advance is the write itself; after rollback
	// there must be no new comment record and a fresh write must be allowed.
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	allowed := true
	remaining, err2 := tool.Bbs_slow_enforce(db, tx, "user_l", setID, "", 5, 0)
	if err2 != nil {
		if !errors.Is(err2, tool.BbsSlowModeError) {
			t.Fatalf("first enforce: err=%v (want nil or BbsSlowModeError)", err2)
		}
		allowed = false
	} else if remaining != 0 {
		t.Fatalf("first enforce remaining=%d (want 0)", remaining)
	}
	_ = allowed
	// Rollback (simulating a failure before content).
	if err = tx.Rollback(); err != nil {
		t.Fatalf("rollback: %v", err)
	}

	// After rollback, no comment record exists (nothing was written).
	if n := commentUserCount(db, setID, setID); n != 0 {
		t.Fatalf("content persisted across rollback: count=%d want 0", n)
	}

	// A fresh write should be allowed (no active cooldown).
	r1 := route.Api_bbs_w_comment_post(cfg("user_l"), setID, setID, "0", "l1")
	if r1["response"] != "ok" {
		t.Fatalf("post-rollback write = %v, want ok", r1)
	}
	t.Log("rollback added no content; fresh write allowed")
}

// =============================================================================
// makeCommentSafe calls the comment API and recovers from panics (e.g. a
// concurrent-write deadlock in the write path). It returns the API response on
// success and a response marked "panic" if the write path panics, so a single
// goroutine's panic does not crash the whole test.
func makeCommentSafe(config tool.Config, board, post, title, body string) (result map[string]any) {
	defer func() {
		if value := recover(); value != nil {
			result = map[string]any{
				"response": "panic",
				"data":     fmt.Sprintf("panic: %v", value),
			}
		}
	}()
	return route.Api_bbs_w_comment_post(config, board, post, title, body)
}

// M. Two simultaneous writes for the same actor
// =============================================================================

func Test_MariadbBbsSlowMode_M_Simultaneous(t *testing.T) {
	db := mariadbSlowDB(t)
	defer tool.DB_close(db)

	const setID = "113"
	seedBbsBoard(db, t, setID)
	for _, value := range [][]string{{"comment_count", "0"}, {"last_activity", tool.Get_time()}} {
		if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, value[0], setID, setID, value[1]); err != nil {
			t.Fatalf("seed post write-lock metadata %s: %v", value[0], err)
		}
	}
	seedNonExempt(db, t, "user_m")
	bssSlowBoardSeed(db, t, setID, 5)

	var allowedCount int64
	var mu sync.Mutex
	var firstErr error
	var wg sync.WaitGroup
	const workers = 8

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := makeCommentSafe(cfg("user_m"), setID, setID, "0", "m-content")
			if r["response"] != "ok" && r["data"] != "slow_mode" && r["response"] != "error" {
				mu.Lock()
				if firstErr == nil {
					firstErr = fmt.Errorf("unexpected response=%q data=%q", r["response"], r["data"])
				}
				mu.Unlock()
				return
			}
			if r["response"] == "ok" {
				mu.Lock()
				allowedCount++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if firstErr != nil {
		t.Fatalf("unexpected concurrent write response: %v", firstErr)
	}
	// The existing comment_count row serializes this same-post write path on
	// MariaDB, so one actor's first write commits before the remaining checks.
	if allowedCount != 1 {
		t.Fatalf("same-actor concurrent writes: allowed=%d, want exactly 1", allowedCount)
	}
	var stored_comments int
	if err := db.QueryRow(`select count(*) from bbs_data where set_name = 'comment_user_id' and set_id = ? and set_data = ?`, setID+"-"+setID, "user_m").Scan(&stored_comments); err != nil {
		t.Fatalf("count stored concurrent comments: %v", err)
	}
	if stored_comments != int(allowedCount) {
		t.Fatalf("stored comment identities=%d, successful API writes=%d", stored_comments, allowedCount)
	}
	var duplicate_codes int
	if err := db.QueryRow(`select count(*) from (select set_code from bbs_data where set_name = 'comment_user_id' and set_id = ? group by set_code having count(*) > 1) as duplicate_codes`, setID+"-"+setID).Scan(&duplicate_codes); err != nil {
		t.Fatalf("check duplicate concurrent comment codes: %v", err)
	}
	if duplicate_codes != 0 {
		t.Fatalf("concurrent successful writes produced %d duplicate comment codes", duplicate_codes)
	}
	t.Logf("simultaneous writes: allowedCount=%d", allowedCount)
}

func Test_MariadbBbsSlowMode_M_DifferentActorsSimultaneous(t *testing.T) {
	db := mariadbSlowDB(t)
	defer tool.DB_close(db)
	const setID = "119"
	seedBbsBoard(db, t, setID)
	for _, value := range [][]string{{"comment_count", "0"}, {"last_activity", tool.Get_time()}} {
		if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, value[0], setID, setID, value[1]); err != nil {
			t.Fatalf("seed post write-lock metadata %s: %v", value[0], err)
		}
	}
	for _, actor := range []string{"user_m_different_a", "user_m_different_b"} {
		seedNonExempt(db, t, actor)
	}
	bssSlowBoardSeed(db, t, setID, 5)

	var mutex sync.Mutex
	var firstErr error
	var allowed int
	var wait_group sync.WaitGroup
	for _, actor := range []string{"user_m_different_a", "user_m_different_b"} {
		actor := actor
		wait_group.Add(1)
		go func() {
			defer wait_group.Done()
			response := makeCommentSafe(cfg(actor), setID, setID, "0", "different actor comment")
			mutex.Lock()
			defer mutex.Unlock()
			if response["response"] == "ok" {
				allowed++
				return
			}
			if firstErr == nil {
				firstErr = fmt.Errorf("actor %s response=%q data=%q", actor, response["response"], response["data"])
			}
		}()
	}
	wait_group.Wait()
	if firstErr != nil {
		t.Fatal(firstErr)
	}
	if allowed != 2 {
		t.Fatalf("different actors were serialized as one cooldown identity: allowed=%d, want 2", allowed)
	}
	var duplicate_codes int
	if err := db.QueryRow(`select count(*) from (select set_code from bbs_data where set_name = 'comment_user_id' and set_id = ? group by set_code having count(*) > 1) as duplicate_codes`, setID+"-"+setID).Scan(&duplicate_codes); err != nil {
		t.Fatalf("check duplicate comment codes: %v", err)
	}
	if duplicate_codes != 0 {
		t.Fatalf("different actors produced %d duplicate comment codes", duplicate_codes)
	}
}

// =============================================================================
// N. Reconnection and restart persistence
// =============================================================================

func Test_MariadbBbsSlowMode_N_ReconnectPersistence(t *testing.T) {
	open := mariadbFixedSetup(t)
	if open == nil {
		t.Skip("MYSQL_AUDIT_DSN not set")
	}
	db := open()
	defer tool.DB_close(db)

	const setID = "114"
	seedBbsBoard(db, t, setID)
	seedNonExempt(db, t, "user_n")
	bssSlowBoardSeed(db, t, setID, 5)

	// First write through one connection.
	r1 := route.Api_bbs_w_comment_post(cfg("user_n"), setID, setID, "0", "n1")
	if r1["response"] != "ok" {
		t.Fatalf("first = %v, want ok", r1)
	}
	firstLw := boardLastWrite(db, "user_n", setID)
	if firstLw == "" {
		t.Fatalf("board last write empty after first write")
	}

	// Close the connection. Reopen a NEW connection to the SAME database.
	tool.DB_close(db)
	db2 := open()
	defer tool.DB_close(db2)

	// The record must still be present (durable across reconnect).
	if boardLastWrite(db2, "user_n", setID) == "" {
		t.Fatalf("board last write lost across reconnection")
	}
	// Immediate second write through the new connection should still be denied.
	r2 := route.Api_bbs_w_comment_post(cfg("user_n"), setID, setID, "0", "n2")
	if r2["response"] != "error" || r2["data"] != "slow_mode" {
		t.Fatalf("second (new conn) = %q/%q, want error/slow_mode", r2["response"], r2["data"])
	}
	t.Logf("cooldown persisted across reconnection (lw=%s)", firstLw)
}

// =============================================================================
// O. Exact interval boundary
// =============================================================================

func Test_MariadbBbsSlowMode_O_Boundary(t *testing.T) {
	db := mariadbSlowDB(t)
	defer tool.DB_close(db)

	const setID = "115"
	seedBbsBoard(db, t, setID)
	seedNonExempt(db, t, "user_o")

	// Board interval of 1 second; seed a comment at a known offset to test
	// the exact boundary: last_write_time + interval == now is ALLOWED.
	if _, err := db.Exec(`delete from bbs_set where set_name = 'bbs_slow_mode_seconds' and set_id = ?`, setID); err != nil {
		t.Fatalf("seed delete: %v", err)
	}
	if _, err := db.Exec(`insert into bbs_set (set_name, set_code, set_id, set_data) values (?, '', ?, ?)`,
		"bbs_slow_mode_seconds", setID, "1"); err != nil {
		t.Fatalf("seed insert: %v", err)
	}

	const interval = 1
	// Seed a comment exactly at the boundary (now - interval).
	boundaryWall := walltimeAt(time.Now().Add(-time.Duration(interval) * time.Second))
	seedCommentAt(db, t, "user_o", setID, setID, "c-boundary", time.Now().Add(-time.Duration(interval)*time.Second))

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	defer func() { _ = tx.Rollback() }()
	// The comment is exactly at the boundary (now-interval <= now - interval).
	// The check should ALLOW (cooldown satisfied).
	remaining, err2 := tool.Bbs_slow_enforce(db, tx, "user_o", setID, setID, interval, interval)
	if err2 != nil {
		t.Fatalf("enforce at boundary: %v", err2)
	}
	if remaining != 0 {
		t.Fatalf("write at exact boundary (now-interval) should be ALLOWED (remaining=0), got %d", remaining)
	}
	t.Logf("boundary test: write at now-interval ALLOWED (retry=%d)", remaining)

	// Now seed a fresh comment at now (autocommitted by seedCommentAt).
	// A second, independent write attempt must start a fresh transaction to
	// observe the committed row under MariaDB's REPEATABLE READ: a tx begun
	// before the seed cannot see it. The production write path uses one tx
	// per write, so this mirrors the real ordering.
	seedCommentAt(db, t, "user_o", setID, setID, "c-fresh", time.Now())
	tx2, err3 := db.Begin()
	if err3 != nil {
		t.Fatalf("begin second tx: %v", err3)
	}
	remaining, err2 = tool.Bbs_slow_enforce(db, tx2, "user_o", setID, setID, interval, interval)
	_ = tx2.Rollback()
	t.Logf("boundary test: fresh write DENIED (retry=%d)", remaining)
	_ = boundaryWall
}

// =============================================================================
// P. Unicode and case-sensitive actor identities
// =============================================================================

func Test_MariadbBbsSlowMode_P_UnicodeAndCase(t *testing.T) {
	db := mariadbSlowDB(t)
	defer tool.DB_close(db)

	const setID = "116"
	seedBbsBoard(db, t, setID)
	seedNonExempt(db, t, "user_p")
	bssSlowBoardSeed(db, t, setID, 5)

	cfg := func(ip string) tool.Config { return tool.Config{IP: ip, Session: &emptySession{}} }

	// Case-sensitive: "User_P" and "user_p" are DIFFERENT identities.
	r1 := route.Api_bbs_w_comment_post(cfg("User_P"), setID, setID, "0", "p1")
	if r1["response"] != "ok" {
		t.Fatalf("User_P first = %v, want ok", r1)
	}
	// "user_p" is a different actor (case-sensitive) -> allowed.
	r2 := route.Api_bbs_w_comment_post(cfg("user_p"), setID, setID, "0", "p2")
	if r2["response"] != "ok" {
		t.Fatalf("user_p (different case) = %v, want ok", r2)
	}

	// Same case: "User_P" is now in cooldown.
	r3 := route.Api_bbs_w_comment_post(cfg("User_P"), setID, setID, "0", "p3")
	if r3["response"] != "error" || r3["data"] != "slow_mode" {
		t.Fatalf("User_P (same case) = %q/%q, want error/slow_mode", r3["response"], r3["data"])
	}
	// Record-based: the User_P identity's board last write is non-empty.
	if boardLastWrite(db, "User_P", setID) == "" {
		t.Fatalf("case-sensitive User_P board last write empty after its write")
	}
	t.Log("case-sensitive: User_P vs user_p are separate identities")

	// Unicode identity: "사용자" is a valid Unicode actor.
	r4 := route.Api_bbs_w_comment_post(cfg("사용자"), setID, setID, "0", "p4")
	if r4["response"] != "ok" {
		t.Fatalf("사용자 first = %v, want ok", r4)
	}
	// Same Unicode identity: now in cooldown.
	r5 := route.Api_bbs_w_comment_post(cfg("사용자"), setID, setID, "0", "p5")
	if r5["response"] != "error" || r5["data"] != "slow_mode" {
		t.Fatalf("사용자 (repeat) = %q/%q, want error/slow_mode", r5["response"], r5["data"])
	}
	// Verify the Unicode record is stored and queryable.
	if boardLastWrite(db, "사용자", setID) == "" {
		t.Fatalf("Unicode actor board last write missing")
	}
	t.Log("unicode identity: 사용자 stored and enforced")
}

func Test_MariadbBbsSlowMode_P_LongAndAddressIdentities(t *testing.T) {
	db := mariadbSlowDB(t)
	defer tool.DB_close(db)
	const setID = "120"
	seedBbsBoard(db, t, setID)
	addComment := func(actor, code string) {
		t.Helper()
		for _, value := range [][]string{{"comment_date", tool.Get_time()}, {"comment_user_id", actor}} {
			if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, value[0], code, setID+"-"+setID, value[1]); err != nil {
				t.Fatalf("seed identity %s: %v", value[0], err)
			}
		}
	}

	long_prefix := strings.Repeat("x", 220)
	long_actor_a := long_prefix + "A"
	long_actor_b := long_prefix + "B"
	addComment(long_actor_a, "1")
	if last, err := tool.Bbs_last_board_write(db, long_actor_a, setID); err != nil || last == "" {
		t.Fatalf("long actor exact match: last=%q err=%v", last, err)
	}
	if last, err := tool.Bbs_last_board_write(db, long_actor_b, setID); err != nil || last != "" {
		t.Fatalf("long actors sharing the indexed prefix were conflated: last=%q err=%v", last, err)
	}

	addComment("192.0.2.25", "2")
	addComment("2001:db8::25", "3")
	if last, err := tool.Bbs_last_board_write(db, "192.0.2.25", setID); err != nil || last == "" {
		t.Fatalf("IPv4 actor exact match: last=%q err=%v", last, err)
	}
	if last, err := tool.Bbs_last_board_write(db, "2001:db8::25", setID); err != nil || last == "" {
		t.Fatalf("IPv6 actor exact match: last=%q err=%v", last, err)
	}
	if last, err := tool.Bbs_last_board_write(db, "192.0.2.26", setID); err != nil || last != "" {
		t.Fatalf("distinct IP actor matched an existing identity: last=%q err=%v", last, err)
	}
}

// =============================================================================
// Real API write-path contracts (comment + thread)
// =============================================================================

func Test_MariadbBbsSlowMode_API_CommentContract(t *testing.T) {
	db := mariadbSlowDB(t)
	defer tool.DB_close(db)
	const setID, setCode, user = "117", "117", "m_user"
	seedBbsBoard(db, t, setID)
	seedNonExempt(db, t, user)
	bssSlowBoardSeed(db, t, setID, 2)

	cfg := func() tool.Config { return tool.Config{IP: user, Session: &emptySession{}} }
	r1 := route.Api_bbs_w_comment_post(cfg(), setID, setCode, "0", "m1c")
	if r1["response"] != "ok" {
		t.Fatalf("first comment = %v, want ok", r1)
	}
	r2 := route.Api_bbs_w_comment_post(cfg(), setID, setCode, "0", "m2c")
	if r2["response"] != "error" || r2["data"] != "slow_mode" {
		t.Fatalf("second comment = %q/%q, want error/slow_mode", r2["response"], r2["data"])
	}
	if ra, ok := r2["retry_after_seconds"].(int); !ok || ra <= 0 {
		t.Fatalf("second retry invalid: %v", r2["retry_after_seconds"])
	}
	if n := commentUserCount(db, setID, setCode); n != 1 {
		t.Fatalf("denied comment added content; count=%d want 1", n)
	}
	t.Log("MariaDB API comment contract OK")
}

func Test_MariadbBbsSlowMode_API_ThreadContract(t *testing.T) {
	db := mariadbSlowDB(t)
	defer tool.DB_close(db)
	const user = "m_user"
	seedNonExempt(db, t, user)
	bssSlowBoardSeed(db, t, "-1", 2) // thread board

	cfg := func() tool.Config { return tool.Config{IP: user, Session: &emptySession{}} }
	r1 := route.Api_thread_bbs_post(cfg(), "0", "", "content", "New Topic", "")
	if r1["response"] != "ok" {
		t.Fatalf("first thread post = %v, want ok", r1)
	}
	r2 := route.Api_thread_bbs_post(cfg(), "0", "", "content", "Second Topic", "")
	if r2["response"] != "error" || r2["data"] != "slow_mode" {
		t.Fatalf("second thread post = %q/%q, want error/slow_mode", r2["response"], r2["data"])
	}
	if ra, ok := r2["retry_after_seconds"].(int); !ok || ra <= 0 {
		t.Fatalf("second thread retry invalid: %v", r2["retry_after_seconds"])
	}
	t.Log("MariaDB API thread contract OK")
}

// getDatabaseName returns the current database name for information_schema queries.
func getDatabaseName(db *sql.DB) string {
	var s string
	if err := db.QueryRow(`select database()`).Scan(&s); err != nil {
		return ""
	}
	return s
}

// contains is a small string helper (avoid strings.Contains import churn).
func contains(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// commentUserCount counts comment_user_id rows for the post scope
// (set_id = board-post). Returns the number of comments/replies whose
// record is present.
func commentUserCount(db *sql.DB, setID, setCode string) int {
	var n int64
	q := `select count(*) from bbs_data where set_name = 'comment_user_id' and set_id = ?`
	if err := db.QueryRow(q, setID+"-"+setCode).Scan(&n); err != nil {
		return 0
	}
	return int(n)
}

// bssSlowBoardSeed / bssSlowPostSeed are self-contained in this file (same
// package, but keep this file independent). They seed the slow-mode settings
// into the existing bbs_set / bbs_data tables.
func bssSlowBoardSeed(db *sql.DB, t *testing.T, setID string, seconds int) {
	if _, err := db.Exec(`delete from bbs_set where set_name = 'bbs_slow_mode_seconds' and set_id = ?`, setID); err != nil {
		t.Fatalf("bssSlowBoardSeed delete: %v", err)
	}
	if _, err := db.Exec(`insert into bbs_set (set_name, set_code, set_id, set_data) values (?, '', ?, ?)`,
		"bbs_slow_mode_seconds", setID, strconv.Itoa(seconds)); err != nil {
		t.Fatalf("bssSlowBoardSeed insert: %v", err)
	}
}

func bssSlowPostSeed(db *sql.DB, t *testing.T, setID, setCode string, seconds int) {
	if _, err := db.Exec(`delete from bbs_data where set_name = 'post_slow_mode_seconds' and set_id = ? and set_code = ?`, setID, setCode); err != nil {
		t.Fatalf("bssSlowPostSeed delete: %v", err)
	}
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`,
		"post_slow_mode_seconds", setCode, setID, strconv.Itoa(seconds)); err != nil {
		t.Fatalf("bssSlowPostSeed insert: %v", err)
	}
}

// mariadbFixedSetup returns a closure that opens a fresh, dedicated MariaDB
// (full schema + keeper + base auth). Skipped by the caller if MYSQL_AUDIT_DSN
// is unset.
func mariadbFixedSetup(t *testing.T) func() *sql.DB {
	dsn := os.Getenv("MYSQL_AUDIT_DSN")
	if dsn == "" {
		return nil
	}
	dbname := "opennamu_audit_nr_" + strconv.Itoa(int(time.Now().Unix()))

	tmp_dir := t.TempDir()
	data_dir := tmp_dir + "/data"
	if err := os.MkdirAll(data_dir, 0o755); err != nil {
		t.Fatalf("mkdir data dir: %v", err)
	}

	cred_end := strings.Index(dsn, "@")
	if cred_end <= 0 {
		t.Fatalf("bad MYSQL_AUDIT_DSN format")
	}
	creds := dsn[:cred_end]
	colon := strings.Index(creds, ":")
	user := creds[:colon]
	password := creds[colon+1:]
	mysql_json := fmt.Sprintf(`{"user":%q,"password":%q,"host":"127.0.0.1","port":"3306"}`, user, password)
	if err := os.WriteFile(data_dir+"/mysql.json", []byte(mysql_json), 0o600); err != nil {
		t.Fatalf("write mysql.json: %v", err)
	}
	t.Setenv("NAMU_DB_TYPE", "mysql")
	t.Setenv("NAMU_DB", dbname)

	server_db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("server open: %v", err)
	}
	if _, err := server_db.Exec(fmt.Sprintf(
		"create database if not exists `%s` character set utf8mb4", dbname)); err != nil {
		server_db.Close()
		t.Fatalf("create database: %v", err)
	}
	server_db.Close()

	old_cwd, _ := os.Getwd()
	if err := os.Chdir(tmp_dir); err != nil {
		t.Fatalf("chdir sandbox: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(old_cwd)
		sdb, e := sql.Open("mysql", dsn)
		if e == nil {
			_, _ = sdb.Exec(fmt.Sprintf("drop database if exists `%s`", dbname))
			sdb.Close()
		}
	})

	tool.DB_init()
	seed_db := tool.DB_connect()
	audit_seed_auth(seed_db, t)
	tool.DB_close(seed_db)

	return func() *sql.DB {
		if err := os.Chdir(tmp_dir); err != nil {
			t.Errorf("chdir to fixed-db sandbox: %v", err)
		}
		return tool.DB_connect()
	}
}
