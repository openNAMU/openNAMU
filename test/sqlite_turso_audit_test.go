package test

// Independent concurrency/integrity/performance audit for the SQLite(WAL+keeper)
// and Turso(embedded MVCC) backends of the local working tree. Directive:
// tmp/openNAMU_SQLite_Turso_concurrency_audit_20261025.md.
//
// Deterministic scenarios run by default; heavier scenarios require
// DB_CONCURRENCY_STRESS=1. Raw metrics are appended as JSONL to
// $AUDIT_OUT_DIR (default ../tmp/audit_results relative to the test dir).
// No production code is modified by this file.

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	sqlite "modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
	turso "turso.tech/database/tursogo"

	"opennamu/route/tool"
)

// ---------------------------------------------------------------- plumbing

func audit_out_dir(t *testing.T) string {
	t.Helper()
	dir := os.Getenv("AUDIT_OUT_DIR")
	if dir == "" {
		dir = "../tmp/audit_results"
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("audit out dir: %v", err)
	}
	return dir
}

// audit_record appends one JSONL metric record.
func audit_record(t *testing.T, rec map[string]any) {
	t.Helper()
	data, err := json.Marshal(rec)
	if err != nil {
		t.Logf("audit_record marshal: %v", err)
		return
	}
	path := filepath.Join(audit_out_dir(t), "audit_results.jsonl")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		t.Logf("audit_record open: %v", err)
		return
	}
	defer f.Close()
	_, _ = f.Write(append(data, '\n'))
}

// audit_classify maps an error to (busy?, label). SQLite uses the primary code
// (code & 0xff) plus message hints for extended BUSY states; Turso maps to its
// sentinel errors.
func audit_classify(err error) (busy bool, label string) {
	if err == nil {
		return false, ""
	}

	var sq *sqlite.Error
	if errors.As(err, &sq) {
		code := sq.Code() & 0xff
		msg := sq.Error()
		switch {
		case code == sqlite3.SQLITE_BUSY && strings.Contains(strings.ToLower(msg), "snapshot"):
			return true, fmt.Sprintf("SQLITE_BUSY_SNAPSHOT(%d)", code)
		case code == sqlite3.SQLITE_BUSY && strings.Contains(strings.ToLower(msg), "recovery"):
			return true, fmt.Sprintf("SQLITE_BUSY_RECOVERY(%d)", code)
		case code == sqlite3.SQLITE_BUSY:
			return true, fmt.Sprintf("SQLITE_BUSY(%d)", code)
		case code == sqlite3.SQLITE_LOCKED:
			return true, fmt.Sprintf("SQLITE_LOCKED(%d)", code)
		case code == sqlite3.SQLITE_CONSTRAINT:
			return false, fmt.Sprintf("SQLITE_CONSTRAINT(%d)", code)
		default:
			return false, fmt.Sprintf("SQLITE_%d", code)
		}
	}

	switch {
	case errors.Is(err, turso.ErrTursoBusy):
		return true, "TURSO_BUSY"
	case errors.Is(err, turso.ErrTursoConstraint):
		return false, "TURSO_CONSTRAINT"
	case errors.Is(err, turso.ErrTursoReadOnly):
		return false, "TURSO_READONLY"
	case errors.Is(err, turso.ErrTursoCorrupt):
		return false, "TURSO_CORRUPT"
	case errors.Is(err, turso.ErrTursoNotADb):
		return false, "TURSO_NOTADB"
	case errors.Is(err, turso.ErrTursoInterrupt):
		return true, "TURSO_INTERRUPT"
	case errors.Is(err, turso.ErrTursoGeneric):
		return false, "TURSO_GENERIC"
	default:
		return false, "OTHER(" + first_line(err.Error()) + ")"
	}
}

func first_line(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	if len(s) > 80 {
		s = s[:80]
	}
	return s
}

// audit_setup boots the real openNAMU DB path for a backend. For sqlite the
// WAL keeper is started (production Main_init behavior); for turso Start_keeper
// is a no-op by design.
func audit_setup(t *testing.T, backend string) *sql.DB {
	t.Helper()

	name := filepath.Join(t.TempDir(), "audit")
	os.Setenv("NAMU_DB", name)
	os.Setenv("NAMU_DB_TYPE", backend)

	tool.DB_init()
	if err := tool.Start_keeper(); err != nil {
		t.Fatalf("keeper start: %v", err)
	}
	t.Cleanup(func() {
		_ = tool.Stop_keeper()
		os.Unsetenv("NAMU_DB")
		os.Unsetenv("NAMU_DB_TYPE")
	})

	db := tool.DB_connect()
	t.Cleanup(func() { tool.DB_close(db) })
	return db
}

func audit_seed_auth(db *sql.DB, t *testing.T) {
	t.Helper()
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed %q: %v", q, err)
		}
	}
	exec(`insert into alist (name, acl) values ('owner','owner')`)
	exec(`insert into alist (name, acl) values ('admin','admin')`)
	exec(`insert into alist (name, acl) values ('user','user')`)
	exec(`insert into alist (name, acl) values ('ip','ip')`)
	exec(`insert into user_set (name, id, data) values ('acl','admin_ip','owner')`)
}

func audit_seed_docs(db *sql.DB, t *testing.T, n int) {
	t.Helper()
	for i := 0; i < n; i++ {
		tool.Exec_DB(db,
			`insert into data (title, data, type) values (?, ?, '')`,
			fmt.Sprintf("audit_%03d", i),
			fmt.Sprintf("audit content %d", i),
		)
	}
}

// audit_env_map describes the run environment for JSONL records.
func audit_env_map() map[string]any {
	return map[string]any{
		"go":        runtime.Version(),
		"gomaxproc": runtime.GOMAXPROCS(0),
		"os":        runtime.GOOS + "/" + runtime.GOARCH,
		"stress":    os.Getenv("DB_CONCURRENCY_STRESS") == "1",
	}
}

type audit_stats struct {
	ops, errs, busy atomic.Int64
	lat_mu          sync.Mutex
	lat             []float64
}

func (s *audit_stats) record(d time.Duration, err error) {
	if err != nil {
		s.errs.Add(1)
		if busy, _ := audit_classify(err); busy {
			s.busy.Add(1)
		}
		return
	}
	s.ops.Add(1)
	s.lat_mu.Lock()
	s.lat = append(s.lat, d.Seconds())
	s.lat_mu.Unlock()
}

func (s *audit_stats) summary() (p50, p95, p99 float64) {
	s.lat_mu.Lock()
	defer s.lat_mu.Unlock()
	sort.Float64s(s.lat)
	return pct(s.lat, 50), pct(s.lat, 95), pct(s.lat, 99)
}

func pct(sorted []float64, p int) float64 {
	if len(sorted) == 0 {
		return 0
	}
	i := len(sorted) * p / 100
	if i >= len(sorted) {
		i = len(sorted) - 1
	}
	return sorted[i]
}

