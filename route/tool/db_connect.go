package tool

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	sqlite "modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
	_ "turso.tech/database/tursogo"
)

var db_set = map[string]string{}

// keeper* holds one long-lived physical connection ("anchor") so per-request
// DB_connect()/DB_close() cycles never close the *last* connection:
//   - SQLite: prevents the WAL+shm cleanup that yields SQLITE_BUSY during a
//     concurrent open (unchanged behavior);
//   - Turso (local embedded MVCC): keeps the local DB instance loaded, so the
//     first request after startup or after an idle period attaches to the
//     in-process instance instead of paying the full per-connection load.
// It holds no transaction, cursor, or write lock (idle connection only);
// keyed by db_type + DB path so it re-binds when either changes instead of
// pinning every temp DB. MySQL is a no-op by design.
var (
	keeperMutex sync.Mutex
	keeperDB    *sql.DB
	keeperConn  *sql.Conn
	keeperKey   string
)

// SQLite busy-retry bounds. A bounded, finite retry replaces the previous
// unbounded retry loop so a permanently locked database surfaces instead of
// hanging the process.
const (
	busyRetryDelay = 10 * time.Millisecond
	busyRetryMax   = 10
)

var busyRetryMaxWait = 30 * time.Second

// busyRetryable reports whether a returned error is a SQLite write-lock busy
// that can be safely retried. Only SQLITE_BUSY is retried; SQLITE_BUSY_SNAPSHOT
// (snapshot read while another transaction holds a write lock) is never
// retried inside the same transaction, and non-busy errors are never retried.
// The error code is the authoritative check; a string fallback covers
// non-sqlite errors.
func busyRetryable(err error) bool {
	if err == nil {
		return false
	}
	var e *sqlite.Error
	if errors.As(err, &e) {
		return e.Code() == sqlite3.SQLITE_BUSY
	}
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "snapshot") {
		return false
	}
	return strings.Contains(msg, "database is locked")
}

type DB_runner interface {
	Prepare(string) (*sql.Stmt, error)
}

