package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Do_add_history(db *sql.DB, doc_name string, data string, date string, ip string, send string, length string, mode string, type_check string) {
	tool.Do_add_history(db, doc_name, data, date, ip, send, length, mode, type_check)
}
