package route

import (
	"strings"

	"opennamu/route/tool"
)

func View_list_history(config tool.Config, doc_name string, set_type string, num string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "history_view", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	if !tool.Check_acl(db, doc_name, "", "render", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	sub := ""
	if set_type == "" {
		set_type = "normal"
	} else {
		sub = " (" + tool.Get_language(db, set_type, true) + ")"
	}

	data_html := ""

	menu_option := []string{"normal", "edit", "move", "delete", "revert", "r1", "setting"}
	for _, option := range menu_option {
		label := tool.Get_language(db, option, true)
		data_html += `<a href="/history_page/1/` + option + `/` + tool.Url_parser(doc_name) + `">(` + label + `)</a> `
	}

	api_data := Api_list_history(config, doc_name, set_type, num)
	api_data_list, _ := api_data["data"].([][]string)

	history_ui, select_ui := Get_ui_history(db, api_data_list)

	data_html += history_ui
	if len(api_data_list) == 0 {
		data_html += tool.Get_language(db, "data_missing", true)
	}

	revision_ids := []string{}
	for _, history_data := range api_data_list {
		if len(history_data) > 1 && history_data[1] != "" {
			revision_ids = append(revision_ids, history_data[0])
		}
	}

	if len(revision_ids) > 1 {
		previous_revision := revision_ids[1]
		after_revision := revision_ids[0]
		previous_select_ui := strings.Replace(
			select_ui,
			`value="`+previous_revision+`"`,
			`value="`+previous_revision+`" selected`,
			1,
		)
		after_select_ui := strings.Replace(
			select_ui,
			`value="`+after_revision+`"`,
			`value="`+after_revision+`" selected`,
			1,
		)

		data_html += `<hr class="main_hr">
		<form method="post">
			<div><label for="history_before">` + tool.Get_language(db, "before_revision", true) + `</label> <select id="history_before" name="b">` + previous_select_ui + `</select></div>
			<hr class="main_hr">
			<div><label for="history_after">` + tool.Get_language(db, "after_revision", true) + `</label> <select id="history_after" name="a">` + after_select_ui + `</select></div>
			<hr class="main_hr">
			<div><button type="submit">` + tool.Get_language(db, "compare", true) + `</button></div>
		</form>`
	}

	data_html += tool.Get_page_control(
		db,
		tool.Str_to_int(num),
		len(api_data_list),
		50,
		"/history_page/{}/"+set_type+"/"+tool.Url_parser(doc_name),
	)

	menu := [][]any{
		{"w/" + tool.Url_parser(doc_name), tool.Get_language(db, "return", true)},
	}
	if tool.Check_permission(db, "history_manage", config.IP) {
		menu = append(menu, []any{"history_add/" + tool.Url_parser(doc_name), tool.Get_language(db, "history_add", true)})
		menu = append(menu, []any{"history_reset/" + tool.Url_parser(doc_name), tool.Get_language(db, "history_reset", true)})
	}

	out := tool.Get_template(
		db,
		config,
		doc_name,
		data_html,
		[]any{"(" + tool.Get_language(db, "history", true) + ")" + sub},
		menu,
		map[string]string{},
	)

	return out
}
