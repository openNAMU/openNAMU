package tool

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var builtin_version_data []byte

func Set_builtin_version_data(data []byte) {
	builtin_version_data = data
}

func DB_table_list() map[string][]string {
	create_data := map[string][]string{}

	// 폐지 예정 (data_set으로 통합)
	create_data["data_set"] = []string{"doc_name", "doc_rev", "set_name", "set_data"}

	create_data["data"] = []string{"title", "data", "type"}
	create_data["history"] = []string{"id", "title", "data", "date", "ip", "send", "leng", "hide", "type"}
	create_data["rc"] = []string{"id", "title", "date", "type"}
	create_data["acl"] = []string{"title", "data", "type"}

	// 개편 예정 (data_link로 변경)
	create_data["back"] = []string{"title", "link", "type", "data"}

	// 폐지 예정 (topic_set으로 통합) [가장 시급]
	create_data["topic_set"] = []string{"thread_code", "set_name", "set_id", "set_data"}

	create_data["rd"] = []string{"title", "sub", "code", "date", "band", "stop", "agree", "acl"}
	create_data["topic"] = []string{"id", "data", "date", "ip", "block", "top", "code"}

	// 폐지 예정 (user_set으로 통합)
	create_data["rb"] = []string{"block", "end", "today", "blocker", "why", "band", "login", "ongoing"}

	// 개편 예정 (wiki_set과 wiki_filter과 wiki_vote으로 변경)
	create_data["other"] = []string{"name", "data", "coverage"}
	create_data["html_filter"] = []string{"html", "kind", "plus", "plus_t"}
	create_data["vote"] = []string{"name", "id", "subject", "data", "user", "type", "acl"}

	// 개편 예정 (auth와 auth_log로 변경)
	create_data["alist"] = []string{"name", "acl"}
	create_data["re_admin"] = []string{"who", "what", "time"}

	// 개편 예정 (user_notice와 user_agent로 변경)
	create_data["ua_d"] = []string{"name", "ip", "ua", "today", "sub"}

	create_data["user_set"] = []string{"name", "id", "data"}
	create_data["user_notice"] = []string{"id", "name", "data", "date", "readme"}

	create_data["bbs_set"] = []string{"set_name", "set_code", "set_id", "set_data"}
	create_data["bbs_data"] = []string{"set_name", "set_code", "set_id", "set_data"}

	return create_data
}

func DB_make_MySQL(db *sql.DB, new_db_set map[string]string) {
	Exec_DB(
		db,
		`create database if not exists `+new_db_set["db_name"]+` default character set utf8mb4`,
	)
}

func DB_column_exists(db *sql.DB, table_name string, column_name string) bool {
	query := fmt.Sprintf(
		"SELECT %s FROM %s LIMIT 1",
		column_name,
		table_name,
	)

	rows, err := db.Query(query)
	if err != nil {
		return false
	}
	defer rows.Close()

	return true
}

func DB_field_text(db_type string) string {
	if db_type == "mysql" {
		return "longtext"
	}

	return "text default ''"
}

func DB_warn_null_column(db *sql.DB, table_name string, column_name string) {
	query := fmt.Sprintf(
		"select count(*) from %s where %s is null",
		table_name,
		column_name,
	)

	var count int
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Printf("[DB WARNING] null check failed: %s.%s: %v", table_name, column_name, err)
		return
	}

	if count > 0 {
		log.Printf("[DB WARNING] %s.%s has %d null values", table_name, column_name, count)
	}
}

func DB_create_table(db *sql.DB, table_name string, field_text string) {
	Exec_DB(
		db,
		fmt.Sprintf(
			"create table if not exists %s (test %s)",
			table_name,
			field_text,
		),
	)
}

func DB_alter_add_column(db *sql.DB, table_name string, column_name string, field_text string) {
	Exec_DB(
		db,
		fmt.Sprintf(
			"alter table %s add column %s %s",
			table_name,
			column_name,
			field_text,
		),
	)
}

