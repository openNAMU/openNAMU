package route

import (
	"database/sql"
	"net/url"
	"strconv"
	"strings"
	"time"

	"opennamu/route/tool"
)

type move_history_row struct {
	id   string
	data string
}

func move_title_exists(db *sql.DB, title string) bool {
	queries := []string{
		"select title from data where title = ? limit 1",
		"select title from history where title = ? limit 1",
		"select title from rd where title = ? limit 1",
		"select doc_name from data_set where doc_name = ? limit 1",
		"select title from acl where title = ? limit 1",
		"select title from back where title = ? limit 1",
		"select link from back where link = ? limit 1",
	}
	for _, query := range queries {
		var value string
		if tool.QueryRow_DB(db, query, []any{&value}, title) {
			return true
		}
	}
	return false
}

func move_temporary_title(db *sql.DB) string {
	base := strconv.FormatInt(time.Now().UnixNano(), 10)
	for index := 0; ; index++ {
		title := "__open_namu_move_" + base + "_" + strconv.Itoa(index)
		if !move_title_exists(db, title) {
			return title
		}
	}
}

func move_document_exists(db *sql.DB, title string) bool {
	var value string
	if tool.QueryRow_DB(db, "select title from data where title = ? limit 1", []any{&value}, title) {
		return true
	}
	return tool.QueryRow_DB(db, "select title from history where title = ? limit 1", []any{&value}, title)
}

func move_topic_exists(db *sql.DB, title string) bool {
	var value string
	return tool.QueryRow_DB(db, "select title from rd where title = ? limit 1", []any{&value}, title)
}

func move_data_value(db *sql.DB, title string) string {
	value := ""
	tool.QueryRow_DB(db, "select data from data where title = ?", []any{&value}, title)
	return value
}

func move_history_rows(db *sql.DB, title string) []move_history_row {
	rows := tool.Query_DB(db, "select id, data from history where title = ? order by id + 0 asc", title)
	defer rows.Close()

	result := []move_history_row{}
	for rows.Next() {
		row := move_history_row{}
		if rows.Scan(&row.id, &row.data) == nil {
			result = append(result, row)
		}
	}
	return result
}

func move_history_rows_valid(rows []move_history_row) bool {
	for _, row := range rows {
		if _, err := strconv.Atoi(row.id); err != nil {
			return false
		}
	}
	return true
}

func move_history_max(db *sql.DB, title string) int {
	rows := tool.Query_DB(db, "select id from history where title = ? order by id + 0 desc", title)
	defer rows.Close()

	max_id := 0
	for rows.Next() {
		id := ""
		if rows.Scan(&id) != nil {
			continue
		}
		value, err := strconv.Atoi(id)
		if err == nil && value > max_id {
			max_id = value
		}
	}
	return max_id
}

func move_backlinks(db *sql.DB, old_name string, new_name string) {
	tool.Exec_DB(db, "update back set title = ? where title = ?", new_name, old_name)
	tool.Exec_DB(db, "update back set link = ? where link = ?", new_name, old_name)
}

func move_document_history(db *sql.DB, old_name string, new_name string) {
	tool.Exec_DB(db, "update history set title = ? where title = ?", new_name, old_name)
	tool.Exec_DB(db, "update rc set title = ? where title = ?", new_name, old_name)
}

func move_document_rotate(db *sql.DB, old_name string, new_name string) {
	temp_name := move_temporary_title(db)
	pairs := [][2]string{
		{old_name, temp_name},
		{new_name, old_name},
		{temp_name, new_name},
	}
	for _, pair := range pairs {
		tool.Exec_DB(db, "update data set title = ? where title = ?", pair[1], pair[0])
		tool.Exec_DB(db, "update back set title = ? where title = ?", pair[1], pair[0])
		tool.Exec_DB(db, "update back set link = ? where link = ?", pair[1], pair[0])
		tool.Exec_DB(db, "update history set title = ? where title = ?", pair[1], pair[0])
		tool.Exec_DB(db, "update rc set title = ? where title = ?", pair[1], pair[0])
	}
}

func move_document_merge(db *sql.DB, config tool.Config, old_name string, new_name string, send string, source_data string, source_history []move_history_row) {
	target_max := move_history_max(db, new_name)

	tool.Exec_DB(db, "delete from data where title = ?", new_name)
	tool.Exec_DB(db, "delete from back where link = ?", new_name)
	tool.Exec_DB(db, "update data set title = ? where title = ?", new_name, old_name)
	move_backlinks(db, old_name, new_name)
	tool.Exec_DB(db, "delete from back where title = ? and type = 'no'", new_name)

	for _, row := range source_history {
		id, _ := strconv.Atoi(row.id)
		new_id := strconv.Itoa(target_max + id)
		tool.Exec_DB(db, "update rc set title = ?, id = ? where title = ? and id = ?", new_name, new_id, old_name, row.id)
		tool.Exec_DB(db, "update history set title = ?, id = ? where title = ? and id = ?", new_name, new_id, old_name, row.id)
	}

	tool.Do_add_history(db, new_name, source_data, tool.Get_time(), config.IP, send, "0", "move", "<a>"+tool.HTML_escape(old_name)+"</a> ↔ <a>"+tool.HTML_escape(new_name)+"</a>")
}

