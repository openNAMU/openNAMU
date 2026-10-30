package test

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"opennamu/route"
	"opennamu/route/tool"
)

func new_test_db(t *testing.T, db_path string) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", db_path)
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}

	exec := func(query string, args ...any) {
		if _, err := db.Exec(query, args...); err != nil {
			t.Fatalf("exec %s: %v", query, err)
		}
	}

	exec(`create table if not exists vote (name, id, subject, data, user, type, acl)`)
	exec(`create table if not exists rb (block, end, today, blocker, why, band, ongoing, login)`)
	exec(`create table if not exists user_set (name, id, data)`)
	exec(`create table if not exists acl (title, data, type)`)
	exec(`create table if not exists alist (name, acl)`)
	exec(`create table if not exists data (title, data, type)`)
	exec(`create table if not exists data_set (doc_name, doc_rev, set_name, set_data)`)
	exec(`create table if not exists ua_d (name, ip, ua, today, sub)`)
	exec(`create table if not exists re_admin (who, what, time)`)
	exec(`create table if not exists other (name, data, coverage)`)

	return db
}

func set_test_setting(db *sql.DB, t *testing.T, name string, value string) {
	t.Helper()
	if _, err := db.Exec(`insert into other (name, data, coverage) values (?, ?, '')`, name, value); err != nil {
		t.Fatalf("set setting %s: %v", name, err)
	}
}

func first_string(db *sql.DB, query string, args ...any) string {
	row := db.QueryRow(query, args...)
	var result string
	if err := row.Scan(&result); err != nil {
		return ""
	}
	return result
}

func first_int(db *sql.DB, query string, args ...any) int {
	row := db.QueryRow(query, args...)
	var result int
	if err := row.Scan(&result); err != nil {
		return 0
	}
	return result
}

func Test_Backup_run_creates_backup(t *testing.T) {
	dir := t.TempDir()
	db_path := filepath.Join(dir, "test.db")

	db := new_test_db(t, db_path)
	defer tool.DB_close(db)

	exec := func(query string, args ...any) {
		if _, err := db.Exec(query, args...); err != nil {
			t.Fatalf("exec %s: %v", query, err)
		}
	}

	set_test_setting(db, t, "back_up", "1")
	set_test_setting(db, t, "backup_count", "1")
	set_test_setting(db, t, "backup_where", db_path)

	exec(`insert into data (title, data, type) values ('Doc', 'content', 'd')`)
	if err := os.WriteFile(filepath.Join(dir, "test_not_a_timestamp.db"), []byte("unrelated"), 0600); err != nil {
		t.Fatalf("create unrelated file: %v", err)
	}

	if err := tool.Backup_run(db, db_path); err != nil {
		t.Fatalf("backup_run: %v", err)
	}

	pattern := regexp.MustCompile(`^test_[0-9]{14}(_[0-9]+)?\.db$`)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}

	backups := []string{}
	for _, entry := range entries {
		if pattern.MatchString(entry.Name()) {
			backups = append(backups, entry.Name())
		}
	}

	if len(backups) != 1 {
		t.Fatalf("expected one backup file, got %v", backups)
	}
	backup_info, err := os.Stat(filepath.Join(dir, backups[0]))
	if err != nil {
		t.Fatalf("stat backup file: %v", err)
	}
	if backup_info.Mode().Perm()&0077 != 0 {
		t.Fatalf("backup file permissions are not private: %04o", backup_info.Mode().Perm())
	}

	exec(`insert into data (title, data, type) values ('Doc2', 'content2', 'd')`)
	if err := tool.Backup_run(db, db_path); err != nil {
		t.Fatalf("backup_run second: %v", err)
	}

	entries, err = os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir second: %v", err)
	}

	backups = []string{}
	for _, entry := range entries {
		if pattern.MatchString(entry.Name()) {
			backups = append(backups, entry.Name())
		}
	}

	if len(backups) != 1 {
		t.Fatalf("expected one retained backup file, got %v", backups)
	}
	if _, err := os.Stat(filepath.Join(dir, "test_not_a_timestamp.db")); err != nil {
		t.Fatalf("retention removed unrelated file: %v", err)
	}
}

func Test_Backup_run_wal_restore_proof(t *testing.T) {
	dir := t.TempDir()
	db_path := filepath.Join(dir, "restore.db")
	db := new_test_db(t, db_path)
	defer tool.DB_close(db)

	if _, err := db.Exec(`PRAGMA journal_mode = WAL`); err != nil {
		t.Fatalf("enable WAL: %v", err)
	}
	if _, err := db.Exec(`create table history (id, title, data, date)`); err != nil {
		t.Fatalf("create history: %v", err)
	}
	if _, err := db.Exec(`create table filler (id, value)`); err != nil {
		t.Fatalf("create filler: %v", err)
	}
	set_test_setting(db, t, "back_up", "1")
	set_test_setting(db, t, "backup_count", "0")
	set_test_setting(db, t, "backup_where", db_path)

	marker := "committed WAL marker — 日本語 — 🧪"
	if _, err := db.Exec(`insert into data (title, data, type) values (?, ?, 'd')`, "文書", marker); err != nil {
		t.Fatalf("insert document: %v", err)
	}
	if _, err := db.Exec(`insert into history (id, title, data, date) values ('1', ?, ?, '2026-10-30')`, "文書", marker); err != nil {
		t.Fatalf("insert history: %v", err)
	}
	filler_value := strings.Repeat("x", 16*1024)
	for index := 0; index < 320; index++ {
		if _, err := db.Exec(`insert into filler (id, value) values (?, ?)`, fmt.Sprint(index), filler_value); err != nil {
			t.Fatalf("insert filler row: %v", err)
		}
	}
	var journal_mode string
	if err := db.QueryRow(`PRAGMA journal_mode`).Scan(&journal_mode); err != nil {
		t.Fatalf("read journal mode: %v", err)
	}
	if journal_mode != "wal" {
		t.Fatalf("expected WAL mode, got %q", journal_mode)
	}
	wal_info, err := os.Stat(db_path + "-wal")
	if err != nil || wal_info.Size() == 0 {
		t.Fatalf("expected nonempty WAL before backup, stat=%v err=%v", wal_info, err)
	}

	writer_started := make(chan struct{})
	backup_started := make(chan struct{})
	stop_writer := make(chan struct{})
	writer_done := make(chan struct{})
	writer_errors := make(chan error, 1)
	var writes_during_backup atomic.Int64
	var writer_start_once sync.Once
	go func() {
		defer close(writer_done)
		for index := 0; ; index++ {
			select {
			case <-stop_writer:
				return
			default:
			}

			tx, err := db.Begin()
			if err != nil {
				writer_start_once.Do(func() { close(writer_started) })
				writer_errors <- err
				return
			}
			title := fmt.Sprintf("Concurrent-%06d", index)
			if _, err := tx.Exec(`insert into data (title, data, type) values (?, ?, 'd')`, title, marker); err != nil {
				_ = tx.Rollback()
				writer_start_once.Do(func() { close(writer_started) })
				writer_errors <- err
				return
			}
			if _, err := tx.Exec(`insert into history (id, title, data, date) values (?, ?, ?, '2026-10-30')`, title, title, marker); err != nil {
				_ = tx.Rollback()
				writer_start_once.Do(func() { close(writer_started) })
				writer_errors <- err
				return
			}
			if err := tx.Commit(); err != nil {
				writer_start_once.Do(func() { close(writer_started) })
				writer_errors <- err
				return
			}
			select {
			case <-backup_started:
				writes_during_backup.Add(1)
			default:
			}
			if index == 0 {
				writer_start_once.Do(func() { close(writer_started) })
			}
			time.Sleep(time.Millisecond)
		}
	}()
	<-writer_started
	close(backup_started)
	backup_err := tool.Backup_run(db, db_path)
	close(stop_writer)
	<-writer_done
	select {
	case err := <-writer_errors:
		t.Fatalf("concurrent writer: %v", err)
	default:
	}
	if writes_during_backup.Load() == 0 {
		t.Fatal("concurrent writer did not commit during backup interval")
	}
	t.Logf("concurrent commits during backup interval: %d", writes_during_backup.Load())
	if backup_err != nil {
		t.Fatalf("backup_run with concurrent WAL writes: %v", backup_err)
	}

	pattern := regexp.MustCompile(`^restore_[0-9]{14}(_[0-9]+)?\.db$`)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read backup directory: %v", err)
	}
	backup_path := ""
	for _, entry := range entries {
		if pattern.MatchString(entry.Name()) {
			backup_path = filepath.Join(dir, entry.Name())
			break
		}
	}
	if backup_path == "" {
		t.Fatal("backup file was not created")
	}

	backup_db, err := sql.Open("sqlite", backup_path)
	if err != nil {
		t.Fatalf("open backup in separate connection: %v", err)
	}
	defer backup_db.Close()
	var integrity string
	if err := backup_db.QueryRow(`PRAGMA integrity_check`).Scan(&integrity); err != nil || integrity != "ok" {
		t.Fatalf("backup integrity_check=%q err=%v", integrity, err)
	}
	if got := first_string(backup_db, `select data from data where title = ?`, "文書"); got != marker {
		t.Fatalf("backup document marker mismatch: %q", got)
	}
	if got := first_string(backup_db, `select data from history where title = ?`, "文書"); got != marker {
		t.Fatalf("backup history marker mismatch: %q", got)
	}
	data_concurrent_count := first_int(backup_db, `select count(*) from data where title like 'Concurrent-%'`)
	history_concurrent_count := first_int(backup_db, `select count(*) from history where title like 'Concurrent-%'`)
	if data_concurrent_count != history_concurrent_count {
		t.Fatalf("backup contains incoherent concurrent rows: documents=%d histories=%d", data_concurrent_count, history_concurrent_count)
	}
	t.Logf("coherent concurrent rows included in backup: %d", data_concurrent_count)
	if count := first_int(backup_db, `select count(*) from data d left join history h on h.title = d.title where d.title like 'Concurrent-%' and h.title is null`); count != 0 {
		t.Fatalf("backup contains %d concurrent document rows without matching history", count)
	}
	if count := first_int(backup_db, `select count(*) from history h left join data d on d.title = h.title where h.title like 'Concurrent-%' and d.title is null`); count != 0 {
		t.Fatalf("backup contains %d concurrent history rows without matching document", count)
	}

	restored_dir := t.TempDir()
	restored_path := filepath.Join(restored_dir, "data.db")
	backup_bytes, err := os.ReadFile(backup_path)
	if err != nil {
		t.Fatalf("read backup for restore: %v", err)
	}
	if err := os.WriteFile(restored_path, backup_bytes, 0600); err != nil {
		t.Fatalf("restore backup to disposable path: %v", err)
	}
	restored_db, err := sql.Open("sqlite", restored_path)
	if err != nil {
		t.Fatalf("open restored database: %v", err)
	}
	defer restored_db.Close()
	if got := first_string(restored_db, `select data from data where title = ?`, "文書"); got != marker {
		t.Fatalf("restored document marker mismatch: %q", got)
	}
}

func Test_Backup_run_restores_openNAMU_startup_environment(t *testing.T) {
	dir := t.TempDir()
	source_dir := filepath.Join(dir, "source")
	if err := os.Mkdir(source_dir, 0700); err != nil {
		t.Fatalf("create source environment: %v", err)
	}
	source_name := filepath.Join(source_dir, "source")
	source_command := exec.Command(os.Args[0], "-test.run=^Test_backup_restore_environment_child$")
	source_command.Env = []string{
		"NAMU_DB=" + source_name,
		"NAMU_DB_TYPE=sqlite",
		"OPENNAMU_AUDIT_RESTORE_MODE=create",
		"OPENNAMU_AUDIT_RESTORE_DIR=" + source_dir,
	}
	if output, err := source_command.CombinedOutput(); err != nil {
		t.Fatalf("source environment child failed: %v\n%s", err, output)
	}

	entries, err := os.ReadDir(source_dir)
	if err != nil {
		t.Fatalf("read source directory: %v", err)
	}
	pattern := regexp.MustCompile(`^source_[0-9]{14}(_[0-9]+)?\.db$`)
	backup_path := ""
	for _, entry := range entries {
		if pattern.MatchString(entry.Name()) {
			backup_path = filepath.Join(source_dir, entry.Name())
			break
		}
	}
	if backup_path == "" {
		t.Fatal("full openNAMU backup was not created")
	}

	restored_dir := filepath.Join(dir, "restored")
	if err := os.Mkdir(restored_dir, 0700); err != nil {
		t.Fatalf("create restored environment: %v", err)
	}
	restored_name := filepath.Join(restored_dir, "data")
	backup_data, err := os.ReadFile(backup_path)
	if err != nil {
		t.Fatalf("read full backup: %v", err)
	}
	if err := os.WriteFile(restored_name+".db", backup_data, 0600); err != nil {
		t.Fatalf("restore full backup: %v", err)
	}
	restore_command := exec.Command(os.Args[0], "-test.run=^Test_backup_restore_environment_child$")
	restore_command.Env = []string{
		"NAMU_DB=" + restored_name,
		"NAMU_DB_TYPE=sqlite",
		"OPENNAMU_AUDIT_RESTORE_MODE=restore",
		"OPENNAMU_AUDIT_RESTORE_DIR=" + restored_dir,
	}
	if output, err := restore_command.CombinedOutput(); err != nil {
		t.Fatalf("restored environment startup child failed: %v\n%s", err, output)
	}
}

