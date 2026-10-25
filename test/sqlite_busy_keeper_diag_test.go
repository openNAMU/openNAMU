package test

import (
	"database/sql"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"opennamu/route/tool"
)

// Test_keeper_panic_diag: focused diagnostic that reproduces the concurrent write+read
// load of the lifecycle test with the keeper up, using the real production paths
// (DB_transaction for writes, Get_user_auth for reads), and captures the actual
// panic/error messages. A goroutine-level recover() also catches panics from
// DB_connect()/DB_close() which sit outside the per-iteration IIFE.
func Test_keeper_panic_diag(t *testing.T) {
	_, _ = setup_busy_env(t)
	if err := tool.Start_keeper(); err != nil {
		t.Fatalf("keeper start: %v", err)
	}
	defer tool.Stop_keeper()

	seedDB := tool.DB_connect()
	seed_auth(seedDB, t)
	tool.DB_close(seedDB)

	const nWorkers, iters = 8, 25
	var wg sync.WaitGroup
	var panicC int64
	var errC int64
	var msgs []string
	var mu sync.Mutex
	addMsg := func(m string) {
		mu.Lock()
		msgs = append(msgs, m)
		mu.Unlock()
	}
	// record a panic and remember the iteration that triggered it.
	panik := func(wid, j int, rec any) {
		atomic.AddInt64(&panicC, 1)
		addMsg(fmt.Sprintf("w%d-%d PANIC: %v", wid, j, rec))
	}

	for i := 0; i < nWorkers; i++ {
		wg.Add(1)
		go func(wid int) {
			defer wg.Done()
			// Goroutine-level recover: catches panics in DB_connect()/DB_close() too.
			defer func() {
				if rec := recover(); rec != nil {
					addMsg(fmt.Sprintf("w%d GO RUNTIME PANIC (loop exit): %v", wid, rec))
				}
			}()
			for j := 0; j < iters; j++ {
				db := tool.DB_connect()
				isReader := wid >= nWorkers/2
				func() {
					defer func() {
						if rec := recover(); rec != nil {
							panik(wid, j, rec)
						}
					}()
					if isReader {
						_ = tool.Get_user_auth(db, "127.0.0.1")
						_ = tool.Get_user_auth(db, "198.104.100.2")
						_ = tool.Get_user_auth(db, "200.20.136.204")
						_ = tool.Get_user_auth(db, "192.168.0.5")
						_ = tool.Get_user_auth(db, "255.255.255.255")
						return
					}
					if e := tool.DB_transaction(db, func(tx *sql.Tx) error {
						_, ex := tx.Exec("insert into data (title) values (?)",
							fmt.Sprintf("w%d-%d", wid, j))
						return ex
					}); e != nil {
						atomic.AddInt64(&errC, 1)
						addMsg(fmt.Sprintf("w%d-%d err: %v", wid, j, e))
					}
				}()
				tool.DB_close(db)
			}
		}(i)
	}
	wg.Wait()

	t.Logf("panics=%d non-zero-write-errors=%d recorded-msgs=%d", panicC, errC, len(msgs))
	n := len(msgs)
	if n > 60 {
		n = 60
	}
	for i := 0; i < n; i++ {
		t.Logf("  %s", msgs[i])
	}
}
