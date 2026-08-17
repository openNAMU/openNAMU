package route

import "opennamu/route/tool"

func View_setting_head_post(config tool.Config, kind string, skin_name string, content string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	name, coverage, _, action, ok := setting_head_info(kind, skin_name)
	if !ok {
		return tool.Get_error_page(db, config, "error")
	}

	setting_save_value(db, name, coverage, content)
	tool.Do_insert_auth_history(db, config.IP, "edit_set ("+name+")")

	if skin_name != "" {
		action += "/" + tool.Url_parser(skin_name)
	}

	return tool.Get_redirect(action)
}