func Test_backup_restore_environment_child(t *testing.T) {
	mode := os.Getenv("OPENNAMU_AUDIT_RESTORE_MODE")
	if mode == "" {
		t.Skip("full restore subprocess helper")
	}
	version_data, err := os.ReadFile(filepath.Join("..", "version.json"))
	if err != nil {
		t.Fatalf("read embedded version data: %v", err)
	}
	t.Chdir(os.Getenv("OPENNAMU_AUDIT_RESTORE_DIR"))
	t.Cleanup(func() { _ = tool.Stop_keeper() })
	tool.Set_builtin_version_data(version_data)
	tool.DB_boot()
	version := tool.Main_init()

	deadline := time.Now().Add(10 * time.Second)
	for !tool.Search_index_ready() && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if !tool.Search_index_ready() {
		t.Fatal("search index did not become ready during startup")
	}

	if mode == "create" {
		source_name := os.Getenv("NAMU_DB")
		db := tool.DB_connect()
		marker := "restored openNAMU marker — 日本語 — 🧪"
		if _, err := db.Exec(`insert into data (title, data, type) values (?, ?, 'd')`, "復元文書", marker); err != nil {
			t.Fatalf("seed source document: %v", err)
		}
		if _, err := db.Exec(`insert into history (id, title, data, date, ip, send, leng, hide, type) values ('audit', ?, ?, '2026-10-30', '', '', '0', '0', 'edit')`, "復元文書", marker); err != nil {
			t.Fatalf("seed source history: %v", err)
		}
		set_test_setting(db, t, "back_up", "1")
		set_test_setting(db, t, "backup_count", "0")
		set_test_setting(db, t, "backup_where", source_name+".db")
		if err := tool.Backup_run(db, source_name+".db"); err != nil {
			t.Fatalf("backup source environment: %v", err)
		}
		tool.DB_close(db)
		return
	}

	if mode != "restore" {
		t.Fatalf("unknown restore probe mode %q", mode)
	}
	if version == "" {
		t.Fatal("restored environment startup did not read the database version")
	}
	restored_db := tool.DB_connect()
	defer tool.DB_close(restored_db)
	marker := "restored openNAMU marker — 日本語 — 🧪"
	if got := first_string(restored_db, `select data from data where title = ?`, "復元文書"); got != marker {
		t.Fatalf("restored openNAMU document read mismatch: %q", got)
	}
	if got := first_string(restored_db, `select data from history where title = ?`, "復元文書"); got != marker {
		t.Fatalf("restored openNAMU history read mismatch: %q", got)
	}
}

func Test_Backup_run_follows_configured_symlink_destination(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation may require elevated privileges on Windows")
	}
	dir := t.TempDir()
	db_path := filepath.Join(dir, "source.db")
	db := new_test_db(t, db_path)
	defer tool.DB_close(db)

	outside_dir := t.TempDir()
	link_dir := filepath.Join(dir, "linked-backups")
	if err := os.Symlink(outside_dir, link_dir); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}
	set_test_setting(db, t, "back_up", "1")
	set_test_setting(db, t, "backup_count", "3")
	set_test_setting(db, t, "backup_where", filepath.Join(link_dir, "source.db"))

	if err := tool.Backup_run(db, db_path); err != nil {
		t.Fatalf("configured symlink destination should remain usable: %v", err)
	}
	entries, err := os.ReadDir(outside_dir)
	if err != nil {
		t.Fatalf("read outside destination: %v", err)
	}
	if len(entries) != 1 || !regexp.MustCompile(`^source_[0-9]{14}(_[0-9]+)?\.db$`).MatchString(entries[0].Name()) {
		t.Fatalf("backup was not written to configured symlink target: %v", entries)
	}
	if _, _, ok := route.Read_root_file(entries[0].Name()); ok {
		t.Fatalf("SQLite backup should not be exposed by the root-file route: %s", entries[0].Name())
	}
}

func Test_Backup_run_preserves_configured_base_name_case(t *testing.T) {
	dir := t.TempDir()
	db_path := filepath.Join(dir, "source.db")
	db := new_test_db(t, db_path)
	defer tool.DB_close(db)
	set_test_setting(db, t, "back_up", "1")
	set_test_setting(db, t, "backup_count", "0")
	set_test_setting(db, t, "backup_where", filepath.Join(dir, "MyBackup.DB"))

	if err := tool.Backup_run(db, db_path); err != nil {
		t.Fatalf("backup with mixed-case base: %v", err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read backup directory: %v", err)
	}
	pattern := regexp.MustCompile(`^MyBackup_[0-9]{14}(_[0-9]+)?\.db$`)
	for _, entry := range entries {
		if pattern.MatchString(entry.Name()) {
			return
		}
	}
	t.Fatalf("configured backup base case was not preserved: %v", entries)
}

func Test_Acquire_auto_maintenance_lock_recovers_stale_lock(t *testing.T) {
	lock_path := filepath.Join(t.TempDir(), ".opennamu_auto_backup_lock")
	if err := os.WriteFile(lock_path, []byte("123|0"), 0600); err != nil {
		t.Fatalf("create stale lock: %v", err)
	}
	stale_time := time.Now().Add(-31 * time.Minute)
	if err := os.Chtimes(lock_path, stale_time, stale_time); err != nil {
		t.Fatalf("age stale lock: %v", err)
	}

	lock, err := tool.Acquire_auto_maintenance_lock(lock_path)
	if err != nil {
		t.Fatalf("acquire should reclaim stale lock: %v", err)
	}
	lock.Release()
}

func Test_Auto_task_retry_state_separates_success_and_failure(t *testing.T) {
	dir := t.TempDir()
	last_success_path := filepath.Join(dir, ".last_success")
	retry_after_path := filepath.Join(dir, ".retry_after")
	start_time := time.Date(2026, 10, 30, 0, 0, 0, 0, time.UTC)
	interval := 24 * time.Hour

	if !tool.Auto_task_is_due(last_success_path, retry_after_path, interval, start_time) {
		t.Fatal("task without prior success should be due")
	}
	if err := tool.Record_auto_task_success(last_success_path, retry_after_path, start_time); err != nil {
		t.Fatalf("record success: %v", err)
	}
	if tool.Auto_task_is_due(last_success_path, retry_after_path, interval, start_time.Add(interval-time.Second)) {
		t.Fatal("successful task repeated before its interval")
	}
	if !tool.Auto_task_is_due(last_success_path, retry_after_path, interval, start_time.Add(interval)) {
		t.Fatal("successful task did not become due at interval")
	}

	failure_time := start_time.Add(interval)
	if err := tool.Record_auto_task_failure(retry_after_path, failure_time); err != nil {
		t.Fatalf("record transient failure: %v", err)
	}
	if tool.Auto_task_is_due(last_success_path, retry_after_path, interval, failure_time.Add(4*time.Minute)) {
		t.Fatal("failed task retried before bounded retry delay")
	}
	retry_time := failure_time.Add(5 * time.Minute)
	if !tool.Auto_task_is_due(last_success_path, retry_after_path, interval, retry_time) {
		t.Fatal("failed task was not eligible after bounded retry delay")
	}
	if err := tool.Record_auto_task_success(last_success_path, retry_after_path, retry_time); err != nil {
		t.Fatalf("record retry success: %v", err)
	}
	if _, err := os.Stat(retry_after_path); !os.IsNotExist(err) {
		t.Fatalf("successful retry did not clear retry state: %v", err)
	}
	if tool.Auto_task_is_due(last_success_path, retry_after_path, interval, retry_time.Add(interval-time.Second)) {
		t.Fatal("successful retry was repeated before next schedule interval")
	}
}

func Test_scheduler_lock_crash_child_process(t *testing.T) {
	if os.Getenv("OPENNAMU_AUDIT_LOCK_CHILD") != "1" {
		t.Skip("lock crash subprocess helper")
	}
	lock, err := tool.Acquire_auto_maintenance_lock(os.Getenv("OPENNAMU_AUDIT_LOCK_PATH"))
	if err != nil {
		t.Fatalf("child lock acquisition: %v", err)
	}
	_ = lock
	os.Exit(0)
}

func Test_scheduler_maintenance_crash_child_process(t *testing.T) {
	if os.Getenv("OPENNAMU_AUDIT_MAINTENANCE_CRASH_CHILD") != "1" {
		t.Skip("maintenance crash subprocess helper")
	}
	tool.DB_boot()
	lock, err := tool.Acquire_auto_maintenance_lock(os.Getenv("OPENNAMU_AUDIT_MAINTENANCE_LOCK_PATH"))
	if err != nil {
		t.Fatalf("maintenance lock acquisition: %v", err)
	}
	_ = lock
	db := tool.DB_connect()
	today := time.Now().Format("2006-01-02")
	if err := tool.Maintenance_vote_close(db, today); err != nil {
		t.Fatalf("maintenance group before simulated process death: %v", err)
	}
	tool.DB_close(db)
	os.Exit(0)
}

func Test_Acquire_auto_maintenance_lock_recovers_after_process_death(t *testing.T) {
	lock_path := filepath.Join(t.TempDir(), ".opennamu_auto_backup_lock")
	command := exec.Command(os.Args[0], "-test.run=^Test_scheduler_lock_crash_child_process$")
	command.Env = []string{
		"OPENNAMU_AUDIT_LOCK_CHILD=1",
		"OPENNAMU_AUDIT_LOCK_PATH=" + lock_path,
	}
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("crash-lock child failed: %v\n%s", err, output)
	}
	if _, err := os.Stat(lock_path); err != nil {
		t.Fatalf("child did not leave a lock file: %v", err)
	}
	stale_time := time.Now().Add(-31 * time.Minute)
	if err := os.Chtimes(lock_path, stale_time, stale_time); err != nil {
		t.Fatalf("age orphaned lock: %v", err)
	}

	lock, err := tool.Acquire_auto_maintenance_lock(lock_path)
	if err != nil {
		t.Fatalf("acquire should recover process-death lock: %v", err)
	}
	lock.Release()
}

func Test_scheduler_recovers_after_maintenance_process_death(t *testing.T) {
	dir := t.TempDir()
	db_name := filepath.Join(dir, "maintenance-crash")
	db_path := db_name + ".db"
	db := new_test_db(t, db_path)
	set_test_setting(db, t, "back_up", "")
	set_test_setting(db, t, "sitemap_auto_make", "")
	set_test_setting(db, t, "ua_expiration_date", "")
	set_test_setting(db, t, "auth_history_expiration_date", "")
	if _, err := db.Exec(`insert into vote (name, id, subject, data, user, type, acl) values ('', 'crash-vote', 'subject', '', '', 'open', '')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into vote (name, id, subject, data, user, type, acl) values ('end_date', 'crash-vote', '', '2000-01-01', '', 'option', '')`); err != nil {
		t.Fatal(err)
	}
	tool.DB_close(db)

	lock_path := filepath.Join(dir, ".opennamu_auto_backup_lock")
	crash_command := exec.Command(os.Args[0], "-test.run=^Test_scheduler_maintenance_crash_child_process$")
	crash_command.Env = []string{
		"NAMU_DB=" + db_name,
		"NAMU_DB_TYPE=sqlite",
		"OPENNAMU_AUDIT_MAINTENANCE_CRASH_CHILD=1",
		"OPENNAMU_AUDIT_MAINTENANCE_LOCK_PATH=" + lock_path,
	}
	if output, err := crash_command.CombinedOutput(); err != nil {
		t.Fatalf("maintenance crash child failed: %v\n%s", err, output)
	}
	if first_string(Must_open_db(t, db_path), `select type from vote where id = 'crash-vote' and user = ''`) != "close" {
		t.Fatal("maintenance child did not commit its scheduled group before process death")
	}
	if _, err := os.Stat(filepath.Join(dir, ".opennamu_auto_maintenance_last_run")); !os.IsNotExist(err) {
		t.Fatalf("crashed maintenance unexpectedly recorded success: %v", err)
	}
	stale_time := time.Now().Add(-31 * time.Minute)
	if err := os.Chtimes(lock_path, stale_time, stale_time); err != nil {
		t.Fatalf("age crashed task lock: %v", err)
	}

	retry_command := exec.Command(os.Args[0], "-test.run=^Test_scheduler_child_process$")
	retry_command.Env = []string{
		"NAMU_DB=" + db_name,
		"NAMU_DB_TYPE=sqlite",
		"OPENNAMU_AUDIT_SCHEDULER_CHILD=1",
	}
	if output, err := retry_command.CombinedOutput(); err != nil {
		t.Fatalf("scheduler recovery child failed: %v\n%s", err, output)
	}
	if _, err := os.Stat(filepath.Join(dir, ".opennamu_auto_maintenance_last_run")); err != nil {
		t.Fatalf("recovered scheduler did not record successful maintenance: %v", err)
	}
	if first_string(Must_open_db(t, db_path), `select type from vote where id = 'crash-vote' and user = ''`) != "close" {
		t.Fatal("repeated maintenance changed already-closed vote")
	}
}

func Test_scheduler_child_process(t *testing.T) {
	if os.Getenv("OPENNAMU_AUDIT_SCHEDULER_CHILD") != "1" {
		t.Skip("scheduler subprocess helper")
	}
	tool.DB_boot()
	tool.Start_auto_backup_and_maintenance()
	time.Sleep(2 * time.Second)
}

