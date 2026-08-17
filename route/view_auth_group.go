package route

import (
	"net/url"
	"strconv"

	"opennamu/route/tool"
)

func View_auth_group(config tool.Config, name string, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if name == "" && values == nil {
		data := `<form method="post"><input name="group_name" placeholder="` + tool.Get_language(db, "name", true) + `"><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "go", true) + `</button></form>`
		return tool.Get_template(db, config, tool.Get_language(db, "add_admin_group", true), data, []any{}, [][]any{{"auth/list", tool.Get_language(db, "return", true)}}, map[string]string{})
	}
	if values != nil && name == "" {
		name = values.Get("group_name")
	}
	if name == "" || tool.HTML_escape(name) != name {
		return tool.Get_error_page(db, config, "error")
	}

	if values != nil {
		tool.Exec_DB(db, "delete from alist where name = ?", name)
		for _, choice := range auth_choices() {
			if values.Get(choice.key) != "" {
				tool.Exec_DB(db, "insert into alist (name, acl) values (?, ?)", name, choice.key)
			}
		}
		tool.Exec_DB(db, "insert into alist (name, acl) values (?, 'nothing')", name)
		tool.Do_insert_auth_history(db, config.IP, "auth_group_save ("+name+")")
		return tool.Get_redirect("/auth/list/add/" + tool.Url_parser(name))
	}

	selected := map[string]bool{}
	rows := tool.Query_DB(db, "select acl from alist where name = ?", name)
	for rows.Next() {
		value := ""
		if rows.Scan(&value) == nil {
			selected[value] = true
		}
	}
	rows.Close()

	data := `<form method="post"><input type="hidden" name="group_name" value="` + tool.HTML_escape(name) + `">`
	for _, choice := range auth_choices() {
		checked := ""
		if selected[choice.key] {
			checked = ` checked`
		}
		label := tool.Get_language(db, choice.lang, true)
		data += `<div class="opennamu_list_1" style="margin-left:` + strconv.Itoa((choice.level-1)*20) + `px"><label><input type="checkbox" name="` + choice.key + `"` + checked + `> ` + label + `</label></div>`
	}
	data += `<hr class="main_hr"><button type="submit">` + tool.Get_language(db, "save", true) + `</button></form>`

	return tool.Get_template(db, config, name, data, []any{"(" + tool.Get_language(db, "admin_group", true) + ")"}, [][]any{{"auth/list", tool.Get_language(db, "return", true)}}, map[string]string{})
}

type auth_choice struct {
	level int
	key   string
	lang  string
}

func auth_choices() []auth_choice {
	return []auth_choice{
		{1, "owner", "owner_authority"},
		{1, "user", "user_authority"},
		{2, "captcha_pass", "captcha_pass_authority"},
		{2, "ip", "ip_authority"},
		{3, "document", "document_authority"},
		{4, "edit", "edit_authority"},
		{4, "move", "move_authority"},
		{4, "new_make", "new_make_authority"},
		{4, "delete", "delete_authority"},
		{4, "view", "view_authority"},
		{3, "discuss", "discuss_authority"},
		{4, "discuss_make_new_thread", "discuss_make_new_thread_authority"},
		{4, "discuss_view", "discuss_view_authority"},
		{3, "upload", "upload_authority"},
		{3, "vote", "vote_authority"},
		{3, "bbs_use", "bbs_authority"},
		{4, "bbs_edit", "bbs_edit_authority"},
		{4, "bbs_comment", "bbs_comment_authority"},
		{4, "bbs_view", "bbs_view_authority"},
		{3, "captcha_one_check_five_pass", "captcha_one_check_five_pass_authority"},
		{3, "edit_filter_view", "edit_filter_view_authority"},
		{1, "admin", "admin_authority"},
		{2, "ban", "ban_authority"},
		{2, "toron", "discussion_authority"},
		{2, "check", "user_analyze_authority"},
		{2, "view_user_watchlist", "view_user_watchlist_authority"},
		{2, "acl", "document_acl_authority"},
		{2, "hidel", "history_hide_authority"},
		{2, "give", "authorization_authority"},
		{2, "bbs", "bbs_management_authority"},
		{2, "vote_fix", "vote_management_authority"},
		{2, "admin_default_feature", "admin_default_feature_authority"},
		{3, "doc_watch_list_view", "doc_watch_list_view_authority"},
		{3, "treat_as_admin", "treat_as_admin_authority"},
		{3, "view_hide_user_name", "view_hide_user_name_authority"},
		{3, "user_name_bold", "user_name_bold_authority"},
		{3, "multiple_upload", "multiple_upload_authority"},
		{3, "slow_edit_pass", "slow_edit_pass_authority"},
		{3, "edit_bottom_compulsion_pass", "edit_bottom_compulsion_pass_authority"},
		{3, "edit_filter_pass", "edit_filter_pass_authority"},
		{3, "nothing", "nothing_authority"},
	}
}
