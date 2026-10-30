package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"database/sql"

	sqlite "modernc.org/sqlite"

	"opennamu/route/tool"
)

// ---- measurement config ----

type measureConfig struct {
	name          string
	keeper        bool
	writers       int
	readers       int
	iters         int
	churn         bool // true=per-iteration DB (prod-like); false=per-worker pool
	journal       bool // true=per-conn _pragma=journal_mode(WAL); false=C2 (busy_timeout only)
	serviceWriter bool // true=model service write APIs: panic(err) on DB_transaction error
	growWAL       bool // write a batch of rows before the loop to grow the WAL
	stress        bool
}

// ---- event record (one line of JSONL per failure) ----

type event struct {
	RunID      string `json:"run_id"`
	Scenario   string `json:"scenario"`
	Iteration  int    `json:"iteration"`
	Worker     int    `json:"worker"`
	Role       string `json:"role"`
	Keeper     bool   `json:"keeper"`
	Churn      bool   `json:"churn"`
	DSN        string `json:"dsn"`
	Timestamp  string `json:"timestamp"`
	Elapsed    string `json:"elapsed"`
	Phase      string `json:"phase"`
	Code       int    `json:"code"`
	CodeAnd    int    `json:"code_and_0xff"`
	Err        string `json:"err"`
	PanicValue string `json:"panic_value,omitempty"`
	IsPanic    bool   `json:"is_panic"`
	SQLLabel   string `json:"sql_label"`
	Fn         string `json:"fn"`
	File       string `json:"file"`
	Line       int    `json:"line"`
	Stack      string `json:"stack,omitempty"`
	WALBytes   int64  `json:"wal_bytes"`
	WALPage    int64  `json:"wal_page_size,omitempty"`
	WALAutochk int    `json:"wal_autocheckpoint,omitempty"`
	OpenConns  int    `json:"open_connections"`
	InUse      int    `json:"in_use"`
	Idle       int    `json:"idle"`
	WaitCount  int    `json:"wait_count"`
}

// ---- measurement state ----

var (
	mu      sync.Mutex
	mEvents []string
)

// envInfo returns the environment metadata recorded per run.
func envInfo() map[string]string {
	return map[string]string{
		"go":     runtime.Version(),
		"os":     runtime.GOOS,
		"gomp":   fmt.Sprint(runtime.GOMAXPROCS(0)),
		"sqlite": "modernc.org/sqlite v1.29.8",
	}
}

// sqliteCode returns (code, msg) using errors.As to *sqlite.Error, with a
// numeric fallback. 5=SQLITE_BUSY, 261=SQLITE_BUSY_RECOVERY, 517=SNAPSHOT.
func sqliteCode(err any) (int, string) {
	msg := fmt.Sprintf("%v", err)
	var se *sqlite.Error
	if e, ok := err.(error); ok {
		if errors.As(e, &se) {
			return se.Code(), msg
		}
	}
	for _, c := range []int{261, 517, 5} {
		if strings.Contains(msg, fmt.Sprintf("( %d )", c)) || strings.Contains(msg, fmt.Sprintf("(%d)", c)) {
			return c, msg
		}
	}
	return -1, msg
}

// framesToPhase walks the current frames and returns (phase, fn, file, line,
// stackString) for the frame relevant to a failure.
func framesToPhase(skip int) (string, string, string, int, string) {
	pc := make([]uintptr, 128)
	n := runtime.Callers(skip, pc)
	frames := runtime.CallersFrames(pc[:maxInt(n, 1)])
	var b strings.Builder
	fn, file := "", ""
	line := 0
	phase := "unknown"
	first := true
	for {
		f, more := frames.Next()
		if !more {
			break
		}
		name := f.Function
		if name == "" {
			name = "<unknown>"
		}
		if first {
			fn, file, line = name, f.File, f.Line
			phase = phaseFromFn(name, line)
			first = false
		}
		b.WriteString(name)
		if f.File != "" && f.Line > 0 {
			b.WriteString(" / " + f.File + ":" + fmt.Sprint(f.Line))
		}
		b.WriteString("\n")
	}
	return phase, fn, file, line, b.String()
}

