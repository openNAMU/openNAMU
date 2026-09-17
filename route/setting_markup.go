package route

import (
	"database/sql"
	"strconv"
	"time"

	"opennamu/route/tool"
	"opennamu/route/tool/markup"
)

func Setting_markup_normalize(markup_name string) string {
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

func Setting_markup_value(db *sql.DB, setting_name string) string {
	return Setting_markup_normalize(
		tool.Get_setting_value(db, setting_name+"_markup", "", ""),
	)
}

func Setting_markup_options() []string {
	return markup.List_markup()
}

func Setting_markup_select_ui(name string, current string, disabled string) string {
	current = Setting_markup_normalize(current)
	data := `<span class="__ON_SELECT_DIV__"><select class="__ON_SELECT__" name="` + tool.HTML_escape(name) + `"` + disabled + `>`
	for _, markup_option := range Setting_markup_options() {
		selected := ""
		if markup_option == current {
			selected = ` selected`
		}
		data += `<option value="` + tool.HTML_escape(markup_option) + `"` + selected + `>` + tool.HTML_escape(markup_option) + `</option>`
	}
	data += `</select></span>`
	return data
}

func Setting_render_markup(db *sql.DB, data string, markup_name string) string {
	markup_name = Setting_markup_normalize(markup_name)
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

func Document_editor_top_render(db *sql.DB, doc_name string) string {
	top_data := tool.Get_document_setting_value_exact(db, doc_name, "document_editor_top", "")
	top_markup := ""
	if top_data != "" {
		top_markup = Setting_markup_normalize(tool.Get_document_setting_value_exact(db, doc_name, "document_editor_top_markup", ""))
	} else {
		top_data = tool.Get_document_setting_value_exact(db, doc_name, "document_top", "")
		if top_data == "" {
			return ""
		}
		top_markup = Setting_markup_normalize(tool.Get_document_setting_value_exact(db, doc_name, "document_top_markup", ""))
	}

	return Setting_render_markup(db, top_data, top_markup)
}
