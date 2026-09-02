package route

import (
	"strconv"
	"time"

	"opennamu/route/tool"
)

func View_bbs_contributor(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "bbs_main_view", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	now := time.Now()
	month_start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	next_month_start := month_start.AddDate(0, 1, 0)
	auth_info := tool.Get_auth_info(db, config.IP)
	bbs_ids := []string{}
	for _, bbs_id := range bbs_list(db) {
		if tool.Check_acl(db, bbs_id, "", "bbs_view", config.IP) && bbs_post_view_allowed(db, bbs_id, "", config.IP, auth_info) {
			bbs_ids = append(bbs_ids, bbs_id)
		}
	}
	if len(bbs_ids) == 0 {
		return list_extra_page(db, config, tool.Get_language(db, "monthly_bbs_contributor", true), tool.Get_language(db, "data_missing", true))
	}

	rows := tool.Get_month_bbs_contributor_rows(
		db,
		month_start.Format("2006-01-02 15:04:05"),
		next_month_start.Format("2006-01-02 15:04:05"),
		bbs_ids,
	)
	defer rows.Close()

	body := ""
	user_name := ""
	count := 0
	if rows.Next() && rows.Scan(&user_name, &count) == nil {
		body = tool.Get_list_ui(
			tool.IP_parser(db, user_name, config.IP),
			tool.Get_language(db, "record_bbs_count", true)+" : "+strconv.Itoa(count),
			"",
			"",
		)
	} else {
		body = tool.Get_language(db, "data_missing", true)
	}

	return list_extra_page(db, config, tool.Get_language(db, "monthly_bbs_contributor", true), body)
}
