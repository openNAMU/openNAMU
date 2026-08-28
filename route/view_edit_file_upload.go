package route

import (
	"strconv"

	"opennamu/route/tool"
)

func View_edit_file_upload(config tool.Config, file_name string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "upload", config.IP) {
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

	upload_help := ""
	tool.QueryRow_DB(db, "select data from other where name = 'upload_help'", []any{&upload_help})

	upload_default := ""
	tool.QueryRow_DB(db, "select data from other where name = 'upload_default'", []any{&upload_default})

	file_max_size := tool.Get_file_max_size(db)
	if file_max_size <= 0 {
		file_max_size = 2
	}

	data_html := `<a href="/filter/file_filter">(` + tool.Get_language(db, "file_filter_list", true) + `)</a> <a href="/filter/extension_filter">(` + tool.Get_language(db, "extension_filter_list", true) + `)</a>`
	if upload_help != "" {
		data_html += `<hr class="main_hr">` + upload_help
	}
	data_html += `<hr class="main_hr">` + tool.Get_language(db, "max_file_size", true) + ` : ` + strconv.Itoa(file_max_size) + `MB`
	data_html += `<hr class="main_hr"><form method="post" enctype="multipart/form-data" accept-charset="utf8">`
	data_html += `<input class="__ON_INPUT__" multiple="multiple" type="file" name="f_data[]" id="file_input">`
	data_html += `<hr class="main_hr">`
	data_html += `<input class="__ON_INPUT__" placeholder="` + tool.Get_language(db, "file_name", false) + `" name="f_name" value="` + tool.HTML_escape(file_name) + `">`
	data_html += `<hr class="main_hr">`
	data_html += `<select class="__ON_INPUT__" name="f_lice_sel">` + license_html + `</select>`
	data_html += `<hr class="main_hr"><textarea class="opennamu_textarea_100" placeholder="` + tool.Get_language(db, "other", false) + `" name="f_lice">` + tool.HTML_escape(upload_default) + `</textarea>`
	data_html += `<hr class="main_hr">` + tool.Get_captcha_ui(db, config)
	data_html += `<button id="opennamu_save_button" type="submit">` + tool.Get_language(db, "save", false) + `</button>`
	data_html += `</form>`

	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "upload", true),
		data_html,
		[]any{},
		[][]any{{"other", tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}
