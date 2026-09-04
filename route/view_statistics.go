package route

import (
	"strconv"
	"strings"
	"time"

	"opennamu/route/tool"
)

func View_statistics(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	get_language := func(name string) string {
		return tool.Get_language(db, name, true)
	}
	statistics_line := func(name string, count int) string {
		return "<li>" + get_language(name) + " : " + strconv.Itoa(count) + "</li>"
	}

	body := strings.Builder{}
	body.WriteString("<ul>")
	body.WriteString(statistics_line("statistics_document_count", tool.Get_statistics_count(db, "document")))
	body.WriteString(statistics_line("statistics_user_count", tool.Get_statistics_count(db, "user")))
	body.WriteString(statistics_line("statistics_edit_count", tool.Get_statistics_count(db, "edit")))
	body.WriteString(statistics_line("statistics_bbs_post_count", tool.Get_statistics_count(db, "bbs_post")))
	body.WriteString(statistics_line("statistics_bbs_comment_count", tool.Get_statistics_count(db, "bbs_comment")))

	now := time.Now()
	month_start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	next_month_start := month_start.AddDate(0, 1, 0)
	month_start_text := month_start.Format("2006-01-02 15:04:05")
	next_month_start_text := next_month_start.Format("2006-01-02 15:04:05")
	body.WriteString(statistics_line(
		"statistics_monthly_edit_count",
		tool.Get_month_history_count(db, month_start_text, next_month_start_text),
	))

	rows := tool.Get_month_contributor_rows(db, month_start_text, next_month_start_text)
	user_name := ""
	count := 0
	if rows.Next() && rows.Scan(&user_name, &count) == nil {
		body.WriteString(
			"<li>" + get_language("statistics_monthly_top_contributor") + " : " +
				tool.IP_parser(db, user_name, config.IP) + " (" + strconv.Itoa(count) + ")</li>",
		)
	} else {
		body.WriteString(
			"<li>" + get_language("statistics_monthly_top_contributor") + " : " +
				get_language("data_missing") + "</li>",
		)
	}
	rows.Close()

	body.WriteString("</ul>")
	return list_extra_page(db, config, get_language("statistics"), body.String())
}
