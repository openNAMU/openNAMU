package route

import (
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_func_language(call_arg []string) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	other_set := make(map[string]interface{})
	json.Unmarshal([]byte(call_arg[0]), &other_set)

	db := tool.DB_connect()
	defer db.Close()

	temp_list := other_set["data"].([]interface{})

	if other_set["legacy"] != "" {
		data_list := map[string][]string{}
		data_list["data"] = []string{}

		for for_a := 0; for_a < len(temp_list); for_a++ {
			data_list["data"] = append(data_list["data"], tool.Get_language(db, temp_list[for_a].(string), false))
		}

		json_data, _ := json.Marshal(data_list)
		return string(json_data)
	} else {
		new_data := make(map[string]interface{})
		new_data["response"] = "ok"

		data_list := map[string]string{}

		for for_a := 0; for_a < len(temp_list); for_a++ {
			data_list[temp_list[for_a].(string)] = tool.Get_language(db, temp_list[for_a].(string), false)
		}

		new_data["data"] = data_list

		json_data, _ := json.Marshal(new_data)
		return string(json_data)
	}
}
