package route

import (
	"opennamu/route/tool"
	"opennamu/route/tool/markup"
)

func Api_list_markup(config tool.Config) string {
    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    data := markup.List_markup()

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = data

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}
