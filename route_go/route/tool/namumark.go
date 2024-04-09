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

type namumark struct {
	db     *sql.DB
	db_set map[string]string
	data   map[string]string

	render_data    string
	render_data_js string

	temp_data       map[string]string
	temp_data_raw   map[string]string
	temp_data_count int

	backlink   [][]string
	link_count int
}

func Namumark_new(db *sql.DB, db_set map[string]string, data map[string]string) *namumark {
	data_string := data["data"]
	data_string = html.EscapeString(data_string)
	data_string = strings.Replace(data_string, "\r", "", -1)
	data_string = "\n" + data_string + "\n"

	return &namumark{
		db,
		db_set,
		data,

		data_string,
		"",

		map[string]string{},
		map[string]string{},
		0,

		[][]string{},
		0,
	}
}

func (class *namumark) func_temp_save(data string, data_raw string) string {
	name := "<temp_save_" + strconv.Itoa(class.temp_data_count) + ">"

	class.temp_data[name] = data
	class.temp_data_raw[name] = data_raw
	class.temp_data_count += 1

	return name
}

func (class namumark) func_temp_restore(data string, to_raw bool) string {
	string_data := data

	if to_raw {
		for k, v := range class.temp_data_raw {
			string_data = strings.Replace(string_data, k, v, 1)
		}
	} else {
		for k, v := range class.temp_data {
			string_data = strings.Replace(string_data, k, v, 1)
		}
	}

	return string_data
}

func (class *namumark) render_text() {
	string_data := class.render_data

	r := regexp2.MustCompile(`&#39;&#39;&#39;((?:(?!&#39;&#39;&#39;).)+)&#39;&#39;&#39;`, 0)
	for {
		if m, _ := r.FindStringMatch(string_data); m != nil {
			gps := m.Groups()
			m_string := m.String()

			temp_name := class.func_temp_save("<b>"+gps[1].Captures[0].String()+"</b>", m_string)
			string_data = strings.Replace(string_data, m_string, temp_name, 1)
		} else {
			break
		}
	}

	r = regexp2.MustCompile(`&#39;&#39;((?:(?!&#39;&#39;).)+)&#39;&#39;`, 0)
	for {
		if m, _ := r.FindStringMatch(string_data); m != nil {
			gps := m.Groups()
			m_string := m.String()

			temp_name := class.func_temp_save("<i>"+gps[1].Captures[0].String()+"</i>", m_string)
			string_data = strings.Replace(string_data, m_string, temp_name, 1)
		} else {
			break
		}
	}

	r = regexp2.MustCompile(`__((?:(?!__).)+)__`, 0)
	for {
		if m, _ := r.FindStringMatch(string_data); m != nil {
			gps := m.Groups()
			m_string := m.String()

			temp_name := class.func_temp_save("<u>"+gps[1].Captures[0].String()+"</u>", m_string)
			string_data = strings.Replace(string_data, m_string, temp_name, 1)
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

				temp_name := class.func_temp_save("<sup>"+gps[1].Captures[0].String()+"</sup>", m_string)
				string_data = strings.Replace(string_data, m_string, temp_name, 1)
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

				temp_name := class.func_temp_save("<sub>"+gps[1].Captures[0].String()+"</sub>", m_string)
				string_data = strings.Replace(string_data, m_string, temp_name, 1)
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

				temp_name := class.func_temp_save("<s>"+gps[1].Captures[0].String()+"</s>", m_string)
				string_data = strings.Replace(string_data, m_string, temp_name, 1)
			} else {
				break
			}
		}
	}

	class.render_data = string_data
}

func (class *namumark) render_last() {
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

func (class *namumark) render_heading() {
	string_data := class.render_data

	r := regexp.MustCompile(`\n(?:(={1,6})(#?) ?([^\n]+))\n`)
	string_data = r.ReplaceAllStringFunc(string_data, func(m string) string {
		match := r.FindStringSubmatch(m)

		r = regexp.MustCompile(` ?(#?={1,6}[^=]*)$`)
		heading_data := r.ReplaceAllString(match[3], "")

		heading_len := strconv.Itoa(len(match[1]))
		heading_render := "<h" + heading_len + ">" + heading_data + "</h" + heading_len + ">"

		temp_name := class.func_temp_save(heading_render, match[0])
		return temp_name
	})

	class.render_data = string_data
}

func (class namumark) main() map[string]interface{} {
	class.render_text()
	class.render_heading()
	class.render_last()

	log.Default().Println(class.render_data)

	end_data := make(map[string]interface{})
	end_data["data"] = class.render_data
	end_data["js_data"] = class.render_data_js
	end_data["backlink"] = class.backlink
	end_data["link_count"] = class.link_count

	return end_data
}
