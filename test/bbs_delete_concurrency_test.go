package test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"opennamu/route"
)

// -----------------------------------------------------------------------------
// Whole-board deletion vs concurrent BBS writes.
//
// Pre-fix audit: Api_bbs_delete used a plain transaction with no identifier lock
// scope, while comment/new-post writes performed existence checks outside their
// transaction. A pre-fix SQLite/Turso race reproduced a writer panic when the
// deletion removed its lock anchor; no pre-fix orphan row was observed. The
// regression cases below verify the shared-lock fix leaves no rows after delete.
//
// Invariants asserted after each race:
//   - a successfully deleted board leaves NO bbs_set/bbs_data rows for it;
//   - no orphaned content (set_id in the deleted board's namespace) remains;
//   - unrelated boards are not blocked by the delete;
//   - board recreation reuses the numeric ID and starts empty (no old content).
// -----------------------------------------------------------------------------

// Bbs_delete_test_setup returns a DB for the deletion-race tests plus seeds the
// required ACL grants: 'owner' (bbs_create/bbs_delete/bbs_comment_manage and the
// admin slow_edit_pass exemption) and 'ip' (bbs_use -> bbs_comment/bbs_view).
func Bbs_delete_test_setup(t *testing.T, backend string) *sql.DB {
	t.Helper()
	var db *sql.DB
	if backend == "mysql" {
		if os.Getenv("MYSQL_AUDIT_DSN") == "" {
			t.Skip("MYSQL_AUDIT_DSN is not configured for the deletion-race test")
		}
		db, _ = mysql_audit_setup(t)
	} else {
		db = audit_setup(t, backend)
		audit_seed_auth(db, t)
	}

	execute := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed %q: %v", q, err)
		}
	}
	// admin_ip: owner (delete/create/manage + slow exemption) + ip (comment/view).
	execute(`insert into user_set (name, id, data) values ('acl','admin_ip','owner')`)
	execute(`insert into user_set (name, id, data) values ('acl','admin_ip','ip')`)

	var anchor_count int
	if err := db.QueryRow("select count(*) from other where name='ver'").Scan(&anchor_count); err != nil || anchor_count == 0 {
		if _, err := db.Exec(`insert into other (name,data,coverage) values ('ver','test','')`); err != nil {
			t.Fatalf("seed other.ver anchor: %v", err)
		}
	}
	// reserved boards 0/-1/-2 (so positive IDs are deletable and reusable).
	for _, board := range [][3]string{{"0", "document_comment", "comment"}, {"-1", "thread", "thread"}, {"-2", "report_bbs", "comment"}} {
		execute(`insert into bbs_set (set_name,set_code,set_id,set_data) values (?,?,?,?)`, "bbs_name", "", board[0], board[1])
		execute(`insert into bbs_set (set_name,set_code,set_id,set_data) values (?,?,?,?)`, "bbs_type", "", board[0], board[2])
	}
	return db
}

// Bbs_delete_create_board makes a deletable board and returns its numeric ID.
func Bbs_delete_create_board(t *testing.T, db *sql.DB, board_name, board_type string) string {
	t.Helper()
	response := route.Api_bbs_make(Identifier_config("admin_ip"), board_name, board_type)
	if response["response"] != "ok" {
		t.Fatalf("Api_bbs_make(%q) response=%v", board_name, response)
	}
	var created_board_id string
	if err := db.QueryRow("select set_id from bbs_set where set_name='bbs_name' and set_code='' and set_data=?", board_name).Scan(&created_board_id); err != nil {
		t.Fatalf("read board %q ID: %v", board_name, err)
	}
	return created_board_id
}

