package route

import (
	"opennamu/route/tool"
	"strconv"
)

func View_user_check(config tool.Config, name string, check_type string, page string, plus_name string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if name == "" {
		data := `<form method="post" action="/manager/3"><label for="check_user_name">` + tool.Get_language(db, "user_name", true) + `</label><input id="check_user_name" name="name"><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "go", true) + `</button></form>`
		return tool.Get_template(db, config, tool.Get_language(db, "check", true), data, []any{}, [][]any{{"manager", tool.Get_language(db, "return", true)}}, map[string]string{})
	}

	page_num, _ := strconv.Atoi(page)
	if page_num < 1 {
		page_num = 1
	}

	result := Api_user_check_rows(db, config.IP, name, check_type, page_num, plus_name)

	response, ok := result["response"].(string)
	if !ok || response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}

	rows := result["rows"].([]map[string]string)

	simple, _ := result["simple"].(bool)
	if simple {
		data := `<a href="/list/user/check/` + tool.Url_parser(name) + `">(` + tool.Get_language(db, "check", true) + `)</a><ul>`
		for _, row := range rows {
			if row["masked"] == "true" {
				data += `<li>` + tool.Get_language(db, "hidden_value", true) + `</li>`
			} else {
				other_name := row["other_name"]
				data += `<li><a href="/list/user/check/` + tool.Url_parser(other_name) + `/simple">` + tool.HTML_escape(other_name) + `</a></li>`
			}
		}
		data += `</ul>`
		data += tool.Get_page_control(db, page_num, len(rows), ua_page_size, "/list/user/check/"+tool.Url_parser(name)+"/simple/{}")
		return tool.Get_template(db, config, name, data, []any{"(" + tool.Get_language(db, "check", true) + ")"}, [][]any{{"manager", tool.Get_language(db, "return", true)}}, map[string]string{})
	}

	same_ip_exist, _ := result["same_ip_exist"].(bool)
	show_qa, _ := result["show_qa"].(bool)
	qa_question, _ := result["qa_question"].(string)
	qa_answer, _ := result["qa_answer"].(string)

	data := ""
	if show_qa {
		qa_data := `<div>Q: ` + tool.HTML_escape(qa_question) + `</div><div>A: ` + tool.HTML_escape(qa_answer) + `</div>`
		data += tool.Get_list_ui(qa_data, "", "", "")
	}
	if same_ip_exist {
		data += tool.Get_language(db, "same_ip_exist", true) + `<hr class="main_hr">`
	}

	delete_flag := "0"
	if tool.IP_or_user(name) {
		delete_flag = "1"
	}
	for _, row := range rows {
		user_name := row["name"]
		user_ip := row["ip"]
		user_agent := row["ua"]
		today := row["today"]
		masked := row["masked"] == "true"

		left := `<a href="/list/user/check/` + tool.Url_parser(user_name) + `">` + tool.HTML_escape(user_name) + `</a>`
		if !masked {
			delete_url := `/list/user/check/delete/` + tool.Url_parser(user_name) + `/` + tool.Url_parser(user_ip) + `/` + tool.Url_parser(today) + `/` + delete_flag
			left += ` <a href="` + delete_url + `">(` + tool.Get_language(db, "delete", true) + `)</a>`
		}

		right := `<strong>` + tool.Get_language(db, "time", true) + `:</strong> ` + tool.HTML_escape(today)

		var bottom string
		if masked {
			bottom = `<div><strong>` + tool.Get_language(db, "ip", true) + `:</strong> ` + tool.Get_language(db, "hidden_value", true) + `</div>`
		} else {
			bottom = `<div><strong>` + tool.Get_language(db, "ip", true) + `:</strong> <a href="/list/user/check/` + tool.Url_parser(user_ip) + `">` + tool.HTML_escape(user_ip) + `</a></div>`
			if user_agent != "" {
				user_agent_html := tool.HTML_escape(user_agent)
				if len(user_agent) > 300 {
					user_agent_html = `<details><summary>(300+)</summary>` + user_agent_html + `</details>`
				}
				bottom += user_agent_html
			}
		}

		data += tool.Get_list_ui(left, right, bottom, "")
	}

	if plus_name != "" {
		data = `(<a href="/list/user/check/` + tool.Url_parser(name) + `">` + tool.HTML_escape(name) + `</a>)&nbsp;(<a href="/list/user/check/` + tool.Url_parser(plus_name) + `">` + tool.HTML_escape(plus_name) + `</a>)<hr class="main_hr">` + data
	} else {
		data = `<a href="/manager/14/` + tool.Url_parser(name) + `">(` + tool.Get_language(db, "compare", true) + `)</a> <a href="/list/user/check/` + tool.Url_parser(name) + `/simple">(` + tool.Get_language(db, "check", true) + `)</a><hr class="main_hr">` + data
	}

	page_url := "/list/user/check/" + tool.Url_parser(name) + "/normal/{}"
	if plus_name != "" {
		page_url += "/" + tool.Url_parser(plus_name)
	}
	data += tool.Get_page_control(db, page_num, len(rows), ua_page_size, page_url)

	return tool.Get_template(db, config, name, data, []any{"(" + tool.Get_language(db, "check", true) + ")"}, [][]any{{"manager", tool.Get_language(db, "return", true)}}, map[string]string{})
}
