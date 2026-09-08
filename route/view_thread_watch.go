package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_thread_watch(config tool.Config, topic_num string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	title := ""
	sub := ""
	if !tool.QueryRow_DB(db, "select title, sub from rd where code = ?", []any{&title, &sub}, topic_num) {
		return tool.Get_redirect("/")
	}

	if values != nil {
		api_data := Api_thread_watch_post(config, topic_num)
		response, _ := api_data["response"].(string)
		if response == "require auth" {
			return tool.Get_redirect("/login")
		}
		if response != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
		return tool.Get_redirect("/thread/" + tool.Url_parser(topic_num))
	}

	data := `<form method="post" action="/thread_watch/` + tool.Url_parser(topic_num) + `">
            <button id="opennamu_save_button" type="submit">` + tool.Get_language(db, "send", true) + `</button>
        </form>`

	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "thread_watchlist", true),
		data,
		[]any{"(" + tool.HTML_escape(sub) + ")"},
		[][]any{
			{"thread/" + tool.Url_parser(topic_num), tool.Get_language(db, "return", true)},
		},
		map[string]string{"title": title},
	)
}
