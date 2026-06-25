package route

import (
	"opennamu/route/tool"
)

func Setting_list() map[string]string {
    setting_acl := map[string]string{}

    setting_acl["manage_404_page"] = ""
    setting_acl["manage_404_page_content"] = ""

    setting_acl["bbs_view_acl_all"] = ""
    setting_acl["bbs_acl_all"] = ""
    setting_acl["bbs_edit_acl_all"] = ""
    setting_acl["bbs_comment_acl_all"] = ""

    setting_acl["rankup_condition"] = ""

    return setting_acl
}

func Api_setting(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    setting_acl := Setting_list()

    if val, ok := setting_acl[other_set["set_name"]]; ok {
        if val != "" {
            if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
                return_data := make(map[string]any)
                return_data["response"] = "require auth"

                json_data, _ := json.Marshal(return_data)
                return string(json_data)
            }
        }

        data_coverage := ""
        if val, ok := other_set["coverage"]; ok {
            data_coverage = val
        }

        return_data := make(map[string]any)
        return_data["response"] = "ok"
        return_data["data"] = tool.Get_setting(db, other_set["set_name"], data_coverage)

        json_data, _ := json.Marshal(return_data)
        return string(json_data)
    } else {
        return_data := make(map[string]any)
        return_data["response"] = "not exist"

        json_data, _ := json.Marshal(return_data)
        return string(json_data)
    }
}
