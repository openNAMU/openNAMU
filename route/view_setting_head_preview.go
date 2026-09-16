package route

import "opennamu/route/tool"

func View_setting_head_preview(config tool.Config, kind string, content string, markup_name string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "setting_head", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	name, coverage, title_key, action, ok := setting_head_info(kind, "")
	if !ok || (kind != "body/top" && kind != "body/bottom") {
		return tool.Get_error_page(db, config, "error")
	}

	markup_name = setting_markup_normalize(markup_name)
	return view_setting_head_data(db, config, kind, "", name, coverage, title_key, action, setting_value(db, name, coverage, ""), markup_name, content, true)
}
