package route

import "opennamu/route/tool"

func Api_w_down(config tool.Config, doc_name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, doc_name, "", "render", config.IP) {
		return map[string]any{"response": "require auth", "data": []string{}}
	}

	rows := tool.Query_DB(
		db,
		"select title from data where title like ?",
		doc_name,
	)
	defer rows.Close()

	title_list := []string{}

	for rows.Next() {
		var title string

		err := rows.Scan(&title)
		if err != nil {
			panic(err)
		}

		if !tool.Check_acl(db, title, "", "render", config.IP) {
			continue
		}

		title_list = append(title_list, title)
	}

	return_data := make(map[string]any)
	return_data["response"] = "ok"
	return_data["data"] = title_list

	return return_data
}
