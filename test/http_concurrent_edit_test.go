package test

import (
	"database/sql"
	stdjson "encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"opennamu/route"
	"opennamu/route/tool"
)

// newEditServer returns a gin test server exposing the two production write
// routes we test here:
//   - POST /w/:doc_name/edit            -> route.Api_edit_post (real doc-edit + backlink path)
//   - POST /setting/backlink_reset      -> route.Api_setting_backlink_reset_post_internal (full rebuild)
//
// Handlers forward form fields to the real APIs with admin_ip, so the DB
// concurrency path is identical to production. Responses are the API's own
// response map as JSON so tests can assert on response / data.
func newEditServer(t *testing.T, adminIP string) *httptest.Server {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/w/:doc_name/edit", func(c *gin.Context) {
		config := tool.Config{IP: adminIP, Session: &emptySession{}}
		out := route.Api_edit_post(
			config,
			c.Param("doc_name"),
			c.PostForm("content"),
			c.PostForm("send"),
			c.PostForm("copyright_agreement"),
			c.PostForm("ver"),
		)
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, out)
	})
	r.POST("/setting/backlink_reset", func(c *gin.Context) {
		config := tool.Config{IP: adminIP, Session: &emptySession{}}
		load := c.PostForm("backlink_reset_load")
		if load == "" {
			load = "normal"
		}
		api := route.Api_setting_backlink_reset_post_internal(config, load)
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, api)
	})
	return httptest.NewServer(r)
}

// postEdit issues a real HTTP POST to the edit route; returns the parsed
// response map.
func postEdit(client *http.Client, base, docName string, content, send, agree, ver string) (map[string]any, error) {
	body := fmt.Sprintf("content=%s&send=%s&copyright_agreement=%s&ver=%s", content, send, agree, ver)
	resp, err := client.Post(base+"/w/"+docName+"/edit", "application/x-www-form-urlencoded", strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	out := make(map[string]any)
	if err := stdjson.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

// postRebuild issues a real HTTP POST to the full-rebuild route.
func postRebuild(client *http.Client, base string, load string) (map[string]any, error) {
	resp, err := http.Post(base+"/setting/backlink_reset", "application/x-www-form-urlencoded",
		strings.NewReader("backlink_reset_load="+load))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	out := make(map[string]any)
	_ = stdjson.NewDecoder(resp.Body).Decode(&out)
	return out, nil
}

// countBacklinks counts back rows for (link, title) with an optional type.
// An empty type matches any type (internal ” + external).
func countBacklinks(db *sql.DB, link, title, typ string) int {
	var n int
	q := "select count(*) from back where link = ? and title = ?"
	args := []any{link, title}
	if typ != "" {
		q += " and type = ?"
		args = append(args, typ)
	}
	if err := db.QueryRow(q, args...).Scan(&n); err != nil {
		return -1
	}
	return n
}

// backlinkTotal returns the number of backlinks of a given link (any type).
func backlinkTotal(db *sql.DB, link string) int {
	var n int
	if err := db.QueryRow("select count(*) from back where link = ?", link).Scan(&n); err != nil {
		return -1
	}
	return n
}

// Test_http_edit_different_docs_concurrent (Scenario A) races edits on two
// DIFFERENT docs, each introducing its own external link. Both must succeed,
// no panic, and each doc's backlink must reflect ONLY its new content (no
// cross-contamination).
func Test_http_edit_different_docs_concurrent(t *testing.T) {
	db := Setup_rebuild_db(t)
	defer tool.DB_close(db)
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed %v: %v", q, err)
		}
	}
	exec("insert into data (title, data) values ('DocA','alpha initial')")
	exec("insert into data (title, data) values ('DocB','beta initial')")

	server := newEditServer(t, "admin_ip")
	defer server.Close()
	client := &http.Client{}

	var wg sync.WaitGroup
	var panicCount int64
	var aOut, bOut map[string]any
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				atomic.AddInt64(&panicCount, 1)
			}
		}()
		out, err := postEdit(client, server.URL, "DocA", "[[https://alpha.example.com/a2|a2]]", "ok", "yes", "")
		if err != nil {
			panic(fmt.Sprintf("edit request failed: %v", err))
		}
		aOut = out
	}()
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				atomic.AddInt64(&panicCount, 1)
			}
		}()
		out, err := postEdit(client, server.URL, "DocB", "[[https://beta.example.com/b2|b2]]", "ok", "yes", "")
		if err != nil {
			panic(fmt.Sprintf("edit request failed: %v", err))
		}
		bOut = out
	}()
	wg.Wait()
	t.Logf("A: %q  B: %q", aOut["response"], bOut["response"])
	if panicCount != 0 {
		t.Fatalf("panic under concurrency: %d", panicCount)
	}
	if aOut["response"] != "ok" || bOut["response"] != "ok" {
		t.Fatalf("edits not ok: A=%q B=%q Afull=%v Bfull=%v", aOut["response"], bOut["response"], aOut, bOut)
	}
	// Each doc must have its own external backlink, and NOT the other's.
	if n := countBacklinks(db, "DocA", "alpha.example.com", "external"); n != 1 {
		t.Errorf("DocA external backlink alpha.example.com count=%d, want 1", n)
	}
	if n := countBacklinks(db, "DocA", "beta.example.com", "external"); n != 0 {
		t.Errorf("DocA cross-contamination beta.example.com count=%d, want 0", n)
	}
	if n := countBacklinks(db, "DocB", "beta.example.com", "external"); n != 1 {
		t.Errorf("DocB external backlink beta.example.com count=%d, want 1", n)
	}
	if n := countBacklinks(db, "DocB", "alpha.example.com", "external"); n != 0 {
		t.Errorf("DocB cross-contamination alpha.example.com count=%d, want 0", n)
	}
}

