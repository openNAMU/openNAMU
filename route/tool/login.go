package tool

import (
	"crypto/sha256"
	"crypto/sha3"
	"database/sql"
	"encoding/hex"
)

func Get_main_encode(db *sql.DB) string {
	encode := ""
	QueryRow_DB(
		db,
		`select data from other where name = "encode"`,
		[]any{&encode},
	)

	if encode == "" {
		encode = "sha3"
	}

	return encode
}

func Password_encode(db *sql.DB, password string, encode string) string {
	if encode == "" {
		encode = Get_main_encode(db)
	}

	switch encode {
	case "sha256":
		hash_data := sha256.Sum256([]byte(password))
		return hex.EncodeToString(hash_data[:])
	case "sha3":
		hash_data := sha3.Sum256([]byte(password))
		return hex.EncodeToString(hash_data[:])
	case "sha3-512":
		hash_data := sha3.Sum512([]byte(password))
		return hex.EncodeToString(hash_data[:])
	}

	salt_key := ""
	QueryRow_DB(
		db,
		`select data from other where name = "salt_key"`,
		[]any{&salt_key},
	)

	salted_password := password + salt_key

	if encode == "sha3-salt" {
		hash_data := sha3.Sum256([]byte(salted_password))
		return hex.EncodeToString(hash_data[:])
	}

	hash_data := sha3.Sum512([]byte(salted_password))
	return hex.EncodeToString(hash_data[:])
}

func Password_check(db *sql.DB, id string, password string) bool {
	db_password := ""
	db_encode := ""

	QueryRow_DB(
		db,
		`select data from user_set where id = ? and name = 'pw'`,
		[]any{&db_password},
		id,
	)

	QueryRow_DB(
		db,
		`select data from user_set where id = ? and name = 'encode'`,
		[]any{&db_encode},
		id,
	)

	password_encode := Password_encode(db, password, db_encode)

	return db_password == password_encode
}
