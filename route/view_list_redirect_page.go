package route

import (
	"strings"

	"opennamu/route/tool"
)

func View_list_redirect_page(config tool.Config, page string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	page_num := list_extra_page_number(page)
	offset := (page_num - 1) * 50
	rows := tool.Get_redirect_problem_rows(db, offset)
	body := strings.Builder{}
	count := 0
	for rows.Next() {
		source, target, next_target := "", "", ""
		if rows.Scan(&source, &target, &next_target) != nil {
			continue
		}

		left := `<a href="/w/` + tool.Url_parser(source) + `">` + tool.HTML_escape(source) + `</a>`
		right := `→ <a href="/w/` + tool.Url_parser(target) + `">` + tool.HTML_escape(target) + `</a>`
		if next_target != "" && target != source {
			right += ` → <a href="/w/` + tool.Url_parser(next_target) + `">` + tool.HTML_escape(next_target) + `</a>`
		}
		body.WriteString(tool.Get_list_ui(left, right, "", ""))
		count++
	}
	rows.Close()

	body.WriteString(tool.Get_page_control(db, page_num, count, 50, "/list/document/redirect/{}"))
	return list_extra_page(db, config, tool.Get_language(db, "redirect_problem_list", true), body.String())
}
