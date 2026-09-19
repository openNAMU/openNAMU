package route

import (
	"database/sql"
	"sort"

	"opennamu/route/tool"
)

func Bbs_list(db *sql.DB) map[string]string {
	rows := tool.Query_DB(
		db,
		"select set_data, set_id from bbs_set where set_name = 'bbs_name'",
	)
	defer rows.Close()

	data_list := map[string]string{}

	for rows.Next() {
		var name string
		var id string

		err := rows.Scan(&name, &id)
		if err != nil {
			panic(err)
		}

		data_list[name] = id
	}

	return data_list
}

type BBS_item struct {
	Id   string
	Name string
	Type string
	Date string
}

func Api_bbs_list(config tool.Config) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "bbs_main_view", config.IP) {
		return map[string]any{"response": "require auth", "data": [][]string{}}
	}

	rows := tool.Query_DB(
		db,
		`select n.set_id, n.set_data, coalesce(t.set_data, ''), coalesce(d.set_data, '') from bbs_set n
		left join bbs_set t on t.set_id = n.set_id and t.set_name = 'bbs_type'
		left join (
			select set_id, max(set_code + 0) as max_code
			from bbs_data
			where set_name = 'date'
			group by set_id
		) latest on latest.set_id = n.set_id
		left join bbs_data d on d.set_id = n.set_id and d.set_name = 'date'
			and d.set_code + 0 = latest.max_code
		where n.set_name = 'bbs_name'`,
	)
	defer rows.Close()

	items := make([]BBS_item, 0, 8)

	for rows.Next() {
		var id string
		var bbs_name string
		var bbs_type string
		var bbs_date string

		if err := rows.Scan(&id, &bbs_name, &bbs_type, &bbs_date); err != nil {
			panic(err)
		}

		if id == "0" {
			bbs_name = tool.Get_language(db, "wiki_comment_bbs", true)
		} else if id == "-1" {
			bbs_name = tool.Get_language(db, "thread_bbs", true)
		} else if id == "-2" {
			bbs_name = tool.Get_language(db, "report_bbs", true)
		}

		if !tool.Check_acl(db, id, "", "bbs_view", config.IP) {
			continue
		}

		items = append(items, BBS_item{
			Id:   id,
			Name: bbs_name,
			Type: bbs_type,
			Date: bbs_date,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Date > items[j].Date
	})

	data_list_sub := make([][]string, 0, len(items))
	for _, item := range items {
		data_list_sub = append(data_list_sub, []string{item.Name, item.Id, item.Type, item.Date})
	}

	return_data := make(map[string]any)
	return_data["response"] = "ok"
	return_data["data"] = data_list_sub

	return return_data
}
