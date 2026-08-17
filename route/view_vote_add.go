package route

import (
	"net/url"
	"opennamu/route/tool"
	"strconv"
	"strings"
)

func View_vote_add(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !tool.Check_acl(db, "", "", "vote", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if values != nil {
		data := values.Get("data")
		options := vote_options(data)
		if len(options) < 2 {
			return tool.Get_error_page(db, config, "error")
		}
		last_id := "0"
		tool.QueryRow_DB(db, "select id from vote where type != 'option' order by id + 0 desc limit 1", []any{&last_id})
		id := strconv.Itoa(tool.Str_to_int(last_id) + 1)
		type_data := "n_open"
		if values.Get("open_select") == "Y" {
			type_data = "open"
		}
		tool.Exec_DB(db, "insert into vote (name, id, subject, data, user, type, acl) values (?, ?, ?, ?, '', ?, ?)", values.Get("name"), id, values.Get("subject"), strings.Join(options, "\n"), type_data, values.Get("acl_select"))
		tool.Exec_DB(db, "insert into vote (name, id, subject, data, user, type, acl) values ('open_user', ?, '', ?, '', 'option', '')", id, config.IP)
		if values.Get("limitless") == "" && values.Get("date") != "" {
			tool.Exec_DB(db, "insert into vote (name, id, subject, data, user, type, acl) values ('end_date', ?, '', ?, '', 'option', '')", id, values.Get("date"))
		}
		return tool.Get_redirect("/vote")
	}
	body := `<form method="post"><input name="name" placeholder="` + tool.Get_language(db, "name", true) + `"><hr class="main_hr"><textarea name="subject"></textarea><hr class="main_hr"><textarea name="data" placeholder="1 line 1 option"></textarea><hr class="main_hr"><label><input type="checkbox" name="open_select" value="Y"> open vote</label><hr class="main_hr"><input type="date" name="date"><label><input type="checkbox" name="limitless" value="Y"> limitless</label><hr class="main_hr"><input name="acl_select"><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "send", true) + `</button></form>`
	return vote_page(db, config, tool.Get_language(db, "add_vote", true), body)
}
