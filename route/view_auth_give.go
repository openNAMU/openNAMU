package route

import (
	"database/sql"
	"net/url"
	"sort"
	"strings"

	"opennamu/route/tool"
)

func auth_groups(db *sql.DB) []string {
	rows := tool.Query_DB(db, "select distinct name from alist order by name asc")
	defer rows.Close()

	groups := []string{}
	for rows.Next() {
		name := ""
		if rows.Scan(&name) == nil {
			groups = append(groups, name)
		}
	}
	sort.Strings(groups)
	return groups
}

func auth_select(name string, groups []string, selected string) string {
	data := `<select name="` + name + `">`
	for _, group := range groups {
		choice := ""
		if group == selected {
			choice = ` selected`
		}
		data += `<option value="` + tool.HTML_escape(group) + `"` + choice + `>` + tool.HTML_escape(group) + `</option>`
	}
	return data + `</select>`
}

func View_auth_give(config tool.Config, mode string, user_name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	groups := auth_groups(db)
	if len(groups) == 0 {
		return tool.Get_error_page(db, config, "error")
	}

	if values != nil {
		change_auth := values.Get("change_auth")
		if !tool.Arr_in_str(groups, change_auth) {
			return tool.Get_error_page(db, config, "error")
		}

		if mode == "total" {
			result := Api_give_auth_patch(config, values.Get("auth"), change_auth, "")
			if result["response"] != "ok" {
				return tool.Get_error_page(db, config, "auth")
			}
		} else {
			names := []string{user_name}
			if user_name == "" {
				names = strings.Split(strings.ReplaceAll(values.Get("user_name"), "\r", ""), "\n")
			}
			for _, name := range names {
				name = strings.TrimSpace(name)
				if name == "" {
					continue
				}
				result := Api_give_auth_patch(config, "", change_auth, name)
				if result["response"] != "ok" {
					return tool.Get_error_page(db, config, "auth")
				}
			}
		}

		if mode == "total" {
			return tool.Get_redirect("/auth/give_total")
		}
		if user_name != "" {
			return tool.Get_redirect("/auth/give/" + tool.Url_parser(user_name))
		}
		return tool.Get_redirect("/auth/give")
	}

	data := `<form method="post">`
	if mode == "total" {
		data += auth_select("auth", groups, "") + `<hr class="main_hr">`
	} else if user_name == "" {
		data += `<textarea class="opennamu_textarea_100" name="user_name" placeholder="` + tool.Get_language(db, "many_delete_help", true) + `"></textarea><hr class="main_hr">`
	} else {
		data += `<div id="opennamu_get_user_info">` + tool.HTML_escape(user_name) + `</div><hr class="main_hr">`
	}

	selected := ""
	if user_name != "" {
		selected = tool.Get_user_auth(db, user_name)
	}
	data += auth_select("change_auth", groups, selected)
	data += `<hr class="main_hr"><button type="submit">` + tool.Get_language(db, "send", true) + `</button></form>`

	title := tool.Get_language(db, "authorize", true)
	if mode == "total" {
		title = tool.Get_language(db, "auth_to_auth", true)
	} else if user_name == "" {
		title = tool.Get_language(db, "multiple_authorize", true)
	}

	return tool.Get_template(db, config, title, data, []any{}, [][]any{{"manager", tool.Get_language(db, "return", true)}}, map[string]string{})
}
