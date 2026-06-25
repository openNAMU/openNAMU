package route

import (
	"opennamu/route/tool"
)

func View_w_random(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    api_data := Api_w_random(config)
    title := api_data["data"]

    redirect := tool.Get_redirect("/w/" + tool.Url_parser(title))

    return redirect
}