func move_topic_normal(db *sql.DB, old_name string, new_name string) {
	tool.Exec_DB(db, "update rd set title = ? where title = ?", new_name, old_name)
}

func move_topic_rotate(db *sql.DB, old_name string, new_name string) {
	temp_name := move_temporary_title(db)
	pairs := [][2]string{
		{old_name, temp_name},
		{new_name, old_name},
		{temp_name, new_name},
	}
	for _, pair := range pairs {
		tool.Exec_DB(db, "update rd set title = ? where title = ?", pair[1], pair[0])
	}
}

func move_data_set_normal(db *sql.DB, old_name string, new_name string) {
	tool.Exec_DB(db, "delete from data_set where doc_name = ?", new_name)
	tool.Exec_DB(db, "delete from acl where title = ?", new_name)
	tool.Exec_DB(db, "update data_set set doc_name = ? where doc_name = ?", new_name, old_name)
	tool.Exec_DB(db, "update acl set title = ? where title = ?", new_name, old_name)
}

func move_data_set_rotate(db *sql.DB, old_name string, new_name string) {
	temp_name := move_temporary_title(db)
	pairs := [][2]string{
		{old_name, temp_name},
		{new_name, old_name},
		{temp_name, new_name},
	}
	for _, pair := range pairs {
		tool.Exec_DB(db, "update data_set set doc_name = ? where doc_name = ?", pair[1], pair[0])
	}
}

func move_document_options(config tool.Config, db *sql.DB, old_name string, new_name string, send string, move_option string, topic_option string, data_set_option string) string {
	if old_name == "" || new_name == "" || old_name == new_name || strings.ContainsAny(new_name, "\r\n") || tool.Get_len(new_name) > 256 {
		return "invalid document"
	}
	if !tool.Do_title_length_check(db, new_name, "document") {
		return "title length"
	}
	if !tool.Check_acl(db, old_name, "", "document_move", config.IP) || !tool.Check_acl(db, new_name, "", "document_move", config.IP) {
		return "auth"
	}

	var source_title string
	if !tool.QueryRow_DB(db, "select title from data where title = ?", []any{&source_title}, old_name) {
		return "not exist"
	}
	if !tool.Arr_in_str([]string{"none", "normal", "reverse", "merge"}, move_option) || !tool.Arr_in_str([]string{"none", "normal", "reverse", "merge"}, topic_option) || !tool.Arr_in_str([]string{"none", "normal", "reverse"}, data_set_option) {
		return "move error"
	}

	owner_auth := tool.Check_acl(db, "", "", "owner_auth", config.IP)
	if (move_option == "merge" || topic_option == "merge" || data_set_option != "none") && !owner_auth {
		return "auth"
	}

	target_exists := move_document_exists(db, new_name)
	if target_exists {
		switch move_option {
		case "normal":
			return "document already exist"
		case "merge":
			source_history := move_history_rows(db, old_name)
			if !move_history_rows_valid(source_history) {
				return "move error"
			}
		case "reverse", "none":
		default:
			return "move error"
		}
	}

	topic_exists := move_topic_exists(db, new_name)
	if topic_exists && topic_option == "normal" {
		return "move error"
	}

	source_data := move_data_value(db, old_name)
	if move_option == "none" {
	} else if target_exists && move_option == "reverse" {
		target_data := move_data_value(db, new_name)
		move_document_rotate(db, old_name, new_name)
		tool.Do_add_history(db, old_name, target_data, tool.Get_time(), config.IP, send, "0", "move", "<a>"+tool.HTML_escape(old_name)+"</a> ⇋ <a>"+tool.HTML_escape(new_name)+"</a>")
		tool.Do_add_history(db, new_name, source_data, tool.Get_time(), config.IP, send, "0", "move", "<a>"+tool.HTML_escape(new_name)+"</a> ⇋ <a>"+tool.HTML_escape(old_name)+"</a>")
	} else if target_exists && move_option == "merge" {
		move_document_merge(db, config, old_name, new_name, send, source_data, move_history_rows(db, old_name))
	} else {
		tool.Exec_DB(db, "update data set title = ? where title = ?", new_name, old_name)
		move_backlinks(db, old_name, new_name)
		move_document_history(db, old_name, new_name)
		tool.Do_add_history(db, new_name, source_data, tool.Get_time(), config.IP, send, "0", "move", "<a>"+tool.HTML_escape(old_name)+"</a> → <a>"+tool.HTML_escape(new_name)+"</a>")
	}

	if topic_option != "none" {
		if topic_exists && topic_option == "reverse" {
			move_topic_rotate(db, old_name, new_name)
		} else if topic_exists && topic_option == "merge" {
			move_topic_normal(db, old_name, new_name)
		} else {
			move_topic_normal(db, old_name, new_name)
		}
	}

	if data_set_option == "reverse" {
		move_data_set_rotate(db, old_name, new_name)
	} else if data_set_option == "normal" {
		move_data_set_normal(db, old_name, new_name)
	}

	return ""
}

