package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"opennamu/route"
	"opennamu/route/tool"
)

func Test_http_w_document_comment_flow(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		t.Run(backend, func(t *testing.T) {
			db := audit_setup(t, backend)
			audit_seed_auth(db, t)

			for _, statement := range []struct {
				query string
				args  []any
			}{
				{query: "insert into data (title, data) values (?, ?)", args: []any{"CommentFlowDoc", "comment flow body"}},
				{query: "insert into other (name, data, coverage) values (?, ?, '')", args: []any{"enable_comment", "1"}},
				{query: "insert into bbs_set (set_name, set_code, set_id, set_data) values (?, '', '0', ?)", args: []any{"bbs_name", "document_comment"}},
				{query: "insert into bbs_set (set_name, set_code, set_id, set_data) values (?, '', '0', ?)", args: []any{"bbs_type", "comment"}},
			} {
				if _, err := db.Exec(statement.query, statement.args...); err != nil {
					t.Fatal(err)
				}
			}

			config := tool.Config{IP: "admin_ip", Session: &emptySession{}}
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.GET("/w/:doc_name", func(c *gin.Context) {
				page, status := route.View_w(c, config, c.Param("doc_name"), "")
				c.Data(status, "text/html; charset=utf-8", []byte(page))
			})
			server := httptest.NewServer(router)
			defer server.Close()
			http_client := &http.Client{Timeout: 10 * time.Second}

			request_page := func() (int, string) {
				response, err := http_client.Get(server.URL + "/w/CommentFlowDoc")
				if err != nil {
					t.Fatal(err)
				}
				defer response.Body.Close()
				body, err := io.ReadAll(response.Body)
				if err != nil {
					t.Fatal(err)
				}
				return response.StatusCode, string(body)
			}

			for request_index := 0; request_index < 2; request_index++ {
				status, body := request_page()
				if status != http.StatusOK {
					t.Fatalf("request %d status=%d, want %d", request_index+1, status, http.StatusOK)
				}
				if !strings.Contains(body, `id="opennamu_bbs_w_post_select"`) {
					t.Fatalf("request %d did not render the document comment form", request_index+1)
				}
			}

			var comment_code string
			if err := db.QueryRow("select set_data from data_set where doc_name = ? and set_name = 'document_comment_code'", "CommentFlowDoc").Scan(&comment_code); err != nil {
				t.Fatalf("read document comment mapping: %v", err)
			}
			var comment_post_count int
			if err := db.QueryRow("select count(*) from bbs_data where set_id = '0' and set_code = ? and set_name = 'title' and set_data = ?", comment_code, "CommentFlowDoc").Scan(&comment_post_count); err != nil {
				t.Fatalf("count document comment posts: %v", err)
			}
			if comment_code == "" || comment_post_count != 1 {
				t.Fatalf("comment code=%q, matching post count=%d; want one stable comment post", comment_code, comment_post_count)
			}

			legacy_ui := route.Api_w_comment_ui(config, "CommentFlowDoc")
			legacy_comment_html, ok := legacy_ui["data"].(string)
			if !ok {
				t.Fatalf("legacy comment UI data has type %T, want string", legacy_ui["data"])
			}
			comment_api := route.Api_w_comment(config, "CommentFlowDoc")
			comment_code, _ = comment_api["data"].(string)
			comment_html := route.View_bbs_in_w_comment(db, config, "0", comment_code, "", 1)
			if comment_html != legacy_comment_html {
				t.Fatal("direct API/View composition differs from the retained legacy UI adapter")
			}
		})
	}
}
