package markup

import (
	"database/sql"
	"regexp"
	"strconv"
	"strings"

	"opennamu/route/tool"
	backlink_parser "opennamu/route/tool/markup/backlink"
)

// markdown converts the small Markdown dialect used by OpenNamu into the
// Macromark intermediate syntax. Keeping this as a string based pass makes
// Markdown and Namumark share the same final HTML renderer.
type markdown struct {
	db   *sql.DB
	data map[string]string

	render_data string

	temp_data       [][]string
	temp_data_count int
}

func Markdown_new(db *sql.DB, data map[string]string) *markdown {
	data_string := data["data"]
	data_string = strings.ReplaceAll(data_string, "\r", "")
	data_string = "\n" + data_string + "\n"

	return &markdown{
		db:              db,
		data:            data,
		render_data:     data_string,
		temp_data:       [][]string{},
		temp_data_count: 0,
	}
}

var (
	markdown_fenced_code_regex = regexp.MustCompile("(?ms)^ {0,3}\\x60\\x60\\x60([^\\n]*)\\n(.*?)^ {0,3}\\x60\\x60\\x60[ \\t]*(?:\\n|$)")
	markdown_tilde_code_regex  = regexp.MustCompile("(?ms)^ {0,3}~~~([^\\n]*)\\n(.*?)^ {0,3}~~~[ \\t]*(?:\\n|$)")
	markdown_inline_code_regex = regexp.MustCompile(string(rune(96)) + "([^" + string(rune(96)) + "\\n]+)" + string(rune(96)))
	markdown_autolink_regex    = regexp.MustCompile("(?i)<(https?://[^<>\\s]+)>")

	markdown_strong_regex      = regexp.MustCompile("\\*\\*([^*\\r\\n]+)\\*\\*")
	markdown_strong_underscore = regexp.MustCompile("__([^_\\r\\n]+)__")
	markdown_strike_regex      = regexp.MustCompile("~~([^~\\r\\n]+)~~")
	markdown_em_regex          = regexp.MustCompile("\\*([^*\\r\\n]+)\\*")
	markdown_em_underscore     = regexp.MustCompile("_([^_\\r\\n]+)_")
	markdown_heading_regex     = regexp.MustCompile("(?m)^ {0,3}(#{1,6})[ \\t]+([^\\n]+?)[ \\t]*(?:\\n|$)")
	markdown_setext_regex      = regexp.MustCompile("(?m)^([^\\n]+)\\n[ \\t]*(=+|-+)[ \\t]*(?:\\n|$)")
	markdown_table_separator   = regexp.MustCompile("^[ \\t]*\\|?[ \\t]*:?-+:?[ \\t]*(?:\\|[ \\t]*:?-+:?[ \\t]*)+\\|?[ \\t]*$")
	markdown_horizontal_regex  = regexp.MustCompile("(?m)^[ \\t]{0,3}(?:(?:\\*[ \\t]*){3,}|(?:-[ \\t]*){3,}|(?:_[ \\t]*){3,})[ \\t]*(?:\\n|$)")
)

func (class *markdown) func_temp_save(data string) string {
	name := "OpenNamuMarkdownTemp" + strconv.Itoa(class.temp_data_count)
	class.temp_data = append(class.temp_data, []string{name, data})
	class.temp_data_count++

	return name
}

func (class *markdown) func_temp_restore(data string) string {
	string_data := data
	for index := len(class.temp_data) - 1; index >= 0; index-- {
		string_data = strings.ReplaceAll(string_data, class.temp_data[index][0], class.temp_data[index][1])
	}

	return string_data
}

func (class *markdown) replace_code_block(regex *regexp.Regexp, data string) string {
	return regex.ReplaceAllStringFunc(data, func(match string) string {
		parts := regex.FindStringSubmatch(match)
		if len(parts) < 3 {
			return match
		}

		language := strings.TrimSpace(parts[1])
		language_class := ""
		if language != "" {
			language = strings.Fields(language)[0]
			language_class = " class=\"language-" + tool.HTML_escape(language) + "\""
		}

		return class.func_temp_save(
			"<pre><code" + language_class + ">" + tool.HTML_escape(parts[2]) + "</code></pre>",
		)
	})
}

