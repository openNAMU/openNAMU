package route

import (
	"database/sql"
	"opennamu/route/tool"
	"strings"
)

func user_email_allowed(db *sql.DB, email string) bool {
	at_index := strings.LastIndex(email, "@")
	if at_index <= 0 || at_index == len(email)-1 {
		return false
	}
	domain := strings.TrimSpace(email[at_index+1:])
	rows := tool.Query_DB(db, "select html from html_filter where kind = 'email'")
	defer rows.Close()
	for rows.Next() {
		allowed_domain := ""
		if rows.Scan(&allowed_domain) != nil {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(allowed_domain), domain) {
			return true
		}
	}
	return false
}
