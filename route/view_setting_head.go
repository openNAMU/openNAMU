package route

import (
	"database/sql"
	"strings"

	"opennamu/route/tool"
)

func View_setting_head(config tool.Config, kind string, skin_name string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	name, coverage, title_key, action, ok := setting_head_info(kind, skin_name)
	if !ok {
		return tool.Get_error_page(db, config, "error")
	}

	return view_setting_head_data(db, config, kind, skin_name, name, coverage, title_key, action, setting_value(db, name, coverage, ""), "", false)
}

func setting_head_info(kind string, skin_name string) (string, string, string, string, bool) {
	switch kind {
	case "body/top":
		return "body", "", "main_body", "/setting/body/top", true
	case "body/bottom":
		return "bottom_body", "", "main_bottom_body", "/setting/body/bottom", true
	case "head":
		return "head", skin_name, "main_head", "/setting/head", true
	default:
		return "", "", "", "", false
	}
}

func view_setting_head_data(db *sql.DB, config tool.Config, kind string, skin_name string, name string, coverage string, title_key string, action string, value string, preview string, is_preview bool) string {
	lang := func(key string) string {
		return tool.Get_language(db, key, true)
	}
	if kind == "head" && skin_name != "" {
		action += "/" + tool.Url_parser(skin_name)
	}

	data := strings.Builder{}
	data.WriteString(`<form method="post" action="` + tool.HTML_escape(action) + `">`)

	if kind == "head" {
		data.WriteString(`<a href="/setting/head">(` + lang("all") + `)</a> `)
		for _, skin := range tool.Get_skin_list("", false) {
			data.WriteString(`<a href="/setting/head/` + tool.Url_parser(skin) + `">(` + tool.HTML_escape(skin) + `)</a> `)
			data.WriteString(`<a href="/setting/head/` + tool.Url_parser(skin+"-cssdark") + `">(` + tool.HTML_escape(skin) + `-cssdark)</a> `)
		}
		data.WriteString(setting_hr())
		data.WriteString(`<span>&lt;style&gt;CSS&lt;/style&gt;<br>&lt;script&gt;JS&lt;/script&gt;</span>` + setting_hr())
	}

	textarea_value := value
	if is_preview {
		textarea_value = preview
	}

	data.WriteString(`<textarea class="opennamu_textarea_500" placeholder="` + lang("enter_html") + `" name="content" id="content">` + tool.HTML_escape(textarea_value) + `</textarea>`)
	data.WriteString(setting_hr())
	if kind == "head" {
		data.WriteString(lang("main_css_warning") + setting_hr())
	}
	data.WriteString(`<button id="opennamu_save_button" type="submit">` + lang("save") + `</button>`)

	if kind == "body/top" || kind == "body/bottom" {
		preview_action := "/setting_preview/" + kind
		data.WriteString(` <button id="opennamu_preview_button" type="submit" formaction="` + preview_action + `">` + lang("preview") + `</button>`)
		if is_preview {
			data.WriteString(setting_hr() + `<div id="opennamu_preview_area">` + preview + `</div>`)
		}
	}
	data.WriteString(`</form>`)

	title := lang(title_key)
	if skin_name != "" {
		title += " (" + tool.HTML_escape(skin_name) + ")"
	}

	return setting_page(db, config, title, data.String(), "setting")
}
