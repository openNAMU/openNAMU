package route

import (
	"database/sql"
	"strconv"
	"strings"

	"opennamu/route/tool"
)

func Api_bbs_report_post(config tool.Config, set_id string, set_code string, comment_code string, reason string, captcha string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "bbs_comment", config.IP) {
		return map[string]any{"response": "require auth"}
	}
	if set_id == report_bbs_id {
		return map[string]any{"response": "not exist", "data": "post"}
	}

	reason = strings.TrimSpace(strings.ReplaceAll(reason, "\r", ""))
	if reason == "" || tool.Get_len(reason) > 1000 {
		return map[string]any{"response": "error", "data": "report reason"}
	}
	if !tool.Captcha_check(db, config.Session, config.IP, captcha) {
		return map[string]any{"response": "error", "data": "recaptcha"}
	}
	if !tool.Do_edit_filter(db, config, "", reason) {
		return map[string]any{"response": "error", "data": "edit filter (content)"}
	}
	if !tool.Do_bbs_max_length_check(db, config, reason) {
		return map[string]any{"response": "error", "data": "bbs overflow max length"}
	}

	target_api := Api_bbs_w(config, set_id, set_code)
	if target_api["response"] != "ok" {
		return map[string]any{"response": target_api["response"], "data": "post"}
	}
	target_data, _ := target_api["data"].(map[string]string)
	target_title := target_data["title"]
	if target_title == "" {
		return map[string]any{"response": "not exist", "data": "post"}
	}

	target_path := "/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code)
	report_prefix := tool.Get_language(db, "report_post", true)
	if comment_code != "" {
		if !bbs_comment_code_regex.MatchString(comment_code) {
			return map[string]any{"response": "not exist", "data": "comment"}
		}
		comment_api := Api_bbs_w_comment_one(config, false, "", set_id+"-"+set_code+"-"+comment_code)
		comment_list, _ := comment_api["data"].([]map[string]string)
		if comment_api["response"] != "ok" || len(comment_list) == 0 || comment_list[0]["comment"] == "" {
			return map[string]any{"response": comment_api["response"], "data": "comment"}
		}
		target_path += "/comment/" + tool.Url_parser(comment_code)
		report_prefix = tool.Get_language(db, "report_comment", true)
	}
	if !tool.Check_daily_limit(db, config.IP, "bbs_edit") {
		return map[string]any{"response": "error", "data": "daily limit"}
	}

	set_code_new, err := Report_create(db, config.IP, target_title, target_path, report_prefix, reason)
	if err != nil {
		panic(err)
	}
	return map[string]any{"response": "ok", "data": set_code_new}
}

func Report_create(db *sql.DB, ip string, target_title string, target_path string, report_prefix string, reason string) (string, error) {
	set_code_new := ""
	err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		last_code := ""
		tool.QueryRow_DB(
			tx,
			"select set_code from bbs_data where set_name = 'title' and set_id = ? order by set_code + 0 desc limit 1",
			[]any{&last_code},
			report_bbs_id,
		)
		set_code_new = strconv.Itoa(tool.Str_to_int(last_code) + 1)
		date_now := tool.Get_time()
		for _, value := range [][]string{
			{"title", target_title},
			{"data", reason},
			{"prefix", report_prefix},
			{"date", date_now},
			{"last_activity", date_now},
			{"user_id", ip},
			{"comment_count", "0"},
			{"report_target", target_path},
		} {
			if _, err := tx.Exec(
				tool.DB_change("insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)"),
				value[0],
				set_code_new,
				report_bbs_id,
				value[1],
			); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	tool.Search_bbs_index_update(db, report_bbs_id, set_code_new)
	return set_code_new, nil
}
