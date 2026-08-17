package route

import "opennamu/route/tool"

func View_list_no_link(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	return list_extra_query(db, config, tool.Get_language(db, "no_link_document", true), "select doc_name, set_data from data_set where set_name = 'link_count' and set_data = '0' order by doc_name")
}
