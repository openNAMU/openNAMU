package test

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"opennamu/route/tool"
)

// setup_busy_env points the real DB path at an isolated temp dir and builds the schema
// via the production DB_init. Returns the dir and the exact .db file path used by
// tool.DB_connect (db_set["db_name"]+".db").
func setup_busy_env(t *testing.T) (string, string) {
	dir := filepath.Join(t.TempDir(), "db")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("mkdir db dir: %v", err)
	}
	os.Setenv("NAMU_DB", dir)
	os.Setenv("NAMU_DB_TYPE", "sqlite")
	defer os.Unsetenv("NAMU_DB")
	defer os.Unsetenv("NAMU_DB_TYPE")
	dbFile := dir + ".db"
	tool.DB_init()
	return dir, dbFile
}

// seed minimal rows so the real acl join SELECT is valid and returns rows.
func seed_auth(db *sql.DB, t *testing.T) {
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed %q: %v", q, err)
		}
	}
	exec("insert into alist (name, acl) values ('ip','ip')")
	exec("insert into alist (name, acl) values ('user','user')")
	exec("insert into user_set (name, data) values ('acl','ip')")
}

// is_busy reports whether the message is a busy/lock error (SQLITE_BUSY/SNAPSHOT).
func is_busy(s string) bool {
	le := strings.ToLower(s)
	return strings.Contains(le, "busy") || strings.Contains(le, "database is locked") || strings.Contains(le, "snapshot")
}

// probe_stage mirrors QueryRow_DB (db_connect.go) at two granularities so a busy can be
// attributed to the Prepare stage (lines 227-230, no retry) or the Scan stage (236-254,
// retryed). Returns ("ok","") when the row is read.
func probe_stage(db *sql.DB, query string) (stage string, err string) {
	stmt, perr := db.Prepare(query)
	if perr != nil {
		return "prepare", fmt.Sprintf("%v", perr)
	}
	defer stmt.Close()
	var out string
	if serr := stmt.QueryRow().Scan(&out); serr != nil {
		return "scan", fmt.Sprintf("%v", serr)
	}
	return "ok", ""
}

// writerWork mirrors the production write path: a *sql.DB write goes through
// DB_transaction -> tx.Prepare (panic on a busy, no retry) -> tx.Exec (retried).
func writerWork(db *sql.DB) (stage string, err string) {
	var recVal string
	defer func() {
		if r := recover(); r != nil {
			recVal = fmt.Sprintf("%v", r)
		}
	}()
	if e := tool.DB_transaction(db, func(tx *sql.Tx) error {
		_, execErr := tx.Exec("insert into data (title) values (?)", "probe")
		return execErr
	}); e != nil {
		if recVal != "" {
			return "panic", recVal
		}
		return "txn", fmt.Sprintf("%v", e)
	}
	if recVal != "" {
		return "panic", recVal
	}
	return "ok", ""
}

type churnResult struct {
	total, ok, busyPrep, busyScan, other, panic int64
	msgs []string
	mu   sync.Mutex
}

func (r *churnResult) record(stage string, err string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.total++
	busy := is_busy(err)
	switch stage {
	case "ok", "auth":
		r.ok++
	case "prepare":
		if busy {
			r.busyPrep++
			if len(r.msgs) < 24 {
				r.msgs = append(r.msgs, "prepare:"+err)
			}
		} else {
			r.other++
		}
	case "scan":
		if busy {
			r.busyScan++
			if len(r.msgs) < 24 {
				r.msgs = append(r.msgs, "scan:"+err)
			}
		} else {
			r.other++
		}
	case "panic":
		r.panic++
		if busy {
			r.busyPrep++ // QueryRow_DB / Exec_DB_internal panics on Prepare before reaching Scan/Exec.
		}
		if len(r.msgs) < 24 {
			r.msgs = append(r.msgs, err)
		}
	case "txn":
		if busy {
			r.other++
		} else {
			r.other++
		}
		if len(r.msgs) < 24 {
			r.msgs = append(r.msgs, "txn:"+err)
		}
	default:
		if busy {
			r.busyScan++
		} else {
			r.other++
		}
	}
}

