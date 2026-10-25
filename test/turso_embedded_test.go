package test

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"opennamu/route/tool"
)

// Turso_setup_db boots the full openNAMU schema on the turso embedded driver
// inside a temp directory and returns a connected pool. It exercises the real
// boot path (DB_boot -> DB_connect_init -> DB_make), including the mvcc pragma.
func Turso_setup_db(t *testing.T) *sql.DB {
	t.Helper()

	db_name := filepath.Join(t.TempDir(), "turso_test")
	os.Setenv("NAMU_DB", db_name)
	os.Setenv("NAMU_DB_TYPE", "turso")
	defer os.Unsetenv("NAMU_DB")
	defer os.Unsetenv("NAMU_DB_TYPE")

	tool.DB_init()

	return tool.DB_connect()
}

// Test_turso_platform_supported asserts the platform gate reports supported on
// the platform the tests run on (CI/dev run on linux; the gate only rejects
// windows/arm64, which cannot be executed in this environment).
func Test_turso_platform_supported(t *testing.T) {
	if !tool.Turso_platform_supported() {
		t.Fatalf("turso should be supported on %s/%s", runtime.GOOS, runtime.GOARCH)
	}
}

// Test_turso_embedded_basic verifies the turso embedded driver end to end:
// boot with mvcc journal mode, CRUD through the tool helpers, transaction
// commit and rollback, and persistence across a reconnect.
func Test_turso_embedded_basic(t *testing.T) {
	db := Turso_setup_db(t)
	defer tool.DB_close(db)

	mode := ""
	if err := db.QueryRow("pragma journal_mode").Scan(&mode); err != nil {
		t.Fatalf("read journal_mode: %v", err)
	}
	if mode != "mvcc" {
		t.Fatalf("journal_mode = %q, want mvcc", mode)
	}

	tool.Exec_DB(db, `insert into other (name, data, coverage) values ('turso_test', 'value1', '')`)

	found := ""
	if !tool.QueryRow_DB(db, `select data from other where name = 'turso_test'`, []any{&found}) {
		t.Fatalf("inserted row not found")
	}
	if found != "value1" {
		t.Fatalf("data = %q, want value1", found)
	}

	rows := tool.Query_DB(db, `select name from other where name like 'turso%'`)
	names := []string{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			rows.Close()
			t.Fatalf("scan: %v", err)
		}
		names = append(names, name)
	}
	rows.Close()
	if len(names) != 1 || names[0] != "turso_test" {
		t.Fatalf("names = %v, want [turso_test]", names)
	}

	// commit path (DB_transaction owns Begin/Commit)
	err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		_, e := tx.Exec(`insert into other (name, data, coverage) values ('turso_tx', 'committed', '')`)
		return e
	})
	if err != nil {
		t.Fatalf("commit tx: %v", err)
	}

	// rollback path
	err = tool.DB_transaction(db, func(tx *sql.Tx) error {
		if _, e := tx.Exec(`insert into other (name, data, coverage) values ('turso_tx', 'rolledback', '')`); e != nil {
			return e
		}
		return errors.New("force rollback")
	})
	if err == nil {
		t.Fatalf("rollback tx: expected error, got nil")
	}

	rolledback := ""
	if tool.QueryRow_DB(db, `select data from other where name = 'turso_tx' and data = 'rolledback'`, []any{&rolledback}) {
		t.Fatalf("rolledback row must not exist")
	}

	committed := ""
	if !tool.QueryRow_DB(db, `select data from other where name = 'turso_tx'`, []any{&committed}) {
		t.Fatalf("committed row not found")
	}
	if committed != "committed" {
		t.Fatalf("committed data = %q", committed)
	}

	// persistence across reconnect
	tool.DB_close(db)
	db2 := tool.DB_connect()
	defer tool.DB_close(db2)

	again := ""
	if !tool.QueryRow_DB(db2, `select data from other where name = 'turso_test'`, []any{&again}) {
		t.Fatalf("row lost after reconnect")
	}
	if again != "value1" {
		t.Fatalf("data after reconnect = %q", again)
	}
}

// Test_turso_embedded_concurrent_writes writes through two independent
// connection pools at the same time. MVCC journal mode must let concurrent
// writers proceed without busy errors, and every write must be durable.
func Test_turso_embedded_concurrent_writes(t *testing.T) {
	db := Turso_setup_db(t)
	defer tool.DB_close(db)

	tool.Exec_DB(db, `create table if not exists turso_conc_test (test text default '')`)

	db2 := tool.DB_connect()
	defer tool.DB_close(db2)

	const writers = 2
	const per_writer = 25

	err_ch := make(chan error, writers)
	var wg sync.WaitGroup
	for writer := 0; writer < writers; writer++ {
		wg.Add(1)
		go func(writer int) {
			defer wg.Done()

			conn := db
			if writer == 1 {
				conn = db2
			}

			for i := 0; i < per_writer; i++ {
				_, err := conn.Exec(
					`insert into turso_conc_test (test) values (?)`,
					fmt.Sprintf("writer%d_%d", writer, i),
				)
				if err != nil {
					err_ch <- err
					return
				}
			}
		}(writer)
	}
	wg.Wait()
	close(err_ch)

	for err := range err_ch {
		t.Fatalf("concurrent write failed: %v", err)
	}

	count := 0
	if !tool.QueryRow_DB(db, `select count(*) from turso_conc_test`, []any{&count}) {
		t.Fatalf("count query failed")
	}
	if count != writers*per_writer {
		t.Fatalf("count = %d, want %d", count, writers*per_writer)
	}
}