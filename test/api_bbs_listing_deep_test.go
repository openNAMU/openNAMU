package test

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_Api_bbs_listing_pinned_pagination_and_activity_cache(t *testing.T) {
	db := apiSlowDB(t)
	defer tool.DB_close(db)

	const board_id = "71"
	for _, setting := range []struct {
		name string
		data string
	}{
		{name: "bbs_name", data: "listing test"},
		{name: "bbs_type", data: "comment"},
	} {
		if _, err := db.Exec(`insert into bbs_set (set_name, set_code, set_id, set_data) values (?, '', ?, ?)`, setting.name, board_id, setting.data); err != nil {
			t.Fatalf("seed board setting %s: %v", setting.name, err)
		}
	}

	for post_code := 1; post_code <= 51; post_code++ {
		Seed_bbs_listing_post(db, t, board_id, fmt.Sprint(post_code), fmt.Sprintf("2026-01-01 00:%02d:00", post_code), "", nil)
	}
	Seed_bbs_listing_post(db, t, board_id, "1001", "2026-01-01 00:59:00", "1", []string{"first", "second"})
	Seed_bbs_listing_post(db, t, board_id, "1002", "2026-01-01 00:58:00", "2", nil)

	config := tool.Config{IP: "admin_ip", Session: &emptySession{}}
	first_page := route.Api_bbs(config, board_id, "1", "activity")
	if first_page["response"] != "ok" {
		t.Fatalf("first page response = %v, want ok", first_page["response"])
	}
	first_data, ok := first_page["data"].([]map[string]string)
	if !ok || len(first_data) != 52 {
		t.Fatalf("first page data = %T(%v), want 52 posts", first_page["data"], first_page["data"])
	}
	if first_page["has_next"] != true {
		t.Fatalf("first page has_next = %v, want true", first_page["has_next"])
	}
	if first_data[0]["set_code"] != "1002" || first_data[1]["set_code"] != "1001" || first_data[0]["pinned"] != "1" {
		t.Fatalf("pinned posts = %q/%q with pinned=%q, want 1002/1001 and pinned=1", first_data[0]["set_code"], first_data[1]["set_code"], first_data[0]["pinned"])
	}
	want_first_codes := []string{"1002", "1001", "1001", "1002"}
	for post_code := 51; post_code >= 4; post_code-- {
		want_first_codes = append(want_first_codes, fmt.Sprint(post_code))
	}
	for index, want_code := range want_first_codes {
		if first_data[index]["set_code"] != want_code {
			t.Fatalf("first page post %d = %q, want %q (pinned posts are also present in the ordinary query)", index, first_data[index]["set_code"], want_code)
		}
	}
	if first_data[4]["activity_date"] != "2026-01-01 00:51:00" {
		t.Fatalf("first ordinary activity date = %q, want 2026-01-01 00:51:00", first_data[4]["activity_date"])
	}
	if first_data[4]["user_id"] != "admin_ip" || first_data[4]["user_id_render"] == "" {
		t.Fatalf("user fields not hydrated: %#v", first_data[4])
	}
	tags := strings.Split(first_data[1]["tags"], ", ")
	seen_tags := map[string]bool{}
	for _, tag := range tags {
		seen_tags[tag] = true
	}
	if len(tags) != 2 || !seen_tags["first"] || !seen_tags["second"] {
		t.Fatalf("pinned post tags = %q, want both tags", first_data[1]["tags"])
	}

	var cached_activity_count int
	if err := db.QueryRow(`select count(*) from bbs_data where set_name = 'last_activity' and set_id = ?`, board_id).Scan(&cached_activity_count); err != nil {
		t.Fatalf("count activity cache rows: %v", err)
	}
	if cached_activity_count != 50 {
		t.Fatalf("activity cache rows after first page = %d, want 50 distinct records selected before the ordinary-list overflow row", cached_activity_count)
	}

	second_page := route.Api_bbs(config, board_id, "2", "activity")
	second_data, ok := second_page["data"].([]map[string]string)
	if !ok || len(second_data) != 5 || second_page["has_next"] != false {
		t.Fatalf("second page = %T(%v), has_next=%v; want 2 pinned + 3 ordinary and false", second_page["data"], second_page["data"], second_page["has_next"])
	}
	if second_data[0]["set_code"] != "1002" || second_data[1]["set_code"] != "1001" || second_data[2]["set_code"] != "3" || second_data[3]["set_code"] != "2" || second_data[4]["set_code"] != "1" {
		t.Fatalf("second page order = %q, %q, %q, %q, %q; want pinned 1002/1001 then ordinary posts 3/2/1", second_data[0]["set_code"], second_data[1]["set_code"], second_data[2]["set_code"], second_data[3]["set_code"], second_data[4]["set_code"])
	}
	if err := db.QueryRow(`select count(*) from bbs_data where set_name = 'last_activity' and set_id = ?`, board_id).Scan(&cached_activity_count); err != nil {
		t.Fatalf("recount activity cache rows: %v", err)
	}
	if cached_activity_count != 53 {
		t.Fatalf("activity cache rows after second page = %d, want 53", cached_activity_count)
	}
}

