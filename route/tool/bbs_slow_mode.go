package tool

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"
)

// bbs_slow_mode_max_seconds is the maximum allowed slow-mode interval.
const bbs_slow_mode_max_seconds = 86400

// Bbs_post_creation_time_name stores an immutable creation timestamp for new
// top-level posts in the existing bbs_data EAV structure. The public `date`
// record remains mutable for sorting and edit behavior.
const Bbs_post_creation_time_name = "creation_date"

// BbsSlowModeError is returned when a slow-mode cooldown is still active so
// the write is denied. Callers treat it as the "retry later" signal and read
// the retry-after from the remaining-seconds value returned by Bbs_slow_enforce.
var BbsSlowModeError = errors.New("bbs slow mode")

// bbs_parse_last_write parses a stored "2006-01-02 15:04:05" write timestamp
// (local wall clock, as produced by Get_time) into a Unix epoch in seconds.
// It returns (0, false) for empty or unparseable values.
func bbs_parse_last_write(raw string) (int64, bool) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return 0, false
	}
	t, err := time.ParseInLocation("2006-01-02 15:04:05", trimmed, time.Local)
	if err != nil {
		return 0, false
	}
	return t.Unix(), true
}

// Bbs_slow_remaining_at returns the remaining cooldown in seconds at a
// supplied evaluation time. Missing or malformed stored timestamps fail open
// to zero, matching legacy behavior. Future timestamps are treated as a full
// interval, preventing clock skew from producing a retry longer than configured.
func Bbs_slow_remaining_at(last string, interval int, now time.Time) int {
	if interval <= 0 || strings.TrimSpace(last) == "" {
		return 0
	}
	lastUnix, ok := bbs_parse_last_write(last)
	if !ok {
		return 0
	}
	elapsed := now.Unix() - lastUnix
	if elapsed < 0 {
		return interval
	}
	if elapsed >= int64(interval) {
		return 0
	}
	return interval - int(elapsed)
}

// Bbs_last_board_write returns the latest applicable write timestamp for
// `actor` in the whole `board`, as a "2006-01-02 15:04:05" string ("" if
// none). For a thread board it uses only the comment/reply records
// (set_name='comment_date' joined to comment_user_id), including the immutable
// initial topic comment, because in a thread the topic's set_name='date' row
// is overwritten on every comment by the thread write paths
// (route/api_bbs_w_comment_post.go, route/api_thread_bbs_post.go); that mutable
// date must not be attributed to the topic author. For a regular board it also
// includes top-level posts joined to user_id, preferring the immutable
// creation_date EAV record when available and falling back to legacy date rows.
// The runner should be the write transaction.
func Bbs_last_board_write(runner DB_runner, actor, board string) (string, error) {
	// Actor equality must be byte-exact (case-sensitive). MariaDB's default
	// table collation (utf8mb4_general_ci) is case-insensitive, so force a
	// binary comparison with the shared helper; SQLite/Turso are already binary.
	// The ordinary equality is a case-insensitive indexable prefilter on
	// MariaDB; the additional binary predicate keeps identity byte-exact.
	actorEq := "u.set_data = ? and " + Exact_binary_eq("u.set_data")
	dt := ""
	// thread_bbs_id is "-1" (route/api_thread_bbs.go). Only that board has its
	// post/topic date overwritten by other actors' comments, so the post part
	// must be excluded for it.
	if board == "-1" {
		q := `
			select coalesce(max(dt), '') from (
				select c.set_data as dt
				from bbs_data c
				join bbs_data u on u.set_id = c.set_id
				and u.set_code = c.set_code
				and u.set_name = 'comment_user_id'
				where c.set_name = 'comment_date'
				and ` + actorEq + `
				and c.set_id like ?
			) as w`
		stmt, err := runner.Prepare(q)
		if err != nil {
			return "", err
		}
		defer stmt.Close()
		if err := stmt.QueryRow(actor, actor, board+"-%").Scan(&dt); err != nil {
			return "", err
		}
		return dt, nil
	}
	q := `
		select coalesce(max(dt), '') from (
			select coalesce(nullif(created.set_data, ''), d.set_data) as dt
			from bbs_data d
			join bbs_data u on u.set_id = d.set_id
			and u.set_code = d.set_code
			and u.set_name = 'user_id'
			left join bbs_data created on created.set_id = d.set_id
			and created.set_code = d.set_code
			and created.set_name = '` + Bbs_post_creation_time_name + `'
			where d.set_id = ? and d.set_name = 'date'
			and ` + actorEq + `
			union all
			select c.set_data as dt
			from bbs_data c
			join bbs_data u on u.set_id = c.set_id
			and u.set_code = c.set_code
			and u.set_name = 'comment_user_id'
			where c.set_name = 'comment_date'
			and ` + actorEq + `
			and c.set_id like ?
		) as w`
	stmt, err := runner.Prepare(q)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	if err := stmt.QueryRow(board, actor, actor, actor, actor, board+"-%").Scan(&dt); err != nil {
		return "", err
	}
	return dt, nil
}

