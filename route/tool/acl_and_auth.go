package tool

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/3th1nk/cidr"
	"github.com/dlclark/regexp2"
)

func List_acl(func_type string) []string {
	if func_type == "user_document" {
		return []string{
			"",
			"owner",
			"admin",
			"user",
			"ip",
		}
	}

	return []string{
		"",
		"owner",
		"admin",
		"user",
		"ip",
		"ban",
		"ban_without_login",
		"ban_without_site",
		"email_verified",
		"up_to_level_10",
		"up_to_level_3",
		"trust_a",
		"trust_b",
		"trust_c",
		"trust_d",
	}
}

func List_auth(db *sql.DB) []string {
	rows := Query_DB(
		db,
		"select distinct name from alist",
	)
	defer rows.Close()

	data_list := []string{}

	for rows.Next() {
		var name string

		err := rows.Scan(&name)
		if err != nil {
			panic(err)
		}
		if Auth_group_name_reserved(name) && !Auth_group_name_default(name) {
			continue
		}

		data_list = append(data_list, name)
	}

	return data_list
}

func Do_insert_auth_history(db *sql.DB, ip string, what string) {
	log_off := ""
	QueryRow_DB(
		db,
		"select data from other where name = 'auth_history_off'",
		[]any{&log_off},
	)

	if log_off == "" {
		time := Get_time()

		Exec_DB(
			db,
			"insert into re_admin (who, what, time) values (?, ?, ?)",
			ip, what, time,
		)
	}
}

func auth_end_active(end string) bool {
	if end == "" || end == "0" {
		return true
	}

	end_time, end_err := time.Parse("2006-01-02 15:04:05", end)
	if end_err != nil {
		end_time, end_err = time.Parse("2006-01-02", end)
	}
	now_time, now_err := time.Parse("2006-01-02 15:04:05", Get_time())
	if end_err != nil || now_err != nil {
		return true
	}

	return now_time.Before(end_time)
}

func get_user_auth_raw(db *sql.DB, ip string) string {
	auth := "ip"
	acl_end := ""
	exist := QueryRow_DB(
		db,
		"select user_set.data, coalesce(acl_end.data, '') from user_set inner join alist on alist.name = user_set.data left join user_set as acl_end on acl_end.id = user_set.id and acl_end.name = 'acl_end' where user_set.id = ? and user_set.name = 'acl' limit 1",
		[]any{&auth, &acl_end},
		ip,
	)

	if exist && !auth_end_active(acl_end) {
		exist = false
	}

	if exist && Auth_group_name_reserved(auth) && !Auth_group_name_default(auth) {
		exist = false
	}
	if !exist {
		if !IP_or_user(ip) {
			auth = "user"
		} else {
			auth = "ip"
		}
	}

	return auth
}

func get_ban_auth_group(db *sql.DB, login string) string {
	switch login {
	case "L", "O", "ban":
		return "ban"
	case "D", "ban_without_site":
		return "ban_without_site"
	case "A", "E", "ban_without_login":
		return "ban_without_login"
	}
	if Auth_group_exists(db, login) {
		return login
	}
	return "ban_without_login"
}

func Get_auth_target_group(db *sql.DB, target string, target_type string) string {
	if target_type == "" || target_type == "normal" {
		return Get_user_auth(db, target)
	}

	auth := ""
	end := ""
	QueryRow_DB(
		db,
		"select login, end from rb where block = ? and band = ? and ongoing = '1' order by today desc limit 1",
		[]any{&auth, &end},
		target,
		target_type,
	)
	if auth == "" || !auth_end_active(end) {
		return "ip"
	}

	return get_ban_auth_group(db, auth)
}

