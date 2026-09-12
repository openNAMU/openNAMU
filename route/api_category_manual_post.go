package route

import (
	"database/sql"
	"strings"

	"opennamu/route/tool"
)

func category_manual_name(value string) string {
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "분류:")
	value = strings.TrimPrefix(value, "category:")
	if value == "" {
		return ""
	}
	return "category:" + value
}

func Api_category_manual_post(config tool.Config, action string, category_name string, doc_name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	category_name = category_manual_name(category_name)
	doc_name = strings.TrimSpace(doc_name)
	if category_name == "" || doc_name == "" {
		return_data["response"] = "error"
		return_data["data"] = "empty data"
		return return_data
	}
	_, doc_exists := tool.Get_data_title(db, doc_name)
	_, category_exists := tool.Get_data_title(db, category_name)
	if !doc_exists || (action == "add" && !category_exists) {
		return_data["response"] = "not exist"
		return return_data
	}
	if !tool.Check_acl(db, doc_name, "", "document_edit", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	if action == "add" {
		value := ""
		if tool.QueryRow_DB(
			db,
			"select link from back where link = ? and title = ? and type = 'cat_manual' limit 1",
			[]any{&value},
			doc_name,
			category_name,
		) {
			return_data["response"] = "ok"
			return return_data
		}

		if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
			_, err := tx.Exec(
				tool.DB_change("insert into back (link, title, type, data) values (?, ?, 'cat_manual', '')"),
				doc_name,
				category_name,
			)
			return err
		}); err != nil {
			panic(err)
		}
	} else if action == "delete" {
		if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
			_, err := tx.Exec(
				tool.DB_change("delete from back where link = ? and title = ? and type = 'cat_manual'"),
				doc_name,
				category_name,
			)
			return err
		}); err != nil {
			panic(err)
		}
	} else {
		return_data["response"] = "error"
		return_data["data"] = "invalid data"
		return return_data
	}

	return_data["response"] = "ok"
	return return_data
}
