package test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	mysql_driver "github.com/go-sql-driver/mysql"

	"opennamu/route"
	"opennamu/route/tool"
)

type bbs_identifier_api_result struct {
	response string
	data     string
	code     string
	panic    bool
}

func Test_Bbs_identifier_lock_name_is_scoped_and_bounded(t *testing.T) {
	first := route.Bbs_identifier_lock_name("database-a", "board:7")
	if len(first) > 64 {
		t.Fatalf("advisory lock name is %d bytes, max is 64", len(first))
	}
	if first != route.Bbs_identifier_lock_name("database-a", "board:7") {
		t.Fatal("same database/scope generated different lock names")
	}
	if first == route.Bbs_identifier_lock_name("database-b", "board:7") {
		t.Fatal("different databases shared an advisory lock name")
	}
	if first == route.Bbs_identifier_lock_name("database-a", "board:8") {
		t.Fatal("different boards shared an advisory lock name")
	}
	if first == route.Bbs_identifier_lock_name("database-a", "post:7:1") {
		t.Fatal("board and post scopes shared an advisory lock name")
	}
}

func Test_Bbs_identifier_concurrent_post_and_thread_creation(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso", "mysql"} {
		t.Run(backend, func(t *testing.T) {
			db := Bbs_identifier_audit_setup(t, backend)
			defer tool.DB_close(db)
			Seed_bbs_identifier_board(t, db, "-1", "thread", "thread", 60)

			for repetition := 0; repetition < 3; repetition++ {
				board := strconv.Itoa(930000 + repetition)
				Seed_bbs_identifier_board(t, db, board, "identifier board", "comment", 60)
				results := Run_bbs_identifier_api_writes(16, false, fmt.Sprintf("post_%s_%d", backend, repetition), "data", func(config tool.Config, index int) map[string]any {
					return route.Api_bbs_w_edit_post(config, board, "", "", fmt.Sprintf("post title %d", index), "post body", "", "", "")
				})
				Assert_bbs_identifier_api_results(t, db, results, board, []string{"title", "data", "user_id"}, 16, 0)
			}

			for repetition := 0; repetition < 3; repetition++ {
				results := Run_bbs_identifier_api_writes(16, false, fmt.Sprintf("thread_%s_%d", backend, repetition), "topic_num", func(config tool.Config, index int) map[string]any {
					return route.Api_thread_bbs_post(config, "0", "identifier-document", "thread body", "identifier-document", fmt.Sprintf("thread title %d", index))
				})
				Assert_bbs_identifier_api_results(t, db, results, "-1", []string{"title", "data", "user_id"}, 16, 0)
			}

			Seed_bbs_identifier_board(t, db, "0", "document comments", "comment", 0)
			document_comment_results := Run_bbs_identifier_api_writes(16, false, "document_comment_"+backend, "data", func(_ tool.Config, index int) map[string]any {
				return route.Api_bbs_w_post(Identifier_config("admin_ip"), "0", fmt.Sprintf("document title %d", index), "document comment")
			})
			Assert_bbs_identifier_api_results(t, db, document_comment_results, "0", []string{"title", "data", "user_id"}, 16, 0)

			document_name := "identifier_document_comment_" + backend
			first_document_comment_results := Run_bbs_identifier_api_writes(16, false, "first_document_comment_"+backend, "data", func(config tool.Config, _ int) map[string]any {
				return route.Api_w_comment(config, document_name)
			})
			var document_comment_code string
			for _, result := range first_document_comment_results {
				if result.panic || result.response != "ok" || result.code == "" {
					t.Fatalf("document comment API result=%+v, want successful idempotent response", result)
				}
				if document_comment_code == "" {
					document_comment_code = result.code
				} else if result.code != document_comment_code {
					t.Fatalf("concurrent document-comment mappings differ: %q vs %q", document_comment_code, result.code)
				}
			}
			for _, field := range []string{"title", "data", "user_id"} {
				var row_count int
				if err := db.QueryRow("select count(*) from bbs_data where set_name=? and set_id='0' and set_code=?", field, document_comment_code).Scan(&row_count); err != nil || row_count != 1 {
					t.Fatalf("document comment field %s rows=%d err=%v, want exactly one row", field, row_count, err)
				}
			}
			var mapped_document_comment_code string
			if err := db.QueryRow("select set_data from data_set where doc_name=? and set_name='document_comment_code' limit 1", document_name).Scan(&mapped_document_comment_code); err != nil || mapped_document_comment_code != document_comment_code {
				t.Fatalf("document comment mapping=%q err=%v, want code %q", mapped_document_comment_code, err, document_comment_code)
			}
			var document_comment_mapping_count int
			if err := db.QueryRow("select count(*) from data_set where doc_name=? and set_name='document_comment_code'", document_name).Scan(&document_comment_mapping_count); err != nil || document_comment_mapping_count != 1 {
				t.Fatalf("document comment mapping rows=%d err=%v, want exactly one", document_comment_mapping_count, err)
			}

			Seed_bbs_identifier_board(t, db, "-2", "report board", "comment", 0)
			report_results := Run_bbs_identifier_api_writes(16, false, "report_"+backend, "data", func(config tool.Config, index int) map[string]any {
				code, err := route.Report_create(db, config.IP, fmt.Sprintf("report target %d", index), "/bbs/w/identifier", "report", "report reason")
				if err != nil {
					return map[string]any{"response": "error", "data": err.Error()}
				}
				return map[string]any{"response": "ok", "data": code}
			})
			Assert_bbs_identifier_api_results(t, db, report_results, "-2", []string{"title", "data", "user_id", "report_target"}, 16, 0)

			var deleted_code string
			if err := db.QueryRow("select set_code from bbs_data where set_name='title' and set_id='930000' order by set_code + 0 desc limit 1").Scan(&deleted_code); err != nil {
				t.Fatalf("query highest post code before deletion: %v", err)
			}
			deleted := route.Api_bbs_w_delete(Identifier_config("admin_ip"), "930000", deleted_code)
			if deleted["response"] != "ok" {
				t.Fatalf("post delete response=%v", deleted)
			}
			recreated := route.Api_bbs_w_edit_post(Identifier_config("post_recreated_after_delete"), "930000", "", "", "recreated post", "post body", "", "", "")
			if recreated["response"] != "ok" || recreated["data"] != deleted_code {
				t.Fatalf("post recreation after deletion = %v, want code %s", recreated, deleted_code)
			}
			for _, field := range []string{"title", "data", "user_id"} {
				var row_count int
				if err := db.QueryRow("select count(*) from bbs_data where set_name=? and set_id='930000' and set_code=?", field, deleted_code).Scan(&row_count); err != nil || row_count != 1 {
					t.Fatalf("post recreation field %s rows=%d err=%v, want one live row", field, row_count, err)
				}
			}

			board := "939900"
			Seed_bbs_identifier_board(t, db, board, "slow identifier board", "comment", 60)
			regular_same_actor := Run_bbs_identifier_api_writes(16, true, "same_regular_"+backend, "data", func(config tool.Config, index int) map[string]any {
				return route.Api_bbs_w_edit_post(config, board, "", "", fmt.Sprintf("same title %d", index), "post body", "", "", "")
			})
			Assert_bbs_identifier_api_results(t, db, regular_same_actor, board, []string{"title", "data", "user_id"}, 1, 15)

			thread_same_actor := Run_bbs_identifier_api_writes(16, true, "same_thread_"+backend, "topic_num", func(config tool.Config, index int) map[string]any {
				return route.Api_thread_bbs_post(config, "0", "identifier-document", "thread body", "identifier-document", fmt.Sprintf("same thread title %d", index))
			})
			Assert_bbs_identifier_api_results(t, db, thread_same_actor, "-1", []string{"title", "data", "user_id"}, 1, 15)
		})
	}
}

