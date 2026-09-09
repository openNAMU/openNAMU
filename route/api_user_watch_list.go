package route

import (
	"strings"

	"opennamu/route/tool"
)

func Api_user_watch_list(config tool.Config, name string, num_str string, do_type string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if do_type != "watchlist" && do_type != "star_doc" && do_type != "thread_watchlist" && do_type != "bbs_watchlist" {
		do_type = "star_doc"
	}

	page := tool.Str_to_int(num_str)
	num := 0
	if page*50 > 0 {
		num = page*50 - 50
	}

	ip := config.IP

	return_data := make(map[string]any)

	if ip != name && !tool.Check_permission(db, "view_user_watchlist", ip) {
		return_data["response"] = "require auth"
		return_data["data"] = []string{}
	} else {
		if do_type == "thread_watchlist" {
			data_list := []string{}
			seen := map[string]bool{}
			for _, watch_type := range []string{"bbs_watchlist", "thread_watchlist"} {
				rows := tool.Query_DB(
					db,
					"select data from user_set where name = ? and id = ?",
					watch_type,
					name,
				)
				for rows.Next() {
					watch_data := ""
					if rows.Scan(&watch_data) != nil {
						continue
					}
					if watch_type == "bbs_watchlist" {
						if !strings.HasPrefix(watch_data, "-1-") {
							continue
						}
						watch_data = strings.TrimPrefix(watch_data, "-1-")
					}
					if seen[watch_data] {
						continue
					}
					seen[watch_data] = true
					data_list = append(data_list, watch_data)
				}
				rows.Close()
			}
			start := num
			if start > len(data_list) {
				start = len(data_list)
			}
			end := start + 50
			if end > len(data_list) {
				end = len(data_list)
			}
			return_data["response"] = "ok"
			return_data["data"] = data_list[start:end]
			return return_data
		}
		query := "select data from user_set where name = ? and id = ? limit ?, 50"
		if do_type == "star_doc" {
			query = `select data from user_set where name = ? and id = ? order by coalesce((select date from history where history.title = user_set.data order by id + 0 desc limit 1), '') desc, data asc limit ?, 50`
		}

		rows := tool.Query_DB(
			db,
			query,
			do_type,
			name,
			num,
		)
		defer rows.Close()

		data_list := []string{}

		for rows.Next() {
			var title_data string

			err := rows.Scan(&title_data)
			if err != nil {
				panic(err)
			}

			data_list = append(data_list, title_data)
		}

		return_data["response"] = "ok"
		return_data["data"] = data_list
	}

	return return_data
}
