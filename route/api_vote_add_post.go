package route

import (
	stdjson "encoding/json"
	"net/url"
	"strconv"
	"strings"
	"time"

	"opennamu/route/tool"
)

func Api_vote_add_post(config tool.Config, data string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "vote", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	values := url.Values{}
	if strings.TrimSpace(data) != "" {
		if err := stdjson.Unmarshal([]byte(data), &values); err != nil {
			parsed, parse_err := url.ParseQuery(data)
			if parse_err != nil {
				return_data["response"] = "error"
				return_data["data"] = "invalid data"
				return return_data
			}
			values = parsed
		}
	}

	name := strings.TrimSpace(values.Get("name"))
	subject := values.Get("subject")
	options := vote_options(values.Get("data"))
	if name == "" || len(options) < 2 {
		return_data["response"] = "error"
		return_data["data"] = "invalid data"
		return return_data
	}

	for index := range options {
		options[index] = strings.TrimSpace(options[index])
		if options[index] == "" {
			return_data["response"] = "error"
			return_data["data"] = "invalid data"
			return return_data
		}
	}

	acl := strings.TrimSpace(values.Get("acl_select"))
	if !acl_value_valid(db, acl) {
		return_data["response"] = "error"
		return_data["data"] = "invalid acl"
		return return_data
	}

	date := strings.TrimSpace(values.Get("date"))
	if date != "" {
		if parsed_date, err := time.Parse("2006-01-02", date); err != nil || parsed_date.Before(time.Now().AddDate(0, 0, -1)) {
			return_data["response"] = "error"
			return_data["data"] = "invalid date"
			return return_data
		}
	}

	last_id := "0"
	tool.QueryRow_DB(db, "select id from vote where type != 'option' order by id + 0 desc limit 1", []any{&last_id})
	id := strconv.Itoa(tool.Str_to_int(last_id) + 1)
	type_data := "n_open"
	if values.Get("open_select") == "Y" || values.Get("open_select") == "on" || values.Get("open") == "1" {
		type_data = "open"
	}

	tool.Exec_DB(db, "insert into vote (name, id, subject, data, user, type, acl) values (?, ?, ?, ?, '', ?, ?)", name, id, subject, strings.Join(options, "\n"), type_data, acl)
	tool.Exec_DB(db, "insert into vote (name, id, subject, data, user, type, acl) values ('open_user', ?, '', ?, '', 'option', '')", id, config.IP)
	if values.Get("limitless") == "" && values.Get("limitless") != "Y" && date != "" {
		tool.Exec_DB(db, "insert into vote (name, id, subject, data, user, type, acl) values ('end_date', ?, '', ?, '', 'option', '')", id, date)
	}

	return_data["response"] = "ok"
	return_data["data"] = id
	return return_data
}