func Test_Bbs_identifier_concurrent_comment_creation(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso", "mysql"} {
		t.Run(backend, func(t *testing.T) {
			db := Bbs_identifier_audit_setup(t, backend)
			defer tool.DB_close(db)
			const board = "940001"
			Seed_bbs_identifier_board(t, db, board, "comment identifier board", "comment", 60)
			for repetition := 0; repetition < 3; repetition++ {
				post_code := strconv.Itoa(repetition + 1)
				Seed_bbs_identifier_post(t, db, board, post_code, fmt.Sprintf("post_%d_owner", repetition))
				results := Run_bbs_identifier_api_writes(16, false, fmt.Sprintf("comment_%s_%d", backend, repetition), "data", func(config tool.Config, index int) map[string]any {
					return route.Api_bbs_w_comment_post(config, board, post_code, "0", fmt.Sprintf("comment body %d", index))
				})
				Assert_bbs_identifier_api_results(t, db, results, board+"-"+post_code, []string{"comment", "comment_date", "comment_user_id"}, 16, 0)
			}

			post_scope_board := "940002"
			Seed_bbs_identifier_board(t, db, post_scope_board, "post lock board", "comment", 0)
			Seed_bbs_identifier_post(t, db, post_scope_board, "1", "post_scope_owner")
			post_scope_results := Run_bbs_identifier_api_writes(16, false, "post_scope_"+backend, "data", func(config tool.Config, index int) map[string]any {
				return route.Api_bbs_w_comment_post(config, post_scope_board, "1", "0", fmt.Sprintf("post-scope comment %d", index))
			})
			Assert_bbs_identifier_api_results(t, db, post_scope_results, post_scope_board+"-1", []string{"comment", "comment_date", "comment_user_id"}, 16, 0)

			legacy_board := "940003"
			Seed_bbs_identifier_board(t, db, legacy_board, "legacy lock board", "comment", 0)
			Seed_bbs_identifier_legacy_post(t, db, legacy_board, "1", "legacy_owner")
			legacy_results := Run_bbs_identifier_api_writes(16, false, "legacy_post_"+backend, "data", func(config tool.Config, index int) map[string]any {
				return route.Api_bbs_w_comment_post(config, legacy_board, "1", "0", fmt.Sprintf("legacy comment %d", index))
			})
			Assert_bbs_identifier_api_results(t, db, legacy_results, legacy_board+"-1", []string{"comment", "comment_date", "comment_user_id"}, 16, 0)
			var legacy_comment_count string
			if err := db.QueryRow("select set_data from bbs_data where set_name='comment_count' and set_id=? and set_code='1'", legacy_board).Scan(&legacy_comment_count); err != nil || legacy_comment_count != "16" {
				t.Fatalf("legacy comment_count after concurrent writes = %q, err=%v, want 16", legacy_comment_count, err)
			}

			Seed_bbs_identifier_board(t, db, "-1", "thread", "thread", 0)
			created_topic := route.Api_thread_bbs_post(Identifier_config("thread_topic_owner"), "0", "identifier-document", "topic body", "identifier-document", "comment identifier topic")
			if created_topic["response"] != "ok" {
				t.Fatalf("seed thread topic response=%v", created_topic)
			}
			topic_code, ok := created_topic["topic_num"].(string)
			if !ok || topic_code == "" {
				t.Fatalf("seed thread topic code=%v", created_topic["topic_num"])
			}
			for repetition := 0; repetition < 3; repetition++ {
				results := Run_bbs_identifier_api_writes(16, false, fmt.Sprintf("thread_comment_%s_%d", backend, repetition), "comment_num", func(config tool.Config, _ int) map[string]any {
					return route.Api_thread_bbs_post(config, topic_code, "", "thread comment body", "", "")
				})
				Assert_bbs_identifier_api_results(t, db, results, "-1-"+topic_code, []string{"comment", "comment_date", "comment_user_id"}, 16, 0)
			}
			deleted_comments := route.Api_bbs_w_comment_all_delete(Identifier_config("admin_ip"), "-1", topic_code)
			if deleted_comments["response"] != "ok" {
				t.Fatalf("thread comment deletion response=%v", deleted_comments)
			}
			recreated_comment := route.Api_thread_bbs_post(Identifier_config("comment_after_delete"), topic_code, "", "comment after deletion", "", "")
			if recreated_comment["response"] != "ok" || recreated_comment["comment_num"] != "1" {
				t.Fatalf("thread comment after deletion = %v, want code 1", recreated_comment)
			}
			for _, field := range []string{"comment", "comment_date", "comment_user_id"} {
				var row_count int
				if err := db.QueryRow("select count(*) from bbs_data where set_name=? and set_id=? and set_code='1'", field, "-1-"+topic_code).Scan(&row_count); err != nil || row_count != 1 {
					t.Fatalf("recreated thread comment field %s rows=%d err=%v, want one row", field, row_count, err)
				}
			}
		})
	}
}

