package route

import (
	"database/sql"
	"net/url"
	"sort"
	"strings"

	"opennamu/route/tool"
)

func auth_groups(db *sql.DB) []string {
	groups := tool.List_auth(db)
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

func auth_target_type_select(db *sql.DB, target_type string, owner bool) string {
	data := `<select name="target_type">`
	for _, target := range []string{"normal", "regex", "cidr"} {
		selected := ""
		if target == target_type {
			selected = ` selected`
		}
		data += `<option value="` + target + `"` + selected + `>` + tool.Get_language(db, target, true) + `</option>`
	}
	if owner {
		selected := ""
		if target_type == "private" {
			selected = ` selected`
		}
		data += `<option value="private"` + selected + `>` + tool.Get_language(db, "private", true) + `</option>`
	}
	return data + `</select>`
}

func View_auth_give(config tool.Config, mode string, user_name string, target_type string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if target_type == "" {
		target_type = "normal"
	}
	required_auth := "give_auth"
	if target_type == "regex" || target_type == "cidr" {
		required_auth = "give_range_auth"
	}
	if target_type == "private" {
		required_auth = "owner_auth"
	}
	if !tool.Check_acl(db, "", "", required_auth, config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	groups := auth_groups(db)
	if len(groups) == 0 {
		return tool.Get_error_page(db, config, "error")
	}

	if values != nil {
		if mode != "total" && values.Get("target_type") != "" {
			target_type = values.Get("target_type")
		}
		change_auth := values.Get("change_auth")
		release := values.Get("action") == "release"
		if !release && !tool.Arr_in_str(groups, change_auth) {
			return tool.Get_error_page(db, config, "error")
		}

		if mode == "total" {
			result := Api_give_auth_patch(config, values.Get("auth"), change_auth, "", "", "normal", "", false)
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
				result := Api_give_auth_patch(config, "", change_auth, name, values.Get("end_date"), target_type, values.Get("why"), release)
				if result["response"] != "ok" {
					return tool.Get_error_page(db, config, "auth")
				}
			}
		}

		if mode == "total" {
			return tool.Get_redirect("/auth/give_total")
		}
		if user_name != "" {
			target_path := tool.Url_parser(user_name)
			if target_type != "normal" {
				target_path = target_type + "/" + target_path
			}
			return tool.Get_redirect("/auth/give/" + target_path)
		}
		return tool.Get_redirect("/auth/give")
	}

	data := `<form method="post">`
	if mode == "total" {
		data += auth_select("auth", groups, "") + `<hr class="main_hr">`
	} else if user_name == "" {
		data += `<textarea class="opennamu_textarea_100" name="user_name" placeholder="` + tool.Get_language(db, "name_or_ip_or_regex_or_cidr_multiple", true) + `"></textarea><hr class="main_hr">`
	} else {
		data += `<div id="opennamu_get_user_info">` + tool.HTML_escape(user_name) + `</div><hr class="main_hr">`
	}

	selected := ""
	end_date := ""
	if user_name != "" {
		selected = tool.Get_auth_target_group(db, user_name, target_type)
		if target_type == "normal" {
			tool.QueryRow_DB(db, "select data from user_set where id = ? and name = 'acl_end' limit 1", []any{&end_date}, user_name)
		} else {
			tool.QueryRow_DB(db, "select end from rb where block = ? and band = ? and ongoing = '1' order by today desc limit 1", []any{&end_date}, user_name, target_type)
		}
		if len(end_date) > 10 {
			end_date = end_date[:10]
		}
	}

	if mode != "total" {
		owner := tool.Check_acl(db, "", "", "owner_auth", config.IP)
		data += auth_target_type_select(db, target_type, owner) + `<hr class="main_hr">`
	}
	data += auth_select("change_auth", groups, selected)
	if mode != "total" {
		data += `<hr class="main_hr"><span>` + tool.Get_language(db, "date", true) + `</span><input type="date" name="end_date" value="` + tool.HTML_escape(end_date) + `">`
		data += `<hr class="main_hr"><input name="why" placeholder="` + tool.Get_language(db, "why", true) + `">`
		data += `<hr class="main_hr"><select name="action"><option value="give">` + tool.Get_language(db, "authorize", true) + `</option><option value="release">` + tool.Get_language(db, "release", true) + `</option></select>`
	}
	data += `<hr class="main_hr"><button type="submit">` + tool.Get_language(db, "send", true) + `</button></form>`

	title := tool.Get_language(db, "authorize", true)
	if mode == "total" {
		title = tool.Get_language(db, "auth_to_auth", true)
	} else if user_name == "" {
		title = tool.Get_language(db, "multiple_authorize", true)
	}

	return tool.Get_template(db, config, title, data, []any{}, [][]any{{"manager", tool.Get_language(db, "return", true)}}, map[string]string{})
}
