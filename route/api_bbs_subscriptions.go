package route

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"

	"opennamu/route/tool"
)

const (
	bbs_bookmark_list_type    = "bbs_bookmark"
	bbs_board_watch_list_type = "bbs_board_watchlist"
	bbs_tag_watch_list_type   = "bbs_tag_watchlist"
)

func Bbs_subscription_board_id_valid(set_id string) bool {
	if set_id == "0" || set_id == thread_bbs_id || set_id == report_bbs_id {
		return true
	}
	board_number, err := strconv.ParseInt(set_id, 10, 64)
	return err == nil && board_number > 0
}

func Bbs_subscription_post_id_valid(set_id string, set_code string) bool {
	if !Bbs_subscription_board_id_valid(set_id) {
		return false
	}
	post_number, err := strconv.ParseInt(set_code, 10, 64)
	return err == nil && post_number >= 0
}

func Bbs_tag_watch_key(set_id string, tag string) string {
	return set_id + "\n" + tag
}

func Bbs_tag_watch_data_matches(data string, tag string) bool {
	if data == tag {
		return true
	}
	_, legacy_tag, is_legacy := strings.Cut(data, "\n")
	return is_legacy && legacy_tag == tag
}

func Bbs_tag_watch_rows(db tool.DB_runner, user_id string, tag string) []string {
	rows := tool.Query_DB(
		db,
		"select data from user_set where name = ? and id = ?",
		bbs_tag_watch_list_type,
		user_id,
	)
	defer rows.Close()

	watch_data_list := []string{}
	for rows.Next() {
		watch_data := ""
		if err := rows.Scan(&watch_data); err != nil {
			panic(err)
		}
		if Bbs_tag_watch_data_matches(watch_data, tag) {
			watch_data_list = append(watch_data_list, watch_data)
		}
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return watch_data_list
}

func Bbs_tag_watch_transaction(db *sql.DB, user_id string, tag string, callback func(*sql.Tx) error) error {
	if tool.Get_DB_type() == "mysql" {
		lock_user_id, err := Bbs_tag_watch_lock_user_id(db, user_id)
		if err != nil {
			return err
		}
		return Bbs_mysql_identifier_transaction(db, []string{"tag-watch:" + lock_user_id + ":" + tag}, callback)
	}
	return tool.DB_transaction_retry(db, func(tx *sql.Tx) error {
		if err := Bbs_local_tag_watch_write_lock(tx, user_id); err != nil {
			return err
		}
		return callback(tx)
	})
}

// MySQL identifier comparisons can be case-insensitive, so use the stored user ID
// to make differently-cased sessions for one account share the same lock scope.
func Bbs_tag_watch_lock_user_id(db *sql.DB, user_id string) (string, error) {
	var canonical_user_id string
	err := db.QueryRow(
		tool.DB_change("select id from user_set where name = 'acl' and id = ? limit 1"),
		user_id,
	).Scan(&canonical_user_id)
	if errors.Is(err, sql.ErrNoRows) {
		return user_id, nil
	}
	if err != nil {
		return "", err
	}
	return canonical_user_id, nil
}

func Bbs_local_tag_watch_write_lock(tx *sql.Tx, user_id string) error {
	result, err := tx.Exec(
		tool.DB_change("update user_set set data = data || char(1) where name = 'acl' and id = ?"),
		user_id,
	)
	if err != nil {
		return err
	}
	rows_affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows_affected > 0 {
		_, err = tx.Exec(
			tool.DB_change("update user_set set data = substr(data, 1, length(data) - 1) where name = 'acl' and id = ?"),
			user_id,
		)
		return err
	}

	// A normal authenticated account has a stable per-user ACL row. Keep the
	// initialized global anchor as a safe fallback for older or incomplete data.
	return Bbs_local_board_write_lock(tx, report_bbs_id)
}

func Bbs_tag_watch_value_valid(tag string) bool {
	tag = strings.TrimSpace(tag)
	if tag == "" || tool.Get_len(tag) > bbs_tag_max_length {
		return false
	}

	tag_list := Bbs_tag_list(tag)
	return len(tag_list) == 1 && tag_list[0] == tag
}

func Bbs_subscription_toggle(db *sql.DB, user_id string, watch_type string, set_id string, set_code string, value string) error {
	toggle := func(tx *sql.Tx) error {
		return User_watch_list_toggle(tx, user_id, watch_type, value)
	}

	return Bbs_board_write_transaction(db, set_id, set_code, toggle)
}

func Api_bbs_bookmark_post(config tool.Config, set_id string, set_code string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if tool.IP_or_user(config.IP) {
		return map[string]any{"response": "require auth"}
	}
	if !Bbs_subscription_post_id_valid(set_id, set_code) {
		return map[string]any{"response": "not exist", "data": "post"}
	}

	key := Bbs_watch_key(set_id, set_code)
	has_bookmark := tool.Get_user_set_data_match(db, config.IP, bbs_bookmark_list_type, key)
	_, post_exists := tool.Get_bbs_data_value(db, set_id, set_code, "title")
	_, can_view := Bbs_post_view_auth(db, set_id, set_code, config.IP)
	if !post_exists || !can_view {
		if !has_bookmark {
			if !post_exists {
				return map[string]any{"response": "not exist", "data": "post"}
			}
			return map[string]any{"response": "require auth"}
		}
		if err := tool.DB_transaction_retry(db, func(tx *sql.Tx) error {
			_, err := tx.Exec(
				tool.DB_change("delete from user_set where name = ? and id = ? and data = ?"),
				bbs_bookmark_list_type,
				config.IP,
				key,
			)
			return err
		}); err != nil {
			panic(err)
		}
		return map[string]any{"response": "ok"}
	}

	if err := Bbs_subscription_toggle(db, config.IP, bbs_bookmark_list_type, set_id, set_code, key); err != nil {
		if errors.Is(err, bbs_board_missing_error) {
			return map[string]any{"response": "not exist", "data": "bbs"}
		}
		panic(err)
	}
	return map[string]any{"response": "ok"}
}

func Api_bbs_bookmark_view(config tool.Config, set_id string, set_code string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if tool.IP_or_user(config.IP) {
		return map[string]any{"response": "require auth"}
	}
	if !Bbs_subscription_post_id_valid(set_id, set_code) {
		return map[string]any{"response": "not exist"}
	}
	key := Bbs_watch_key(set_id, set_code)
	bookmarked := tool.Get_user_set_data_match(db, config.IP, bbs_bookmark_list_type, key)
	bbs_name := ""
	title := ""
	bbs_exists := tool.QueryRow_DB(
		db,
		"select set_data from bbs_set where set_name = 'bbs_name' and set_id = ?",
		[]any{&bbs_name},
		set_id,
	)
	post_exists := tool.QueryRow_DB(
		db,
		"select set_data from bbs_data where set_name = 'title' and set_id = ? and set_code = ?",
		[]any{&title},
		set_id,
		set_code,
	)
	_, can_view := Bbs_post_view_auth(db, set_id, set_code, config.IP)
	if !bbs_exists || !post_exists || !can_view {
		if !bookmarked {
			if !bbs_exists || !post_exists {
				return map[string]any{"response": "not exist"}
			}
			return map[string]any{"response": "require auth"}
		}
		return map[string]any{
			"response": "ok",
			"data": map[string]string{
				"available": "false",
				"watched":   "true",
			},
		}
	}

	return map[string]any{
		"response": "ok",
		"data": map[string]string{
			"available": "true",
			"bbs_name":  bbs_name,
			"title":     title,
			"watched":   strconv.FormatBool(bookmarked),
		},
	}
}

func Api_bbs_board_watch_post(config tool.Config, set_id string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if tool.IP_or_user(config.IP) {
		return map[string]any{"response": "require auth"}
	}
	if !Bbs_subscription_board_id_valid(set_id) {
		return map[string]any{"response": "not exist", "data": "bbs"}
	}
	watching := tool.Get_user_set_data_match(db, config.IP, bbs_board_watch_list_type, set_id)
	bbs_name := ""
	bbs_exists := tool.QueryRow_DB(
		db,
		"select set_data from bbs_set where set_name = 'bbs_name' and set_id = ?",
		[]any{&bbs_name},
		set_id,
	)
	can_view := bbs_exists && tool.Check_acl(db, set_id, "", "bbs_view", config.IP)
	if !can_view {
		if !watching {
			if !bbs_exists {
				return map[string]any{"response": "not exist", "data": "bbs"}
			}
			return map[string]any{"response": "require auth"}
		}
		if err := tool.DB_transaction_retry(db, func(tx *sql.Tx) error {
			_, err := tx.Exec(
				tool.DB_change("delete from user_set where name = ? and id = ? and data = ?"),
				bbs_board_watch_list_type,
				config.IP,
				set_id,
			)
			return err
		}); err != nil {
			panic(err)
		}
		return map[string]any{"response": "ok"}
	}

	if err := Bbs_subscription_toggle(db, config.IP, bbs_board_watch_list_type, set_id, "", set_id); err != nil {
		if errors.Is(err, bbs_board_missing_error) {
			return map[string]any{"response": "not exist", "data": "bbs"}
		}
		panic(err)
	}
	return map[string]any{"response": "ok"}
}

func Api_bbs_board_watch_view(config tool.Config, set_id string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if tool.IP_or_user(config.IP) {
		return map[string]any{"response": "require auth"}
	}
	if !Bbs_subscription_board_id_valid(set_id) {
		return map[string]any{"response": "not exist"}
	}
	watching := tool.Get_user_set_data_match(db, config.IP, bbs_board_watch_list_type, set_id)
	bbs_name := ""
	bbs_exists := tool.QueryRow_DB(
		db,
		"select set_data from bbs_set where set_name = 'bbs_name' and set_id = ?",
		[]any{&bbs_name},
		set_id,
	)
	if !bbs_exists || !tool.Check_acl(db, set_id, "", "bbs_view", config.IP) {
		if !watching {
			if !bbs_exists {
				return map[string]any{"response": "not exist"}
			}
			return map[string]any{"response": "require auth"}
		}
		return map[string]any{
			"response": "ok",
			"data": map[string]string{
				"available": "false",
				"watched":   "true",
			},
		}
	}
	return map[string]any{
		"response": "ok",
		"data": map[string]string{
			"available": "true",
			"bbs_name":  bbs_name,
			"watched":   strconv.FormatBool(watching),
		},
	}
}

func Api_bbs_tag_watch_post(config tool.Config, set_id string, tag string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if tool.IP_or_user(config.IP) {
		return map[string]any{"response": "require auth"}
	}
	tag = strings.TrimSpace(tag)
	if !Bbs_tag_watch_value_valid(tag) {
		return map[string]any{"response": "error", "data": "tag"}
	}
	if set_id != "" {
		if !Bbs_subscription_board_id_valid(set_id) {
			return map[string]any{"response": "not exist", "data": "bbs"}
		}
		bbs_name := ""
		if !tool.QueryRow_DB(
			db,
			"select set_data from bbs_set where set_name = 'bbs_name' and set_id = ?",
			[]any{&bbs_name},
			set_id,
		) {
			return map[string]any{"response": "not exist", "data": "bbs"}
		}
		if !tool.Check_acl(db, set_id, "", "bbs_view", config.IP) {
			return map[string]any{"response": "require auth"}
		}
	}

	if err := Bbs_tag_watch_transaction(db, config.IP, tag, func(tx *sql.Tx) error {
		watch_data_list := Bbs_tag_watch_rows(tx, config.IP, tag)
		if len(watch_data_list) > 0 {
			delete_query := "delete from user_set where name = ? and id = ? and data = ?"
			if tool.Get_DB_type() == "mysql" {
				delete_query = "delete from user_set where name = ? and id = ? and binary data = binary ?"
			}
			for _, watch_data := range watch_data_list {
				if _, err := tx.Exec(
					tool.DB_change(delete_query),
					bbs_tag_watch_list_type,
					config.IP,
					watch_data,
				); err != nil {
					return err
				}
			}
			return nil
		}
		_, err := tx.Exec(
			tool.DB_change("insert into user_set (id, name, data) values (?, ?, ?)"),
			config.IP,
			bbs_tag_watch_list_type,
			tag,
		)
		return err
	}); err != nil {
		panic(err)
	}
	return map[string]any{"response": "ok"}
}

func Api_bbs_tag_watch_view(config tool.Config, set_id string, tag string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if tool.IP_or_user(config.IP) {
		return map[string]any{"response": "require auth"}
	}
	tag = strings.TrimSpace(tag)
	if !Bbs_tag_watch_value_valid(tag) {
		return map[string]any{"response": "not exist"}
	}
	bbs_name := ""
	if set_id != "" {
		if !Bbs_subscription_board_id_valid(set_id) {
			return map[string]any{"response": "not exist"}
		}
		if !tool.QueryRow_DB(
			db,
			"select set_data from bbs_set where set_name = 'bbs_name' and set_id = ?",
			[]any{&bbs_name},
			set_id,
		) {
			return map[string]any{"response": "not exist"}
		}
		if !tool.Check_acl(db, set_id, "", "bbs_view", config.IP) {
			return map[string]any{"response": "require auth"}
		}
	}
	watching := len(Bbs_tag_watch_rows(db, config.IP, tag)) > 0
	return map[string]any{
		"response": "ok",
		"data": map[string]string{
			"available": "true",
			"bbs_name":  bbs_name,
			"tag":       tag,
			"watched":   strconv.FormatBool(watching),
		},
	}
}

func Bbs_post_watch_notify(db *sql.DB, config tool.Config, set_id string, set_code string, bbs_name string, title string, post_user string, tags []string) {
	targets := []string{}
	seen := map[string]bool{}
	mark_existing_watchers := func(watch_type string, value string) error {
		rows, err := db.Query(
			tool.DB_change("select id, data from user_set where name = ? and data = ?"),
			watch_type,
			value,
		)
		if err != nil {
			return err
		}
		for rows.Next() {
			target := ""
			stored_value := ""
			if err := rows.Scan(&target, &stored_value); err != nil {
				rows.Close()
				return err
			}
			if stored_value == value && target != "" {
				seen[target] = true
			}
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			return err
		}
		return rows.Close()
	}
	append_watchers := func(watch_type string, value string) error {
		rows, err := db.Query(
			tool.DB_change("select id, data from user_set where name = ? and data = ?"),
			watch_type,
			value,
		)
		if err != nil {
			return err
		}
		for rows.Next() {
			target := ""
			stored_value := ""
			if err := rows.Scan(&target, &stored_value); err != nil {
				rows.Close()
				return err
			}
			if stored_value != value || target == "" || seen[target] {
				continue
			}
			seen[target] = true
			targets = append(targets, target)
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			return err
		}
		return rows.Close()
	}

	if set_id == thread_bbs_id {
		thread_watch_key := Bbs_watch_key(set_id, set_code)
		for _, watch_type := range []string{"bbs_watchlist", "thread_watchlist"} {
			watch_values := []string{thread_watch_key}
			if watch_type == "thread_watchlist" {
				watch_values = append(watch_values, set_code)
			}
			for _, watch_value := range watch_values {
				if err := mark_existing_watchers(watch_type, watch_value); err != nil {
					log.Printf("[BBS WATCH] existing thread watcher lookup failed: %v", err)
					return
				}
			}
		}
	}

	if err := append_watchers(bbs_board_watch_list_type, set_id); err != nil {
		log.Printf("[BBS WATCH] board watcher lookup failed: %v", err)
		return
	}
	tag_set := map[string]bool{}
	for _, tag := range tags {
		tag_set[tag] = true
	}
	if len(tag_set) > 0 {
		rows, err := db.Query(
			tool.DB_change("select id, data from user_set where name = ?"),
			bbs_tag_watch_list_type,
		)
		if err != nil {
			log.Printf("[BBS WATCH] tag watcher lookup failed: %v", err)
			return
		}
		for rows.Next() {
			target := ""
			watch_data := ""
			if err := rows.Scan(&target, &watch_data); err != nil {
				rows.Close()
				log.Printf("[BBS WATCH] tag watcher scan failed: %v", err)
				return
			}
			matched := tag_set[watch_data]
			if !matched {
				_, legacy_tag, is_legacy := strings.Cut(watch_data, "\n")
				matched = is_legacy && tag_set[legacy_tag]
			}
			if !matched || target == "" || seen[target] {
				continue
			}
			seen[target] = true
			targets = append(targets, target)
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			log.Printf("[BBS WATCH] tag watcher lookup failed: %v", err)
			return
		}
		if err := rows.Close(); err != nil {
			log.Printf("[BBS WATCH] tag watcher cursor close failed: %v", err)
			return
		}
	}

	for _, target := range targets {
		if target == post_user {
			continue
		}
		var allowed bool
		if err := Bbs_notification_db_call(func() {
			_, allowed = Bbs_post_view_auth(db, set_id, set_code, target)
		}); err != nil {
			log.Printf("[BBS WATCH] visibility lookup failed for recipient %q: %v", target, err)
			continue
		}
		if !allowed {
			continue
		}

		alarm := `BBS <a href="/bbs/w/` + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code) + `">` + tool.HTML_escape(bbs_name) + " - " + tool.HTML_escape(title) + `</a>`
		if err := Bbs_notification_db_call(func() {
			tool.Send_alarm(db, config.IP, target, alarm)
		}); err != nil {
			log.Printf("[BBS WATCH] notification delivery failed for recipient %q: %v", target, err)
		}
	}
}
