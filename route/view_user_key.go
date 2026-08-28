package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_user_key(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	api_data := Api_user_key_post(config)
	response, _ := api_data["response"].(string)
	if response != "ok" {
		if response == "require auth" {
			return tool.Get_redirect("/user")
		}
		return tool.Get_error_page(db, config, "error")
	}
	return tool.Get_redirect("/change")
}
