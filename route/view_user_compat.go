package route

import (
	"database/sql"
	"net/url"
	"opennamu/route/tool"
	"regexp"
	"strconv"
	"strings"
)

func user_field_values(values url.Values, field string) url.Values {
	if values == nil {
		return nil
	}
	copy_values := url.Values{}
	for key, list := range values {
		copy_values[key] = append([]string{}, list...)
	}
	if copy_values.Get("data") == "" {
		if value := copy_values.Get(field); value != "" {
			copy_values.Set("data", value)
		}
	}
	if field == "user_name" && copy_values.Get("data") == "" {
		copy_values.Set("data", copy_values.Get("new_user_name"))
	}
	return copy_values
}

func user_skin_choice(set_list map[string][][]string, field string, value string) bool {
	for _, choice := range set_list[field] {
		if len(choice) > 0 && choice[0] == value {
			return true
		}
	}
	return false
}

func user_skin_main_render_simple_set(db *sql.DB, data string) string {
	heading_regex := regexp.MustCompile("(?s)<h([1-6])>([^<>]+)</h[1-6]>")
	matches := heading_regex.FindAllStringSubmatch(data, -1)
	if len(matches) == 0 {
		return data
	}

	heading_stack := make([]int, 6)
	toc_data := strings.Builder{}
	toc_data.WriteString(`<div class="opennamu_TOC" id="toc"><span class="opennamu_TOC_title">` + tool.Get_language(db, "toc", true) + `</span><br>`)
	for _, match := range matches {
		heading_level, err := strconv.Atoi(match[1])
		if err != nil || heading_level < 1 || heading_level > 6 {
			continue
		}
		heading_stack[heading_level-1]++
		for i := heading_level; i < len(heading_stack); i++ {
			heading_stack[i] = 0
		}

		heading_number := strings.Builder{}
		for _, count := range heading_stack {
			if count == 0 {
				continue
			}
			if heading_number.Len() > 0 {
				heading_number.WriteString(".")
			}
			heading_number.WriteString(strconv.Itoa(count))
		}
		number := heading_number.String()
		indent := strings.Count(number, ".")
		toc_data.WriteString(`<br><span class="opennamu_TOC_list">` + strings.Repeat(`<span style="margin-left: 10px;"></span>`, indent) + `<a href="#s-` + number + `">` + number + `.</a> ` + match[2] + `</span>`)
		heading := `<h` + match[1] + ` id="s-` + number + `"><a href="#toc">` + number + `.</a> ` + match[2] + `</h` + match[1] + `>`
		data = strings.Replace(data, match[0], heading, 1)
	}
	toc_data.WriteString(`</div>`)
	return toc_data.String() + data
}
