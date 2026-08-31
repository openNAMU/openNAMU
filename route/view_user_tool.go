package route

import (
	"opennamu/route/tool"
)

func View_user_tool(config tool.Config, user_name string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if user_name == "" {
		user_name = config.IP
	}

	menu := tool.IP_menu(db, user_name, config.IP, "")
	body := ""
	for title, items := range menu {
		body += "<h2>" + tool.HTML_escape(title) + "</h2><ul>"
		for _, item := range items {
			if len(item) < 2 {
				continue
			}
			body += `<li><a href="` + tool.HTML_escape(item[0]) + `">` + tool.HTML_escape(item[1]) + `</a></li>`
		}
		body += "</ul>"
	}

	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "user_tool", true),
		body,
		[]any{},
		[][]any{{"user/" + tool.Url_parser(user_name), tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}