func Test_scheduler_backend_child_process(t *testing.T) {
	if os.Getenv("OPENNAMU_AUDIT_BACKEND_CHILD") != "1" {
		t.Skip("non-SQLite scheduler subprocess helper")
	}
	tool.DB_boot()
	tool.Start_auto_backup_and_maintenance()
	time.Sleep(2 * time.Second)
}

func Test_scheduler_mysql_child_process(t *testing.T) {
	if os.Getenv("OPENNAMU_AUDIT_MYSQL_SCHEDULER_CHILD") != "1" {
		t.Skip("MariaDB scheduler subprocess helper")
	}
	tool.DB_boot()
	tool.Start_auto_backup_and_maintenance()
	time.Sleep(2 * time.Second)
}

func Test_scheduler_unavailable_mysql_child_process(t *testing.T) {
	if os.Getenv("OPENNAMU_AUDIT_MYSQL_FAILURE_CHILD") != "1" {
		t.Skip("unavailable MySQL scheduler subprocess helper")
	}
	tool.DB_boot()
	tool.Start_auto_backup_and_maintenance()
	time.Sleep(2 * time.Second)
}

func Test_scheduler_multi_process_startup(t *testing.T) {
	dir := t.TempDir()
	db_name := filepath.Join(dir, "scheduler")
	db_path := db_name + ".db"
	db := new_test_db(t, db_path)
	set_test_setting(db, t, "back_up", "1")
	set_test_setting(db, t, "backup_count", "0")
	set_test_setting(db, t, "backup_where", db_path)
	set_test_setting(db, t, "sitemap_auto_make", "")
	set_test_setting(db, t, "ua_expiration_date", "")
	set_test_setting(db, t, "auth_history_expiration_date", "")
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', 'scheduled-user', 'admin')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('auth_date', 'scheduled-user', '2000-01-01')`); err != nil {
		t.Fatal(err)
	}
	tool.DB_close(db)

	var wait_group sync.WaitGroup
	outputs := make([][]byte, 2)
	errors_list := make([]error, 2)
	for index := range outputs {
		wait_group.Add(1)
		go func(index int) {
			defer wait_group.Done()
			command := exec.Command(os.Args[0], "-test.run=^Test_scheduler_child_process$")
			command.Env = []string{
				"NAMU_DB=" + db_name,
				"NAMU_DB_TYPE=sqlite",
				"OPENNAMU_AUDIT_SCHEDULER_CHILD=1",
			}
			outputs[index], errors_list[index] = command.CombinedOutput()
		}(index)
	}
	wait_group.Wait()
	for index, err := range errors_list {
		if err != nil {
			t.Fatalf("scheduler child %d failed: %v\n%s", index, err, outputs[index])
		}
	}
	staggered_db_name := db_name
	if runtime.GOOS != "windows" {
		alias_dir := filepath.Join(dir, "alias")
		if err := os.Mkdir(alias_dir, 0700); err != nil {
			t.Fatalf("create alias directory: %v", err)
		}
		staggered_db_name = filepath.Join(alias_dir, "scheduler")
		if err := os.Symlink(db_path, staggered_db_name+".db"); err != nil {
			t.Logf("source DB symlink unavailable; staggered process uses canonical path: %v", err)
			staggered_db_name = db_name
		}
	}
	staggered_command := exec.Command(os.Args[0], "-test.run=^Test_scheduler_child_process$")
	staggered_command.Env = []string{
		"NAMU_DB=" + staggered_db_name,
		"NAMU_DB_TYPE=sqlite",
		"OPENNAMU_AUDIT_SCHEDULER_CHILD=1",
	}
	if output, err := staggered_command.CombinedOutput(); err != nil {
		t.Fatalf("staggered scheduler child failed: %v\n%s", err, output)
	}

	pattern := regexp.MustCompile(`^scheduler_[0-9]{14}(_[0-9]+)?\.db$`)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read scheduler directory: %v", err)
	}
	backup_count := 0
	for _, entry := range entries {
		if pattern.MatchString(entry.Name()) {
			backup_count++
		}
	}
	if backup_count != 1 {
		t.Fatalf("expected one scheduled backup across two processes, found %d; entries=%v", backup_count, entries)
	}
	if got := first_string(Must_open_db(t, db_path), `select data from user_set where id = 'scheduled-user' and name = 'acl'`); got != "user" {
		t.Fatalf("daily maintenance did not expire legacy grant: %q", got)
	}
	if _, err := os.Stat(filepath.Join(dir, ".opennamu_auto_maintenance_last_run")); err != nil {
		t.Fatalf("daily maintenance completion state missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, ".opennamu_auto_backup_last_run")); err != nil {
		t.Fatalf("backup completion state missing: %v", err)
	}
}

func Test_scheduler_failed_backup_does_not_record_success(t *testing.T) {
	dir := t.TempDir()
	db_name := filepath.Join(dir, "scheduler-failure")
	db_path := db_name + ".db"
	db := new_test_db(t, db_path)
	set_test_setting(db, t, "back_up", "1")
	set_test_setting(db, t, "backup_count", "0")
	blocker_path := filepath.Join(dir, "not-a-directory")
	if err := os.WriteFile(blocker_path, []byte("block"), 0600); err != nil {
		t.Fatalf("create destination blocker: %v", err)
	}
	set_test_setting(db, t, "backup_where", filepath.Join(blocker_path, "source.db"))
	set_test_setting(db, t, "sitemap_auto_make", "")
	set_test_setting(db, t, "ua_expiration_date", "")
	set_test_setting(db, t, "auth_history_expiration_date", "")
	tool.DB_close(db)

	command := exec.Command(os.Args[0], "-test.run=^Test_scheduler_child_process$")
	command.Env = []string{
		"NAMU_DB=" + db_name,
		"NAMU_DB_TYPE=sqlite",
		"OPENNAMU_AUDIT_SCHEDULER_CHILD=1",
	}
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("scheduler child failed: %v\n%s", err, output)
	}
	last_success_path := filepath.Join(dir, ".opennamu_auto_backup_last_run")
	if _, err := os.Stat(last_success_path); !os.IsNotExist(err) {
		t.Fatalf("failed backup incorrectly recorded success: stat error=%v", err)
	}
	retry_after_path := filepath.Join(dir, ".opennamu_auto_backup_retry_after")
	retry_after, exists, err := tool.Read_auto_run_time(retry_after_path)
	if err != nil || !exists {
		t.Fatalf("failed backup retry state missing: exists=%v err=%v", exists, err)
	}
	if tool.Auto_task_is_due(last_success_path, retry_after_path, time.Hour, time.Now().Add(4*time.Minute)) {
		t.Fatal("failed backup retry became due before retry_after")
	}
	if !tool.Auto_task_is_due(last_success_path, retry_after_path, time.Hour, retry_after) {
		t.Fatal("failed backup retry was not due at retry_after")
	}
}

func Test_scheduler_temporary_backup_permission_failure_retries_without_success_state(t *testing.T) {
	dir := t.TempDir()
	db_name := filepath.Join(dir, "scheduler-permission-failure")
	db_path := db_name + ".db"
	backup_dir := filepath.Join(dir, "backups")
	if err := os.Mkdir(backup_dir, 0700); err != nil {
		t.Fatalf("create backup directory: %v", err)
	}
	db := new_test_db(t, db_path)
	set_test_setting(db, t, "back_up", "1")
	set_test_setting(db, t, "backup_count", "1")
	set_test_setting(db, t, "backup_where", filepath.Join(backup_dir, "source.db"))
	set_test_setting(db, t, "sitemap_auto_make", "")
	set_test_setting(db, t, "ua_expiration_date", "")
	set_test_setting(db, t, "auth_history_expiration_date", "")
	if err := tool.Backup_run(db, db_path); err != nil {
		t.Fatalf("create initial valid backup: %v", err)
	}
	tool.DB_close(db)

	if err := os.Chmod(backup_dir, 0500); err != nil {
		t.Fatalf("make backup directory temporarily unwritable: %v", err)
	}
	probe, probe_err := os.CreateTemp(backup_dir, "permission_probe")
	if probe_err == nil {
		_ = probe.Close()
		_ = os.Remove(probe.Name())
		_ = os.Chmod(backup_dir, 0700)
		t.Skip("filesystem permissions are not enforced for this test user")
	}
	defer os.Chmod(backup_dir, 0700)

	command := exec.Command(os.Args[0], "-test.run=^Test_scheduler_child_process$")
	command.Env = []string{
		"NAMU_DB=" + db_name,
		"NAMU_DB_TYPE=sqlite",
		"OPENNAMU_AUDIT_SCHEDULER_CHILD=1",
	}
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("scheduler child failed: %v\n%s", err, output)
	}
	if _, err := os.Stat(filepath.Join(dir, ".opennamu_auto_backup_last_run")); !os.IsNotExist(err) {
		t.Fatalf("failed backup attempt advanced successful-run state: %v", err)
	}
	retry_after_path := filepath.Join(dir, ".opennamu_auto_backup_retry_after")
	retry_after, exists, err := tool.Read_auto_run_time(retry_after_path)
	if err != nil || !exists {
		t.Fatalf("temporary permission failure retry state missing: exists=%v err=%v", exists, err)
	}
	if tool.Auto_task_is_due(filepath.Join(dir, ".opennamu_auto_backup_last_run"), retry_after_path, time.Hour, time.Now().Add(4*time.Minute)) {
		t.Fatal("temporary permission failure retried before bounded delay")
	}
	if !tool.Auto_task_is_due(filepath.Join(dir, ".opennamu_auto_backup_last_run"), retry_after_path, time.Hour, retry_after) {
		t.Fatal("temporary permission failure was not eligible at retry deadline")
	}
	entries, err := os.ReadDir(backup_dir)
	if err != nil {
		t.Fatalf("read backup directory after failure: %v", err)
	}
	backup_count := 0
	for _, entry := range entries {
		if regexp.MustCompile(`^source_[0-9]{14}(_[0-9]+)?\.db$`).MatchString(entry.Name()) {
			backup_count++
		}
	}
	if backup_count != 1 {
		t.Fatalf("temporary failure damaged prior valid backup: found %d", backup_count)
	}
}

func Test_scheduler_failed_maintenance_does_not_record_success(t *testing.T) {
	dir := t.TempDir()
	db_name := filepath.Join(dir, "scheduler-maintenance-failure")
	db_path := db_name + ".db"
	db := new_test_db(t, db_path)
	set_test_setting(db, t, "back_up", "")
	set_test_setting(db, t, "sitemap_auto_make", "1")
	if _, err := db.Exec(`drop table data`); err != nil {
		t.Fatalf("inject sitemap query failure: %v", err)
	}
	tool.DB_close(db)

	command := exec.Command(os.Args[0], "-test.run=^Test_scheduler_child_process$")
	command.Env = []string{
		"NAMU_DB=" + db_name,
		"NAMU_DB_TYPE=sqlite",
		"OPENNAMU_AUDIT_SCHEDULER_CHILD=1",
	}
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("scheduler child failed: %v\n%s", err, output)
	}
	last_success_path := filepath.Join(dir, ".opennamu_auto_maintenance_last_run")
	if _, err := os.Stat(last_success_path); !os.IsNotExist(err) {
		t.Fatalf("partially failed maintenance incorrectly recorded success: stat error=%v", err)
	}
	retry_after_path := filepath.Join(dir, ".opennamu_auto_maintenance_retry_after")
	retry_after, exists, err := tool.Read_auto_run_time(retry_after_path)
	if err != nil || !exists {
		t.Fatalf("failed maintenance retry state missing: exists=%v err=%v", exists, err)
	}
	if tool.Auto_task_is_due(last_success_path, retry_after_path, 24*time.Hour, time.Now()) {
		t.Fatal("failed maintenance retry was due before the persisted retry deadline")
	}
	if !tool.Auto_task_is_due(last_success_path, retry_after_path, 24*time.Hour, retry_after) {
		t.Fatal("failed maintenance retry was not due at the persisted retry deadline")
	}
}

func Test_scheduler_lock_contention_does_not_advance_run_state(t *testing.T) {
	dir := t.TempDir()
	db_name := filepath.Join(dir, "scheduler-lock-contention")
	db_path := db_name + ".db"
	db := new_test_db(t, db_path)
	set_test_setting(db, t, "back_up", "1")
	set_test_setting(db, t, "backup_count", "0")
	set_test_setting(db, t, "backup_where", db_path)
	set_test_setting(db, t, "sitemap_auto_make", "")
	set_test_setting(db, t, "ua_expiration_date", "")
	set_test_setting(db, t, "auth_history_expiration_date", "")
	tool.DB_close(db)

	lock_path := filepath.Join(dir, ".opennamu_auto_backup_lock")
	if err := os.WriteFile(lock_path, []byte("other-process|active"), 0600); err != nil {
		t.Fatalf("create active lock fixture: %v", err)
	}
	command := exec.Command(os.Args[0], "-test.run=^Test_scheduler_child_process$")
	command.Env = []string{
		"NAMU_DB=" + db_name,
		"NAMU_DB_TYPE=sqlite",
		"OPENNAMU_AUDIT_SCHEDULER_CHILD=1",
	}
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("scheduler child failed: %v\n%s", err, output)
	}
	for _, state_file := range []string{
		".opennamu_auto_backup_last_run",
		".opennamu_auto_backup_retry_after",
		".opennamu_auto_maintenance_last_run",
		".opennamu_auto_maintenance_retry_after",
	} {
		if _, err := os.Stat(filepath.Join(dir, state_file)); !os.IsNotExist(err) {
			t.Fatalf("lock contention advanced scheduler state %s: %v", state_file, err)
		}
	}
}

