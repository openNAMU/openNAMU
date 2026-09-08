package main

import (
	"net/http"

	"opennamu/route"
	"opennamu/route/tool"

	"github.com/gin-gonic/gin"
)

func register_user_multiple_routes(r *gin.Engine) {
	render := func(c *gin.Context, page string, sort string, search string) {
		data := route.View_user_multiple(make_route_config(c), page, sort, search)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	}

	r.GET("/list/user/multiple", func(c *gin.Context) {
		render(c, "1", "recent", "")
	})
	r.POST("/list/user/multiple", func(c *gin.Context) {
		render(c, "1", c.PostForm("sort"), c.PostForm("search"))
	})
	r.GET("/list/user/multiple/:sort/:page", func(c *gin.Context) {
		render(c, c.Param("page"), c.Param("sort"), "")
	})
	r.GET("/list/user/multiple/:sort/:page/:search", func(c *gin.Context) {
		search, err := tool.Get_base64_decode(c.Param("search"))
		if err != nil {
			search = ""
		}
		render(c, c.Param("page"), c.Param("sort"), search)
	})
}
