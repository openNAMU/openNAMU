package route

import (
	"database/sql"
	"errors"
	"strconv"

	"opennamu/route/tool"
)

func Bbs_tabom_vote_type_valid(vote_type string) bool {
	return vote_type == "up" || vote_type == "down"
}

func Bbs_tabom_vote_state(db tool.DB_runner, user string, set_id string, set_code string) (string, error) {
	statement, err := db.Prepare(tool.DB_change("select set_name from bbs_data where set_name in ('tabom_list', 'tabom_down_list') and set_data = ? and set_id = ? and set_code = ?"))
	if err != nil {
		return "", err
	}
	defer statement.Close()

	rows, err := statement.Query(user, set_id, set_code)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	has_up := false
	has_down := false
	for rows.Next() {
		var set_name string
		if err := rows.Scan(&set_name); err != nil {
			return "", err
		}
		switch set_name {
		case "tabom_list":
			has_up = true
		case "tabom_down_list":
			has_down = true
		}
	}
	if err := rows.Err(); err != nil {
		return "", err
	}

	return Bbs_tabom_vote_state_from_presence(has_up, has_down), nil
}

func Bbs_tabom_vote_state_from_presence(has_up bool, has_down bool) string {
	if has_up && has_down {
		return "both"
	}
	if has_up {
		return "up"
	}
	if has_down {
		return "down"
	}
	return "none"
}

func Bbs_tabom_unique_vote_count(tx *sql.Tx, list_name string, set_id string, set_code string) (int, error) {
	rows, err := tx.Query(
		tool.DB_change("select set_data from bbs_data where set_name = ? and set_id = ? and set_code = ?"),
		list_name,
		set_id,
		set_code,
	)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	users := map[string]struct{}{}
	for rows.Next() {
		var user sql.NullString
		if err := rows.Scan(&user); err != nil {
			return 0, err
		}
		if user.Valid {
			users[user.String] = struct{}{}
		}
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	return len(users), nil
}

func Bbs_tabom_counter_reconcile(tx *sql.Tx, counter_name string, set_id string, set_code string, count int) error {
	var counter_exists int
	err := tx.QueryRow(
		tool.DB_change("select 1 from bbs_data where set_name = ? and set_id = ? and set_code = ? limit 1"),
		counter_name,
		set_id,
		set_code,
	).Scan(&counter_exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if errors.Is(err, sql.ErrNoRows) {
		if count == 0 {
			return nil
		}
		_, err = tx.Exec(
			tool.DB_change("insert into bbs_data (set_name, set_data, set_id, set_code) values (?, ?, ?, ?)"),
			counter_name,
			strconv.Itoa(count),
			set_id,
			set_code,
		)
		return err
	}

	_, err = tx.Exec(
		tool.DB_change("update bbs_data set set_data = ? where set_name = ? and set_id = ? and set_code = ?"),
		strconv.Itoa(count),
		counter_name,
		set_id,
		set_code,
	)
	return err
}

func Api_bbs_tabom_post(tx *sql.Tx, user string, set_id string, set_code string, vote_type string) (map[string]any, error) {
	if !Bbs_tabom_vote_type_valid(vote_type) {
		return map[string]any{"response": "error", "data": "vote_type"}, nil
	}

	current_vote, err := Bbs_tabom_vote_state(tx, user, set_id, set_code)
	if err != nil {
		return nil, err
	}

	next_vote := vote_type
	switch current_vote {
	case vote_type:
		next_vote = "none"
	case "both":
		// Legacy dual votes are repaired to the requested single direction.
	}

	if _, err := tx.Exec(
		tool.DB_change("delete from bbs_data where set_name in ('tabom_list', 'tabom_down_list') and set_data = ? and set_id = ? and set_code = ?"),
		user,
		set_id,
		set_code,
	); err != nil {
		return nil, err
	}

	if next_vote != "none" {
		list_name := "tabom_list"
		if next_vote == "down" {
			list_name = "tabom_down_list"
		}
		if _, err := tx.Exec(
			tool.DB_change("insert into bbs_data (set_name, set_data, set_id, set_code) values (?, ?, ?, ?)"),
			list_name,
			user,
			set_id,
			set_code,
		); err != nil {
			return nil, err
		}
	}

	up_count, err := Bbs_tabom_unique_vote_count(tx, "tabom_list", set_id, set_code)
	if err != nil {
		return nil, err
	}
	down_count, err := Bbs_tabom_unique_vote_count(tx, "tabom_down_list", set_id, set_code)
	if err != nil {
		return nil, err
	}
	if err := Bbs_tabom_counter_reconcile(tx, "tabom_count", set_id, set_code, up_count); err != nil {
		return nil, err
	}
	if err := Bbs_tabom_counter_reconcile(tx, "tabom_down_count", set_id, set_code, down_count); err != nil {
		return nil, err
	}

	return map[string]any{
		"response": "ok",
		"data": map[string]any{
			"vote":       next_vote,
			"up_count":   up_count,
			"down_count": down_count,
		},
	}, nil
}
