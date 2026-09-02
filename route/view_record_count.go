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

	count_range := func(start time.Time, end time.Time) (int, int) {
		rows := tool.Get_history_length_range_rows(
			db,
			user_name,
			start.Format("2006-01-02 15:04:05"),
			end.Format("2006-01-02 15:04:05"),
		)
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

	period_days := func(start time.Time, end time.Time) int {
		days := int(end.Sub(start) / (24 * time.Hour))
		if days < 1 {
			return 1
		}
		return days
	}

	now := time.Now()
	today_start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tomorrow_start := today_start.AddDate(0, 0, 1)
	yesterday_start := today_start.AddDate(0, 0, -1)
	week_start := today_start.AddDate(0, 0, -(int(today_start.Weekday())+6)%7)
	last_week_start := week_start.AddDate(0, 0, -7)
	month_start := time.Date(today_start.Year(), today_start.Month(), 1, 0, 0, 0, 0, today_start.Location())
	last_month_start := month_start.AddDate(0, -1, 0)

	today_count, today_length := count_range(today_start, tomorrow_start)
	yesterday_count, yesterday_length := count_range(yesterday_start, today_start)
	week_count, week_length := count_range(week_start, tomorrow_start)
	last_week_count, last_week_length := count_range(last_week_start, week_start)
	month_count, month_length := count_range(month_start, tomorrow_start)
	last_month_count, last_month_length := count_range(last_month_start, month_start)

	record_line := func(label string, count int, length int) string {
		return `<li>` + label + ` : ` + tool.Get_language(db, "record_edit_count", true) + ` ` + strconv.Itoa(count) + `, ` + tool.Get_language(db, "record_length", true) + ` ` + strconv.Itoa(length) + `</li>`
	}
	record_average_line := func(label string, count int, length int, start time.Time, end time.Time) string {
		average := length / period_days(start, end)
		return `<li>` + label + ` : ` + tool.Get_language(db, "record_edit_count", true) + ` ` + strconv.Itoa(count) + `, ` + tool.Get_language(db, "record_length", true) + ` ` + strconv.Itoa(length) + `, ` + tool.Get_language(db, "record_daily_average", true) + ` ` + strconv.Itoa(average) + `</li>`
	}
	record_diff_line := func(label string, length int) string {
		return `<li>` + label + ` : ` + tool.Get_language(db, "record_length_diff", true) + ` ` + strconv.Itoa(length) + `</li>`
	}

	user_url := tool.Url_parser(user_name)
	body := `<ul><li><a href="/record/` + user_url + `">` + tool.Get_language(db, "edit_record", true) + `</a> : ` + history + `</li><li><a href="/record/topic/` + user_url + `">` + tool.Get_language(db, "discussion_record", true) + `</a> : ` + topic + `</li><li>bbs : ` + bbs + `</li><hr>`
	body += record_line(tool.Get_language(db, "record_today", true), today_count, today_length)
	body += record_line(tool.Get_language(db, "record_yesterday", true), yesterday_count, yesterday_length)
	body += record_diff_line(tool.Get_language(db, "record_today", true)+" - "+tool.Get_language(db, "record_yesterday", true), today_length-yesterday_length)
	body += `<hr>`
	body += record_average_line(tool.Get_language(db, "record_this_week", true), week_count, week_length, week_start, tomorrow_start)
	body += record_line(tool.Get_language(db, "record_last_week", true), last_week_count, last_week_length)
	body += record_diff_line(tool.Get_language(db, "record_this_week", true)+" - "+tool.Get_language(db, "record_last_week", true), week_length-last_week_length)
	body += `<hr>`
	body += record_average_line(tool.Get_language(db, "record_this_month", true), month_count, month_length, month_start, tomorrow_start)
	body += record_line(tool.Get_language(db, "record_last_month", true), last_month_count, last_month_length)
	body += record_diff_line(tool.Get_language(db, "record_this_month", true)+" - "+tool.Get_language(db, "record_last_month", true), month_length-last_month_length)
	body += `</ul>`
	return list_extra_page(db, config, tool.Get_language(db, "count", true), body)
}
