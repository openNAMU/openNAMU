package route

import (
	"database/sql"
	"strings"

	"opennamu/route/tool"
)

func Api_setting_rankup_post(config tool.Config, form map[string]string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "rankup_manage", config.IP) {
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

	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		for _, rankup_group := range tool.Rankup_group_list() {
			if _, err := tx.Exec(tool.DB_change("delete from other where name = 'rankup_condition' and coverage = ?"), rankup_group); err != nil {
				return err
			}
			for _, condition := range condition_map[rankup_group] {
				if _, err := tx.Exec(tool.DB_change("insert into other (name, data, coverage) values ('rankup_condition', ?, ?)"), condition, rankup_group); err != nil {
					return err
				}
			}
		}
		tool.Do_insert_auth_history(tx, config.IP, "rankup_condition")
		return nil
	}); err != nil {
		panic(err)
	}
	return_data["response"] = "ok"
	return return_data
}
