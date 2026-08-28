package route

import (
	"net/url"

	"opennamu/route/tool"
)

func Api_app_submit_post(config tool.Config, values url.Values) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_acl(db, "", "", "application_auth", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	approved := values.Get("approve")
	declined := values.Get("decline")
	if approved != "" {
		raw := ""
		if !tool.QueryRow_DB(db, "select data from user_set where id = ? and name = 'application'", []any{&raw}, approved) {
			return_data["response"] = "error"
			return return_data
		}
		application := map[string]string{}
		if err := json.Unmarshal([]byte(raw), &application); err != nil {
			return_data["response"] = "error"
			return return_data
		}
		user_id := application["id"]
		result := map[string]any{}
		if application["pw_hash"] != "" {
			result = Api_add_user_hash(config, user_id, application["pw_hash"], application["email"], application["encode"])
		} else if application["pw"] != "" {
			result = Api_add_user(config, user_id, application["pw"], application["email"], application["encode"])
		} else {
			return_data["response"] = "error"
			return return_data
		}
		if result["response"] != "ok" {
			return_data["response"] = "error"
			return return_data
		}
		tool.Exec_DB(db, "insert into user_set (name, id, data) values ('approval_question', ?, ?)", user_id, application["question"])
		tool.Exec_DB(db, "insert into user_set (name, id, data) values ('approval_question_answer', ?, ?)", user_id, application["answer"])
		tool.Exec_DB(db, "delete from user_set where id = ? and name = 'application'", approved)
		tool.Do_insert_auth_history(db, config.IP, "application_approve ("+approved+")")
	} else if declined != "" {
		tool.Exec_DB(db, "delete from user_set where id = ? and name = 'application'", declined)
		tool.Do_insert_auth_history(db, config.IP, "application_decline ("+declined+")")
	}

	return_data["response"] = "ok"
	return return_data
}
