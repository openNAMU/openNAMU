package route

import "opennamu/route/tool"

func View_login_logout(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if tool.IP_or_user(config.IP) {
		return tool.Get_redirect("/user")
	}

	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "logout", true),
		`<form method="post">
            <button class="__ON_BUTTON__" type="submit">` + tool.Get_language(db, "logout", true) + `</button>
        </form>`,
		[]any{},
		[][]any{
			{"user", tool.Get_language(db, "return", true)},
		},
		map[string]string{},
	)
}
