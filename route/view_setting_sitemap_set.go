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

	return View_setting_sitemap_set_data(db, config, Setting_load_fields(db, Setting_sitemap_fields()))
}

func Setting_sitemap_fields() []setting_field {
	return []setting_field{
		{name: "sitemap_auto_exclude_domain"},
		{name: "sitemap_auto_exclude_user_page"},
		{name: "sitemap_auto_exclude_file_page"},
		{name: "sitemap_auto_exclude_category_page"},
		{name: "sitemap_auto_make"},
		{name: "indexnow_key"},
	}
}

func View_setting_sitemap_set_data(db *sql.DB, config tool.Config, values map[string]string) string {
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

	data.WriteString(Main_hr() + `<form method="post">`)
	data.WriteString(`<a href="/setting/sitemap">(` + lang("sitemap_manual_create") + `)</a>` + Main_hr())
	data.WriteString(`<label><input type="checkbox" name="sitemap_auto_make" ` + Setting_checked(values["sitemap_auto_make"]) + `> ` + lang("sitemap_auto_make") + `</label>` + Main_hr())
	data.WriteString(`<span>` + lang("indexnow_key") + `</span>` + Main_hr())
	data.WriteString(`<sup>` + lang("indexnow_key_help") + `</sup>` + Main_hr())
	data.WriteString(Setting_input("indexnow_key", values["indexnow_key"], "text") + Main_hr())
	data.WriteString(`<label><input type="checkbox" name="sitemap_auto_exclude_domain" ` + Setting_checked(values["sitemap_auto_exclude_domain"]) + `> ` + lang("stiemap_exclude_domain") + `</label>` + Main_hr())
	data.WriteString(`<label><input type="checkbox" name="sitemap_auto_exclude_user_page" ` + Setting_checked(values["sitemap_auto_exclude_user_page"]) + `> ` + lang("stiemap_exclude_user_page") + `</label>` + Main_hr())
	data.WriteString(`<label><input type="checkbox" name="sitemap_auto_exclude_file_page" ` + Setting_checked(values["sitemap_auto_exclude_file_page"]) + `> ` + lang("stiemap_exclude_file_page") + `</label>` + Main_hr())
	data.WriteString(`<label><input type="checkbox" name="sitemap_auto_exclude_category_page" ` + Setting_checked(values["sitemap_auto_exclude_category_page"]) + `> ` + lang("stiemap_exclude_category_page") + `</label>` + Main_hr())
	data.WriteString(`<button id="opennamu_save_button" type="submit">` + lang("save") + `</button></form>`)

	return Setting_page(db, config, lang("sitemap_management"), data.String(), "setting")
}
