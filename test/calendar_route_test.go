package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"

	"github.com/gin-gonic/gin"
)

func Test_calendar_http_routes(t *testing.T) {
	db := audit_setup(t, "sqlite")
	audit_seed_auth(db, t)

	admin := calendar_test_config("admin_ip")
	added := route.Api_calendar_add(admin, "공지", "설명")
	id, _ := added["data"].(string)
	event := route.Api_calendar_event_add(admin, id, "2024-02-29", "윤년 행사", "메모")
	code, _ := event["data"].(string)
	if added["response"] != "ok" || event["response"] != "ok" {
		t.Fatalf("seed calendar/event = %#v / %#v", added, event)
	}

	prev_mode := gin.Mode()
	gin.SetMode(gin.TestMode)
	t.Cleanup(func() { gin.SetMode(prev_mode) })

	r := gin.New()
	r.Use(tool.Session_middleware())
	route.Register_calendar_routes(r)

	do := func(method string, path string, values url.Values) *httptest.ResponseRecorder {
		var body io.Reader
		if values != nil {
			body = strings.NewReader(values.Encode())
		}
		req := httptest.NewRequest(method, path, body)
		if values != nil {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}

	w := do(http.MethodGet, "/calendar", nil)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "공지") {
		t.Fatalf("GET /calendar = %d, body missing calendar name", w.Code)
	}

	w = do(http.MethodGet, "/calendar/"+id+"/month/2024-02", nil)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "윤년 행사") {
		t.Fatalf("GET calendar month = %d, body missing event", w.Code)
	}
	if !strings.Contains(w.Body.String(), `id="event-`+code+`"`) {
		t.Fatalf("GET calendar month missing event anchor for %q", code)
	}

	w = do(http.MethodGet, "/calendar/"+id, nil)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "공지") {
		t.Fatalf("GET calendar = %d, body missing calendar name", w.Code)
	}

	w = do(http.MethodGet, "/calendar/"+id+"/month/2024-13", nil)
	if w.Code != http.StatusFound || w.Header().Get("Location") != "/calendar/"+id {
		t.Fatalf("GET invalid month = %d %q, want redirect to calendar", w.Code, w.Header().Get("Location"))
	}

	w = do(http.MethodGet, "/calendar/999", nil)
	if w.Code != http.StatusFound || w.Header().Get("Location") != "/calendar" {
		t.Fatalf("GET missing calendar = %d %q, want redirect to /calendar", w.Code, w.Header().Get("Location"))
	}

	w = do(http.MethodGet, "/calendar/add", nil)
	if w.Code != http.StatusOK || strings.Contains(w.Body.String(), `id="calendar_name"`) {
		t.Fatalf("anonymous GET /calendar/add = %d, want auth page without form", w.Code)
	}

	do(http.MethodPost, "/calendar/add", url.Values{"name": {"새 캘린더"}})
	listed := route.Api_calendar_list(admin)
	list, _ := listed["data"].([][]string)
	if len(list) != 1 {
		t.Fatalf("anonymous calendar add created a calendar: %#v", list)
	}

	w = do(http.MethodPost, "/calendar/"+id+"/event/add", url.Values{"date": {"2024-02-01"}, "title": {"새 일정"}})
	if w.Code != http.StatusOK {
		t.Fatalf("POST event add = %d %q, want auth page on the add route", w.Code, w.Header().Get("Location"))
	}
	viewed := route.Api_calendar_view(admin, id, "2024-02")
	events, _ := viewed["events"].([][]string)
	if len(events) != 1 {
		t.Fatalf("anonymous event add created an event: %#v", events)
	}

	w = do(http.MethodGet, "/calendar/"+id+"/event/"+code, nil)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "윤년 행사") {
		t.Fatalf("GET event = %d, body missing event", w.Code)
	}

	w = do(http.MethodGet, "/calendar/"+id+"/event/999", nil)
	if w.Code != http.StatusFound || w.Header().Get("Location") != "/calendar/"+id {
		t.Fatalf("GET missing event = %d %q, want redirect to calendar", w.Code, w.Header().Get("Location"))
	}

	w = do(http.MethodGet, "/calendar/"+id+"/delete", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("anonymous GET calendar delete = %d, want auth page", w.Code)
	}

	do(http.MethodPost, "/calendar/"+id+"/delete", nil)
	if got := route.Api_calendar_view(admin, id, ""); got["response"] != "ok" {
		t.Fatalf("anonymous calendar delete removed the calendar: %#v", got)
	}
}