func (r *churnResult) report(t *testing.T) {
	t.Logf("total=%d ok=%d busy@prepare=%d busy@scan=%d other=%d panic=%d",
		r.total, r.ok, r.busyPrep, r.busyScan, r.other, r.panic)
	for _, m := range r.msgs {
		t.Logf("  %s", m)
	}
}

// runChurn drives N workers over I iterations of (open DB_connect / work / DB_close),
// the per-request pattern. authMode uses the real Get_user_auth (middleware path);
// writerShare/10 of workers do real writes; otherwise they read via probe_stage.
func runChurn(t *testing.T, workers, iters, writerShare int, authMode bool) *churnResult {
	r := new(churnResult)
	start := time.Now()
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			isWriter := id%10 < writerShare
			for i := 0; i < iters; i++ {
				db := tool.DB_connect()
				var stage, err string
				if authMode {
					func() {
						defer func() {
							if rec := recover(); rec != nil {
								stage = "panic"
								err = fmt.Sprintf("%v", rec)
							}
						}()
						_ = tool.Get_user_auth(db, "probe_ip")
						stage = "auth"
					}()
				} else if isWriter {
					stage, err = writerWork(db)
				} else {
					stage, err = probe_stage(db, "select 1")
				}
				tool.DB_close(db)
				r.record(stage, err)
			}
		}(w)
	}
	wg.Wait()
	t.Logf("workers=%d iters=%d writerShare=%d authMode=%v elapsed=%v", workers, iters, writerShare, authMode, time.Since(start))
	r.report(t)
	return r
}

// Test_verify_wal_and_busy_timeout confirms the production DSN actually puts the DB in
// WAL and sets busy_timeout=5000 — determining whether a Prepare-stage SQLITE_BUSY should
// have been rescued by busy_timeout (and therefore why it surfaced as a panic).
func Test_verify_wal_and_busy_timeout(t *testing.T) {
	_, _ = setup_busy_env(t)
	db := tool.DB_connect()
	defer tool.DB_close(db)
	var journal, bt string
	if err := db.QueryRow("PRAGMA journal_mode").Scan(&journal); err != nil {
		t.Fatalf("journal_mode: %v", err)
	}
	if err := db.QueryRow("PRAGMA busy_timeout").Scan(&bt); err != nil {
		t.Fatalf("busy_timeout: %v", err)
	}
	t.Logf("journal_mode=%q busy_timeout=%q", journal, bt)
	if journal != "wal" {
		t.Errorf("expected journal_mode=wal got %q", journal)
	}
	if bt != "5000" {
		t.Errorf("expected busy_timeout=5000 got %q", bt)
	}
}

// Test_A_read_churn: read-only per-request churn (open/query/close), no writers.
func Test_A_read_churn(t *testing.T) {
	_, _ = setup_busy_env(t)
	seed_db := tool.DB_connect()
	seed_auth(seed_db, t)
	tool.DB_close(seed_db)
	runChurn(t, 16, 2000, 0, false)
}

// Test_B_read_write_mix: readers + a minority of real writers (per-request).
func Test_B_read_write_mix(t *testing.T) {
	_, _ = setup_busy_env(t)
	seed_db := tool.DB_connect()
	seed_auth(seed_db, t)
	tool.DB_close(seed_db)
	runChurn(t, 16, 2000, 2, false)
}

// Test_C_long_lived_shared: a single shared long-lived pool (all workers share one DB)
// as a control against the per-request churn.
func Test_C_long_lived_shared(t *testing.T) {
	_, _ = setup_busy_env(t)
	seed_db := tool.DB_connect()
	seed_auth(seed_db, t)
	defer tool.DB_close(seed_db)
	db := seed_db
	workers, iters := 16, 2000
	var busyPrep, busyScan, other, total int64
	var mu sync.Mutex
	start := time.Now()
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < iters; i++ {
				stage, err := probe_stage(db, "select 1")
				mu.Lock()
				total++
				if stage == "ok" {
				} else if busy := is_busy(err); busy {
					if stage == "prepare" {
						busyPrep++
					} else {
						busyScan++
					}
				} else {
					other++
				}
				mu.Unlock()
				_ = id
				_ = i
			}
		}(w)
	}
	wg.Wait()
	t.Logf("C shared pool total=%d busy@prepare=%d busy@scan=%d other=%d elapsed=%v", total, busyPrep, busyScan, other, time.Since(start))
}

