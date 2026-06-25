package markup

import (
	"database/sql"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/dlclark/regexp2"
)

type namumark struct {
    db   *sql.DB
    data map[string]string

    render_data string
}

func Namumark_new(db *sql.DB, data map[string]string) *namumark {
    data_string := data["data"]
    data_string = "\n" + data_string + "\n"
    data_string = strings.ReplaceAll(data_string, "\r", "")

    return &namumark{
        db,
        data,

        data_string,
    }
}

type replacer struct {
    re *regexp2.Regexp
    prefix string
}

func (class *namumark) render_text() {
    string_data := class.render_data

    replacers := []replacer{
        { regexp2.MustCompile(`'''((?:(?!''').)+)'''`, 0), "b" },
        { regexp2.MustCompile(`''((?:(?!'').)+)''`, 0), "i" },
        { regexp2.MustCompile(`__((?:(?!__).)+)__`, 0), "u" },
        { regexp2.MustCompile(`\^\^\^((?:(?!\^\^\^).)+)\^\^\^`, 0), "sup" },
        { regexp2.MustCompile(`\^\^((?:(?!\^\^).)+)\^\^`, 0), "sup" },
        { regexp2.MustCompile(`,,,((?:(?!,,,).)+),,,`, 0), "sub" },
        { regexp2.MustCompile(`,,((?:(?!,,).)+),,`, 0), "sub" },
        { regexp2.MustCompile(`--((?:(?!--).)+)--`, 0), "s" },
        { regexp2.MustCompile(`~~((?:(?!~~).)+)~~`, 0), "s" },
    }

    for _, rep := range replacers {
        for {
            m, _ := rep.re.FindStringMatch(string_data)
            if m == nil {
                break
            }
            
            gps := m.Groups()

            start := m.Index
            end := start + len(m.String())
            replacement := "[" + rep.prefix + "(" + gps[1].Captures[0].String() + ")]"

            string_data = string_data[:start] + replacement + string_data[end:]
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

        heading_len := strconv.Itoa(len(match[1]))
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

func (class *namumark) render_link() {
    string_data := class.render_data

    r := regexp2.MustCompile(`\[\[((?:(?!\[\[|\]\]|\|).)+)(?:\|((?:(?!\[\[|\]\]).)+))?\]\]`, 0)

    string_data, _ = r.ReplaceFunc(
        string_data,
        func(m regexp2.Match) string {
            target := m.GroupByNumber(1).String()
            label := m.GroupByNumber(2).String()
            if label == "" {
                label = target
            }

            tag_name := "a"
            if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
                tag_name = "ex"
            }

            log.Default().Println(tag_name, target)

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

    log.Default().Println(class.render_data)

    class.data["data"] = class.render_data

    render_data_class := Macromark_new(class.db, class.data, "html")
    render_data := render_data_class.main()

    return render_data
}
