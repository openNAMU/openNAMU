package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_bbs_w_set_put(config tool.Config, set_id string, set_name string, data string, coverage string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	auth_info := tool.Check_permission(db, "bbs_setting", config.IP)

	setting_acl := BBS_w_set_list()
	return_data := make(map[string]any)

	if _, ok := setting_acl[set_name]; ok {
		if tool.Arr_in_str(bbs_set_fields, set_name) && !acl_value_valid(db, data) {
			return_data["response"] = "error"
			return_data["data"] = "invalid acl"
			return return_data
		}
		if auth_info {
			if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
				if coverage == "" {
					if _, err := tx.Exec(tool.DB_change("delete from bbs_set where set_name = ? and set_id = ?"), set_name, set_id); err != nil {
						return err
					}
				} else if _, err := tx.Exec(tool.DB_change("delete from bbs_set where set_name = ? and set_code = ? and set_id = ?"), set_name, coverage, set_id); err != nil {
					return err
				}
				if _, err := tx.Exec(tool.DB_change("insert into bbs_set (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)"), set_name, coverage, set_id, data); err != nil {
					return err
				}
				tool.Do_insert_auth_history(tx, config.IP, "bbs_set ("+set_id+"/"+set_name+")")
				return nil
			}); err != nil {
				panic(err)
			}

			return_data["response"] = "ok"
		} else {
			return_data["response"] = "require auth"
		}
	} else {
		return_data["response"] = "not exist"
	}

	return return_data
}
