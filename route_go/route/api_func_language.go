package route

import (
	"encoding/json"
	"fmt"
	"opennamu/route/tool"
)

func Api_func_language(call_arg []string) {
	db_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &db_set)

	other_set := make(map[string]interface{})
	json.Unmarshal([]byte(call_arg[1]), &other_set)

	db := tool.DB_connect(db_set)
	if db == nil {
		return
	}
	defer db.Close()

	data_list := map[string][]string{}
	data_list["data"] = []string{}

	temp_list := other_set["data"].([]interface{})

	for for_a := 0; for_a < len(temp_list); for_a++ {
		data_list["data"] = append(data_list["data"], tool.Get_language(db, db_set, temp_list[for_a].(string), false))
	}

	json_data, _ := json.Marshal(data_list)
	fmt.Print(string(json_data))
}
