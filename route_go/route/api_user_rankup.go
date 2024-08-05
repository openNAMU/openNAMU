package route

import (
	"opennamu/route/tool"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
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

func Api_user_rankup(call_arg []string) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &other_set)

	db := tool.DB_connect()
	defer db.Close()

	ip := other_set["ip"]
	if !tool.IP_or_user(ip) {
		coverage := ""
		if val, ok := other_set["rankup_name"]; ok {
			coverage = val
		}

		return_data_api := tool.Get_setting(db, "rankup_condition", coverage)

		end_data := make(map[string]interface{})
		if len(return_data_api) != 0 {
			for k := range return_data_api {
				rank_name := string(return_data_api[k][0])
				rank_data := string(return_data_api[k][1])

				split_data := strings.Split(rank_data, " ")
				if len(split_data) == 2 {
					type_data := User_rankup_condition(split_data[0])
					if type_data != "" {
						if _, ok := end_data[rank_name]; !ok {
							end_data[rank_name] = make(map[string]interface{})
						}

						if type_data == "int" {
							split_int, err := strconv.Atoi(split_data[1])
							if err == nil {
								end_data[rank_name].(map[string]int)[split_data[0]] = split_int
							}
						} else {
							end_data[rank_name].(map[string]string)[split_data[0]] = split_data[1]
						}
					}
				}
			}
		}

		json_data_end, _ := json.Marshal(end_data)
		return string(json_data_end)
	} else {
		return_data := make(map[string]interface{})
		return_data["response"] = "require auth"
		return_data["language"] = map[string]string{
			"authority_error": tool.Get_language(db, "authority_error", false),
		}

		json_data, _ := json.Marshal(return_data)
		return string(json_data)
	}
}
