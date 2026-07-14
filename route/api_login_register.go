package route

import "opennamu/route/tool"

func Api_login_register(config tool.Config, id string, password string, password_check string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    return_data := make(map[string]any)

    if !tool.IP_or_user(config.IP) {
        return_data["response"] = "error"
        return_data["data"] = "login user"

        return return_data
    }

    ban_data := tool.Get_user_ban(db, config.IP, "register")
    if ban_data[0] == "true" {
        return_data["response"] = "error"
        return_data["data"] = "ban"
        return_data["ban_type"] = ban_data[1]

        return return_data
    }

    if password != password_check {
        return_data["response"] = "error"
        return_data["data"] = "password error"

        return return_data
    }

    if password == "" {
        return_data["response"] = "error"
        return_data["data"] = "empty password"

        return return_data
    }

    password_length_limit_str := ""
    tool.QueryRow_DB(
        db,
        `select data from other where name = 'password_min_length'`,
        []any{ &password_length_limit_str },
    )

    password_length_limit := tool.Str_to_int(password_length_limit_str)
    if len(password) < password_length_limit {
        return_data["response"] = "error"
        return_data["data"] = "password too short"

        return return_data
    }

    if !tool.Get_user_name_check(db, id) {
        return_data["response"] = "error"
        return_data["data"] = "user name error"

        return return_data
    }

    Api_add_user(config, id, password, "", "")

    return_data["response"] = "ok"

    return return_data
}