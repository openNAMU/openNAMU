package route

import "opennamu/route/tool"

func Api_login_login(config tool.Config, id string, password string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    if !tool.Password_check(db, id, password) {
        return_data := make(map[string]any)
        return_data["response"] = "error"
        return_data["data"] = "password error"

        return return_data
    }

    config.Session.Set("id", id)

    return_data := make(map[string]any)

    if err := config.Session.Save(); err != nil {
        return_data["response"] = "error"
        return_data["data"] = "session save error"

        return return_data
    }

    return_data["response"] = "ok"

    return return_data
}