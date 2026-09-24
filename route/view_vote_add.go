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
	body := `<form method="post">
<label for="vote_name">` + tool.Get_language(db, "name", true) + `</label> <input id="vote_name" name="name" placeholder="` + tool.Get_language(db, "name", true) + `">
<hr class="main_hr"><div><label for="vote_subject">` + tool.Get_language(db, "explanation", true) + `</label></div><textarea id="vote_subject" name="subject"></textarea>
<hr class="main_hr"><div><label for="vote_options">` + tool.Get_language(db, "vote_options", true) + `</label></div><textarea id="vote_options" name="data" placeholder="` + tool.Get_language(db, "1_line_1_q", true) + `"></textarea>
<hr class="main_hr"><label><input type="checkbox" name="open_select" value="Y"> ` + tool.Get_language(db, "open_vote", true) + `</label>
<hr class="main_hr"><label for="vote_date">` + tool.Get_language(db, "date", true) + `</label> <input id="vote_date" type="date" name="date">
<hr class="main_hr"><label><input type="checkbox" name="limitless" value="Y"> ` + tool.Get_language(db, "limitless", true) + `</label>
<hr class="main_hr"><label for="acl_select">` + tool.Get_language(db, "vote_acl", true) + `</label> ` + tool.Build_select("acl_select", Acl_value_list(db, ""), "", tool.Get_language(db, "normal", true)) + `
<hr class="main_hr"><button type="submit">` + tool.Get_language(db, "send", true) + `</button></form>`
	return Vote_page(db, config, tool.Get_language(db, "add_vote", true), body)
}
