package route

import (
	"opennamu/route/tool"
)

func View_w_random(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_w_random(config)
	response, _ := api_data["response"].(string)
	if response != "ok" {
		return tool.Get_error_page(db, config, "not found")
	}

	title, _ := api_data["data"].(string)

	redirect := tool.Get_redirect("/w/" + tool.Url_parser(title))

	return redirect
}

func View_w_random_category(config tool.Config, category_name string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_w_random_category(config, category_name)
	title, _ := api_data["data"].(string)
	if title == "" {
		return tool.Get_redirect("/w/" + tool.Url_parser(category_name))
	}

	redirect := tool.Get_redirect("/w/" + tool.Url_parser(title))

	return redirect
}
