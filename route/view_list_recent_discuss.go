package route

import (
	"database/sql"
	"opennamu/route/tool"
)

func Get_ui_recent_discuss(db *sql.DB, config tool.Config, data_all [][]string) string {
    auth_name := tool.Get_user_auth(db, config.IP)
    auth_info := tool.Get_auth_group_info(db, auth_name)

    data_html := ""

    for_count := 1
    for _, in_data := range data_all {
        // for_count_str := strconv.Itoa(for_count)
        for_count += 1

        left := `<a href="/thread/` + in_data[3] + `">` + tool.HTML_escape(in_data[1]) + `</a> `
        left += `<a href="/w/` + tool.Url_parser(in_data[0]) + `">(` + tool.HTML_escape(in_data[0]) + `)</a> `

        if _, ok := auth_info["hidel"]; ok {
            left = `<label><input type="checkbox"> ` + left + `</label>`
        }

        right := ``
        switch in_data[4] {
        case "O":
            right += tool.Get_language(db, "closed", true) + " | "
        case "S":
            right += tool.Get_language(db, "stop", true) + " | "
        }

        if in_data[8] != "" {
            right += tool.Get_language(db, "agreed_discussion", true) + " | "
        }

        right += `<a href="/thread/` + in_data[3] + `#` + in_data[7] + `">#` + in_data[7] + `</a> | `
        right += in_data[6] + " | " + in_data[2]

        data_html += tool.Get_list_ui(left, right, "", "")
    }

    return data_html
}

func View_list_recent_discuss(config tool.Config, limit string, num string, set_type string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    data_html := ""
    sub := ""

    menu_option := [][]string{ 
        { "normal", tool.Get_language(db, "normal", true) },
        { "close", tool.Get_language(db, "close_discussion", true) },
        { "open", tool.Get_language(db, "open_discussion", true) },
    }
    for _, option := range menu_option {
        data_html += `<a href="/recent_discuss/1/` + option[0] + `">(` + option[1] + `)</a> `

        if option[0] == set_type {
            sub = "(" + option[1] + ")"
        }
    }

    data_html += "<hr class=\"main_hr\">"

    api_data := Api_list_recent_discuss(config, limit, num, set_type)
    api_data_list := api_data["data"].([][]string)

    data_html += Get_ui_recent_discuss(db, config, api_data_list)
    data_html += tool.Get_page_control(
        db,
        tool.Str_to_int(num),
        len(api_data_list),
        tool.Str_to_int(limit),
        "/recent_discuss/{}/" + set_type,
    )

    out := tool.Get_template(
        db,
        config,
        tool.Get_language(db, "recent_discussion", true),
        data_html,
        []any{ sub },
        [][]any{},
        map[string]string{},
    )

    return out
}
