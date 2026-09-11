package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_bbs_w_comment_all_delete(config tool.Config, set_id string, set_code string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if set_id != "0" && set_id != thread_bbs_id {
		return_data["response"] = "not allowed"
		return return_data
	}
	if _, exists := tool.Get_bbs_data_value(db, set_id, set_code, "title"); !exists {
		return_data["response"] = "not exist"
		return_data["data"] = "post"
		return return_data
	}
	if !tool.Check_permission(db, "bbs_comment_manage", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	comment_set_id := set_id + "-" + set_code
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		if _, err := tx.Exec(tool.DB_change("delete from bbs_data where set_id = ? or set_id like ?"), comment_set_id, comment_set_id+"-%"); err != nil {
			return err
		}
		if _, err := tx.Exec(tool.DB_change("delete from bbs_data where set_name = 'comment_count' and set_id = ? and set_code = ?"), set_id, set_code); err != nil {
			return err
		}
		bbs_post_last_activity_rebuild(tx, set_id, set_code)
		return nil
	}); err != nil {
		panic(err)
	}
	tool.Search_bbs_index_delete_comments(db, set_id, set_code)

	return_data["response"] = "ok"
	return return_data
}
