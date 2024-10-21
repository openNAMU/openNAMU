package tool

import (
    "database/sql"
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

    _, err = db.Exec("pragma journal_mode = WAL")
    if err != nil {
        log.Fatal(err)
    }

    return db
}

func DB_connect() *sql.DB {
    m_db := Temp_DB_connect()
    defer m_db.Close()

    db_set := map[string]string{}

    rows, err := m_db.Query("select name, data from temp where name in ('db_name', 'db_type')")
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

    if db_set["db_type"] == "sqlite" {
        db, err := sql.Open("sqlite", db_set["db_name"]+".db")
        if err != nil {
            log.Fatal(err)
        }

        _, err = db.Exec("pragma journal_mode = WAL")
        if err != nil {
            log.Fatal(err)
        }

        return db
    } else {
        rows, err := m_db.Query("select name, data from temp where name in ('db_mysql_host', 'db_mysql_user', 'db_mysql_pw', 'db_mysql_port')")
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

        db, err := sql.Open("mysql", db_set["db_mysql_user"]+":"+db_set["db_mysql_pw"]+"@tcp("+db_set["db_mysql_host"]+":"+db_set["db_mysql_port"]+")/"+db_set["db_name"])
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

    err := m_db.QueryRow("select data from temp where name = 'db_type'").Scan(&db_set_type)
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
