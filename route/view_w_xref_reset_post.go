package route

import (
	"opennamu/route/tool"
)

func View_w_xref_reset_post(config tool.Config, doc_name string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_w_raw(config, doc_name, "", "")
	response, _ := api_data["response"].(string)
	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	} else if response != "ok" {
		return tool.Get_redirect("/w/" + tool.Url_parser(doc_name))
	}

	render_data := Api_w_render(config, doc_name, api_data["data"].(string), "backlink", "")
	render_response, _ := render_data["response"].(string)
	if render_response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}
	if render_response != "ok" {
		return tool.Get_redirect("/w/" + tool.Url_parser(doc_name))
	}

	return tool.Get_redirect("/xref/" + tool.Url_parser(doc_name))
}
