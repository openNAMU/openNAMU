package markup

import (
	"opennamu/route/tool"

	"database/sql"
	"strconv"
	"time"
)

func List_markup() []string {
    return []string{
        "namumark",
        "namumark_beta",
        "macromark",
        "markdown",
        "custom",
        "raw",
    }
}

func Get_render(db *sql.DB, doc_name string, data string, render_type string) map[string]string {
    markup := ""
    if render_type == "api_view" || render_type == "api_from" || render_type == "api_include" || render_type == "backlink" {
        markup = tool.Get_document_markup(db, doc_name, "document")
    } else {
        markup = tool.Get_document_markup(db, doc_name, "")
    }

    now_time := time.Now().UnixNano()
    render_name := strconv.Itoa(int(now_time))

    render_data := Get_render_direct(db, doc_name, data, markup, render_name, render_type)

    return render_data
}

func Get_render_direct(db *sql.DB, doc_name string, data string, markup string, render_name string, render_type string) map[string]string {
    from := ""
    include := ""
    backlink := ""
    
    switch render_type {
    case "api_include":
        include = "1"
    case "api_from":
        from = "1"
    case "backlink":
        backlink = "1"
    }

    if render_type == "api_view" || render_type == "api_from" || render_type == "api_include" || render_type == "backlink" {
        render_type = "view"
    }

    doc_data_set := map[string]string{
        "doc_name" : doc_name,
        "data" : data,
        "render_name" : render_name,
        "render_type" : render_type,
        "from" : from,
        "include" : include,
    }

    render_data := make(map[string]any)
    switch markup {
    case "namumark":
        render_data_class := Namumark_new(db, doc_data_set)
        render_data = render_data_class.main()
    case "markdown":
        render_data = Markdown(db, doc_data_set)
    case "macromark":
        render_data_class := Macromark_new(db, doc_data_set, "html")
        render_data = render_data_class.main()
    default:
        render_data["data"] = data
        render_data["js_data"] = ""
        render_data["backlink"] = [][]string{}
    }

    if backlink == "1" {
        tool.Exec_DB(
            db,
            "delete from back where link = ?",
            doc_name,
        )

        tool.Exec_DB(
            db,
            "delete from back where title = ? and type = 'no'",
            doc_name,
        )
        
        tool.Exec_DB(
            db,
            "delete from data_set where doc_name = ? and set_name = 'link_count'",
            doc_name,
        )

        tool.Exec_DB(
            db,
            "delete from data_set where doc_name = ? and set_name = 'doc_type'",
            doc_name,
        )

        end_backlink := render_data["backlink"].([][]string)
        for for_a := 0; for_a < len(end_backlink); for_a++ {
            tool.Exec_DB(
                db,
                "insert into back (link, title, type, data) values (?, ?, ?, ?)",
                end_backlink[0], end_backlink[1], end_backlink[2],
            )
        }

        tool.Exec_DB(
            db,
            "insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', 'link_count', ?)",
            doc_name, render_data["link_count"].(int),
        )

        tool.Exec_DB(
            db,
            "insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', 'doc_type', ?)",
            doc_name, "",
        )
    }

    return map[string]string{
        "data" : "<div id=\"opennamu_render_complete\">" + render_data["data"].(string) + "</div>",
        "js_data" : render_data["js_data"].(string),
    }
}
