package route

import (
	"opennamu/route/tool"
	"strings"
)

func Api_func_language(config tool.Config, data string, safe_string string, legacy string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    temp_list := strings.Split(data, " ")

    safe := false
    if safe_string != "" {
        safe = true
    }

    if legacy != "" {
        data_list := make(map[string]any)
        tmp_list := []string{}

        for for_a := 0; for_a < len(temp_list); for_a++ {
            tmp_list = append(tmp_list, tool.Get_language(db, temp_list[for_a], safe))
        }

        data_list["data"] = tmp_list

        return data_list
    } else {
        new_data := make(map[string]any)
        new_data["response"] = "ok"

        data_list := map[string]string{}

        for for_a := 0; for_a < len(temp_list); for_a++ {
            data_list[temp_list[for_a]] = tool.Get_language(db, temp_list[for_a], safe)
        }

        new_data["data"] = data_list

        return new_data
    }
}
