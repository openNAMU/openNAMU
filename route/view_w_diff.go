package route

import (
	"opennamu/route/tool"
	"strconv"
	"strings"
)

type diff_line_part struct {
    diff_type string
    data string
}

type diff_line struct {
    line_number int
    part_list []diff_line_part
}

func View_w_diff(config tool.Config, doc_name string, before_rev string, after_rev string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    api_data := Api_w_diff(config, doc_name, before_rev, after_rev)
    response := api_data["response"].(string)
    
    if response == "require auth" {
        return tool.Get_error_page(db, config, "auth")
    } else if response != "ok" {
        return tool.Get_redirect("/history/" + tool.Url_parser(doc_name))
    }

    diff_data_list := api_data["data"].([]map[string]string)
    diff_line_list := Get_diff_line_list(diff_data_list)

    data_html := ""
    if len(diff_line_list) > 0 {
        data_html = `<table style="width: 100%; white-space: pre-wrap;"><tr><td colspan="2">r` + tool.HTML_escape(before_rev) + ` &rarr; r` + tool.HTML_escape(after_rev) + `</td></tr>`

        for _, line_data := range diff_line_list {
            line_html := ""

            for _, part_data := range line_data.part_list {
                part_html := tool.HTML_escape(part_data.data)
                switch part_data.diff_type {
                case "insert":
                    part_html = `<span class="opennamu_diff_green">` + part_html + `</span>`
                case "delete":
                    part_html = `<span class="opennamu_diff_red">` + part_html + `</span>`
                }

                line_html += part_html
            }

            if line_html == "" {
                line_html = "&nbsp;"
            }

            data_html += `<tr>
                <td style="width: 40px; user-select: none;">` + strconv.Itoa(line_data.line_number) + `</td>
                <td>` + line_html + `</td>
            </tr>`
        }

        data_html += `</table>`
    }

    return tool.Get_template(
        db,
        config,
        doc_name,
        data_html,
        []any{ "(" + tool.Get_language(db, "compare", true) + ")" },
        [][]any{
            { "history/" + tool.Url_parser(doc_name), tool.Get_language(db, "return", true) },
        },
        map[string]string{},
    )
}

func Get_diff_line_list(diff_data_list []map[string]string) []diff_line {
    return_data := []diff_line{}
    part_list := []diff_line_part{}
    line_number := 1
    line_changed := false

    for _, diff_data := range diff_data_list {
        split_data := strings.Split(diff_data["data"], "\n")

        for index, part_data := range split_data {
            if part_data != "" {
                part_list = append(part_list, diff_line_part{
                    diff_type: diff_data["type"],
                    data: part_data,
                })
            }

            if diff_data["type"] != "equal" && (part_data != "" || index < len(split_data)-1) {
                line_changed = true
            }

            if index < len(split_data) - 1 {
                if line_changed {
                    return_data = append(return_data, diff_line{
                        line_number: line_number,
                        part_list: part_list,
                    })
                }

                part_list = []diff_line_part{}
                line_number += 1
                line_changed = false
            }
        }
    }

    if line_changed {
        return_data = append(return_data, diff_line{
            line_number: line_number,
            part_list: part_list,
        })
    }

    return return_data
}
