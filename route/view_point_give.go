package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_point_give(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "admin", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	if values != nil {
		result := Api_point_give_post(config, values.Get("user_name"), values.Get("amount"))
		if result["response"] == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if result["response"] != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
		return tool.Get_redirect("/point/give")
	}

	data := `<form method="post"><input name="user_name" placeholder="` + tool.Get_language(db, "user_name", true) + `"><hr class="main_hr"><input type="number" name="amount" placeholder="` + tool.Get_language(db, "point", true) + `"><hr class="main_hr"><button type="submit">` + tool.Get_language(db, "send", true) + `</button></form>`
	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "point_give", true),
		data,
		[]any{},
		[][]any{{"manager", tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}
