package tool

import (
	"database/sql"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/3th1nk/cidr"
	"github.com/dlclark/regexp2"
)

func IP_or_user(ip string) bool {
	match, _ := regexp.MatchString("(\\.|:)", ip)
	if match {
		return true
	} else {
		return false
	}
}

func Get_level(db *sql.DB, ip string) []string {
	var level string
	var exp string
	var max_exp string

	stmt, err := db.Prepare(DB_change("select data from user_set where id = ? and name = 'level'"))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(ip).Scan(&level)
	if err != nil {
		if err == sql.ErrNoRows {
			level = "0"
		} else {
			log.Fatal(err)
		}
	}

	stmt, err = db.Prepare(DB_change("select data from user_set where id = ? and name = 'experience'"))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(ip).Scan(&exp)
	if err != nil {
		if err == sql.ErrNoRows {
			exp = "0"
		} else {
			log.Fatal(err)
		}
	}

	level_int, _ := strconv.Atoi(level)
	max_exp = strconv.Itoa(level_int*50 + 500)

	return []string{level, exp, max_exp}
}

func IP_preprocess(db *sql.DB, ip string, my_ip string) []string {
	var ip_view string
	var user_name_view string

	ip_split := strings.Split(ip, ":")
	if len(ip_split) != 1 && ip_split[0] == "tool" {
		return []string{ip, ""}
	}

	err := db.QueryRow(DB_change("select data from other where name = 'ip_view'")).Scan(&ip_view)
	if err != nil {
		if err == sql.ErrNoRows {
			ip_view = ""
		} else {
			log.Fatal(err)
		}
	}

	err = db.QueryRow(DB_change("select data from other where name = 'user_name_view'")).Scan(&user_name_view)
	if err != nil {
		if err == sql.ErrNoRows {
			user_name_view = ""
		} else {
			log.Fatal(err)
		}
	}

	if Check_acl(db, "", "", "view_hide_user_name", my_ip) {
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
			var sub_user_name string

			stmt, err := db.Prepare(DB_change("select data from user_set where id = ? and name = 'sub_user_name'"))
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()

			err = stmt.QueryRow(ip).Scan(&sub_user_name)
			if err != nil {
				if err == sql.ErrNoRows {
					sub_user_name = Get_language(db, "member", false)
				} else {
					log.Fatal(err)
				}
			}

			if sub_user_name == "" {
				sub_user_name = Get_language(db, "member", false)
			}

			ip = sub_user_name
			ip_change = "true"
		} else {
			var user_name string

			stmt, err := db.Prepare(DB_change("select data from user_set where name = 'user_name' and id = ?"))
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()

			err = stmt.QueryRow(ip).Scan(&user_name)
			if err != nil {
				if err == sql.ErrNoRows {
					user_name = ip
				} else {
					log.Fatal(err)
				}
			}

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
		stmt, err := db.Prepare(DB_change("select count(*) from user_notice where name = ? and readme = ''"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		var alarm_count string

		err = stmt.QueryRow(my_ip).Scan(&alarm_count)
		if err != nil {
			if err == sql.ErrNoRows {
				alarm_count = "0"
			} else {
				log.Fatal(err)
			}
		}

		if IP_or_user(my_ip) {
			menu[Get_language(db, "login", false)] = [][]string{
				{"/login", Get_language(db, "login", false)},
				{"/register", Get_language(db, "register", false)},
				{"/change", Get_language(db, "user_setting", false)},
				{"/login/find", Get_language(db, "password_search", false)},
				{"/alarm" + Url_parser(my_ip), Get_language(db, "alarm", false) + " (" + alarm_count + ")"},
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
				{"/alarm" + Url_parser(my_ip), Get_language(db, "alarm", false) + " (" + alarm_count + ")"},
			}
		}
	}

	auth_name := Check_acl(db, "", "", "ban_auth", my_ip)
	if auth_name {
		menu[Get_language(db, "admin", false)] = [][]string{
			{"/auth/ban/" + Url_parser(ip), Get_language(db, "ban", false)},
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

func Get_user_ban_type(ban_type string) string {
	if ban_type == "O" {
		return "1"
	} else if ban_type == "E" {
		return "2"
	} else if ban_type == "A" {
		return "3"
	} else if ban_type == "D" {
		return "4"
	} else if ban_type == "L" {
		return "5"
	} else {
		return ""
	}
}

func Get_user_ban(db *sql.DB, ip string, tool string) []string {
	rows, err := db.Query(DB_change("select login, block from rb where band = 'regex' and ongoing = '1'"))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var login string
		var block string

		err := rows.Scan(&login, &block)
		if err != nil {
			log.Fatal(err)
		}

		ban_type := Get_user_ban_type(login)

		r := regexp2.MustCompile(block, 0)
		if m, _ := r.FindStringMatch(ip); m != nil {
			if tool == "login" {
				if ban_type != "1" && ban_type != "5" {
					return []string{"true", "a" + ban_type}
				}
			} else if tool == "register" {
				if ban_type != "5" {
					return []string{"true", "a" + ban_type}
				}
			} else if tool == "edit_request" {
				if ban_type != "2" {
					return []string{"true", "a" + ban_type}
				}
			} else {
				return []string{"true", "a" + ban_type}
			}
		}
	}

	if IP_or_user(ip) {
		rows, err = db.Query(DB_change("select login, block from rb where band = 'cidr' and ongoing = '1'"))
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var login string
			var block string

			err := rows.Scan(&login, &block)
			if err != nil {
				log.Fatal(err)
			}

			ban_type := Get_user_ban_type(login)

			c, err := cidr.Parse(block)
			if err != nil {
				continue
			} else if c.Contains(ip) {
				if tool == "login" {
					if ban_type != "1" && ban_type != "5" {
						return []string{"true", "b" + ban_type}
					}
				} else if tool == "register" {
					if ban_type != "5" {
						return []string{"true", "b" + ban_type}
					}
				} else if tool == "edit_request" {
					if ban_type != "2" {
						return []string{"true", "b" + ban_type}
					}
				} else {
					return []string{"true", "b" + ban_type}
				}
			}
		}
	}

	stmt, err := db.Prepare(DB_change("select login from rb where block = ? and band = '' and ongoing = '1'"))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var login string

	err = stmt.QueryRow(ip).Scan(&login)
	if err != nil {
		if err == sql.ErrNoRows {

		} else {
			log.Fatal(err)
		}
	} else {
		ban_type := Get_user_ban_type(login)

		if tool == "login" {
			if ban_type != "1" && ban_type != "5" {
				return []string{"true", ban_type}
			}
		} else if tool == "register" {
			if ban_type != "5" {
				return []string{"true", ban_type}
			}
		} else if tool == "edit_request" {
			if ban_type != "2" {
				return []string{"true", ban_type}
			}
		} else {
			return []string{"true", ban_type}
		}
	}

	stmt, err = db.Prepare(DB_change("select data from user_set where id = ? and name = 'acl'"))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var data string

	err = stmt.QueryRow(ip).Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {

		} else {
			log.Fatal(err)
		}
	} else {
		if data == "ban" {
			return []string{"true", "c"}
		}
	}

	return []string{"", ""}
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
			var user_name_level string
			var user_title string

			err := db.QueryRow(DB_change("select data from other where name = 'user_name_level'")).Scan(&user_name_level)
			if err != nil {
				if err == sql.ErrNoRows {
					user_name_level = ""
				} else {
					log.Fatal(err)
				}
			}

			if user_name_level != "" {
				level_data := Get_level(db, raw_ip)
				ip += "<sup>" + level_data[0] + "</sup>"
			}

			ip = "<a href=\"/w/" + Url_parser("user:"+raw_ip) + "\">" + ip + "</a>"

			stmt, err := db.Prepare(DB_change("select data from user_set where name = 'user_title' and id = ?"))
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()

			err = stmt.QueryRow(raw_ip).Scan(&user_title)
			if err != nil {
				if err == sql.ErrNoRows {
					user_title = ""
				} else {
					log.Fatal(err)
				}
			}

			if Check_acl(db, "", "", "user_name_bold", raw_ip) {
				ip = "<b>" + ip + "</b>"
			}

			ip = user_title + ip
		}

		ban := Get_user_ban(db, raw_ip, "")
		if ban[0] == "true" {
			ip = "<sup>" + ban[1] + "</sup><s>" + ip + "</s>"
		}

		ip += "<a href=\"javascript:void(0);\" name=\"" + Url_parser(raw_ip) + "\" onclick=\"opennamu_do_ip_click(this);\"><span class=\"opennamu_svg opennamu_svg_tool\">&nbsp;</span></a>"

		return ip
	}
}
