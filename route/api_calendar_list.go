package route

import (
	"opennamu/route/tool"
)

func Api_calendar_list(config tool.Config) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !Calendar_view_allowed(db, config.IP) {
		return_data["response"] = "require auth"
		return_data["data"] = [][]string{}
		return return_data
	}

	rows := tool.Query_DB(
		db,
		`select n.set_id, n.set_data, coalesce(s.set_data, '')
		from calendar n
		left join calendar s on s.set_id = n.set_id and s.set_code = '' and s.set_name = 'subject'
		where n.set_code = '' and n.set_name = 'name'
		order by n.set_id + 0 desc`,
	)
	defer rows.Close()

	data_list := [][]string{}
	for rows.Next() {
		id := ""
		name := ""
		subject := ""
		if err := rows.Scan(&id, &name, &subject); err != nil {
			panic(err)
		}

		data_list = append(data_list, []string{id, name, subject})
	}

	return_data["response"] = "ok"
	return_data["data"] = data_list
	return return_data
}
