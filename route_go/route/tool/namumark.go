package tool

import (
    "database/sql"
    "regexp"
    "log"
    "strconv"
    "strings"

    "github.com/dlclark/regexp2"
)

type namumark struct {
    db *sql.DB
    data map[string]string

    render_data string
}

func Namumark_new(db *sql.DB, data map[string]string) *namumark {
    data_string := data["data"]
    data_string = "\n" + data_string + "\n"
    data_string = strings.Replace(data_string, "\r", "", -1)

    return &namumark{
        db,
        data,

        data_string,
    }
}

func (class *namumark) render_text() {
    string_data := class.render_data

    r := regexp2.MustCompile(`'''((?:(?!''').)+)'''`, 0)
    for {
        if m, _ := r.FindStringMatch(string_data); m != nil {
            gps := m.Groups()
            m_string := m.String()

            string_data = strings.Replace(string_data, m_string, "[b(" + gps[1].Captures[0].String() + ")]", 1)
        } else {
            break
        }
    }

    r = regexp2.MustCompile(`''((?:(?!'').)+)''`, 0)
    for {
        if m, _ := r.FindStringMatch(string_data); m != nil {
            gps := m.Groups()
            m_string := m.String()

            string_data = strings.Replace(string_data, m_string, "[i(" + gps[1].Captures[0].String() + ")]", 1)
        } else {
            break
        }
    }

    r = regexp2.MustCompile(`__((?:(?!__).)+)__`, 0)
    for {
        if m, _ := r.FindStringMatch(string_data); m != nil {
            gps := m.Groups()
            m_string := m.String()

            string_data = strings.Replace(string_data, m_string, "[u(" + gps[1].Captures[0].String() + ")]", 1)
        } else {
            break
        }
    }

    r_list := []string{`\^\^\^((?:(?!\^\^\^).)+)\^\^\^`, `\^\^((?:(?!\^\^).)+)\^\^`}
    for for_a := 0; for_a < len(r_list); for_a++ {
        r = regexp2.MustCompile(r_list[for_a], 0)
        for {
            if m, _ := r.FindStringMatch(string_data); m != nil {
                gps := m.Groups()
                m_string := m.String()

                string_data = strings.Replace(string_data, m_string, "[sup(" + gps[1].Captures[0].String() + ")]", 1)
            } else {
                break
            }
        }
    }

    r_list = []string{`,,,((?:(?!,,,).)+),,,`, `,,((?:(?!,,).)+),,`}
    for for_a := 0; for_a < len(r_list); for_a++ {
        r = regexp2.MustCompile(r_list[for_a], 0)
        for {
            if m, _ := r.FindStringMatch(string_data); m != nil {
                gps := m.Groups()
                m_string := m.String()

                string_data = strings.Replace(string_data, m_string, "[sub(" + gps[1].Captures[0].String() + ")]", 1)
            } else {
                break
            }
        }
    }

    r_list = []string{`--((?:(?!--).)+)--`, `~~((?:(?!~~).)+)~~`}
    for for_a := 0; for_a < len(r_list); for_a++ {
        r = regexp2.MustCompile(r_list[for_a], 0)
        for {
            if m, _ := r.FindStringMatch(string_data); m != nil {
                gps := m.Groups()
                m_string := m.String()

                string_data = strings.Replace(string_data, m_string, "[s(" + gps[1].Captures[0].String() + ")]", 1)
            } else {
                break
            }
        }
    }

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

    string_data = strings.Replace(string_data, "\n", "<br>", -1)

    class.render_data = string_data
}


func (class *namumark) render_heading() {
    string_data := class.render_data

    r := regexp.MustCompile(`\n(?:(={1,6})(#?) ?([^\n]+))\n`)
    string_data = r.ReplaceAllStringFunc(string_data, func(m string) string {
        match := r.FindStringSubmatch(m)

        r = regexp.MustCompile(` ?(#?={1,6}[^=]*)$`)
        heading_data := r.ReplaceAllString(match[3], "")

        heading_len := strconv.Itoa(len(match[1]))
        heading_render := "[h" + heading_len + "(" + heading_data + ")]"

        return heading_render
    })

    class.render_data = string_data
}

func (class namumark) main() map[string]interface{} {
    class.render_text()
    class.render_heading()
    class.render_last()

    log.Default().Println(class.render_data)

    class.data["data"] = class.render_data

    render_data_class := Macromark_new(class.db, class.data)

    render_data := make(map[string]interface{})
    render_data = render_data_class.main()

    return render_data
}
