package markup

import (
	"database/sql"
	"html"
	"log"
	"regexp"
	"strconv"
	"strings"

	"opennamu/route/tool"

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

    toc_struct []int
    toc_result map[string]string

    backlink [][]string
    link_count int

    result_markup string
}

func Macromark_new(db *sql.DB, data map[string]string, result_markup string) *macromark {
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

        []int{0, 0, 0, 0, 0, 0},
        map[string]string{},

        [][]string{},
        0,

        result_markup,
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

type macro_data struct {
    function map[string]macro_transform_func    
}

type macro_transform_func func(class *macromark, macro_name string, macro_data string, m_string string)

var heading_markdown = func(class *macromark, macro_name string, macro_data string, m_string string) {
    temp_name := class.func_temp_save("\n## " + macro_data + "\n", m_string)
    class.render_data = strings.Replace(class.render_data, m_string, temp_name, 1)
}

var heading_html = func(class *macromark, macro_name string, macro_data string, m_string string) {
    heading_n := macro_name[1:]
    heading_n_int := tool.Str_to_int(heading_n)

    class.toc_struct[heading_n_int - 1] += 1
    for for_a := heading_n_int; for_a < 6; for_a++ {
        class.toc_struct[for_a] = 0
    }

    toc_string := ""
    temp_string := ""
    for for_a := 0; for_a < 6; for_a++ {
        heading_str := strconv.Itoa(class.toc_struct[for_a])
        
        if class.toc_struct[for_a] == 0 {
            temp_string += "0."
            continue
        } 
        
        if toc_string == "" {
            toc_string += heading_str + "."
            temp_string = ""
        } else {
            toc_string += temp_string + heading_str + "."
            temp_string = ""
        }
    }

    if toc_string != "" {
        toc_string = toc_string[:len(toc_string) - 1]
    }

    class.toc_result[toc_string] = macro_data

    temp_name := class.func_temp_save("<" + macro_name + "><a href=\"#toc\">" + toc_string + ". </a>" + macro_data + "</" + macro_name + "><back_br>", m_string)
    class.render_data = strings.Replace(class.render_data, m_string, temp_name, 1)
}
var simple_html = func(class *macromark, macro_name string, macro_data string, m_string string) {
    temp_name := class.func_temp_save("<" + macro_name + ">" + macro_data + "</" + macro_name + ">", m_string)
    class.render_data = strings.Replace(class.render_data, m_string, temp_name, 1)
}

var macro_transform_map = map[string]macro_data{
    "namumark" : {
    },
    "markdown" : {
        function : map[string]macro_transform_func{
            "h1" : heading_markdown,
            "h2" : heading_markdown,
            "h3" : heading_markdown,            
            "h4" : heading_markdown,
            "h5" : heading_markdown,
            "h6" : heading_markdown,
        },
    },
    "html" : {
        function : map[string]macro_transform_func{
            "nowiki" : func(class *macromark, macro_name string, macro_data string, m_string string) {
                temp_name := class.func_temp_save(class.func_temp_restore(macro_data, true), m_string)
                class.render_data = strings.Replace(class.render_data, m_string, temp_name, 1)
            },
            "h1" : heading_html,
            "h2" : heading_html,
            "h3" : heading_html,
            "h4" : heading_html,
            "h5" : heading_html,
            "h6" : heading_html,
            "ul" : func(class *macromark, macro_name string, macro_data string, m_string string) {
                temp_name := class.func_temp_save("<ul><back_br>" + macro_data + "</ul><back_br>", m_string)
                class.render_data = strings.Replace(class.render_data, m_string, temp_name, 1)
            },
            "li" : func(class *macromark, macro_name string, macro_data string, m_string string) {
                temp_name := class.func_temp_save("<li>" + macro_data + "</li><back_br>", m_string)
                class.render_data = strings.Replace(class.render_data, m_string, temp_name, 1)
            },
            "a" : func(class *macromark, macro_name string, macro_data string, m_string string) {
                a_data := class.func_temp_restore(macro_data, true)
                a_data = strings.ReplaceAll(a_data, ",,", "<temp>")

                part := strings.SplitN(a_data, ",", 3)

                a_data_link := tool.HTML_unescape(part[0])
                a_data_view := a_data_link
                if len(part) > 1 {
                    a_data_view = part[1]
                }

                a_data_hash := ""
                if len(part) > 2 {
                    a_data_hash = part[2]
                }
                
                a_data_link = strings.ReplaceAll(a_data_link, "<temp>", ",")
                a_data_view = strings.ReplaceAll(a_data_view, "<temp>", ",")
                a_data_hash = strings.ReplaceAll(a_data_hash, "<temp>", ",")

                temp := ""
                exist := tool.QueryRow_DB(
                    class.db,
                    "select title from data where title = ?",
                    []any{ &temp },
                    a_data_link,
                )

                exist_link := ""
                if !exist {
                    exist_link = "class=\"opennamu_not_exist_link\""
                }

                temp_name := class.func_temp_save("<a " + exist_link + " href=\"/w/" + tool.Url_parser(a_data_link) + a_data_hash + "\">" + a_data_view + "</a>", m_string)
                class.render_data = strings.Replace(class.render_data, m_string, temp_name, 1)
            },
            "ex" : func(class *macromark, macro_name string, macro_data string, m_string string) {
                a_data := class.func_temp_restore(macro_data, true)
                a_data = strings.ReplaceAll(a_data, ",,", "<temp>")

                part := strings.SplitN(a_data, ",", 2)

                a_data_link := tool.HTML_unescape(part[0])
                a_data_view := a_data_link
                if len(part) > 1 {
                    a_data_view = part[1]
                }

                a_data_link = strings.ReplaceAll(a_data_link, "<temp>", ",")
                a_data_view = strings.ReplaceAll(a_data_view, "<temp>", ",")
                
                temp_name := class.func_temp_save("<a class=\"opennamu_link_out\" href=\"" + tool.HTML_escape(a_data_link) + "\">" + a_data_view + "</a>", m_string)
                class.render_data = strings.Replace(class.render_data, m_string, temp_name, 1)
            },
            "b" : simple_html,
            "i" : simple_html,
            "u" : simple_html,
            "s" : simple_html,
            "sup" : simple_html,
            "sub" : simple_html,
            "toc" : func(class *macromark, macro_name string, macro_data string, m_string string) {
                temp_name := class.func_temp_save("<toc_data>", m_string)
                class.render_data = strings.Replace(class.render_data, m_string, temp_name, 1)
            },
        },
    },
}

func (class *macromark) render_text() {
    r := regexp2.MustCompile(`\[([^[(\]]+)\(((?:(?!\(|\)\])[\s\S])+)?\)\]`, 0)
    for {
        if m, _ := r.FindStringMatch(class.render_data); m != nil {
            gps := m.Groups()
            m_string := m.String()

            macro_name := gps[1].Captures[0].String()
            macro_data := ""
            if len(gps) > 2 && gps[2].Captures != nil && len(gps[2].Captures) > 0 {
                macro_data = gps[2].Captures[0].String()
            }

            reg, ok := macro_transform_map[class.result_markup]
            if !ok || reg.function == nil {
                class.render_data = strings.Replace(class.render_data, m_string, macro_data, 1)
                continue
            }

            fn, ok := reg.function[macro_name]
            if !ok {
                class.render_data = strings.Replace(class.render_data, m_string, macro_data, 1)
                continue
            }

            fn(class, macro_name, macro_data, m_string)
        } else {
            break
        }
    }
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

    toc_data := ""
    for k, v := range class.toc_result {
        toc_data += "<br>"
        
        count_comma := strings.Count(k, ".")
        toc_space := strings.Repeat("<span style=\"margin-left: 10px;\"></span>", count_comma)

        toc_data += "<span class=\"opennamu_TOC_list\">" + toc_space + "<a href=\"#s-" + k + "\">" + k + ". </a>" + v + "</span>"
    }

    if toc_data != "" {
        toc_data = "<div class=\"opennamu_TOC\" id=\"toc\"><span class=\"opennamu_TOC_title\">" + tool.Get_language(class.db, "toc", true) + "</span><br>" + toc_data + "</div>"
    }

    r = regexp.MustCompile(`<toc_data>`)
    string_data = r.ReplaceAllString(string_data, toc_data)

    string_data = strings.Replace(string_data, "\n", "<br>", -1)

    class.render_data = string_data
}

func (class macromark) main() map[string]any {
    class.render_text()
    class.render_last()

    log.Default().Println(class.toc_result)

    end_data := make(map[string]any)
    end_data["data"] = class.render_data
    end_data["js_data"] = class.render_data_js
    end_data["backlink"] = class.backlink
    end_data["link_count"] = class.link_count

    return end_data
}