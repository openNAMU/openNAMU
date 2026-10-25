package test

import (
	"database/sql"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"opennamu/route/tool"
)

// Test_keeper_lifecycle_diag: identical structure to Test_keeper_lifecycle_integrity
// but captures the actual panic message from the IIFE recover(), to identify the
// source of the panics (writer DB_transaction vs reader Get_user_auth).
func Test_keeper_lifecycle_diag(t *testing.T) {
	_, _ = setup_busy_env(t)
	if err := tool.Start_keeper(); err != nil {
		t.Fatalf("keeper start: %v", err)
	}
	defer tool.Stop_keeper()

	const nWriters, nReaders, iters = 4, 4, 25
	var wg sync.WaitGroup
	var panicC int64
	var panicsMsgs []string
	var mu sync.Mutex
	addPanicMsg := func(wid, j int, rec any) {
		atomic.AddInt64(&panicC, 1)
		mu.Lock()
		panicsMsgs = append(panicsMsgs, fmt.Sprintf("w%d-%d: %v", wid, j, rec))
		mu.Unlock()
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
							addPanicMsg(wid, j, rec)
						}
					}()
					if e := tool.DB_transaction(db, func(tx *sql.Tx) error {
						_, ex := tx.Exec("insert into data (title) values (?)",
							fmt.Sprintf("lw-%d-%d", wid, j))
						return ex
					}); e != nil {
						return
					}
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
							addPanicMsg(rid, j, rec)
						}
					}()
					_ = tool.Get_user_auth(db, fmt.Sprintf("lr-%d-%d", rid, j))
				}()
				tool.DB_close(db)
			}
		}(i)
	}
	wg.Wait()

	t.Logf("panics=%d", panicC)
	n := len(panicsMsgs)
	if n > 40 {
		n = 40
	}
	for i := 0; i < n; i++ {
		t.Logf("  %s", panicsMsgs[i])
	}
}
