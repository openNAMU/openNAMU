package route

import (
	"opennamu/route/tool"
)

func View_wiki_access(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	body := `<form method="post"><div><label for="access_password">` + tool.Get_language(db, "password", true) + `</label> <input id="access_password" type="password" name="password"></div><hr class="main_hr"><div><button type="submit">` + tool.Get_language(db, "ok", true) + `</button></div></form>`
	return tool.Get_template(db, config, tool.Get_language(db, "error_password_require_for_wiki_access", true), body, []any{}, [][]any{}, map[string]string{})
}

func Check_wiki_access(password string) bool {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	saved_password := tool.Get_setting_value(db, "wiki_access_password", "", "")
	return saved_password != "" && saved_password == password
}
