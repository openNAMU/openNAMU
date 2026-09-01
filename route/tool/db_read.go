package tool

import "database/sql"

func Get_setting_value(db *sql.DB, name string, coverage string, default_value string) string {
	data := Get_setting(db, name, coverage)
	if len(data) == 0 || len(data[0]) == 0 {
		return default_value
	}
	return data[0][0]
}

func Get_setting_value_exact(db *sql.DB, name string, coverage string, default_value string) string {
	data := ""
	if !QueryRow_DB(
		db,
		"select data from other where name = ? and coverage = ?",
		[]any{&data},
		name,
		coverage,
	) {
		return default_value
	}
	return data
}

func Get_document_setting_value(db *sql.DB, doc_name string, set_name string, doc_rev string) string {
	data := Get_document_setting(db, doc_name, set_name, doc_rev)
	if len(data) == 0 || len(data[0]) == 0 {
		return ""
	}
	return data[0][0]
}

func Get_document_setting_value_exact(db *sql.DB, doc_name string, set_name string, doc_rev string) string {
	data := ""
	QueryRow_DB(
		db,
		"select set_data from data_set where doc_name = ? and doc_rev = ? and set_name = ? limit 1",
		[]any{&data},
		doc_name,
		doc_rev,
		set_name,
	)
	return data
}

func Get_user_set_data(db *sql.DB, id string, name string) string {
	data := ""
	QueryRow_DB(
		db,
		"select data from user_set where id = ? and name = ?",
		[]any{&data},
		id,
		name,
	)
	return data
}

func Get_user_set_id(db *sql.DB, name string, data string) (string, bool) {
	id := ""
	exists := QueryRow_DB(
		db,
		"select id from user_set where name = ? and data = ?",
		[]any{&id},
		name,
		data,
	)
	return id, exists
}

func Get_user_set_exists(db *sql.DB, id string, name string) bool {
	value := ""
	return QueryRow_DB(
		db,
		"select id from user_set where id = ? and name = ?",
		[]any{&value},
		id,
		name,
	)
}

func Get_user_notice_unread_count(db *sql.DB, id string) string {
	count := "0"
	QueryRow_DB(
		db,
		"select count(*) from user_notice where name = ? and readme = ''",
		[]any{&count},
		id,
	)
	return count
}

func Get_bbs_set_data(db *sql.DB, set_id string, set_name string) string {
	data := ""
	QueryRow_DB(
		db,
		"select set_data from bbs_set where set_id = ? and set_name = ? order by set_code + 0 desc limit 1",
		[]any{&data},
		set_id,
		set_name,
	)
	return data
}

func Get_bbs_set_first_data(db *sql.DB, set_id string, set_name string) string {
	data := ""
	QueryRow_DB(
		db,
		"select set_data from bbs_set where set_name = ? and set_id = ?",
		[]any{&data},
		set_name,
		set_id,
	)
	return data
}

func Get_bbs_data_value(db *sql.DB, set_id string, set_code string, set_name string) (string, bool) {
	data := ""
	exists := QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = ? and set_id = ? and set_code = ? limit 1",
		[]any{&data},
		set_name,
		set_id,
		set_code,
	)
	return data, exists
}

func Get_data_title(db *sql.DB, title string) (string, bool) {
	data := ""
	exists := QueryRow_DB(
		db,
		"select title from data where title = ?",
		[]any{&data},
		title,
	)
	return data, exists
}

func Get_history_content(db *sql.DB, title string, revision string) (string, string, bool) {
	data := ""
	hide := ""
	exists := QueryRow_DB(
		db,
		"select data, hide from history where title = ? and id = ?",
		[]any{&data, &hide},
		title,
		revision,
	)
	return data, hide, exists
}

func Get_history_hide(db *sql.DB, title string, revision string) string {
	data := ""
	QueryRow_DB(
		db,
		"select hide from history where title = ? and id = ?",
		[]any{&data},
		title,
		revision,
	)
	return data
}

func Get_history_send(db *sql.DB, title string, revision string) string {
	data := ""
	QueryRow_DB(
		db,
		"select send from history where title = ? and id = ?",
		[]any{&data},
		title,
		revision,
	)
	return data
}

func Get_history_exists(db *sql.DB, title string) bool {
	value := ""
	return QueryRow_DB(
		db,
		"select title from history where title = ? limit 1",
		[]any{&value},
		title,
	)
}

func Get_history_last_revision(db *sql.DB, title string) string {
	value := ""
	QueryRow_DB(
		db,
		"select id from history where title = ? order by date desc limit 1",
		[]any{&value},
		title,
	)
	return value
}

