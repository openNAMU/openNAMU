package tool

import (
	"database/sql"
	"strconv"
	"time"
)

func Get_render(db *sql.DB, db_set map[string]string, doc_name string, data string, render_type string) map[string]string {
	var markup string

	if render_type == "api_view" || render_type == "api_from" || render_type == "api_include" {
		stmt, err := db.Prepare(DB_change(db_set, "select set_data from data_set where doc_name = ? and set_name = 'document_markup'"))
		if err != nil {
			return map[string]string{}
		}
		defer stmt.Close()

		err = stmt.QueryRow(doc_name).Scan(&markup)
		if err != nil {
			if err == sql.ErrNoRows {
				markup = ""
			} else {
				return map[string]string{}
			}
		}
	}

	if markup == "" {
		err := db.QueryRow(DB_change(db_set, "select data from other where name = 'markup'")).Scan(&markup)
		if err != nil {
			if err == sql.ErrNoRows {
				markup = ""
			} else {
				return map[string]string{}
			}
		}
	}

	if markup == "" {
		markup = "namumark"
	}

	now_time := time.Now().UnixNano()
	render_name := strconv.Itoa(int(now_time))

	render_data := Get_render_direct(db, db_set, doc_name, data, markup, render_name, render_type)

	return render_data
}

func Get_render_direct(db *sql.DB, db_set map[string]string, doc_name string, data string, markup string, render_name string, render_type string) map[string]string {
	doc_data_set := map[string]string{
		"doc_name":    doc_name,
		"data":        data,
		"render_name": render_name,
		"render_type": render_type,
	}

	render_data := make(map[string]interface{})
	if markup == "namumark" {
		render_data = Namumark()
	} else if markup == "markdown" {
		render_data = Markdown(db, db_set, doc_data_set)
	} else {
		render_data["data"] = data
		render_data["js_data"] = ""
		render_data["backlink"] = []string{}
	}

	return map[string]string{
		"data":    render_data["data"].(string),
		"data_js": render_data["js_data"].(string),
	}
}