func get_pattern_auth_group(db *sql.DB, ip string) string {
	rows := Query_DB(
		db,
		"select login, block, end from rb where band = 'regex' and ongoing = '1'",
	)
	for rows.Next() {
		login := ""
		block := ""
		end := ""
		if rows.Scan(&login, &block, &end) != nil {
			continue
		}
		if !auth_end_active(end) {
			continue
		}
		r, err := regexp2.Compile(block, 0)
		if err == nil {
			if match, _ := r.FindStringMatch(ip); match != nil {
				rows.Close()
				return get_ban_auth_group(db, login)
			}
		}
	}
	rows.Close()

	if IP_or_user(ip) {
		rows := Query_DB(
			db,
			"select login, block, end from rb where band = 'cidr' and ongoing = '1'",
		)
		for rows.Next() {
			login := ""
			block := ""
			end := ""
			if rows.Scan(&login, &block, &end) != nil {
				continue
			}
			if !auth_end_active(end) {
				continue
			}
			c, err := cidr.Parse(block)
			if err == nil && c.Contains(ip) {
				rows.Close()
				return get_ban_auth_group(db, login)
			}
		}
		rows.Close()
	}
	return ""
}

func Get_user_auth(db *sql.DB, ip string) string {
	auth := get_user_auth_raw(db, ip)
	if Auth_group_name_ban(auth) {
		return auth
	}
	if pattern_auth := get_pattern_auth_group(db, ip); pattern_auth != "" {
		return pattern_auth
	}
	return auth
}

func Get_auth_permission_list(db *sql.DB, auth string) []string {
	rows := Query_DB(
		db,
		"select acl from alist where name = ?",
		auth,
	)
	defer rows.Close()

	data_list := []string{}
	for rows.Next() {
		name := ""
		if err := rows.Scan(&name); err != nil {
			panic(err)
		}
		data_list = append(data_list, name)
	}
	return data_list
}

func Get_auth_user_list(db *sql.DB, offset int, limit int) [][]string {
	query := "select id, data from user_set where name = 'acl' and data != 'user' order by id"
	args := []any{}
	if limit > 0 {
		query += " limit ?, ?"
		args = append(args, offset, limit)
	}

	rows := Query_DB(db, query, args...)
	defer rows.Close()

	data_list := [][]string{}
	for rows.Next() {
		data := []string{"", ""}
		if rows.Scan(&data[0], &data[1]) != nil {
			continue
		}
		auth := Get_user_auth(db, data[0])
		if auth == "user" {
			continue
		}
		data[1] = auth
		data_list = append(data_list, data)
	}
	return data_list
}

func Get_auth_group_info(db *sql.DB, auth string) map[string]bool {
	data_list := map[string]bool{}

	for _, name := range Get_auth_permission_list(db, auth) {
		data_list[name] = true
	}

	if len(data_list) == 0 {
		data_list["nothing"] = true
	}
	if auth != "" {
		data_list["group:"+auth] = true
	}

	return Check_auth(data_list)
}

func Get_acl_data_list(db *sql.DB, title string, acl_type string) []string {
	rows := Query_DB(
		db,
		"select data from acl where title = ? and type = ? and data != '' and data != 'normal'",
		title, acl_type,
	)
	defer rows.Close()

	data_list := []string{}
	seen := map[string]bool{}
	for rows.Next() {
		data := ""
		if rows.Scan(&data) == nil && !seen[data] {
			seen[data] = true
			data_list = append(data_list, data)
		}
	}

	return data_list
}

func Check_acl_group(db *sql.DB, acl_data string, auth_info map[string]bool) bool {
	if acl_data == "" || acl_data == "normal" || acl_data == "nothing" {
		return false
	}
	if Auth_group_exists(db, acl_data) && auth_info["group:"+acl_data] {
		return true
	}
	if Auth_permission_name(acl_data) {
		return auth_info[acl_data]
	}
	return false
}

func get_document_filter_acl(data string, action string) (string, bool) {
	data = strings.ReplaceAll(data, "\r", "")
	if !strings.Contains(data, "=") {
		if action == "edit" {
			return strings.TrimSpace(data), true
		}
		return "", false
	}

	for _, line := range strings.Split(data, "\n") {
		parts := strings.SplitN(strings.TrimSpace(line), "=", 2)
		if len(parts) == 2 && strings.TrimSpace(parts[0]) == action {
			return strings.TrimSpace(parts[1]), true
		}
	}

	return "", false
}

