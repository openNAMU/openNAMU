package route

import (
	"database/sql"
	"regexp"
	"strings"

	"opennamu/route/tool"
)

var render_external_link_regex = regexp.MustCompile(`(?is)<a\b([^>]*\bclass="[^"]*\bopennamu_link_out\b[^"]*"[^>]*)>`)
var render_target_regex = regexp.MustCompile(`(?i)\btarget\s*=`)

func move_render_category_top(data string) string {
	category_start := strings.LastIndex(data, `<div class="opennamu_category" id="cate">`)
	if category_start < 0 {
		return data
	}

	category_end_offset := strings.Index(data[category_start:], `</div>`)
	if category_end_offset < 0 {
		return data
	}
	category_end := category_start + category_end_offset + len(`</div>`)
	separator := strings.LastIndex(data[:category_start], `<hr class="main_hr">`)
	opening_end := strings.Index(data, ">")
	if separator < 0 || opening_end < 0 || opening_end >= separator {
		return data
	}

	return data[:opening_end+1] + data[category_start:category_end] + `<hr class="main_hr">` + data[opening_end+1:separator] + data[category_end:]
}

func add_render_external_link_target(data string) string {
	return render_external_link_regex.ReplaceAllStringFunc(data, func(value string) string {
		match := render_external_link_regex.FindStringSubmatch(value)
		if len(match) < 2 || render_target_regex.MatchString(match[1]) {
			return value
		}
		return `<a` + match[1] + ` target="_blank">`
	})
}

func apply_render_setting_data(db *sql.DB, config tool.Config, data string) string {
	if tool.Get_main_skin_set(db, config, "main_css_category_set") != "bottom" {
		data = move_render_category_top(data)
	}
	return add_render_external_link_target(data)
}

func get_render_other_value(db *sql.DB, name string) string {
	data := ""
	tool.QueryRow_DB(
		db,
		"select data from other where name = ? and coverage = ''",
		[]any{&data},
		name,
	)
	return data
}

func get_render_setting_parameter(db *sql.DB, config tool.Config) map[string]string {
	setting := map[string]string{}
	for _, name := range []string{
		"main_css_bold",
		"main_css_strike",
		"main_css_include_link",
		"main_css_category_set",
		"main_css_toc_set",
		"main_css_footnote_set",
		"main_css_footnote_number",
		"main_css_view_real_footnote_num",
		"main_css_list_view_change",
	} {
		setting[name] = tool.Get_main_skin_set(db, config, name)
	}
	cookie_data := tool.Get_cookie_header(config.Cookies)
	darkmode := cookie_data["main_css_darkmode"]
	if _, exists := cookie_data["main_css_darkmode"]; !exists && config.Session != nil {
		darkmode, _ = config.Session.Get("main_css_darkmode").(string)
	}
	setting["main_css_darkmode"] = darkmode
	return setting
}

func get_render_setting_css(db *sql.DB, config tool.Config) string {
	data := strings.Builder{}

	font_size := tool.Get_main_skin_set(db, config, "main_css_font_size")
	switch font_size {
	case "10", "12", "14", "16", "18", "20", "22":
		data.WriteString(`.opennamu_render_complete { font-size: ` + font_size + `px !important; }`)
	}

	if get_render_other_value(db, "namumark_compatible") != "" {
		data.WriteString(`.opennamu_render_complete {
    font-size: 15px !important;
    line-height: 1.5;
}
.opennamu_render_complete td {
    padding: 5px 10px !important;
    word-break: break-all;
}
.opennamu_render_complete summary {
    list-style: none !important;
    font-weight: bold !important;
}
.opennamu_render_complete .opennamu_folding {
    margin-bottom: 5px;
}
.opennamu_render_complete .opennamu_footnote {
    padding-bottom: 30px;
}
.opennamu_render_complete iframe {
    display: block;
}`)
	}

	if tool.Get_main_skin_set(db, config, "main_css_table_scroll") == "on" {
		data.WriteString(`.table_safe { overflow-x: scroll; white-space: nowrap; }`)
	}
	if tool.Get_main_skin_set(db, config, "main_css_view_joke") == "off" {
		data.WriteString(`.opennamu_joke { display: none; }`)
	}
	if tool.Get_main_skin_set(db, config, "main_css_math_scroll") == "on" {
		data.WriteString(`.katex .base { overflow-x: scroll; }`)
	}
	if tool.Get_main_skin_set(db, config, "main_css_table_transparent") == "on" {
		data.WriteString(`.table_safe td {
    background: transparent !important;
    color: inherit !important;
}`)
	}
	if tool.Get_main_skin_set(db, config, "main_css_toc_set") == "off" {
		data.WriteString(`.opennamu_render_complete .opennamu_TOC { display: none !important; }`)
	}

	if data.Len() == 0 {
		return ""
	}
	return `<style>` + data.String() + `</style>`
}
