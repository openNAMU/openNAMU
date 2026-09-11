package route

import (
	"database/sql"
	"strconv"

	"opennamu/route/tool"
)

func Api_bbs_w_post(config tool.Config, set_id string, title string, data string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "bbs_comment", config.IP) {
		return_data := make(map[string]any)
		return_data["response"] = "require auth"

		return return_data
	}

	set_code_str := ""
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		last_code := ""
		tool.QueryRow_DB(
			tx,
			"select set_code from bbs_data where set_name = 'title' and set_id = ? order by set_code + 0 desc",
			[]any{&last_code},
			set_id,
		)
		set_code_str = strconv.Itoa(tool.Str_to_int(last_code) + 1)
		date_now := tool.Get_time()

		for _, value := range [][]string{
			{"title", title},
			{"data", data},
			{"date", date_now},
			{"last_activity", date_now},
			{"user_id", config.IP},
			{"comment_count", "0"},
		} {
			if _, err := tx.Exec(
				tool.DB_change("insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)"),
				value[0],
				set_code_str,
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
	tool.Search_bbs_index_update(db, set_id, set_code_str)

	return_data := make(map[string]any)
	return_data["response"] = "ok"
	return_data["data"] = set_code_str

	return return_data
}
