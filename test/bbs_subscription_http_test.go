package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_bbs_subscription_registered_http_routes(t *testing.T) {
	db := new_record_db(t, "sqlite")
	defer tool.DB_close(db)
	Seed_bbs_subscription_board(db, t, "88", "normal")
	for _, user_id := range []string{"http_watcher_one", "http_watcher_two"} {
		if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', ?, 'owner')`, user_id); err != nil {
			t.Fatalf("seed HTTP test user %q: %v", user_id, err)
		}
	}

	admin_config := tool.Config{IP: "admin_ip", Session: &emptySession{}}
	initial_post := route.Api_bbs_w_edit_post(admin_config, "88", "", "", "HTTP initial post", "initial body", "", "", "")
	if initial_post["response"] != "ok" || initial_post["data"] != "1" {
		t.Fatalf("create initial post: %#v", initial_post)
	}

	previous_mode := gin.Mode()
	gin.SetMode(gin.TestMode)
	t.Cleanup(func() { gin.SetMode(previous_mode) })
	t.Setenv("NAMU_SESSION_KEY", "bbs_subscription_http_test_session_key")

	session_middleware := tool.Session_middleware()
	user_one_cookie := Bbs_subscription_http_session_cookie(t, "http_watcher_one", session_middleware)
	user_two_cookie := Bbs_subscription_http_session_cookie(t, "http_watcher_two", session_middleware)
	admin_cookie := Bbs_subscription_http_session_cookie(t, "admin_ip", session_middleware)

	http_router := gin.New()
	http_router.Use(session_middleware)
	http_router.Use(route.Site_view_middleware())
	http_router.GET("/_test/identity", func(c *gin.Context) { c.String(http.StatusOK, tool.Get_IP(c)) })
	route.Register_bbs_subscription_routes(http_router)
	route.Register_bbs_subscription_api_routes(http_router)
	route.Register_bbs_post_create_route(http_router)
	server := httptest.NewServer(http_router)
	t.Cleanup(server.Close)
	identity_status, identity_body := Bbs_subscription_http_request(t, server.URL, http.MethodGet, "/_test/identity", user_one_cookie, nil)
	if identity_status != http.StatusOK || string(identity_body) != "http_watcher_one" {
		t.Fatalf("session identity = %d %s, want authenticated user", identity_status, identity_body)
	}

	anonymous_status, anonymous_body := Bbs_subscription_http_request(t, server.URL, http.MethodPost, "/api/v2/bbs_bookmark/88/1", nil, nil)
	if anonymous_status != http.StatusForbidden {
		var response map[string]any
		if anonymous_status != http.StatusOK || json.Unmarshal(anonymous_body, &response) != nil || response["response"] != "require auth" {
			t.Fatalf("anonymous bookmark mutation status/body = %d %s, want denial", anonymous_status, anonymous_body)
		}
	}
	if count := Bbs_subscription_count(db, t, "bbs_bookmark", "http_watcher_one", "88-1"); count != 0 {
		t.Fatalf("anonymous request created %d bookmark rows", count)
	}

	status, body := Bbs_subscription_http_request(t, server.URL, http.MethodPost, "/api/v2/bbs_bookmark/88/1", user_one_cookie, nil)
	Bbs_subscription_http_expect_response(t, status, body, "ok")
	if count := Bbs_subscription_count(db, t, "bbs_bookmark", "http_watcher_one", "88-1"); count != 1 {
		t.Fatalf("authenticated bookmark add left %d rows, want 1", count)
	}
	status, body = Bbs_subscription_http_request(t, server.URL, http.MethodGet, "/bbs_bookmark/88/1", user_one_cookie, nil)
	if status != http.StatusOK || !strings.Contains(string(body), tool.Get_language(db, "bbs_post_bookmark_remove", true)) {
		t.Fatalf("first user's current bookmark state = %d %s, want remove action", status, body)
	}
	status, body = Bbs_subscription_http_request(t, server.URL, http.MethodGet, "/bbs_bookmark_list", user_one_cookie, nil)
	if status != http.StatusOK || !strings.Contains(string(body), "HTTP initial post") || !strings.Contains(string(body), "/bbs_bookmark/88/1") {
		t.Fatalf("first user's bookmark management list omitted the saved post: %d %s", status, body)
	}

	// GET on JSON mutation routes must not toggle state.
	status, _ = Bbs_subscription_http_request(t, server.URL, http.MethodGet, "/api/v2/bbs_bookmark/88/1", user_one_cookie, nil)
	if status < http.StatusBadRequest {
		t.Fatalf("GET bookmark mutation route returned status %d, want method/path denial", status)
	}
	if count := Bbs_subscription_count(db, t, "bbs_bookmark", "http_watcher_one", "88-1"); count != 1 {
		t.Fatalf("GET on bookmark mutation route changed state: %d rows", count)
	}

	// The second authenticated session sees its own empty list and cannot
	// remove the first user's bookmark by issuing the same toggle.
	status, body = Bbs_subscription_http_request(t, server.URL, http.MethodGet, "/bbs_bookmark_list", user_two_cookie, nil)
	if status != http.StatusOK || strings.Contains(string(body), "HTTP initial post") {
		t.Fatalf("second user's bookmark list leaked another user's post: %d %s", status, body)
	}
	status, body = Bbs_subscription_http_request(t, server.URL, http.MethodGet, "/bbs_bookmark/88/1", user_two_cookie, nil)
	if status != http.StatusOK || !strings.Contains(string(body), tool.Get_language(db, "bbs_post_bookmark", true)) {
		t.Fatalf("second user's current bookmark state = %d %s, want add action", status, body)
	}
	status, body = Bbs_subscription_http_request(t, server.URL, http.MethodPost, "/api/v2/bbs_bookmark/88/1", user_two_cookie, nil)
	Bbs_subscription_http_expect_response(t, status, body, "ok")
	if count := Bbs_subscription_count(db, t, "bbs_bookmark", "http_watcher_one", "88-1"); count != 1 {
		t.Fatalf("second user's toggle changed first user's bookmark: %d rows", count)
	}
	if count := Bbs_subscription_count(db, t, "bbs_bookmark", "http_watcher_two", "88-1"); count != 1 {
		t.Fatalf("second user's own bookmark row count = %d, want 1", count)
	}
	status, body = Bbs_subscription_http_request(t, server.URL, http.MethodPost, "/api/v2/bbs_bookmark/88/1", user_two_cookie, nil)
	Bbs_subscription_http_expect_response(t, status, body, "ok")

	status, _ = Bbs_subscription_http_request(t, server.URL, http.MethodGet, "/api/v2/bbs_board_watch/88", user_one_cookie, nil)
	if status < http.StatusBadRequest || Bbs_subscription_count(db, t, "bbs_board_watchlist", "http_watcher_one", "88") != 0 {
		t.Fatalf("GET board-watch API route changed state or was accepted: status=%d", status)
	}
	status, body = Bbs_subscription_http_request(t, server.URL, http.MethodPost, "/api/v2/bbs_board_watch/88", user_one_cookie, nil)
	Bbs_subscription_http_expect_response(t, status, body, "ok")
	if count := Bbs_subscription_count(db, t, "bbs_board_watchlist", "http_watcher_one", "88"); count != 1 {
		t.Fatalf("HTTP board-watch add left %d rows, want 1", count)
	}
	status, body = Bbs_subscription_http_request(t, server.URL, http.MethodGet, "/bbs_board_watch/88", user_one_cookie, nil)
	if status != http.StatusOK || !strings.Contains(string(body), tool.Get_language(db, "bbs_board_watch_remove", true)) {
		t.Fatalf("board-watch current state = %d %s, want remove action", status, body)
	}
	status, body = Bbs_subscription_http_request(t, server.URL, http.MethodGet, "/bbs_board_watch_list", user_one_cookie, nil)
	if status != http.StatusOK || !strings.Contains(string(body), "Subscription board") || !strings.Contains(string(body), "/bbs_board_watch/88") {
		t.Fatalf("board-watch management list omitted the board: %d %s", status, body)
	}

	tag := "Release & 支援"
	tag_token := tool.Base64_encode(tag)
	tag_route := "/api/v2/bbs_tag_watch/" + tag_token
	status, body = Bbs_subscription_http_request(t, server.URL, http.MethodGet, tag_route, user_one_cookie, nil)
	if status < http.StatusBadRequest {
		t.Fatalf("GET tag mutation route returned status %d, want method/path denial", status)
	}
	if count := Bbs_subscription_count(db, t, "bbs_tag_watchlist", "http_watcher_one", tag); count != 0 {
		t.Fatalf("GET tag mutation route changed state: %d rows", count)
	}
	status, body = Bbs_subscription_http_request(t, server.URL, http.MethodPost, tag_route, user_one_cookie, nil)
	Bbs_subscription_http_expect_response(t, status, body, "ok")
	if count := Bbs_subscription_count(db, t, "bbs_tag_watchlist", "http_watcher_one", tag); count != 1 {
		t.Fatalf("HTTP global tag-watch add left %d rows, want 1", count)
	}
	for _, path := range []string{
		"/bbs_tag_watch/" + tag_token,
		"/bbs_tag_watch/88/" + tag_token,
	} {
		status, body = Bbs_subscription_http_request(t, server.URL, http.MethodGet, path, user_one_cookie, nil)
		if status != http.StatusOK || !strings.Contains(string(body), tool.Get_language(db, "bbs_tag_watch_remove", true)) {
			t.Fatalf("global or legacy contextual tag route %q state = %d %s", path, status, body)
		}
	}
	status, body = Bbs_subscription_http_request(t, server.URL, http.MethodGet, "/bbs_tag_watch_list", user_one_cookie, nil)
	if status != http.StatusOK || !strings.Contains(string(body), "Release &amp; 支援") {
		t.Fatalf("tag management route omitted escaped global tag: %d %s", status, body)
	}

	// The actual production BBS post-creation route is registered here as well.
	// A recipient subscribed to both the board and matching tag gets one notice.
	post_form := url.Values{
		"title":   {"HTTP & tagged post"},
		"content": {"post body"},
		"tags":    {tag},
	}
	status, body, location := Bbs_subscription_http_request_with_location(t, server.URL, http.MethodPost, "/bbs/edit/88", admin_cookie, post_form)
	if status != http.StatusFound || location != "/bbs/w/88/2" {
		t.Fatalf("HTTP post creation status/body/location = %d %s %q, want redirect to new post", status, body, location)
	}
	var created_post_count int
	if err := db.QueryRow(`select count(*) from bbs_data where set_name = 'title' and set_id = '88' and set_code = '2'`).Scan(&created_post_count); err != nil || created_post_count != 1 {
		t.Fatalf("HTTP-created tagged post count = %d, err=%v", created_post_count, err)
	}
	if count := Bbs_user_notice_count(db, t, "http_watcher_one"); count != 1 {
		t.Fatalf("board+tag duplicate HTTP notification count = %d, want 1", count)
	}
	var notice_data string
	if err := db.QueryRow(`select data from user_notice where name = 'http_watcher_one'`).Scan(&notice_data); err != nil {
		t.Fatalf("read HTTP notification: %v", err)
	}
	if !strings.Contains(notice_data, "HTTP &amp; tagged post") || strings.Contains(notice_data, "HTTP & tagged post") {
		t.Fatalf("notification title was not HTML-escaped: %q", notice_data)
	}

	if comment := route.Api_bbs_w_comment_post(admin_config, "88", "2", "", "ordinary comment"); comment["response"] != "ok" {
		t.Fatalf("create comment after tagged post: %#v", comment)
	}
	if edited := route.Api_bbs_w_edit_post(admin_config, "88", "2", "", "HTTP tagged post edited", "edited body", "", tag, ""); edited["response"] != "ok" {
		t.Fatalf("edit tagged post: %#v", edited)
	}
	if count := Bbs_user_notice_count(db, t, "http_watcher_one"); count != 1 {
		t.Fatalf("comment or edit triggered another board/tag notification: %d, want 1", count)
	}

	// Exercise the existing HTML POST controls and the old contextual tag URL.
	for _, test_case := range []struct {
		path     string
		location string
	}{
		{"/bbs_bookmark/88/1", "/bbs/w/88/1"},
		{"/bbs_board_watch/88", "/bbs/in/88"},
		{"/bbs_tag_watch/88/" + tag_token, "/bbs/in/88/filter/tag/" + tool.Url_parser(tag) + "/1"},
	} {
		status, _, location = Bbs_subscription_http_request_with_location(t, server.URL, http.MethodPost, test_case.path, user_one_cookie, url.Values{})
		if status != http.StatusFound || location != test_case.location {
			t.Fatalf("HTML removal route %q returned %d redirect %q, want %q", test_case.path, status, location, test_case.location)
		}
	}
	for _, entry := range []struct {
		watch_type string
		value      string
	}{
		{"bbs_bookmark", "88-1"},
		{"bbs_board_watchlist", "88"},
		{"bbs_tag_watchlist", tag},
	} {
		if count := Bbs_subscription_count(db, t, entry.watch_type, "http_watcher_one", entry.value); count != 0 {
			t.Errorf("HTTP remove left %s %q row count %d", entry.watch_type, entry.value, count)
		}
	}
	status, body = Bbs_subscription_http_request(t, server.URL, http.MethodGet, "/bbs_tag_watch/"+tag_token, user_one_cookie, nil)
	if status != http.StatusOK || !strings.Contains(string(body), tool.Get_language(db, "bbs_tag_watch", true)) {
		t.Errorf("global tag state after legacy contextual removal = %d %s, want add action", status, body)
	}
}