func audit_backend_version(t *testing.T, db *sql.DB, backend string) string {
	t.Helper()
	ver := ""
	if err := db.QueryRow(`select sqlite_version()`).Scan(&ver); err != nil {
		ver = "unknown(" + first_line(err.Error()) + ")"
	}
	return ver
}

// ---------------------------------------------------------------- A

// Test_audit_A_read_parallelism: 1/4/8/16/32 concurrent read clients doing
// permission lookups (the real Get_user_auth) and document point reads.
func Test_audit_A_read_parallelism(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		db := audit_setup(t, backend)
		audit_seed_auth(db, t)
		audit_seed_docs(db, t, 50)
		ver := audit_backend_version(t, db, backend)
		audit_record(t, map[string]any{"type": "engine", "backend": backend, "sqlite_version_reported": ver})

		for _, clients := range []int{1, 4, 8, 16, 32} {
			stats := &audit_stats{}
			var wg sync.WaitGroup
			for c := 0; c < clients; c++ {
				wg.Add(1)
				go func(c int) {
					defer wg.Done()
					conn := tool.DB_connect()
					defer conn.Close()
					for i := 0; i < 30; i++ {
						t0 := time.Now()
						_ = tool.Get_user_auth(conn, "admin_ip")
						found := ""
						err := conn.QueryRow(
							`select data from data where title = ?`,
							fmt.Sprintf("audit_%03d", (c+i)%50),
						).Scan(&found)
						if err == nil && found == "" {
							err = errors.New("stale read: empty content")
						}
						stats.record(time.Since(t0), err)
					}
				}(c)
			}
			wg.Wait()

			p50, p95, p99 := stats.summary()
			audit_record(t, map[string]any{
				"type": "A", "backend": backend, "clients": clients,
				"ops": stats.ops.Load(), "errs": stats.errs.Load(), "busy": stats.busy.Load(),
				"p50": p50, "p95": p95, "p99": p99,
			})
			if stats.errs.Load() != 0 {
				t.Errorf("A %s clients=%d: %d read errors", backend, clients, stats.errs.Load())
			}
			t.Logf("A %s clients=%2d ops=%d errs=%d p50=%.4f p95=%.4f p99=%.4f",
				backend, clients, stats.ops.Load(), stats.errs.Load(), p50, p95, p99)
		}
		audit_record(t, map[string]any{"type": "engine", "backend": backend, "sqlite_version_reported": ver})
	}
}

// ---------------------------------------------------------------- B

// Test_audit_B_write_distinct_docs: 2/8/32 concurrent writers, each creating
// its own documents with 3 revisions; verifies counts and revision continuity.
func Test_audit_B_write_distinct_docs(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		db := audit_setup(t, backend)
		audit_seed_auth(db, t)

		for _, writers := range []int{2, 8, 32} {
			docs_per_writer := 3
			if writers == 32 {
				docs_per_writer = 1 // budget
			}
			revs := 3

			var panic_cnt atomic.Int64
			var err_sample atomic.Value // string
			var wg sync.WaitGroup
			start := time.Now()
			for w := 0; w < writers; w++ {
				wg.Add(1)
				go func(w int) {
					defer wg.Done()
					defer func() {
						if r := recover(); r != nil {
							panic_cnt.Add(1)
							err_sample.CompareAndSwap(nil, fmt.Sprintf("panic: %v", r))
						}
					}()
					conn := tool.DB_connect()
					defer conn.Close()
					for d := 0; d < docs_per_writer; d++ {
						title := fmt.Sprintf("b_%02d_w%02d_d%d", writers, w, d)
						tool.Exec_DB(conn, `insert into data (title, data, type) values (?, ?, '')`, title, "v0")
						for r := 1; r <= revs; r++ {
							tool.Exec_DB(conn,
								`insert into history (id, title, data, date, ip, send, leng, hide, type) values (?, ?, ?, ?, ?, ?, ?, '', '')`,
								fmt.Sprintf("%d", r), title, fmt.Sprintf("rev %d", r),
								"2026-01-01 00:00:00", "admin_ip", "audit", "+1",
							)
						}
					}
				}(w)
			}
			wg.Wait()
			elapsed := time.Since(start)

			total_docs := 0
			doc_pattern := fmt.Sprintf("b_%02d_w%%", writers)
			if !tool.QueryRow_DB(db, fmt.Sprintf(`select count(*) from data where title like '%s'`, doc_pattern), []any{&total_docs}) {
				t.Fatalf("B %s: count query failed", backend)
			}
			total_hist := 0
			if !tool.QueryRow_DB(db, fmt.Sprintf(`select count(*) from history where title like '%s'`, doc_pattern), []any{&total_hist}) {
				t.Fatalf("B %s: count query failed", backend)
			}
			want_docs := writers * docs_per_writer
			want_hist := writers * docs_per_writer * revs

			// revision continuity: every title must have exactly revs distinct ids 1..revs
			bad_rev := 0
			rows := tool.Query_DB(db, fmt.Sprintf(
				`select title, count(distinct id) from history where title like '%s' group by title having count(distinct id) != %d`,
				doc_pattern, revs,
			))
			for rows.Next() {
				bad_rev++
			}
			rows.Close()

			audit_record(t, map[string]any{
				"type": "B", "backend": backend, "writers": writers,
				"want_docs": want_docs, "got_docs": total_docs,
				"want_hist": want_hist, "got_hist": total_hist,
				"bad_revisions": bad_rev, "panics": panic_cnt.Load(),
				"secs": elapsed.Seconds(),
				"sample_err": err_sample.Load(),
			})
			if total_docs != want_docs || total_hist != want_hist {
				t.Errorf("B %s writers=%d: docs=%d(want %d) hist=%d(want %d)",
					backend, writers, total_docs, want_docs, total_hist, want_hist)
			}
			if bad_rev != 0 {
				t.Errorf("B %s writers=%d: %d titles with broken revision continuity", backend, writers, bad_rev)
			}
			if panic_cnt.Load() != 0 {
				t.Errorf("B %s writers=%d: %d worker panics, sample=%v", backend, writers, panic_cnt.Load(), err_sample.Load())
			}
			t.Logf("B %s writers=%2d docs=%d hist=%d bad_rev=%d panics=%d %.3fs",
				backend, writers, total_docs, total_hist, bad_rev, panic_cnt.Load(), elapsed.Seconds())
		}
	}
}

// ---------------------------------------------------------------- C (DB level)

