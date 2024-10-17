package route

import (
    "log"
    "opennamu/route/tool"
    "strconv"

    jsoniter "github.com/json-iterator/go"
)

func Api_bbs_w_post(call_arg []string) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(call_arg[0]), &other_set)

    db := tool.DB_connect()
    defer db.Close()

    if tool.Check_acl(db, "", "", "bbs_comment", other_set["ip"]) {
        return_data := make(map[string]interface{})
        return_data["response"] = "require auth"

        json_data, _ := json.Marshal(return_data)
        return string(json_data)
    }
    
    stmt, err := db.Prepare(tool.DB_change("select set_code from bbs_data where set_name = 'title' and set_id = ? order by set_code + 0 desc"))
    if err != nil {
        log.Fatal(err)
    }
    defer stmt.Close()

    var set_code string

    err = stmt.QueryRow(other_set["set_id"]).Scan(&set_code)
    if err != nil {
        log.Fatal(err)
    }

    set_code_int, _ := strconv.Atoi(set_code)
    set_code_int += 1

    set_code_str := strconv.Itoa(set_code_int)

    date_now := tool.Get_time()

    insert_db := [][]string{
        []string{"title", other_set["title"]},
        []string{"data", other_set["data"]},
        []string{"date", date_now},
        []string{"user_id", other_set["ip"]},
    }
    for _, v := range insert_db {
        stmt, err := db.Prepare(tool.DB_change("insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)"))
        if err != nil {
            log.Fatal(err)
        }
        defer stmt.Close()

        _, err = stmt.Exec(v[0], set_code_str, other_set["set_id"], v[1])
        if err != nil {
            log.Fatal(err)
        }
    }

    return_data := make(map[string]interface{})
    return_data["response"] = "ok"
    return_data["data"] = set_code_str

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}