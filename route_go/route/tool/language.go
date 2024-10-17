package tool

import (
    "database/sql"
    "log"
    "os"

    jsoniter "github.com/json-iterator/go"
)

func Get_language(db *sql.DB, data string, safe bool) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    m_db := Temp_DB_connect()
    defer m_db.Close()

    var language string

    err := db.QueryRow(DB_change("select data from other where name = 'language'")).Scan(&language)
    if err != nil {
        if err == sql.ErrNoRows {
            language = "ko-KR"
        } else {
            log.Fatal(err)
        }
    }

    var language_data string

    stmt, err := m_db.Prepare("select data from temp where name = ?")
    if err != nil {
        log.Fatal(err)
    }
    defer stmt.Close()

    err = stmt.QueryRow("lang_" + language + "_" + data).Scan(&language_data)
    if err != nil {
        if err == sql.ErrNoRows {
            language_data = ""
        } else {
            log.Fatal(err)
        }
    }

    if language_data != "" {
        if safe {
            return language_data
        } else {
            return HTML_escape(language_data)
        }
    } else {
        file, err := os.Open("./lang/" + language + ".json")
        if err != nil {
            log.Fatal(err)
        }
        defer file.Close()

        lang_data := map[string]string{}

        decoder := json.NewDecoder(file)
        if err := decoder.Decode(&lang_data); err != nil {
            log.Fatal(err)
        }

        if _, ok := lang_data[data]; ok {
            if safe {
                return lang_data[data]
            } else {
                return HTML_escape(lang_data[data])
            }
        } else {
            log.Default().Println(data + " (" + language + ")")

            return data + " (" + language + ")"
        }
    }
}
