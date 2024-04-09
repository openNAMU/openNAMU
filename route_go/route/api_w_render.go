package route

import (
	"encoding/json"
	"opennamu/route/tool"
)

func Api_w_render(call_arg []string) string {
	db_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &db_set)

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[1]), &other_set)

	db := tool.DB_connect(db_set)
	defer db.Close()

	data := tool.Get_render(db, db_set, other_set["doc_name"], other_set["data"], other_set["render_type"])

	json_data, _ := json.Marshal(data)
	return string(json_data)
}
