package route

import (
	"database/sql"
	"strings"

	"opennamu/route/tool"
)

func view_w_recent_documents(config tool.Config) []string {
	if config.Session == nil {
		return []string{}
	}

	documents, ok := config.Session.Get("lastest_document").([]string)
	if !ok {
		return []string{}
	}
	return append([]string{}, documents...)
}

func view_w_add_recent_document(config tool.Config, doc_name string) []string {
	documents := view_w_recent_documents(config)
	if doc_name == "" || config.Session == nil {
		return documents
	}
	if len(documents) >= 10 {
		documents = documents[len(documents)-9:]
	}
	documents = append(documents, doc_name)

	seen := map[string]bool{}
	unique_documents := []string{}
	for index := len(documents) - 1; index >= 0; index-- {
		if seen[documents[index]] {
			continue
		}
		seen[documents[index]] = true
		unique_documents = append(unique_documents, documents[index])
	}
	for left, right := 0, len(unique_documents)-1; left < right; left, right = left+1, right-1 {
		unique_documents[left], unique_documents[right] = unique_documents[right], unique_documents[left]
	}

	config.Session.Set("lastest_document", unique_documents)
	_ = config.Session.Save()
	return unique_documents
}

func view_w_redirect_trace(db *sql.DB, doc_name string, documents []string) string {
	last_page := ""
	for index := len(documents) - 1; index >= 0; index-- {
		last_page = documents[index]
		redirect_exists := ""
		if tool.QueryRow_DB(db, "select link from back where (title = ? or link = ?) and type = 'redirect' limit 1", []any{&redirect_exists}, last_page, last_page) {
			break
		}
	}
	if last_page == "" || last_page == doc_name {
		return ""
	}

	redirect_text := "{0} ➤ {1}"
	tool.QueryRow_DB(db, "select data from other where name = 'redirect_text'", []any{&redirect_text})
	if redirect_text == "" {
		redirect_text = "{0} ➤ {1}"
	}
	redirect_text = strings.ReplaceAll(redirect_text, "{0}", `<a href="/w_from/`+tool.Url_parser(last_page)+`">`+tool.HTML_escape(last_page)+`</a>`)
	redirect_text = strings.ReplaceAll(redirect_text, "{1}", `<b>`+tool.HTML_escape(doc_name)+`</b>`)
	return `<div class="opennamu_redirect" id="redirect">` + redirect_text + `</div><hr class="main_hr">`
}

func view_w_trace(db *sql.DB, config tool.Config, documents []string) string {
	if tool.Get_main_skin_set(db, config, "main_css_view_history") != "on" {
		return ""
	}

	trace_list := []string{}
	for index := len(documents) - 1; index >= 0; index-- {
		trace_list = append(trace_list, `<a href="/w/`+tool.Url_parser(documents[index])+`">`+tool.HTML_escape(documents[index])+`</a>`)
	}
	return `<div class="opennamu_trace"><a class="opennamu_trace_button" href="javascript:opennamu_do_trace_spread();"> (+)</a> ` + tool.Get_language(db, "trace", true) + ` : ` + strings.Join(trace_list, " ← ") + `</div><hr class="main_hr">`
}
