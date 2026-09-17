package main

import (
	"net/http"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func Register_category_routes(r *gin.Engine) {
	r.GET("/category/add/:add_type/:value", func(c *gin.Context) {
		data := route.View_category_manual_add(Make_route_config(c), c.Param("add_type"), c.Param("value"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/category/delete/:category/:document/:return_name", func(c *gin.Context) {
		data := route.View_category_manual_delete(Make_route_config(c), c.Param("category"), c.Param("document"), c.Param("return_name"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/category/add", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_category_manual_post(Make_route_config(c), "add", c.Request.PostForm)))
	})
	r.POST("/category/delete", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_category_manual_post(Make_route_config(c), "delete", c.Request.PostForm)))
	})
}