func (class *markdown) render_code() {
	string_data := class.render_data
	string_data = class.replace_code_block(markdown_fenced_code_regex, string_data)
	string_data = class.replace_code_block(markdown_tilde_code_regex, string_data)

	string_data = markdown_inline_code_regex.ReplaceAllStringFunc(string_data, func(match string) string {
		parts := markdown_inline_code_regex.FindStringSubmatch(match)
		if len(parts) < 2 {
			return match
		}

		return class.func_temp_save("<code>" + tool.HTML_escape(parts[1]) + "</code>")
	})

	class.render_data = string_data
}

func is_markdown_escape(data byte) bool {
	escape_characters := "\\" + string(rune(96)) + "*_{}[]()#+-.!<>~"
	return strings.ContainsRune(escape_characters, rune(data))
}

func (class *markdown) render_escape() {
	string_data := class.render_data
	var result strings.Builder

	for index := 0; index < len(string_data); index++ {
		if string_data[index] == '\\' && index+1 < len(string_data) && is_markdown_escape(string_data[index+1]) {
			result.WriteString(class.func_temp_save(tool.HTML_escape(string(string_data[index+1]))))
			index++
			continue
		}

		result.WriteByte(string_data[index])
	}

	class.render_data = result.String()
}

func markdown_destination(data string) string {
	data = strings.TrimSpace(data)
	if strings.HasPrefix(data, "<") && strings.Contains(data, ">") {
		data = data[1:strings.Index(data, ">")]
	} else if space_index := strings.IndexAny(data, " \t\r\n"); space_index >= 0 {
		data = data[:space_index]
	}

	return data
}

func markdown_find_link(data string, cursor int, image bool) (int, int, string, string, bool) {
	for cursor < len(data) {
		relative_start := strings.IndexByte(data[cursor:], '[')
		if relative_start < 0 {
			return 0, 0, "", "", false
		}

		bracket_start := cursor + relative_start
		link_start := bracket_start
		if image {
			if bracket_start == 0 || data[bracket_start-1] != '!' {
				cursor = bracket_start + 1
				continue
			}
			link_start = bracket_start - 1
		} else if bracket_start > 0 && data[bracket_start-1] == '!' {
			cursor = bracket_start + 1
			continue
		}

		bracket_end_relative := strings.IndexByte(data[bracket_start+1:], ']')
		if bracket_end_relative < 0 {
			return 0, 0, "", "", false
		}
		bracket_end := bracket_start + 1 + bracket_end_relative
		if bracket_end+1 >= len(data) || data[bracket_end+1] != '(' {
			cursor = bracket_start + 1
			continue
		}

		depth := 1
		destination_end := -1
		for index := bracket_end + 2; index < len(data); index++ {
			switch data[index] {
			case '(':
				depth++
			case ')':
				depth--
				if depth == 0 {
					destination_end = index
				}
			}
			if destination_end >= 0 {
				break
			}
		}

		if destination_end < 0 {
			return 0, 0, "", "", false
		}

		return link_start, destination_end + 1,
			data[bracket_start+1 : bracket_end],
			data[bracket_end+2 : destination_end],
			true
	}

	return 0, 0, "", "", false
}

func markdown_replace_link_syntax(data string, image bool, replace func(label string, destination string) string) string {
	var result strings.Builder
	last := 0
	cursor := 0

	for {
		start, end, label, destination, ok := markdown_find_link(data, cursor, image)
		if !ok {
			break
		}

		result.WriteString(data[last:start])
		result.WriteString(replace(label, destination))
		last = end
		cursor = end
	}

	result.WriteString(data[last:])
	return result.String()
}

func markdown_macro_argument(data string) string {
	return strings.ReplaceAll(data, ",", ",,")
}

func markdown_is_external_link(data string) bool {
	data = strings.ToLower(strings.TrimSpace(data))
	return strings.HasPrefix(data, "http://") || strings.HasPrefix(data, "https://")
}

func (class *markdown) render_image() {
	class.render_data = markdown_replace_link_syntax(class.render_data, true, func(label string, destination string) string {
		return "[img(" + markdown_macro_argument(markdown_destination(destination)) + "," + markdown_macro_argument(label) + ")]"
	})
}

