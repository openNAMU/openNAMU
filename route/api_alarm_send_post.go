package route

import (
	"database/sql"
	"strings"

	"opennamu/route/tool"
)

func Api_alarm_send_post(config tool.Config, target string, data string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	result := map[string]any{"response": "error"}
	if !tool.Check_permission(db, "alarm_send", config.IP) {
		result["response"] = "require auth"
		return result
	}

	target = strings.TrimSpace(target)
	data = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(data, "\r\n", "\n"), "\r", "\n"))
	if target == "" || target == config.IP || tool.IP_or_user(target) || !tool.Get_user_set_exists(db, target, "pw") {
		result["response"] = "invalid"
		return result
	}
	if data == "" || tool.Get_len(data) > 1000 {
		result["response"] = "invalid"
		return result
	}

	data = strings.ReplaceAll(tool.HTML_escape(data), "\n", "<br>")
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		tool.Send_alarm(tx, config.IP, target, data)
		tool.Do_insert_auth_history(tx, config.IP, "alarm_send ("+target+")")
		return nil
	}); err != nil {
		panic(err)
	}
	result["response"] = "ok"
	return result
}
