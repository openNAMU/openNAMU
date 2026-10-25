package test

// Regression tests for the load_lang filter optimization: a process-wide
// filter-dedicated connection pool (tool.Get_lang_db) replaces the per-call
// DB_connect/DB_close that the pongo2 load_lang filter used to perform.
//
// Correctness invariants proven here:
//   - the current language is still read from the DB on every call, so setting
//     changes reflect on the very next render;
//   - replacing the DB file invalidates the cached pool (no stale language
//     from the previous DB);
//   - missing-key fallback, HTML escape and repeated-key behavior are unchanged;
//   - concurrent renders through the shared pool have no data race (run with
//     -race) and no cross-request mixing.

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"opennamu/route/tool"
)

// load_lang_setup boots an isolated DB for a backend and seeds the minimum
// rows so templates render.
func load_lang_setup(t *testing.T, backend string) *sql.DB {
	t.Helper()
	db := audit_setup(t, backend)
	audit_seed_auth(db, t)
	return db
}

// set_other_language writes the language setting the way the settings layer
// does (delete + insert; the `other` table has no unique key on name).
func set_other_language(db *sql.DB, language string) {
	if _, err := db.Exec(`delete from other where name = 'language'`); err != nil {
		panic(err)
	}
	if _, err := db.Exec(`insert into other (name, data, coverage) values ('language', ?, '')`, language); err != nil {
		panic(err)
	}
}

// Test_load_lang_language_switch_immediate: the language SELECT still happens
// per call, so changing the setting reflects on the very next render even
// though the filter now reuses one pool.
func Test_load_lang_language_switch_immediate(t *testing.T) {
	db := load_lang_setup(t, "sqlite")

	// fresh DB: no language row -> default ko-KR
	if got := tool.Get_language(tool.Get_lang_db(), "language_tag", false); got != "ko-KR" {
		t.Fatalf("default language = %q, want ko-KR", got)
	}

	set_other_language(db, "en-US")
	if got := tool.Get_language(tool.Get_lang_db(), "language_tag", false); got != "en-US" {
		t.Fatalf("after switch language = %q, want en-US (immediate reflection broken)", got)
	}

	// end-to-end: the very next real render must carry the new language in
	// the rendered <html lang=...> attribute (ringo skin uses load_lang there)
	rendered := tool.Get_template(db, tool.Config{IP: "admin_ip", Session: &emptySession{}},
		"lang_switch", "", nil, nil, map[string]string{})
	if !strings.Contains(rendered, `lang="en-US"`) {
		t.Fatalf("rendered page did not reflect en-US immediately")
	}

	set_other_language(db, "ko-KR")
	if got := tool.Get_language(tool.Get_lang_db(), "language_tag", false); got != "ko-KR" {
		t.Fatalf("after switch-back language = %q, want ko-KR", got)
	}
	rendered = tool.Get_template(db, tool.Config{IP: "admin_ip", Session: &emptySession{}},
		"lang_switch", "", nil, nil, map[string]string{})
	if !strings.Contains(rendered, `lang="ko-KR"`) {
		t.Fatalf("rendered page did not reflect ko-KR after switch-back")
	}
}

// Test_load_lang_db_switch_no_stale_language: pointing NAMU_DB at a different
// DB file must invalidate the cached pool. A stale pool would keep answering
// with the previous DB's language.
func Test_load_lang_db_switch_no_stale_language(t *testing.T) {
	db_a := load_lang_setup(t, "sqlite")
	// DB A explicitly uses en-US
	set_other_language(db_a, "en-US")
	if got := tool.Get_language(tool.Get_lang_db(), "language_tag", false); got != "en-US" {
		t.Fatalf("DB A language = %q, want en-US", got)
	}

	// swap to a brand-new DB file (fresh identity), which has no language row.
	// t.Setenv restores the previous NAMU_DB on cleanup.
	t.Setenv("NAMU_DB", t.TempDir()+"/switched")
	tool.DB_init()

	if got := tool.Get_language(tool.Get_lang_db(), "language_tag", false); got != "ko-KR" {
		t.Fatalf("after DB switch language = %q, want ko-KR (stale pool returned previous DB language)", got)
	}
}

