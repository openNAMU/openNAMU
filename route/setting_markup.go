package route

import (
	"database/sql"
	"strconv"
	"time"

	"opennamu/route/tool"
	"opennamu/route/tool/markup"
)

func setting_markup_normalize(markup_name string) string {
	switch markup_name {
	case "", "custom":
		return "html"
	case "raw":
		return "plain"
	}

	for _, available_markup := range markup.List_markup() {
		if markup_name == available_markup {
			return markup_name
		}
	}

	return "html"
}

func setting_markup_value(db *sql.DB, setting_name string) string {
	return setting_markup_normalize(
		tool.Get_setting_value(db, setting_name+"_markup", "", ""),
	)
}

func setting_markup_options() []string {
	return markup.List_markup()
}

func setting_render_markup(db *sql.DB, data string, markup_name string) string {
	markup_name = setting_markup_normalize(markup_name)
	if markup_name == "html" {
		return data
	}

	render_data := markup.Get_render_direct(
		db,
		"",
		data,
		markup_name,
		strconv.FormatInt(time.Now().UnixNano(), 10),
		"setting",
	)

	return render_data["data"]
}

func document_editor_top_render(db *sql.DB, doc_name string) string {
	top_data := tool.Get_document_setting_value_exact(db, doc_name, "document_editor_top", "")
	top_markup := ""
	if top_data != "" {
		top_markup = setting_markup_normalize(tool.Get_document_setting_value_exact(db, doc_name, "document_editor_top_markup", ""))
	} else {
		top_data = tool.Get_document_setting_value_exact(db, doc_name, "document_top", "")
		if top_data == "" {
			return ""
		}
		top_markup = setting_markup_normalize(tool.Get_document_setting_value_exact(db, doc_name, "document_top_markup", ""))
	}

	return setting_render_markup(db, top_data, top_markup)
}
