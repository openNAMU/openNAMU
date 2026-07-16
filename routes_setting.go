package main

import (
	"net/http"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_setting_routes(r *gin.Engine) {
	r.GET("/setting", func(c *gin.Context) {
		route_data := route.View_setting(make_route_config(c))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/setting/backlink_reset", func(c *gin.Context) {
		route_data := route.View_setting_backlink_reset(make_route_config(c))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/setting/backlink_reset", func(c *gin.Context) {
		route_data := route.View_setting_backlink_reset_post(make_route_config(c))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
}
