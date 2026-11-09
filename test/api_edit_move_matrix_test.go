package test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

type move_option_snapshot map[string][]string

func Test_Move_document_options_mode_matrix_against_existing_contract(t *testing.T) {
	db := apiSlowDB(t)
	defer tool.DB_close(db)

	move_options := []string{"none", "normal", "reverse", "merge"}
	topic_options := []string{"none", "normal", "reverse", "merge"}
	data_set_options := []string{"none", "normal", "reverse"}
	actors := []struct {
		name string
		ip   string
	}{
		{name: "owner", ip: "admin_ip"},
		{name: "non_owner", ip: "192.0.2.88"},
	}

	case_number := 0
	for _, actor := range actors {
		for _, move_option := range move_options {
			for _, topic_option := range topic_options {
				for _, data_set_option := range data_set_options {
					case_number++
					source_name := fmt.Sprintf("Phase9_move_%03d_source", case_number)
					target_name := fmt.Sprintf("Phase9_move_%03d_target", case_number)
					t.Run(fmt.Sprintf("%s/move=%s/topic=%s/settings=%s", actor.name, move_option, topic_option, data_set_option), func(t *testing.T) {
						Seed_move_option_document(db, t, source_name, "source body", "1")
						if !tool.Check_acl(db, source_name, "", "document_move", actor.ip) || !tool.Check_acl(db, target_name, "", "document_move", actor.ip) {
							t.Fatalf("fixture does not grant base document-move ACL to %s", actor.ip)
						}

						owner_auth := tool.Check_permission(db, "document_move_manage", actor.ip)
						requires_owner := move_option == "merge" || topic_option == "merge" || data_set_option == "reverse" ||
							(move_option == "normal" && data_set_option != "normal") ||
							(move_option != "normal" && data_set_option != "none")
						want_error := ""
						if !owner_auth && requires_owner {
							want_error = "auth"
						}

						before := Snapshot_move_option_state(t, db, source_name, target_name)
						got_error := route.Move_document_options(
							tool.Config{IP: actor.ip}, db, source_name, target_name, "matrix",
							move_option, topic_option, data_set_option,
						)
						if got_error != want_error {
							t.Fatalf("Move_document_options() error = %q, want %q", got_error, want_error)
						}
						if want_error != "" {
							if after := Snapshot_move_option_state(t, db, source_name, target_name); !reflect.DeepEqual(after, before) {
								t.Fatalf("rejected mode combination mutated state:\nbefore=%v\nafter=%v", before, after)
							}
							return
						}

						Assert_move_option_absent_target_result(t, db, source_name, target_name, move_option, topic_option, data_set_option)
					})
				}
			}
		}
	}
}

