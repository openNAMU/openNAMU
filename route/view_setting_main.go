package route

import (
	"database/sql"
	"strings"

	"opennamu/route/tool"
)

func View_setting_main(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "setting_main", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	return View_setting_main_data(db, config, Setting_load_fields(db, Setting_main_fields()))
}

type setting_field struct {
	name          string
	default_value string
}

func Setting_value(db *sql.DB, name string, coverage string, default_value string) string {
	return tool.Get_setting_value_exact(db, name, coverage, default_value)
}

func Setting_load_fields(db *sql.DB, fields []setting_field) map[string]string {
	data := make(map[string]string, len(fields))

	for _, field := range fields {
		data[field.name] = Setting_value(db, field.name, "", field.default_value)
	}

	return data
}

func Setting_form_value(form map[string]string, name string, default_value string) string {
	value, exists := form[name]
	if !exists {
		return default_value
	}

	return value
}

func Setting_checked(value string) string {
	if value != "" {
		return `checked="checked"`
	}

	return ""
}

func Setting_options(current string, values []string, labels map[string]string) string {
	data := strings.Builder{}

	for _, value := range values {
		selected := ""
		if value == current {
			selected = ` selected="selected"`
		}

		label := value
		if labels != nil {
			if label_data, ok := labels[value]; ok {
				label = label_data
			}
		}

		data.WriteString(`<option value="`)
		data.WriteString(tool.HTML_escape(value))
		data.WriteString(`"`)
		data.WriteString(selected)
		data.WriteString(`>`)
		data.WriteString(tool.HTML_escape(label))
		data.WriteString(`</option>`)
	}

	return data.String()
}

func Setting_input(name string, value string, input_type string) string {
	if input_type == "" {
		input_type = "text"
	}

	return `<input type="` + tool.HTML_escape(input_type) + `" name="` + tool.HTML_escape(name) + `" value="` + tool.HTML_escape(value) + `">`
}

func Setting_textarea(name string, value string, class_name string) string {
	if class_name == "" {
		class_name = "opennamu_textarea_100"
	}

	return `<textarea class="` + tool.HTML_escape(class_name) + `" name="` + tool.HTML_escape(name) + `">` + tool.HTML_escape(value) + `</textarea>`
}

func Main_hr() string {
	return `<hr class="main_hr">`
}

func Setting_page(db *sql.DB, config tool.Config, title string, data string, return_path string) string {
	menu := [][]any{}
	if return_path != "" {
		menu = append(menu, []any{return_path, tool.Get_language(db, "return", true)})
	}

	return tool.Get_template(
		db,
		config,
		title,
		data,
		[]any{},
		menu,
		map[string]string{},
	)
}

