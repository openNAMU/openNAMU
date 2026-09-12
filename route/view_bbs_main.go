package route

import (
	"database/sql"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"opennamu/route/tool"
)

func search_highlight(data string, keyword string) string {
	data = tool.HTML_escape(data)
	keyword_list := []string{}
	for _, keyword_data := range strings.Fields(keyword) {
		keyword_data = tool.HTML_escape(keyword_data)
		if keyword_data != "" {
			keyword_list = append(keyword_list, regexp.QuoteMeta(keyword_data))
		}
	}
	if len(keyword_list) == 0 {
		return data
	}
	sort.Slice(keyword_list, func(i int, j int) bool {
		return len(keyword_list[i]) > len(keyword_list[j])
	})
	pattern, err := regexp.Compile("(?i)" + strings.Join(keyword_list, "|"))
	if err != nil {
		return data
	}
	return pattern.ReplaceAllStringFunc(data, func(value string) string {
		return "<mark>" + value + "</mark>"
	})
}

func search_snippet(data string, keyword string) string {
	data = strings.ReplaceAll(strings.TrimSpace(data), "\r", " ")
	data = strings.ReplaceAll(data, "\n", " ")
	if data == "" {
		return ""
	}

	start := 0
	end := tool.Get_len(data)
	if end > 160 {
		end = 160
	}
	for _, keyword_data := range strings.Fields(keyword) {
		pattern, err := regexp.Compile("(?i)" + regexp.QuoteMeta(keyword_data))
		if err != nil {
			continue
		}
		match := pattern.FindStringIndex(data)
		if match == nil {
			continue
		}
		match_start := tool.Get_len(data[:match[0]])
		start = match_start - 60
		if start < 0 {
			start = 0
		}
		end = start + 160
		if end > tool.Get_len(data) {
			end = tool.Get_len(data)
		}
		break
	}

	result := tool.Get_slice(data, start, end)
	if start > 0 {
		result = "..." + result
	}
	if end < tool.Get_len(data) {
		result += "..."
	}
	return search_highlight(result, keyword)
}

func bbs_tags_html(bbs_id string, tags string, keyword string) string {
	data_html := ""
	for _, tag := range bbs_tag_list(tags) {
		if data_html != "" {
			data_html += ", "
		}
		data_html += `<a href="/bbs/in/` + tool.Url_parser(bbs_id) + `/filter/tag/` + tool.Url_parser(tag) + `/1">#` + search_highlight(tag, keyword) + `</a>`
	}
	return data_html
}

func bbs_list_example_ui(db *sql.DB) string {
	left := tool.Get_language(db, "title", true) + " [" + tool.Get_language(db, "statistics_bbs_comment_count", true) + "] [+" + tool.Get_language(db, "upvote", true) + "] [-" + tool.Get_language(db, "downvote", true) + "]"
	right := tool.Get_language(db, "page_view", true) + " | " + tool.Get_language(db, "user_name", true) + " | " + tool.Get_language(db, "date", true)
	return tool.Get_list_ui(left, right, "", "")
}

