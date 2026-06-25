package route

import (
	"opennamu/route/tool"
)

func Api_give_auth_patch(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    new_data := make(map[string]any)

    ip := config.IP
    user_name := other_set["user_name"]

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
            auth_data = tool.Get_auth_group_info(db, other_set["change_auth"])

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
                    user_name, other_set["change_auth"],
                )
                
                new_data["response"] = "ok"
            }
        }
    } else {
        auth_check := false
        auth_data := tool.Get_auth_group_info(db, other_set["auth"])

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
            auth_data = tool.Get_auth_group_info(db, other_set["change_auth"])

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
                    other_set["change_auth"], other_set["auth"],
                )

                new_data["response"] = "ok"
            }
        }
    }

    json_data, _ := json.Marshal(new_data)
    return string(json_data)
}
