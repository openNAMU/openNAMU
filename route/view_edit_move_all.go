package route

import (
	"database/sql"
	"net/url"
	"strings"

	"opennamu/route/tool"
)

type move_all_document struct {
	old_name string
	new_name string
	data     string
}

func move_all_like_escape(data string) string {
	return strings.NewReplacer("!", "!!", "%", "!%", "_", "!_").Replace(data)
}

func move_all_pattern(source string, match_type string) string {
	switch match_type {
	case "end":
		return "%" + move_all_like_escape(source)
	case "include":
		return "%" + move_all_like_escape(source) + "%"
	default:
		return move_all_like_escape(source) + "%"
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
	documents := move_all_documents(db, source, target, match_type)

	body := "<h2>" + tool.Get_language(db, "move", true) + " preview</h2><ul>"
	for _, document := range documents {
		body += "<li><a href=\"/w/" + tool.Url_parser(document.old_name) + "\">" + tool.HTML_escape(document.old_name) + "</a> → " + tool.HTML_escape(document.new_name) + "</li>"
	}
	body += "</ul>"
	if len(documents) == 0 {
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
	send := ""
	execute := false
	if values != nil {
		source = strings.TrimSpace(values.Get("source"))
		target = strings.TrimSpace(values.Get("target"))
		send = values.Get("send")
		execute = values.Get("execute") == "1"
		match_type = values.Get("match")
		if match_type != "end" && match_type != "include" {
			match_type = "start"
		}
	}

	body := "<p>" + tool.Get_language(db, "multiple_move", true) + " preview</p>"
	body += "<form method=\"get\"><input name=\"source\" value=\"" + tool.HTML_escape(source) + "\" placeholder=\"" + tool.Get_language(db, "document_name", true) + "\">"
	body += "<hr class=\"main_hr\"><input name=\"target\" value=\"" + tool.HTML_escape(target) + "\" placeholder=\"" + tool.Get_language(db, "document_name", true) + "\">"
	body += "<hr class=\"main_hr\"><input name=\"send\" value=\"" + tool.HTML_escape(send) + "\" placeholder=\"" + tool.Get_language(db, "why", true) + "\">"
	body += "<hr class=\"main_hr\"><select name=\"match\">"
	body += "<option value=\"start\"" + move_all_selected(match_type, "start") + ">" + tool.Get_language(db, "title_start_document", true) + "</option>"
	body += "<option value=\"end\"" + move_all_selected(match_type, "end") + ">" + tool.Get_language(db, "title_end_document", true) + "</option>"
	body += "<option value=\"include\"" + move_all_selected(match_type, "include") + ">" + tool.Get_language(db, "title_include_document", true) + "</option>"
	body += "</select><hr class=\"main_hr\"><button type=\"submit\">" + tool.Get_language(db, "move", true) + "</button></form>"
	if execute {
		api_data := Api_edit_move_all_post(config, source, target, match_type, send)
		response, _ := api_data["response"].(string)
		if response == "ok" {
			body += "<hr class=\"main_hr\"><p>" + tool.Get_language(db, "multiple_move_complete", true) + "</p>"
		} else {
			body += "<hr class=\"main_hr\"><p>" + tool.Get_language(db, "move_error", true) + "</p>"
		}
	}
	if !execute && values != nil {
		body += "<hr class=\"main_hr\">" + move_all_rows(db, source, target, match_type)
		body += "<form method=\"post\"><input type=\"hidden\" name=\"source\" value=\"" + tool.HTML_escape(source) + "\">"
		body += "<input type=\"hidden\" name=\"target\" value=\"" + tool.HTML_escape(target) + "\"><input type=\"hidden\" name=\"match\" value=\"" + tool.HTML_escape(match_type) + "\">"
		body += "<input type=\"hidden\" name=\"send\" value=\"" + tool.HTML_escape(send) + "\"><input type=\"hidden\" name=\"execute\" value=\"1\">"
		body += "<button type=\"submit\">" + tool.Get_language(db, "move", true) + "</button></form>"
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

func move_all_match(title string, source string, match_type string) bool {
	switch match_type {
	case "end":
		return strings.HasSuffix(title, source)
	case "include":
		return strings.Contains(title, source)
	default:
		return title == source || strings.HasPrefix(title, source+"/")
	}
}

func move_all_documents(db *sql.DB, source string, target string, match_type string) []move_all_document {
	if source == "" || target == "" {
		return []move_all_document{}
	}

	rows := tool.Get_move_document_rows(db, move_all_pattern(source, match_type))
	defer rows.Close()

	result := []move_all_document{}
	for rows.Next() {
		title := ""
		if rows.Scan(&title) != nil || !move_all_match(title, source, match_type) {
			continue
		}

		new_title := move_all_title(title, source, target, match_type)
		if title == new_title || new_title == "" {
			continue
		}
		result = append(result, move_all_document{old_name: title, new_name: new_title})
	}
	return result
}
