package tool

import (
	"database/sql"
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

func List_auth(db *sql.DB) []string {
    rows := Query_DB(
        db,
        "select distinct name from alist",
    )
    defer rows.Close()

    data_list := []string{}

    for rows.Next() {
        var name string

        err := rows.Scan(&name)
        if err != nil {
            panic(err)
        }

        data_list = append(data_list, name)
    }

    return data_list
}

func Do_insert_auth_history(db *sql.DB, ip string, what string) {
    log_off := ""
    QueryRow_DB(
        db,
        "select data from other where name = 'auth_history_off'",
        []any{ &log_off },
    )

    if log_off == "" {
        time := Get_time()
        
        Exec_DB(
            db,
            "insert into re_admin (who, what, time) values (?, ?, ?)",
            ip, what, time,
        )
    }
}

func Get_user_auth(db *sql.DB, ip string) string {
    auth := "ip"
    exist := QueryRow_DB(
        db,
        "select data from user_set where id = ? and name = 'acl'",
        []any{ &auth },
        ip,
    )

    if !exist {
        if !IP_or_user(ip) {
            auth = "user"
        } else {
            auth = "ip"
        }
    }

    return auth
}

func Get_auth_group_info(db *sql.DB, auth string) map[string]bool {
    rows := Query_DB(
        db,
        "select acl from alist where name = ?",
        auth,
    )
    defer rows.Close()

    data_list := map[string]bool{}

    for rows.Next() {
        var name string

        err := rows.Scan(&name)
        if err != nil {
            panic(err)
        }

        data_list[name] = true
    }

    if len(data_list) == 0 {
        data_list["nothing"] = true
    }

    return Check_auth(data_list)
}

func Get_auth_date(db *sql.DB, user_name string) string {
    data := ""

    QueryRow_DB(
        db,
        "select data from user_set where id = ? and name = 'auth_date'",
        []any{ &data },
        user_name,
    )

    if data == "" {
        data = "0"
    }

    return data
}

func Auth_include_upper_auth(auth_info map[string]bool) bool {
    return auth_info["owner"]
}

func Check_auth(auth_info map[string]bool) map[string]bool {
    if _, ok := auth_info["owner"]; ok {
        auth_info["admin"] = true
    }

    admin_auth := []string{"ban", "toron", "check", "acl", "hidel", "give", "bbs", "vote_fix"}

    if _, ok := auth_info["admin"]; ok {
        for _, v := range admin_auth {
            auth_info[v] = true
        }
    }

    if _, ok := auth_info["check"]; ok {
        auth_info["view_user_watchlist"] = true
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

    admin_default_feature := []string{"treat_as_admin", "user_name_bold", "multiple_upload", "slow_edit_pass", "edit_bottom_compulsion_pass", "view_hide_user_name", "doc_watch_list_view", "edit_filter_pass", "user"}

    if _, ok := auth_info["admin_default_feature"]; ok {
        for _, v := range admin_default_feature {
            auth_info[v] = true
        }
    }

    user_default := []string{"captcha_pass", "ip"}

    if _, ok := auth_info["user"]; ok {
        for _, v := range user_default {
            auth_info[v] = true
        }
    }

    ip_default := []string{"document", "discuss", "upload", "vote", "bbs_use", "captcha_one_check_five_pass", "edit_filter_view"}

    if _, ok := auth_info["ip"]; ok {
        for _, v := range ip_default {
            auth_info[v] = true
        }
    }

    document_default := []string{"edit", "edit_request", "move", "new_make", "delete"}

    if _, ok := auth_info["document"]; ok {
        for _, v := range document_default {
            auth_info[v] = true
        }
    }

    check = false
    for _, v := range document_default {
        if _, ok := auth_info[v]; ok {
            check = true
            break
        }
    }

    if check {
        auth_info["view"] = true
    }

    topic_default := []string{"discuss_view", "discuss_make_new_thread"}

    if _, ok := auth_info["discuss"]; ok {
        for _, v := range topic_default {
            auth_info[v] = true
        }
    }

    bbs_default := []string{"bbs_edit", "bbs_comment"}

    if _, ok := auth_info["bbs_use"]; ok {
        for _, v := range bbs_default {
            auth_info[v] = true
        }
    }

    check = false
    for _, v := range bbs_default {
        if _, ok := auth_info[v]; ok {
            check = true
            break
        }
    }

    if check {
        auth_info["bbs_view"] = true
    }

    return auth_info
}

// PASS is TRUE
func Check_acl(db *sql.DB, name string, topic_number string, tool string, ip string) bool {
    auth_name := Get_user_auth(db, ip)
    auth_info := Get_auth_group_info(db, auth_name)

    ip_or_user := IP_or_user(ip)
    level := "0"
    if !ip_or_user {
        level = Get_level(db, ip)[0]
    }

    level_int := Str_to_int(level)

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
        switch ban_type_len {
        case 1:
            ban_type = string(ban_type[0])
        case 2:
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

            acl_data := ""
            QueryRow_DB(
                db,
                "select data from acl where title = ? and type = 'decu'",
                []any{ &acl_data },
                name,
            )

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
        if !Check_acl(db, name, topic_number, "render", ip) {
            return false
        } else if !Check_acl(db, name, topic_number, "", ip) {
            return false
        }
    } else if Arr_in_str([]string{"bbs_edit", "bbs_comment"}, tool) {
        if !Check_acl(db, name, topic_number, "bbs_view", ip) {
            return false
        }
    } else if Arr_in_str([]string{"topic"}, tool) {
        if !Check_acl(db, name, topic_number, "topic_view", ip) {
            return false
        }
    }

    if Arr_in_str([]string{"topic", "topic_view"}, tool) {
        if name == "" {
            name = "test"
            QueryRow_DB(
                db,
                "select title from rd where code = ?",
                []any{ &name },
                topic_number,
            )
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
            acl_pass_auth = "ban"
            acl_data = "owner"
        } else if tool == "bbs_auth" {
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
        } else if tool == "vote_auth" {
            acl_pass_auth = "vote_fix"
            acl_data = "owner"
        } else if tool == "" {
            acl_pass_auth = "acl"

            if for_a == 0 {
                end_number += 1

                acl_data = ""
                QueryRow_DB(
                    db,
                    "select data from acl where title = ? and type = 'decu'",
                    []any{ &acl_data },
                    name,
                )
            } else {
                if auth_info["document"] {
                    acl_data = ""
                } else {
                    acl_data = "owner"
                }
            }
        } else if tool == "document_move" {
            acl_pass_auth = "acl"

            if for_a == 0 {
                end_number += 1

                acl_data = ""
                QueryRow_DB(
                    db,
                    "select data from acl where title = ? and type = 'document_move_acl'",
                    []any{ &acl_data },
                    name,
                )
            } else {
                if auth_info["move"] {
                    acl_data = ""
                } else {
                    acl_data = "owner"
                }
            }
        } else if tool == "document_edit" {
            acl_pass_auth = "acl"

            if for_a == 0 {
                end_number += 1

                acl_data = ""
                QueryRow_DB(
                    db,
                    "select data from acl where title = ? and type = 'document_edit_acl'",
                    []any{ &acl_data },
                    name,
                )
            } else {
                if auth_info["edit"] {
                    acl_data = ""
                } else {
                    acl_data = "owner"
                }
            }
        } else if tool == "document_delete" {
            acl_pass_auth = "acl"

            if for_a == 0 {
                end_number += 1

                acl_data = ""
                QueryRow_DB(
                    db,
                    "select data from acl where title = ? and type = 'document_delete_acl'",
                    []any{ &acl_data },
                    name,
                )
            } else {
                if auth_info["delete"] {
                    acl_data = ""
                } else {
                    acl_data = "owner"
                }
            }
        } else if tool == "topic" {
            acl_pass_auth = "topic"

            switch for_a {
            case 0:
                end_number += 1

                acl_data = ""
                QueryRow_DB(
                    db,
                    "select acl from rd where code = ?",
                    []any{ &acl_data },
                    topic_number,
                )
            case 1:
                end_number += 1

                acl_data = ""
                QueryRow_DB(
                    db,
                    "select data from acl where title = ? and type = 'dis'",
                    []any{ &acl_data },
                    name,
                )
            default:
                if auth_info["discuss"] {
                    acl_data = ""
                } else {
                    acl_data = "owner"
                }
            }
        } else if tool == "topic_view" {
            acl_pass_auth = "topic"

            if for_a == 0 {
                end_number += 1

                acl_data = ""
                QueryRow_DB(
                    db,
                    "select set_data from topic_set where thread_code = ? and set_name = 'thread_view_acl'",
                    []any{ &acl_data },
                    topic_number,
                )
            } else {
                if auth_info["discuss_view"] {
                    acl_data = ""
                } else {
                    acl_data = "owner"
                }
            }
        } else if tool == "upload" {
            acl_pass_auth = "admin_default_feature"

            if auth_info["upload"] {
                acl_data = ""
            } else {
                acl_data = "owner"
            }
        } else if tool == "many_upload" {
            acl_pass_auth = "admin_default_feature"

            if auth_info["multiple_upload"] {
                acl_data = ""
            } else {
                acl_data = "owner"
            }
        } else if tool == "vote" {
            acl_pass_auth = "vote_fix"

            if for_a == 0 {
                end_number += 1

                if topic_number != "" {
                    acl_data = ""
                    QueryRow_DB(
                        db,
                        "select acl from vote where id = ? and user = ''",
                        []any{ &acl_data },
                        topic_number,
                    )
                } else {
                    continue
                }
            } else {
                if auth_info["vote"] {
                    acl_data = ""
                } else {
                    acl_data = "owner"
                }
            }
        } else if tool == "slow_edit" {
            acl_pass_auth = "admin_default_feature"

            if auth_info["slow_edit_pass"] {
                acl_data = ""
            } else {
                acl_data = "owner"
            }
        } else if tool == "edit_bottom_compulsion" {
            acl_pass_auth = "admin_default_feature"

            if auth_info["edit_bottom_compulsion_pass"] {
                acl_data = ""
            } else {
                acl_data = "owner"
            }
        } else if tool == "bbs_edit" {
            acl_pass_auth = "bbs"

            switch for_a {
            case 0:
                end_number += 1

                acl_data = ""
                QueryRow_DB(
                    db,
                    "select set_data from bbs_set where set_name = 'bbs_edit_acl' and set_id = ?",
                    []any{ &acl_data },
                    name,
                )
            case 1:
                end_number += 1
                
                acl_data = ""
                QueryRow_DB(
                    db,
                    "select set_data from bbs_set where set_name = 'bbs_acl' and set_id = ?",
                    []any{ &acl_data },
                    name,
                )
            case 2:
                end_number += 1
                
                acl_data = ""
                QueryRow_DB(
                    db,
                    "select set_data from bbs_set where set_name = 'bbs_edit_acl_all'",
                    []any{ &acl_data },
                )
            default:
                if auth_info["bbs_edit"] {
                    acl_data = ""
                } else {
                    acl_data = "owner"
                }
            }
        } else if tool == "bbs_comment" {
            acl_pass_auth = "bbs"

            switch for_a {
            case 0:
                end_number += 1
                
                acl_data = ""
                QueryRow_DB(
                    db,
                    "select set_data from bbs_set where set_name = 'bbs_comment_acl' and set_id = ?",
                    []any{ &acl_data },
                    name,
                )
            case 1:
                end_number += 1
                
                acl_data = ""
                QueryRow_DB(
                    db,
                    "select set_data from bbs_set where set_name = 'bbs_acl' and set_id = ?",
                    []any{ &acl_data },
                    name,
                )
            case 2:
                end_number += 1
                
                acl_data = ""
                QueryRow_DB(
                    db,
                    "select set_data from bbs_set where set_name = 'bbs_comment_acl_all'",
                    []any{ &acl_data },
                )
            default:
                if auth_info["bbs_comment"] {
                    acl_data = ""
                } else {
                    acl_data = "owner"
                }
            }
        } else if tool == "bbs_view" {
            acl_pass_auth = "bbs"

            if for_a == 0 {
                end_number += 1
                
                acl_data = ""
                QueryRow_DB(
                    db,
                    "select set_data from bbs_set where set_name = 'bbs_view_acl' and set_id = ?",
                    []any{ &acl_data },
                    name,
                )
            } else {
                if auth_info["bbs_view"] {
                    acl_data = ""
                } else {
                    acl_data = "owner"
                }
            }
        } else if tool == "discuss_make_new_thread" {
            acl_pass_auth = "toron"

            if auth_info["discuss_make_new_thread"] {
                acl_data = ""
            } else {
                acl_data = "owner"
            }
        } else if tool == "recaptcha" {
            acl_pass_auth = "admin_default_feature"

            if auth_info["captcha_pass"] {
                acl_data = ""
            } else {
                acl_data = "owner"
            }
        } else if tool == "recaptcha_five_pass" {
            acl_pass_auth = "admin_default_feature"

            if auth_info["captcha_one_check_five_pass"] {
                acl_data = ""
            } else {
                acl_data = "owner"
            }
        } else if tool == "edit_filter_pass" {
            acl_pass_auth = "admin_default_feature"

            if auth_info["edit_filter_pass"] {
                acl_data = ""
            } else {
                acl_data = "owner"
            }
        } else if tool == "edit_filter_view" {
            acl_pass_auth = "edit_filter_pass"

            if auth_info["edit_filter_view"] {
                acl_data = ""
            } else {
                acl_data = "owner"
            }
        } else if tool == "view_hide_user_name" {
            acl_pass_auth = "admin_default_feature"

            if auth_info["view_hide_user_name"] {
                acl_data = ""
            } else {
                acl_data = "owner"
            }
        } else if tool == "user_name_bold" {
            acl_pass_auth = "admin_default_feature"

            if auth_info["user_name_bold"] {
                acl_data = ""
            } else {
                acl_data = "owner"
            }
        } else if tool == "doc_watch_list_view" {
            acl_pass_auth = "admin_default_feature"

            if auth_info["doc_watch_list_view"] {
                acl_data = ""
            } else {
                acl_data = "owner"
            }
        } else if tool == "document_edit_request" {
            acl_pass_auth = "acl"

            if for_a == 0 {
                end_number += 1
                
                acl_data = ""
                QueryRow_DB(
                    db,
                    "select data from acl where title = ? and type = 'document_edit_request_acl'",
                    []any{ &acl_data },
                    name,
                )
            } else {
                if auth_info["edit_request"] {
                    acl_data = ""
                } else {
                    acl_data = "owner"
                }
            }
        } else if tool == "document_make_acl" {
            acl_pass_auth = "acl"

            if auth_info["new_make"] {
                acl_data = ""
            } else {
                acl_data = "owner"
            }
        } else {
            // tool == "render"
            acl_pass_auth = "acl"

            if for_a == 0 {
                end_number += 1

                acl_data = ""
                QueryRow_DB(
                    db,
                    "select data from acl where title = ? and type = 'view'",
                    []any{ &acl_data },
                    name,
                )
            } else {
                if auth_info["view"] {
                    acl_data = ""
                } else {
                    acl_data = "owner"
                }
            }
        }

        if auth_info[acl_pass_auth] {
            return true
        } else if ban_type == "4" {
            return false
        }

        if acl_data == "" {
            acl_data = "normal"
        }

        except_ban_tool_list := []string{"render", "topic_view", "bbs_view"}
        if acl_data != "normal" {
            if (acl_data != "ban" && acl_data != "ban_admin") || ban_type == "3" {
                if !Arr_in_str(except_ban_tool_list, tool) {
                    if get_ban == "true" {
                        return false
                    }
                }
            }

            switch acl_data {
            case "all", "ban":
                return true
            case "user":
                if !ip_or_user {
                    return true
                }
            case "admin":
                if auth_info["treat_as_admin"] {
                    return true
                }
            case "50_edit":
                if !ip_or_user {
                    count := 0
                    QueryRow_DB(
                        db,
                        "select count(*) from history where ip = ?",
                        []any{ &count },
                        ip,
                    )

                    if count >= 50 {
                        return true
                    }
                }
            case "before":
                exist := ""
                QueryRow_DB(
                    db,
                    "select ip from history where title = ? and ip = ? and type != 'edit_request'",
                    []any{ &exist },
                    name, ip,
                )

                if exist != "" {
                    return true
                }
            case "30_day", "90_day":
                if !ip_or_user {
                    signup_date := Get_time()
                    QueryRow_DB(
                        db,
                        "select data from user_set where id = ? and name = 'date'",
                        []any{ &signup_date },
                        ip,
                    )

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
            case "email":
                if !ip_or_user {
                    exist := ""
                    QueryRow_DB(
                        db,
                        "select data from user_set where id = ? and name = 'email'",
                        []any{ &exist },
                        ip,
                    )

                    if exist != "" {
                        return true
                    }
                }
            case "owner":
                if auth_info["owner"] {
                    return true
                }
            case "ban_admin":
                if auth_info["treat_as_admin"] || get_ban == "true" {
                    return true
                }
            case "not_all":
                return false
            case "up_to_level_3", "up_to_level_10":
                switch acl_data {
                case "up_to_level_3":
                    if level_int >= 3 {
                        return true
                    }
                case "up_to_level_10":
                    if level_int >= 10 {
                        return true
                    }
                }
            case "30_day_50_edit":
                if !ip_or_user {
                    signup_date := Get_time()
                    QueryRow_DB(
                        db,
                        "select data from user_set where id = ? and name = 'date'",
                        []any{ &signup_date },
                        ip,
                    )

                    time_1, _ := time.Parse("2006-01-02 15:04:05", signup_date)
                    time_1 = time_1.AddDate(0, 0, 30)

                    time_2, _ := time.Parse("2006-01-02 15:04:05", Get_time())
                    if time_2.After(time_1) {
                        count := 0
                        QueryRow_DB(
                            db,
                            "select count(*) from history where ip = ?",
                            []any{ &count },
                            ip,
                        )

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
                topic_state := ""
                QueryRow_DB(
                    db,
                    "select title from rd where code = ? and stop != ''",
                    []any{ &topic_state },
                    topic_number,
                )

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