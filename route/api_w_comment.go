package route

import (
	"database/sql"
	"strconv"

	"opennamu/route/tool"
)

func Api_w_comment(config tool.Config, doc_name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	db_code := tool.Get_document_setting(db, doc_name, "document_comment_code", "")

	db_code_str := ""
	if len(db_code) >= 1 {
		db_code_str = db_code[0][0]
	}

	if db_code_str != "" {
		comment_title := ""
		if !tool.QueryRow_DB(
			db,
			"select set_data from bbs_data where set_name = 'title' and set_id = '0' and set_code = ? limit 1",
			[]any{&comment_title},
			db_code_str,
		) || comment_title != doc_name {
			tool.Exec_DB(
				db,
				"delete from data_set where doc_name = ? and set_name = 'document_comment_code'",
				doc_name,
			)
			db_code_str = ""
		}
	}

	if db_code_str == "" {
		if !tool.Check_permission(db, "bbs_comment", "Tool:System") {
			return map[string]any{"response": "require auth"}
		}

		if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
			last_code := ""
			tool.QueryRow_DB(
				tx,
				"select set_code from bbs_data where set_name = ? and set_id = ? order by set_code + 0 desc",
				[]any{&last_code},
				"title",
				"0",
			)
			set_code := strconv.Itoa(tool.Str_to_int(last_code) + 1)
			date_now := tool.Get_time()
			for _, value := range [][]string{
				{"title", doc_name},
				{"data", ""},
				{"date", date_now},
				{"last_activity", date_now},
				{"user_id", "Tool:System"},
				{"comment_count", "0"},
			} {
				if _, err := tx.Exec(
					tool.DB_change("insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)"),
					value[0],
					set_code,
					"0",
					value[1],
				); err != nil {
					return err
				}
			}
			db_code_str = set_code
			_, err := tx.Exec(
				tool.DB_change("insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, ?, ?, ?)"),
				doc_name,
				"",
				"document_comment_code",
				db_code_str,
			)
			return err
		}); err != nil {
			panic(err)
		}
		tool.Search_bbs_index_update(db, "0", db_code_str)
	}

	return_data := make(map[string]any)
	return_data["response"] = "ok"
	return_data["data"] = db_code_str

	return return_data
}
