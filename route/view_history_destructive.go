package route

import (
	"database/sql"
	"opennamu/route/tool"
	"strconv"
	"strings"
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
