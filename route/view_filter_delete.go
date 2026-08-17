package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_filter_delete(config tool.Config, kind string, name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	spec, ok := get_filter_spec(kind)
	if !ok {
		return tool.Get_error_page(db, config, "error")
	}
	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if values != nil {
		tool.Exec_DB(db, "delete from html_filter where html = ? and kind = ?", name, spec.db_kind)
		if kind == "inter_wiki" {
			tool.Exec_DB(db, "delete from html_filter where html = ? and kind = 'inter_wiki_sub'", name)
		}
		tool.Do_insert_auth_history(db, config.IP, "filter_delete ("+kind+")")
		return tool.Get_redirect("/filter/" + kind)
	}

	body := `<form method="post"><span>` + tool.Get_language(db, "delete_warning", true) + `</span><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "delete", true) + `</button></form>`
	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "delete", true),
		body,
		[]any{"(" + tool.HTML_escape(name) + ")"},
		[][]any{{"filter/" + kind, tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}
