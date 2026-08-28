package route

import (
	"database/sql"
	"opennamu/route/tool"
	"strings"
)

func Get_ui_recent_block(db *sql.DB, config tool.Config, data_all [][]string) string {
	data_html := ""
	give_auth := tool.Check_acl(db, "", "", "give_auth", config.IP)

	for_count := 1
	for _, in_data := range data_all {
		// for_count_str := strconv.Itoa(for_count)
		for_count += 1

		left := ""

		ip := in_data[1]

		target_path := "/auth/give/"
		switch in_data[7] {
		case "private":
			target_path += "private/"
			ip += " (" + tool.Get_language(db, "private", true) + ")"
		case "cidr":
			target_path += "cidr/"
			ip += " (" + tool.Get_language(db, "cidr", true) + ")"
		case "regex":
			target_path += "regex/"
			ip += " (" + tool.Get_language(db, "regex", true) + ")"
		}

		if give_auth {
			ip = `<a href="` + target_path + tool.Url_parser(in_data[1]) + `">` + ip + `</a>`
		}

		if in_data[8] == "1" {
			ip = `<s>` + ip + `</s>`
		}

		left += ip + ` ← ` + in_data[4]

		end := ""
		switch in_data[5] {
		case "release":
			end = tool.Get_language(db, "release", true)
		case "":
			end = tool.Get_language(db, "limitless", true)
		default:
			end = in_data[5]
		}

		right := end + "<br>" + in_data[6]

		bottom := ""
		if in_data[0] != "" {
			if in_data[0] == "edit filter" {
				bottom = `<a href="/edit_filter/` + tool.Url_parser(in_data[1]) + `">edit filter</a>`
			} else {
				bottom = Get_safe_send_data(tool.HTML_escape(in_data[0]))
			}
		}

		data_html += tool.Get_list_ui(left, right, bottom, "")
	}

	return data_html
}

func View_list_recent_block(config tool.Config, num string, set_type string, why string, user_name string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	data_html := ""
	sub := ""

	if set_type == "" {
		set_type = "all"
	}

	menu_option := [][]string{
		{"all", tool.Get_language(db, "all", true)},
		{"regex", tool.Get_language(db, "regex", true)},
		{"cidr", tool.Get_language(db, "cidr", true)},
		{"private", tool.Get_language(db, "private", true)},
		{"ongoing", tool.Get_language(db, "in_progress", true)},
	}
	for _, option := range menu_option {
		data_html += `<a href="/recent_block/` + option[0] + `/1">(` + option[1] + `)</a> `

		if option[0] == set_type && set_type != "all" {
			sub = "(" + option[1] + ")"
		}
	}

	menu_option = [][]string{
		{"/manager/11", tool.Get_language(db, "blocked", true)},
		{"/manager/12", tool.Get_language(db, "admin", true)},
		{"/manager/19", tool.Get_language(db, "why", true)},
	}
	for _, option := range menu_option {
		data_html += `<a href="` + option[0] + `">(` + option[1] + `)</a> `
	}

	login_menu_option := [][]string{
		{"normal", tool.Get_language(db, "normal", true)},
		{"login_able", tool.Get_language(db, "login_able", true)},
		{"login_able_and_regsiter_disable", tool.Get_language(db, "login_able_and_regsiter_disable", true)},
		{"edit_request_able", tool.Get_language(db, "edit_request_able", true)},
		{"completely_ban", tool.Get_language(db, "completely_ban", true)},
		{"dont_come_this_site", tool.Get_language(db, "dont_come_this_site", true)},
	}
	for _, option := range login_menu_option {
		data_html += `<a href="/recent_block/login/` + option[0] + `/1">(` + option[1] + `)</a> `

		if set_type == option[0] {
			sub = "(" + option[1] + ")"
		}
	}

	data_html += "<hr class=\"main_hr\">"

	api_data := Api_list_recent_block(config, num, set_type, why, user_name)
	api_data_list := api_data["data"].([][]string)

	data_html += Get_ui_recent_block(db, config, api_data_list)

	base_url := "/recent_block/" + tool.Url_parser(set_type)
	if set_type == "normal" || strings.HasPrefix(set_type, "login_") {
		base_url = "/recent_block/login/" + tool.Url_parser(set_type)
	} else if user_name != "" {
		base_url += "/" + tool.Url_parser(user_name)
	}

	base_url += "/{}"

	if why != "" {
		base_url += "/" + tool.Url_parser(why)
	}

	data_html += tool.Get_page_control(
		db,
		tool.Str_to_int(num),
		len(api_data_list),
		50,
		base_url,
	)

	out := tool.Get_template(
		db,
		config,
		tool.Get_language(db, "recent_ban", true),
		data_html,
		[]any{sub},
		[][]any{
			{"other", tool.Get_language(db, "return", true)},
		},
		map[string]string{},
	)

	return out
}
