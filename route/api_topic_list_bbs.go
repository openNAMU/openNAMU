package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func api_topic_list_bbs(config tool.Config, num string, doc_name string, do_type string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, doc_name, "", "render", config.IP) || !tool.Check_acl(db, thread_bbs_id, "", "bbs_view", config.IP) {
		return map[string]any{"response": "require auth", "data": [][]string{}}
	}

	page := tool.Str_to_int(num)
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * 50

	condition := ""
	switch do_type {
	case "close":
		condition = " and coalesce((select set_data from bbs_data prefix_data where prefix_data.set_name = 'prefix' and prefix_data.set_id = document_data.set_id and prefix_data.set_code = document_data.set_code limit 1), '') = '닫힘'"
	case "agree":
		condition = " and coalesce((select set_data from bbs_data prefix_data where prefix_data.set_name = 'prefix' and prefix_data.set_id = document_data.set_id and prefix_data.set_code = document_data.set_code limit 1), '') = '합의'"
	default:
		condition = " and coalesce((select set_data from bbs_data prefix_data where prefix_data.set_name = 'prefix' and prefix_data.set_id = document_data.set_id and prefix_data.set_code = document_data.set_code limit 1), '') != '닫힘'"
	}

	rows := tool.Query_DB(
		db,
		"select document_data.set_code, coalesce((select set_data from bbs_data title_data where title_data.set_name = 'title' and title_data.set_id = document_data.set_id and title_data.set_code = document_data.set_code limit 1), ''), case when coalesce((select set_data from bbs_data prefix_data where prefix_data.set_name = 'prefix' and prefix_data.set_id = document_data.set_id and prefix_data.set_code = document_data.set_code limit 1), '') = '닫힘' then 'O' else '' end, case when coalesce((select set_data from bbs_data agree_prefix_data where agree_prefix_data.set_name = 'prefix' and agree_prefix_data.set_id = document_data.set_id and agree_prefix_data.set_code = document_data.set_code limit 1), '') = '합의' then 'O' else '' end, coalesce((select set_data from bbs_data date_data where date_data.set_name = 'date' and date_data.set_id = document_data.set_id and date_data.set_code = document_data.set_code limit 1), '') as topic_date from bbs_data document_data where document_data.set_id = ? and document_data.set_name = 'document' and document_data.set_data = ?"+condition+" order by topic_date desc limit ?, 50",
		thread_bbs_id,
		doc_name,
		offset,
	)
	defer rows.Close()

	data_list := [][]string{}
	for rows.Next() {
		var code string
		var sub string
		var stop string
		var agree string
		var date string
		if rows.Scan(&code, &sub, &stop, &agree, &date) != nil {
			continue
		}

		comment_set_id := thread_bbs_id + "-" + code
		ip := ""
		id := ""
		tool.QueryRow_DB(db, "select set_data from bbs_data where set_name = 'comment_user_id' and set_id = ? order by set_code + 0 desc limit 1", []any{&ip}, comment_set_id)
		tool.QueryRow_DB(db, "select set_code from bbs_data where set_name = 'comment' and set_id = ? order by set_code + 0 desc limit 1", []any{&id}, comment_set_id)

		ip_pre := ""
		ip_render := ""
		if ip != "" {
			ip_pre = tool.IP_preprocess(db, ip, config.IP)[0]
			ip_render = tool.IP_parser(db, ip, config.IP)
		}
		data_list = append(data_list, []string{code, sub, stop, agree, ip_pre, ip_render, date, id})
	}

	return map[string]any{"response": "ok", "data": data_list}
}

func thread_bbs_document_exists(db *sql.DB, doc_name string) bool {
	value := ""
	return tool.QueryRow_DB(
		db,
		"select set_code from bbs_data where set_id = ? and set_name = 'document' and set_data = ? limit 1",
		[]any{&value},
		thread_bbs_id,
		doc_name,
	)

}
