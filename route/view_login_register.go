package route

import (
	"opennamu/route/tool"
)

func View_login_register(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    contract := ""
    tool.QueryRow_DB(
        db,
        `select data from other where name = "contract"`,
        []any{ &contract },
    )

    if contract != "" {
        contract += `<hr class="main_hr">`
    }

    password_length_limit := ""
    tool.QueryRow_DB(
        db,
        `select data from other where name = 'password_min_length'`,
        []any{ &password_length_limit },
    )
    
    if password_length_limit != "" {
        password_length_limit = " (" + tool.Get_language(db, "password_min_length", true) + " : " + password_length_limit + ")"
    }

    return tool.Get_template(
        db,
        config,
        tool.Get_language(db, "register", true),
        `<form method="post">
            ` + contract + `

            <input class="__ON_INPUT__" placeholder="` + tool.Get_language(db, "id", true) + `" name="id" type="text">
            <hr class="main_hr">

            <input class="__ON_INPUT__" placeholder="` + tool.Get_language(db, "password", true) + `" name="password" type="password">
            <hr class="main_hr">

            <input class="__ON_INPUT__" placeholder="` + tool.Get_language(db, "password_confirm", true) + `" name="password_check" type="password">
            <hr class="main_hr">

            ` + tool.Get_captcha_ui(db, config) + `

            <button class="__ON_BUTTON__" type="submit">` + tool.Get_language(db, "save", true) + `</button>

            ` + tool.Get_http_warning(db) + `
        </form>`,
        []any{},
        [][]any{
            { "user", tool.Get_language(db, "return", true) },
        },
        map[string]string{},
    )
}