package route

import (
	"database/sql"

	"opennamu/route/tool"
	"opennamu/route/tool/markup"
)

func get_thread_ui(db *sql.DB, user_name string, date string, data string, code string, color string, blind string, add_style string, topic_num string, config tool.Config) string {
	rendered_data := ""
	if data != "" {
		parameter_data := map[string]any{}
		parameter_data["__opennamu_skin_set"] = get_render_setting_parameter(db, config)
		if config.IP != "" {
			parameter_data["ip"] = config.IP
		}
		rendered_data = markup.Get_render(db, "", data, "thread", parameter_data)["data"]
		rendered_data = render_topic_reference(rendered_data, topic_num, "", "", "thread")
		rendered_data = add_render_external_link_target(rendered_data)
	}
	return get_thread_ui_with_render(db, user_name, date, rendered_data, code, color, blind, add_style, topic_num)
}

func get_thread_ui_with_render(db *sql.DB, user_name string, date string, rendered_data string, code string, color string, blind string, add_style string, topic_num string) string {
	color_b := ""
	class_b := ""

	if blind == "O" {
		if rendered_data == "" {
			color_b = "opennamu_comment_blind"
		} else {
			color_b = "opennamu_comment_blind_admin"
		}

		class_b = "opennamu_comment_blind_js opennamu_list_hidden"
	} else {
		color_b = "opennamu_comment_blind_not"
	}

	return `
        <span class="` + class_b + `">
            <table class="opennamu_comment" style="` + add_style + `">
                <tr>
                    <td class="opennamu_comment_color_` + color + `">
                        <a href="#thread_shortcut" id="` + code + `">#` + code + `</a>
                        ` + user_name + `
                        <span style="float: right;">` + date + `</span>
                    </td>
                </tr>
                <tr>
                    <td class="` + color_b + ` opennamu_comment_data_main" id="thread_` + code + `">
                        <div class="opennamu_comment_scroll" id="opennamu_thread_render_` + code + `">` + rendered_data + `</div>
                    </td>
                </tr>
            </table>
            <hr class="main_hr">
        </span>
    `
}
