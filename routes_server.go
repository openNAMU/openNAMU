package main

import (
	"net/http"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_server_routes(r *gin.Engine) {
	r.GET("/restart", func(c *gin.Context) {
		data := route.View_server_action(make_route_config(c), "restart", false)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/restart", func(c *gin.Context) {
		data := route.View_server_action(make_route_config(c), "restart", true)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/shutdown", func(c *gin.Context) {
		data := route.View_server_action(make_route_config(c), "shutdown", false)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/shutdown", func(c *gin.Context) {
		data := route.View_server_action(make_route_config(c), "shutdown", true)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/update", func(c *gin.Context) {
		data := route.View_server_action(make_route_config(c), "update", false)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/update", func(c *gin.Context) {
		data := route.View_server_action(make_route_config(c), "update", true)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
}