func Setting_main_fields() []setting_field {
	return []setting_field{
		{name: "name", default_value: "Wiki"},
		{name: "frontpage", default_value: "FrontPage"},
		{name: "frontpage_type", default_value: "document"},
		{name: "upload", default_value: "2"},
		{name: "skin", default_value: ""},
		{name: "reg", default_value: ""},
		{name: "ip_view", default_value: ""},
		{name: "back_up", default_value: ""},
		{name: "port", default_value: "3000"},
		{name: "key", default_value: tool.Get_random_key(128)},
		{name: "update", default_value: "stable"},
		{name: "auto_update", default_value: ""},
		{name: "encode", default_value: "sha3"},
		{name: "host", default_value: "0.0.0.0"},
		{name: "slow_edit", default_value: ""},
		{name: "requires_approval", default_value: ""},
		{name: "application_expiration_date", default_value: ""},
		{name: "application_expiration_action", default_value: ""},
		{name: "backup_where", default_value: ""},
		{name: "domain", default_value: ""},
		{name: "ua_get", default_value: ""},
		{name: "enable_comment", default_value: ""},
		{name: "bbs_excellent_min", default_value: ""},
		{name: "edit_bottom_compulsion", default_value: ""},
		{name: "http_select", default_value: "http"},
		{name: "title_max_length", default_value: ""},
		{name: "title_topic_max_length", default_value: ""},
		{name: "password_min_length", default_value: ""},
		{name: "wiki_access_password_need", default_value: ""},
		{name: "wiki_access_password", default_value: ""},
		{name: "history_recording_off", default_value: ""},
		{name: "namumark_compatible", default_value: ""},
		{name: "user_name_view", default_value: ""},
		{name: "user_document_view_acl_all", default_value: ""},
		{name: "invite_required", default_value: ""},
		{name: "link_case_insensitive", default_value: ""},
		{name: "move_with_redirect", default_value: ""},
		{name: "slow_thread", default_value: ""},
		{name: "edit_timeout", default_value: "5"},
		{name: "document_content_max_length", default_value: ""},
		{name: "bbs_content_max_length", default_value: ""},
		{name: "backup_count", default_value: ""},
		{name: "ua_expiration_date", default_value: ""},
		{name: "auth_history_expiration_date", default_value: ""},
		{name: "auth_history_off", default_value: ""},
		{name: "user_name_level", default_value: ""},
		{name: "load_ip_select", default_value: ""},
		{name: "not_use_view_count", default_value: ""},
	}
}

