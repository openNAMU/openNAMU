package route

import (
	stdjson "encoding/json"

	"opennamu/route/tool"
)

func Api_register_submit_post(config tool.Config, id string, pw string, email string, question string, answer string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if id == "" || pw == "" || question == "" {
		return_data["response"] = "error"
		return return_data
	}
	encode := tool.Get_main_encode(db)
	application, err := stdjson.Marshal(map[string]string{
		"id":       id,
		"pw_hash":  tool.Password_encode(db, pw, encode),
		"email":    email,
		"encode":   encode,
		"question": question,
		"answer":   answer,
	})
	if err != nil {
		return_data["response"] = "error"
		return return_data
	}
	tool.Exec_DB(db, "delete from user_set where id = ? and name = 'application'", id)
	tool.Exec_DB(db, "insert into user_set (id, name, data) values (?, 'application', ?)", id, string(application))

	return_data["response"] = "ok"
	return return_data
}
