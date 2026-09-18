package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func View_setting_top_menu(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "setting_top_menu", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	return View_setting_top_menu_data(db, config, Setting_value(db, "top_menu", "", ""))
}

func View_setting_top_menu_data(db *sql.DB, config tool.Config, value string) string {
	data := `<pre>EX)
ONTS
https://2du.pythonanywhere.com/
FrontPage
/w/FrontPage</pre>`
	data += Main_hr() + tool.Get_language(db, "not_support_skin_warning", true) + Main_hr()
	data += `<form method="post">`
	data += `<textarea class="opennamu_textarea_500" placeholder="` + tool.Get_language(db, "enter_top_menu_setting", true) + `" name="content" id="content">` + tool.HTML_escape(value) + `</textarea>`
	data += Main_hr() + `<button id="opennamu_save_button" type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`

	return Setting_page(db, config, tool.Get_language(db, "top_menu_setting", true), data, "setting")
}
