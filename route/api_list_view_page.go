package route

import "opennamu/route/tool"

func Api_list_view_page(config tool.Config, num string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	page_int := tool.Str_to_int(num)
	if page_int > 0 {
		page_int = (page_int * 50) - 50
	} else {
		page_int = 0
	}

	rows := tool.Query_DB(
		db,
		"select doc_name, set_data from data_set where set_name = 'view_count' and doc_rev = '' and "+tool.Get_except_document_name_SQL("doc_name")+" order by set_data + 0 desc, doc_name asc limit ?, 50",
		page_int,
	)
	defer rows.Close()

	data_list := [][]string{}
	for rows.Next() {
		doc_name := ""
		view_count := ""
		if rows.Scan(&doc_name, &view_count) != nil {
			continue
		}
		data_list = append(data_list, []string{doc_name, view_count})
	}

	return map[string]any{
		"response": "ok",
		"data":     data_list,
	}
}