func DB_create_index(db *sql.DB) {
	queries := []string{
		"create index history_index on history (title, ip)",
		"create index history_title_id_index on history (title, id)",
		"create index history_ip_date_index on history (ip, date)",
		"create index bbs_data_index on bbs_data (set_id, set_code, set_name)",
		"create index acl_title_type_index on acl (title, type)",
		"create index data_set_document_index on data_set (doc_name, set_name, doc_rev)",
		"create index data_set_list_index on data_set (set_name, doc_rev)",
		"create index data_title_index on data (title)",
		"create index back_title_type_index on back (title, type)",
		"create index back_type_title_index on back (type, title)",
		"create index back_link_type_index on back (link, type)",
		"create index user_set_id_name_index on user_set (id, name)",
		"create index topic_code_id_index on topic (code, id)",
	}
	if Get_DB_type() == "mysql" {
		queries = []string{
			"create index history_index on history (title(191), ip(191))",
			"create index history_title_id_index on history (title(191), id(191))",
			"create index history_ip_date_index on history (ip(191), date(191))",
			"create index bbs_data_index on bbs_data (set_id(191), set_code(191), set_name(191))",
			"create index acl_title_type_index on acl (title(191), type(191))",
			"create index data_set_document_index on data_set (doc_name(191), set_name(191), doc_rev(191))",
			"create index data_set_list_index on data_set (set_name(191), doc_rev(191))",
			"create index data_title_index on data (title(191))",
			"create index back_title_type_index on back (title(191), type(191))",
			"create index back_type_title_index on back (type(191), title(191))",
			"create index back_link_type_index on back (link(191), type(191))",
			"create index user_set_id_name_index on user_set (id(191), name(191))",
			"create index topic_code_id_index on topic (code(191), id(191))",
		}
	}

	for _, query := range queries {
		_, err := db.Exec(DB_change(query))
		if err != nil {
			continue
		}
	}
}

func DB_make(db *sql.DB, new_db_set map[string]string) error {
	if new_db_set["db_type"] == "mysql" {
		DB_make_MySQL(db, new_db_set)
	} else {
		Exec_DB(
			db,
			`pragma journal_mode = WAL`,
		)
	}

	db = DB_connect()
	defer DB_close(db)

	if err := db.Ping(); err != nil {
		return err
	}

	field_text := DB_field_text(new_db_set["db_type"])

	table_list := DB_table_list()
	for table_name, table_data := range table_list {
		DB_create_table(db, table_name, field_text)

		columns := append([]string{"test"}, table_data...)

		for _, column_name := range columns {
			if !DB_column_exists(db, table_name, column_name) {
				DB_alter_add_column(db, table_name, column_name, field_text)
			}

			DB_warn_null_column(db, table_name, column_name)
		}
	}

	DB_create_index(db)

	return nil
}

func DB_init() {
	new_db_set := DB_boot()

	db, err := DB_connect_init()
	if err != nil {
		panic(fmt.Errorf("DB connection failed: %w", err))
	}
	defer DB_close(db)

	if err := DB_make(db, new_db_set); err != nil {
		panic(fmt.Errorf("DB setup failed: %w", err))
	}
}

func check_update_version() {
	DB_boot()

	db := DB_connect()
	defer DB_close(db)

	if !DB_column_exists(db, "other", "name") {
		return
	}

	now_version := ""
	if !QueryRow_DB(
		db,
		`select coalesce(data, '') from other where name = "ver"`,
		[]any{&now_version},
	) {
		log.Fatal("[DB ERROR] DB version is missing; run the Python version first")
	}

	last_version := Get_last_version()
	if now_version == "" || now_version == last_version["c_ver"] {
		return
	}

	if now_version == "20250527" || now_version == "20250529" || now_version == "20251201" {
		return
	}

	log.Fatalf(
		"[DB ERROR] unsupported DB version %s; run Python openNAMU v3.6.0-Beta-v56 or v58 first",
		now_version,
	)
}

func Main_init() {
	check_update_version()
	DB_init()
	DB_boot()

	db := DB_connect()
	defer DB_close(db)

	now_version := ""
	QueryRow_DB(
		db,
		`select data from other where name = "ver"`,
		[]any{&now_version},
	)

	last_version := Get_last_version()
	if now_version == "" {
		First_init(db)
	} else {
		if now_version != last_version["c_ver"] {
			log.Println("[DB] updating database")
			Update_init(db)
			log.Println("[DB] database update complete")
		}
	}

	Always_init(db, last_version["c_ver"])
	Get_ip_select(db)
}

