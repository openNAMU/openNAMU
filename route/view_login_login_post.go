package route

import "opennamu/route/tool"

func View_login_login_post(config tool.Config, id string, password string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    return_data := Api_login_login(config, id, password)
    if return_data["response"].(string) == "error" {
        return tool.Get_error_page(db, config, return_data["data"].(string))
    }

    return tool.Get_redirect("/user")
}