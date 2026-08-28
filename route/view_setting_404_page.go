package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func View_setting_404_page(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "setting_404", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	page := setting_value(db, "manage_404_page", "", "404_page")
	content := setting_value(db, "manage_404_page_content", "", "")
	return view_setting_404_page_data(db, config, page, content, false)
}

func view_setting_404_page_data(db *sql.DB, config tool.Config, page string, content string, preview bool) string {
	lang := func(name string) string {
		return tool.Get_language(db, name, true)
	}

	options := ""
	for _, one := range [][]string{
		{"404_page", lang("404_file")},
		{"404_file", lang("404_page")},
	} {
		selected := ""
		if page == one[0] {
			selected = ` selected="selected"`
		}
		options += `<option value="` + one[0] + `"` + selected + `>` + tool.HTML_escape(one[1]) + `</option>`
	}

	data := `<form method="post">`
	data += `<select name="page">` + options + `</select>` + setting_hr()
	data += `<textarea class="opennamu_textarea_500" name="content">` + tool.HTML_escape(content) + `</textarea>` + setting_hr()
	data += `<button id="opennamu_save_button" name="action" value="save" type="submit">` + lang("save") + `</button> `
	data += `<button name="action" value="preview" type="submit">` + lang("preview") + `</button></form>`

	if preview {
		data += setting_hr() + `<div id="opennamu_setting_404_page_preview">` + content + `</div>`
	}

	return setting_page(db, config, lang("404_page_setting"), data, "setting")
}
