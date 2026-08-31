package route

import "opennamu/route/tool"

func View_setting_rankup(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "rankup_manage", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	condition_map := map[string]string{}
	for _, condition_data := range tool.Get_setting(db, "rankup_condition", "") {
		if len(condition_data) < 2 || !tool.Auth_group_name_rankup(condition_data[1]) {
			continue
		}

		if condition_map[condition_data[1]] != "" {
			condition_map[condition_data[1]] += "\n"
		}
		condition_map[condition_data[1]] += condition_data[0]
	}

	lang := func(name string) string {
		return tool.Get_language(db, name, true)
	}

	data := `<form method="post"><p>` + lang("rankup_condition_help") + `</p>`
	for _, rankup_group := range tool.Rankup_group_list() {
		data += `<h3>` + lang(rankup_group+"_acl") + `</h3>`
		data += setting_textarea(rankup_group, condition_map[rankup_group], "opennamu_textarea_100") + setting_hr()
	}
	data += `<button id="opennamu_save_button" type="submit">` + lang("save") + `</button></form>`

	return setting_page(db, config, lang("rankup_setting"), data, "setting")
}
