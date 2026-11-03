package test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"opennamu/route"
)

type bbs_board_creation_result struct {
	name        string
	board_type  string
	response    string
	panic_value string
}

func Test_Bbs_board_creation_sequential_reserved_ids_and_reuse(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso", "mysql"} {
		t.Run(backend, func(t *testing.T) {
			db := Bbs_board_creation_audit_setup(t, backend, true)

			Bbs_assert_board_creation_response(t, route.Api_bbs_make(Identifier_config("admin_ip"), "board_seq_make", "thread"))
			Bbs_assert_board_id_and_type(t, db, "board_seq_make", "1", "thread")

			Bbs_assert_board_creation_response(t, route.Api_bbs_list_post(Identifier_config("admin_ip"), "board_seq_list_post", "comment"))
			Bbs_assert_board_id_and_type(t, db, "board_seq_list_post", "2", "comment")

			Bbs_assert_board_creation_response(t, route.Api_bbs_make(Identifier_config("admin_ip"), "board_seq_fallback_type", "invalid"))
			Bbs_assert_board_id_and_type(t, db, "board_seq_fallback_type", "3", "comment")

			deleted := route.Api_bbs_delete(Identifier_config("admin_ip"), "3")
			if deleted["response"] != "ok" {
				t.Fatalf("delete highest board response=%v, want ok", deleted)
			}
			Bbs_assert_board_creation_response(t, route.Api_bbs_list_post(Identifier_config("admin_ip"), "board_after_delete_reuses_id", "thread"))
			Bbs_assert_board_id_and_type(t, db, "board_after_delete_reuses_id", "3", "thread")

			special_delete := route.Api_bbs_delete(Identifier_config("admin_ip"), "-1")
			if special_delete["response"] != "error" || special_delete["data"] != "not allowed" {
				t.Fatalf("delete reserved thread board response=%v, want not allowed", special_delete)
			}
			for _, reserved := range [][3]string{{"0", "document_comment", "comment"}, {"-1", "thread", "thread"}, {"-2", "report_bbs", "comment"}} {
				Bbs_assert_board_id_and_type(t, db, reserved[1], reserved[0], reserved[2])
			}
		})
	}
}

func Test_Bbs_board_creation_empty_namespace_defaults(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso", "mysql"} {
		t.Run(backend, func(t *testing.T) {
			db := Bbs_board_creation_audit_setup(t, backend, false)
			var name_rows int
			if err := db.QueryRow("select count(*) from bbs_set where set_name = 'bbs_name'").Scan(&name_rows); err != nil || name_rows != 0 {
				t.Fatalf("empty namespace has %d names, err=%v", name_rows, err)
			}

			Bbs_assert_board_creation_response(t, route.Api_bbs_make(Identifier_config("admin_ip"), "empty_namespace_make", "thread"))
			Bbs_assert_board_id_and_type(t, db, "empty_namespace_make", "2", "thread")

			// Test the other legacy default independently while retaining the
			// same initialized local lock anchor (`other.ver`).
			if _, err := db.Exec("delete from bbs_set"); err != nil {
				t.Fatal(err)
			}
			Bbs_assert_board_creation_response(t, route.Api_bbs_list_post(Identifier_config("admin_ip"), "empty_namespace_list_post", "comment"))
			Bbs_assert_board_id_and_type(t, db, "empty_namespace_list_post", "1", "comment")
		})
	}
}