func Bbs_subscription_http_session_cookie(t *testing.T, user_id string, session_middleware gin.HandlerFunc) *http.Cookie {
	t.Helper()
	issuer := gin.New()
	issuer.Use(session_middleware)
	issuer.GET("/session/:user_id", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("id", c.Param("user_id"))
		if err := session.Save(); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/session/"+user_id, nil)
	request.RemoteAddr = "192.0.2.10:12345"
	recorder := httptest.NewRecorder()
	issuer.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("issue session for %q: status=%d body=%s", user_id, recorder.Code, recorder.Body.String())
	}
	for _, cookie := range recorder.Result().Cookies() {
		if cookie.Name == "opennamu_session" {
			return cookie
		}
	}
	t.Fatalf("session middleware did not issue a cookie for %q", user_id)
	return nil
}

func Bbs_subscription_http_request(t *testing.T, base_url string, method string, path string, cookie *http.Cookie, form url.Values) (int, []byte) {
	t.Helper()
	status, body, _ := Bbs_subscription_http_request_with_location(t, base_url, method, path, cookie, form)
	return status, body
}

func Bbs_subscription_http_request_with_location(t *testing.T, base_url string, method string, path string, cookie *http.Cookie, form url.Values) (int, []byte, string) {
	t.Helper()
	var request_body io.Reader
	if form != nil {
		request_body = strings.NewReader(form.Encode())
	}
	request, err := http.NewRequest(method, base_url+path, request_body)
	if err != nil {
		t.Fatalf("create HTTP request %s %s: %v", method, path, err)
	}
	if form != nil {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if cookie != nil {
		request.AddCookie(cookie)
	}
	client := &http.Client{CheckRedirect: func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }}
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("send HTTP request %s %s: %v", method, path, err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("read HTTP response %s %s: %v", method, path, err)
	}
	return response.StatusCode, body, response.Header.Get("Location")
}

func Bbs_subscription_http_expect_response(t *testing.T, status int, body []byte, expected string) {
	t.Helper()
	var response map[string]any
	if status != http.StatusOK || json.Unmarshal(body, &response) != nil || response["response"] != expected {
		t.Fatalf("HTTP API response = %d %s, want status 200 response %q", status, body, expected)
	}
}
