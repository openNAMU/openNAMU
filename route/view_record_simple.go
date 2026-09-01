package route

import (
	"strings"

	"opennamu/route/tool"
)

func View_record_simple(config tool.Config, user_name string, record_type string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if user_name != "" && user_name != config.IP && !tool.Check_permission(db, "hidel", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if user_name == "" {
		user_name = config.IP
	}
	body := strings.Builder{}
	switch record_type {
	case "topic":
		rows := tool.Get_topic_record_rows(db, user_name, 0, false)
		for rows.Next() {
			code, data, date := "", "", ""
			if rows.Scan(&code, &data, &date) == nil {
				body.WriteString(tool.Get_list_ui(`<a href="/thread/`+tool.Url_parser(code)+`">`+tool.HTML_escape(code)+`</a>`, tool.HTML_escape(date), tool.HTML_escape(data), ""))
			}
		}
		rows.Close()
	default:
		rows := tool.Get_history_record_rows(db, user_name, "", 0, false)
		for rows.Next() {
			title, date, send := "", "", ""
			if rows.Scan(&title, &date, &send) == nil {
				body.WriteString(tool.Get_list_ui(`<a href="/w/`+tool.Url_parser(title)+`">`+tool.HTML_escape(title)+`</a>`, tool.HTML_escape(date), tool.HTML_escape(send), ""))
			}
		}
		rows.Close()
	}
	menu := [][]any{{"other", tool.Get_language(db, "return", true)}}
	if tool.Check_permission(db, "record_manage", config.IP) {
		menu = append(menu, []any{"record/reset/" + tool.Url_parser(user_name), tool.Get_language(db, "record_reset", true)})
	}
	return tool.Get_template(db, config, tool.Get_language(db, record_type+"_record", true), body.String(), []any{}, menu, map[string]string{})
}
