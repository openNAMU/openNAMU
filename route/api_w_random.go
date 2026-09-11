package route

import (
	"opennamu/route/tool"
)

func Api_w_random(config tool.Config) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	title := "Test"
	tool.QueryRow_DB(
		db,
		"select title from data where title not like 'user:%' and title not like 'category:%' and title not like 'file:%' order by random() limit 1",
		[]any{&title},
	)

	new_data := map[string]any{}
	new_data["response"] = "ok"
	new_data["data"] = title

	return new_data
}

func Api_w_random_category(config tool.Config, category_name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	title := ""
	tool.QueryRow_DB(
		db,
		"select distinct title from back where link = ? and (type = 'cat' or type = '') and title not like 'user:%' and title not like 'category:%' and title not like 'file:%' order by random() limit 1",
		[]any{&title},
		category_name,
	)

	new_data := map[string]any{}
	new_data["response"] = "ok"
	new_data["data"] = title

	return new_data
}
