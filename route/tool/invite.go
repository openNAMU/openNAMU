package tool

import (
	"database/sql"
	"strings"
	"time"
)

type Invite_record struct {
	Hash   string `json:"-"`
	Issuer string `json:"issuer"`
	Date   string `json:"date"`
	End    string `json:"end"`
}

func Invite_hash(invite string) string {
	return Sha224(strings.TrimSpace(invite))
}

func Invite_required(db *sql.DB) bool {
	if Get_setting_value(db, "invite_required", "", "") != "on" {
		return false
	}

	id := ""
	return QueryRow_DB(db, "select id from user_set limit 1", []any{&id})
}

func Invite_valid(db DB_runner, invite_hash string) bool {
	invite_hash = strings.TrimSpace(invite_hash)
	if invite_hash == "" {
		return false
	}

	data := ""
	if !QueryRow_DB(
		db,
		"select data from user_set where id = ? and name = 'invite' limit 1",
		[]any{&data},
		invite_hash,
	) {
		return false
	}

	record := Invite_record{}
	if json.Unmarshal([]byte(data), &record) != nil {
		return false
	}
	if record.End == "" {
		return true
	}

	end, err := time.ParseInLocation("2006-01-02 15:04:05", record.End, time.Local)
	return err == nil && time.Now().Before(end)
}

func Invite_consume(tx *sql.Tx, invite_hash string) bool {
	if !Invite_valid(tx, invite_hash) {
		return false
	}

	result, err := tx.Exec(
		DB_change("delete from user_set where id = ? and name = 'invite'"),
		strings.TrimSpace(invite_hash),
	)
	if err != nil {
		return false
	}

	count, err := result.RowsAffected()
	return err == nil && count == 1
}

func Get_invite_list(db *sql.DB) []Invite_record {
	rows := Query_DB(db, "select id, data from user_set where name = 'invite'")
	defer rows.Close()

	data_list := []Invite_record{}
	for rows.Next() {
		record := Invite_record{}
		data := ""
		if rows.Scan(&record.Hash, &data) != nil {
			continue
		}
		if json.Unmarshal([]byte(data), &record) != nil {
			continue
		}
		data_list = append(data_list, record)
	}

	return data_list
}
