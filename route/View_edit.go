package route

import (
	"opennamu/route/tool"
)

func View_edit(config tool.Config, doc_name string, load_doc_name string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    var raw_data map[string]any
    raw_data_get := ""
    if load_doc_name == "" {
        raw_data = Api_w_raw(config, doc_name, "", "")
    } else {
        raw_data = Api_w_raw(config, load_doc_name, "", "")
    }

    if raw_data["response"].(string) == "ok" {
        raw_data_get = raw_data["data"].(string)
    }

    check_box := tool.Get_edit_check_box_ui(db)
    bottom_text := tool.Get_edit_bottom_text_ui(db, "edit")

    editor_top_text := ""
    if load_doc_name == "" {
        editor_top_text += `<a href="/manager/15/` + tool.Url_parser(doc_name) + `">(` + tool.Get_language(db, "load", true) + `)</a> `
    }

    if editor_top_text != "" {
        editor_top_text += `<hr class="main_hr">`
    }

    form_data := editor_top_text + `<form action="/edit/` + tool.Url_parser(doc_name) + `" method="post">
        <input class="__ON_INPUT__" type="text" name="send" placeholder="` + tool.Get_language(db, "why", true) + `">
        <hr class="main_hr">
        ` + tool.Get_editor_ui(db, config, raw_data_get, "edit", check_box + bottom_text, doc_name) + `
    </form>`

    out := tool.Get_template(
        db,
        config,
        doc_name,
        form_data,
        []any{ "(" + tool.Get_language(db, "edit", true) + ")" },
        [][]any{},
        map[string]string{},
    )

    return out
}