func check_document_filter_acl(db *sql.DB, doc_name string, action string, auth_info map[string]bool) (bool, bool) {
	rows := Query_DB(
		db,
		"select plus, plus_t from html_filter where kind = 'document'",
	)
	defer rows.Close()

	matched := false
	for rows.Next() {
		pattern := ""
		acl_data := ""
		if rows.Scan(&pattern, &acl_data) != nil {
			continue
		}

		regex, err := regexp2.Compile(pattern, 0)
		if err != nil {
			continue
		}

		match, err := regex.MatchString(doc_name)
		if err != nil || !match {
			continue
		}

		acl_data, ok := get_document_filter_acl(acl_data, action)
		if !ok {
			continue
		}

		matched = true
		if acl_data == "" || acl_data == "normal" || Check_acl_group(db, acl_data, auth_info) {
			return true, true
		}
	}

	return matched, false
}

func Rankup_group_list() []string {
	return []string{"trust_a", "trust_b", "trust_c", "trust_d"}
}

func Auth_group_name_rankup(name string) bool {
	return Arr_in_str(Rankup_group_list(), name)
}

func Rankup_condition_type(data string) string {
	if data == "edit" || data == "time" {
		return "int"
	}

	return ""
}

func Rankup_condition_list(data string) ([]string, bool) {
	data = strings.ReplaceAll(data, "\r", "")
	condition_list := []string{}
	for _, line := range strings.Split(data, "\n") {
		parts := strings.Fields(strings.TrimSpace(line))
		if len(parts) == 0 {
			continue
		}
		if len(parts) != 2 || Rankup_condition_type(strings.ToLower(parts[0])) == "" {
			return nil, false
		}

		condition_value, err := strconv.Atoi(parts[1])
		if err != nil || condition_value < 0 {
			return nil, false
		}

		condition := strings.ToLower(parts[0]) + " " + strconv.Itoa(condition_value)
		if !Arr_in_str(condition_list, condition) {
			condition_list = append(condition_list, condition)
		}
	}

	if len(condition_list) == 0 {
		return nil, false
	}

	return condition_list, true
}

func Get_rankup_auth_info(db *sql.DB, ip string) map[string]bool {
	auth_info := map[string]bool{}
	if IP_or_user(ip) {
		return auth_info
	}

	condition_map := map[string][]string{}
	invalid_map := map[string]bool{}
	for _, condition_data := range Get_setting(db, "rankup_condition", "") {
		if len(condition_data) < 2 || !Auth_group_name_rankup(condition_data[1]) {
			continue
		}

		condition_list, ok := Rankup_condition_list(condition_data[0])
		if !ok {
			invalid_map[condition_data[1]] = true
			continue
		}

		condition_map[condition_data[1]] = append(condition_map[condition_data[1]], condition_list...)
	}

	edit_required := false
	time_required := false
	for _, condition_list := range condition_map {
		for _, condition := range condition_list {
			parts := strings.Fields(condition)
			if len(parts) != 2 {
				continue
			}
			if parts[0] == "edit" {
				edit_required = true
			}
			if parts[0] == "time" {
				time_required = true
			}
		}
	}

	edit_count := 0
	if edit_required {
		QueryRow_DB(
			db,
			"select count(*) from history where ip = ?",
			[]any{&edit_count},
			ip,
		)
	}

	signup_date := ""
	if time_required {
		QueryRow_DB(
			db,
			"select data from user_set where id = ? and name = 'date'",
			[]any{&signup_date},
			ip,
		)
	}

	var date time.Time
	var date_err error
	if time_required {
		date, date_err = time.Parse("2006-01-02 15:04:05", signup_date)
		if date_err != nil {
			date, date_err = time.Parse("2006-01-02", signup_date)
		}
	}
	now, now_err := time.Parse("2006-01-02 15:04:05", Get_time())

	for _, rankup_group := range Rankup_group_list() {
		condition_list := condition_map[rankup_group]
		if len(condition_list) == 0 || invalid_map[rankup_group] {
			continue
		}

		passed := true
		for _, condition := range condition_list {
			parts := strings.Fields(condition)
			condition_value := Str_to_int(parts[1])
			switch parts[0] {
			case "edit":
				passed = edit_count >= condition_value
			case "time":
				passed = date_err == nil && now_err == nil && !now.Before(date.AddDate(0, 0, condition_value))
			}
			if !passed {
				break
			}
		}

		if passed {
			auth_info[rankup_group] = true
		}
	}

	return auth_info
}

