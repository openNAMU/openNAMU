package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_bbs_watch(config tool.Config, set_id string, set_code string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	bbs_name := ""
	title := ""
	if !tool.QueryRow_DB(
		db,
		"select set_data from bbs_set where set_name = 'bbs_name' and set_id = ?",
		[]any{&bbs_name},
		set_id,
	) || !tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'title' and set_id = ? and set_code = ?",
		[]any{&title},
		set_id,
		set_code,
	) {
		return tool.Get_redirect("/bbs/main")
	}

	if values != nil {
		api_data := Api_bbs_watch_post(config, set_id, set_code)
		response, _ := api_data["response"].(string)
		if response == "require auth" {
			return tool.Get_redirect("/login")
		}
		if response != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
		return tool.Get_redirect("/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code))
	}

	data := `<form method="post" action="/bbs_watch/` + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code) + `">
            <button id="opennamu_save_button" type="submit">` + tool.Get_language(db, "send", true) + `</button>
        </form>`

	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "bbs_watchlist", true),
		data,
		[]any{"(" + tool.HTML_escape(bbs_name) + " - " + tool.HTML_escape(title) + ")"},
		[][]any{
			{"bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code), tool.Get_language(db, "return", true)},
		},
		map[string]string{"title": title},
	)
}
