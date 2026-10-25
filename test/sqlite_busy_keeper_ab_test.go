package test

import (
	"database/sql"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"opennamu/route/tool"
)

// abCounts aggregates churn results across A/B repeats.
type abCounts struct {
	total, busyPrep, busyScan, other, panic int64
}

func (c *abCounts) add(r *churnResult) {
	c.total += r.total
	c.busyPrep += r.busyPrep
	c.busyScan += r.busyScan
	c.other += r.other
	c.panic += r.panic
}

// runABChurn runs the per-request churn (open/work/close) on a fresh isolated DB.
// mode=true starts the long-lived keeper (patched behavior); mode=false leaves it
// off (baseline / existing behavior). The keeper is always stopped when the run ends.
func runABChurn(t *testing.T, mode bool, writerShare int, authMode bool, iters int) *churnResult {
	_, _ = setup_busy_env(t)
	seedDB := tool.DB_connect()
	seed_auth(seedDB, t)
	tool.DB_close(seedDB)

	var stopKeeper func()
	if mode {
		if err := tool.Start_keeper(); err != nil {
			t.Fatalf("keeper start: %v", err)
		}
		stopKeeper = func() {
			if err := tool.Stop_keeper(); err != nil {
				t.Errorf("keeper stop: %v", err)
			}
		}
	}
	defer func() {
		if stopKeeper != nil {
			stopKeeper()
		}
	}()

	return runChurn(t, 16, iters, writerShare, authMode)
}

// heavyAB runs the per-request churn A/B (baseline no-keeper vs patched keeper) for
// repeats of workers×iters, with the given writerShare and authMode. It returns the
// aggregated counts for both the baseline and patched runs.
func heavyAB(t *testing.T, workers, iters, repeats, writerShare int, authMode bool) (abCounts, abCounts) {
	var base, patch abCounts
	for rep := 1; rep <= repeats; rep++ {
		base.add(runABChurn(t, false, writerShare, authMode, iters))
		patch.add(runABChurn(t, true, writerShare, authMode, iters))
	}
	return base, patch
}

// logAB prints the baseline/patched comparison for an A/B run.
func logAB(t *testing.T, label string, base, patch abCounts) {
	t.Logf("%s baseline: total=%d busy@prepare=%d busy@scan=%d other=%d panic=%d",
		label, base.total, base.busyPrep, base.busyScan, base.other, base.panic)
	t.Logf("%s patched : total=%d busy@prepare=%d busy@scan=%d other=%d panic=%d",
		label, patch.total, patch.busyPrep, patch.busyScan, patch.other, patch.panic)
}

// Test_keeper_ab_read_churn: A/B read-only per-request churn. The patched (keeper)
// run must not produce a busy@prepare or panic, because the keeper prevents the
// last-connection-close WAL cleanup.
func Test_keeper_ab_read_churn(t *testing.T) {
	const iters, repeats = 300, 3
	var base, patch abCounts
	start := time.Now()
	for rep := 1; rep <= repeats; rep++ {
		base.add(runABChurn(t, false, 0, false, iters))
		patch.add(runABChurn(t, true, 0, false, iters))
	}
	elapsed := time.Since(start)
	t.Logf("read-churn A/B elapsed=%s iters=%d repeats=%d", elapsed, iters, repeats)
	t.Logf("  baseline: total=%d busy@prepare=%d busy@scan=%d other=%d panic=%d",
		base.total, base.busyPrep, base.busyScan, base.other, base.panic)
	t.Logf("  patched : total=%d busy@prepare=%d busy@scan=%d other=%d panic=%d",
		patch.total, patch.busyPrep, patch.busyScan, patch.other, patch.panic)

	if patch.panic != 0 {
		t.Fatalf("patched keeper read churn panicked: %d", patch.panic)
	}
	if patch.busyPrep != 0 {
		t.Errorf("patched keeper read churn still had busy@prepare=%d (want 0)", patch.busyPrep)
	}
	if patch.busyScan != 0 {
		t.Errorf("patched keeper read churn still had busy@scan=%d (want 0)", patch.busyScan)
	}
}

