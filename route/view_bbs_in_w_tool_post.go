package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_bbs_in_w_tool_post(config tool.Config, set_id string, set_code string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if _, allowed := Bbs_post_view_auth(db, set_id, set_code, config.IP); !allowed {
		return tool.Get_error_page(db, config, "auth")
	}
	_, title_exists := tool.Get_bbs_data_value(db, set_id, set_code, "title")
	if !title_exists {
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
	case "comment_delete_all":
		api_data := Api_bbs_w_comment_all_delete(config, set_id, set_code)
		if api_data["response"] == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
	case "excellent_exclusion":
		excluded_value := values.Get("excluded")
		if excluded_value != "0" && excluded_value != "1" {
			return tool.Get_error_page(db, config, "error")
		}
		api_data := Api_bbs_excellent_exclusion_set(config, set_id, set_code, excluded_value == "1")
		if api_data["response"] == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if api_data["response"] != "ok" {
			if api_data["response"] == "not exist" {
				return tool.Get_redirect("/bbs/main")
			}
			return tool.Get_error_page(db, config, "error")
		}
	case "slow_mode_set":
		api_data := Api_bbs_w_post_slow_mode_put(config, set_id, set_code, values.Get("slow_mode_seconds"))
		if api_data["response"] == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if api_data["response"] != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
	}

	return tool.Get_redirect("/bbs/tool/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code))
}
