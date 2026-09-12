package route

import (
	"database/sql"
	stdjson "encoding/json"
	"strconv"
	"strings"

	"opennamu/route/tool"
	"opennamu/route/tool/markup"
)

type document_backup_history struct {
	ID     string `json:"id"`
	Data   string `json:"data"`
	Date   string `json:"date"`
	IP     string `json:"ip"`
	Send   string `json:"send"`
	Length string `json:"length"`
	Hide   string `json:"hide"`
	Type   string `json:"type"`
}

type document_backup_document struct {
	Title   string                    `json:"title"`
	Data    string                    `json:"data"`
	Type    string                    `json:"type"`
	History []document_backup_history `json:"history"`
}

type document_backup_file struct {
	Version   int                        `json:"version"`
	Documents []document_backup_document `json:"documents"`
}

func document_backup_string(value sql.NullString) string {
	if value.Valid {
		return value.String
	}

	return ""
}

func Api_document_backup_export(config tool.Config) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "admin", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	backup := document_backup_file{
		Version:   1,
		Documents: []document_backup_document{},
	}
	rows := tool.Query_DB(db, "select title, data, type from data order by title")
	defer rows.Close()

	for rows.Next() {
		var title sql.NullString
		var data sql.NullString
		var document_type sql.NullString
		if err := rows.Scan(&title, &data, &document_type); err != nil {
			panic(err)
		}

		document := document_backup_document{
			Title:   document_backup_string(title),
			Data:    document_backup_string(data),
			Type:    document_backup_string(document_type),
			History: []document_backup_history{},
		}
		history_rows := tool.Query_DB(
			db,
			"select id, data, date, ip, send, leng, hide, type from history where title = ? order by id + 0 asc",
			document.Title,
		)

		for history_rows.Next() {
			var id sql.NullString
			var history_data sql.NullString
			var date sql.NullString
			var ip sql.NullString
			var send sql.NullString
			var length sql.NullString
			var hide sql.NullString
			var history_type sql.NullString
			if err := history_rows.Scan(&id, &history_data, &date, &ip, &send, &length, &hide, &history_type); err != nil {
				history_rows.Close()
				panic(err)
			}

			document.History = append(document.History, document_backup_history{
				ID:     document_backup_string(id),
				Data:   document_backup_string(history_data),
				Date:   document_backup_string(date),
				IP:     document_backup_string(ip),
				Send:   document_backup_string(send),
				Length: document_backup_string(length),
				Hide:   document_backup_string(hide),
				Type:   document_backup_string(history_type),
			})
		}
		history_rows.Close()
		backup.Documents = append(backup.Documents, document)
	}

	raw_data, err := stdjson.MarshalIndent(backup, "", "  ")
	if err != nil {
		return_data["response"] = "error"
		return return_data
	}

	return_data["response"] = "ok"
	return_data["data"] = raw_data
	return return_data
}

func Api_document_backup_import(config tool.Config, raw_data []byte) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "admin", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	backup := document_backup_file{}
	if err := stdjson.Unmarshal(raw_data, &backup); err != nil || backup.Version != 1 {
		return_data["response"] = "error"
		return return_data
	}

	seen := map[string]bool{}
	for _, document := range backup.Documents {
		if document.Title == "" || strings.ContainsAny(document.Title, "\r\n") || seen[document.Title] {
			return_data["response"] = "error"
			return return_data
		}
		seen[document.Title] = true
		for _, history := range document.History {
			if history.ID == "" {
				return_data["response"] = "error"
				return return_data
			}
		}
	}

	imported := 0
	skipped := 0
	imported_titles := []string{}
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		for _, document := range backup.Documents {
			var current_title string
			err := tx.QueryRow(
				tool.DB_change("select title from data where title = ? limit 1"),
				document.Title,
			).Scan(&current_title)
			if err == nil {
				skipped++
				continue
			}
			if err != sql.ErrNoRows {
				return err
			}

			if _, err := tx.Exec(
				tool.DB_change("insert into data (title, data, type) values (?, ?, ?)"),
				document.Title,
				document.Data,
				document.Type,
			); err != nil {
				return err
			}

			last_id := ""
			history_exists := tx.QueryRow(
				tool.DB_change("select id from history where title = ? order by id + 0 desc limit 1"),
				document.Title,
			).Scan(&last_id) == nil
			next_id := tool.Str_to_int(last_id) + 1
			for _, history := range document.History {
				history_id := history.ID
				if history_exists {
					history_id = strconv.Itoa(next_id)
					next_id++
				}
				if _, err := tx.Exec(
					tool.DB_change("insert into history (id, title, data, date, ip, send, leng, hide, type) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"),
					history_id,
					document.Title,
					history.Data,
					history.Date,
					history.IP,
					history.Send,
					history.Length,
					history.Hide,
					history.Type,
				); err != nil {
					return err
				}
			}

			imported++
			imported_titles = append(imported_titles, document.Title)
		}
		return nil
	}); err != nil {
		return_data["response"] = "error"
		return return_data
	}

	for _, title := range imported_titles {
		data := ""
		if tool.QueryRow_DB(db, "select data from data where title = ?", []any{&data}, title) {
			markup.Get_render(db, title, data, "backlink")
			tool.Search_index_sync(db, title)
		}
	}

	return_data["response"] = "ok"
	return_data["imported"] = imported
	return_data["skipped"] = skipped
	return return_data
}
