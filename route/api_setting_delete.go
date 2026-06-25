package route

import (
	"opennamu/route/tool"
)

func Api_setting_delete(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    auth_info := tool.Check_acl(db, "", "", "owner_auth", config.IP)

    setting_acl := Setting_list()
    return_data := make(map[string]any)

    if _, ok := setting_acl[other_set["set_name"]]; ok {
        if auth_info {
            tool.Exec_DB(
                db,
                "delete from other where name = ?",
                other_set["set_name"],
            )
        } else {
            return_data["response"] = "require auth"
        }
    } else {
        return_data["response"] = "not exist"
    }

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}
