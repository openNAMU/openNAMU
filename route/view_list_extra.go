package route

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"opennamu/route/tool"
)

func list_extra_page(db *sql.DB, config tool.Config, title string, body string) string {
	return tool.Get_template(db, config, title, body, []any{}, [][]any{{"other", tool.Get_language(db, "return", true)}}, map[string]string{})
}

func list_extra_query(db *sql.DB, config tool.Config, title string, query string, args ...any) string {
	rows := tool.Query_DB(db, query, args...)
	body := strings.Builder{}
	for rows.Next() {
		name, value := "", ""
		if rows.Scan(&name, &value) != nil {
			continue
		}
		body.WriteString(tool.Get_list_ui(`<a href="/w/`+tool.Url_parser(name)+`">`+tool.HTML_escape(name)+`</a>`, tool.HTML_escape(value), "", ""))
	}
	rows.Close()
	return list_extra_page(db, config, title, body.String())
}

func View_list_document_all(config tool.Config, page string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	page_num := list_extra_page_number(page)
	offset := (page_num - 1) * 50
	rows := tool.Query_DB(db, "select title from data order by title asc limit ?, 50", offset)
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
		all_title := ""
		tool.QueryRow_DB(db, "select data from other where name = 'count_all_title'", []any{&all_title})
		if all_title != "" {
			total := tool.Str_to_int(all_title)
			data += `<ul><li>` + tool.Get_language(db, "all", true) + ` : ` + strconv.Itoa(total) + `</li></ul>`
			if total < 30000 {
				counts := []int{}
				for _, prefix := range []string{"category:", "user:", "file:"} {
					value := "0"
					tool.QueryRow_DB(db, "select count(*) from data where title like ?", []any{&value}, prefix+"%")
					counts = append(counts, tool.Str_to_int(value))
				}
				other_count := total - counts[0] - counts[1] - counts[2]
				data += `<ul><li>` + tool.Get_language(db, "category", true) + ` : ` + strconv.Itoa(counts[0]) + `</li><li>` + tool.Get_language(db, "user_document", true) + ` : ` + strconv.Itoa(counts[1]) + `</li><li>` + tool.Get_language(db, "file", true) + ` : ` + strconv.Itoa(counts[2]) + `</li><li>` + tool.Get_language(db, "other", true) + ` : ` + strconv.Itoa(other_count) + `</li></ul>`
			}
		}
	}

	data += tool.Get_page_control(db, page_num, count, 50, "/list/document/all/{}")
	return list_extra_page(db, config, tool.Get_language(db, "all_document_list", true), data)
}

func View_list_need(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	return list_extra_query(db, config, tool.Get_language(db, "need_document", true), "select title, link from back where type = 'no' order by title")
}

func View_list_no_link(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	return list_extra_query(db, config, tool.Get_language(db, "no_link_document", true), "select doc_name, set_data from data_set where set_name = 'link_count' and set_data = '0' order by doc_name")
}

func View_list_file(config tool.Config, image_only bool) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	body := strings.Builder{}
	rows := tool.Query_DB(db, "select title, type from data where title like 'file:%' order by title")
	for rows.Next() {
		name, kind := "", ""
		if rows.Scan(&name, &kind) != nil {
			continue
		}
		if image_only && !strings.HasPrefix(kind, "image") && !strings.Contains(name, ".png") && !strings.Contains(name, ".jpg") && !strings.Contains(name, ".jpeg") && !strings.Contains(name, ".gif") && !strings.Contains(name, ".webp") {
			continue
		}
		body.WriteString(`<li><a href="/w/` + tool.Url_parser(name) + `">` + tool.HTML_escape(name) + `</a></li>`)
	}
	rows.Close()
	return list_extra_page(db, config, tool.Get_language(db, func() string {
		if image_only {
			return "image_list"
		}
		return "file_list"
	}(), true), `<ul>`+body.String()+`</ul>`)
}

func View_list_user(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	body := strings.Builder{}
	rows := tool.Query_DB(db, "select id, data from user_set where name = 'date' order by data desc")
	for rows.Next() {
		name, date := "", ""
		if rows.Scan(&name, &date) == nil {
			body.WriteString(tool.Get_list_ui(`<a href="/user/`+tool.Url_parser(name)+`">`+tool.IP_parser(db, name, config.IP)+`</a>`, tool.HTML_escape(date), "", ""))
		}
	}
	rows.Close()
	return list_extra_page(db, config, tool.Get_language(db, "user_list", true), body.String())
}

