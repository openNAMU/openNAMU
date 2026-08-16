package route

import (
	"opennamu/route/tool"
)

func View_edit(config tool.Config, doc_name string, load_doc_name string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Do_title_length_check(db, doc_name, "document") {
		return tool.Get_error_page(db, config, "title length")
	}

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

	revision := tool.Get_document_revision(db, doc_name)
	form_data := editor_top_text + `<form action="/edit/` + tool.Url_parser(doc_name) + `" method="post">
        <input type="hidden" name="ver" value="` + tool.HTML_escape(revision) + `">
        <input class="__ON_INPUT__" type="text" name="send" placeholder="` + tool.Get_language(db, "why", true) + `">
        <hr class="main_hr">
        ` + tool.Get_editor_ui(db, config, raw_data_get, "edit", check_box+bottom_text, doc_name) + `
    </form>`

	out := tool.Get_template(
		db,
		config,
		doc_name,
		form_data,
		[]any{"(" + tool.Get_language(db, "edit", true) + ")"},
		[][]any{
			{"w/" + tool.Url_parser(doc_name), tool.Get_language(db, "return", true)},
			{"delete/" + tool.Url_parser(doc_name), tool.Get_language(db, "delete", true)},
			{"move/" + tool.Url_parser(doc_name), tool.Get_language(db, "move", true)},
		},
		map[string]string{},
	)

	return out
}
