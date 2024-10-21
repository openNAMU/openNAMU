package tool

import (
    "database/sql"
    "log"
    "strings"

    _ "github.com/go-sql-driver/mysql"
    _ "modernc.org/sqlite"
)

var db_set = map[string]string{}

func Temp_DB_connect() *sql.DB {
    db, err := sql.Open("sqlite", "./data/temp.db")
    if err != nil {
        log.Fatal(err)
    }

    return db
}

func DB_init() {
    m_db := Temp_DB_connect()
    defer m_db.Close()

    rows, err := m_db.Query("select name, data from temp")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    for rows.Next() {
        var name string
        var data string

        err := rows.Scan(&name, &data)
        if err != nil {
            log.Fatal(err)
        }

        db_set[name] = data
    }
}

func DB_connect() *sql.DB {
    if db_set["db_type"] == "sqlite" {
        db, err := sql.Open("sqlite", db_set["db_name"]+".db")
        if err != nil {
            log.Fatal(err)
        }

        return db
    } else {
        db, err := sql.Open("mysql", db_set["db_mysql_user"]+":"+db_set["db_mysql_pw"]+"@tcp("+db_set["db_mysql_host"]+":"+db_set["db_mysql_port"]+")/"+db_set["db_name"])
        if err != nil {
            log.Fatal(err)
        }

        return db
    }
}

func Get_DB_type() string {
    return db_set["db_type"]
}

func DB_change(data string) string {
    if Get_DB_type() == "mysql" {
        data = strings.Replace(data, "random()", "rand()", -1)
        data = strings.Replace(data, "collate nocase", "collate utf8mb4_general_ci", -1)
    }

    return data
}
