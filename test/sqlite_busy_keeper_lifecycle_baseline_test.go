package test

import (
	"database/sql"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"opennamu/route/tool"
)

// Test_keeper_lifecycle_baseline: IDENTICAL concurrent write+read load to
// Test_keeper_lifecycle_integrity (4 writers + 4 readers, each opening a fresh
// connection per iteration), but WITHOUT the keeper. Used to compare the before/
// after panic rate and confirm the remaining panics with the keeper are not
// caused by the keeper itself.
func Test_keeper_lifecycle_baseline(t *testing.T) {
	_, _ = setup_busy_env(t)
	// No keeper started (baseline).

	const nWriters, nReaders, iters = 4, 4, 25
	var wg sync.WaitGroup
	var panicC int64
	var writeOK, readOK int64
	var panicMsgs []string
	var panicMu sync.Mutex
	recordPanic := func(wid, j int, rec any) {
		atomic.AddInt64(&panicC, 1)
		panicMu.Lock()
		panicsMsgs := &panicMsgs
		*panicsMsgs = append(*panicsMsgs, fmt.Sprintf("w%d-%d: %v", wid, j, rec))
		panicMu.Unlock()
	}

	seedDB := tool.DB_connect()
	seed_auth(seedDB, t)
	tool.DB_close(seedDB)

	for i := 0; i < nWriters; i++ {
		wg.Add(1)
		go func(wid int) {
			defer wg.Done()
			for j := 0; j < iters; j++ {
				db := tool.DB_connect()
				func() {
					defer func() {
						if rec := recover(); rec != nil {
							recordPanic(wid, j, rec)
						}
					}()
					if e := tool.DB_transaction(db, func(tx *sql.Tx) error {
						_, ex := tx.Exec("insert into data (title) values (?)",
							fmt.Sprintf("lw-%d-%d", wid, j))
						return ex
					}); e != nil {
						return
					}
					atomic.AddInt64(&writeOK, 1)
				}()
				tool.DB_close(db)
			}
		}(i)
	}

	for i := 0; i < nReaders; i++ {
		wg.Add(1)
		go func(rid int) {
			defer wg.Done()
			for j := 0; j < iters; j++ {
				db := tool.DB_connect()
				func() {
					defer func() {
						if rec := recover(); rec != nil {
							recordPanic(rid, j, rec)
						}
					}()
					_ = tool.Get_user_auth(db, fmt.Sprintf("lr-%d-%d", rid, j))
					atomic.AddInt64(&readOK, 1)
				}()
				tool.DB_close(db)
			}
		}(i)
	}
	wg.Wait()

	t.Logf("BASELINE (no keeper): writes=%d reads=%d panic=%d", writeOK, readOK, panicC)
	if panicMsgs != nil {
		n := len(panicMsgs)
		if n > 10 {
			n = 10
		}
		t.Logf("panic messages (first %d): %v", n, panicMsgs[:n])
	}
}