func Get_rd_data(db *sql.DB, code string) (map[string]string, bool) {
	title := ""
	sub := ""
	stop := ""
	agree := ""
	acl := ""
	exists := QueryRow_DB(
		db,
		"select title, sub, stop, agree, acl from rd where code = ?",
		[]any{&title, &sub, &stop, &agree, &acl},
		code,
	)
	return map[string]string{
		"title": title,
		"sub":   sub,
		"stop":  stop,
		"agree": agree,
		"acl":   acl,
	}, exists
}

func Get_acl_exists(db *sql.DB, title string) bool {
	value := ""
	return QueryRow_DB(
		db,
		"select title from acl where title = ? limit 1",
		[]any{&value},
		title,
	)
}

func Get_back_link(db *sql.DB, title string, link_type string) (string, bool) {
	value := ""
	exists := QueryRow_DB(
		db,
		"select link from back where title = ? and type = ? limit 1",
		[]any{&value},
		title,
		link_type,
	)
	return value, exists
}

func Get_back_title(db *sql.DB, link string, link_type string) (string, bool) {
	value := ""
	exists := QueryRow_DB(
		db,
		"select title from back where link = ? and type = ? limit 1",
		[]any{&value},
		link,
		link_type,
	)
	return value, exists
}

func Get_html_filter_value(db *sql.DB, html string, kind string) []string {
	data := []string{"", "", ""}
	QueryRow_DB(
		db,
		"select html, plus, plus_t from html_filter where html = ? and kind = ? limit 1",
		[]any{&data[0], &data[1], &data[2]},
		html,
		kind,
	)
	return data
}

func Get_back_redirect(db *sql.DB, value string) bool {
	data := ""
	return QueryRow_DB(
		db,
		"select link from back where (title = ? or link = ?) and type = 'redirect' limit 1",
		[]any{&data},
		value,
		value,
	)
}

func Get_topic_exists(db *sql.DB, code string, id string) bool {
	value := ""
	return QueryRow_DB(
		db,
		"select id from topic where code = ? and id = ?",
		[]any{&value},
		code,
		id,
	)
}

func Get_topic_block(db *sql.DB, code string, id string) string {
	value := ""
	QueryRow_DB(
		db,
		"select block from topic where code = ? and id = ?",
		[]any{&value},
		code,
		id,
	)
	return value
}

func Get_topic_top(db *sql.DB, code string, id string) string {
	value := ""
	QueryRow_DB(
		db,
		"select top from topic where code = ? and id = ?",
		[]any{&value},
		code,
		id,
	)
	return value
}

func Get_vote_data(db *sql.DB, id string) (map[string]string, bool) {
	name := ""
	subject := ""
	data := ""
	type_data := ""
	exists := QueryRow_DB(
		db,
		"select name, subject, data, type from vote where id = ? and user = ''",
		[]any{&name, &subject, &data, &type_data},
		id,
	)
	return map[string]string{
		"name":    name,
		"subject": subject,
		"data":    data,
		"type":    type_data,
	}, exists
}

func Get_vote_value(db *sql.DB, id string, name string) string {
	data := ""
	QueryRow_DB(
		db,
		"select data from vote where id = ? and name = ? and type = 'option'",
		[]any{&data},
		id,
		name,
	)
	return data
}

func Get_vote_user_exists(db *sql.DB, id string, user string) bool {
	value := ""
	return QueryRow_DB(
		db,
		"select user from vote where id = ? and user = ?",
		[]any{&value},
		id,
		user,
	)
}

func Get_user_set_data_match(db *sql.DB, id string, name string, data string) bool {
	value := ""
	return QueryRow_DB(
		db,
		"select data from user_set where id = ? and name = ? and data = ? limit 1",
		[]any{&value},
		id,
		name,
		data,
	)
}

func Get_rd_active_title(db *sql.DB, title string) bool {
	value := ""
	return QueryRow_DB(
		db,
		"select title from rd where title = ? and not stop = 'O' order by date desc limit 1",
		[]any{&value},
		title,
	)
}

func Get_topic_set_data(db *sql.DB, thread_code string, set_name string) string {
	data := ""
	QueryRow_DB(
		db,
		"select set_data from topic_set where thread_code = ? and set_name = ?",
		[]any{&data},
		thread_code,
		set_name,
	)
	return data
}

