package route

import "opennamu/route/tool"

func View_setting_sitemap(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	data := `<form method="post"><button id="opennamu_save_button" type="submit">` + tool.Get_language(db, "create", true) + `</button></form>`
	return setting_page(db, config, tool.Get_language(db, "sitemap_manual_create", true), data, "setting/sitemap_set")
}
