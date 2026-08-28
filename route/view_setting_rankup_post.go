package route

import (
	"strings"

	"opennamu/route/tool"
)

func View_setting_rankup_post(config tool.Config, form map[string]string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	condition_map := map[string][]string{}
	for _, rankup_group := range tool.Rankup_group_list() {
		data := strings.TrimSpace(form[rankup_group])
		if data == "" {
			continue
		}

		condition_list, ok := tool.Rankup_condition_list(data)
		if !ok {
			return tool.Get_error_page(db, config, "error")
		}
		condition_map[rankup_group] = condition_list
	}

	for _, rankup_group := range tool.Rankup_group_list() {
		tool.Exec_DB(db, "delete from other where name = 'rankup_condition' and coverage = ?", rankup_group)
		for _, condition := range condition_map[rankup_group] {
			tool.Exec_DB(
				db,
				"insert into other (name, data, coverage) values ('rankup_condition', ?, ?)",
				condition,
				rankup_group,
			)
		}
	}
	tool.Do_insert_auth_history(db, config.IP, "rankup_condition")

	return tool.Get_redirect("/setting/rankup")
}
