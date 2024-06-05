package route

import (
	"encoding/json"
	"opennamu/route/tool"
)

func Api_func_language(call_arg []string) string {
	other_set := make(map[string]interface{})
	json.Unmarshal([]byte(call_arg[0]), &other_set)

	db := tool.DB_connect()
	defer db.Close()

	data_list := map[string][]string{}
	data_list["data"] = []string{}

	temp_list := other_set["data"].([]interface{})

	for for_a := 0; for_a < len(temp_list); for_a++ {
		data_list["data"] = append(data_list["data"], tool.Get_language(db, temp_list[for_a].(string), false))
	}

	json_data, _ := json.Marshal(data_list)
	return string(json_data)
}
