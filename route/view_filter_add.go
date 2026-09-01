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
	if values == nil && !tool.Check_permission(db, "filter_manage", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	if values != nil {
		result := Api_filter_add_post(config, kind, name, values)
		if result["response"] == "redirect" {
			return tool.Get_redirect("/filter/document")
		}
		if result["response"] != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
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
			sub_type := tool.Get_html_filter_inter_wiki_sub(db, name)
			if sub_type == "under_bar" {
				inter_type = "under_bar"
			}
			form += `<hr class="main_hr">` + filter_select("inter_type", []string{"url_encode", "under_bar"}, inter_type)
		}
	case "external_image":
		form += filter_input(tool.Get_language(db, "domain", true), "title", value[0])
	case "html":
		form += filter_input(tool.Get_language(db, "tag", true), "title", value[0])
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
		form += `<hr class="main_hr"><span>` + tool.Get_language(db, "acl", true) + `</span><hr class="main_hr"><textarea name="acl" placeholder="view=normal&#10;edit=trust_a&#10;move=owner&#10;delete=owner&#10;new_make=trust_a">` + tool.HTML_escape(value[2]) + `</textarea>`
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
		"html":             {"html", "html_filter_list"},
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

var html_filter_tag_regex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-]*$`)

var html_filter_blocked_tags = map[string]bool{
	"embed": true, "object": true, "script": true, "style": true,
}

func filter_value(db *sql.DB, kind string, name string) []string {
	return tool.Get_html_filter_value(db, name, kind)
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

func document_filter_acl_data(db *sql.DB, data string) (string, bool) {
	data = strings.ReplaceAll(data, "\r", "")
	if !strings.Contains(data, "=") {
		data = strings.TrimSpace(data)
		if data == "normal" || acl_value_valid(db, data) {
			return data, true
		}
		return "", false
	}

	valid_action := map[string]bool{
		"view":     true,
		"edit":     true,
		"move":     true,
		"delete":   true,
		"new_make": true,
	}
	seen := map[string]bool{}
	acl_list := []string{}
	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return "", false
		}
		action := strings.TrimSpace(parts[0])
		acl := strings.TrimSpace(parts[1])
		if !valid_action[action] || seen[action] || (acl != "normal" && !acl_value_valid(db, acl)) {
			return "", false
		}

		seen[action] = true
		acl_list = append(acl_list, action+"="+acl)
	}

	if len(acl_list) == 0 {
		return "", false
	}
	return strings.Join(acl_list, "\n"), true
}
