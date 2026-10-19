package test

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"opennamu/route"
	"opennamu/route/tool"
)

// Test_setting_backlink_reset_clean_fail_preserves_stale_index verifies the fix
// that the rebuild's destructive backup/delete never runs on a failed run. The
// doc-list enumeration (which runs before the backup) must succeed first; a
// concurrent long-held write lock makes the atomic backup/delete fail, so the
// pre-existing (stale) backlink + link_count rows must be left intact.
func Test_setting_backlink_reset_clean_fail_preserves_stale_index(t *testing.T) {
	db := Setup_rebuild_db(t)
	defer tool.DB_close(db)
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed exec %v: %v", q, err)
		}
	}
	// Stale auto backlink + link_count that a successful rebuild would wipe.
	exec("insert into back (link, title, type, data) values ('Example', 'stale.example.com', 'external', '')")
	exec("insert into data_set (doc_name, doc_rev, set_name, set_data) values ('Example', '', 'link_count', '9')")

	var (
		wg        sync.WaitGroup
		panicCount int64
		writerMu  sync.Mutex
		writerErr error
	)
	var rebuildOut map[string]any
	start := time.Now()

	// Hold a write transaction long enough that the rebuild's backup/delete
	// (which needs the write lock) cannot complete; it must clean-fail.
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				atomic.AddInt64(&panicCount, 1)
			}
		}()
		tx, err := db.Begin()
		if err != nil {
			writerMu.Lock()
			writerErr = err
			writerMu.Unlock()
			return
		}
		defer func() { _ = tx.Rollback() }()
		_, _ = tx.Exec(`update data set data = data || 'lock' where title = 'Example'`)
		// Longer than the DSN busy_timeout (5s) so the rebuild's write
		// attempt cannot win the lock and must clean-fail.
		time.Sleep(6 * time.Second)
	}()

	// Run the rebuild while the write lock is held; it must clean-fail
	// within the bounded retry window without wiping the index.
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				atomic.AddInt64(&panicCount, 1)
			}
		}()
		out := route.Api_setting_backlink_reset_post_internal(tool.Config{IP: "admin_ip"}, "normal")
		rebuildOut = out
	}()

	wg.Wait()
	elapsed := time.Since(start)
	writerMu.Lock()
	if writerErr != nil {
		writerMu.Unlock()
		t.Fatalf("writer begin failed: %v", writerErr)
	}
	writerMu.Unlock()
	t.Logf("rebuild: response=%q elapsed=%s", rebuildOut["response"], elapsed)
	if panicCount != 0 {
		t.Fatalf("panic under clean-fail rebuild: %d", panicCount)
	}
	if rebuildOut == nil || rebuildOut["response"] == "" {
		t.Fatalf("rebuild did not reach a terminal response")
	}
	if elapsed > 15*time.Second {
		t.Errorf("rebuild hung: elapsed=%s (want bounded)", elapsed)
	}

	// The stale index rows must survive a clean-failed rebuild.
	var stale int
	if err := db.QueryRow("select count(*) from back where link = 'Example' and title = 'stale.example.com' and type = 'external'").Scan(&stale); err != nil {
		t.Fatalf("count stale: %v", err)
	}
	if stale != 1 {
		t.Errorf("stale backlink deleted by a clean-failed rebuild: count=%d, want 1", stale)
	}
	var lc int
	if err := db.QueryRow("select count(*) from data_set where doc_name = 'Example' and set_name = 'link_count'").Scan(&lc); err != nil {
		t.Fatalf("count link_count: %v", err)
	}
	if lc != 1 {
		t.Errorf("stale link_count deleted by a clean-failed rebuild: count=%d, want 1", lc)
	}
}

// Test_api_w_render_backlink_serialized_with_edit races a backlink render
// (Api_w_render, which under backlink_mutex rebuilds the index from the saved
// source) against a concurrent edit on the same doc. After both finish the index
// must reflect the final committed content, and the render's unsaved input must
// never be recorded — i.e. the render path cannot bypass the lock or write
// unsaved data.
func Test_api_w_render_backlink_serialized_with_edit(t *testing.T) {
	db := Setup_rebuild_db(t)
	defer tool.DB_close(db)
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed exec %v: %v", q, err)
		}
	}
	exec("insert into data (title, data) values ('RenderRace', '[[https://alpha.example.com/1|one]]')")

	server := newEditServer(t, "admin_ip")
	defer server.Close()
	client := &http.Client{}

	var wg sync.WaitGroup
	var panicCount int64

	// Concurrent edit (final committed content is gamma).
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				atomic.AddInt64(&panicCount, 1)
			}
		}()
		out, err := postEdit(client, server.URL, "RenderRace", "[[https://gamma.example.com/2|two]]", "ok", "yes", "")
		if err != nil {
			panic(fmt.Sprintf("edit request failed: %v", err))
		}
		if out["response"] != "ok" {
			panic(fmt.Sprintf("edit not ok: %v", out))
		}
	}()

	// Concurrent backlink render carrying an UNSAVED external link (beta).
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				atomic.AddInt64(&panicCount, 1)
			}
		}()
		out := route.Api_w_render(tool.Config{IP: "admin_ip"}, "RenderRace", "[[https://beta.example.com/unsaved|unsaved]]", "backlink", "")
		if out["response"] != "ok" {
			panic(fmt.Sprintf("render not ok: %v", out))
		}
	}()

	wg.Wait()
	if panicCount != 0 {
		t.Fatalf("panic under render-vs-edit: %d", panicCount)
	}

	// The final committed content must be gamma (the edit).
	var data string
	if err := db.QueryRow("select data from data where title = 'RenderRace'").Scan(&data); err != nil {
		t.Fatalf("read content: %v", err)
	}
	if !strings.Contains(data, "gamma.example.com") {
		t.Errorf("final content lacks gamma (edit lost): %q", data)
	}
	// The unsaved render input (beta) must never have been recorded in the index.
	if n := countBacklinks(db, "RenderRace", "beta.example.com", "external"); n != 0 {
		t.Errorf("unsaved render backlink recorded: beta.example.com count=%d, want 0", n)
	}
}
