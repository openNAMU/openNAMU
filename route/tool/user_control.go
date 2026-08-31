package tool

import (
	"database/sql"
	"regexp"
	"strconv"
	"strings"
)

// IP is TRUE
func IP_or_user(ip string) bool {
	match, _ := regexp.MatchString("(\\.|:)", ip)
	if match {
		return true
	} else {
		return false
	}
}

// PASS is TRUE
func Get_user_name_check(db *sql.DB, user_name string) bool {
	if user_name == "" {
		return false
	}

	if HTML_escape(user_name) != user_name {
		return false
	}

	if IP_or_user(user_name) {
		return false
	}

	user_name_arr := strings.Split(user_name, "")
	if Arr_in_str(user_name_arr, "/") {
		return false
	}

	rows := Query_DB(
		db,
		`select html from html_filter where kind = "name"`,
	)
	defer rows.Close()

	for rows.Next() {
		var html string

		err := rows.Scan(&html)
		if err != nil {
			panic(err)
		}

		check_regex, err := regexp.Compile("(?i)" + html)
		if err != nil {
			panic(err)
		}

		if check_regex.MatchString(user_name) {
			return false
		}
	}

	if Get_len(user_name) > 128 {
		return false
	}

	data := ""
	QueryRow_DB(
		db,
		`select id from user_set where name = 'user_name' and data = ?`,
		[]any{&data},
		user_name,
	)

	if data != "" {
		return false
	}

	QueryRow_DB(
		db,
		`select id from user_set where id = ?`,
		[]any{&data},
		user_name,
	)

	if data != "" {
		return false
	}

	return true
}

func Get_user_document(db *sql.DB, user_name string) bool {
	data := ""

	QueryRow_DB(
		db,
		"select title from data where title = ?",
		[]any{&data},
		"user:"+user_name,
	)

	return data != ""
}

func Get_user_title(db *sql.DB, user_name string) string {
	user_title := ""

	QueryRow_DB(
		db,
		"select data from user_set where name = 'user_title' and id = ?",
		[]any{&user_title},
		user_name,
	)

	return user_title
}

func Get_level(db *sql.DB, ip string) []string {
	level := "0"
	QueryRow_DB(
		db,
		"select data from user_set where id = ? and name = 'level'",
		[]any{&level},
		ip,
	)

	exp := "0"
	QueryRow_DB(
		db,
		"select data from user_set where id = ? and name = 'experience'",
		[]any{&exp},
		ip,
	)

	level_int := Str_to_int(level)
	max_exp := strconv.Itoa(level_int*50 + 500)

	return []string{level, exp, max_exp}
}

func IP_preprocess(db *sql.DB, ip string, my_ip string) []string {
	ip_split := strings.Split(ip, ":")
	if len(ip_split) != 1 && ip_split[0] == "tool" {
		return []string{ip, ""}
	}

	ip_view := ""
	QueryRow_DB(
		db,
		"select data from other where name = 'ip_view'",
		[]any{&ip_view},
	)

	user_name_view := ""
	QueryRow_DB(
		db,
		"select data from other where name = 'user_name_view'",
		[]any{&user_name_view},
	)

	if Check_permission(db, "view_hide_user_name", my_ip) {
		ip_view = ""
		user_name_view = ""
	}

	ip_change := ""
	if IP_or_user(ip) {
		if ip_view != "" && ip != my_ip {
			hash_ip := Sha224(ip)
			ip = hash_ip[:10]
			ip_change = "true"
		}
	} else {
		if user_name_view != "" {
			sub_user_name := ""
			QueryRow_DB(
				db,
				"select data from user_set where id = ? and name = 'sub_user_name'",
				[]any{&sub_user_name},
				ip,
			)

			if sub_user_name == "" {
				sub_user_name = Get_language(db, "member", false)
			}

			ip = sub_user_name
			ip_change = "true"
		} else {
			user_name := ""
			QueryRow_DB(
				db,
				"select data from user_set where name = 'user_name' and id = ?",
				[]any{&user_name},
				ip,
			)

			if user_name == "" {
				user_name = ip
			}

			ip = user_name
		}
	}

	return []string{ip, ip_change}
}

func IP_menu(db *sql.DB, ip string, my_ip string, option string) map[string][][]string {
	menu := map[string][][]string{}

	if ip == my_ip && option == "" {
		alarm_count := "0"
		QueryRow_DB(
			db,
			"select count(*) from user_notice where name = ? and readme = ''",
			[]any{&alarm_count},
			my_ip,
		)

		if IP_or_user(my_ip) {
			menu[Get_language(db, "login", false)] = [][]string{
				{"/login", Get_language(db, "login", false)},
				{"/register", Get_language(db, "register", false)},
				{"/change", Get_language(db, "user_setting", false)},
				{"/login/find", Get_language(db, "password_search", false)},
				{"/alarm", Get_language(db, "alarm", false) + " (" + alarm_count + ")"},
			}
		} else {
			menu[Get_language(db, "login", false)] = [][]string{
				{"/logout", Get_language(db, "logout", false)},
				{"/change", Get_language(db, "user_setting", false)},
			}

			menu[Get_language(db, "tool", false)] = [][]string{
				{"/watch_list", Get_language(db, "watchlist", false)},
				{"/star_doc", Get_language(db, "star_doc", false)},
				{"/challenge", Get_language(db, "challenge_and_level_manage", false)},
				{"/acl/user:" + Url_parser(my_ip), Get_language(db, "user_document_acl", false)},
				{"/alarm", Get_language(db, "alarm", false) + " (" + alarm_count + ")"},
			}
		}
	}

	auth_name := Check_permission(db, "give", my_ip)
	if auth_name {
		menu[Get_language(db, "admin", false)] = [][]string{
			{"/auth/give/" + Url_parser(ip), Get_language(db, "ban", false)},
			{"/list/user/check_submit/" + Url_parser(ip), Get_language(db, "check", false)},
		}
	}

	menu[Get_language(db, "other", false)] = [][]string{
		{"/record/" + Url_parser(ip), Get_language(db, "edit_record", false)},
		{"/record/topic/" + Url_parser(ip), Get_language(db, "discussion_record", false)},
		{"/record/bbs/" + Url_parser(ip), Get_language(db, "bbs_record", false)},
		{"/record/bbs_comment/" + Url_parser(ip), Get_language(db, "bbs_comment_record", false)},
		{"/topic/user:" + Url_parser(ip), Get_language(db, "user_discussion", false)},
		{"/count/" + Url_parser(ip), Get_language(db, "count", false)},
	}

	return menu
}

