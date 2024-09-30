package tool

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"html/template"
	"log"
	"net/url"
	"time"
	"html"
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

func HTML_unescape(data string) string {
	return html.UnescapeString(data)
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

func Get_month() string {
	return time.Now().Format("2006-01")
}

func Get_document_setting(db *sql.DB, doc_name string, set_name string, doc_rev string) [][]string {
	var rows *sql.Rows

	if doc_rev != "" {
		stmt, err := db.Prepare(DB_change("select set_data, doc_rev from data_set where doc_name = ? and doc_rev = ? and set_name = ?"))
		if err != nil {
			log.Fatal(err)
		}

		defer stmt.Close()

		rows, err = stmt.Query(doc_name, doc_rev, set_name)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		stmt, err := db.Prepare(DB_change("select set_data, doc_rev from data_set where doc_name = ? and set_name = ?"))
		if err != nil {
			log.Fatal(err)
		}

		defer stmt.Close()

		rows, err = stmt.Query(doc_name, set_name)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer rows.Close()

	data_list := [][]string{}

	for rows.Next() {
		var set_data string
		var doc_rev string

		err := rows.Scan(&set_data, &doc_rev)
		if err != nil {
			log.Fatal(err)
		}

		data_list = append(data_list, []string{set_data, doc_rev})
	}

	return data_list
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
