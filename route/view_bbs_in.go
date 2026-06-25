package route

import (
	"opennamu/route/tool"
)

func View_bbs_in(config tool.Config, set_id string, page_num string) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    bbs_name := Api_bbs_num_to_name(db, set_id)

    data_api := Api_bbs(config, set_id, page_num)
    data_api_in := data_api["data"].([]map[string]string)

    data_html := Get_bbs_list_ui(config, data_api_in, map[string]string{})

    out := tool.Get_template(
        db,
        config,
        bbs_name,
        data_html,
        []any{},
        [][]any{
            { "bbs/main", tool.Get_language(db, "return", true) },
            { "bbs/edit/" + tool.Url_parser(set_id), tool.Get_language(db, "add", true) },
            { "bbs/set/" + tool.Url_parser(set_id), tool.Get_language(db, "bbs_set", true) },
        },
        map[string]string{},
    )

    return out
}
