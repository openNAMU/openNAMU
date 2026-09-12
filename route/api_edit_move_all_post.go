package route

import (
	"database/sql"
	"strings"

	"opennamu/route/tool"
)

func move_bbs_index_update(db *sql.DB, doc_name string) {
	rows := tool.Query_DB(
		db,
		"select distinct set_id, set_code from bbs_data where (set_id = '0' and set_name = 'title' and set_data = ?) or (set_id = ? and set_name in ('document', 'tag') and set_data = ?)",
		doc_name,
		thread_bbs_id,
		doc_name,
	)

	set_list := [][2]string{}
	for rows.Next() {
		set_data := [2]string{}
		if rows.Scan(&set_data[0], &set_data[1]) == nil {
			set_list = append(set_list, set_data)
		}
	}
	rows.Close()

	for _, set_data := range set_list {
		tool.Search_bbs_index_update(db, set_data[0], set_data[1])
	}
}

func Api_edit_move_all_post(config tool.Config, source string, target string, match_type string, send string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "document_move_manage", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	documents := move_all_documents(db, source, target, match_type)
	if len(documents) == 0 {
		return_data["response"] = "error"
		return_data["data"] = "move error"
		return return_data
	}

	target_list := map[string]bool{}
	for _, document := range documents {
		if !tool.Check_acl(db, document.old_name, "", "document_move", config.IP) || !tool.Check_acl(db, document.new_name, "", "document_move", config.IP) {
			return_data["response"] = "require auth"
			return return_data
		}
		if document.old_name == document.new_name || strings.ContainsAny(document.new_name, "\r\n") || tool.Get_len(document.new_name) > 256 {
			return_data["response"] = "error"
			return_data["data"] = "invalid document"
			return return_data
		}
		if !tool.Do_title_length_check(db, document.new_name, "document") {
			return_data["response"] = "error"
			return_data["data"] = "title length"
			return return_data
		}
		if target_list[document.new_name] {
			return_data["response"] = "error"
			return_data["data"] = "document already exist"
			return return_data
		}
		target_list[document.new_name] = true

		target_exists, target_history_only := move_document_exists(db, document.new_name)
		if target_exists {
			if target_history_only {
				return_data["response"] = "error"
				return_data["data"] = "document history exist"
				return return_data
			}
			return_data["response"] = "error"
			return_data["data"] = "document already exist"
			return return_data
		}
		if move_document_settings_exists(db, document.new_name) {
			return_data["response"] = "error"
			return_data["data"] = "move error"
			return return_data
		}
	}

	for index := range documents {
		documents[index].data = move_data_value(db, documents[index].old_name)
	}

	date := tool.Get_time()
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		for _, document := range documents {
			move_document_normal_data(tx, document.old_name, document.new_name)
			move_topic_normal(tx, document.old_name, document.new_name)
			move_data_set_normal(tx, document.old_name, document.new_name)
			type_check := "<a>" + tool.HTML_escape(document.old_name) + "</a> → <a>" + tool.HTML_escape(document.new_name) + "</a>"
			tool.Do_add_history(tx, document.new_name, document.data, date, config.IP, send, "0", "move", type_check)
		}
		return nil
	}); err != nil {
		return_data["response"] = "error"
		return_data["data"] = "move error"
		return return_data
	}

	for _, document := range documents {
		tool.Search_index_sync(db, document.old_name)
		tool.Search_index_sync(db, document.new_name)
		move_bbs_index_update(db, document.new_name)
	}
	_ = tool.Search_bbs_index_mark_rebuild()

	return_data["response"] = "ok"
	return_data["data"] = len(documents)
	return return_data
}
