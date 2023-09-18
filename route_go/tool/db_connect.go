package tool

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)

func DB_connect(db_set map[string]string) *sql.DB {
	if db_set["type"] == "sqlite" {
		db, err := sql.Open("sqlite", db_set["name"]+".db")
		if err != nil {
			fmt.Println(err)
			return nil
		}

		return db
	} else {
		db, err := sql.Open("mysql", db_set["mysql_user"]+":"+db_set["mysql_pw"]+"@tcp("+db_set["mysql_host"]+":"+db_set["mysql_port"]+")/"+db_set["name"])
		if err != nil {
			fmt.Println(err)
			return nil
		}

		return db
	}
}

func DB_change(db_set map[string]string, data string) string {
	if db_set["type"] == "mysql" {
		data = strings.Replace(data, "random()", "rand()", -1)
		data = strings.Replace(data, "collate nocase", "collate utf8mb4_general_ci", -1)
	}

	return data
}