func Test_Bbs_board_creation_permission_and_transaction_rollback(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso", "mysql"} {
		t.Run(backend, func(t *testing.T) {
			db := Bbs_board_creation_audit_setup(t, backend, true)
			var before_count int
			if err := db.QueryRow("select count(*) from bbs_set").Scan(&before_count); err != nil {
				t.Fatal(err)
			}
			for _, response := range []map[string]any{
				route.Api_bbs_make(Identifier_config("unauthorized_board_creator"), "board_denied_make", "comment"),
				route.Api_bbs_list_post(Identifier_config("unauthorized_board_creator"), "board_denied_list_post", "thread"),
			} {
				if response["response"] != "require auth" {
					t.Fatalf("unauthorized board creation response=%v, want require auth", response)
				}
			}
			var after_count int
			if err := db.QueryRow("select count(*) from bbs_set").Scan(&after_count); err != nil || after_count != before_count {
				t.Fatalf("permission-denied calls changed bbs_set rows: before=%d after=%d err=%v", before_count, after_count, err)
			}

			transaction_error := errors.New("board allocation rollback probe")
			if err := route.Bbs_board_id_allocation_transaction(db, func(tx *sql.Tx) error {
				if _, err := tx.Exec("insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_name', '', '777', 'board_rollback_name')"); err != nil {
					return err
				}
				if _, err := tx.Exec("insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_type', '', '777', 'comment')"); err != nil {
					return err
				}
				return transaction_error
			}); !errors.Is(err, transaction_error) {
				t.Fatalf("allocation callback error=%v, want original error", err)
			}
			Bbs_assert_no_board_records(t, db, "board_rollback_name", "777")

			panic_recovered := false
			func() {
				defer func() { panic_recovered = recover() != nil }()
				_ = route.Bbs_board_id_allocation_transaction(db, func(tx *sql.Tx) error {
					if _, err := tx.Exec("insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_name', '', '778', 'board_panic_name')"); err != nil {
						return err
					}
					if _, err := tx.Exec("insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_type', '', '778', 'comment')"); err != nil {
						return err
					}
					panic("board allocation panic cleanup probe")
				})
			}()
			if !panic_recovered {
				t.Fatal("board allocation callback panic was not propagated")
			}
			Bbs_assert_no_board_records(t, db, "board_panic_name", "778")

			Bbs_assert_board_creation_response(t, route.Api_bbs_make(Identifier_config("admin_ip"), "board_after_rollback", "comment"))
			Bbs_assert_board_id_and_type(t, db, "board_after_rollback", "1", "comment")
		})
	}
}

func Test_Bbs_board_creation_concurrent_api_paths(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso", "mysql"} {
		t.Run(backend, func(t *testing.T) {
			db := Bbs_board_creation_audit_setup(t, backend, true)
			for _, actor_mode := range []string{"distinct", "same"} {
				for _, api_path := range []string{"make", "list_post", "mixed"} {
					for _, actor_count := range []int{2, 8, 16, 32} {
						for repetition := 0; repetition < 2; repetition++ {
							group := fmt.Sprintf("board_%s_%s_%d_%d", backend, actor_mode, actor_count, repetition)
							results := Bbs_run_concurrent_board_creation(db, api_path, actor_count, actor_mode == "same", group)
							Bbs_assert_concurrent_board_creation(t, db, group, results)
						}
					}
				}
			}
		})
	}
}

