package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_user_check_delete(config tool.Config, user_name string, user_ip string, today string, return_type string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if values == nil && !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if values != nil {
		api_data := Api_user_check_delete_post(config, user_name, user_ip, today)
		response, _ := api_data["response"].(string)
		if response != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
		return tool.Get_redirect("/list/user/check/" + tool.Url_parser(func() string {
			if return_type == "0" {
				return user_name
			}
			return user_ip
		}()))
	}
	data := tool.Get_language(db, "name", true) + " : " + tool.HTML_escape(user_name) + `<hr class="main_hr">` + tool.Get_language(db, "ip", true) + " : " + tool.HTML_escape(user_ip) + `<hr class="main_hr">` + tool.Get_language(db, "time", true) + " : " + tool.HTML_escape(today) + `<hr class="main_hr"><form method="post"><button type="submit">` + tool.Get_language(db, "delete", true) + `</button></form>`
	return tool.Get_template(db, config, tool.Get_language(db, "check", true), data, []any{}, [][]any{{"list/user/check/" + tool.Url_parser(user_name), tool.Get_language(db, "return", true)}}, map[string]string{})
}
