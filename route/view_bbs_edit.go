package route

import "opennamu/route/tool"

func View_bbs_edit(config tool.Config, set_id string, set_code string, comment_code string, document_data string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	document := ""
	if document_data != "" {
		decoded, err := tool.Get_base64_decode(document_data)
		if err != nil || set_id != thread_bbs_id || set_code != "" || comment_code != "" {
			return tool.Get_redirect("/bbs/main")
		}
		document = decoded
		if !tool.Check_acl(db, document, "", "render", config.IP) {
			return tool.Get_error_page(db, config, "auth")
		}
	}

	api_data := Api_bbs_w_edit_view(config, set_id, set_code, comment_code)
	response, _ := api_data["response"].(string)
	if response != "ok" {
		if response == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}

		return tool.Get_redirect("/bbs/main")
	}

	data, ok := api_data["data"].(map[string]string)
	if !ok {
		return tool.Get_redirect("/bbs/main")
	}

	bbs_placeholder := bbs_set_value(db, set_id, "bbs_placeholder")
	bbs_comment_placeholder := bbs_set_value(db, set_id, "bbs_comment_placeholder")
	editor_type := "bbs"
	editor_help := bbs_placeholder
	if comment_code != "" {
		editor_type = "bbs_comment"
		editor_help = bbs_comment_placeholder
	}

	path := "/bbs/edit/" + tool.Url_parser(set_id)
	if document != "" {
		path += "/document/" + tool.Base64_encode(document)
	} else if set_code != "" {
		path += "/" + tool.Url_parser(set_code)
	}
	if comment_code != "" {
		path += "/" + tool.Url_parser(comment_code)
	}

	title := tool.Get_language(db, "post_edit", true)
	if comment_code != "" {
		title = tool.Get_language(db, "bbs_comment_edit", true)
	} else if set_code == "" {
		title = tool.Get_language(db, "post_add", true)
	}

	title_style := ""
	if comment_code != "" {
		title_style = ` style="display: none;"`
	}

	prefix_html := ""
	tag_html := ""
	if comment_code == "" {
		prefix_list := bbs_prefix_list(db, set_id)
		if len(prefix_list) > 0 {
			prefix_html = "<select class=\"__ON_INPUT__\" name=\"prefix\">"
			if set_id != thread_bbs_id {
				prefix_html += "<option value=\"\">" + tool.Get_language(db, "empty", true) + "</option>"
			}
			for _, prefix := range prefix_list {
				selected := ""
				if prefix == data["prefix"] {
					selected = " selected"
				}
				prefix_html += "<option value=\"" + tool.HTML_escape(prefix) + "\"" + selected + ">" + tool.HTML_escape(prefix) + "</option>"
			}
			prefix_html += "</select><hr class=\"main_hr\">"
		}
		tag_html = `<input class="__ON_INPUT__" placeholder="` + tool.Get_language(db, "tag", true) + `" name="tags" value="` + tool.HTML_escape(data["tags"]) + `"><hr class="main_hr">`
	}

	document_html := ""
	if document != "" {
		document_html = `<div>` + tool.Get_language(db, "document", true) + `: ` + tool.HTML_escape(document) + `</div><hr class="main_hr">`
	}

	data_html := `<a href="/filter/edit_filter">(` + tool.Get_language(db, "edit_filter_rule", true) + `)</a><hr class="main_hr">` + document_html + `
        <form action="` + path + `" method="post">
            <input class="__ON_INPUT__"` + title_style + ` placeholder="` + tool.Get_language(db, "title", true) + `" name="title" value="` + tool.HTML_escape(data["title"]) + `">
            <hr` + title_style + ` class="main_hr">
            ` + prefix_html + tag_html + tool.Get_editor_ui(db, config, data["data"], editor_type, "", "", editor_help) + `
        </form>`

	return tool.Get_template(
		db,
		config,
		title,
		data_html,
		[]any{},
		[][]any{
			{"bbs/in/" + tool.Url_parser(set_id), tool.Get_language(db, "return", true)},
		},
		map[string]string{},
	)
}
