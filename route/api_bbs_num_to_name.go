package route

import (
	"database/sql"
	"opennamu/route/tool"
)

func Api_bbs_num_to_name(db *sql.DB, set_id string) map[string]any {
	bbs_name := ""

	tool.QueryRow_DB(
		db,
		"select set_data from bbs_set where set_id = ? and set_name = 'bbs_name'",
		[]any{&bbs_name},
		set_id,
	)

	return map[string]any{
		"response": "ok",
		"data":     bbs_name,
	}
}