func Test_Bbs_identifier_transaction_mariadb_lock_lifecycle(t *testing.T) {
	if os.Getenv("MYSQL_AUDIT_DSN") == "" {
		t.Skip("MYSQL_AUDIT_DSN is not configured")
	}
	db, _ := mysql_audit_setup(t)
	defer tool.DB_close(db)
	var database_name string
	if err := db.QueryRow("select database()").Scan(&database_name); err != nil {
		t.Fatal(err)
	}
	const board = "lifecycle_probe"
	lock_name := route.Bbs_identifier_lock_name(database_name, "board:"+board)
	lock_connection, err := db.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer lock_connection.Close()

	entered := make(chan struct{})
	release := make(chan struct{})
	transaction_done := make(chan error, 1)
	go func() {
		transaction_done <- route.Bbs_identifier_transaction(db, board, "", true, func(*sql.Tx) error {
			close(entered)
			<-release
			return nil
		})
	}()
	<-entered
	var lock_result sql.NullInt64
	if err := lock_connection.QueryRowContext(context.Background(), "select get_lock(?, 0)", lock_name).Scan(&lock_result); err != nil {
		close(release)
		t.Fatal(err)
	}
	if !lock_result.Valid || lock_result.Int64 != 0 {
		close(release)
		t.Fatalf("MariaDB advisory lock was not held through transaction callback: %v", lock_result)
	}
	close(release)
	if err := <-transaction_done; err != nil {
		t.Fatalf("identifier transaction: %v", err)
	}
	Assert_bbs_mariadb_lock_acquirable(t, lock_connection, lock_name)

	callback_error := errors.New("callback error probe")
	if err := route.Bbs_identifier_transaction(db, board, "", true, func(*sql.Tx) error { return callback_error }); !errors.Is(err, callback_error) {
		t.Fatalf("callback error = %v, want original callback error", err)
	}
	Assert_bbs_mariadb_lock_acquirable(t, lock_connection, lock_name)

	for _, deadlock_number := range []uint16{1213, 1205} {
		attempts := 0
		err := route.Bbs_identifier_transaction(db, board, "", true, func(*sql.Tx) error {
			attempts++
			if attempts == 1 {
				return &mysql_driver.MySQLError{Number: deadlock_number, Message: "synthetic retry probe"}
			}
			return nil
		})
		if err != nil || attempts != 2 {
			t.Fatalf("MariaDB retry for error %d: attempts=%d err=%v, want two attempts and success", deadlock_number, attempts, err)
		}
		Assert_bbs_mariadb_lock_acquirable(t, lock_connection, lock_name)
	}

	panic_attempts := 0
	if err := route.Bbs_identifier_transaction(db, board, "", true, func(*sql.Tx) error {
		panic_attempts++
		if panic_attempts == 1 {
			panic(&mysql_driver.MySQLError{Number: 1213, Message: "synthetic panic retry probe"})
		}
		return nil
	}); err != nil || panic_attempts != 2 {
		t.Fatalf("MariaDB retry for a driver-error panic: attempts=%d err=%v, want two attempts and success", panic_attempts, err)
	}
	Assert_bbs_mariadb_lock_acquirable(t, lock_connection, lock_name)

	panic_recovered := false
	func() {
		defer func() { panic_recovered = recover() != nil }()
		_ = route.Bbs_identifier_transaction(db, board, "", true, func(*sql.Tx) error { panic("lock cleanup probe") })
	}()
	if !panic_recovered {
		t.Fatal("identifier transaction callback panic was not propagated")
	}
	Assert_bbs_mariadb_lock_acquirable(t, lock_connection, lock_name)

	first_entered := make(chan struct{})
	first_release := make(chan struct{})
	first_done := make(chan error, 1)
	go func() {
		first_done <- route.Bbs_identifier_transaction(db, "independent_board_a", "", true, func(*sql.Tx) error {
			close(first_entered)
			<-first_release
			return nil
		})
	}()
	<-first_entered
	second_entered := make(chan struct{})
	second_release := make(chan struct{})
	second_done := make(chan error, 1)
	go func() {
		second_done <- route.Bbs_identifier_transaction(db, "independent_board_b", "", true, func(*sql.Tx) error {
			close(second_entered)
			<-second_release
			return nil
		})
	}()
	select {
	case <-second_entered:
		close(second_release)
	case err := <-second_done:
		close(second_release)
		close(first_release)
		<-first_done
		t.Fatalf("independent board lock transaction failed: %v", err)
	case <-time.After(5 * time.Second):
		close(second_release)
		close(first_release)
		<-first_done
		<-second_done
		t.Fatal("distinct MariaDB board scopes serialized behind one global lock")
	}
	close(first_release)
	if err := <-first_done; err != nil {
		t.Fatalf("first independent board transaction: %v", err)
	}
	if err := <-second_done; err != nil {
		t.Fatalf("second independent board transaction: %v", err)
	}
}

