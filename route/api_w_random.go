package route

import (
	"opennamu/route/tool"
)

func Api_w_random(config tool.Config) map[string]string {
    db := tool.DB_connect()
    defer tool.DB_close(db)
    
    title := "Test"
    tool.QueryRow_DB(
        db,
        "select title from data where title not like 'user:%' and title not like 'category:%' and title not like 'file:%' order by random() limit 1",
        []any{ &title },
    )

    new_data := map[string]string{}
    new_data["data"] = title

    return new_data
}
