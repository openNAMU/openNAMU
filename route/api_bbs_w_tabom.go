package route

import (
	"opennamu/route/tool"
)

func Api_bbs_w_tabom(config tool.Config, set_id string, set_code string) map[string]string {
    db := tool.DB_connect()
    defer tool.DB_close(db)
    
    return_data := make(map[string]string)

    if !tool.Check_acl(db, "", "", "bbs_comment", config.IP) {
        return_data["response"] = "require auth"
        return_data["data"] = "0"
    } else {
        tabom_count := "0"
        tool.QueryRow_DB(
            db,
            "select set_data from bbs_data where set_name = 'tabom_count' and set_id = ? and set_code = ?",
            []any{ &tabom_count },
            set_id,
            set_code,
        )
    
        return_data["response"] = "ok"
        return_data["data"] = tabom_count
    }

    return return_data
}