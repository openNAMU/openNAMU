package test

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func calendar_test_config(ip string) tool.Config {
	return tool.Config{IP: ip, Session: &emptySession{}}
}

func Test_calendar_api_crud_and_permissions(t *testing.T) {
	db := audit_setup(t, "sqlite")
	audit_seed_auth(db, t)
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', 'normal_user', 'user')`); err != nil {
		t.Fatalf("seed normal user: %v", err)
	}

	admin := calendar_test_config("admin_ip")
	user := calendar_test_config("normal_user")
	anon := calendar_test_config("192.0.2.10")

	if got := route.Api_calendar_add(anon, "공지", "설명"); got["response"] != "require auth" {
		t.Fatalf("anonymous add response = %#v, want require auth", got)
	}

	added := route.Api_calendar_add(admin, "공지", "설명")
	if added["response"] != "ok" || added["data"] != "1" {
		t.Fatalf("add calendar response = %#v, want ok with id 1", added)
	}

	listed := route.Api_calendar_list(anon)
	list, _ := listed["data"].([][]string)
	if listed["response"] != "ok" || len(list) != 1 || list[0][1] != "공지" || list[0][2] != "설명" {
		t.Fatalf("calendar list = %#v, want one 공지/설명 calendar", listed)
	}

	viewed := route.Api_calendar_view(anon, "1", "2024-02")
	calendar_data, _ := viewed["data"].(map[string]string)
	if viewed["response"] != "ok" || calendar_data["name"] != "공지" || calendar_data["month"] != "2024-02" {
		t.Fatalf("calendar view = %#v, want 2024-02 공지", viewed)
	}
	if calendar_data["prev_month"] != "2024-01" || calendar_data["next_month"] != "2024-03" {
		t.Fatalf("calendar navigation = %#v", calendar_data)
	}

	if got := route.Api_calendar_event_add(anon, "1", "2024-02-29", "윤년 행사", "메모"); got["response"] != "require auth" {
		t.Fatalf("anonymous event add response = %#v, want require auth", got)
	}
	if got := route.Api_calendar_event_add(admin, "1", "2024-02-30", "잘못된 날짜", ""); got["response"] != "error" || got["data"] != "invalid date" {
		t.Fatalf("invalid date response = %#v, want invalid date", got)
	}
	if got := route.Api_calendar_event_add(admin, "1", "2024-02-01", "", ""); got["response"] != "error" || got["data"] != "empty data" {
		t.Fatalf("empty title response = %#v, want empty data", got)
	}

	event_added := route.Api_calendar_event_add(admin, "1", "2024-02-29", "윤년 행사", "메모")
	if event_added["response"] != "ok" || event_added["data"] != "1" || event_added["month"] != "2024-02" {
		t.Fatalf("event add response = %#v, want ok code 1 in 2024-02", event_added)
	}

	viewed = route.Api_calendar_view(anon, "1", "2024-02")
	events, _ := viewed["events"].([][]string)
	if viewed["response"] != "ok" || len(events) != 1 {
		t.Fatalf("february events = %#v, want one event", viewed)
	}
	if events[0][1] != "2024-02-29" || events[0][2] != "윤년 행사" || events[0][3] != "메모" || events[0][4] != "admin_ip" {
		t.Fatalf("event row = %#v", events[0])
	}

	march := route.Api_calendar_view(anon, "1", "2024-03")
	if march_events, _ := march["events"].([][]string); len(march_events) != 0 {
		t.Fatalf("march events = %#v, want none", march_events)
	}

	event := route.Api_calendar_event(anon, "1", "1")
	event_data, _ := event["data"].(map[string]string)
	editable, _ := event["editable"].(bool)
	if event["response"] != "ok" || event_data["title"] != "윤년 행사" || editable {
		t.Fatalf("anonymous event view = %#v, want read-only event", event)
	}

	event = route.Api_calendar_event(admin, "1", "1")
	editable, _ = event["editable"].(bool)
	if !editable {
		t.Fatalf("owner event view editable = false, want true")
	}

	if got := route.Api_calendar_event_edit(user, "1", "1", "2024-02-28", "수정", ""); got["response"] != "require auth" {
		t.Fatalf("non-owner event edit response = %#v, want require auth", got)
	}

	updated := route.Api_calendar_event_edit(admin, "1", "1", "2024-02-28", "수정된 행사", "수정 메모")
	if updated["response"] != "ok" || updated["month"] != "2024-02" {
		t.Fatalf("event edit response = %#v, want ok", updated)
	}

	user_event := route.Api_calendar_event_add(user, "1", "2024-02-10", "사용자 일정", "")
	if user_event["response"] != "ok" || user_event["data"] != "2" {
		t.Fatalf("user event add response = %#v, want ok code 2", user_event)
	}
	if got := route.Api_calendar_event_edit(user, "1", "2", "2024-02-11", "사용자 수정", ""); got["response"] != "ok" {
		t.Fatalf("author event edit response = %#v, want ok", got)
	}
	if got := route.Api_calendar_event_delete(user, "1", "1"); got["response"] != "require auth" {
		t.Fatalf("non-owner event delete response = %#v, want require auth", got)
	}

	deleted := route.Api_calendar_event_delete(admin, "1", "1")
	if deleted["response"] != "ok" || deleted["month"] != "2024-02" {
		t.Fatalf("event delete response = %#v, want ok", deleted)
	}
	viewed = route.Api_calendar_view(anon, "1", "2024-02")
	events, _ = viewed["events"].([][]string)
	if len(events) != 1 || events[0][2] != "사용자 수정" {
		t.Fatalf("events after delete = %#v, want only the user event", events)
	}

	if got := route.Api_calendar_edit(user, "1", "x", ""); got["response"] != "require auth" {
		t.Fatalf("non-owner calendar edit response = %#v, want require auth", got)
	}
	if got := route.Api_calendar_edit(admin, "1", "새 이름", "새 설명"); got["response"] != "ok" {
		t.Fatalf("calendar edit response = %#v, want ok", got)
	}
	viewed = route.Api_calendar_view(anon, "1", "2024-02")
	calendar_data, _ = viewed["data"].(map[string]string)
	if calendar_data["name"] != "새 이름" || calendar_data["subject"] != "새 설명" {
		t.Fatalf("edited calendar data = %#v", calendar_data)
	}

	if got := route.Api_calendar_delete(user, "1"); got["response"] != "require auth" {
		t.Fatalf("non-owner calendar delete response = %#v, want require auth", got)
	}
	if got := route.Api_calendar_delete(admin, "1"); got["response"] != "ok" {
		t.Fatalf("calendar delete response = %#v, want ok", got)
	}
	if got := route.Api_calendar_view(anon, "1", "2024-02"); got["response"] != "not exist" {
		t.Fatalf("view deleted calendar = %#v, want not exist", got)
	}

	var remaining int
	if err := db.QueryRow(`select count(*) from calendar where set_id = '1'`).Scan(&remaining); err != nil {
		t.Fatalf("count deleted calendar rows: %v", err)
	}
	if remaining != 0 {
		t.Fatalf("deleted calendar rows = %d, want 0", remaining)
	}
}

