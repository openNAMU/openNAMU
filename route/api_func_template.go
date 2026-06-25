package route

import (
	"opennamu/route/tool"
)

func Api_func_template(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    other_set := map[string]any{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    sub_menu := []any{}
    if v, ok := other_set["sub"]; ok {
        switch x := v.(type) {
        case []any:
            sub_menu = x
        default:
        }
    }

    menu := [][]any{}
    if v, ok := other_set["menu"]; ok {
        switch x := v.(type) {
        case [][]any:
            menu = x
        case []any:
            for _, row := range x {
                if r, ok := row.([]any); ok {
                    menu = append(menu, r)
                }
            }
        default:
        }
    }

    option := map[string]string{}
    if v, ok := other_set["option"]; ok {
        for k, v := range v.(map[string]any) {
            option[k] = v.(string)
        }
    }

    return tool.Get_template(
        db,
        config,
        other_set["name"].(string),
        other_set["data"].(string),
        sub_menu,
        menu,
        option,
    )
}
