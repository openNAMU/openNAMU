package route

import (
	"opennamu/route/tool"
)

func Document_set_list() map[string]string {
	set_list := map[string]string{}

	set_list["document_markup"] = ""
	set_list["document_top"] = ""
	set_list["document_editor_top"] = ""
	set_list["document_comment_code"] = ""

	return set_list
}

func Api_w_set(config tool.Config, doc_name string, set_name string, doc_rev string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	set_list := Document_set_list()
	return_data := make(map[string]any)

	if _, ok := set_list[set_name]; ok {
		return_data["data"] = tool.Get_document_setting(db, doc_name, set_name, doc_rev)
		return_data["response"] = "ok"
	} else {
		return_data["response"] = "not exist"
	}

	return return_data
}