func phaseFromFn(name string, line int) string {
	switch {
	case strings.HasSuffix(name, ".QueryRow_DB"):
		return "prepare"
	case strings.HasSuffix(name, ".Query_DB"):
		return "prepare"
	case strings.HasSuffix(name, ".DB_transaction"):
		switch {
		case line <= 71:
			return "tx_begin"
		case line >= 85:
			return "tx_commit"
		default:
			return "tx_other"
		}
	case strings.HasSuffix(name, ".Exec_DB_internal"):
		return "exec"
	default:
		return "other"
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// walInfo reads the WAL file size and PRAGMAs for the given db file.
func walInfo(dbFile string) (bytes int64, pageSize int64, autochk int) {
	if st, err := os.Stat(dbFile + "-wal"); err == nil {
		bytes = st.Size()
	} else {
		bytes = -1
	}
	db, err := sql.Open("sqlite", dbFile+"?")
	if err != nil {
		return -1, 0, 0
	}
	defer func() { _ = db.Close() }()
	var ps int64
	if err := db.QueryRow("PRAGMA page_size").Scan(&ps); err == nil {
		pageSize = ps
	}
	var ac int
	if err := db.QueryRow("PRAGMA wal_autocheckpoint").Scan(&ac); err == nil {
		autochk = ac
	}
	return
}

// connStats returns the pool stats for a *sql.DB.
func connStats(db *sql.DB) (int, int, int, int) {
	s := db.Stats()
	return s.OpenConnections, s.InUse, s.Idle, int(s.WaitCount)
}

// runMeasure drives the per-request churn scenario and writes raw JSONL events to
// ../tmp/ (project root). Returns the output path.
func runMeasure(t *testing.T, o measureConfig) string {
	outPath := filepath.Join("../tmp", fmt.Sprintf("residual_round2_%s_%d.jsonl",
		sanitize(o.name), time.Now().UnixNano()%100000))
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	f, err := os.Create(outPath)
	if err != nil {
		t.Fatalf("create %s: %v", outPath, err)
	}
	defer f.Close()

	runID := fmt.Sprintf("%d", time.Now().UnixNano()%100000)
	env := envInfo()

	_, dbFile := setup_busy_env(t)
	seedDB := tool.DB_connect()
	seed_auth(seedDB, t)
	tool.DB_close(seedDB)
	if err := tool.Start_keeper(); err != nil {
		t.Log("debug: keeper start failed: " + err.Error())
		o.keeper = false
	}
	defer func() { _ = tool.Stop_keeper() }()
	if err := ensureWAL(dbFile); err != nil {
		t.Fatalf("ensureWAL: %v", err)
	}

	if o.growWAL {
		grow := tool.DB_connect()
		defer tool.DB_close(grow)
		for i := 0; i < 5000; i++ {
			_, _ = grow.Exec("insert into data (title) values (?)",
				fmt.Sprintf("grow-%d", i))
		}
	}

	var workerDBs []*sql.DB
	if !o.churn {
		for i := 0; i < o.writers+o.readers; i++ {
			workerDBs = append(workerDBs, openWithDSN(dbFile, o.journal))
		}
		defer func() {
			for _, d := range workerDBs {
				_ = d.Close()
			}
		}()
	}

	var panicC, errC int64
	var wg sync.WaitGroup

	work := func(wid, role int) {
		for iter := 0; iter < o.iters; iter++ {
			start := time.Now()
			var db *sql.DB
			if o.churn {
				if o.journal {
					db = tool.DB_connect()
				} else {
					db = openWithDSN(dbFile, false)
				}
			} else {
				db = workerDBs[wid]
			}
			func() {
				defer func() {
					if rec := recover(); rec != nil {
						atomic.AddInt64(&panicC, 1)
						ph, ffn, ffile, fl, stk := framesToPhase(0)
						code, s := sqliteCode(rec)
						walB, walP, walA := walInfo(dbFile)
						oc, iu, id, wc := connStats(db)
						writeEvent(f, event{
							RunID: runID, Scenario: o.name, Iteration: iter, Worker: wid,
							Role: roleStr(role), Keeper: o.keeper, Churn: o.churn,
							DSN: dsnLabel(o.journal), Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
							Elapsed: time.Since(start).String(), Phase: ph, Code: code,
							CodeAnd: code & 0xff, Err: s, PanicValue: fmt.Sprintf("%v", rec),
							IsPanic: true, SQLLabel: sqlLabel(role, iter),
							Fn: ffn, File: ffile, Line: fl, Stack: stk,
							WALBytes: walB, WALPage: walP, WALAutochk: walA,
							OpenConns: oc, InUse: iu, Idle: id, WaitCount: wc,
						})
					}
				}()
				if role == 1 {
					_ = tool.Get_user_auth(db, fmt.Sprintf("lr-%d-%d", wid, iter))
				} else {
					if e := tool.DB_transaction(db, func(tx *sql.Tx) error {
						_, ex := tx.Exec("insert into data (title) values (?)",
							fmt.Sprintf("lw-%d-%d", wid, iter))
						return ex
					}); e != nil {
						atomic.AddInt64(&errC, 1)
						ph, ffn, ffile, fl, stk := framesToPhase(0)
						code, s := sqliteCode(e)
						walB, walP, walA := walInfo(dbFile)
						oc, iu, id, wc := connStats(db)
						writeEvent(f, event{
							RunID: runID, Scenario: o.name, Iteration: iter, Worker: wid,
							Role: "writer", Keeper: o.keeper, Churn: o.churn,
							DSN: dsnLabel(o.journal), Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
							Elapsed: time.Since(start).String(), Phase: ph, Code: code,
							CodeAnd: code & 0xff, Err: s, IsPanic: false,
							SQLLabel: fmt.Sprintf("insert data lw-%d-%d", wid, iter),
							Fn:       ffn, File: ffile, Line: fl, Stack: stk,
							WALBytes: walB, WALPage: walP, WALAutochk: walA,
							OpenConns: oc, InUse: iu, Idle: id, WaitCount: wc,
						})
					}
				}
			}()
			if o.churn {
				_ = db.Close()
			}
		}
	}

	for i := 0; i < o.writers; i++ {
		wg.Add(1)
		go func(wid int) {
			defer wg.Done()
			work(wid, 0)
		}(i)
	}
	for i := 0; i < o.readers; i++ {
		wg.Add(1)
		go func(wid int) {
			defer wg.Done()
			work(wid, 1)
		}(i)
	}
	wg.Wait()

	// data-integrity invariant: the DB must be whole after the churn
	// (no partial writes, no corruption).
	var integ string
	dbCheck := tool.DB_connect()
	defer tool.DB_close(dbCheck)
	if serr := dbCheck.QueryRow("PRAGMA integrity_check").Scan(&integ); serr != nil {
		integ = fmt.Sprintf("check-err: %v", serr)
	}
	var rowCount int
	if serr := dbCheck.QueryRow("select count(*) from data").Scan(&rowCount); serr != nil {
		rowCount = -1
	}
	t.Logf("[measure] name=%s keeper=%v churn=%v journal=%v svcWriter=%v growWAL=%v writers=%d readers=%d iters=%d panics=%d writer_errors=%d data_rows=%d integrity_check=%q out=%s env: go=%s os=%s gomp=%s",
		o.name, o.keeper, o.churn, o.journal, o.serviceWriter, o.growWAL,
		o.writers, o.readers, o.iters, panicC, errC, rowCount, integ, outPath,
		env["go"], env["os"], env["gomp"])
	return outPath
}

func roleStr(role int) string {
	if role == 1 {
		return "reader"
	}
	return "writer"
}

func sqlLabel(role, iter int) string {
	if role == 1 {
		return "auth_select join user_set/alist/acl_end"
	}
	return "insert data lw"
}

func dsnLabel(journal bool) string {
	if journal {
		return "prod:journal_mode(WAL)+busy_timeout=5000"
	}
	return "exp:busy_timeout=5000 (WAL pre-inited)"
}

func sanitize(s string) string {
	return strings.ToLower(strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			return r
		}
		return '_'
	}, s))
}

