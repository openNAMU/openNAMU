package tool

import (
	"database/sql"
	"regexp"
	"strconv"
)

func IP_or_user(ip string) bool {
	match, _ := regexp.MatchString("(\\.|:)", ip)
	if match {
		return true
	} else {
		return false
	}
}

func Get_level(db_set map[string]string, ip string) []string {
	db := DB_connect(db_set)
	if db == nil {
		return []string{"", "", ""}
	}
	defer db.Close()

	var level string
	var exp string
	var max_exp string

	stmt, err := db.Prepare(DB_change(db_set, "select data from user_set where id = ? and name = 'level'"))
	if err != nil {
		return []string{"", "", ""}
	}
	defer stmt.Close()

	err = stmt.QueryRow(ip).Scan(&level)
	if err != nil {
		if err == sql.ErrNoRows {
			level = "0"
		} else {
			return []string{"", "", ""}
		}
	}

	stmt, err = db.Prepare(DB_change(db_set, "select data from user_set where id = ? and name = 'experience'"))
	if err != nil {
		return []string{"", "", ""}
	}
	defer stmt.Close()

	err = stmt.QueryRow(ip).Scan(&exp)
	if err != nil {
		if err == sql.ErrNoRows {
			exp = "0"
		} else {
			return []string{"", "", ""}
		}
	}

	level_int, _ := strconv.Atoi(level)
	max_exp = strconv.Itoa(level_int*50 + 500)

	return []string{level, exp, max_exp}
}

func Get_admin_auth(db_set map[string]string, ip string) string {
	db := DB_connect(db_set)
	if db == nil {
		return ""
	}
	defer db.Close()

	if !IP_or_user(ip) {
		var auth string

		stmt, err := db.Prepare(DB_change(db_set, "select data from user_set where id = ? and name = 'acl'"))
		if err != nil {
			return ""
		}
		defer stmt.Close()

		err = stmt.QueryRow(ip).Scan(&auth)
		if err != nil {
			if err == sql.ErrNoRows {
				auth = "user"
			} else {
				return ""
			}
		}

		if auth != "user" {
			return auth
		} else {
			return ""
		}
	}

	return ""
}

func IP_preprocess(db_set map[string]string, ip string, my_ip string) []string {
	db := DB_connect(db_set)
	if db == nil {
		return []string{"", ""}
	}
	defer db.Close()

	var ip_view string
	var user_name_view string

	err := db.QueryRow(DB_change(db_set, "select data from other where name = 'ip_view'")).Scan(&ip_view)
	if err != nil {
		if err == sql.ErrNoRows {
			ip_view = ""
		} else {
			return []string{"", ""}
		}
	}

	err = db.QueryRow(DB_change(db_set, "select data from other where name = 'user_name_view'")).Scan(&user_name_view)
	if err != nil {
		if err == sql.ErrNoRows {
			user_name_view = ""
		} else {
			return []string{"", ""}
		}
	}

	if Get_admin_auth(db_set, my_ip) != "" {
		ip_view = ""
		user_name_view = ""
	}

	ip_change := ""
	if IP_or_user(ip) {
		if ip_view != "" {
			hash_ip := Sha224(ip)
			ip = hash_ip[:10]
			ip_change = "true"
		}
	} else {
		if user_name_view != "" {
			var sub_user_name string

			stmt, err := db.Prepare(DB_change(db_set, "select data from user_set where id = ? and name = 'sub_user_name'"))
			if err != nil {
				return []string{"", ""}
			}
			defer stmt.Close()

			err = stmt.QueryRow(ip).Scan(&sub_user_name)
			if err != nil {
				if err == sql.ErrNoRows {
					sub_user_name = "user"
				} else {
					return []string{"", ""}
				}
			}

			ip = sub_user_name
			ip_change = "true"
		} else {
			var user_name string

			stmt, err := db.Prepare(DB_change(db_set, "select data from user_set where name = 'user_name' and id = ?"))
			if err != nil {
				return []string{"", ""}
			}
			defer stmt.Close()

			err = stmt.QueryRow(ip).Scan(&user_name)
			if err != nil {
				if err == sql.ErrNoRows {
					user_name = ip
				} else {
					return []string{"", ""}
				}
			}

			ip = user_name
		}
	}

	return []string{ip, ip_change}
}

func IP_parser(db_set map[string]string, ip string, my_ip string) string {
	db := DB_connect(db_set)
	if db == nil {
		return ""
	}
	defer db.Close()

	ip_pre_data := IP_preprocess(db_set, ip, my_ip)
	if ip_pre_data[0] == "" {
		return ""
	}

	if ip_pre_data[1] != "" {
		return ip_pre_data[0]
	} else {
		raw_ip := ip
		ip = ip_pre_data[0]

		if !IP_or_user(raw_ip) {
			var user_name_level string
			var user_title string

			err := db.QueryRow(DB_change(db_set, "select data from other where name = 'user_name_view'")).Scan(&user_name_level)
			if err != nil {
				if err == sql.ErrNoRows {
					user_name_level = ""
				} else {
					return ""
				}
			}

			if user_name_level != "" {
				level_data := Get_level(db_set, raw_ip)
				ip += "<sup>" + level_data[0] + "</sup>"
			}

			ip = "<a href=\"/w/" + Url_parser("user:"+raw_ip) + "\">" + ip + "</a>"

			stmt, err := db.Prepare(DB_change(db_set, "select data from user_set where name = 'user_title' and id = ?"))
			if err != nil {
				return ""
			}
			defer stmt.Close()

			err = stmt.QueryRow(raw_ip).Scan(&user_title)
			if err != nil {
				if err == sql.ErrNoRows {
					user_title = ""
				} else {
					return ""
				}
			}

			if Get_admin_auth(db_set, ip) != "" {
				ip = "<b>" + ip + "</b>"
			}

			ip = user_title + ip
		}

		ip += " <a href=\"/user/" + Url_parser(raw_ip) + "\">(" + Get_language(db_set, "tool", false) + ")</a>"

		return ip
	}
}
