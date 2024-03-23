package tool

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
)

func Get_language(db *sql.DB, db_set map[string]string, data string, safe bool) string {
	var language string

	err := db.QueryRow(DB_change(db_set, "select data from other where name = 'language'")).Scan(&language)
	if err != nil {
		if err == sql.ErrNoRows {
			language = "ko-KR"
		} else {
			log.Fatal(err)
		}
	}

	file, err := os.Open("./lang/" + language + ".json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lang_data := map[string]string{}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&lang_data); err != nil {
		return ""
	}

	if safe {
		return lang_data[data]
	} else {
		return HTML_escape(lang_data[data])
	}
}