func Test_scheduler_turso_startup_runs_daily_maintenance(t *testing.T) {
	if !tool.Turso_platform_supported() {
		t.Skip("embedded Turso is unsupported on this platform")
	}
	dir := t.TempDir()
	db_name := filepath.Join(dir, "scheduler-turso")
	t.Cleanup(func() { tool.DB_boot() })
	t.Setenv("NAMU_DB", db_name)
	t.Setenv("NAMU_DB_TYPE", "turso")
	tool.DB_init()
	db := tool.DB_connect()
	if _, err := db.Exec(`insert into other (name, data, coverage) values ('sitemap_auto_make', '', '')`); err != nil {
		t.Fatalf("disable sitemap generation: %v", err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', 'turso-scheduled-user', 'admin')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('auth_date', 'turso-scheduled-user', '2000-01-01')`); err != nil {
		t.Fatal(err)
	}
	tool.DB_close(db)

	var wait_group sync.WaitGroup
	outputs := make([][]byte, 2)
	errors_list := make([]error, 2)
	for index := range outputs {
		wait_group.Add(1)
		go func(index int) {
			defer wait_group.Done()
			command := exec.Command(os.Args[0], "-test.run=^Test_scheduler_backend_child_process$")
			command.Env = []string{
				"NAMU_DB=" + db_name,
				"NAMU_DB_TYPE=turso",
				"OPENNAMU_AUDIT_BACKEND_CHILD=1",
			}
			outputs[index], errors_list[index] = command.CombinedOutput()
		}(index)
	}
	wait_group.Wait()
	for index, err := range errors_list {
		if err != nil {
			t.Fatalf("Turso scheduler child %d failed: %v\n%s", index, err, outputs[index])
		}
	}
	maintenance_state_path := filepath.Join(dir, ".opennamu_auto_maintenance_last_run")
	state_before_staggered_start, err := os.ReadFile(maintenance_state_path)
	if err != nil {
		t.Fatalf("Turso maintenance state missing after simultaneous start: %v", err)
	}
	staggered_command := exec.Command(os.Args[0], "-test.run=^Test_scheduler_backend_child_process$")
	staggered_command.Env = []string{
		"NAMU_DB=" + db_name,
		"NAMU_DB_TYPE=turso",
		"OPENNAMU_AUDIT_BACKEND_CHILD=1",
	}
	if output, err := staggered_command.CombinedOutput(); err != nil {
		t.Fatalf("staggered Turso scheduler child failed: %v\n%s", err, output)
	}
	state_after_staggered_start, err := os.ReadFile(maintenance_state_path)
	if err != nil {
		t.Fatalf("Turso maintenance state after staggered start: %v", err)
	}
	if string(state_after_staggered_start) != string(state_before_staggered_start) {
		t.Fatal("staggered Turso process repeated daily maintenance within the 24-hour interval")
	}

	tool.DB_boot()
	db = tool.DB_connect()
	defer tool.DB_close(db)
	if got := first_string(db, `select data from user_set where id = 'turso-scheduled-user' and name = 'acl'`); got != "user" {
		t.Fatalf("Turso scheduler did not run daily maintenance: %q", got)
	}
	if _, err := os.Stat(maintenance_state_path); err != nil {
		t.Fatalf("Turso daily maintenance state missing: %v", err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read Turso scheduler directory: %v", err)
	}
	for _, entry := range entries {
		if regexp.MustCompile(`^scheduler-turso_[0-9]{14}(_[0-9]+)?\.db$`).MatchString(entry.Name()) {
			t.Fatalf("Turso scheduler unexpectedly created a SQLite backup: %s", entry.Name())
		}
	}
}

func Must_open_db(t *testing.T, db_path string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", db_path)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func Test_Backup_run_disabled_and_invalid(t *testing.T) {
	t.Run("empty interval", func(t *testing.T) {
		dir := t.TempDir()
		db_path := filepath.Join(dir, "test.db")
		db := new_test_db(t, db_path)
		defer tool.DB_close(db)

		set_test_setting(db, t, "back_up", "")
		set_test_setting(db, t, "backup_count", "3")
		set_test_setting(db, t, "backup_where", "test.db")

		if err := tool.Backup_run(db, db_path); err != nil {
			t.Fatalf("empty interval should be disabled, got: %v", err)
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("read dir: %v", err)
		}
		for _, entry := range entries {
			if entry.Name() != "test.db" {
				t.Fatalf("unexpected file: %s", entry.Name())
			}
		}
	})

	t.Run("invalid interval", func(t *testing.T) {
		dir := t.TempDir()
		db_path := filepath.Join(dir, "test.db")
		db := new_test_db(t, db_path)
		defer tool.DB_close(db)

		set_test_setting(db, t, "back_up", "abc")
		set_test_setting(db, t, "backup_count", "3")
		set_test_setting(db, t, "backup_where", "test.db")

		if err := tool.Backup_run(db, db_path); err == nil {
			t.Fatalf("expected invalid interval to fail")
		}
	})

	t.Run("invalid count", func(t *testing.T) {
		dir := t.TempDir()
		db_path := filepath.Join(dir, "test.db")
		db := new_test_db(t, db_path)
		defer tool.DB_close(db)

		set_test_setting(db, t, "back_up", "1")
		set_test_setting(db, t, "backup_count", "-1")
		set_test_setting(db, t, "backup_where", "test.db")

		if err := tool.Backup_run(db, db_path); err == nil {
			t.Fatalf("expected invalid backup_count to fail")
		}
	})
}

func Test_Backup_parse_interval_rejects_nonfinite_and_overflow_values(t *testing.T) {
	for _, value := range []string{"NaN", "+Inf", "1e100"} {
		if _, err := tool.Backup_parse_interval_hours(value); err == nil {
			t.Errorf("expected %q to be rejected", value)
		}
	}
}

func Test_Backup_parse_count_default(t *testing.T) {
	count, err := tool.Backup_parse_count("")
	if err != nil {
		t.Fatalf("parse default backup count: %v", err)
	}
	if count != 3 {
		t.Fatalf("default backup count = %d, want 3", count)
	}
}

func Test_Backup_run_default_retention_keeps_latest_three(t *testing.T) {
	dir := t.TempDir()
	db_path := filepath.Join(dir, "default.db")
	db := new_test_db(t, db_path)
	defer tool.DB_close(db)
	set_test_setting(db, t, "back_up", "1")
	set_test_setting(db, t, "backup_count", "")
	set_test_setting(db, t, "backup_where", db_path)

	for index := 0; index < 4; index++ {
		if err := tool.Backup_run(db, db_path); err != nil {
			t.Fatalf("backup %d: %v", index+1, err)
		}
	}
	pattern := regexp.MustCompile(`^default_[0-9]{14}(_[0-9]+)?\.db$`)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read backups: %v", err)
	}
	backup_count := 0
	for _, entry := range entries {
		if pattern.MatchString(entry.Name()) {
			backup_count++
		}
	}
	if backup_count != 3 {
		t.Fatalf("default retention kept %d backups, want 3", backup_count)
	}
}

func Test_Backup_run_unlimited_and_failed_attempt_preserves_backups(t *testing.T) {
	dir := t.TempDir()
	db_path := filepath.Join(dir, "source.db")
	db := new_test_db(t, db_path)
	defer tool.DB_close(db)
	set_test_setting(db, t, "back_up", "1")
	set_test_setting(db, t, "backup_count", "0")
	set_test_setting(db, t, "backup_where", db_path)

	for index := 0; index < 2; index++ {
		if err := tool.Backup_run(db, db_path); err != nil {
			t.Fatalf("backup %d: %v", index+1, err)
		}
	}
	pattern := regexp.MustCompile(`^source_[0-9]{14}(_[0-9]+)?\.db$`)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read backups: %v", err)
	}
	backup_count := 0
	for _, entry := range entries {
		if pattern.MatchString(entry.Name()) {
			backup_count++
		}
	}
	if backup_count != 2 {
		t.Fatalf("backup_count=0 should retain both backups, got %d", backup_count)
	}

	blocker_path := filepath.Join(dir, "not-a-directory")
	if err := os.WriteFile(blocker_path, []byte("block"), 0600); err != nil {
		t.Fatalf("create destination blocker: %v", err)
	}
	if _, err := db.Exec(`update other set data = ? where name = 'backup_where'`, filepath.Join(blocker_path, "source.db")); err != nil {
		t.Fatalf("change destination: %v", err)
	}
	if err := tool.Backup_run(db, db_path); err == nil {
		t.Fatal("expected malformed destination to fail")
	}

	entries, err = os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read backups after failed attempt: %v", err)
	}
	backup_count = 0
	for _, entry := range entries {
		if pattern.MatchString(entry.Name()) {
			backup_count++
		}
	}
	if backup_count != 2 {
		t.Fatalf("failed attempt damaged valid backups: found %d", backup_count)
	}
}

func Test_Backup_cleanup_temporary_files_removes_only_stale_temporary_files(t *testing.T) {
	dir := t.TempDir()
	old_temp := filepath.Join(dir, "source.backup.tmp1234567890")
	recent_temp := filepath.Join(dir, "source.backup.tmp1234567891")
	valid_backup := filepath.Join(dir, "source_20261030000000.db")
	for _, path := range []string{old_temp, recent_temp, valid_backup} {
		if err := os.WriteFile(path, []byte("fixture"), 0600); err != nil {
			t.Fatalf("create fixture %q: %v", path, err)
		}
	}
	old_time := time.Now().Add(-31 * time.Minute)
	if err := os.Chtimes(old_temp, old_time, old_time); err != nil {
		t.Fatalf("age temporary file: %v", err)
	}

	if err := tool.Backup_cleanup_temporary_files(dir, "source", time.Now()); err != nil {
		t.Fatalf("cleanup stale temporary files: %v", err)
	}
	if _, err := os.Stat(old_temp); !os.IsNotExist(err) {
		t.Fatalf("stale temporary file remains: %v", err)
	}
	for _, path := range []string{recent_temp, valid_backup} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("unexpectedly removed %q: %v", path, err)
		}
	}
}

