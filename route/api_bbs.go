package route

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"opennamu/route/tool"
)

const report_bbs_id = "-2"

func Bbs_is_special_board(set_id string) bool {
	return tool.Str_to_int(set_id) <= 0
}

func Bbs_post_blind(db *sql.DB, set_id string, set_code string) bool {
	blind, exists := tool.Get_bbs_data_value(db, set_id, set_code, "blind")
	return exists && blind == "O"
}

func Bbs_post_blind_allowed(db *sql.DB, set_id string, set_code string, ip string, auth_info map[string]bool) bool {
	if !Bbs_post_blind(db, set_id, set_code) {
		return true
	}
	if auth_info == nil {
		auth_info = tool.Get_auth_info(db, ip)
	}
	return auth_info["bbs_post_manage"]
}

func Bbs_post_blind_sql(row_alias string) string {
	return "not exists (select 1 from bbs_data blind_data where blind_data.set_name = 'blind' and blind_data.set_data = 'O' and blind_data.set_id = " + row_alias + ".set_id and blind_data.set_code = " + row_alias + ".set_code)"
}

func Bbs_post_view_allowed(db *sql.DB, set_id string, set_code string, user_id string, ip string, auth_info map[string]bool) bool {
	if !Bbs_post_blind_allowed(db, set_id, set_code, ip, auth_info) {
		return false
	}

	acl_data := Bbs_set_value(db, set_id, "bbs_only_my_data_view_acl")
	if acl_data == "" || acl_data == "normal" || user_id == ip {
		return true
	}
	if auth_info == nil {
		auth_info = tool.Get_auth_info(db, ip)
	}
	if auth_info["bbs"] {
		return true
	}
	return tool.Check_acl_group(db, acl_data, auth_info)
}

func Bbs_post_view_sql(db *sql.DB, set_id string, ip string, row_alias string) (string, []any) {
	auth_info := tool.Get_auth_info(db, ip)
	blind_sql := ""
	if !auth_info["bbs_post_manage"] {
		blind_sql = Bbs_post_blind_sql(row_alias)
	}

	if set_id != "" {
		acl_data := Bbs_set_value(db, set_id, "bbs_only_my_data_view_acl")
		if acl_data == "" || acl_data == "normal" || auth_info["bbs"] || tool.Check_acl_group(db, acl_data, auth_info) {
			return blind_sql, nil
		}

		view_sql := "exists (select 1 from bbs_data author where author.set_name = 'user_id' and author.set_id = " + row_alias + ".set_id and author.set_code = " + row_alias + ".set_code and author.set_data = ?)"
		if blind_sql != "" {
			view_sql = blind_sql + " and " + view_sql
		}
		return view_sql, []any{ip}
	}

	rows := tool.Query_DB(
		db,
		"select set_id, set_data from bbs_set where set_name = 'bbs_only_my_data_view_acl' and set_code = ''",
	)
	defer rows.Close()

	allowed_set_id := []string{}
	private_count := 0
	for rows.Next() {
		setting_set_id := ""
		acl_data := ""
		if rows.Scan(&setting_set_id, &acl_data) != nil || acl_data == "" || acl_data == "normal" {
			continue
		}

		private_count += 1
		if auth_info["bbs"] || tool.Check_acl_group(db, acl_data, auth_info) {
			if !tool.Arr_in_str(allowed_set_id, setting_set_id) {
				allowed_set_id = append(allowed_set_id, setting_set_id)
			}
		}
	}

	if private_count == 0 {
		return blind_sql, nil
	}

	view_sql := "(not exists (select 1 from bbs_set only_view where only_view.set_name = 'bbs_only_my_data_view_acl' and only_view.set_code = '' and only_view.set_id = " + row_alias + ".set_id and only_view.set_data != '' and only_view.set_data != 'normal')"
	view_values := []any{}
	if len(allowed_set_id) > 0 {
		placeholders := strings.Repeat("?,", len(allowed_set_id))
		placeholders = strings.TrimSuffix(placeholders, ",")
		view_sql += " or " + row_alias + ".set_id in (" + placeholders + ")"
		for _, allowed_id := range allowed_set_id {
			view_values = append(view_values, allowed_id)
		}
	}
	view_sql += " or exists (select 1 from bbs_data author where author.set_name = 'user_id' and author.set_id = " + row_alias + ".set_id and author.set_code = " + row_alias + ".set_code and author.set_data = ?))"
	view_values = append(view_values, ip)
	if blind_sql != "" {
		view_sql = "(" + blind_sql + " and " + view_sql + ")"
	}

	return view_sql, view_values
}

