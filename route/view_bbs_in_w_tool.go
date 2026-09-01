package route

import (
	"strings"

	"opennamu/route/tool"
)

func View_bbs_in_w_tool(config tool.Config, set_id string, set_code string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	pinned_name := "pinned"
	if _, exists := tool.Get_bbs_data_value(db, set_id, set_code, "pinned"); exists {
		pinned_name = "pinned_release"
	}

	data_html := `
        <h2>` + tool.Get_language(db, "tool", true) + `</h2>
        <ul>
            <li><a href="/bbs/raw/` + tool.Url_parser(set_id) + `/` + tool.Url_parser(set_code) + `">` + tool.Get_language(db, "raw", true) + `</a></li>
        </ul>
    `

	if tool.Check_permission(db, "bbs_pin", config.IP) {
		data_html += `<h3>` + tool.Get_language(db, "admin", true) + `</h3><ul><li><a href="/bbs/pinned/` + tool.Url_parser(set_id) + `/` + tool.Url_parser(set_code) + `">` + tool.Get_language(db, pinned_name, true) + `</a></li></ul>`
	}
	if tool.Check_permission(db, "bbs_delete", config.IP) {
		data_html += `<h3>` + tool.Get_language(db, "owner", true) + `</h3><ul><li><a href="/bbs/delete/` + tool.Url_parser(set_id) + `/` + tool.Url_parser(set_code) + `">` + tool.Get_language(db, "delete", true) + `</a></li></ul>`
	}

	if tool.Check_permission(db, "bbs_comment_manage", config.IP) {
		comment_closed := bbs_comment_closed(db, set_id, set_code)
		comment_closed_value := "1"
		comment_state := "comment_close"
		if comment_closed {
			comment_closed_value = "0"
			comment_state = "comment_open"
		}

		comment_form := ""
		if set_id == "0" {
			comment_form = `
            <form method="post">
                <input type="hidden" name="action" value="comment_delete_all">
                <span>` + tool.Get_language(db, "delete_warning", true) + `</span><br>
                <button class="__ON_BUTTON__" type="submit">` + tool.Get_language(db, "comment_delete_all", true) + `</button>
            </form>
            `
		} else {
			comment_data := Api_bbs_w_comment_all(config, set_id+"-"+set_code, true, "")
			comments, _ := comment_data["data"].([]map[string]string)
			comment_html := ""
			for _, comment := range comments {
				if comment["comment"] == "" {
					continue
				}

				full_code := comment["id"] + "-" + comment["code"]
				comment_code := strings.TrimPrefix(full_code, set_id+"-"+set_code+"-")
				if comment_code == full_code || !bbs_comment_code_regex.MatchString(comment_code) {
					continue
				}

				comment_html += `<label><input type="checkbox" name="comment_code" value="` + comment_code + `"> #` + comment_code + ` ` + comment["comment_user_id_render"] + ` ` + tool.HTML_escape(comment["comment_date"]) + ` ` + tool.HTML_escape(comment["comment"]) + `</label><br>`
			}

			if comment_html == "" {
				comment_html = `<span>` + tool.Get_language(db, "empty", true) + `</span>`
			} else {
				comment_html += `<br><button class="__ON_BUTTON__" type="submit">` + tool.Get_language(db, "comment_bulk_delete", true) + `</button>`
			}
			comment_form = `<form method="post"><input type="hidden" name="action" value="comment_delete">` + comment_html + `</form>`
		}

		data_html += `
            <h3>` + tool.Get_language(db, "comment_manage", true) + `</h3>
            <form method="post">
                <input type="hidden" name="action" value="comment_close">
                <input type="hidden" name="comment_closed" value="` + comment_closed_value + `">
                <button class="__ON_BUTTON__" type="submit">` + tool.Get_language(db, comment_state, true) + `</button>
            </form>
            ` + comment_form + `
        `
	}

	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "bbs_post_tool", true),
		data_html,
		[]any{},
		[][]any{
			{"bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(set_code), tool.Get_language(db, "return", true)},
		},
		map[string]string{},
	)
}