func Test_Maintenance_daily(t *testing.T) {
	dir := t.TempDir()
	db_path := filepath.Join(dir, "test.db")
	db := new_test_db(t, db_path)
	defer tool.DB_close(db)

	exec := func(query string, args ...any) {
		if _, err := db.Exec(query, args...); err != nil {
			t.Fatalf("exec %s: %v", query, err)
		}
	}

	set_test_setting(db, t, "ua_expiration_date", "1")
	set_test_setting(db, t, "auth_history_expiration_date", "1")
	set_test_setting(db, t, "sitemap_auto_make", "")

	exec(`insert into vote (name, id, subject, data, user, type, acl) values ('', 'vote1', 's', '', '', 'open', '')`)
	exec(`insert into vote (name, id, subject, data, user, type, acl) values ('end_date', 'vote1', '', '2020-01-01', '', 'option', '')`)
	exec(`insert into vote (name, id, subject, data, user, type, acl) values ('', 'vote2', 's2', '', '', 'n_open', '')`)
	exec(`insert into vote (name, id, subject, data, user, type, acl) values ('end_date', 'vote2', '', '9999-01-01', '', 'option', '')`)
	today := time.Now().Format("2006-01-02")
	exec(`insert into vote (name, id, subject, data, user, type, acl) values ('', 'vote3', 's3', '', '', 'open', '')`)
	exec(`insert into vote (name, id, subject, data, user, type, acl) values ('end_date', 'vote3', '', ?, '', 'option', '')`, today)
	exec(`insert into vote (name, id, subject, data, user, type, acl) values ('', 'vote4', 's4', '', '', 'open', '')`)
	exec(`insert into vote (name, id, subject, data, user, type, acl) values ('end_date', 'vote4', '', 'invalid-date', '', 'option', '')`)
	exec(`insert into vote (name, id, subject, data, user, type, acl) values ('', 'vote5', 's5', '', '', 'n_open', '')`)
	exec(`insert into vote (name, id, subject, data, user, type, acl) values ('end_date', 'vote5', '', '2020-01-01', '', 'option', '')`)

	exec(`insert into rb (block, end, today, blocker, why, band, ongoing, login) values ('user1', '2020-01-01', '2026-10-01 00:00:00', 'admin', '', 'private', '1', 'A')`)
	exec(`insert into rb (block, end, today, blocker, why, band, ongoing, login) values ('user2', '9999-01-01', '2026-10-01 00:00:00', 'admin', '', 'private', '1', 'A')`)
	exec(`insert into rb (block, end, today, blocker, why, band, ongoing, login) values ('malformed-ban', '0000-invalid-date', '2026-10-01 00:00:00', 'admin', '', 'private', '1', 'A')`)

	exec(`insert into user_set (name, id, data) values ('acl', 'U1', 'admin')`)
	exec(`insert into user_set (name, id, data) values ('auth_date', 'U1', '2020-01-01')`)
	exec(`insert into user_set (name, id, data) values ('acl', 'U2', 'admin')`)
	exec(`insert into user_set (name, id, data) values ('auth_date', 'U2', '9999-01-01')`)

	exec(`insert into acl (title, data, type) values ('Doc1', 'admin', 'view')`)
	exec(`insert into data_set (doc_name, doc_rev, set_name, set_data) values ('Doc1', 'view', 'acl_date', '2020-01-01')`)

	exec(`insert into ua_d (name, ip, ua, today, sub) values ('ua1', '1.1.1.1', 'agent', '2020-01-01 00:00:00', '')`)
	exec(`insert into re_admin (who, what, time) values ('admin1', 'grant', '2020-01-01 00:00:00')`)
	exec(`insert into re_admin (who, what, time) values ('admin2', 'grant', '9999-01-01 00:00:00')`)

	exec(`insert into user_set (name, id, data) values ('user_title', 'U2', '✅')`)
	exec(`insert into user_set (name, id, data) values ('acl', 'U3', 'test_admin')`)
	exec(`insert into alist (name, acl) values ('test_admin', 'all_admin_auth')`)
	exec(`insert into user_set (name, id, data) values ('user_title', 'U3', '✅')`)

	err := tool.Maintenance_daily(db)
	if err != nil {
		t.Logf("maintenance errors: %v", err)
	}

	if first_string(db, `select type from vote where id = ?`, "vote1") != "close" {
		t.Fatalf("expected vote1 to be closed")
	}
	if first_string(db, `select type from vote where id = ?`, "vote2") != "n_open" {
		t.Fatalf("expected vote2 to remain n_open")
	}
	if first_string(db, `select type from vote where id = ?`, "vote3") != "open" {
		t.Fatalf("expected same-day vote3 to remain open")
	}
	if first_string(db, `select type from vote where id = ?`, "vote4") != "close" {
		t.Fatalf("expected malformed-date vote4 to fail closed")
	}
	if first_string(db, `select type from vote where id = ?`, "vote5") != "n_close" {
		t.Fatalf("expected expired vote5 to be n_close")
	}
	if first_string(db, `select ongoing from rb where block = ?`, "user1") != "" {
		t.Fatalf("expected rb user1 ongoing to be empty")
	}
	if first_string(db, `select ongoing from rb where block = ?`, "user2") != "1" {
		t.Fatalf("expected rb user2 to remain ongoing")
	}
	if first_string(db, `select ongoing from rb where block = ?`, "malformed-ban") != "1" {
		t.Fatalf("malformed restrictive ban was deactivated")
	}
	if got := tool.Get_auth_target_group(db, "user1", "private"); got != "ip" {
		t.Fatalf("expired ban still affects authorization: %q", got)
	}
	if got := tool.Get_auth_target_group(db, "user2", "private"); got != "ban_without_login" {
		t.Fatalf("future ban lost authorization effect: %q", got)
	}
	if got := tool.Get_auth_target_group(db, "malformed-ban", "private"); got != "ban_without_login" {
		t.Fatalf("malformed restrictive ban lost authorization effect: %q", got)
	}
	if first_string(db, `select data from user_set where id = ? and name = ?`, "U1", "acl") != "user" {
		t.Fatalf("expected U1 acl to be user")
	}
	if first_string(db, `select data from user_set where id = ? and name = ?`, "U2", "acl") != "admin" {
		t.Fatalf("expected U2 acl to remain admin")
	}
	if first_int(db, `select count(*) from data_set where doc_name = ? and doc_rev = ? and set_name = ?`, "Doc1", "view", "acl_date") != 0 {
		t.Fatalf("expected doc acl_date row to be removed")
	}
	if first_int(db, `select count(*) from acl where title = ? and type = ?`, "Doc1", "view") != 0 {
		t.Fatalf("expected acl row to be removed")
	}
	if first_int(db, `select count(*) from ua_d`) != 0 {
		t.Fatalf("expected ua_d rows to be removed")
	}
	if first_int(db, `select count(*) from re_admin where time < '9999-01-01'`) != 0 {
		t.Fatalf("expected expired re_admin row to be removed")
	}
	if first_string(db, `select data from user_set where id = ? and name = ?`, "U3", "user_title") != "☑️" {
		t.Fatalf("expected U3 user_title to become checked admin")
	}
}

func Test_Maintenance_daily_closed_database_returns_error(t *testing.T) {
	db := new_test_db(t, filepath.Join(t.TempDir(), "closed.db"))
	if err := db.Close(); err != nil {
		t.Fatalf("close database: %v", err)
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("maintenance panicked after a database failure: %v", recovered)
		}
	}()
	if err := tool.Maintenance_daily(db); err == nil {
		t.Fatal("expected database errors from daily maintenance")
	}
}

func Test_Maintenance_rb_expire_retries_after_sqlite_busy(t *testing.T) {
	db_path := filepath.Join(t.TempDir(), "retry-busy.db")
	db := new_test_db(t, db_path)
	defer tool.DB_close(db)
	if _, err := db.Exec(`insert into rb (block, end, today, blocker, why, band, ongoing, login) values ('expired-ban', '2000-01-01', '1999-01-01', 'admin', '', 'private', '1', 'A')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into rb (block, end, today, blocker, why, band, ongoing, login) values ('lock-holder', '9999-01-01', '2000-01-01', 'admin', '', 'private', '', 'A')`); err != nil {
		t.Fatal(err)
	}

	lock_db, err := sql.Open("sqlite", db_path)
	if err != nil {
		t.Fatalf("open lock connection: %v", err)
	}
	defer lock_db.Close()
	lock_tx, err := lock_db.Begin()
	if err != nil {
		t.Fatalf("begin lock transaction: %v", err)
	}
	if _, err := lock_tx.Exec(`update rb set ongoing = '1' where block = 'lock-holder'`); err != nil {
		_ = lock_tx.Rollback()
		t.Fatalf("acquire SQLite write lock: %v", err)
	}

	busy_err := tool.Maintenance_rb_expire(db)
	if busy_err == nil || !strings.Contains(strings.ToLower(busy_err.Error()), "locked") {
		_ = lock_tx.Rollback()
		t.Fatalf("expected temporary SQLITE_BUSY, got %v", busy_err)
	}
	if got := first_string(db, `select ongoing from rb where block = 'expired-ban'`); got != "1" {
		_ = lock_tx.Rollback()
		t.Fatalf("failed transaction partially changed expired ban: %q", got)
	}
	if err := lock_tx.Commit(); err != nil {
		t.Fatalf("release SQLite write lock: %v", err)
	}
	if err := tool.Maintenance_rb_expire(db); err != nil {
		t.Fatalf("retry after SQLITE_BUSY: %v", err)
	}
	if got := first_string(db, `select ongoing from rb where block = 'expired-ban'`); got != "" {
		t.Fatalf("successful retry left expired ban active: %q", got)
	}
	if err := tool.Maintenance_rb_expire(db); err != nil {
		t.Fatalf("idempotent second retry: %v", err)
	}
}

func Test_Maintenance_vote_close_rolls_back_and_retries(t *testing.T) {
	db := new_test_db(t, filepath.Join(t.TempDir(), "vote-retry.db"))
	defer tool.DB_close(db)
	if _, err := db.Exec(`insert into vote (name, id, subject, data, user, type, acl) values ('', 'retry-vote', 'subject', '', '', 'open', '')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into vote (name, id, subject, data, user, type, acl) values ('end_date', 'retry-vote', '', '2000-01-01', '', 'option', '')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`create trigger fail_vote_close before update on vote when old.id = 'retry-vote' begin select raise(abort, 'temporary maintenance failure'); end`); err != nil {
		t.Fatalf("create transient failure trigger: %v", err)
	}
	today := time.Now().Format("2006-01-02")
	if err := tool.Maintenance_vote_close(db, today); err == nil {
		t.Fatal("expected temporary vote maintenance error")
	}
	if got := first_string(db, `select type from vote where id = 'retry-vote' and user = ''`); got != "open" {
		t.Fatalf("failed vote transaction left partial state %q", got)
	}
	if _, err := db.Exec(`drop trigger fail_vote_close`); err != nil {
		t.Fatalf("remove transient failure trigger: %v", err)
	}
	if err := tool.Maintenance_vote_close(db, today); err != nil {
		t.Fatalf("retry vote maintenance: %v", err)
	}
	if got := first_string(db, `select type from vote where id = 'retry-vote' and user = ''`); got != "close" {
		t.Fatalf("successful vote retry left state %q", got)
	}
	if err := tool.Maintenance_vote_close(db, today); err != nil {
		t.Fatalf("idempotent vote maintenance: %v", err)
	}
}

func Test_Maintenance_invalid_retention_settings_preserve_rows(t *testing.T) {
	db := new_test_db(t, filepath.Join(t.TempDir(), "invalid-retention.db"))
	defer tool.DB_close(db)
	set_test_setting(db, t, "ua_expiration_date", "-1")
	set_test_setting(db, t, "auth_history_expiration_date", "not-an-integer")
	if _, err := db.Exec(`insert into ua_d (name, ip, ua, today, sub) values ('keep-ua', '1.1.1.1', 'agent', '2000-01-01', '')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into re_admin (who, what, time) values ('audit', 'keep-history', '2000-01-01 00:00:00')`); err != nil {
		t.Fatal(err)
	}
	if err := tool.Maintenance_ua_expire(db); err == nil {
		t.Fatal("negative UA retention setting should report invalid configuration")
	}
	if err := tool.Maintenance_auth_history_expire(db); err == nil {
		t.Fatal("malformed history retention setting should report invalid configuration")
	}
	if first_int(db, `select count(*) from ua_d where name = 'keep-ua'`) != 1 {
		t.Fatal("invalid UA retention setting deleted data")
	}
	if first_int(db, `select count(*) from re_admin where what = 'keep-history'`) != 1 {
		t.Fatal("invalid history retention setting deleted data")
	}

	cutoff, enabled, err := tool.Maintenance_expiration_cutoff("0", time.Now())
	if err != nil || !enabled || cutoff.Hour() != 0 {
		t.Fatalf("zero-day legacy retention behavior changed: cutoff=%v enabled=%v err=%v", cutoff, enabled, err)
	}
	if _, enabled, err := tool.Maintenance_expiration_cutoff("", time.Now()); err != nil || enabled {
		t.Fatalf("empty retention setting should remain disabled: enabled=%v err=%v", enabled, err)
	}
}

func Test_Maintenance_daily_turso_backend(t *testing.T) {
	working_dir := t.TempDir()
	t.Chdir(working_dir)
	db := Turso_setup_db(t)
	defer tool.DB_close(db)

	set_test_setting(db, t, "ua_expiration_date", "1")
	set_test_setting(db, t, "auth_history_expiration_date", "1")
	set_test_setting(db, t, "sitemap_auto_make", "1")
	if _, err := db.Exec(`insert into alist (name, acl) values ('turso_admin_title', 'all_admin_auth')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into alist (name, acl) values ('turso_viewer', 'document')`); err != nil {
		t.Fatal(err)
	}
	today := time.Now().Format("2006-01-02")
	for _, vote := range []struct {
		id  string
		end string
	}{
		{id: "turso-expired-vote", end: "2000-01-01"},
		{id: "turso-future-vote", end: "9999-01-01"},
	} {
		if _, err := db.Exec(`insert into vote (name, id, subject, data, user, type, acl) values ('', ?, 'subject', '', '', 'open', '')`, vote.id); err != nil {
			t.Fatal(err)
		}
		if _, err := db.Exec(`insert into vote (name, id, subject, data, user, type, acl) values ('end_date', ?, '', ?, '', 'option', '')`, vote.id, vote.end); err != nil {
			t.Fatal(err)
		}
	}
	for _, ban := range []struct {
		block string
		end   string
	}{
		{block: "turso-expired-ban", end: "2000-01-01"},
		{block: "turso-future-ban", end: "9999-01-01"},
		{block: "turso-malformed-ban", end: "0000-invalid"},
	} {
		if _, err := db.Exec(`insert into rb (block, end, today, blocker, why, band, ongoing, login) values (?, ?, '2026-10-01', 'admin', '', 'private', '1', 'A')`, ban.block, ban.end); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', 'turso-expired-user', 'admin')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('auth_date', 'turso-expired-user', '2000-01-01')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', 'turso-active-modern', 'turso_admin_title')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('auth_date', 'turso-active-modern', '2000-01-01')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl_end', 'turso-active-modern', '9999-01-01')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', 'turso-boundary-user', 'admin')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('auth_date', 'turso-boundary-user', ?)`, today); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('user_title', 'turso-active-modern', '✅')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', 'turso-ordinary-title', 'turso_viewer')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('user_title', 'turso-ordinary-title', '✅')`); err != nil {
		t.Fatal(err)
	}
	for _, doc := range []struct {
		title string
		date  string
	}{
		{title: "TursoExpiredDoc", date: "2000-01-01"},
		{title: "TursoFutureDoc", date: "9999-01-01"},
		{title: "TursoBoundaryDoc", date: today},
	} {
		if _, err := db.Exec(`insert into data (title, data, type) values (?, 'content', 'd')`, doc.title); err != nil {
			t.Fatal(err)
		}
		if _, err := db.Exec(`insert into acl (title, data, type) values (?, 'turso_admin_title', 'view')`, doc.title); err != nil {
			t.Fatal(err)
		}
		if _, err := db.Exec(`insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, 'view', 'acl_date', ?)`, doc.title, doc.date); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', 'turso-viewer', 'turso_viewer')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into ua_d (name, ip, ua, today, sub) values ('turso-old-ua', '1.1.1.1', 'agent', '2000-01-01 00:00:00', '')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into ua_d (name, ip, ua, today, sub) values ('turso-future-ua', '1.1.1.2', 'agent', '9999-01-01 00:00:00', '')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into re_admin (who, what, time) values ('audit', 'turso-old-history', '2000-01-01 00:00:00')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into re_admin (who, what, time) values ('audit', 'turso-future-history', '9999-01-01 00:00:00')`); err != nil {
		t.Fatal(err)
	}
	if err := tool.Maintenance_daily(db); err != nil {
		t.Fatalf("daily maintenance on Turso: %v", err)
	}
	for _, vote := range []struct {
		id   string
		want string
	}{
		{id: "turso-expired-vote", want: "close"},
		{id: "turso-future-vote", want: "open"},
	} {
		if got := first_string(db, `select type from vote where id = ? and name = ''`, vote.id); got != vote.want {
			t.Fatalf("Turso vote %s type = %q, want %q", vote.id, got, vote.want)
		}
	}
	for _, user := range []struct {
		id   string
		want string
	}{
		{id: "turso-expired-user", want: "user"},
		{id: "turso-active-modern", want: "turso_admin_title"},
		{id: "turso-boundary-user", want: "admin"},
	} {
		if got := first_string(db, `select data from user_set where id = ? and name = 'acl'`, user.id); got != user.want {
			t.Fatalf("Turso user %s ACL = %q, want %q", user.id, got, user.want)
		}
	}
	for _, ban := range []struct {
		block string
		want  string
	}{
		{block: "turso-expired-ban", want: ""},
		{block: "turso-future-ban", want: "1"},
		{block: "turso-malformed-ban", want: "1"},
	} {
		if got := first_string(db, `select ongoing from rb where block = ?`, ban.block); got != ban.want {
			t.Fatalf("Turso ban %s ongoing = %q, want %q", ban.block, got, ban.want)
		}
	}
	if first_int(db, `select count(*) from ua_d where name = 'turso-old-ua'`) != 0 || first_int(db, `select count(*) from ua_d where name = 'turso-future-ua'`) != 1 {
		t.Fatal("Turso UA retention result mismatch")
	}
	if first_int(db, `select count(*) from re_admin where what = 'turso-old-history'`) != 0 || first_int(db, `select count(*) from re_admin where what = 'turso-future-history'`) != 1 {
		t.Fatal("Turso auth-history retention result mismatch")
	}
	if first_int(db, `select count(*) from acl where title = 'TursoExpiredDoc'`) != 0 || first_int(db, `select count(*) from acl where title = 'TursoFutureDoc'`) != 1 || first_int(db, `select count(*) from acl where title = 'TursoBoundaryDoc'`) != 1 {
		t.Fatal("Turso document ACL expiration result mismatch")
	}
	if !tool.Check_permission(db, "all_admin_auth", "turso-active-modern") {
		t.Fatal("Turso active modern grant became ineffective")
	}
	if tool.Check_permission(db, "all_admin_auth", "turso-expired-user") {
		t.Fatal("Turso expired legacy grant remained effective")
	}
	if got := first_string(db, `select data from user_set where id = 'turso-active-modern' and name = 'user_title'`); got != "☑️" {
		t.Fatalf("Turso former admin title = %q", got)
	}
	if got := first_string(db, `select data from user_set where id = 'turso-ordinary-title' and name = 'user_title'`); got != "✅" {
		t.Fatalf("Turso ordinary title changed: %q", got)
	}
	for _, doc := range []struct {
		title string
		want  bool
	}{
		{title: "TursoExpiredDoc", want: true},
		{title: "TursoFutureDoc", want: false},
		{title: "TursoBoundaryDoc", want: false},
	} {
		if got := tool.Check_acl(db, doc.title, "", "render", "turso-viewer"); got != doc.want {
			t.Fatalf("Turso viewer access for %s = %v, want %v", doc.title, got, doc.want)
		}
	}
	if _, err := os.Stat("sitemap.xml"); err != nil {
		t.Fatalf("Turso automatic sitemap missing: %v", err)
	}
	if err := tool.Maintenance_daily(db); err != nil {
		t.Fatalf("idempotent second Turso maintenance pass: %v", err)
	}
}

