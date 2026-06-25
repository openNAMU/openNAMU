package route

import (
	"database/sql"
	"opennamu/route/tool"
)

func Get_main_skin_set_list(db *sql.DB) map[string][][]string {
    need_keys := []string{
        "default", "off", "change_to_normal", "delete", "use",
        "bottom", "top", "normal", "spread", "popup", "not_working", "popover",
        "change_to_link", "click_load", "all_off", "in_content",
        "self_tab", "all", "only_number",
    }

    lang := make(map[string]string, len(need_keys))
    for _, k := range need_keys {
        lang[k] = tool.Get_language(db, k, true)
    }

    popup_not_working := lang["popup"] + " (" + lang["not_working"] + ")"

    set_list := map[string][][]string{
        "main_css_strike": {
            { "default", lang["default"] },
            { "normal", lang["off"] },
            { "change", lang["change_to_normal"] },
            { "delete", lang["delete"] },
        },
        "main_css_bold": {
            { "default", lang["default"] },
            { "normal", lang["off"] },
            { "change", lang["change_to_normal"] },
            { "delete", lang["delete"] },
        },
        "main_css_include_link": {
            { "default", lang["default"] },
            { "normal", lang["off"] },
            { "use", lang["use"] },
        },
        "main_css_category_set": {
            { "default", lang["default"] },
            { "bottom", lang["bottom"] },
            { "top", lang["top"] },
        },
        "main_css_footnote_set": {
            { "default", lang["default"] },
            { "normal", lang["normal"] },
            { "spread", lang["spread"] },
            { "popup", popup_not_working },
            { "popover", lang["popover"] },
        },
        "main_css_image_set": {
            { "default", lang["default"] },
            { "normal", lang["normal"] },
            { "click", lang["change_to_link"] },
            { "new_click", lang["click_load"] },
        },
        "main_css_toc_set": {
            { "default", lang["default"] },
            { "normal", lang["normal"] },
            { "off", lang["all_off"] },
            { "half_off", lang["in_content"] },
        },
        "main_css_monaco": {
            { "default", lang["default"] },
            { "normal", lang["off"] },
            { "use", lang["use"] },
        },
        "main_css_exter_link": {
            { "default", lang["default"] },
            { "blank", lang["normal"] },
            { "self", lang["self_tab"] },
        },
        "main_css_link_delimiter": {
            { "default", lang["default"] },
            { "normal", lang["off"] },
            { "use", lang["use"] },
        },
        "main_css_darkmode": {
            { "default", lang["default"] },
            { "0", lang["off"] },
            { "1", lang["use"] },
        },
        "main_css_footnote_number": {
            { "default", lang["default"] },
            { "all", lang["all"] },
            { "only_number", lang["only_number"] },
        },
        "main_css_view_real_footnote_num": {
            { "default", lang["default"] },
            { "off", lang["off"] },
            { "on", lang["use"] },
        },
        "main_css_table_scroll": {
            { "default", lang["default"] },
            { "off", lang["off"] },
            { "on", lang["use"] },
        },
        "main_css_category_change_title": {
            { "default", lang["default"] },
            { "off", lang["off"] },
            { "on", lang["use"] },
        },
        "main_css_list_view_change": {
            { "default", lang["default"] },
            { "off", lang["off"] },
            { "on", lang["use"] },
        },
        "main_css_view_joke": {
            { "default", lang["default"] },
            { "on", lang["use"] },
            { "off", lang["off"] },
        },
        "main_css_math_scroll": {
            { "default", lang["default"] },
            { "off", lang["off"] },
            { "on", lang["use"] },
        },
        "main_css_view_history": {
            { "default", lang["default"] },
            { "off", lang["off"] },
            { "on", lang["use"] },
        },
        "main_css_table_transparent": {
            { "default", lang["default"] },
            { "off", lang["off"] },
            { "on", lang["use"] },
        },
        "main_css_font_size": {
            { "default", lang["default"] },
            { "10", "10" },
            { "12", "12" },
            { "14", "14" },
            { "16", "16" },
            { "18", "18" },
            { "20", "20" },
            { "22", "22" },
        },
    }

    return set_list
}

func Api_user_setting_skin_set_main_post(config tool.Config, user_set_list map[string]string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    return_data := make(map[string]any)

    if tool.IP_or_user(config.IP) {
        return_data["response"] = "require auth"

        return return_data
    }

    return_data["response"] = "ok"

    set_list := Get_main_skin_set_list(db)
    for k := range set_list {
        if val, ok := user_set_list[k]; ok {
            tool.Exec_DB(
                db,
                "delete from user_set where id = ? and name = ? and data = ?",
                config.IP,
                k,
                val,
            )
            tool.Exec_DB(
                db,
                "insert into user_set (name, id, data) values (?, ?, ?)",
                k,
                config.IP,
                val,
            )
        }
    }

    return return_data
}