func Get_auth_info(db *sql.DB, ip string) map[string]bool {
	auth_name := Get_user_auth(db, ip)
	auth_info := Get_auth_group_info(db, auth_name)
	if Auth_group_name_ban(auth_name) {
		return auth_info
	}
	if !IP_or_user(ip) && auth_info["do_email_verified"] && !auth_info["email_verified"] {
		email := ""
		QueryRow_DB(
			db,
			"select data from user_set where id = ? and name = 'email'",
			[]any{&email},
			ip,
		)
		if email != "" {
			auth_info["email_verified"] = true
		}
	}
	if IP_or_user(ip) || auth_info["admin"] || auth_info["owner"] {
		return auth_info
	}

	if auth_info["rankup"] {
		for auth := range Get_rankup_auth_info(db, ip) {
			auth_info[auth] = true
			for name := range Get_auth_group_info(db, auth) {
				auth_info[name] = true
			}
		}
	}

	return Check_auth(auth_info)
}

func Get_auth_date(db *sql.DB, user_name string) string {
	data := ""

	QueryRow_DB(
		db,
		"select data from user_set where id = ? and name = 'auth_date'",
		[]any{&data},
		user_name,
	)

	if data == "" {
		data = "0"
	}

	return data
}

func Get_auth_level(auth_info map[string]bool) int {
	if auth_info["owner"] {
		return 7
	}
	if auth_info["admin"] {
		return 6
	}
	if auth_info["up_to_level_10"] {
		return 5
	}
	if auth_info["trust_d"] {
		return 5
	}
	if auth_info["trust_c"] {
		return 4
	}
	if auth_info["up_to_level_3"] {
		return 3
	}
	if auth_info["trust_a"] || auth_info["trust_b"] || auth_info["email_verified"] {
		return 3
	}
	if auth_info["user"] {
		return 2
	}
	if auth_info["ip"] {
		return 1
	}
	return 0
}

func Auth_group_name_default(name string) bool {
	return Arr_in_str([]string{
		"owner",
		"admin",
		"user",
		"ip",
		"ban",
		"ban_without_login",
		"ban_without_site",
		"email_verified",
		"up_to_level_10",
		"up_to_level_3",
		"trust_a",
		"trust_b",
		"trust_c",
		"trust_d",
	}, name)
}

func Auth_group_name_ban(name string) bool {
	return name == "ban" || name == "ban_without_login" || name == "ban_without_site"
}

func Auth_permission_name(name string) bool {
	for _, choice := range Auth_choices() {
		if choice.Key == name {
			return true
		}
	}
	return false
}

func Auth_group_name_reserved(name string) bool {
	return name == "normal" || Auth_permission_name(name)
}

func Auth_group_exists(db *sql.DB, auth string) bool {
	if auth == "" || (Auth_group_name_reserved(auth) && !Auth_group_name_default(auth)) {
		return false
	}

	name := ""
	return QueryRow_DB(
		db,
		"select name from alist where name = ? limit 1",
		[]any{&name},
		auth,
	)
}

func Auth_group_in_use(db *sql.DB, auth string) bool {
	queries := []string{
		"select id from user_set where name = 'acl' and data = ? limit 1",
		"select title from acl where data = ? limit 1",
		"select code from rd where acl = ? limit 1",
		"select thread_code from topic_set where set_name = 'thread_view_acl' and set_data = ? limit 1",
		"select set_id from bbs_set where set_name in ('bbs_view_acl', 'bbs_acl', 'bbs_edit_acl', 'bbs_comment_acl', 'bbs_view_acl_all', 'bbs_acl_all', 'bbs_edit_acl_all', 'bbs_comment_acl_all') and set_data = ? limit 1",
		"select name from other where name in ('bbs_view_acl_all', 'bbs_acl_all', 'bbs_edit_acl_all', 'bbs_comment_acl_all') and coverage = '' and data = ? limit 1",
		"select id from vote where user = '' and acl = ? limit 1",
	}

	for _, query := range queries {
		value := ""
		if QueryRow_DB(db, query, []any{&value}, auth) {
			return true
		}
	}
	return false
}

