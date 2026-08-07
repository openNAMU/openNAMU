package route

import "opennamu/route/tool"

func View_bbs_in_w_tool(config tool.Config, set_id string, set_code string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    pinned := ""
    pinned_name := "pinned"
    if tool.QueryRow_DB(
        db,
        "select set_data from bbs_data where set_name = 'pinned' and set_id = ? and set_code = ?",
        []any{ &pinned },
        set_id,
        set_code,
    ) {
        pinned_name = "pinned_release"
    }

    data_html := `
        <h2>` + tool.Get_language(db, "tool", true) + `</h2>
        <ul>
            <li><a href="/bbs/raw/` + tool.Url_parser(set_id) + `/` + tool.Url_parser(set_code) + `">` + tool.Get_language(db, "raw", true) + `</a></li>
        </ul>
    `

    if tool.Check_acl(db, "", "", "bbs_auth", config.IP) {
        data_html += `
            <h3>` + tool.Get_language(db, "admin", true) + `</h3>
            <ul>
                <li><a href="/bbs/pinned/` + tool.Url_parser(set_id) + `/` + tool.Url_parser(set_code) + `">` + tool.Get_language(db, pinned_name, true) + `</a></li>
            </ul>
            <h3>` + tool.Get_language(db, "owner", true) + `</h3>
            <ul>
                <li><a href="/bbs/delete/` + tool.Url_parser(set_id) + `/` + tool.Url_parser(set_code) + `">` + tool.Get_language(db, "delete", true) + `</a></li>
            </ul>
        `
    }

    return tool.Get_template(
        db,
        config,
        tool.Get_language(db, "bbs_post_tool", true),
        data_html,
        []any{},
        [][]any{
            { "bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code), tool.Get_language(db, "return", true) },
        },
        map[string]string{},
    )
}
