package route

import "opennamu/route/tool"

func Api_add_user(config tool.Config, id string, password string, email string, encode string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if encode == "" {
		encode = tool.Get_main_encode(db)
	}

	password_hash := tool.Password_encode(db, password, encode)

	temp := ""
	tool.QueryRow_DB(
		db,
		`select id from user_set limit 1`,
		[]any{&temp},
	)

	auth := "user"
	if temp == "" {
		auth = "owner"
	}

	tool.Exec_DB(
		db,
		`insert into user_set (id, name, data) values (?, 'pw', ?)`,
		id,
		password_hash,
	)
	tool.Exec_DB(
		db,
		`insert into user_set (id, name, data) values (?, 'acl', ?)`,
		id,
		auth,
	)
	tool.Exec_DB(
		db,
		`insert into user_set (id, name, data) values (?, 'date', ?)`,
		id,
		tool.Get_time(),
	)
	tool.Exec_DB(
		db,
		`insert into user_set (id, name, data) values (?, 'encode', ?)`,
		id,
		encode,
	)
	if email != "" {
		tool.Exec_DB(
			db,
			`insert into user_set (id, name, data) values (?, 'email', ?)`,
			id,
			email,
		)
	}

	return_data := make(map[string]any)
	return_data["response"] = "ok"

	return return_data
}

func Api_add_user_hash(config tool.Config, id string, password_hash string, email string, encode string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if encode == "" {
		encode = tool.Get_main_encode(db)
	}

	temp := ""
	tool.QueryRow_DB(
		db,
		`select id from user_set limit 1`,
		[]any{&temp},
	)

	auth := "user"
	if temp == "" {
		auth = "owner"
	}

	tool.Exec_DB(
		db,
		`insert into user_set (id, name, data) values (?, 'pw', ?)`,
		id,
		password_hash,
	)
	tool.Exec_DB(
		db,
		`insert into user_set (id, name, data) values (?, 'acl', ?)`,
		id,
		auth,
	)
	tool.Exec_DB(
		db,
		`insert into user_set (id, name, data) values (?, 'date', ?)`,
		id,
		tool.Get_time(),
	)
	tool.Exec_DB(
		db,
		`insert into user_set (id, name, data) values (?, 'encode', ?)`,
		id,
		encode,
	)
	if email != "" {
		tool.Exec_DB(
			db,
			`insert into user_set (id, name, data) values (?, 'email', ?)`,
			id,
			email,
		)
	}

	return_data := make(map[string]any)
	return_data["response"] = "ok"

	return return_data
}