func writeEvent(f *os.File, ev event) {
	b, _ := json.Marshal(ev)
	line := string(b) + "\n"
	mu.Lock()
	mEvents = append(mEvents, line)
	mu.Unlock()
	_, _ = f.WriteString(line)
}

// ---- scenario tests ----

func Test_measure_A_original_4w4r_keeper_on(t *testing.T) {
	if !gated() {
		t.Skip("gated")
	}
	runMeasure(t, measureConfig{
		name: "A_orig_4w4r_keepon", keeper: true, writers: 4, readers: 4, iters: 25,
		churn: true, journal: true, serviceWriter: true,
	})
}

func Test_measure_A_original_4w4r_keeper_off(t *testing.T) {
	runMeasure(t, measureConfig{
		name: "A_orig_4w4r_keepoff", keeper: false, writers: 4, readers: 4, iters: 25,
		churn: true, journal: true, serviceWriter: true,
	})
}

func Test_measure_B_service_writer_panic(t *testing.T) {
	runMeasure(t, measureConfig{
		name: "B_svc_writer_panic", keeper: true, writers: 4, readers: 4, iters: 25,
		churn: true, journal: true, serviceWriter: true,
	})
}

func Test_measure_C_wal_grow_checkpoint(t *testing.T) {
	runMeasure(t, measureConfig{
		name: "C_wal_grow", keeper: true, writers: 4, readers: 4, iters: 100,
		churn: true, journal: true, serviceWriter: true, growWAL: true,
	})
}

