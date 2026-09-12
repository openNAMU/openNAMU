package main

import (
	"net/http"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_category_routes(r *gin.Engine) {
	r.POST("/category/add", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_category_manual_post(make_route_config(c), "add", c.Request.PostForm)))
	})
	r.POST("/category/delete", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_category_manual_post(make_route_config(c), "delete", c.Request.PostForm)))
	})
}