func Test_Bbs_board_id_allocation_mariadb_lock_lifecycle_and_scope(t *testing.T) {
	if os.Getenv("MYSQL_AUDIT_DSN") == "" {
		t.Skip("MYSQL_AUDIT_DSN is not configured")
	}
	db, _ := mysql_audit_setup(t)
	var database_name string
	if err := db.QueryRow("select database()").Scan(&database_name); err != nil {
		t.Fatal(err)
	}
	allocation_lock_name := route.Bbs_identifier_lock_name(database_name, "board-id-allocation")
	board_zero_lock_name := route.Bbs_identifier_lock_name(database_name, "board:0")
	if allocation_lock_name == board_zero_lock_name {
		t.Fatal("board allocation scope unexpectedly aliases reserved board 0")
	}
	observer_connection, err := db.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer observer_connection.Close()

	callback_error := errors.New("board allocation callback error")
	if err := route.Bbs_board_id_allocation_transaction(db, func(*sql.Tx) error { return callback_error }); !errors.Is(err, callback_error) {
		t.Fatalf("callback error=%v, want original error", err)
	}
	Assert_bbs_mariadb_lock_acquirable(t, observer_connection, allocation_lock_name)

	panic_recovered := false
	func() {
		defer func() { panic_recovered = recover() != nil }()
		_ = route.Bbs_board_id_allocation_transaction(db, func(*sql.Tx) error { panic("board lock panic cleanup") })
	}()
	if !panic_recovered {
		t.Fatal("board allocation panic was not propagated")
	}
	Assert_bbs_mariadb_lock_acquirable(t, observer_connection, allocation_lock_name)

	// A post on reserved board 0 uses a different lock and must be able to enter
	// while the dedicated namespace-allocation lock is held.
	allocation_entered := make(chan struct{})
	allocation_release := make(chan struct{})
	allocation_done := make(chan error, 1)
	go func() {
		allocation_done <- route.Bbs_board_id_allocation_transaction(db, func(*sql.Tx) error {
			close(allocation_entered)
			<-allocation_release
			return nil
		})
	}()
	<-allocation_entered
	board_zero_entered := make(chan struct{})
	board_zero_done := make(chan error, 1)
	go func() {
		board_zero_done <- route.Bbs_identifier_transaction(db, "0", "", true, func(*sql.Tx) error {
			close(board_zero_entered)
			return nil
		})
	}()
	select {
	case <-board_zero_entered:
	case err := <-board_zero_done:
		close(allocation_release)
		<-allocation_done
		t.Fatalf("board 0 transaction while allocation lock held: %v", err)
	case <-time.After(5 * time.Second):
		close(allocation_release)
		<-allocation_done
		<-board_zero_done
		t.Fatal("board-ID allocation scope serialized ordinary board 0 writes")
	}
	close(allocation_release)
	if err := <-allocation_done; err != nil {
		t.Fatalf("board allocation transaction: %v", err)
	}
	if err := <-board_zero_done; err != nil {
		t.Fatalf("board 0 transaction: %v", err)
	}
	Assert_bbs_mariadb_lock_acquirable(t, observer_connection, allocation_lock_name)
}

func Test_Bbs_board_id_allocation_mariadb_lock_timeout(t *testing.T) {
	if os.Getenv("MYSQL_AUDIT_DSN") == "" {
		t.Skip("MYSQL_AUDIT_DSN is not configured")
	}
	db, _ := mysql_audit_setup(t)
	var database_name string
	if err := db.QueryRow("select database()").Scan(&database_name); err != nil {
		t.Fatal(err)
	}
	lock_name := route.Bbs_identifier_lock_name(database_name, "board-id-allocation")
	lock_owner, err := db.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer lock_owner.Close()
	var acquired sql.NullInt64
	if err := lock_owner.QueryRowContext(context.Background(), "select get_lock(?, 0)", lock_name).Scan(&acquired); err != nil || !acquired.Valid || acquired.Int64 != 1 {
		t.Fatalf("reserve board allocation lock: acquired=%v err=%v", acquired, err)
	}

	callback_called := atomic.Bool{}
	transaction_done := make(chan error, 1)
	started_at := time.Now()
	go func() {
		transaction_done <- route.Bbs_board_id_allocation_transaction(db, func(*sql.Tx) error {
			callback_called.Store(true)
			return nil
		})
	}()
	select {
	case err := <-transaction_done:
		elapsed := time.Since(started_at)
		if err == nil || !strings.Contains(err.Error(), "timed out") || callback_called.Load() {
			var released sql.NullInt64
			_ = lock_owner.QueryRowContext(context.Background(), "select release_lock(?)", lock_name).Scan(&released)
			t.Fatalf("lock timeout result err=%v callback_called=%v elapsed=%s", err, callback_called.Load(), elapsed)
		}
		if elapsed > 35*time.Second {
			t.Fatalf("lock acquisition timeout was not bounded: %s", elapsed)
		}
	case <-time.After(36 * time.Second):
		var released sql.NullInt64
		_ = lock_owner.QueryRowContext(context.Background(), "select release_lock(?)", lock_name).Scan(&released)
		err := <-transaction_done
		t.Fatalf("lock acquisition exceeded 36 seconds: %v", err)
	}
	var released sql.NullInt64
	if err := lock_owner.QueryRowContext(context.Background(), "select release_lock(?)", lock_name).Scan(&released); err != nil || !released.Valid || released.Int64 != 1 {
		t.Fatalf("release test-owned board allocation lock: result=%v err=%v", released, err)
	}
	if err := route.Bbs_board_id_allocation_transaction(db, func(*sql.Tx) error { return nil }); err != nil {
		t.Fatalf("allocation after timed-out lock: %v", err)
	}
}

