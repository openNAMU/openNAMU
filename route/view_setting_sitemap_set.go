package route

import (
	"database/sql"
	"strconv"
	"strings"

	"opennamu/route/tool"
)

func View_setting_sitemap_set(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "setting_sitemap", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	return view_setting_sitemap_set_data(db, config, setting_load_fields(db, setting_sitemap_fields()))
}

func setting_sitemap_fields() []setting_field {
	return []setting_field{
		{name: "sitemap_auto_exclude_domain"},
		{name: "sitemap_auto_exclude_user_page"},
		{name: "sitemap_auto_exclude_file_page"},
		{name: "sitemap_auto_exclude_category_page"},
		{name: "sitemap_auto_make"},
	}
}

func view_setting_sitemap_set_data(db *sql.DB, config tool.Config, values map[string]string) string {
	lang := func(name string) string {
		return tool.Get_language(db, name, true)
	}

	data := strings.Builder{}
	if tool.File_exist_check("sitemap.xml") {
		data.WriteString(`<a href="/sitemap.xml">(` + lang("view") + `)</a>`)
		for i := 0; ; i++ {
			name := "sitemap_" + strconv.Itoa(i) + ".xml"
			if !tool.File_exist_check(name) {
				break
			}
			data.WriteString(` <a href="/` + name + `">(` + name + `)</a>`)
		}
	}

	data.WriteString(setting_hr() + `<form method="post">`)
	data.WriteString(`<a href="/setting/sitemap">(` + lang("sitemap_manual_create") + `)</a>` + setting_hr())
	data.WriteString(`<label><input type="checkbox" name="sitemap_auto_make" ` + setting_checked(values["sitemap_auto_make"]) + `> ` + lang("sitemap_auto_make") + `</label>` + setting_hr())
	data.WriteString(`<label><input type="checkbox" name="sitemap_auto_exclude_domain" ` + setting_checked(values["sitemap_auto_exclude_domain"]) + `> ` + lang("stiemap_exclude_domain") + `</label>` + setting_hr())
	data.WriteString(`<label><input type="checkbox" name="sitemap_auto_exclude_user_page" ` + setting_checked(values["sitemap_auto_exclude_user_page"]) + `> ` + lang("stiemap_exclude_user_page") + `</label>` + setting_hr())
	data.WriteString(`<label><input type="checkbox" name="sitemap_auto_exclude_file_page" ` + setting_checked(values["sitemap_auto_exclude_file_page"]) + `> ` + lang("stiemap_exclude_file_page") + `</label>` + setting_hr())
	data.WriteString(`<label><input type="checkbox" name="sitemap_auto_exclude_category_page" ` + setting_checked(values["sitemap_auto_exclude_category_page"]) + `> ` + lang("stiemap_exclude_category_page") + `</label>` + setting_hr())
	data.WriteString(`<button id="opennamu_save_button" type="submit">` + lang("save") + `</button></form>`)

	return setting_page(db, config, lang("sitemap_management"), data.String(), "setting")
}
