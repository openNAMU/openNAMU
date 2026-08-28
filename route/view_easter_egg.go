package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_easter_egg(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if values != nil {
		result := Api_easter_egg_post(config)
		if result["response"] != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
	}

	return tool.Get_template(
		db,
		config,
		"Easter Egg",
		`<div class="opennamu_easter_egg">🥚</div><form method="post"><button type="submit">🥚</button></form>`,
		[]any{},
		[][]any{{"other", tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}
