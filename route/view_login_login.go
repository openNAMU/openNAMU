package route

import "opennamu/route/tool"

func View_login_login(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.IP_or_user(config.IP) {
		return tool.Get_redirect("/user")
	}

	if !tool.Get_auth_info(db, config.IP)["login_available"] {
		return tool.Get_error_page(db, config, "ban")
	}

	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "login", true),
		`<form method="post">
            <input class="__ON_INPUT__" placeholder="`+tool.Get_language(db, "id", true)+`" name="id" type="text">
            <hr class="main_hr">
            <input class="__ON_INPUT__" placeholder="`+tool.Get_language(db, "password", true)+`" name="password" type="password">
            <hr class="main_hr">
            `+tool.Get_captcha_ui(db, config)+`
            <button class="__ON_BUTTON__" type="submit">`+tool.Get_language(db, "login", true)+`</button>
            `+tool.Get_http_warning(db)+`
        </form>`,
		[]any{},
		[][]any{
			{"user", tool.Get_language(db, "return", true)},
		},
		map[string]string{},
	)
}
