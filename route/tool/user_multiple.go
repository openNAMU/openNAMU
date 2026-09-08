package tool

import (
	"database/sql"
	"strings"
)

type UA_multiple_data struct {
	IP    string
	Count int
	Date  string
	Names []string
}

func Get_ua_multiple_rows(db *sql.DB, search string, sort string, offset int) []UA_multiple_data {
	query := "select base.ip, count(distinct base.name), max(base.today) from ua_d as base where base.ip != '' and base.name != ''"
	values := []any{}

	if search != "" {
		query += " and (base.ip like ? or exists (select 1 from ua_d as matched where matched.ip = base.ip and matched.name like ?))"
		search_value := "%" + search + "%"
		values = append(values, search_value, search_value)
	}

	query += " group by base.ip having count(distinct base.name) >= 2"
	if sort == "count" {
		query += " order by count(distinct base.name) desc, max(base.today) desc, base.ip asc"
	} else {
		query += " order by max(base.today) desc, count(distinct base.name) desc, base.ip asc"
	}
	query += " limit ?, 50"
	values = append(values, offset)

	rows := Query_DB(db, query, values...)
	defer rows.Close()

	result := []UA_multiple_data{}
	for rows.Next() {
		data := UA_multiple_data{}
		if rows.Scan(&data.IP, &data.Count, &data.Date) == nil {
			result = append(result, data)
		}
	}

	if len(result) == 0 {
		return result
	}

	placeholder_list := make([]string, len(result))
	value_list := make([]any, len(result))
	index_list := map[string]int{}
	for index, data := range result {
		placeholder_list[index] = "?"
		value_list[index] = data.IP
		index_list[data.IP] = index
	}

	name_query := "select ip, name, max(today) from ua_d where ip in (" + strings.Join(placeholder_list, ",") + ") and name != '' group by ip, name order by max(today) desc, name asc"
	name_rows := Query_DB(db, name_query, value_list...)
	defer name_rows.Close()

	for name_rows.Next() {
		ip, name, date := "", "", ""
		if name_rows.Scan(&ip, &name, &date) != nil {
			continue
		}
		index, ok := index_list[ip]
		if !ok {
			continue
		}
		result[index].Names = append(result[index].Names, name)
	}

	return result
}