// Test_keeper_ab_auth_churn: A/B through the real Get_user_auth (QueryRow_DB)
// middleware path — the exact path that produced the production panic.
func Test_keeper_ab_auth_churn(t *testing.T) {
	const iters, repeats = 200, 3
	var base, patch abCounts
	start := time.Now()
	for rep := 1; rep <= repeats; rep++ {
		base.add(runABChurn(t, false, 0, true, iters))
		patch.add(runABChurn(t, true, 0, true, iters))
	}
	elapsed := time.Since(start)
	t.Logf("auth-churn A/B elapsed=%s iters=%d repeats=%d", elapsed, iters, repeats)
	t.Logf("  baseline: total=%d busy@prepare=%d busy@scan=%d other=%d panic=%d",
		base.total, base.busyPrep, base.busyScan, base.other, base.panic)
	t.Logf("  patched : total=%d busy@prepare=%d busy@scan=%d other=%d panic=%d",
		patch.total, patch.busyPrep, patch.busyScan, patch.other, patch.panic)

	if patch.panic != 0 {
		t.Fatalf("patched keeper auth churn panicked: %d", patch.panic)
	}
	if patch.busyPrep != 0 {
		t.Errorf("patched keeper auth churn still had busy@prepare=%d (want 0)", patch.busyPrep)
	}
	if patch.busyScan != 0 {
		t.Errorf("patched keeper auth churn still had busy@scan=%d (want 0)", patch.busyScan)
	}
}

// Test_keeper_pragmas: with the keeper active, the DB is still in WAL with
// busy_timeout=5000 and the per-request connection path still works.
func Test_keeper_pragmas(t *testing.T) {
	_, _ = setup_busy_env(t)
	if err := tool.Start_keeper(); err != nil {
		t.Fatalf("keeper start: %v", err)
	}
	defer tool.Stop_keeper()

	db := tool.DB_connect()
	defer tool.DB_close(db)
	var journal, bt string
	if err := db.QueryRow("PRAGMA journal_mode").Scan(&journal); err != nil {
		t.Fatalf("journal_mode: %v", err)
	}
	if err := db.QueryRow("PRAGMA busy_timeout").Scan(&bt); err != nil {
		t.Fatalf("busy_timeout: %v", err)
	}
	t.Logf("keeper active: journal_mode=%q busy_timeout=%q", journal, bt)
	if journal != "wal" {
		t.Errorf("expected journal_mode=wal got %q", journal)
	}
	if bt != "5000" {
		t.Errorf("expected busy_timeout=5000 got %q", bt)
	}
}

// Test_keeper_read_parallelism: with the keeper active, concurrent readers and
// writers complete within a bounded time with no panic. This confirms the keeper
// does not serialize reads (no whole-DB lock / single-connection forcing) and reads
// are not blocked by the keeper or by writers.
func Test_keeper_read_parallelism(t *testing.T) {
	_, _ = setup_busy_env(t)
	seedDB := tool.DB_connect()
	seed_auth(seedDB, t)
	if _, err := seedDB.Exec("insert into data (title) values ('Para')"); err != nil {
		tool.DB_close(seedDB)
		t.Fatalf("seed para: %v", err)
	}
	tool.DB_close(seedDB)

	if err := tool.Start_keeper(); err != nil {
		t.Fatalf("keeper start: %v", err)
	}
	defer tool.Stop_keeper()

	const readers, writers, iters = 8, 2, 64
	var wg sync.WaitGroup
	var panicCount int64
	var readOK int64
	start := time.Now()
	for i := 0; i < readers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iters; j++ {
				db := tool.DB_connect()
				func() {
					defer func() {
						if rec := recover(); rec != nil {
							atomic.AddInt64(&panicCount, 1)
						}
					}()
					_ = tool.Get_user_auth(db, fmt.Sprintf("r%d-%d", i, j))
					atomic.AddInt64(&readOK, 1)
				}()
				tool.DB_close(db)
			}
		}()
	}
	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iters; j++ {
				db := tool.DB_connect()
				_ = tool.DB_transaction(db, func(tx *sql.Tx) error {
					_, e := tx.Exec("insert into data (title) values (?)", fmt.Sprintf("w%d-%d", i, j))
					return e
				})
				tool.DB_close(db)
			}
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)
	t.Logf("keeper read parallelism readers=%d writers=%d iters=%d elapsed=%s readOK=%d panic=%d",
		readers, writers, iters, elapsed, readOK, panicCount)
	if panicCount != 0 {
		t.Errorf("keeper read parallelism panicked: %d", panicCount)
	}
	if elapsed > 60*time.Second {
		t.Errorf("keeper read parallelism hung: %s (want < 60s)", elapsed)
	}
}