func Get_back_redirect_data(db *sql.DB, doc_name string) (string, string, bool) {
	target := ""
	anchor := ""
	exists := QueryRow_DB(
		db,
		"select title, data from back where link = ? and type = 'redirect' limit 1",
		[]any{&target, &anchor},
		doc_name,
	)
	return target, anchor, exists
}

func Get_data_title_like(db *sql.DB, title string) bool {
	value := ""
	return QueryRow_DB(
		db,
		"select title from data where title like ? limit 1",
		[]any{&value},
		title,
	)
}

func Get_back_data(db *sql.DB, title string, link string, link_type string) (string, bool) {
	data := ""
	exists := QueryRow_DB(
		db,
		"select data from back where title = ? and link = ? and type = ? limit 1",
		[]any{&data},
		title,
		link,
		link_type,
	)
	return data, exists
}

func Get_history_date(db *sql.DB, title string) string {
	data := ""
	QueryRow_DB(
		db,
		"select date from history where title = ? order by id + 0 desc limit 1",
		[]any{&data},
		title,
	)
	return data
}

func Get_topic_data(db *sql.DB, code string, id string) (map[string]string, bool) {
	comment_id := ""
	data := ""
	date := ""
	ip := ""
	block := ""
	top := ""
	exists := QueryRow_DB(
		db,
		"select id, data, date, ip, block, top from topic where code = ? and id = ?",
		[]any{&comment_id, &data, &date, &ip, &block, &top},
		code,
		id,
	)
	return map[string]string{
		"id":    comment_id,
		"data":  data,
		"date":  date,
		"ip":    ip,
		"block": block,
		"top":   top,
	}, exists
}

func Get_topic_rows(db *sql.DB, code string) *sql.Rows {
	return Query_DB(
		db,
		"select id, data, date, ip, block, top from topic where code = ? order by id + 0 asc",
		code,
	)
}

func Get_vote_count(db *sql.DB, id string, data string) string {
	count := "0"
	QueryRow_DB(
		db,
		"select count(*) from vote where id = ? and user != '' and data = ?",
		[]any{&count},
		id,
		data,
	)
	return count
}

func Get_vote_users(db *sql.DB, id string, data string) *sql.Rows {
	return Query_DB(
		db,
		"select user from vote where id = ? and user != '' and data = ?",
		id,
		data,
	)
}

func Get_html_filter_rows(db *sql.DB, kind string) *sql.Rows {
	return Query_DB(db, "select html, plus, plus_t from html_filter where kind = ?", kind)
}

func Get_html_filter_html_rows(db *sql.DB, kind string) *sql.Rows {
	return Query_DB(db, "select html from html_filter where kind = ?", kind)
}

func Get_html_filter_plus_rows(db *sql.DB, kind string) *sql.Rows {
	return Query_DB(db, "select plus from html_filter where kind = ? and plus != ''", kind)
}

func Get_html_filter_inter_wiki_sub(db *sql.DB, name string) string {
	value := ""
	QueryRow_DB(
		db,
		"select plus_t from html_filter where html = ? and kind = 'inter_wiki_sub'",
		[]any{&value},
		name,
	)
	return value
}

func Get_application_rows(db *sql.DB) *sql.Rows {
	return Query_DB(db, "select id, data from user_set where name = 'application'")
}

func Get_user_notice_rows(db *sql.DB, name string, offset int) *sql.Rows {
	return Query_DB(
		db,
		"select id, data, date, readme from user_notice where name = ? order by date desc limit ?, 50",
		name,
		offset,
	)
}

func Get_re_admin_rows(db *sql.DB, limit int) *sql.Rows {
	return Query_DB(db, "select who, what from re_admin order by time desc limit ?", limit)
}

func Get_re_admin_page_rows(db *sql.DB, search string, offset int) *sql.Rows {
	if search == "" {
		return Query_DB(db, "select who, what, time from re_admin order by time desc limit ?, 50", offset)
	}
	return Query_DB(
		db,
		"select who, what, time from re_admin where what like ? order by time desc limit ?, 50",
		search+"%",
		offset,
	)
}

func Get_re_admin_last_time(db *sql.DB, value string) string {
	data := ""
	QueryRow_DB(
		db,
		"select time from re_admin where what like ? order by time desc limit 1",
		[]any{&data},
		value,
	)
	return data
}

func Get_document_acl_rows(db *sql.DB, offset int) *sql.Rows {
	return Query_DB(
		db,
		"select distinct title, data, type from acl where data != '' and title not like 'user:%' order by title desc limit ?, 50",
		offset,
	)
}

