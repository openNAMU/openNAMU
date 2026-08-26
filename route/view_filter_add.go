package route

import (
	"database/sql"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"opennamu/route/tool"
)

func View_filter_add(config tool.Config, kind string, name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	spec, ok := get_filter_spec(kind)
	if !ok {
		return tool.Get_error_page(db, config, "error")
	}
	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	if values != nil {
		title := values.Get("title")
		if title == "" {
			title = values.Get("name")
		}
		if title == "" {
			title = name
		}
		if title == "" && kind != "external_image" {
			title = "test"
		}

		if kind == "external_image" {
			title = strings.ToLower(strings.TrimSpace(title))
			parsed, err := url.Parse("https://" + title)
			if title == "" || err != nil || parsed.Host != title || parsed.Hostname() == "" || parsed.Port() != "" || parsed.User != nil {
				return tool.Get_error_page(db, config, "error")
			}
			if name != "" && name != title {
				tool.Exec_DB(db, "delete from html_filter where html = ? and kind = ?", name, spec.db_kind)
			}
			tool.Exec_DB(db, "delete from html_filter where html = ? and kind = ?", title, spec.db_kind)
			tool.Exec_DB(db, "insert into html_filter (html, kind, plus, plus_t) values (?, ?, '', '')", title, spec.db_kind)
		} else if kind == "inter_wiki" || kind == "outer_link" {
			if name != "" && name != title {
				tool.Exec_DB(db, "delete from html_filter where html = ? and kind = ?", name, spec.db_kind)
			}
			tool.Exec_DB(db, "delete from html_filter where html = ? and kind = ?", title, spec.db_kind)
			tool.Exec_DB(db, "insert into html_filter (html, plus, plus_t, kind) values (?, ?, ?, ?)", title, values.Get("link"), values.Get("icon"), spec.db_kind)
			if kind == "inter_wiki" {
				if name != "" && name != title {
					tool.Exec_DB(db, "delete from html_filter where html = ? and kind = 'inter_wiki_sub'", name)
				}
				tool.Exec_DB(db, "delete from html_filter where html = ? and kind = 'inter_wiki_sub'", title)
				tool.Exec_DB(db, "insert into html_filter (html, plus, plus_t, kind) values (?, 'inter_wiki_type', ?, 'inter_wiki_sub')", title, values.Get("inter_type"))
			}
		} else if kind == "edit_filter" {
			content := values.Get("content")
			if _, err := regexp.Compile(content); err != nil {
				return tool.Get_error_page(db, config, "error")
			}
			filter_name := name
			if filter_name == "" {
				filter_name = title
			}
			end := "X"
			if days, err := strconv.Atoi(strings.TrimSpace(values.Get("day"))); err == nil && days > 0 {
				end = strconv.Itoa(days * 24 * 60 * 60)
			}
			tool.Exec_DB(db, "delete from html_filter where html = ? and kind = 'regex_filter'", filter_name)
			tool.Exec_DB(db, "insert into html_filter (html, plus, plus_t, kind) values (?, ?, ?, 'regex_filter')", filter_name, content, end)
		} else if kind == "document" {
			doc_name := values.Get("name")
			if doc_name == "" {
				return tool.Get_redirect("/filter/document")
			}
			if _, err := regexp.Compile(values.Get("regex")); err != nil {
				return tool.Get_error_page(db, config, "error")
			}
			if name != "" && name != doc_name {
				tool.Exec_DB(db, "delete from html_filter where html = ? and kind = 'document'", name)
			}
			tool.Exec_DB(db, "delete from html_filter where html = ? and kind = 'document'", doc_name)
			tool.Exec_DB(db, "insert into html_filter (html, kind, plus, plus_t) values (?, 'document', ?, ?)", doc_name, values.Get("regex"), values.Get("acl"))
		} else {
			plus := ""
			switch kind {
			case "name_filter", "file_filter":
				if _, err := regexp.Compile(title); err != nil {
					return tool.Get_error_page(db, config, "error")
				}
			case "extension_filter":
				plus = regexp.MustCompile(`[^0-9]`).ReplaceAllString(values.Get("max_file_size"), "")
			case "template":
				plus = values.Get("exp")
			case "edit_top":
				plus = values.Get("markup")
			}
			if name != "" && name != title {
				tool.Exec_DB(db, "delete from html_filter where html = ? and kind = ?", name, spec.db_kind)
			}
			tool.Exec_DB(db, "delete from html_filter where html = ? and kind = ?", title, spec.db_kind)
			tool.Exec_DB(db, "insert into html_filter (html, kind, plus, plus_t) values (?, ?, ?, '')", title, spec.db_kind, plus)
		}

		tool.Do_insert_auth_history(db, config.IP, "filter_save ("+kind+")")
		return tool.Get_redirect("/filter/" + kind)
	}

	value := filter_value(db, spec.db_kind, name)
	title := tool.Get_language(db, spec.title, true)
	form := `<form method="post">`
	switch kind {
	case "inter_wiki", "outer_link":
		form += filter_input(tool.Get_language(db, "name", true), "title", value[0])
		form += `<hr class="main_hr">` + filter_input(tool.Get_language(db, "link", true), "link", value[1])
		form += `<hr class="main_hr">` + filter_input(tool.Get_language(db, "icon", true), "icon", value[2])
		if kind == "inter_wiki" {
			inter_type := "url_encode"
			sub_type := ""
			tool.QueryRow_DB(db, "select plus_t from html_filter where html = ? and kind = 'inter_wiki_sub'", []any{&sub_type}, name)
			if sub_type == "under_bar" {
				inter_type = "under_bar"
			}
			form += `<hr class="main_hr">` + filter_select("inter_type", []string{"url_encode", "under_bar"}, inter_type)
		}
	case "external_image":
		form += filter_input(tool.Get_language(db, "domain", true), "title", value[0])
	case "edit_filter":
		end := ""
		if value[2] != "" && value[2] != "X" {
			seconds, _ := strconv.Atoi(value[2])
			end = strconv.Itoa(seconds / (24 * 60 * 60))
		}
		if name == "" {
			form += filter_input(tool.Get_language(db, "name", true), "title", "")
		} else {
			form += `<input type="hidden" name="title" value="` + tool.HTML_escape(name) + `">`
		}
		form += filter_input(tool.Get_language(db, "day", true), "day", end)
		form += `<hr class="main_hr">` + filter_input(tool.Get_language(db, "regex", true), "content", value[1])
	case "document":
		form += filter_input(tool.Get_language(db, "name", true), "name", name)
		form += `<hr class="main_hr">` + filter_input(tool.Get_language(db, "regex", true), "regex", value[1])
		form += `<hr class="main_hr">` + filter_select("acl", tool.List_acl("normal"), value[2])
	case "name_filter", "file_filter":
		form += filter_input(tool.Get_language(db, "regex", true), "title", name)
	case "email_filter":
		form += filter_input(tool.Get_language(db, "email", true), "title", name)
	case "image_license":
		form += filter_input(tool.Get_language(db, "license", true), "title", name)
	case "extension_filter":
		form += filter_input(tool.Get_language(db, "extension", true), "title", name)
		form += `<hr class="main_hr">` + filter_input(tool.Get_language(db, "max_file_size", true), "max_file_size", value[1])
	case "template":
		form += filter_input(tool.Get_language(db, "template", true), "title", name)
		form += `<hr class="main_hr">` + filter_input(tool.Get_language(db, "explanation", true), "exp", value[1])
	default:
		form += filter_input(tool.Get_language(db, "title", true), "title", name)
		form += `<hr class="main_hr">` + filter_input(tool.Get_language(db, "markup", true), "markup", value[1])
	}
	form += `<hr class="main_hr"><button type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`

	return tool.Get_template(db, config, title, form, []any{}, [][]any{{"filter/" + kind, tool.Get_language(db, "return", true)}}, map[string]string{})
}

