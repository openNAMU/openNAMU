package tool

import (
	"database/sql"
	"log"
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

func Get_level(db *sql.DB, db_set map[string]string, ip string) []string {
	var level string
	var exp string
	var max_exp string

	stmt, err := db.Prepare(DB_change(db_set, "select data from user_set where id = ? and name = 'level'"))
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

	stmt, err = db.Prepare(DB_change(db_set, "select data from user_set where id = ? and name = 'experience'"))
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

func Get_user_auth(db *sql.DB, db_set map[string]string, ip string) string {
	if !IP_or_user(ip) {
		var auth string

		stmt, err := db.Prepare(DB_change(db_set, "select data from user_set where id = ? and name = 'acl'"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		err = stmt.QueryRow(ip).Scan(&auth)
		if err != nil {
			if err == sql.ErrNoRows {
				auth = "user"
			} else {
				log.Fatal(err)
			}
		}

		if auth != "user" && auth != "ban" {
			return auth
		} else {
			return ""
		}
	}

	return ""
}

func Get_auth_group_info(db *sql.DB, db_set map[string]string, auth string) map[string]bool {
	stmt, err := db.Prepare(DB_change(db_set, "select name from alist where name = ?"))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(auth)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	data_list := map[string]bool{}

	for rows.Next() {
		var name string

		err := rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}

		data_list[name] = true
	}

	return data_list
}

func IP_preprocess(db *sql.DB, db_set map[string]string, ip string, my_ip string) []string {
	var ip_view string
	var user_name_view string

	err := db.QueryRow(DB_change(db_set, "select data from other where name = 'ip_view'")).Scan(&ip_view)
	if err != nil {
		if err == sql.ErrNoRows {
			ip_view = ""
		} else {
			log.Fatal(err)
		}
	}

	err = db.QueryRow(DB_change(db_set, "select data from other where name = 'user_name_view'")).Scan(&user_name_view)
	if err != nil {
		if err == sql.ErrNoRows {
			user_name_view = ""
		} else {
			log.Fatal(err)
		}
	}

	if Get_user_auth(db, db_set, my_ip) != "" {
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

			stmt, err := db.Prepare(DB_change(db_set, "select data from user_set where id = ? and name = 'sub_user_name'"))
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()

			err = stmt.QueryRow(ip).Scan(&sub_user_name)
			if err != nil {
				if err == sql.ErrNoRows {
					sub_user_name = Get_language(db, db_set, "member", false)
				} else {
					log.Fatal(err)
				}
			}

			if sub_user_name == "" {
				sub_user_name = Get_language(db, db_set, "member", false)
			}

			ip = sub_user_name
			ip_change = "true"
		} else {
			var user_name string

			stmt, err := db.Prepare(DB_change(db_set, "select data from user_set where name = 'user_name' and id = ?"))
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

func IP_parser(db *sql.DB, db_set map[string]string, ip string, my_ip string) string {
	ip_pre_data := IP_preprocess(db, db_set, ip, my_ip)
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

			err := db.QueryRow(DB_change(db_set, "select data from other where name = 'user_name_level'")).Scan(&user_name_level)
			if err != nil {
				if err == sql.ErrNoRows {
					user_name_level = ""
				} else {
					log.Fatal(err)
				}
			}

			if user_name_level != "" {
				level_data := Get_level(db, db_set, raw_ip)
				ip += "<sup>" + level_data[0] + "</sup>"
			}

			ip = "<a href=\"/w/" + Url_parser("user:"+raw_ip) + "\">" + ip + "</a>"

			stmt, err := db.Prepare(DB_change(db_set, "select data from user_set where name = 'user_title' and id = ?"))
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

			if Get_user_auth(db, db_set, ip) != "" {
				ip = "<b>" + ip + "</b>"
			}

			ip = user_title + ip
		}

		ip += " <a href=\"/user/" + Url_parser(raw_ip) + "\">(" + Get_language(db, db_set, "tool", false) + ")</a>"

		return ip
	}
}
