package route

import (
	"strings"

	"opennamu/route/tool"
	"opennamu/route/tool/markup"
)

func Api_edit_post(config tool.Config, doc_name string, data string, send string, agree string, expected_revision string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)

	date := tool.Get_time()
	data = strings.ReplaceAll(data, "\r", "")
	current_revision := tool.Get_document_revision(db, doc_name)

	if doc_name == "" {
		return_data["response"] = "error"
		return_data["data"] = "empty title"

		return return_data
	} else if !tool.Do_title_length_check(db, doc_name, "document") {
		return_data["response"] = "error"
		return_data["data"] = "title length"

		return return_data
	} else if !tool.Check_acl(db, doc_name, "", "document_edit", config.IP) {
		return_data["response"] = "error"
		return_data["data"] = "permission denied"

		return return_data
	} else if current_revision == "0" && !tool.Check_acl(db, doc_name, "", "document_make_acl", config.IP) {
		return_data["response"] = "error"
		return_data["data"] = "permission denied"

		return return_data
	} else if expected_revision != "" && expected_revision != current_revision {
		return_data["response"] = "error"
		return_data["data"] = "edit conflict"

		return return_data
	} else if !tool.Do_edit_slow_check(db, config, "edit") {
		return_data["response"] = "error"
		return_data["data"] = "slow edit limit"

		return return_data
	} else if !tool.Do_edit_filter(db, config, doc_name, data) {
		return_data["response"] = "error"
		return_data["data"] = "edit filter (content)"

		return return_data
	} else if !tool.Do_edit_filter(db, config, doc_name, send) {
		return_data["response"] = "error"
		return_data["data"] = "edit filter (send)"

		return return_data
	} else if !tool.Do_edit_send_require_check(db, config, send) {
		return_data["response"] = "error"
		return_data["data"] = "send require"

		return return_data
	} else if !tool.Do_edit_text_checkbox_check(db, config, agree) {
		return_data["response"] = "error"
		return_data["data"] = "checkbox check require"

		return return_data
	} else if !tool.Do_edit_max_length_check(db, config, data) {
		return_data["response"] = "error"
		return_data["data"] = "overflow max length"

		return return_data
	}

	var old_data string

	tool.QueryRow_DB(
		db,
		`select data from data where title = ?`,
		[]any{&old_data},
		doc_name,
	)

	length := tool.Get_edit_length_diff(old_data, data)

	tool.Exec_DB(
		db,
		`delete from data where title = ?`,
		doc_name,
	)
	tool.Exec_DB(
		db,
		`insert into data (title, data) values (?, ?)`,
		doc_name,
		data,
	)

	tool.Do_watchlist_alarm_send(db, config, doc_name)

	tool.Do_add_history(
		db,
		doc_name,
		data,
		date,
		config.IP,
		send,
		length,
		"",
		"",
	)

	markup.Get_render(db, doc_name, data, "backlink")

	return_data["response"] = "ok"

	return return_data
}
