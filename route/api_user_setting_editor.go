package route

import (
	"opennamu/route/tool"
)

func Api_user_setting_editor(config tool.Config) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    ip := config.IP
    if !tool.IP_or_user(ip) {
        rows := tool.Query_DB(
            db,
            "select data from user_set where id = ? and name = 'user_editor_top'",
            ip,
        )
        defer rows.Close()

        data_list := []string{}

        for rows.Next() {
            var data string

            err := rows.Scan(&data)
            if err != nil {
                panic(err)
            }

            data_list = append(data_list, data)
        }

        return_data := make(map[string]any)
        return_data["response"] = "ok"
        return_data["data"] = data_list

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
