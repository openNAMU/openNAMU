package route

import (
	"database/sql"
	"opennamu/route/tool"
)

func Get_ui_recent_block(db *sql.DB, config tool.Config, data_all [][]string) string {
    auth_name := tool.Get_user_auth(db, config.IP)
    auth_info := tool.Get_auth_group_info(db, auth_name)

    data_html := ""

    for_count := 1
    for _, in_data := range data_all {
        // for_count_str := strconv.Itoa(for_count)
        for_count += 1
    
        left := ""

        ban_auth := false
        if _, ok := auth_info["owner"]; ok {
            ban_auth = true
        } else if  _, ok := auth_info["ban"]; ok {
            ban_auth = true
        }

        ip := in_data[1]

        url_match := ""
        switch in_data[7] {
        case "":
            url_match = "ban"
        case "private":
            url_match = "ban"
            ip += " (" + tool.Get_language(db, "private", true) + ")"
        case "cidr":
            url_match = "ban_cidr"
            ip += " (" + tool.Get_language(db, "cidr", true) + ")"
        default:
            url_match = "ban_regex"
            ip += " (" + tool.Get_language(db, "regex", true) + ")"
        }

        if ban_auth {
            ip = `<a href="/auth/` + url_match + `/` + tool.Url_parser(in_data[1]) + `">` + ip + `</a>`;
        }

        if in_data[8] == "1" {
            ip = `<s>` + ip + `</s>`
        }

        left += ip + ` ← ` + in_data[4]

        end := ""
        switch in_data[5] {
        case "release":
            end = tool.Get_language(db, "release", true)
        case "":
            end = tool.Get_language(db, "limitless", true)
        default:
            end = in_data[5]
        }

        right := end + "<br>" + in_data[6]

        bottom := ""
        if in_data[0] != "" {
            if in_data[0] == "edit filter" {
                bottom = `<a href="/edit_filter/` + tool.Url_parser(in_data[1]) + `">edit filter</a>`
            } else {
                bottom = Get_safe_send_data(tool.HTML_escape(in_data[0]))
            }
        }

        data_html += tool.Get_list_ui(left, right, bottom, "")
    }

    return data_html
}

func View_list_recent_block(config tool.Config, num string, set_type string, why string, user_name string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    data_html := ""
    sub := ""

    if set_type == "" {
        set_type = "all"
    }

    menu_option := [][]string{ 
        { "all", tool.Get_language(db, "all", true) },
        { "regex", tool.Get_language(db, "regex", true) },
        { "cidr", tool.Get_language(db, "cidr", true) },
        { "private", tool.Get_language(db, "private", true) },
        { "ongoing", tool.Get_language(db, "in_progress", true) },
    }
    for _, option := range menu_option {
        data_html += `<a href="/recent_block/` + option[0] + `/1">(` + option[1] + `)</a> `

        if option[0] == set_type && set_type != "all" {
            sub = "(" + option[1] + ")"
        }
    }

    menu_option = [][]string{ 
        { "/manager/11", tool.Get_language(db, "blocked", true) },
        { "/manager/12", tool.Get_language(db, "admin", true) },
        { "/manager/19", tool.Get_language(db, "why", true) },
    }
    for _, option := range menu_option {
        data_html += `<a href="` + option[0] + `">(` + option[1] + `)</a> `
    }

    data_html += "<hr class=\"main_hr\">"

    api_data := Api_list_recent_block(config, num, set_type, why, user_name)
    api_data_list := api_data["data"].([][]string)

    data_html += Get_ui_recent_block(db, config, api_data_list)

    base_url := "/recent_block/" + tool.Url_parser(set_type)

    if user_name != "" {
        base_url += "/" + tool.Url_parser(user_name)
    }

    base_url += "/{}/"

    if why != "" {
        base_url += "/" + tool.Url_parser(why)
    }

    data_html += tool.Get_page_control(
        db,
        tool.Str_to_int(num),
        len(api_data_list),
        50,
        base_url,
    )

    out := tool.Get_template(
        db,
        config,
        tool.Get_language(db, "recent_ban", true),
        data_html,
        []any{ sub },
        [][]any{
            { "other", tool.Get_language(db, "return", true) },
        },
        map[string]string{},
    )

    return out
}
