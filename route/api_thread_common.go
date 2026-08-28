package route

import (
	"database/sql"
	"strconv"

	"opennamu/route/tool"
)

func thread_next_code(db *sql.DB) string {
	last_code := "0"
	tool.QueryRow_DB(
		db,
		"select code from topic order by code + 0 desc limit 1",
		[]any{&last_code},
	)
	return strconv.Itoa(tool.Str_to_int(last_code) + 1)
}

func thread_next_id(db *sql.DB, topic_num string) string {
	last_id := "0"
	tool.QueryRow_DB(
		db,
		"select id from topic where code = ? order by id + 0 desc limit 1",
		[]any{&last_id},
		topic_num,
	)
	return strconv.Itoa(tool.Str_to_int(last_id) + 1)
}

func thread_add(db *sql.DB, topic_num string, id string, data string, ip string, top string) {
	tool.Exec_DB(
		db,
		"insert into topic (id, data, date, ip, block, top, code) values (?, ?, ?, ?, '', ?, ?)",
		id,
		data,
		tool.Get_time(),
		ip,
		top,
		topic_num,
	)
}
