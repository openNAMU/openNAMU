package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_user_key(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !user_auth(db, config) {
		return tool.Get_redirect("/user")
	}
	value := ""
	for value == "" {
		value = tool.Get_random_key(128)
		existing := ""
		if tool.QueryRow_DB(db, "select data from user_set where name = ? and data = ?", []any{&existing}, "random_key", value) {
			value = ""
		}
	}
	user_save(db, config.IP, "random_key", value)
	return tool.Get_redirect("/change")
}