func Auth_can_change_auth(db *sql.DB, ip string, before_auth string, after_auth string) bool {
	if !Auth_group_exists(db, before_auth) || !Auth_group_exists(db, after_auth) {
		return false
	}
	if !Check_acl(db, "", "", "give_auth", ip) {
		return false
	}
	if before_auth == after_auth {
		return true
	}

	auth_info := Get_auth_info(db, ip)
	if auth_info["owner"] {
		return true
	}

	auth_level := Get_auth_level(auth_info)
	if auth_level <= Get_auth_level(Get_auth_group_info(db, before_auth)) {
		return false
	}
	return auth_level > Get_auth_level(Get_auth_group_info(db, after_auth))
}

func Check_auth(auth_info map[string]bool) map[string]bool {
	if _, ok := auth_info["owner"]; ok {
		auth_info["admin"] = true
	}

	if _, ok := auth_info["admin"]; ok {
		auth_info["email_verified"] = true
		auth_info["up_to_level_10"] = true
		auth_info["trust_d"] = true
	}

	if _, ok := auth_info["up_to_level_10"]; ok {
		auth_info["up_to_level_3"] = true
	}

	if _, ok := auth_info["email_verified"]; ok {
		auth_info["user"] = true
	}

	if _, ok := auth_info["up_to_level_3"]; ok {
		auth_info["user"] = true
	}

	if _, ok := auth_info["trust_d"]; ok {
		auth_info["trust_c"] = true
	}

	if _, ok := auth_info["trust_c"]; ok {
		auth_info["trust_a"] = true
		auth_info["trust_b"] = true
	}

	if _, ok := auth_info["trust_a"]; ok {
		auth_info["large_edit"] = true
	}

	if _, ok := auth_info["give_range"]; ok {
		auth_info["give"] = true
	}

	if auth_info["trust_a"] || auth_info["trust_b"] {
		auth_info["user"] = true
	}

	admin_auth := []string{"toron", "check", "acl", "hidel", "give_range", "give", "bbs", "vote_fix"}

	if _, ok := auth_info["admin"]; ok {
		for _, v := range admin_auth {
			auth_info[v] = true
		}
		auth_info["edit_filter_manage"] = true
		auth_info["application_manage"] = true
	}

	if _, ok := auth_info["check"]; ok {
		auth_info["view_user_watchlist"] = true
	}

	check := false
	for _, v := range admin_auth {
		if _, ok := auth_info[v]; ok {
			check = true
			break
		}
	}

	if check {
		auth_info["admin_default_feature"] = true
	}

	admin_default_feature := []string{"treat_as_admin", "user_name_bold", "multiple_upload", "slow_edit_pass", "edit_bottom_compulsion_pass", "view_hide_user_name", "doc_watch_list_view", "edit_filter_pass", "user"}

	if _, ok := auth_info["admin_default_feature"]; ok {
		for _, v := range admin_default_feature {
			auth_info[v] = true
		}
	}

	user_default := []string{"rankup", "do_email_verified", "captcha_pass", "ip"}

	if _, ok := auth_info["user"]; ok {
		for _, v := range user_default {
			auth_info[v] = true
		}
	}

	ip_default := []string{"document", "discuss", "upload", "vote", "bbs_use", "captcha_one_check_five_pass", "edit_filter_view", "login_available", "register_available"}

	if _, ok := auth_info["ip"]; ok {
		for _, v := range ip_default {
			auth_info[v] = true
		}
	}

	document_default := []string{"edit", "move", "new_make", "delete"}

	if _, ok := auth_info["document"]; ok {
		for _, v := range document_default {
			auth_info[v] = true
		}
	}

	check = false
	for _, v := range document_default {
		if _, ok := auth_info[v]; ok {
			check = true
			break
		}
	}

	if check {
		auth_info["view"] = true
	}

	if auth_info["view"] {
		auth_info["site_view"] = true
	}

	topic_default := []string{"discuss_view", "discuss_make_new_thread"}

	if _, ok := auth_info["discuss"]; ok {
		for _, v := range topic_default {
			auth_info[v] = true
		}
	}

	bbs_default := []string{"bbs_edit", "bbs_comment"}

	if _, ok := auth_info["bbs_use"]; ok {
		for _, v := range bbs_default {
			auth_info[v] = true
		}
	}

	check = false
	for _, v := range bbs_default {
		if _, ok := auth_info[v]; ok {
			check = true
			break
		}
	}

	if check {
		auth_info["bbs_view"] = true
	}

	return auth_info
}