func Test_Maintenance_daily_mysql_backend(t *testing.T) {
	db, _ := mysql_audit_setup(t)
	set_setting := func(name string, value string) {
		t.Helper()
		if _, err := db.Exec(`delete from other where name = ? and coverage = ''`, name); err != nil {
			t.Fatalf("clear MySQL setting %s: %v", name, err)
		}
		if _, err := db.Exec(`insert into other (name, data, coverage) values (?, ?, '')`, name, value); err != nil {
			t.Fatalf("set MySQL setting %s: %v", name, err)
		}
	}
	exec := func(query string, args ...any) {
		t.Helper()
		if _, err := db.Exec(query, args...); err != nil {
			t.Fatalf("MySQL fixture SQL %q: %v", query, err)
		}
	}

	today := time.Now()
	today_text := today.Format("2006-01-02")
	cutoff := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()).Add(-24 * time.Hour).Format("2006-01-02 15:04:05")
	set_setting("ua_expiration_date", "1")
	set_setting("auth_history_expiration_date", "1")
	set_setting("sitemap_auto_make", "1")

	exec(`insert into alist (name, acl) values ('maintenance_admin', 'all_admin_auth')`)
	exec(`insert into alist (name, acl) values ('maintenance_viewer', 'document')`)
	exec(`insert into alist (name, acl) values ('ban', 'ban')`)

	for _, vote := range []struct {
		id      string
		kind    string
		endDate string
	}{
		{id: "expired-open", kind: "open", endDate: "2000-01-01"},
		{id: "expired-n-open", kind: "n_open", endDate: "2000-01-01"},
		{id: "future-open", kind: "open", endDate: "9999-01-01"},
		{id: "malformed-open", kind: "open", endDate: "invalid-date"},
	} {
		exec(`insert into vote (name, id, subject, data, user, type, acl) values ('', ?, 'subject', '', '', ?, '')`, vote.id, vote.kind)
		exec(`insert into vote (name, id, subject, data, user, type, acl) values ('end_date', ?, '', ?, '', 'option', '')`, vote.id, vote.endDate)
	}

	for _, ban := range []struct {
		block string
		end   string
	}{
		{block: "expired-ban", end: "2000-01-01 00:00:00"},
		{block: "future-ban", end: "9999-01-01 00:00:00"},
		{block: "malformed-ban", end: "0000-invalid-date"},
	} {
		exec(`insert into rb (block, end, today, blocker, why, band, login, ongoing) values (?, ?, '2026-10-01 00:00:00', 'admin', '', 'private', 'A', '1')`, ban.block, ban.end)
	}

	for _, user := range []struct {
		id        string
		group     string
		legacyEnd string
		modernEnd string
	}{
		{id: "expired-legacy", group: "maintenance_admin", legacyEnd: "2000-01-01"},
		{id: "malformed-legacy", group: "maintenance_admin", legacyEnd: "not-a-date"},
		{id: "future-legacy", group: "maintenance_admin", legacyEnd: "9999-01-01"},
		{id: "active-modern", group: "maintenance_admin", legacyEnd: "2000-01-01", modernEnd: "9999-01-01"},
		{id: "unlimited-modern", group: "maintenance_admin", legacyEnd: "2000-01-01", modernEnd: ""},
		{id: "expired-modern", group: "maintenance_admin", legacyEnd: "9999-01-01", modernEnd: "2000-01-01"},
		{id: "malformed-modern", group: "maintenance_admin", legacyEnd: "9999-01-01", modernEnd: "invalid-end"},
		{id: "viewer", group: "maintenance_viewer"},
	} {
		exec(`insert into user_set (name, id, data) values ('acl', ?, ?)`, user.id, user.group)
		if user.legacyEnd != "" {
			exec(`insert into user_set (name, id, data) values ('auth_date', ?, ?)`, user.id, user.legacyEnd)
		}
		if user.modernEnd != "" || user.id == "unlimited-modern" {
			exec(`insert into user_set (name, id, data) values ('acl_end', ?, ?)`, user.id, user.modernEnd)
		}
	}

	for _, doc := range []struct {
		title string
		type_ string
		date  string
	}{
		{title: "ExpiredDoc", type_: "view", date: "2000-01-01"},
		{title: "MalformedDoc", type_: "view", date: "invalid-date"},
		{title: "FutureDoc", type_: "view", date: "9999-01-01"},
		{title: "BoundaryDoc", type_: "view", date: today_text},
	} {
		exec(`insert into acl (title, data, type) values (?, 'maintenance_admin', ?)`, doc.title, doc.type_)
		exec(`insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, ?, 'acl_date', ?)`, doc.title, doc.type_, doc.date)
	}
	exec(`insert into acl (title, data, type) values ('UnrelatedDoc', 'maintenance_admin', 'edit')`)
	exec(`insert into data (title, data, type) values ('MariaDBSitemapDoc', 'content', 'd')`)

	exec(`insert into ua_d (name, ip, ua, today, sub) values ('old-ua', '1.1.1.1', 'agent', '2000-01-01 00:00:00', '')`)
	exec(`insert into ua_d (name, ip, ua, today, sub) values ('boundary-ua', '1.1.1.2', 'agent', ?, '')`, cutoff)
	exec(`insert into ua_d (name, ip, ua, today, sub) values ('future-ua', '1.1.1.3', 'agent', '9999-01-01 00:00:00', '')`)
	exec(`insert into re_admin (who, what, time) values ('audit', 'old-history', '2000-01-01 00:00:00')`)
	exec(`insert into re_admin (who, what, time) values ('audit', 'boundary-history', ?)`, cutoff)
	exec(`insert into re_admin (who, what, time) values ('audit', 'future-history', '9999-01-01 00:00:00')`)

	exec(`insert into user_set (name, id, data) values ('acl', 'former-admin-title', 'maintenance_admin')`)
	exec(`insert into user_set (name, id, data) values ('user_title', 'former-admin-title', '✅')`)
	exec(`insert into user_set (name, id, data) values ('acl', 'ordinary-title', 'maintenance_viewer')`)
	exec(`insert into user_set (name, id, data) values ('user_title', 'ordinary-title', '✅')`)

	for _, user := range []struct {
		id   string
		want bool
	}{
		{id: "expired-legacy", want: false},
		{id: "malformed-legacy", want: false},
		{id: "future-legacy", want: true},
		{id: "active-modern", want: true},
		{id: "unlimited-modern", want: true},
		{id: "expired-modern", want: false},
		{id: "malformed-modern", want: false},
	} {
		if got := tool.Check_permission(db, "all_admin_auth", user.id); got != user.want {
			t.Fatalf("MariaDB effective permission before maintenance for %s = %v, want %v", user.id, got, user.want)
		}
	}
	for _, title := range []string{"ExpiredDoc", "MalformedDoc", "FutureDoc", "BoundaryDoc"} {
		if tool.Check_acl(db, title, "", "render", "viewer") {
			t.Fatalf("expected ACL to restrict viewer before maintenance: %s", title)
		}
	}

	if err := tool.Maintenance_daily(db); err != nil {
		t.Fatalf("daily maintenance on MariaDB: %v", err)
	}

	for _, vote := range []struct {
		id   string
		want string
	}{
		{id: "expired-open", want: "close"},
		{id: "expired-n-open", want: "n_close"},
		{id: "malformed-open", want: "close"},
		{id: "future-open", want: "open"},
	} {
		if got := first_string(db, `select type from vote where id = ? and name = '' limit 1`, vote.id); got != vote.want {
			t.Fatalf("MariaDB vote %s type = %q, want %q", vote.id, got, vote.want)
		}
	}
	for _, ban := range []struct {
		block string
		want  string
	}{
		{block: "expired-ban", want: ""},
		{block: "future-ban", want: "1"},
		{block: "malformed-ban", want: "1"},
	} {
		if got := first_string(db, `select ongoing from rb where block = ?`, ban.block); got != ban.want {
			t.Fatalf("MariaDB ban %s ongoing = %q, want %q", ban.block, got, ban.want)
		}
	}
	if got := tool.Get_auth_target_group(db, "expired-ban", "private"); got != "ip" {
		t.Fatalf("expired MariaDB ban still affects authorization: %q", got)
	}
	if got := tool.Get_auth_target_group(db, "malformed-ban", "private"); got != "ban_without_login" {
		t.Fatalf("malformed restrictive MariaDB ban lost authorization effect: %q", got)
	}
	for _, user := range []struct {
		id   string
		want bool
	}{
		{id: "expired-legacy", want: false},
		{id: "malformed-legacy", want: false},
		{id: "future-legacy", want: true},
		{id: "active-modern", want: true},
		{id: "unlimited-modern", want: true},
		{id: "expired-modern", want: false},
		{id: "malformed-modern", want: false},
	} {
		if got := tool.Check_permission(db, "all_admin_auth", user.id); got != user.want {
			t.Fatalf("MariaDB effective permission for %s = %v, want %v", user.id, got, user.want)
		}
	}
	for _, user_id := range []string{"expired-legacy", "malformed-legacy"} {
		if first_int(db, `select count(*) from user_set where id = ? and name = 'auth_date'`, user_id) != 0 {
			t.Fatalf("MariaDB expired auth_date remained for %s", user_id)
		}
	}
	for _, doc := range []struct {
		title string
		want  bool
	}{
		{title: "ExpiredDoc", want: true},
		{title: "MalformedDoc", want: true},
		{title: "FutureDoc", want: false},
		{title: "BoundaryDoc", want: false},
	} {
		if got := tool.Check_acl(db, doc.title, "", "render", "viewer"); got != doc.want {
			t.Fatalf("MariaDB viewer access for %s = %v, want %v", doc.title, got, doc.want)
		}
	}
	if first_int(db, `select count(*) from acl where title = 'UnrelatedDoc' and type = 'edit'`) != 1 {
		t.Fatal("MariaDB maintenance removed unrelated document ACL")
	}
	for _, row := range []struct {
		name string
		want int
	}{
		{name: "old-ua", want: 0},
		{name: "boundary-ua", want: 1},
		{name: "future-ua", want: 1},
	} {
		if got := first_int(db, `select count(*) from ua_d where name = ?`, row.name); got != row.want {
			t.Fatalf("MariaDB UA retention %s count = %d, want %d", row.name, got, row.want)
		}
	}
	for _, row := range []struct {
		what string
		want int
	}{
		{what: "old-history", want: 0},
		{what: "boundary-history", want: 1},
		{what: "future-history", want: 1},
	} {
		if got := first_int(db, `select count(*) from re_admin where what = ?`, row.what); got != row.want {
			t.Fatalf("MariaDB auth history %s count = %d, want %d", row.what, got, row.want)
		}
	}
	if first_string(db, `select data from user_set where id = 'former-admin-title' and name = 'user_title'`) != "☑️" {
		t.Fatal("MariaDB former admin title was not updated")
	}
	if first_string(db, `select data from user_set where id = 'ordinary-title' and name = 'user_title'`) != "✅" {
		t.Fatal("MariaDB ordinary user title was changed")
	}
	if _, err := os.Stat("sitemap.xml"); err != nil {
		t.Fatalf("MariaDB automatic sitemap missing: %v", err)
	}
	if _, err := os.Stat("sitemap_0.xml"); err != nil {
		t.Fatalf("MariaDB automatic sitemap shard missing: %v", err)
	}
	if err := tool.Maintenance_daily(db); err != nil {
		t.Fatalf("second MariaDB maintenance pass: %v", err)
	}
	if first_string(db, `select type from vote where id = 'expired-open' and name = '' limit 1`) != "close" {
		t.Fatal("repeated MariaDB maintenance changed expired vote state")
	}

	exec(`insert into acl (title, data, type) values ('RollbackDoc', 'maintenance_admin', 'view')`)
	exec(`insert into data_set (doc_name, doc_rev, set_name, set_data) values ('RollbackDoc', 'view', 'acl_date', '2000-01-01')`)
	if _, err := db.Exec(`CREATE TRIGGER maintenance_acl_reject BEFORE DELETE ON acl FOR EACH ROW SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'audit rollback'`); err != nil {
		t.Fatalf("create rollback trigger: %v", err)
	}
	if err := tool.Maintenance_document_acl_expire(db, today_text); err == nil {
		t.Fatal("expected MariaDB document ACL delete error")
	}
	if first_int(db, `select count(*) from acl where title = 'RollbackDoc' and type = 'view'`) != 1 || first_int(db, `select count(*) from data_set where doc_name = 'RollbackDoc' and set_name = 'acl_date'`) != 1 {
		t.Fatal("failed MariaDB ACL expiration did not roll back both writes")
	}
	if _, err := db.Exec(`DROP TRIGGER maintenance_acl_reject`); err != nil {
		t.Fatalf("drop rollback trigger: %v", err)
	}
	if err := tool.Maintenance_document_acl_expire(db, today_text); err != nil {
		t.Fatalf("retry MariaDB ACL expiration: %v", err)
	}
	if first_int(db, `select count(*) from acl where title = 'RollbackDoc'`) != 0 || first_int(db, `select count(*) from data_set where doc_name = 'RollbackDoc' and set_name = 'acl_date'`) != 0 {
		t.Fatal("successful MariaDB ACL retry did not remove both rows")
	}
}

