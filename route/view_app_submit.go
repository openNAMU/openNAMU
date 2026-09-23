package route

import (
	"net/url"
	"strconv"

	"opennamu/route/tool"
)

func View_app_submit(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if values == nil && !tool.Check_permission(db, "application_manage", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if values == nil {
		Api_app_submit_expire(config)
	}

	if values != nil {
		result := Api_app_submit_post(config, values)
		if result["response"] == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if result["response"] != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
		return tool.Get_redirect("/app_submit")
	}

	data := ""
	requirement := tool.Get_setting_value(db, "requires_approval", "", "")
	if requirement != "on" {
		data += tool.Get_language(db, "approval_requirement_disabled", true) + `<hr class="main_hr">`
	}

	rows := tool.Get_application_rows(db)
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
		data += tool.Get_language(db, "all_register_num", true) + " : " + strconv.Itoa(len(applications))
		if !tool.Check_permission(db, "application_view", config.IP) {
			data += `<hr class="main_hr">` + tool.Get_language(db, "application_detail_hidden", true)
		} else {
			data += `<hr class="main_hr">`
			for _, application := range applications {
				user_id := application["id"]
				left := `<strong>` + tool.Get_language(db, "id", true) + `:</strong> ` + tool.HTML_escape(user_id)
				right := `<strong>` + tool.Get_language(db, "application_time", true) + `:</strong> ` + tool.HTML_escape(application["date"])
				bottom := `<div><strong>` + tool.Get_language(db, "email", true) + `:</strong> ` + tool.HTML_escape(application["email"]) + `</div>`
				bottom += `<div><strong>` + tool.Get_language(db, "approval_question", true) + `:</strong> ` + tool.HTML_escape(application["question"]) + `</div>`
				bottom += `<div><strong>` + tool.Get_language(db, "answer", true) + `:</strong> ` + tool.HTML_escape(application["answer"]) + `</div>`
				bottom += `<form method="post"><button name="approve" value="` + tool.HTML_escape(user_id) + `">` + tool.Get_language(db, "approve", true) + `</button> <button name="decline" value="` + tool.HTML_escape(user_id) + `">` + tool.Get_language(db, "decline", true) + `</button></form>`
				data += tool.Get_list_ui(left, right, bottom, "")
			}
		}
	}

	return tool.Get_template(db, config, tool.Get_language(db, "application_list", true), data, []any{}, [][]any{{"manager", tool.Get_language(db, "return", true)}}, map[string]string{})
}
