package route

import (
	"database/sql"
	"strconv"
	"strings"

	"opennamu/route/tool"
)

func Api_history_report_post(config tool.Config, doc_name string, rev string, reason string, captcha string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	revision, response, error_data := History_report_target(db, config, doc_name, rev)
	if response != "ok" {
		return map[string]any{"response": response, "data": error_data}
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
	if !tool.Check_daily_limit(db, config.IP, "bbs_edit") {
		return map[string]any{"response": "error", "data": "daily limit"}
	}

	target_path := "/history_tool/" + revision + "/" + tool.Url_parser(doc_name)
	set_code, err := Report_create(
		db,
		config.IP,
		doc_name,
		target_path,
		tool.Get_language(db, "report_revision", true)+" r"+revision,
		reason,
	)
	if err != nil {
		panic(err)
	}
	return map[string]any{"response": "ok", "data": set_code}
}

func Api_history_report_target(config tool.Config, doc_name string, rev string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	revision, response, error_data := History_report_target(db, config, doc_name, rev)
	if response != "ok" {
		return map[string]any{"response": response, "data": error_data}
	}
	return map[string]any{"response": "ok", "data": revision}
}

func History_report_target(db *sql.DB, config tool.Config, doc_name string, rev string) (string, string, string) {
	if !tool.Check_permission(db, "bbs_comment", config.IP) || !tool.Check_permission(db, "history_view", config.IP) {
		return "", "require auth", ""
	}
	if doc_name == "" || !tool.Check_acl(db, doc_name, "", "render", config.IP) {
		return "", "require auth", ""
	}

	revision_number, err := strconv.Atoi(strings.TrimSpace(rev))
	if err != nil || revision_number < 1 {
		return "", "error", "invalid revision"
	}
	revision := strconv.Itoa(revision_number)
	hide, exists := tool.Get_history_revision_hide(db, doc_name, revision)
	if !exists {
		return "", "not exist", "revision"
	}
	if hide != "" && !tool.Check_permission(db, "hidel", config.IP) {
		return "", "require auth", ""
	}
	return revision, "ok", ""
}
