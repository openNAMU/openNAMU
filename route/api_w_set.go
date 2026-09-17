package route

import (
	"database/sql"
	"strings"

	"opennamu/route/tool"
)

func Document_set_list() map[string]string {
	set_list := map[string]string{}
	set_list["document_markup"] = ""
	set_list["document_top"] = ""
	set_list["document_top_markup"] = ""
	set_list["document_editor_top"] = ""
	set_list["document_editor_top_markup"] = ""
	set_list["document_comment_code"] = ""
	return set_list
}

func Document_set_read_list() map[string]string {
	set_list := Document_set_list()
	set_list["last_edit"] = ""
	return set_list
}

func Api_w_set_put(config tool.Config, doc_name string, set_name string, data string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	set_list := Document_set_list()
	if _, ok := set_list[set_name]; !ok {
		return map[string]any{"response": "not exist"}
	}
	if doc_name == "" {
		return map[string]any{"response": "error", "data": "invalid document"}
	}

	allowed := tool.Check_permission(db, "document_acl_manage", config.IP)
	if strings.HasPrefix(doc_name, "user:") && strings.TrimPrefix(doc_name, "user:") == config.IP {
		allowed = true
	}
	if !allowed {
		return map[string]any{"response": "require auth"}
	}

	if tool.Arr_in_str([]string{"document_top", "document_top_markup", "document_editor_top", "document_editor_top_markup"}, set_name) && !tool.Check_permission(db, "owner", config.IP) {
		return map[string]any{"response": "require auth"}
	}

	if set_name == "document_markup" || set_name == "document_top_markup" || set_name == "document_editor_top_markup" {
		switch data {
		case "custom":
			data = "html"
		case "raw":
			data = "plain"
		}
		if data != "" && !tool.Arr_in_str(Setting_markup_options(), data) {
			return map[string]any{"response": "error", "data": "invalid data"}
		}
	}

	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		if _, err := tx.Exec(tool.DB_change("delete from data_set where doc_name = ? and doc_rev = '' and set_name = ?"), doc_name, set_name); err != nil {
			return err
		}
		_, err := tx.Exec(tool.DB_change("insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', ?, ?)"), doc_name, set_name, data)
		return err
	}); err != nil {
		panic(err)
	}

	return map[string]any{"response": "ok"}
}

func Api_w_set(config tool.Config, doc_name string, set_name string, doc_rev string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	set_list := Document_set_read_list()
	return_data := make(map[string]any)

	if _, ok := set_list[set_name]; ok {
		return_data["data"] = tool.Get_document_setting_value(db, doc_name, set_name, doc_rev)
		return_data["response"] = "ok"
	} else {
		return_data["response"] = "not exist"
	}

	return return_data
}
