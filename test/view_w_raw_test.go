package test

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

// Setup_raw_test_db builds a schema with an admin user, a view-only regular
// user, and documents covering normal, restricted, revision, hidden-revision
// and HTML-escape scenarios.
func Setup_raw_test_db(t *testing.T) *sql.DB {
	t.Helper()

	db_dir := filepath.Join(t.TempDir(), "db")
	os.MkdirAll(db_dir, 0755)
	os.Setenv("NAMU_DB", db_dir)
	os.Setenv("NAMU_DB_TYPE", "sqlite")
	defer os.Unsetenv("NAMU_DB")
	defer os.Unsetenv("NAMU_DB_TYPE")

	tool.DB_init()
	db := tool.DB_connect()

	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed %q: %v", q, err)
		}
	}

	// owner/admin grant full access; regular_group has view + history_view but
	// no hidel (cannot unhide) and no admin.
	exec("insert into alist (name, acl) values ('owner','owner')")
	exec("insert into alist (name, acl) values ('admin','admin')")
	exec("insert into alist (name, acl) values ('regular_group','view')")
	exec("insert into alist (name, acl) values ('regular_group','history_view')")

	exec("insert into user_set (name, id, data) values ('acl','admin_ip','owner')")
	exec("insert into user_set (name, id, data) values ('acl','regular_ip','regular_group')")

	// Public document: viewable by anyone with site_view.
	exec("insert into data (title, data) values (?, ?)", "PublicDoc", "Public document content.")

	// Restricted document: viewable only by the owner group.
	exec("insert into data (title, data) values (?, ?)", "RestrictedDoc", "Restricted document content.")
	exec("insert into acl (title, data, type) values ('RestrictedDoc','owner','view')")

	// Document containing HTML special characters (escape check).
	exec("insert into data (title, data) values (?, ?)", "HtmlDoc", "Test<b>Bold</b> & <i>Italic</i>")

	// Revision document: visible history revision.
	exec("insert into data (title, data) values (?, ?)", "RevDoc", "Rev current content.")
	exec("insert into history (id, title, data, date, ip, send, leng, hide, type) values (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"1", "RevDoc", "Revision one content", "2026-01-01 00:00:00", "seed_ip", "edit", "+1", "", "")

	// Hidden revision document: revision is hidden (no hidel => require auth).
	exec("insert into data (title, data) values (?, ?)", "HiddenDoc", "Hidden current content.")
	exec("insert into history (id, title, data, date, ip, send, leng, hide, type) values (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"1", "HiddenDoc", "Hidden content.", "2026-01-01 00:00:00", "seed_ip", "edit", "+1", "O", "")

	return db
}

func Test_view_w_raw(t *testing.T) {
	db := Setup_raw_test_db(t)
	defer tool.DB_close(db)

	admin := tool.Config{IP: "admin_ip", Session: &emptySession{}}
	regular := tool.Config{IP: "regular_ip", Session: &emptySession{}}

	// not_found_error is the message a missing resource must show;
	// email_error is the wrong message the pre-fix code could show.
	notFoundMsg := tool.Get_language(db, "not_found_error", true)
	authMsg := tool.Get_language(db, "authority_error", true)
	emailErrorMsg := tool.Get_language(db, "email_error", true)

	// 1. 정상 문서 RAW: admin -> 정상 RAW 화면(내용 표시).
	out := route.View_w_raw(admin, "PublicDoc", "", "")
	if !strings.Contains(out, "Public document content.") {
		t.Fatalf("normal doc: missing content, out: %q", out)
	}
	if !strings.Contains(out, "opennamu_preview_area") {
		t.Fatalf("normal doc: raw screen missing (rendered as error page?)")
	}

	// 2. 정상 리비전 RAW: admin, rev=1 -> 정상 RAW 화면(리비전 내용 표시).
	out = route.View_w_raw(admin, "RevDoc", "1", "")
	if !strings.Contains(out, "Revision one content") {
		t.Fatalf("revision: missing content, out: %q", out)
	}
	if !strings.Contains(out, "opennamu_preview_area") {
		t.Fatalf("revision: raw screen missing (rendered as error page?)")
	}

	// 3. 존재하지 않는 문서: admin -> not_found_error 안내(RAW 화면 아님, 이메일 오류 메시지 아님).
	out = route.View_w_raw(admin, "NoSuchDoc", "", "")
	if !strings.Contains(out, notFoundMsg) {
		t.Fatalf("nonexistent doc: missing not-found message, out: %q", out)
	}
	if strings.Contains(out, "opennamu_preview_area") {
		t.Fatalf("nonexistent doc: still rendered raw screen")
	}
	if strings.Contains(out, emailErrorMsg) {
		t.Fatalf("nonexistent doc: showed email error message (original bug), out: %q", out)
	}

	// 4. 존재하지 않는 리비전: admin, rev=99999 -> not_found_error 안내(RAW 화면 아님, 이메일 오류 메시지 아님).
	out = route.View_w_raw(admin, "PublicDoc", "99999", "")
	if !strings.Contains(out, notFoundMsg) {
		t.Fatalf("nonexistent rev: missing not-found message, out: %q", out)
	}
	if strings.Contains(out, "opennamu_preview_area") {
		t.Fatalf("nonexistent rev: still rendered raw screen")
	}
	if strings.Contains(out, emailErrorMsg) {
		t.Fatalf("nonexistent rev: showed email error message (original bug), out: %q", out)
	}

	// 5. 조회 권한 거부: regular, owner만 볼 수 있는 문서 -> auth 안내(RAW 화면 아님).
	out = route.View_w_raw(regular, "RestrictedDoc", "", "")
	if !strings.Contains(out, authMsg) {
		t.Fatalf("restricted doc: missing auth message, out: %q", out)
	}
	if strings.Contains(out, "opennamu_preview_area") {
		t.Fatalf("restricted doc: still rendered raw screen")
	}

	// 6. 숨김 리비전 권한 제한: regular(no hidel), rev=1 -> auth 안내(RAW 화면 아님).
	out = route.View_w_raw(regular, "HiddenDoc", "1", "")
	if !strings.Contains(out, authMsg) {
		t.Fatalf("hidden rev: missing auth message, out: %q", out)
	}
	if strings.Contains(out, "opennamu_preview_area") {
		t.Fatalf("hidden rev: still rendered raw screen")
	}

	// 7. 정상 RAW HTML 특수문자 이스케이프 유지: admin -> 이스케이프된 내용 표시.
	htmlContent := "Test<b>Bold</b> & <i>Italic</i>"
	out = route.View_w_raw(admin, "HtmlDoc", "", "")
	escapedContent := tool.HTML_escape(htmlContent)
	if !strings.Contains(out, escapedContent) {
		t.Fatalf("html escape: missing escaped content, out: %q", out)
	}
	if strings.Contains(out, htmlContent) {
		t.Fatalf("html escape: raw content present (not escaped), out: %q", out)
	}

	// 8. /raw_acl/ 안내 화면 유지: admin, do_type=document_acl -> 안내 + RAW 화면.
	out = route.View_w_raw(admin, "PublicDoc", "", "document_acl")
	if !strings.Contains(out, authMsg) {
		t.Fatalf("acl path: missing authority notice")
	}
	if !strings.Contains(out, "Public document content.") {
		t.Fatalf("acl path: missing raw content")
	}
	if !strings.Contains(out, "opennamu_preview_area") {
		t.Fatalf("acl path: missing preview area")
	}
}