// Bbs_delete_seed_post_and_comment adds a top-level post and one comment to the
// given board so a concurrent comment has an existing post to attach to. Returns
// the seed post code and the seed comment code.
func Bbs_delete_seed_post_and_comment(t *testing.T, db *sql.DB, board_id string) (post_code, comment_code string) {
	t.Helper()
	edit_response := route.Api_bbs_w_edit_post(Identifier_config("admin_ip"), board_id, "", "", "", "seed post body", "", "", "")
	if edit_response["response"] != "ok" {
		t.Fatalf("Api_bbs_w_edit_post(new post %q) response=%v", board_id, edit_response)
	}
	var created_post_code string
	if err := db.QueryRow("select set_code from bbs_data where set_name='title' and set_id=?", board_id).Scan(&created_post_code); err != nil {
		t.Fatalf("read new top-level post code for %q: %v", board_id, err)
	}
	post_code = created_post_code
	comment_response := route.Api_bbs_w_comment_post(Identifier_config("admin_ip"), board_id, post_code, "", "seeded comment for deletion race")
	if comment_response["response"] != "ok" {
		t.Fatalf("seeded comment %q: response=%v", board_id, comment_response)
	}
	// derive the seed comment code.
	var comment_code_from_db string
	if err := db.QueryRow("select set_code from bbs_data where set_name='comment' and set_id=? limit 1", board_id+"-"+post_code).Scan(&comment_code_from_db); err != nil {
		t.Fatalf("read seed comment code for %q-%s: %v", board_id, post_code, err)
	}
	comment_code = comment_code_from_db
	return
}

// Bbs_board_content counts all bbs_set and bbs_data rows for a
// board. Returns 0 only when the board namespace is fully empty.
func Bbs_board_content(db *sql.DB, board_id string) (int, error) {
	var board_rows int
	if err := db.QueryRow(`select count(*) from bbs_set where set_id = ?`, board_id).Scan(&board_rows); err != nil {
		return 0, err
	}
	var content int
	if err := db.QueryRow(`select count(*) from bbs_data where set_id = ? or set_id like ?`, board_id, board_id+"-%").Scan(&content); err != nil {
		return 0, err
	}
	return board_rows + content, nil
}

// Bbs_board_data_count returns the number of content rows in a board namespace.
func Bbs_board_data_count(db *sql.DB, board_id string) (int, error) {
	var count int
	if err := db.QueryRow(`select count(*) from bbs_data where set_id = ? or set_id like ?`, board_id, board_id+"-%").Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// Bbs_board_has_seed_post reports whether the board still has its top-level post.
func Bbs_board_has_seed_post(db *sql.DB, board_id string) (bool, error) {
	var post_count int
	if err := db.QueryRow("select count(*) from bbs_data where set_name='title' and set_id=?", board_id).Scan(&post_count); err != nil {
		return false, err
	}
	return post_count > 0, nil
}

// Bbs_check_goroutine_panic drains panic reports after all workers have exited.
// The nonblocking default is the termination condition because the channel is
// intentionally left open for the lifetime of the race.
func Bbs_check_goroutine_panic(t *testing.T, panic_ch chan error) {
	t.Helper()
	for {
		select {
		case panic_err, channel_open := <-panic_ch:
			if !channel_open {
				return
			}
			t.Errorf("goroutine panic: %v", panic_err)
		default:
			return
		}
	}
}

// bbs_race_timeout bounds each race so a hung worker cannot block the test.
const bbs_race_timeout = 30 * time.Second

// Bbs_delete_race runs the deletion race: one thread deletes board_id while the
// other performs a concurrent write. mode "comment" -> Api_bbs_w_comment_post on
// the seed post; mode "post" -> a new top-level post. Returns the delete result,
// remaining board rows, and the writer response.
//
// Workers report their response (or a recovered panic) through shared,
// mutex-guarded state. The caller always returns once all workers have
// completed or the race times out, so a hung or panicking worker never blocks
// the test main goroutine on a channel receive.
func Bbs_delete_race(t *testing.T, db *sql.DB, board_id, post_code, write_mode string) (delete_ok bool, orphan_count int, writer_response string) {
	t.Helper()
	start := make(chan struct{})
	var wait_group sync.WaitGroup
	wait_group.Add(2)
	var mutex sync.Mutex
	var delete_resp, writer_resp map[string]any
	var panic_errors []error

	go func() {
		defer wait_group.Done()
		defer func() {
			if panic_value := recover(); panic_value != nil {
				mutex.Lock()
				panic_errors = append(panic_errors, fmt.Errorf("delete goroutine panicked: %v", panic_value))
				mutex.Unlock()
			}
		}()
		<-start
		response := route.Api_bbs_delete(Identifier_config("admin_ip"), board_id)
		mutex.Lock()
		delete_resp = response
		mutex.Unlock()
	}()
	go func() {
		defer wait_group.Done()
		defer func() {
			if panic_value := recover(); panic_value != nil {
				mutex.Lock()
				panic_errors = append(panic_errors, fmt.Errorf("writer goroutine panicked: %v", panic_value))
				mutex.Unlock()
			}
		}()
		<-start
		var response map[string]any
		if write_mode == "post" {
			response = route.Api_bbs_w_edit_post(Identifier_config("admin_ip"), board_id, "", "", "", "racing top-level post", "", "", "")
		} else {
			response = route.Api_bbs_w_comment_post(Identifier_config("admin_ip"), board_id, post_code, "", "racing concurrent comment")
		}
		mutex.Lock()
		writer_resp = response
		mutex.Unlock()
	}()
	close(start)
	done := make(chan struct{})
	go func() {
		wait_group.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(bbs_race_timeout):
		t.Fatalf("race timed out after %v", bbs_race_timeout)
	}
	mutex.Lock()
	for _, panic_err := range panic_errors {
		t.Errorf("goroutine panic: %v", panic_err)
	}
	delete_ok = delete_resp["response"] == "ok"
	if writer_resp != nil {
		writer_response, _ = writer_resp["response"].(string)
	}
	if writer_response != "ok" && writer_response != "not exist" && writer_response != "require auth" {
		t.Errorf("writer returned unexpected response %q (map=%v)", writer_response, writer_resp)
	}
	mutex.Unlock()
	count, _ := Bbs_board_content(db, board_id)
	orphan_count = count
	return
}

// -----------------------------------------------------------------------------
// Scenario A: delete vs new comment; assert no post-delete orphan remains.
// -----------------------------------------------------------------------------
func Test_Bss_delete_vs_comment_orphan(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso", "mysql"} {
		t.Run(backend, func(t *testing.T) {
			db := Bbs_delete_test_setup(t, backend)
			const reps = 40
			var failures int
			for i := 0; i < reps; i++ {
				name := fmt.Sprintf("delc_%s_r%03d", backend, i)
				board_id := Bbs_delete_create_board(t, db, name, "comment")
				post_code, _ := Bbs_delete_seed_post_and_comment(t, db, board_id)
				delete_ok, orphan_count, _ := Bbs_delete_race(t, db, board_id, post_code, "comment")
				if !delete_ok {
					t.Fatalf("race %d (board %q) delete did not succeed", i, board_id)
				}
				if orphan_count > 0 {
					failures++
					t.Fatalf("race %d (board %q) left %d orphaned row(s) after a successful delete", i, board_id, orphan_count)
				}
			}
			t.Logf("backend=%s reps=%d failures=%d", backend, reps, failures)
		})
	}
}

