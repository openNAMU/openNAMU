package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_w_set_reset(config tool.Config, doc_name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	ip := config.IP

	if tool.Check_permission(db, "owner", ip) {
		if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
			if _, err := tx.Exec(tool.DB_change("delete from acl where title = ?"), doc_name); err != nil {
				return err
			}
			if _, err := tx.Exec(tool.DB_change("delete from data_set where doc_name = ? and set_name = 'acl_date'"), doc_name); err != nil {
				return err
			}
			for _, set_name := range []string{"document_markup", "document_top", "document_editor_top"} {
				if _, err := tx.Exec(tool.DB_change("delete from data_set where doc_name = ? and set_name = ?"), doc_name, set_name); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			panic(err)
		}

		return_data := make(map[string]any)
		return_data["response"] = "ok"
		return_data["language"] = map[string]string{
			"reset": tool.Get_language(db, "reset", false),
		}

		return return_data
	} else {
		return_data := make(map[string]any)
		return_data["response"] = "require auth"
		return_data["language"] = map[string]string{
			"authority_error": tool.Get_language(db, "authority_error", false),
		}

		return return_data
	}
}
