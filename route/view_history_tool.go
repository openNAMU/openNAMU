package route

import (
	"strconv"

	"opennamu/route/tool"
)

func history_tool_revision(value string) (string, int, bool) {
	revision, err := strconv.Atoi(value)
	if err != nil {
		return "", 0, false
	}
	return strconv.Itoa(revision), revision, true
}

func View_history_tool(config tool.Config, doc_name string, rev string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "history_view", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	revision, revision_number, ok := history_tool_revision(rev)
	if !ok || doc_name == "" {
		return tool.Get_error_page(db, config, "error")
	}

	doc_name_url := tool.Url_parser(doc_name)
	data := `<h2>` + tool.Get_language(db, "tool", true) + `</h2><ul>`
	data += `<li><a href="/render/` + revision + `/` + doc_name_url + `">` + tool.Get_language(db, "view", true) + `</a></li>`
	data += `<li><a href="/raw_rev/` + revision + `/` + doc_name_url + `">` + tool.Get_language(db, "raw", true) + `</a></li>`
	data += `<li><a href="/revert/` + revision + `/` + doc_name_url + `">` + tool.Get_language(db, "revert", true) + ` (r` + revision + `)</a></li>`

	if revision_number > 1 {
		previous := strconv.Itoa(revision_number - 1)
		data += `<li><a href="/revert/` + previous + `/` + doc_name_url + `">` + tool.Get_language(db, "revert", true) + ` (r` + previous + `)</a></li>`
		data += `<li><a href="/diff/` + previous + `/` + revision + `/` + doc_name_url + `">` + tool.Get_language(db, "compare", true) + `</a></li>`
	}

	data += `<li><a href="/history/` + doc_name_url + `">` + tool.Get_language(db, "history", true) + `</a></li></ul>`

	if tool.Check_permission(db, "hidel", config.IP) {
		hide := ""
		tool.QueryRow_DB(db, "select hide from history where title = ? and id = ?", []any{&hide}, doc_name, revision)
		hide_name := "hide"
		if hide == "O" {
			hide_name = "hide_release"
		}
		data += `<h3>` + tool.Get_language(db, "admin", true) + `</h3><ul><li><a href="/history_hidden/` + revision + `/` + doc_name_url + `">` + tool.Get_language(db, hide_name, true) + `</a></li></ul>`
	}

	if tool.Check_permission(db, "history_manage", config.IP) {
		data += `<h3>` + tool.Get_language(db, "owner", true) + `</h3><ul>`
		data += `<li><a href="/history_delete/` + revision + `/` + doc_name_url + `">` + tool.Get_language(db, "history_delete", true) + `</a></li>`
		data += `<li><a href="/history_send/` + revision + `/` + doc_name_url + `">` + tool.Get_language(db, "send_edit", true) + `</a></li></ul>`
	}

	return tool.Get_template(
		db,
		config,
		doc_name,
		data,
		[]any{`(r` + revision + `)`},
		[][]any{{`history/` + tool.Url_parser(doc_name), tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}