// Scenario B: delete vs new top-level post; assert no post-delete orphan remains.
func Test_Bss_delete_vs_new_post_orphan(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso", "mysql"} {
		t.Run(backend, func(t *testing.T) {
			db := Bbs_delete_test_setup(t, backend)
			const reps = 120
			var failures int
			for i := 0; i < reps; i++ {
				name := fmt.Sprintf("delq_%s_r%03d", backend, i)
				board_id := Bbs_delete_create_board(t, db, name, "comment")
				post_code, _ := Bbs_delete_seed_post_and_comment(t, db, board_id)
				delete_ok, orphan_count, _ := Bbs_delete_race(t, db, board_id, post_code, "post")
				if !delete_ok {
					t.Fatalf("race %d (board %q) delete did not succeed", i, board_id)
				}
				if orphan_count > 0 {
					failures++
					t.Fatalf("race %d (board %q) left %d orphaned row(s) after a successful delete", i, board_id, orphan_count)
				}
			}
			t.Logf("backend=%s reps=%d failures=%d", backend, reps, failures)
		})
	}
}

// Scenario F: board recreation + numeric ID reuse. The new board must start
// empty (no content from the old board) and the ID must be reusable.
func Test_Bss_delete_recreate_id_reuse(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso", "mysql"} {
		t.Run(backend, func(t *testing.T) {
			db := Bbs_delete_test_setup(t, backend)
			admin := Identifier_config("admin_ip")

			// old board with content.
			old_id := Bbs_delete_create_board(t, db, "reuse_old", "comment")
			Bbs_delete_seed_post_and_comment(t, db, old_id)
			if content_count, _ := Bbs_board_content(db, old_id); content_count == 0 {
				t.Fatalf("seeded old board %q has no content", old_id)
			}

			// delete it.
			if response := route.Api_bbs_delete(admin, old_id); response["response"] != "ok" {
				t.Fatalf("delete old board %q: response=%v", old_id, response)
			}
			if remaining_content, _ := Bbs_board_content(db, old_id); remaining_content > 0 {
				t.Fatalf("old board %q not fully removed: content=%d", old_id, remaining_content)
			}

			// recreate; the ID should be reusable and the new board starts empty.
			new_id := Bbs_delete_create_board(t, db, "reuse_new", "comment")
			if new_id != old_id {
				t.Fatalf("ID reuse: expected recreated board to reuse %q, got %q", old_id, new_id)
			}
			if remaining_content, _ := Bbs_board_data_count(db, new_id); remaining_content != 0 {
				t.Fatalf("recreated board %q not empty: stale_content=%d", new_id, remaining_content)
			}

			// A legitimate write to the recreated board must produce content.
			Bbs_delete_seed_post_and_comment(t, db, new_id)
		})
	}
}

