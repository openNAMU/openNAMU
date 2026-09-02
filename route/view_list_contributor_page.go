package route

import (
	"strconv"
	"time"

	"opennamu/route/tool"
)

func View_list_contributor_page(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	now := time.Now()
	month_start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	next_month_start := month_start.AddDate(0, 1, 0)
	rows := tool.Get_month_contributor_rows(
		db,
		month_start.Format("2006-01-02 15:04:05"),
		next_month_start.Format("2006-01-02 15:04:05"),
	)
	defer rows.Close()

	body := ""
	user_name := ""
	count := 0
	if rows.Next() && rows.Scan(&user_name, &count) == nil {
		body = tool.Get_list_ui(
			tool.IP_parser(db, user_name, config.IP),
			tool.Get_language(db, "record_edit_count", true)+" : "+strconv.Itoa(count),
			"",
			"",
		)
	} else {
		body = tool.Get_language(db, "data_missing", true)
	}

	return list_extra_page(db, config, tool.Get_language(db, "monthly_top_contributor", true), body)
}