func View_list_admin(config tool.Config, auth_use bool) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if auth_use && !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	body := strings.Builder{}
	if auth_use {
		rows := tool.Query_DB(db, "select who, what from re_admin order by time desc limit 200")
		for rows.Next() {
			who, what := "", ""
			if rows.Scan(&who, &what) == nil {
				body.WriteString(tool.Get_list_ui(tool.IP_parser(db, who, config.IP), tool.HTML_escape(what), "", ""))
			}
		}
		rows.Close()
		return list_extra_page(db, config, tool.Get_language(db, "auth_use", true), body.String())
	}
	rows := tool.Query_DB(db, "select id, data from user_set where name = 'acl' and data != 'user' order by id")
	for rows.Next() {
		name, auth := "", ""
		if rows.Scan(&name, &auth) == nil {
			body.WriteString(tool.Get_list_ui(`<a href="/user/`+tool.Url_parser(name)+`">`+tool.IP_parser(db, name, config.IP)+`</a>`, tool.HTML_escape(auth), "", ""))
		}
	}
	rows.Close()
	return list_extra_page(db, config, tool.Get_language(db, "admin_list", true), body.String())
}

func View_record_simple(config tool.Config, user_name string, record_type string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if user_name != "" && user_name != config.IP && !tool.Check_acl(db, "", "", "hidel_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if user_name == "" {
		user_name = config.IP
	}
	body := strings.Builder{}
	switch record_type {
	case "topic":
		rows := tool.Query_DB(db, "select code, data, date from topic where ip = ? order by date desc limit 100", user_name)
		for rows.Next() {
			code, data, date := "", "", ""
			if rows.Scan(&code, &data, &date) == nil {
				body.WriteString(tool.Get_list_ui(`<a href="/thread/`+tool.Url_parser(code)+`">`+tool.HTML_escape(code)+`</a>`, tool.HTML_escape(date), tool.HTML_escape(data), ""))
			}
		}
		rows.Close()
	default:
		rows := tool.Query_DB(db, "select title, date, send from history where ip = ? order by id + 0 desc limit 100", user_name)
		for rows.Next() {
			title, date, send := "", "", ""
			if rows.Scan(&title, &date, &send) == nil {
				body.WriteString(tool.Get_list_ui(`<a href="/w/`+tool.Url_parser(title)+`">`+tool.HTML_escape(title)+`</a>`, tool.HTML_escape(date), tool.HTML_escape(send), ""))
			}
		}
		rows.Close()
	}
	menu := [][]any{{"other", tool.Get_language(db, "return", true)}}
	if tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		menu = append(menu, []any{"record/reset/" + tool.Url_parser(user_name), tool.Get_language(db, "record_reset", true)})
	}
	return tool.Get_template(db, config, tool.Get_language(db, record_type+"_record", true), body.String(), []any{}, menu, map[string]string{})
}
func View_record_count(config tool.Config, user_name string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if user_name != "" && user_name != config.IP && !tool.Check_acl(db, "", "", "hidel_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if user_name == "" {
		user_name = config.IP
	}
	history, topic, bbs := "0", "0", "0"
	tool.QueryRow_DB(db, "select count(*) from history where ip = ?", []any{&history}, user_name)
	tool.QueryRow_DB(db, "select count(*) from topic where ip = ?", []any{&topic}, user_name)
	tool.QueryRow_DB(db, "select count(*) from bbs_data where set_name = 'comment_user_id' and set_data = ?", []any{&bbs}, user_name)

	count_day := func(day string) (int, int) {
		rows := tool.Query_DB(db, "select leng from history where date like ? and ip = ?", day+"%", user_name)
		defer rows.Close()
		count := 0
		length := 0
		for rows.Next() {
			value := ""
			if rows.Scan(&value) != nil {
				continue
			}
			value = strings.TrimLeft(value, "+-")
			length += tool.Str_to_int(value)
			count++
		}
		return count, length
	}

	today_count, today_length := count_day(tool.Get_date())
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	yesterday_count, yesterday_length := count_day(yesterday)
	user_url := tool.Url_parser(user_name)
	body := `<ul><li><a href="/record/` + user_url + `">` + tool.Get_language(db, "edit_record", true) + `</a> : ` + history + `</li><li><a href="/record/topic/` + user_url + `">` + tool.Get_language(db, "discussion_record", true) + `</a> : ` + topic + `</li><li>bbs : ` + bbs + `</li><hr><li>(` + tool.Get_language(db, "beta", true) + `) TODAY : ` + strconv.Itoa(today_count) + `</li><li>(` + tool.Get_language(db, "beta", true) + `) TODAY LEN : ` + strconv.Itoa(today_length) + `</li><li>(` + tool.Get_language(db, "beta", true) + `) TODAY DIFF : ` + strconv.Itoa(today_length-yesterday_length) + `</li><hr><li>(` + tool.Get_language(db, "beta", true) + `) YESTERDAY : ` + strconv.Itoa(yesterday_count) + `</li><li>(` + tool.Get_language(db, "beta", true) + `) YESTERDAY LEN : ` + strconv.Itoa(yesterday_length) + `</li></ul>`
	return list_extra_page(db, config, tool.Get_language(db, "count", true), body)
}

func View_list_user_check_submit(config tool.Config, user_name string) string {
	return View_user_check(config, user_name, "check", "1", "")
}

func list_page_path(base string, page int) string {
	return base + "/" + strconv.Itoa(page)
}
