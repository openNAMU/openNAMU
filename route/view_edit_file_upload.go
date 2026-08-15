package route

import (
	"opennamu/route/tool"
)

func View_edit_file_upload(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "upload", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	license_html := `<option value="direct_input">` + tool.Get_language(db, "direct_input", false) + `</option>`
	rows := tool.Query_DB(db, "select html from html_filter where kind = 'image_license'")
	for rows.Next() {
		license := ""
		if rows.Scan(&license) == nil {
			license_html += `<option value="` + tool.HTML_escape(license) + `">` + tool.HTML_escape(license) + `</option>`
		}
	}
	rows.Close()

	data_html := `<form method="post" enctype="multipart/form-data" accept-charset="utf8">`
	data_html += `<input class="__ON_INPUT__" multiple="multiple" type="file" name="f_data[]" id="file_input">`
	data_html += `<hr class="main_hr">`
	data_html += `<input class="__ON_INPUT__" placeholder="` + tool.Get_language(db, "file_name", false) + `" name="f_name" value="">`
	data_html += `<hr class="main_hr">`
	data_html += `<select class="__ON_INPUT__" name="f_lice_sel">` + license_html + `</select>`
	data_html += `<hr class="main_hr"><textarea class="opennamu_textarea_100" placeholder="` + tool.Get_language(db, "other", false) + `" name="f_lice"></textarea>`
	data_html += `<hr class="main_hr">` + tool.Get_captcha_ui(db, config)
	data_html += `<button id="opennamu_save_button" type="submit">` + tool.Get_language(db, "save", false) + `</button>`
	data_html += `</form>`

	out := tool.Get_template(
		db,
		config,
		tool.Get_language(db, "upload", true),
		data_html,
		[]any{},
		[][]any{},
		map[string]string{},
	)

	return out
}
