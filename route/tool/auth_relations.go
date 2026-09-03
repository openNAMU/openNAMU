package tool

const (
	Auth_relation_group      = "group"
	Auth_relation_permission = "permission"
	Auth_relation_auto       = "auto"
)

type Auth_relation struct {
	From string
	To   string
	Type string
}

func Auth_relations() []Auth_relation {
	result := []Auth_relation{}
	add := func(from string, to string, relation_type string) {
		result = append(result, Auth_relation{from, to, relation_type})
	}

	add("owner", "admin", Auth_relation_group)

	add("admin", "email_verified", Auth_relation_permission)
	add("admin", "up_to_level_10", Auth_relation_permission)
	add("admin", "trust_d", Auth_relation_permission)
	add("admin", "history_view", Auth_relation_permission)
	add("admin", "manager_view", Auth_relation_permission)
	add("up_to_level_10", "up_to_level_3", Auth_relation_permission)
	add("email_verified", "user", Auth_relation_permission)
	add("up_to_level_3", "user", Auth_relation_permission)
	add("trust_d", "trust_c", Auth_relation_permission)
	add("trust_c", "trust_a", Auth_relation_permission)
	add("trust_c", "trust_b", Auth_relation_permission)
	add("trust_a", "large_edit", Auth_relation_permission)
	add("give_range", "give", Auth_relation_permission)
	add("trust_a", "user", Auth_relation_permission)
	add("trust_b", "user", Auth_relation_permission)
	add("trust_a", "edit_limit_unlimited", Auth_relation_permission)
	add("trust_b", "edit_limit_unlimited", Auth_relation_permission)
	add("trust_a", "bbs_edit_limit_unlimited", Auth_relation_permission)
	add("trust_b", "bbs_edit_limit_unlimited", Auth_relation_permission)
	add("trust_a", "bbs_comment_limit_unlimited", Auth_relation_permission)
	add("trust_b", "bbs_comment_limit_unlimited", Auth_relation_permission)

	for _, permission := range []string{"toron", "check", "acl", "hidel", "give_range", "give", "bbs", "vote_fix"} {
		add("admin", permission, Auth_relation_permission)
	}
	for _, permission := range []string{"edit_filter_manage", "application_manage", "application_view"} {
		add("admin", permission, Auth_relation_permission)
	}

	for _, permission := range []string{
		"setting_main",
		"setting_main_logo",
		"setting_skin",
		"setting_head",
		"setting_top_menu",
		"setting_phrase",
		"setting_external",
		"setting_robot",
		"setting_sitemap",
		"setting_404",
		"setting_backlink",
		"setting_delete",
		"setting_manage",
		"setting_email_test",
		"server_action",
		"auth_group_manage",
		"auth_fix",
		"rankup_manage",
		"history_manage",
		"record_manage",
		"user_manage",
		"file_delete",
		"document_move_manage",
		"filter_manage",
		"user_edit_filter_manage",
		"bbs_post_manage",
		"auth_private_give",
		"bbs_create",
		"bbs_setting",
		"bbs_delete",
		"bbs_comment_manage",
		"thread_change",
		"thread_delete",
		"thread_comment_delete",
	} {
		add("owner", permission, Auth_relation_permission)
	}

	add("history_manage", "history_view", Auth_relation_permission)
	add("acl", "document_acl_manage", Auth_relation_permission)
	add("acl", "document_bulk_delete", Auth_relation_permission)
	add("bbs", "bbs_manage", Auth_relation_permission)
	add("bbs_manage", "bbs_pin", Auth_relation_permission)
	add("bbs_manage", "bbs_main_view", Auth_relation_permission)
	add("toron", "thread_manage", Auth_relation_permission)
	for _, permission := range []string{"thread_setting", "thread_acl", "thread_comment_manage"} {
		add("thread_manage", permission, Auth_relation_permission)
	}
	add("vote_fix", "vote_manage", Auth_relation_permission)
	add("check", "view_user_watchlist", Auth_relation_permission)

	for _, permission := range []string{"toron", "check", "acl", "hidel", "give_range", "give", "bbs", "vote_fix"} {
		add(permission, "admin_default_feature", Auth_relation_permission)
	}
	for _, permission := range []string{
		"treat_as_admin",
		"user_name_bold",
		"multiple_upload",
		"slow_edit_pass",
		"edit_bottom_compulsion_pass",
		"view_hide_user_name",
		"doc_watch_list_view",
		"edit_filter_pass",
		"user",
	} {
		add("admin_default_feature", permission, Auth_relation_permission)
	}

	add("document", "edit", Auth_relation_permission)
	add("document", "move", Auth_relation_permission)
	add("document", "new_make", Auth_relation_permission)
	add("document", "delete", Auth_relation_permission)
	add("edit", "view", Auth_relation_permission)
	add("move", "view", Auth_relation_permission)
	add("new_make", "view", Auth_relation_permission)
	add("delete", "view", Auth_relation_permission)
	add("view", "site_view", Auth_relation_permission)

	add("discuss", "discuss_view", Auth_relation_permission)
	add("discuss", "discuss_make_new_thread", Auth_relation_permission)

	add("bbs_use", "bbs_edit", Auth_relation_permission)
	add("bbs_use", "bbs_comment", Auth_relation_permission)
	add("bbs_use", "bbs_main_view", Auth_relation_permission)
	add("bbs_edit", "bbs_view", Auth_relation_permission)
	add("bbs_comment", "bbs_view", Auth_relation_permission)

	for _, prefix := range []string{"edit", "bbs_edit", "bbs_comment"} {
		add(prefix, prefix+"_limit_10", Auth_relation_permission)
		add(prefix+"_limit_unlimited", prefix+"_limit_100", Auth_relation_permission)
		add(prefix+"_limit_100", prefix+"_limit_50", Auth_relation_permission)
		add(prefix+"_limit_50", prefix+"_limit_10", Auth_relation_permission)
	}

	for _, permission := range []string{
		"rankup",
		"do_email_verified",
		"captcha_pass",
		"ip",
		"edit_limit_100",
		"bbs_edit_limit_100",
		"bbs_comment_limit_100",
	} {
		add("user", permission, Auth_relation_auto)
	}
	for _, permission := range []string{
		"document",
		"discuss",
		"upload",
		"vote",
		"bbs_use",
		"captcha_one_check_five_pass",
		"edit_filter_view",
		"login_available",
		"register_available",
		"history_view",
		"image_view",
		"edit_day",
		"edit_night",
		"edit_limit_50",
		"bbs_edit_limit_50",
		"bbs_comment_limit_50",
	} {
		add("ip", permission, Auth_relation_auto)
	}

	return result
}
