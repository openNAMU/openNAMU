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

func Test_view_w_error_status_and_message(t *testing.T) {
	db := Setup_raw_test_db(t)
	defer tool.DB_close(db)

	make_context := func(path string) *gin.Context {
		context, _ := gin.CreateTestContext(httptest.NewRecorder())
		context.Request = httptest.NewRequest(http.MethodGet, path, nil)
		return context
	}

	regular := tool.Config{IP: "regular_ip", Session: &emptySession{}}
	restricted_output, restricted_status := route.View_w(make_context("/w/RestrictedDoc"), regular, "RestrictedDoc", "")
	if restricted_status != http.StatusUnauthorized {
		t.Fatalf("restricted document status = %d, want %d", restricted_status, http.StatusUnauthorized)
	}
	if !strings.Contains(restricted_output, tool.Get_language(db, "authority_error", true)) {
		t.Fatalf("restricted document output does not show the authorization message: %q", restricted_output)
	}

	admin := tool.Config{IP: "admin_ip", Session: &emptySession{}}
	missing_output, missing_status := route.View_w(make_context("/w/NoSuchDoc"), admin, "NoSuchDoc", "")
	if missing_status != http.StatusNotFound {
		t.Fatalf("missing document status = %d, want %d", missing_status, http.StatusNotFound)
	}
	if !strings.Contains(missing_output, tool.Get_language(db, "document_404_error", true)) {
		t.Fatalf("missing document output does not show the not-found message: %q", missing_output)
	}
}
