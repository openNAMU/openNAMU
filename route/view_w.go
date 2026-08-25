package route

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"opennamu/route/tool"

	"github.com/gin-gonic/gin"
)

func View_w(c *gin.Context, config tool.Config, doc_name string, view_type string) (string, int) {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	do_type := view_type

	Api_w_page_view_post(config, doc_name)

	var render_data string
	status := http.StatusOK
	raw_data := ""
	length_doc := ""
	description := ""

	raw_data_api := Api_w_raw(config, doc_name, "", "")
	raw_response, _ := raw_data_api["response"].(string)
	if raw_response == "require auth" {
		error_data := ""
		tool.QueryRow_DB(db, `select data from other where name = "error_401"`, []any{&error_data})
		if error_data != "" {
			render_data = "<h2>" + tool.Get_language(db, "error", true) + "</h2><ul><li>" + error_data + "</li></ul>"
		} else {
			render_data = "<h2>" + tool.Get_language(db, "error", true) + "</h2><ul><li>" + tool.Get_language(db, "authority_error", true) + "</li></ul>"
		}
		status = http.StatusUnauthorized
	} else if raw_response != "ok" {
		error_data := ""
		tool.QueryRow_DB(db, `select data from other where name = "error_404"`, []any{&error_data})
		if error_data != "" {
			render_data = "<h2>" + tool.Get_language(db, "error", true) + "</h2><ul><li>" + error_data + "</li></ul>"
		} else {
			render_data = "<h2>" + tool.Get_language(db, "error", true) + "</h2><ul><li>" + tool.Get_language(db, "document_404_error", true) + "</li></ul>"
		}
		status = http.StatusNotFound
	} else {
		if do_type != "from" {
			redirect_target, redirect_anchor := view_w_redirect_target(db, doc_name)
			if redirect_target != "" {
				view_w_add_recent_document(config, doc_name)
				return tool.Get_redirect("/w_from/" + tool.Url_parser(redirect_target) + redirect_anchor), http.StatusOK
			}
		}

		raw_data, _ = raw_data_api["data"].(string)
		render_data_api := Api_w_render(config, doc_name, raw_data, "normal", "")
		render_data, _ = render_data_api["data"].(string)
		length_doc = strconv.Itoa(tool.Get_len(raw_data))
		description = tool.Get_slice(strings.ReplaceAll(strings.ReplaceAll(raw_data, "\r", ""), "\n", " "), 0, 200)
	}

	recent_documents := view_w_recent_documents(config)
	if do_type == "from" {
		render_data = view_w_redirect_trace(db, doc_name, recent_documents) + render_data
	}
	recent_documents = view_w_add_recent_document(config, doc_name)
	render_data = view_w_trace(db, config, recent_documents) + render_data
	document_type := ""
	if strings.HasPrefix(doc_name, "user:") {
		document_type = "special"
		render_data = view_w_user_data(db, doc_name) + render_data
	} else if strings.HasPrefix(doc_name, "category:") {
		document_type = "special"
		render_data += view_w_category_data(db, config, doc_name)
	} else if strings.HasPrefix(doc_name, "file:") {
		document_type = "special"
		render_data = view_w_file_data(db, doc_name) + render_data
	} else {
		include_link := ""
		if tool.QueryRow_DB(db, "select link from back where title = ? and type = 'include' limit 1", []any{&include_link}, doc_name) {
			document_type = "include"
		}
		redirect_title := ""
		if tool.QueryRow_DB(db, "select title from back where link = ? and type = 'redirect' limit 1", []any{&redirect_title}, doc_name) {
			document_type = "redirect"
		}
	}

	last_edit := ""
	tool.QueryRow_DB(
		db,
		"select set_data from data_set where doc_name = ? and set_name = 'last_edit' and doc_rev = '' limit 1",
		[]any{&last_edit},
		doc_name,
	)

	if document_type == "" && last_edit != "" {
		warning_days := ""
		tool.QueryRow_DB(db, "select data from other where name = 'outdated_doc_warning_date'", []any{&warning_days})
		warning_days_int := tool.Str_to_int(warning_days)
		if warning_days_int > 0 {
			last_edit_time, parse_error := time.Parse("2006-01-02 15:04:05", last_edit)
			if parse_error == nil && time.Now().After(last_edit_time.AddDate(0, 0, warning_days_int)) {
				warning_text := ""
				tool.QueryRow_DB(db, "select data from other where name = 'outdated_doc_warning'", []any{&warning_text})
				if warning_text == "" {
					warning_text = tool.Get_language(db, "old_page_warning", true)
				}
				render_data = warning_text + `<hr class="main_hr">` + render_data
			}
		}
	}

	body := ""
	tool.QueryRow_DB(db, "select data from other where name = 'body'", []any{&body})
	render_data = body + render_data

	bottom_body := ""
	tool.QueryRow_DB(db, "select data from other where name = 'bottom_body'", []any{&bottom_body})
	render_data += bottom_body

	document_top := ""
	tool.QueryRow_DB(
		db,
		"select set_data from data_set where doc_name = ? and doc_rev = '' and set_name = 'document_top' limit 1",
		[]any{&document_top},
		doc_name,
	)
	render_data = document_top + render_data

	topic := 0
	topic_name := ""
	if tool.QueryRow_DB(db, "select title from rd where title = ? and not stop = 'O' order by date desc limit 1", []any{&topic_name}, doc_name) {
		topic = 1
	}
	history_color := 0
	if status == http.StatusNotFound {
		history_title := ""
		if tool.QueryRow_DB(db, "select title from history where title = ? limit 1", []any{&history_title}, doc_name) {
			history_color = 1
		}
	}
	acl_color := 0
	acl_title := ""
	if tool.QueryRow_DB(db, "select title from acl where title = ? limit 1", []any{&acl_title}, doc_name) {
		acl_color = 1
	}
	menu_acl := 0
	if tool.Check_acl(db, doc_name, "", "document_edit", config.IP) {
		menu_acl = 1
	}

	watch_list := 0
	menu := [][]any{
		{"edit/" + tool.Url_parser(doc_name), tool.Get_language(db, "edit", true), menu_acl},
		{"topic/" + tool.Url_parser(doc_name), tool.Get_language(db, "discussion", true), topic},
		{"history/" + tool.Url_parser(doc_name), tool.Get_language(db, "history", true), history_color},
		{"xref/" + tool.Url_parser(doc_name), tool.Get_language(db, "backlink", true)},
		{"acl/" + tool.Url_parser(doc_name), tool.Get_language(db, "setting", true), acl_color},
	}
	if status == http.StatusNotFound {
		menu[0][1] = tool.Get_language(db, "create", true)
	}
	if do_type == "from" {
		menu = append(menu, []any{"w/" + tool.Url_parser(doc_name), tool.Get_language(db, "pass", true)})
	}
	if strings.HasPrefix(doc_name, "user:") && !strings.Contains(strings.TrimPrefix(doc_name, "user:"), "/") {
		menu = append(menu, []any{"w/" + tool.Url_parser(doc_name+"/"+tool.Get_date()), tool.Get_language(db, "today_doc", true)})
	}
	if strings.HasPrefix(doc_name, "file:") {
		menu = append(menu, []any{"delete_file/" + tool.Url_parser(doc_name), tool.Get_language(db, "file_delete", true)})
	}
	if slash_index := strings.LastIndex(doc_name, "/"); slash_index > 0 {
		menu = append(menu, []any{"w/" + tool.Url_parser(doc_name[:slash_index]), tool.Get_language(db, "upper", true)})
	}
	if view_w_child_exists(db, doc_name) {
		menu = append(menu, []any{"down/" + tool.Url_parser(doc_name), tool.Get_language(db, "sub", true)})
	}
	if !tool.IP_or_user(config.IP) {
		star_doc := ""
		if tool.QueryRow_DB(db, "select data from user_set where name = 'star_doc' and id = ? and data = ? limit 1", []any{&star_doc}, config.IP, doc_name) {
			watch_list = 2
		} else {
			watch_list = 1
		}
		star_text := "☆"
		if watch_list == 2 {
			star_text = "★"
		}
		menu = append(menu, []any{"star_doc_from/" + tool.Url_parser(doc_name), star_text + " " + tool.Get_language(db, "star_doc", true), watch_list - 1})
	}
	if tool.Check_acl(db, "", "", "doc_watch_list_view", config.IP) {
		menu = append(menu, []any{"doc_watch_list/1/" + tool.Url_parser(doc_name), tool.Get_language(db, "watch_user_list", true)})
	}

	enable_comment := tool.Get_setting(db, "enable_comment", "")
	if status == http.StatusOK && len(enable_comment) > 0 && enable_comment[0][0] != "" {
		comment_api := Api_w_comment_ui(config, doc_name)
		if comment_data, ok := comment_api["data"].(string); ok && comment_data != "" {
			render_data += `<div class="opennamu_clearfix"></div>` + comment_data
		}
	}

	view_count := 0
	view_count_api := Api_w_page_view(config, doc_name)
	if view_count_data, ok := view_count_api["data"].(int); ok {
		view_count = view_count_data
	}

	out := tool.Get_template(
		db,
		config,
		doc_name,
		render_data,
		[]any{"", last_edit, watch_list, description, view_count},
		menu,
		map[string]string{
			"path":       c.Request.URL.Path,
			"length_doc": length_doc,
		},
	)

	return out, status
}
