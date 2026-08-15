package route

import (
	"database/sql"

	"opennamu/route/tool"
)

// Api_thread_bbs exposes the legacy topic table in the BBS comment shape.
func Api_thread_bbs(config tool.Config, tool_name string, topic_num string, s_num string, e_num string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", topic_num, "topic_view", config.IP) {
		return map[string]any{
			"response": "require auth",
			"data":     []map[string]string{},
		}
	}

	if tool_name == "length" {
		length := "0"
		tool.QueryRow_DB(
			db,
			"select count(*) from topic where code = ?",
			[]any{&length},
			topic_num,
		)

		return map[string]any{
			"response": "ok",
			"comment":  length,
			"reply":    "0",
			"data":     tool.Str_to_int(length),
		}
	}

	rows := thread_bbs_query(db, topic_num, tool_name, s_num, e_num)
	defer rows.Close()

	admin_auth := tool.Check_acl(db, "", "", "toron_auth", config.IP)
	ip_parser_temp := map[string][]string{}
	data_list := []map[string]string{}

	for rows.Next() {
		var id string
		var data string
		var date string
		var ip string
		var block string
		var top string

		if err := rows.Scan(&id, &data, &date, &ip, &block, &top); err != nil {
			panic(err)
		}

		if block == "O" && !admin_auth {
			data = ""
		}

		ip_pre := ""
		ip_render := ""
		if cached, ok := ip_parser_temp[ip]; ok {
			ip_pre = cached[0]
			ip_render = cached[1]
		} else {
			ip_pre = tool.IP_preprocess(db, ip, config.IP)[0]
			ip_render = tool.IP_parser(db, ip, config.IP)
			ip_parser_temp[ip] = []string{ip_pre, ip_render}
		}

		data_list = append(data_list, map[string]string{
			"id":                     topic_num,
			"code":                   id,
			"comment":                data,
			"comment_date":           date,
			"comment_user_id":        ip_pre,
			"comment_user_id_render": ip_render,
			"blind":                  block,
			"top":                    top,
		})
	}

	return map[string]any{
		"response": "ok",
		"data":     data_list,
	}
}

func thread_bbs_query(db *sql.DB, topic_num string, tool_name string, s_num string, e_num string) *sql.Rows {
	if tool_name == "top" {
		return tool.Query_DB(
			db,
			"select id, data, date, ip, block, top from topic where code = ? and top = 'O' order by id + 0 asc",
			topic_num,
		)
	}

	if s_num != "" && e_num != "" {
		return tool.Query_DB(
			db,
			"select id, data, date, ip, block, top from topic where code = ? and ? + 0 <= id + 0 and id + 0 <= ? + 0 order by id + 0 asc",
			topic_num,
			s_num,
			e_num,
		)
	}

	return tool.Query_DB(
		db,
		"select id, data, date, ip, block, top from topic where code = ? order by id + 0 asc",
		topic_num,
	)
}
