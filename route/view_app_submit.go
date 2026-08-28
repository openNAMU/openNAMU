package route

import (
	"net/url"
	"strconv"

	"opennamu/route/tool"
)

func View_app_submit(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "application_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	if values != nil {
		approved := values.Get("approve")
		declined := values.Get("decline")
		if approved != "" {
			raw := ""
			if !tool.QueryRow_DB(db, "select data from user_set where id = ? and name = 'application'", []any{&raw}, approved) {
				return tool.Get_error_page(db, config, "error")
			}
			application := map[string]string{}
			if err := json.Unmarshal([]byte(raw), &application); err != nil {
				return tool.Get_error_page(db, config, "error")
			}
			user_id := application["id"]
			result := map[string]any{}
			if application["pw_hash"] != "" {
				result = Api_add_user_hash(config, user_id, application["pw_hash"], application["email"], application["encode"])
			} else if application["pw"] != "" {
				result = Api_add_user(config, user_id, application["pw"], application["email"], application["encode"])
			} else {
				return tool.Get_error_page(db, config, "error")
			}
			if result["response"] != "ok" {
				return tool.Get_error_page(db, config, "error")
			}
			tool.Exec_DB(db, "insert into user_set (name, id, data) values ('approval_question', ?, ?)", user_id, application["question"])
			tool.Exec_DB(db, "insert into user_set (name, id, data) values ('approval_question_answer', ?, ?)", user_id, application["answer"])
			tool.Exec_DB(db, "delete from user_set where id = ? and name = 'application'", approved)
			tool.Do_insert_auth_history(db, config.IP, "application_approve ("+approved+")")
		} else if declined != "" {
			tool.Exec_DB(db, "delete from user_set where id = ? and name = 'application'", declined)
			tool.Do_insert_auth_history(db, config.IP, "application_decline ("+declined+")")
		}
		return tool.Get_redirect("/app_submit")
	}

	data := ""
	requirement := ""
	tool.QueryRow_DB(db, "select data from other where name = 'requires_approval'", []any{&requirement})
	if requirement != "on" {
		data += tool.Get_language(db, "approval_requirement_disabled", true) + `<hr class="main_hr">`
	}

	rows := tool.Query_DB(db, "select id, data from user_set where name = 'application'")
	applications := []map[string]string{}
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
		if application["id"] == "" {
			application["id"] = user_id
		}
		applications = append(applications, application)
	}
	rows.Close()

	if len(applications) == 0 {
		data += tool.Get_language(db, "no_applications_now", true)
	} else {
		data += tool.Get_language(db, "all_register_num", true) + " : " + strconv.Itoa(len(applications)) + `<hr class="main_hr"><table id="main_table_set"><tr id="main_table_top_tr"><td>` + tool.Get_language(db, "id", true) + `</td><td>` + tool.Get_language(db, "email", true) + `</td><td>` + tool.Get_language(db, "approve_or_decline", true) + `</td></tr>`
		for _, application := range applications {
			user_id := application["id"]
			data += `<tr><td>` + tool.HTML_escape(user_id) + `<br>` + tool.HTML_escape(application["question"]) + `<br>` + tool.HTML_escape(application["answer"]) + `</td><td>` + tool.HTML_escape(application["email"]) + `</td><td><form method="post"><button name="approve" value="` + tool.HTML_escape(user_id) + `">` + tool.Get_language(db, "approve", true) + `</button> <button name="decline" value="` + tool.HTML_escape(user_id) + `">` + tool.Get_language(db, "decline", true) + `</button></form></td></tr>`
		}
		data += `</table>`
	}

	return tool.Get_template(db, config, tool.Get_language(db, "application_list", true), data, []any{}, [][]any{{"manager", tool.Get_language(db, "return", true)}}, map[string]string{})
}