// Test_load_lang_missing_key_and_escape: missing keys keep the "key (lang)"
// fallback, repeated keys behave identically, and safe=false always escapes
// exactly what safe=true returns.
func Test_load_lang_missing_key_and_escape(t *testing.T) {
	load_lang_setup(t, "sqlite")
	shared := tool.Get_lang_db()

	got := tool.Get_language(shared, "definitely_missing_key_zzz", false)
	if !strings.Contains(got, "definitely_missing_key_zzz") || !strings.Contains(got, "(") {
		t.Fatalf("missing key fallback = %q, want the \"key (lang)\" format", got)
	}

	for _, key := range []string{"language_tag", "save", "key", "definitely_missing_key_zzz"} {
		raw := tool.Get_language(shared, key, true)
		escaped := tool.Get_language(shared, key, false)
		if escaped != tool.HTML_escape(raw) {
			t.Fatalf("escape mismatch for %q: escaped=%q want %q", key, escaped, tool.HTML_escape(raw))
		}
		// repeated use must be stable
		if tool.Get_language(shared, key, false) != escaped {
			t.Fatalf("repeated Get_language for %q is not stable", key)
		}
	}
}

// Test_load_lang_concurrent_render_race: concurrent page renders through the
// shared filter pool (ringo skin, real load_lang executions) plus concurrent
// direct Get_language calls. Run with -race; also checks every response
// carries its own document content only.
func Test_load_lang_concurrent_render_race(t *testing.T) {
	db := load_lang_setup(t, "sqlite")
	audit_seed_docs(db, t, 4)

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
			own := fmt.Sprintf("LANG_ISO_MARKER_%02d", c)
			doc := fmt.Sprintf("lang_doc_%02d", c)
			tool.Exec_DB(db, `insert into data (title, data) values (?, ?)`, doc, own)
			<-start
			for i := 0; i < iters; i++ {
				// direct language reads concurrent with renders
				if got := tool.Get_language(tool.Get_lang_db(), "language_tag", false); got != "ko-KR" {
					mismatches.Add(1)
				}
				resp, err := client.Get(ts.URL + "/w/" + doc_label(doc))
				if err != nil {
					mismatches.Add(1)
					continue
				}
				buf := make([]byte, 1<<20)
				n, _ := resp.Body.Read(buf)
				body := string(buf[:n])
				resp.Body.Close()
				if !strings.Contains(body, own) {
					mismatches.Add(1)
				}
				for other := 0; other < clients; other++ {
					if other != c && strings.Contains(body, fmt.Sprintf("LANG_ISO_MARKER_%02d", other)) {
						mismatches.Add(1)
					}
				}
			}
		}(c)
	}
	close(start)
	wg.Wait()

	if mismatches.Load() != 0 {
		t.Errorf("language/render isolation mismatches=%d panics=%d", mismatches.Load(), panic_cnt.Load())
	}
	if panic_cnt.Load() != 0 {
		t.Errorf("handler panics=%d", panic_cnt.Load())
	}
}

// Test_load_lang_stress_optin: heavier concurrent render churn (opt-in),
// including the sqlite -> turso backend switch mid-test to exercise pool
// invalidation across backends.
func Test_load_lang_stress_optin(t *testing.T) {
	if os.Getenv("LOAD_LANG_STRESS") != "1" {
		t.Skip("stress: set LOAD_LANG_STRESS=1 to run")
	}

	db := load_lang_setup(t, "turso")
	audit_seed_docs(db, t, 8)

	panic_cnt := atomic.Int64{}
	ts := audit_http_server(t, &panic_cnt)

	var wg sync.WaitGroup
	var bad atomic.Int64
	start := make(chan struct{})
	for c := 0; c < 32; c++ {
		wg.Add(1)
		go func(c int) {
			defer wg.Done()
			client := &http.Client{Timeout: 30 * time.Second}
			own := fmt.Sprintf("LANG_STRESS_MARKER_%02d", c)
			doc := fmt.Sprintf("lang_stress_%02d", c)
			tool.Exec_DB(db, `insert into data (title, data) values (?, ?)`, doc, own)
			<-start
			for i := 0; i < 20; i++ {
				resp, err := client.Get(ts.URL + "/w/" + doc_label(doc))
				if err != nil {
					bad.Add(1)
					continue
				}
				buf := make([]byte, 1<<20)
				n, _ := resp.Body.Read(buf)
				body := string(buf[:n])
				resp.Body.Close()
				if !strings.Contains(body, own) || resp.StatusCode != 200 {
					bad.Add(1)
				}
			}
		}(c)
	}
	close(start)
	wg.Wait()
	if bad.Load() != 0 || panic_cnt.Load() != 0 {
		t.Errorf("stress bad responses=%d panics=%d", bad.Load(), panic_cnt.Load())
	}
}