func move_document_normal(config tool.Config, db *sql.DB, old_name string, new_name string, send string) string {
	return move_document_options(config, db, old_name, new_name, send, "normal", "normal", "normal")
}

func move_captcha_normal(values url.Values) string {
	return tool.Captcha_response(values.Get("g-recaptcha"), values.Get("g-recaptcha-response"), values.Get("h-captcha-response"), values.Get("cf-turnstile-response"))
}

func View_edit_move(config tool.Config, doc_name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, doc_name, "", "document_move", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	var source_title string
	if !tool.QueryRow_DB(db, "select title from data where title = ?", []any{&source_title}, doc_name) {
		return tool.Get_redirect("/w/" + tool.Url_parser(doc_name))
	}

	if values != nil {
		if !tool.Captcha_check(db, config.Session, config.IP, move_captcha_normal(values)) {
			return tool.Get_error_page(db, config, "recaptcha")
		}
		if !tool.Do_edit_slow_check(db, config, "edit") {
			return tool.Get_error_page(db, config, "slow edit limit")
		}
		if !tool.Do_edit_send_require_check(db, config, values.Get("send")) {
			return tool.Get_error_page(db, config, "send require")
		}
		if !tool.Do_edit_text_checkbox_check(db, config, values.Get("copyright_agreement")) {
			return tool.Get_error_page(db, config, "checkbox check require")
		}

		move_option := "normal"
		topic_option := "normal"
		data_set_option := "normal"
		if _, ok := values["move_option"]; ok {
			move_option = values.Get("move_option")
			topic_option = "none"
			data_set_option = "none"
			if _, ok := values["move_topic_option"]; ok {
				topic_option = values.Get("move_topic_option")
			}
			if _, ok := values["document_set_option"]; ok {
				data_set_option = values.Get("document_set_option")
			}
		}

		new_name := strings.TrimSpace(values.Get("title"))
		if err_name := move_document_options(config, db, doc_name, new_name, values.Get("send"), move_option, topic_option, data_set_option); err_name != "" {
			return tool.Get_error_page(db, config, err_name)
		}
		return tool.Get_redirect("/w/" + tool.Url_parser(new_name))
	}

	owner_auth := tool.Check_acl(db, "", "", "owner_auth", config.IP)
	body := "<form method=\"post\">"
	body += "<input name=\"title\" value=\"" + tool.HTML_escape(doc_name) + "\" placeholder=\"" + tool.Get_language(db, "document_name", true) + "\"><hr class=\"main_hr\">"
	body += "<input name=\"send\" placeholder=\"" + tool.Get_language(db, "why", true) + "\"><hr class=\"main_hr\">"

	body += "<h2>" + tool.Get_language(db, "document", true) + "</h2>"
	body += "<select name=\"move_option\">"
	body += "<option value=\"normal\" selected>" + tool.Get_language(db, "normal", true) + "</option>"
	body += "<option value=\"none\">" + tool.Get_language(db, "dont_move", true) + "</option>"
	body += "<option value=\"reverse\">" + tool.Get_language(db, "replace_move", true) + "</option>"
	if owner_auth {
		body += "<option value=\"merge\">" + tool.Get_language(db, "merge_move", true) + "</option>"
	}
	body += "</select><hr class=\"main_hr\">"

	body += "<h2>" + tool.Get_language(db, "discussion", true) + "</h2>"
	body += "<select name=\"move_topic_option\">"
	body += "<option value=\"none\" selected>" + tool.Get_language(db, "dont_move", true) + "</option>"
	body += "<option value=\"normal\">" + tool.Get_language(db, "normal", true) + "</option>"
	body += "<option value=\"reverse\">" + tool.Get_language(db, "replace_move", true) + "</option>"
	if owner_auth {
		body += "<option value=\"merge\">" + tool.Get_language(db, "merge_move", true) + "</option>"
	}
	body += "</select><hr class=\"main_hr\">"

	if owner_auth {
		body += "<h2>" + tool.Get_language(db, "document_set", true) + "</h2>"
		body += "<select name=\"document_set_option\">"
		body += "<option value=\"none\" selected>" + tool.Get_language(db, "dont_move", true) + "</option>"
		body += "<option value=\"normal\">" + tool.Get_language(db, "normal", true) + "</option>"
		body += "<option value=\"reverse\">" + tool.Get_language(db, "replace_move", true) + "</option>"
		body += "</select><hr class=\"main_hr\">"
	}

	body += tool.Get_captcha_ui(db, config) + tool.Get_IP_warning_ui(db, config) + tool.Get_edit_check_box_ui(db) + tool.Get_edit_bottom_text_ui(db, "move")
	body += "<button type=\"submit\">" + tool.Get_language(db, "move", true) + "</button></form>"
	return tool.Get_template(db, config, doc_name, body, []any{"(" + tool.Get_language(db, "move", true) + ")"}, [][]any{{"w/" + tool.Url_parser(doc_name), tool.Get_language(db, "return", true)}, {"move_all", tool.Get_language(db, "multiple_move", true)}}, map[string]string{})
}