func Test_calendar_api_month_and_isolation(t *testing.T) {
	db := audit_setup(t, "sqlite")
	audit_seed_auth(db, t)
	admin := calendar_test_config("admin_ip")
	anon := calendar_test_config("192.0.2.11")

	first := route.Api_calendar_add(admin, "첫 번째", "")
	second := route.Api_calendar_add(admin, "두 번째", "")
	if first["response"] != "ok" || second["response"] != "ok" {
		t.Fatalf("add calendars = %#v / %#v", first, second)
	}

	first_id, _ := first["data"].(string)
	second_id, _ := second["data"].(string)

	first_event := route.Api_calendar_event_add(admin, first_id, "2024-02-01", "첫 일정", "")
	second_event := route.Api_calendar_event_add(admin, second_id, "2024-02-01", "둘째 일정", "")
	december_event := route.Api_calendar_event_add(admin, first_id, "2024-12-31", "연말 일정", "")
	if first_event["response"] != "ok" || second_event["response"] != "ok" || december_event["response"] != "ok" {
		t.Fatalf("add events = %#v / %#v / %#v", first_event, second_event, december_event)
	}

	first_view := route.Api_calendar_view(anon, first_id, "2024-02")
	first_events, _ := first_view["events"].([][]string)
	if len(first_events) != 1 || first_events[0][2] != "첫 일정" {
		t.Fatalf("first calendar february events = %#v", first_events)
	}

	second_view := route.Api_calendar_view(anon, second_id, "2024-02")
	second_events, _ := second_view["events"].([][]string)
	if len(second_events) != 1 || second_events[0][2] != "둘째 일정" {
		t.Fatalf("second calendar february events = %#v", second_events)
	}

	january_view := route.Api_calendar_view(anon, first_id, "2025-01")
	if january_events, _ := january_view["events"].([][]string); len(january_events) != 0 {
		t.Fatalf("january events = %#v, want none", january_events)
	}

	december_view := route.Api_calendar_view(anon, first_id, "2024-12")
	december_events, _ := december_view["events"].([][]string)
	if len(december_events) != 1 || december_events[0][2] != "연말 일정" {
		t.Fatalf("december events = %#v", december_events)
	}

	december_code, _ := december_event["data"].(string)
	if got := route.Api_calendar_event(anon, second_id, december_code); got["response"] != "not exist" {
		t.Fatalf("cross-calendar event view = %#v, want not exist", got)
	}

	if got := route.Api_calendar_view(anon, first_id, "2024-13"); got["response"] != "error" || got["data"] != "invalid date" {
		t.Fatalf("invalid month view = %#v, want invalid date", got)
	}
	if got := route.Api_calendar_view(anon, "abc", ""); got["response"] != "not exist" {
		t.Fatalf("invalid calendar id view = %#v, want not exist", got)
	}
}

