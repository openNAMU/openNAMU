package route

import (
	"opennamu/route/tool"
	"strconv"
)

func Get_bbs_list_ui(config tool.Config, bbs_all_data []map[string]string, bbs_id_to_name map[string]string) string {
    count := 0
    data_html := ""

    for _, in_data := range bbs_all_data {
        count_str := strconv.Itoa(count)
        count += 1

        bbs_title := in_data["title"]
        bbs_id := in_data["set_id"]
        bbs_code := in_data["set_code"]
        bbs_date := in_data["date"]
        bbs_user_id := in_data["user_id_render"]

        bbs_comment_length_api := Api_bbs_w_comment(config, "length", bbs_id + "-" + bbs_code)
        bbs_comment_length := bbs_comment_length_api["data"].(int)

        bbs_comment_length_str := strconv.Itoa(bbs_comment_length)
        
        bbs_view_count := "0"
        if _, ok := in_data["view_count"]; ok {
            bbs_view_count = in_data["view_count"]
        }

        bbs_name := ""
        if len(bbs_id_to_name) != 0 {
            bbs_name = bbs_id_to_name[bbs_id]
        }

        left := ""
        left += `<a href="/bbs/w/` + bbs_id + `/` + bbs_code + `">` + tool.HTML_escape(bbs_title) + `</a>`

        if bbs_name != "" {
            left += ` <a href="/bbs/in/` + bbs_id + `">(` + bbs_name + `)</a>`
        }

        left += ` [` + bbs_comment_length_str + `]`

        right := ""
        right += `<span id="opennamu_bbs_comment_` + count_str + `"></span>`
        right += bbs_view_count + " | "
        right += bbs_user_id + " | "
        right += bbs_date

        data_html += tool.Get_list_ui(left, right, "", "")
    }

    return data_html
}

func View_bbs_main(config tool.Config, page string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

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

    bbs_api_data := Api_bbs(config, "", page)
    data_html += Get_bbs_list_ui(config, bbs_api_data["data"].([]map[string]string), bbs_id_to_name)

    out := tool.Get_template(
        db,
        config,
        tool.Get_language(db, "bbs_main", true),
        data_html,
        []any{},
        [][]any{
            { "other", tool.Get_language(db, "other_tool", false) },
            { "bbs/make", tool.Get_language(db, "add", false) },
        },
        map[string]string{},
    )

    return out
}
