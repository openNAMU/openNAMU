package route

import (
	"database/sql"
	"strings"

	"opennamu/route/tool"
)

func setting_main_fields() []setting_field {
	return []setting_field{
		{name: "name", default_value: "Wiki"},
		{name: "frontpage", default_value: "FrontPage"},
		{name: "upload", default_value: "2"},
		{name: "skin", default_value: ""},
		{name: "reg", default_value: ""},
		{name: "ip_view", default_value: ""},
		{name: "back_up", default_value: ""},
		{name: "port", default_value: "3000"},
		{name: "key", default_value: tool.Get_random_key(128)},
		{name: "update", default_value: "stable"},
		{name: "encode", default_value: "sha3"},
		{name: "host", default_value: "0.0.0.0"},
		{name: "slow_edit", default_value: ""},
		{name: "requires_approval", default_value: ""},
		{name: "backup_where", default_value: ""},
		{name: "domain", default_value: ""},
		{name: "ua_get", default_value: ""},
		{name: "enable_comment", default_value: ""},
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
		{name: "link_case_insensitive", default_value: ""},
		{name: "move_with_redirect", default_value: ""},
		{name: "slow_thread", default_value: ""},
		{name: "edit_timeout", default_value: "5"},
		{name: "document_content_max_length", default_value: ""},
		{name: "backup_count", default_value: ""},
		{name: "ua_expiration_date", default_value: ""},
		{name: "auth_history_expiration_date", default_value: ""},
		{name: "auth_history_off", default_value: ""},
		{name: "user_name_level", default_value: ""},
		{name: "load_ip_select", default_value: ""},
		{name: "not_use_view_count", default_value: ""},
	}
}

func View_setting_main(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	return view_setting_main_data(db, config, setting_load_fields(db, setting_main_fields()))
}

func View_setting_main_post(config tool.Config, form map[string]string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	fields := setting_main_fields()
	setting_save_fields(db, fields, form)
	tool.Do_insert_auth_history(db, config.IP, "edit_set (main)")

	return tool.Get_redirect("/setting/main")
}

