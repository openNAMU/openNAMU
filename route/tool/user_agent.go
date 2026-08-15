package tool

import "database/sql"

func Record_user_agent(db *sql.DB, user_name string, user_ip string, user_agent string, date string) {
	ua_get := ""
	QueryRow_DB(db, "select data from other where name = 'ua_get'", []any{&ua_get})
	if ua_get != "" {
		return
	}
	Exec_DB(db, "insert into ua_d (name, ip, ua, today, sub) values (?, ?, ?, ?, '')", user_name, user_ip, user_agent, date)
}
