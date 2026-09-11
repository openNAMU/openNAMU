package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_filter_delete_post(config tool.Config, kind string, name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	spec, ok := get_filter_spec(kind)
	if !ok {
		return_data["response"] = "error"
		return return_data
	}
	if !tool.Check_permission(db, "filter_manage", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		if _, err := tx.Exec(tool.DB_change("delete from html_filter where html = ? and kind = ?"), name, spec.db_kind); err != nil {
			return err
		}
		if kind == "inter_wiki" {
			if _, err := tx.Exec(tool.DB_change("delete from html_filter where html = ? and kind = 'inter_wiki_sub'"), name); err != nil {
				return err
			}
		}
		tool.Do_insert_auth_history(tx, config.IP, "filter_delete ("+kind+")")
		return nil
	}); err != nil {
		panic(err)
	}

	return_data["response"] = "ok"
	return return_data
}
