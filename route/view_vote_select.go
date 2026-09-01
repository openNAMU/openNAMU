package route

import (
	"net/url"
	"opennamu/route/tool"
	"strconv"
	"strings"
)

func View_vote_select(config tool.Config, id string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	vote_data, vote_exists := tool.Get_vote_data(db, id)
	if !vote_exists {
		return tool.Get_redirect("/vote")
	}
	name := vote_data["name"]
	subject := vote_data["subject"]
	data := vote_data["data"]
	type_data := vote_data["type"]
	end_date := ""
	if type_data == "close" || type_data == "n_close" || (values == nil && !tool.Check_acl(db, "", id, "vote", config.IP)) {
		return tool.Get_redirect("/vote/end/" + tool.Url_parser(id))
	}
	if tool.Get_vote_user_exists(db, id, config.IP) {
		return tool.Get_redirect("/vote/end/" + tool.Url_parser(id))
	}
	end_date = tool.Get_vote_value(db, id, "end_date")
	if end_date != "" && strings.HasPrefix(tool.Get_time(), end_date) == false && strings.Split(tool.Get_time(), " ")[0] > strings.Split(end_date, " ")[0] {
		return tool.Get_redirect("/vote/end/" + tool.Url_parser(id))
	}
	options := vote_options(data)
	if values != nil {
		api_data := Api_vote_select_post(config, id, values.Get("vote_data"))
		response, _ := api_data["response"].(string)
		if response == "not exist" {
			return tool.Get_redirect("/vote")
		}
		if response == "error" {
			return tool.Get_redirect("/vote/" + tool.Url_parser(id))
		}
		return tool.Get_redirect("/vote/end/" + tool.Url_parser(id))
	}
	body := `<h2>` + tool.HTML_escape(name) + `</h2>`
	if subject != "" {
		body += `<b>` + tool.HTML_escape(subject) + `</b><hr class="main_hr">`
	}
	if end_date != "" {
		body += `<span>~ ` + tool.HTML_escape(end_date) + `</span><hr class="main_hr">`
	}
	body += `<form method="post"><select name="vote_data">`
	for index, option := range options {
		body += `<option value="` + strconv.Itoa(index) + `">` + tool.HTML_escape(option) + `</option>`
	}
	body += `</select><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "send", true) + `</button></form>`
	return vote_page(db, config, tool.Get_language(db, "vote", true), body)
}
