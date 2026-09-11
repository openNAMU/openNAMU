package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_bbs_w_page_view_post(config tool.Config, set_id string, set_code string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	return_data["response"] = "ok"
	if _, allowed := bbs_post_view_auth(db, set_id, set_code, config.IP); !allowed {
		return_data["response"] = "require auth"
		return return_data
	}

	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		result, err := tx.Exec(
			tool.DB_change("update bbs_data set set_data = set_data + 1 where set_name = 'view_count' and set_id = ? and set_code = ?"),
			set_id,
			set_code,
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
			tool.DB_change("insert into bbs_data (set_name, set_id, set_code, set_data) values ('view_count', ?, ?, ?)"),
			set_id,
			set_code,
			"1",
		)
		return err
	}); err != nil {
		panic(err)
	}

	return return_data
}
