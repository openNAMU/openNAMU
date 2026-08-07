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

func Api_setting(config tool.Config, set_name string, coverage string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	setting_acl := Setting_list()
	return_data := make(map[string]any)

	if val, ok := setting_acl[set_name]; ok {
		if val != "" {
			if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
				return_data["response"] = "require auth"

				return return_data
			}
		}

		return_data["response"] = "ok"
		return_data["data"] = tool.Get_setting(db, set_name, coverage)

		return return_data
	} else {
		return_data["response"] = "error"
		return_data["data"] = "not exist"

		return return_data
	}
}
