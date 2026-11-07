package test

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"opennamu/route"
	"opennamu/route/tool"
	"opennamu/route/tool/markup"
)

func Test_namumark_include_acl_uses_server_requester_not_render_parameters(t *testing.T) {
	db := Setup_raw_test_db(t)
	defer tool.DB_close(db)

	const secret_marker = "RESTRICTED_INCLUDE_MARKER"
	if _, err := db.Exec(`update data set data = ? where title = 'RestrictedDoc'`, secret_marker); err != nil {
		t.Fatalf("seed restricted document: %v", err)
	}
	if _, err := db.Exec(`insert into other (name, data, coverage) values ('link_case_insensitive', '1', '')`); err != nil {
		t.Fatalf("seed case-insensitive document links: %v", err)
	}
	if _, err := db.Exec(`insert into data (title, data) values ('BridgeDoc', ?)`, "{{{#!wiki\n[include(RestrictedDoc)]\n}}}"); err != nil {
		t.Fatalf("seed nested include document: %v", err)
	}

	if tool.Check_acl(db, "RestrictedDoc", "", "render", "regular_ip") {
		t.Fatal("fixture error: regular_ip unexpectedly has view access")
	}
	if !tool.Check_acl(db, "RestrictedDoc", "", "render", "admin_ip") {
		t.Fatal("fixture error: admin_ip unexpectedly lacks view access")
	}

	regular_config := tool.Config{IP: "regular_ip", Session: &emptySession{}}
	admin_config := tool.Config{IP: "admin_ip", Session: &emptySession{}}
	include_cases := []struct {
		name   string
		data   string
		option string
	}{
		{
			name:   "client option identity",
			data:   "before [include(RestrictedDoc)] after",
			option: `{"ip":"admin_ip"}`,
		},
		{
			name: "include argument identity",
			data: "before [include(RestrictedDoc,ip=admin_ip)] after",
		},
		{
			name:   "case-insensitive target uses actual document ACL",
			data:   "before [include(restricteddoc)] after",
			option: `{"ip":"admin_ip"}`,
		},
		{
			name:   "interwiki child inherits trusted requester",
			data:   "{{{#!wiki\n[include(RestrictedDoc,ip=admin_ip)]\n}}}",
			option: `{"ip":"admin_ip"}`,
		},
		{
			name: "included interwiki child cannot bypass ACL",
			data: "[include(BridgeDoc)]",
		},
	}

	for _, test_case := range include_cases {
		t.Run(test_case.name, func(t *testing.T) {
			result := route.Api_w_render(regular_config, "PublicDoc", test_case.data, "api_view", test_case.option)
			if result["response"] != "ok" {
				t.Fatalf("render failed: %v", result)
			}
			rendered, ok := result["data"].(string)
			if !ok {
				t.Fatalf("render data has unexpected type %T", result["data"])
			}
			if strings.Contains(rendered, secret_marker) {
				t.Fatalf("restricted include leaked to regular_ip: %s", rendered)
			}
		})
	}

	admin_result := route.Api_w_render(admin_config, "PublicDoc", "before [include(BridgeDoc)] after", "api_view", "")
	admin_rendered, ok := admin_result["data"].(string)
	if admin_result["response"] != "ok" || !ok || !strings.Contains(admin_rendered, secret_marker) {
		t.Fatalf("authorized requester should retain include access, result=%v", admin_result)
	}

	for _, identity := range []string{"", tool.Invalid_remote_address} {
		result := route.Api_w_render(
			tool.Config{IP: identity, Session: &emptySession{}},
			"PublicDoc",
			"[include(RestrictedDoc)]",
			"api_view",
			`{"ip":"admin_ip"}`,
		)
		rendered, ok := result["data"].(string)
		if result["response"] != "ok" || !ok || strings.Contains(rendered, secret_marker) {
			t.Fatalf("untrusted requester %q should fail closed, result=%v", identity, result)
		}
	}

	thread_result := route.Api_w_render(regular_config, "PublicDoc", "[include(RestrictedDoc)]", "api_thread", "")
	thread_rendered, ok := thread_result["data"].(string)
	if thread_result["response"] != "ok" || !ok || strings.Contains(thread_rendered, secret_marker) || !strings.Contains(thread_rendered, "/include/") {
		t.Fatalf("thread mode should retain a non-expanding include link: %v", thread_result)
	}
	bbs_rendered := route.Get_bbs_render(db, "unset_board", "[include(RestrictedDoc)]", "bbs", regular_config)
	if strings.Contains(bbs_rendered, secret_marker) || !strings.Contains(bbs_rendered, "/include/") {
		t.Fatalf("BBS mode should retain a non-expanding include link: %s", bbs_rendered)
	}
}