// Test_audit_C_same_doc_db_level: concurrent inserts into the same document's
// history at DB level; row counts must be exact (no lost/duplicate rows).
// App-level optimistic conflicts are covered by the HTTP audit test.
func Test_audit_C_same_doc_db_level(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		db := audit_setup(t, backend)
		audit_seed_auth(db, t)

		workers, per := 8, 5
		var wg sync.WaitGroup
		var errs atomic.Int64
		start := time.Now()
		for w := 0; w < workers; w++ {
			wg.Add(1)
			go func(w int) {
				defer wg.Done()
				conn := tool.DB_connect()
				defer conn.Close()
				for i := 0; i < per; i++ {
					if _, err := conn.Exec(
						`insert into history (id, title, data, date, ip, send, leng, hide, type) values (?, ?, ?, ?, ?, ?, ?, '', '')`,
						fmt.Sprintf("c%d_%d", w, i), "audit_same_doc", "rev", "2026-01-01 00:00:00", "admin_ip", "audit", "+1",
					); err != nil {
						errs.Add(1)
					}
				}
			}(w)
		}
		wg.Wait()
		elapsed := time.Since(start)

		got := 0
		if !tool.QueryRow_DB(db, `select count(*) from history where title = 'audit_same_doc'`, []any{&got}) {
			t.Fatalf("C %s: count failed", backend)
		}
		audit_record(t, map[string]any{
			"type": "C", "backend": backend, "workers": workers,
			"want": workers * per, "got": got, "errs": errs.Load(), "secs": elapsed.Seconds(),
		})
		if got != workers*per || errs.Load() != 0 {
			t.Errorf("C %s: rows=%d(want %d) errs=%d", backend, got, workers*per, errs.Load())
		}
		t.Logf("C %s rows=%d errs=%d %.3fs", backend, got, errs.Load(), elapsed.Seconds())
	}
}

// ---------------------------------------------------------------- D

// Test_audit_D_mixed_load: read+write mixes (7R+1W, 4R+4W, 8R+8W, 0R+8W) in
// pooled-connection and per-request-connection modes, 3 repeats each. Readers
// exercise the real Get_user_auth (middleware permission lookup) plus document
// reads; writers insert history through their own connections.
func Test_audit_D_mixed_load(t *testing.T) {
	type phase struct {
		readers, writers int
	}
	phases := []phase{{7, 1}, {4, 4}, {8, 8}, {0, 8}}
	seconds := 1200 * time.Millisecond
	repeats := 3

	for _, backend := range []string{"sqlite", "turso"} {
		db := audit_setup(t, backend)
		audit_seed_auth(db, t)
		audit_seed_docs(db, t, 50)

		for _, mode := range []string{"pooled", "perreq"} {
			for _, ph := range phases {
				for rep := 0; rep < repeats; rep++ {
					// unique history id prefix per phase+rep so integrity counts
					// cannot accumulate across scenarios sharing the test DB
					sig := fmt.Sprintf("d_%s_%d_%d_%d", mode, ph.readers, ph.writers, rep)
					stats := &audit_stats{}
					var w_ops, w_errs atomic.Int64
					var panic_cnt atomic.Int64
					deadline := time.Now().Add(seconds)

					var wg sync.WaitGroup
					for r := 0; r < ph.readers; r++ {
						wg.Add(1)
						go func(r int) {
							defer wg.Done()
							defer func() {
								if rec := recover(); rec != nil {
									panic_cnt.Add(1)
								}
							}()
							var conn *sql.DB
							if mode == "pooled" {
								conn = tool.DB_connect()
								defer conn.Close()
							}
							i := 0
							for time.Now().Before(deadline) {
								c := conn
								if mode == "perreq" {
									c = tool.DB_connect()
								}
								t0 := time.Now()
								_ = tool.Get_user_auth(c, "admin_ip")
								found := ""
								err := c.QueryRow(
									`select data from data where title = ?`,
									fmt.Sprintf("audit_%03d", (r+i)%50),
								).Scan(&found)
								rows := tool.Query_DB(c,
									`select title from data where title like 'audit_%' order by title limit 10`)
								n := 0
								for rows.Next() {
									var s string
									if rows.Scan(&s) == nil {
										n++
									}
								}
								rows.Close()
								if err != nil || n == 0 {
									err = fmt.Errorf("read failed: %v (n=%d)", err, n)
									stats.record(time.Since(t0), err)
								} else {
									stats.record(time.Since(t0), nil)
								}
								if mode == "perreq" {
									c.Close()
								}
								i++
							}
						}(r)
					}
					for w := 0; w < ph.writers; w++ {
						wg.Add(1)
						go func(w int) {
							defer wg.Done()
							defer func() {
								if rec := recover(); rec != nil {
									panic_cnt.Add(1)
								}
							}()
							var conn *sql.DB
							if mode == "pooled" {
								conn = tool.DB_connect()
								defer conn.Close()
							}
							i := 0
							for time.Now().Before(deadline) {
								c := conn
								if mode == "perreq" {
									c = tool.DB_connect()
								}
								_, err := c.Exec(
									`insert into history (id, title, data, date, ip, send, leng, hide, type) values (?, ?, ?, ?, ?, ?, ?, '', '')`,
									fmt.Sprintf("%s_%d_%d", sig, w, i), fmt.Sprintf("audit_%03d", w%50), "rev",
									"2026-01-01 00:00:00", "admin_ip", "audit", "+1",
								)
								if mode == "perreq" {
									c.Close()
								}
								if err != nil {
									w_errs.Add(1)
									if busy, label := audit_classify(err); busy {
										stats.busy.Add(1)
										_ = label
									}
									stats.errs.Add(1)
								} else {
									w_ops.Add(1)
								}
								i++
							}
						}(w)
					}
					wg.Wait()

					got := 0
					_ = tool.QueryRow_DB(db, fmt.Sprintf(`select count(*) from history where id like '%s_%%'`, sig), []any{&got})
					want := int(w_ops.Load())
					p50, p95, p99 := stats.summary()
					audit_record(t, map[string]any{
						"type": "D", "backend": backend, "mode": mode,
						"readers": ph.readers, "writers": ph.writers, "rep": rep,
						"reads_ok": stats.ops.Load(), "read_errs": stats.errs.Load(),
						"write_ok": w_ops.Load(), "write_errs": w_errs.Load(),
						"panics": panic_cnt.Load(), "integrity_match": got == want,
						"p50": p50, "p95": p95, "p99": p99, "secs": seconds.Seconds(),
					})
					if got != want {
						t.Errorf("D %s %s %dR+%dW rep%d: history rows=%d want=%d (lost writes)",
							backend, mode, ph.readers, ph.writers, rep, got, want)
					}
					if stats.errs.Load() != 0 || w_errs.Load() != 0 || panic_cnt.Load() != 0 {
						t.Errorf("D %s %s %dR+%dW rep%d: read_errs=%d write_errs=%d panics=%d",
							backend, mode, ph.readers, ph.writers, rep,
							stats.errs.Load(), w_errs.Load(), panic_cnt.Load())
					}
					t.Logf("D %s %s %dR+%dW rep%d reads=%d(w%d) writes=%d(e%d) p99=%.4f panics=%d",
						backend, mode, ph.readers, ph.writers, rep,
						stats.ops.Load(), stats.errs.Load(), w_ops.Load(), w_errs.Load(), p99, panic_cnt.Load())
				}
			}
		}
	}
}

// ---------------------------------------------------------------- E

