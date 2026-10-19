package route

import (
	"database/sql"
	"sync"
	"time"

	"opennamu/route/tool"
	"opennamu/route/tool/markup"
)

// rebuildStateLock protects rebuildRunning. rebuildRunning is true while a full
// rebuild is running so a second concurrent rebuild (same process) is rejected
// rather than double-wiping the index.
var (
	rebuildStateLock sync.Mutex
	rebuildRunning   bool
)

// backlink_backup_row is one original backlink row captured before the rebuild
// deletes all non-manual backlinks.
type backlink_backup_row struct {
	link     string
	title    string
	row_type string
	data     string
}

// backlink_backup holds the original non-manual backlinks and link_count rows,
// keyed by document (link / doc_name), so failed docs can be restored to their
// pre-rebuild state.
type backlink_backup struct {
	rows        map[string][]backlink_backup_row
	link_counts map[string]string
}

// Try_start_rebuild begins a full backlink rebuild if none is running (same
// process). It returns false — without holding anything — when a rebuild is
// already in progress, so the caller can reject the second request.
func Try_start_rebuild() bool {
	rebuildStateLock.Lock()
	defer rebuildStateLock.Unlock()
	if rebuildRunning {
		return false
	}
	rebuildRunning = true
	return true
}

// finish_rebuild clears the in-progress flag. Defer it immediately after
// Try_start_rebuild returns true so the flag is cleared on error or panic.
func finish_rebuild() {
	rebuildStateLock.Lock()
	defer rebuildStateLock.Unlock()
	rebuildRunning = false
}

// Take_backlink_backup captures all non-manual backlinks and link_count rows and
// deletes them inside a single transaction, so the snapshot and the wipe are
// atomic. Returns the snapshot used to restore failed docs.
func Take_backlink_backup(db *sql.DB) (*backlink_backup, error) {
	backup := &backlink_backup{
		rows:        make(map[string][]backlink_backup_row),
		link_counts: make(map[string]string),
	}
	err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		candidates, err := tx.Query(tool.DB_change("select link, title, type, data from back where type != 'cat_manual'"))
		if err != nil {
			return err
		}
		defer candidates.Close()
		for candidates.Next() {
			var row backlink_backup_row
			if err := candidates.Scan(&row.link, &row.title, &row.row_type, &row.data); err != nil {
				return err
			}
			backup.rows[row.link] = append(backup.rows[row.link], row)
		}
		// A non-nil error after Next() stops means the snapshot capture failed;
		// we must abort (rollback) before wiping, never delete on an incomplete
		// snapshot.
		if err := candidates.Err(); err != nil {
			return err
		}

		lc_rows, err := tx.Query(tool.DB_change("select doc_name, set_data from data_set where set_name = 'link_count'"))
		if err != nil {
			return err
		}
		defer lc_rows.Close()
		for lc_rows.Next() {
			doc_name := ""
			set_data := ""
			if err := lc_rows.Scan(&doc_name, &set_data); err != nil {
				return err
			}
			backup.link_counts[doc_name] = set_data
		}
		// Same contract as the back snapshot: abort before the wipe on any
		// error surfaced during the link_count iteration.
		if err := lc_rows.Err(); err != nil {
			return err
		}

		if _, err := tx.Exec(tool.DB_change("delete from back where type != 'cat_manual'")); err != nil {
			return err
		}
		_, err = tx.Exec(tool.DB_change("delete from data_set where set_name = 'link_count'"))
		return err
	})
	if err != nil {
		return nil, err
	}
	return backup, nil
}

// Restore_backlink_backup restores one failed doc to its original backlinks and
// link_count (captured by Take_backlink_backup). Deletes are case-sensitive
// (Exact_binary_eq) so case-sibling docs (Page / page) are never merged.
func Restore_backlink_backup(db *sql.DB, doc_name string, backup *backlink_backup) error {
	return tool.DB_transaction(db, func(tx *sql.Tx) error {
		if _, err := tx.Exec(tool.DB_change("delete from back where "+tool.Exact_binary_eq("link")+" and type != 'cat_manual'"), doc_name); err != nil {
			return err
		}
		if _, err := tx.Exec(tool.DB_change("delete from data_set where "+tool.Exact_binary_eq("doc_name")+" and set_name = 'link_count'"), doc_name); err != nil {
			return err
		}
		for _, row := range backup.rows[doc_name] {
			if _, err := tx.Exec(
				tool.DB_change("insert into back (link, title, type, data) values (?, ?, ?, ?)"),
				doc_name,
				row.title,
				row.row_type,
				row.data,
			); err != nil {
				return err
			}
		}
		if link_count, ok := backup.link_counts[doc_name]; ok {
			if _, err := tx.Exec(
				tool.DB_change("insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', 'link_count', ?)"),
				doc_name,
				link_count,
			); err != nil {
				return err
			}
		}
		return nil
	})
}

