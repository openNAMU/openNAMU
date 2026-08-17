package route

import (
	"database/sql"
	"strings"

	"opennamu/route/tool"
)

func View_setting_skin_set(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	fields, _ := setting_skin_fields(db)
	values := map[string]string{}
	for _, field := range fields {
		values[field.name] = setting_value(db, field.name, "", "")
	}

	return view_setting_skin_set_data(db, config, fields, values)
}

type setting_skin_field struct {
	name  string
	label string
}

func setting_skin_fields(db *sql.DB) ([]setting_skin_field, map[string][][]string) {
	set_list := Get_main_skin_set_list(db)
	fields := []setting_skin_field{
		{name: "main_css_strike", label: "strike"},
		{name: "main_css_bold", label: "bold"},
		{name: "main_css_category_set", label: "category"},
		{name: "main_css_category_change_title", label: "category_change_title"},
		{name: "main_css_footnote_set", label: "footnote_render"},
		{name: "main_css_footnote_number", label: "footnote_number"},
		{name: "main_css_view_real_footnote_num", label: "footnote_real_num_view"},
		{name: "main_css_include_link", label: "include_link"},
		{name: "main_css_image_set", label: "image"},
		{name: "main_css_toc_set", label: "toc"},
		{name: "main_css_exter_link", label: "exter_link"},
		{name: "main_css_link_delimiter", label: "link_delimiter"},
		{name: "main_css_darkmode", label: "force_darkmode"},
		{name: "main_css_table_scroll", label: "table_scroll"},
		{name: "main_css_table_transparent", label: "table_transparent"},
		{name: "main_css_list_view_change", label: "list_view_change"},
		{name: "main_css_view_joke", label: "view_joke"},
		{name: "main_css_math_scroll", label: "math_scroll"},
		{name: "main_css_view_history", label: "view_history"},
		{name: "main_css_font_size", label: "font_size"},
		{name: "main_css_monaco", label: "monaco_editor"},
	}

	return fields, set_list
}

func view_setting_skin_set_data(db *sql.DB, config tool.Config, fields []setting_skin_field, values map[string]string) string {
	_, set_list := setting_skin_fields(db)
	data := strings.Builder{}
	data.WriteString(`<form method="post"><h2>` + tool.Get_language(db, "render", true) + `</h2>`)

	for _, field := range fields {
		choices := set_list[field.name]
		choice_values := make([]string, 0, len(choices))
		choice_labels := make(map[string]string, len(choices))
		for _, choice := range choices {
			if len(choice) < 2 {
				continue
			}
			choice_values = append(choice_values, choice[0])
			choice_labels[choice[0]] = choice[1]
		}

		data.WriteString(`<h3>` + tool.Get_language(db, field.label, true) + `</h3>`)
		data.WriteString(`<select name="` + tool.HTML_escape(field.name) + `">` + setting_options(values[field.name], choice_values, choice_labels) + `</select>` + setting_hr())
	}

	data.WriteString(`<button type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`)
	return setting_page(db, config, tool.Get_language(db, "main_skin_set_default", true), data.String(), "setting")
}
