package route

import "opennamu/route/tool"

func View_bbs_excellent(config tool.Config, page string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	api_data := Api_bbs(config, "", page, "excellent")
	if api_data["response"] == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}
	if api_data["response"] != "ok" {
		return tool.Get_error_page(db, config, "error")
	}

	posts, ok := api_data["data"].([]map[string]string)
	if !ok {
		return tool.Get_error_page(db, config, "error")
	}
	has_next, ok := api_data["has_next"].(bool)
	if !ok {
		return tool.Get_error_page(db, config, "error")
	}

	board_api_data := Api_bbs_list(config)
	if board_api_data["response"] == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	}
	boards, ok := board_api_data["data"].([][]string)
	if !ok || board_api_data["response"] != "ok" {
		return tool.Get_error_page(db, config, "error")
	}
	board_names := make(map[string]string, len(boards))
	for _, board := range boards {
		if len(board) >= 2 {
			board_names[board[1]] = board[0]
		}
	}

	page_num := tool.Str_to_int(page)
	if page_num < 1 {
		page_num = 1
	}
	page_path := "/bbs/excellent/page/{}"
	data_html := Bbs_list_example_ui(db)
	if len(posts) == 0 {
		data_html += `<p>` + tool.Get_language(db, "no_excellent_posts", true) + `</p>`
	} else {
		data_html += Get_bbs_list_ui(db, config, posts, board_names)
	}
	data_html += tool.Get_page_control(db, page_num, len(posts), 50, page_path, has_next)

	global_excellent_name := tool.Get_language(db, "global_excellent_posts", true)
	return tool.Get_template(
		db,
		config,
		global_excellent_name,
		data_html,
		[]any{},
		[][]any{
			{"bbs/main", tool.Get_language(db, "return", true)},
			{"bbs/search", tool.Get_language(db, "search", true)},
		},
		map[string]string{},
	)
}