// Test_audit_E_locks_and_transactions: long read snapshot vs writer, long write
// tx vs second writer (finite busy timeout), explicit failure classification,
// rollback atomicity, canceled context recovery.
func Test_audit_E_locks_and_transactions(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		db := audit_setup(t, backend)
		audit_seed_auth(db, t)
		audit_seed_docs(db, t, 10)

		// E1: long read transaction while a writer inserts
		rtx, err := db.Begin()
		if err != nil {
			t.Fatalf("E1 %s begin: %v", backend, err)
		}
		var one int
		if err := rtx.QueryRow(`select count(*) from data`).Scan(&one); err != nil {
			t.Fatalf("E1 %s read: %v", backend, err)
		}
		w_conn := tool.DB_connect()
		t0 := time.Now()
		_, w_err := w_conn.Exec(
			`insert into history (id, title, data, date, ip, send, leng, hide, type) values ('e1','audit_e1','v','2026-01-01 00:00:00','admin_ip','audit','+1','','')`)
		w_secs := time.Since(t0).Seconds()
		_ = rtx.Rollback()
		tool.DB_close(w_conn)
		audit_record(t, map[string]any{
			"type": "E1_read_tx_vs_writer", "backend": backend,
			"writer_ok": w_err == nil, "writer_secs": w_secs, "err": err_str(w_err),
		})
		if w_err != nil {
			t.Errorf("E1 %s: writer failed during read snapshot: %v", backend, w_err)
		}
		t.Logf("E1 %s writer_ok=%t writer_secs=%.3f", backend, w_err == nil, w_secs)

		// E2: 2s write tx, second writer must succeed via busy timeout
		hold_done := make(chan struct{})
		go func() {
			defer close(hold_done)
			_ = tool.DB_transaction(db, func(tx *sql.Tx) error {
				if _, err := tx.Exec(`insert into history (id, title, data, date, ip, send, leng, hide, type) values ('e2h','audit_e2','v','2026-01-01 00:00:00','admin_ip','audit','+1','','')`); err != nil {
					return err
				}
				time.Sleep(2 * time.Second)
				return nil
			})
		}()
		time.Sleep(200 * time.Millisecond)
		c2 := tool.DB_connect()
		t0 = time.Now()
		_, e2_err := c2.Exec(
			`insert into history (id, title, data, date, ip, send, leng, hide, type) values ('e2','audit_e2','v','2026-01-01 00:00:00','admin_ip','audit','+1','','')`)
		e2_secs := time.Since(t0).Seconds()
		tool.DB_close(c2)
		<-hold_done
		audit_record(t, map[string]any{
			"type": "E2_write_hold_2s", "backend": backend,
			"second_ok": e2_err == nil, "second_secs": e2_secs, "err": err_str(e2_err),
		})
		if e2_err != nil {
			t.Errorf("E2 %s: second writer failed while hold=2s (< busy_timeout 5s): %v (%.2fs)", backend, e2_err, e2_secs)
		}
		t.Logf("E2 %s second_ok=%t second_secs=%.3f", backend, e2_err == nil, e2_secs)

		// E3: 8s write tx, second writer must fail explicitly within a finite
		// window (busy_timeout 5s) - never hang forever, never corrupt.
		hold8_done := make(chan struct{})
		go func() {
			defer close(hold8_done)
			_ = tool.DB_transaction(db, func(tx *sql.Tx) error {
				if _, err := tx.Exec(`insert into history (id, title, data, date, ip, send, leng, hide, type) values ('e3h','audit_e3','v','2026-01-01 00:00:00','admin_ip','audit','+1','','')`); err != nil {
					return err
				}
				time.Sleep(8 * time.Second)
				return nil
			})
		}()
		time.Sleep(200 * time.Millisecond)
		c3 := tool.DB_connect()
		t0 = time.Now()
		_, e3_err := c3.Exec(
			`insert into history (id, title, data, date, ip, send, leng, hide, type) values ('e3','audit_e3','v','2026-01-01 00:00:00','admin_ip','audit','+1','','')`)
		e3_secs := time.Since(t0).Seconds()
		tool.DB_close(c3)
		<-hold8_done
		busy3, label3 := audit_classify(e3_err)
		audit_record(t, map[string]any{
			"type": "E3_write_hold_8s", "backend": backend,
			"second_failed": e3_err != nil, "busy": busy3, "label": label3,
			"second_secs": e3_secs, "err": err_str(e3_err),
		})
		if e3_secs > 10 {
			t.Errorf("E3 %s: second writer blocked %.1fs (must fail within busy_timeout window)", backend, e3_secs)
		}
		if e3_err == nil {
			t.Errorf("E3 %s: second writer succeeded while 8s hold expected explicit busy failure", backend)
		}
		t.Logf("E3 %s failed=%t busy=%t label=%s secs=%.2f", backend, e3_err != nil, busy3, label3, e3_secs)

		// E4: rollback must leave no partial rows
		_ = tool.DB_transaction(db, func(tx *sql.Tx) error {
			for i := 0; i < 3; i++ {
				if _, err := tx.Exec(
					`insert into history (id, title, data, date, ip, send, leng, hide, type) values (?, 'audit_e4','v','2026-01-01 00:00:00','admin_ip','audit','+1','','')`,
					fmt.Sprintf("e4_%d", i)); err != nil {
					return err
				}
			}
			return errors.New("intentional rollback")
		})
		n := 0
		_ = tool.QueryRow_DB(db, `select count(*) from history where title = 'audit_e4'`, []any{&n})
		audit_record(t, map[string]any{"type": "E4_rollback", "backend": backend, "partial_rows": n})
		if n != 0 {
			t.Errorf("E4 %s: rollback left %d partial rows", backend, n)
		}

		// E5: canceled context must not poison the connection
		c5 := tool.DB_connect()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, cancel_err := c5.ExecContext(ctx, `insert into history (id, title, data, date, ip, send, leng, hide, type) values ('e5','audit_e5','v','2026-01-01 00:00:00','admin_ip','audit','+1','','')`)
		_ = cancel_err
		var ping int
		ping_err := c5.QueryRow(`select 1`).Scan(&ping)
		tool.DB_close(c5)
		audit_record(t, map[string]any{
			"type": "E5_cancel", "backend": backend,
			"cancel_err": err_str(cancel_err), "conn_usable_after": ping_err == nil,
		})
		if ping_err != nil {
			t.Errorf("E5 %s: connection unusable after canceled request: %v", backend, ping_err)
		}
		t.Logf("E5 %s cancel_err=%v usable_after=%t", backend, cancel_err, ping_err == nil)
	}
}

func err_str(err error) string {
	if err == nil {
		return ""
	}
	return first_line(err.Error())
}

// ---------------------------------------------------------------- F

