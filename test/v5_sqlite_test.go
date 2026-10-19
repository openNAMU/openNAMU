package test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"opennamu/route"
	"opennamu/route/tool"
)

// Test_v5_sqlite_read_during_write runs the V5 external-host read paths while a
// write transaction holds the lock, and asserts the reads complete within a
// bounded time with a valid response. WAL read/write concurrency means the read
// path must not stall on the write lock; a hang or a panic (a 500-style
// failure) indicates a regression.
func Test_v5_sqlite_read_during_write(t *testing.T) {
	db := Setup_query_api_db(t)
	defer tool.DB_close(db)

	// Bind admin_ip to the owner group (grants setting_backlink, needed by
	// route.Api_external_all_hosts) — mirrors the existing query-API tests.
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed %v: %v", q, err)
		}
	}
	exec("insert into user_set (name, id, data) values ('acl','admin_ip','owner')")

	anon := tool.Config{}
	admin := tool.Config{IP: "admin_ip"}
	hold := 2 * time.Second

	var wg sync.WaitGroup
	var panicCount int64
	start := time.Now()

	// Concurrent writer: hold a write transaction for `hold`.
	wg.Add(1)
	go func() {
		defer wg.Done()
		tx, err := db.Begin()
		if err != nil {
			t.Errorf("writer begin: %v", err)
			return
		}
		defer tx.Rollback()
		_, _ = tx.Exec(`update data set data = data || 'w' where title = 'Alpha'`)
		time.Sleep(hold)
	}()

	// V5 read paths race the writer. Each read opens its own DB connection,
	// mirroring the HTTP handler path.
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				atomic.AddInt64(&panicCount, 1)
			}
		}()
		// exact-host + case normalization (V5 fix 1-2).
		r := route.Api_external_docs_for_host(anon, "1", "https://example.com/Page")
		if r["response"] != "ok" {
			t.Errorf("docs_for_host response=%v (want ok)", r)
		}
		r = route.Api_external_docs_for_host(anon, "1", "https://EXAMPL.eCom/page")
		if r["response"] != "ok" {
			t.Errorf("docs_for_host case response=%v (want ok)", r)
		}
		// all-hosts / batched IN paging (V5 fix 4; requires setting_backlink).
		r = route.Api_external_all_hosts(admin, "1")
		if r["response"] != "ok" {
			t.Errorf("all_hosts response=%v (want ok)", r)
		}
	}()

	wg.Wait()
	elapsed := time.Since(start)
	t.Logf("elapsed=%s panics=%d", elapsed, panicCount)

	if panicCount != 0 {
		t.Fatalf("V5 read panicked during write: %d", panicCount)
	}
	if elapsed > 12*time.Second {
		t.Fatalf("V5 read hung during write: elapsed=%s", elapsed)
	}
}
