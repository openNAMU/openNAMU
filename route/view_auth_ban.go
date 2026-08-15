package route

import (
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"opennamu/route/tool"
)

func selected_attr(value string, target string) string {
	if value == target {
		return ` selected`
	}
	return ""
}

func View_auth_ban(config tool.Config, name string, ban_type string, multiple bool, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	ban_data := tool.Get_user_ban(db, config.IP, "login")
	can_ban := tool.Check_acl(db, "", "", "ban_auth", config.IP)
	if ban_data[0] == "true" && (tool.IP_or_user(config.IP) || tool.Check_acl(db, "", "", "all_admin_auth", config.IP)) {
		can_ban = false
	}
	if !can_ban {
		return tool.Get_error_page(db, config, "auth")
	}

	if values != nil {
		end := "0"
		if values.Get("date_type") == "date" {
			date_value := values.Get("date")
			if _, err := time.Parse("2006-01-02", date_value); err == nil {
				end = date_value + " 00:00:00"
			}
		} else {
			days, _ := strconv.Atoi(values.Get("date_days"))
			if days <= 0 {
				days = 1
			}
			end = time.Now().AddDate(0, 0, days).Format("2006-01-02 15:04:05")
		}

		login := ""
		switch values.Get("ban_option") {
		case "login_able":
			login = "L"
		case "login_able_and_regsiter_disable":
			login = "O"
		case "completely_ban":
			login = "A"
		case "dont_come_this_site":
			login = "D"
		}
		release := values.Get("ban_option") == "release"
		type_name := ""
		switch values.Get("do_ban_type") {
		case "regex":
			type_name = "regex"
		case "cidr":
			type_name = "cidr"
		case "private":
			type_name = "private"
			if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
				return tool.Get_error_page(db, config, "auth")
			}
		}

		if !multiple && name == "" {
			name = values.Get("name")
		}
		names := []string{name}
		if multiple {
			names = strings.Split(strings.ReplaceAll(values.Get("name"), "\r", ""), "\n")
		}
		for _, target := range names {
			target = strings.TrimSpace(target)
			if target == "" {
				continue
			}
			if type_name == "regex" {
				if _, err := regexp.Compile(target); err != nil {
					return tool.Get_error_page(db, config, "error")
				}
			} else if type_name == "cidr" {
				if _, _, err := net.ParseCIDR(target); err != nil {
					return tool.Get_error_page(db, config, "error")
				}
			}
			if type_name != "private" && target == config.IP {
				if !tool.Check_acl(db, "", "", "all_admin_auth", config.IP) {
					return tool.Get_error_page(db, config, "auth")
				}
			} else if type_name != "private" && !tool.Check_acl(db, "", "", "ban_auth", config.IP) {
				return tool.Get_error_page(db, config, "auth")
			}
			tool.Do_ban_insert(db, target, end, values.Get("why"), login, config.IP, type_name, release)
			tool.Do_insert_auth_history(db, config.IP, "ban ("+target+")")
		}
		return tool.Get_redirect("/recent_block")
	}

	end_date := ""
	if name != "" {
		tool.QueryRow_DB(db, "select end from rb where block = ? and ongoing = '1' limit 1", []any{&end_date}, name)
		if len(end_date) > 10 {
			end_date = end_date[:10]
		}
	}
	target_input := `<input name="name" value="` + tool.HTML_escape(name) + `" placeholder="` + tool.Get_language(db, "name_or_ip_or_regex_or_cidr", true) + `">`
	if multiple {
		target_input = `<textarea class="opennamu_textarea_500" name="name" placeholder="` + tool.Get_language(db, "name_or_ip_or_regex_or_cidr_multiple", true) + `"></textarea>`
	}
	type_options := `<option value="normal">` + tool.Get_language(db, "normal", true) + `</option><option value="regex"` + selected_attr(ban_type, "regex") + `>` + tool.Get_language(db, "regex", true) + `</option><option value="cidr"` + selected_attr(ban_type, "cidr") + `>` + tool.Get_language(db, "cidr", true) + `</option>`
	if tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		type_options += `<option value="private"` + selected_attr(ban_type, "private") + `>` + tool.Get_language(db, "private", true) + `</option>`
	}
	data := `<form method="post"><h2>` + tool.Get_language(db, "method", true) + `</h2>` + target_input + `<hr class="main_hr"><select name="do_ban_type">` + type_options + `</select><hr class="main_hr"><select name="ban_option"><option value="">` + tool.Get_language(db, "default", true) + `</option><option value="login_able">` + tool.Get_language(db, "login_able", true) + `</option><option value="login_able_and_regsiter_disable">` + tool.Get_language(db, "login_able_and_regsiter_disable", true) + `</option><option value="completely_ban">` + tool.Get_language(db, "completely_ban", true) + `</option><option value="dont_come_this_site">` + tool.Get_language(db, "dont_come_this_site", true) + `</option><option value="release">` + tool.Get_language(db, "release", true) + `</option></select>` +
		`<h2>` + tool.Get_language(db, "date", true) + `</h2><select name="date_type"><option value="date">` + tool.Get_language(db, "date", true) + `</option><option value="days">` + tool.Get_language(db, "day", true) + `</option></select><hr class="main_hr"><input name="date_days" placeholder="` + tool.Get_language(db, "day", true) + `"><hr class="main_hr"><input type="date" name="date" value="` + end_date + `"><h2>` + tool.Get_language(db, "other", true) + `</h2><input name="why" placeholder="` + tool.Get_language(db, "why", true) + `"><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`
	title := tool.Get_language(db, "ban", true)
	if multiple {
		title = tool.Get_language(db, "multiple_ban", true)
	}
	return tool.Get_template(db, config, title, data, []any{}, [][]any{{"manager", tool.Get_language(db, "return", true)}}, map[string]string{})
}
