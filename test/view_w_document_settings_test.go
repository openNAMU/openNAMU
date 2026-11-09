package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"opennamu/route"
	"opennamu/route/tool"
)

func Test_View_w_uses_existing_document_settings_and_markup(t *testing.T) {
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
			exec(`insert into other (name, data, coverage) values ('related_document_enabled', '0', '')`)
			exec(`insert into other (name, data, coverage) values ('outdated_doc_warning_date', '1', '')`)
			exec(`insert into other (name, data, coverage) values ('outdated_doc_warning', 'OUTDATED_VIEW_WARNING', '')`)
			exec(`insert into data (title, data, type) values ('ViewSettingsDoc', 'document body', '')`)
			exec(`insert into data_set (doc_name, doc_rev, set_name, set_data) values ('ViewSettingsDoc', '', 'last_edit', '2000-01-01 00:00:00')`)
			exec(`insert into data_set (doc_name, doc_rev, set_name, set_data) values ('ViewSettingsDoc', '', 'document_top', '<b>VIEW_TOP_MARKER</b>')`)
			exec(`insert into data_set (doc_name, doc_rev, set_name, set_data) values ('ViewSettingsDoc', '', 'document_top_markup', 'plain')`)

			gin.SetMode(gin.TestMode)
			context, _ := gin.CreateTestContext(httptest.NewRecorder())
			context.Request = httptest.NewRequest(http.MethodGet, "/w/ViewSettingsDoc", nil)
			config := tool.Config{IP: "admin_ip", Session: &emptySession{}}
			page, status := route.View_w(context, config, "ViewSettingsDoc", "")
			if status != http.StatusOK {
				t.Fatalf("View_w status = %d, want %d", status, http.StatusOK)
			}
			if !strings.Contains(page, "OUTDATED_VIEW_WARNING") {
				t.Fatal("View_w did not apply the stored last-edit warning")
			}
			if !strings.Contains(page, "&lt;b&gt;VIEW_TOP_MARKER&lt;/b&gt;") || strings.Contains(page, "<b>VIEW_TOP_MARKER</b>") {
				t.Fatalf("View_w did not apply the stored plain document-top markup: %s", page)
			}
		})
	}
}
