package tool

import (
    "database/sql"
    "html"
    "log"
    "regexp"
    "strconv"
    "strings"

    "github.com/dlclark/regexp2"
)

type macromark struct {
    db *sql.DB
    data map[string]string

    render_data string
    render_data_js string

    temp_data [][]string
    temp_data_raw [][]string
    temp_data_count int

    backlink [][]string
    link_count int
}

func Macromark_new(db *sql.DB, data map[string]string) *macromark {
    data_string := data["data"]
    data_string = html.EscapeString(data_string)
    data_string = strings.Replace(data_string, "\r", "", -1)
    data_string = "\n" + data_string + "\n"

    return &macromark{
        db,
        data,

        data_string,
        "",

        [][]string{},
        [][]string{},
        0,

        [][]string{},
        0,
    }
}

func (class *macromark) func_temp_save(data string, data_raw string) string {
    name := "<temp_save_" + strconv.Itoa(class.temp_data_count) + ">"

    class.temp_data = append(class.temp_data, []string{name, data})
    class.temp_data_raw = append(class.temp_data_raw, []string{name, data_raw})

    class.temp_data_count += 1

    return name
}

func (class macromark) func_temp_restore(data string, to_raw bool) string {
    string_data := data
 
    if to_raw {
        for for_a := len(class.temp_data_raw) - 1; for_a >= 0; for_a-- {
            string_data = strings.Replace(string_data, class.temp_data_raw[for_a][0], class.temp_data_raw[for_a][1], 1)
        }
    } else {
        for for_a := len(class.temp_data) - 1; for_a >= 0; for_a-- {
            string_data = strings.Replace(string_data, class.temp_data[for_a][0], class.temp_data[for_a][1], 1)
        }
    }

    return string_data
}

func (class *macromark) render_text() {
    string_data := class.render_data

    r := regexp2.MustCompile(`\[([^[(\]]+)\(((?:(?!\(|\)\])[\s\S])+)?\)\]`, 0)
    for {
        if m, _ := r.FindStringMatch(string_data); m != nil {
            gps := m.Groups()
            m_string := m.String()

            macro_name := gps[1].Captures[0].String()
            macro_data := ""
            if len(gps) > 2 {
                macro_data = gps[2].Captures[0].String()
            }

            switch macro_name {
                case "nowiki":
                    temp_name := class.func_temp_save(class.func_temp_restore(macro_data, true), m_string)
                    string_data = strings.Replace(string_data, m_string, temp_name, 1)
                case "h1":
                    temp_name := class.func_temp_save("<h1>" + macro_data + "</h1><back_br>", m_string)
                    string_data = strings.Replace(string_data, m_string, temp_name, 1)
                case "h2":
                    temp_name := class.func_temp_save("<h2>" + macro_data + "</h2><back_br>", m_string)
                    string_data = strings.Replace(string_data, m_string, temp_name, 1)
                case "h3":
                    temp_name := class.func_temp_save("<h3>" + macro_data + "</h3><back_br>", m_string)
                    string_data = strings.Replace(string_data, m_string, temp_name, 1)
                case "h4":
                    temp_name := class.func_temp_save("<h4>" + macro_data + "</h4><back_br>", m_string)
                    string_data = strings.Replace(string_data, m_string, temp_name, 1)
                case "h5":
                    temp_name := class.func_temp_save("<h5>" + macro_data + "</h5><back_br>", m_string)
                    string_data = strings.Replace(string_data, m_string, temp_name, 1)
                case "h6":
                    temp_name := class.func_temp_save("<h6>" + macro_data + "</h6><back_br>", m_string)
                    string_data = strings.Replace(string_data, m_string, temp_name, 1)
                case "ul":
                    temp_name := class.func_temp_save("<ul><back_br>" + macro_data + "</ul><back_br>", m_string)
                    string_data = strings.Replace(string_data, m_string, temp_name, 1)
                case "li":
                    temp_name := class.func_temp_save("<li>" + macro_data + "</li><back_br>", m_string)
                    string_data = strings.Replace(string_data, m_string, temp_name, 1)
                case "a":
                    a_data := class.func_temp_restore(macro_data, true)
                    a_data = strings.ReplaceAll(a_data, ",,", "<temp>")

                    part := strings.SplitN(a_data, ",", 2)

                    a_data_link := HTML_unescape(part[0])
                    a_data_view := a_data_link
                    if len(part) > 1 {
                        a_data_view = part[1]
                    }
                    
                    a_data_link = strings.ReplaceAll(a_data_link, "<temp>", ",")
                    a_data_view = strings.ReplaceAll(a_data_view, "<temp>", ",")

                    temp_name := class.func_temp_save("<a href=\"/w/" + Url_parser(a_data_link) + "\">" + a_data_view + "</a>", m_string)
                    string_data = strings.Replace(string_data, m_string, temp_name, 1)
                case "b":
                    temp_name := class.func_temp_save("<b>" + macro_data + "</b>", m_string)
                    string_data = strings.Replace(string_data, m_string, temp_name, 1)
                case "i":
                    temp_name := class.func_temp_save("<i>" + macro_data + "</i>", m_string)
                    string_data = strings.Replace(string_data, m_string, temp_name, 1)
                case "u":
                    temp_name := class.func_temp_save("<u>" + macro_data + "</u>", m_string)
                    string_data = strings.Replace(string_data, m_string, temp_name, 1)
                case "s":
                    temp_name := class.func_temp_save("<s>" + macro_data + "</s>", m_string)
                    string_data = strings.Replace(string_data, m_string, temp_name, 1)
                case "sup":
                    temp_name := class.func_temp_save("<sup>" + macro_data + "</sup>", m_string)
                    string_data = strings.Replace(string_data, m_string, temp_name, 1)
                case "sub":
                    temp_name := class.func_temp_save("<sub>" + macro_data + "</sub>", m_string)
                    string_data = strings.Replace(string_data, m_string, temp_name, 1)
                default:
                    temp_name := class.func_temp_save("", m_string)
                    string_data = strings.Replace(string_data, m_string, temp_name, 1)
            }
        } else {
            break
        }
    }

    class.render_data = string_data
}

func (class *macromark) render_last() {
    string_data := class.render_data

    string_data = class.func_temp_restore(string_data, false)

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
    class.render_data_js += "opennamu_do_toc();"
}

func (class macromark) main() map[string]interface{} {
    class.render_text()
    class.render_last()

    log.Default().Println(class.render_data)

    end_data := make(map[string]interface{})
    end_data["data"] = class.render_data
    end_data["js_data"] = class.render_data_js
    end_data["backlink"] = class.backlink
    end_data["link_count"] = class.link_count

    return end_data
}