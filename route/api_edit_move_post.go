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

type move_history_event struct {
	doc_name   string
	data       string
	send       string
	type_check string
}

func move_title_exists(db *sql.DB, title string) bool {
	queries := []string{
		"select title from data where title = ? limit 1",
		"select title from history where title = ? limit 1",
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

func move_document_exists(db *sql.DB, title string) (bool, bool) {
	var value string
	data_exists := tool.QueryRow_DB(db, "select title from data where title = ? limit 1", []any{&value}, title)
	history_exists := tool.QueryRow_DB(db, "select title from history where title = ? limit 1", []any{&value}, title)
	topic_exists := move_topic_exists(db, title)
	return data_exists || history_exists || topic_exists, history_exists && !data_exists
}

func move_topic_exists(db *sql.DB, title string) bool {
	var value string
	return tool.QueryRow_DB(db, "select set_code from bbs_data where set_id = ? and set_name = 'document' and set_data = ? limit 1", []any{&value}, thread_bbs_id, title)
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

func move_backlinks(executor tool.DB_runner, old_name string, new_name string) {
	tool.Exec_DB(executor, "update back set title = ? where title = ?", new_name, old_name)
	tool.Exec_DB(executor, "update back set link = ? where link = ?", new_name, old_name)
}

func move_document_normal_data(executor tool.DB_runner, old_name string, new_name string) {
	tool.Exec_DB(executor, "update data set title = ? where title = ?", new_name, old_name)
	move_backlinks(executor, old_name, new_name)
	move_document_history(executor, old_name, new_name)
}

func move_document_history(executor tool.DB_runner, old_name string, new_name string) {
	tool.Exec_DB(executor, "update history set title = ? where title = ?", new_name, old_name)
	tool.Exec_DB(executor, "update rc set title = ? where title = ?", new_name, old_name)
}

func move_document_rotate(executor tool.DB_runner, temp_name string, old_name string, new_name string) {
	pairs := [][2]string{
		{old_name, temp_name},
		{new_name, old_name},
		{temp_name, new_name},
	}
	for _, pair := range pairs {
		tool.Exec_DB(executor, "update data set title = ? where title = ?", pair[1], pair[0])
		tool.Exec_DB(executor, "update back set title = ? where title = ?", pair[1], pair[0])
		tool.Exec_DB(executor, "update back set link = ? where link = ?", pair[1], pair[0])
		tool.Exec_DB(executor, "update history set title = ? where title = ?", pair[1], pair[0])
		tool.Exec_DB(executor, "update rc set title = ? where title = ?", pair[1], pair[0])
	}
}

func move_document_merge(executor tool.DB_runner, target_max int, old_name string, new_name string, source_history []move_history_row) {
	tool.Exec_DB(executor, "delete from data where title = ?", new_name)
	tool.Exec_DB(executor, "delete from back where link = ? and type != 'cat_manual'", new_name)
	tool.Exec_DB(executor, "update data set title = ? where title = ?", new_name, old_name)
	move_backlinks(executor, old_name, new_name)
	tool.Exec_DB(executor, "delete from back where title = ? and type = 'no'", new_name)

	for _, row := range source_history {
		id, _ := strconv.Atoi(row.id)
		new_id := strconv.Itoa(target_max + id)
		tool.Exec_DB(executor, "update rc set title = ?, id = ? where title = ? and id = ?", new_name, new_id, old_name, row.id)
		tool.Exec_DB(executor, "update history set title = ?, id = ? where title = ? and id = ?", new_name, new_id, old_name, row.id)
	}
}

func move_topic_normal(executor tool.DB_runner, old_name string, new_name string) {
	tool.Exec_DB(executor, "update bbs_data set set_data = ? where set_id = ? and set_name = 'document' and set_data = ?", new_name, thread_bbs_id, old_name)
	tool.Exec_DB(executor, "update bbs_data set set_data = ? where set_id = ? and set_name = 'tag' and set_data = ?", new_name, thread_bbs_id, old_name)
}

func move_topic_rotate(executor tool.DB_runner, temp_name string, old_name string, new_name string) {
	pairs := [][2]string{
		{old_name, temp_name},
		{new_name, old_name},
		{temp_name, new_name},
	}
	for _, pair := range pairs {
		tool.Exec_DB(executor, "update bbs_data set set_data = ? where set_id = ? and set_name = 'document' and set_data = ?", pair[1], thread_bbs_id, pair[0])
		tool.Exec_DB(executor, "update bbs_data set set_data = ? where set_id = ? and set_name = 'tag' and set_data = ?", pair[1], thread_bbs_id, pair[0])
	}
}

func move_data_set_normal(executor tool.DB_runner, old_name string, new_name string) {
	tool.Exec_DB(executor, "delete from data_set where doc_name = ?", new_name)
	tool.Exec_DB(executor, "delete from acl where title = ?", new_name)
	tool.Exec_DB(executor, "update data_set set doc_name = ? where doc_name = ?", new_name, old_name)
	tool.Exec_DB(executor, "update acl set title = ? where title = ?", new_name, old_name)
	tool.Exec_DB(executor, "update bbs_data set set_data = ? where set_id = '0' and set_name = 'title' and set_data = ?", new_name, old_name)
}

func move_data_set_rotate(executor tool.DB_runner, temp_name string, old_name string, new_name string) {
	pairs := [][2]string{
		{old_name, temp_name},
		{new_name, old_name},
		{temp_name, new_name},
	}
	for _, pair := range pairs {
		tool.Exec_DB(executor, "update data_set set doc_name = ? where doc_name = ?", pair[1], pair[0])
		tool.Exec_DB(executor, "update acl set title = ? where title = ?", pair[1], pair[0])
	}
}

func move_document_settings_exists(db *sql.DB, title string) bool {
	var value string
	if tool.QueryRow_DB(db, "select doc_name from data_set where doc_name = ? limit 1", []any{&value}, title) {
		return true
	}
	return tool.QueryRow_DB(db, "select title from acl where title = ? limit 1", []any{&value}, title)
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

	owner_auth := tool.Check_permission(db, "document_move_manage", config.IP)
	if !owner_auth && (move_option == "merge" || topic_option == "merge" || data_set_option == "reverse" || (move_option == "normal" && data_set_option != "normal") || (move_option != "normal" && data_set_option != "none")) {
		return "auth"
	}

	target_exists, target_history_only := move_document_exists(db, new_name)
	source_history := []move_history_row{}
	if target_exists {
		switch move_option {
		case "normal":
			if target_history_only {
				return "document history exist"
			}
			return "document already exist"
		case "merge":
			source_history = move_history_rows(db, old_name)
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
	if move_option == "normal" && data_set_option == "normal" && !owner_auth && move_document_settings_exists(db, new_name) {
		return "move error"
	}

	document_temp := ""
	if target_exists && move_option == "reverse" {
		document_temp = move_temporary_title(db)
	}
	topic_temp := ""
	if topic_exists && topic_option == "reverse" {
		topic_temp = move_temporary_title(db)
	}
	data_set_temp := ""
	if data_set_option == "reverse" {
		data_set_temp = move_temporary_title(db)
	}

	source_data := move_data_value(db, old_name)
	target_data := ""
	if target_exists && move_option == "reverse" {
		target_data = move_data_value(db, new_name)
	}
	target_max := 0
	if move_option == "merge" {
		target_max = move_history_max(db, new_name)
	}

	tx, err := db.Begin()
	if err != nil {
		return "move error"
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	history_events := []move_history_event{}
	if move_option == "none" {
	} else if target_exists && move_option == "reverse" {
		move_document_rotate(tx, document_temp, old_name, new_name)
		history_events = append(history_events,
			move_history_event{old_name, target_data, send, "<a>" + tool.HTML_escape(old_name) + "</a> ⇋ <a>" + tool.HTML_escape(new_name) + "</a>"},
			move_history_event{new_name, source_data, send, "<a>" + tool.HTML_escape(new_name) + "</a> ⇋ <a>" + tool.HTML_escape(old_name) + "</a>"},
		)
	} else if target_exists && move_option == "merge" {
		move_document_merge(tx, target_max, old_name, new_name, source_history)
		history_events = append(history_events, move_history_event{new_name, source_data, send, "<a>" + tool.HTML_escape(old_name) + "</a> ↔ <a>" + tool.HTML_escape(new_name) + "</a>"})
	} else {
		move_document_normal_data(tx, old_name, new_name)
		history_events = append(history_events, move_history_event{new_name, source_data, send, "<a>" + tool.HTML_escape(old_name) + "</a> → <a>" + tool.HTML_escape(new_name) + "</a>"})
	}

	if topic_option != "none" {
		if topic_exists && topic_option == "reverse" {
			move_topic_rotate(tx, topic_temp, old_name, new_name)
		} else {
			move_topic_normal(tx, old_name, new_name)
		}
	}

	if data_set_option == "reverse" {
		move_data_set_rotate(tx, data_set_temp, old_name, new_name)
	} else if data_set_option == "normal" {
		move_data_set_normal(tx, old_name, new_name)
	}

	for _, event := range history_events {
		tool.Do_add_history(tx, event.doc_name, event.data, tool.Get_time(), config.IP, event.send, "0", "move", event.type_check)
	}

	if err := tx.Commit(); err != nil {
		return "move error"
	}
	committed = true
	tool.Search_index_sync(db, old_name)
	tool.Search_index_sync(db, new_name)
	return ""
}

func move_document_normal(config tool.Config, db *sql.DB, old_name string, new_name string, send string) string {
	return move_document_options(config, db, old_name, new_name, send, "normal", "normal", "normal")
}

func move_captcha_normal(values url.Values) string {
	return tool.Captcha_response(values.Get("g-recaptcha"), values.Get("g-recaptcha-response"), values.Get("h-captcha-response"), values.Get("cf-turnstile-response"), values.Get("altcha"))
}

func Api_edit_move_post(config tool.Config, doc_name string, values url.Values) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_acl(db, doc_name, "", "document_move", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if !tool.Captcha_check(db, config.Session, config.IP, move_captcha_normal(values)) {
		return_data["response"] = "error"
		return_data["data"] = "recaptcha"
		return return_data
	}
	if !tool.Do_edit_slow_check(db, config, "edit") {
		return_data["response"] = "error"
		return_data["data"] = "slow edit limit"
		return return_data
	}
	if !tool.Do_edit_send_require_check(db, config, values.Get("send")) {
		return_data["response"] = "error"
		return_data["data"] = "send require"
		return return_data
	}
	if !tool.Do_edit_text_checkbox_check(db, config, values.Get("copyright_agreement")) {
		return_data["response"] = "error"
		return_data["data"] = "checkbox check require"
		return return_data
	}

	move_option := "normal"
	topic_option := "normal"
	data_set_option := "normal"
	if _, ok := values["move_option"]; ok {
		move_option = values.Get("move_option")
		topic_option = "none"
		if move_option != "normal" {
			data_set_option = "none"
		}
		if _, ok := values["move_topic_option"]; ok {
			topic_option = values.Get("move_topic_option")
		}
		if _, ok := values["document_set_option"]; ok {
			data_set_option = values.Get("document_set_option")
		}
	}

	new_name := strings.TrimSpace(values.Get("title"))
	if err_name := move_document_options(config, db, doc_name, new_name, values.Get("send"), move_option, topic_option, data_set_option); err_name != "" {
		return_data["response"] = "error"
		return_data["data"] = err_name
		return return_data
	}

	return_data["response"] = "ok"
	return_data["data"] = new_name
	return return_data
}
