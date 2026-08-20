package route

import (
	stdjson "encoding/json"
	"sort"

	"opennamu/route/tool"
)

func View_main_manager(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	lang := func(name string) string {
		return tool.Get_language(db, name, true)
	}

	data := `<h2>` + lang("admin") + `</h2>
<ul>
<li><a href="/manager/2">` + lang("document_setting") + `</a></li>
<li><a href="/acl_multiple">` + lang("mutiple_document_setting") + `</a></li>
<li><a href="/manager/3">` + lang("check_user") + `</a></li>
<li><a href="/auth/ban">` + lang("ban") + `</a></li>
<li><a href="/auth/ban/multiple">` + lang("multiple_ban") + `</a></li>
<li><a href="/manager/5">` + lang("authorize") + `</a></li>
<li><a href="/auth/give">` + lang("multiple_authorize") + `</a></li>
<li><a href="/auth/give_total">` + lang("auth_to_auth") + `</a></li>
<li><a href="/delete_multiple">` + lang("many_delete") + `</a></li>
<li><a href="/app_submit">` + lang("application_list") + `</a></li>
</ul>
<h2>` + lang("owner") + `</h2>
<ul>
<li><a href="/auth/list">` + lang("admin_group_list") + `</a></li>
<li><a href="/register">` + lang("add_user") + `</a></li>
<li><a href="/setting">` + lang("setting") + `</a></li>
<li><a href="/manager/18">` + lang("user_fix") + `</a></li>
</ul>
<h3>` + lang("filter") + `</h3>
<ul>
<li><a href="/filter/edit_filter">` + lang("edit_filter_list") + `</a></li>
<li><a href="/filter/inter_wiki">` + lang("interwiki_list") + `</a></li>
<li><a href="/filter/edit_top">` + lang("edit_tool_list") + `</a></li>
<li><a href="/filter/image_license">` + lang("image_license_list") + `</a></li>
<li><a href="/filter/email_filter">` + lang("email_filter_list") + `</a></li>
<li><a href="/filter/name_filter">` + lang("id_filter_list") + `</a></li>
<li><a href="/filter/file_filter">` + lang("file_filter_list") + `</a></li>
<li><a href="/filter/extension_filter">` + lang("extension_filter_list") + `</a></li>
<li><a href="/filter/document">` + lang("document_filter_list") + `</a></li>
<li><a href="/filter/outer_link">` + lang("outer_link_filter_list") + `</a> (` + lang("beta") + `)</li>
<li><a href="/filter/template">` + lang("template_document_list") + `</a> (` + lang("beta") + `)</li>
</ul>
<h3>` + lang("server") + `</h3>
<ul>
<li><a href="/restart">` + lang("wiki_restart") + `</a></li>
<li><a href="/shutdown">` + lang("wiki_shutdown") + `</a></li>
<li><a href="/update">` + lang("update") + `</a></li>
</ul>`

	version_list := tool.Get_last_version()
	latest_version := get_remote_version(get_version_branch(db))
	data += `<h2>` + lang("version") + `</h2><ul><li>` + lang("version") + ` : ` + tool.HTML_escape(version_list["r_ver"]) + `</li><li>` + lang("lastest") + ` : ` + tool.HTML_escape(latest_version) + `</li></ul>`

	skin_html := `<h3>` + lang("skin_info") + `</h3><ul><li><a href="/api/skin_info/all">` + lang("skin_info") + `</a></li>`
	if raw_skin_data, ok := Api_skin_info_all(config); ok {
		skin_data := map[string]map[string]any{}
		if stdjson.Unmarshal(raw_skin_data, &skin_data) == nil {
			skin_names := []string{}
			for name := range skin_data {
				skin_names = append(skin_names, name)
			}
			sort.Strings(skin_names)
			for _, name := range skin_names {
				info := skin_data[name]
				display_name := name
				if value, ok := info["name"].(string); ok && value != "" {
					display_name = value
				}
				version, _ := info["skin_ver"].(string)
				if latest, ok := info["lastest_version"].(map[string]any); ok {
					if latest_version, ok := latest["skin_ver"].(string); ok && latest_version != "" {
						version += " (" + latest_version + ")"
					}
				}
				skin_html += `<li>` + tool.HTML_escape(display_name) + ` : ` + tool.HTML_escape(version) + `</li>`
			}
		}
	}
	skin_html += `</ul>`
	data += skin_html

	return tool.Get_template(
		db,
		config,
		lang("admin_tool"),
		data,
		[]any{},
		[][]any{{"other", lang("return")}},
		map[string]string{},
	)
}
