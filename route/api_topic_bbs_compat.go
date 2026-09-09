package route

import "opennamu/route/tool"

func api_topic_bbs(config tool.Config, tool_name string, topic_num string, s_num string, e_num string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	thread_data := Api_thread_bbs(config, tool_name, topic_num, s_num, e_num)
	if tool_name == "length" {
		length, _ := thread_data["comment"].(string)
		return map[string]any{
			"response": thread_data["response"],
			"length":   length,
		}
	}

	comments, _ := thread_data["data"].([]map[string]string)
	data_list := []map[string]string{}
	for _, comment := range comments {
		data_list = append(data_list, map[string]string{
			"id":        comment["code"],
			"data":      comment["comment"],
			"date":      comment["comment_date"],
			"ip":        comment["comment_user_id"],
			"ip_render": comment["comment_user_id_render"],
			"blind":     comment["blind"],
		})
	}

	return map[string]any{
		"response": thread_data["response"],
		"data":     data_list,
		"language": map[string]string{
			"tool":   tool.Get_language(db, "tool", false),
			"render": tool.Get_language(db, "render", false),
		},
	}
}
