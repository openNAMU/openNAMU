package tool

import (
	"database/sql"

	"github.com/russross/blackfriday/v2"
)

func Markdown(db *sql.DB, db_set map[string]string, data map[string]string) map[string]interface{} {
	input := []byte(data["data"])

	html := blackfriday.Run(input)

	end_data := make(map[string]interface{})
	end_data["data"] = string(html)
	end_data["js_data"] = ""
	end_data["backlink"] = []string{}

	return end_data
}
