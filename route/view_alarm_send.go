package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_alarm_send(config tool.Config, user_name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if values != nil {
		result := Api_alarm_send_post(config, user_name, values.Get("data"))
		switch result["response"] {
		case "require auth":
			return tool.Get_error_page(db, config, "auth")
		case "ok":
			body := "<p>" + tool.Get_language(db, "alarm_send_done", true) + "</p>"
			return tool.Get_template(
				db,
				config,
				tool.Get_language(db, "alarm_send", true),
				body,
				[]any{},
				[][]any{{"user/" + tool.Url_parser(user_name), tool.Get_language(db, "return", true)}},
				map[string]string{},
			)
		default:
			return tool.Get_error_page(db, config, "alarm_send_invalid")
		}
	}

	if !tool.Check_permission(db, "alarm_send", config.IP) ||
		user_name == "" || user_name == config.IP || tool.IP_or_user(user_name) ||
		!tool.Get_user_set_exists(db, user_name, "pw") {
		return tool.Get_error_page(db, config, "auth")
	}

	action := "/alarm_user/" + tool.Url_parser(user_name) + "/send"
	body := "<form method=\"post\" action=\"" + action + "\">" +
		"<p>" + tool.Get_language(db, "user_name", true) + ": " + tool.HTML_escape(user_name) + "</p>" +
		"<textarea name=\"data\" rows=\"8\" required></textarea><hr class=\"main_hr\">" +
		"<button type=\"submit\">" + tool.Get_language(db, "send", true) + "</button></form>"
	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "alarm_send", true),
		body,
		[]any{},
		[][]any{{"user/" + tool.Url_parser(user_name), tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}
