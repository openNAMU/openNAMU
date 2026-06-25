package tool

import (
	"database/sql"
	"regexp"
	"strconv"
	"strings"

	"github.com/3th1nk/cidr"
	"github.com/dlclark/regexp2"
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

func Get_user_document(db *sql.DB, user_name string) bool {
    data := ""

    QueryRow_DB(
        db,
        "select title from data where title = ?",
        []any{ &data },
        "user:" + user_name,
    )

    return data != ""
}

func Get_user_title(db *sql.DB, user_name string) string {
    user_title := ""

    QueryRow_DB(
        db,
        "select data from user_set where name = 'user_title' and id = ?",
        []any{ &user_title },
        user_name,
    )

    return user_title
}

func Get_level(db *sql.DB, ip string) []string {
    level := "0"
    QueryRow_DB(
        db,
        "select data from user_set where id = ? and name = 'level'",
        []any{ &level },
        ip,
    )

    exp := "0"
    QueryRow_DB(
        db,
        "select data from user_set where id = ? and name = 'experience'",
        []any{ &exp },
        ip,
    )

    level_int := Str_to_int(level)
    max_exp := strconv.Itoa(level_int * 50 + 500)

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
        []any{ &ip_view },
    )

    user_name_view := ""
    QueryRow_DB(
        db,
        "select data from other where name = 'user_name_view'",
        []any{ &user_name_view },
    )

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
            sub_user_name := ""
            QueryRow_DB(
                db,
                "select data from user_set where id = ? and name = 'sub_user_name'",
                []any{ &sub_user_name },
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
                []any{ &user_name },
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
            []any{ &alarm_count },
            my_ip,
        )

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
    switch ban_type {
    case "O":
        return "1"
    case "E":
        return "2"
    case "A":
        return "3"
    case "D":
        return "4"
    case "L":
        return "5"
    default:
        return ""
    }
}

// Get_user_ban : login, register, edit_request, ""
// Return : []string{"true", "a" + ban_type}
func Get_user_ban(db *sql.DB, ip string, tool string) []string {
    rows := Query_DB(
        db,
        "select login, block from rb where band = 'regex' and ongoing = '1'",
    )
    defer rows.Close()

    for rows.Next() {
        var login string
        var block string

        err := rows.Scan(&login, &block)
        if err != nil {
            panic(err)
        }

        ban_type := Get_user_ban_type(login)

        r := regexp2.MustCompile(block, 0)
        if m, _ := r.FindStringMatch(ip); m != nil {
            switch tool {
            case "login":
                if ban_type != "1" && ban_type != "5" {
                    return []string{"true", "a" + ban_type}
                }
            case "register":
                if ban_type != "5" {
                    return []string{"true", "a" + ban_type}
                }
            case "edit_request":
                if ban_type != "2" {
                    return []string{"true", "a" + ban_type}
                }
            default:
                return []string{"true", "a" + ban_type}
            }
        }
    }

    if IP_or_user(ip) {
        rows := Query_DB(
            db,
            "select login, block from rb where band = 'cidr' and ongoing = '1'",
        )
        defer rows.Close()

        for rows.Next() {
            var login string
            var block string

            err := rows.Scan(&login, &block)
            if err != nil {
                panic(err)
            }

            ban_type := Get_user_ban_type(login)

            c, err := cidr.Parse(block)
            if err != nil {
                continue
            } else if c.Contains(ip) {
                switch tool {
                case "login":
                    if ban_type != "1" && ban_type != "5" {
                        return []string{"true", "b" + ban_type}
                    }
                case "register":
                    if ban_type != "5" {
                        return []string{"true", "b" + ban_type}
                    }
                case "edit_request":
                    if ban_type != "2" {
                        return []string{"true", "b" + ban_type}
                    }
                default:
                    return []string{"true", "b" + ban_type}
                }
            }
        }
    }

    login := ""
    exist := QueryRow_DB(
        db,
        "select login from rb where block = ? and (band = '' or band = 'private') and ongoing = '1'",
        []any{ &login },
        ip,
    )

    if exist {
        ban_type := Get_user_ban_type(login)

        switch tool {
        case "login":
            if ban_type != "1" && ban_type != "5" {
                return []string{"true", ban_type}
            }
        case "register":
            if ban_type != "5" {
                return []string{"true", ban_type}
            }
        case "edit_request":
            if ban_type != "2" {
                return []string{"true", ban_type}
            }
        default:
            return []string{"true", ban_type}
        }
    }

    data := ""
    exist = QueryRow_DB(
        db,
        "select data from user_set where id = ? and name = 'acl'",
        []any{ &data },
        ip,
    )

    if exist {
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
            user_name_level := ""
            QueryRow_DB(
                db,
                "select data from other where name = 'user_name_level'",
                []any{ &user_name_level },
            )

            if user_name_level != "" {
                level_data := Get_level(db, raw_ip)
                ip += "<sup>" + level_data[0] + "</sup>"
            }

            ip = "<a href=\"/w/" + Url_parser("user:" + raw_ip) + "\">" + ip + "</a>"
            user_title := Get_user_title(db, raw_ip)

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

func Do_ban_insert(db *sql.DB, user_name string, end_date string, reason string, login string, blocker string, do_type string, release bool) {
    now_time := Get_time()

    Exec_DB(
        db,
        "update rb set ongoing = '' where block = ? and band = ? and ongoing = '1'",
        user_name,
        do_type,
    )
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
}

func Get_main_skin_set(db *sql.DB, config Config, set_name string) string {
    set_data := ""

    if !IP_or_user(config.IP) {
        QueryRow_DB(
            db,
            "select data from user_set where name = ? and id = ?",
            []any{ &set_data },
            set_name,
            config.IP,
        )
    }

    if set_data == "default" || set_data == "" {
        QueryRow_DB(
            db,
            "select data from other where name = ?",
            []any{ &set_data },
            set_name,
        )
    }

    if set_data == "" {
        set_data = "default"
    }

    return set_data
}