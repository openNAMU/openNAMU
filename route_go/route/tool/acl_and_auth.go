package tool

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
	"time"
)

func List_acl(func_type string) []string {
	if func_type == "user_document" {
		return []string{
			"",
			"user",
			"all",
		}
	} else {
		return []string{
			"",
			"all",
			"user",
			"admin",
			"owner",
			"50_edit",
			"email",
			"ban",
			"before",
			"30_day",
			"90_day",
			"ban_admin",
			"not_all",
			"up_to_level_3",
			"up_to_level_10",
			"30_day_50_edit",
		}
	}
}

func Do_insert_auth_history(db *sql.DB, ip string, what string) {
	var log_off string

	err := db.QueryRow(DB_change("select data from other where name = 'auth_history_off'")).Scan(&log_off)
	if err != nil {
		if err == sql.ErrNoRows {
			log_off = ""
		} else {
			log.Fatal(err)
		}
	}

	if log_off == "" {
		stmt, err := db.Prepare(DB_change("insert into re_admin (who, what, time) values (?, ?, ?)"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		time := Get_time()

		_, err = stmt.Exec(ip, what, time)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func Get_user_auth(db *sql.DB, ip string) string {
	if !IP_or_user(ip) {
		var auth string

		stmt, err := db.Prepare(DB_change("select data from user_set where id = ? and name = 'acl'"))
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

func Get_auth_group_info(db *sql.DB, auth string) map[string]bool {
	stmt, err := db.Prepare(DB_change("select acl from alist where name = ?"))
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

	return Check_auth(data_list)
}

func Check_auth(auth_info map[string]bool) map[string]bool {
	if _, ok := auth_info["owner"]; ok {
		auth_info["admin"] = true
	}

	admin_auth := []string{"ban", "toron", "check", "acl", "hidel", "give", "bbs"}

	if _, ok := auth_info["admin"]; ok {
		for _, v := range admin_auth {
			auth_info[v] = true
		}
	}

	check := false
	for _, v := range admin_auth {
		if _, ok := auth_info[v]; ok {
			check = true
			break
		}
	}

	if check {
		auth_info["admin_default_feature"] = true
	}

	admin_default_feature := []string{"treat_as_admin", "user_name_bold", "multiple_upload", "slow_edit_pass", "edit_bottom_compulsion_pass"}

	if _, ok := auth_info["admin_default_feature"]; ok {
		for _, v := range admin_default_feature {
			auth_info[v] = true
		}

		auth_info["user"] = true
	}

	return auth_info
}

func Check_acl(db *sql.DB, name string, topic_number string, tool string, ip string) bool {
	auth_name := Get_user_auth(db, ip)
	auth_info := Get_auth_group_info(db, auth_name)

	ip_or_user := IP_or_user(ip)
	level := "0"
	if !ip_or_user {
		level = Get_level(db, ip)[0]
	}

	level_int, _ := strconv.Atoi(level)

	get_ban := ""
	ban_type := ""
	if tool == "document_edit_request" {
		temp_arr := Get_user_ban(db, ip, "edit_request")

		get_ban = temp_arr[0]
		ban_type = temp_arr[1]
	} else {
		temp_arr := Get_user_ban(db, ip, "")

		get_ban = temp_arr[0]
		ban_type = temp_arr[1]
	}

	if ban_type != "" {
		ban_type_len := len(ban_type)
		if ban_type_len == 1 {
			ban_type = string(ban_type[0])
		} else if ban_type_len == 2 {
			ban_type = string(ban_type[1])
		}
	}

	if tool == "" && name != "" {
		if !Check_acl(db, name, "", "render", ip) {
			return false
		}

		if strings.HasPrefix(name, "user:") {
			user_page_str := name[5:]
			if slash_index := strings.Index(user_page_str, "/"); slash_index != -1 {
				user_page_str = user_page_str[:slash_index]
			}

			if auth_info["acl"] {
				return true
			}

			if get_ban == "true" {
				return false
			}

			stmt, err := db.Prepare(DB_change("select data from acl where title = ? and type = 'decu'"))
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()

			var acl_data string

			err = stmt.QueryRow(name).Scan(&acl_data)
			if err != nil {
				if err == sql.ErrNoRows {
					acl_data = ""
				} else {
					log.Fatal(err)
				}
			}

			if acl_data == "all" {
				return true
			} else if acl_data == "user" {
				if !ip_or_user {
					return true
				}
			} else if ip == user_page_str {
				if !ip_or_user {
					return true
				}
			}

			return false
		}
	}

	if Arr_in_str([]string{"document_edit", "document_edit_request", "document_move", "document_delete"}, tool) {
		if !Check_acl(db, name, topic_number, "", ip) {
			return false
		}
	} else if Arr_in_str([]string{"bbs_edit", "bbs_comment"}, tool) {
		if !Check_acl(db, name, topic_number, "bbs_view", ip) {
			return false
		}
	}

	if tool == "topic" {
		if name == "" {
			stmt, err := db.Prepare(DB_change("select title from rd where code = ?"))
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()

			err = stmt.QueryRow(topic_number).Scan(&name)
			if err != nil {
				if err == sql.ErrNoRows {
					name = "test"
				} else {
					log.Fatal(err)
				}
			}
		}
	}

	end_number := 1
	for for_a := 0; for_a < end_number; for_a++ {
		acl_data := ""
		acl_pass_auth := ""

		if tool == "all_admin_auth" {
			acl_pass_auth = "treat_as_admin"
			acl_data = "owner"
		} else if tool == "owner_auth" {
			acl_pass_auth = "owner"
			acl_data = "owner"
		} else if tool == "ban_auth" {
			acl_pass_auth = "bbs"
			acl_data = "owner"
		} else if tool == "toron_auth" {
			acl_pass_auth = "toron"
			acl_data = "owner"
		} else if tool == "check_auth" {
			acl_pass_auth = "check"
			acl_data = "owner"
		} else if tool == "acl_auth" {
			acl_pass_auth = "acl"
			acl_data = "owner"
		} else if tool == "hidel_auth" {
			acl_pass_auth = "hidel"
			acl_data = "owner"
		} else if tool == "give_auth" {
			acl_pass_auth = "give"
			acl_data = "owner"
		} else if tool == "" {
			acl_pass_auth = "acl"

			if for_a == 0 {
				end_number += 1

				stmt, err := db.Prepare(DB_change("select data from acl where title = ? and type = 'decu'"))
				if err != nil {
					log.Fatal(err)
				}
				defer stmt.Close()

				err = stmt.QueryRow(name).Scan(&acl_data)
				if err != nil {
					if err == sql.ErrNoRows {
						acl_data = ""
					} else {
						log.Fatal(err)
					}
				}
			} else {
				err := db.QueryRow(DB_change("select data from other where name = 'edit'")).Scan(&acl_data)
				if err != nil {
					if err == sql.ErrNoRows {
						acl_data = ""
					} else {
						log.Fatal(err)
					}
				}
			}
		} else if tool == "document_move" {
			acl_pass_auth = "acl"

			if for_a == 0 {
				end_number += 1

				stmt, err := db.Prepare(DB_change("select data from acl where title = ? and type = 'document_move_acl'"))
				if err != nil {
					log.Fatal(err)
				}
				defer stmt.Close()

				err = stmt.QueryRow(name).Scan(&acl_data)
				if err != nil {
					if err == sql.ErrNoRows {
						acl_data = ""
					} else {
						log.Fatal(err)
					}
				}
			} else {
				err := db.QueryRow(DB_change("select data from other where name = 'document_move_acl'")).Scan(&acl_data)
				if err != nil {
					if err == sql.ErrNoRows {
						acl_data = ""
					} else {
						log.Fatal(err)
					}
				}
			}
		} else if tool == "document_edit" {
			acl_pass_auth = "acl"

			if for_a == 0 {
				end_number += 1

				stmt, err := db.Prepare(DB_change("select data from acl where title = ? and type = 'document_edit_acl'"))
				if err != nil {
					log.Fatal(err)
				}
				defer stmt.Close()

				err = stmt.QueryRow(name).Scan(&acl_data)
				if err != nil {
					if err == sql.ErrNoRows {
						acl_data = ""
					} else {
						log.Fatal(err)
					}
				}
			} else {
				err := db.QueryRow(DB_change("select data from other where name = 'document_edit_acl'")).Scan(&acl_data)
				if err != nil {
					if err == sql.ErrNoRows {
						acl_data = ""
					} else {
						log.Fatal(err)
					}
				}
			}
		} else if tool == "document_edit" {
			acl_pass_auth = "acl"

			if for_a == 0 {
				end_number += 1

				stmt, err := db.Prepare(DB_change("select data from acl where title = ? and type = 'document_delete_acl'"))
				if err != nil {
					log.Fatal(err)
				}
				defer stmt.Close()

				err = stmt.QueryRow(name).Scan(&acl_data)
				if err != nil {
					if err == sql.ErrNoRows {
						acl_data = ""
					} else {
						log.Fatal(err)
					}
				}
			} else {
				err := db.QueryRow(DB_change("select data from other where name = 'document_delete_acl'")).Scan(&acl_data)
				if err != nil {
					if err == sql.ErrNoRows {
						acl_data = ""
					} else {
						log.Fatal(err)
					}
				}
			}
		} else if tool == "topic" {
			acl_pass_auth = "topic"

			if for_a == 0 {
				end_number += 1

				stmt, err := db.Prepare(DB_change("select acl from rd where code = ?"))
				if err != nil {
					log.Fatal(err)
				}
				defer stmt.Close()

				err = stmt.QueryRow(topic_number).Scan(&acl_data)
				if err != nil {
					if err == sql.ErrNoRows {
						acl_data = ""
					} else {
						log.Fatal(err)
					}
				}
			} else if for_a == 1 {
				end_number += 1

				stmt, err := db.Prepare(DB_change("select data from acl where title = ? and type = 'dis'"))
				if err != nil {
					log.Fatal(err)
				}
				defer stmt.Close()

				err = stmt.QueryRow(topic_number).Scan(&acl_data)
				if err != nil {
					if err == sql.ErrNoRows {
						acl_data = ""
					} else {
						log.Fatal(err)
					}
				}
			} else {
				err := db.QueryRow(DB_change("select data from other where name = 'discussion'")).Scan(&acl_data)
				if err != nil {
					if err == sql.ErrNoRows {
						acl_data = ""
					} else {
						log.Fatal(err)
					}
				}
			}
		} else if tool == "topic_view" {
			acl_pass_auth = "topic"

			stmt, err := db.Prepare(DB_change("select set_data from topic_set where thread_code = ? and set_name = 'thread_view_acl'"))
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()

			err = stmt.QueryRow(topic_number).Scan(&acl_data)
			if err != nil {
				if err == sql.ErrNoRows {
					acl_data = ""
				} else {
					log.Fatal(err)
				}
			}
		} else if tool == "upload" {
			acl_pass_auth = "multiple_upload"

			err := db.QueryRow(DB_change("select data from other where name = 'upload_acl'")).Scan(&acl_data)
			if err != nil {
				if err == sql.ErrNoRows {
					acl_data = ""
				} else {
					log.Fatal(err)
				}
			}
		} else if tool == "many_upload" {
			acl_pass_auth = "multiple_upload"

			err := db.QueryRow(DB_change("select data from other where name = 'many_upload_acl'")).Scan(&acl_data)
			if err != nil {
				if err == sql.ErrNoRows {
					acl_data = ""
				} else {
					log.Fatal(err)
				}
			}
		} else if tool == "vote" {
			acl_pass_auth = "owner"

			if for_a == 0 {
				end_number += 1

				if topic_number != "" {
					stmt, err := db.Prepare(DB_change("select acl from vote where id = ? and user = ''"))
					if err != nil {
						log.Fatal(err)
					}
					defer stmt.Close()

					err = stmt.QueryRow(topic_number).Scan(&acl_data)
					if err != nil {
						if err == sql.ErrNoRows {
							acl_data = ""
						} else {
							log.Fatal(err)
						}
					}
				} else {
					continue
				}
			} else {
				err := db.QueryRow(DB_change("select data from other where name = 'vote_acl'")).Scan(&acl_data)
				if err != nil {
					if err == sql.ErrNoRows {
						acl_data = ""
					} else {
						log.Fatal(err)
					}
				}
			}
		} else if tool == "slow_edit" {
			acl_pass_auth = "slow_edit_pass"

			err := db.QueryRow(DB_change("select data from other where name = 'slow_edit_acl'")).Scan(&acl_data)
			if err != nil {
				if err == sql.ErrNoRows {
					acl_data = ""
				} else {
					log.Fatal(err)
				}
			}
		} else if tool == "edit_bottom_compulsion" {
			acl_pass_auth = "edit_bottom_compulsion_pass"

			err := db.QueryRow(DB_change("select data from other where name = 'edit_bottom_compulsion_acl'")).Scan(&acl_data)
			if err != nil {
				if err == sql.ErrNoRows {
					acl_data = ""
				} else {
					log.Fatal(err)
				}
			}
		} else if tool == "bbs_edit" {
			acl_pass_auth = "bbs"

			if for_a == 0 {
				end_number += 1

				stmt, err := db.Prepare(DB_change("select set_data from bbs_set where set_name = 'bbs_edit_acl' and set_id = ?"))
				if err != nil {
					log.Fatal(err)
				}
				defer stmt.Close()

				err = stmt.QueryRow(name).Scan(&acl_data)
				if err != nil {
					if err == sql.ErrNoRows {
						acl_data = ""
					} else {
						log.Fatal(err)
					}
				}
			} else if for_a == 1 {
				end_number += 1

				stmt, err := db.Prepare(DB_change("select set_data from bbs_set where set_name = 'bbs_acl' and set_id = ?"))
				if err != nil {
					log.Fatal(err)
				}
				defer stmt.Close()

				err = stmt.QueryRow(name).Scan(&acl_data)
				if err != nil {
					if err == sql.ErrNoRows {
						acl_data = ""
					} else {
						log.Fatal(err)
					}
				}
			} else if for_a == 2 {
				end_number += 1

				err := db.QueryRow(DB_change("select set_data from bbs_set where set_name = 'bbs_edit_acl_all'")).Scan(&acl_data)
				if err != nil {
					if err == sql.ErrNoRows {
						acl_data = ""
					} else {
						log.Fatal(err)
					}
				}
			} else {
				err := db.QueryRow(DB_change("select set_data from bbs_set where set_name = 'bbs_acl_all'")).Scan(&acl_data)
				if err != nil {
					if err == sql.ErrNoRows {
						acl_data = ""
					} else {
						log.Fatal(err)
					}
				}
			}
		} else if tool == "bbs_comment" {
			acl_pass_auth = "bbs"

			if for_a == 0 {
				end_number += 1

				stmt, err := db.Prepare(DB_change("select set_data from bbs_set where set_name = 'bbs_comment_acl' and set_id = ?"))
				if err != nil {
					log.Fatal(err)
				}
				defer stmt.Close()

				err = stmt.QueryRow(name).Scan(&acl_data)
				if err != nil {
					if err == sql.ErrNoRows {
						acl_data = ""
					} else {
						log.Fatal(err)
					}
				}
			} else if for_a == 1 {
				end_number += 1

				stmt, err := db.Prepare(DB_change("select set_data from bbs_set where set_name = 'bbs_acl' and set_id = ?"))
				if err != nil {
					log.Fatal(err)
				}
				defer stmt.Close()

				err = stmt.QueryRow(name).Scan(&acl_data)
				if err != nil {
					if err == sql.ErrNoRows {
						acl_data = ""
					} else {
						log.Fatal(err)
					}
				}
			} else if for_a == 2 {
				end_number += 1

				err := db.QueryRow(DB_change("select set_data from bbs_set where set_name = 'bbs_comment_acl_all'")).Scan(&acl_data)
				if err != nil {
					if err == sql.ErrNoRows {
						acl_data = ""
					} else {
						log.Fatal(err)
					}
				}
			} else {
				err := db.QueryRow(DB_change("select set_data from bbs_set where set_name = 'bbs_acl_all'")).Scan(&acl_data)
				if err != nil {
					if err == sql.ErrNoRows {
						acl_data = ""
					} else {
						log.Fatal(err)
					}
				}
			}
		} else if tool == "bbs_view" {
			acl_pass_auth = "bbs"

			if for_a == 0 {
				end_number += 1

				stmt, err := db.Prepare(DB_change("select set_data from bbs_set where set_name = 'bbs_view_acl' and set_id = ?"))
				if err != nil {
					log.Fatal(err)
				}
				defer stmt.Close()

				err = stmt.QueryRow(name).Scan(&acl_data)
				if err != nil {
					if err == sql.ErrNoRows {
						acl_data = ""
					} else {
						log.Fatal(err)
					}
				}
			} else {
				err := db.QueryRow(DB_change("select set_data from bbs_set where set_name = 'bbs_view_acl_all'")).Scan(&acl_data)
				if err != nil {
					if err == sql.ErrNoRows {
						acl_data = ""
					} else {
						log.Fatal(err)
					}
				}
			}
		} else if tool == "recaptcha" {
			acl_pass_auth = "captcha_pass"

			err := db.QueryRow(DB_change("select data from other where name = 'recaptcha_pass_acl'")).Scan(&acl_data)
			if err != nil {
				if err == sql.ErrNoRows {
					acl_data = ""
				} else {
					log.Fatal(err)
				}
			}
		} else if tool == "recaptcha_five_pass" {
			acl_pass_auth = "captcha_one_check_five_pass"

			err := db.QueryRow(DB_change("select data from other where name = 'recaptcha_one_check_five_pass_acl'")).Scan(&acl_data)
			if err != nil {
				if err == sql.ErrNoRows {
					acl_data = ""
				} else {
					log.Fatal(err)
				}
			}
		} else if tool == "document_edit_request" {
			acl_pass_auth = "acl"

			if for_a == 0 {
				end_number += 1

				stmt, err := db.Prepare(DB_change("select data from acl where title = ? and type = 'document_edit_request_acl'"))
				if err != nil {
					log.Fatal(err)
				}
				defer stmt.Close()

				err = stmt.QueryRow(name).Scan(&acl_data)
				if err != nil {
					if err == sql.ErrNoRows {
						acl_data = ""
					} else {
						log.Fatal(err)
					}
				}
			} else {
				err := db.QueryRow(DB_change("select data from other where name = 'document_edit_request_acl'")).Scan(&acl_data)
				if err != nil {
					if err == sql.ErrNoRows {
						acl_data = ""
					} else {
						log.Fatal(err)
					}
				}
			}
		} else if tool == "document_make_acl" {
			acl_pass_auth = "acl"

			err := db.QueryRow(DB_change("select data from other where name = 'document_make_acl'")).Scan(&acl_data)
			if err != nil {
				if err == sql.ErrNoRows {
					acl_data = ""
				} else {
					log.Fatal(err)
				}
			}
		} else {
			// tool == "render"
			acl_pass_auth = "acl"

			if for_a == 0 {
				end_number += 1

				stmt, err := db.Prepare(DB_change("select data from acl where title = ? and type = 'view'"))
				if err != nil {
					log.Fatal(err)
				}
				defer stmt.Close()

				err = stmt.QueryRow(name).Scan(&acl_data)
				if err != nil {
					if err == sql.ErrNoRows {
						acl_data = ""
					} else {
						log.Fatal(err)
					}
				}
			} else {
				err := db.QueryRow(DB_change("select data from other where name = 'all_view_acl'")).Scan(&acl_data)
				if err != nil {
					if err == sql.ErrNoRows {
						acl_data = ""
					} else {
						log.Fatal(err)
					}
				}
			}
		}

		if auth_info[acl_pass_auth] {
			return true
		} else if ban_type == "4" {
			return false
		}

		if acl_data == "" {
			if tool == "recaptcha" {
				acl_data = "admin"
			} else if tool == "slow_edit" || tool == "edit_bottom_compulsion" {
				acl_data = "not_all"
			} else {
				acl_data = "normal"
			}
		}

		except_ban_tool_list := []string{"render", "topic_view", "bbs_view"}
		if acl_data != "normal" {
			if !(acl_data == "ban" || acl_data == "ban_admin") || ban_type == "3" {
				if !Arr_in_str(except_ban_tool_list, tool) {
					if get_ban == "true" {
						return false
					}
				}
			}

			if acl_data == "all" || acl_data == "ban" {
				return true
			} else if acl_data == "user" {
				if !ip_or_user {
					return true
				}
			} else if acl_data == "admin" {
				if auth_info["treat_as_admin"] {
					return true
				}
			} else if acl_data == "50_edit" {
				if !ip_or_user {
					stmt, err := db.Prepare(DB_change("select count(*) from history where ip = ?"))
					if err != nil {
						log.Fatal(err)
					}
					defer stmt.Close()

					var count int

					err = stmt.QueryRow(ip).Scan(&count)
					if err != nil {
						if err == sql.ErrNoRows {
							count = 0
						} else {
							log.Fatal(err)
						}
					}

					if count >= 50 {
						return true
					}
				}
			} else if acl_data == "before" {
				stmt, err := db.Prepare(DB_change("select ip from history where title = ? and ip = ?"))
				if err != nil {
					log.Fatal(err)
				}
				defer stmt.Close()

				var exist string

				err = stmt.QueryRow(name, ip).Scan(&exist)
				if err != nil {
					if err == sql.ErrNoRows {
						exist = ""
					} else {
						log.Fatal(err)
					}
				}

				if exist != "" {
					return true
				}
			} else if acl_data == "30_day" || acl_data == "90_day" {
				if !ip_or_user {
					stmt, err := db.Prepare(DB_change("select data from user_set where id = ? and name = 'date'"))
					if err != nil {
						log.Fatal(err)
					}
					defer stmt.Close()

					var signup_date string

					err = stmt.QueryRow(ip).Scan(&signup_date)
					if err != nil {
						if err == sql.ErrNoRows {
							signup_date = Get_time()
						} else {
							log.Fatal(err)
						}
					}

					time_1, _ := time.Parse("2006-01-02 15:04:05", signup_date)
					if acl_data == "30_day" {
						time_1 = time_1.AddDate(0, 0, 30)
					} else {
						time_1 = time_1.AddDate(0, 0, 90)
					}

					time_2, _ := time.Parse("2006-01-02 15:04:05", Get_time())
					if time_2.After(time_1) {
						return true
					}
				}
			} else if acl_data == "email" {
				if !ip_or_user {
					stmt, err := db.Prepare(DB_change("select data from user_set where id = ? and name = 'email'"))
					if err != nil {
						log.Fatal(err)
					}
					defer stmt.Close()

					var exist string

					err = stmt.QueryRow(ip).Scan(&exist)
					if err != nil {
						if err == sql.ErrNoRows {
							exist = ""
						} else {
							log.Fatal(err)
						}
					}

					if exist != "" {
						return true
					}
				}
			} else if acl_data == "owner" {
				if auth_info["owner"] {
					return true
				}
			} else if acl_data == "ban_admin" {
				if auth_info["treat_as_admin"] || get_ban == "true" {
					return true
				}
			} else if acl_data == "not_all" {
				return false
			} else if acl_data == "up_to_level_3" || acl_data == "up_to_level_10" {
				if acl_data == "up_to_level_3" {
					if level_int >= 3 {
						return true
					}
				} else if acl_data == "up_to_level_10" {
					if level_int >= 10 {
						return true
					}
				}
			} else if acl_data == "30_day_50_edit" {
				if !ip_or_user {
					stmt, err := db.Prepare(DB_change("select data from user_set where id = ? and name = 'date'"))
					if err != nil {
						log.Fatal(err)
					}
					defer stmt.Close()

					var signup_date string

					err = stmt.QueryRow(ip).Scan(&signup_date)
					if err != nil {
						if err == sql.ErrNoRows {
							signup_date = Get_time()
						} else {
							log.Fatal(err)
						}
					}

					time_1, _ := time.Parse("2006-01-02 15:04:05", signup_date)
					time_1 = time_1.AddDate(0, 0, 30)

					time_2, _ := time.Parse("2006-01-02 15:04:05", Get_time())
					if time_2.After(time_1) {
						stmt, err := db.Prepare(DB_change("select count(*) from history where ip = ?"))
						if err != nil {
							log.Fatal(err)
						}
						defer stmt.Close()

						var count int

						err = stmt.QueryRow(ip).Scan(&count)
						if err != nil {
							if err == sql.ErrNoRows {
								count = 0
							} else {
								log.Fatal(err)
							}
						}

						if count >= 50 {
							return true
						}
					}
				}
			}

			return false
		} else if for_a == end_number-1 {
			if !Arr_in_str(except_ban_tool_list, tool) {
				if get_ban == "true" {
					return false
				}
			}

			if tool == "topic" {
				stmt, err := db.Prepare(DB_change("select title from rd where code = ? and stop != ''"))
				if err != nil {
					log.Fatal(err)
				}
				defer stmt.Close()

				var topic_state string

				err = stmt.QueryRow(topic_number).Scan(&topic_state)
				if err != nil {
					if err == sql.ErrNoRows {
						topic_state = ""
					} else {
						log.Fatal(err)
					}
				}

				if topic_state != "" {
					if auth_info["topic"] {
						return true
					} else {
						return false
					}
				} else {
					return true
				}
			} else {
				return true
			}
		}
	}

	return false
}
