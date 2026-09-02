package route

import (
	"database/sql"
	"net/url"
	"time"

	"opennamu/route/tool"
)

func app_submit_action(db *sql.DB, config tool.Config, user_id string, approve bool) bool {
	if !approve {
		tool.Exec_DB(db, "delete from user_set where id = ? and name = 'application'", user_id)
		tool.Do_insert_auth_history(db, config.IP, "application_decline ("+user_id+")")
		return true
	}

	raw := ""
	if !tool.QueryRow_DB(db, "select data from user_set where id = ? and name = 'application'", []any{&raw}, user_id) {
		return false
	}
	application := map[string]string{}
	if err := json.Unmarshal([]byte(raw), &application); err != nil {
		return false
	}

	application_id := application["id"]
	if application_id == "" {
		application_id = user_id
	}
	result := map[string]any{}
	if application["pw_hash"] != "" {
		result = Api_add_user_hash(config, application_id, application["pw_hash"], application["email"], application["encode"])
	} else if application["pw"] != "" {
		result = Api_add_user(config, application_id, application["pw"], application["email"], application["encode"])
	} else {
		return false
	}
	if result["response"] != "ok" {
		return false
	}

	tool.Exec_DB(db, "insert into user_set (name, id, data) values ('approval_question', ?, ?)", application_id, application["question"])
	tool.Exec_DB(db, "insert into user_set (name, id, data) values ('approval_question_answer', ?, ?)", application_id, application["answer"])
	tool.Exec_DB(db, "delete from user_set where id = ? and name = 'application'", user_id)
	tool.Do_insert_auth_history(db, config.IP, "application_approve ("+user_id+")")
	return true
}

func Api_app_submit_expire(config tool.Config) {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "application_manage", config.IP) {
		return
	}
	days := tool.Str_to_int(tool.Get_setting_value(db, "application_expiration_date", "", "0"))
	action := tool.Get_setting_value(db, "application_expiration_action", "", "")
	if days <= 0 || !tool.Arr_in_str([]string{"approve", "decline"}, action) {
		return
	}

	user_list := []string{}
	rows := tool.Get_application_rows(db)
	for rows.Next() {
		user_id := ""
		raw := ""
		if rows.Scan(&user_id, &raw) != nil {
			continue
		}
		application := map[string]string{}
		if json.Unmarshal([]byte(raw), &application) != nil {
			continue
		}
		if application["id"] != "" {
			user_id = application["id"]
		}
		application_date, err := time.ParseInLocation("2006-01-02 15:04:05", application["date"], time.Local)
		if err != nil || time.Now().Before(application_date.AddDate(0, 0, days)) {
			continue
		}
		user_list = append(user_list, user_id)
	}
	rows.Close()

	for _, user_id := range user_list {
		app_submit_action(db, config, user_id, action == "approve")
	}
}

func Api_app_submit_post(config tool.Config, values url.Values) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "application_manage", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if values.Get("approve") != "" || values.Get("decline") != "" {
		if !tool.Check_permission(db, "application_view", config.IP) {
			return_data["response"] = "require auth"
			return return_data
		}
	}

	approved := values.Get("approve")
	declined := values.Get("decline")
	if approved != "" {
		if !app_submit_action(db, config, approved, true) {
			return_data["response"] = "error"
			return return_data
		}
	} else if declined != "" {
		app_submit_action(db, config, declined, false)
	}

	return_data["response"] = "ok"
	return return_data
}
