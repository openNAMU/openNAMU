package route

import (
	"strconv"
	"strings"

	"opennamu/route/tool"
)

func View_user_multiple(config tool.Config, page string, sort string, search string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "check", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	page_num := list_extra_page_number(page)
	if sort != "count" {
		sort = "recent"
	}

	data_list := tool.Get_ua_multiple_rows(db, search, sort, (page_num-1)*50)
	body := strings.Builder{}
	recent_selected := ""
	count_selected := ""
	if sort == "count" {
		count_selected = ` selected="selected"`
	} else {
		recent_selected = ` selected="selected"`
	}

	body.WriteString(`<form method="post" action="/list/user/multiple"><input name="search" value="` + tool.HTML_escape(search) + `" placeholder="` + tool.Get_language(db, "name", true) + ` / ` + tool.Get_language(db, "ip", true) + `"><select name="sort"><option value="recent"` + recent_selected + `>` + tool.Get_language(db, "recent", true) + `</option><option value="count"` + count_selected + `>` + tool.Get_language(db, "account_count", true) + `</option></select><button type="submit">` + tool.Get_language(db, "search", true) + `</button></form><hr class="main_hr">`)
	body.WriteString(`<p>` + tool.Get_language(db, "multiple_account_note", true) + `</p>`)
	if tool.Get_setting_value(db, "ua_get", "", "") != "" {
		body.WriteString(`<p>` + tool.Get_language(db, "ua_collection_off_notice", true) + `</p>`)
	}

	count := 0
	for _, item := range data_list {
		names := strings.Builder{}
		for index, name := range item.Names {
			if index > 0 {
				names.WriteString(", ")
			}
			names.WriteString(`<a href="/user/` + tool.Url_parser(name) + `">` + tool.HTML_escape(name) + `</a>`)
		}

		left := `<a href="/list/user/check/` + tool.Url_parser(item.IP) + `">` + tool.HTML_escape(item.IP) + `</a> (` + strconv.Itoa(item.Count) + `)`
		right := tool.Get_language(db, "account_count", true) + ` : ` + strconv.Itoa(item.Count) + ` | ` + tool.Get_language(db, "last_login", true) + ` : ` + tool.HTML_escape(item.Date) + `<br>` + names.String()
		body.WriteString(tool.Get_list_ui(left, right, "", ""))
		count++
	}

	if count == 0 {
		body.WriteString(`<p>` + tool.Get_language(db, "data_missing", true) + `</p>`)
	}

	page_url := "/list/user/multiple/" + sort + "/{}"
	if search != "" {
		page_url += "/" + tool.Url_parser(tool.Base64_encode(search))
	}
	body.WriteString(tool.Get_page_control(db, page_num, count, 50, page_url))

	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "multiple_account", true),
		body.String(),
		[]any{},
		[][]any{{"manager", tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}
