package route

import (
	"database/sql"
	"net/url"

	"opennamu/route/tool"
)

func Api_auth_group_post(config tool.Config, name string, values url.Values) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "auth_group_manage", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if name == "" {
		return_data["response"] = "error"
		return return_data
	}

	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		if _, err := tx.Exec(tool.DB_change("delete from alist where name = ?"), name); err != nil {
			return err
		}
		for _, choice := range tool.Auth_choices() {
			if values.Get(choice.Key) != "" {
				if _, err := tx.Exec(tool.DB_change("insert into alist (name, acl) values (?, ?)"), name, choice.Key); err != nil {
					return err
				}
			}
		}
		if _, err := tx.Exec(tool.DB_change("insert into alist (name, acl) values (?, 'nothing')"), name); err != nil {
			return err
		}
		tool.Do_insert_auth_history(tx, config.IP, "auth_group_save ("+name+")")
		return nil
	}); err != nil {
		panic(err)
	}

	return_data["response"] = "ok"
	return return_data
}