func Bbs_board_creation_audit_setup(t *testing.T, backend string, seed_reserved bool) *sql.DB {
	t.Helper()
	var db *sql.DB
	if backend == "mysql" {
		if os.Getenv("MYSQL_AUDIT_DSN") == "" {
			t.Skip("MYSQL_AUDIT_DSN is not configured")
		}
		db, _ = mysql_audit_setup(t)
	} else {
		db = audit_setup(t, backend)
		audit_seed_auth(db, t)
	}

	var version_count int
	if err := db.QueryRow("select count(*) from other where name='ver'").Scan(&version_count); err != nil {
		t.Fatalf("query initialized other.ver anchor: %v", err)
	}
	if version_count == 0 {
		if _, err := db.Exec("insert into other (name, data, coverage) values ('ver', 'test', '')"); err != nil {
			t.Fatalf("seed existing other.ver lock anchor: %v", err)
		}
	}
	if seed_reserved {
		for _, board := range [][3]string{{"0", "document_comment", "comment"}, {"-1", "thread", "thread"}, {"-2", "report_bbs", "comment"}} {
			for _, setting := range [][2]string{{"bbs_name", board[1]}, {"bbs_type", board[2]}} {
				if _, err := db.Exec("insert into bbs_set (set_name, set_code, set_id, set_data) values (?, '', ?, ?)", setting[0], board[0], setting[1]); err != nil {
					t.Fatalf("seed reserved BBS board %q/%s: %v", board[0], setting[0], err)
				}
			}
		}
	}
	return db
}

func Bbs_run_concurrent_board_creation(db *sql.DB, api_path string, actor_count int, same_actor bool, group string) []bbs_board_creation_result {
	start := make(chan struct{})
	results := make([]bbs_board_creation_result, actor_count)
	actor_ids := make([]string, actor_count)
	for index := range actor_ids {
		if same_actor {
			actor_ids[index] = "admin_ip"
		} else {
			actor_ids[index] = fmt.Sprintf("board_actor_%s_%02d", group, index)
			if _, err := db.Exec("insert into user_set (name, id, data) values ('acl', ?, 'owner')", actor_ids[index]); err != nil {
				panic(err)
			}
		}
	}

	var wait_group sync.WaitGroup
	for index := 0; index < actor_count; index++ {
		wait_group.Add(1)
		go func(index int) {
			defer wait_group.Done()
			<-start
			name := fmt.Sprintf("%s_board_%02d", group, index)
			board_type := "comment"
			if index%2 == 1 {
				board_type = "thread"
			}
			path := api_path
			if path == "mixed" {
				if index%2 == 0 {
					path = "make"
				} else {
					path = "list_post"
				}
			}
			results[index] = bbs_board_creation_result{name: name, board_type: board_type}
			defer func() {
				if recovered := recover(); recovered != nil {
					results[index].panic_value = fmt.Sprint(recovered)
				}
			}()
			config := Identifier_config(actor_ids[index])
			var response map[string]any
			if path == "make" {
				response = route.Api_bbs_make(config, name, board_type)
			} else {
				response = route.Api_bbs_list_post(config, name, board_type)
			}
			results[index].response, _ = response["response"].(string)
		}(index)
	}
	close(start)
	wait_group.Wait()
	if !same_actor {
		for _, actor_id := range actor_ids {
			_, _ = db.Exec("delete from user_set where id = ? and name = 'acl'", actor_id)
		}
	}
	return results
}

