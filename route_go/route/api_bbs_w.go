package route

import (
	"log"
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_bbs_w(call_arg []string) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	
	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &other_set)

	db := tool.DB_connect()
	defer db.Close()

	stmt, err := db.Prepare(tool.DB_change("select set_name, set_data from bbs_data where set_id = ? and set_code = ?"))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(other_set["bbs_num"], other_set["post_num"])
	if err != nil {
		log.Fatal(err)
	}

	data_list := map[string]string{}

	for rows.Next() {
		var set_name string
		var set_data string

		data_list[set_name] = set_data
	}

	return_data := make(map[string]interface{})
	return_data["language"] = map[string]string{}
	return_data["data"] = data_list

	json_data, _ := json.Marshal(return_data)
	return string(json_data)
}
