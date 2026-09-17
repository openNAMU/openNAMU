package markup

import (
	"database/sql"
	"regexp"
	"strconv"
	"strings"

	"opennamu/route/tool"
	backlink_parser "opennamu/route/tool/markup/backlink"

	"github.com/dlclark/regexp2"
)

type markdown struct {
	db   *sql.DB
	data map[string]string

	render_data string

	temp_data       [][]string
	temp_data_count int
}

func Markdown_new(db *sql.DB, data map[string]string) *markdown {
	data_string := data["data"]
	data_string = "\n" + strings.ReplaceAll(data_string, "\r", "") + "\n"

	return &markdown{
		db,
		data,
		data_string,
		[][]string{},
		0,
	}
}

var (
	markdown_fenced_code_regex = regexp2.MustCompile("(?ms)^ {0,3}\\x60\\x60\\x60([^\\n]*)\\n(.*?)^ {0,3}\\x60\\x60\\x60[ \\t]*(?:\\n|$)", 0)
	markdown_tilde_code_regex  = regexp2.MustCompile("(?ms)^ {0,3}~~~([^\\n]*)\\n(.*?)^ {0,3}~~~[ \\t]*(?:\\n|$)", 0)
	markdown_inline_code_regex = regexp2.MustCompile("`([^`\\n]+)`", 0)
	markdown_autolink_regex    = regexp2.MustCompile("(?i)<(https?://[^<>\\s]+)>", 0)
	markdown_image_regex       = regexp2.MustCompile(`!\[([^\]\r\n]*)\]\(((?:[^()\r\n]+|\([^()\r\n]*\))+?)\)`, 0)
	markdown_link_regex        = regexp2.MustCompile(`\[([^\]\r\n]*)\]\(((?:[^()\r\n]+|\([^()\r\n]*\))+?)\)`, 0)

	markdown_strong_regex      = regexp2.MustCompile(`\*\*([^*\r\n]+)\*\*`, 0)
	markdown_strong_underscore = regexp2.MustCompile(`__([^_\r\n]+)__`, 0)
	markdown_strike_regex      = regexp2.MustCompile(`~~([^~\r\n]+)~~`, 0)
	markdown_em_regex          = regexp2.MustCompile(`\*([^*\r\n]+)\*`, 0)
	markdown_em_underscore     = regexp2.MustCompile(`_([^_\r\n]+)_`, 0)

	markdown_heading_regex    = regexp.MustCompile(`(?m)^ {0,3}(#{1,6})[ \t]+([^\n]+?)[ \t]*(?:\n|$)`)
	markdown_setext_regex     = regexp.MustCompile(`(?m)^([^\n]+)\n[ \t]*(=+|-+)[ \t]*(?:\n|$)`)
	markdown_horizontal_regex = regexp.MustCompile(`(?m)^[ \t]{0,3}(?:(?:\*[ \t]*){3,}|(?:-[ \t]*){3,}|(?:_[ \t]*){3,})[ \t]*(?:\n|$)`)
	markdown_quote_regex      = regexp.MustCompile(`(?m)(?:^ {0,3}> ?[^\n]*(?:\n|$))+`)
	markdown_table_separator  = regexp.MustCompile(`^[ \t]*\|?[ \t]*:?-+:?[ \t]*(?:\|[ \t]*:?-+:?[ \t]*)+\|?[ \t]*$`)
)

type markdown_replacer struct {
	regex      *regexp2.Regexp
	macro_name string
}

func (class *markdown) Func_temp_save(data string) string {
	name := "OpenNamuMarkdownTemp" + strconv.Itoa(class.temp_data_count)
	class.temp_data = append(class.temp_data, []string{name, data})
	class.temp_data_count++

	return name
}

func (class *markdown) Func_temp_restore(data string) string {
	string_data := data
	for for_a := len(class.temp_data) - 1; for_a >= 0; for_a-- {
		string_data = strings.ReplaceAll(string_data, class.temp_data[for_a][0], class.temp_data[for_a][1])
	}

	return string_data
}

func Markdown_replace(regex *regexp2.Regexp, data string, fn func(regexp2.Match) string) string {
	data, _ = regex.ReplaceFunc(data, fn, -1, -1)
	return data
}

func Markdown_group(match regexp2.Match, number int) string {
	return match.GroupByNumber(number).String()
}

func Markdown_destination(data string) string {
	data = strings.TrimSpace(data)
	if strings.HasPrefix(data, "<") && strings.Contains(data, ">") {
		return data[1:strings.Index(data, ">")]
	}

	if index := strings.IndexAny(data, " \t\r\n"); index >= 0 {
		return data[:index]
	}

	return data
}

func Markdown_macro_argument(data string) string {
	return strings.ReplaceAll(data, ",", ",,")
}

func Markdown_is_external_link(data string) bool {
	data = strings.ToLower(strings.TrimSpace(data))
	return strings.HasPrefix(data, "http://") || strings.HasPrefix(data, "https://")
}