func Bbs_assert_concurrent_board_creation(t *testing.T, db *sql.DB, group string, results []bbs_board_creation_result) {
	t.Helper()
	id_to_names := make(map[string][]string)
	for _, result := range results {
		if result.panic_value != "" {
			t.Fatalf("concurrent board creation panicked for %q: %s", result.name, result.panic_value)
		}
		if result.response != "ok" {
			t.Fatalf("concurrent board creation for %q returned %q", result.name, result.response)
		}
		var board_id string
		if err := db.QueryRow("select set_id from bbs_set where set_name='bbs_name' and set_code='' and set_data=?", result.name).Scan(&board_id); err != nil {
			t.Fatalf("read created board %q ID: %v", result.name, err)
		}
		id_to_names[board_id] = append(id_to_names[board_id], result.name)

		var name_rows, type_rows int
		var actual_type string
		if err := db.QueryRow("select count(*) from bbs_set where set_name='bbs_name' and set_id=? and set_code=''", board_id).Scan(&name_rows); err != nil {
			t.Fatalf("count names for board ID %q: %v", board_id, err)
		}
		if err := db.QueryRow("select count(*), coalesce(max(set_data), '') from bbs_set where set_name='bbs_type' and set_id=? and set_code=''", board_id).Scan(&type_rows, &actual_type); err != nil {
			t.Fatalf("count types for board ID %q: %v", board_id, err)
		}
		var setting_rows int
		if err := db.QueryRow("select count(*) from bbs_set where set_id=? and set_code='' and set_name in ('bbs_name','bbs_type')", board_id).Scan(&setting_rows); err != nil {
			t.Fatalf("count EAV settings for board ID %q: %v", board_id, err)
		}
		if name_rows != 1 || type_rows != 1 || actual_type != result.board_type || setting_rows != 2 {
			t.Fatalf("created board %q/%s has name_rows=%d type_rows=%d type=%q setting_rows=%d", result.name, board_id, name_rows, type_rows, actual_type, setting_rows)
		}
	}
	if len(id_to_names) != len(results) {
		for board_id, names := range id_to_names {
			if len(names) > 1 {
				t.Fatalf("group %s assigned duplicate logical board ID %q to names %v", group, board_id, names)
			}
		}
		t.Fatalf("group %s returned %d unique IDs for %d successful requests", group, len(id_to_names), len(results))
	}
}

func Bbs_assert_board_creation_response(t *testing.T, response map[string]any) {
	t.Helper()
	if response["response"] != "ok" {
		t.Fatalf("board creation response=%v, want ok", response)
	}
}

func Bbs_assert_board_id_and_type(t *testing.T, db *sql.DB, board_name string, expected_id string, expected_type string) {
	t.Helper()
	var board_id string
	if err := db.QueryRow("select set_id from bbs_set where set_name='bbs_name' and set_code='' and set_data=?", board_name).Scan(&board_id); err != nil || board_id != expected_id {
		t.Fatalf("board %q ID=%q err=%v, want %q", board_name, board_id, err, expected_id)
	}
	var board_type string
	if err := db.QueryRow("select set_data from bbs_set where set_name='bbs_type' and set_code='' and set_id=?", board_id).Scan(&board_type); err != nil || board_type != expected_type {
		t.Fatalf("board %q/%s type=%q err=%v, want %q", board_name, board_id, board_type, err, expected_type)
	}
}

func Bbs_assert_no_board_records(t *testing.T, db *sql.DB, board_name string, board_id string) {
	t.Helper()
	var name_rows int
	if err := db.QueryRow("select count(*) from bbs_set where set_name='bbs_name' and set_data=?", board_name).Scan(&name_rows); err != nil || name_rows != 0 {
		t.Fatalf("failed board transaction left %d name records for %q, err=%v", name_rows, board_name, err)
	}
	var type_rows int
	if err := db.QueryRow("select count(*) from bbs_set where set_name='bbs_type' and set_id=?", board_id).Scan(&type_rows); err != nil || type_rows != 0 {
		t.Fatalf("failed board transaction left %d type records for ID %s, err=%v", type_rows, board_id, err)
	}
}
