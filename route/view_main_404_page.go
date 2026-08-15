package route

import "opennamu/route/tool"

func Get_frontpage_url() string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	frontpage := "FrontPage"

	tool.QueryRow_DB(
		db,
		`select data from other where name = "frontpage"`,
		[]any{&frontpage},
	)

	return "/w/" + tool.Url_parser(frontpage)
}

func View_main_404_page(config tool.Config, url string) string {
	if url == "/" {
		return tool.Get_redirect(Get_frontpage_url())
	}

	db := tool.DB_connect()
	defer tool.DB_close(db)

	page_404_set := ""
	tool.QueryRow_DB(
		db,
		`select data from other where name = "manage_404_page"`,
		[]any{&page_404_set},
	)

	data_html := ""

	page_404_dir := "404.html"
	if tool.File_exist_check(page_404_dir) && page_404_set == "404_file" {
		data_html = tool.File_text_read(page_404_dir)
	} else {
		db_data := ""

		tool.QueryRow_DB(
			db,
			`select data from other where name = "manage_404_page_content"`,
			[]any{&db_data},
		)

		if db_data != "" {
			data_html = tool.Get_template(
				db,
				config,
				"404",
				db_data,
				[]any{},
				[][]any{},
				map[string]string{},
			)
		} else {
			data_in := tool.Get_language(db, "func_404_error", true)
			data_in += "<hr class=\"main_hr\">"
			data_in += "Path : " + tool.HTML_escape(url)

			data_html = tool.Get_template(
				db,
				config,
				"404",
				data_in,
				[]any{},
				[][]any{},
				map[string]string{},
			)
		}
	}

	return data_html
}
