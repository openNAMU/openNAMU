package route

import (
	"database/sql"
	"strings"

	"opennamu/route/tool"
)

func list_extra_page_number(value string) int {
	page := tool.Str_to_int(value)
	if page < 1 {
		return 1
	}
	return page
}

func View_list_document_acl(config tool.Config, page string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "acl_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	page_num := list_extra_page_number(page)
	offset := (page_num - 1) * 50
	rows := tool.Query_DB(db, "select distinct title, data, type from acl where data != '' and title not like 'user:%' order by title desc limit ?, 50", offset)
	body := strings.Builder{}
	count := 0
	for rows.Next() {
		title, value, field := "", "", ""
		if rows.Scan(&title, &value, &field) != nil {
			continue
		}
		last_change := ""
		tool.QueryRow_DB(db, "select time from re_admin where what like ? order by time desc limit 1", []any{&last_change}, "acl ("+title+")%")
		why := ""
		tool.QueryRow_DB(db, "select data from acl where title = ? and type = 'why' limit 1", []any{&why}, title)
		body.WriteString(`<li>` + tool.HTML_escape(last_change) + ` | <a href="/acl/` + tool.Url_parser(title) + `">` + tool.HTML_escape(title) + `</a> | ` + tool.HTML_escape(value) + ` (` + tool.HTML_escape(field) + `)`)
		if why != "" {
			body.WriteString(` | ` + tool.HTML_escape(why))
		}
		body.WriteString(`</li>`)
		count++
	}
	rows.Close()
	body.WriteString(tool.Get_page_control(db, page_num, count, 50, "/list/document/acl/{}"))
	return list_extra_page(db, config, tool.Get_language(db, "acl_document_list", true), `<ul>`+body.String()+`</ul>`)
}

func View_list_need_page(config tool.Config, page string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	page_num := list_extra_page_number(page)
	offset := (page_num - 1) * 50
	rows := tool.Query_DB(db, "select title, link from back where type = 'no' order by title limit ?, 50", offset)
	body := strings.Builder{}
	count := 0
	for rows.Next() {
		name, link := "", ""
		if rows.Scan(&name, &link) == nil {
			body.WriteString(tool.Get_list_ui(`<a href="/w/`+tool.Url_parser(name)+`">`+tool.HTML_escape(name)+`</a>`, tool.HTML_escape(link), "", ""))
			count++
		}
	}
	rows.Close()
	body.WriteString(tool.Get_page_control(db, page_num, count, 50, "/list/document/need/{}"))
	return list_extra_page(db, config, tool.Get_language(db, "need_document", true), body.String())
}

func View_list_no_link_page(config tool.Config, page string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	page_num := list_extra_page_number(page)
	offset := (page_num - 1) * 50
	rows := tool.Query_DB(db, "select doc_name, set_data from data_set where set_name = 'link_count' and set_data = '0' order by doc_name limit ?, 50", offset)
	body := strings.Builder{}
	count := 0
	for rows.Next() {
		name, value := "", ""
		if rows.Scan(&name, &value) == nil {
			body.WriteString(tool.Get_list_ui(`<a href="/w/`+tool.Url_parser(name)+`">`+tool.HTML_escape(name)+`</a>`, tool.HTML_escape(value), "", ""))
			count++
		}
	}
	rows.Close()
	body.WriteString(tool.Get_page_control(db, page_num, count, 50, "/list/document/no_link/{}"))
	return list_extra_page(db, config, tool.Get_language(db, "no_link_document", true), body.String())
}

func View_list_file_page(config tool.Config, page string, image_only bool) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	page_num := list_extra_page_number(page)
	offset := (page_num - 1) * 50
	rows := tool.Query_DB(db, "select title, type from data where title like 'file:%' order by title limit ?, 50", offset)
	body := strings.Builder{}
	count := 0
	for rows.Next() {
		name, kind := "", ""
		if rows.Scan(&name, &kind) != nil {
			continue
		}
		if image_only && !strings.HasPrefix(kind, "image") && !strings.Contains(name, ".png") && !strings.Contains(name, ".jpg") && !strings.Contains(name, ".jpeg") && !strings.Contains(name, ".gif") && !strings.Contains(name, ".webp") {
			continue
		}
		body.WriteString(`<li><a href="/w/` + tool.Url_parser(name) + `">` + tool.HTML_escape(name) + `</a></li>`)
		count++
	}
	rows.Close()
	body.WriteString(tool.Get_page_control(db, page_num, count, 50, func() string {
		if image_only {
			return "/list/image/{}"
		}
		return "/list/file/{}"
	}()))
	title := "file_list"
	if image_only {
		title = "image_list"
	}
	return list_extra_page(db, config, tool.Get_language(db, title, true), `<ul>`+body.String()+`</ul>`)
}

func View_list_user_page(config tool.Config, page string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	page_num := list_extra_page_number(page)
	offset := (page_num - 1) * 50
	rows := tool.Query_DB(db, "select id, data from user_set where name = 'date' order by data desc limit ?, 50", offset)
	body := strings.Builder{}
	count := 0
	for rows.Next() {
		name, date := "", ""
		if rows.Scan(&name, &date) == nil {
			body.WriteString(tool.Get_list_ui(`<a href="/user/`+tool.Url_parser(name)+`">`+tool.IP_parser(db, name, config.IP)+`</a>`, tool.HTML_escape(date), "", ""))
			count++
		}
	}
	rows.Close()
	body.WriteString(tool.Get_page_control(db, page_num, count, 50, "/list/user/{}"))
	return list_extra_page(db, config, tool.Get_language(db, "user_list", true), body.String())
}

func View_list_admin_page(config tool.Config, page string, auth_use bool, search string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if auth_use && !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	page_num := list_extra_page_number(page)
	offset := (page_num - 1) * 50
	body := strings.Builder{}
	count := 0
	if auth_use {
		if search == "" || search == "normal" {
			search = ""
		}
		var rows *sql.Rows
		if search == "" {
			rows = tool.Query_DB(db, "select who, what, time from re_admin order by time desc limit ?, 50", offset)
		} else {
			rows = tool.Query_DB(db, "select who, what, time from re_admin where what like ? order by time desc limit ?, 50", search+"%", offset)
		}
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
	rows := tool.Query_DB(db, "select id, data from user_set where name = 'acl' and data != 'user' order by id limit ?, 50", offset)
	for rows.Next() {
		name, auth := "", ""
		if rows.Scan(&name, &auth) == nil {
			body.WriteString(tool.Get_list_ui(`<a href="/user/`+tool.Url_parser(name)+`">`+tool.IP_parser(db, name, config.IP)+`</a>`, tool.HTML_escape(auth), "", ""))
			count++
		}
	}
	rows.Close()
	body.WriteString(tool.Get_page_control(db, page_num, count, 50, "/list/admin/{}"))
	return list_extra_page(db, config, tool.Get_language(db, "admin_list", true), body.String())
}
