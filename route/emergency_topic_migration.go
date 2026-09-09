package route

import (
	"database/sql"
	"fmt"
	"strconv"

	"opennamu/route/tool"
)

const thread_bbs_set_id = "-1"

type emergency_topic_comment struct {
	id    string
	data  string
	date  string
	ip    string
	block string
	top   string
}

type emergency_topic_row struct {
	title    string
	sub      string
	code     string
	date     string
	band     string
	stop     string
	agree    string
	comments []emergency_topic_comment
}

func emergency_topic_exec(db *sql.DB, query string, values ...any) error {
	_, err := db.Exec(tool.DB_change(query), values...)
	return err
}

func emergency_topic_insert(tx *sql.Tx, set_name string, set_code string, set_id string, set_data string) error {
	_, err := tx.Exec(
		tool.DB_change("insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)"),
		set_name,
		set_code,
		set_id,
		set_data,
	)
	return err
}

func emergency_topic_bbs_init(db *sql.DB) error {
	var name_count int
	err := db.QueryRow(
		tool.DB_change("select count(*) from bbs_set where set_id = ? and set_name = 'bbs_name'"),
		thread_bbs_set_id,
	).Scan(&name_count)
	if err != nil {
		return err
	}
	if name_count > 1 {
		return fmt.Errorf("duplicate thread BBS set: %s", thread_bbs_set_id)
	}
	if name_count == 0 {
		if err := emergency_topic_exec(
			db,
			"insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_name', '', ?, 'thread')",
			thread_bbs_set_id,
		); err != nil {
			return err
		}
	} else if err := emergency_topic_exec(
		db,
		"update bbs_set set set_data = 'thread' where set_id = ? and set_name = 'bbs_name'",
		thread_bbs_set_id,
	); err != nil {
		return err
	}

	var type_count int
	err = db.QueryRow(
		tool.DB_change("select count(*) from bbs_set where set_id = ? and set_name = 'bbs_type'"),
		thread_bbs_set_id,
	).Scan(&type_count)
	if err != nil {
		return err
	}
	if type_count > 1 {
		return fmt.Errorf("duplicate thread BBS type: %s", thread_bbs_set_id)
	}
	if type_count == 0 {
		return emergency_topic_exec(
			db,
			"insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_type', '', ?, 'thread')",
			thread_bbs_set_id,
		)
	}

	if err := emergency_topic_exec(
		db,
		"update bbs_set set set_data = 'thread' where set_id = ? and set_name = 'bbs_type'",
		thread_bbs_set_id,
	); err != nil {
		return err
	}
	return emergency_topic_prefix_init(db)
}

func emergency_topic_prefix_init(db *sql.DB) error {
	var count int
	if err := db.QueryRow(
		tool.DB_change("select count(*) from bbs_set where set_id = ? and set_name = 'bbs_prefix'"),
		thread_bbs_set_id,
	).Scan(&count); err != nil {
		return err
	}
	if count > 1 {
		return fmt.Errorf("duplicate thread BBS prefix: %s", thread_bbs_set_id)
	}
	if count == 0 {
		return emergency_topic_exec(
			db,
			"insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_prefix', '', ?, ?)",
			thread_bbs_set_id,
			"열림\n닫힘\n합의",
		)
	}
	return nil
}