type bbs_user_display_data struct {
	ip_pre    string
	ip_render string
}

type bbs_post_activity_update struct {
	set_id   string
	set_code string
	date     string
}

type bbs_filter struct {
	comment_min          int
	commented            int
	tabom_min            int
	mine                 bool
	participate          bool
	comment_user         string
	comment_user_invalid bool
	tabom_user           bool
	author               string
	author_invalid       bool
	prefix               string
	tag                  string
}

func Bbs_filter_number(data string) int {
	num := tool.Str_to_int(data)
	if num < 0 {
		return 0
	}

	return num
}

func Bbs_filter_parse(data string) bbs_filter {
	data = strings.Trim(data, "/")
	parts := strings.Split(data, "/")
	filter := bbs_filter{}

	for i := 0; i+1 < len(parts); {
		if parts[i] == "tag" {
			filter.tag = strings.TrimSpace(strings.Join(parts[i+1:], "/"))
			break
		}

		value := Bbs_filter_number(parts[i+1])
		switch parts[i] {
		case "comment":
			filter.comment_min = value
		case "commented":
			if parts[i+1] == "1" {
				filter.commented = 1
			} else if parts[i+1] == "0" {
				filter.commented = 2
			}
		case "tabom":
			filter.tabom_min = value
		case "mine":
			filter.mine = parts[i+1] == "1"
		case "participate":
			filter.participate = parts[i+1] == "1"
		case "comment_user":
			decoded, err := tool.Get_base64_decode(parts[i+1])
			if err == nil {
				filter.comment_user = decoded
			} else {
				filter.comment_user_invalid = true
			}
		case "tabom_user":
			filter.tabom_user = parts[i+1] == "1"
		case "user":
			decoded, err := tool.Get_base64_decode(parts[i+1])
			if err == nil {
				filter.author = decoded
			} else {
				filter.author_invalid = true
			}
		case "prefix":
			filter.prefix = strings.TrimSpace(parts[i+1])
		}
		i += 2
	}

	return filter
}

func Bbs_filter_path(filter bbs_filter) string {
	path := []string{}
	if filter.comment_min > 0 {
		path = append(path, "comment", strconv.Itoa(filter.comment_min))
	}
	if filter.commented == 1 {
		path = append(path, "commented", "1")
	} else if filter.commented == 2 {
		path = append(path, "commented", "0")
	}
	if filter.comment_user != "" {
		path = append(path, "comment_user", tool.Base64_encode(filter.comment_user))
	}
	if filter.tabom_min > 0 {
		path = append(path, "tabom", strconv.Itoa(filter.tabom_min))
	}
	if filter.mine {
		path = append(path, "mine", "1")
	}
	if filter.participate {
		path = append(path, "participate", "1")
	}
	if filter.tabom_user {
		path = append(path, "tabom_user", "1")
	}
	if filter.author != "" {
		path = append(path, "user", tool.Base64_encode(filter.author))
	}
	if filter.prefix != "" {
		path = append(path, "prefix", tool.Url_parser(filter.prefix))
	}
	if filter.tag != "" {
		path = append(path, "tag", tool.Url_parser(filter.tag))
	}

	return strings.Join(path, "/")
}

