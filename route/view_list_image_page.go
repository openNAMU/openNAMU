package route

import (
	"strconv"
	"strings"
	"time"

	"opennamu/route/tool"
	"opennamu/route/tool/markup"
)

func View_list_image_page(config tool.Config, page string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	page_num := list_extra_page_number(page)
	offset := (page_num - 1) * 50
	body := strings.Builder{}
	body.WriteString(`<a href="/list/file">(` + tool.Get_language(db, "normal", true) + `)</a><hr class="main_hr">`)

	rows := tool.Get_data_file_rows(db, offset, false)
	render_data := strings.Builder{}
	sub_data := strings.Builder{}
	count := 0
	for rows.Next() {
		name := ""
		if rows.Scan(&name) != nil {
			continue
		}
		if count != 0 && count%4 == 0 {
			render_data.WriteString("||\n")
			render_data.WriteString(sub_data.String())
			render_data.WriteString("||\n")
			sub_data.Reset()
		}
		render_data.WriteString("|| [[" + name + "]] ")
		sub_data.WriteString("|| [[:")
		sub_data.WriteString(name)
		sub_data.WriteString("]] ")
		count++
	}
	rows.Close()

	parameter_data := map[string]any{}
	parameter_data["__opennamu_skin_set"] = get_render_setting_parameter(db, config)
	if config.IP != "" {
		parameter_data["ip"] = config.IP
	}
	render_name := strconv.FormatInt(time.Now().UnixNano(), 10)
	rendered_data := markup.Get_render_direct(db, "", render_data.String(), "namumark", render_name, "view", parameter_data)["data"]
	body.WriteString(get_render_setting_css(db, config) + apply_render_setting_data(db, config, rendered_data))
	body.WriteString(tool.Get_page_control(db, page_num, count, 50, "/list/image/{}"))
	return list_extra_page(db, config, tool.Get_language(db, "image_file_list", true), body.String())
}
