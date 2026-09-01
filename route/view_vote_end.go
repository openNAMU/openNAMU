package route

import (
	"opennamu/route/tool"
	"strconv"
)

func View_vote_end(config tool.Config, id string) string {
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
	end_date := tool.Get_vote_value(db, id, "end_date")
	body := `<a href="/vote/close/` + tool.Url_parser(id) + `">` + tool.Get_language(db, "close_vote", true) + `</a><h2>` + tool.HTML_escape(name) + `</h2>`
	if subject != "" {
		body += `<b>` + tool.HTML_escape(subject) + `</b><hr class="main_hr">`
	}
	if end_date != "" {
		body += `<span>~ ` + tool.HTML_escape(end_date) + `</span><hr class="main_hr">`
	}
	for index, option := range vote_options(data) {
		count := "0"
		count = tool.Get_vote_count(db, id, strconv.Itoa(index))
		body += `<h3>` + tool.HTML_escape(option) + `</h3><p>` + count + `</p>`
		if type_data == "open" || type_data == "close" {
			rows := tool.Get_vote_users(db, id, strconv.Itoa(index))
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
