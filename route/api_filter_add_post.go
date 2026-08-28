package route

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"opennamu/route/tool"

	"github.com/dlclark/regexp2"
)

func Api_filter_add_post(config tool.Config, kind string, name string, values url.Values) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	spec, ok := get_filter_spec(kind)
	if !ok {
		return_data["response"] = "error"
		return return_data
	}
	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	title := values.Get("title")
	if title == "" {
		title = values.Get("name")
	}
	if title == "" {
		title = name
	}
	if title == "" && kind != "external_image" && kind != "html" {
		title = "test"
	}

	if kind == "external_image" {
		title = strings.ToLower(strings.TrimSpace(title))
		parsed, err := url.Parse("https://" + title)
		if title == "" || err != nil || parsed.Host != title || parsed.Hostname() == "" || parsed.Port() != "" || parsed.User != nil {
			return_data["response"] = "error"
			return return_data
		}
		if name != "" && name != title {
			tool.Exec_DB(db, "delete from html_filter where html = ? and kind = ?", name, spec.db_kind)
		}
		tool.Exec_DB(db, "delete from html_filter where html = ? and kind = ?", title, spec.db_kind)
		tool.Exec_DB(db, "insert into html_filter (html, kind, plus, plus_t) values (?, ?, '', '')", title, spec.db_kind)
	} else if kind == "html" {
		title = strings.ToLower(strings.TrimSpace(title))
		if !html_filter_tag_regex.MatchString(title) || html_filter_blocked_tags[title] {
			return_data["response"] = "error"
			return return_data
		}
		if name != "" && name != title {
			tool.Exec_DB(db, "delete from html_filter where html = ? and kind = ?", name, spec.db_kind)
		}
		tool.Exec_DB(db, "delete from html_filter where html = ? and kind = ?", title, spec.db_kind)
		tool.Exec_DB(db, "insert into html_filter (html, kind, plus, plus_t) values (?, ?, '', '')", title, spec.db_kind)
	} else if kind == "inter_wiki" || kind == "outer_link" {
		if name != "" && name != title {
			tool.Exec_DB(db, "delete from html_filter where html = ? and kind = ?", name, spec.db_kind)
		}
		tool.Exec_DB(db, "delete from html_filter where html = ? and kind = ?", title, spec.db_kind)
		tool.Exec_DB(db, "insert into html_filter (html, plus, plus_t, kind) values (?, ?, ?, ?)", title, values.Get("link"), values.Get("icon"), spec.db_kind)
		if kind == "inter_wiki" {
			if name != "" && name != title {
				tool.Exec_DB(db, "delete from html_filter where html = ? and kind = 'inter_wiki_sub'", name)
			}
			tool.Exec_DB(db, "delete from html_filter where html = ? and kind = 'inter_wiki_sub'", title)
			tool.Exec_DB(db, "insert into html_filter (html, plus, plus_t, kind) values (?, 'inter_wiki_type', ?, 'inter_wiki_sub')", title, values.Get("inter_type"))
		}
	} else if kind == "edit_filter" {
		content := values.Get("content")
		if _, err := regexp.Compile(content); err != nil {
			return_data["response"] = "error"
			return return_data
		}
		filter_name := name
		if filter_name == "" {
			filter_name = title
		}
		end := "X"
		if days, err := strconv.Atoi(strings.TrimSpace(values.Get("day"))); err == nil && days > 0 {
			end = strconv.Itoa(days * 24 * 60 * 60)
		}
		tool.Exec_DB(db, "delete from html_filter where html = ? and kind = 'regex_filter'", filter_name)
		tool.Exec_DB(db, "insert into html_filter (html, plus, plus_t, kind) values (?, ?, ?, 'regex_filter')", filter_name, content, end)
	} else if kind == "document" {
		acl_data, acl_ok := document_filter_acl_data(db, values.Get("acl"))
		if !acl_ok {
			return_data["response"] = "error"
			return return_data
		}
		doc_name := values.Get("name")
		if doc_name == "" {
			return_data["response"] = "redirect"
			return return_data
		}
		if _, err := regexp2.Compile(values.Get("regex"), 0); err != nil {
			return_data["response"] = "error"
			return return_data
		}
		if name != "" && name != doc_name {
			tool.Exec_DB(db, "delete from html_filter where html = ? and kind = 'document'", name)
		}
		tool.Exec_DB(db, "delete from html_filter where html = ? and kind = 'document'", doc_name)
		tool.Exec_DB(db, "insert into html_filter (html, kind, plus, plus_t) values (?, 'document', ?, ?)", doc_name, values.Get("regex"), acl_data)
	} else {
		plus := ""
		switch kind {
		case "name_filter", "file_filter":
			if _, err := regexp.Compile(title); err != nil {
				return_data["response"] = "error"
				return return_data
			}
		case "extension_filter":
			plus = regexp.MustCompile(`[^0-9]`).ReplaceAllString(values.Get("max_file_size"), "")
		case "template":
			plus = values.Get("exp")
		case "edit_top":
			plus = values.Get("markup")
		}
		if name != "" && name != title {
			tool.Exec_DB(db, "delete from html_filter where html = ? and kind = ?", name, spec.db_kind)
		}
		tool.Exec_DB(db, "delete from html_filter where html = ? and kind = ?", title, spec.db_kind)
		tool.Exec_DB(db, "insert into html_filter (html, kind, plus, plus_t) values (?, ?, ?, '')", title, spec.db_kind, plus)
	}

	tool.Do_insert_auth_history(db, config.IP, "filter_save ("+kind+")")
	return_data["response"] = "ok"
	return return_data
}
