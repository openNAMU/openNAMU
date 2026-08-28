package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_vote_add(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if values == nil && !tool.Check_permission(db, "vote", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if values != nil {
		api_data := Api_vote_add_post(config, values.Encode())
		response, _ := api_data["response"].(string)
		if response == "ok" {
			return tool.Get_redirect("/vote")
		}
		if response == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		return tool.Get_error_page(db, config, "error")
	}
	body := `<form method="post"><input name="name" placeholder="` + tool.Get_language(db, "name", true) + `"><hr class="main_hr"><textarea name="subject"></textarea><hr class="main_hr"><textarea name="data" placeholder="1 line 1 option"></textarea><hr class="main_hr"><label><input type="checkbox" name="open_select" value="Y"> open vote</label><hr class="main_hr"><input type="date" name="date"><label><input type="checkbox" name="limitless" value="Y"> limitless</label><hr class="main_hr">` + bbs_set_select(db, "acl_select", "", acl_value_list(db, "")) + `<hr class="main_hr"><button type="submit">` + tool.Get_language(db, "send", true) + `</button></form>`
	return vote_page(db, config, tool.Get_language(db, "add_vote", true), body)
}