func Get_last_version() map[string]string {
	if len(builtin_version_data) == 0 {
		panic("builtin version data is empty")
	}

	version_json := map[string]string{}
	if err := json.Unmarshal(builtin_version_data, &version_json); err != nil {
		panic(err)
	}

	return version_json
}

func First_init(db *sql.DB) {
	email := ""
	exists := QueryRow_DB(
		db,
		`select html from html_filter where kind = 'email'`,
		[]any{&email},
	)
	if !exists {
		for _, v := range []string{"naver.com", "gmail.com", "daum.net", "kakao.com"} {
			Exec_DB(
				db,
				`insert into html_filter (html, kind, plus, plus_t) values (?, 'email', '', '')`,
				v,
			)
		}
	}

	extension := ""
	exists = QueryRow_DB(
		db,
		`select html from html_filter where kind = 'extension'`,
		[]any{&extension},
	)
	if !exists {
		for _, v := range []string{"jpg", "jpeg", "png", "gif", "webp"} {
			Exec_DB(
				db,
				`insert into html_filter (html, kind, plus, plus_t) values (?, 'extension', '', '')`,
				v,
			)
		}
	}

	smtp_server := ""
	exists = QueryRow_DB(
		db,
		`select data from other where name = "smtp_server"`,
		[]any{&smtp_server},
	)
	if !exists {
		for _, v := range [][]string{
			{"smtp_server", "smtp.gmail.com"},
			{"smtp_port", "587"},
			{"smtp_security", "starttls"},
		} {
			Exec_DB(
				db,
				`insert into other (name, data, coverage) values (?, ?, '')`,
				v[0],
				v[1],
			)
		}
	}

	name_filter := ""
	exists = QueryRow_DB(
		db,
		`select html from html_filter where kind = 'name'`,
		[]any{&name_filter},
	)
	if !exists {
		Exec_DB(
			db,
			`insert into html_filter (html, kind, plus, plus_t) values (?, "name", "", "")`,
			`(?:[^A-Za-zㄱ-ㅣ가-힣0-9])`,
		)
	}
}

func legacy_acl_values(title string, acl_type string, value string) ([]string, bool) {
	if strings.HasPrefix(title, "user:") && acl_type == "decu" {
		switch value {
		case "all":
			return []string{"site_view"}, true
		case "ban":
			return []string{"owner"}, true
		}
	}

	switch value {
	case "", "normal", "all", "ban":
		return []string{}, true
	case "user":
		return []string{"user"}, true
	case "admin":
		return []string{"treat_as_admin"}, true
	case "owner":
		return []string{"owner"}, true
	case "email":
		return []string{"email_verified"}, true
	case "50_edit", "before":
		return []string{"trust_a"}, true
	case "30_day":
		return []string{"trust_b"}, true
	case "30_day_50_edit":
		return []string{"trust_c"}, true
	case "90_day":
		return []string{"trust_d"}, true
	case "ban_admin":
		return []string{"treat_as_admin", "ban", "ban_without_login"}, true
	case "not_all":
		return []string{"owner"}, true
	case "up_to_level_3", "up_to_level_10":
		return []string{value}, true
	}

	return nil, false
}

func legacy_acl_single_value(value string) (string, bool) {
	values, ok := legacy_acl_values("", "", value)
	if !ok {
		return "", false
	}
	if len(values) == 0 {
		return "", true
	}
	return values[0], true
}

func migrate_legacy_acl_rows(db *sql.DB, select_query string, update_query string, key_count int) {
	rows := Query_DB(db, select_query)
	data_list := [][]string{}
	for rows.Next() {
		data := make([]string, key_count+1)
		scan_list := []any{}
		for i := range data {
			scan_list = append(scan_list, &data[i])
		}
		if rows.Scan(scan_list...) == nil {
			data_list = append(data_list, data)
		}
	}
	rows.Close()

	for _, data := range data_list {
		value, ok := legacy_acl_single_value(data[key_count])
		if !ok {
			continue
		}

		values := []any{value}
		for i := 0; i < key_count; i++ {
			values = append(values, data[i])
		}
		Exec_DB(db, update_query, values...)
	}
}

