package route

import (
	"database/sql"
	"strings"

	"opennamu/route/tool"
)

type User_watch_item struct {
	Data     string
	Display  string
	Date     string
	Set_id   string
	Set_code string
}

func User_watch_list_items(db *sql.DB, data_list []string, do_type string) []User_watch_item {
	item_list := make([]User_watch_item, 0, len(data_list))
	for _, data := range data_list {
		item := User_watch_item{
			Data:    data,
			Display: data,
		}
		if do_type == "watchlist" {
			item.Date = tool.Get_history_date(db, data)
		} else if do_type == "bbs_watchlist" {
			if strings.HasPrefix(data, "-1-") {
				item.Set_id = "-1"
				item.Set_code = strings.TrimPrefix(data, "-1-")
			} else {
				watch_data := strings.SplitN(data, "-", 2)
				if len(watch_data) != 2 {
					continue
				}
				item.Set_id = watch_data[0]
				item.Set_code = watch_data[1]
			}

			post_title, _ := tool.Get_bbs_data_value(db, item.Set_id, item.Set_code, "title")
			bbs_name := tool.Get_bbs_set_data(db, item.Set_id, "bbs_name")
			if post_title != "" {
				item.Display = post_title
				if bbs_name != "" {
					item.Display = bbs_name + " - " + item.Display
				}
			}
			item.Date, _ = tool.Get_bbs_data_value(db, item.Set_id, item.Set_code, "date")
		}
		item_list = append(item_list, item)
	}

	return item_list
}

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
		return_data["data"] = []User_watch_item{}
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
			return_data["data"] = User_watch_list_items(db, data_list[start:end], do_type)
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
		return_data["data"] = User_watch_list_items(db, data_list, do_type)
	}

	return return_data
}
