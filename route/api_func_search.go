package route

import (
	"strings"

	"opennamu/route/tool"
)

func Api_func_search_ui(config tool.Config, title_list []string, keyword string, search_type string) []map[string]string {
	data_list := []map[string]string{}
	if search_type != "data" {
		for _, title := range title_list {
			data_list = append(data_list, map[string]string{
				"title":      title,
				"title_html": search_highlight(title, keyword),
			})
		}
		return data_list
	}

	db := tool.DB_connect()
	defer tool.DB_close(db)

	data_map := map[string]string{}
	if len(title_list) > 0 {
		placeholder_list := []string{}
		values := []any{}
		for _, title := range title_list {
			placeholder_list = append(placeholder_list, "?")
			values = append(values, title)
		}

		rows := tool.Query_DB(
			db,
			"select title, coalesce(data, '') from data where title in ("+strings.Join(placeholder_list, ", ")+")",
			values...,
		)
		defer rows.Close()

		for rows.Next() {
			title := ""
			data := ""
			if rows.Scan(&title, &data) != nil {
				continue
			}
			data_map[title] = data
		}
	}

	for _, title := range title_list {
		temp_data := map[string]string{
			"title":      title,
			"title_html": search_highlight(title, keyword),
		}
		if tool.Check_acl(db, title, "", "render", config.IP) {
			temp_data["search_snippet_html"] = search_snippet(data_map[title], keyword)
		}
		data_list = append(data_list, temp_data)
	}
	return data_list
}

func Api_func_search(config tool.Config, keyword string, num_str string, search_type string) map[string]any {
	page := tool.Str_to_int(num_str)
	num := 0
	if page*50 > 0 {
		num = page*50 - 50
	}

	name := keyword
	query := ""

	if search_type == "title" {
		name = tool.Do_remove_spaces(name)
		query = "select title from data where replace(title, ' ', '') collate nocase like ? order by title limit ?, 50"
	} else {
		query = "select title from data where data collate nocase like ? order by title limit ?, 50"
	}

	if keyword != "" {
		if title_list, ok := tool.Search_index_search(name, search_type, num, 50); ok {
			return map[string]any{
				"response": "ok",
				"data":     title_list,
			}
		}
	}

	db := tool.DB_connect()
	defer tool.DB_close(db)

	title_list := []string{}
	rows := tool.Query_DB(
		db,
		query,
		"%"+name+"%", num,
	)
	defer rows.Close()

	for rows.Next() {
		var title string

		err := rows.Scan(&title)
		if err != nil {
			panic(err)
		}

		title_list = append(title_list, title)
	}

	return_data := make(map[string]any)
	return_data["response"] = "ok"
	return_data["data"] = title_list

	return return_data
}
