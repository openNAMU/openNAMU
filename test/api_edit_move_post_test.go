package test

import (
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_Move_document_normal_preserves_document_and_backlinks(t *testing.T) {
	db := apiSlowDB(t)
	defer tool.DB_close(db)

	if _, err := db.Exec(`insert into data (title, data, type) values (?, ?, '')`, "MoveSource", "source body"); err != nil {
		t.Fatalf("seed source: %v", err)
	}
	if _, err := db.Exec(`insert into back (link, title, type, data) values (?, ?, '', '')`, "MoveSource", "BacklinkSource"); err != nil {
		t.Fatalf("seed backlink: %v", err)
	}

	config := tool.Config{IP: "admin_ip"}
	if result := route.Move_document_normal(config, db, "MoveSource", "MoveTarget", "rename"); result != "" {
		t.Fatalf("Move_document_normal() = %q, want success", result)
	}

	var source_count int
	if err := db.QueryRow(`select count(*) from data where title = 'MoveSource'`).Scan(&source_count); err != nil {
		t.Fatalf("count source: %v", err)
	}
	if source_count != 0 {
		t.Fatalf("source document remains after move")
	}

	var moved_body string
	if err := db.QueryRow(`select data from data where title = 'MoveTarget'`).Scan(&moved_body); err != nil {
		t.Fatalf("read moved document: %v", err)
	}
	if moved_body != "source body" {
		t.Fatalf("moved body = %q, want source body", moved_body)
	}

	var moved_backlink string
	if err := db.QueryRow(`select link from back where title = 'BacklinkSource'`).Scan(&moved_backlink); err != nil {
		t.Fatalf("read moved backlink: %v", err)
	}
	if moved_backlink != "MoveTarget" {
		t.Fatalf("backlink target = %q, want MoveTarget", moved_backlink)
	}

	var history_count int
	if err := db.QueryRow(`select count(*) from history where title = 'MoveTarget'`).Scan(&history_count); err != nil {
		t.Fatalf("count move history: %v", err)
	}
	if history_count != 1 {
		t.Fatalf("move history rows = %d, want 1", history_count)
	}

	var history_data string
	var history_ip string
	var history_send string
	if err := db.QueryRow(`select data, ip, send from history where title = 'MoveTarget'`).Scan(&history_data, &history_ip, &history_send); err != nil {
		t.Fatalf("read move history: %v", err)
	}
	if history_data != "source body" || history_ip != config.IP || history_send != "rename (<a>MoveSource</a> → <a>MoveTarget</a>)" {
		t.Fatalf("move history = data %q, ip %q, send %q; want source body, %q, and the legacy move label", history_data, history_ip, history_send, config.IP)
	}
}

func Test_Move_document_options_none_preserves_existing_documents(t *testing.T) {
	db := apiSlowDB(t)
	defer tool.DB_close(db)

	for _, document := range []struct {
		title string
		data  string
	}{
		{title: "NoneSource", data: "source body"},
		{title: "NoneTarget", data: "target body"},
	} {
		if _, err := db.Exec(`insert into data (title, data, type) values (?, ?, '')`, document.title, document.data); err != nil {
			t.Fatalf("seed %s: %v", document.title, err)
		}
	}

	config := tool.Config{IP: "admin_ip"}
	if result := route.Move_document_options(config, db, "NoneSource", "NoneTarget", "", "none", "none", "none"); result != "" {
		t.Fatalf("Move_document_options(none) = %q, want success", result)
	}
	for _, document := range []struct {
		title string
		data  string
	}{
		{title: "NoneSource", data: "source body"},
		{title: "NoneTarget", data: "target body"},
	} {
		var got string
		if err := db.QueryRow(`select data from data where title = ?`, document.title).Scan(&got); err != nil {
			t.Fatalf("read %s after no-op: %v", document.title, err)
		}
		if got != document.data {
			t.Errorf("%s after no-op = %q, want %q", document.title, got, document.data)
		}
	}
	var history_count int
	if err := db.QueryRow(`select count(*) from history where title in ('NoneSource', 'NoneTarget')`).Scan(&history_count); err != nil {
		t.Fatalf("count history after no-op: %v", err)
	}
	if history_count != 0 {
		t.Fatalf("no-op move history rows = %d, want 0", history_count)
	}
}

func Test_Move_document_normal_rolls_back_on_history_failure(t *testing.T) {
	db := apiSlowDB(t)
	defer tool.DB_close(db)

	if _, err := db.Exec(`insert into data (title, data, type) values (?, ?, '')`, "RollbackSource", "rollback body"); err != nil {
		t.Fatalf("seed source: %v", err)
	}
	if _, err := db.Exec(`create trigger fail_move_history before insert on history when new.title = 'RollbackTarget' and new.ip = 'admin_ip' begin select raise(abort, 'forced history failure'); end`); err != nil {
		t.Fatalf("create history failure trigger: %v", err)
	}

	config := tool.Config{IP: "admin_ip"}
	func() {
		defer func() {
			if recovered := recover(); recovered == nil {
				t.Fatal("Move_document_normal() did not panic after the forced history failure")
			}
		}()
		route.Move_document_normal(config, db, "RollbackSource", "RollbackTarget", "rename")
	}()

	var source_count int
	if err := db.QueryRow(`select count(*) from data where title = 'RollbackSource' and data = 'rollback body'`).Scan(&source_count); err != nil {
		t.Fatalf("check source after rollback: %v", err)
	}
	if source_count != 1 {
		t.Fatalf("source document rows after rollback = %d, want 1", source_count)
	}
	var target_count int
	if err := db.QueryRow(`select count(*) from data where title = 'RollbackTarget'`).Scan(&target_count); err != nil {
		t.Fatalf("check target after rollback: %v", err)
	}
	if target_count != 0 {
		t.Fatalf("target document rows after rollback = %d, want 0", target_count)
	}
}