func Bbs_filter_path_data(data string) (string, string) {
	parts := strings.Split(strings.Trim(data, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		return "1", ""
	}

	page := parts[len(parts)-1]
	if tool.Str_to_int(page) < 1 {
		page = "1"
	}

	return page, strings.Join(parts[:len(parts)-1], "/")
}

func Bbs_post_comment_count_sql(row_alias string) string {
	return "coalesce((select set_data from bbs_data comment_count_data where comment_count_data.set_name = 'comment_count' and comment_count_data.set_id = " + row_alias + ".set_id and comment_count_data.set_code = " + row_alias + ".set_code limit 1), '0') + 0"
}

func Bbs_post_comment_count_update(tx *sql.Tx, set_id string, set_code string, change int) {
	if change == 0 {
		return
	}

	result, err := tx.Exec(
		tool.DB_change("update bbs_data set set_data = case when set_data + ? < 0 then 0 else set_data + ? end where set_name = 'comment_count' and set_id = ? and set_code = ?"),
		change,
		change,
		set_id,
		set_code,
	)
	if err != nil {
		panic(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}
	if rows == 0 && change > 0 {
		if _, err := tx.Exec(
			tool.DB_change("insert into bbs_data (set_name, set_id, set_code, set_data) values ('comment_count', ?, ?, ?)"),
			set_id,
			set_code,
			change,
		); err != nil {
			panic(err)
		}
	}
}

func Bbs_post_tabom_count_sql(row_alias string) string {
	return "coalesce((select set_data from bbs_data tabom_data where tabom_data.set_name = 'tabom_count' and tabom_data.set_id = " + row_alias + ".set_id and tabom_data.set_code = " + row_alias + ".set_code limit 1), '0') + 0"
}

func Bbs_post_tabom_down_count_sql(row_alias string) string {
	return "coalesce((select set_data from bbs_data tabom_down_data where tabom_down_data.set_name = 'tabom_down_count' and tabom_down_data.set_id = " + row_alias + ".set_id and tabom_down_data.set_code = " + row_alias + ".set_code limit 1), '0') + 0"
}

func Bbs_post_tabom_score_sql(row_alias string) string {
	return "(" + Bbs_post_tabom_count_sql(row_alias) + " - " + Bbs_post_tabom_down_count_sql(row_alias) + ")"
}

func Bbs_post_excellent_excluded_sql(row_alias string) string {
	return "not exists (select 1 from bbs_data excellent_exclusion where excellent_exclusion.set_name = 'excellent_excluded' and excellent_exclusion.set_data = '1' and excellent_exclusion.set_id = " + row_alias + ".set_id and excellent_exclusion.set_code = " + row_alias + ".set_code)"
}

func Bbs_excellent_min_from_values(board_value string, global_value string) int {
	for _, value := range []string{board_value, global_value} {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}

		min, err := strconv.Atoi(value)
		if err == nil && min >= 0 {
			return min
		}
	}

	return 5
}

func Bbs_excellent_min(db *sql.DB, set_id string) int {
	return Bbs_excellent_min_from_values(
		Bbs_set_value(db, set_id, "bbs_excellent_min"),
		tool.Get_setting_value_exact(db, "bbs_excellent_min", "", ""),
	)
}

func Bbs_global_excellent_thresholds(db *sql.DB, board_ids []string) map[string]int {
	global_value := tool.Get_setting_value_exact(db, "bbs_excellent_min", "", "")
	board_values := make(map[string]string, len(board_ids))

	rows := tool.Query_DB(
		db,
		"select set_id, set_data from bbs_set where set_name = 'bbs_excellent_min' order by set_code + 0 desc",
	)
	defer rows.Close()
	for rows.Next() {
		var set_id string
		var set_data string
		if err := rows.Scan(&set_id, &set_data); err != nil {
			panic(err)
		}
		if _, exists := board_values[set_id]; !exists {
			board_values[set_id] = set_data
		}
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

	thresholds := make(map[string]int, len(board_ids))
	for _, set_id := range board_ids {
		thresholds[set_id] = Bbs_excellent_min_from_values(board_values[set_id], global_value)
	}
	return thresholds
}

func Bbs_global_excellent_board_ids(db *sql.DB, ip string) []string {
	rows := tool.Query_DB(db, "select set_id from bbs_set where set_name = 'bbs_name'")
	defer rows.Close()

	board_ids := make([]string, 0)
	seen := make(map[string]bool)
	for rows.Next() {
		var set_id string
		if err := rows.Scan(&set_id); err != nil {
			panic(err)
		}
		board_number, err := strconv.ParseInt(set_id, 10, 64)
		if err != nil || board_number <= 0 || strconv.FormatInt(board_number, 10) != set_id || seen[set_id] || Bbs_is_special_board(set_id) {
			continue
		}
		seen[set_id] = true
		if tool.Check_acl(db, set_id, "", "bbs_view", ip) {
			board_ids = append(board_ids, set_id)
		}
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return board_ids
}

func Bbs_global_excellent_hydrate(db *sql.DB, config tool.Config, data_list []map[string]string, display_cache map[string]bbs_user_display_data) {
	if len(data_list) == 0 {
		return
	}

	posts_by_key := make(map[string]map[string]string, len(data_list))
	conditions := make([]string, 0, len(data_list))
	values := make([]any, 0, len(data_list)*2)
	for _, post_data := range data_list {
		set_id := post_data["set_id"]
		set_code := post_data["set_code"]
		posts_by_key[set_id+"\x00"+set_code] = post_data
		conditions = append(conditions, "(set_id = ? and set_code = ?)")
		values = append(values, set_id, set_code)
	}

	query := "select set_name, set_data, set_id, set_code from bbs_data where set_name in ('title', 'date', 'user_id', 'prefix', 'comment_count', 'tabom_count', 'tabom_down_count', 'view_count', 'tag', 'blind') and (" + strings.Join(conditions, " or ") + ")"
	rows := tool.Query_DB(db, query, values...)
	defer rows.Close()
	for rows.Next() {
		var set_name string
		var set_data string
		var set_id string
		var set_code string
		if err := rows.Scan(&set_name, &set_data, &set_id, &set_code); err != nil {
			panic(err)
		}
		post_data := posts_by_key[set_id+"\x00"+set_code]
		if post_data == nil {
			continue
		}

		if set_name == "user_id" {
			user_display, exists := display_cache[set_data]
			if !exists {
				user_display = bbs_user_display_data{
					ip_pre:    tool.IP_preprocess(db, set_data, config.IP)[0],
					ip_render: tool.Get_user_profile_image_ui(db, set_data) + tool.IP_parser(db, set_data, config.IP),
				}
				display_cache[set_data] = user_display
			}
			post_data["user_id"] = user_display.ip_pre
			post_data["user_id_render"] = user_display.ip_render
			continue
		}

		if set_name == "tag" {
			if post_data["tags"] != "" {
				post_data["tags"] += ", "
			}
			post_data["tags"] += set_data
			continue
		}
		post_data[set_name] = set_data
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
}

func Bbs_comment_set_id_sql(row_alias string, suffix string) string {
	if tool.Get_DB_type() == "mysql" {
		return "concat(" + row_alias + ".set_id, '-', " + row_alias + ".set_code, '" + suffix + "')"
	}

	return row_alias + ".set_id || '-' || " + row_alias + ".set_code || '" + suffix + "'"
}

func Bbs_post_last_activity_sql(row_alias string) string {
	comment_set_id := Bbs_comment_set_id_sql(row_alias, "")
	comment_set_id_nested := Bbs_comment_set_id_sql(row_alias, "-%")

	return "coalesce((select nullif(last_activity_data.set_data, '') from bbs_data last_activity_data where last_activity_data.set_name = 'last_activity' and last_activity_data.set_id = " + row_alias + ".set_id and last_activity_data.set_code = " + row_alias + ".set_code limit 1), (select nullif(max(comment_date_data.set_data), '') from bbs_data comment_date_data where comment_date_data.set_name = 'comment_date' and (comment_date_data.set_id = " + comment_set_id + " or comment_date_data.set_id like " + comment_set_id_nested + ")), (select nullif(date_data.set_data, '') from bbs_data date_data where date_data.set_name = 'date' and date_data.set_id = " + row_alias + ".set_id and date_data.set_code = " + row_alias + ".set_code limit 1))"
}

func Bbs_post_last_activity_missing_sql(row_alias string) string {
	return "case when exists (select 1 from bbs_data last_activity_check where last_activity_check.set_name = 'last_activity' and last_activity_check.set_id = " + row_alias + ".set_id and last_activity_check.set_code = " + row_alias + ".set_code and last_activity_check.set_data != '') then 0 else 1 end"
}

func Bbs_post_last_activity_update(tx *sql.Tx, set_id string, set_code string, date string) {
	if date == "" {
		return
	}

	result, err := tx.Exec(
		tool.DB_change("update bbs_data set set_data = case when set_data is null or set_data < ? then ? else set_data end where set_name = 'last_activity' and set_id = ? and set_code = ?"),
		date,
		date,
		set_id,
		set_code,
	)
	if err != nil {
		panic(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}
	if rows > 0 {
		return
	}
	last_activity := ""
	if tool.QueryRow_DB(
		tx,
		"select set_data from bbs_data where set_name = 'last_activity' and set_id = ? and set_code = ? limit 1",
		[]any{&last_activity},
		set_id,
		set_code,
	) {
		return
	}

	if _, err := tx.Exec(
		tool.DB_change("insert into bbs_data (set_name, set_id, set_code, set_data) values ('last_activity', ?, ?, ?)"),
		set_id,
		set_code,
		date,
	); err != nil {
		panic(err)
	}
}

func Bbs_post_last_activity_rebuild(tx *sql.Tx, set_id string, set_code string) {
	comment_set_id := set_id + "-" + set_code
	last_activity := ""
	tool.QueryRow_DB(
		tx,
		"select coalesce(max(set_data), '') from bbs_data where set_name = 'comment_date' and (set_id = ? or set_id like ?)",
		[]any{&last_activity},
		comment_set_id,
		comment_set_id+"-%",
	)
	if last_activity == "" {
		tool.QueryRow_DB(
			tx,
			"select coalesce(set_data, '') from bbs_data where set_name = 'date' and set_id = ? and set_code = ? limit 1",
			[]any{&last_activity},
			set_id,
			set_code,
		)
	}
	if last_activity == "" {
		if _, err := tx.Exec(tool.DB_change("delete from bbs_data where set_name = 'last_activity' and set_id = ? and set_code = ?"), set_id, set_code); err != nil {
			panic(err)
		}
		return
	}
	Bbs_post_last_activity_update(tx, set_id, set_code, last_activity)
}

func Bbs_filter_sql(filter bbs_filter, row_alias string, user_id string) (string, []any) {
	filter_sql := ""
	filter_values := []any{}

	if filter.comment_min > 0 {
		filter_sql += " and " + Bbs_post_comment_count_sql(row_alias) + " >= ?"
		filter_values = append(filter_values, filter.comment_min)
	}
	if filter.commented == 1 {
		filter_sql += " and " + Bbs_post_comment_count_sql(row_alias) + " >= 1"
	} else if filter.commented == 2 {
		filter_sql += " and " + Bbs_post_comment_count_sql(row_alias) + " = 0"
	}
	if filter.tabom_min > 0 {
		filter_sql += " and " + Bbs_post_tabom_count_sql(row_alias) + " >= ?"
		filter_values = append(filter_values, filter.tabom_min)
	}
	if filter.mine {
		filter_sql += " and exists (select 1 from bbs_data author_data where author_data.set_name = 'user_id' and author_data.set_id = " + row_alias + ".set_id and author_data.set_code = " + row_alias + ".set_code and author_data.set_data = ?)"
		filter_values = append(filter_values, user_id)
	}
	if filter.participate {
		comment_set_id := Bbs_comment_set_id_sql(row_alias, "")
		comment_set_id_nested := Bbs_comment_set_id_sql(row_alias, "-%")
		filter_sql += " and (exists (select 1 from bbs_data author_data where author_data.set_name = 'user_id' and author_data.set_id = " + row_alias + ".set_id and author_data.set_code = " + row_alias + ".set_code and author_data.set_data = ?) or exists (select 1 from bbs_data comment_user_data where comment_user_data.set_name = 'comment_user_id' and comment_user_data.set_data = ? and (comment_user_data.set_id = " + comment_set_id + " or comment_user_data.set_id like " + comment_set_id_nested + ")))"
		filter_values = append(filter_values, user_id, user_id)
	}
	if filter.comment_user_invalid || filter.author_invalid {
		filter_sql += " and 1 = 0"
	}
	if filter.comment_user != "" {
		comment_set_id := Bbs_comment_set_id_sql(row_alias, "")
		comment_set_id_nested := Bbs_comment_set_id_sql(row_alias, "-%")
		filter_sql += " and exists (select 1 from bbs_data comment_user_data where comment_user_data.set_name = 'comment_user_id' and comment_user_data.set_data = ? and (comment_user_data.set_id = " + comment_set_id + " or comment_user_data.set_id like " + comment_set_id_nested + "))"
		filter_values = append(filter_values, filter.comment_user)
	}
	if filter.tabom_user {
		filter_sql += " and exists (select 1 from bbs_data tabom_user_data where tabom_user_data.set_name = 'tabom_list' and tabom_user_data.set_id = " + row_alias + ".set_id and tabom_user_data.set_code = " + row_alias + ".set_code and tabom_user_data.set_data = ?)"
		filter_values = append(filter_values, user_id)
	}
	if filter.author != "" {
		filter_sql += " and exists (select 1 from bbs_data author_data where author_data.set_name = 'user_id' and author_data.set_id = " + row_alias + ".set_id and author_data.set_code = " + row_alias + ".set_code and author_data.set_data = ?)"
		filter_values = append(filter_values, filter.author)
	}
	if filter.prefix != "" {
		filter_sql += " and exists (select 1 from bbs_data prefix_data where prefix_data.set_name = 'prefix' and prefix_data.set_id = " + row_alias + ".set_id and prefix_data.set_code = " + row_alias + ".set_code and prefix_data.set_data = ?)"
		filter_values = append(filter_values, filter.prefix)
	}
	if filter.tag != "" {
		filter_sql += " and exists (select 1 from bbs_data tag_data where tag_data.set_name = 'tag' and tag_data.set_id = " + row_alias + ".set_id and tag_data.set_code = " + row_alias + ".set_code and tag_data.set_data = ?)"
		filter_values = append(filter_values, filter.tag)
	}

	return filter_sql, filter_values
}

func Api_bbs(config tool.Config, bbs_num string, page string, sort_type string) map[string]any {
	return Api_bbs_internal(config, bbs_num, page, sort_type, bbs_filter{})
}

func Api_bbs_filter(config tool.Config, bbs_num string, filter_data string) map[string]any {
	page, filter_path := Bbs_filter_path_data(filter_data)
	return Api_bbs_internal(config, bbs_num, page, "", Bbs_filter_parse(filter_path))
}

func Api_bbs_internal(config tool.Config, bbs_num string, page string, sort_type string, filter bbs_filter) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if bbs_num == "" && !tool.Check_permission(db, "bbs_main_view", config.IP) {
		return map[string]any{"response": "require auth", "data": []map[string]string{}, "has_next": false}
	}

	if bbs_num != "" && !tool.Check_acl(db, bbs_num, "", "bbs_view", config.IP) {
		return map[string]any{"response": "require auth", "data": []map[string]string{}, "has_next": false}
	}

	append_post_view_and_filter_sql := func(query string, row_alias string, values []any) (string, []any) {
		view_sql, view_values := Bbs_post_view_sql(db, bbs_num, config.IP, row_alias)
		filter_sql, filter_values := Bbs_filter_sql(filter, row_alias, config.IP)
		if view_sql != "" {
			query += " and " + view_sql
			values = append(values, view_values...)
		}
		query += filter_sql
		values = append(values, filter_values...)
		return query, values
	}

	global_excellent_listing := bbs_num == "" && sort_type == "excellent"
	var main_listing_rows *sql.Rows
	var pinned_listing_rows *sql.Rows
	var ordinary_listing_rows *sql.Rows
	defer func() {
		if main_listing_rows != nil {
			main_listing_rows.Close()
		}
		if pinned_listing_rows != nil {
			pinned_listing_rows.Close()
		}
		if ordinary_listing_rows != nil {
			ordinary_listing_rows.Close()
		}
	}()
	if bbs_num == "" {
		if global_excellent_listing {
			board_ids := Bbs_global_excellent_board_ids(db, config.IP)
			if len(board_ids) > 0 {
				thresholds := Bbs_global_excellent_thresholds(db, board_ids)
				score_sql := Bbs_post_tabom_score_sql("title")
				board_ids_sql := make([]string, 0, len(board_ids))
				threshold_cases := make([]string, 0, len(board_ids))
				for _, set_id := range board_ids {
					board_ids_sql = append(board_ids_sql, "'"+set_id+"'")
					threshold_cases = append(threshold_cases, "when '"+set_id+"' then "+strconv.Itoa(thresholds[set_id]))
				}

				query := "select distinct title.set_code, title.set_id, case when exists (select 1 from bbs_data pinned_data where pinned_data.set_name = 'pinned' and pinned_data.set_id = title.set_id and pinned_data.set_code = title.set_code) then '1' else '0' end, " + score_sql + " as excellent_score from bbs_data title where title.set_name = 'title' and title.set_id in (" + strings.Join(board_ids_sql, ",") + ") and " + score_sql + " > case title.set_id " + strings.Join(threshold_cases, " ") + " else 5 end"
				query, values := append_post_view_and_filter_sql(query, "title", nil)
				query += " and " + Bbs_post_excellent_excluded_sql("title")
				query += " order by excellent_score desc, title.set_id asc, title.set_code asc limit ?, 51"
				page_num := tool.Str_to_int(page)
				num := 0
				if page_num*50 > 0 {
					num = page_num*50 - 50
				}
				values = append(values, num)
				main_listing_rows = tool.Query_DB(db, query, values...)
			}
		} else {
			view_sql, view_values := Bbs_post_view_sql(db, bbs_num, config.IP, "bbs_data")
			query := "select set_code, set_id, '0' from bbs_data where set_name = 'date' and set_id + 0 > 0 and " + tool.Get_except_set_id_SQL()
			if view_sql != "" {
				query += " and " + view_sql
			}
			query += " order by set_data desc limit 50"
			main_listing_rows = tool.Query_DB(db, query, view_values...)
		}
	} else {
		page := tool.Str_to_int(page)
		num := 0
		if page*50 > 0 {
			num = page*50 - 50
		}

		query := "select set_code, set_id, '1'"
		if sort_type == "activity" {
			query += ", " + Bbs_post_last_activity_sql("bbs_data") + " as activity_date, " + Bbs_post_last_activity_missing_sql("bbs_data") + " as activity_cache_missing"
		}
		query += " from bbs_data where set_name = 'pinned' and set_id like ?"
		query, values := append_post_view_and_filter_sql(query, "bbs_data", []any{bbs_num})
		if sort_type == "excellent" {
			query += " and " + Bbs_post_excellent_excluded_sql("bbs_data")
			query += " and " + Bbs_post_tabom_score_sql("bbs_data") + " > ?"
			values = append(values, Bbs_excellent_min(db, bbs_num))
		}
		query += " order by set_data desc"
		pinned_listing_rows = tool.Query_DB(db, query, values...)

		row_alias := "title"
		query = "select title.set_code, title.set_id, '0'"
		if sort_type == "activity" {
			query += ", " + Bbs_post_last_activity_sql(row_alias) + " as activity_date, " + Bbs_post_last_activity_missing_sql(row_alias) + " as activity_cache_missing"
		}
		query += " from bbs_data title"
		if sort_type == "view" {
			query += " left join bbs_data view_data on view_data.set_name = 'view_count' and view_data.set_id = title.set_id and view_data.set_code = title.set_code"
		}
		query += " where title.set_name = 'title' and title.set_id like ?"

		query, values = append_post_view_and_filter_sql(query, row_alias, []any{bbs_num})

		order_by := "title.set_code + 0 desc"
		switch sort_type {
		case "activity":
			order_by = "activity_date desc, title.set_code + 0 desc"
		case "view":
			order_by = "coalesce(view_data.set_data, '0') + 0 desc, title.set_code + 0 desc"
		case "comment":
			order_by = Bbs_post_comment_count_sql(row_alias) + " desc, title.set_code + 0 desc"
		case "tabom":
			order_by = Bbs_post_tabom_count_sql(row_alias) + " desc, title.set_code + 0 desc"
		case "excellent":
			score_sql := Bbs_post_tabom_score_sql(row_alias)
			query += " and " + Bbs_post_excellent_excluded_sql(row_alias)
			query += " and " + score_sql + " > ?"
			values = append(values, Bbs_excellent_min(db, bbs_num))
			order_by = score_sql + " desc, title.set_code + 0 desc"
		}
		query += " order by " + order_by + " limit ?, 51"
		values = append(values, num)
		ordinary_listing_rows = tool.Query_DB(db, query, values...)
	}

	data_list := []map[string]string{}
	has_next := false
	ip_parser_temp := map[string]bbs_user_display_data{}
	activity_update_list := []bbs_post_activity_update{}

	hydrate_post_data := func(temp_data map[string]string, set_code string, set_id string) {
		post_rows := tool.Query_DB(
			db,
			"select set_name, set_data from bbs_data where set_code = ? and set_id = ?",
			set_code,
			set_id,
		)
		defer post_rows.Close()

		for post_rows.Next() {
			var set_name string
			var set_data string
			if err := post_rows.Scan(&set_name, &set_data); err != nil {
				panic(err)
			}
			if set_name == tool.Bbs_post_creation_time_name {
				continue
			}

			if set_name == "user_id" {
				user_display, exists := ip_parser_temp[set_data]
				if !exists {
					user_display = bbs_user_display_data{
						ip_pre:    tool.IP_preprocess(db, set_data, config.IP)[0],
						ip_render: tool.Get_user_profile_image_ui(db, set_data) + tool.IP_parser(db, set_data, config.IP),
					}
					ip_parser_temp[set_data] = user_display
				}

				set_data = user_display.ip_pre
				temp_data["user_id_render"] = user_display.ip_render
			}

			if set_name == "tag" {
				if temp_data["tags"] != "" {
					temp_data["tags"] += ", "
				}
				temp_data["tags"] += set_data
				continue
			}

			if set_name != "data" && set_name != "pinned" {
				temp_data[set_name] = set_data
			}
		}
		if err := post_rows.Err(); err != nil {
			panic(err)
		}
	}

	process_listing_rows := func(listing_rows *sql.Rows, is_ordinary_list bool) {
		normal_count := 0
		for listing_rows.Next() {
			temp_data := make(map[string]string)

			var set_code string
			var set_id string
			var pinned string
			activity_date := ""
			activity_cache_missing := 0
			var excellent_score any
			var err error
			if sort_type == "activity" {
				err = listing_rows.Scan(&set_code, &set_id, &pinned, &activity_date, &activity_cache_missing)
			} else if global_excellent_listing {
				err = listing_rows.Scan(&set_code, &set_id, &pinned, &excellent_score)
			} else {
				err = listing_rows.Scan(&set_code, &set_id, &pinned)
			}
			if err != nil {
				panic(err)
			}
			if is_ordinary_list {
				normal_count++
				if normal_count > 50 {
					has_next = true
					continue
				}
			}
			if sort_type == "activity" && activity_cache_missing == 1 && activity_date != "" {
				activity_update_list = append(activity_update_list, bbs_post_activity_update{
					set_id:   set_id,
					set_code: set_code,
					date:     activity_date,
				})
			}

			if bbs_num == "" && !global_excellent_listing && !tool.Check_acl(db, set_id, "", "bbs_view", config.IP) {
				continue
			}

			temp_data["set_code"] = set_code
			temp_data["set_id"] = set_id
			temp_data["pinned"] = pinned
			if activity_date != "" {
				temp_data["activity_date"] = activity_date
			}
			if !global_excellent_listing {
				hydrate_post_data(temp_data, set_code, set_id)
			}
			data_list = append(data_list, temp_data)
		}
		if err := listing_rows.Err(); err != nil {
			panic(err)
		}
		listing_rows.Close()
	}

	if bbs_num == "" {
		if main_listing_rows != nil {
			process_listing_rows(main_listing_rows, global_excellent_listing)
			main_listing_rows = nil
		}
		if global_excellent_listing {
			Bbs_global_excellent_hydrate(db, config, data_list, ip_parser_temp)
		}
	} else {
		process_listing_rows(pinned_listing_rows, false)
		pinned_listing_rows = nil
		process_listing_rows(ordinary_listing_rows, true)
		ordinary_listing_rows = nil
	}

	if len(activity_update_list) > 0 {
		err := Bbs_board_write_transaction(db, bbs_num, "", func(tx *sql.Tx) error {
			for _, activity_data := range activity_update_list {
				Bbs_post_last_activity_update(tx, activity_data.set_id, activity_data.set_code, activity_data.date)
			}
			return nil
		})
		if err != nil && !errors.Is(err, bbs_board_missing_error) {
			panic(err)
		}
	}

	return_data := make(map[string]any)
	return_data["response"] = "ok"
	return_data["data"] = data_list

	return_data["has_next"] = has_next
	return return_data
}

func Bbs_post_view_auth(db *sql.DB, set_id string, set_code string, ip string) (string, bool) {
	_, title_exists := tool.Get_bbs_data_value(db, set_id, set_code, "title")
	user_id, user_exists := tool.Get_bbs_data_value(db, set_id, set_code, "user_id")
	if !title_exists || !user_exists {
		return "", false
	}
	if !tool.Check_acl(db, set_id, "", "bbs_view", ip) || !Bbs_post_view_allowed(db, set_id, set_code, user_id, ip, nil) {
		return "", false
	}
	return user_id, true
}
