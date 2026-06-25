package route

import "opennamu/route/tool"

func View_list_history_post(config tool.Config, doc_name string, a string, b string) string {
    return tool.Get_redirect("/diff/" + tool.Url_parser(b) + "/" + tool.Url_parser(a) + "/" + tool.Url_parser(doc_name))
}