func Bbs_identifier_audit_setup(t *testing.T, backend string) *sql.DB {
	t.Helper()
	if backend == "mysql" {
		db, _ := mysql_audit_setup(t)
		return db
	}
	db := audit_setup(t, backend)
	audit_seed_auth(db, t)
	return db
}

func Seed_bbs_identifier_board(t *testing.T, db *sql.DB, set_id string, name string, board_type string, slow_seconds int) {
	t.Helper()
	for _, row := range [][4]string{
		{"bbs_name", "", set_id, name},
		{"bbs_type", "", set_id, board_type},
		{"bbs_slow_mode_seconds", "", set_id, strconv.Itoa(slow_seconds)},
	} {
		if _, err := db.Exec("insert into bbs_set (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)", row[0], row[1], row[2], row[3]); err != nil {
			t.Fatalf("seed BBS identifier board %q: %v", set_id, err)
		}
	}
}

func Seed_bbs_identifier_post(t *testing.T, db *sql.DB, board string, post_code string, owner string) {
	t.Helper()
	for _, row := range [][4]string{
		{"title", post_code, board, "identifier test post"},
		{"data", post_code, board, "identifier test body"},
		{"date", post_code, board, tool.Get_time()},
		{"last_activity", post_code, board, tool.Get_time()},
		{"user_id", post_code, board, owner},
		{"comment_count", post_code, board, "0"},
	} {
		if _, err := db.Exec("insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)", row[0], row[1], row[2], row[3]); err != nil {
			t.Fatalf("seed BBS identifier post %q/%q: %v", board, post_code, err)
		}
	}
}