func Get_acl_why(db *sql.DB, title string) string {
	data := ""
	QueryRow_DB(
		db,
		"select data from acl where title = ? and type = 'why' limit 1",
		[]any{&data},
		title,
	)
	return data
}

func Get_need_document_rows(db *sql.DB, offset int) *sql.Rows {
	return Query_DB(
		db,
		"select b.title, count(*) from back b where b.type = 'no' and not exists (select 1 from data d where d.title = b.title) group by b.title order by count(*) desc, b.title asc limit ?, 50",
		offset,
	)
}

func Get_data_rows(db *sql.DB, offset int) *sql.Rows {
	return Query_DB(db, "select title from data order by title asc limit ?, 50", offset)
}

func Get_data_prefix_count(db *sql.DB, prefix string) string {
	count := "0"
	QueryRow_DB(
		db,
		"select count(*) from data where title like ?",
		[]any{&count},
		prefix+"%",
	)
	return count
}

func Get_data_file_rows(db *sql.DB, offset int, order bool) *sql.Rows {
	if order {
		return Query_DB(db, "select title from data where title like 'file:%' order by title limit ?, 50", offset)
	}
	return Query_DB(db, "select title from data where title like 'file:%' limit ?, 50", offset)
}

func Get_file_rows(db *sql.DB) *sql.Rows {
	return Query_DB(db, "select title, type from data where title like 'file:%' order by title")
}

func Get_user_date_rows(db *sql.DB, offset int, limit bool) *sql.Rows {
	if limit {
		return Query_DB(db, "select id, data from user_set where name = 'date' order by data desc limit ?, 50", offset)
	}
	return Query_DB(db, "select id, data from user_set where name = 'date' order by data desc")
}

func Get_no_link_rows(db *sql.DB) *sql.Rows {
	return Query_DB(db, "select doc_name, set_data from data_set where set_name = 'link_count' and set_data = '0' order by doc_name")
}

func Get_no_link_page_rows(db *sql.DB, offset int) *sql.Rows {
	return Query_DB(db, "select doc_name, set_data from data_set where set_name = 'link_count' and set_data = '0' order by doc_name limit ?, 50", offset)
}

func Get_move_document_rows(db *sql.DB, pattern string) *sql.Rows {
	return Query_DB(
		db,
		"select title from data where title not like 'file:%' and title not like 'category:%' and title like ? order by title limit 100",
		pattern,
	)
}

func Get_bbs_record_rows(db *sql.DB, user_name string, offset int) *sql.Rows {
	return Query_DB(
		db,
		`select d.set_code, d.set_id, d.set_data
		 from bbs_data d
		 where d.set_name = 'date'
		 and exists (
			 select 1 from bbs_data u
			 where u.set_name = 'user_id'
			 and u.set_id = d.set_id
			 and u.set_code = d.set_code
			 and u.set_data = ?
		 )
		 order by d.set_data desc limit ?, 50`,
		user_name,
		offset,
	)
}

func Get_bbs_comment_record_rows(db *sql.DB, user_name string, offset int) *sql.Rows {
	return Query_DB(
		db,
		`select d.set_id, d.set_code, d.set_data
		 from bbs_data d
		 where d.set_name = 'comment_date'
		 and exists (
			 select 1 from bbs_data u
			 where u.set_name = 'comment_user_id'
			 and u.set_id = d.set_id
			 and u.set_code = d.set_code
			 and u.set_data = ?
		 )
		 order by d.set_data desc limit ?, 50`,
		user_name,
		offset,
	)
}

func Get_bbs_last_comment_date(db *sql.DB, root_id string) string {
	date := ""
	QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'comment_date' and (set_id = ? or set_id like ?) order by set_data desc limit 1",
		[]any{&date},
		root_id,
		root_id+"-%",
	)
	return date
}

func Get_topic_record_rows(db *sql.DB, user_name string, offset int, limit bool) *sql.Rows {
	if limit {
		return Query_DB(db, "select code, id, date from topic where ip = ? order by date desc limit ?, 50", user_name, offset)
	}
	return Query_DB(db, "select code, data, date from topic where ip = ? order by date desc limit 100", user_name)
}

func Get_history_record_rows(db *sql.DB, user_name string, record_type string, offset int, limit bool) *sql.Rows {
	if !limit {
		return Query_DB(db, "select title, date, send from history where ip = ? order by id + 0 desc limit 100", user_name)
	}
	if record_type == "" || record_type == "normal" || record_type == "edit" {
		return Query_DB(db, "select id, title, date, ip, send, leng, hide, type from history where ip = ? order by date desc limit ?, 50", user_name, offset)
	}
	return Query_DB(db, "select id, title, date, ip, send, leng, hide, type from history where ip = ? and type = ? order by date desc limit ?, 50", user_name, record_type, offset)
}

