package route

import (
	"opennamu/route/tool"
	"strconv"
)

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