func Test_Move_document_options_destination_and_history_states(t *testing.T) {
	db := apiSlowDB(t)
	defer tool.DB_close(db)
	config := tool.Config{IP: "admin_ip"}

	destination_case_number := 0
	for _, test_case := range []struct {
		name            string
		target_kind     string
		move_option     string
		topic_option    string
		data_set_option string
		want_error      string
	}{
		{name: "normal rejects existing document", target_kind: "full", move_option: "normal", topic_option: "none", data_set_option: "none", want_error: "document already exist"},
		{name: "normal rejects history-only destination", target_kind: "history", move_option: "normal", topic_option: "none", data_set_option: "none", want_error: "document history exist"},
		{name: "normal rejects topic-only destination", target_kind: "topic", move_option: "normal", topic_option: "none", data_set_option: "none", want_error: "document already exist"},
		{name: "topic normal rejects occupied topic", target_kind: "topic", move_option: "none", topic_option: "normal", data_set_option: "none", want_error: "move error"},
		{name: "non-owner rejects destination settings collision", target_kind: "settings", move_option: "normal", topic_option: "none", data_set_option: "normal", want_error: "move error"},
		{name: "non-owner privilege gate precedes settings collision", target_kind: "settings", move_option: "normal", topic_option: "none", data_set_option: "none", want_error: "auth"},
		{name: "merge rejects malformed source history id", target_kind: "full", move_option: "merge", topic_option: "none", data_set_option: "none", want_error: "move error"},
	} {
		t.Run(test_case.name, func(t *testing.T) {
			destination_case_number++
			source_name := fmt.Sprintf("Phase9_destination_%02d_source", destination_case_number)
			target_name := fmt.Sprintf("Phase9_destination_%02d_target", destination_case_number)
			Seed_move_option_document(db, t, source_name, "source body", "1")
			Seed_move_option_destination(db, t, target_name, test_case.target_kind)
			if test_case.name == "merge rejects malformed source history id" {
				if _, err := db.Exec(`insert into history (id, title, data, date, ip, send, leng, hide, type) values ('bad', ?, 'invalid id', '2026-01-02 00:00:00', 'admin_ip', '', '0', '', 'edit')`, source_name); err != nil {
					t.Fatalf("seed malformed source history %s: %v", source_name, err)
				}
			}

			actor_ip := config.IP
			if test_case.name == "non-owner rejects destination settings collision" || test_case.name == "non-owner privilege gate precedes settings collision" {
				actor_ip = "192.0.2.88"
			}
			before := Snapshot_move_option_state(t, db, source_name, target_name)
			got_error := route.Move_document_options(
				tool.Config{IP: actor_ip}, db, source_name, target_name, "destination case",
				test_case.move_option, test_case.topic_option, test_case.data_set_option,
			)
			if got_error != test_case.want_error {
				t.Fatalf("Move_document_options() error = %q, want %q", got_error, test_case.want_error)
			}
			if after := Snapshot_move_option_state(t, db, source_name, target_name); !reflect.DeepEqual(after, before) {
				t.Fatalf("rejected destination case mutated state:\nbefore=%v\nafter=%v", before, after)
			}
		})
	}

	t.Run("reverse swaps documents topics and settings", func(t *testing.T) {
		source_name := "Phase9_reverse_source"
		target_name := "Phase9_reverse_target"
		Seed_move_option_document(db, t, source_name, "source body", "1")
		Seed_move_option_document(db, t, target_name, "target body", "7")

		if result := route.Move_document_options(config, db, source_name, target_name, "swap", "reverse", "reverse", "reverse"); result != "" {
			t.Fatalf("reverse move error = %q, want success", result)
		}
		Assert_move_document_body(t, db, source_name, "target body")
		Assert_move_document_body(t, db, target_name, "source body")
		Assert_move_history_count(t, db, source_name, 2)
		Assert_move_history_count(t, db, target_name, 2)
		Assert_move_setting(t, db, source_name, "matrix_setting", "target setting")
		Assert_move_setting(t, db, target_name, "matrix_setting", "source setting")
		Assert_move_reference_count(t, db, "-1", source_name, 2)
		Assert_move_reference_count(t, db, "-1", target_name, 2)
	})

	t.Run("merge combines document and history-only destination", func(t *testing.T) {
		source_name := "Phase9_merge_source"
		target_name := "Phase9_merge_target"
		Seed_move_option_document(db, t, source_name, "source body", "1")
		Seed_move_option_destination(db, t, target_name, "history")

		if result := route.Move_document_options(config, db, source_name, target_name, "combine", "merge", "none", "none"); result != "" {
			t.Fatalf("merge into history-only target error = %q, want success", result)
		}
		Assert_move_document_body(t, db, source_name, "")
		Assert_move_document_body(t, db, target_name, "source body")
		Assert_move_history_count(t, db, source_name, 0)
		Assert_move_history_count(t, db, target_name, 3)
		var event_count int
		if err := db.QueryRow(`select count(*) from history where title = ? and type = 'move'`, target_name).Scan(&event_count); err != nil {
			t.Fatalf("count merge event: %v", err)
		}
		if event_count != 1 {
			t.Fatalf("merge event count = %d, want 1", event_count)
		}
		var merged_history_id string
		if err := db.QueryRow(`select id from history where title = ? and data = 'revision:source body'`, target_name).Scan(&merged_history_id); err != nil {
			t.Fatalf("read merged source revision ID: %v", err)
		}
		if merged_history_id != "8" {
			t.Fatalf("merged source revision ID = %q, want 8", merged_history_id)
		}
	})
}

