package route

import (
	"encoding/json"
	"fmt"

	"opennamu/route/tool"
)

func Api_func_sha224(call_arg []string) {
	data := call_arg[0]

	hash_str := tool.Sha224(data)

	new_data := map[string]string{}
	new_data["data"] = hash_str

	json_data, _ := json.Marshal(new_data)
	fmt.Print(string(json_data))
}
