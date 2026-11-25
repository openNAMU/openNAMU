package test

import (
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_calendar_permission_groups(t *testing.T) {
	db := audit_setup(t, "sqlite")
	audit_seed_auth(db, t)

	exec := func(query string, values ...any) {
		t.Helper()
		if _, err := db.Exec(query, values...); err != nil {
			t.Fatalf("seed %q: %v", query, err)
		}
	}

	// Custom groups: document grants site_view. The writer has only
	// calendar_create; the manager has only calendar_manage (which implies
	// calendar_create at runtime). The limited group has no calendar write
	// permission at all.
	exec(`insert into alist (name, acl) values ('calendar_writer', 'document')`)
	exec(`insert into alist (name, acl) values ('calendar_writer', 'calendar_create')`)
	exec(`insert into alist (name, acl) values ('calendar_manager', 'document')`)
	exec(`insert into alist (name, acl) values ('calendar_manager', 'calendar_manage')`)
	exec(`insert into alist (name, acl) values ('calendar_limited', 'document')`)
	for user, group := range map[string]string{
		"normal_user":  "user",
		"writer_user":  "calendar_writer",
		"manager_user": "calendar_manager",
		"limited_user": "calendar_limited",
	} {
		exec(`insert into user_set (name, id, data) values ('acl', ?, ?)`, user, group)
	}

	admin := calendar_test_config("admin_ip")
	normal := calendar_test_config("normal_user")
	writer := calendar_test_config("writer_user")
	manager := calendar_test_config("manager_user")
	limited := calendar_test_config("limited_user")

	if !tool.Check_permission(db, "calendar_create", normal.IP) || tool.Check_permission(db, "calendar_manage", normal.IP) {
		t.Fatalf("default user calendar permissions are wrong: create=%v manage=%v",
			tool.Check_permission(db, "calendar_create", normal.IP),
			tool.Check_permission(db, "calendar_manage", normal.IP))
	}
	if !tool.Check_permission(db, "calendar_create", admin.IP) || !tool.Check_permission(db, "calendar_manage", admin.IP) {
		t.Fatalf("admin calendar permissions are wrong")
	}
	if !tool.Check_permission(db, "calendar_create", writer.IP) || tool.Check_permission(db, "calendar_manage", writer.IP) {
		t.Fatalf("writer calendar permissions are wrong")
	}
	if !tool.Check_permission(db, "calendar_create", manager.IP) || !tool.Check_permission(db, "calendar_manage", manager.IP) {
		t.Fatalf("manager calendar permissions are wrong")
	}
	if tool.Check_permission(db, "calendar_create", limited.IP) || tool.Check_permission(db, "calendar_manage", limited.IP) {
		t.Fatalf("limited user unexpectedly has calendar write permission")
	}
	if !tool.Check_permission(db, "site_view", limited.IP) {
		t.Fatalf("limited user does not have site_view")
	}
	if tool.Check_permission(db, "treat_as_admin", manager.IP) {
		t.Fatalf("calendar_manage must not grant admin default features")
	}

	added := route.Api_calendar_add(admin, "관리 대상", "")
	id, _ := added["data"].(string)
	if added["response"] != "ok" {
		t.Fatalf("admin calendar add = %#v", added)
	}

	if got := route.Api_calendar_view(limited, id, ""); got["response"] != "ok" {
		t.Fatalf("limited user cannot read calendar: %#v", got)
	}
	if got := route.Api_calendar_add(limited, "불가", ""); got["response"] != "require auth" {
		t.Fatalf("limited calendar add = %#v, want require auth", got)
	}
	if got := route.Api_calendar_event_add(limited, id, "2024-01-01", "불가", ""); got["response"] != "require auth" {
		t.Fatalf("limited event add = %#v, want require auth", got)
	}

	admin_event := route.Api_calendar_event_add(admin, id, "2024-01-01", "관리자 일정", "")
	if admin_event["response"] != "ok" {
		t.Fatalf("admin event add = %#v, want ok", admin_event)
	}
	if got := route.Api_calendar_add(writer, "작성자 캘린더", ""); got["response"] != "ok" {
		t.Fatalf("writer calendar add = %#v, want ok", got)
	}
	writer_event := route.Api_calendar_event_add(writer, id, "2024-01-02", "작성자 일정", "")
	if writer_event["response"] != "ok" {
		t.Fatalf("writer event add = %#v, want ok", writer_event)
	}
	if got := route.Api_calendar_edit(writer, id, "불가", ""); got["response"] != "require auth" {
		t.Fatalf("writer calendar edit = %#v, want require auth", got)
	}
	if got := route.Api_calendar_event_edit(writer, id, "1", "2024-01-03", "불가", ""); got["response"] != "require auth" {
		t.Fatalf("writer editing another author's event = %#v, want require auth", got)
	}
	if got := route.Api_calendar_event_edit(writer, id, "2", "2024-01-03", "작성자 수정", ""); got["response"] != "ok" {
		t.Fatalf("writer editing own event = %#v, want ok", got)
	}

	if got := route.Api_calendar_edit(manager, id, "관리자 수정", ""); got["response"] != "ok" {
		t.Fatalf("manager calendar edit = %#v, want ok", got)
	}
	if got := route.Api_calendar_event_edit(manager, id, "1", "2024-01-02", "관리자 수정 일정", ""); got["response"] != "ok" {
		t.Fatalf("manager event edit = %#v, want ok", got)
	}
	if got := route.Api_calendar_event_delete(manager, id, "1"); got["response"] != "ok" {
		t.Fatalf("manager event delete = %#v, want ok", got)
	}
	if got := route.Api_calendar_delete(manager, id); got["response"] != "ok" {
		t.Fatalf("manager calendar delete = %#v, want ok", got)
	}
}
