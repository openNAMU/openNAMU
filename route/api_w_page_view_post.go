package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_w_page_view_post(config tool.Config, doc_name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	pv_continue := tool.Get_setting(db, "not_use_view_count", "")
	if len(pv_continue) == 0 || pv_continue[0][0] == "" {
		if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
			increase_count := func(doc_rev string) error {
				result, err := tx.Exec(
					tool.DB_change("update data_set set set_data = set_data + 1 where doc_name = ? and set_name = 'view_count' and doc_rev = ?"),
					doc_name,
					doc_rev,
				)
				if err != nil {
					return err
				}
				rows, err := result.RowsAffected()
				if err != nil {
					return err
				}
				if rows > 0 {
					return nil
				}
				_, err = tx.Exec(
					tool.DB_change("insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, ?, 'view_count', '1')"),
					doc_name,
					doc_rev,
				)
				return err
			}

			if err := increase_count(""); err != nil {
				return err
			}

			for _, now_date := range []string{tool.Get_month(), tool.Get_date()} {
				if err := increase_count(now_date); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			panic(err)
		}

	}

	return_data := make(map[string]any)
	return_data["response"] = "ok"

	return return_data
}
