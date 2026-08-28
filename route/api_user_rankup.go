package route

import (
	"strings"

	"opennamu/route/tool"
)

func User_rankup_condition(data string) string {
	return tool.Rankup_condition_type(data)
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

		end_data := make(map[string]any)
		for _, condition_data := range tool.Get_setting(db, "rankup_condition", coverage) {
			if len(condition_data) < 2 || !tool.Auth_group_name_rankup(condition_data[1]) {
				continue
			}

			condition_list, ok := tool.Rankup_condition_list(condition_data[0])
			if !ok {
				continue
			}

			rank_data_map, ok := end_data[condition_data[1]].(map[string]any)
			if !ok {
				rank_data_map = make(map[string]any)
				end_data[condition_data[1]] = rank_data_map
			}

			for _, condition := range condition_list {
				split_data := strings.Fields(condition)
				if len(split_data) == 2 && User_rankup_condition(split_data[0]) == "int" {
					rank_data_map[split_data[0]] = tool.Str_to_int(split_data[1])
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
