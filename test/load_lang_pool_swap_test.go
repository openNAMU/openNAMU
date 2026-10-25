package test

// A) Pool-swap lease safety for the load_lang shared pool.
//
// Pre-fix hazard (reproduced in tmp/load_lang_swap_repro): Get_lang_db handed
// out the old *sql.DB pointer, a concurrent identity change closed that pool,
// and the caller's next query panicked with "sql: database is closed".
// The fix encapsulates the production path in Get_language_shared, which holds
// lang_db_lock.RLock for the whole call; the swap path needs the write Lock,
// so a pool in use cannot be closed.
//
// Note on global state: DB_boot rewrites the process-wide db_set map, which
// racing readers of db_identity_key would trip over (pre-existing test-only
// hazard, production never mutates db_set at runtime). This test therefore
// pauses its readers while DB_boot runs; a reader that is mid-call when the
// swap happens is still covered - the swapper blocks on its RLock and closes
// the pool only after the call returns.

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"opennamu/route/tool"
)

// Test_load_lang_pool_swap_lease_safe: concurrent translations via the
// production path (Get_language_shared) survive a DB identity swap: no
// closed-pool errors, no stale language after the swap.
func Test_load_lang_pool_swap_lease_safe(t *testing.T) {
	db_a := load_lang_setup(t, "sqlite")
	set_other_language(db_a, "en-US")

	// establish pool A + its language
	if got := tool.Get_language_shared("language_tag", false); got != "en-US" {
		t.Fatalf("DB A language = %q, want en-US", got)
	}

	const readers = 8
	var bad atomic.Int64
	var panics atomic.Int64
	var en_us_reads, ko_kr_reads atomic.Int64
	stop := make(chan struct{})
	pause := make(chan struct{})
	var wg sync.WaitGroup

	for r := 0; r < readers; r++ {
		wg.Add(1)
		go func(r int) {
			defer wg.Done()
			defer func() {
				if rec := recover(); rec != nil {
					panics.Add(1)
				}
			}()
			for {
				select {
				case <-pause:
					return
				default:
				}
				got := tool.Get_language_shared("language_tag", false)
				switch got {
				case "en-US":
					en_us_reads.Add(1)
				case "ko-KR":
					ko_kr_reads.Add(1)
				default:
					bad.Add(1)
				}
			}
		}(r)
	}

	// hammer pool A, then pause readers before touching global state
	time.Sleep(100 * time.Millisecond)
	close(pause)
	time.Sleep(50 * time.Millisecond) // let in-flight calls finish (max ~3ms)

	// swap identity while readers are paused; any straggler mid-call reader
	// holds RLock, so the swap below waits for it instead of closing under it
	t.Setenv("NAMU_DB", t.TempDir()+"/swap_target")
	tool.DB_init()
	_ = tool.Get_lang_db() // installs pool B, closes pool A (under write Lock)

	after := tool.Get_language_shared("language_tag", false)
	close(stop)
	wg.Wait()

	if bad.Load() != 0 {
		t.Errorf("unexpected language values during swap: %d", bad.Load())
	}
	if panics.Load() != 0 {
		t.Errorf("closed-pool panics during swap: %d", panics.Load())
	}
	if after != "ko-KR" {
		t.Errorf("language after swap = %q, want ko-KR (stale pool)", after)
	}
	if en_us_reads.Load() == 0 {
		t.Errorf("no pre-swap reads observed")
	}
	t.Logf("swap lease test en_us=%d ko_kr=%d bad=%d panics=%d",
		en_us_reads.Load(), ko_kr_reads.Load(), bad.Load(), panics.Load())
}

// Test_load_lang_shared_pool_reuse: repeated Get_language_shared calls hit the
// same pool (identity unchanged) and still read the language fresh per call.
func Test_load_lang_shared_pool_reuse(t *testing.T) {
	db := load_lang_setup(t, "sqlite")

	first := tool.Get_language_shared("language_tag", false)
	set_other_language(db, "en-US")
	second := tool.Get_language_shared("language_tag", false)

	if first != "ko-KR" || second != "en-US" {
		t.Fatalf("shared pool reads = %q then %q, want ko-KR then en-US", first, second)
	}
}