package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_w_page_view(config tool.Config, doc_name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return map[string]any{
		"response": "ok",
		"data":     Document_page_view_read(db, doc_name),
	}
}

func Document_page_view_read(db *sql.DB, doc_name string) int {
	pv_continue := tool.Get_setting(db, "not_use_view_count", "")
	if len(pv_continue) > 0 && pv_continue[0][0] != "" {
		return 0
	}

	page_view_str := ""
	tool.QueryRow_DB(
		db,
		"select set_data from data_set where doc_name = ? and set_name = 'view_count' and doc_rev = ''",
		[]any{&page_view_str},
		doc_name,
	)

	return tool.Str_to_int(page_view_str)
}
