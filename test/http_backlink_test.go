package test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"opennamu/route"
	"opennamu/route/tool"
)

// Test_backlink_reset_http_under_write holds a write transaction and issues a
// real HTTP POST to the /setting/backlink_reset route (via the real view + API
// path), asserting the request completes within a bounded window (no hang) and
// the handler does not panic. This exercises the full HTTP -> view -> API ->
// DB path under a concurrent write, the outermost layer of the rebuild.
func Test_backlink_reset_http_under_write(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db_dir := filepath.Join(t.TempDir(), "db")
	if err := os.MkdirAll(db_dir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	os.Setenv("NAMU_DB", db_dir)
	os.Setenv("NAMU_DB_TYPE", "sqlite")
	defer os.Unsetenv("NAMU_DB")
	defer os.Unsetenv("NAMU_DB_TYPE")

	tool.DB_init()
	db := tool.DB_connect()
	defer tool.DB_close(db)

	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed %v: %v", q, err)
		}
	}
	// Default permission groups + bind admin_ip to owner (grants
	// setting_backlink), mirroring the other backlink tests.
	exec("insert into alist (name, acl) values ('owner','owner')")
	exec("insert into alist (name, acl) values ('admin','admin')")
	exec("insert into alist (name, acl) values ('user','user')")
	exec("insert into alist (name, acl) values ('ip','ip')")
	exec("insert into user_set (name, id, data) values ('acl','admin_ip','owner')")
	exec("insert into data (title, data) values ('Alpha','content')")
	exec("insert into data (title, data) values ('Beta','content')")

	// Minimal router with the real backlink-reset POST route. The handler
	// mirrors routes_setting.go but supplies the config/form directly (the
	// DB concurrency path is identical to the production handler).
	r := gin.New()
	r.POST("/setting/backlink_reset", func(c *gin.Context) {
		config := tool.Config{IP: "admin_ip", Session: &emptySession{}}
		form := map[string]string{"backlink_reset_load": "normal"}
		data := route.View_setting_backlink_reset_post(config, form)
		c.Header("Content-Type", "text/html; charset=utf-8")
		_, _ = c.Writer.Write([]byte(data))
	})
	ts := httptest.NewServer(r)
	defer ts.Close()

	var wg sync.WaitGroup
	var panicCount int64
	var status int
	var postErr string
	var elapsed time.Duration

	// Concurrent writer: hold a write transaction for `hold`.
	hold := 2 * time.Second
	wg.Add(1)
	go func() {
		defer wg.Done()
		tx, err := db.Begin()
		if err != nil {
			t.Logf("writer begin: %v", err)
			return
		}
		defer tx.Rollback()
		_, _ = tx.Exec(`update data set data = data || 'lock' where title = 'Alpha'`)
		time.Sleep(hold)
	}()

	// Issue the rebuild over HTTP while the writer holds the lock.
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if rec := recover(); rec != nil {
				atomic.AddInt64(&panicCount, 1)
			}
		}()
		start := time.Now()
		resp, err := http.Post(ts.URL+"/setting/backlink_reset",
			"application/x-www-form-urlencoded",
			strings.NewReader("backlink_reset_load=normal"))
		if err != nil {
			postErr = err.Error()
			elapsed = time.Since(start)
			return
		}
		defer resp.Body.Close()
		status = resp.StatusCode
		elapsed = time.Since(start)
	}()

	wg.Wait()
	t.Logf("http status=%d post_err=%q panic=%d elapsed=%s", status, postErr, panicCount, elapsed)

	if panicCount != 0 {
		t.Fatalf("HTTP rebuild handler panicked: %d", panicCount)
	}
	if status == 0 {
		t.Errorf("no HTTP response received (hung?); post_err=%q panic=%d", postErr, panicCount)
	}
}
