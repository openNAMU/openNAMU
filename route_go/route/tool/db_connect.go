package tool

import (
	"database/sql"
	"encoding/json"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)

func Temp_DB_connect() *sql.DB {
	db, err := sql.Open("sqlite", "./data/temp.db")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func DB_connect() *sql.DB {
	m_db := Temp_DB_connect()
	defer m_db.Close()

	var db_set_str string

	err := m_db.QueryRow("select data from temp where name = 'db_set'").Scan(&db_set_str)
	if err != nil {
		if err == sql.ErrNoRows {
			db_set_str = "{}"
		} else {
			log.Fatal(err)
		}
	}

	db_set := map[string]string{}
	json.Unmarshal([]byte(db_set_str), &db_set)

	if db_set["type"] == "sqlite" {
		db, err := sql.Open("sqlite", db_set["name"]+".db")
		if err != nil {
			log.Fatal(err)
		}

		return db
	} else {
		db, err := sql.Open("mysql", db_set["mysql_user"]+":"+db_set["mysql_pw"]+"@tcp("+db_set["mysql_host"]+":"+db_set["mysql_port"]+")/"+db_set["name"])
		if err != nil {
			log.Fatal(err)
		}

		return db
	}
}

func DB_change(data string) string {
	m_db := Temp_DB_connect()
	defer m_db.Close()

	var db_set_type string

	err := m_db.QueryRow("select data from temp where name = 'db_set_type'").Scan(&db_set_type)
	if err != nil {
		if err == sql.ErrNoRows {
			db_set_type = "sqlite"
		} else {
			log.Fatal(err)
		}
	}

	if db_set_type == "mysql" {
		data = strings.Replace(data, "random()", "rand()", -1)
		data = strings.Replace(data, "collate nocase", "collate utf8mb4_general_ci", -1)
	}

	return data
}