// Scenario G: unrelated-board independence. Deleting board A must not block or
// corrupt writes to board B.
func Test_Bss_delete_unrelated_board_independent(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso", "mysql"} {
		t.Run(backend, func(t *testing.T) {
			db := Bbs_delete_test_setup(t, backend)
			admin := Identifier_config("admin_ip")
			const reps = 40

			for i := 0; i < reps; i++ {
				a_id := Bbs_delete_create_board(t, db, fmt.Sprintf("deua_%s_a%03d", backend, i), "comment")
				b_id := Bbs_delete_create_board(t, db, fmt.Sprintf("deub_%s_b%03d", backend, i), "comment")
				Bbs_delete_seed_post_and_comment(t, db, a_id)
				_, b_post := Bbs_delete_seed_post_and_comment(t, db, b_id)

				start := make(chan struct{})
				var wait_group sync.WaitGroup
				wait_group.Add(2)
				err_ch := make(chan error, 2)
				go func() {
					defer wait_group.Done()
					defer func() {
						if panic_value := recover(); panic_value != nil {
							err_ch <- fmt.Errorf("delete panic: %v", panic_value)
						}
					}()
					<-start
					_ = route.Api_bbs_delete(admin, a_id)
				}()
				go func() {
					defer wait_group.Done()
					defer func() {
						if panic_value := recover(); panic_value != nil {
							err_ch <- fmt.Errorf("writer panic: %v", panic_value)
						}
					}()
					<-start
					response := route.Api_bbs_w_comment_post(admin, b_id, b_post, "", "unrelated board comment")
					if response["response"] != "ok" {
						err_ch <- fmt.Errorf("unrelated write response=%v", response)
					}
				}()
				close(start)
				wait_group.Wait()
				close(err_ch)
				for err := range err_ch {
					t.Fatalf("unrelated board %q: %v", b_id, err)
				}
				if ok, err := Bbs_board_has_seed_post(db, b_id); !ok || err != nil {
					t.Fatalf("unrelated board %q lost its seed post after concurrent delete of %q: ok=%v err=%v", b_id, a_id, ok, err)
				}
				if remaining_content, _ := Bbs_board_content(db, a_id); remaining_content > 0 {
					t.Fatalf("deleted board %q left content: content=%d", a_id, remaining_content)
				}
			}
		})
	}
}