func View_setting_main_data(db *sql.DB, config tool.Config, values map[string]string) string {
	lang := func(name string) string {
		return tool.Get_language(db, name, true)
	}

	encode_values := []string{}
	init_values := tool.Get_init_set_list("encode")
	if encode_set, ok := init_values["encode"]; ok {
		if values, ok := encode_set["list"].([]string); ok {
			encode_values = append(encode_values, values...)
		}
	}
	encode_values = append(encode_values, "sha256")

	skin_value := values["skin"]
	if skin_value == "" {
		skin_value = "ringo"
	}
	skin_values := tool.Get_skin_list(skin_value, true)
	if !tool.Arr_in_str(skin_values, skin_value) {
		skin_values = append([]string{skin_value}, skin_values...)
	}

	ip_values := []string{"default", "HTTP_X_REAL_IP", "HTTP_CF_CONNECTING_IP", "REMOTE_ADDR"}
	ip_labels := map[string]string{"default": lang("default")}

	sqlite_only := ""
	if tool.Get_DB_set()["db_type"] != "sqlite" {
		sqlite_only = ` style="display:none;"`
	}

	data := strings.Builder{}
	data.WriteString(`<form method="post">`)

	data.WriteString(`<h2>` + lang("basic_set") + `</h2>`)
	data.WriteString(`<span>` + lang("wiki_name") + `</span>` + Main_hr())
	data.WriteString(Setting_input("name", values["name"], "text") + Main_hr())
	data.WriteString(`<span><a href="/setting/main/logo">(` + lang("wiki_logo") + `)</a></span>` + Main_hr())
	data.WriteString(`<span>` + lang("main_page") + `</span>` + Main_hr())
	data.WriteString(`<select name="frontpage_type">` + Setting_options(values["frontpage_type"], []string{"document", "bbs"}, map[string]string{"document": lang("document"), "bbs": lang("bbs_main")}) + `</select>` + Main_hr())
	data.WriteString(Setting_input("frontpage", values["frontpage"], "text") + Main_hr())
	data.WriteString(`<span>` + lang("tls_method") + `</span>` + Main_hr())
	data.WriteString(`<select name="http_select">` + Setting_options(values["http_select"], []string{"http", "https"}, nil) + `</select>` + Main_hr())
	data.WriteString(`<span>` + lang("domain") + `</span> (EX : 2du.pythonanywhere.com) (` + lang("off") + ` : ` + lang("empty") + `)` + Main_hr())
	data.WriteString(Setting_input("domain", values["domain"], "text") + Main_hr())
	data.WriteString(`<span>` + lang("wiki_host") + `</span>` + Main_hr())
	data.WriteString(Setting_input("host", values["host"], "text") + Main_hr())
	data.WriteString(`<span>` + lang("wiki_port") + `</span>` + Main_hr())
	data.WriteString(Setting_input("port", values["port"], "text") + Main_hr())
	data.WriteString(`<span>` + lang("wiki_secret_key") + `</span>` + Main_hr())
	data.WriteString(Setting_input("key", values["key"], "password") + Main_hr())
	data.WriteString(`<label><input type="checkbox" name="wiki_access_password_need" ` + Setting_checked(values["wiki_access_password_need"]) + `> ` + lang("set_wiki_access_password_need") + ` (` + lang("restart_required") + `)</label>` + Main_hr())
	data.WriteString(`<span>` + lang("set_wiki_access_password") + `</span> (` + lang("restart_required") + `)` + Main_hr())
	data.WriteString(Setting_input("wiki_access_password", values["wiki_access_password"], "password") + Main_hr())
	data.WriteString(`<span>` + lang("wiki_load_ip_select") + `</span> (` + lang("restart_required") + `)` + Main_hr())
	data.WriteString(`<select name="load_ip_select">` + Setting_options(values["load_ip_select"], ip_values, ip_labels) + `</select>` + Main_hr())
	data.WriteString(`<h3>` + lang("authority_use_list") + `</h3>`)
	data.WriteString(`<label><input type="checkbox" name="auth_history_off" ` + Setting_checked(values["auth_history_off"]) + `> ` + lang("authority_use_list_off") + `</label>` + Main_hr())
	data.WriteString(`<span>` + lang("authority_use_list_expiration_date") + `</span> (` + lang("day") + `) (` + lang("off") + ` : ` + lang("empty") + `)` + Main_hr())
	data.WriteString(Setting_input("auth_history_expiration_date", values["auth_history_expiration_date"], "text") + Main_hr())
	data.WriteString(`<h3>` + lang("communication_set") + `</h3>`)
	data.WriteString(`<label><input type="checkbox" name="enable_comment" ` + Setting_checked(values["enable_comment"]) + `> ` + lang("enable_comment_function") + `</label>` + Main_hr())
	data.WriteString(`<label><input type="checkbox" name="user_name_level" ` + Setting_checked(values["user_name_level"]) + `> ` + lang("display_level_in_user_name") + `</label>` + Main_hr())
	data.WriteString(`<label><input type="checkbox" name="not_use_view_count" ` + Setting_checked(values["not_use_view_count"]) + `> ` + lang("not_use_view_count") + `</label>` + Main_hr())
	data.WriteString(`<span>` + lang("bbs_excellent_min") + ` (` + lang("empty") + ` : 5)</span>` + Main_hr())
	data.WriteString(Setting_input("bbs_excellent_min", values["bbs_excellent_min"], "number") + Main_hr())

	data.WriteString(`<h2>` + lang("design_set") + `</h2>`)
	data.WriteString(`<span>` + lang("wiki_skin") + `</span>` + Main_hr())
	data.WriteString(`<select name="skin">` + Setting_options(skin_value, skin_values, nil) + `</select>` + Main_hr())

	data.WriteString(`<h2>` + lang("render_set") + `</h2>`)
	data.WriteString(`<label><input type="checkbox" name="namumark_compatible" ` + Setting_checked(values["namumark_compatible"]) + `> ` + lang("namumark_fully_compatible_mode") + `</label>` + Main_hr())
	data.WriteString(`<label><input type="checkbox" name="link_case_insensitive" ` + Setting_checked(values["link_case_insensitive"]) + `> ` + lang("link_case_insensitive") + `</label>` + Main_hr())

	data.WriteString(`<h2>` + lang("login_set") + `</h2>`)
	data.WriteString(`<label><input type="checkbox" name="reg" ` + Setting_checked(values["reg"]) + `> ` + lang("no_register") + `</label>` + Main_hr())
	data.WriteString(`<label><input type="checkbox" name="ip_view" ` + Setting_checked(values["ip_view"]) + `> ` + lang("hide_ip") + `</label>` + Main_hr())
	data.WriteString(`<label><input type="checkbox" name="user_name_view" ` + Setting_checked(values["user_name_view"]) + `> ` + lang("hide_user_name") + `</label>` + Main_hr())
	data.WriteString(`<span>` + lang("user_document_view_acl_all") + ` (` + lang("off") + ` : ` + lang("empty") + `)</span>` + Main_hr())
	data.WriteString(tool.Build_select("user_document_view_acl_all", Acl_value_list(db, values["user_document_view_acl_all"]), values["user_document_view_acl_all"], lang("normal")) + Main_hr())
	data.WriteString(`<label><input type="checkbox" name="invite_required" ` + Setting_checked(values["invite_required"]) + `> ` + lang("invite_required") + `</label>` + Main_hr())
	data.WriteString(`<label><input type="checkbox" name="requires_approval" ` + Setting_checked(values["requires_approval"]) + `> ` + lang("requires_approval") + `</label>` + Main_hr())
	data.WriteString(`<span>` + lang("application_expiration_date") + ` (` + lang("day") + `) (` + lang("off") + ` : ` + lang("empty") + `)</span>` + Main_hr())
	data.WriteString(`<sup>` + lang("application_expiration_help") + `</sup>` + Main_hr())
	data.WriteString(Setting_input("application_expiration_date", values["application_expiration_date"], "text") + Main_hr())
	data.WriteString(`<span>` + lang("application_expiration_action") + `</span>` + Main_hr())
	data.WriteString(`<select name="application_expiration_action">` + Setting_options(values["application_expiration_action"], []string{"", "approve", "decline"}, map[string]string{"": lang("off"), "approve": lang("approve"), "decline": lang("decline")}) + `</select>` + Main_hr())
	data.WriteString(`<span>` + lang("password_min_length") + `</span> (` + lang("off") + ` : ` + lang("empty") + `)` + Main_hr())
	data.WriteString(Setting_input("password_min_length", values["password_min_length"], "text") + Main_hr())
	data.WriteString(`<span>` + lang("encryption_method") + `</span>` + Main_hr())
	data.WriteString(`<select name="encode">` + Setting_options(values["encode"], encode_values, nil) + `</select>` + Main_hr())

	data.WriteString(`<h3>` + lang("ua") + `</h3>`)
	data.WriteString(`<label><input type="checkbox" name="ua_get" ` + Setting_checked(values["ua_get"]) + `> ` + lang("ua_get_off") + `</label>` + Main_hr())
	data.WriteString(`<span>` + lang("ua_expiration_date") + `</span> (` + lang("day") + `) (` + lang("off") + ` : ` + lang("empty") + `)` + Main_hr())
	data.WriteString(Setting_input("ua_expiration_date", values["ua_expiration_date"], "text") + Main_hr())

	data.WriteString(`<h2>` + lang("server_set") + `</h2>`)
	data.WriteString(`<span>` + lang("update_branch") + `</span>` + Main_hr())
	data.WriteString(`<select name="update">` + Setting_options(values["update"], []string{"stable", "beta"}, nil) + `</select>` + Main_hr())
	data.WriteString(`<label><input type="checkbox" name="auto_update" ` + Setting_checked(values["auto_update"]) + `> ` + lang("auto_update") + `</label>` + Main_hr())
	data.WriteString(`<span` + sqlite_only + `>`)
	data.WriteString(`<h3>` + lang("backup") + ` (` + lang("sqlite_only") + `)</h3>`)
	data.WriteString(`<span>` + lang("backup_warning") + ` (EX : data_YYYYMMDDHHMMSS.db)</span>` + Main_hr())
	data.WriteString(`<span>` + lang("backup_interval") + ` (` + lang("hour") + `) (` + lang("off") + ` : ` + lang("empty") + `)</span>` + Main_hr())
	data.WriteString(Setting_input("back_up", values["back_up"], "text") + Main_hr())
	data.WriteString(`<span>` + lang("backup_where") + ` (` + lang("default") + ` : ` + lang("empty") + `) (` + lang("example") + ` : ./data/backup.db)</span>` + Main_hr())
	data.WriteString(Setting_input("backup_where", values["backup_where"], "text") + Main_hr())
	data.WriteString(`<span>` + lang("backup_count") + ` (` + lang("default") + ` : ` + lang("empty") + `)</span>` + Main_hr())
	data.WriteString(Setting_input("backup_count", values["backup_count"], "text") + Main_hr())
	data.WriteString(`</span>`)

	data.WriteString(`<h2>` + lang("edit_set") + `</h2>`)
	data.WriteString(`<span>` + lang("slow_edit") + ` (` + lang("second") + `) (` + lang("off") + ` : ` + lang("empty") + `)</span>` + Main_hr())
	data.WriteString(Setting_input("slow_edit", values["slow_edit"], "text") + Main_hr())
	data.WriteString(`<label><input type="checkbox" name="edit_bottom_compulsion" ` + Setting_checked(values["edit_bottom_compulsion"]) + `> ` + lang("edit_bottom_compulsion") + `</label>` + Main_hr())
	data.WriteString(`<span>` + lang("title_max_length") + ` (` + lang("off") + ` : ` + lang("empty") + `)</span>` + Main_hr())
	data.WriteString(Setting_input("title_max_length", values["title_max_length"], "text") + Main_hr())
	data.WriteString(`<span>` + lang("title_topic_max_length") + ` (` + lang("off") + ` : ` + lang("empty") + `)</span>` + Main_hr())
	data.WriteString(Setting_input("title_topic_max_length", values["title_topic_max_length"], "text") + Main_hr())
	data.WriteString(`<span>` + lang("max_file_size") + ` (MB)</span>` + Main_hr())
	data.WriteString(Setting_input("upload", values["upload"], "text") + Main_hr())
	data.WriteString(`<label><input type="checkbox" name="history_recording_off" ` + Setting_checked(values["history_recording_off"]) + `> ` + lang("set_history_recording_off") + `</label>` + Main_hr())
	data.WriteString(`<label><input type="checkbox" name="move_with_redirect" ` + Setting_checked(values["move_with_redirect"]) + `> ` + lang("move_with_redirect") + ` (` + lang("not_working") + `)</label>` + Main_hr())
	data.WriteString(`<span>` + lang("slow_thread") + ` (` + lang("second") + `) (` + lang("off") + ` : ` + lang("empty") + `)</span>` + Main_hr())
	data.WriteString(Setting_input("slow_thread", values["slow_thread"], "text") + Main_hr())
	data.WriteString(`<span>` + lang("edit_timeout") + ` (` + lang("second") + `) (` + lang("off") + ` : ` + lang("empty") + `) (` + lang("linux_only") + `)</span>` + Main_hr())
	data.WriteString(Setting_input("edit_timeout", values["edit_timeout"], "text") + Main_hr())
	data.WriteString(`<span>` + lang("document_content_max_length") + ` (` + lang("off") + ` : ` + lang("empty") + `)</span>` + Main_hr())
	data.WriteString(Setting_input("document_content_max_length", values["document_content_max_length"], "text") + Main_hr())
	data.WriteString(`<span>` + lang("bbs_content_max_length") + ` (` + lang("off") + ` : ` + lang("empty") + `)</span>` + Main_hr())
	data.WriteString(Setting_input("bbs_content_max_length", values["bbs_content_max_length"], "text") + Main_hr())

	data.WriteString(`<button id="opennamu_save_button" type="submit">` + lang("save") + `</button></form>`)

	return Setting_page(db, config, lang("main_setting"), data.String(), "setting")
}