// Bbs_last_post_write returns the latest comment/reply timestamp for `actor`
// within post `post` of `board` (top-level and nested), as a
// "2006-01-02 15:04:05" string ("" if none). Only comments/replies count; a
// post edit never appears here. The runner should be the write transaction.
func Bbs_last_post_write(runner DB_runner, actor, board, post string) (string, error) {
	dt := ""
	comment_set_id := board + "-" + post
	comment_set_id_nested := comment_set_id + "-%"
	actor_eq := "u.set_data = ? and " + Exact_binary_eq("u.set_data")
	q := `
		select coalesce(max(dt), '') from (
			select c.set_data as dt
			from bbs_data c
			join bbs_data u on u.set_id = c.set_id
			and u.set_code = c.set_code
			and u.set_name = 'comment_user_id'
			where c.set_name = 'comment_date'
			and ` + actor_eq + `
			and c.set_id = ?
			union all
			select c.set_data as dt
			from bbs_data c
			join bbs_data u on u.set_id = c.set_id
			and u.set_code = c.set_code
			and u.set_name = 'comment_user_id'
			where c.set_name = 'comment_date'
			and ` + actor_eq + `
			and c.set_id like ?
		) as w`
	stmt, err := runner.Prepare(q)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	if err := stmt.QueryRow(actor, actor, comment_set_id, actor, actor, comment_set_id_nested).Scan(&dt); err != nil {
		return "", err
	}
	return dt, nil
}

// Bbs_slow_check computes the remaining combined cooldown for an actor writing
// in board/post, considering board_interval and post_interval (seconds). It
// returns the max of the two applicable remaining waits; 0 means allowed.
// The runner is normally the write transaction so the check reads the pre-insert
// state.
func Bbs_slow_check(runner DB_runner, actor, board, post string, board_interval, post_interval int) (int, error) {
	evaluation_time := time.Now()
	remaining := 0
	if board_interval > 0 {
		last, err := Bbs_last_board_write(runner, actor, board)
		if err != nil {
			return 0, err
		}
		if r := Bbs_slow_remaining_at(last, board_interval, evaluation_time); r > remaining {
			remaining = r
		}
	}
	if post_interval > 0 {
		last, err := Bbs_last_post_write(runner, actor, board, post)
		if err != nil {
			return 0, err
		}
		if r := Bbs_slow_remaining_at(last, post_interval, evaluation_time); r > remaining {
			remaining = r
		}
	}
	return remaining, nil
}

// Bbs_slow_enforce is the single entry point for slow-mode enforcement in
// write paths. It applies the slow_edit_pass exemption and then runs the
// record-based cooldown check. It returns the remaining cooldown in seconds:
// 0 means allowed, >0 means the write is denied and must retry after that many
// seconds.
func Bbs_slow_enforce(db *sql.DB, tx *sql.Tx, actor, board, post string, board_interval, post_interval int) (int, error) {
	if board_interval <= 0 && post_interval <= 0 {
		return 0, nil
	}
	if actor == "" {
		// Fail closed: an identityless write cannot be rate-limited.
		return 1, nil
	}
	if Check_permission(db, "slow_edit_pass", actor) {
		return 0, nil
	}
	return Bbs_slow_check(tx, actor, board, post, board_interval, post_interval)
}

// Bbs_slow_mode_normalize trims surrounding whitespace from a slow-mode setting.
func Bbs_slow_mode_normalize(raw string) string {
	return strings.TrimSpace(raw)
}

// Bbs_slow_mode_valid reports whether raw is an acceptable slow-mode setting
// value. Empty (->0/disabled) and integers in [0, bbs_slow_mode_max_seconds]
// are valid. Anything else is invalid.
func Bbs_slow_mode_valid(raw string) bool {
	trimmed := Bbs_slow_mode_normalize(raw)
	if trimmed == "" {
		return true
	}
	n, err := strconv.Atoi(trimmed)
	if err != nil {
		return false
	}
	if n < 0 || n > bbs_slow_mode_max_seconds {
		return false
	}
	return true
}

// Bbs_slow_mode_seconds parses raw into a slow-mode interval in seconds.
// It returns 0 for empty, malformed, out-of-range, or unparseable input.
// Use Bbs_slow_mode_valid at the write boundary to reject invalid values
// before this is relied upon.
func Bbs_slow_mode_seconds(raw string) int {
	trimmed := Bbs_slow_mode_normalize(raw)
	if trimmed == "" {
		return 0
	}
	n, err := strconv.Atoi(trimmed)
	if err != nil {
		return 0
	}
	if n < 0 || n > bbs_slow_mode_max_seconds {
		return 0
	}
	return n
}

// Bbs_board_slow_interval returns the board-level slow-mode interval (seconds)
// for the given board, or 0 if unset/invalid.
func Bbs_board_slow_interval(db *sql.DB, set_id string) int {
	return Bbs_slow_mode_seconds(Get_bbs_set_data(db, set_id, "bbs_slow_mode_seconds"))
}

// Bbs_post_slow_interval returns the post-level slow-mode interval (seconds)
// for the given post, or 0 if unset/invalid.
func Bbs_post_slow_interval(db *sql.DB, set_id, set_code string) int {
	raw, _ := Get_bbs_data_value(db, set_id, set_code, "post_slow_mode_seconds")
	return Bbs_slow_mode_seconds(raw)
}