func Seed_bbs_identifier_legacy_post(t *testing.T, db *sql.DB, board string, post_code string, owner string) {
	t.Helper()
	for _, row := range [][4]string{
		{"title", post_code, board, "legacy identifier test post"},
		{"data", post_code, board, "legacy identifier test body"},
		{"date", post_code, board, tool.Get_time()},
		{"last_activity", post_code, board, tool.Get_time()},
		{"user_id", post_code, board, owner},
	} {
		if _, err := db.Exec("insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)", row[0], row[1], row[2], row[3]); err != nil {
			t.Fatalf("seed legacy BBS identifier post %q/%q: %v", board, post_code, err)
		}
	}
}

func Run_bbs_identifier_api_writes(actor_count int, same_actor bool, group string, code_field string, write func(tool.Config, int) map[string]any) []bbs_identifier_api_result {
	start := make(chan struct{})
	results := make([]bbs_identifier_api_result, actor_count)
	var wait_group sync.WaitGroup
	for index := 0; index < actor_count; index++ {
		wait_group.Add(1)
		go func(index int) {
			defer wait_group.Done()
			<-start
			actor := fmt.Sprintf("identifier_%s_%d", group, index)
			if same_actor {
				actor = "identifier_same_" + group
			}
			defer func() {
				if recover() != nil {
					results[index].panic = true
				}
			}()
			response := write(Identifier_config(actor), index)
			results[index].response, _ = response["response"].(string)
			results[index].data, _ = response["data"].(string)
			results[index].code, _ = response[code_field].(string)
		}(index)
	}
	close(start)
	wait_group.Wait()
	return results
}

