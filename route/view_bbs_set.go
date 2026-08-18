package route

import (
	"database/sql"
	"net/url"
	"strings"

	"opennamu/route/tool"
	"opennamu/route/tool/markup"
)

var bbs_set_fields = []string{
	"bbs_view_acl",
	"bbs_acl",
	"bbs_edit_acl",
	"bbs_comment_acl",
}

func bbs_set_value(db *sql.DB, set_id string, set_name string) string {
	value := ""
	tool.QueryRow_DB(
		db,
		"select set_data from bbs_set where set_id = ? and set_name = ? order by set_code + 0 desc limit 1",
		[]any{&value},
		set_id,
		set_name,
	)
	return value
}

func bbs_prefix_list(db *sql.DB, set_id string) []string {
	value := strings.ReplaceAll(bbs_set_value(db, set_id, "bbs_prefix"), "\r", "")
	prefix_list := []string{}
	for _, prefix := range strings.Split(value, "\n") {
		prefix = strings.TrimSpace(prefix)
		if prefix != "" && !tool.Arr_in_str(prefix_list, prefix) {
			prefix_list = append(prefix_list, prefix)
		}
	}
	return prefix_list
}

func bbs_prefix_check(db *sql.DB, set_id string, prefix string) string {
	prefix = strings.TrimSpace(prefix)
	if tool.Arr_in_str(bbs_prefix_list(db, set_id), prefix) {
		return prefix
	}
	return ""
}

func bbs_tag_list(data string) []string {
	data = strings.ReplaceAll(data, "\r", "")
	data = strings.ReplaceAll(data, "\n", ",")

	tag_list := []string{}
	for _, tag := range strings.Split(data, ",") {
		tag = strings.TrimSpace(tag)
		if tag != "" && !tool.Arr_in_str(tag_list, tag) {
			tag_list = append(tag_list, tag)
		}
	}

	return tag_list
}

const bbs_title_max_length = 128
const bbs_tag_max_length = 64

func bbs_set_select(db *sql.DB, name string, selected string, values []string) string {
	data := `<select name="` + tool.HTML_escape(name) + `">`
	for _, value := range values {
		choice := ""
		if value == selected {
			choice = ` selected`
		}
		label := value
		if label == "" {
			label = tool.Get_language(db, "normal", true)
		}
		data += `<option value="` + tool.HTML_escape(value) + `"` + choice + `>` + tool.HTML_escape(label) + `</option>`
	}
	return data + `</select>`
}

func View_bbs_set(config tool.Config, set_id string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	bbs_name := bbs_set_value(db, set_id, "bbs_name")
	if bbs_name == "" {
		return tool.Get_redirect("/bbs/main")
	}
	if values != nil && !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if values != nil {
		for _, field := range bbs_set_fields {
			Api_bbs_w_set_put(config, set_id, field, values.Get(field), "")
		}
		Api_bbs_w_set_put(config, set_id, "bbs_markup", values.Get("bbs_markup"), "")
		Api_bbs_w_set_put(config, set_id, "bbs_name", values.Get("bbs_name"), "")
		Api_bbs_w_set_put(config, set_id, "bbs_prefix", values.Get("bbs_prefix"), "")
		return tool.Get_redirect("/bbs/set/" + tool.Url_parser(set_id))
	}

	acl_values := tool.List_acl("normal")
	data := `<form method="post">`
	for _, field := range bbs_set_fields {
		data += `<h3>` + tool.Get_language(db, field, true) + `</h3>`
		data += bbs_set_select(db, field, bbs_set_value(db, set_id, field), acl_values)
		data += `<hr class="main_hr">`
	}

	markup_values := markup.List_markup()
	data += `<h3>` + tool.Get_language(db, "markup", true) + `</h3>`
	data += bbs_set_select(db, "bbs_markup", bbs_set_value(db, set_id, "bbs_markup"), markup_values)
	data += `<hr class="main_hr"><h3>` + tool.Get_language(db, "bbs_name", true) + `</h3>`
	data += `<input name="bbs_name" value="` + tool.HTML_escape(bbs_name) + `"><hr class="main_hr">`
	data += "<h3>" + tool.Get_language(db, "bbs_prefix", true) + "</h3>"
	data += "<textarea class=\"opennamu_textarea_100\" name=\"bbs_prefix\">" + tool.HTML_escape(bbs_set_value(db, set_id, "bbs_prefix")) + "</textarea><hr class=\"main_hr\">"
	data += `<button type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`

	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "bbs_set", true),
		data,
		[]any{"(" + bbs_name + ")"},
		[][]any{{"bbs/in/" + tool.Url_parser(set_id), tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}