func Test_Move_document_options_validation_and_rollback(t *testing.T) {
	db := apiSlowDB(t)
	defer tool.DB_close(db)
	config := tool.Config{IP: "admin_ip"}

	for _, test_case := range []struct {
		name            string
		move_option     string
		topic_option    string
		data_set_option string
		want_error      string
	}{
		{name: "invalid document mode", move_option: "bad", topic_option: "none", data_set_option: "none", want_error: "move error"},
		{name: "invalid topic mode", move_option: "none", topic_option: "bad", data_set_option: "none", want_error: "move error"},
		{name: "invalid settings mode", move_option: "none", topic_option: "none", data_set_option: "bad", want_error: "move error"},
	} {
		t.Run(test_case.name, func(t *testing.T) {
			source_name := "Phase9_invalid_" + test_case.name + "_source"
			target_name := "Phase9_invalid_" + test_case.name + "_target"
			Seed_move_option_document(db, t, source_name, "source body", "1")
			before := Snapshot_move_option_state(t, db, source_name, target_name)
			got_error := route.Move_document_options(config, db, source_name, target_name, "", test_case.move_option, test_case.topic_option, test_case.data_set_option)
			if got_error != test_case.want_error {
				t.Fatalf("Move_document_options() error = %q, want %q", got_error, test_case.want_error)
			}
			if after := Snapshot_move_option_state(t, db, source_name, target_name); !reflect.DeepEqual(after, before) {
				t.Fatalf("invalid mode mutated state:\nbefore=%v\nafter=%v", before, after)
			}
		})
	}

	t.Run("invalid and missing inputs keep legacy errors", func(t *testing.T) {
		for case_number, test_case := range []struct {
			name       string
			kind       string
			actor      string
			want_error string
		}{
			{name: "equal titles", kind: "equal", actor: "admin_ip", want_error: "invalid document"},
			{name: "newline target", kind: "newline", actor: "admin_ip", want_error: "invalid document"},
			{name: "missing source", kind: "missing", actor: "admin_ip", want_error: "not exist"},
			{name: "authorization denied", kind: "denied", actor: tool.Invalid_remote_address, want_error: "auth"},
			{name: "oversized target", kind: "long", actor: "admin_ip", want_error: "invalid document"},
		} {
			t.Run(test_case.name, func(t *testing.T) {
				source_name := fmt.Sprintf("Phase9_input_%02d_source", case_number)
				target_name := fmt.Sprintf("Phase9_input_%02d_target", case_number)
				if test_case.kind == "equal" {
					target_name = source_name
				}
				if test_case.kind == "newline" {
					target_name += "\nbad"
				}
				if test_case.kind == "long" {
					target_name = strings.Repeat("x", 257)
				}
				if test_case.kind == "newline" || test_case.kind == "denied" {
					Seed_move_option_document(db, t, source_name, "source body", "1")
				}
				before := Snapshot_move_option_state(t, db, source_name, target_name)
				got_error := route.Move_document_options(tool.Config{IP: test_case.actor}, db, source_name, target_name, "", "normal", "none", "normal")
				if got_error != test_case.want_error {
					t.Fatalf("Move_document_options() error = %q, want %q", got_error, test_case.want_error)
				}
				if after := Snapshot_move_option_state(t, db, source_name, target_name); !reflect.DeepEqual(after, before) {
					t.Fatalf("rejected input mutated state:\nbefore=%v\nafter=%v", before, after)
				}
			})
		}
	})

	t.Run("combined reverse transaction rolls back after history failure", func(t *testing.T) {
		source_name := "Phase9_rollback_source"
		target_name := "Phase9_rollback_target"
		Seed_move_option_document(db, t, source_name, "source body", "1")
		Seed_move_option_document(db, t, target_name, "target body", "7")
		if _, err := db.Exec(`create trigger phase9_fail_reverse_move before insert on history when new.title = 'Phase9_rollback_source' and new.ip = 'admin_ip' begin select raise(abort, 'forced reverse history failure'); end`); err != nil {
			t.Fatalf("create reverse-move trigger: %v", err)
		}
		before := Snapshot_move_option_state(t, db, source_name, target_name)

		var recovered any
		func() {
			defer func() { recovered = recover() }()
			route.Move_document_options(config, db, source_name, target_name, "swap", "reverse", "reverse", "reverse")
		}()
		if recovered == nil {
			t.Fatal("reverse move did not panic after forced history failure")
		}
		if after := Snapshot_move_option_state(t, db, source_name, target_name); !reflect.DeepEqual(after, before) {
			t.Fatalf("failed reverse transaction left partial changes:\nbefore=%v\nafter=%v", before, after)
		}
	})
}

