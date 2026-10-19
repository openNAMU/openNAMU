package route

import (
	"opennamu/route/tool"
)

func Api_w_random(config tool.Config) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	title := ""
	exist := tool.QueryRow_DB(
		db,
		"select title from data where title not like 'user:%' and title not like 'category:%' and title not like 'file:%' and not exists (select 1 from back where back.type = 'redirect' and back.link = data.title) order by random() limit 1",
		[]any{&title},
	)

	new_data := map[string]any{}
	if exist {
		new_data["response"] = "ok"
		new_data["data"] = title
	} else {
		new_data["response"] = "not exist"
		new_data["data"] = ""
	}

	return new_data
}

func Api_w_random_category(config tool.Config, category_name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	title := ""
	tool.QueryRow_DB(
		db,
		"select distinct title from back where link = ? and (type = 'cat' or type = 'cat_manual' or type = '') and title not like 'user:%' and title not like 'category:%' and title not like 'file:%' order by random() limit 1",
		[]any{&title},
		category_name,
	)

	new_data := map[string]any{}
	new_data["response"] = "ok"
	new_data["data"] = title

	return new_data
}
