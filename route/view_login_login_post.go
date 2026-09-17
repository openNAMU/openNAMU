package route

import "opennamu/route/tool"

func View_login_login_post(config tool.Config, id string, password string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := Api_login_login(config, id, password)
	response, _ := return_data["response"].(string)
	if response == "error" {
		error_name, _ := return_data["data"].(string)
		return tool.Get_error_page(db, config, error_name)
	}

	return tool.Get_redirect("/user")
}