func view_setting_main_data(db *sql.DB, config tool.Config, values map[string]string) string {
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
	data.WriteString(`<span>` + lang("wiki_name") + `</span>` + setting_hr())
	data.WriteString(setting_input("name", values["name"], "text") + setting_hr())
	data.WriteString(`<span><a href="/setting/main/logo">(` + lang("wiki_logo") + `)</a></span>` + setting_hr())
	data.WriteString(`<span>` + lang("main_page") + `</span>` + setting_hr())
	data.WriteString(setting_input("frontpage", values["frontpage"], "text") + setting_hr())
	data.WriteString(`<span>` + lang("tls_method") + `</span>` + setting_hr())
	data.WriteString(`<select name="http_select">` + setting_options(values["http_select"], []string{"http", "https"}, nil) + `</select>` + setting_hr())
	data.WriteString(`<span>` + lang("domain") + `</span> (EX : 2du.pythonanywhere.com) (` + lang("off") + ` : ` + lang("empty") + `)` + setting_hr())
	data.WriteString(setting_input("domain", values["domain"], "text") + setting_hr())
	data.WriteString(`<span>` + lang("wiki_host") + `</span>` + setting_hr())
	data.WriteString(setting_input("host", values["host"], "text") + setting_hr())
	data.WriteString(`<span>` + lang("wiki_port") + `</span>` + setting_hr())
	data.WriteString(setting_input("port", values["port"], "text") + setting_hr())
	data.WriteString(`<span>` + lang("wiki_secret_key") + `</span>` + setting_hr())
	data.WriteString(setting_input("key", values["key"], "password") + setting_hr())
	data.WriteString(`<label><input type="checkbox" name="wiki_access_password_need" ` + setting_checked(values["wiki_access_password_need"]) + `> ` + lang("set_wiki_access_password_need") + ` (` + lang("restart_required") + `)</label>` + setting_hr())
	data.WriteString(`<span>` + lang("set_wiki_access_password") + `</span> (` + lang("restart_required") + `)` + setting_hr())
	data.WriteString(setting_input("wiki_access_password", values["wiki_access_password"], "password") + setting_hr())
	data.WriteString(`<span>` + lang("wiki_load_ip_select") + `</span> (` + lang("restart_required") + `)` + setting_hr())
	data.WriteString(`<select name="load_ip_select">` + setting_options(values["load_ip_select"], ip_values, ip_labels) + `</select>` + setting_hr())
	data.WriteString(`<h3>` + lang("authority_use_list") + `</h3>`)
	data.WriteString(`<label><input type="checkbox" name="auth_history_off" ` + setting_checked(values["auth_history_off"]) + `> ` + lang("authority_use_list_off") + `</label>` + setting_hr())
	data.WriteString(`<span>` + lang("authority_use_list_expiration_date") + `</span> (` + lang("day") + `) (` + lang("off") + ` : ` + lang("empty") + `)` + setting_hr())
	data.WriteString(setting_input("auth_history_expiration_date", values["auth_history_expiration_date"], "text") + setting_hr())
	data.WriteString(`<h3>` + lang("communication_set") + `</h3>`)
	data.WriteString(`<label><input type="checkbox" name="enable_comment" ` + setting_checked(values["enable_comment"]) + `> ` + lang("enable_comment_function") + `</label>` + setting_hr())
	data.WriteString(`<label><input type="checkbox" name="user_name_level" ` + setting_checked(values["user_name_level"]) + `> ` + lang("display_level_in_user_name") + `</label>` + setting_hr())
	data.WriteString(`<label><input type="checkbox" name="not_use_view_count" ` + setting_checked(values["not_use_view_count"]) + `> ` + lang("not_use_view_count") + `</label>` + setting_hr())

	data.WriteString(`<h2>` + lang("design_set") + `</h2>`)
	data.WriteString(`<span>` + lang("wiki_skin") + `</span>` + setting_hr())
	data.WriteString(`<select name="skin">` + setting_options(skin_value, skin_values, nil) + `</select>` + setting_hr())

	data.WriteString(`<h2>` + lang("render_set") + `</h2>`)
	data.WriteString(`<label><input type="checkbox" name="namumark_compatible" ` + setting_checked(values["namumark_compatible"]) + `> ` + lang("namumark_fully_compatible_mode") + `</label>` + setting_hr())
	data.WriteString(`<label><input type="checkbox" name="link_case_insensitive" ` + setting_checked(values["link_case_insensitive"]) + `> ` + lang("link_case_insensitive") + `</label>` + setting_hr())

	data.WriteString(`<h2>` + lang("login_set") + `</h2>`)
	data.WriteString(`<label><input type="checkbox" name="reg" ` + setting_checked(values["reg"]) + `> ` + lang("no_register") + `</label>` + setting_hr())
	data.WriteString(`<label><input type="checkbox" name="ip_view" ` + setting_checked(values["ip_view"]) + `> ` + lang("hide_ip") + `</label>` + setting_hr())
	data.WriteString(`<label><input type="checkbox" name="user_name_view" ` + setting_checked(values["user_name_view"]) + `> ` + lang("hide_user_name") + `</label>` + setting_hr())
	data.WriteString(`<label><input type="checkbox" name="requires_approval" ` + setting_checked(values["requires_approval"]) + `> ` + lang("requires_approval") + `</label>` + setting_hr())
	data.WriteString(`<span>` + lang("password_min_length") + `</span> (` + lang("off") + ` : ` + lang("empty") + `)` + setting_hr())
	data.WriteString(setting_input("password_min_length", values["password_min_length"], "text") + setting_hr())
	data.WriteString(`<span>` + lang("encryption_method") + `</span>` + setting_hr())
	data.WriteString(`<select name="encode">` + setting_options(values["encode"], encode_values, nil) + `</select>` + setting_hr())

	data.WriteString(`<h3>` + lang("ua") + `</h3>`)
	data.WriteString(`<label><input type="checkbox" name="ua_get" ` + setting_checked(values["ua_get"]) + `> ` + lang("ua_get_off") + `</label>` + setting_hr())
	data.WriteString(`<span>` + lang("ua_expiration_date") + `</span> (` + lang("day") + `) (` + lang("off") + ` : ` + lang("empty") + `)` + setting_hr())
	data.WriteString(setting_input("ua_expiration_date", values["ua_expiration_date"], "text") + setting_hr())

	data.WriteString(`<h2>` + lang("server_set") + `</h2>`)
	data.WriteString(`<span>` + lang("update_branch") + `</span>` + setting_hr())
	data.WriteString(`<select name="update">` + setting_options(values["update"], []string{"stable", "dev", "beta"}, nil) + `</select>` + setting_hr())
	data.WriteString(`<span` + sqlite_only + `>`)
	data.WriteString(`<h3>` + lang("backup") + ` (` + lang("sqlite_only") + `)</h3>`)
	data.WriteString(`<span>` + lang("backup_warning") + ` (EX : data_YYYYMMDDHHMMSS.db)</span>` + setting_hr())
	data.WriteString(`<span>` + lang("backup_interval") + ` (` + lang("hour") + `) (` + lang("off") + ` : ` + lang("empty") + `)</span>` + setting_hr())
	data.WriteString(setting_input("back_up", values["back_up"], "text") + setting_hr())
	data.WriteString(`<span>` + lang("backup_where") + ` (` + lang("default") + ` : ` + lang("empty") + `) (` + lang("example") + ` : ./data/backup.db)</span>` + setting_hr())
	data.WriteString(setting_input("backup_where", values["backup_where"], "text") + setting_hr())
	data.WriteString(`<span>` + lang("backup_count") + ` (` + lang("default") + ` : ` + lang("empty") + `)</span>` + setting_hr())
	data.WriteString(setting_input("backup_count", values["backup_count"], "text") + setting_hr())
	data.WriteString(`</span>`)

	data.WriteString(`<h2>` + lang("edit_set") + `</h2>`)
	data.WriteString(`<span>` + lang("slow_edit") + ` (` + lang("second") + `) (` + lang("off") + ` : ` + lang("empty") + `)</span>` + setting_hr())
	data.WriteString(setting_input("slow_edit", values["slow_edit"], "text") + setting_hr())
	data.WriteString(`<label><input type="checkbox" name="edit_bottom_compulsion" ` + setting_checked(values["edit_bottom_compulsion"]) + `> ` + lang("edit_bottom_compulsion") + `</label>` + setting_hr())
	data.WriteString(`<span>` + lang("title_max_length") + ` (` + lang("off") + ` : ` + lang("empty") + `)</span>` + setting_hr())
	data.WriteString(setting_input("title_max_length", values["title_max_length"], "text") + setting_hr())
	data.WriteString(`<span>` + lang("title_topic_max_length") + ` (` + lang("off") + ` : ` + lang("empty") + `)</span>` + setting_hr())
	data.WriteString(setting_input("title_topic_max_length", values["title_topic_max_length"], "text") + setting_hr())
	data.WriteString(`<span>` + lang("max_file_size") + ` (MB)</span>` + setting_hr())
	data.WriteString(setting_input("upload", values["upload"], "text") + setting_hr())
	data.WriteString(`<label><input type="checkbox" name="history_recording_off" ` + setting_checked(values["history_recording_off"]) + `> ` + lang("set_history_recording_off") + `</label>` + setting_hr())
	data.WriteString(`<label><input type="checkbox" name="move_with_redirect" ` + setting_checked(values["move_with_redirect"]) + `> ` + lang("move_with_redirect") + ` (` + lang("not_working") + `)</label>` + setting_hr())
	data.WriteString(`<span>` + lang("slow_thread") + ` (` + lang("second") + `) (` + lang("off") + ` : ` + lang("empty") + `)</span>` + setting_hr())
	data.WriteString(setting_input("slow_thread", values["slow_thread"], "text") + setting_hr())
	data.WriteString(`<span>` + lang("edit_timeout") + ` (` + lang("second") + `) (` + lang("off") + ` : ` + lang("empty") + `) (` + lang("linux_only") + `)</span>` + setting_hr())
	data.WriteString(setting_input("edit_timeout", values["edit_timeout"], "text") + setting_hr())
	data.WriteString(`<span>` + lang("document_content_max_length") + ` (` + lang("off") + ` : ` + lang("empty") + `)</span>` + setting_hr())
	data.WriteString(setting_input("document_content_max_length", values["document_content_max_length"], "text") + setting_hr())

	data.WriteString(`<button id="opennamu_save_button" type="submit">` + lang("save") + `</button></form>`)

	return setting_page(db, config, lang("main_setting"), data.String(), "setting")
}
