package route

import (
	"database/sql"
	"strings"

	"opennamu/route/tool"
)

func setting_logo_skins() []string {
	return append([]string{"default"}, tool.Get_skin_list("", false)...)
}

func View_setting_main_logo(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	values := map[string]string{}
	for _, skin := range setting_logo_skins() {
		coverage := ""
		field_name := "main_css"
		if skin != "default" {
			coverage = skin
			field_name = skin
		}

		values[field_name] = setting_value(db, "logo", coverage, "")
	}

	return view_setting_main_logo_data(db, config, values)
}

func View_setting_main_logo_post(config tool.Config, form map[string]string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	for _, skin := range setting_logo_skins() {
		coverage := ""
		field_name := "main_css"
		if skin != "default" {
			coverage = skin
			field_name = skin
		}

		setting_save_value(db, "logo", coverage, setting_form_value(form, field_name, ""))
	}
	tool.Do_insert_auth_history(db, config.IP, "edit_set (logo)")

	return tool.Get_redirect("/setting/main/logo")
}

func view_setting_main_logo_data(db *sql.DB, config tool.Config, values map[string]string) string {
	data := strings.Builder{}
	data.WriteString(`<form method="post">`)

	for _, skin := range setting_logo_skins() {
		field_name := "main_css"
		label := tool.Get_language(db, "wiki_logo", true)
		if skin != "default" {
			field_name = skin
			label += " (" + tool.HTML_escape(skin) + ")"
		}

		data.WriteString(`<span>` + label + ` (HTML)</span>` + setting_hr())
		data.WriteString(setting_input(field_name, values[field_name], "text") + setting_hr())
	}

	data.WriteString(`<button id="opennamu_save_button" type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`)

	return setting_page(db, config, tool.Get_language(db, "wiki_logo", true), data.String(), "setting/main")
}
