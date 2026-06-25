package route

import "opennamu/route/tool"

func View_bbs_in_w_comment_tool(config tool.Config, set_id string, set_code string, comment_id string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    data_html := `
        <h2>` + tool.Get_language(db, "tool", true) + `</h2>
        <ul>
            <li><a href="/bbs/raw/` + tool.Url_parser(set_id) + `/` + tool.Url_parser(set_code) + `/` + tool.Url_parser(comment_id) + `">` + tool.Get_language(db, "raw", true) + `</a></li>
            <li><a href="/bbs/edit/` + tool.Url_parser(set_id) + `/` + tool.Url_parser(set_code) + `/` + tool.Url_parser(comment_id) + `">` + tool.Get_language(db, "edit", true) + `</a></li>
        </ul>
    `

    if tool.Check_acl(db, "", "", "owner_auth", config.IP) {
        data_html += `
            <h3>` + tool.Get_language(db, "owner", true) + `</h3>
            <ul>
                <li><a href="/bbs/delete/` + tool.Url_parser(set_id) + `/` + tool.Url_parser(set_code) + `/` + tool.Url_parser(comment_id) + `">` + tool.Get_language(db, "delete", true) + `</a></li>
            </ul>
        `
    }

    return tool.Get_template(
        db,
        config,
        tool.Get_language(db, "bbs_comment_tool", true),
        data_html,
        []any{},
        [][]any{
            { "bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code) + "#" + tool.Url_parser(comment_id), tool.Get_language(db, "return", true) },
        },
        map[string]string{},
    )
}