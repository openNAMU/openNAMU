package route

import "opennamu/route/tool"

func View_bbs_make(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    data_html := `
        <form method="post">
            <input class="__ON_INPUT__" placeholder="` + tool.Get_language(db, "bbs_name", true) + `" name="bbs_name">
            <hr class="main_hr">
            
            <span class="__ON_SELECT_DIV__">
                <select class="__ON_SELECT__" name="bbs_type">
                    <option value="comment">` + tool.Get_language(db, "comment_base", true) + `</option>
                    <option value="thread">` + tool.Get_language(db, "thread_base", true) + `</option>
                </select>
            </span>
            <hr class="main_hr">
            
            <button class="__ON_BUTTON__" type="submit">` + tool.Get_language(db, "save", true) + `</button>
        </form>
    `

    return tool.Get_template(
        db,
        config,
        tool.Get_language(db, "bbs_make", true),
        data_html,
        []any{},
        [][]any{
            { "bbs/main", tool.Get_language(db, "return", true) },
        },
        map[string]string{},
    )
}