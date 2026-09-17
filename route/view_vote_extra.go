package route

import (
	"database/sql"
	"opennamu/route/tool"
	"strings"
)

func Vote_page(db *sql.DB, config tool.Config, title string, body string) string {
	return tool.Get_template(db, config, title, body, []any{}, [][]any{{"vote", tool.Get_language(db, "return", true)}}, map[string]string{})
}

func Vote_options(data string) []string {
	lines := strings.Split(strings.ReplaceAll(data, "\r", ""), "\n")
	options := []string{}
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			options = append(options, line)
		}
	}
	return options
}
