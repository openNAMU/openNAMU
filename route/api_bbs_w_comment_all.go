package route

import (
	"database/sql"

	"opennamu/route/tool"
)

type bbs_comment_page_state struct {
	offset int
	limit  int
	data   []map[string]string
}

func Api_bbs_w_comment_page_data(config tool.Config, sub_code string, already_auth_check bool, do_type string, state *bbs_comment_page_state) {
	if len(state.data) >= state.limit {
		return
	}

	db := tool.DB_connect()
	defer tool.DB_close(db)
	Bbs_comment_page_data(db, config, sub_code, already_auth_check, do_type, state)
}

func Bbs_comment_page_data(db *sql.DB, config tool.Config, sub_code string, already_auth_check bool, do_type string, state *bbs_comment_page_state) {
	if len(state.data) >= state.limit {
		return
	}
	if !already_auth_check && Bbs_comment_read_access(db, config, sub_code) != "" {
		return
	}

	for _, comment := range Bbs_comment_one_data(db, config, do_type, sub_code) {
		if state.offset > 0 {
			state.offset--
		} else {
			state.data = append(state.data, comment)
		}

		if len(state.data) >= state.limit {
			return
		}
		Bbs_comment_page_data(db, config, sub_code+"-"+comment["code"], already_auth_check, do_type, state)
		if len(state.data) >= state.limit {
			return
		}
	}
}

func Api_bbs_w_comment_page(config tool.Config, sub_code string, already_auth_check bool, do_type string, page int) map[string]any {
	if page < 1 {
		page = 1
	}

	state := &bbs_comment_page_state{
		offset: (page - 1) * 50,
		limit:  50,
		data:   []map[string]string{},
	}
	Api_bbs_w_comment_page_data(config, sub_code, already_auth_check, do_type, state)

	return map[string]any{
		"response": "ok",
		"data":     state.data,
	}
}

func Api_bbs_w_comment_pinned_data(config tool.Config, set_id string, set_code string) []map[string]string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	root := set_id + "-" + set_code
	rows := tool.Query_DB(
		db,
		"select set_id, set_code from bbs_data where set_name = 'pinned' and (set_id = ? or set_id like ?) order by set_id, set_code",
		root,
		root+"-%",
	)
	defer rows.Close()

	data := []map[string]string{}
	for rows.Next() {
		comment_set_id := ""
		comment_set_code := ""
		if rows.Scan(&comment_set_id, &comment_set_code) != nil {
			continue
		}
		post_set_id, post_set_code, _, valid := Bbs_comment_storage_location(comment_set_id, comment_set_code)
		if !valid || post_set_id != set_id || post_set_code != set_code {
			continue
		}

		comment_list := Bbs_comment_one_data(db, config, "", comment_set_id+"-"+comment_set_code)
		if len(comment_list) > 0 {
			data = append(data, comment_list[0])
		}
	}

	return data
}

func Bbs_comment_all_data(db *sql.DB, config tool.Config, sub_code string, already_auth_check bool, do_type string) []map[string]string {
	data := []map[string]string{}
	if !already_auth_check && Bbs_comment_read_access(db, config, sub_code) != "" {
		return data
	}

	for _, comment := range Bbs_comment_one_data(db, config, do_type, sub_code) {
		data = append(data, comment)
		data = append(data, Bbs_comment_all_data(db, config, sub_code+"-"+comment["code"], already_auth_check, do_type)...)
	}
	return data
}

func Api_bbs_w_comment_all(config tool.Config, sub_code string, already_auth_check bool, do_type string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return map[string]any{
		"response": "ok",
		"data":     Bbs_comment_all_data(db, config, sub_code, already_auth_check, do_type),
	}
}