// Test_http_edit_same_doc_conflict (Scenario B, deterministic) exercises the
// optimistic-concurrency guard: an edit whose expected_revision is stale must be
// rejected with "edit conflict" and the doc content must remain the earlier
// committed value.
func Test_http_edit_same_doc_conflict(t *testing.T) {
	db := Setup_rebuild_db(t)
	defer tool.DB_close(db)
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed %v: %v", q, err)
		}
	}
	exec("insert into data (title, data) values ('TestDoc','TestDoc v0 [[https://zeta.example.com/v0|v0]]')")

	server := newEditServer(t, "admin_ip")
	defer server.Close()
	client := &http.Client{}

	// Step 1: first edit (expected_revision="") creates history + a real revision.
	out, err := postEdit(client, server.URL, "TestDoc", "TestDoc v1 [[https://zeta.example.com/v1|v1]]", "ok", "yes", "")
	if err != nil {
		t.Fatalf("first edit request failed: %v", err)
	}
	if out["response"] != "ok" {
		t.Fatalf("first edit not ok: %v", out)
	}
	var rev string
	if err := db.QueryRow("select id from history where title = 'TestDoc' order by id desc limit 1").Scan(&rev); err != nil || rev == "" {
		t.Fatalf("no revision after first edit: %v", err)
	}
	t.Logf("revision after v1 = %q", rev)

	// Step 2: an edit with that (current) revision commits v2.
	out, err = postEdit(client, server.URL, "TestDoc", "TestDoc v2 [[https://zeta.example.com/v2|v2]]", "ok", "yes", rev)
	if err != nil {
		t.Fatalf("v2 edit request failed: %v", err)
	}
	if out["response"] != "ok" {
		t.Fatalf("v2 edit (matching revision) not ok: %v", out)
	}
	var rev2 string
	if err := db.QueryRow("select id from history where title = 'TestDoc' order by id desc limit 1").Scan(&rev2); err != nil {
		t.Fatalf("read rev2: %v", err)
	}
	if rev2 == rev {
		t.Errorf("revision did not advance after v2: rev=%q rev2=%q", rev, rev2)
	}

	// Step 3: an edit with the STALE revision must be rejected.
	out, err = postEdit(client, server.URL, "TestDoc", "TestDoc v3 CORRUPT [[https://zeta.example.com/v3|v3]]", "ok", "yes", rev)
	if err != nil {
		t.Fatalf("stale edit request failed: %v", err)
	}
	if out["response"] != "error" {
		t.Fatalf("stale edit not rejected: %v", out)
	}
	if d, ok := out["data"].(string); ok && d != "edit conflict" {
		t.Errorf("stale edit data=%q, want 'edit conflict'", d)
	}
	// Content must still be v2 (v3 never committed).
	var data string
	if err := db.QueryRow("select data from data where title = 'TestDoc'").Scan(&data); err != nil || !strings.HasPrefix(data, "TestDoc v2") {
		t.Errorf("content after rejected stale edit: %q (want v2)", data)
	}
	// Backlink must reflect v2 (zeta v2), not the rejected v3.
	if n := countBacklinks(db, "TestDoc", "zeta.example.com", "external"); n != 1 {
		t.Errorf("TestDoc external backlink count=%d, want 1", n)
	}
}

