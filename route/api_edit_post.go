package route

import (
	"database/sql"
	"strings"

	"opennamu/route/tool"
	"opennamu/route/tool/markup"
)

func Edit_post_content_validation_error(db *sql.DB, config tool.Config, doc_name string, data string, send string, agree string) string {
	if !tool.Do_edit_filter(db, config, doc_name, data) {
		return "edit filter (content)"
	}
	if !tool.Do_edit_filter(db, config, doc_name, send) {
		return "edit filter (send)"
	}
	if !tool.Do_edit_send_require_check(db, config, send) {
		return "send require"
	}
	if !tool.Do_edit_text_checkbox_check(db, config, agree) {
		return "checkbox check require"
	}
	if !tool.Do_edit_max_length_check(db, config, data) {
		return "overflow max length"
	}
	return ""
}

func Api_edit_post(config tool.Config, doc_name string, data string, send string, agree string, expected_revision string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	date := tool.Get_time()
	data = strings.ReplaceAll(data, "\r", "")
	data = tool.Do_edit_replace(db, data)
	current_revision := tool.Get_document_revision(db, doc_name)

	if doc_name == "" {
		return map[string]any{"response": "error", "data": "empty title"}
	}
	if !tool.Do_title_length_check(db, doc_name, "document") {
		return map[string]any{"response": "error", "data": "title length"}
	}
	if !tool.Check_acl(db, doc_name, "", "document_edit", config.IP) {
		return map[string]any{"response": "error", "data": "permission denied"}
	}
	if current_revision == "0" && !tool.Check_acl(db, doc_name, "", "document_make_acl", config.IP) {
		return map[string]any{"response": "error", "data": "permission denied"}
	}
	if !tool.Check_edit_time(db, config.IP) {
		return map[string]any{"response": "error", "data": "edit time"}
	}
	if !tool.Do_edit_slow_check(db, config, "edit") {
		return map[string]any{"response": "error", "data": "slow edit limit"}
	}
	if !tool.Check_daily_limit(db, config.IP, "edit") {
		return map[string]any{"response": "error", "data": "daily limit"}
	}
	if error_name := Edit_post_content_validation_error(db, config, doc_name, data, send, agree); error_name != "" {
		return map[string]any{"response": "error", "data": error_name}
	}

	// The document save (content commit) and its backlink index update are one
	// critical section. The common mutex keeps a concurrent edit from committing
	// between saving content and writing backlinks. The holder owns the whole
	// section, so nothing inside re-enters it (no nested locking).
	backlink_mutex.Lock()
	defer backlink_mutex.Unlock()

	// Re-check the revision under the lock, just before the commit. The
	// pre-lock read above is only needed for the new-doc make_acl check; the
	// commit-time revision may have changed, so it must be validated here.
	current_revision = tool.Get_document_revision(db, doc_name)
	if expected_revision != "" && expected_revision != current_revision {
		return map[string]any{"response": "error", "data": "edit conflict"}
	}

	var old_data string
	tool.QueryRow_DB(
		db,
		`select data from data where title = ?`,
		[]any{&old_data},
		doc_name,
	)

	length := tool.Get_edit_length_diff(old_data, data)

	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		if _, err := tx.Exec(tool.DB_change(`delete from data where title = ?`), doc_name); err != nil {
			return err
		}
		if _, err := tx.Exec(tool.DB_change(`insert into data (title, data) values (?, ?)`), doc_name, data); err != nil {
			return err
		}
		tool.Do_add_history(tx, doc_name, data, date, config.IP, send, length, "", "")
		return nil
	}); err != nil {
		panic(err)
	}
	tool.Search_index_update(doc_name, data)

	tool.Do_watchlist_alarm_send(db, config, doc_name)

	markup.Get_render(db, doc_name, data, "backlink")
	Notify_indexnow(db, doc_name)

	return map[string]any{"response": "ok"}
}