func (class *markdown) render_link() {
	string_data := markdown_autolink_regex.ReplaceAllStringFunc(class.render_data, func(match string) string {
		parts := markdown_autolink_regex.FindStringSubmatch(match)
		if len(parts) < 2 {
			return match
		}

		return "[ex(" + markdown_macro_argument(parts[1]) + "," + markdown_macro_argument(parts[1]) + ")]"
	})

	string_data = markdown_replace_link_syntax(string_data, false, func(label string, destination string) string {
		target := markdown_destination(destination)
		if label == "" {
			label = target
		}

		if strings.HasPrefix(target, "#") {
			return "[an(" + markdown_macro_argument(target) + "," + markdown_macro_argument(label) + ")]"
		}

		if strings.HasPrefix(target, "/") ||
			strings.HasPrefix(target, "./") ||
			strings.HasPrefix(target, "../") {
			return "[in(" + markdown_macro_argument(target) + "," + markdown_macro_argument(label) + ")]"
		}

		if strings.HasPrefix(strings.ToLower(target), "mailto:") {
			return label
		}

		if markdown_is_external_link(target) {
			return "[ex(" + markdown_macro_argument(target) + "," + markdown_macro_argument(label) + ")]"
		}

		link := target
		hash := ""
		if hash_index := strings.Index(link, "#"); hash_index >= 0 {
			hash = link[hash_index:]
			link = link[:hash_index]
		}

		macro := "[a(" + markdown_macro_argument(link) + "," + markdown_macro_argument(label)
		if hash != "" {
			macro += "," + markdown_macro_argument(hash)
		}

		return macro + ")]"
	})

	class.render_data = string_data
}

func markdown_replace_inline(data string, regex *regexp.Regexp, macro_name string) string {
	for {
		changed := false
		data = regex.ReplaceAllStringFunc(data, func(match string) string {
			parts := regex.FindStringSubmatch(match)
			if len(parts) < 2 {
				return match
			}

			changed = true
			return "[" + macro_name + "(" + parts[1] + ")]"
		})

		if !changed {
			return data
		}
	}
}

func (class *markdown) render_text() {
	string_data := class.render_data
	string_data = markdown_replace_inline(string_data, regexp.MustCompile("\\*\\*\\*([^*\\r\\n]+)\\*\\*\\*"), "i")
	string_data = markdown_replace_inline(string_data, markdown_strong_regex, "b")
	string_data = markdown_replace_inline(string_data, markdown_strong_underscore, "b")
	string_data = markdown_replace_inline(string_data, markdown_strike_regex, "s")
	string_data = markdown_replace_inline(string_data, markdown_em_regex, "i")
	string_data = markdown_replace_inline(string_data, markdown_em_underscore, "i")

	class.render_data = string_data
}

func markdown_line_indent(data string) int {
	trimmed := strings.TrimLeft(data, " ")
	return len(data) - len(trimmed)
}

func markdown_list_marker(data string) (string, string, bool) {
	if markdown_line_indent(data) > 3 {
		return "", "", false
	}

	data = strings.TrimLeft(data, " ")
	if len(data) >= 2 && (data[0] == '-' || data[0] == '+' || data[0] == '*') && (data[1] == ' ' || data[1] == '\t') {
		return "ul", strings.TrimSpace(data[2:]), true
	}

	index := 0
	for index < len(data) && data[index] >= '0' && data[index] <= '9' {
		index++
	}
	if index > 0 && index+1 < len(data) && (data[index] == '.' || data[index] == ')') && (data[index+1] == ' ' || data[index+1] == '\t') {
		return "ol", strings.TrimSpace(data[index+1:]), true
	}

	return "", "", false
}

func (class *markdown) render_lists() {
	lines := strings.Split(class.render_data, "\n")
	result := make([]string, 0, len(lines))

	for index := 0; index < len(lines); {
		list_type, item, ok := markdown_list_marker(lines[index])
		if !ok {
			result = append(result, lines[index])
			index++
			continue
		}

		items := []string{"[li(" + item + ")]"}
		index++
		for index < len(lines) {
			next_type, next_item, next_ok := markdown_list_marker(lines[index])
			if !next_ok || next_type != list_type {
				break
			}

			items = append(items, "[li("+next_item+")]")
			index++
		}

		result = append(result, "["+list_type+"("+strings.Join(items, "")+")]")
	}

	class.render_data = strings.Join(result, "\n")
}

func (class *markdown) render_heading() {
	string_data := markdown_heading_regex.ReplaceAllStringFunc(class.render_data, func(match string) string {
		parts := markdown_heading_regex.FindStringSubmatch(match)
		if len(parts) < 3 {
			return match
		}

		heading_data := strings.TrimSpace(parts[2])
		heading_data = regexp.MustCompile("[ \\t]+#+[ \\t]*$").ReplaceAllString(heading_data, "")

		return "[h" + strconv.Itoa(len(parts[1])) + "(" + heading_data + ")]\n"
	})

	string_data = markdown_setext_regex.ReplaceAllStringFunc(string_data, func(match string) string {
		parts := markdown_setext_regex.FindStringSubmatch(match)
		if len(parts) < 3 {
			return match
		}

		heading_len := 2
		if strings.HasPrefix(parts[2], "=") {
			heading_len = 1
		}

		return "[h" + strconv.Itoa(heading_len) + "(" + strings.TrimSpace(parts[1]) + ")]\n"
	})

	class.render_data = string_data
}

