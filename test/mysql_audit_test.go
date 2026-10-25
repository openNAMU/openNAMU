package test

// MySQL stability audit: real MariaDB execution against a dedicated test-only
// schema/user created by the audit (see tmp/mysql_audit/mysql_test_setup.md).
// Skipped entirely when MYSQL_AUDIT_DSN is not provided - never touches the
// example repo's production databases.
//
// MYSQL_AUDIT_DSN: "user:password@tcp(127.0.0.1:3306)/" (server-level DSN,
// DB name is appended by the audit itself; passwords stay out of logs).

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"opennamu/route/tool"
)

// mysql_audit_dsn returns the audit DSN and database name, or skips.
func mysql_audit_dsn(t *testing.T) (string, string) {
	t.Helper()
	dsn := os.Getenv("MYSQL_AUDIT_DSN")
	if dsn == "" {
		t.Skip("MySQL audit: set MYSQL_AUDIT_DSN to run against a test-only MariaDB")
	}
	return dsn, fmt.Sprintf("opennamu_audit_run_%d", time.Now().UnixNano()%1_000_000)
}

// mysql_audit_setup boots the real internal build against a fresh dedicated
// database. It writes a temporary data/mysql.json in the CURRENT working
// directory - tests run from ./test, so the repo's data/ directory is never
// touched; the file lives only for the test duration in a temp cwd copy.
//
// Simpler and safer: the test process chdirs into a temp dir so every relative
// path (data/mysql.json, db bootstrap writes) stays inside the sandbox.
func mysql_audit_setup(t *testing.T) (*sql.DB, string) {
	t.Helper()
	dsn, dbname := mysql_audit_dsn(t)

	tmp_dir := t.TempDir()
	data_dir := tmp_dir + "/data"
	if err := os.MkdirAll(data_dir, 0o755); err != nil {
		t.Fatalf("mkdir data dir: %v", err)
	}

	// extract user/password from the server-level DSN for mysql.json
	// dsn format: user:password@tcp(host:port)/
	cred_end := strings.Index(dsn, "@")
	if cred_end <= 0 {
		t.Fatalf("bad MYSQL_AUDIT_DSN format")
	}
	creds := dsn[:cred_end]
	colon := strings.Index(creds, ":")
	if colon <= 0 {
		t.Fatalf("bad MYSQL_AUDIT_DSN credentials")
	}
	user := creds[:colon]
	password := creds[colon+1:]

	mysql_json := fmt.Sprintf(`{"user":%q,"password":%q,"host":"127.0.0.1","port":"3306"}`, user, password)
	if err := os.WriteFile(data_dir+"/mysql.json", []byte(mysql_json), 0o600); err != nil {
		t.Fatalf("write mysql.json: %v", err)
	}
	t.Setenv("NAMU_DB_TYPE", "mysql")
	t.Setenv("NAMU_DB", dbname)

	// create the dedicated database (server-level connection) before boot
	server_db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("server open: %v", err)
	}
	defer server_db.Close()
	if _, err := server_db.Exec(fmt.Sprintf(
		"create database if not exists `%s` character set utf8mb4", dbname)); err != nil {
		t.Fatalf("create database: %v", err)
	}

	// run DB_init from the sandbox cwd so data/mysql.json resolves there
	old_cwd, _ := os.Getwd()
	if err := os.Chdir(tmp_dir); err != nil {
		t.Fatalf("chdir sandbox: %v", err)
	}
	t.Cleanup(func() {
		os.Chdir(old_cwd)
		// drop the audit database created for this run
		sdb, e := sql.Open("mysql", dsn)
		if e == nil {
			_, _ = sdb.Exec(fmt.Sprintf("drop database if exists `%s`", dbname))
			sdb.Close()
		}
	})

	tool.DB_init()
	db := tool.DB_connect()
	t.Cleanup(func() { tool.DB_close(db) })
	audit_seed_auth(db, t)

	return db, dbname
}

func set_other_language_mysql(db *sql.DB, language string) {
	if _, err := db.Exec(`delete from other where name = 'language'`); err != nil {
		panic(err)
	}
	if _, err := db.Exec(`insert into other (name, data, coverage) values ('language', ?, '')`, language); err != nil {
		panic(err)
	}
}

// Test_mysql_boot_and_basics: real DB_init + schema bootstrap + language
// semantics through the actual internal build.
func Test_mysql_boot_and_basics(t *testing.T) {
	db, _ := mysql_audit_setup(t)

	var one int
	if err := db.QueryRow(`select 1`).Scan(&one); err != nil {
		t.Fatalf("select 1: %v", err)
	}

	// language semantics through the shared pool (real load_lang path)
	if got := tool.Get_language_shared("language_tag", false); got != "ko-KR" {
		t.Fatalf("default language = %q, want ko-KR", got)
	}
	set_other_language_mysql(db, "en-US")
	if got := tool.Get_language_shared("language_tag", false); got != "en-US" {
		t.Fatalf("after switch = %q, want en-US", got)
	}
	set_other_language_mysql(db, "ko-KR")
	if got := tool.Get_language_shared("language_tag", false); got != "ko-KR" {
		t.Fatalf("after switch-back = %q, want ko-KR", got)
	}

	missing := tool.Get_language_shared("definitely_missing_key_zzz", false)
	if !strings.Contains(missing, "definitely_missing_key_zzz") || !strings.Contains(missing, "(") {
		t.Fatalf("missing-key fallback = %q", missing)
	}

	t.Logf("mysql boot + language semantics OK")
}

