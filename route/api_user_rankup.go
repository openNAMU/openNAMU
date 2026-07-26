package route

import (
	"opennamu/route/tool"
	"strings"
)

func User_rankup_condition(data string) string {
    can_set := map[string]string{
        "edit": "int",
        "time": "int",
    }

    if val, ok := can_set[data]; ok {
        return val
    } else {
        return ""
    }
}

func Api_user_rankup(config tool.Config, rankup_name string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    ip := config.IP
    if !tool.IP_or_user(ip) {
        coverage := ""
        if rankup_name != "" {
            coverage = rankup_name
        }

        return_data_api := tool.Get_setting(db, "rankup_condition", coverage)

        end_data := make(map[string]any)
        if len(return_data_api) != 0 {
            for k := range return_data_api {
                rank_name := string(return_data_api[k][0])
                rank_data := string(return_data_api[k][1])

                split_data := strings.Split(rank_data, " ")
                if len(split_data) == 2 {
                    type_data := User_rankup_condition(split_data[0])
                    if type_data != "" {
                        rank_data_map, ok := end_data[rank_name].(map[string]any)
                        if !ok {
                            rank_data_map = make(map[string]any)
                            end_data[rank_name] = rank_data_map
                        }

                        if type_data == "int" {
                            split_int := tool.Str_to_int(split_data[1])
                            rank_data_map[split_data[0]] = split_int
                        } else {
                            rank_data_map[split_data[0]] = split_data[1]
                        }
                    }
                }
            }
        }

        end_data["response"] = "ok"
        return end_data
    } else {
        return_data := make(map[string]any)
        return_data["response"] = "require auth"
        return_data["language"] = map[string]string{
            "authority_error": tool.Get_language(db, "authority_error", false),
        }

        return return_data
    }
}
