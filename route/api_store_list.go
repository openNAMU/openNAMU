package route

import "opennamu/route/tool"

func Api_store_list(config tool.Config) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return map[string]any{
		"response": "ok",
		"data": map[string]any{
			"point": tool.Get_user_point(db, config.IP),
			"items": []map[string]string{},
		},
	}
}
