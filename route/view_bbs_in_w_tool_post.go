package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_bbs_in_w_tool_post(config tool.Config, set_id string, set_code string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

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
		tool.Exec_DB(
			db,
			"delete from bbs_data where set_name = 'comment_close' and set_id = ? and set_code = ?",
			set_id,
			set_code,
		)
		if values.Get("comment_closed") == "1" {
			tool.Exec_DB(
				db,
				"insert into bbs_data (set_name, set_code, set_id, set_data) values ('comment_close', ?, ?, '1')",
				set_code,
				set_id,
			)
		}
	case "comment_delete":
		for _, comment_code := range values["comment_code"] {
			if !bbs_comment_code_regex.MatchString(comment_code) {
				continue
			}
			Api_bbs_w_comment_one_delete(config, set_id, set_code+"-"+comment_code)
		}
	}

	return tool.Get_redirect("/bbs/tool/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code))
}