func Seed_move_option_document(db *sql.DB, t *testing.T, title string, body string, history_id string) {
	t.Helper()
	if _, err := db.Exec(`insert into data (title, data, type) values (?, ?, '')`, title, body); err != nil {
		t.Fatalf("seed document %s: %v", title, err)
	}
	if history_id != "" {
		if _, err := db.Exec(`insert into history (id, title, data, date, ip, send, leng, hide, type) values (?, ?, ?, ?, ?, ?, ?, '', 'edit')`, history_id, title, "revision:"+body, "2026-01-01 00:00:00", "admin_ip", "original", "0"); err != nil {
			t.Fatalf("seed history %s: %v", title, err)
		}
		if _, err := db.Exec(`insert into rc (id, title, date, type) values (?, ?, ?, 'edit')`, history_id, title, "2026-01-01 00:00:00"); err != nil {
			t.Fatalf("seed recent change %s: %v", title, err)
		}
	}
	if _, err := db.Exec(`insert into back (title, link, type, data) values (?, ?, 'matrix', 'backlink')`, title, title); err != nil {
		t.Fatalf("seed backlinks %s: %v", title, err)
	}
	setting_value := strings.TrimSuffix(body, " body") + " setting"
	if _, err := db.Exec(`insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', 'matrix_setting', ?)`, title, setting_value); err != nil {
		t.Fatalf("seed document setting %s: %v", title, err)
	}
	if _, err := db.Exec(`insert into acl (title, data, type) values (?, ?, 'matrix_acl')`, title, setting_value); err != nil {
		t.Fatalf("seed document ACL %s: %v", title, err)
	}
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('title', ?, '0', ?)`, "catalog_"+title, title); err != nil {
		t.Fatalf("seed document title reference %s: %v", title, err)
	}
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('document', ?, '-1', ?), ('tag', ?, '-1', ?)`, "topic_"+title, title, "topic_"+title, title); err != nil {
		t.Fatalf("seed topic references %s: %v", title, err)
	}
}

func Seed_move_option_destination(db *sql.DB, t *testing.T, title string, target_kind string) {
	t.Helper()
	switch target_kind {
	case "none":
		return
	case "full":
		Seed_move_option_document(db, t, title, "target body", "7")
	case "history":
		if _, err := db.Exec(`insert into history (id, title, data, date, ip, send, leng, hide, type) values ('7', ?, 'history only', '2026-01-02 00:00:00', 'admin_ip', '', '0', '', 'edit')`, title); err != nil {
			t.Fatalf("seed history-only destination %s: %v", title, err)
		}
		if _, err := db.Exec(`insert into rc (id, title, date, type) values ('7', ?, '2026-01-02 00:00:00', 'edit')`, title); err != nil {
			t.Fatalf("seed history-only recent change %s: %v", title, err)
		}
	case "topic":
		if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('document', ?, '-1', ?), ('tag', ?, '-1', ?)`, "topic_"+title, title, "topic_"+title, title); err != nil {
			t.Fatalf("seed topic-only destination %s: %v", title, err)
		}
	case "settings":
		if _, err := db.Exec(`insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', 'matrix_setting', 'target setting')`, title); err != nil {
			t.Fatalf("seed destination setting %s: %v", title, err)
		}
		if _, err := db.Exec(`insert into acl (title, data, type) values (?, 'target setting', 'matrix_acl')`, title); err != nil {
			t.Fatalf("seed destination ACL %s: %v", title, err)
		}
	default:
		t.Fatalf("unknown destination kind %q", target_kind)
	}
}

