package route

import (
	"net/url"
	"opennamu/route/tool"
	"strings"
)

func View_alarm(config tool.Config, user_name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if user_name == "" {
		user_name = config.IP
	}
	is_change := values != nil && (values.Get("all") != "" || values.Get("delete") != "")
	if !is_change && user_name != config.IP && !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if is_change {
		result := Api_alarm_delete_post(config, user_name, values.Get("delete"), values.Get("all") != "")
		if result["response"] == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if result["response"] != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
		return_path := "/alarm"
		if user_name != config.IP {
			return_path = "/alarm_user/" + tool.Url_parser(user_name)
		}
		return tool.Get_redirect(return_path)
	}

	delete_url := "/alarm/delete"
	if user_name != config.IP {
		delete_url = "/alarm_user/" + tool.Url_parser(user_name) + "/delete"
	}
	page := 1
	if values != nil {
		page = tool.Str_to_int(values.Get("num"))
		if page < 1 {
			page = 1
		}
	}
	offset := (page - 1) * 50
	rows := tool.Query_DB(db, "select id, data, date, readme from user_notice where name = ? order by date desc limit ?, 50", user_name, offset)
	body := `<ul>`
	row_count := 0
	for rows.Next() {
		id, data, date, readme := "", "", "", ""
		if rows.Scan(&id, &data, &date, &readme) != nil {
			continue
		}
		row_count++
		data_split := strings.Split(data, " | ")
		data_html := tool.IP_parser(db, data_split[0], config.IP)
		if len(data_split) > 1 {
			data_html += " | " + strings.Join(data_split[1:], " | ")
		}
		data_style := ""
		if readme == "1" {
			data_style = ` style="opacity: 0.75;"`
		}
		body += `<li` + data_style + `>` + data_html + ` | ` + tool.HTML_escape(date) + ` <a href="` + delete_url + `/` + tool.Url_parser(id) + `">(` + tool.Get_language(db, "delete", true) + `)</a></li>`
	}
	rows.Close()
	body += `</ul>`
	read_url := "/alarm/read"
	if user_name != config.IP {
		read_url = "/alarm_user/" + tool.Url_parser(user_name) + "/read"
	}
	read_form := `<form method="post" action="` + read_url + `"><button type="submit">` + tool.Get_language(db, "read_all", true) + `</button></form><hr class="main_hr">`
	if row_count > 0 {
		body = read_form + `<a href="` + delete_url + `">(` + tool.Get_language(db, "delete", true) + `)</a><hr class="main_hr">` + body
	} else {
		body = read_form + body
	}
	page_url := "/alarm/page/{}"
	if user_name != config.IP {
		page_url = "/alarm_user/" + tool.Url_parser(user_name) + "/page/{}"
	}
	body += tool.Get_page_control(db, page, row_count, 50, page_url)
	return user_form_page(db, config, tool.Get_language(db, "alarm", true), body)
}