func Test_calendar_view_escapes_user_content(t *testing.T) {
	db := audit_setup(t, "sqlite")
	audit_seed_auth(db, t)
	admin := calendar_test_config("admin_ip")

	added := route.Api_calendar_add(admin, "공지", "")
	id, _ := added["data"].(string)
	event := route.Api_calendar_event_add(admin, id, "2024-02-01", `<script>alert(1)</script>`, `<b>memo</b>`)
	code, _ := event["data"].(string)
	if event["response"] != "ok" {
		t.Fatalf("add escaping event = %#v", event)
	}

	html := route.View_calendar(admin, id, "2024-02")
	if strings.Contains(html, `<script>alert(1)</script>`) {
		t.Fatalf("calendar month HTML contains raw script title")
	}
	if !strings.Contains(html, `&lt;script&gt;alert(1)&lt;/script&gt;`) {
		t.Fatalf("calendar month HTML does not contain escaped title")
	}
	if strings.Contains(html, `<b>memo</b>`) {
		t.Fatalf("calendar month HTML contains raw memo HTML")
	}

	event_html := route.View_calendar_event(admin, id, code)
	if strings.Contains(event_html, `<script>alert(1)</script>`) {
		t.Fatalf("event HTML contains raw script title")
	}
	if !strings.Contains(event_html, `&lt;script&gt;alert(1)&lt;/script&gt;`) {
		t.Fatalf("event HTML does not contain escaped title")
	}
	if strings.Contains(event_html, `<b>memo</b>`) {
		t.Fatalf("event HTML contains raw memo HTML")
	}
}

func Test_calendar_event_concurrent_add(t *testing.T) {
	db := audit_setup(t, "sqlite")
	audit_seed_auth(db, t)
	admin := calendar_test_config("admin_ip")

	added := route.Api_calendar_add(admin, "동시성", "")
	id, _ := added["data"].(string)

	const workers = 4
	codes := make([]string, workers)
	failures := make([]any, workers)
	var wait sync.WaitGroup
	for i := 0; i < workers; i++ {
		wait.Add(1)
		go func(index int) {
			defer wait.Done()
			result := route.Api_calendar_event_add(admin, id, "2024-02-01", "동시 일정", "")
			codes[index], _ = result["data"].(string)
			if result["response"] != "ok" {
				failures[index] = result
			}
		}(i)
	}
	wait.Wait()

	seen := map[string]bool{}
	for i, code := range codes {
		if failures[i] != nil {
			t.Fatalf("concurrent add %d = %#v", i, failures[i])
		}
		if code == "" || seen[code] {
			t.Fatalf("concurrent codes = %#v, want unique non-empty codes", codes)
		}
		seen[code] = true
	}

	viewed := route.Api_calendar_view(admin, id, "2024-02")
	events, _ := viewed["events"].([][]string)
	if len(events) != workers {
		t.Fatalf("concurrent event count = %d, want %d", len(events), workers)
	}
}

func Test_calendar_concurrent_add(t *testing.T) {
	db := audit_setup(t, "sqlite")
	audit_seed_auth(db, t)
	admin := calendar_test_config("admin_ip")

	const workers = 4
	ids := make([]string, workers)
	failures := make([]any, workers)
	var wait sync.WaitGroup
	for i := 0; i < workers; i++ {
		wait.Add(1)
		go func(index int) {
			defer wait.Done()
			result := route.Api_calendar_add(admin, fmt.Sprintf("동시 캘린더 %d", index), "")
			ids[index], _ = result["data"].(string)
			if result["response"] != "ok" {
				failures[index] = result
			}
		}(i)
	}
	wait.Wait()

	seen := map[string]bool{}
	for i, id := range ids {
		if failures[i] != nil {
			t.Fatalf("concurrent calendar add %d = %#v", i, failures[i])
		}
		if id == "" || seen[id] {
			t.Fatalf("concurrent calendar ids = %#v, want unique non-empty ids", ids)
		}
		seen[id] = true
	}

	listed := route.Api_calendar_list(admin)
	list, _ := listed["data"].([][]string)
	if len(list) != workers {
		t.Fatalf("concurrent calendar count = %d, want %d", len(list), workers)
	}
}