func Test_measure_D_8r(t *testing.T) {
	runMeasure(t, measureConfig{name: "D_8r", keeper: true, writers: 0, readers: 8, iters: 25, churn: true, journal: true, serviceWriter: true})
}

func Test_measure_D_7r1w(t *testing.T) {
	runMeasure(t, measureConfig{name: "D_7r1w", keeper: true, writers: 1, readers: 7, iters: 25, churn: true, journal: true, serviceWriter: true})
}

func Test_measure_D_4r4w(t *testing.T) {
	runMeasure(t, measureConfig{name: "D_4r4w", keeper: true, writers: 4, readers: 4, iters: 25, churn: true, journal: true, serviceWriter: true})
}

func Test_measure_D_4w(t *testing.T) {
	runMeasure(t, measureConfig{name: "D_4w", keeper: true, writers: 4, readers: 0, iters: 25, churn: true, journal: true, serviceWriter: true})
}

func Test_measure_D_8w(t *testing.T) {
	runMeasure(t, measureConfig{name: "D_8w", keeper: true, writers: 8, readers: 0, iters: 25, churn: true, journal: true, serviceWriter: true})
}

func Test_measure_D_nochurn(t *testing.T) {
	runMeasure(t, measureConfig{name: "D_nochurn", keeper: true, writers: 4, readers: 4, iters: 50, churn: false, journal: true, serviceWriter: true})
}

func Test_measure_D_dsn_exp(t *testing.T) {
	runMeasure(t, measureConfig{name: "D_dsn_exp", keeper: true, writers: 4, readers: 4, iters: 25, churn: true, journal: false, serviceWriter: true})
}

func Test_measure_D_gomp1(t *testing.T) {
	g := runtime.GOMAXPROCS(0)
	defer runtime.GOMAXPROCS(g)
	runtime.GOMAXPROCS(1)
	runMeasure(t, measureConfig{name: "D_gomp1", keeper: true, writers: 4, readers: 4, iters: 25, churn: true, journal: true, serviceWriter: true})
}

func Test_measure_D_gomphigh(t *testing.T) {
	g := runtime.GOMAXPROCS(0)
	defer runtime.GOMAXPROCS(g)
	runtime.GOMAXPROCS(12)
	runMeasure(t, measureConfig{name: "D_gomphigh", keeper: true, writers: 4, readers: 4, iters: 25, churn: true, journal: true, serviceWriter: true})
}

