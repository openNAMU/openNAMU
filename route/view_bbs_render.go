package route

import (
	"database/sql"
	"strconv"
	"time"

	"opennamu/route/tool"
	"opennamu/route/tool/markup"
)

func Get_bbs_render(db *sql.DB, set_id string, data string, render_type string, config tool.Config) string {
	parameter_data := map[string]any{}
	parameter_data["__opennamu_skin_set"] = get_render_setting_parameter(db, config)
	if config.IP != "" {
		parameter_data["ip"] = config.IP
	}
	bbs_markup := tool.Get_bbs_set_first_data(db, set_id, "bbs_markup")

	if bbs_markup == "" || bbs_markup == "normal" || bbs_markup == "namumark_beta" {
		rendered_data := markup.Get_render(db, "", data, render_type, parameter_data)["data"]
		return add_render_external_link_target(rendered_data)
	}

	render_name := strconv.FormatInt(time.Now().UnixNano(), 10)
	rendered_data := markup.Get_render_direct(db, "", data, bbs_markup, render_name, render_type, parameter_data)["data"]
	return add_render_external_link_target(rendered_data)
}
