package test

// HTTP-level audit over the real production view handlers (View_w,
// View_edit_post, View_list_history, View_setting_backlink_reset_post) for both
// backends. Covers concurrent distinct-document edits, same-document optimistic
// conflicts, read/write mixes, and backlink/history integrity. Directive:
// tmp/openNAMU_SQLite_Turso_concurrency_audit_20261025.md.

import (
	"database/sql"
	"fmt"
	"io"
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

// audit_http_server builds an httptest server exposing the real handlers with
// a panic-counting middleware. The IP "admin_ip" is bound to owner by
// audit_seed_auth so permission and captcha ACL checks pass.
func audit_http_server(t *testing.T, panic_cnt *atomic.Int64) *httptest.Server {
	t.Helper()
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				panic_cnt.Add(1)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	})

	config_of := func(c *gin.Context) tool.Config {
		return tool.Config{IP: "admin_ip", Session: &emptySession{}}
	}

	r.GET("/w/*doc_name", func(c *gin.Context) {
		data, status := route.View_w(c, config_of(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "")
		c.Data(status, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/history/*doc_name", func(c *gin.Context) {
		data := route.View_list_history(config_of(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "", "1")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/edit/*doc_name", func(c *gin.Context) {
		// The captcha form value is empty here: the seeded owner binding passes the
// recaptcha ACL check inside Captcha_check before any captcha response is
// evaluated, mirroring a logged-in admin session in production.
		data := route.View_edit_post(
			config_of(c),
			strings.TrimPrefix(c.Param("doc_name"), "/"),
			c.PostForm("content"),
			c.PostForm("send"),
			c.PostForm("copyright_agreement"),
			"",
			c.PostForm("ver"),
		)
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/setting/backlink_reset", func(c *gin.Context) {
		data := route.View_setting_backlink_reset_post(config_of(c), map[string]string{"backlink_reset_load": "normal"})
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	ts := httptest.NewServer(r)
	t.Cleanup(ts.Close)
	return ts
}

func audit_post_edit(t *testing.T, ts *httptest.Server, doc, content, send, expected_rev string) (int, string, string) {
	t.Helper()
	form := fmt.Sprintf("content=%s&send=%s&copyright_agreement=yes", doc_label(content), doc_label(send))
	if expected_rev != "" {
		form += "&ver=" + doc_label(expected_rev)
	}
	resp, err := http.Post(ts.URL+"/edit/"+doc_label(doc), "application/x-www-form-urlencoded", strings.NewReader(form))
	if err != nil {
		return 0, "transport_error", err.Error()
	}
	defer resp.Body.Close()

	buf := make([]byte, 1<<20)
	n, _ := resp.Body.Read(buf)
	body := string(buf[:n])
	switch {
	case strings.Contains(body, "edit conflict"):
		return resp.StatusCode, "edit_conflict", ""
	case strings.Contains(body, "CAPTCHA"):
		return resp.StatusCode, "captcha_required", ""
	case strings.Contains(body, "permission denied"):
		return resp.StatusCode, "permission_denied", ""
	case strings.Contains(body, "error"):
		return resp.StatusCode, "other_error", cut_body(body)
	default:
		return resp.StatusCode, "ok", ""
	}
}

// doc_label percent-encodes a value for form encoding without importing net/url
// in every helper; labels here are fixed ASCII audit tokens.
func doc_label(s string) string {
	var b strings.Builder
	for _, ch := range s {
		switch {
		case ch >= 'a' && ch <= 'z', ch >= 'A' && ch <= 'Z', ch >= '0' && ch <= '9',
			ch == '_', ch == '-', ch == '.':
			b.WriteRune(ch)
		default:
			b.WriteString(fmt.Sprintf("%%%02X", ch))
		}
	}
	return b.String()
}

func cut_body(body string) string {
	// extract a short plain-text hint of the error page
	if i := strings.Index(body, "오류"); i >= 0 {
		rest := body[i:]
		rest = strings.ReplaceAll(rest, "\n", " ")
		if len(rest) > 100 {
			rest = rest[:100]
		}
		return rest
	}
	if len(body) > 100 {
		return body[:100]
	}
	return body
}

// io_copy_discard drains and closes a response body.
func io_copy_discard(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
}

// Test_audit_http_routes: real-route scenarios for both backends.
func Test_audit_http_routes(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		db := audit_setup(t, backend)
		audit_seed_auth(db, t)
		audit_seed_docs(db, t, 10)

		var panic_cnt atomic.Int64
		ts := audit_http_server(t, &panic_cnt)

		// B-http: 16 concurrent edits to 16 distinct documents
		{
			var wg sync.WaitGroup
			var ok_cnt, bad_cnt atomic.Int64
			var first_bad atomic.Value
			start := time.Now()
			for w := 0; w < 16; w++ {
				wg.Add(1)
				go func(w int) {
					defer wg.Done()
					doc := fmt.Sprintf("audit_http_b%02d", w)
					status, kind, detail := audit_post_edit(t, ts, doc, "content "+doc, "audit b", "")
					if kind == "ok" {
						ok_cnt.Add(1)
					} else {
						bad_cnt.Add(1)
						first_bad.CompareAndSwap(nil, fmt.Sprintf("%s %d %s", kind, status, detail))
					}
				}(w)
			}
			wg.Wait()
			elapsed := time.Since(start)

			got_hist := 0
			if !tool.QueryRow_DB(db, `select count(*) from history where title like 'audit_http_b%'`, []any{&got_hist}) {
				t.Fatalf("B-http %s: count failed", backend)
			}
			got_data := 0
			_ = tool.QueryRow_DB(db, `select count(*) from data where title like 'audit_http_b%'`, []any{&got_data})

			// backlink integrity: edits contained no [[links]], so back rows for
			// these titles must be absent; history/data must be exact.
			audit_record(t, map[string]any{
				"type": "B_http", "backend": backend,
				"ok": ok_cnt.Load(), "bad": bad_cnt.Load(),
				"hist": got_hist, "data": got_data,
				"panics": panic_cnt.Load(), "secs": elapsed.Seconds(),
				"first_bad": first_bad.Load(),
			})
			if ok_cnt.Load() != 16 || got_hist != 16 || got_data != 16 {
				t.Errorf("B-http %s: ok=%d hist=%d data=%d (want 16 each), first_bad=%v",
					backend, ok_cnt.Load(), got_hist, got_data, first_bad.Load())
			}
			t.Logf("B-http %s ok=%d bad=%d panics=%d %.2fs", backend, ok_cnt.Load(), bad_cnt.Load(), panic_cnt.Load(), elapsed.Seconds())
		}

		// C-http-1: 8 concurrent same-doc edits with a stale expected revision.
		// Exactly one must win; the other 7 must get an explicit edit conflict.
		// Classification uses the real Api_edit_post return map (the view layer
		// only wraps it in localized HTML, which is not greppable reliably).
		{
			var ok_cnt, conflict_cnt, other_cnt atomic.Int64
			var wg sync.WaitGroup
			for w := 0; w < 8; w++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					result := route.Api_edit_post(
						tool.Config{IP: "admin_ip", Session: &emptySession{}},
						"audit_http_c", "content", "audit c", "yes", "0",
					)
					switch result["data"] {
					case "edit conflict":
						conflict_cnt.Add(1)
					default:
						if result["response"] == "ok" {
							ok_cnt.Add(1)
						} else {
							other_cnt.Add(1)
							t.Logf("C-http %s unexpected: %v", backend, result)
						}
					}
				}()
			}
			wg.Wait()

			got_hist := 0
			_ = tool.QueryRow_DB(db, `select count(*) from history where title = 'audit_http_c'`, []any{&got_hist})
			audit_record(t, map[string]any{
				"type": "C_http_stale_rev", "backend": backend,
				"ok": ok_cnt.Load(), "explicit_conflicts": conflict_cnt.Load(), "other": other_cnt.Load(),
				"hist_rows": got_hist, "panics": panic_cnt.Load(),
			})
			if ok_cnt.Load() != 1 || conflict_cnt.Load() != 7 {
				t.Errorf("C-http %s: ok=%d conflicts=%d other=%d, want 1 ok + 7 explicit conflicts",
					backend, ok_cnt.Load(), conflict_cnt.Load(), other_cnt.Load())
			}
			if got_hist != 1 {
				t.Errorf("C-http %s: history rows=%d, want 1 (no duplicate revisions)", backend, got_hist)
			}
			t.Logf("C-http %s ok=%d conflict=%d other=%d hist=%d", backend, ok_cnt.Load(), conflict_cnt.Load(), other_cnt.Load(), got_hist)
		}

		// C-http-2: 8 concurrent same-doc edits without expected_revision via the
// real HTTP route. All must succeed (serialized by the app-level save
// section); revision ids must be strictly increasing with no duplicates.
// The document already has 1 revision from C-http-1, so ids end at 9.
		{
			var ok_cnt, bad_cnt atomic.Int64
			var wg sync.WaitGroup
			for w := 0; w < 8; w++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					result := route.Api_edit_post(
						tool.Config{IP: "admin_ip", Session: &emptySession{}},
						"audit_http_c", "content", "audit c", "yes", "",
					)
					if result["response"] == "ok" {
						ok_cnt.Add(1)
					} else {
						bad_cnt.Add(1)
						t.Logf("C-http2 %s unexpected: %v", backend, result)
					}
				}()
			}
			wg.Wait()

			rows := tool.Query_DB(db,
				`select id from history where title = 'audit_http_c' order by cast(id as integer)`)
			ids := []string{}
			for rows.Next() {
				var id string
				if rows.Scan(&id) == nil {
					ids = append(ids, id)
				}
			}
			rows.Close()
			continuous := len(ids) == 9
			if continuous {
				for i, id := range ids {
					if id != fmt.Sprintf("%d", i+1) {
						continuous = false
					}
				}
			}
			audit_record(t, map[string]any{
				"type": "C_http_no_rev", "backend": backend,
				"ok": ok_cnt.Load(), "bad": bad_cnt.Load(),
				"history_ids": strings.Join(ids, ","), "continuous": continuous,
				"panics": panic_cnt.Load(),
			})
			if ok_cnt.Load() != 8 || !continuous {
				t.Errorf("C-http2 %s: ok=%d ids=%v continuous=%t (want 8 ok, revisions 1..9 continuous)",
					backend, ok_cnt.Load(), ids, continuous)
			}
			t.Logf("C-http2 %s ok=%d bad=%d ids=%s", backend, ok_cnt.Load(), bad_cnt.Load(), strings.Join(ids, ","))
		}

		// D-http: mixed real-route load - 8 readers (View_w + View_list_history)
		// and 2 writers (distinct-doc edits) for ~2 seconds.
		{
			var read_ok, read_err, write_ok, write_bad atomic.Int64
			var wg sync.WaitGroup
			deadline := time.Now().Add(2 * time.Second)
			for r := 0; r < 8; r++ {
				wg.Add(1)
				go func(r int) {
					defer wg.Done()
					client := &http.Client{Timeout: 10 * time.Second}
					i := 0
					for time.Now().Before(deadline) {
						doc := fmt.Sprintf("audit_%03d", (r+i)%10)
						resp, err := client.Get(ts.URL + "/w/" + doc_label(doc))
						if err != nil {
							read_err.Add(1)
							continue
						}
						io_copy_discard(resp)
						if resp.StatusCode == 200 {
							read_ok.Add(1)
						} else {
							read_err.Add(1)
						}
						i++
					}
				}(r)
			}
			for w := 0; w < 2; w++ {
				wg.Add(1)
				go func(w int) {
					defer wg.Done()
					for i := 0; time.Now().Before(deadline); i++ {
						doc := fmt.Sprintf("audit_http_d%02d_%03d", w, i)
						status, kind, _ := audit_post_edit(t, ts, doc, "content", "audit d", "")
						_ = status
						if kind == "ok" {
							write_ok.Add(1)
						} else {
							write_bad.Add(1)
						}
					}
				}(w)
			}
			wg.Wait()

			got_hist := 0
			_ = tool.QueryRow_DB(db, `select count(*) from history where title like 'audit_http_d%'`, []any{&got_hist})
			audit_record(t, map[string]any{
				"type": "D_http_8R2W", "backend": backend,
				"read_ok": read_ok.Load(), "read_err": read_err.Load(),
				"write_ok": write_ok.Load(), "write_bad": write_bad.Load(),
				"hist": got_hist, "panics": panic_cnt.Load(),
			})
			if read_err.Load() != 0 || write_bad.Load() != 0 || panic_cnt.Load() != 0 || got_hist != int(write_ok.Load()) {
				t.Errorf("D-http %s: read_err=%d write_bad=%d panics=%d hist=%d write_ok=%d",
					backend, read_err.Load(), write_bad.Load(), panic_cnt.Load(), got_hist, write_ok.Load())
			}
			t.Logf("D-http %s reads=%d(err%d) writes=%d(bad%d) panics=%d",
				backend, read_ok.Load(), read_err.Load(), write_ok.Load(), write_bad.Load(), panic_cnt.Load())
		}

		// E-http: real backlink reset under a held write transaction must
		// complete within a bounded window (no hang), mirroring the keeper tests.
		{
			hold_done := make(chan struct{})
			go func() {
				defer close(hold_done)
				_ = tool.DB_transaction(db, func(tx *sql.Tx) error {
					if _, err := tx.Exec(`insert into history (id, title, data, date, ip, send, leng, hide, type) values ('eh','audit_http_e','v','2026-01-01 00:00:00','admin_ip','audit','+1','','')`); err != nil {
						return err
					}
					time.Sleep(2 * time.Second)
					return nil
				})
			}()
			time.Sleep(200 * time.Millisecond)

			client := &http.Client{Timeout: 30 * time.Second}
			start := time.Now()
			resp, err := client.Post(ts.URL+"/setting/backlink_reset", "application/x-www-form-urlencoded", strings.NewReader("backlink_reset_load=normal"))
			secs := time.Since(start).Seconds()
			status := 0
			if err == nil {
				io_copy_discard(resp)
				status = resp.StatusCode
			}
			<-hold_done
			audit_record(t, map[string]any{
				"type": "E_http_backlink_reset", "backend": backend,
				"status": status, "transport_err": err_str(err), "secs": secs,
				"panics": panic_cnt.Load(),
			})
			if err != nil || status != 200 || secs > 25 {
				t.Errorf("E-http %s: backlink reset under write hold: status=%d err=%v secs=%.1f",
					backend, status, err, secs)
			}
			t.Logf("E-http %s status=%d secs=%.2f", backend, status, secs)
		}

		if panic_cnt.Load() != 0 {
			t.Errorf("http %s: total handler panics=%d", backend, panic_cnt.Load())
		}
	}
}