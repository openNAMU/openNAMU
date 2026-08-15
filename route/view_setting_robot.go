package route

import (
	"database/sql"
	"os"
	"strings"

	"opennamu/route/tool"
)

func View_setting_robot(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	return view_setting_robot_data(db, config, setting_value(db, "robot", "", ""), setting_value(db, "robot_default", "", ""))
}

func View_setting_robot_post(config tool.Config, form map[string]string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	setting_save_value(db, "robot", "", setting_form_value(form, "content", ""))
	setting_save_value(db, "robot_default", "", setting_form_value(form, "default", ""))
	tool.Do_insert_auth_history(db, config.IP, "edit_set (robot)")

	return tool.Get_redirect("/setting/robot")
}

func view_setting_robot_data(db *sql.DB, config tool.Config, value string, default_value string) string {
	data := `<a href="/robots.txt">(` + tool.Get_language(db, "view", true) + `)</a>` + setting_hr()
	data += `<form method="post">`
	data += `<textarea class="opennamu_textarea_500" name="content">` + tool.HTML_escape(value) + `</textarea>` + setting_hr()
	data += `<label><input type="checkbox" name="default" ` + setting_checked(default_value) + `> ` + tool.Get_language(db, "default", true) + `</label>` + setting_hr()
	data += `<button id="opennamu_save_button" type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`

	return setting_page(db, config, "robots.txt", data, "setting")
}

func View_robots_txt(config tool.Config, request_host string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	robot_default := setting_value(db, "robot_default", "", "")
	robot := setting_value(db, "robot", "", "")
	if robot_default == "" && robot != "" {
		return robot
	}

	domain := tool.Get_domain(db, true)
	if domain == "" || domain == "http://" || domain == "https://" {
		domain = "http://" + request_host
	}

	data := strings.Builder{}
	data.WriteString("User-agent: *\n")
	data.WriteString("Disallow: /\n")
	data.WriteString("Allow: /$\n")
	data.WriteString("Allow: /w/\n")
	data.WriteString("Allow: /bbs/w/\n")
	data.WriteString("Allow: /sitemap.xml$\n")
	data.WriteString("Allow: /sitemap_*.xml$")
	if _, err := os.Stat("sitemap.xml"); err == nil {
		data.WriteString("\nSitemap: " + domain + "/sitemap.xml")
	}

	return data.String()
}
