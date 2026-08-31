package route

import (
	"database/sql"
	"regexp"
	"strconv"
	"strings"

	"opennamu/route/tool"
)

var bbs_comment_code_regex = regexp.MustCompile(`^[0-9]+(?:-[0-9]+)*$`)

func bbs_comment_closed(db *sql.DB, set_id string, set_code string) bool {
	closed := ""
	return tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'comment_close' and set_id = ? and set_code = ?",
		[]any{&closed},
		set_id,
		set_code,
	) && closed == "1"
}

func bbs_comment_parent(db *sql.DB, set_id string, set_code string, comment_select string) (string, string, bool) {
	base_id := set_id + "-" + set_code
	if comment_select == "" || comment_select == "0" {
		return base_id, "", true
	}

	if !bbs_comment_code_regex.MatchString(comment_select) {
		return "", "", false
	}

	parts := strings.Split(comment_select, "-")
	parent_id := base_id
	if len(parts) > 1 {
		parent_id += "-" + strings.Join(parts[:len(parts)-1], "-")
	}

	parent_user := ""
	if !tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'comment_user_id' and set_id = ? and set_code = ?",
		[]any{&parent_user},
		parent_id,
		parts[len(parts)-1],
	) {
		return "", "", false
	}

	return base_id + "-" + comment_select, parent_user, true
}

func Api_bbs_w_comment_post(config tool.Config, set_id string, set_code string, comment_select string, data string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)

	bbs_name := ""
	bbs_type := "comment"
	title := ""
	post_user := ""
	if !tool.QueryRow_DB(
		db,
		"select set_data from bbs_set where set_name = 'bbs_name' and set_id = ?",
		[]any{&bbs_name},
		set_id,
	) || !tool.QueryRow_DB(
		db,
		"select set_data from bbs_set where set_name = 'bbs_type' and set_id = ?",
		[]any{&bbs_type},
		set_id,
	) || !tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'title' and set_id = ? and set_code = ?",
		[]any{&title},
		set_id,
		set_code,
	) || !tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'user_id' and set_id = ? and set_code = ?",
		[]any{&post_user},
		set_id,
		set_code,
	) {
		return_data["response"] = "not exist"
		return_data["data"] = "post"
		return return_data
	}

	if !tool.Check_acl(db, set_id, "", "bbs_comment", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if bbs_comment_closed(db, set_id, set_code) {
		return_data["response"] = "error"
		return_data["data"] = "comment_closed"
		return return_data
	}

	data = strings.ReplaceAll(data, "\r", "")
	if data == "" {
		return_data["response"] = "error"
		return_data["data"] = "empty data"
		return return_data
	}
	if !tool.Do_bbs_max_length_check(db, config, data) {
		return_data["response"] = "error"
		return_data["data"] = "bbs overflow max length"
		return return_data
	}
	if !tool.Do_edit_filter(db, config, "", data) {
		return_data["response"] = "error"
		return_data["data"] = "edit filter (content)"
		return return_data
	}

	parent_id := set_id + "-" + set_code
	parent_user := ""
	if bbs_type != "thread" {
		var ok bool
		parent_id, parent_user, ok = bbs_comment_parent(db, set_id, set_code, comment_select)
		if !ok {
			return_data["response"] = "not exist"
			return_data["data"] = "comment"
			return return_data
		}
	}

	last_code := ""
	tool.QueryRow_DB(
		db,
		"select set_code from bbs_data where set_name = 'comment' and set_id = ? order by set_code + 0 desc limit 1",
		[]any{&last_code},
		parent_id,
	)
	comment_code := strconv.Itoa(tool.Str_to_int(last_code) + 1)
	date := tool.Get_time()

	insert_data := [][]string{
		{"comment", data},
		{"comment_date", date},
		{"comment_user_id", config.IP},
	}
	for _, value := range insert_data {
		tool.Exec_DB(
			db,
			"insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)",
			value[0],
			comment_code,
			parent_id,
			value[1],
		)
	}
	bbs_post_comment_count_update(db, set_id, set_code, 1)

	end_code := comment_code
	if bbs_type != "thread" && comment_select != "" && comment_select != "0" {
		end_code = comment_select + "-" + comment_code
	}
	alarm := "BBS <a href=\"/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code) + "#" + tool.Url_parser(end_code) + "\">" + tool.HTML_escape(bbs_name) + " - " + tool.HTML_escape(title) + "#" + tool.HTML_escape(end_code) + "</a>"
	tool.Send_alarm(db, config.IP, post_user, alarm)
	if parent_user != "" {
		tool.Send_alarm(db, config.IP, parent_user, alarm)
	}
	topic_reference_notify(db, config, data, end_code, set_code, set_id, bbs_name, title, "bbs")

	return_data["response"] = "ok"
	return_data["data"] = end_code
	return return_data
}
