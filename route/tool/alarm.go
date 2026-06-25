package tool

import (
	"database/sql"
)

func Send_alarm(db *sql.DB, from string, target string, data string) {
    if from != target {
        data = from + " | " + data

        now_time := Get_time()

        count := "1"
        QueryRow_DB(
            db,
            "select id from user_notice where name = ? order by id + 0 desc limit 1",
            []any{ &count },
            target,
        )

        count_int := Str_to_int(count)
        count_int += 1

        Exec_DB(
            db,
            "insert into user_notice (id, name, data, date, readme) values (?, ?, ?, ?, '')",
            count_int, target, data, now_time,
        )
    }
}