func Api_setting_backlink_reset_post(config tool.Config) map[string]any {
	return Api_setting_backlink_reset_post_internal(config, "")
}

func Api_setting_backlink_reset_post_internal(config tool.Config, load string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "setting_backlink", config.IP) {
		return_data["response"] = "require auth"

		return return_data
	}

	// Reject a concurrent full rebuild (same process). Two rebuilds wiping the
	// index at once would double-delete and race; only one may run at a time.
	if !Try_start_rebuild() {
		return_data["response"] = "partial_failure"
		return return_data
	}
	defer finish_rebuild()

	// Hold the common backlink mutex for the whole rebuild. Edits (and other
	// backlink re-indexes) acquire the same mutex, so general editing is blocked
	// while the index is being rebuilt, eliminating the rebuild-vs-edit race. The
	// flag above (released before the mutex is taken) is separate, so there is
	// no nested locking.
	backlink_mutex.Lock()
	defer backlink_mutex.Unlock()

	total_count := 0
	tool.QueryRow_DB(db, "select count(*) from data", []any{&total_count})

	// Enumerate every document BEFORE any destructive wipe. If enumeration fails
	// the index is untouched; report failure without tearing it down. A genuine
	// 0-doc database gets its stale auto backlinks wiped by the backup step.
	all_titles := []string{}
	page := 0
	for {
		rows := tool.Query_DB(db, "select title from data order by title limit ?, 50", page*50)
		title_list := []string{}
		for rows.Next() {
			doc_name := ""
			if err := rows.Scan(&doc_name); err != nil {
				rows.Close()
				return_data["response"] = "error"
				return return_data
			}
			title_list = append(title_list, doc_name)
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			return_data["response"] = "error"
			return return_data
		}
		rows.Close()

		all_titles = append(all_titles, title_list...)
		if len(title_list) < 50 {
			break
		}
		page++
	}

	backup, err := Take_backlink_backup(db)
	if err != nil {
		return_data["response"] = "error"
		return return_data
	}

	document_count := 0
	error_count := 0
	failed_docs := []string{}
	delay := map[string]time.Duration{
		"normal": 10 * time.Millisecond,
		"slow":   100 * time.Millisecond,
	}[load]

	for _, doc_name := range all_titles {
		raw_data := Api_w_raw(config, doc_name, "", "")
		response, ok := raw_data["response"].(string)
		if !ok || response != "ok" {
			error_count++
			failed_docs = append(failed_docs, doc_name)
			continue
		}
		parameter_data := map[string]any{}
		parameter_data["__opennamu_skin_set"] = Get_render_setting_parameter(db, config)
		markup.Get_render(db, doc_name, raw_data["data"].(string), "backlink", parameter_data)
		document_count++
		if delay > 0 {
			time.Sleep(delay)
		}
	}

	backup_err := false
	// Restore each failed doc. A doc deleted while the rebuild ran has no
	// content; restoring it would re-create backlinks for a deleted doc, so it
	// is skipped.
	for _, doc_name := range failed_docs {
		var exists int
		tool.QueryRow_DB(db,
			"select count(*) from data where title = ?",
			[]any{&exists}, doc_name)
		if exists == 0 {
			continue
		}
		if err := Restore_backlink_backup(db, doc_name, backup); err != nil {
			backup_err = true
		}
	}

	if backup_err {
		return_data["response"] = "error"
	} else if total_count == 0 {
		// Genuine 0-doc: the backup step wiped the stale auto backlinks and
		// there is nothing left to rebuild; nothing was lost.
		return_data["response"] = "ok"
	} else if error_count == 0 {
		return_data["response"] = "ok"
	} else if document_count == 0 {
		// Documents existed but every rebuild failed.
		return_data["response"] = "failure"
	} else {
		return_data["response"] = "partial_failure"
	}
	return_data["document_total"] = total_count
	return_data["document_count"] = document_count
	return_data["error_count"] = error_count

	return return_data
}
