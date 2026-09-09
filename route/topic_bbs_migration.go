package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func normalize_thread_bbs_prefix(db *sql.DB) error {
	rows, err := db.Query(
		tool.DB_change("select set_code from bbs_data where set_id = ? and set_name = 'document'"),
		thread_bbs_set_id,
	)
	if err != nil {
		return err
	}

	codes := []string{}
	for rows.Next() {
		code := ""
		if err := rows.Scan(&code); err != nil {
			rows.Close()
			return err
		}
		codes = append(codes, code)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}
	rows.Close()

	changed := false
	for _, code := range codes {
		prefix := ""
		prefix_exists := tool.QueryRow_DB(
			db,
			"select set_data from bbs_data where set_id = ? and set_code = ? and set_name = 'prefix' limit 1",
			[]any{&prefix},
			thread_bbs_set_id,
			code,
		)
		if prefix_exists && prefix != "" {
			continue
		}

		stop := ""
		tool.QueryRow_DB(
			db,
			"select set_data from bbs_data where set_id = ? and set_code = ? and set_name = 'topic_stop' limit 1",
			[]any{&stop},
			thread_bbs_set_id,
			code,
		)
		prefix = "열림"
		if stop == "O" {
			prefix = "닫힘"
		}

		if prefix_exists {
			tool.Exec_DB(
				db,
				"update bbs_data set set_data = ? where set_id = ? and set_code = ? and set_name = 'prefix'",
				prefix,
				thread_bbs_set_id,
				code,
			)
		} else {
			tool.Exec_DB(
				db,
				"insert into bbs_data (set_name, set_code, set_id, set_data) values ('prefix', ?, ?, ?)",
				code,
				thread_bbs_set_id,
				prefix,
			)
		}
		changed = true
	}

	result, err := db.Exec(
		tool.DB_change("delete from bbs_data where set_id = ? and set_name in ('topic_stop', 'comment_close')"),
		thread_bbs_set_id,
	)
	if err != nil {
		return err
	}
	deleted, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if deleted > 0 {
		changed = true
	}

	if changed {
		if err := tool.Search_bbs_index_mark_rebuild(); err != nil {
			return err
		}
	}
	return nil
}

func migrate_thread_bbs_document_tags(db *sql.DB) error {
	rows, err := db.Query(
		tool.DB_change("select set_code, set_data from bbs_data where set_id = ? and set_name = 'document' and set_data != ''"),
		thread_bbs_set_id,
	)
	if err != nil {
		return err
	}

	document_data := [][]string{}
	for rows.Next() {
		data := []string{"", ""}
		if err := rows.Scan(&data[0], &data[1]); err != nil {
			rows.Close()
			return err
		}
		document_data = append(document_data, data)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}
	rows.Close()

	changed := false
	for _, data := range document_data {
		var tag_count int
		if err := db.QueryRow(
			tool.DB_change("select count(*) from bbs_data where set_id = ? and set_code = ? and set_name = 'tag' and set_data = ?"),
			thread_bbs_set_id,
			data[0],
			data[1],
		).Scan(&tag_count); err != nil {
			return err
		}
		if tag_count > 0 {
			continue
		}

		if _, err := db.Exec(
			tool.DB_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('tag', ?, ?, ?)"),
			data[0],
			thread_bbs_set_id,
			data[1],
		); err != nil {
			return err
		}
		changed = true
	}

	if changed {
		return tool.Search_bbs_index_mark_rebuild()
	}
	return nil
}

func Migrate_topic_to_bbs(previous_version string) error {
	if previous_version == "" || previous_version >= "20260902" {
		return nil
	}

	db := tool.DB_connect()
	defer tool.DB_close(db)

	if previous_version < "20260901" {
		if err := emergency_migrate_topic_to_bbs(db); err != nil {
			return err
		}
	}
	if previous_version < "20260902" {
		return migrate_thread_bbs_document_tags(db)
	}
	return nil
}