func Seed_bbs_listing_post(db *sql.DB, t *testing.T, board_id string, post_code string, date string, pinned string, tags []string) {
	t.Helper()
	for _, value := range []struct {
		name string
		data string
	}{
		{name: "title", data: "post " + post_code},
		{name: "date", data: date},
		{name: "user_id", data: "admin_ip"},
	} {
		if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, value.name, post_code, board_id, value.data); err != nil {
			t.Fatalf("seed %s for post %s: %v", value.name, post_code, err)
		}
	}
	if pinned != "" {
		if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('pinned', ?, ?, ?)`, post_code, board_id, pinned); err != nil {
			t.Fatalf("seed pinned post %s: %v", post_code, err)
		}
	}
	for _, tag := range tags {
		if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('tag', ?, ?, ?)`, post_code, board_id, tag); err != nil {
			t.Fatalf("seed tag %q for post %s: %v", tag, post_code, err)
		}
	}
}

func Test_Api_bbs_listing_closes_cursors_after_metadata_scan_panic(t *testing.T) {
	db := apiSlowDB(t)
	defer tool.DB_close(db)

	bbs_comment_seed_board(db, t, "73", "1", "scan failure", "admin_ip")
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('scan_failure', '1', '73', NULL)`); err != nil {
		t.Fatalf("seed null metadata: %v", err)
	}
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('metadata_after_failure', '1', '73', 'still open')`); err != nil {
		t.Fatalf("seed metadata after scan failure: %v", err)
	}

	var recovered any
	func() {
		defer func() { recovered = recover() }()
		route.Api_bbs(tool.Config{IP: "admin_ip", Session: &emptySession{}}, "73", "1", "date")
	}()
	if recovered == nil {
		t.Fatal("Api_bbs did not panic after scanning NULL metadata into a string")
	}

	var busy int
	var log_frames int
	var checkpointed_frames int
	if err := db.QueryRow(`pragma wal_checkpoint(TRUNCATE)`).Scan(&busy, &log_frames, &checkpointed_frames); err != nil {
		t.Fatalf("checkpoint after listing panic: %v", err)
	}
	if busy != 0 {
		t.Fatalf("WAL checkpoint remained busy after panic (log=%d checkpointed=%d); listing or metadata cursor leaked", log_frames, checkpointed_frames)
	}
}

func Test_sql_db_close_does_not_close_active_rows(t *testing.T) {
	db := audit_setup(t, "sqlite")
	rows, err := db.Query(`select 1 union all select 2`)
	if err != nil {
		t.Fatalf("open rows: %v", err)
	}

	tool.DB_close(db)
	defer rows.Close()
	for expected := 1; expected <= 2; expected++ {
		if !rows.Next() {
			t.Fatalf("active rows stopped at value %d after DB_close; err=%v", expected, rows.Err())
		}
		var actual int
		if err := rows.Scan(&actual); err != nil {
			t.Fatalf("scan active row after DB_close: %v", err)
		}
		if actual != expected {
			t.Fatalf("active row = %d, want %d", actual, expected)
		}
	}
}