func Snapshot_move_option_state(t *testing.T, db *sql.DB, source_name string, target_name string) move_option_snapshot {
	t.Helper()
	queries := map[string]string{
		"data":     `select title, data, type from data where title in (?, ?) order by title, data, type`,
		"history":  `select id, title, data, date, ip, send, leng, hide, type from history where title in (?, ?) order by title, id + 0, id`,
		"rc":       `select id, title, date, type from rc where title in (?, ?) order by title, id + 0, id`,
		"back":     `select title, link, type, data from back where title in (?, ?) or link in (?, ?) order by title, link, type, data`,
		"data_set": `select doc_name, doc_rev, set_name, set_data from data_set where doc_name in (?, ?) order by doc_name, doc_rev, set_name, set_data`,
		"acl":      `select title, data, type from acl where title in (?, ?) order by title, data, type`,
		"bbs_data": `select set_name, set_code, set_id, set_data from bbs_data where set_id in ('0', '-1') and set_data in (?, ?) order by set_id, set_name, set_code, set_data`,
		"other":    `select name, data, coverage from other where name = 'count_all_title' order by name, coverage`,
	}
	table_names := make([]string, 0, len(queries))
	for table_name := range queries {
		table_names = append(table_names, table_name)
	}
	sort.Strings(table_names)

	snapshot := move_option_snapshot{}
	for _, table_name := range table_names {
		query := queries[table_name]
		arguments := []any{source_name, target_name}
		if table_name == "back" {
			arguments = []any{source_name, target_name, source_name, target_name}
		}
		rows, err := db.Query(query, arguments...)
		if err != nil {
			t.Fatalf("snapshot %s: %v", table_name, err)
		}
		columns, err := rows.Columns()
		if err != nil {
			rows.Close()
			t.Fatalf("read %s columns: %v", table_name, err)
		}
		for rows.Next() {
			values := make([]sql.NullString, len(columns))
			scan_values := make([]any, len(columns))
			for index := range values {
				scan_values[index] = &values[index]
			}
			if err := rows.Scan(scan_values...); err != nil {
				rows.Close()
				t.Fatalf("scan %s snapshot: %v", table_name, err)
			}
			row := make([]string, len(values))
			for index, value := range values {
				if value.Valid {
					row[index] = value.String
				} else {
					row[index] = "<NULL>"
				}
			}
			encoded_row, err := json.Marshal(row)
			if err != nil {
				rows.Close()
				t.Fatalf("encode %s snapshot: %v", table_name, err)
			}
			snapshot[table_name] = append(snapshot[table_name], string(encoded_row))
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			t.Fatalf("read %s snapshot: %v", table_name, err)
		}
		if err := rows.Close(); err != nil {
			t.Fatalf("close %s snapshot: %v", table_name, err)
		}
	}
	return snapshot
}

func Assert_move_option_absent_target_result(t *testing.T, db *sql.DB, source_name string, target_name string, move_option string, topic_option string, data_set_option string) {
	t.Helper()
	moves_document := move_option != "none"
	if moves_document {
		Assert_move_document_body(t, db, source_name, "")
		Assert_move_document_body(t, db, target_name, "source body")
		Assert_move_history_count(t, db, source_name, 0)
		Assert_move_history_count(t, db, target_name, 2)
		Assert_move_count(t, db, 1, `select count(*) from history where title = ? and type = 'move'`, target_name)
		Assert_move_count(t, db, 0, `select count(*) from rc where title = ?`, source_name)
		Assert_move_count(t, db, 3, `select count(*) from rc where title = ?`, target_name)
		Assert_move_reference_count(t, db, "", source_name, 0)
		Assert_move_reference_count(t, db, "", target_name, 1)
	} else {
		Assert_move_document_body(t, db, source_name, "source body")
		Assert_move_document_body(t, db, target_name, "")
		Assert_move_history_count(t, db, source_name, 1)
		Assert_move_history_count(t, db, target_name, 0)
		Assert_move_count(t, db, 0, `select count(*) from history where title = ? and type = 'move'`, target_name)
		Assert_move_count(t, db, 1, `select count(*) from rc where title = ?`, source_name)
		Assert_move_count(t, db, 0, `select count(*) from rc where title = ?`, target_name)
		Assert_move_reference_count(t, db, "", source_name, 1)
		Assert_move_reference_count(t, db, "", target_name, 0)
	}

	if topic_option == "none" {
		Assert_move_reference_count(t, db, "-1", source_name, 2)
		Assert_move_reference_count(t, db, "-1", target_name, 0)
	} else {
		Assert_move_reference_count(t, db, "-1", source_name, 0)
		Assert_move_reference_count(t, db, "-1", target_name, 2)
	}

	if data_set_option == "none" {
		Assert_move_setting(t, db, source_name, "matrix_setting", "source setting")
		Assert_move_setting(t, db, target_name, "matrix_setting", "")
		Assert_move_acl_setting(t, db, source_name, "source setting")
		Assert_move_acl_setting(t, db, target_name, "")
		Assert_move_reference_count(t, db, "0", source_name, 1)
		Assert_move_reference_count(t, db, "0", target_name, 0)
	} else {
		Assert_move_setting(t, db, source_name, "matrix_setting", "")
		Assert_move_setting(t, db, target_name, "matrix_setting", "source setting")
		Assert_move_acl_setting(t, db, source_name, "")
		Assert_move_acl_setting(t, db, target_name, "source setting")
		if data_set_option == "normal" {
			Assert_move_reference_count(t, db, "0", source_name, 0)
			Assert_move_reference_count(t, db, "0", target_name, 1)
		} else {
			Assert_move_reference_count(t, db, "0", source_name, 1)
			Assert_move_reference_count(t, db, "0", target_name, 0)
		}
	}
}

