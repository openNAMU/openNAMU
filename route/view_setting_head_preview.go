package route

import "opennamu/route/tool"

func View_setting_head_preview(config tool.Config, kind string, content string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	name, coverage, title_key, action, ok := setting_head_info(kind, "")
	if !ok || (kind != "body/top" && kind != "body/bottom") {
		return tool.Get_error_page(db, config, "error")
	}

	return view_setting_head_data(db, config, kind, "", name, coverage, title_key, action, setting_value(db, name, coverage, ""), content, true)
}
