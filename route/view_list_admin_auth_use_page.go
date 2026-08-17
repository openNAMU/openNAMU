package route

import (
	"strings"

	"opennamu/route/tool"
)

func View_list_admin_auth_use_page(config tool.Config, page string, search string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if search == "" || search == "normal" {
		search = ""
	}
	page_num := list_extra_page_number(page)
	offset := (page_num - 1) * 50
	query := "select who, what, time from re_admin order by time desc limit ?, 50"
	args := []any{offset}
	if search != "" {
		query = "select who, what, time from re_admin where what like ? order by time desc limit ?, 50"
		args = []any{search + "%", offset}
	}
	rows := tool.Query_DB(db, query, args...)
	body := strings.Builder{}
	count := 0
	for rows.Next() {
		who, what, date := "", "", ""
		if rows.Scan(&who, &what, &date) == nil {
			body.WriteString(tool.Get_list_ui(tool.IP_parser(db, who, config.IP), tool.HTML_escape(what), tool.HTML_escape(date), ""))
			count++
		}
	}
	rows.Close()
	page_url := "/list/admin/auth_use_page/{}"
	if search != "" {
		page_url += "/" + tool.Url_parser(search)
	}
	body.WriteString(tool.Get_page_control(db, page_num, count, 50, page_url))
	return list_extra_page(db, config, tool.Get_language(db, "auth_use", true), `<form method="post"><input name="search" value="`+tool.HTML_escape(search)+`"><button type="submit">`+tool.Get_language(db, "search", true)+`</button></form><hr class="main_hr">`+body.String())
}
