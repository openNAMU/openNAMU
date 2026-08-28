package route

import (
	"database/sql"
	"strings"

	"opennamu/route/tool"
)

func View_setting_main_logo(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "setting_main_logo", config.IP) {
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

func setting_logo_skins() []string {
	return append([]string{"default"}, tool.Get_skin_list("", false)...)
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