func emergency_topic_comments(db *sql.DB, topic_code string) ([]emergency_topic_comment, error) {
	rows, err := db.Query(
		tool.DB_change("select coalesce(id, ''), coalesce(data, ''), coalesce(date, ''), coalesce(ip, ''), coalesce(block, ''), coalesce(top, '') from topic where code = ? order by id + 0 asc"),
		topic_code,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []emergency_topic_comment{}
	for rows.Next() {
		comment := emergency_topic_comment{}
		if err := rows.Scan(&comment.id, &comment.data, &comment.date, &comment.ip, &comment.block, &comment.top); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func emergency_migrate_topic_to_bbs(db *sql.DB) (return_err error) {
	if err := emergency_topic_bbs_init(db); err != nil {
		return err
	}

	rows, err := db.Query(
		tool.DB_change("select coalesce(title, ''), coalesce(sub, ''), coalesce(code, ''), coalesce(date, ''), coalesce(band, ''), coalesce(stop, ''), coalesce(agree, '') from rd order by code + 0 asc"),
	)
	if err != nil {
		return err
	}

	topics := []emergency_topic_row{}
	for rows.Next() {
		topic := emergency_topic_row{}
		if err := rows.Scan(&topic.title, &topic.sub, &topic.code, &topic.date, &topic.band, &topic.stop, &topic.agree); err != nil {
			rows.Close()
			return err
		}
		topics = append(topics, topic)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}
	rows.Close()

	index_changed := false
	defer func() {
		if !index_changed {
			return
		}
		if err := tool.Search_bbs_index_mark_rebuild(); err != nil && return_err == nil {
			return_err = err
		}
	}()

	migrated_count := 0
	skipped_count := 0
	for _, topic := range topics {
		if topic.code == "" {
			continue
		}

		var source_count int
		err := db.QueryRow(
			tool.DB_change("select count(*) from bbs_data where set_name = 'topic_source' and set_id = ? and set_code = ?"),
			thread_bbs_set_id,
			topic.code,
		).Scan(&source_count)
		if err != nil {
			return err
		}
		if source_count > 0 {
			skipped_count++
			continue
		}

		var existing_count int
		err = db.QueryRow(
			tool.DB_change("select count(*) from bbs_data where set_id = ? and set_code = ?"),
			thread_bbs_set_id,
			topic.code,
		).Scan(&existing_count)
		if err != nil {
			return err
		}
		if existing_count > 0 {
			return fmt.Errorf("thread BBS post already exists without migration marker: %s", topic.code)
		}

		comments, err := emergency_topic_comments(db, topic.code)
		if err != nil {
			return err
		}
		topic.comments = comments

		user_id := ""
		if len(topic.comments) > 0 {
			user_id = topic.comments[0].ip
			if topic.date == "" {
				topic.date = topic.comments[0].date
			}
		}

		tx, err := db.Begin()
		if err != nil {
			return err
		}

		prefix := "열림"
		if topic.stop == "O" {
			prefix = "닫힘"
		}
		if topic.agree == "O" {
			prefix = "합의"
		}
		root_data := []struct {
			name string
			data string
		}{
			{"title", topic.sub},
			{"data", ""},
			{"date", topic.date},
			{"last_activity", topic.date},
			{"user_id", user_id},
			{"comment_count", strconv.Itoa(len(topic.comments))},
			{"document", topic.title},
			{"topic_source", "rd"},
			{"prefix", prefix},
		}
		if topic.band != "" {
			root_data = append(root_data, struct {
				name string
				data string
			}{"topic_band", topic.band})
		}
		for _, data := range root_data {
			if err := emergency_topic_insert(tx, data.name, topic.code, thread_bbs_set_id, data.data); err != nil {
				_ = tx.Rollback()
				return err
			}
		}

		for _, comment := range topic.comments {
			comment_data := []struct {
				name string
				data string
			}{
				{"comment", comment.data},
				{"comment_date", comment.date},
				{"comment_user_id", comment.ip},
			}
			if comment.block != "" {
				comment_data = append(comment_data, struct {
					name string
					data string
				}{"blind", comment.block})
			}
			if comment.top != "" {
				comment_data = append(comment_data, struct {
					name string
					data string
				}{"top", comment.top})
			}
			for _, data := range comment_data {
				if err := emergency_topic_insert(tx, data.name, comment.id, thread_bbs_set_id+"-"+topic.code, data.data); err != nil {
					_ = tx.Rollback()
					return err
				}
			}
		}

		if err := tx.Commit(); err != nil {
			return err
		}
		index_changed = true
		migrated_count++
		if migrated_count%100 == 0 {
			fmt.Printf("Topic migration: %d/%d\n", migrated_count, len(topics))
		}
	}

	if err := normalize_thread_bbs_prefix(db); err != nil {
		return err
	}

	fmt.Printf("Topic migration complete: migrated=%d skipped=%d total=%d\n", migrated_count, skipped_count, len(topics))
	return nil
}