func migrate_legacy_acl(db *sql.DB) {
	rows := Query_DB(db, "select title, data, type from acl")
	data_list := [][]string{}
	for rows.Next() {
		data := []string{"", "", ""}
		if rows.Scan(&data[0], &data[1], &data[2]) == nil {
			data_list = append(data_list, data)
		}
	}
	rows.Close()

	for _, data := range data_list {
		title := data[0]
		value := data[1]
		acl_type := data[2]
		if acl_type == "why" {
			continue
		}
		if acl_type == "document_edit_request_acl" {
			Exec_DB(db, "delete from acl where title = ? and data = ? and type = ?", title, value, acl_type)
			continue
		}

		values, ok := legacy_acl_values(title, acl_type, value)
		if !ok {
			continue
		}

		Exec_DB(db, "delete from acl where title = ? and data = ? and type = ?", title, value, acl_type)
		for _, new_value := range values {
			if new_value == "" {
				continue
			}
			exist := ""
			if !QueryRow_DB(db, "select data from acl where title = ? and data = ? and type = ? limit 1", []any{&exist}, title, new_value, acl_type) {
				Exec_DB(db, "insert into acl (title, data, type) values (?, ?, ?)", title, new_value, acl_type)
			}
		}
	}

	migrate_legacy_acl_rows(db, "select code, acl from rd", "update rd set acl = ? where code = ?", 1)
	migrate_legacy_acl_rows(db, "select thread_code, set_id, set_data from topic_set where set_name = 'thread_view_acl'", "update topic_set set set_data = ? where thread_code = ? and set_id = ? and set_name = 'thread_view_acl'", 2)
	migrate_legacy_acl_rows(db, "select set_id, set_name, set_code, set_data from bbs_set where set_name in ('bbs_view_acl', 'bbs_acl', 'bbs_edit_acl', 'bbs_comment_acl', 'bbs_view_acl_all', 'bbs_acl_all', 'bbs_edit_acl_all', 'bbs_comment_acl_all')", "update bbs_set set set_data = ? where set_id = ? and set_name = ? and set_code = ?", 3)
	migrate_legacy_acl_rows(db, "select name, coverage, data from other where name in ('bbs_view_acl_all', 'bbs_acl_all', 'bbs_edit_acl_all', 'bbs_comment_acl_all')", "update other set data = ? where name = ? and coverage = ?", 2)
	migrate_legacy_acl_rows(db, "select id, acl from vote where user = '' and type != 'option'", "update vote set acl = ? where id = ? and user = '' and type != 'option'", 1)
}

func Update_init(db *sql.DB) {
	migrate_legacy_acl(db)
}

func init_rankup_conditions(db *sql.DB) {
	initialized := ""
	if QueryRow_DB(
		db,
		"select data from other where name = 'rankup_condition_initialized' limit 1",
		[]any{&initialized},
	) {
		return
	}

	if len(Get_setting(db, "rankup_condition", "")) == 0 {
		for _, condition_data := range [][]string{
			{"trust_a", "edit 50"},
			{"trust_b", "time 30"},
			{"trust_c", "edit 50"},
			{"trust_c", "time 30"},
			{"trust_d", "edit 100"},
			{"trust_d", "time 90"},
		} {
			Exec_DB(
				db,
				"insert into other (name, data, coverage) values ('rankup_condition', ?, ?)",
				condition_data[1],
				condition_data[0],
			)
		}
	}

	Exec_DB(
		db,
		"insert into other (name, data, coverage) values ('rankup_condition_initialized', '1', '')",
	)
}

