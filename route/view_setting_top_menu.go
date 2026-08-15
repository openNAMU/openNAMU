package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func View_setting_top_menu(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	return view_setting_top_menu_data(db, config, setting_value(db, "top_menu", "", ""))
}

func View_setting_top_menu_post(config tool.Config, form map[string]string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	setting_save_value(db, "top_menu", "", setting_form_value(form, "content", ""))
	tool.Do_insert_auth_history(db, config.IP, "edit_set (top_menu)")

	return tool.Get_redirect("/setting/top_menu")
}

func view_setting_top_menu_data(db *sql.DB, config tool.Config, value string) string {
	data := `<span>
        EX)
        <br>
        ONTS
        <br>
        https://2du.pythonanywhere.com/
        <br>
        FrontPage
        <br>
        /w/FrontPage
    </span>`
	data += setting_hr() + tool.Get_language(db, "not_support_skin_warning", true) + setting_hr()
	data += `<form method="post">`
	data += `<textarea class="opennamu_textarea_500" placeholder="` + tool.Get_language(db, "enter_top_menu_setting", true) + `" name="content" id="content">` + tool.HTML_escape(value) + `</textarea>`
	data += setting_hr() + `<button id="opennamu_save_button" type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`

	return setting_page(db, config, tool.Get_language(db, "top_menu_setting", true), data, "setting")
}