// PASS is TRUE
func Check_acl(db *sql.DB, name string, topic_number string, tool string, ip string) bool {
	auth_info := Get_auth_info(db, ip)

	ip_or_user := IP_or_user(ip)
	if !auth_info["site_view"] {
		return false
	}

	document_filter_action := ""
	switch tool {
	case "render":
		document_filter_action = "view"
	case "document_edit":
		document_filter_action = "edit"
	case "document_move":
		document_filter_action = "move"
	case "document_delete":
		document_filter_action = "delete"
	case "document_make_acl":
		document_filter_action = "new_make"
	}
	if document_filter_action != "" && !auth_info["acl"] {
		if matched, allowed := check_document_filter_acl(db, name, document_filter_action, auth_info); matched {
			return allowed
		}
	}

	if tool == "" && name != "" {
		if !Check_acl(db, name, "", "render", ip) {
			return false
		}

		if strings.HasPrefix(name, "user:") {
			user_page_str := name[5:]
			if slash_index := strings.Index(user_page_str, "/"); slash_index != -1 {
				user_page_str = user_page_str[:slash_index]
			}

			if auth_info["acl"] {
				return true
			}

			acl_data_list := Get_acl_data_list(db, name, "decu")
			if len(acl_data_list) == 0 {
				return false
			}

			for _, acl_data := range acl_data_list {
				if Check_acl_group(db, acl_data, auth_info) {
					return true
				}
				if ip == user_page_str && !ip_or_user {
					return true
				}
			}

			return false
		}
	}

	if Arr_in_str([]string{"document_edit", "document_move", "document_delete"}, tool) {
		if !Check_acl(db, name, topic_number, "render", ip) {
			return false
		} else if !Check_acl(db, name, topic_number, "", ip) {
			return false
		}
	} else if Arr_in_str([]string{"bbs_edit", "bbs_comment"}, tool) {
		if !Check_acl(db, name, topic_number, "bbs_view", ip) {
			return false
		}
	} else if Arr_in_str([]string{"topic"}, tool) {
		if !Check_acl(db, name, topic_number, "topic_view", ip) {
			return false
		}
	}

	if Arr_in_str([]string{"topic", "topic_view"}, tool) {
		if name == "" {
			name = "test"
			QueryRow_DB(
				db,
				"select title from rd where code = ?",
				[]any{&name},
				topic_number,
			)
		}
	}

	end_number := 1
	for for_a := 0; for_a < end_number; for_a++ {
		acl_data := ""
		acl_data_list := []string{}
		acl_pass_auth := ""

		if tool == "all_admin_auth" {
			acl_pass_auth = "treat_as_admin"
			acl_data = "owner"
		} else if tool == "owner_auth" {
			acl_pass_auth = "owner"
			acl_data = "owner"
		} else if tool == "edit_filter_auth" {
			acl_pass_auth = "edit_filter_manage"
			acl_data = "owner"
		} else if tool == "application_auth" {
			acl_pass_auth = "application_manage"
			acl_data = "owner"
		} else if tool == "bbs_auth" {
			acl_pass_auth = "bbs"
			acl_data = "owner"
		} else if tool == "toron_auth" {
			acl_pass_auth = "toron"
			acl_data = "owner"
		} else if tool == "check_auth" {
			acl_pass_auth = "check"
			acl_data = "owner"
		} else if tool == "acl_auth" {
			acl_pass_auth = "acl"
			acl_data = "owner"
		} else if tool == "hidel_auth" {
			acl_pass_auth = "hidel"
			acl_data = "owner"
		} else if tool == "give_auth" {
			acl_pass_auth = "give"
			acl_data = "owner"
		} else if tool == "give_range_auth" {
			acl_pass_auth = "give_range"
			acl_data = "owner"
		} else if tool == "vote_auth" {
			acl_pass_auth = "vote_fix"
			acl_data = "owner"
		} else if tool == "" {
			acl_pass_auth = "acl"

			if for_a == 0 {
				end_number += 1

				acl_data_list = Get_acl_data_list(db, name, "decu")
			} else {
				if auth_info["document"] {
					acl_data = ""
				} else {
					acl_data = "owner"
				}
			}
		} else if tool == "document_move" {
			acl_pass_auth = "acl"

			if for_a == 0 {
				end_number += 1

				acl_data_list = Get_acl_data_list(db, name, "document_move_acl")
			} else {
				if auth_info["move"] {
					acl_data = ""
				} else {
					acl_data = "owner"
				}
			}
		} else if tool == "document_edit" {
			acl_pass_auth = "acl"

			if for_a == 0 {
				end_number += 1

				acl_data_list = Get_acl_data_list(db, name, "document_edit_acl")
			} else {
				if auth_info["edit"] {
					acl_data = ""
				} else {
					acl_data = "owner"
				}
			}
		} else if tool == "document_delete" {
			acl_pass_auth = "acl"

			if for_a == 0 {
				end_number += 1

				acl_data_list = Get_acl_data_list(db, name, "document_delete_acl")
			} else {
				if auth_info["delete"] {
					acl_data = ""
				} else {
					acl_data = "owner"
				}
			}
		} else if tool == "topic" {
			acl_pass_auth = "topic"

			switch for_a {
			case 0:
				end_number += 1

				QueryRow_DB(
					db,
					"select acl from rd where code = ?",
					[]any{&acl_data},
					topic_number,
				)
			case 1:
				end_number += 1

				acl_data_list = Get_acl_data_list(db, name, "dis")
			default:
				if auth_info["discuss"] {
					acl_data = ""
				} else {
					acl_data = "owner"
				}
			}
		} else if tool == "topic_view" {
			acl_pass_auth = "topic"

			if for_a == 0 {
				end_number += 1

				QueryRow_DB(
					db,
					"select set_data from topic_set where thread_code = ? and set_name = 'thread_view_acl'",
					[]any{&acl_data},
					topic_number,
				)
			} else {
				if auth_info["discuss_view"] {
					acl_data = ""
				} else {
					acl_data = "owner"
				}
			}
		} else if tool == "upload" {
			acl_pass_auth = "admin_default_feature"

			if auth_info["upload"] {
				acl_data = ""
			} else {
				acl_data = "owner"
			}
		} else if tool == "many_upload" {
			acl_pass_auth = "admin_default_feature"

			if auth_info["multiple_upload"] {
				acl_data = ""
			} else {
				acl_data = "owner"
			}
		} else if tool == "vote" {
			acl_pass_auth = "vote_fix"

			if for_a == 0 {
				end_number += 1

				if topic_number != "" {
					QueryRow_DB(
						db,
						"select acl from vote where id = ? and user = ''",
						[]any{&acl_data},
						topic_number,
					)
				} else {
					continue
				}
			} else {
				if auth_info["vote"] {
					acl_data = ""
				} else {
					acl_data = "owner"
				}
			}
		} else if tool == "slow_edit" {
			acl_pass_auth = "admin_default_feature"

			if auth_info["slow_edit_pass"] {
				acl_data = ""
			} else {
				acl_data = "owner"
			}
		} else if tool == "edit_bottom_compulsion" {
			acl_pass_auth = "admin_default_feature"

			if auth_info["edit_bottom_compulsion_pass"] {
				acl_data = ""
			} else {
				acl_data = "owner"
			}
		} else if tool == "bbs_edit" {
			acl_pass_auth = "bbs"

			switch for_a {
			case 0:
				end_number += 1

				QueryRow_DB(
					db,
					"select set_data from bbs_set where set_name = 'bbs_edit_acl' and set_id = ?",
					[]any{&acl_data},
					name,
				)
			case 1:
				end_number += 1

				QueryRow_DB(
					db,
					"select set_data from bbs_set where set_name = 'bbs_acl' and set_id = ?",
					[]any{&acl_data},
					name,
				)
			case 2:
				end_number += 1

				QueryRow_DB(
					db,
					"select set_data from bbs_set where set_name = 'bbs_edit_acl_all'",
					[]any{&acl_data},
				)
			default:
				if auth_info["bbs_edit"] {
					acl_data = ""
				} else {
					acl_data = "owner"
				}
			}
		} else if tool == "bbs_comment" {
			acl_pass_auth = "bbs"

			switch for_a {
			case 0:
				end_number += 1

				QueryRow_DB(
					db,
					"select set_data from bbs_set where set_name = 'bbs_comment_acl' and set_id = ?",
					[]any{&acl_data},
					name,
				)
			case 1:
				end_number += 1

				QueryRow_DB(
					db,
					"select set_data from bbs_set where set_name = 'bbs_acl' and set_id = ?",
					[]any{&acl_data},
					name,
				)
			case 2:
				end_number += 1

				QueryRow_DB(
					db,
					"select set_data from bbs_set where set_name = 'bbs_comment_acl_all'",
					[]any{&acl_data},
				)
			default:
				if auth_info["bbs_comment"] {
					acl_data = ""
				} else {
					acl_data = "owner"
				}
			}
		} else if tool == "bbs_view" {
			acl_pass_auth = "bbs"

			if for_a == 0 {
				end_number += 1

				QueryRow_DB(
					db,
					"select set_data from bbs_set where set_name = 'bbs_view_acl' and set_id = ?",
					[]any{&acl_data},
					name,
				)
			} else {
				if auth_info["bbs_view"] {
					acl_data = ""
				} else {
					acl_data = "owner"
				}
			}
		} else if tool == "discuss_make_new_thread" {
			acl_pass_auth = "toron"

			if auth_info["discuss_make_new_thread"] {
				acl_data = ""
			} else {
				acl_data = "owner"
			}
		} else if tool == "recaptcha" {
			acl_pass_auth = "admin_default_feature"

			if auth_info["captcha_pass"] {
				acl_data = ""
			} else {
				acl_data = "owner"
			}
		} else if tool == "recaptcha_five_pass" {
			acl_pass_auth = "admin_default_feature"

			if auth_info["captcha_one_check_five_pass"] {
				acl_data = ""
			} else {
				acl_data = "owner"
			}
		} else if tool == "edit_filter_pass" {
			acl_pass_auth = "admin_default_feature"

			if auth_info["edit_filter_pass"] {
				acl_data = ""
			} else {
				acl_data = "owner"
			}
		} else if tool == "edit_filter_view" {
			acl_pass_auth = "edit_filter_pass"

			if auth_info["edit_filter_view"] {
				acl_data = ""
			} else {
				acl_data = "owner"
			}
		} else if tool == "view_hide_user_name" {
			acl_pass_auth = "admin_default_feature"

			if auth_info["view_hide_user_name"] {
				acl_data = ""
			} else {
				acl_data = "owner"
			}
		} else if tool == "user_name_bold" {
			acl_pass_auth = "admin_default_feature"

			if auth_info["user_name_bold"] {
				acl_data = ""
			} else {
				acl_data = "owner"
			}
		} else if tool == "doc_watch_list_view" {
			acl_pass_auth = "admin_default_feature"

			if auth_info["doc_watch_list_view"] {
				acl_data = ""
			} else {
				acl_data = "owner"
			}
		} else if tool == "document_make_acl" {
			acl_pass_auth = "acl"

			if auth_info["new_make"] {
				acl_data = ""
			} else {
				acl_data = "owner"
			}
		} else {
			// tool == "render"
			acl_pass_auth = "acl"

			if for_a == 0 {
				end_number += 1

				acl_data_list = Get_acl_data_list(db, name, "view")
			} else {
				if auth_info["view"] {
					acl_data = ""
				} else {
					acl_data = "owner"
				}
			}
		}

		if len(acl_data_list) > 0 {
			if auth_info[acl_pass_auth] {
				return true
			}
			for _, value := range acl_data_list {
				if Check_acl_group(db, value, auth_info) {
					return true
				}
			}
			return false
		}

		if auth_info[acl_pass_auth] {
			return true
		}

		if acl_data == "" {
			acl_data = "normal"
		}

		if acl_data != "normal" {
			if Check_acl_group(db, acl_data, auth_info) {
				return true
			}
			return false
		} else if for_a == end_number-1 {
			if tool == "topic" {
				topic_state := ""
				QueryRow_DB(
					db,
					"select title from rd where code = ? and stop != ''",
					[]any{&topic_state},
					topic_number,
				)

				if topic_state != "" {
					if auth_info["topic"] {
						return true
					} else {
						return false
					}
				} else {
					return true
				}
			} else {
				return true
			}
		}
	}

	return false
}
