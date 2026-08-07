package route

import (
	"opennamu/route/tool"
)

func Api_give_auth_patch(config tool.Config, auth string, change_auth string, user_name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	new_data := make(map[string]any)

	ip := config.IP

	if user_name != "" {
		auth_check := false
		auth_name := tool.Get_user_auth(db, user_name)
		auth_data := tool.Get_auth_group_info(db, auth_name)

		if tool.Auth_include_upper_auth(auth_data) {
			if tool.Check_acl(db, "", "", "owner_auth", ip) {
				auth_check = true
			}
		} else {
			if tool.Check_acl(db, "", "", "give_auth", ip) {
				auth_check = true
			}
		}

		if !auth_check {
			new_data["response"] = "require auth"
		} else {
			auth_check = false
			auth_data = tool.Get_auth_group_info(db, change_auth)

			if tool.Auth_include_upper_auth(auth_data) {
				if tool.Check_acl(db, "", "", "owner_auth", ip) {
					auth_check = true
				}
			} else {
				if tool.Check_acl(db, "", "", "give_auth", ip) {
					auth_check = true
				}
			}

			if !auth_check {
				new_data["response"] = "require auth"
			} else {
				tool.Exec_DB(
					db,
					"delete from user_set where id = ? and name = 'acl'",
					user_name,
				)

				tool.Exec_DB(
					db,
					"insert into user_set (id, name, data) values (?, 'acl', ?)",
					user_name, change_auth,
				)

				new_data["response"] = "ok"
			}
		}
	} else {
		auth_check := false
		auth_data := tool.Get_auth_group_info(db, auth)

		if tool.Auth_include_upper_auth(auth_data) {
			if tool.Check_acl(db, "", "", "owner_auth", config.IP) {
				auth_check = true
			}
		} else {
			if tool.Check_acl(db, "", "", "give_auth", config.IP) {
				auth_check = true
			}
		}

		if !auth_check {
			new_data["response"] = "require auth"
		} else {
			auth_check = false
			auth_data = tool.Get_auth_group_info(db, change_auth)

			if tool.Auth_include_upper_auth(auth_data) {
				if tool.Check_acl(db, "", "", "owner_auth", config.IP) {
					auth_check = true
				}
			} else {
				if tool.Check_acl(db, "", "", "give_auth", config.IP) {
					auth_check = true
				}
			}

			if !auth_check {
				new_data["response"] = "require auth"
			} else {
				tool.Exec_DB(
					db,
					"update user_set set data = ? where name = 'acl' and data = ?",
					change_auth, auth,
				)

				new_data["response"] = "ok"
			}
		}
	}

	return new_data
}
