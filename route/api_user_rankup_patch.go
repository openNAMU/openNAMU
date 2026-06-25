package route

import (
	"opennamu/route/tool"
)

func Api_user_rankup_patch(config tool.Config) string {
    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    return "{}"
}
