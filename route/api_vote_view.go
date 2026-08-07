package route

import (
	"regexp"
	"strings"

	"opennamu/route/tool"
)

func Api_vote_view(config tool.Config, set_id string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	name := ""
	subject := ""
	data := ""
	type_str := ""

	tool.QueryRow_DB(
		db,
		`select name, subject, data, type from vote where id = ? and user = ""`,
		[]any{&name, &subject, &data, &type_str},
		set_id,
	)

	if name == "" && subject == "" && data == "" && type_str == "" {
		return map[string]any{
			"response": "error",
			"data":     "not found",
		}
	}

	end_date := ""
	tool.QueryRow_DB(
		db,
		`select data from vote where id = ? and name = "end_date" and type = "option"`,
		[]any{&end_date},
		set_id,
	)

	time_db := ""
	time_field := strings.Fields(end_date)
	if len(time_field) > 0 {
		time_db = time_field[0]
	}

	re := regexp.MustCompile(`[^\n]+`)
	vote_data := re.FindAllString(
		strings.ReplaceAll(data, "\r", ""),
		-1,
	)

	result := make(map[string]any)
	result["name"] = name
	result["subject"] = subject
	result["data"] = vote_data
	result["type"] = type_str
	result["end_date"] = time_db

	return_data := make(map[string]any)
	return_data["response"] = "ok"
	return_data["data"] = result

	return return_data
}