// Scenario C: delete vs post/comment edit and deletion. These operations only
// touch existing content and must not leave orphaned rows.
func Test_Bss_delete_vs_comment_edit_or_delete(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso", "mysql"} {
		t.Run(backend, func(t *testing.T) {
			db := Bbs_delete_test_setup(t, backend)
			admin := Identifier_config("admin_ip")
			const reps = 60
			for i := 0; i < reps; i++ {
				board_id := Bbs_delete_create_board(t, db, fmt.Sprintf("delcd_%s_r%03d", backend, i), "comment")
				post_code, comment_code := Bbs_delete_seed_post_and_comment(t, db, board_id)
				start := make(chan struct{})
				var wait_group sync.WaitGroup
				wait_group.Add(2)
				err_ch := make(chan error, 2)
				var delete_ok bool
				go func() {
					defer wait_group.Done()
					defer func() {
						if panic_value := recover(); panic_value != nil {
							err_ch <- fmt.Errorf("delete panic: %v", panic_value)
						}
					}()
					<-start
					response := route.Api_bbs_delete(admin, board_id)
					delete_ok = response["response"] == "ok"
				}()
				go func() {
					defer wait_group.Done()
					defer func() {
						if panic_value := recover(); panic_value != nil {
							err_ch <- fmt.Errorf("writer panic: %v", panic_value)
						}
					}()
					<-start
					var response map[string]any
					switch i % 4 {
					case 0:
						response = route.Api_bbs_w_edit_post(admin, board_id, post_code, "", "edited title", "edited post body", "", "", "")
					case 1:
						response = route.Api_bbs_w_delete(admin, board_id, post_code)
					case 2:
						response = route.Api_bbs_w_edit_post(admin, board_id, post_code, comment_code, "", "edited comment", "", "", "")
					default:
						response = route.Api_bbs_w_comment_one_delete(admin, board_id, post_code+"-"+comment_code)
					}
					if response["response"] != "ok" && response["response"] != "not exist" && response["response"] != "require auth" {
						err_ch <- fmt.Errorf("writer response=%v", response)
					}
				}()
				close(start)
				wait_group.Wait()
				close(err_ch)
				for err := range err_ch {
					t.Fatalf("race %d failed: %v", i, err)
				}
				if !delete_ok {
					t.Fatalf("race %d delete did not succeed (board %q)", i, board_id)
				}
				if remaining_content, _ := Bbs_board_content(db, board_id); remaining_content > 0 {
					t.Fatalf("race %d left torn content after delete (board %q): content=%d", i, board_id, remaining_content)
				}
			}
		})
	}
}

// Scenario H: delete vs concurrent second delete. Must converge to a fully
// empty board with no partial state.
func Test_Bss_delete_vs_delete(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso", "mysql"} {
		t.Run(backend, func(t *testing.T) {
			db := Bbs_delete_test_setup(t, backend)
			admin := Identifier_config("admin_ip")
			const reps = 60
			for i := 0; i < reps; i++ {
				board_id := Bbs_delete_create_board(t, db, fmt.Sprintf("deldd_%s_r%03d", backend, i), "comment")
				Bbs_delete_seed_post_and_comment(t, db, board_id)
				start := make(chan struct{})
				var wait_group sync.WaitGroup
				panic_ch := make(chan error, 2)
				responses := make(chan map[string]any, 2)
				for d := 0; d < 2; d++ {
					wait_group.Add(1)
					go func() {
						defer wait_group.Done()
						defer func() {
							if panic_value := recover(); panic_value != nil {
								panic_ch <- fmt.Errorf("concurrent delete panicked: %v", panic_value)
							}
						}()
						<-start
						responses <- route.Api_bbs_delete(admin, board_id)
					}()
				}
				close(start)
				wait_group.Wait()
				close(responses)
				Bbs_check_goroutine_panic(t, panic_ch)
				success_count := 0
				for response := range responses {
					switch response["response"] {
					case "ok":
						success_count++
					case "not exist":
					default:
						t.Fatalf("concurrent delete returned unexpected response: %v", response)
					}
				}
				if success_count != 1 {
					t.Fatalf("concurrent deletion success count=%d, want exactly one", success_count)
				}
				if remaining_content, _ := Bbs_board_content(db, board_id); remaining_content > 0 {
					t.Fatalf("concurrent deletes left torn content (board %q): content=%d", board_id, remaining_content)
				}
			}
		})
	}
}

