package route

import (
	"database/sql"
	"strings"

	"opennamu/route/tool"
)

type setting_field struct {
	name          string
	default_value string
}

func setting_value(db *sql.DB, name string, coverage string, default_value string) string {
	data := ""
	exists := tool.QueryRow_DB(
		db,
		"select data from other where name = ? and coverage = ?",
		[]any{&data},
		name,
		coverage,
	)

	if !exists {
		setting_save_value(db, name, coverage, default_value)
		return default_value
	}

	return data
}

func setting_save_value(db *sql.DB, name string, coverage string, data string) {
	old_data := ""
	exists := tool.QueryRow_DB(
		db,
		"select data from other where name = ? and coverage = ?",
		[]any{&old_data},
		name,
		coverage,
	)

	if exists {
		tool.Exec_DB(
			db,
			"update other set data = ? where name = ? and coverage = ?",
			data,
			name,
			coverage,
		)
		return
	}

	tool.Exec_DB(
		db,
		"insert into other (name, data, coverage) values (?, ?, ?)",
		name,
		data,
		coverage,
	)
}

func setting_load_fields(db *sql.DB, fields []setting_field) map[string]string {
	data := make(map[string]string, len(fields))

	for _, field := range fields {
		data[field.name] = setting_value(db, field.name, "", field.default_value)
	}

	return data
}

func setting_save_fields(db *sql.DB, fields []setting_field, form map[string]string) {
	for _, field := range fields {
		value, exists := form[field.name]
		if !exists {
			value = field.default_value
		}

		setting_save_value(db, field.name, "", value)
	}
}

func setting_form_value(form map[string]string, name string, default_value string) string {
	value, exists := form[name]
	if !exists {
		return default_value
	}

	return value
}

func setting_checked(value string) string {
	if value != "" {
		return `checked="checked"`
	}

	return ""
}

func setting_options(current string, values []string, labels map[string]string) string {
	data := strings.Builder{}

	for _, value := range values {
		selected := ""
		if value == current {
			selected = ` selected="selected"`
		}

		label := value
		if labels != nil {
			if label_data, ok := labels[value]; ok {
				label = label_data
			}
		}

		data.WriteString(`<option value="`)
		data.WriteString(tool.HTML_escape(value))
		data.WriteString(`"`)
		data.WriteString(selected)
		data.WriteString(`>`)
		data.WriteString(tool.HTML_escape(label))
		data.WriteString(`</option>`)
	}

	return data.String()
}

func setting_input(name string, value string, input_type string) string {
	if input_type == "" {
		input_type = "text"
	}

	return `<input type="` + tool.HTML_escape(input_type) + `" name="` + tool.HTML_escape(name) + `" value="` + tool.HTML_escape(value) + `">`
}

func setting_textarea(name string, value string, class_name string) string {
	if class_name == "" {
		class_name = "opennamu_textarea_100"
	}

	return `<textarea class="` + tool.HTML_escape(class_name) + `" name="` + tool.HTML_escape(name) + `">` + tool.HTML_escape(value) + `</textarea>`
}

func setting_hr() string {
	return `<hr class="main_hr">`
}

func setting_page(db *sql.DB, config tool.Config, title string, data string, return_path string) string {
	menu := [][]any{}
	if return_path != "" {
		menu = append(menu, []any{return_path, tool.Get_language(db, "return", true)})
	}

	return tool.Get_template(
		db,
		config,
		title,
		data,
		[]any{},
		menu,
		map[string]string{},
	)
}
