package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_bbs_w_comment_pinned(config tool.Config, set_id string, set_code string, comment_code string, toggle bool) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)

	if comment_code == "" {
		return_data["response"] = "not exist"
		return_data["data"] = "comment"

		return return_data
	}

	comment_set_id, comment_set_code, exists := bbs_search_comment_location(set_id, set_code, comment_code)
	if !exists {
		return_data["response"] = "not exist"
		return_data["data"] = "comment"

		return return_data
	}

	comment := ""
	if !tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'comment' and set_id = ? and set_code = ?",
		[]any{&comment},
		comment_set_id,
		comment_set_code,
	) || comment == "" {
		return_data["response"] = "not exist"
		return_data["data"] = "comment"

		return return_data
	}

	if !tool.Check_permission(db, "bbs_comment_manage", config.IP) {
		return_data["response"] = "require auth"

		return return_data
	}

	if toggle {
		if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
			pinned := ""
			pinned_exist := tool.QueryRow_DB(tx, "select set_data from bbs_data where set_name = 'pinned' and set_id = ? and set_code = ?", []any{&pinned}, comment_set_id, comment_set_code)
			if pinned_exist {
				_, err := tx.Exec(tool.DB_change("delete from bbs_data where set_name = 'pinned' and set_id = ? and set_code = ?"), comment_set_id, comment_set_code)
				return err
			}
			_, err := tx.Exec(tool.DB_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('pinned', ?, ?, ?)"), comment_set_code, comment_set_id, tool.Get_time())
			return err
		}); err != nil {
			panic(err)
		}
	}

	pinned := ""
	pinned_exist := tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'pinned' and set_id = ? and set_code = ?",
		[]any{&pinned},
		comment_set_id,
		comment_set_code,
	)

	return_data["response"] = "ok"
	if pinned_exist {
		return_data["data"] = "pinned_release"
	} else {
		return_data["data"] = "pinned"
	}

	return return_data
}
