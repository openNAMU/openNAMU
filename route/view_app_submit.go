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
		data += tool.Get_language(db, "all_register_num", true) + " : " + strconv.Itoa(len(applications)) + `<hr class="main_hr"><table id="main_table_set"><tr id="main_table_top_tr"><td>` + tool.Get_language(db, "id", true) + `</td><td>` + tool.Get_language(db, "email", true) + `</td><td>` + tool.Get_language(db, "approve_or_decline", true) + `</td></tr>`
		for _, application := range applications {
			user_id := application["id"]
			data += `<tr><td>` + tool.HTML_escape(user_id) + `<br>` + tool.HTML_escape(application["question"]) + `<br>` + tool.HTML_escape(application["answer"]) + `</td><td>` + tool.HTML_escape(application["email"]) + `</td><td><form method="post"><button name="approve" value="` + tool.HTML_escape(user_id) + `">` + tool.Get_language(db, "approve", true) + `</button> <button name="decline" value="` + tool.HTML_escape(user_id) + `">` + tool.Get_language(db, "decline", true) + `</button></form></td></tr>`
		}
		data += `</table>`
	}

	return tool.Get_template(db, config, tool.Get_language(db, "application_list", true), data, []any{}, [][]any{{"manager", tool.Get_language(db, "return", true)}}, map[string]string{})
}