func Test_Maintenance_mysql_connection_failure_returns_retryable_error(t *testing.T) {
	db, db_name := mysql_audit_setup(t)
	server_dsn := os.Getenv("MYSQL_AUDIT_DSN")
	credential_end := strings.Index(server_dsn, "@")
	if credential_end <= 0 {
		t.Fatal("test DSN does not contain a server-level credential prefix")
	}
	unavailable_dsn := server_dsn[:credential_end] + "@tcp(127.0.0.1:1)/" + db_name
	unavailable_db, err := sql.Open("mysql", unavailable_dsn)
	if err != nil {
		t.Fatalf("open unavailable MySQL pool: %v", err)
	}
	defer unavailable_db.Close()

	maintenance_err := tool.Run_auto_scheduled_task(func() error {
		return tool.Maintenance_daily(unavailable_db)
	})
	if maintenance_err == nil {
		t.Fatal("unreachable MySQL endpoint should return a retryable maintenance error")
	}
	if err := tool.Maintenance_daily(db); err != nil {
		t.Fatalf("maintenance did not recover on the valid MySQL connection: %v", err)
	}
}

func Test_scheduler_mysql_startup_runs_daily_maintenance(t *testing.T) {
	db, db_name := mysql_audit_setup(t)
	set_setting := func(name string, value string) {
		t.Helper()
		if _, err := db.Exec(`delete from other where name = ? and coverage = ''`, name); err != nil {
			t.Fatalf("clear MySQL setting %s: %v", name, err)
		}
		if _, err := db.Exec(`insert into other (name, data, coverage) values (?, ?, '')`, name, value); err != nil {
			t.Fatalf("set MySQL setting %s: %v", name, err)
		}
	}
	set_setting("sitemap_auto_make", "")
	set_setting("ua_expiration_date", "")
	set_setting("auth_history_expiration_date", "")
	if _, err := db.Exec(`insert into alist (name, acl) values ('scheduler_admin', 'all_admin_auth')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', 'scheduler-mysql-user', 'scheduler_admin')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('auth_date', 'scheduler-mysql-user', '2000-01-01')`); err != nil {
		t.Fatal(err)
	}
	if got := first_string(db, `select data from user_set where id = 'scheduler-mysql-user' and name = 'acl'`); got != "scheduler_admin" {
		t.Fatalf("seeded MySQL ACL = %q, want scheduler_admin", got)
	}
	if first_int(db, `select count(*) from user_set where id = 'scheduler-mysql-user' and name = 'auth_date'`) != 1 {
		t.Fatal("seeded MySQL legacy expiration row missing")
	}
	working_dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("read isolated MySQL working directory: %v", err)
	}

	var wait_group sync.WaitGroup
	outputs := make([][]byte, 2)
	errors_list := make([]error, 2)
	for index := range outputs {
		wait_group.Add(1)
		go func(index int) {
			defer wait_group.Done()
			command := exec.Command(os.Args[0], "-test.run=^Test_scheduler_mysql_child_process$")
			command.Dir = working_dir
			command.Env = []string{
				"NAMU_DB=" + db_name,
				"NAMU_DB_TYPE=mysql",
				"OPENNAMU_AUDIT_MYSQL_SCHEDULER_CHILD=1",
			}
			outputs[index], errors_list[index] = command.CombinedOutput()
		}(index)
	}
	wait_group.Wait()
	for index, err := range errors_list {
		if err != nil {
			t.Fatalf("MariaDB scheduler child %d failed: %v\n%s", index, err, outputs[index])
		}
	}
	state_path := filepath.Join(working_dir, ".opennamu_auto_maintenance_last_run")
	state_before_staggered_start, err := os.ReadFile(state_path)
	if err != nil {
		t.Fatalf("MariaDB scheduler state missing after simultaneous start: %v", err)
	}
	staggered_command := exec.Command(os.Args[0], "-test.run=^Test_scheduler_mysql_child_process$")
	staggered_command.Dir = working_dir
	staggered_command.Env = []string{
		"NAMU_DB=" + db_name,
		"NAMU_DB_TYPE=mysql",
		"OPENNAMU_AUDIT_MYSQL_SCHEDULER_CHILD=1",
	}
	if output, err := staggered_command.CombinedOutput(); err != nil {
		t.Fatalf("staggered MariaDB scheduler child failed: %v\n%s", err, output)
	}
	state_after_staggered_start, err := os.ReadFile(state_path)
	if err != nil {
		t.Fatalf("read MariaDB scheduler state after staggered start: %v", err)
	}
	if string(state_after_staggered_start) != string(state_before_staggered_start) {
		t.Fatal("staggered MariaDB process repeated successful maintenance within the 24-hour interval")
	}
	if tool.Check_permission(db, "all_admin_auth", "scheduler-mysql-user") {
		t.Fatal("MySQL scheduler left the expired legacy grant effective")
	}
	if got := first_string(db, `select data from user_set where id = 'scheduler-mysql-user' and name = 'acl'`); got != "user" {
		t.Fatalf("MySQL scheduler ACL = %q, want user", got)
	}
	if _, err := os.Stat(state_path); err != nil {
		t.Fatalf("MySQL scheduler completion state missing: %v", err)
	}
}

func Test_scheduler_mysql_connection_failure_records_bounded_retry(t *testing.T) {
	dir := t.TempDir()
	data_dir := filepath.Join(dir, "data")
	if err := os.Mkdir(data_dir, 0700); err != nil {
		t.Fatalf("create MySQL config directory: %v", err)
	}
	mysql_config := `{"user":"audit-test","password":"dummy","host":"127.0.0.1","port":"1"}`
	if err := os.WriteFile(filepath.Join(data_dir, "mysql.json"), []byte(mysql_config), 0600); err != nil {
		t.Fatalf("write unavailable MySQL config: %v", err)
	}
	db_name := "maintenance_unavailable"
	command := exec.Command(os.Args[0], "-test.run=^Test_scheduler_unavailable_mysql_child_process$")
	command.Dir = dir
	command.Env = []string{
		"NAMU_DB=" + db_name,
		"NAMU_DB_TYPE=mysql",
		"OPENNAMU_AUDIT_MYSQL_FAILURE_CHILD=1",
	}
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("unavailable MySQL scheduler child panicked: %v\n%s", err, output)
	}
	last_success_path := filepath.Join(dir, ".opennamu_auto_maintenance_last_run")
	if _, err := os.Stat(last_success_path); !os.IsNotExist(err) {
		t.Fatalf("unavailable MySQL attempt recorded success: %v", err)
	}
	retry_after_path := filepath.Join(dir, ".opennamu_auto_maintenance_retry_after")
	retry_after, exists, err := tool.Read_auto_run_time(retry_after_path)
	if err != nil || !exists {
		t.Fatalf("unavailable MySQL retry state missing: exists=%v err=%v", exists, err)
	}
	if tool.Auto_task_is_due(last_success_path, retry_after_path, 24*time.Hour, time.Now().Add(4*time.Minute)) {
		t.Fatal("unavailable MySQL maintenance retried before backoff")
	}
	if !tool.Auto_task_is_due(last_success_path, retry_after_path, 24*time.Hour, retry_after) {
		t.Fatal("unavailable MySQL maintenance was not due at retry deadline")
	}
}