func Assert_bbs_identifier_api_results(t *testing.T, db *sql.DB, results []bbs_identifier_api_result, set_id string, fields []string, expected_ok int, expected_slow_denied int) {
	t.Helper()
	ok_count, slow_denied_count := 0, 0
	code_counts := make(map[string]int)
	for _, result := range results {
		if result.panic {
			t.Fatal("BBS identifier API panicked")
		}
		switch result.response {
		case "ok":
			ok_count++
			if result.code == "" {
				t.Fatalf("successful API response has no allocated identifier: %+v", result)
			}
			code_counts[result.code]++
		case "error":
			if result.data == "slow_mode" {
				slow_denied_count++
			} else {
				t.Fatalf("unexpected BBS identifier API error: %+v", result)
			}
		default:
			t.Fatalf("unexpected BBS identifier API response: %+v", result)
		}
	}
	if ok_count != expected_ok || slow_denied_count != expected_slow_denied {
		t.Fatalf("API outcome counts: ok=%d slow_denied=%d, want ok=%d slow_denied=%d", ok_count, slow_denied_count, expected_ok, expected_slow_denied)
	}
	for code, count := range code_counts {
		if count != 1 {
			t.Fatalf("identifier %q was returned %d times", code, count)
		}
		for _, field := range fields {
			var row_count int
			if err := db.QueryRow("select count(*) from bbs_data where set_name = ? and set_id = ? and set_code = ?", field, set_id, code).Scan(&row_count); err != nil {
				t.Fatalf("count %s record for %s/%s: %v", field, set_id, code, err)
			}
			if row_count != 1 {
				t.Fatalf("identifier %s/%s has %d %s rows, want exactly one", set_id, code, row_count, field)
			}
		}
	}
}

func Identifier_config(actor string) tool.Config {
	return tool.Config{IP: actor, Session: &emptySession{}}
}

func Assert_bbs_mariadb_lock_acquirable(t *testing.T, conn *sql.Conn, lock_name string) {
	t.Helper()
	var acquired sql.NullInt64
	if err := conn.QueryRowContext(context.Background(), "select get_lock(?, 0)", lock_name).Scan(&acquired); err != nil {
		t.Fatalf("acquire released MariaDB lock: %v", err)
	}
	if !acquired.Valid || acquired.Int64 != 1 {
		t.Fatalf("MariaDB identifier lock was not released: %v", acquired)
	}
	var released sql.NullInt64
	if err := conn.QueryRowContext(context.Background(), "select release_lock(?)", lock_name).Scan(&released); err != nil {
		t.Fatalf("release test MariaDB lock: %v", err)
	}
	if !released.Valid || released.Int64 != 1 {
		t.Fatalf("release test MariaDB lock returned %v", released)
	}
}
