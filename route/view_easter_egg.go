package route

import (
	"net/url"

	"opennamu/route/tool"
)

func View_easter_egg(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if values != nil && !tool.IP_or_user(config.IP) {
		var found string
		if !tool.QueryRow_DB(db, "select id from user_set where id = ? and name = 'get_🥚' limit 1", []any{&found}, config.IP) {
			tool.Exec_DB(db, "insert into user_set (name, id, data) values ('get_🥚', ?, 'Y')", config.IP)
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
