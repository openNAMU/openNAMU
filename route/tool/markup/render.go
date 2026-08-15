package markup

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"opennamu/route/tool"
	backlink_parser "opennamu/route/tool/markup/backlink"
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

func Get_render(db *sql.DB, doc_name string, data string, render_type string, parameters ...map[string]any) map[string]string {
	parameter_data := map[string]any{}
	if len(parameters) > 0 && parameters[0] != nil {
		parameter_data = parameters[0]
	}
	markup := tool.Get_document_markup(db, doc_name, "")
	if render_type == "normal" || render_type == "api_view" || render_type == "api_from" || render_type == "api_include" || render_type == "backlink" {
		markup = tool.Get_document_markup(db, doc_name, "document")
	}
	now_time := time.Now().UnixNano()
	render_name := strconv.Itoa(int(now_time))

	render_data := Get_render_direct(db, doc_name, data, markup, render_name, render_type, parameter_data)

	return render_data
}

func Get_render_direct(db *sql.DB, doc_name string, data string, markup string, render_name string, render_type string, parameters ...map[string]any) map[string]string {
	parameter_data := map[string]any{}
	if len(parameters) > 0 && parameters[0] != nil {
		parameter_data = parameters[0]
	}

	from := ""
	include := ""
	backlink_mode := false

	switch render_type {
	case "api_include":
		include = "1"
	case "api_from":
		from = "1"
	case "backlink":
		backlink_mode = true
	}

	if render_type == "api_view" || render_type == "api_from" || render_type == "api_include" || render_type == "backlink" {
		render_type = "view"
	} else if render_type == "api_thread" {
		render_type = "thread"
	}

	doc_data_set := map[string]string{
		"doc_name":    doc_name,
		"data":        data,
		"render_name": render_name,
		"render_type": render_type,
		"from":        from,
		"include":     include,
	}

	if !backlink_mode && (markup == "namumark" || markup == "namumark_beta") {
		compat_render_type := render_type
		if include == "1" {
			compat_render_type = "include"
		}
		compat_data := render_namumark_compat(db, doc_name, data, compat_render_type, parameter_data, 0, false)
		return map[string]string{
			"data":    `<div class="opennamu_render_complete">` + compat_data["data"].(string) + `</div>`,
			"js_data": compat_data["js_data"].(string),
		}
	}

	if backlink_mode && (markup == "namumark" || markup == "namumark_beta") {
		return render_namumark_compat_backlink(db, doc_name, data)
	}

	render_data := make(map[string]any)
	backlink_list := map[string][]string{}
	backlink_count := 0
	backlink_supported := false

	if backlink_mode {
		backlink_list, backlink_count, backlink_supported = backlink_parser.Get_backlink(data, markup)
		render_data["data"] = ""
		render_data["js_data"] = ""
	} else {
		switch markup {
		case "namumark", "namumark_beta":
			render_data_class := Namumark_new(db, doc_data_set, parameter_data)
			render_data = render_data_class.main()
		case "markdown":
			render_data = Markdown(db, doc_data_set)
		case "raw":
			render_data["data"] = strings.ReplaceAll(tool.HTML_escape(data), "\n", "<br>")
			render_data["js_data"] = ""
		case "macromark":
			render_data_class := Macromark_new(db, doc_data_set, "html")
			render_data = render_data_class.main()
		default:
			render_data["data"] = data
			render_data["js_data"] = ""
		}
	}

	if backlink_mode && backlink_supported {
		tool.Exec_DB(
			db,
			"delete from back where link = ?",
			doc_name,
		)
		tool.Exec_DB(db, "delete from back where title = ? and type = 'no'", doc_name)

		tool.Exec_DB(
			db,
			"delete from data_set where doc_name = ? and set_name = 'link_count'",
			doc_name,
		)

		for link, link_type_list := range backlink_list {
			for _, link_type := range link_type_list {
				tool.Exec_DB(
					db,
					"insert into back (link, title, type, data) values (?, ?, ?, '')",
					doc_name,
					link,
					link_type,
				)
			}
		}

		tool.Exec_DB(
			db,
			"insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', 'link_count', ?)",
			doc_name, backlink_count,
		)
	}

	return map[string]string{
		"data":    "<div id=\"opennamu_render_complete\" class=\"opennamu_render_complete\">" + render_data["data"].(string) + "</div>",
		"js_data": render_data["js_data"].(string),
	}
}
