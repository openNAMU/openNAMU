package test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"opennamu/route"
	"opennamu/route/tool"
)

func Test_api_w_page_view_standalone_contract(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		t.Run(backend, func(t *testing.T) {
			db := audit_setup(t, backend)
			audit_seed_auth(db, t)

			config := tool.Config{IP: "admin_ip", Session: &emptySession{}}
			const doc_name = "PageViewApiTarget"

			post_result := route.Api_w_page_view_post(config, doc_name)
			if post_result["response"] != "ok" || len(post_result) != 1 {
				t.Fatalf("page-view post response = %#v, want only response=ok", post_result)
			}
			if result := route.Api_w_page_view(config, doc_name); result["response"] != "ok" || result["data"] != 1 {
				t.Fatalf("standalone page-view read = %#v, want ok with count 1", result)
			}

			post_result = route.Api_w_page_view_post(config, doc_name)
			if post_result["response"] != "ok" {
				t.Fatalf("second page-view post = %#v, want ok", post_result)
			}
			if result := route.Api_w_page_view(config, doc_name); result["response"] != "ok" || result["data"] != 2 {
				t.Fatalf("page-view read after second post = %#v, want count 2", result)
			}

			counter_for := func(doc_rev string) string {
				var counter string
				tool.QueryRow_DB(db,
					"select set_data from data_set where doc_name = ? and set_name = 'view_count' and doc_rev = ?",
					[]any{&counter}, doc_name, doc_rev,
				)
				return counter
			}
			counter_revisions := []string{"", tool.Get_month(), tool.Get_date()}
			for _, doc_rev := range counter_revisions {
				if got := counter_for(doc_rev); got != "2" {
					t.Errorf("counter for revision %q = %q, want 2", doc_rev, got)
				}
			}

			if _, err := db.Exec(`delete from other where name = 'not_use_view_count'`); err != nil {
				t.Fatalf("clear view-count setting: %v", err)
			}
			if _, err := db.Exec(`insert into other (name, data, coverage) values ('not_use_view_count', '1', '')`); err != nil {
				t.Fatalf("disable view counts: %v", err)
			}
			if post_result := route.Api_w_page_view_post(config, doc_name); post_result["response"] != "ok" {
				t.Fatalf("disabled page-view post = %#v, want ok", post_result)
			}
			if result := route.Api_w_page_view(config, doc_name); result["response"] != "ok" || result["data"] != 0 {
				t.Fatalf("disabled page-view read = %#v, want ok with count 0", result)
			}
			for _, doc_rev := range counter_revisions {
				if got := counter_for(doc_rev); got != "2" {
					t.Errorf("disabled counter for revision %q = %q, want unchanged value 2", doc_rev, got)
				}
			}
		})
	}
}

func Test_api_w_page_view_post_concurrent_updates(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		t.Run(backend, func(t *testing.T) {
			db := audit_setup(t, backend)
			audit_seed_auth(db, t)

			const doc_name = "ConcurrentPageViewTarget"
			const writer_count = 8
			config := tool.Config{IP: "admin_ip", Session: &emptySession{}}
			write_errors := make(chan string, writer_count)
			var wait_group sync.WaitGroup
			for range writer_count {
				wait_group.Add(1)
				go func() {
					defer wait_group.Done()
					defer func() {
						if recovered := recover(); recovered != nil {
							write_errors <- "page-view post panicked"
						}
					}()
					result := route.Api_w_page_view_post(config, doc_name)
					if result["response"] != "ok" {
						write_errors <- "page-view post did not return ok"
					}
				}()
			}
			wait_group.Wait()
			close(write_errors)
			for write_error := range write_errors {
				t.Error(write_error)
			}

			result := route.Api_w_page_view(config, doc_name)
			if result["response"] != "ok" || result["data"] != writer_count {
				t.Fatalf("page-view count after concurrent posts = %#v, want %d", result, writer_count)
			}
			for _, doc_rev := range []string{tool.Get_month(), tool.Get_date()} {
				var counter string
				tool.QueryRow_DB(db,
					"select set_data from data_set where doc_name = ? and set_name = 'view_count' and doc_rev = ?",
					[]any{&counter}, doc_name, doc_rev,
				)
				if counter != strconv.Itoa(writer_count) {
					t.Errorf("concurrent counter for revision %q = %q, want %d", doc_rev, counter, writer_count)
				}
			}
		})
	}
}

