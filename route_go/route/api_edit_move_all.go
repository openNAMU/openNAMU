package route

import (
	"encoding/json"
	"opennamu/route/tool"
)

func Api_edit_move_all(call_arg []string) string {
	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &other_set)

	db := tool.DB_connect()
	defer db.Close()

	if other_set["select"] == "include" {

	} else if other_set["select"] == "start" {

	} else {

	}

	return ""
}
