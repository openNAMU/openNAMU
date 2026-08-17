package route

import (
	"opennamu/route/tool"
)

func View_wiki_access(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	body := "<form method='post'><input type='password' name='password'><button type='submit'>submit</button></form>"
	return tool.Get_template(db, config, tool.Get_language(db, "error_password_require_for_wiki_access", true), body, []any{}, [][]any{}, map[string]string{})
}

func Check_wiki_access(password string) bool {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	saved_password := ""
	tool.QueryRow_DB(db, "select data from other where name = 'wiki_access_password'", []any{&saved_password})
	return saved_password != "" && saved_password == password
}