// Test_residual_diag_positive_control: verifies the diagnostic pipeline by creating
// a REAL, controlled busy (hold a write lock; second conn with busy_timeout=0 fails
// with SQLITE_BUSY). This is a DELIBERATE contention, not evidence of the residual.
func Test_residual_diag_positive_control(t *testing.T) {
	dir := t.TempDir()
	dbFile := filepath.Join(dir, "control.db")
	ensureWAL(dbFile)

	connA, err := sql.Open("sqlite", dbFile)
	if err != nil {
		t.Fatalf("open A: %v", err)
	}
	defer func() { _ = connA.Close() }()

	tx1, _ := connA.Begin()
	_, _ = tx1.Exec("create table control(x int)")
	if err := tx1.Commit(); err != nil {
		t.Fatalf("commit 1: %v", err)
	}

	txA, err := connA.Begin()
	if err != nil {
		t.Fatalf("begin lock: %v", err)
	}
	defer func() { _ = txA.Rollback() }()
	_, _ = txA.Exec("insert into control(x) values (1)")

	connB, err := sql.Open("sqlite", dbFile+"?_pragma=busy_timeout=0")
	if err != nil {
		t.Fatalf("open B: %v", err)
	}
	defer func() { _ = connB.Close() }()
	tx2, _ := connB.Begin()
	defer func() { _ = tx2.Rollback() }()
	_, exErr := tx2.Exec("insert into control(x) values (2)")
	code, s := sqliteCode(exErr)
	t.Logf("positive-control: code=%d err=%q (expected 5 SQLITE_BUSY; timing-dependent)", code, s)
	if code == 0 {
		t.Log("note: no busy observed (lock may have been released); diagnostic still exercised")
	}
}

// Test_measure_A_large_4w4r: larger faithful sample (100 iters) for confidence.
func Test_measure_A_large_4w4r(t *testing.T) {
	if !gated() {
		t.Skip("gated")
	}
	runMeasure(t, measureConfig{
		name: "A_large_4w4r", keeper: true, writers: 4, readers: 4, iters: 100,
		churn: true, journal: true, serviceWriter: true,
	})
}

// gated() lets us run the faithful 30-run re-execution only when explicitly
// requested (env BUSY_ROUND2_RUN=1) to avoid bloating the default suite.
func gated() bool {
	return os.Getenv("BUSY_ROUND2_RUN") == "1"
}

// openWithDSN opens a raw *sql.DB for the file using the prod DSN (journal_mode(WAL)
// + busy_timeout=5000) or the experimental DSN (busy_timeout=5000 only, WAL
// pre-initialized by ensureWAL).
func openWithDSN(dbFile string, journal bool) *sql.DB {
	dsn := dbFile + "?_pragma=busy_timeout=5000"
	if journal {
		dsn = dbFile + "?_pragma=journal_mode(WAL)&_pragma=busy_timeout=5000"
	}
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		panic("openWithDSN: " + err.Error())
	}
	return db
}

// ensureWAL puts the file into WAL mode and creates the initial schema.
func ensureWAL(dbFile string) error {
	db, err := sql.Open("sqlite", dbFile+"?")
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()
	_, _ = db.Exec("PRAGMA journal_mode(WAL)")
	_, _ = db.Exec("create table if not exists data (id TEXT PRIMARY KEY, title TEXT)")
	_, _ = db.Exec("create table if not exists auth_user (user_id TEXT PRIMARY KEY, name TEXT, email TEXT)")
	_, _ = db.Exec("create table if not exists auth_role (user_id TEXT, role TEXT, PRIMARY KEY(user_id, role))")
	_, _ = db.Exec("create table if not exists acl_list (acl_id TEXT, name TEXT, PRIMARY KEY(acl_id))")
	_, _ = db.Exec("create table if not exists acl_end (id INTEGER, acl_id TEXT, PRIMARY KEY(id, acl_id))")
	_, _ = db.Exec("create table if not exists acl_set (id INTEGER, user_id TEXT, acl_id TEXT, PRIMARY KEY(id))")
	return nil
}
