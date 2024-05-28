package route

import (
	"encoding/json"
	"opennamu/route/tool"
)

func Api_edit_move_all(call_arg []string) string {
	db_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &db_set)

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[1]), &other_set)

	db := tool.DB_connect(db_set)
	defer db.Close()

	if other_set["select"] == "include" {

	} else if other_set["select"] == "start" {

	} else {

	}
}
