package route

import (
	"opennamu/route/tool"
	"opennamu/route/tool/markup"
)

func Api_list_markup(config tool.Config) map[string]any {
    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    data := markup.List_markup()

    return_data := make(map[string]any)
    return_data["response"] = "ok"
    return_data["data"] = data

    return return_data
}
