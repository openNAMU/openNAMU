package route

import (
	"database/sql"
	"net/url"
	"strconv"
	"strings"

	"opennamu/route/tool"
)

func history_destructive_page(db *sql.DB, config tool.Config, title string, action string, return_path string) string {
	body := `<form method="post"><span>` + tool.Get_language(db, "delete_warning", true) + `</span><hr class="main_hr"><button type="submit">` + action + `</button></form>`
	return tool.Get_template(db, config, title, body, []any{}, [][]any{{return_path, tool.Get_language(db, "return", true)}}, map[string]string{})
}

func history_revision_value(value string) (string, bool) {
	revision := strings.TrimSpace(value)
	if revision == "" {
		return "", false
	}
	if _, err := strconv.Atoi(revision); err != nil {
		return "", false
	}
	return revision, true
}

func View_history_delete(config tool.Config, doc_name string, rev string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	revision, ok := history_revision_value(rev)
	if !ok || doc_name == "" {
		return tool.Get_error_page(db, config, "error")
	}
	if values != nil {
		tool.Exec_DB(db, "delete from history where id = ? and title = ?", revision, doc_name)
		tool.Do_insert_auth_history(db, config.IP, "history_delete ("+doc_name+" r"+revision+")")
		return tool.Get_redirect("/history/" + tool.Url_parser(doc_name))
	}
	return history_destructive_page(db, config, tool.Get_language(db, "history_delete", true)+" (r"+revision+")", tool.Get_language(db, "delete", true), "history/"+tool.Url_parser(doc_name))
}

func View_history_reset(config tool.Config, doc_name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if doc_name == "" {
		return tool.Get_error_page(db, config, "error")
	}
	if values != nil {
		tool.Exec_DB(db, "delete from history where title = ?", doc_name)
		tool.Do_insert_auth_history(db, config.IP, "history_reset ("+doc_name+")")
		return tool.Get_redirect("/history/" + tool.Url_parser(doc_name))
	}
	return history_destructive_page(db, config, tool.Get_language(db, "history_reset", true), tool.Get_language(db, "reset", true), "history/"+tool.Url_parser(doc_name))
}

func View_record_reset(config tool.Config, user_name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if user_name == "" {
		return tool.Get_error_page(db, config, "error")
	}
	if values != nil {
		tool.Exec_DB(db, "delete from history where ip = ?", user_name)
		tool.Do_insert_auth_history(db, config.IP, "record_reset ("+user_name+")")
		return tool.Get_redirect("/record/" + tool.Url_parser(user_name))
	}
	return history_destructive_page(db, config, tool.Get_language(db, "record_reset", true), tool.Get_language(db, "reset", true), "record/"+tool.Url_parser(user_name))
}
