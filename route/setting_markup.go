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
		"view",
	)

	return render_data["data"]
}
