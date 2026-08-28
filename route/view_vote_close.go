package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_vote_close(config tool.Config, id string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if values == nil && !tool.Check_permission(db, "vote", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	type_data := ""
	if !tool.QueryRow_DB(db, "select type from vote where id = ? and user = ''", []any{&type_data}, id) {
		return tool.Get_redirect("/vote")
	}
	owner := ""
	tool.QueryRow_DB(db, "select data from vote where id = ? and name = 'open_user' and type = 'option'", []any{&owner}, id)
	if values == nil && owner != config.IP && !tool.Check_permission(db, "vote_manage", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	next := "n_close"
	if type_data == "n_close" {
		next = "n_open"
	} else if type_data == "close" {
		next = "open"
	} else if type_data == "open" {
		next = "close"
	}
	if values == nil {
		action := "close_vote"
		if next == "open" || next == "n_open" {
			action = "open_vote"
		}
		body := `<form method="post"><button type="submit">` + tool.Get_language(db, action, true) + `</button></form>`
		return vote_page(db, config, tool.Get_language(db, action, true), body)
	}
	api_data := Api_vote_close_post(config, id)
	if api_data["response"] == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}
	if api_data["response"] != "ok" {
		return tool.Get_redirect("/vote")
	}
	next, _ = api_data["data"].(string)
	if next == "open" || next == "n_open" {
		return tool.Get_redirect("/vote/end/" + tool.Url_parser(id))
	}
	return tool.Get_redirect("/vote")
}