// Test_http_edit_same_doc_concurrent (Scenario B2) races N edits on the SAME
// doc. Invariants: no panic, at least one succeeds, and the final content is
// exactly one of the submitted contents (no corruption / interleaving).
func Test_http_edit_same_doc_concurrent(t *testing.T) {
	db := Setup_rebuild_db(t)
	defer tool.DB_close(db)
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed %v: %v", q, err)
		}
	}
	exec("insert into data (title, data) values ('RaceDoc','RaceDoc v0')")

	server := newEditServer(t, "admin_ip")
	defer server.Close()
	client := &http.Client{}

	n := 12
	contents := make([]string, n)
	for i := 0; i < n; i++ {
		// Distinct external domain per version so a stale index (a backlink from
		// an earlier version) is visible as a backlink to the wrong host.
		contents[i] = fmt.Sprintf("RaceDoc v%d [[https://host%03d.example.com/doc|v%d]]", i, i, i)
	}
	var wg sync.WaitGroup
	var okCount, panicCount int64
	start := time.Now()
	for i := 0; i < n; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					atomic.AddInt64(&panicCount, 1)
				}
			}()
			out, err := postEdit(client, server.URL, "RaceDoc", contents[i], "ok", "yes", "")
			if err != nil {
				panic(fmt.Sprintf("edit request failed: %v", err))
			}
			if out["response"] == "ok" {
				atomic.AddInt64(&okCount, 1)
			}
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)
	t.Logf("same-doc concurrent: n=%d ok=%d panic=%d elapsed=%s", n, okCount, panicCount, elapsed)
	if panicCount != 0 {
		t.Fatalf("panic under same-doc concurrency: %d", panicCount)
	}
	if okCount == 0 {
		t.Fatalf("no edit succeeded: n=%d ok=0", n)
	}
	// Final content must be exactly one of the submitted contents, and the backlink
	// must agree with THAT content. Each version uses a distinct external domain, so
	// a stale index would show as a backlink to the wrong (non-winning) host.
	var data string
	if err := db.QueryRow("select data from data where title = 'RaceDoc'").Scan(&data); err != nil {
		t.Fatalf("read final content: %v", err)
	}
	win := -1
	for i, c := range contents {
		if data == c {
			win = i
		}
	}
	if win < 0 {
		t.Fatalf("final content %q is not one of the submitted contents (possible corruption)", data)
	}
	wantHost := fmt.Sprintf("host%03d.example.com", win)
	if n := countBacklinks(db, "RaceDoc", wantHost, "external"); n != 1 {
		t.Errorf("RaceDoc backlink must match the winning content (win=%d) but host %q count=%d", win, wantHost, n)
	}
	for i := range contents {
		if i == win {
			continue
		}
		other := fmt.Sprintf("host%03d.example.com", i)
		if n := countBacklinks(db, "RaceDoc", other, "external"); n != 0 {
			t.Errorf("stale RaceDoc backlink to non-winning host %q (win=%d): count=%d", other, win, n)
		}
	}
}

