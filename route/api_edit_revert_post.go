package route

import (
	"database/sql"
	"strings"

	"opennamu/route/tool"
	"opennamu/route/tool/markup"
)

func Api_edit_revert_post(config tool.Config, doc_name string, rev string, send string, agree string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)

	data := ""
	hide := ""
	exist := tool.QueryRow_DB(
		db,
		"select data, hide from history where title = ? and id = ?",
		[]any{&data, &hide},
		doc_name,
		rev,
	)

	data = strings.ReplaceAll(data, "\r", "")
	data = tool.Do_edit_replace(db, data)
	if !exist {
		return_data["response"] = "not exist"

		return return_data
	} else if hide != "" && !tool.Check_permission(db, "hidel", config.IP) {
		return_data["response"] = "require auth"

		return return_data
	} else if !tool.Check_acl(db, doc_name, "", "document_edit", config.IP) {
		return_data["response"] = "require auth"

		return return_data
	} else if !tool.Do_edit_slow_check(db, config, "edit") {
		return_data["response"] = "error"
		return_data["data"] = "slow edit limit"

		return return_data
	}
	if error_name := Edit_post_content_validation_error(db, config, doc_name, data, send, agree); error_name != "" {
		return map[string]any{"response": "error", "data": error_name}
	}

	// The revert (content commit) and its backlink index update are one critical
	// section. The common mutex keeps a concurrent edit from committing between
	// saving content and writing backlinks (no nested locking — holder owns all).
	backlink_mutex.Lock()
	defer backlink_mutex.Unlock()

	var old_data string
	old_exist := tool.QueryRow_DB(
		db,
		"select data from data where title = ?",
		[]any{&old_data},
		doc_name,
	)

	length := tool.Get_edit_length_diff(old_data, data)
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		if old_exist {
			if _, err := tx.Exec(tool.DB_change("update data set data = ? where title = ?"), data, doc_name); err != nil {
				return err
			}
		} else if _, err := tx.Exec(tool.DB_change("insert into data (title, data) values (?, ?)"), doc_name, data); err != nil {
			return err
		}
		tool.Do_add_history(tx, doc_name, data, tool.Get_time(), config.IP, send, length, "revert", "r"+rev)
		return nil
	}); err != nil {
		panic(err)
	}
	tool.Search_index_sync(db, doc_name)

	markup.Get_render(db, doc_name, data, "backlink")

	return_data["response"] = "ok"

	return return_data
}
