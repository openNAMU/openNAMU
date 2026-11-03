package route

import (
	"database/sql"

	"opennamu/route/tool"
)

const post_slow_mode_name = "post_slow_mode_seconds"

// Api_bbs_w_post_slow_mode returns the current post-level slow-mode interval
// for a post. The value is empty when no post-level slow mode is set.
func Api_bbs_w_post_slow_mode(config tool.Config, set_id string, set_code string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "bbs_post_manage", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	value, _ := tool.Get_bbs_data_value(db, set_id, set_code, post_slow_mode_name)
	return_data["response"] = "ok"
	return_data["data"] = value
	return return_data
}

// Api_bbs_w_post_slow_mode_put validates and stores the post-level slow-mode
// interval. It requires bbs_post_manage permission and rejects any value that
// is not an integer in [0, 86400]. On failure the previous value is kept.
func Api_bbs_w_post_slow_mode_put(config tool.Config, set_id string, set_code string, data string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "bbs_post_manage", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if Bbs_is_special_board(set_id) && set_id != "0" && set_id != thread_bbs_id {
		return_data["response"] = "not allowed"
		return return_data
	}
	data = tool.Bbs_slow_mode_normalize(data)
	if !tool.Bbs_slow_mode_valid(data) {
		return_data["response"] = "error"
		return_data["data"] = "invalid bbs slow mode"
		return return_data
	}
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		if _, err := tx.Exec(tool.DB_change("delete from bbs_data where set_name = ? and set_id = ? and set_code = ?"), post_slow_mode_name, set_id, set_code); err != nil {
			return err
		}
		if _, err := tx.Exec(tool.DB_change("insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)"), post_slow_mode_name, set_code, set_id, data); err != nil {
			return err
		}
		tool.Do_insert_auth_history(tx, config.IP, "bbs_set ("+set_id+"/"+set_code+"/"+post_slow_mode_name+")")
		return nil
	}); err != nil {
		panic(err)
	}
	return_data["response"] = "ok"
	return return_data
}
