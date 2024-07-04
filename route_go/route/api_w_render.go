package route

import (
	"opennamu/route/tool"

	jsoniter "github.com/json-iterator/go"
)

func Api_w_render(call_arg []string) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &other_set)

	db := tool.DB_connect()
	defer db.Close()

	data := tool.Get_render(db, other_set["doc_name"], other_set["data"], other_set["render_type"])

	json_data, _ := json.Marshal(data)
	return string(json_data)
}
