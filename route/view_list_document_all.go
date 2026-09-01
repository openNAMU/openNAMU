package route

import (
	"strconv"
	"strings"

	"opennamu/route/tool"
)

func View_list_document_all(config tool.Config, page string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	page_num := list_extra_page_number(page)
	offset := (page_num - 1) * 50
	rows := tool.Get_data_rows(db, offset)
	body := strings.Builder{}
	count := 0
	for rows.Next() {
		title := ""
		if rows.Scan(&title) != nil {
			continue
		}
		body.WriteString(`<li>` + strconv.Itoa(offset+count+1) + `. <a href="/w/` + tool.Url_parser(title) + `">` + tool.HTML_escape(title) + `</a></li>`)
		count++
	}
	rows.Close()

	data := `<ul>` + body.String() + `</ul>`
	if page_num == 1 {
		all_title := tool.Get_setting_value(db, "count_all_title", "", "")
		if all_title != "" {
			total := tool.Str_to_int(all_title)
			data += `<ul><li>` + tool.Get_language(db, "all", true) + ` : ` + strconv.Itoa(total) + `</li></ul>`
			if total < 30000 {
				counts := []int{}
				for _, prefix := range []string{"category:", "user:", "file:"} {
					counts = append(counts, tool.Str_to_int(tool.Get_data_prefix_count(db, prefix)))
				}
				other_count := total - counts[0] - counts[1] - counts[2]
				data += `<ul><li>` + tool.Get_language(db, "category", true) + ` : ` + strconv.Itoa(counts[0]) + `</li><li>` + tool.Get_language(db, "user_document", true) + ` : ` + strconv.Itoa(counts[1]) + `</li><li>` + tool.Get_language(db, "file", true) + ` : ` + strconv.Itoa(counts[2]) + `</li><li>` + tool.Get_language(db, "other", true) + ` : ` + strconv.Itoa(other_count) + `</li></ul>`
			}
		}
	}

	data += tool.Get_page_control(db, page_num, count, 50, "/list/document/all/{}")
	return list_extra_page(db, config, tool.Get_language(db, "all_document_list", true), data)
}

func list_extra_page_number(value string) int {
	page := tool.Str_to_int(value)
	if page < 1 {
		return 1
	}
	return page
}
