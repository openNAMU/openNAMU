package route

import (
	"database/sql"
	"net/url"

	"opennamu/route/tool"
)

func thread_user_render(db *sql.DB, config tool.Config, ip string) string {
	return tool.IP_parser(db, ip, config.IP)
}

func thread_comments(db *sql.DB, config tool.Config, topic_num string) string {
	rows := tool.Query_DB(
		db,
		"select id, data, date, ip, block, top from topic where code = ? order by id + 0 asc",
		topic_num,
	)
	defer rows.Close()

	data_html := ""
	shortcut := `<div class="opennamu_thread_shortcut" id="thread_shortcut">`
	admin_auth := tool.Check_permission(db, "thread_manage", config.IP)

	for rows.Next() {
		var id string
		var data string
		var date_value string
		var ip string
		var block string
		var top string
		if rows.Scan(&id, &data, &date_value, &ip, &block, &top) != nil {
			continue
		}

		shortcut += `<a href="#` + tool.Url_parser(id) + `">#` + tool.HTML_escape(id) + `</a> `
		if block == "O" && !admin_auth {
			data = ""
		}

		color := "default"
		if id == "1" {
			color = "red"
		} else if ip == config.IP {
			color = "green"
		}
		date := `<a href="/thread/` + tool.Url_parser(topic_num) + `/comment/` + tool.Url_parser(id) + `/tool">(` + tool.Get_language(db, "tool", true) + `)</a> <a href="/thread/` + tool.Url_parser(topic_num) + `/comment/` + tool.Url_parser(id) + `/raw">(` + tool.Get_language(db, "raw", true) + `)</a> ` + tool.HTML_escape(date_value)
		data_html += get_thread_ui(db, thread_user_render(db, config, ip), date, data, id, color, block, "", topic_num, config)
	}

	shortcut += `</div>`
	return shortcut + data_html
}

func View_thread(config tool.Config, topic_num string, doc_name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	name := doc_name
	sub := ""
	if topic_num == "0" {
		if name == "" {
			name = "Test"
		}
		sub = tool.Get_language(db, "make_new_topic", true)
		if values == nil && (!tool.Check_acl(db, name, "0", "topic_view", config.IP) || !tool.Check_permission(db, "discuss_make_new_thread", config.IP)) {
			return tool.Get_error_page(db, config, "auth")
		}
	} else {
		if !tool.QueryRow_DB(db, "select title, sub from rd where code = ?", []any{&name, &sub}, topic_num) {
			return tool.Get_redirect("/")
		}
		if values == nil && !tool.Check_acl(db, name, topic_num, "topic_view", config.IP) {
			return tool.Get_error_page(db, config, "auth")
		}
	}

	if values != nil {
		api_data := Api_thread_post(config, topic_num, name, values.Get("content"), values.Get("topic"), values.Get("title"))
		response, _ := api_data["response"].(string)
		if response == "empty data" {
			if topic_num == "0" {
				return tool.Get_redirect("/thread/0/" + tool.Url_parser(name))
			}
			return tool.Get_redirect("/thread/" + tool.Url_parser(topic_num))
		}
		if response == "not exist" {
			return tool.Get_redirect("/")
		}
		if response == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if response != "ok" {
			error_name, _ := api_data["data"].(string)
			if error_name == "" {
				error_name = "error"
			}
			return tool.Get_error_page(db, config, error_name)
		}
		posted_topic_num, topic_ok := api_data["topic_num"].(string)
		comment_num, comment_ok := api_data["comment_num"].(string)
		if !topic_ok || !comment_ok {
			return tool.Get_error_page(db, config, "error")
		}
		return tool.Get_redirect("/thread/" + tool.Url_parser(posted_topic_num) + "#" + tool.Url_parser(comment_num))
	}

	data_html := get_render_setting_css(db, config) + `<style id="opennamu_list_hidden_style">.opennamu_list_hidden { display: none; }</style>`
	data_html += `<label><input type="checkbox" onclick="opennamu_list_hidden_remove();" checked> ` + tool.Get_language(db, "remove_hidden", true) + `</label><hr class="main_hr">`
	if topic_num != "0" {
		data_html += thread_comments(db, config, topic_num)
	}

	data_html += `<h2>` + tool.HTML_escape(sub) + `</h2>`
	can_post := topic_num == "0" || tool.Check_acl(db, name, topic_num, "topic", config.IP)
	if can_post {
		path := "/thread/" + tool.Url_parser(topic_num)
		if topic_num == "0" {
			path = "/thread/0/" + tool.Url_parser(doc_name)
		}
		data_html += `<form method="post" action="` + path + `">`
		if topic_num == "0" {
			data_html += `<input name="topic" value="` + tool.HTML_escape(name) + `" placeholder="` + tool.Get_language(db, "document_name", true) + `"><hr class="main_hr">`
			data_html += `<input name="title" value="" placeholder="` + tool.Get_language(db, "discussion_name", true) + `"><hr class="main_hr">`
		}
		data_html += tool.Get_editor_ui(db, config, "", "thread", "", "") + `</form>`
	}

	menu := [][]any{}
	if doc_name != "" {
		menu = append(menu, []any{"topic/" + tool.Url_parser(doc_name), tool.Get_language(db, "list", true)})
	}
	if topic_num != "0" {
		menu = append(menu, []any{"thread/" + tool.Url_parser(topic_num) + "/tool", tool.Get_language(db, "tool", true)})
	}

	return tool.Get_template(db, config, name, data_html, []any{"(" + tool.Get_language(db, "discussion", true) + ")"}, menu, map[string]string{})
}
