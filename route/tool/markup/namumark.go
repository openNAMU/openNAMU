package markup

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"opennamu/route/tool"

	"github.com/dlclark/regexp2"
)

type namumark struct {
	db   *sql.DB
	data map[string]string

	render_data string
}

func render_parameter_value(value any) string {
	if value == nil {
		return ""
	}
	if value_string, ok := value.(string); ok {
		return value_string
	}
	return fmt.Sprint(value)
}

func render_parameter_data(data string, parameter map[string]any) string {
	if len(parameter) == 0 {
		return data
	}

	parameter_regexp := regexp.MustCompile("@([ㄱ-ㅣ가-힣a-zA-Z0-9_]+)(?:=([^@\\n]+))?@")
	return parameter_regexp.ReplaceAllStringFunc(data, func(value string) string {
		match := parameter_regexp.FindStringSubmatch(value)
		parameter_value, ok := parameter[match[1]]
		if !ok {
			return match[2]
		}
		return render_parameter_value(parameter_value)
	})
}

func Namumark_new(db *sql.DB, data map[string]string, parameters ...map[string]any) *namumark {
	parameter_data := map[string]any{}
	if len(parameters) > 0 && parameters[0] != nil {
		parameter_data = parameters[0]
	}
	data_string := data["data"]
	data_string = render_parameter_data(data_string, parameter_data)
	data_string = "\n" + data_string + "\n"
	data_string = strings.ReplaceAll(data_string, "\r", "")

	return &namumark{
		db,
		data,

		data_string,
	}
}

type replacer struct {
	re     *regexp.Regexp
	prefix string
}

func (class *namumark) render_text() {
	string_data := class.render_data

	replacers := []replacer{
		{regexp.MustCompile(`'''([^\n]+?)'''`), "b"},
		{regexp.MustCompile(`''([^\n]+?)''`), "i"},
		{regexp.MustCompile(`__([^\n]+?)__`), "u"},
		{regexp.MustCompile(`\^\^\^([^\n]+?)\^\^\^`), "sup"},
		{regexp.MustCompile(`\^\^([^\n]+?)\^\^`), "sup"},
		{regexp.MustCompile(`,,,([^\n]+?),,,`), "sub"},
		{regexp.MustCompile(`,,([^\n]+?),,`), "sub"},
		{regexp.MustCompile(`--([^\n]+?)--`), "s"},
		{regexp.MustCompile(`~~([^\n]+?)~~`), "s"},
	}

	for _, rep := range replacers {
		for {
			next_data := rep.re.ReplaceAllString(string_data, "["+rep.prefix+"($1)]")
			if next_data == string_data {
				break
			}
			string_data = next_data
		}
	}

	class.render_data = string_data
}

func (class *namumark) render_heading() {
	string_data := class.render_data

	r := regexp.MustCompile(`\n(?:(={1,6})(#?) ?([^\n]+))\n`)
	r_sub := regexp.MustCompile(` ?(#?={1,6}[^=]*)$`)
	string_data = r.ReplaceAllStringFunc(string_data, func(m string) string {
		match := r.FindStringSubmatch(m)

		heading_data := r_sub.ReplaceAllString(match[3], "")

		heading_len := strconv.Itoa(tool.Get_len(match[1]))
		heading_render := "[h" + heading_len + "(" + heading_data + ")]"

		return heading_render
	})

	class.render_data = string_data
}

func (class *namumark) render_macro() {
	string_data := class.render_data

	r := regexp2.MustCompile(`\[([^\[\]\(\)]+)\]`, 0)

	string_data, _ = r.ReplaceFunc(
		string_data,
		func(m regexp2.Match) string {
			macro_name := m.GroupByNumber(1).String()

			if macro_name == "toc" || macro_name == "목차" {
				return "[toc()]"
			} else {
				return ""
			}
		},
		-1,
		-1,
	)

	class.render_data = string_data
}

func normalize_namumark_link(target string) string {
	target = strings.TrimSpace(target)
	prefix_list := []struct {
		prefix string
		name   string
	}{
		{":category:", "category:"}, {"category:", "category:"}, {":분류:", "category:"}, {"분류:", "category:"},
		{":file:", "file:"}, {"file:", "file:"}, {":파일:", "file:"}, {"파일:", "file:"},
		{"user:", "user:"}, {"사용자:", "user:"},
	}
	lower_target := strings.ToLower(target)
	for _, prefix_data := range prefix_list {
		if strings.HasPrefix(lower_target, strings.ToLower(prefix_data.prefix)) {
			return prefix_data.name + target[len(prefix_data.prefix):]
		}
	}
	return target
}

func (class *namumark) render_link() {
	string_data := class.render_data

	r := regexp2.MustCompile(`\[\[((?:(?!\[\[|\]\]|\|).)+)(?:\|((?:(?!\[\[|\]\]).)+))?\]\]`, 0)

	string_data, _ = r.ReplaceFunc(
		string_data,
		func(m regexp2.Match) string {
			target := normalize_namumark_link(m.GroupByNumber(1).String())
			label := m.GroupByNumber(2).String()
			if label == "" {
				label = target
			}

			tag_name := "a"
			if strings.HasPrefix(strings.ToLower(target), "http://") || strings.HasPrefix(strings.ToLower(target), "https://") {
				tag_name = "ex"
			}
			if tag_name == "a" {
				anchor := ""
				if hash_index := strings.Index(target, "#"); hash_index >= 0 {
					anchor = "#" + tool.Url_parser(target[hash_index+1:])
					target = target[:hash_index]
				}
				if anchor != "" {
					return "[" + tag_name + "(" + target + "," + label + "," + anchor + ")]"
				}
			}

			return "[" + tag_name + "(" + target + "," + label + ")]"
		},
		-1,
		-1,
	)

	class.render_data = string_data
}

func (class *namumark) render_last() {
	string_data := class.render_data

	r := regexp.MustCompile(`(\n| )+$`)
	string_data = r.ReplaceAllString(string_data, "")

	r = regexp.MustCompile(`^(\n| )+`)
	string_data = r.ReplaceAllString(string_data, "")

	r = regexp.MustCompile(`\n?<front_br>`)
	string_data = r.ReplaceAllString(string_data, "")

	r = regexp.MustCompile(`<back_br>\n?`)
	string_data = r.ReplaceAllString(string_data, "")

	class.render_data = string_data
}

func (class *namumark) main() map[string]any {
	class.render_text()
	class.render_link()
	class.render_heading()
	class.render_macro()
	class.render_last()

	class.data["data"] = class.render_data

	render_data_class := Macromark_new(class.db, class.data, "html")
	render_data := render_data_class.main()

	return render_data
}