func Get_bbs_list_ui(db *sql.DB, config tool.Config, bbs_all_data []map[string]string, bbs_id_to_name map[string]string) string {
	count := 0
	data_html := ""
	date_heading := ""

	for _, in_data := range bbs_all_data {
		count_str := strconv.Itoa(count)
		count += 1

		bbs_title := in_data["title"]
		bbs_title_html := tool.HTML_escape(bbs_title)
		if in_data["title_html"] != "" {
			bbs_title_html = in_data["title_html"]
		}
		bbs_id := in_data["set_id"]
		bbs_code := in_data["set_code"]
		bbs_date := in_data["date"]
		if in_data["activity_date"] != "" {
			bbs_date = in_data["activity_date"]
		}
		bbs_user_id := in_data["user_id_render"]

		bbs_comment_length := tool.Str_to_int(in_data["comment_count"])

		bbs_comment_length_str := strconv.Itoa(bbs_comment_length)
		bbs_tabom_length_str := strconv.Itoa(tool.Str_to_int(in_data["tabom_count"]))
		bbs_tabom_down_length_str := strconv.Itoa(tool.Str_to_int(in_data["tabom_down_count"]))

		bbs_view_count := "0"
		if _, ok := in_data["view_count"]; ok {
			bbs_view_count = in_data["view_count"]
		}

		bbs_name := ""
		if len(bbs_id_to_name) != 0 {
			bbs_name = bbs_id_to_name[bbs_id]
		}

		bbs_link := "/bbs/w/" + bbs_id + "/" + bbs_code
		if in_data["comment_code"] != "" {
			bbs_link += "#" + tool.Url_parser(in_data["comment_code"])
		}

		left := ""
		if in_data["prefix"] != "" {
			prefix_html := tool.HTML_escape(in_data["prefix"])
			if in_data["prefix_html"] != "" {
				prefix_html = in_data["prefix_html"]
			}
			left += "[" + prefix_html + "] "
		}
		if in_data["blind"] == "O" {
			left = "[" + tool.Get_language(db, "blind_post", true) + "] " + left
		}
		left += `<a href="` + bbs_link + `">` + bbs_title_html + `</a>`

		if bbs_name != "" {
			left += ` <a href="/bbs/in/` + bbs_id + `">(` + bbs_name + `)</a>`
		}

		left += ` [` + bbs_comment_length_str + `] [+` + bbs_tabom_length_str + `] [-` + bbs_tabom_down_length_str + `]`

		right := ""
		right += `<span id="opennamu_bbs_comment_` + count_str + `"></span>`
		right += bbs_view_count + " | "
		right += bbs_user_id + " | "
		date_ui, date_text, new_date_heading := tool.Get_date_list_ui(bbs_date, date_heading)
		date_heading = new_date_heading
		right += date_text

		class_name := ""
		if in_data["pinned"] == "1" {
			class_name = "opennamu_comment_color_red"
			left = "<strong>" + left + "</strong>"
		}

		bottom := in_data["search_snippet_html"]
		if bottom != "" {
			bottom += "<br>"
		}
		tags_html := in_data["tags_html"]
		if tags_html == "" {
			tags_html = bbs_tags_html(bbs_id, in_data["tags"], "")
		}
		bottom += tags_html
		data_html += date_ui
		data_html += tool.Get_list_ui(left, right, bottom, class_name)
	}

	return data_html
}

func View_bbs_main(config tool.Config, page string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "bbs_main_view", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	bbs_list_api_data := Api_bbs_list(config)

	bbs_id_to_name := map[string]string{}

	data_html := "<ul>"
	for _, in_data := range bbs_list_api_data["data"].([][]string) {
		bbs_name := in_data[0]
		bbs_id := in_data[1]
		bbs_type := in_data[2]
		bbs_date := in_data[3]

		bbs_id_to_name[bbs_id] = bbs_name

		data_html += "<li>"
		data_html += "<a href=\"/bbs/in/" + tool.Url_parser(bbs_id) + "\">"
		data_html += tool.HTML_escape(bbs_name)
		data_html += "</a>"

		if bbs_type == "comment" {
			data_html += " (" + tool.Get_language(db, "comment_base", false) + ")"
		} else {
			data_html += " (" + tool.Get_language(db, "thread_base", false) + ")"
		}

		if bbs_date != "" {
			data_html += " (" + bbs_date + ")"
		}

		data_html += "</li>"
	}

	data_html += "</ul><hr class=\"main_hr\">"
	data_html += bbs_list_example_ui(db)

	bbs_api_data := Api_bbs(config, "", page, "")
	data_html += Get_bbs_list_ui(db, config, bbs_api_data["data"].([]map[string]string), bbs_id_to_name)

	menu := [][]any{
		{"other", tool.Get_language(db, "other_tool", false)},
		{"bbs/search", tool.Get_language(db, "search", false)},
		{"bbs/contributor", tool.Get_language(db, "monthly_bbs_contributor", false)},
	}
	if tool.Check_permission(db, "bbs_create", config.IP) {
		menu = append(menu, []any{"bbs/make", tool.Get_language(db, "add", false)})
	}

	out := tool.Get_template(
		db,
		config,
		tool.Get_language(db, "bbs_main", true),
		data_html,
		[]any{},
		menu,
		map[string]string{},
	)

	return out
}
