package route

import "opennamu/route/tool"

func View_setting(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	set_list := [][]string{
		{"main", tool.Get_language(db, "main_setting", true)},
		{"phrase", tool.Get_language(db, "text_setting", true)},
		{"robot", "robots.txt"},
		{"external", tool.Get_language(db, "ext_api_req_set", true)},
		{"head", tool.Get_language(db, "main_head", false)},
		{"body/top", tool.Get_language(db, "main_body", true)},
		{"body/bottom", tool.Get_language(db, "main_bottom_body", true)},
		{"sitemap_set", tool.Get_language(db, "sitemap_management", true)},
		{"top_menu", tool.Get_language(db, "top_menu_setting", true)},
		{"skin_set", tool.Get_language(db, "main_skin_set_default", true)},
		{"404_page", tool.Get_language(db, "404_page_setting", true)},
		{"backlink_reset", tool.Get_language(db, "reset_all_backlink", true)},
	}

	set_data := "<ul>"
	for _, li := range set_list {
		set_data += `<li><a href="/setting/` +
			li[0] +
			`">` +
			li[1] +
			`</a></li>`
	}
	
	set_data += "</ul>"

	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "setting", true),
		set_data,
		[]any{},
		[][]any{
			{"manager", tool.Get_language(db, "return", true)},
		},
		map[string]string{},
	)
}
