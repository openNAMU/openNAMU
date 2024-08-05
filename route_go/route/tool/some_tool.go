package tool

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"html/template"
	"log"
	"net/url"
	"time"
)

func Sha224(data string) string {
	hasher := sha256.New224()
	hasher.Write([]byte(data))
	hash_byte := hasher.Sum(nil)
	hash_str := hex.EncodeToString(hash_byte)

	return hash_str
}

func Url_parser(data string) string {
	return url.QueryEscape(data)
}

func HTML_escape(data string) string {
	return template.HTMLEscapeString(data)
}

func Arr_in_str(arr []string, data string) bool {
	for _, v := range arr {
		if v == data {
			return true
		}
	}

	return false
}

func Get_time() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func Get_date() string {
	return time.Now().Format("2006-01-02")
}

func Get_setting(db *sql.DB, set_name string, data_coverage string) [][]string {
	var rows *sql.Rows

	if data_coverage != "" {
		stmt, err := db.Prepare(DB_change("select data, coverage from other where name = ? and coverage = ?"))
		if err != nil {
			log.Fatal(err)
		}

		defer stmt.Close()

		rows, err = stmt.Query(set_name, data_coverage)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		stmt, err := db.Prepare(DB_change("select data, coverage from other where name = ?"))
		if err != nil {
			log.Fatal(err)
		}

		defer stmt.Close()

		rows, err = stmt.Query(set_name)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer rows.Close()

	data_list := [][]string{}

	for rows.Next() {
		var set_data string
		var set_coverage string

		err := rows.Scan(&set_data, &set_coverage)
		if err != nil {
			log.Fatal(err)
		}

		data_list = append(data_list, []string{set_data, set_coverage})
	}

	return data_list
}
