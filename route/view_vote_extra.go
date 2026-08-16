package route

import (
	"database/sql"
	"net/url"
	"strconv"
	"strings"

	"opennamu/route/tool"
)

func vote_page(db *sql.DB, config tool.Config, title string, body string) string {
	return tool.Get_template(db, config, title, body, []any{}, [][]any{{"vote", tool.Get_language(db, "return", true)}}, map[string]string{})
}

func vote_options(data string) []string {
	lines := strings.Split(strings.ReplaceAll(data, "\r", ""), "\n")
	options := []string{}
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			options = append(options, line)
		}
	}
	return options
}

func View_vote_select(config tool.Config, id string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	name, subject, data, type_data, end_date := "", "", "", "", ""
	if !tool.QueryRow_DB(db, "select name, subject, data, type from vote where id = ? and user = ''", []any{&name, &subject, &data, &type_data}, id) {
		return tool.Get_redirect("/vote")
	}
	if type_data == "close" || type_data == "n_close" || !tool.Check_acl(db, "", id, "vote", config.IP) {
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
		choice, err := strconv.Atoi(values.Get("vote_data"))
		if err != nil || choice < 0 || choice >= len(options) {
			return tool.Get_redirect("/vote/" + tool.Url_parser(id))
		}
		tool.Exec_DB(db, "insert into vote (name, id, subject, data, user, type) values ('', ?, '', ?, ?, 'select')", id, strconv.Itoa(choice), config.IP)
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

func View_vote_end(config tool.Config, id string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	name, subject, data, type_data := "", "", "", ""
	if !tool.QueryRow_DB(db, "select name, subject, data, type from vote where id = ? and user = ''", []any{&name, &subject, &data, &type_data}, id) {
		return tool.Get_redirect("/vote")
	}
	end_date := ""
	tool.QueryRow_DB(db, "select data from vote where id = ? and name = 'end_date' and type = 'option'", []any{&end_date}, id)
	body := `<a href="/vote/close/` + tool.Url_parser(id) + `">` + tool.Get_language(db, "close_vote", true) + `</a><h2>` + tool.HTML_escape(name) + `</h2>`
	if subject != "" {
		body += `<b>` + tool.HTML_escape(subject) + `</b><hr class="main_hr">`
	}
	if end_date != "" {
		body += `<span>~ ` + tool.HTML_escape(end_date) + `</span><hr class="main_hr">`
	}
	for index, option := range vote_options(data) {
		count := "0"
		tool.QueryRow_DB(db, "select count(*) from vote where id = ? and user != '' and data = ?", []any{&count}, id, strconv.Itoa(index))
		body += `<h3>` + tool.HTML_escape(option) + `</h3><p>` + count + `</p>`
		if type_data == "open" || type_data == "close" {
			rows := tool.Query_DB(db, "select user from vote where id = ? and user != '' and data = ?", id, strconv.Itoa(index))
			for rows.Next() {
				user := ""
				if rows.Scan(&user) == nil {
					body += `<div>` + tool.IP_parser(db, user, config.IP) + `</div>`
				}
			}
			rows.Close()
		}
	}
	return vote_page(db, config, tool.Get_language(db, "result_vote", true), body)
}

func View_vote_close(config tool.Config, id string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !tool.Check_acl(db, "", "", "vote", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	type_data := ""
	if !tool.QueryRow_DB(db, "select type from vote where id = ? and user = ''", []any{&type_data}, id) {
		return tool.Get_redirect("/vote")
	}
	owner := ""
	tool.QueryRow_DB(db, "select data from vote where id = ? and name = 'open_user' and type = 'option'", []any{&owner}, id)
	if owner != config.IP && !tool.Check_acl(db, "", "", "vote_auth", config.IP) {
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
	tool.Exec_DB(db, "update vote set type = ? where id = ? and user = ''", next, id)
	if next == "open" || next == "n_open" {
		tool.Exec_DB(db, "delete from vote where id = ? and name = 'end_date' and type = 'option'", id)
	}
	if next == "open" || next == "n_open" {
		return tool.Get_redirect("/vote/end/" + tool.Url_parser(id))
	}
	return tool.Get_redirect("/vote")
}
