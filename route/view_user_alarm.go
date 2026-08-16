package route

import (
	"database/sql"
	"net/url"
	"strconv"
	"strings"

	"opennamu/route/tool"
)

func View_alarm_read(config tool.Config, user_name string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if user_name == "" {
		user_name = config.IP
	}
	if user_name != config.IP && !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	tool.Exec_DB(db, "update user_notice set readme = '1' where name = ?", user_name)
	return tool.Get_redirect("/alarm/" + tool.Url_parser(user_name))
}

func View_alarm(config tool.Config, user_name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if user_name == "" {
		user_name = config.IP
	}
	if user_name != config.IP && !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if values != nil && (values.Get("all") != "" || values.Get("delete") != "") {
		if values.Get("all") != "" {
			tool.Exec_DB(db, "delete from user_notice where name = ?", user_name)
		} else if id := values.Get("delete"); id != "" {
			tool.Exec_DB(db, "delete from user_notice where name = ? and id = ?", user_name, id)
		}
		return tool.Get_redirect("/alarm/" + tool.Url_parser(user_name))
	}

	delete_url := "/alarm/delete"
	if user_name != config.IP {
		delete_url = "/alarm/" + tool.Url_parser(user_name) + "/delete"
	}
	page := 1
	if values != nil {
		page = tool.Str_to_int(values.Get("num"))
		if page < 1 {
			page = 1
		}
	}
	offset := (page - 1) * 50
	rows := tool.Query_DB(db, "select id, data, date, readme from user_notice where name = ? order by date desc limit ?, 50", user_name, offset)
	body := `<ul>`
	row_count := 0
	for rows.Next() {
		id, data, date, readme := "", "", "", ""
		if rows.Scan(&id, &data, &date, &readme) != nil {
			continue
		}
		row_count++
		data_split := strings.Split(data, " | ")
		data_html := tool.IP_parser(db, data_split[0], config.IP)
		if len(data_split) > 1 {
			data_html += " | " + strings.Join(data_split[1:], " | ")
		}
		data_style := ""
		if readme == "1" {
			data_style = ` style="opacity: 0.75;"`
		}
		body += `<li` + data_style + `>` + data_html + ` | ` + tool.HTML_escape(date) + ` <a href="` + delete_url + `/` + tool.Url_parser(id) + `">(` + tool.Get_language(db, "delete", true) + `)</a></li>`
	}
	rows.Close()
	body += `</ul>`
	read_url := "/alarm/read"
	if user_name != config.IP {
		read_url = "/alarm/" + tool.Url_parser(user_name) + "/read"
	}
	read_form := `<form method="post" action="` + read_url + `"><button type="submit">` + tool.Get_language(db, "read_all", true) + `</button></form><hr class="main_hr">`
	if row_count > 0 {
		body = read_form + `<a href="` + delete_url + `">(` + tool.Get_language(db, "delete", true) + `)</a><hr class="main_hr">` + body
	} else {
		body = read_form + body
	}
	page_url := "/alarm?num={}"
	if user_name != config.IP {
		page_url = "/alarm/" + tool.Url_parser(user_name) + "?num={}"
	}
	body += tool.Get_page_control(db, page, row_count, 50, page_url)
	return user_form_page(db, config, tool.Get_language(db, "alarm", true), body)
}

