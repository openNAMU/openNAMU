package route

import "opennamu/route/tool"

func Api_list_auth_give(config tool.Config, num string, include_default bool) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	return_data["response"] = "ok"
	return_data["data"] = [][]string{}

	if !tool.Check_permission(db, "give_range", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	page_int := tool.Str_to_int(num)
	if page_int > 0 {
		page_int = (page_int - 1) * 50
	} else {
		page_int = 0
	}

	default_filter := ""
	if !include_default {
		default_filter = " and a.data not in ('user', 'ip')"
	}

	query := `select
		a.id,
		a.data,
		'normal',
		coalesce((select blocker from rb where block = a.id and band = '' and ongoing = '1' order by today desc limit 1), ''),
		coalesce((select today from rb where block = a.id and band = '' and ongoing = '1' order by today desc limit 1), (select data from user_set where id = a.id and name = 'date' limit 1), ''),
		coalesce((select data from user_set where id = a.id and name = 'acl_end' limit 1), ''),
		coalesce((select why from rb where block = a.id and band = '' and ongoing = '1' order by today desc limit 1), ''),
		'1',
		case when a.data in ('user', 'ip') then '1' else '0' end
	from user_set as a
	where a.name = 'acl' and a.data != ''` + default_filter + `
	union all
	select block, login, band, blocker, today, end, why, ongoing, '0'
	from rb
	where ongoing = '1' and (band = 'regex' or band = 'cidr')
	order by 5 desc, 1
	limit ?, 50`

	rows := tool.Query_DB(db, query, page_int)
	defer rows.Close()

	data_list := [][]string{}
	for rows.Next() {
		data := []string{"", "", "", "", "", "", "", "", ""}
		if err := rows.Scan(&data[0], &data[1], &data[2], &data[3], &data[4], &data[5], &data[6], &data[7], &data[8]); err != nil {
			panic(err)
		}

		if data[2] == "regex" || data[2] == "cidr" {
			data[1] = tool.Get_auth_target_group(db, data[0], data[2])
		}

		data_list = append(data_list, data)
	}

	return_data["data"] = data_list
	return return_data
}