type filter_spec struct {
	db_kind string
	title   string
}

func get_filter_spec(kind string) (filter_spec, bool) {
	list := map[string]filter_spec{
		"inter_wiki":       {"inter_wiki", "interwiki_list"},
		"outer_link":       {"outer_link", "outer_link_filter_list"},
		"external_image":   {"external_image", "external_image_filter_list"},
		"document":         {"document", "document_filter_list"},
		"edit_top":         {"edit_top", "edit_tool_list"},
		"image_license":    {"image_license", "image_license_list"},
		"template":         {"template", "template_document_list"},
		"edit_filter":      {"regex_filter", "edit_filter_list"},
		"email_filter":     {"email", "email_filter_list"},
		"file_filter":      {"file", "file_filter_list"},
		"name_filter":      {"name", "id_filter_list"},
		"extension_filter": {"extension", "extension_filter_list"},
	}

	spec, ok := list[kind]
	return spec, ok
}

func filter_value(db *sql.DB, kind string, name string) []string {
	html := ""
	plus := ""
	plus_t := ""
	if tool.QueryRow_DB(db, "select html, plus, plus_t from html_filter where html = ? and kind = ? limit 1", []any{&html, &plus, &plus_t}, name, kind) {
		return []string{html, plus, plus_t}
	}
	return []string{"", "", ""}
}

func filter_safe_link(value string) string {
	parsed, err := url.Parse(value)
	if err != nil || (parsed.Scheme != "" && parsed.Scheme != "http" && parsed.Scheme != "https") {
		return "#"
	}
	return tool.HTML_escape(value)
}

func filter_input(label string, name string, value string) string {
	return `<span>` + label + `</span><hr class="main_hr"><input name="` + name + `" value="` + tool.HTML_escape(value) + `">`
}

func filter_select(name string, values []string, selected string) string {
	data := `<select name="` + name + `">`
	for _, value := range values {
		option_value := tool.HTML_escape(value)
		option_name := option_value
		if value == "" {
			option_name = "normal"
		}
		selected_text := ""
		if value == selected {
			selected_text = ` selected`
		}
		data += `<option value="` + option_value + `"` + selected_text + `>` + option_name + `</option>`
	}
	return data + `</select>`
}
