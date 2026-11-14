package route

import (
	"database/sql"
	"errors"

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
	deleted_comments := []string{}
	if err := Bbs_identifier_transaction(db, set_id, set_code, false, func(tx *sql.Tx) error {
		rows, err := tx.Query(tool.DB_change("select set_id, set_code from bbs_data where set_name = 'comment' and (set_id = ? or set_id like ?)"), comment_set_id, comment_set_id+"-%")
		if err != nil {
			return err
		}
		for rows.Next() {
			storage_set_id := ""
			storage_set_code := ""
			if err := rows.Scan(&storage_set_id, &storage_set_code); err != nil {
				return errors.Join(err, rows.Close())
			}
			comment_post_id, comment_post_code, comment_code, valid := tool.Search_bbs_index_comment_location(storage_set_id, storage_set_code)
			if valid && comment_post_id == set_id && comment_post_code == set_code {
				deleted_comments = append(deleted_comments, comment_code)
			}
		}
		if err := rows.Err(); err != nil {
			return errors.Join(err, rows.Close())
		}
		if err := rows.Close(); err != nil {
			return err
		}
		if _, err := tx.Exec(tool.DB_change("delete from bbs_data where set_id = ? or set_id like ?"), comment_set_id, comment_set_id+"-%"); err != nil {
			return err
		}
		if _, err := tx.Exec(tool.DB_change("delete from bbs_data where set_name = 'comment_count' and set_id = ? and set_code = ?"), set_id, set_code); err != nil {
			return err
		}
		Bbs_post_last_activity_rebuild(tx, set_id, set_code)
		return nil
	}); err != nil {
		panic(err)
	}
	for _, comment_code := range deleted_comments {
		tool.Search_bbs_index_delete_comment(set_id, set_code, comment_code)
	}

	return_data["response"] = "ok"
	return return_data
}
