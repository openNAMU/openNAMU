package route

import (
	"encoding/json"
	"opennamu/route/tool"
)

func Api_bbs_set_update(call_arg []string) string {
	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &other_set)

	db := tool.DB_connect()
	defer db.Close()

	return ""
}
