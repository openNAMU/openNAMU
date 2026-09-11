package route

import (
	"database/sql"
	"strconv"

	"opennamu/route/tool"
)

func Api_bbs_make(config tool.Config, bbs_name string, bbs_type string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)

	if !tool.Check_permission(db, "bbs_create", config.IP) {
		return_data["response"] = "require auth"

		return return_data
	}

	if !tool.Arr_in_str([]string{"comment", "thread"}, bbs_type) {
		bbs_type = "comment"
	}
	set_id := ""
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		last_id := "1"
		tool.QueryRow_DB(
			tx,
			`select set_id from bbs_set where set_name = "bbs_name" order by set_id + 0 desc`,
			[]any{&last_id},
		)
		set_id = strconv.Itoa(tool.Str_to_int(last_id) + 1)
		for _, value := range [][]any{
			{"bbs_name", bbs_name},
			{"bbs_type", bbs_type},
		} {
			if _, err := tx.Exec(
				tool.DB_change("insert into bbs_set (set_name, set_code, set_id, set_data) values (?, '', ?, ?)"),
				value[0],
				set_id,
				value[1],
			); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}

	return_data["response"] = "ok"

	return return_data
}
