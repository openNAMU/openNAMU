package route

import (
	"strconv"
	"strings"
	"time"

	"opennamu/route/tool"
)

func View_record_count(config tool.Config, user_name string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if user_name != "" && user_name != config.IP && !tool.Check_permission(db, "hidel", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if user_name == "" {
		user_name = config.IP
	}
	history := tool.Get_history_count(db, user_name)
	topic := tool.Get_topic_count(db, user_name)
	bbs := tool.Get_bbs_comment_count(db, user_name)

	count_day := func(day string) (int, int) {
		rows := tool.Get_history_length_rows(db, day, user_name)
		defer rows.Close()
		count := 0
		length := 0
		for rows.Next() {
			value := ""
			if rows.Scan(&value) != nil {
				continue
			}
			value = strings.TrimLeft(value, "+-")
			length += tool.Str_to_int(value)
			count++
		}
		return count, length
	}

	today_count, today_length := count_day(tool.Get_date())
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	yesterday_count, yesterday_length := count_day(yesterday)
	user_url := tool.Url_parser(user_name)
	body := `<ul><li><a href="/record/` + user_url + `">` + tool.Get_language(db, "edit_record", true) + `</a> : ` + history + `</li><li><a href="/record/topic/` + user_url + `">` + tool.Get_language(db, "discussion_record", true) + `</a> : ` + topic + `</li><li>bbs : ` + bbs + `</li><hr><li>(` + tool.Get_language(db, "beta", true) + `) TODAY : ` + strconv.Itoa(today_count) + `</li><li>(` + tool.Get_language(db, "beta", true) + `) TODAY LEN : ` + strconv.Itoa(today_length) + `</li><li>(` + tool.Get_language(db, "beta", true) + `) TODAY DIFF : ` + strconv.Itoa(today_length-yesterday_length) + `</li><hr><li>(` + tool.Get_language(db, "beta", true) + `) YESTERDAY : ` + strconv.Itoa(yesterday_count) + `</li><li>(` + tool.Get_language(db, "beta", true) + `) YESTERDAY LEN : ` + strconv.Itoa(yesterday_length) + `</li></ul>`
	return list_extra_page(db, config, tool.Get_language(db, "count", true), body)
}