func Assert_move_document_body(t *testing.T, db *sql.DB, title string, want string) {
	t.Helper()
	if want == "" {
		Assert_move_count(t, db, 0, `select count(*) from data where title = ?`, title)
		return
	}
	var got string
	if err := db.QueryRow(`select data from data where title = ?`, title).Scan(&got); err != nil {
		t.Fatalf("read document %s: %v", title, err)
	}
	if got != want {
		t.Fatalf("document %s body = %q, want %q", title, got, want)
	}
}

func Assert_move_history_count(t *testing.T, db *sql.DB, title string, want int) {
	t.Helper()
	Assert_move_count(t, db, want, `select count(*) from history where title = ?`, title)
}

func Assert_move_setting(t *testing.T, db *sql.DB, title string, setting_name string, want string) {
	t.Helper()
	var got string
	err := db.QueryRow(`select set_data from data_set where doc_name = ? and set_name = ?`, title, setting_name).Scan(&got)
	if want == "" {
		if err == nil {
			t.Fatalf("document setting %s/%s unexpectedly exists with value %q", title, setting_name, got)
		}
		if err != sql.ErrNoRows {
			t.Fatalf("read absent setting %s/%s: %v", title, setting_name, err)
		}
		return
	}
	if err != nil {
		t.Fatalf("read setting %s/%s: %v", title, setting_name, err)
	}
	if got != want {
		t.Fatalf("document setting %s/%s = %q, want %q", title, setting_name, got, want)
	}
}

func Assert_move_acl_setting(t *testing.T, db *sql.DB, title string, want string) {
	t.Helper()
	var got string
	err := db.QueryRow(`select data from acl where title = ? and type = 'matrix_acl'`, title).Scan(&got)
	if want == "" {
		if err == nil {
			t.Fatalf("document ACL %s unexpectedly exists with value %q", title, got)
		}
		if err != sql.ErrNoRows {
			t.Fatalf("read absent ACL %s: %v", title, err)
		}
		return
	}
	if err != nil {
		t.Fatalf("read ACL %s: %v", title, err)
	}
	if got != want {
		t.Fatalf("document ACL %s = %q, want %q", title, got, want)
	}
}

func Assert_move_reference_count(t *testing.T, db *sql.DB, set_id string, title string, want int) {
	t.Helper()
	var query string
	if set_id == "" {
		query = `select count(*) from back where title = ? or link = ?`
		Assert_move_count(t, db, want, query, title, title)
		return
	}
	query = `select count(*) from bbs_data where set_id = ? and set_name in ('title', 'document', 'tag') and set_data = ?`
	Assert_move_count(t, db, want, query, set_id, title)
}

func Assert_move_count(t *testing.T, db *sql.DB, want int, query string, arguments ...any) {
	t.Helper()
	var got int
	if err := db.QueryRow(query, arguments...).Scan(&got); err != nil {
		t.Fatalf("count query %q: %v", query, err)
	}
	if got != want {
		t.Fatalf("count query %q = %d, want %d", query, got, want)
	}
}
