package route

import (
	"database/sql"
	"html"
	"opennamu/route/tool"
	"regexp"
	"strconv"
	"strings"
)

var re_esc_a = regexp.MustCompile(`&lt;a&gt;\s*([^\r\n]*?)\s*&lt;/a&gt;`)

func Get_safe_send_data(data string) string {    
    escaped_data := tool.HTML_escape(data)

    return re_esc_a.ReplaceAllStringFunc(escaped_data, func(match string) string {
        inner_esc := re_esc_a.FindStringSubmatch(match)[1]
        inner_text := strings.TrimSpace(html.UnescapeString(inner_esc))

        if inner_text == "" || strings.ContainsAny(inner_text, "<>") {
            return match
        }
        
        return `<a href="/w/` + tool.Url_parser(inner_text) + `">` + strings.TrimSpace(inner_esc) + `</a>`
    })
}

func Get_ui_history(db *sql.DB, config tool.Config, data_all [][]string) (string, string) {
    auth_name := tool.Get_user_auth(db, config.IP)
    auth_info := tool.Get_auth_group_info(db, auth_name)

    date_heading := ""
    data_html := ""
    data_select := ""

    for_count := 1
    for _, in_data := range data_all {
        for_count_str := strconv.Itoa(for_count)
        for_count += 1

        data_select = `<option value="` + for_count_str + `">` + for_count_str + `</option>` + data_select

        if in_data[6] != "" && in_data[1] == "" {
            if date_heading != "----" {
                data_html += "<h2>----</h2>"
                date_heading = "----"
            }

            data_html += tool.Get_list_ui("----", "", "", "")
            continue
        }

        doc_name := in_data[1]
        doc_name_url := tool.Url_parser(doc_name)
        rev_str := in_data[0]

        left := `<a href="/w/` + doc_name_url + `">` + tool.HTML_escape(doc_name) + `</a> `
        rev := ""

        if in_data[6] != "" {
            rev = `<span style="color: red;">r` + rev_str + `</span>`
        } else {
            rev = `r` + rev_str
        }

        rev_int := tool.Str_to_int(rev_str)
        if rev_int > 1 {
            before_rev := rev_int - 1
            before_rev_str := strconv.Itoa(before_rev)

            rev = `<a href="/diff/` + before_rev_str + `/` + rev_str + `/` + doc_name_url + `">` + rev + `</a>`
        }

        right := ""
        right += `<span id="opennamu_list_history_` + for_count_str + `_over">`
        right += `<a id="opennamu_list_history_` + for_count_str + `" href="javascript:void(0);">`
        right += `<span class="opennamu_svg opennamu_svg_tool">&nbsp;</span></a>`
        right += `<span class="opennamu_popup_footnote" id="opennamu_list_history_` + for_count_str + `_load" style="display: none;"></span>`
        right += `</span> | `
        right += rev + " | "

        diff_size := in_data[5]
        if diff_size == "0" {
            right += `<span style="color: gray;">` + diff_size + `</span>`
        } else if strings.Contains(diff_size, "+") {
            right += `<span style="color: green;">` + diff_size + `</span>`
        } else {
            right += `<span style="color: red;">` + diff_size + `</span>`
        }

        right += " | "
        right += in_data[7] + " | "

        edit_type := "edit"
        if in_data[8] != "" {
            edit_type = in_data[8]
        }

        right += tool.Get_language(db, edit_type, true) + " | "

        time_split := strings.Split(in_data[2], " ")
        if date_heading != time_split[0] {
            data_html += "<h2>" + time_split[0] + "</h2>"
            date_heading = time_split[0]
        }

        if len(time_split) > 1 {
            right += time_split[1]
        }

        right += `<span style="display: none;" id="opennamu_history_tool_` + for_count_str + `">`

        right += `<a href="/render/` + rev_str + `/` + doc_name_url + `">` + tool.Get_language(db, "view", true) + `</a>`
        right += ` | <a href="/raw_rev/` + rev_str + `/` + doc_name_url + `">` + tool.Get_language(db, "raw", true) + `</a>`
        right += ` | <a href="/revert/` + rev_str + `/` + doc_name_url + `">` + tool.Get_language(db, "revert", true) + ` (r` + rev_str + `)</a>`

        if rev_int > 1 {
            before_rev := rev_int - 1
            before_rev_str := strconv.Itoa(before_rev)

            right += ` | <a href="/revert/` + before_rev_str + `/` + doc_name_url + `">` + tool.Get_language(db, "revert", true) + ` (r` + before_rev_str + `)</a>`
            right += ` | <a href="/diff/` + before_rev_str + `/` + rev_str +  `/` + doc_name_url + `">` + tool.Get_language(db, "compare", true) + `</a>`
        }

        right += ` | <a href="/history/` + doc_name_url + `">` + tool.Get_language(db, "history", true) + `</a>`

        if _, ok := auth_info["owner"]; ok {
            right += ` | <a href="/history_hidden/` + rev_str + `/` + doc_name_url + `">` + tool.Get_language(db, "hide", true) + `</a>`
            right += ` | <a href="/history_delete/` + rev_str + `/` + doc_name_url + `">` + tool.Get_language(db, "history_delete", true) + `</a>`
            right += ` | <a href="/history_send/` + rev_str + `/` + doc_name_url + `">` + tool.Get_language(db, "send_edit", true) + `</a>`
        } else if _, ok := auth_info["hidel"]; ok {
            right += ` | <a href="/history_hidden/` + rev_str + `/` + doc_name_url + `">` + tool.Get_language(db, "hide", true) + `</a>`
        }

        right += `</span>`

        bottom := ``
        if in_data[4] != "" {
            bottom = Get_safe_send_data(in_data[4])
        }

        data_html += tool.Get_list_ui(left, right, bottom, "")

        data_html += `<script>
            document.getElementById('opennamu_list_history_` + for_count_str + `').addEventListener("click", function() {{
                opennamu_do_footnote_popover('opennamu_list_history_` + for_count_str + `', '', 'opennamu_history_tool_` + for_count_str + `', 'open');
            }});
            document.addEventListener("click", function() {{
                opennamu_do_footnote_popover('opennamu_list_history_` + for_count_str + `', '', 'opennamu_history_tool_` + for_count_str + `', 'close');
            }});
        </script>`
    }

    return data_html, data_select
}

func View_list_recent_change(config tool.Config, set_type string, limit string, num string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    sub := ""
    if set_type == "" {
        set_type = "normal"
    } else {
        sub = "(" + tool.Get_language(db, set_type, true) + ")"
    }

    data_html := ""

    menu_option := []string{ "normal", "edit", "move", "delete", "revert", "r1", "edit_request", "user", "file", "category" }
    for _, option := range menu_option {
        label := tool.Get_language(db, option, true)
        data_html += `<a href="/recent_change/1/` + option + `">(` + label + `)</a> `
    }

    api_data := Api_list_recent_change(config, set_type, limit, num)
    api_data_list := api_data["data"].([][]string)

    history_ui, _ := Get_ui_history(db, config, api_data_list)

    data_html += history_ui
    data_html += tool.Get_page_control(
        db,
        tool.Str_to_int(num),
        len(api_data_list),
        tool.Str_to_int(limit),
        "/recent_change/{}/" + set_type,
    )

    out := tool.Get_template(
        db,
        config,
        tool.Get_language(db, "recent_change", true),
        data_html,
        []any{ sub },
        [][]any{},
        map[string]string{},
    )

    return out
}
