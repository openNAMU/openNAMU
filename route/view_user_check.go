package route

import (
	"opennamu/route/tool"
	"strconv"
)

func View_user_check(config tool.Config, name string, check_type string, page string, plus_name string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if name == "" {
		data := `<form method="post" action="/manager/3"><input name="name" placeholder="` + tool.Get_language(db, "user_name", true) + `"><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "go", true) + `</button></form>`
		return tool.Get_template(db, config, tool.Get_language(db, "check", true), data, []any{}, [][]any{{"manager", tool.Get_language(db, "return", true)}}, map[string]string{})
	}
	if !tool.Check_permission(db, "check", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	page_num, _ := strconv.Atoi(page)
	if page_num < 1 {
		page_num = 1
	}
	offset := (page_num - 1) * 50
	if check_type == "simple" {
		row_count := 0
		rows := tool.Get_ua_simple_rows(db, name, tool.IP_or_user(name), offset)
		defer rows.Close()
		data := `<a href="/list/user/check/` + tool.Url_parser(name) + `">(` + tool.Get_language(db, "check", true) + `)</a><ul>`
		for rows.Next() {
			other_name := ""
			if rows.Scan(&other_name) == nil {
				data += `<li><a href="/list/user/check/` + tool.Url_parser(other_name) + `/simple">` + tool.HTML_escape(other_name) + `</a></li>`
				row_count++
			}
		}
		data += `</ul>`
		data += tool.Get_page_control(
			db,
			page_num,
			row_count,
			50,
			"/list/user/check/"+tool.Url_parser(name)+"/simple/{}",
		)
		return tool.Get_template(db, config, name, data, []any{"(" + tool.Get_language(db, "check", true) + ")"}, [][]any{{"manager", tool.Get_language(db, "return", true)}}, map[string]string{})
	}

	rows := tool.Get_ua_rows(db, name, plus_name, tool.IP_or_user(name), tool.IP_or_user(plus_name), offset)
	defer rows.Close()

	data := ""
	if !tool.IP_or_user(name) {
		question, answer := tool.Get_user_approval_data(db, name)
		if question != "" && answer != "" {
			data += `<table id="main_table_set"><tr id="main_table_top_tr"><td>Q</td><td>` + tool.HTML_escape(question) + `</td><td>A</td><td>` + tool.HTML_escape(answer) + `</td></tr></table><hr class="main_hr">`
		}
	}
	if plus_name != "" && page_num == 1 {
		all_count := len(tool.Get_ua_distinct_ips(db, name, plus_name, tool.IP_or_user(name), tool.IP_or_user(plus_name)))
		name_count := len(tool.Get_ua_distinct_ips(db, name, "", tool.IP_or_user(name), false))
		plus_count := len(tool.Get_ua_distinct_ips(db, plus_name, "", tool.IP_or_user(plus_name), false))
		if name_count+plus_count != all_count {
			data += tool.Get_language(db, "same_ip_exist", true) + `<hr class="main_hr">`
		}
	}
	data += `<table id="main_table_set"><tr id="main_table_top_tr"><td>` + tool.Get_language(db, "name", true) + `</td><td>` + tool.Get_language(db, "ip", true) + `</td><td>` + tool.Get_language(db, "time", true) + `</td></tr>`
	row_count := 0
	for rows.Next() {
		user_name := ""
		user_ip := ""
		user_agent := ""
		today := ""
		if rows.Scan(&user_name, &user_ip, &user_agent, &today) != nil {
			continue
		}
		row_count++
		delete_url := `/list/user/check/delete/` + tool.Url_parser(user_name) + `/` + tool.Url_parser(user_ip) + `/` + tool.Url_parser(today) + `/` + func() string {
			if tool.IP_or_user(name) {
				return "1"
			}
			return "0"
		}()
		user_agent_html := `<br>`
		if user_agent != "" {
			user_agent_html = tool.HTML_escape(user_agent)
			if len(user_agent) > 300 {
				user_agent_html = `<details><summary>(300+)</summary>` + user_agent_html + `</details>`
			}
		}
		data += `<tr><td><a href="/list/user/check/` + tool.Url_parser(user_name) + `">` + tool.HTML_escape(user_name) + `</a> <a href="` + delete_url + `">(` + tool.Get_language(db, "delete", true) + `)</a></td><td><a href="/list/user/check/` + tool.Url_parser(user_ip) + `">` + tool.HTML_escape(user_ip) + `</a></td><td>` + tool.HTML_escape(today) + `</td></tr><tr><td colspan="3">` + user_agent_html + `</td></tr>`
	}
	data += `</table>`
	if plus_name != "" {
		data = `<a href="/list/user/check/` + tool.Url_parser(name) + `">` + tool.HTML_escape(name) + `</a> <a href="/list/user/check/` + tool.Url_parser(plus_name) + `">` + tool.HTML_escape(plus_name) + `</a><hr class="main_hr">` + data
	} else {
		data = `<a href="/manager/14/` + tool.Url_parser(name) + `">(` + tool.Get_language(db, "compare", true) + `)</a> <a href="/list/user/check/` + tool.Url_parser(name) + `/simple">(` + tool.Get_language(db, "check", true) + `)</a><hr class="main_hr">` + data
	}

	page_url := "/list/user/check/" + tool.Url_parser(name) + "/normal/{}"
	if plus_name != "" {
		page_url += "/" + tool.Url_parser(plus_name)
	}
	data += tool.Get_page_control(db, page_num, row_count, 50, page_url)

	return tool.Get_template(db, config, name, data, []any{"(" + tool.Get_language(db, "check", true) + ")"}, [][]any{{"manager", tool.Get_language(db, "return", true)}}, map[string]string{})
}