// Scenario I: lock lifecycle on MariaDB. A deletion transaction holds the board
// lock; while it is held, an unrelated board must still be writable (independence),
// and once the transaction completes the board lock must be released.
func Test_bbs_delete_mariadb_lock_lifecycle(t *testing.T) {
	if os.Getenv("MYSQL_AUDIT_DSN") == "" {
		t.Skip("MYSQL_AUDIT_DSN is not configured")
	}
	db, _ := mysql_audit_setup(t)
	admin := Identifier_config("admin_ip")
	var database_name string
	if err := db.QueryRow("select database()").Scan(&database_name); err != nil {
		t.Fatal(err)
	}
	execute := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed %q: %v", q, err)
		}
	}
	execute(`insert into user_set (name,id,data) values ('acl','admin_ip','owner')`)
	execute(`insert into user_set (name,id,data) values ('acl','admin_ip','ip')`)
	for _, board := range [][3]string{{"0", "document_comment", "comment"}, {"-1", "thread", "thread"}, {"-2", "report_bbs", "comment"}} {
		execute(`insert into bbs_set (set_name,set_code,set_id,set_data) values (?,?,?,?)`, "bbs_name", "", board[0], board[1])
		execute(`insert into bbs_set (set_name,set_code,set_id,set_data) values (?,?,?,?)`, "bbs_type", "", board[0], board[2])
	}
	observer, err := db.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer observer.Close()

	// Create two boards: deletion target A and unrelated board B.
	a_id := Bbs_delete_create_board(t, db, "lock_lifecycle_a", "comment")
	Bbs_delete_seed_post_and_comment(t, db, a_id)
	b_id := Bbs_delete_create_board(t, db, "lock_lifecycle_b", "comment")
	b_post, _ := Bbs_delete_seed_post_and_comment(t, db, b_id)

	a_lock := route.Bbs_identifier_lock_name(database_name, "board:"+a_id)
	b_lock := route.Bbs_identifier_lock_name(database_name, "board:"+b_id)
	if a_lock == b_lock {
		t.Fatalf("boards A(%q) and B(%q) share a lock scope: %q", a_id, b_id, a_lock)
	}

	// While a board-A write transaction holds the board lock, deletion of A
	// must wait, while an unrelated board B write must still proceed.
	held := make(chan struct{})
	release := make(chan struct{})
	transaction_done := make(chan error, 1)
	go func() {
		var transaction_err error
		defer func() {
			if panic_value := recover(); panic_value != nil {
				transaction_done <- fmt.Errorf("board write transaction panicked: %v", panic_value)
				return
			}
			transaction_done <- transaction_err
		}()
		transaction_err = route.Bbs_board_write_transaction(db, a_id, "", func(*sql.Tx) error {
			close(held)
			<-release
			return nil
		})
	}()
	select {
	case <-held:
	case transaction_err := <-transaction_done:
		close(release)
		t.Fatalf("board write transaction ended before acquiring its lock: %v", transaction_err)
	case <-time.After(bbs_race_timeout):
		close(release)
		t.Fatalf("board write transaction did not acquire its lock within %v", bbs_race_timeout)
	}
	type delete_call_result struct {
		response    map[string]any
		panic_error error
	}
	delete_done := make(chan delete_call_result, 1)
	go func() {
		result := delete_call_result{}
		defer func() {
			if panic_value := recover(); panic_value != nil {
				result.panic_error = fmt.Errorf("board delete panicked: %v", panic_value)
			}
			delete_done <- result
		}()
		result.response = route.Api_bbs_delete(admin, a_id)
	}()
	select {
	case delete_result := <-delete_done:
		close(release)
		<-transaction_done
		t.Fatalf("board A deletion passed an active board write lock: response=%v panic=%v", delete_result.response, delete_result.panic_error)
	case <-time.After(100 * time.Millisecond):
	}
	response := route.Api_bbs_w_comment_post(admin, b_id, b_post, "", "unrelated while A lock held")
	if response["response"] != "ok" {
		close(release)
		<-transaction_done
		<-delete_done
		t.Fatalf("unrelated board B write while A lock held: response=%v", response)
	}
	close(release)
	select {
	case transaction_err := <-transaction_done:
		if transaction_err != nil {
			t.Fatalf("board A transaction while holding lock: %v", transaction_err)
		}
	case <-time.After(bbs_race_timeout):
		t.Fatalf("board A transaction did not finish within %v", bbs_race_timeout)
	}
	var delete_result delete_call_result
	select {
	case delete_result = <-delete_done:
	case <-time.After(bbs_race_timeout):
		t.Fatalf("board A deletion did not finish within %v", bbs_race_timeout)
	}
	if delete_result.panic_error != nil {
		t.Fatalf("board A delete panicked: %v", delete_result.panic_error)
	}
	if delete_result.response["response"] != "ok" {
		t.Fatalf("board A delete response after lock release: %v", delete_result.response)
	}
	if content, err := Bbs_board_content(db, a_id); err != nil || content != 0 {
		t.Fatalf("board A remains after deletion: content=%d err=%v", content, err)
	}

	// After the transaction completed, board A's lock must be released.
	Assert_bbs_mariadb_lock_acquirable(t, observer, a_lock)
}