func IP_parser(db *sql.DB, ip string, my_ip string) string {
	ip_pre_data := IP_preprocess(db, ip, my_ip)
	if ip_pre_data[0] == "" {
		return ""
	}

	if ip_pre_data[1] != "" {
		return ip_pre_data[0]
	} else {
		raw_ip := ip
		ip = HTML_escape(ip_pre_data[0])

		if !IP_or_user(raw_ip) {
			user_name_level := ""
			QueryRow_DB(
				db,
				"select data from other where name = 'user_name_level'",
				[]any{&user_name_level},
			)

			if user_name_level != "" {
				level_data := Get_level(db, raw_ip)
				ip += "<sup>" + level_data[0] + "</sup>"
			}

			ip = "<a href=\"/w/" + Url_parser("user:"+raw_ip) + "\">" + ip + "</a>"
			user_title := Get_user_title(db, raw_ip)

			if Check_permission(db, "user_name_bold", raw_ip) {
				ip = "<b>" + ip + "</b>"
			}

			ip = user_title + ip
		}

		auth_name := Get_user_auth(db, raw_ip)
		if Auth_group_name_ban(auth_name) {
			ip = "<sup>" + HTML_escape(auth_name) + "</sup><s>" + ip + "</s>"
		}

		ip += "<a href=\"/user_tool/" + Url_parser(raw_ip) + "\"><span class=\"opennamu_svg opennamu_svg_tool\">&nbsp;</span></a>"

		return ip
	}
}

func Do_auth_insert(db *sql.DB, user_name string, end_date string, reason string, login string, blocker string, do_type string, release bool) {
	now_time := Get_time()

	if do_type == "" {
		Exec_DB(
			db,
			"update rb set ongoing = '' where block = ? and (band = '' or band = 'private') and ongoing = '1'",
			user_name,
		)
	} else {
		Exec_DB(
			db,
			"update rb set ongoing = '' where block = ? and band = ? and ongoing = '1'",
			user_name,
			do_type,
		)
	}
	if release {
		Exec_DB(
			db,
			`insert into rb (block, end, today, blocker, why, band, ongoing, login) values (?, ?, ?, ?, ?, ?, '', '')`,
			user_name,
			"release",
			now_time,
			blocker,
			reason,
			do_type,
		)
	} else {
		if end_date == "0" {
			end_date = ""
		}

		Exec_DB(
			db,
			`insert into rb (block, end, today, blocker, why, band, ongoing, login) values (?, ?, ?, ?, ?, ?, '1', ?)`,
			user_name,
			end_date,
			now_time,
			blocker,
			reason,
			do_type,
			login,
		)
	}

	if do_type != "regex" && do_type != "cidr" {
		if release {
			Exec_DB(
				db,
				"delete from user_set where id = ? and name = 'acl'",
				user_name,
			)
			Exec_DB(
				db,
				"delete from user_set where id = ? and name = 'acl_end'",
				user_name,
			)
		} else {
			user_id := ""
			user_exists := IP_or_user(user_name)
			if !user_exists {
				user_exists = QueryRow_DB(
					db,
					"select id from user_set where id = ? limit 1",
					[]any{&user_id},
					user_name,
				)
			}
			if user_exists {
				auth := get_ban_auth_group(db, login)
				Exec_DB(
					db,
					"delete from user_set where id = ? and name = 'acl'",
					user_name,
				)
				Exec_DB(
					db,
					"insert into user_set (id, name, data) values (?, 'acl', ?)",
					user_name,
					auth,
				)
				Exec_DB(
					db,
					"delete from user_set where id = ? and name = 'acl_end'",
					user_name,
				)
				if end_date != "" {
					Exec_DB(
						db,
						"insert into user_set (id, name, data) values (?, 'acl_end', ?)",
						user_name,
						end_date,
					)
				}
			}
		}
	}
}

func Get_main_skin_set(db *sql.DB, config Config, set_name string) string {
	set_data := ""

	if !IP_or_user(config.IP) {
		QueryRow_DB(
			db,
			"select data from user_set where name = ? and id = ?",
			[]any{&set_data},
			set_name,
			config.IP,
		)
	} else if config.Session != nil {
		set_data, _ = config.Session.Get(set_name).(string)
	}

	if set_data == "default" || set_data == "" {
		QueryRow_DB(
			db,
			"select data from other where name = ?",
			[]any{&set_data},
			set_name,
		)
	}

	if set_data == "" {
		set_data = "default"
	}

	return set_data
}