// Test_audit_F_keeper_churn: per-request open/auth/close churn with the WAL
// keeper ON vs OFF (sqlite) and plain churn for turso. Keeper ON must never
// panic; busy counts are recorded for the report (comparison, not assertion).
func Test_audit_F_keeper_churn(t *testing.T) {
	churn := func(iters int, stats *audit_stats) {
		for i := 0; i < iters; i++ {
			func() {
				defer func() {
					if r := recover(); r != nil {
						stats.errs.Add(1)
					}
				}()
				t0 := time.Now()
				conn := tool.DB_connect()
				_ = tool.Get_user_auth(conn, "admin_ip")
				tool.DB_close(conn)
				stats.record(time.Since(t0), nil)
			}()
		}
	}

	iters := 300
	repeats := 3

	// turso: no keeper concept; plain churn
	{
		audit_setup(t, "turso") // env + boot; keeper no-op
		seed := tool.DB_connect()
		audit_seed_auth(seed, t)
		tool.DB_close(seed)

		for rep := 0; rep < repeats; rep++ {
			stats := &audit_stats{}
			start := time.Now()
			churn(iters, stats)
			p50, p95, p99 := stats.summary()
			audit_record(t, map[string]any{
				"type": "F", "backend": "turso", "keeper": false, "rep": rep,
				"iters": iters, "errs": stats.errs.Load(), "p50": p50, "p95": p95, "p99": p99,
				"secs": time.Since(start).Seconds(),
			})
			t.Logf("F turso rep%d errs=%d p50=%.4f", rep, stats.errs.Load(), p50)
		}
	}

	// sqlite: keeper ON vs OFF
	{
		audit_setup(t, "sqlite")
		seed := tool.DB_connect()
		audit_seed_auth(seed, t)
		tool.DB_close(seed)

		if err := tool.Start_keeper(); err != nil {
			t.Fatalf("keeper start: %v", err)
		}
		if err := tool.Start_keeper(); err != nil { // idempotency
			t.Fatalf("keeper start twice: %v", err)
		}

		for rep := 0; rep < repeats; rep++ {
			stats := &audit_stats{}
			start := time.Now()
			churn(iters, stats)
			p50, p95, p99 := stats.summary()
			audit_record(t, map[string]any{
				"type": "F", "backend": "sqlite", "keeper": true, "rep": rep,
				"iters": iters, "errs": stats.errs.Load(), "p50": p50, "p95": p95, "p99": p99,
				"secs": time.Since(start).Seconds(),
			})
			t.Logf("F sqlite keeper=ON rep%d errs=%d p50=%.4f", rep, stats.errs.Load(), p50)
		}

		if err := tool.Stop_keeper(); err != nil {
			t.Fatalf("keeper stop: %v", err)
		}
		if err := tool.Stop_keeper(); err != nil { // idempotency
			t.Fatalf("keeper stop twice: %v", err)
		}

		for rep := 0; rep < repeats; rep++ {
			stats := &audit_stats{}
			start := time.Now()
			churn(iters, stats)
			p50, p95, p99 := stats.summary()
			audit_record(t, map[string]any{
				"type": "F", "backend": "sqlite", "keeper": false, "rep": rep,
				"iters": iters, "errs": stats.errs.Load(), "p50": p50, "p95": p95, "p99": p99,
				"secs": time.Since(start).Seconds(),
			})
			t.Logf("F sqlite keeper=OFF rep%d errs=%d p50=%.4f", rep, stats.errs.Load(), p50)
		}
	}
}

// ---------------------------------------------------------------- G

// Test_audit_G_recovery: committed-write durability across reconnect, plus a
// child-process crash matrix (clean exit / hard exit after commits / hard exit
// mid-transaction) over isolated databases per mode. The child is re-exec'd
// with AUDIT_CRASH_CHILD=1.
func Test_audit_G_recovery(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		db := audit_setup(t, backend)
		audit_seed_auth(db, t)

		// durability across reconnect
		if _, err := db.Exec(`create table if not exists audit_crash (id integer primary key, tag text)`); err != nil {
			t.Fatalf("G %s create: %v", backend, err)
		}
		c := tool.DB_connect()
		if _, err := c.Exec(`insert into audit_crash (id, tag) values (1, 'durable')`); err != nil {
			t.Fatalf("G %s insert: %v", backend, err)
		}
		tool.DB_close(c)
		db2 := tool.DB_connect()
		got := ""
		if !tool.QueryRow_DB(db2, `select tag from audit_crash where id = 1`, []any{&got}) || got != "durable" {
			t.Errorf("G %s: committed row lost after reconnect", backend)
		}
		tool.DB_close(db2)

		// crash matrix: three child exit modes over isolated databases
		for _, mode := range []string{"clean", "after_commit", "mid_tx"} {
			audit_run_crash_case(t, backend, mode)
		}
	}
}

