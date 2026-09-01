package route

import "opennamu/route/tool"

func View_list_no_link(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	return list_extra_rows(db, config, tool.Get_language(db, "no_link_document", true), tool.Get_no_link_rows(db))
}
