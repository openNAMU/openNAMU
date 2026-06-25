package route

import (
	"opennamu/route/tool"
)

func Api_user_setting_editor_post(config tool.Config, data string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    ip := config.IP

    if !tool.IP_or_user(ip) {
        tool.Exec_DB(
            db,
            "insert into user_set (id, name, data) values (?, 'user_editor_top', ?)",
            ip,
            data,
        )
        
        return_data := make(map[string]any)
        return_data["response"] = "ok"
        return_data["language"] = map[string]string{
            "save": tool.Get_language(db, "save", false),
        }

        return return_data
    } else {
        return_data := make(map[string]any)
        return_data["response"] = "require auth"
        return_data["language"] = map[string]string{
            "authority_error": tool.Get_language(db, "authority_error", false),
        }

        return return_data
    }
}