func DB_transaction(db *sql.DB, callback func(*sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	if err := callback(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	committed = true
	return nil
}

func Get_DB_set() map[string]string {
	new_db_set := map[string]string{}

	db_env, db_env_exist := os.LookupEnv("NAMU_DB")
	db_env_type, db_env_type_exist := os.LookupEnv("NAMU_DB_TYPE")
	if db_env_exist || db_env_type_exist {
		new_db_set["db_name"] = Choose(db_env, "data")
		new_db_set["db_type"] = Choose(db_env_type, "sqlite")

		return new_db_set
	}

	path_dir := filepath.Join("data", "set.json")
	if File_exist_check(path_dir) {
		raw, err := os.ReadFile(path_dir)
		if err == nil {
			tmp := map[string]string{}
			if err := json.Unmarshal(raw, &tmp); err == nil {
				if v, ok := tmp["db"]; ok {
					new_db_set["db_name"] = v
				} else {
					new_db_set["db_name"] = "data"
				}

				if v, ok := tmp["db_type"]; ok {
					new_db_set["db_type"] = v
				} else {
					new_db_set["db_type"] = "sqlite"
				}

				return new_db_set
			}
		}
	}

	new_db_set["db_name"] = "data"
	new_db_set["db_type"] = "sqlite"

	return new_db_set
}

func Get_DB_set_MySQL(new_db_set map[string]string) map[string]string {
	path := filepath.Join("data", "mysql.json")
	if !File_exist_check(path) {
		return new_db_set
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return new_db_set
	}

	tmp := new_db_set
	if err := json.Unmarshal(raw, &tmp); err != nil {
		return tmp
	}

	if host, ok := tmp["host"]; ok && host != "" {
		tmp["db_mysql_host"] = host
	} else {
		tmp["db_mysql_host"] = "127.0.0.1"
	}

	if port, ok := tmp["port"]; ok && port != "" {
		tmp["db_mysql_port"] = port
	} else {
		tmp["db_mysql_port"] = "3306"
	}

	if user, ok := tmp["user"]; ok {
		tmp["db_mysql_user"] = user
	}

	if pw, ok := tmp["password"]; ok {
		tmp["db_mysql_pw"] = pw
	}

	return tmp
}

func Exec_DB_internal(db DB_runner, query string, values ...any) {
	stmt, err := db.Prepare(DB_change(query))
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	attempt := 0
	start := time.Now()
	for {
		_, execErr := stmt.Exec(values...)
		if execErr == nil {
			return
		}

		if !busyRetryable(execErr) ||
			attempt >= busyRetryMax ||
			time.Since(start) >= busyRetryMaxWait {
			panic(execErr)
		}

		attempt++
		time.Sleep(busyRetryDelay)
	}
}

func Exec_DB(db DB_runner, query string, values ...any) {
	if sql_db, ok := db.(*sql.DB); ok && !strings.HasPrefix(strings.ToLower(strings.TrimSpace(query)), "pragma ") {
		err := DB_transaction(sql_db, func(tx *sql.Tx) error {
			Exec_DB_internal(tx, query, values...)
			return nil
		})
		if err != nil {
			panic(err)
		}
		return
	}

	Exec_DB_internal(db, query, values...)
}

func Query_DB(db DB_runner, query string, values ...any) *sql.Rows {
	stmt, err := db.Prepare(DB_change(query))
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	attempt := 0
	start := time.Now()
	for {
		rows, queryErr := stmt.Query(values...)
		if queryErr == nil {
			return rows
		}

		if !busyRetryable(queryErr) ||
			attempt >= busyRetryMax ||
			time.Since(start) >= busyRetryMaxWait {
			panic(queryErr)
		}

		attempt++
		time.Sleep(busyRetryDelay)
	}
}

// QueryRow_DB 이래서 포인터를 배우는구나...
func QueryRow_DB(db DB_runner, query string, var_list []any, values ...any) bool {
	stmt, err := db.Prepare(DB_change(query))
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	attempt := 0
	start := time.Now()
	for {
		row := stmt.QueryRow(values...)

		scanErr := row.Scan(var_list...)
		switch scanErr {
		case nil:
			return true
		case sql.ErrNoRows:
			return false
		}

		if !busyRetryable(scanErr) ||
			attempt >= busyRetryMax ||
			time.Since(start) >= busyRetryMaxWait {
			panic(scanErr)
		}

		attempt++
		time.Sleep(busyRetryDelay)
	}
}

func DB_boot() map[string]string {
	new_db_set := Get_DB_set()
	if new_db_set["db_type"] == "mysql" {
		new_db_set = Get_DB_set_MySQL(new_db_set)
	}

	db_set = new_db_set

	return new_db_set
}

func DB_connect_init() (*sql.DB, error) {
	if db_set["db_type"] == "sqlite" {
		db, err := sql.Open("sqlite", filepath.Join(db_set["db_name"]+".db")+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout=5000")
		if err != nil {
			return nil, err
		}

		if err := db.Ping(); err != nil {
			db.Close()
			return nil, err
		}

		return db, nil
	} else if db_set["db_type"] == "turso" {
		if !Turso_platform_supported() {
			return nil, fmt.Errorf("turso is not supported on %s/%s; use sqlite or mysql db_type", runtime.GOOS, runtime.GOARCH)
		}

		db, err := sql.Open("turso", filepath.Join(db_set["db_name"]+".db"))
		if err != nil {
			return nil, err
		}

		if err := db.Ping(); err != nil {
			db.Close()
			return nil, err
		}

		return db, nil
	} else {
		db, err := sql.Open("mysql", db_set["db_mysql_user"]+":"+db_set["db_mysql_pw"]+"@tcp("+db_set["db_mysql_host"]+":"+db_set["db_mysql_port"]+")/")
		if err != nil {
			return nil, err
		}

		if err := db.Ping(); err != nil {
			db.Close()
			return nil, err
		}

		return db, nil
	}
}

func DB_connect() *sql.DB {
	// log.Default().Println("DB open")

	if db_set["db_type"] == "sqlite" {
		db, err := sql.Open("sqlite", filepath.Join(db_set["db_name"]+".db")+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout=5000")
		if err != nil {
			panic(err)
		}

		return db
	} else if db_set["db_type"] == "turso" {
		if !Turso_platform_supported() {
			panic(fmt.Errorf("turso is not supported on %s/%s; use sqlite or mysql db_type", runtime.GOOS, runtime.GOARCH))
		}

		db, err := sql.Open("turso", filepath.Join(db_set["db_name"]+".db"))
		if err != nil {
			panic(err)
		}

		return db
	} else {
		db, err := sql.Open("mysql", db_set["db_mysql_user"]+":"+db_set["db_mysql_pw"]+"@tcp("+db_set["db_mysql_host"]+":"+db_set["db_mysql_port"]+")/"+db_set["db_name"])
		if err != nil {
			panic(err)
		}

		return db
	}
}

// Turso_platform_supported reports whether the turso driver can run on this
// platform. The tursogo driver loads its native library lazily on first open,
// but windows/arm64 has no precompiled library at all, so turso db_type is
// rejected up front there instead of failing with a driver-internal error.
func Turso_platform_supported() bool {
	return runtime.GOOS != "windows" || runtime.GOARCH != "arm64"
}

func DB_close(db *sql.DB) {
	db.Close()

	// log.Default().Println("DB close")
}

func Start_keeper() error {
	keeperMutex.Lock()
	defer keeperMutex.Unlock()
	return startKeeper_locked()
}

func Stop_keeper() error {
	keeperMutex.Lock()
	defer keeperMutex.Unlock()
	return stopKeeper_locked()
}

func startKeeper_locked() error {
	db_type := db_set["db_type"]
	if db_type != "sqlite" && db_type != "turso" {
		// Keeper-less backends (mysql and unknown types) hold no keeper of
		// their own, but a keeper armed for a previous backend must still be
		// released here - otherwise the old sqlite/turso connection survives
		// a backend switch and keeps the old DB pinned. A release error is
		// propagated so the caller sees a failed keeper arm instead of a
		// silently stale connection.
		return stopKeeper_locked()
	}
	path := filepath.Join(db_set["db_name"]+".db")
	key := db_type + "|" + path
	if keeperDB != nil && keeperConn != nil && keeperKey == key {
		return nil
	}
	stopKeeper_locked()

	ctx := context.Background()
	var db *sql.DB
	var err error
	if db_type == "turso" {
		if !Turso_platform_supported() {
			return fmt.Errorf("keeper: turso is not supported on %s/%s", runtime.GOOS, runtime.GOARCH)
		}
		db, err = sql.Open("turso", path)
	} else {
		dsn := path + "?_pragma=journal_mode(WAL)&_pragma=busy_timeout=5000"
		db, err = sql.Open("sqlite", dsn)
	}
	if err != nil {
		return fmt.Errorf("keeper open: %w", err)
	}
	conn, cerr := db.Conn(ctx)
	if cerr != nil {
		db.Close()
		return fmt.Errorf("keeper conn: %w", cerr)
	}

	// Warm the HELD connection itself. For turso the first query on this
	// *sql.Conn physically loads the local DB instance into the process, which
	// is exactly what the keeper exists to keep alive; a pool-level query
	// would not prove the held connection works. SQLite keeps its existing
	// pool-level ping unchanged.
	var one string
	var serr error
	if db_type == "turso" {
		serr = conn.QueryRowContext(ctx, "select 1").Scan(&one)
	} else {
		serr = db.QueryRow("select 1").Scan(&one)
	}
	if serr != nil {
		conn.Close()
		db.Close()
		return fmt.Errorf("keeper ping: %w", serr)
	}
	keeperDB, keeperConn, keeperKey = db, conn, key
	return nil
}

func stopKeeper_locked() error {
	if keeperConn != nil {
		_ = keeperConn.Close()
	}
	var err error
	if keeperDB != nil {
		err = keeperDB.Close()
	}
	keeperDB, keeperConn = nil, nil
	keeperKey = ""
	return err
}

func Set_db_version(version string) {
	db := DB_connect()
	defer DB_close(db)
	Exec_DB(db, `update other set data = ? where name = "ver"`, version)
}

func Get_DB_type() string {
	return db_set["db_type"]
}

func DB_change(data string) string {
	if Get_DB_type() == "mysql" {
		data = strings.Replace(data, "random()", "rand()", -1)
		data = strings.Replace(data, "collate nocase", "collate utf8mb4_general_ci", -1)
	}

	return data
}

func Exact_binary_eq(column string) string {
	if Get_DB_type() == "mysql" {
		return "CAST(" + column + " AS BINARY) = ?"
	}
	return column + " = ?"
}
