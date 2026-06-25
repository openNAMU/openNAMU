package tool

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
)

var global_lang_data = map[string]string{}

func Get_language(db *sql.DB, data string, safe bool) string {
    language := "ko-KR"
    QueryRow_DB(
        db,
        "select data from other where name = 'language'",
        []any{ &language },
    )

    if _, ok := global_lang_data[language + "_" + data]; ok {
        if safe {
            return global_lang_data[language + "_" + data]
        } else {
            return HTML_escape(global_lang_data[language + "_" + data])
        }
    } else {
        file, err := os.Open(filepath.Join("..", "lang", language + ".json"))
        if err != nil {
            panic(err)
        }
        defer file.Close()

        lang_data := map[string]string{}

        decoder := json.NewDecoder(file)
        if err := decoder.Decode(&lang_data); err != nil {
            panic(err)
        }

        for k, v := range lang_data {
            global_lang_data[language + "_" + k] = v
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
