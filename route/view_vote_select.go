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
	name, subject, data, type_data, end_date := "", "", "", "", ""
	if !tool.QueryRow_DB(db, "select name, subject, data, type from vote where id = ? and user = ''", []any{&name, &subject, &data, &type_data}, id) {
		return tool.Get_redirect("/vote")
	}
	if type_data == "close" || type_data == "n_close" || (values == nil && !tool.Check_acl(db, "", id, "vote", config.IP)) {
		return tool.Get_redirect("/vote/end/" + tool.Url_parser(id))
	}
	voted := ""
	if tool.QueryRow_DB(db, "select user from vote where id = ? and user = ?", []any{&voted}, id, config.IP) {
		return tool.Get_redirect("/vote/end/" + tool.Url_parser(id))
	}
	tool.QueryRow_DB(db, "select data from vote where id = ? and name = 'end_date' and type = 'option'", []any{&end_date}, id)
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
