package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_thread(config tool.Config, topic_num string, doc_name string, page string, values url.Values) string {
	new_topic := topic_num == "0"
	if values == nil {
		if new_topic {
			return tool.Get_redirect("/bbs/in/" + tool.Url_parser(thread_bbs_id))
		}
		return tool.Get_redirect("/bbs/w/" + tool.Url_parser(thread_bbs_id) + "/" + tool.Url_parser(topic_num))
	}

	if new_topic && doc_name == "" {
		doc_name = "Test"
	}

	api_data := Api_thread_bbs_post(config, topic_num, doc_name, values.Get("content"), values.Get("topic"), values.Get("title"))
	response, _ := api_data["response"].(string)
	if new_topic && response == "empty data" {
		return tool.Get_redirect("/bbs/in/" + tool.Url_parser(thread_bbs_id))
	}

	error_page := func(error_name string) string {
		db := tool.DB_connect()
		defer tool.DB_close(db)
		return tool.Get_error_page(db, config, error_name)
	}
	if response == "require auth" {
		return error_page("auth")
	}
	if response != "ok" {
		error_name, _ := api_data["data"].(string)
		if error_name == "" {
			error_name = "error"
		}
		return error_page(error_name)
	}

	if new_topic {
		topic_num, _ = api_data["topic_num"].(string)
	}
	comment_num, _ := api_data["comment_num"].(string)
	return tool.Get_redirect("/bbs/w/" + tool.Url_parser(thread_bbs_id) + "/" + tool.Url_parser(topic_num) + "#" + tool.Url_parser(comment_num))
}
