package route

import (
	"database/sql"
	"net/url"
	"strings"

	"opennamu/route/tool"
)

func move_all_pattern(source string, match_type string) string {
	switch match_type {
	case "end":
		return "%" + source
	case "include":
		return "%" + source + "%"
	default:
		return source + "%"
	}
}

func move_all_title(title string, source string, target string, match_type string) string {
	switch match_type {
	case "end":
		return strings.TrimSuffix(title, source) + target
	case "include":
		return strings.ReplaceAll(title, source, target)
	default:
		return target + strings.TrimPrefix(title, source)
	}
}

func move_all_selected(current string, option string) string {
	if current == option {
		return " selected"
	}
	return ""
}

func move_all_rows(db *sql.DB, source string, target string, match_type string) string {
	if source == "" || target == "" {
		return ""
	}

	rows := tool.Get_move_document_rows(db, move_all_pattern(source, match_type))
	defer rows.Close()

	body := "<h2>" + tool.Get_language(db, "move", true) + " preview</h2><ul>"
	count := 0
	for rows.Next() {
		title := ""
		if rows.Scan(&title) != nil {
			continue
		}
		new_title := move_all_title(title, source, target, match_type)
		if title == new_title || new_title == "" {
			continue
		}
		body += "<li><a href=\"/w/" + tool.Url_parser(title) + "\">" + tool.HTML_escape(title) + "</a> → " + tool.HTML_escape(new_title) + "</li>"
		count++
	}
	body += "</ul>"
	if count == 0 {
		return "<p>" + tool.Get_language(db, "document_404_error", true) + "</p>"
	}
	return body
}

func View_edit_move_all(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "document_move_manage", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	source := ""
	target := ""
	match_type := "start"
	if values != nil {
		source = strings.TrimSpace(values.Get("source"))
		target = strings.TrimSpace(values.Get("target"))
		match_type = values.Get("match")
		if match_type != "end" && match_type != "include" {
			match_type = "start"
		}
	}

	body := "<p>" + tool.Get_language(db, "multiple_move", true) + " preview</p>"
	body += "<form method=\"get\"><input name=\"source\" value=\"" + tool.HTML_escape(source) + "\" placeholder=\"" + tool.Get_language(db, "document_name", true) + "\">"
	body += "<hr class=\"main_hr\"><input name=\"target\" value=\"" + tool.HTML_escape(target) + "\" placeholder=\"" + tool.Get_language(db, "document_name", true) + "\">"
	body += "<hr class=\"main_hr\"><select name=\"match\">"
	body += "<option value=\"start\"" + move_all_selected(match_type, "start") + ">" + tool.Get_language(db, "title_start_document", true) + "</option>"
	body += "<option value=\"end\"" + move_all_selected(match_type, "end") + ">" + tool.Get_language(db, "title_end_document", true) + "</option>"
	body += "<option value=\"include\"" + move_all_selected(match_type, "include") + ">" + tool.Get_language(db, "title_include_document", true) + "</option>"
	body += "</select><hr class=\"main_hr\"><button type=\"submit\">" + tool.Get_language(db, "move", true) + "</button></form>"
	if values != nil {
		body += "<hr class=\"main_hr\">" + move_all_rows(db, source, target, match_type)
	}

	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "multiple_move", true),
		body,
		[]any{},
		[][]any{{"other", tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}
