package route

import "opennamu/route/tool"

func Api_vote_add_post(config tool.Config, data string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)

	return return_data
}