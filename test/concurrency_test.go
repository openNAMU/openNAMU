package test

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"opennamu/route"
	"opennamu/route/tool"
)

// Test_long_write_lock_finish_or_clean_fail races a rebuild against a
// long-held write transaction and verifies the safety property: the rebuild
// must terminate (either complete or clean-fail) within a bounded time, with
// no panic.
//
// BEGIN opens a deferred transaction (no write lock is taken), so the rebuild's
// BEGIN can "succeed" while a concurrent writer still holds the write lock;
// the rebuild's first write (DELETE) then fails on the lock upgrade, and the
// (Stage 4) no-callback-retry policy makes it clean-fail. Conversely, if the
// lock is released in time, the rebuild finishes. Both outcomes are correct
// under the new policy; the invariant is bounded termination + no panic + no
// corruption.
func Test_long_write_lock_finish_or_clean_fail(t *testing.T) {
	cases := []struct {
		name string
		hold time.Duration
	}{
		{"short_hold", 250 * time.Millisecond},
		{"long_hold", 2 * time.Second},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			db := Setup_rebuild_db(t)
			defer tool.DB_close(db)

			admin := tool.Config{IP: "admin_ip"}

			var wg sync.WaitGroup
			var rebuildResponse string
			var rebuildPanic string
			start := time.Now()

			// Concurrent write: hold a write transaction for `hold`.
			wg.Add(1)
			go func() {
				defer wg.Done()
				tx, err := db.Begin()
				if err != nil {
					t.Errorf("writer begin: %v", err)
					return
				}
				defer tx.Rollback()
				_, _ = tx.Exec(`update data set data = data || 'lock' where title = 'Example'`)
				time.Sleep(c.hold)
			}()

			// Rebuild races the writer.
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer func() {
					if r := recover(); r != nil {
						rebuildPanic = fmt.Sprintf("%v", r)
					}
				}()
				out := route.Api_setting_backlink_reset_post_internal(admin, "fast")
				if resp, ok := out["response"].(string); ok {
					rebuildResponse = resp
				}
			}()

			wg.Wait()
			elapsed := time.Since(start)
			t.Logf("response=%q elapsed=%s", rebuildResponse, elapsed)

			if rebuildPanic != "" {
				t.Fatalf("rebuild panicked: %s (elapsed=%s)", rebuildPanic, elapsed)
			}
			if rebuildResponse == "" || rebuildResponse == "require auth" {
				t.Fatalf("rebuild did not reach a terminal response (elapsed=%s)", elapsed)
			}
			// Bounded termination: must not hang. The rebuild is bounded by the
			// Begin retry + callback work (or by the writer hold); a generous
			// cap catches any infinite-retry regression.
			if elapsed > 8*time.Second {
				t.Errorf("rebuild hung: elapsed=%s (want bounded)", elapsed)
			}
		})
	}
}

// Test_multi_write runs N concurrent rebuild workers against one DB and asserts
// that none panic and the total completes within a bounded window. Because each
// rebuild's first write cannot acquire the lock while a peer holds it, some
// workers will clean-fail (response "error") under contention; that is a
// correct terminal outcome. The invariants are: no panic, bounded completion,
// and at least one worker reaching a terminal state.
func Test_multi_write(t *testing.T) {
	for _, workers := range []int{10, 30, 50, 100} {
		t.Run("workers="+strconv.Itoa(workers), func(t *testing.T) {
			db := Setup_rebuild_db(t)
			defer tool.DB_close(db)

			admin := tool.Config{IP: "admin_ip"}

			var wg sync.WaitGroup
			var okCount, errCount, panicCount int64
			start := time.Now()
			for i := 0; i < workers; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					defer func() {
						if r := recover(); r != nil {
							atomic.AddInt64(&panicCount, 1)
						}
					}()
					out := route.Api_setting_backlink_reset_post_internal(admin, "fast")
					resp, _ := out["response"].(string)
					switch resp {
					case "ok", "partial_failure", "failure":
						atomic.AddInt64(&okCount, 1)
					default:
						atomic.AddInt64(&errCount, 1)
					}
				}()
			}
			wg.Wait()
			elapsed := time.Since(start)
			t.Logf("workers=%d ok=%d err=%d panic=%d elapsed=%s",
				workers, okCount, errCount, panicCount, elapsed)
			if panicCount != 0 {
				t.Fatalf("rebuild panicked under %d workers: %d", workers, panicCount)
			}
			// Bounded: must not hang. With a 500ms Begin retry and fast per-rebuild
			// work, all workers terminate within ~ a few seconds even under 100
			//-way contention.
			if elapsed > 60*time.Second {
				t.Fatalf("multi_write %d hung: elapsed=%s", workers, elapsed)
			}
			if okCount+errCount == 0 {
				t.Fatalf("multi_write %d: no worker reached a terminal state (ok=%d err=%d)", workers, okCount, errCount)
			}
		})
	}
}
