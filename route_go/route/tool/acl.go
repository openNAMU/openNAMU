package tool

func List_acl(func_type string) []string {
	if func_type == "user_document" {
		return []string{"", "user", "all"}
	} else {
		return []string{"", "all", "user", "admin", "owner", "50_edit", "email", "ban", "before", "30_day", "90_day", "ban_admin", "not_all", "up_to_level_3", "up_to_level_10"}
	}
}