func View_alarm_delete(config tool.Config, user_name string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if user_name == "" {
		user_name = config.IP
	}
	if user_name != config.IP && !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	return_path := "/alarm"
	if user_name != config.IP {
		return_path = "/alarm/" + tool.Url_parser(user_name)
	}
	body := `<form method="post"><span>` + tool.Get_language(db, "delete_warning", true) + `</span><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "delete", true) + `</button></form>`
	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "delete", true),
		body,
		[]any{},
		[][]any{{return_path, tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}

func challenge_count(db *sql.DB, query string, id string) int {
	count := 0
	tool.QueryRow_DB(db, query, []any{&count}, id)
	return count
}

func challenge_is_complete(db *sql.DB, id string, name string) bool {
	return user_value(db, id, name) != ""
}

func challenge_design(image string, title string, info string, complete bool) string {
	border := "red"
	if complete {
		border = "green"
	}

	return `<table id="main_table_set" style="border: 2px solid ` + border + `">
		<tr>
			<td id="main_table_width_quarter" rowspan="2"><span style="font-size: 64px;">` + image + `</span></td>
			<td><span style="font-size: 32px;">` + title + `</span></td>
		</tr>
		<tr><td>` + info + `</td></tr>
	</table>
	<hr class="main_hr">`
}

func View_challenge(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !user_auth(db, config) {
		return tool.Get_redirect("/user")
	}

	if values != nil {
		edit_count := challenge_count(db, "select count(*) from history where ip = ?", config.IP)
		topic_count := challenge_count(db, "select count(*) from topic where ip = ?", config.IP)
		experience := 5 * edit_count

		for _, challenge := range []struct {
			name   string
			count  int
			reward int
		}{
			{"challenge_first_contribute", 1, 500},
			{"challenge_tenth_contribute", 10, 1000},
			{"challenge_hundredth_contribute", 100, 3000},
			{"challenge_thousandth_contribute", 1000, 10000},
		} {
			if edit_count >= challenge.count {
				user_save(db, config.IP, challenge.name, "1")
				experience += challenge.reward
			}
		}

		experience += 5 * topic_count
		for _, challenge := range []struct {
			name   string
			count  int
			reward int
		}{
			{"challenge_first_discussion", 1, 500},
			{"challenge_tenth_discussion", 10, 1000},
			{"challenge_hundredth_discussion", 100, 3000},
			{"challenge_thousandth_discussion", 1000, 10000},
		} {
			if topic_count >= challenge.count {
				user_save(db, config.IP, challenge.name, "1")
				experience += challenge.reward
			}
		}

		if tool.Check_acl(db, "", "", "all_admin_auth", config.IP) || challenge_is_complete(db, config.IP, "challenge_admin") {
			user_save(db, config.IP, "challenge_admin", "1")
			experience += 10000
		}

		level := 0
		for experience >= 500+level*50 {
			experience -= 500 + level*50
			level++
		}
		user_save(db, config.IP, "level", strconv.Itoa(level))
		user_save(db, config.IP, "experience", strconv.Itoa(experience))
		return tool.Get_redirect("/challenge")
	}

	challenge_list := []struct {
		image    string
		name     string
		complete bool
	}{
		{"🌳", "register", true},
		{"🔰", "first_contribute", challenge_is_complete(db, config.IP, "challenge_first_contribute")},
		{"📝", "tenth_contribute", challenge_is_complete(db, config.IP, "challenge_tenth_contribute")},
		{"🖊️", "hundredth_contribute", challenge_is_complete(db, config.IP, "challenge_hundredth_contribute")},
		{"🏅", "thousandth_contribute", challenge_is_complete(db, config.IP, "challenge_thousandth_contribute")},
		{"💬", "first_discussion", challenge_is_complete(db, config.IP, "challenge_first_discussion")},
		{"💡", "tenth_discussion", challenge_is_complete(db, config.IP, "challenge_tenth_discussion")},
		{"📢", "hundredth_discussion", challenge_is_complete(db, config.IP, "challenge_hundredth_discussion")},
		{"📜", "thousandth_discussion", challenge_is_complete(db, config.IP, "challenge_thousandth_discussion")},
		{"☑️", "admin", challenge_is_complete(db, config.IP, "challenge_admin")},
	}

	green_html := ""
	red_html := ""
	for _, challenge := range challenge_list {
		design := challenge_design(
			challenge.image,
			tool.Get_language(db, "challenge_title_"+challenge.name, true),
			tool.Get_language(db, "challenge_info_"+challenge.name, true),
			challenge.complete,
		)
		if challenge.complete {
			green_html += design
		} else {
			red_html += design
		}
	}

	body := green_html + red_html + `<form method="post">
		<div id="opennamu_get_user_info">` + tool.HTML_escape(config.IP) + `</div>
		<hr class="main_hr">
		<button id="opennamu_save_button" type="submit">` + tool.Get_language(db, "reload", true) + `</button>
	</form>`
	return user_form_page(db, config, tool.Get_language(db, "challenge_and_level_manage", true), body)
}