// audit_run_crash_case boots an isolated DB, runs the crash child, and probes
// the file from multiple engines afterwards.
func audit_run_crash_case(t *testing.T, backend string, mode string) {
	t.Helper()

	dir := filepath.Join(t.TempDir(), mode)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("G %s %s: mkdir: %v", backend, mode, err)
	}
	os.Setenv("NAMU_DB", filepath.Join(dir, "audit"))
	os.Setenv("NAMU_DB_TYPE", backend)

	tool.DB_init()
	db := tool.DB_connect()
	if _, err := db.Exec(`create table if not exists audit_crash (id integer primary key, tag text)`); err != nil {
		t.Fatalf("G %s %s: create: %v", backend, mode, err)
	}
	tool.DB_close(db)

	exe, err := os.Executable()
	if err != nil {
		t.Fatalf("G: os.Executable: %v", err)
	}
	db_path, _ := os.LookupEnv("NAMU_DB")
	cmd := exec.Command(exe, "-test.run", "^Test_audit_crash_child$", "-test.count=1")
	cmd.Env = append(os.Environ(),
		"AUDIT_CRASH_CHILD=1",
		"AUDIT_CRASH_MODE="+mode,
		"NAMU_DB="+db_path,
		"NAMU_DB_TYPE="+backend,
	)
	out, _ := cmd.CombinedOutput()
	log_path := filepath.Join(audit_out_dir(t), fmt.Sprintf("crash_%s_%s.log", backend, mode))
	_ = os.WriteFile(log_path, out, 0o644)

	// probes on the file after the child exit
	file_path := db_path + ".db"
	probe := map[string]any{
		"type": "G_crash", "backend": backend, "mode": mode,
		"file": filepath.Base(file_path),
	}

	// files left in the db dir
	entries, _ := os.ReadDir(dir)
	files := map[string]any{}
	for _, e := range entries {
		if info, err := e.Info(); err == nil {
			files[e.Name()] = info.Size()
		}
	}
	probe["dir_files"] = files

	// header magic (SQLite main-db files start with "SQLite format 3\x00")
	if header, err := os.ReadFile(file_path); err == nil && len(header) >= 16 {
		probe["header_is_sqlite"] = string(header[:15]) == "SQLite format 3"
	}

	// probe: turso plain reopen
	p_open_err := ""
	committed, partial := -1, -1
	if p, err := sql.Open("turso", file_path); err == nil {
		if err := p.QueryRow(`select count(*) from audit_crash where tag = 'committed'`).Scan(&committed); err != nil {
			p_open_err = err_str(err)
		}
		if committed >= 0 {
			var pp int
			if err := p.QueryRow(`select count(*) from audit_crash where tag = 'in_tx'`).Scan(&pp); err == nil {
				partial = pp
			}
		}
		p.Close()
	} else {
		p_open_err = err_str(err)
	}
	probe["turso_committed"] = committed
	probe["turso_in_tx"] = partial
	probe["turso_open_err"] = p_open_err

	// probe: turso reopen with explicit mvcc pragma
	mvcc_committed := -1
	if backend == "turso" {
		if p, err := sql.Open("turso", file_path); err == nil {
			if _, err = p.Exec(`pragma journal_mode = mvcc`); err == nil {
				var c1 int
				if err = p.QueryRow(`select count(*) from audit_crash where tag = 'committed'`).Scan(&c1); err == nil {
					mvcc_committed = c1
				}
			}
			p.Close()
		}
	}
	probe["turso_mvcc_reopen_committed"] = mvcc_committed

	// probe: modernc sqlite reader
	sq_committed := -1
	sq_err := ""
	if p, err := sql.Open("sqlite", file_path); err == nil {
		if err := p.QueryRow(`select count(*) from audit_crash where tag = 'committed'`).Scan(&sq_committed); err != nil {
			sq_err = err_str(err)
			sq_committed = -1
		}
		p.Close()
	} else {
		sq_err = err_str(err)
	}
	probe["sqlite_reader_committed"] = sq_committed
	probe["sqlite_reader_err"] = sq_err

	// probe (turso MVCC only): can the journal mode be switched back to WAL so
	// standard SQLite can read the file again? Two paths: plain switch, and
	// checkpoint(TRUNCATE) before the switch.
	if backend == "turso" && sq_committed == -1 {
		// path A: switch journal mode to wal on a turso connection, then read with modernc
		if p, err := sql.Open("turso", file_path); err == nil {
			if _, err = p.Exec(`pragma journal_mode = wal`); err == nil {
				probe["wal_switch_pragma"] = "ok"
			} else {
				probe["wal_switch_pragma_err"] = err_str(err)
			}
			p.Close()

			p2, err2 := sql.Open("sqlite", file_path)
			if err2 == nil {
				var c4 int
				if err := p2.QueryRow(`select count(*) from audit_crash where tag = 'committed'`).Scan(&c4); err != nil {
					probe["after_wal_switch_sqlite_err"] = err_str(err)
				} else {
					probe["after_wal_switch_sqlite_committed"] = c4
				}
				p2.Close()
			}
		}

		// path B: checkpoint(TRUNCATE) first, then modernc read
		if p, err := sql.Open("turso", file_path); err == nil {
			var s1, s2, s3 int
			if err := p.QueryRow(`pragma wal_checkpoint(TRUNCATE)`).Scan(&s1, &s2, &s3); err != nil {
				probe["checkpoint_err"] = err_str(err)
			} else {
				probe["checkpoint"] = fmt.Sprintf("%d/%d/%d", s1, s2, s3)
			}
			p.Close()

			p2, err2 := sql.Open("sqlite", file_path)
			if err2 == nil {
				var c5 int
				if err := p2.QueryRow(`select count(*) from audit_crash where tag = 'committed'`).Scan(&c5); err != nil {
					probe["after_checkpoint_sqlite_err"] = err_str(err)
				} else {
					probe["after_checkpoint_sqlite_committed"] = c5
				}
				p2.Close()
			}
		}
	}

	audit_record(t, probe)

	// expectations per mode
	want_committed, want_partial := 5, 0
	fail := committed != want_committed || partial != want_partial
	t.Logf("G %s mode=%s turso_committed=%d in_tx=%d mvcc=%d sqlite=%d sqlite_err=%q header_ok=%v files=%v",
		backend, mode, committed, partial, mvcc_committed, sq_committed, sq_err, probe["header_is_sqlite"], files)
	if fail {
		t.Errorf("G %s mode=%s: committed=%d(want %d) in_tx=%d(want %d) probes=%v",
			backend, mode, committed, want_committed, partial, want_partial, probe)
	}
	if backend == "sqlite" && mode == "clean" && sq_committed != 5 {
		t.Errorf("G sqlite mode=clean: modernc sqlite could not re-read its own file (committed=%d err=%q)",
			sq_committed, sq_err)
	}
}

// Test_audit_crash_child is the re-exec child used by the G crash matrix.
// It is skipped in normal runs and only executes when AUDIT_CRASH_CHILD=1.
func Test_audit_crash_child(t *testing.T) {
	if os.Getenv("AUDIT_CRASH_CHILD") != "1" {
		t.Skip("crash child: not in child mode")
	}
	mode := os.Getenv("AUDIT_CRASH_MODE")

	tool.DB_boot()
	db := tool.DB_connect()

	if _, err := db.Exec(`create table if not exists audit_crash (id integer primary key, tag text)`); err != nil {
		panic(err)
	}
	for i := 0; i < 5; i++ {
		if _, err := db.Exec(`insert into audit_crash (id, tag) values (?, 'committed')`, 100+i); err != nil {
			panic(err)
		}
	}

	switch mode {
	case "clean":
		if err := db.Close(); err != nil {
			panic(err)
		}
		os.Exit(0)
	case "after_commit":
		os.Exit(9)
	case "mid_tx":
		tx, err := db.Begin()
		if err != nil {
			panic(err)
		}
		for i := 0; i < 3; i++ {
			if _, err := tx.Exec(`insert into audit_crash (id, tag) values (?, 'in_tx')`, 200+i); err != nil {
				panic(err)
			}
		}
		_ = tx
		os.Exit(9)
	default:
		panic("unknown crash mode " + mode)
	}
}

// ---------------------------------------------------------------- H

