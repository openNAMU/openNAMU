package route

import (
	"database/sql"
	"errors"
	"strings"

	"opennamu/route/tool"
)

func Bbs_comment_vote_target_check(tx *sql.Tx, set_id string, set_code string, comment_manage bool) (string, error) {
	var comment string
	err := tx.QueryRow(
		tool.DB_change("select set_data from bbs_data where set_name = 'comment' and set_id = ? and set_code = ? limit 1"),
		set_id,
		set_code,
	).Scan(&comment)
	if errors.Is(err, sql.ErrNoRows) || (err == nil && strings.TrimSpace(comment) == "") {
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
	if err == nil && !comment_manage {
		return "require auth", nil
	}

	return "ok", nil
}

func Api_bbs_w_comment_tabom_post(config tool.Config, set_id string, set_code string, comment_code string, vote_type string) map[string]any {
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

	comment_set_id, comment_set_code, exists := Bbs_search_comment_location(set_id, set_code, comment_code)
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
	) || strings.TrimSpace(comment) == "" {
		return_data["response"] = "not exist"
		return_data["data"] = "comment"
		return return_data
	}
	blind := ""
	tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'blind' and set_id = ? and set_code = ?",
		[]any{&blind},
		comment_set_id,
		comment_set_code,
	)
	comment_manage := tool.Check_permission(db, "bbs_comment_manage", config.IP)
	if blind == "O" && !comment_manage {
		return_data["response"] = "require auth"
		return return_data
	}

	var result map[string]any
	if err := Bbs_board_write_transaction(db, set_id, set_code, func(tx *sql.Tx) error {
		target_response, target_err := Bbs_comment_vote_target_check(tx, comment_set_id, comment_set_code, comment_manage)
		if target_err != nil {
			return target_err
		}
		if target_response != "ok" {
			result = map[string]any{"response": target_response}
			if target_response == "not exist" {
				result["data"] = "comment"
			}
			return nil
		}
		result, target_err = Api_bbs_tabom_post(tx, config.IP, comment_set_id, comment_set_code, vote_type)
		return target_err
	}); err != nil {
		if Bbs_board_missing_response(return_data, err, "bbs") {
			return return_data
		}
		return map[string]any{"response": "error", "data": "database"}
	}

	if result["response"] == "ok" {
		if comment_page, exists := Bbs_comment_page_for_code(db, set_id, set_code, comment_code); exists {
			if result_data, ok := result["data"].(map[string]any); ok {
				result_data["page"] = comment_page
			}
		}
	}
	return result
}
