package tool

import (
	"database/sql"
	"log"
	"strconv"
	"time"
)

func Get_render(db *sql.DB, db_set map[string]string, doc_name string, data string, render_type string) map[string]string {
	var markup string

	if render_type == "api_view" || render_type == "api_from" || render_type == "api_include" || render_type == "backlink" {
		stmt, err := db.Prepare(DB_change(db_set, "select set_data from data_set where doc_name = ? and set_name = 'document_markup'"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		err = stmt.QueryRow(doc_name).Scan(&markup)
		if err != nil {
			if err == sql.ErrNoRows {
				markup = ""
			} else {
				log.Fatal(err)
			}
		}
	}

	if markup == "" {
		err := db.QueryRow(DB_change(db_set, "select data from other where name = 'markup'")).Scan(&markup)
		if err != nil {
			if err == sql.ErrNoRows {
				markup = ""
			} else {
				log.Fatal(err)
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
	from := ""
	include := ""
	backlink := ""
	if render_type == "api_include" {
		include = "1"
	} else if render_type == "api_from" {
		from = "1"
	} else if render_type == "backlink" {
		backlink = "1"
	}

	if render_type == "api_view" || render_type == "api_from" || render_type == "api_include" || render_type == "backlink" {
		render_type = "view"
	}

	doc_data_set := map[string]string{
		"doc_name":    doc_name,
		"data":        data,
		"render_name": render_name,
		"render_type": render_type,
		"from":        from,
		"include":     include,
	}

	render_data := make(map[string]interface{})
	if markup == "namumark" {
		render_data = Namumark()
	} else if markup == "markdown" {
		render_data = Markdown(db, db_set, doc_data_set)
	} else {
		render_data["data"] = data
		render_data["js_data"] = ""
		render_data["backlink"] = [][]string{}
	}

	if backlink == "1" {
		stmt, err := db.Prepare(DB_change(db_set, "delete from back where link = ?"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(doc_name)
		if err != nil {
			log.Fatal(err)
		}

		stmt, err = db.Prepare(DB_change(db_set, "delete from back where title = ? and type = 'no'"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(doc_name)
		if err != nil {
			log.Fatal(err)
		}

		stmt, err = db.Prepare(DB_change(db_set, "delete from data_set where doc_name = ? and set_name = 'link_count'"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(doc_name)
		if err != nil {
			log.Fatal(err)
		}

		stmt, err = db.Prepare(DB_change(db_set, "delete from data_set where doc_name = ? and set_name = 'doc_type'"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(doc_name)
		if err != nil {
			log.Fatal(err)
		}

		end_backlink := render_data["backlink"].([][]string)
		for for_a := 0; for_a < len(end_backlink); for_a++ {
			stmt, err := db.Prepare(DB_change(db_set, "insert into back (link, title, type, data) values (?, ?, ?, ?)"))
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()

			_, err = stmt.Exec(end_backlink[0], end_backlink[1], end_backlink[2])
			if err != nil {
				log.Fatal(err)
			}
		}

		stmt, err = db.Prepare(DB_change(db_set, "insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', 'link_count', ?)"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(doc_name, render_data["link_count"].(int))
		if err != nil {
			log.Fatal(err)
		}

		stmt, err = db.Prepare(DB_change(db_set, "insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', 'doc_type', ?)"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(doc_name, "")
		if err != nil {
			log.Fatal(err)
		}
	}

	return map[string]string{
		"data":    "<div id=\"opennamu_render_complete\">" + render_data["data"].(string) + "</div>",
		"js_data": render_data["js_data"].(string),
	}
}
