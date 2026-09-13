package main

import (
	"net/http"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_category_routes(r *gin.Engine) {
	r.GET("/category/add/:add_type/:value", func(c *gin.Context) {
		data := route.View_category_manual_add(make_route_config(c), c.Param("add_type"), c.Param("value"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/category/delete/:category/:document/:return_name", func(c *gin.Context) {
		data := route.View_category_manual_delete(make_route_config(c), c.Param("category"), c.Param("document"), c.Param("return_name"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/category/add", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_category_manual_post(make_route_config(c), "add", c.Request.PostForm)))
	})
	r.POST("/category/delete", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_category_manual_post(make_route_config(c), "delete", c.Request.PostForm)))
	})
}
