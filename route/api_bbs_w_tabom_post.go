package route

import (
	"database/sql"
	"errors"

	"opennamu/route/tool"
)

func Bbs_post_vote_target_check(tx *sql.Tx, db *sql.DB, set_id string, set_code string, user string, auth_info map[string]bool) (string, error) {
	var title_exists int
	err := tx.QueryRow(
		tool.DB_change("select 1 from bbs_data where set_name = 'title' and set_id = ? and set_code = ? limit 1"),
		set_id,
		set_code,
	).Scan(&title_exists)
	if errors.Is(err, sql.ErrNoRows) {
		return "not exist", nil
	}
	if err != nil {
		return "", err
	}

	var post_user string
	err = tx.QueryRow(
		tool.DB_change("select set_data from bbs_data where set_name = 'user_id' and set_id = ? and set_code = ? limit 1"),
		set_id,
		set_code,
	).Scan(&post_user)
	if errors.Is(err, sql.ErrNoRows) {
		return "not exist", nil
	}
	if err != nil {
		return "", err
	}

	var blind_exists int
	err = tx.QueryRow(
		tool.DB_change("select 1 from bbs_data where set_name = 'blind' and set_data = 'O' and set_id = ? and set_code = ? limit 1"),
		set_id,
		set_code,
	).Scan(&blind_exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}
	if err == nil && !auth_info["bbs_post_manage"] {
		return "require auth", nil
	}

	var private_acl string
	err = tx.QueryRow(
		tool.DB_change("select set_data from bbs_set where set_name = 'bbs_only_my_data_view_acl' and set_code = '' and set_id = ? limit 1"),
		set_id,
	).Scan(&private_acl)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}
	if private_acl != "" && private_acl != "normal" && post_user != user && !auth_info["bbs"] && !tool.Check_acl_group(db, private_acl, auth_info) {
		return "require auth", nil
	}

	return "ok", nil
}

func Api_bbs_w_tabom_post(config tool.Config, set_id string, set_code string, vote_type string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if Bbs_is_special_board(set_id) && set_id != "0" && set_id != thread_bbs_id {
		return_data["response"] = "not allowed"
		return return_data
	}

	if _, allowed := Bbs_post_view_auth(db, set_id, set_code, config.IP); !allowed {
		return_data["response"] = "require auth"
		return return_data
	}

	if !tool.Check_acl(db, set_id, "", "bbs_comment", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if !Bbs_tabom_vote_type_valid(vote_type) {
		return_data["response"] = "error"
		return_data["data"] = "vote_type"
		return return_data
	}

	auth_info := tool.Get_auth_info(db, config.IP)
	var result map[string]any
	if err := Bbs_board_write_transaction(db, set_id, set_code, func(tx *sql.Tx) error {
		target_response, target_err := Bbs_post_vote_target_check(tx, db, set_id, set_code, config.IP, auth_info)
		if target_err != nil {
			return target_err
		}
		if target_response != "ok" {
			result = map[string]any{"response": target_response}
			if target_response == "not exist" {
				result["data"] = "post"
			}
			return nil
		}
		result, target_err = Api_bbs_tabom_post(tx, config.IP, set_id, set_code, vote_type)
		return target_err
	}); err != nil {
		if Bbs_board_missing_response(return_data, err, "bbs") {
			return return_data
		}
		return map[string]any{"response": "error", "data": "database"}
	}
	return result
}