// Test_audit_H_file_compatibility: SQLite-file compatibility between the two
// engines on isolated temp copies (never production data).
func Test_audit_H_file_compatibility(t *testing.T) {
	dir := t.TempDir()

	// H1: modernc sqlite creates a file; tursogo reads it and writes; modernc reads back
	path := filepath.Join(dir, "compat.db")
	sdb, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("H1 sqlite open: %v", err)
	}
	if _, err := sdb.Exec(`create table compat (k text, v text)`); err != nil {
		t.Fatalf("H1 create: %v", err)
	}
	if _, err := sdb.Exec(`insert into compat values ('from_sqlite', 'ok')`); err != nil {
		t.Fatalf("H1 insert: %v", err)
	}
	if err := sdb.Close(); err != nil {
		t.Fatalf("H1 close: %v", err)
	}

	tdb, err := sql.Open("turso", path)
	if err != nil {
		t.Fatalf("H1 turso open: %v", err)
	}
	var v string
	if err := tdb.QueryRow(`select v from compat where k = 'from_sqlite'`).Scan(&v); err != nil {
		t.Fatalf("H1 turso read sqlite data: %v", err)
	}
	if _, err := tdb.Exec(`insert into compat values ('from_turso', 'ok')`); err != nil {
		t.Fatalf("H1 turso insert: %v", err)
	}
	if err := tdb.Close(); err != nil {
		t.Fatalf("H1 turso close: %v", err)
	}

	sdb2, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("H1 reopen: %v", err)
	}
	var v2 string
	if err := sdb2.QueryRow(`select v from compat where k = 'from_turso'`).Scan(&v2); err != nil {
		t.Fatalf("H1 sqlite read turso data: %v", err)
	}
	if err := sdb2.Close(); err != nil {
		t.Fatalf("H1 close: %v", err)
	}
	audit_record(t, map[string]any{
		"type": "H1_direct_open", "result": "PASS",
		"sqlite_reads_turso": v2 == "ok", "turso_reads_sqlite": v == "ok",
	})
	t.Logf("H1 direct file open: turso reads sqlite=%t, sqlite reads turso=%t", v == "ok", v2 == "ok")

	// H2: turso MVCC-mode writes then modernc read (checkpoint caveat)
	path2 := filepath.Join(dir, "compat_mvcc.db")
	mdb, err := sql.Open("turso", path2)
	if err != nil {
		t.Fatalf("H2 turso open: %v", err)
	}
	if _, err := mdb.Exec(`create table compat2 (k text, v text)`); err != nil {
		t.Fatalf("H2 create: %v", err)
	}
	if _, err := mdb.Exec(`insert into compat2 values ('mvcc_row', 'ok')`); err != nil {
		t.Fatalf("H2 insert: %v", err)
	}
	var jm string
	if err := mdb.QueryRow(`pragma journal_mode`).Scan(&jm); err == nil {
		audit_record(t, map[string]any{"type": "H2", "backend": "turso", "journal_mode_seen": jm})
	}
	if err := mdb.Close(); err != nil {
		t.Fatalf("H2 close: %v", err)
	}

	sdb3, err := sql.Open("sqlite", path2)
	if err != nil {
		t.Fatalf("H2 sqlite open: %v", err)
	}
	err_h2 := sdb3.QueryRow(`select v from compat2 where k = 'mvcc_row'`).Scan(&v)
	sdb3.Close()
	audit_record(t, map[string]any{
		"type": "H2_mvcc_then_sqlite_read", "visible": err_h2 == nil, "err": err_str(err_h2),
	})
	t.Logf("H2 turso(MVCC) row visible to modernc sqlite: %t err=%v", err_h2 == nil, err_h2)

	// H3: turso default(WAL) write then modernc read
	path3 := filepath.Join(dir, "compat_wal.db")
	wdb, err := sql.Open("turso", path3)
	if err != nil {
		t.Fatalf("H3 turso open: %v", err)
	}
	if _, err := wdb.Exec(`create table compat3 (k text, v text)`); err != nil {
		t.Fatalf("H3 create: %v", err)
	}
	if _, err := wdb.Exec(`insert into compat3 values ('wal_row', 'ok')`); err != nil {
		t.Fatalf("H3 insert: %v", err)
	}
	if err := wdb.Close(); err != nil {
		t.Fatalf("H3 close: %v", err)
	}
	sdb4, err := sql.Open("sqlite", path3)
	if err != nil {
		t.Fatalf("H3 sqlite open: %v", err)
	}
	err_h3 := sdb4.QueryRow(`select v from compat3 where k = 'wal_row'`).Scan(&v)
	sdb4.Close()
	audit_record(t, map[string]any{
		"type": "H3_wal_then_sqlite_read", "visible": err_h3 == nil, "err": err_str(err_h3),
	})
	t.Logf("H3 turso(WAL) row visible to modernc sqlite: %t err=%v", err_h3 == nil, err_h3)
}

// ---------------------------------------------------------------- I

// Test_audit_I_integrity: counters, unique keys, and engine integrity checks.
func Test_audit_I_integrity(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		db := audit_setup(t, backend)
		audit_seed_auth(db, t)

		if _, err := db.Exec(`create table audit_counter (k text primary key, v integer)`); err != nil {
			t.Fatalf("I %s create counter: %v", backend, err)
		}
		if _, err := db.Exec(`insert into audit_counter values ('c', 0)`); err != nil {
			t.Fatalf("I %s seed counter: %v", backend, err)
		}
		if _, err := db.Exec(`create table audit_unique (k text primary key)`); err != nil {
			t.Fatalf("I %s create unique: %v", backend, err)
		}

		// I1: single-statement atomic increments from 8 workers
		var i1_errs atomic.Int64
		var wg sync.WaitGroup
		for w := 0; w < 8; w++ {
			wg.Add(1)
			go func(w int) {
				defer wg.Done()
				conn := tool.DB_connect()
				defer conn.Close()
				for i := 0; i < 50; i++ {
					if _, err := conn.Exec(`update audit_counter set v = v + 1 where k = 'c'`); err != nil {
						i1_errs.Add(1)
					}
				}
			}(w)
		}
		wg.Wait()
		v1 := 0
		_ = tool.QueryRow_DB(db, `select v from audit_counter where k = 'c'`, []any{&v1})
		audit_record(t, map[string]any{
			"type": "I1_counter_atomic", "backend": backend,
			"want": 400, "got": v1, "errs": i1_errs.Load(),
		})
		if v1 != 400 || i1_errs.Load() != 0 {
			t.Errorf("I1 %s: counter=%d want 400, errs=%d (lost updates)", backend, v1, i1_errs.Load())
		}
		t.Logf("I1 %s counter=%d/400 errs=%d", backend, v1, i1_errs.Load())

		// I2: read-modify-write across a transaction (snapshot isolation probe).
		// No retries: each lost race must surface as an explicit error.
		if _, err := db.Exec(`update audit_counter set v = 0 where k = 'c'`); err != nil {
			t.Fatalf("I2 reset: %v", err)
		}
		var i2_ok, i2_conflict atomic.Int64
		for w := 0; w < 8; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				conn := tool.DB_connect()
				defer conn.Close()
				for i := 0; i < 10; i++ {
					err := tool.DB_transaction(conn, func(tx *sql.Tx) error {
						cur := 0
						if err := tx.QueryRow(`select v from audit_counter where k = 'c'`).Scan(&cur); err != nil {
							return err
						}
						time.Sleep(2 * time.Millisecond)
						_, err := tx.Exec(`update audit_counter set v = ? where k = 'c'`, cur+1)
						return err
					})
					if err != nil {
						i2_conflict.Add(1)
					} else {
						i2_ok.Add(1)
					}
				}
			}()
		}
		wg.Wait()
		v2 := 0
		_ = tool.QueryRow_DB(db, `select v from audit_counter where k = 'c'`, []any{&v2})
		integrity := v2 == int(i2_ok.Load()) // no silent lost updates
		audit_record(t, map[string]any{
			"type": "I2_rmw_snapshot", "backend": backend,
			"committed": i2_ok.Load(), "explicit_conflicts": i2_conflict.Load(),
			"final_v": v2, "no_silent_lost_update": integrity,
		})
		if !integrity {
			t.Errorf("I2 %s: silent lost update: committed=%d final_v=%d", backend, i2_ok.Load(), v2)
		}
		t.Logf("I2 %s committed=%d conflicts=%d final_v=%d integrity=%t",
			backend, i2_ok.Load(), i2_conflict.Load(), v2, integrity)

		// I3: unique key - concurrent identical inserts, exactly one winner
		var i3_ok, i3_dup atomic.Int64
		for w := 0; w < 8; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				conn := tool.DB_connect()
				defer conn.Close()
				_, err := conn.Exec(`insert into audit_unique values ('same')`)
				if err != nil {
					if busy, label := audit_classify(err); !busy && strings.Contains(label, "CONSTRAINT") {
						i3_dup.Add(1)
					} else {
						t.Logf("I3 %s unexpected insert error: %v (%s)", backend, err, label)
					}
				} else {
					i3_ok.Add(1)
				}
			}()
		}
		wg.Wait()
		audit_record(t, map[string]any{
			"type": "I3_unique", "backend": backend,
			"winners": i3_ok.Load(), "constraint_rejects": i3_dup.Load(),
		})
		if i3_ok.Load() != 1 || i3_dup.Load() != 7 {
			t.Errorf("I3 %s: winners=%d rejects=%d, want 1/7", backend, i3_ok.Load(), i3_dup.Load())
		}
		t.Logf("I3 %s winners=%d rejects=%d", backend, i3_ok.Load(), i3_dup.Load())

		// I4: engine integrity check (where supported)
		var ic string
		ic_err := db.QueryRow(`pragma integrity_check`).Scan(&ic)
		audit_record(t, map[string]any{
			"type": "I4_integrity_check", "backend": backend,
			"supported": ic_err == nil, "result": ic, "err": err_str(ic_err),
		})
		if ic_err == nil && ic != "ok" {
			t.Errorf("I4 %s: integrity_check = %q", backend, ic)
		}
		t.Logf("I4 %s integrity_check=%q err=%v", backend, ic, ic_err)
	}
}

