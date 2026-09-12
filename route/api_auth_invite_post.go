package route

import (
	"database/sql"
	"time"

	"opennamu/route/tool"
)

func Api_auth_invite_post(config tool.Config, action string, value string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := map[string]any{}
	if !tool.Check_permission(db, "user_manage", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	switch action {
	case "make":
		token := tool.Get_random_key(32)
		record := tool.Invite_record{
			Issuer: config.IP,
			Date:   tool.Get_time(),
			End:    time.Now().AddDate(0, 1, 0).Format("2006-01-02 15:04:05"),
		}
		raw, err := json.Marshal(record)
		if err != nil {
			return_data["response"] = "error"
			return return_data
		}

		if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
			if _, err := tx.Exec(
				tool.DB_change("insert into user_set (id, name, data) values (?, 'invite', ?)"),
				tool.Invite_hash(token),
				string(raw),
			); err != nil {
				return err
			}
			tool.Do_insert_auth_history(tx, config.IP, "invite_make")
			return nil
		}); err != nil {
			return_data["response"] = "error"
			return return_data
		}

		return_data["response"] = "ok"
		return_data["token"] = token
		return return_data

	case "revoke":
		revoked := false
		if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
			result, err := tx.Exec(
				tool.DB_change("delete from user_set where id = ? and name = 'invite'"),
				value,
			)
			if err != nil {
				return err
			}
			count, err := result.RowsAffected()
			if err != nil {
				return err
			}
			if count == 1 {
				tool.Do_insert_auth_history(tx, config.IP, "invite_revoke")
				revoked = true
			}
			return nil
		}); err != nil {
			return_data["response"] = "error"
			return return_data
		}
		if !revoked {
			return_data["response"] = "error"
			return return_data
		}
		return_data["response"] = "ok"
		return return_data
	}

	return_data["response"] = "error"
	return return_data
}
