package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_bbs_in_w_tool_post(config tool.Config, set_id string, set_code string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	title := ""
	if !tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'title' and set_id = ? and set_code = ?",
		[]any{&title},
		set_id,
		set_code,
	) {
		return tool.Get_redirect("/bbs/main")
	}

	switch values.Get("action") {
	case "comment_close":
		api_data := Api_bbs_w_comment_close(config, set_id, set_code, values.Get("comment_closed") == "1")
		if api_data["response"] == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if api_data["response"] != "ok" {
			return tool.Get_redirect("/bbs/main")
		}
	case "comment_delete":
		for _, comment_code := range values["comment_code"] {
			if !bbs_comment_code_regex.MatchString(comment_code) {
				continue
			}
			api_data := Api_bbs_w_comment_one_delete(config, set_id, set_code+"-"+comment_code)
			if api_data["response"] == "require auth" {
				return tool.Get_error_page(db, config, "auth")
			}
		}
	}

	return tool.Get_redirect("/bbs/tool/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code))
}
