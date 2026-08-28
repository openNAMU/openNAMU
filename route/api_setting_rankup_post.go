package route

import (
	"strings"

	"opennamu/route/tool"
)

func Api_setting_rankup_post(config tool.Config, form map[string]string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	condition_map := map[string][]string{}
	for _, rankup_group := range tool.Rankup_group_list() {
		data := strings.TrimSpace(form[rankup_group])
		if data == "" {
			continue
		}
		condition_list, ok := tool.Rankup_condition_list(data)
		if !ok {
			return_data["response"] = "error"
			return_data["data"] = "error"
			return return_data
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
	return_data["response"] = "ok"
	return return_data
}
