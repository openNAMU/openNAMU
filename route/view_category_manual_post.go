package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_category_manual_post(config tool.Config, action string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	category_name := values.Get("category")
	doc_name := values.Get("document")
	return_name := values.Get("return")
	if return_name == "" {
		return_name = doc_name
	}

	api_data := Api_category_manual_post(config, action, category_name, doc_name)
	response, _ := api_data["response"].(string)
	switch response {
	case "ok":
		return tool.Get_redirect("/w/" + tool.Url_parser(return_name))
	case "require auth":
		return tool.Get_error_page(db, config, "auth")
	case "not exist":
		return tool.Get_error_page(db, config, "not found")
	}

	error_name, _ := api_data["data"].(string)
	if error_name == "" {
		error_name = "error"
	}
	return tool.Get_error_page(db, config, error_name)
}