// Test_D_exclusive_lock_prepare: deterministically holds an exclusive lock (as the
// close-cleanup briefly does) and shows a fresh connection's Prepare returns SQLITE_BUSY,
// which QueryRow_DB then turns into a panic (no retry on Prepare).
func Test_D_exclusive_lock_prepare(t *testing.T) {
	_, dbFile := setup_busy_env(t)
	exDB, err := sql.Open("sqlite", dbFile)
	if err != nil {
		t.Fatalf("open exclusive: %v", err)
	}
	defer exDB.Close()
	if _, e := exDB.Exec("PRAGMA locking_mode=EXCLUSIVE"); e != nil {
		t.Logf("note: locking_mode=EXCLUSIVE failed: %v", e)
	} else {
		t.Log("held exclusive lock")
	}
	probeDB, perr := sql.Open("sqlite", dbFile+"?_pragma=busy_timeout=0")
	if perr != nil {
		t.Fatalf("open probe: %v", perr)
	}
	defer probeDB.Close()
	var stage, pErr string
	func() {
		defer func() {
			if rec := recover(); rec != nil {
				stage = "prepare(panic)"
				pErr = fmt.Sprintf("%v", rec)
			}
		}()
		stmt, perr := probeDB.Prepare("select 1")
		if perr != nil {
			stage = "prepare"
			pErr = fmt.Sprintf("%v", perr)
		} else {
			defer stmt.Close()
			var out string
			if serr := stmt.QueryRow().Scan(&out); serr != nil {
				stage = "scan"
				pErr = fmt.Sprintf("%v", serr)
			} else {
				stage = "ok"
			}
		}
	}()
	// Diagnostic: in WAL mode a separately-held EXCLUSIVE lock did NOT block a
	// concurrent read Prepare (stage="ok"), so this is not a faithful model of the
	// close-cleanup lock. The churn tests (A/E) and the timing test reproduce the
	// real condition instead.
	t.Logf("D exclusive-lock probe stage=%q err=%q", stage, pErr)
}

// Test_E_real_auth_path: the real Get_user_auth (QueryRow_DB) through per-request churn,
// the exact middleware path that produced the production panic.
func Test_E_real_auth_path(t *testing.T) {
	_, _ = setup_busy_env(t)
	seed_db := tool.DB_connect()
	seed_auth(seed_db, t)
	tool.DB_close(seed_db)
	runChurn(t, 16, 2000, 0, true)
}

// Test_timing_prepare_busy measures how long a busy Prepare actually lasted during churn.
// This disambiguates two hypotheses: (a) busy_timeout=5000 applied and a lock was held
// for ~5s, vs (b) the busy happened before/without the busy_timeout pragma taking effect
// (sub-millisecond).
func Test_timing_prepare_busy(t *testing.T) {
	_, _ = setup_busy_env(t)
	seed_db := tool.DB_connect()
	seed_auth(seed_db, t)
	tool.DB_close(seed_db)
	const workers, iters = 16, 2500
	var busyTimes []time.Duration
	var mu sync.Mutex
	var wg sync.WaitGroup
	start := time.Now()
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < iters; i++ {
				db := tool.DB_connect()
				pstart := time.Now()
				_, perr := db.Prepare("select 1")
				elapsed := time.Since(pstart)
				tool.DB_close(db)
				if perr != nil && is_busy(fmt.Sprintf("%v", perr)) {
					mu.Lock()
					busyTimes = append(busyTimes, elapsed)
					mu.Unlock()
				}
			}
		}(w)
	}
	wg.Wait()
	tot := time.Since(start)
	sort.Slice(busyTimes, func(a, b int) bool { return busyTimes[a] < busyTimes[b] })
	if len(busyTimes) > 0 {
		var sum time.Duration
		for _, d := range busyTimes {
			sum += d
		}
		t.Logf("busy count=%d total=%v mean=%v min=%v max=%v",
			len(busyTimes), tot, sum/time.Duration(len(busyTimes)), busyTimes[0], busyTimes[len(busyTimes)-1])
	} else {
		t.Logf("no busy errors observed in %d workers x %d iters (%v)", workers, iters, tot)
	}
}
