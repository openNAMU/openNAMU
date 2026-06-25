package route

import "opennamu/route/tool"

func View_w_raw(config tool.Config, doc_name string, rev string, do_type string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    sub := "(" + tool.Get_language(db, "raw", true) + ")"
    if rev != "" {
        sub += " (" + rev + ")"
    }

    var menu [][]any

    if rev != "" {
        menu = [][]any{
            { "history_tool/" + tool.Url_parser(rev) + "/" + tool.Url_parser(doc_name), tool.Get_language(db, "return", true) },
        }
    } else {
        menu = [][]any{
            { "w/" + tool.Url_parser(doc_name), tool.Get_language(db, "return", true) },
        }
    }

    data_api_in := ""

    data_api := Api_w_raw(config, doc_name, "", rev)
    if data_api["response"].(string) == "ok" {
        data_api_in = data_api["data"].(string)
    }

    pre_data := `
        <div id="opennamu_preview_area">
            <textarea readonly id="opennamu_edit_textarea" class="opennamu_textarea_500 __ON_TEXTAREA__">` + tool.HTML_escape(data_api_in) + `</textarea>
        </div>
    `

    if do_type == "document_acl" {
        pre_data = `
            ` + tool.Get_language(db, "authority_error", true) + `
            <hr class="main_hr">
        ` + pre_data

        sub = "(" + tool.Get_language(db, "edit", true) + ")"
    }

    out := tool.Get_template(
        db,
        config,
        doc_name,
        pre_data,
        []any{ sub },
        menu,
        map[string]string{},
    )

    return out
}
