package route

import (
    "database/sql"
    "log"

    "opennamu/route/tool"

    jsoniter "github.com/json-iterator/go"
)

func Api_w_random(call_arg []string) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    db := tool.DB_connect()
    defer db.Close()

    var title string

    err := db.QueryRow(tool.DB_change("select title from data where title not like 'user:%' and title not like 'category:%' and title not like 'file:%' order by random() limit 1")).Scan(&title)
    if err != nil {
        if err == sql.ErrNoRows {
            title = ""
        } else {
            log.Fatal(err)
        }
    }

    new_data := map[string]string{}
    new_data["data"] = title

    json_data, _ := json.Marshal(new_data)
    return string(json_data)
}