func init_bbs_comment_count(db *sql.DB) {
	initialized := ""
	if QueryRow_DB(
		db,
		"select data from other where name = 'bbs_comment_count_initialized' limit 1",
		[]any{&initialized},
	) {
		return
	}
	log.Println("[DB] updating BBS comment count")

	rows := Query_DB(
		db,
		"select set_id, set_code from bbs_data where set_name = 'title'",
	)
	post_list := [][]string{}
	for rows.Next() {
		post_data := []string{"", ""}
		if rows.Scan(&post_data[0], &post_data[1]) == nil {
			post_list = append(post_list, post_data)
		}
	}
	rows.Close()

	for _, post_data := range post_list {
		post_id := post_data[0] + "-" + post_data[1]
		comment_count := 0
		QueryRow_DB(
			db,
			"select count(*) from bbs_data where set_name = 'comment' and set_data != '' and (set_id = ? or set_id like ?)",
			[]any{&comment_count},
			post_id,
			post_id+"-%",
		)

		current := ""
		if QueryRow_DB(
			db,
			"select set_data from bbs_data where set_name = 'comment_count' and set_id = ? and set_code = ? limit 1",
			[]any{&current},
			post_data[0],
			post_data[1],
		) {
			Exec_DB(
				db,
				"update bbs_data set set_data = ? where set_name = 'comment_count' and set_id = ? and set_code = ?",
				strconv.Itoa(comment_count),
				post_data[0],
				post_data[1],
			)
		} else {
			Exec_DB(
				db,
				"insert into bbs_data (set_name, set_id, set_code, set_data) values ('comment_count', ?, ?, ?)",
				post_data[0],
				post_data[1],
				strconv.Itoa(comment_count),
			)
		}
	}

	Exec_DB(
		db,
		"insert into other (name, data, coverage) values ('bbs_comment_count_initialized', '1', '')",
	)
	log.Printf("[DB] BBS comment count update complete: %d posts", len(post_list))
}

