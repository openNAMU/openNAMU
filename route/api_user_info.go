package route

import (
	"opennamu/route/tool"
)

func Api_user_info(config tool.Config, ip string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    data_result := map[string]any{}

    ip_render := tool.IP_parser(db, ip, config.IP)
    auth_name := tool.Get_user_auth(db, ip)
    level_data := tool.Get_level(db, ip)
    ban_check := tool.Get_user_ban(db, ip, "")
    user_document := tool.Get_user_document(db, ip)

    data_result["render"] = ip_render
    
    data_result["auth"] = auth_name
    data_result["auth_date"] = tool.Get_auth_date(db, ip)

    data_result["level"] = level_data[0]
    data_result["exp"] = level_data[1]
    data_result["max_exp"] = level_data[2]

    data_result["ban"] = func(ban_check []string) any {
        if ban_check[0] == "" {
            return "0"
        } else {
            return ban_check
        }
    }(ban_check)

    data_result["document"] = func(user_document bool) string {
        if user_document {
            return "1"
        } else {
            return "0"
        }
    }(user_document)

    data_result["user_title"] = tool.Get_user_title(db, ip)

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = data_result

    return return_data
}