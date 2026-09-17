package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func View_setting_robot(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "setting_robot", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	return View_setting_robot_data(db, config, Setting_value(db, "robot", "", ""), Setting_value(db, "robot_default", "", ""))
}

func View_setting_robot_data(db *sql.DB, config tool.Config, value string, default_value string) string {
	data := `<a href="/robots.txt">(` + tool.Get_language(db, "view", true) + `)</a>` + Main_hr()
	data += `<form method="post">`
	data += `<textarea class="opennamu_textarea_500" name="content">` + tool.HTML_escape(value) + `</textarea>` + Main_hr()
	data += `<label><input type="checkbox" name="default" ` + Setting_checked(default_value) + `> ` + tool.Get_language(db, "default", true) + `</label>` + Main_hr()
	data += `<button id="opennamu_save_button" type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`

	return Setting_page(db, config, "robots.txt", data, "setting")
}
