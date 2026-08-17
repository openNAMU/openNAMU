package route

import (
	"database/sql"
	"opennamu/route/tool"
)

func document_safe_page(db *sql.DB, config tool.Config, title string, body string) string {
	return tool.Get_template(db, config, title, body, []any{}, [][]any{{"other", tool.Get_language(db, "return", true)}}, map[string]string{})
}
