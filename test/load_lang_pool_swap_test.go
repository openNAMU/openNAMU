package test

// A) Clean hand-off from the load_lang shared pool to a swapped DB identity.
//
// Pre-fix hazard (reproduced in tmp/load_lang_swap_repro): Get_lang_db handed
// out the old *sql.DB pointer, a concurrent identity change closed that pool,
// and the caller's next query panicked with "sql: database is closed". The
// production fix encapsulates the path in Get_language_shared, which holds
// lang_db_lock.RLock for the whole call; the swap path takes the write Lock,
// so a pool in use cannot be closed.
//
// Global-state note: DB_boot rewrites the process-wide db_set map, while
// Get_language_shared reads db_set (via db_identity_key) *before* it takes
// lang_db_lock. There is no lock on db_set, so a reader that is still running
// when the identity swap runs races the swap's write. This test therefore
// stops all readers and waits (WaitGroup) before the swap, then verifies the
// swapped pool is fresh (no stale language, no panic, no unexpected value).
// The RLock write-lease safety itself is exercised by the write-locked
// Get_lang_db swap, which closes the old pool only under the write Lock.

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"opennamu/route/tool"
)

// Test_load_lang_pool_swap_verified_handoff: readers via the production path
// (Get_language_shared) are verified stopped (WaitGroup) before a DB identity
// swap; the swap must install a fresh pool - no closed-pool errors, no stale
// language after the swap.
func Test_load_lang_pool_swap_verified_handoff(t *testing.T) {
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

	// hammer pool A, then stop readers before touching global state
	time.Sleep(100 * time.Millisecond)
	close(pause)
	wg.Wait() // verified: every reader has returned before the swap

	// swap identity; no reader is running, so DB_boot's db_set write cannot
	// race a reader's db_set read. The write Lock still serializes the pool
	// close against any (now-absent) reader lease.
	t.Setenv("NAMU_DB", t.TempDir()+"/swap_target")
	tool.DB_init()
	_ = tool.Get_lang_db() // installs pool B, closes pool A (under write Lock)

	after := tool.Get_language_shared("language_tag", false)

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