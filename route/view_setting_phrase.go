package route

import (
	"database/sql"
	"strings"

	"opennamu/route/tool"
)

type setting_phrase_field struct {
	name       string
	label      string
	input      bool
	class_name string
}

func setting_phrase_fields() []setting_phrase_field {
	return []setting_phrase_field{
		{name: "contract", label: "register_text"},
		{name: "no_login_warning", label: "non_login_alert"},
		{name: "copyright_checkbox_text", label: "copyright_checkbox_text"},
		{name: "check_key_text", label: "check_key_text"},
		{name: "email_title", label: "email_title"},
		{name: "email_text", label: "email_text"},
		{name: "email_insert_text", label: "email_insert_text"},
		{name: "password_search_text", label: "password_search_text"},
		{name: "reset_user_text", label: "reset_user_text"},
		{name: "error_401", label: "error_401"},
		{name: "error_404", label: "error_404"},
		{name: "approval_question", label: "approval_question"},
		{name: "edit_help", label: "edit_help"},
		{name: "upload_help", label: "upload_help"},
		{name: "upload_default", label: "upload_default"},
		{name: "license", label: "bottom_text"},
		{name: "topic_text", label: "topic_text"},
		{name: "phrase_user_page_admin", label: "phrase_user_page_admin"},
		{name: "phrase_user_page_owner", label: "phrase_user_page_owner"},
		{name: "phrase_old_page_warning", label: "phrase_old_page_warning"},
		{name: "bbs_help", label: "bbs_help"},
		{name: "bbs_comment_help", label: "bbs_comment_help"},
		{name: "outdated_doc_warning", label: "outdated_doc_warning"},
		{name: "outdated_doc_warning_date", label: "period", input: true},
		{name: "category_text", label: "category", input: true},
		{name: "redirect_text", label: "redirect", input: true},
		{name: "template_var_1", label: "template_var_1"},
		{name: "template_var_2", label: "template_var_2"},
		{name: "template_var_3", label: "template_var_3"},
		{name: "edit_bottom_text", label: "edit_bottom_text"},
		{name: "edit_only_bottom_text", label: "edit_only_bottom_text"},
		{name: "move_bottom_text", label: "move_bottom_text"},
		{name: "delete_bottom_text", label: "delete_bottom_text"},
		{name: "revert_bottom_text", label: "revert_bottom_text"},
	}
}

func View_setting_phrase(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	fields := setting_phrase_fields()
	values := make([]setting_field, 0, len(fields))
	for _, field := range fields {
		values = append(values, setting_field{name: field.name})
	}

	return view_setting_phrase_data(db, config, setting_load_fields(db, values))
}

func View_setting_phrase_post(config tool.Config, form map[string]string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	fields := setting_phrase_fields()
	for _, field := range fields {
		setting_save_value(db, field.name, "", setting_form_value(form, field.name, ""))
	}
	tool.Do_insert_auth_history(db, config.IP, "edit_set (phrase)")

	return tool.Get_redirect("/setting/phrase")
}

func view_setting_phrase_data(db *sql.DB, config tool.Config, values map[string]string) string {
	lang := func(name string) string {
		return tool.Get_language(db, name, true)
	}

	data := strings.Builder{}
	data.WriteString(`<form method="post">`)

	for _, field := range setting_phrase_fields() {
		data.WriteString(`<h2>` + lang(field.label) + `</h2>`)

		switch field.name {
		case "approval_question":
			data.WriteString(`<sup><a href="/setting/main">` + lang("approval_question_visible_only_when_approval_on") + `</a></sup>` + setting_hr())
		case "outdated_doc_warning":
			data.WriteString(`<span>` + lang("period") + ` (` + lang("day") + `) (` + lang("off") + ` : ` + lang("empty") + `)</span>` + setting_hr())
		case "redirect_text":
			data.WriteString(`<span>EX : {0} ➤ {1}</span>` + setting_hr())
		}

		if field.input {
			data.WriteString(setting_input(field.name, values[field.name], "text"))
		} else {
			class_name := field.class_name
			if class_name == "" {
				class_name = "opennamu_textarea_100"
			}
			data.WriteString(setting_textarea(field.name, values[field.name], class_name))
		}

		data.WriteString(setting_hr())
	}

	data.WriteString(`<button id="opennamu_save_button" type="submit">` + lang("save") + `</button></form>`)

	return setting_page(db, config, lang("text_setting"), data.String(), "setting")
}