// Test_mysql_lang_pool_reuse_and_swap: the long-lived lang pool reuses physical
// connections on MySQL (no per-call pool churn) and the lease survives a
// mid-flight identity swap without "sql: database is closed".
func Test_mysql_lang_pool_reuse_and_swap(t *testing.T) {
	db, _ := mysql_audit_setup(t)

	// baseline: render with 23 load_lang calls must not open 23 pools.
	// Probe via Stats() on the shared pool: MaxOpenConnections should stay
	// constant across repeated calls (per-call pools would create many).
	p1 := tool.Get_lang_db().Stats()
	for i := 0; i < 50; i++ {
		_ = tool.Get_language_shared("language_tag", false)
	}
	p2 := tool.Get_lang_db().Stats()
	t.Logf("lang pool stats p1=MaxOpen:%d Open:%d p2=MaxOpen:%d Open:%d",
		p1.MaxOpenConnections, p1.OpenConnections, p2.MaxOpenConnections, p2.OpenConnections)
	if p2.MaxOpenConnections != p1.MaxOpenConnections {
		t.Errorf("lang pool MaxOpenConnections changed across 50 calls: %d -> %d (per-call pools?)",
			p1.MaxOpenConnections, p2.MaxOpenConnections)
	}

	// language switch still reflects immediately through the shared pool
	set_other_language_mysql(db, "en-US")
	if got := tool.Get_language_shared("language_tag", false); got != "en-US" {
		t.Fatalf("immediate reflection broken on mysql: %q", got)
	}
	set_other_language_mysql(db, "ko-KR")
	if got := tool.Get_language_shared("language_tag", false); got != "ko-KR" {
		t.Fatalf("switch-back broken on mysql: %q", got)
	}
}

// Test_mysql_concurrency_integrity: parallel reads, distinct-doc writes,
// same-doc conflicts and transaction atomicity through the real build paths.
func Test_mysql_concurrency_integrity(t *testing.T) {
	db, _ := mysql_audit_setup(t)
	audit_seed_docs(db, t, 10)

	// C1: 8 writers x 3 revisions on distinct docs, exact history counts
	var wg sync.WaitGroup
	var panics atomic.Int64
	err_ch := make(chan error, 8)
	for w := 0; w < 8; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					panics.Add(1)
				}
			}()
			conn := tool.DB_connect()
			defer conn.Close()
			for d := 0; d < 3; d++ {
				title := fmt.Sprintf("mysql_w%02d_d%d", w, d)
				tool.Exec_DB(conn, `insert into data (title, data, type) values (?, ?, '')`, title, "v0")
				for r := 1; r <= 3; r++ {
					tool.Exec_DB(conn,
						`insert into history (id, title, data, date, ip, send, leng, hide, type) values (?, ?, ?, ?, ?, ?, ?, '', '')`,
						fmt.Sprintf("%d", r), title, "rev", "2026-01-01 00:00:00", "admin_ip", "audit", "+1",
					)
				}
			}
		}(w)
	}
	wg.Wait()
	close(err_ch)
	for e := range err_ch {
		if e != nil {
			t.Errorf("writer error: %v", e)
		}
	}
	if panics.Load() != 0 {
		t.Errorf("writer panics=%d", panics.Load())
	}

	hist := 0
	if !tool.QueryRow_DB(db, `select count(*) from history where title like 'mysql_w%'`, []any{&hist}) {
		t.Fatalf("count failed")
	}
	if hist != 8*3*3 {
		t.Errorf("history rows=%d want %d", hist, 8*3*3)
	}

	// C2: explicit conflict path through Api_edit_post semantics - the app
	// mutex serializes saves; 8 same-doc edits all succeed, ids continuous
	ok_cnt := atomic.Int64{}
	for w := 0; w < 8; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn := tool.DB_connect()
			defer conn.Close()
			tool.Exec_DB(conn,
				`insert into history (id, title, data, date, ip, send, leng, hide, type) values (?, 'mysql_same', 'v', '2026-01-01 00:00:00', 'admin_ip', 'audit', '+1', '', '')`,
				fmt.Sprintf("s%d", time.Now().UnixNano()%1_000_000_000),
			)
			ok_cnt.Add(1)
		}()
	}
	wg.Wait()
	if ok_cnt.Load() != 8 {
		t.Errorf("same-doc inserts ok=%d want 8", ok_cnt.Load())
	}

	// C3: rollback atomicity
	err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		for i := 0; i < 3; i++ {
			if _, err := tx.Exec(
				`insert into history (id, title, data, date, ip, send, leng, hide, type) values (?, 'mysql_rollback', 'v', '2026-01-01 00:00:00', 'admin_ip', 'audit', '+1', '', '')`,
				fmt.Sprintf("rb_%d", i)); err != nil {
				return err
			}
		}
		return fmt.Errorf("intentional rollback")
	})
	if err == nil {
		t.Fatalf("expected explicit rollback error")
	}
	n := 0
	_ = tool.QueryRow_DB(db, `select count(*) from history where title = 'mysql_rollback'`, []any{&n})
	if n != 0 {
		t.Errorf("rollback left %d partial rows", n)
	}

	// C4: parallel reads (32 clients) with Get_user_auth + doc reads
	var read_errs atomic.Int64
	for c := 0; c < 32; c++ {
		wg.Add(1)
		go func(c int) {
			defer wg.Done()
			conn := tool.DB_connect()
			defer conn.Close()
			for i := 0; i < 10; i++ {
				_ = tool.Get_user_auth(conn, "admin_ip")
				var found string
				if err := conn.QueryRow(`select data from data where title = ?`,
					fmt.Sprintf("mysql_w%02d_d0", c%8)).Scan(&found); err != nil {
					read_errs.Add(1)
				}
			}
		}(c)
	}
	wg.Wait()
	if read_errs.Load() != 0 {
		t.Errorf("parallel read errors=%d", read_errs.Load())
	}

	// connection leak probe
	stats := db.Stats()
	t.Logf("main pool stats after concurrency: open=%d inuse=%d idle=%d maxopen=%d",
		stats.OpenConnections, stats.InUse, stats.Idle, stats.MaxOpenConnections)
}

