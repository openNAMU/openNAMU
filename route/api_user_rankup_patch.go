package route

import (
	"net/url"
	"strings"

	"opennamu/route/tool"
)

func Api_user_rankup_patch(config tool.Config, values url.Values) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "rankup_manage", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	coverage := strings.TrimSpace(values.Get("coverage"))
	if coverage == "" {
		coverage = strings.TrimSpace(values.Get("name"))
	}
	if coverage == "" {
		coverage = strings.TrimSpace(values.Get("rankup_name"))
	}
	if coverage == "" {
		return_data["response"] = "error"
		return_data["data"] = "invalid name"
		return return_data
	}
	if !tool.Auth_group_name_rankup(coverage) {
		return_data["response"] = "error"
		return_data["data"] = "invalid name"
		return return_data
	}

	condition := strings.TrimSpace(values.Get("data"))
	if condition == "" {
		condition_type := strings.TrimSpace(values.Get("type"))
		condition_value := strings.TrimSpace(values.Get("value"))
		if condition_type != "" && condition_value != "" {
			condition = condition_type + " " + condition_value
		}
	}

	if values.Get("delete") != "" {
		tool.Exec_DB(db, "delete from other where name = 'rankup_condition' and coverage = ?", coverage)
		tool.Do_insert_auth_history(db, config.IP, "rankup_condition ("+coverage+")")
		return_data["response"] = "ok"
		return_data["data"] = coverage
		return return_data
	}

	condition_list, ok := tool.Rankup_condition_list(condition)
	if !ok {
		return_data["response"] = "error"
		return_data["data"] = "invalid condition"
		return return_data
	}

	tool.Exec_DB(db, "delete from other where name = 'rankup_condition' and coverage = ?", coverage)
	for _, condition_data := range condition_list {
		tool.Exec_DB(db, "insert into other (name, data, coverage) values ('rankup_condition', ?, ?)", condition_data, coverage)
	}

	tool.Do_insert_auth_history(db, config.IP, "rankup_condition ("+coverage+")")
	return_data["response"] = "ok"
	return_data["data"] = coverage
	return return_data
}