func Test_Maintenance_expiry_fail_closed_and_preserve_active_modern_grant(t *testing.T) {
	dir := t.TempDir()
	db_path := filepath.Join(dir, "test.db")
	db := new_test_db(t, db_path)
	defer tool.DB_close(db)

	today := time.Now().Format("2006-01-02")
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', 'malformed-auth', 'admin')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('auth_date', 'malformed-auth', 'not-a-date')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', 'malformed-modern', 'admin')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('auth_date', 'malformed-modern', '2020-01-01')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl_end', 'malformed-modern', 'invalid-modern-date')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', 'malformed-modern-ban', 'ban')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('auth_date', 'malformed-modern-ban', '2020-01-01')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl_end', 'malformed-modern-ban', 'invalid-ban-date')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', 'active-modern', 'admin')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('auth_date', 'active-modern', '2020-01-01')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl_end', 'active-modern', '9999-01-01')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', 'active-modern-zero', 'admin')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('auth_date', 'active-modern-zero', '2020-01-01')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl_end', 'active-modern-zero', '0')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', 'active-modern-unlimited', 'admin')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('auth_date', 'active-modern-unlimited', '2020-01-01')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl_end', 'active-modern-unlimited', '')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', 'boundary-auth', 'admin')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('auth_date', 'boundary-auth', ?)`, today); err != nil {
		t.Fatal(err)
	}

	if _, err := db.Exec(`insert into acl (title, data, type) values ('MalformedDoc', 'admin', 'view')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into data_set (doc_name, doc_rev, set_name, set_data) values ('MalformedDoc', 'view', 'acl_date', 'invalid-date')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into acl (title, data, type) values ('BoundaryDoc', 'admin', 'view')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into data_set (doc_name, doc_rev, set_name, set_data) values ('BoundaryDoc', 'view', 'acl_date', ?)`, today); err != nil {
		t.Fatal(err)
	}

	if err := tool.Maintenance_auth_date_expire(db, today); err != nil {
		t.Fatalf("auth date expiration: %v", err)
	}
	if err := tool.Maintenance_document_acl_expire(db, today); err != nil {
		t.Fatalf("document ACL expiration: %v", err)
	}

	for _, id := range []string{"malformed-auth", "malformed-modern"} {
		if got := first_string(db, `select data from user_set where id = ? and name = 'acl'`, id); got != "user" {
			t.Fatalf("malformed expired grant %q retained ACL %q", id, got)
		}
	}
	if got := first_string(db, `select data from user_set where id = 'malformed-modern-ban' and name = 'acl'`); got != "ban" {
		t.Fatalf("malformed restrictive ban end was revoked: %q", got)
	}
	if got := first_string(db, `select data from user_set where id = 'active-modern' and name = 'acl'`); got != "admin" {
		t.Fatalf("active modern grant was revoked: %q", got)
	}
	if got := first_string(db, `select data from user_set where id = 'active-modern-unlimited' and name = 'acl'`); got != "admin" {
		t.Fatalf("unlimited modern grant was revoked: %q", got)
	}
	if got := first_string(db, `select data from user_set where id = 'active-modern-zero' and name = 'acl'`); got != "admin" {
		t.Fatalf("zero-expiry modern grant was revoked: %q", got)
	}
	if got := first_string(db, `select data from user_set where id = 'boundary-auth' and name = 'acl'`); got != "admin" {
		t.Fatalf("same-day legacy grant expired at boundary: %q", got)
	}
	if first_int(db, `select count(*) from acl where title = 'MalformedDoc'`) != 0 {
		t.Fatal("malformed temporary document ACL was not removed")
	}
	if first_int(db, `select count(*) from acl where title = 'BoundaryDoc'`) != 1 {
		t.Fatal("same-day temporary document ACL expired at boundary")
	}
}

func Test_Maintenance_expiration_changes_effective_authorization(t *testing.T) {
	db := new_test_db(t, filepath.Join(t.TempDir(), "effective-auth.db"))
	defer tool.DB_close(db)
	if _, err := db.Exec(`insert into alist (name, acl) values ('maintenance_admin', 'all_admin_auth')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into alist (name, acl) values ('maintenance_viewer', 'document')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`create table html_filter (html, kind, plus, plus_t)`); err != nil {
		t.Fatal(err)
	}
	today := time.Now().Format("2006-01-02")

	for _, user := range []struct {
		id        string
		legacyEnd string
		modernEnd string
	}{
		{id: "legacy-expired", legacyEnd: "2000-01-01"},
		{id: "legacy-malformed", legacyEnd: "bad-expiration"},
		{id: "legacy-active", legacyEnd: "9999-01-01"},
		{id: "modern-active", legacyEnd: "2000-01-01", modernEnd: "9999-01-01"},
		{id: "modern-unlimited", legacyEnd: "2000-01-01", modernEnd: ""},
		{id: "modern-expired", legacyEnd: "9999-01-01", modernEnd: "2000-01-01"},
		{id: "modern-malformed", legacyEnd: "9999-01-01", modernEnd: "bad-expiration"},
		{id: "legacy-boundary", legacyEnd: time.Now().Format("2006-01-02")},
	} {
		if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', ?, 'maintenance_admin')`, user.id); err != nil {
			t.Fatal(err)
		}
		if _, err := db.Exec(`insert into user_set (name, id, data) values ('auth_date', ?, ?)`, user.id, user.legacyEnd); err != nil {
			t.Fatal(err)
		}
		if user.modernEnd != "" {
			if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl_end', ?, ?)`, user.id, user.modernEnd); err != nil {
				t.Fatal(err)
			}
		} else if user.id == "modern-unlimited" {
			if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl_end', ?, '')`, user.id); err != nil {
				t.Fatal(err)
			}
		}
	}
	for _, user := range []struct {
		id   string
		want bool
	}{
		{id: "legacy-expired", want: false},
		{id: "legacy-malformed", want: false},
		{id: "legacy-active", want: true},
		{id: "modern-active", want: true},
		{id: "modern-unlimited", want: true},
		{id: "modern-expired", want: false},
		{id: "modern-malformed", want: false},
		{id: "legacy-boundary", want: true},
	} {
		if got := tool.Check_permission(db, "all_admin_auth", user.id); got != user.want {
			t.Fatalf("effective permission before cleanup for %s = %v, want %v", user.id, got, user.want)
		}
	}

	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', 'effective-viewer', 'maintenance_viewer')`); err != nil {
		t.Fatal(err)
	}
	for _, doc := range []struct {
		title string
		date  string
	}{
		{title: "ExpiredAclDoc", date: "2000-01-01"},
		{title: "MalformedAclDoc", date: "not-a-date"},
		{title: "FutureAclDoc", date: "9999-01-01"},
		{title: "BoundaryAclDoc", date: today},
	} {
		if _, err := db.Exec(`insert into acl (title, data, type) values (?, 'maintenance_admin', 'view')`, doc.title); err != nil {
			t.Fatal(err)
		}
		if _, err := db.Exec(`insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, 'view', 'acl_date', ?)`, doc.title, doc.date); err != nil {
			t.Fatal(err)
		}
		if tool.Check_acl(db, doc.title, "", "render", "effective-viewer") {
			t.Fatalf("viewer unexpectedly accessed restricted document before expiration: %s", doc.title)
		}
	}

	if err := tool.Maintenance_auth_date_expire(db, today); err != nil {
		t.Fatalf("expire legacy grants: %v", err)
	}
	if err := tool.Maintenance_document_acl_expire(db, today); err != nil {
		t.Fatalf("expire document ACLs: %v", err)
	}
	for _, user := range []struct {
		id   string
		want bool
	}{
		{id: "legacy-expired", want: false},
		{id: "legacy-malformed", want: false},
		{id: "legacy-active", want: true},
		{id: "modern-active", want: true},
		{id: "modern-unlimited", want: true},
		{id: "modern-expired", want: false},
		{id: "modern-malformed", want: false},
		{id: "legacy-boundary", want: true},
	} {
		if got := tool.Check_permission(db, "all_admin_auth", user.id); got != user.want {
			t.Fatalf("effective permission after cleanup for %s = %v, want %v", user.id, got, user.want)
		}
	}
	for _, doc := range []struct {
		title string
		want  bool
	}{
		{title: "ExpiredAclDoc", want: true},
		{title: "MalformedAclDoc", want: true},
		{title: "FutureAclDoc", want: false},
		{title: "BoundaryAclDoc", want: false},
	} {
		if got := tool.Check_acl(db, doc.title, "", "render", "effective-viewer"); got != doc.want {
			t.Fatalf("effective viewer access after cleanup for %s = %v, want %v", doc.title, got, doc.want)
		}
	}
}

func Test_Expired_legacy_auth_date_is_ineffective_before_daily_cleanup(t *testing.T) {
	db := new_test_db(t, filepath.Join(t.TempDir(), "legacy-auth-effective.db"))
	defer tool.DB_close(db)
	if _, err := db.Exec(`insert into alist (name, acl) values ('maintenance_admin', 'all_admin_auth')`); err != nil {
		t.Fatal(err)
	}
	for _, user := range []struct {
		id        string
		legacyEnd string
		modernEnd string
	}{
		{id: "already-expired", legacyEnd: "2000-01-01"},
		{id: "malformed-expiry", legacyEnd: "not-a-date"},
		{id: "still-valid", legacyEnd: "9999-01-01"},
		{id: "modern-extension", legacyEnd: "2000-01-01", modernEnd: "9999-01-01"},
	} {
		if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', ?, 'maintenance_admin')`, user.id); err != nil {
			t.Fatal(err)
		}
		if _, err := db.Exec(`insert into user_set (name, id, data) values ('auth_date', ?, ?)`, user.id, user.legacyEnd); err != nil {
			t.Fatal(err)
		}
		if user.modernEnd != "" {
			if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl_end', ?, ?)`, user.id, user.modernEnd); err != nil {
				t.Fatal(err)
			}
		}
	}

	for _, user := range []struct {
		id   string
		want bool
	}{
		{id: "already-expired", want: false},
		{id: "malformed-expiry", want: false},
		{id: "still-valid", want: true},
		{id: "modern-extension", want: true},
	} {
		if got := tool.Check_permission(db, "all_admin_auth", user.id); got != user.want {
			t.Fatalf("effective legacy/modern permission for %s = %v, want %v before daily cleanup", user.id, got, user.want)
		}
	}
}

func Test_Check_permission_database_failure_does_not_grant(t *testing.T) {
	db := new_test_db(t, filepath.Join(t.TempDir(), "permission-db-failure.db"))
	if err := db.Close(); err != nil {
		t.Fatalf("close database: %v", err)
	}
	granted := false
	panicked := false
	func() {
		defer func() {
			if recover() != nil {
				panicked = true
			}
		}()
		granted = tool.Check_permission(db, "all_admin_auth", "permission-user")
	}()
	if granted {
		t.Fatal("database lookup failure granted elevated permission")
	}
	if !panicked {
		t.Fatal("expected the database error to surface instead of returning an authorization decision")
	}
}

func Test_Maintenance_malformed_document_acl_expiration_fails_closed(t *testing.T) {
	dir := t.TempDir()
	db := new_test_db(t, filepath.Join(dir, "test.db"))
	defer tool.DB_close(db)

	if _, err := db.Exec(`insert into acl (title, data, type) values ('MalformedDoc', 'admin', 'view')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into data_set (doc_name, doc_rev, set_name, set_data) values ('MalformedDoc', 'view', 'acl_date', 'invalid-date')`); err != nil {
		t.Fatal(err)
	}
	if err := tool.Maintenance_document_acl_expire(db, time.Now().Format("2006-01-02")); err != nil {
		t.Fatalf("document ACL expiration: %v", err)
	}
	if first_int(db, `select count(*) from acl where title = 'MalformedDoc'`) != 0 {
		t.Fatal("malformed temporary document ACL was not removed")
	}
}

func Test_Auth_malformed_grant_end_fails_closed(t *testing.T) {
	dir := t.TempDir()
	db := new_test_db(t, filepath.Join(dir, "test.db"))
	defer tool.DB_close(db)
	today := time.Now().Format("2006-01-02")
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	if tool.Auth_permission_end_active(today) {
		t.Fatal("date-only modern end equal to today should follow existing midnight-expiry semantics")
	}
	if !tool.Auth_permission_end_active(tomorrow) {
		t.Fatal("date-only modern end tomorrow should remain active")
	}

	if _, err := db.Exec(`insert into alist (name, acl) values ('admin', 'all_admin_auth')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into alist (name, acl) values ('ban', '')`); err != nil {
		t.Fatal(err)
	}
	for _, grant := range []struct {
		id    string
		end   string
		group string
	}{
		{id: "malformed", end: "not-a-date"},
		{id: "expired", end: "2000-01-01"},
		{id: "unlimited-empty", end: ""},
		{id: "unlimited-zero", end: "0"},
		{id: "malformed-ban", end: "not-a-date", group: "ban"},
	} {
		group := grant.group
		if group == "" {
			group = "admin"
		}
		if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', ?, ?)`, grant.id, group); err != nil {
			t.Fatal(err)
		}
		if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl_end', ?, ?)`, grant.id, grant.end); err != nil {
			t.Fatal(err)
		}
	}
	for _, id := range []string{"malformed", "expired"} {
		if got := tool.Get_user_auth_raw(db, id); got != "user" {
			t.Fatalf("grant %q with invalid/expired end remained %q", id, got)
		}
	}
	for _, id := range []string{"unlimited-empty", "unlimited-zero"} {
		if got := tool.Get_user_auth_raw(db, id); got != "admin" {
			t.Fatalf("unlimited grant %q was rejected: %q", id, got)
		}
	}
	if got := tool.Get_user_auth_raw(db, "malformed-ban"); got != "ban" {
		t.Fatalf("malformed restrictive ban end was rejected: %q", got)
	}
}

func Test_Maintenance_sitemap_auto(t *testing.T) {
	t.Chdir(t.TempDir())

	dir := t.TempDir()
	t.Chdir(dir)

	db_path := filepath.Join(dir, "test.db")
	db := new_test_db(t, db_path)
	defer tool.DB_close(db)

	exec := func(query string, args ...any) {
		if _, err := db.Exec(query, args...); err != nil {
			t.Fatalf("exec %s: %v", query, err)
		}
	}

	set_test_setting(db, t, "sitemap_auto_make", "1")
	set_test_setting(db, t, "ua_expiration_date", "")
	set_test_setting(db, t, "auth_history_expiration_date", "")

	exec(`insert into data (title, data, type) values ('Doc', 'content', 'd')`)

	if err := tool.Maintenance_daily(db); err != nil {
		t.Fatalf("maintenance: %v", err)
	}

	_, err := os.Stat("sitemap.xml")
	if err != nil {
		t.Fatalf("sitemap.xml missing: %v", err)
	}

	_, err = os.Stat("sitemap_0.xml")
	if err != nil {
		t.Fatalf("sitemap_0.xml missing: %v", err)
	}
}

func Test_Api_setting_sitemap_post_preserves_file_error(t *testing.T) {
	dir := t.TempDir()
	db_name := filepath.Join(dir, "sitemap-api")
	db_path := db_name + ".db"
	db := new_test_db(t, db_path)
	if _, err := db.Exec(`insert into alist (name, acl) values ('admin', 'setting_sitemap')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', 'audit-admin', 'admin')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`insert into data (title, data, type) values ('Doc', 'content', 'd')`); err != nil {
		t.Fatal(err)
	}
	tool.DB_close(db)

	t.Cleanup(func() { tool.DB_boot() })
	t.Setenv("NAMU_DB", db_name)
	t.Setenv("NAMU_DB_TYPE", "sqlite")
	tool.DB_boot()
	if err := os.Mkdir(filepath.Join(dir, "sitemap.xml"), 0700); err != nil {
		t.Fatalf("create write-failure fixture: %v", err)
	}
	t.Chdir(dir)

	result := route.Api_setting_sitemap_post(tool.Config{IP: "audit-admin"})
	if result["response"] != "error" {
		t.Fatalf("expected sitemap API error, got %#v", result)
	}
	if result["error"] != "file write error : sitemap.xml" {
		t.Fatalf("sitemap API error contract changed: %#v", result["error"])
	}
}