func Test_namumark_include_acl_protects_stored_history_and_include_views(t *testing.T) {
	db := Setup_raw_test_db(t)
	defer tool.DB_close(db)

	const secret_marker = "RESTRICTED_INCLUDE_MARKER"
	if _, err := db.Exec(`update data set data = ? where title = 'RestrictedDoc'`, secret_marker); err != nil {
		t.Fatalf("seed restricted document: %v", err)
	}
	if _, err := db.Exec(`update data set data = ? where title = 'PublicDoc'`, "stored [include(RestrictedDoc)]"); err != nil {
		t.Fatalf("seed stored include: %v", err)
	}
	if _, err := db.Exec(`update history set data = ? where title = 'RevDoc' and id = '1'`, "historical [include(RestrictedDoc)]"); err != nil {
		t.Fatalf("seed historical include: %v", err)
	}
	if _, err := db.Exec(`insert into other (name, data, coverage) values ('body', '[include(RestrictedDoc)]', '')`); err != nil {
		t.Fatalf("seed body setting: %v", err)
	}
	if _, err := db.Exec(`insert into other (name, data, coverage) values ('body_markup', 'namumark_beta', '')`); err != nil {
		t.Fatalf("seed body markup: %v", err)
	}
	if _, err := db.Exec(`insert into data_set (doc_name, doc_rev, set_name, set_data) values ('PublicDoc', '', 'document_editor_top', '[include(RestrictedDoc)]')`); err != nil {
		t.Fatalf("seed editor-top content: %v", err)
	}
	if _, err := db.Exec(`insert into data_set (doc_name, doc_rev, set_name, set_data) values ('PublicDoc', '', 'document_editor_top_markup', 'namumark_beta')`); err != nil {
		t.Fatalf("seed editor-top markup: %v", err)
	}

	regular_config := tool.Config{IP: "regular_ip", Session: &emptySession{}}
	admin_config := tool.Config{IP: "admin_ip", Session: &emptySession{}}

	view_page := func(config tool.Config) string {
		t.Helper()
		gin.SetMode(gin.TestMode)
		context, _ := gin.CreateTestContext(httptest.NewRecorder())
		context.Request = httptest.NewRequest(http.MethodGet, "/w/PublicDoc", nil)
		page, status := route.View_w(context, config, "PublicDoc", "")
		if status != http.StatusOK {
			t.Fatalf("View_w returned status %d", status)
		}
		return page
	}

	if page := view_page(regular_config); strings.Contains(page, secret_marker) {
		t.Fatalf("stored page or global setting leaked restricted include: %s", page)
	}
	if page := view_page(admin_config); !strings.Contains(page, secret_marker) {
		t.Fatalf("authorized requester lost stored/global include content: %s", page)
	}

	regular_history := route.View_render(regular_config, "RevDoc", "1")
	if strings.Contains(regular_history, secret_marker) {
		t.Fatalf("historical rendering leaked restricted include: %s", regular_history)
	}
	admin_history := route.View_render(admin_config, "RevDoc", "1")
	if !strings.Contains(admin_history, secret_marker) {
		t.Fatalf("authorized historical rendering lost include content: %s", admin_history)
	}

	include_payload := `{"version":1,"name":"RestrictedDoc","parameters":{"ip":"admin_ip"}}`
	regular_include := route.View_include(regular_config, tool.Base64_encode(include_payload))
	if strings.Contains(regular_include, secret_marker) {
		t.Fatalf("include view accepted a client-supplied identity: %s", regular_include)
	}
	admin_include := route.View_include(admin_config, tool.Base64_encode(include_payload))
	if !strings.Contains(admin_include, secret_marker) {
		t.Fatalf("authorized include view lost target content: %s", admin_include)
	}

	regular_editor := route.View_edit_page(db, regular_config, "PublicDoc", "", "", "", "", "")
	if strings.Contains(regular_editor, secret_marker) {
		t.Fatalf("editor-top markup leaked restricted include: %s", regular_editor)
	}
	admin_editor := route.View_edit_page(db, admin_config, "PublicDoc", "", "", "", "", "")
	if !strings.Contains(admin_editor, secret_marker) {
		t.Fatalf("authorized editor-top rendering lost include content: %s", admin_editor)
	}
}

