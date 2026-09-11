package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func change_bbs_tabom_count(tx *sql.Tx, set_name string, set_id string, set_code string, amount int) {
	result, err := tx.Exec(
		tool.DB_change("update bbs_data set set_data = case when set_data + ? < 0 then 0 else set_data + ? end where set_name = ? and set_id = ? and set_code = ?"),
		amount,
		amount,
		set_name,
		set_id,
		set_code,
	)
	if err != nil {
		panic(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}
	if rows == 0 && amount > 0 {
		if _, err := tx.Exec(
			tool.DB_change("insert into bbs_data (set_name, set_data, set_id, set_code) values (?, ?, ?, ?)"),
			set_name,
			amount,
			set_id,
			set_code,
		); err != nil {
			panic(err)
		}
	}
}

func bbs_tabom_user_exists(db tool.DB_runner, set_name string, user string, set_id string, set_code string) bool {
	data := ""
	return tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = ? and set_data = ? and set_id = ? and set_code = ?",
		[]any{&data},
		set_name,
		user,
		set_id,
		set_code,
	)
}

func Api_bbs_w_tabom_post(config tool.Config, set_id string, set_code string, vote_type string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)

	if _, allowed := bbs_post_view_auth(db, set_id, set_code, config.IP); !allowed {
		return_data["response"] = "require auth"
		return return_data
	}

	if !tool.Check_acl(db, set_id, "", "bbs_comment", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	var result map[string]any
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		result = api_bbs_tabom_post(tx, config.IP, set_id, set_code, vote_type)
		return nil
	}); err != nil {
		panic(err)
	}
	return result
}
