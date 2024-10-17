package tool

import (
    "database/sql"
    "log"
    "strconv"
)

func Send_alarm(db *sql.DB, from string, target string, data string) {
    if from != target {
        data = from + " | " + data

        now_time := Get_time()

        var count string

        stmt, err := db.Prepare(DB_change("select id from user_notice where name = ? order by id + 0 desc limit 1"))
        if err != nil {
            log.Fatal(err)
        }
        defer stmt.Close()

        err = stmt.QueryRow(target).Scan(&count)
        if err != nil {
            if err == sql.ErrNoRows {
                count = "1"
            } else {
                log.Fatal(err)
            }
        }

        count_int, _ := strconv.Atoi(count)
        count_int += 1

        stmt, err = db.Prepare(DB_change("insert into user_notice (id, name, data, date, readme) values (?, ?, ?, ?, '')"))
        if err != nil {
            log.Fatal(err)
        }
        defer stmt.Close()

        _, err = stmt.Exec(count_int, target, data, now_time)
        if err != nil {
            log.Fatal(err)
        }
    }
}