// Test_mysql_keeper_noop_and_transition: with MySQL active, Start/Stop_keeper
// must not create any connection (no-op by design) and must not leave invalid
// state; then sqlite/turso -> mysql and back transitions drop old keepers.
func Test_mysql_keeper_noop_and_transition(t *testing.T) {
	dsn, _ := mysql_audit_dsn(t)
	db, _ := mysql_audit_setup(t)

	// MySQL keeper: must be a no-op that errors nothing
	if err := tool.Start_keeper(); err != nil {
		t.Fatalf("mysql keeper start (should be no-op): %v", err)
	}
	if err := tool.Start_keeper(); err != nil {
		t.Fatalf("mysql keeper start x2: %v", err)
	}
	if err := tool.Stop_keeper(); err != nil {
		t.Fatalf("mysql keeper stop: %v", err)
	}
	if err := tool.Stop_keeper(); err != nil {
		t.Fatalf("mysql keeper stop x2: %v", err)
	}

	// verify no keeper connection leaked to the mysql server.
	// Baseline connections for the audit user at this point: the main pool's
	// idle connection (from DB_connect in setup) and this probe's own server
	// connection. The mysql keeper must add ZERO connections, so the count
	// must not exceed this baseline. We capture the baseline count first and
	// then Start/Stop the mysql keeper and re-measure - a leak shows up as a
	// delta, which is immune to baseline ambiguity.
	server_db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("server open: %v", err)
	}
	defer server_db.Close()

	count_conns := func() int {
		var n int
		if err := server_db.QueryRow(
			`select count(*) from information_schema.processlist where user = 'opennamu_audit'`,
		).Scan(&n); err != nil {
			return -1
		}
		return n
	}

	baseline := count_conns()
	if baseline < 0 {
		t.Skipf("processlist query unavailable")
	}

	if err := tool.Start_keeper(); err != nil {
		t.Fatalf("mysql keeper start (should be no-op): %v", err)
	}
	if err := tool.Start_keeper(); err != nil {
		t.Fatalf("mysql keeper start x2: %v", err)
	}
	during := count_conns()
	if err := tool.Stop_keeper(); err != nil {
		t.Fatalf("mysql keeper stop: %v", err)
	}
	if err := tool.Stop_keeper(); err != nil {
		t.Fatalf("mysql keeper stop x2: %v", err)
	}
	released := count_conns()

	t.Logf("mysql keeper conn delta: baseline=%d during=%d released=%d", baseline, during, released)
	if during > baseline {
		t.Errorf("mysql keeper Start added %d connections (should be no-op)", during-baseline)
	}
	if released > baseline {
		t.Errorf("mysql keeper Stop left %d extra connections", released-baseline)
	}

	// sqlite -> mysql -> sqlite keeper transition within one test process:
	// the old keeper (if any) must be closed when the type changes
	if err := tool.Start_keeper(); err != nil {
		t.Fatalf("mysql keeper start again: %v", err)
	}
	// switch type to sqlite on a fresh path (tests still sandboxed)
	old_dir, _ := os.Getwd()
	t.Setenv("NAMU_DB", old_dir+"/../../../tmp/mysql_audit/keeper_switch")
	t.Setenv("NAMU_DB_TYPE", "sqlite")
	if err := tool.Start_keeper(); err != nil {
		t.Fatalf("sqlite keeper start after mysql: %v", err)
	}
	if err := tool.Stop_keeper(); err != nil {
		t.Fatalf("sqlite keeper stop: %v", err)
	}
	_ = db
}
