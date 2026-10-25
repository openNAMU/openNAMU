package test

// Focused regression for the Get_template() concurrent-render data race.
//
// Root cause (pre-fix): pongo2's (*TemplateSet).FromString writes the shared
// DefaultSet.firstTemplateCreated bool without synchronization
// (template_sets.go:198) - both racing stacks in the baseline log are writes
// to that field. The fix wraps only the pongo2.FromString call in ui.go with
// a mutex; Execute stays parallel.
//
// Under -race the direct test below must fail before the fix and pass after.

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"opennamu/route/tool"
)

// Test_template_concurrent_render_direct hits the FromString race site
// directly: 32 goroutines (start-barrier synchronized) x 10 renders of the
// real Get_template.
func Test_template_concurrent_render_direct(t *testing.T) {
	db := audit_setup(t, "sqlite")
	audit_seed_auth(db, t)

	panics := atomic.Int64{}
	empty := atomic.Int64{}
	start := make(chan struct{})
	var wg sync.WaitGroup
	for w := 0; w < 32; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					panics.Add(1)
				}
			}()
			<-start
			config := tool.Config{IP: "admin_ip", Session: &emptySession{}}
			for i := 0; i < 10; i++ {
				out := tool.Get_template(db, config,
					"race_test",
					fmt.Sprintf("<div>worker %d iter %d</div>", w, i),
					nil, nil, map[string]string{},
				)
				if out == "" {
					empty.Add(1)
				}
			}
		}(w)
	}
	close(start)
	wg.Wait()

	if panics.Load() != 0 {
		t.Errorf("concurrent Get_template panics=%d", panics.Load())
	}
	if empty.Load() != 0 {
		t.Errorf("concurrent Get_template empty outputs=%d", empty.Load())
	}
}

// Test_template_concurrent_render_http renders real pages (View_w -> Get_template)
// through the actual route handlers with 16 concurrent clients.
func Test_template_concurrent_render_http(t *testing.T) {
	db := audit_setup(t, "sqlite")
	audit_seed_auth(db, t)
	audit_seed_docs(db, t, 5)

	panic_cnt := atomic.Int64{}
	ts := audit_http_server(t, &panic_cnt)

	start := make(chan struct{})
	var wg sync.WaitGroup
	var statuses atomic.Int64
	for c := 0; c < 16; c++ {
		wg.Add(1)
		go func(c int) {
			defer wg.Done()
			client := &http.Client{Timeout: 15 * time.Second}
			<-start
			for i := 0; i < 5; i++ {
				resp, err := client.Get(ts.URL + "/w/audit_" + fmt.Sprintf("%03d", (c+i)%5))
				if err != nil {
					t.Logf("http: %v", err)
					continue
				}
				io_copy_discard(resp)
				if resp.StatusCode == 200 {
					statuses.Add(1)
				}
			}
		}(c)
	}
	close(start)
	wg.Wait()

	if got := statuses.Load(); got != 16*5 {
		t.Errorf("successful renders=%d, want %d (panics=%d)", got, 16*5, panic_cnt.Load())
	}
	if panic_cnt.Load() != 0 {
		t.Errorf("handler panics=%d", panic_cnt.Load())
	}
}

// Test_template_concurrent_render_isolation proves per-request template
// contexts never mix: 16 concurrent clients each render their own document
// through the real route and every response must contain only that
// document's own marker.
func Test_template_concurrent_render_isolation(t *testing.T) {
	db := audit_setup(t, "sqlite")
	audit_seed_auth(db, t)

	panic_cnt := atomic.Int64{}
	ts := audit_http_server(t, &panic_cnt)

	const clients = 16
	const iters = 5
	var wg sync.WaitGroup
	var mismatches atomic.Int64
	start := make(chan struct{})
	for c := 0; c < clients; c++ {
		wg.Add(1)
		go func(c int) {
			defer wg.Done()
			client := &http.Client{Timeout: 15 * time.Second}
			own_marker := fmt.Sprintf("ISOLATE_MARKER_%02d", c)
			doc := fmt.Sprintf("isolate_doc_%02d", c)
			tool.Exec_DB(db, `insert into data (title, data) values (?, ?)`, doc, own_marker)
			<-start
			for i := 0; i < iters; i++ {
				resp, err := client.Get(ts.URL + "/w/" + doc_label(doc))
				if err != nil {
					mismatches.Add(1)
					continue
				}
				buf := make([]byte, 1<<20)
				n, _ := resp.Body.Read(buf)
				body := string(buf[:n])
				resp.Body.Close()
				if !strings.Contains(body, own_marker) {
					mismatches.Add(1)
				}
				for other := 0; other < clients; other++ {
					if other == c {
						continue
					}
					if strings.Contains(body, fmt.Sprintf("ISOLATE_MARKER_%02d", other)) {
						mismatches.Add(1)
					}
				}
			}
		}(c)
	}
	close(start)
	wg.Wait()

	if mismatches.Load() != 0 {
		t.Errorf("response isolation mismatches=%d (cross-request context mixing?)", mismatches.Load())
	}
	if panic_cnt.Load() != 0 {
		t.Errorf("handler panics=%d", panic_cnt.Load())
	}
}

// Test_template_skin_dynamic_reflect proves the render path still reads the
// skin file per request: edits to a test-only copy skin (under the test
// working directory, never a repository skin) are reflected immediately.
func Test_template_skin_dynamic_reflect(t *testing.T) {
	db := audit_setup(t, "sqlite")
	audit_seed_auth(db, t)

	skin_dir := filepath.Join("views", "race_test_skin")
	if err := os.MkdirAll(skin_dir, 0o755); err != nil {
		t.Fatalf("mkdir skin: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(skin_dir) })

	index_path := filepath.Join(skin_dir, "index.html")
	if err := os.WriteFile(index_path, []byte(`<html><body>RACE_TEST_SKIN_A {{ data }}</body></html>`), 0o644); err != nil {
		t.Fatalf("write index: %v", err)
	}
	if err := os.WriteFile(filepath.Join(skin_dir, "set.json"), []byte(`{}`), 0o644); err != nil {
		t.Fatalf("write set.json: %v", err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('skin','admin_ip','race_test_skin')`); err != nil {
		t.Fatalf("select skin: %v", err)
	}

	config := tool.Config{IP: "admin_ip", Session: &emptySession{}}
	out1 := tool.Get_template(db, config, "t", "CONTENT_A", nil, nil, map[string]string{})
	if !strings.Contains(out1, "RACE_TEST_SKIN_A") || !strings.Contains(out1, "CONTENT_A") {
		t.Fatalf("first render missing skin/data marker")
	}

	if err := os.WriteFile(index_path, []byte(`<html><body>RACE_TEST_SKIN_B {{ data }}</body></html>`), 0o644); err != nil {
		t.Fatalf("rewrite index: %v", err)
	}
	out2 := tool.Get_template(db, config, "t", "CONTENT_B", nil, nil, map[string]string{})
	if !strings.Contains(out2, "RACE_TEST_SKIN_B") || !strings.Contains(out2, "CONTENT_B") {
		t.Fatalf("skin file edit not reflected on re-render")
	}
	if strings.Contains(out2, "RACE_TEST_SKIN_A") {
		t.Fatalf("stale template served after skin edit")
	}
}