func Get_history_count(db *sql.DB, user_name string) string {
	count := "0"
	QueryRow_DB(db, "select count(*) from history where ip = ?", []any{&count}, user_name)
	return count
}

func Get_topic_count(db *sql.DB, user_name string) string {
	count := "0"
	QueryRow_DB(db, "select count(*) from topic where ip = ?", []any{&count}, user_name)
	return count
}

func Get_bbs_comment_count(db *sql.DB, user_name string) string {
	count := "0"
	QueryRow_DB(
		db,
		"select count(*) from bbs_data where set_name = 'comment_user_id' and set_data = ?",
		[]any{&count},
		user_name,
	)
	return count
}

func Get_history_length_rows(db *sql.DB, day string, user_name string) *sql.Rows {
	return Query_DB(db, "select leng from history where date like ? and ip = ?", day+"%", user_name)
}

func Get_ua_simple_rows(db *sql.DB, name string, name_is_ip bool, offset int) *sql.Rows {
	name_column := "name"
	other_column := "ip"
	if name_is_ip {
		name_column = "ip"
		other_column = "name"
	}
	return Query_DB(
		db,
		"select distinct "+other_column+" from ua_d where "+name_column+" = ? order by today desc limit ?, 50",
		name,
		offset,
	)
}

func Get_ua_rows(db *sql.DB, name string, plus_name string, name_is_ip bool, plus_is_ip bool, offset int) *sql.Rows {
	name_column := "name"
	plus_column := "name"
	if name_is_ip {
		name_column = "ip"
	}
	if plus_is_ip {
		plus_column = "ip"
	}
	query := "select name, ip, ua, today from ua_d where " + name_column + " = ?"
	args := []any{name}
	if plus_name != "" {
		query += " or " + plus_column + " = ?"
		args = append(args, plus_name)
	}
	query += " order by today desc limit ?, 50"
	args = append(args, offset)
	return Query_DB(db, query, args...)
}

func Get_ua_distinct_ips(db *sql.DB, name string, plus_name string, name_is_ip bool, plus_is_ip bool) []string {
	name_column := "name"
	plus_column := "name"
	if name_is_ip {
		name_column = "ip"
	}
	if plus_is_ip {
		plus_column = "ip"
	}
	query := "select distinct ip from ua_d where " + name_column + " = ?"
	args := []any{name}
	if plus_name != "" {
		query += " or " + plus_column + " = ?"
		args = append(args, plus_name)
	}
	rows := Query_DB(db, query, args...)
	defer rows.Close()
	values := []string{}
	for rows.Next() {
		value := ""
		if rows.Scan(&value) == nil {
			values = append(values, value)
		}
	}
	return values
}

func Get_user_approval_data(db *sql.DB, id string) (string, string) {
	question := Get_user_set_data(db, id, "approval_question")
	answer := Get_user_set_data(db, id, "approval_question_answer")
	return question, answer
}

func Get_other_data(db *sql.DB, name string) string {
	data := ""
	QueryRow_DB(db, "select data from other where name = ?", []any{&data}, name)
	return data
}

func Get_category_rows(db *sql.DB, doc_name string) *sql.Rows {
	return Query_DB(db, "select link, data from back where title = ? and type = 'cat' order by link", doc_name)
}

func Get_category_back_rows(db *sql.DB, doc_name string) *sql.Rows {
	return Query_DB(db, "select title, data from back where link = ? and (type = 'cat' or type = '') order by title", doc_name)
}

func Get_category_meta(db *sql.DB, doc_name string, category_name string, meta_type string) (string, bool) {
	data := ""
	exists := QueryRow_DB(
		db,
		"select data from back where title = ? and link = ? and type = ? limit 1",
		[]any{&data},
		doc_name,
		category_name,
		meta_type,
	)
	return data, exists
}

func Get_rb_end(db *sql.DB, block string, band string) string {
	end_date := ""
	QueryRow_DB(
		db,
		"select end from rb where block = ? and band = ? and ongoing = '1' order by today desc limit 1",
		[]any{&end_date},
		block,
		band,
	)
	return end_date
}

func Get_file_license_rows(db *sql.DB) *sql.Rows {
	return Get_html_filter_html_rows(db, "image_license")
}

func Get_email_filter_rows(db *sql.DB) *sql.Rows {
	return Get_html_filter_html_rows(db, "email")
}
