package route

import "opennamu/route/tool"

func View_login_register_post(config tool.Config, id string, password string, password_check string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := Api_login_register(config, id, password, password_check)
	if return_data["response"].(string) == "error" {
		error_name := return_data["data"].(string)

		switch error_name {
		case "login user":
			return tool.Get_redirect("/user")
		case "password error":
			error_name = "password different"
		}

		return tool.Get_error_page(db, config, error_name)
	}

	return tool.Get_redirect("/login")
}
