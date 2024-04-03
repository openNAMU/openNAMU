package tool

import (
	"database/sql"
	"html"
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
	temp_data_count int

	backlink   [][]string
	link_count int
}

func Namumark_new(db *sql.DB, db_set map[string]string, data map[string]string) *namumark {
	return &namumark{
		db,
		db_set,
		data,

		html.EscapeString(data["data"]),
		"",

		map[string]string{},
		0,

		[][]string{},
		0,
	}
}

func (class *namumark) func_temp_save(data string) string {
	name := "<temp_save_" + strconv.Itoa(class.temp_data_count) + ">"

	class.temp_data[name] = data
	class.temp_data_count += 1

	return name
}

func (class namumark) func_temp_restore(data string) string {
	string_data := data

	for k, v := range class.temp_data {
		string_data = strings.Replace(string_data, k, v, 1)
	}

	return string_data
}

func (class *namumark) render_text() {
	string_data := class.render_data

	r := regexp2.MustCompile(`&#39;&#39;&#39;((?:(?!&#39;&#39;&#39;).)+)&#39;&#39;&#39;`, 0)
	for {
		if m, _ := r.FindStringMatch(string_data); m != nil {
			gps := m.Groups()

			temp_name := class.func_temp_save("<b>" + gps[1].Captures[0].String() + "</b>")
			string_data = strings.Replace(string_data, m.String(), temp_name, 1)
		} else {
			break
		}
	}

	class.render_data = string_data
}

func (class *namumark) render_last() {
	string_data := class.render_data

	string_data = class.func_temp_restore(string_data)

	class.render_data = string_data
}

func (class namumark) main() map[string]interface{} {
	class.render_text()
	class.render_last()

	end_data := make(map[string]interface{})
	end_data["data"] = class.render_data
	end_data["js_data"] = class.render_data_js
	end_data["backlink"] = class.backlink
	end_data["link_count"] = class.link_count

	return end_data
}
