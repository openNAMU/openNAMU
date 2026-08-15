package route

import (
	"opennamu/route/tool"
)

func Api_user_info(config tool.Config, ip string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	data_result := map[string]any{}

	ip_render := tool.IP_parser(db, ip, config.IP)
	auth_name := tool.Get_user_auth(db, ip)
	level_data := tool.Get_level(db, ip)
	ban_check := tool.Get_user_ban(db, ip, "")
	user_document := tool.Get_user_document(db, ip)

	data_result["render"] = ip_render

	data_result["auth"] = auth_name
	data_result["auth_date"] = tool.Get_auth_date(db, ip)

	data_result["level"] = level_data[0]
	data_result["exp"] = level_data[1]
	data_result["max_exp"] = level_data[2]

	ban_data := any("0")
	if len(ban_check) > 0 && ban_check[0] != "" {
		ban_data = ban_check
	}
	data_result["ban"] = ban_data

	document_data := "0"
	if user_document {
		document_data = "1"
	}
	data_result["document"] = document_data

	data_result["user_title"] = tool.Get_user_title(db, ip)

	return_data := make(map[string]any)
	return_data["response"] = "ok"
	language_name_list := []string{
		"user_name",
		"authority",
		"state",
		"member",
		"normal",
		"blocked",
		"type",
		"regex",
		"period",
		"limitless",
		"login_able",
		"why",
		"band_blocked",
		"ip",
		"ban",
		"level",
		"option",
		"edit_request_able",
		"cidr",
	}
	language_data := map[string]string{}
	for _, name := range language_name_list {
		language_data[name] = tool.Get_language(db, name, false)
	}
	return_data["language"] = language_data
	return_data["data"] = data_result

	return return_data
}
