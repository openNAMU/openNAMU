package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_bbs_w_comment_blind_post(config tool.Config, set_id string, set_code string, comment_code string, blind bool) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "bbs_comment_manage", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	comment_set_id, comment_set_code, exists := bbs_search_comment_location(set_id, set_code, comment_code)
	if !exists {
		return_data["response"] = "not exist"
		return_data["data"] = "comment"
		return return_data
	}

	comment, comment_exists := tool.Get_bbs_data_value(db, comment_set_id, comment_set_code, "comment")
	if !comment_exists || comment == "" {
		return_data["response"] = "not exist"
		return_data["data"] = "comment"
		return return_data
	}

	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		blind_value := ""
		blind_exists := tool.QueryRow_DB(tx, "select set_data from bbs_data where set_name = 'blind' and set_id = ? and set_code = ?", []any{&blind_value}, comment_set_id, comment_set_code)
		if blind {
			if blind_exists {
				_, err := tx.Exec(tool.DB_change("update bbs_data set set_data = 'O' where set_name = 'blind' and set_id = ? and set_code = ?"), comment_set_id, comment_set_code)
				return err
			}
			_, err := tx.Exec(tool.DB_change("insert into bbs_data (set_name, set_id, set_code, set_data) values ('blind', ?, ?, 'O')"), comment_set_id, comment_set_code)
			return err
		}
		_, err := tx.Exec(tool.DB_change("delete from bbs_data where set_name = 'blind' and set_id = ? and set_code = ?"), comment_set_id, comment_set_code)
		return err
	}); err != nil {
		panic(err)
	}

	return_data["response"] = "ok"
	if blind {
		return_data["data"] = "O"
	} else {
		return_data["data"] = ""
	}
	return return_data
}
