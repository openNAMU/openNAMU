package route

import (
	"opennamu/route/tool"
)

func Api_list_old_page(config tool.Config, num string, set_type string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	page_int := tool.Str_to_int(num)
	if page_int > 0 {
		page_int = (page_int * 50) - 50
	} else {
		page_int = 0
	}

	query := ""
	if set_type == "old" {
		query = "select doc_name, set_data, case when exists (select 1 from back where back.link = data_set.doc_name and back.type = 'redirect') then '1' else '' end from data_set where set_name = 'last_edit' and doc_rev = '' and " + tool.Get_except_document_name_SQL("doc_name") + " order by set_data asc limit ?, 50"
	} else {
		query = "select doc_name, set_data, case when exists (select 1 from back where back.link = data_set.doc_name and back.type = 'redirect') then '1' else '' end from data_set where set_name = 'last_edit' and doc_rev = '' and " + tool.Get_except_document_name_SQL("doc_name") + " order by set_data desc limit ?, 50"
	}

	rows := tool.Query_DB(
		db,
		query,
		page_int,
	)
	defer rows.Close()

	data_list := [][]string{}

	for rows.Next() {
		var doc_name string
		var date string
		var is_redirect string

		err := rows.Scan(&doc_name, &date, &is_redirect)
		if err != nil {
			panic(err)
		}

		data_list = append(data_list, []string{doc_name, date, is_redirect})
	}

	return_data := make(map[string]any)
	return_data["response"] = "ok"
	return_data["data"] = data_list

	return return_data
}