// Test_http_edit_during_rebuild (Scenario C) races a full rebuild against a
// doc edit, then verifies the backlink index matches the doc's FINAL content
// (not a stale snapshot) with no duplicates.
func Test_http_edit_during_rebuild(t *testing.T) {
	db := Setup_rebuild_db(t)
	defer tool.DB_close(db)
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed %v: %v", q, err)
		}
	}
	// Initial content links to alpha (an external host).
	exec("insert into data (title, data) values ('RaceTarget','[[https://alpha.example.com/1|one]]')")
	// Sentinel with a pre-existing external backlink, sorted LAST in the rebuild
	// (lowercase 'z' sorts after every other seed title). The rebuild's atomic
	// wipe deletes its backlink, and because the sentinel is processed last it
	// stays 0 for the entire rebuild. Polling for 0 deterministically means "the
	// rebuild is running" — no magic-number sleep.
	exec("insert into data (title, data) values ('zSentinel','[[https://sentry.example.com/doc|sentry]]')")
	exec("insert into back (link, title, type, data) values ('zSentinel','sentry.example.com','external','')")
	// Pad with many docs so the (slow) rebuild runs long enough to overlap the edit.
	for i := 0; i < 30; i++ {
		exec(fmt.Sprintf("insert into data (title, data) values ('Pad%03d','pad content %d')", i, i))
	}

	server := newEditServer(t, "admin_ip")
	defer server.Close()
	client := &http.Client{}

	var wg sync.WaitGroup
	var panicCount int64
	var rebuildOut map[string]any
	var rebuildErr error
	wg.Add(2)
	// Rebuild (slow) — the long-running write.
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				atomic.AddInt64(&panicCount, 1)
			}
		}()
		out, err := postRebuild(client, server.URL, "slow")
		rebuildOut = out
		rebuildErr = err
	}()
	// Deterministically wait until the rebuild has wiped the index (the sentinel's
	// pre-existing backlink is gone), then fire the edit while the rebuild is
	// guaranteed still running. This replaces the old 600ms sleep, which only
	// hoped for overlap.
	for i := 0; i < 300; i++ {
		if backlinkTotal(db, "zSentinel") == 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				atomic.AddInt64(&panicCount, 1)
			}
		}()
		out, err := postEdit(client, server.URL, "RaceTarget", "[[https://gamma.example.com/2|two]]", "ok", "yes", "")
		// A request error is a test failure (the HTTP call itself failed), not a
		// silent no-op.
		if err != nil {
			panic(fmt.Sprintf("edit request failed: %v", err))
		}
		_ = out
	}()
	wg.Wait()
	if panicCount != 0 {
		t.Fatalf("panic under rebuild-vs-edit: %d", panicCount)
	}
	if rebuildErr != nil {
		t.Fatalf("rebuild request failed: %v", rebuildErr)
	}
	if rebuildOut == nil {
		t.Fatalf("no rebuild response")
	}
	t.Logf("rebuild response=%q error_count=%v", rebuildOut["response"], rebuildOut["error_count"])

	// Final: RaceTarget's backlink must reflect the FINAL content (gamma), not
	// the stale initial (alpha).
	if n := countBacklinks(db, "RaceTarget", "gamma.example.com", "external"); n != 1 {
		t.Errorf("RaceTarget external backlink gamma.example.com count=%d, want 1 (stale index? alpha survived?)", n)
	}
	if n := countBacklinks(db, "RaceTarget", "alpha.example.com", "external"); n != 0 {
		t.Errorf("RaceTarget stale external backlink alpha.example.com count=%d, want 0", n)
	}
	// No duplicate backlinks for RaceTarget.
	if n := backlinkTotal(db, "RaceTarget"); n != 1 {
		t.Errorf("RaceTarget total backlinks=%d, want 1 (no duplicates)", n)
	}
}

