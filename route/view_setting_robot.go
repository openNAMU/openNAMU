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

	return view_setting_robot_data(db, config, setting_value(db, "robot", "", ""), setting_value(db, "robot_default", "", ""))
}

func view_setting_robot_data(db *sql.DB, config tool.Config, value string, default_value string) string {
	data := `<a href="/robots.txt">(` + tool.Get_language(db, "view", true) + `)</a>` + setting_hr()
	data += `<form method="post">`
	data += `<textarea class="opennamu_textarea_500" name="content">` + tool.HTML_escape(value) + `</textarea>` + setting_hr()
	data += `<label><input type="checkbox" name="default" ` + setting_checked(default_value) + `> ` + tool.Get_language(db, "default", true) + `</label>` + setting_hr()
	data += `<button id="opennamu_save_button" type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`

	return setting_page(db, config, "robots.txt", data, "setting")
}
