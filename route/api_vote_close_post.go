package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Api_vote_close_post(config tool.Config, id string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)

	type_data := ""
	if !tool.QueryRow_DB(
		db,
		"select type from vote where id = ? and user = ''",
		[]any{&type_data},
		id,
	) {
		return_data["response"] = "not exist"
		return_data["data"] = "vote"
		return return_data
	}

	if !tool.Check_permission(db, "vote", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	owner := ""
	tool.QueryRow_DB(
		db,
		"select data from vote where id = ? and name = 'open_user' and type = 'option'",
		[]any{&owner},
		id,
	)
	if owner != config.IP && !tool.Check_permission(db, "vote_manage", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	next := "n_close"
	if type_data == "n_close" {
		next = "n_open"
	} else if type_data == "close" {
		next = "open"
	} else if type_data == "open" {
		next = "close"
	}

	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		if _, err := tx.Exec(
			tool.DB_change("update vote set type = ? where id = ? and user = ''"),
			next,
			id,
		); err != nil {
			return err
		}
		if next == "open" || next == "n_open" {
			_, err := tx.Exec(
				tool.DB_change("delete from vote where id = ? and name = 'end_date' and type = 'option'"),
				id,
			)
			return err
		}
		return nil
	}); err != nil {
		panic(err)
	}

	return_data["response"] = "ok"
	return_data["data"] = next
	return return_data
}