// Test_http_rebuild_concurrent_reject (Stage 2-1) races two full rebuilds and
// verifies no panic, both reach a terminal response, and the final index is
// consistent (the seed doc's external backlink exists exactly once).
func Test_http_rebuild_concurrent_reject(t *testing.T) {
	db := Setup_rebuild_db(t)
	defer tool.DB_close(db)
	server := newEditServer(t, "admin_ip")
	defer server.Close()
	client := &http.Client{}

	var wg sync.WaitGroup
	var panicCount int64
	var out1, out2 map[string]any
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				atomic.AddInt64(&panicCount, 1)
			}
		}()
		out1, _ = postRebuild(client, server.URL, "fast")
	}()
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				atomic.AddInt64(&panicCount, 1)
			}
		}()
		out2, _ = postRebuild(client, server.URL, "fast")
	}()
	wg.Wait()
	t.Logf("rebuild1=%q rebuild2=%q", out1["response"], out2["response"])
	if panicCount != 0 {
		t.Fatalf("panic under concurrent rebuild: %d", panicCount)
	}
	for _, out := range []map[string]any{out1, out2} {
		if out == nil || out["response"] == "" {
			t.Fatal("a rebuild did not reach a terminal response")
		}
	}
	// Final index consistent: Example's external backlink to example.com exists
	// exactly once (no duplicate from a double-wipe).
	if n := countBacklinks(db, "Example", "example.com", "external"); n != 1 {
		t.Errorf("Example external backlink example.com count=%d, want 1 (no duplicate)", n)
	}
}

// Test_http_write_bounded_under_lock (Scenario E) issues a doc edit over HTTP
// while a long-held write transaction blocks it, asserting the request terminates
// within a bounded window (no infinite lock wait) and does not panic.
func Test_http_write_bounded_under_lock(t *testing.T) {
	db := Setup_rebuild_db(t)
	defer tool.DB_close(db)
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed %v: %v", q, err)
		}
	}
	exec("insert into data (title, data) values ('LockedDoc','[[https://locked.example.com/doc|locked]]')")

	server := newEditServer(t, "admin_ip")
	defer server.Close()
	client := &http.Client{}

	var wg sync.WaitGroup
	var panicCount int64
	var elapsed time.Duration
	wg.Add(2)
	// Hold a write transaction longer than the edit should take.
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				atomic.AddInt64(&panicCount, 1)
			}
		}()
		tx, err := db.Begin()
		if err != nil {
			return
		}
		defer tx.Rollback()
		_, _ = tx.Exec(`update data set data = data || 'lock' where title = 'LockedDoc'`)
		time.Sleep(3 * time.Second)
	}()
	// Issue the edit while the lock is held.
	var editErr error
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				atomic.AddInt64(&panicCount, 1)
			}
		}()
		start := time.Now()
		resp, err := http.Post(server.URL+"/w/LockedDoc/edit", "application/x-www-form-urlencoded",
			strings.NewReader("content=LockedDoc locked [[https://locked.example.com/doc|locked]]&send=ok&copyright_agreement=yes&ver="))
		if err != nil {
			editErr = err
			return
		}
		defer resp.Body.Close()
		elapsed = time.Since(start)
	}()
	wg.Wait()
	if editErr != nil {
		t.Fatalf("edit request under lock failed: %v", editErr)
	}
	t.Logf("edit under lock: panic=%d elapsed=%s", panicCount, elapsed)
	if panicCount != 0 {
		t.Fatalf("edit under lock panicked: %d", panicCount)
	}
	// Bounded termination: must not hang. DSN busy_timeout=5000 + Stage 4 bounded
	// retries guarantee a finite wait; a generous cap catches any infinite-regression.
	if elapsed > 15*time.Second {
		t.Errorf("edit under lock hung: elapsed=%s (want bounded)", elapsed)
	}
	_ = client
}
