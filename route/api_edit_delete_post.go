package route

import (
	"strconv"

	"opennamu/route/tool"
)

func Api_edit_delete_post(config tool.Config, doc_name string, send string, agree string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_acl(db, doc_name, "", "document_delete", config.IP) {
		return_data["response"] = "require auth"

		return return_data
	}

	raw_data := ""
	exist := tool.QueryRow_DB(
		db,
		"select data from data where title = ?",
		[]any{&raw_data},
		doc_name,
	)
	if !exist {
		return_data["response"] = "not exist"

		return return_data
	} else if !tool.Do_edit_slow_check(db, config, "edit") {
		return_data["response"] = "error"
		return_data["data"] = "slow edit limit"

		return return_data
	} else if !tool.Do_edit_send_require_check(db, config, send) {
		return_data["response"] = "error"
		return_data["data"] = "send require"

		return return_data
	} else if !tool.Do_edit_text_checkbox_check(db, config, agree) {
		return_data["response"] = "error"
		return_data["data"] = "checkbox check require"

		return return_data
	}

	tool.Do_add_history(
		db,
		doc_name,
		"",
		tool.Get_time(),
		config.IP,
		send,
		"-"+strconv.Itoa(len(raw_data)),
		"delete",
		"",
	)

	rows := tool.Query_DB(
		db,
		"select title, link from back where title = ? and not type = 'cat' and not type = 'no'",
		doc_name,
	)
	link_list := [][]string{}
	for rows.Next() {
		title := ""
		link := ""

		err := rows.Scan(&title, &link)
		if err != nil {
			panic(err)
		}

		link_list = append(link_list, []string{title, link})
	}
	rows.Close()

	for _, link_data := range link_list {
		tool.Exec_DB(
			db,
			"insert into back (title, link, type, data) values (?, ?, 'no', '')",
			link_data[0],
			link_data[1],
		)
	}

	tool.Exec_DB(
		db,
		"delete from back where link = ?",
		doc_name,
	)
	tool.Exec_DB(
		db,
		"delete from data where title = ?",
		doc_name,
	)

	return_data["response"] = "ok"

	return return_data
}