func (class *markdown) render_horizontal() {
	class.render_data = markdown_horizontal_regex.ReplaceAllString(class.render_data, "[hr()]\n")
}

func (class *markdown) render_quote() {
	lines := strings.Split(class.render_data, "\n")
	result := make([]string, 0, len(lines))

	for index := 0; index < len(lines); {
		trimmed := strings.TrimLeft(lines[index], " ")
		if !strings.HasPrefix(trimmed, ">") || (len(trimmed) > 1 && trimmed[1] != ' ') {
			result = append(result, lines[index])
			index++
			continue
		}

		content := []string{}
		for index < len(lines) {
			next := strings.TrimLeft(lines[index], " ")
			if !strings.HasPrefix(next, ">") || (len(next) > 1 && next[1] != ' ') {
				break
			}

			next = strings.TrimPrefix(strings.TrimPrefix(next, ">"), " ")
			content = append(content, next)
			index++
		}

		result = append(result, "[q("+strings.Join(content, "\n")+")]")
	}

	class.render_data = strings.Join(result, "\n")
}

func markdown_table_cells(data string) []string {
	data = strings.TrimSpace(data)
	data = strings.TrimPrefix(data, "|")
	data = strings.TrimSuffix(data, "|")
	parts := strings.Split(data, "|")
	for index := range parts {
		parts[index] = strings.TrimSpace(parts[index])
	}

	return parts
}

func (class *markdown) render_table() {
	lines := strings.Split(class.render_data, "\n")
	result := make([]string, 0, len(lines))

	for index := 0; index < len(lines); {
		if index+1 >= len(lines) || !strings.Contains(lines[index], "|") || !markdown_table_separator.MatchString(lines[index+1]) {
			result = append(result, lines[index])
			index++
			continue
		}

		head := markdown_table_cells(lines[index])
		rows := []string{}
		head_cells := make([]string, 0, len(head))
		for _, cell := range head {
			head_cells = append(head_cells, "[th("+cell+")]")
		}
		rows = append(rows, "[tr("+strings.Join(head_cells, "")+")]")
		index += 2

		for index < len(lines) && strings.Contains(lines[index], "|") && strings.TrimSpace(lines[index]) != "" {
			cells := markdown_table_cells(lines[index])
			row_cells := make([]string, 0, len(cells))
			for _, cell := range cells {
				row_cells = append(row_cells, "[td("+cell+")]")
			}
			rows = append(rows, "[tr("+strings.Join(row_cells, "")+")]")
			index++
		}

		result = append(result, "[table("+strings.Join(rows, "")+")]")
	}

	class.render_data = strings.Join(result, "\n")
}

func (class *markdown) render_last() {
	class.render_data = strings.Trim(class.render_data, "\n ")

	render_data_class := Macromark_new(class.db, map[string]string{
		"data":        class.render_data,
		"doc_name":    class.data["doc_name"],
		"render_name": class.data["render_name"],
		"render_type": class.data["render_type"],
		"from":        class.data["from"],
		"include":     class.data["include"],
	}, "html")
	render_data := render_data_class.main()
	class.render_data = class.func_temp_restore(render_data["data"].(string))
}

func (class *markdown) main() map[string]any {
	class.render_code()
	class.render_escape()
	class.render_image()
	class.render_link()
	class.render_text()
	class.render_table()
	class.render_lists()
	class.render_quote()
	class.render_heading()
	class.render_horizontal()
	class.render_last()

	backlink_list, link_count, _ := backlink_parser.Get_backlink(class.data["data"], "markdown")
	end_backlink := [][]string{}
	for link, link_type_list := range backlink_list {
		for _, link_type := range link_type_list {
			end_backlink = append(end_backlink, []string{
				class.data["doc_name"],
				link,
				link_type,
				"",
			})
		}
	}

	return map[string]any{
		"data":       class.render_data,
		"js_data":    "",
		"backlink":   end_backlink,
		"link_count": link_count,
	}
}

func Markdown(db *sql.DB, data map[string]string) map[string]any {
	return Markdown_new(db, data).main()
}
