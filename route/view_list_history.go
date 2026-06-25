package route

import (
	"opennamu/route/tool"
)

func View_list_history(config tool.Config, doc_name string, set_type string, num string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    sub := ""
    if set_type == "" {
        set_type = "normal"
    } else {
        sub = " (" + tool.Get_language(db, set_type, true) + ")"
    }

    data_html := ""

    menu_option := []string{ "normal", "edit", "move", "delete", "revert", "r1", "setting" }
    for _, option := range menu_option {
        label := tool.Get_language(db, option, true)
        data_html += `<a href="/recent_change/1/` + option + `">(` + label + `)</a> `
    }

    api_data := Api_list_history(config, doc_name, set_type, num)
    api_data_list := api_data["data"].([][]string)

    history_ui, select_ui := Get_ui_history(db, config, api_data_list)

    data_html += history_ui
    data_html += tool.Get_page_control(
        db,
        tool.Str_to_int(num),
        len(api_data_list),
        50,
        "/history_page/{}/" + set_type + "/" + tool.Url_parser(doc_name),
    )

    data_html = `
        <form method="post">
            <select name="a">` + select_ui + `</select> 
            <select name="b">` + select_ui + `</select> 
            <button type="submit">` + tool.Get_language(db, "compare", true) + `</button>
        </form>
        <hr class="main_hr">
    ` + data_html

    out := tool.Get_template(
        db,
        config,
        doc_name,
        data_html,
        []any{ "(" + tool.Get_language(db, "history", true) + ")" + sub },
        [][]any{
            { "w/" + tool.Url_parser(doc_name), tool.Get_language(db, "return", true) },
        },
        map[string]string{},
    )

    return out
}