func (class *markdown) Render_code() {
	for _, regex := range []*regexp2.Regexp{
		markdown_fenced_code_regex,
		markdown_tilde_code_regex,
	} {
		class.render_data = Markdown_replace(regex, class.render_data, func(match regexp2.Match) string {
			language := strings.TrimSpace(Markdown_group(match, 1))
			language_class := ""
			if language != "" {
				language_class = " class=\"language-" + tool.HTML_escape(strings.Fields(language)[0]) + "\""
			}

			return class.Func_temp_save(
				"<pre><code" + language_class + ">" + tool.HTML_escape(Markdown_group(match, 2)) + "</code></pre>",
			)
		})
	}

	class.render_data = Markdown_replace(markdown_inline_code_regex, class.render_data, func(match regexp2.Match) string {
		return class.Func_temp_save("<code>" + tool.HTML_escape(Markdown_group(match, 1)) + "</code>")
	})
}

func (class *markdown) Render_escape() {
	string_data := class.render_data
	escape_characters := "\\`*_{}[]()#+-.!<>~"

	for index := 0; index < len(string_data)-1; index++ {
		if string_data[index] != '\\' || !strings.ContainsRune(escape_characters, rune(string_data[index+1])) {
			continue
		}

		string_data = string_data[:index] + class.Func_temp_save(tool.HTML_escape(string(string_data[index+1]))) + string_data[index+2:]
	}

	class.render_data = string_data
}

func (class *markdown) Render_image() {
	class.render_data = Markdown_replace(markdown_image_regex, class.render_data, func(match regexp2.Match) string {
		return "[img(" + Markdown_macro_argument(Markdown_destination(Markdown_group(match, 2))) + "," + Markdown_macro_argument(Markdown_group(match, 1)) + ")]"
	})
}

func (class *markdown) Render_link() {
	class.render_data = Markdown_replace(markdown_autolink_regex, class.render_data, func(match regexp2.Match) string {
		link := Markdown_group(match, 1)
		return "[ex(" + Markdown_macro_argument(link) + "," + Markdown_macro_argument(link) + ")]"
	})

	class.render_data = Markdown_replace(markdown_link_regex, class.render_data, func(match regexp2.Match) string {
		label := Markdown_group(match, 1)
		target := Markdown_destination(Markdown_group(match, 2))
		if label == "" {
			label = target
		}

		if strings.HasPrefix(strings.ToLower(target), "mailto:") {
			return label
		}
		if strings.HasPrefix(target, "#") {
			return "[an(" + Markdown_macro_argument(target) + "," + Markdown_macro_argument(label) + ")]"
		}
		if strings.HasPrefix(target, "/") || strings.HasPrefix(target, "./") || strings.HasPrefix(target, "../") {
			return "[in(" + Markdown_macro_argument(target) + "," + Markdown_macro_argument(label) + ")]"
		}
		if Markdown_is_external_link(target) {
			return "[ex(" + Markdown_macro_argument(target) + "," + Markdown_macro_argument(label) + ")]"
		}

		link := target
		hash := ""
		if index := strings.Index(link, "#"); index >= 0 {
			hash = link[index:]
			link = link[:index]
		}

		result := "[a(" + Markdown_macro_argument(link) + "," + Markdown_macro_argument(label)
		if hash != "" {
			result += "," + Markdown_macro_argument(hash)
		}

		return result + ")]"
	})
}

func (class *markdown) Render_text() {
	replacers := []markdown_replacer{
		{regexp2.MustCompile(`\*\*\*([^*\r\n]+)\*\*\*`, 0), "i"},
		{markdown_strong_regex, "b"},
		{markdown_strong_underscore, "b"},
		{markdown_strike_regex, "s"},
		{markdown_em_regex, "i"},
		{markdown_em_underscore, "i"},
	}

	for _, replacer := range replacers {
		class.render_data = Markdown_replace(replacer.regex, class.render_data, func(match regexp2.Match) string {
			return "[" + replacer.macro_name + "(" + Markdown_group(match, 1) + ")]"
		})
	}
}

func (class *markdown) Render_heading() {
	string_data := markdown_heading_regex.ReplaceAllStringFunc(class.render_data, func(match string) string {
		data := markdown_heading_regex.FindStringSubmatch(match)
		heading := strings.TrimSpace(regexp.MustCompile(`[ \t]+#+[ \t]*$`).ReplaceAllString(data[2], ""))
		return "[h" + strconv.Itoa(len(data[1])) + "(" + heading + ")]\n"
	})

	class.render_data = markdown_setext_regex.ReplaceAllStringFunc(string_data, func(match string) string {
		data := markdown_setext_regex.FindStringSubmatch(match)
		heading := 2
		if strings.HasPrefix(data[2], "=") {
			heading = 1
		}
		return "[h" + strconv.Itoa(heading) + "(" + strings.TrimSpace(data[1]) + ")]\n"
	})
}

