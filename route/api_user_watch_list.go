package route

import (
	"database/sql"
	"sort"
	"strings"

	"opennamu/route/tool"
)

type User_watch_item struct {
	Data      string
	Display   string
	Date      string
	Set_id    string
	Set_code  string
	Tag       string `json:"Tag,omitempty"`
	Available bool   `json:"Available,omitempty"`
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
		} else if do_type == "bbs_watchlist" || do_type == bbs_bookmark_list_type {
			if do_type == "bbs_watchlist" {
				if strings.HasPrefix(data, "-1-") {
					item.Set_id = thread_bbs_id
					item.Set_code = strings.TrimPrefix(data, "-1-")
				} else {
					watch_data := strings.SplitN(data, "-", 2)
					if len(watch_data) == 1 {
						item.Set_id = thread_bbs_id
						item.Set_code = watch_data[0]
					} else {
						item.Set_id = watch_data[0]
						item.Set_code = watch_data[1]
					}
				}
			} else {
				separator := strings.LastIndex(data, "-")
				if separator > 0 && separator < len(data)-1 {
					item.Set_id = data[:separator]
					item.Set_code = data[separator+1:]
				} else {
					item.Set_id = thread_bbs_id
					item.Set_code = data
				}
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
		} else if do_type == bbs_board_watch_list_type {
			item.Set_id = data
			if bbs_name := tool.Get_bbs_set_data(db, item.Set_id, "bbs_name"); bbs_name != "" {
				item.Display = bbs_name
			}
		} else if do_type == bbs_tag_watch_list_type {
			item.Tag = data
			if _, legacy_tag, is_legacy := strings.Cut(data, "\n"); is_legacy {
				item.Tag = legacy_tag
			}
			item.Tag = strings.TrimSpace(item.Tag)
			if item.Tag == "" || !Bbs_tag_watch_value_valid(item.Tag) {
				continue
			}
			item.Data = item.Tag
			item.Display = "#" + item.Tag
			item.Available = true
		}
		item_list = append(item_list, item)
	}

	return item_list
}

func watch_key(watch_data string) string {
	if strings.HasPrefix(watch_data, "-1-") {
		return strings.TrimPrefix(watch_data, "-1-")
	}
	return watch_data
}

func Bbs_watch_merged_data(db *sql.DB, name string) []string {
	bbs_rows := []string{}
	rows := tool.Query_DB(db, "select data from user_set where name = ? and id = ?", "bbs_watchlist", name)
	for rows.Next() {
		watch_data := ""
		if rows.Scan(&watch_data) != nil {
			continue
		}
		bbs_rows = append(bbs_rows, watch_data)
	}
	rows.Close()

	thread_rows := []string{}
	rows = tool.Query_DB(db, "select data from user_set where name = ? and id = ?", "thread_watchlist", name)
	for rows.Next() {
		watch_data := ""
		if rows.Scan(&watch_data) != nil {
			continue
		}
		thread_rows = append(thread_rows, watch_data)
	}
	rows.Close()

	seen := map[string]bool{}
	bbs_keys := []string{}
	thread_keys := []string{}

	for _, watch_data := range bbs_rows {
		key := watch_key(watch_data)
		if seen[key] {
			continue
		}
		seen[key] = true
		bbs_keys = append(bbs_keys, key)
	}

	for _, watch_data := range thread_rows {
		key := watch_key(watch_data)
		if seen[key] {
			continue
		}
		seen[key] = true
		thread_keys = append(thread_keys, key)
	}

	sort.Strings(bbs_keys)
	sort.Strings(thread_keys)

	data_list := []string{}
	data_list = append(data_list, bbs_keys...)
	data_list = append(data_list, thread_keys...)
	return data_list
}

func bbs_watch_merged_count(db *sql.DB, name string) int {
	count := 0
	query := `
		select count(distinct merge_key)
		from (
			select case when coalesce(data, '') like '-1-%' then substr(coalesce(data, ''), 4) else coalesce(data, '') end as merge_key
			from user_set
			where name = 'bbs_watchlist' and id = ?
			union
			select case when coalesce(data, '') like '-1-%' then substr(coalesce(data, ''), 4) else coalesce(data, '') end as merge_key
			from user_set
			where name = 'thread_watchlist' and id = ?
		) as merged
	`
	tool.QueryRow_DB(db, query, []any{&count}, name, name)
	return count
}

func Api_user_watch_list(config tool.Config, name string, num_str string, do_type string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if do_type != "watchlist" && do_type != "star_doc" && do_type != "thread_watchlist" && do_type != "bbs_watchlist" && do_type != bbs_bookmark_list_type && do_type != bbs_board_watch_list_type && do_type != bbs_tag_watch_list_type {
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
		if do_type == "thread_watchlist" || do_type == "bbs_watchlist" {
			data_list := []string{}
			if do_type == "thread_watchlist" {
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
			} else {
				data_list = Bbs_watch_merged_data(db, name)
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
		if do_type == bbs_tag_watch_list_type {
			rows := tool.Query_DB(
				db,
				"select data from user_set where name = ? and id = ? order by data asc",
				do_type,
				name,
			)
			tag_list := []string{}
			seen_tags := map[string]bool{}
			for rows.Next() {
				data := ""
				if err := rows.Scan(&data); err != nil {
					rows.Close()
					panic(err)
				}
				item_list := User_watch_list_items(db, []string{data}, do_type)
				if len(item_list) == 0 || seen_tags[item_list[0].Tag] {
					continue
				}
				seen_tags[item_list[0].Tag] = true
				tag_list = append(tag_list, item_list[0].Tag)
			}
			if err := rows.Err(); err != nil {
				rows.Close()
				panic(err)
			}
			rows.Close()
			sort.Strings(tag_list)

			start := num
			if start < 0 {
				start = 0
			}
			if start > len(tag_list) {
				start = len(tag_list)
			}
			end := start + 50
			if end > len(tag_list) {
				end = len(tag_list)
			}
			return_data["response"] = "ok"
			return_data["data"] = User_watch_list_items(db, tag_list[start:end], do_type)
			return_data["has_next"] = end < len(tag_list)
			return return_data
		}
		query := "select data from user_set where name = ? and id = ? limit ?, 50"
		if do_type == "star_doc" {
			query = `select data from user_set where name = ? and id = ? order by coalesce((select date from history where history.title = user_set.data order by id + 0 desc limit 1), '') desc, data asc limit ?, 50`
		} else if do_type == bbs_bookmark_list_type || do_type == bbs_board_watch_list_type {
			query = "select data from user_set where name = ? and id = ? order by data asc limit ?, 51"
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

		item_list := User_watch_list_items(db, data_list, do_type)
		has_next := false
		if do_type == bbs_bookmark_list_type || do_type == bbs_board_watch_list_type {
			has_next = len(item_list) > 50
			if has_next {
				item_list = item_list[:50]
			}
			for index := range item_list {
				item := &item_list[index]
				if do_type == bbs_bookmark_list_type {
					_, item.Available = Bbs_post_view_auth(db, item.Set_id, item.Set_code, ip)
				} else {
					item.Available = tool.Check_acl(db, item.Set_id, "", "bbs_view", ip)
				}
				if !item.Available {
					item.Display = ""
					item.Date = ""
					if name != ip {
						item.Data = ""
						item.Set_id = ""
						item.Set_code = ""
					}
				}
			}
		}
		return_data["response"] = "ok"
		return_data["data"] = item_list
		if do_type == bbs_bookmark_list_type || do_type == bbs_board_watch_list_type {
			return_data["has_next"] = has_next
		}
	}

	return return_data
}