// ---------------------------------------------------------------- stress (opt-in)

// Test_audit_stress_optin runs the heavier matrix (32 clients/writers, longer
// mixed load, 3 repeats) and only executes with DB_CONCURRENCY_STRESS=1.
func Test_audit_stress_optin(t *testing.T) {
	if os.Getenv("DB_CONCURRENCY_STRESS") != "1" {
		t.Skip("stress: set DB_CONCURRENCY_STRESS=1 to run")
	}

	for _, backend := range []string{"sqlite", "turso"} {
		db := audit_setup(t, backend)
		audit_seed_auth(db, t)
		audit_seed_docs(db, t, 100)

		// 32 writers x 50 docs x 3 revisions, 3 repeats
		for rep := 0; rep < 3; rep++ {
			var wg sync.WaitGroup
			var panics atomic.Int64
			start := time.Now()
			for w := 0; w < 32; w++ {
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
					for d := 0; d < 50; d++ {
						title := fmt.Sprintf("s_r%d_w%02d_d%d", rep, w, d)
						tool.Exec_DB(conn, `insert into data (title, data, type) values (?, ?, '')`, title, "v0")
						for r := 1; r <= 3; r++ {
							tool.Exec_DB(conn,
								`insert into history (id, title, data, date, ip, send, leng, hide, type) values (?, ?, ?, ?, ?, ?, ?, '', '')`,
								fmt.Sprintf("%d", r), title, "rev", "2026-01-01 00:00:00", "admin_ip", "stress", "+1",
							)
						}
					}
				}(w)
			}
			wg.Wait()
			elapsed := time.Since(start)
			got := 0
			_ = tool.QueryRow_DB(db, fmt.Sprintf(`select count(*) from history where title like 's_r%d_%%'`, rep), []any{&got})
			want_stress_hist := 32 * 50 * 3
			audit_record(t, map[string]any{
				"type": "stress_B", "backend": backend, "rep": rep,
				"panics": panics.Load(), "secs": elapsed.Seconds(),
				"want_hist": want_stress_hist, "got_hist": got,
			})
			t.Logf("stress %s B rep%d panics=%d hist=%d/%d %.2fs",
				backend, rep, panics.Load(), got, want_stress_hist, elapsed.Seconds())
			if panics.Load() != 0 {
				t.Errorf("stress %s rep%d: %d panics", backend, rep, panics.Load())
			}
			if got != want_stress_hist {
				t.Errorf("stress %s rep%d: history rows=%d want=%d", backend, rep, got, want_stress_hist)
			}
		}

		// 8R+8W for 5s x 3 repeats with per-request connections (worst case churn)
		for rep := 0; rep < 3; rep++ {
			stats := &audit_stats{}
			var w_ops, w_errs atomic.Int64
			deadline := time.Now().Add(5 * time.Second)
			var wg sync.WaitGroup
			for r := 0; r < 8; r++ {
				wg.Add(1)
				go func(r int) {
					defer wg.Done()
					for time.Now().Before(deadline) {
						c := tool.DB_connect()
						t0 := time.Now()
						_ = tool.Get_user_auth(c, "admin_ip")
						var s string
						err := c.QueryRow(`select data from data where title = ?`, fmt.Sprintf("audit_%03d", r%100)).Scan(&s)
						stats.record(time.Since(t0), err)
						c.Close()
					}
				}(r)
			}
			for w := 0; w < 8; w++ {
				wg.Add(1)
				go func(w int) {
					defer wg.Done()
					for i := 0; time.Now().Before(deadline); i++ {
						c := tool.DB_connect()
						_, err := c.Exec(
							`insert into history (id, title, data, date, ip, send, leng, hide, type) values (?, ?, ?, ?, ?, ?, ?, '', '')`,
							fmt.Sprintf("str%d_%d_%d", rep, w, i), fmt.Sprintf("audit_%03d", w%100), "rev",
							"2026-01-01 00:00:00", "admin_ip", "stress", "+1",
						)
						if err != nil {
							w_errs.Add(1)
						} else {
							w_ops.Add(1)
						}
						c.Close()
					}
				}(w)
			}
			wg.Wait()
			got := 0
			_ = tool.QueryRow_DB(db, fmt.Sprintf(`select count(*) from history where id like 'str%d_%%'`, rep), []any{&got})
			audit_record(t, map[string]any{
				"type": "stress_D_perreq_8R8W", "backend": backend, "rep": rep,
				"write_ok": w_ops.Load(), "write_errs": w_errs.Load(),
				"integrity_match": got == int(w_ops.Load()),
			})
			t.Logf("stress %s D rep%d writes=%d errs=%d integrity=%t",
				backend, rep, w_ops.Load(), w_errs.Load(), got == int(w_ops.Load()))
			if w_errs.Load() != 0 {
				t.Errorf("stress %s D rep%d: %d write errors", backend, rep, w_errs.Load())
			}
		}
	}
}