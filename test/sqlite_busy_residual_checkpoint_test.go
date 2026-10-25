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
// goroutine that periodically runs `PRAGMA checkpoint`. This isolates whether an
// explicit checkpoint while readers/writers run produces the residual SQLITE_BUSY
// (5/261) panics.
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

	stopC := make(chan struct{})
	go func() {
		defer close(stopC)
		ckDB := openWithDSN(dbFile, false)
		defer ckDB.Close()
		ticker := time.NewTicker(25 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			_, _ = ckDB.Exec("PRAGMA checkpoint")
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
	close(stopC)
	t.Logf("[checkpoint] panics=%d writer_errors=%d", panicC, errC)
}
