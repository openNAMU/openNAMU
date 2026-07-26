package route

import (
	"opennamu/route/tool"
)

func Api_bbs_w_set_put(config tool.Config, set_id string, set_name string, data string, coverage string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    auth_info := tool.Check_acl(db, "", "", "owner_auth", config.IP)

    setting_acl := BBS_w_set_list()
    return_data := make(map[string]any)

    if _, ok := setting_acl[set_name]; ok {
        if auth_info {
            if coverage == "" {
                tool.Exec_DB(
                    db,
                    "delete from bbs_set where set_name = ? and set_id = ?",
                    set_name, set_id,
                )
            }

            tool.Exec_DB(
                db,
                "insert into bbs_set (set_name, set_code, set_id, set_data) values (?, '', ?, ?)",
                set_name, set_id, data,
            )

            return_data["response"] = "ok"
        } else {
            return_data["response"] = "require auth"
        }
    } else {
        return_data["response"] = "not exist"
    }

    return return_data
}