func Always_init(db *sql.DB, version string) {
	// 버전 기입
	Exec_DB(
		db,
		`delete from other where name = "ver"`,
	)
	Exec_DB(
		db,
		`insert into other (name, data, coverage) values ("ver", ?, "")`,
		version,
	)

	// legacy ban management permissions
	legacy_groups := []string{}
	legacy_rows := Query_DB(
		db,
		"select distinct name from alist where acl in ('ban_manage', 'ban') and name not in ('ban', 'ban_without_login', 'ban_without_site')",
	)
	for legacy_rows.Next() {
		name := ""
		if legacy_rows.Scan(&name) == nil && name != "" {
			legacy_groups = append(legacy_groups, name)
		}
	}
	legacy_rows.Close()

	for _, name := range legacy_groups {
		for _, acl := range []string{"edit_filter_manage", "application_manage"} {
			current := ""
			QueryRow_DB(
				db,
				"select acl from alist where name = ? and acl = ? limit 1",
				[]any{&current},
				name,
				acl,
			)
			if current == "" {
				Exec_DB(
					db,
					"insert into alist (name, acl) values (?, ?)",
					name,
					acl,
				)
			}
		}
	}
	Exec_DB(db, "delete from alist where acl = 'ban_manage'")
	Exec_DB(db, "delete from alist where acl = 'ban' and name not in ('ban', 'ban_without_login', 'ban_without_site')")

	Exec_DB(
		db,
		`delete from alist where name = "owner"`,
	)
	Exec_DB(
		db,
		`insert into alist (name, acl) values ("owner", "owner")`,
	)

	for _, auth_data := range [][]string{{"admin", "admin"}, {"user", "user"}, {"ip", "ip"}} {
		acl := ""
		QueryRow_DB(
			db,
			"select acl from alist where name = ? and acl = ? limit 1",
			[]any{&acl},
			auth_data[0],
			auth_data[1],
		)
		if acl == "" {
			Exec_DB(
				db,
				"insert into alist (name, acl) values (?, ?)",
				auth_data[0],
				auth_data[1],
			)
		}
	}

	for _, ban_data := range [][]string{{"ip", "image_view"}, {"ban", "view"}, {"ban", "login_available"}, {"ban", "image_view"}, {"ban_without_login", "view"}, {"ban_without_login", "image_view"}, {"ban_without_site", "nothing"}} {
		ban_acl := ""
		QueryRow_DB(
			db,
			"select acl from alist where name = ? and acl = ? limit 1",
			[]any{&ban_acl},
			ban_data[0], ban_data[1],
		)
		if ban_acl == "" {
			Exec_DB(
				db,
				"insert into alist (name, acl) values (?, ?)",
				ban_data[0], ban_data[1],
			)
		}
	}
	for _, trust_name := range []string{"email_verified", "up_to_level_10", "up_to_level_3", "trust_a", "trust_b", "trust_c", "trust_d"} {
		trust_acl := ""
		QueryRow_DB(
			db,
			"select acl from alist where name = ? and acl = ? limit 1",
			[]any{&trust_acl},
			trust_name, trust_name,
		)
		if trust_acl == "" {
			Exec_DB(
				db,
				"insert into alist (name, acl) values (?, ?)",
				trust_name, trust_name,
			)
		}
	}
	init_rankup_conditions(db)
	init_bbs_comment_count(db)
	legacy_user_bans := [][]string{}
	rows := Query_DB(
		db,
		"select block, login, end from rb where (band = '' or band = 'private') and ongoing = '1' order by today",
	)
	for rows.Next() {
		ban_data := []string{"", "", ""}
		if rows.Scan(&ban_data[0], &ban_data[1], &ban_data[2]) == nil {
			legacy_user_bans = append(legacy_user_bans, ban_data)
		}
	}
	rows.Close()
	for _, ban_data := range legacy_user_bans {
		user_id := ""
		user_exists := IP_or_user(ban_data[0])
		if !user_exists {
			user_exists = QueryRow_DB(db, "select id from user_set where id = ? limit 1", []any{&user_id}, ban_data[0])
		}
		if !user_exists {
			continue
		}
		auth := get_ban_auth_group(db, ban_data[1])
		Exec_DB(db, "delete from user_set where id = ? and name = 'acl'", ban_data[0])
		Exec_DB(db, "insert into user_set (id, name, data) values (?, 'acl', ?)", ban_data[0], auth)
		Exec_DB(db, "delete from user_set where id = ? and name = 'acl_end'", ban_data[0])
		if ban_data[2] != "" && ban_data[2] != "release" {
			Exec_DB(db, "insert into user_set (id, name, data) values (?, 'acl_end', ?)", ban_data[0], ban_data[2])
		}
	}

	length := 0
	QueryRow_DB(
		db,
		`select count(*) from bbs_set where set_id = "0" and set_name = "bbs_name"`,
		[]any{&length},
	)

	if length > 1 {
		Exec_DB(
			db,
			`delete from bbs_set where set_id = "0" and set_name = "bbs_name"`,
		)
		Exec_DB(
			db,
			`delete from bbs_set where set_id = "0" and set_name = "bbs_type"`,
		)

		length = 0
	}

	if length == 0 {
		Exec_DB(
			db,
			`insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_name', '', '0', 'document_comment')`,
		)
		Exec_DB(
			db,
			`insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_type', '', '0', 'comment')`,
		)
	}

	image_url := Get_image_url(db)
	exists_folder := false

	_, err := os.Stat(image_url)
	if err == nil {
		exists_folder = true
	}

	if !exists_folder {
		if err := os.MkdirAll(image_url, 0o755); err != nil {
			panic(fmt.Errorf("image directory creation failed: %w", err))
		}
	}

	key := ""
	exists := QueryRow_DB(
		db,
		`select data from other where name = "key"`,
		[]any{&key},
	)
	if !exists {
		Exec_DB(
			db,
			`insert into other (name, data, coverage) values ("key", ?, "")`,
			Get_random_key(128),
		)
	}

	session_key := ""
	exists = QueryRow_DB(
		db,
		`select data from other where name = "session_key"`,
		[]any{&session_key},
	)
	if !exists {
		Exec_DB(
			db,
			`insert into other (name, data, coverage) values ("session_key", ?, "")`,
			Get_random_key(128),
		)
	}

	salt := ""
	exists = QueryRow_DB(
		db,
		`select data from other where name = "salt_key"`,
		[]any{&salt},
	)
	if !exists {
		Exec_DB(
			db,
			`insert into other (name, data, coverage) values ("salt_key", ?, "")`,
			Get_random_key(4),
		)
	}

	document_count := ""
	exists = QueryRow_DB(
		db,
		`select data from other where name = "count_all_title"`,
		[]any{&document_count},
	)
	if !exists {
		Exec_DB(
			db,
			`insert into other (name, data, coverage) values ("count_all_title", "0", "")`,
		)
	}
}
		"create index bbs_data_name_id_code_index on bbs_data (set_name, set_id, set_code)",
		"create index bbs_data_name_data_index on bbs_data (set_name, set_data)",
			"create index bbs_data_name_id_code_index on bbs_data (set_name(64), set_id(191), set_code(191))",
			"create index bbs_data_name_data_index on bbs_data (set_name(64), set_data(191))",
