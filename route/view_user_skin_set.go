package route

import "opennamu/route/tool"

func View_user_skin_set(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	data := `<div id="main_skin_set"><h2>` + tool.Get_language(db, "error", true) + `</h2><ul><li>` + tool.Get_language(db, "error_skin_set", true) + `<br>` + tool.Get_language(db, "error_skin_set_old", true) + ` <a href="/skin_set">(` + tool.Get_language(db, "go", true) + `)</a></li></ul></div>`
	menu := [][]any{
		{"change", tool.Get_language(db, "user_setting", true)},
		{"change/skin_set/main", tool.Get_language(db, "main_skin_set", true)},
	}

	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "skin_set", true),
		data,
		[]any{},
		menu,
		map[string]string{},
	)
}
