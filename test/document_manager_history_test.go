package test

import (
	"database/sql"
	"reflect"
	"testing"

	"opennamu/route/tool"
)

func Test_Do_add_history_id_assignment_modes(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		t.Run(backend, func(t *testing.T) {
			db := audit_setup(t, backend)
			defer tool.DB_close(db)

			for _, mode := range []string{"add", "setting"} {
				document_name := "History_" + mode
				tool.Do_add_history(db, document_name, "first", "2026-11-01 00:00:00", "history_user", "", "0", mode, "")
				tool.Do_add_history(db, document_name, "second", "2026-11-01 00:00:01", "history_user", "", "0", mode, "")
				Assert_history_ids_and_modes(t, db, document_name, []string{"-1", "0"}, []string{mode, mode})
				Assert_recent_history_types(t, db, document_name, map[string]int{})
			}

			tool.Do_add_history(db, "History_edit", "first", "2026-11-01 00:00:02", "history_user", "", "0", "edit", "")
			tool.Do_add_history(db, "History_edit", "second", "2026-11-01 00:00:03", "history_user", "", "0", "edit", "")
			Assert_history_ids_and_modes(t, db, "History_edit", []string{"1", "2"}, []string{"r1", "edit"})
			Assert_recent_history_types(t, db, "History_edit", map[string]int{"normal": 2, "r1": 1, "edit": 1})
			var last_edit string
			var edit_length string
			if err := db.QueryRow(`select set_data from data_set where doc_name = 'History_edit' and set_name = 'last_edit'`).Scan(&last_edit); err != nil {
				t.Fatalf("read last_edit setting: %v", err)
			}
			if err := db.QueryRow(`select set_data from data_set where doc_name = 'History_edit' and set_name = 'length'`).Scan(&edit_length); err != nil {
				t.Fatalf("read length setting: %v", err)
			}
			if last_edit != "2026-11-01 00:00:03" || edit_length != "6" {
				t.Fatalf("edit settings = last_edit %q, length %q; want latest date and 6", last_edit, edit_length)
			}

			tool.Do_add_history(db, "user:History_user", "first", "2026-11-01 00:00:04", "history_user", "", "0", "edit", "")
			Assert_history_ids_and_modes(t, db, "user:History_user", []string{"1"}, []string{"user"})
			Assert_recent_history_types(t, db, "user:History_user", map[string]int{"user": 1})
		})
	}
}

func Assert_history_ids_and_modes(t *testing.T, db *sql.DB, document_name string, want_ids []string, want_modes []string) {
	t.Helper()
	rows, err := db.Query(`select id, type from history where title = ? order by id + 0 asc`, document_name)
	if err != nil {
		t.Fatalf("query %s history: %v", document_name, err)
	}
	defer rows.Close()

	ids := []string{}
	modes := []string{}
	for rows.Next() {
		var id string
		var mode string
		if err := rows.Scan(&id, &mode); err != nil {
			t.Fatalf("scan %s history: %v", document_name, err)
		}
		ids = append(ids, id)
		modes = append(modes, mode)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("read %s history: %v", document_name, err)
	}
	if !reflect.DeepEqual(ids, want_ids) || !reflect.DeepEqual(modes, want_modes) {
		t.Fatalf("%s history = ids %v, modes %v; want ids %v, modes %v", document_name, ids, modes, want_ids, want_modes)
	}
}

func Assert_recent_history_types(t *testing.T, db *sql.DB, document_name string, want map[string]int) {
	t.Helper()
	rows, err := db.Query(`select type, count(*) from rc where title = ? group by type`, document_name)
	if err != nil {
		t.Fatalf("query %s recent history: %v", document_name, err)
	}
	defer rows.Close()

	got := map[string]int{}
	for rows.Next() {
		var mode string
		var count int
		if err := rows.Scan(&mode, &count); err != nil {
			t.Fatalf("scan %s recent history: %v", document_name, err)
		}
		got[mode] = count
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("read %s recent history: %v", document_name, err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s recent history = %v, want %v", document_name, got, want)
	}
}
