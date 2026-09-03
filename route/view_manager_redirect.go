package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func manager_redirect_list(db *sql.DB) map[int][]string {
	lang := func(name string) string {
		return tool.Get_language(db, name, true)
	}

	return map[int][]string{
		0:  {lang("document_name"), "/acl", lang("document_setting")},
		1:  {"", "/list/user/check", lang("check")},
		2:  {lang("file_name"), "/filter/file_filter/add", lang("file_filter_add")},
		3:  {"", "/auth/give", lang("authorize")},
		4:  {"", "/user", lang("user_tool")},
		6:  {lang("name"), "/auth/list/add", lang("add_admin_group")},
		7:  {lang("name"), "/filter/edit_filter/add", lang("edit_filter_add")},
		8:  {lang("document_name"), "/search", lang("search")},
		9:  {"", "/recent_block/user", lang("blocked_user")},
		10: {"", "/recent_block/admin", lang("blocked_admin")},
		11: {lang("document_name"), "/watch_list", lang("add_watchlist")},
		12: {lang("compare_target"), "/list/user/check", lang("compare_target")},
		13: {lang("document_name"), "/edit", lang("load")},
		14: {lang("document_name"), "/star_doc", lang("add_star_doc")},
		16: {"", "/auth/give/fix", lang("user_fix")},
		17: {lang("search"), "/auth/give_list", lang("search")},
	}
}

func View_manager_redirect(config tool.Config, num int, add_2 string, name string, post bool) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if num == 1 {
		if !tool.Check_permission(db, "site_view", config.IP) {
			return tool.Get_error_page(db, config, "auth")
		}
		return tool.Get_redirect("/manager")
	}

	if !tool.Check_permission(db, "manager_view", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	num -= 2
	item, ok := manager_redirect_list(db)[num]
	if !ok {
		return tool.Get_redirect("/manager")
	}

	if post {
		if name == "" {
			name = "test"
		}

		if add_2 != "" && num == 12 {
			return tool.Get_redirect(item[1] + "/" + tool.Url_parser(add_2) + "/normal/1/" + tool.Url_parser(name))
		}
		if add_2 != "" {
			return tool.Get_redirect("/edit_from_load/" + tool.Base64_encode(name) + "/" + tool.Url_parser(add_2))
		}

		return tool.Get_redirect(item[1] + "/" + tool.Url_parser(name))
	}

	placeholder := item[0]
	if placeholder == "" {
		placeholder = tool.Get_language(db, "user_name", true)
	}

	data := `<form method="post">`
	data += `<input placeholder="` + tool.HTML_escape(placeholder) + `" name="name" type="text">`
	data += `<hr class="main_hr"><button type="submit">` + tool.Get_language(db, "go", true) + `</button></form>`

	return tool.Get_template(
		db,
		config,
		item[2],
		data,
		[]any{},
		[][]any{{"manager", tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}
