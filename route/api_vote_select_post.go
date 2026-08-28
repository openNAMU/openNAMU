package route

import (
	"strconv"
	"strings"

	"opennamu/route/tool"
)

func Api_vote_select_post(config tool.Config, id string, vote_data string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	name, subject, data, type_data := "", "", "", ""
	if !tool.QueryRow_DB(
		db,
		"select name, subject, data, type from vote where id = ? and user = ''",
		[]any{&name, &subject, &data, &type_data},
		id,
	) {
		return_data["response"] = "not exist"
		return_data["data"] = "vote"
		return return_data
	}

	if type_data == "close" || type_data == "n_close" {
		return_data["response"] = "closed"
		return return_data
	}
	if !tool.Check_acl(db, "", id, "vote", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	voted := ""
	if tool.QueryRow_DB(db, "select user from vote where id = ? and user = ?", []any{&voted}, id, config.IP) {
		return_data["response"] = "already voted"
		return return_data
	}

	end_date := ""
	tool.QueryRow_DB(
		db,
		"select data from vote where id = ? and name = 'end_date' and type = 'option'",
		[]any{&end_date},
		id,
	)
	now_time := tool.Get_time()
	if end_date != "" && !strings.HasPrefix(now_time, end_date) && strings.Split(now_time, " ")[0] > strings.Split(end_date, " ")[0] {
		return_data["response"] = "closed"
		return return_data
	}

	options := vote_options(data)
	choice, err := strconv.Atoi(vote_data)
	if err != nil || choice < 0 || choice >= len(options) {
		return_data["response"] = "error"
		return_data["data"] = "invalid data"
		return return_data
	}

	tool.Exec_DB(
		db,
		"insert into vote (name, id, subject, data, user, type) values ('', ?, '', ?, ?, 'select')",
		id,
		strconv.Itoa(choice),
		config.IP,
	)

	return_data["response"] = "ok"
	return return_data
}
