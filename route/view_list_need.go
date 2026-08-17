package route

import "opennamu/route/tool"

func View_list_need(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	return list_extra_query(db, config, tool.Get_language(db, "need_document", true), "select title, link from back where type = 'no' order by title")
}
