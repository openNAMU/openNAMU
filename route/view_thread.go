package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_thread(config tool.Config, topic_num string, doc_name string, page string, values url.Values) string {
	if topic_num != "0" {
		if values != nil {
			api_data := Api_thread_post(config, topic_num, doc_name, values.Get("content"), values.Get("topic"), values.Get("title"))
			response, _ := api_data["response"].(string)
			if response == "require auth" {
				db := tool.DB_connect()
				defer tool.DB_close(db)
				return tool.Get_error_page(db, config, "auth")
			}
			if response != "ok" {
				db := tool.DB_connect()
				defer tool.DB_close(db)
				error_name, _ := api_data["data"].(string)
				if error_name == "" {
					error_name = "error"
				}
				return tool.Get_error_page(db, config, error_name)
			}
			comment_num, _ := api_data["comment_num"].(string)
			return tool.Get_redirect("/bbs/w/" + tool.Url_parser(thread_bbs_id) + "/" + tool.Url_parser(topic_num) + "#" + tool.Url_parser(comment_num))
		}
		return tool.Get_redirect("/bbs/w/" + tool.Url_parser(thread_bbs_id) + "/" + tool.Url_parser(topic_num))
	}

	if values == nil {
		return tool.Get_redirect("/bbs/in/" + tool.Url_parser(thread_bbs_id))
	}

	db := tool.DB_connect()
	defer tool.DB_close(db)

	if doc_name == "" {
		doc_name = "Test"
	}

	api_data := Api_thread_post(config, "0", doc_name, values.Get("content"), values.Get("topic"), values.Get("title"))
	response, _ := api_data["response"].(string)
	if response == "empty data" {
		return tool.Get_redirect("/bbs/in/" + tool.Url_parser(thread_bbs_id))
	}
	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}
	if response != "ok" {
		error_name, _ := api_data["data"].(string)
		if error_name == "" {
			error_name = "error"
		}
		return tool.Get_error_page(db, config, error_name)
	}
	topic_num, _ = api_data["topic_num"].(string)
	comment_num, _ := api_data["comment_num"].(string)
	return tool.Get_redirect("/bbs/w/" + tool.Url_parser(thread_bbs_id) + "/" + tool.Url_parser(topic_num) + "#" + tool.Url_parser(comment_num))
}