func Test_View_w_concurrent_page_view_counts(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		t.Run(backend, func(t *testing.T) {
			db := audit_setup(t, backend)
			audit_seed_auth(db, t)
			if _, err := db.Exec(`delete from other where name = 'related_document_enabled'`); err != nil {
				t.Fatalf("clear related-document setting: %v", err)
			}
			if _, err := db.Exec(`insert into other (name, data, coverage) values ('related_document_enabled', '0', '')`); err != nil {
				t.Fatalf("disable related-document lookup: %v", err)
			}
			if _, err := db.Exec(`insert into data (title, data, type) values ('ConcurrentViewDoc', 'view body', '')`); err != nil {
				t.Fatalf("seed document: %v", err)
			}

			gin.SetMode(gin.TestMode)
			const view_count = 4
			view_errors := make(chan string, view_count)
			var wait_group sync.WaitGroup
			for range view_count {
				wait_group.Add(1)
				go func() {
					defer wait_group.Done()
					defer func() {
						if recovered := recover(); recovered != nil {
							view_errors <- "View_w panicked during a concurrent page view"
						}
					}()
					context, _ := gin.CreateTestContext(httptest.NewRecorder())
					context.Request = httptest.NewRequest(http.MethodGet, "/w/ConcurrentViewDoc", nil)
					config := tool.Config{IP: "admin_ip", Session: &emptySession{}}
					_, status := route.View_w(context, config, "ConcurrentViewDoc", "")
					if status != http.StatusOK {
						view_errors <- "concurrent View_w did not return HTTP 200"
					}
				}()
			}
			wait_group.Wait()
			close(view_errors)
			for view_error := range view_errors {
				t.Error(view_error)
			}

			count := tool.Get_document_setting_value(db, "ConcurrentViewDoc", "view_count", "")
			if count != strconv.Itoa(view_count) {
				t.Fatalf("page-view count after concurrent View_w calls = %q, want %d", count, view_count)
			}
		})
	}
}

func Test_View_w_records_page_views_for_normal_and_redirect_requests(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		t.Run(backend, func(t *testing.T) {
			db := audit_setup(t, backend)
			audit_seed_auth(db, t)

			exec := func(query string, args ...any) {
				t.Helper()
				if _, err := db.Exec(query, args...); err != nil {
					t.Fatalf("seed %q: %v", query, err)
				}
			}
			exec(`delete from other where name = 'related_document_enabled'`)
			exec(`insert into other (name, data, coverage) values ('related_document_enabled', '0', '')`)
			exec(`insert into data (title, data, type) values ('ViewCountDoc', 'page view test', '')`)
			exec(`insert into data (title, data, type) values ('RedirectViewCountDoc', 'redirect source', '')`)
			exec(`insert into data (title, data, type) values ('RedirectViewTarget', 'redirect target', '')`)
			exec(`insert into back (title, link, data, type) values ('RedirectViewTarget', 'RedirectViewCountDoc', '', 'redirect')`)

			gin.SetMode(gin.TestMode)
			config := tool.Config{IP: "admin_ip", Session: &emptySession{}}
			view := func(doc_name string) (string, int) {
				context, _ := gin.CreateTestContext(httptest.NewRecorder())
				context.Request = httptest.NewRequest(http.MethodGet, "/w/"+doc_name, nil)
				return route.View_w(context, config, doc_name, "")
			}

			for expected_count := 1; expected_count <= 2; expected_count++ {
				_, status := view("ViewCountDoc")
				if status != http.StatusOK {
					t.Fatalf("normal View_w status = %d, want %d", status, http.StatusOK)
				}
				count := tool.Get_document_setting_value(db, "ViewCountDoc", "view_count", "")
				if count != strconv.Itoa(expected_count) {
					t.Fatalf("normal page view count = %q, want %d", count, expected_count)
				}
			}

			redirect_output, redirect_status := view("RedirectViewCountDoc")
			redirect_target, redirected := tool.Get_redirect_target(redirect_output)
			if redirect_status != http.StatusOK || !redirected || redirect_target != "/w_from/RedirectViewTarget" {
				t.Fatalf("redirect View_w result = (%d, %q), want redirect to /w_from/RedirectViewTarget", redirect_status, redirect_output)
			}
			if count := tool.Get_document_setting_value(db, "RedirectViewCountDoc", "view_count", ""); count != "1" {
				t.Fatalf("redirect page view count = %q, want 1", count)
			}
		})
	}
}
