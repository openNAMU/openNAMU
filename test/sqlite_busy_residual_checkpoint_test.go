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

// Test_measure_checkpoint: 4W+4R keeper-on, per-request churn, plus a background
// goroutine that periodically runs `PRAGMA wal_checkpoint(PASSIVE)`. This isolates
// whether an explicit WAL checkpoint while readers/writers run produces SQLITE_BUSY.
func Test_measure_checkpoint(t *testing.T) {
	_, dbFile := setup_busy_env(t)
	seedDB := tool.DB_connect()
	seed_auth(seedDB, t)
	tool.DB_close(seedDB)
	if err := tool.Start_keeper(); err != nil {
		t.Fatalf("keeper: %v", err)
	}
	defer tool.Stop_keeper()
	if err := ensureWAL(dbFile); err != nil {
		t.Fatalf("ensureWAL: %v", err)
	}

	stop_checkpoint := make(chan struct{})
	checkpoint_done := make(chan struct{})
	var checkpoint_calls, checkpoint_busy, checkpoint_errors, checkpointed_frames int64
	go func() {
		defer close(checkpoint_done)
		ckDB := openWithDSN(dbFile, false)
		defer ckDB.Close()
		ticker := time.NewTicker(25 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-stop_checkpoint:
				return
			case <-ticker.C:
				var busy, log_frames, checkpointed int64
				if err := ckDB.QueryRow(`PRAGMA wal_checkpoint(PASSIVE)`).Scan(&busy, &log_frames, &checkpointed); err != nil {
					atomic.AddInt64(&checkpoint_errors, 1)
					continue
				}
				atomic.AddInt64(&checkpoint_calls, 1)
				if busy != 0 {
					atomic.AddInt64(&checkpoint_busy, 1)
				}
				atomic.AddInt64(&checkpointed_frames, checkpointed)
			}
		}
	}()

	var panicC, errC int64
	var wg sync.WaitGroup
	const nWriters, nReaders, iters = 4, 4, 25
	runID := time.Now().UnixNano() % 1000000000

	work := func(wid, role int) {
		for iter := 0; iter < iters; iter++ {
			db := openWithDSN(dbFile, true)
			func() {
				defer func() {
					if rec := recover(); rec != nil {
						atomic.AddInt64(&panicC, 1)
						t.Logf("PANIC run=%d wid=%d role=%d iter=%d: %v", runID, wid, role, iter, rec)
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
						t.Logf("WRITER ERR wid=%d iter=%d: %v", wid, iter, e)
					}
				}
			}()
			_ = db.Close()
		}
	}
	for i := 0; i < nWriters; i++ {
		wg.Add(1)
		go func(wid int) {
			defer wg.Done()
			work(wid, 0)
		}(i)
	}
	for i := 0; i < nReaders; i++ {
		wg.Add(1)
		go func(wid int) {
			defer wg.Done()
			work(wid, 1)
		}(i)
	}
	wg.Wait()
	close(stop_checkpoint)
	<-checkpoint_done
	t.Logf("[checkpoint] calls=%d busy_results=%d errors=%d checkpointed_frames=%d panics=%d writer_errors=%d",
		checkpoint_calls, checkpoint_busy, checkpoint_errors, checkpointed_frames, panicC, errC)
	if checkpoint_calls == 0 {
		t.Fatal("checkpoint worker did not execute a WAL checkpoint")
	}
}
