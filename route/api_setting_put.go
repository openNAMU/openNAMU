package route

import (
	"opennamu/route/tool"
)

func Api_setting_put(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    auth_info := tool.Check_acl(db, "", "", "owner_auth", config.IP)

    setting_acl := Setting_list()
    return_data := make(map[string]any)

    if _, ok := setting_acl[other_set["set_name"]]; ok {
        if auth_info {
            if _, ok := other_set["coverage"]; !ok {
                tool.Exec_DB(
                    db,
                    "delete from other where name = ?",
                    other_set["set_name"],
                )
            }

            data_coverage := ""
            if val, ok := other_set["coverage"]; ok {
                data_coverage = val
            }

            tool.Exec_DB(
                db,
                "insert into other (name, data, coverage) values (?, ?, ?)",
                other_set["set_name"], other_set["data"], data_coverage,
            )
            
            return_data["response"] = "ok"
        } else {
            return_data["response"] = "require auth"
        }
    } else {
        return_data["response"] = "not exist"
    }

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}