func (class *markdown) Render_lists() {
	lines := strings.Split(class.render_data, "\n")
	result := []string{}

	for index := 0; index < len(lines); index++ {
		line := strings.TrimLeft(lines[index], " ")
		list_type := ""
		if len(line) > 2 && (line[0] == '-' || line[0] == '+' || line[0] == '*') && (line[1] == ' ' || line[1] == '\t') {
			list_type = "ul"
		}

		marker := 0
		for marker < len(line) && line[marker] >= '0' && line[marker] <= '9' {
			marker++
		}
		if marker > 0 && marker+1 < len(line) && (line[marker] == '.' || line[marker] == ')') && (line[marker+1] == ' ' || line[marker+1] == '\t') {
			list_type = "ol"
		}

		if list_type == "" {
			result = append(result, lines[index])
			continue
		}

		items := []string{}
		for index < len(lines) {
			item_line := strings.TrimLeft(lines[index], " ")
			item := ""
			if list_type == "ul" && len(item_line) > 2 && (item_line[0] == '-' || item_line[0] == '+' || item_line[0] == '*') {
				item = strings.TrimSpace(item_line[2:])
			} else if list_type == "ol" {
				end := 0
				for end < len(item_line) && item_line[end] >= '0' && item_line[end] <= '9' {
					end++
				}
				if end+1 < len(item_line) && (item_line[end] == '.' || item_line[end] == ')') {
					item = strings.TrimSpace(item_line[end+1:])
				}
			}
			if item == "" {
				break
			}
			items = append(items, "[li("+item+")]")
			index++
		}

		result = append(result, "["+list_type+"("+strings.Join(items, "")+")]")
		index--
	}

	class.render_data = strings.Join(result, "\n")
}

func (class *markdown) Render_quote() {
	class.render_data = markdown_quote_regex.ReplaceAllStringFunc(class.render_data, func(match string) string {
		lines := strings.Split(strings.TrimSuffix(match, "\n"), "\n")
		for index := range lines {
			lines[index] = strings.TrimPrefix(strings.TrimPrefix(strings.TrimLeft(lines[index], " "), ">"), " ")
		}
		return "[q(" + strings.Join(lines, "\n") + ")]\n"
	})
}

func Markdown_table_cells(data string) []string {
	data = strings.Trim(strings.TrimSpace(data), "|")
	cells := strings.Split(data, "|")
	for index := range cells {
		cells[index] = strings.TrimSpace(cells[index])
	}
	return cells
}

func (class *markdown) Render_table() {
	lines := strings.Split(class.render_data, "\n")
	result := []string{}

	for index := 0; index < len(lines); index++ {
		if index+1 >= len(lines) || !strings.Contains(lines[index], "|") || !markdown_table_separator.MatchString(lines[index+1]) {
			result = append(result, lines[index])
			continue
		}

		rows := []string{}
		for _, cell := range Markdown_table_cells(lines[index]) {
			rows = append(rows, "[th("+cell+")]")
		}
		table := []string{"[tr(" + strings.Join(rows, "") + ")]"}
		index += 2

		for index < len(lines) && strings.Contains(lines[index], "|") && strings.TrimSpace(lines[index]) != "" {
			rows = []string{}
			for _, cell := range Markdown_table_cells(lines[index]) {
				rows = append(rows, "[td("+cell+")]")
			}
			table = append(table, "[tr("+strings.Join(rows, "")+")]")
			index++
		}
		result = append(result, "[table("+strings.Join(table, "")+")]")
		index--
	}

	class.render_data = strings.Join(result, "\n")
}

func (class *markdown) Render_last() {
	class.render_data = strings.Trim(class.render_data, "\n ")

	renderer := Macromark_new(class.db, map[string]string{
		"data":        class.render_data,
		"doc_name":    class.data["doc_name"],
		"render_name": class.data["render_name"],
		"render_type": class.data["render_type"],
		"from":        class.data["from"],
		"include":     class.data["include"],
	}, "html")
	result := renderer.Main()
	class.render_data = class.Func_temp_restore(result["data"].(string))
}

func (class *markdown) Main() map[string]any {
	class.Render_code()
	class.Render_escape()
	class.Render_image()
	class.Render_link()
	class.Render_text()
	class.Render_table()
	class.Render_lists()
	class.Render_quote()
	class.Render_heading()
	class.render_data = markdown_horizontal_regex.ReplaceAllString(class.render_data, "[hr()]\n")
	class.Render_last()

	backlink_list, link_count, _ := backlink_parser.Get_backlink(class.data["data"], "markdown")
	backlinks := [][]string{}
	for link, link_type_list := range backlink_list {
		for _, link_type := range link_type_list {
			backlinks = append(backlinks, []string{class.data["doc_name"], link, link_type, ""})
		}
	}

	return map[string]any{
		"data":       class.render_data,
		"js_data":    "",
		"backlink":   backlinks,
		"link_count": link_count,
	}
}

func Markdown(db *sql.DB, data map[string]string) map[string]any {
	return Markdown_new(db, data).Main()
}
