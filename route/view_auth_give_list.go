package route

import (
	"database/sql"
	"opennamu/route/tool"
)

func Get_ui_auth_give_list(db *sql.DB, config tool.Config, data_all [][]string) string {
	data_html := ""
	date_heading := ""

	for _, in_data := range data_all {
		target := in_data[0]
		group := in_data[1]
		target_type := in_data[2]
		blocker := in_data[3]

		target_ui := tool.HTML_escape(target)
		target_path := "/auth/give/"
		if target_type == "normal" {
			target_ui = tool.IP_parser(db, target, config.IP)
		} else {
			target_path += target_type + "/"
			target_ui += " (" + tool.Get_language(db, target_type, true) + ")"
		}

		can_give := tool.Check_acl(db, "", "", "give_auth", config.IP)
		if target_type == "regex" || target_type == "cidr" {
			can_give = tool.Check_acl(db, "", "", "give_range_auth", config.IP)
		}
		if can_give {
			target_ui = `<a href="` + target_path + tool.Url_parser(target) + `">` + target_ui + `</a>`
		}

		if in_data[8] == "1" {
			target_ui += " (" + tool.Get_language(db, "default", true) + ")"
		}

		left := target_ui + " ← " + tool.HTML_escape(group)

		blocker_ui := tool.Get_language(db, "system", true)
		if blocker != "" {
			blocker_ui = tool.IP_parser(db, blocker, config.IP)
		}

		date_ui, date_text, new_date_heading := tool.Get_date_list_ui(in_data[4], date_heading)
		date_heading = new_date_heading

		right := tool.Get_language(db, "grantor", true) + " : " + blocker_ui
		if date_text != "" {
			right += " | " + tool.Get_language(db, "granted_at", true) + " : " + date_text
		}

		end := in_data[5]
		if end == "" || end == "0" {
			end = tool.Get_language(db, "limitless", true)
		}
		right += "<br>" + tool.Get_language(db, "date", true) + " : " + tool.HTML_escape(end)

		bottom := ""
		if in_data[6] != "" {
			bottom = Get_safe_send_data(in_data[6])
		}

		data_html += date_ui + tool.Get_list_ui(left, right, bottom, "")
	}

	return data_html
}

func View_auth_give_list(config tool.Config, num string, include_default bool) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "give_range_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	data_html := `<a href="/auth/give_list/1">(` + tool.Get_language(db, "auth_give_list", true) + `)</a> `
	data_html += `<a href="/auth/give_list/all/1">(` + tool.Get_language(db, "include_default_auth", true) + `)</a> `
	data_html += `<a href="/auth/give">(` + tool.Get_language(db, "authorize", true) + `)</a> `
	data_html += `<a href="/auth/list">(` + tool.Get_language(db, "admin_group_list", true) + `)</a>`
	data_html += `<hr class="main_hr">`

	api_data := Api_list_auth_give(config, num, include_default)
	api_data_list, ok := api_data["data"].([][]string)
	if !ok {
		return tool.Get_error_page(db, config, "error")
	}
	data_html += Get_ui_auth_give_list(db, config, api_data_list)

	base_url := "/auth/give_list/{}"
	if include_default {
		base_url = "/auth/give_list/all/{}"
	}
	data_html += tool.Get_page_control(db, tool.Str_to_int(num), len(api_data_list), 50, base_url)

	sub := ""
	if include_default {
		sub = "(" + tool.Get_language(db, "include_default_auth", true) + ")"
	}

	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "auth_give_list", true),
		data_html,
		[]any{sub},
		[][]any{{"manager", tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}