func Test_namumark_include_acl_sqlite_and_embedded_turso(t *testing.T) {
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
			exec(`insert into alist (name, acl) values ('regular_group', 'view')`)
			exec(`insert into user_set (name, id, data) values ('acl', 'regular_ip', 'regular_group')`)
			exec(`insert into data (title, data, type) values ('PublicDoc', 'public', '')`)
			exec(`insert into data (title, data, type) values ('RestrictedDoc', 'RESTRICTED_INCLUDE_MARKER', '')`)
			exec(`insert into acl (title, data, type) values ('RestrictedDoc', 'owner', 'view')`)

			result := route.Api_w_render(
				tool.Config{IP: "regular_ip", Session: &emptySession{}},
				"PublicDoc",
				"[include(RestrictedDoc)]",
				"api_view",
				`{"ip":"admin_ip"}`,
			)
			rendered, ok := result["data"].(string)
			if result["response"] != "ok" || !ok || strings.Contains(rendered, "RESTRICTED_INCLUDE_MARKER") {
				t.Fatalf("%s leaked restricted include through a forged option: %v", backend, result)
			}
		})
	}
}

func Test_namumark_include_acl_http_end_to_end(t *testing.T) {
	db := Setup_raw_test_db(t)
	defer tool.DB_close(db)
	t.Setenv("NAMU_SESSION_KEY", "synthetic-ghsa-include-test-secret")

	const secret_marker = "RESTRICTED_INCLUDE_MARKER"
	const public_marker = "PUBLIC_INCLUDE_MARKER"
	const filter_marker = "FILTERED_INCLUDE_MARKER"
	const user_marker = "USER_PAGE_INCLUDE_MARKER"
	exec := func(query string, args ...any) {
		t.Helper()
		if _, err := db.Exec(query, args...); err != nil {
			t.Fatalf("seed %q: %v", query, err)
		}
	}
	exec(`update data set data = ? where title = 'RestrictedDoc'`, secret_marker)
	exec(`insert into data (title, data) values ('PublicInclude', ?)`, public_marker)
	exec(`insert into data (title, data) values ('FilteredInclude', ?)`, filter_marker)
	exec(`insert into data (title, data) values ('user:alice', ?)`, user_marker)
	exec(`insert into data (title, data) values ('RedirectInclude', '#redirect RestrictedDoc')`)
	exec(`insert into data (title, data) values ('BridgeDoc', '[include(RestrictedDoc)]')`)
	exec(`insert into data (title, data) values ('StoredParent', 'stored [include(RestrictedDoc)] and [include(PublicInclude)]')`)
	exec(`insert into data (title, data) values ('PreviewParent', 'preview base')`)
	exec(`insert into html_filter (html, kind, plus, plus_t) values ('FilteredInclude', 'document', '^FilteredInclude$', 'view=owner')`)
	exec(`delete from other where name = 'user_document_view_acl_all'`)
	exec(`insert into other (name, data, coverage) values ('user_document_view_acl_all', 'owner', '')`)
	exec(`insert into alist (name, acl) values ('ip', 'view')`)
	exec(`insert into alist (name, acl) values ('regular_group', 'edit')`)
	exec(`insert into acl (title, data, type) values ('PreviewParent', 'regular_group', 'decu')`)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(tool.Session_middleware())
	router.GET("/__audit/session/:id", func(c *gin.Context) {
		session := tool.Get_session(c)
		session.Set("id", c.Param("id"))
		if err := session.Save(); err != nil {
			c.String(http.StatusInternalServerError, "session save failed")
			return
		}
		c.Status(http.StatusNoContent)
	})
	request_config := func(c *gin.Context) tool.Config {
		return tool.Config{
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   tool.Get_session(c),
			UserAgent: c.GetHeader("User-Agent"),
		}
	}
	router.POST("/api/render", func(c *gin.Context) {
		c.JSON(http.StatusOK, route.Api_w_render(request_config(c), c.PostForm("name"), c.PostForm("data"), "api_view", c.PostForm("option")))
	})
	router.POST("/api/render/:tool", func(c *gin.Context) {
		render_type := "api_thread"
		switch c.Param("tool") {
		case "from":
			render_type = "api_from"
		case "include":
			render_type = "api_include"
		case "backlink":
			render_type = "backlink"
		}
		c.JSON(http.StatusOK, route.Api_w_render(request_config(c), c.PostForm("name"), c.PostForm("data"), render_type, c.PostForm("option")))
	})
	router.GET("/w/*doc_name", func(c *gin.Context) {
		page, status := route.View_w(c, request_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "")
		c.Data(status, "text/html; charset=utf-8", []byte(page))
	})
	router.POST("/edit_preview/:mode/*doc_name", func(c *gin.Context) {
		page := route.View_edit_preview(request_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.PostForm("content"), c.Param("mode"), c.PostForm("send"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(page))
	})
	router.GET("/include/:payload", func(c *gin.Context) {
		page := route.View_include(request_config(c), c.Param("payload"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(page))
	})

	server := httptest.NewServer(router)
	defer server.Close()
	new_client := func() *http.Client {
		jar, err := cookiejar.New(nil)
		if err != nil {
			t.Fatalf("create cookie jar: %v", err)
		}
		return &http.Client{Jar: jar}
	}
	request := func(client *http.Client, method string, path string, form url.Values) (int, string) {
		t.Helper()
		var request_body io.Reader
		if form != nil {
			request_body = strings.NewReader(form.Encode())
		}
		req, err := http.NewRequest(method, server.URL+path, request_body)
		if err != nil {
			t.Fatalf("create %s request for %s: %v", method, path, err)
		}
		if form != nil {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
		response, err := client.Do(req)
		if err != nil {
			t.Fatalf("send %s request for %s: %v", method, path, err)
		}
		defer response.Body.Close()
		response_data, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("read %s response for %s: %v", method, path, err)
		}
		return response.StatusCode, string(response_data)
	}
	anonymous_client := new_client()
	regular_client := new_client()
	admin_client := new_client()
	for _, identity := range []struct {
		client *http.Client
		id     string
	}{{regular_client, "regular_ip"}, {admin_client, "admin_ip"}} {
		if status, _ := request(identity.client, http.MethodGet, "/__audit/session/"+identity.id, nil); status != http.StatusNoContent {
			t.Fatalf("establish %s session: HTTP %d", identity.id, status)
		}
	}

	api_form := func(data string, option string) url.Values {
		return url.Values{"name": {"PublicDoc"}, "data": {data}, "option": {option}}
	}
	for _, identity := range []struct {
		name   string
		client *http.Client
	}{{"anonymous", anonymous_client}, {"regular", regular_client}} {
		for _, test_case := range []struct {
			name   string
			data   string
			option string
		}{
			{"option identity spoof", "[include(RestrictedDoc)]", `{"ip":"admin_ip"}`},
			{"include argument identity spoof", "[include(RestrictedDoc,ip=admin_ip)]", ""},
			{"case-insensitive target", "[include(restricteddoc)]", `{"ip":"admin_ip"}`},
			{"nested include", "[include(BridgeDoc)]", ""},
		} {
			t.Run(identity.name+"/"+test_case.name, func(t *testing.T) {
				status, body := request(identity.client, http.MethodPost, "/api/render", api_form(test_case.data, test_case.option))
				if status != http.StatusOK || strings.Contains(body, secret_marker) {
					t.Fatalf("HTTP render leaked restricted content (status=%d): %s", status, body)
				}
			})
		}
		status, body := request(identity.client, http.MethodPost, "/api/render", api_form("[include(PublicInclude)]", ""))
		if status != http.StatusOK || !strings.Contains(body, public_marker) {
			t.Fatalf("%s requester lost public include access (status=%d): %s", identity.name, status, body)
		}
	}

	for _, target := range []struct {
		name   string
		marker string
	}{
		{"RestrictedDoc", secret_marker},
		{"FilteredInclude", filter_marker},
		{"user:alice", user_marker},
	} {
		for _, identity := range []struct {
			name   string
			client *http.Client
		}{{"anonymous", anonymous_client}, {"regular", regular_client}} {
			status, body := request(identity.client, http.MethodPost, "/api/render", api_form("[include("+target.name+")]", ""))
			if status != http.StatusOK || strings.Contains(body, target.marker) {
				t.Fatalf("%s requester read %s include (status=%d): %s", identity.name, target.name, status, body)
			}
		}
		status, body := request(admin_client, http.MethodPost, "/api/render", api_form("[include("+target.name+")]", ""))
		if status != http.StatusOK || !strings.Contains(body, target.marker) {
			t.Fatalf("admin lost authorized %s include (status=%d): %s", target.name, status, body)
		}
	}
	if status, body := request(anonymous_client, http.MethodPost, "/api/render", api_form("[include(RedirectInclude)]", "")); status != http.StatusOK || strings.Contains(body, secret_marker) {
		t.Fatalf("redirecting include exposed its restricted target (status=%d): %s", status, body)
	}
	if status, body := request(admin_client, http.MethodPost, "/api/render", api_form("[include(RestrictedDoc)]", "")); status != http.StatusOK || !strings.Contains(body, secret_marker) {
		t.Fatalf("authorized include was denied (status=%d): %s", status, body)
	}

	for _, identity := range []struct {
		name   string
		client *http.Client
	}{
		{"anonymous", anonymous_client},
		{"regular", regular_client},
		{"admin", admin_client},
	} {
		status, body := request(identity.client, http.MethodGet, "/w/StoredParent", nil)
		if status != http.StatusOK {
			t.Fatalf("%s stored page returned HTTP %d: %s", identity.name, status, body)
		}
		if identity.name == "admin" {
			if !strings.Contains(body, secret_marker) || !strings.Contains(body, public_marker) {
				t.Fatalf("admin stored page lost authorized include content: %s", body)
			}
		} else if strings.Contains(body, secret_marker) || !strings.Contains(body, public_marker) {
			t.Fatalf("%s stored page include authorization mismatch: %s", identity.name, body)
		}
	}

	preview_form := url.Values{
		"content": {"preview [include(RestrictedDoc)] and [include(PublicInclude)]"},
		"send":    {""},
	}
	for _, identity := range []struct {
		name   string
		client *http.Client
	}{
		{"regular", regular_client},
		{"admin", admin_client},
	} {
		status, body := request(identity.client, http.MethodPost, "/edit_preview/normal/PreviewParent", preview_form)
		if status != http.StatusOK {
			t.Fatalf("%s preview returned HTTP %d", identity.name, status)
		}
		if identity.name == "regular" {
			if strings.Contains(body, secret_marker) || !strings.Contains(body, public_marker) {
				t.Fatalf("regular preview include authorization mismatch: %s", body)
			}
		} else if !strings.Contains(body, secret_marker) || !strings.Contains(body, public_marker) {
			t.Fatalf("admin preview lost authorized include content: %s", body)
		}
	}

	include_payload := tool.Base64_encode(`{"version":1,"name":"RestrictedDoc","parameters":{"ip":"admin_ip"}}`)
	for _, identity := range []struct {
		name   string
		client *http.Client
	}{
		{"anonymous", anonymous_client},
		{"regular", regular_client},
	} {
		status, body := request(identity.client, http.MethodGet, "/include/"+include_payload, nil)
		if status != http.StatusOK || strings.Contains(body, secret_marker) {
			t.Fatalf("%s /include view exposed restricted content (status=%d): %s", identity.name, status, body)
		}
	}
	if status, body := request(admin_client, http.MethodGet, "/include/"+include_payload, nil); status != http.StatusOK || !strings.Contains(body, secret_marker) {
		t.Fatalf("admin /include view lost authorized content (status=%d): %s", status, body)
	}

	for _, mode := range []string{"from", "include", "thread"} {
		status, body := request(regular_client, http.MethodPost, "/api/render/"+mode, api_form("[include(RestrictedDoc)]", `{"ip":"admin_ip"}`))
		if status != http.StatusOK || strings.Contains(body, secret_marker) {
			t.Fatalf("/api/render/%s leaked restricted content (status=%d): %s", mode, status, body)
		}
		if mode == "thread" && !strings.Contains(body, "/include/") {
			t.Fatalf("thread mode should keep a non-expanding include link: %s", body)
		}
	}
}

func Test_namumark_include_acl_fails_closed_without_a_trusted_requester(t *testing.T) {
	db := Setup_raw_test_db(t)
	defer tool.DB_close(db)

	const secret_marker = "RESTRICTED_INCLUDE_MARKER"
	if _, err := db.Exec(`update data set data = ? where title = 'RestrictedDoc'`, secret_marker); err != nil {
		t.Fatalf("seed restricted document: %v", err)
	}

	parameter_data := map[string]any{"ip": "admin_ip"}
	for _, requester := range []string{"", tool.Invalid_remote_address} {
		rendered := markup.Get_render_with_requester(
			db,
			"PublicDoc",
			"[include(RestrictedDoc)]",
			"api_view",
			parameter_data,
			requester,
		)["data"]
		if strings.Contains(rendered, secret_marker) {
			t.Fatalf("requester %q leaked restricted include: %s", requester, rendered)
		}
	}

	// The legacy renderer API has no server requester argument. Its ordinary
	// parameter map must not be treated as an authorization identity.
	legacy_rendered := markup.Get_render(
		db,
		"PublicDoc",
		"[include(RestrictedDoc)]",
		"api_view",
		parameter_data,
	)["data"]
	if strings.Contains(legacy_rendered, secret_marker) {
		t.Fatalf("legacy renderer treated parameter ip as trusted identity: %s", legacy_rendered)
	}
}
