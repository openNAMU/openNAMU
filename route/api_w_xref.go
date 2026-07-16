package route

import (
	"opennamu/route/tool"
)

func Api_w_xref(config tool.Config, num_str string, doc_name string, do_type string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_acl(db, doc_name, "", "render", config.IP) {
		return_data["response"] = "require auth"
		return_data["data"] = [][]string{}

		return return_data
	}

	page := tool.Str_to_int(num_str)
	num := 0
	if page*50 > 0 {
		num = page*50 - 50
	}

	link_case_insensitive := ""
	tool.QueryRow_DB(
		db,
		"select data from other where name = 'link_case_insensitive'",
		[]any{&link_case_insensitive},
	)

	if link_case_insensitive != "" {
		link_case_insensitive = " collate nocase"
	}

	query := ""
	if do_type == "1" {
		query = "select distinct link, type from back where title" + link_case_insensitive + " = ? and not type = 'no' and not type = 'nothing' order by type asc, link asc limit ?, 50"
	} else {
		query = "select distinct title, type from back where link" + link_case_insensitive + " = ? and not type = 'no' and not type = 'nothing' order by type asc, title asc limit ?, 50"
	}

	rows := tool.Query_DB(
		db,
		query,
		doc_name,
		num,
	)

	data_list := [][]string{}
	for rows.Next() {
		name := ""
		type_data := ""

		err := rows.Scan(&name, &type_data)
		if err != nil {
			panic(err)
		}

		data_list = append(data_list, []string{name, type_data})
	}
	rows.Close()

	for i, data := range data_list {
		include_doc_name := ""
		tool.QueryRow_DB(
			db,
			"select title from back where title = ? and type = 'include' limit 1",
			[]any{&include_doc_name},
			data[0],
		)
		data_list[i] = append(data_list[i], include_doc_name)
	}

	if do_type != "1" {
		link_count := ""
		tool.QueryRow_DB(
			db,
			"select set_data from data_set where doc_name = ? and set_name = 'link_count'",
			[]any{&link_count},
			doc_name,
		)
		return_data["link_count"] = link_count
	}

	return_data["response"] = "ok"
	return_data["data"] = data_list

	return return_data
}
