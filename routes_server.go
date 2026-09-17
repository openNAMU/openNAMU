package main

import (
	"net/http"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func Register_server_routes(r *gin.Engine) {
	r.GET("/restart", func(c *gin.Context) {
		data := route.View_server_action(Make_route_config(c), "restart", false)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/restart", func(c *gin.Context) {
		data := route.View_server_action(Make_route_config(c), "restart", true)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/shutdown", func(c *gin.Context) {
		data := route.View_server_action(Make_route_config(c), "shutdown", false)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/shutdown", func(c *gin.Context) {
		data := route.View_server_action(Make_route_config(c), "shutdown", true)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/update", func(c *gin.Context) {
		data := route.View_server_action(Make_route_config(c), "update", false)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/update", func(c *gin.Context) {
		data := route.View_server_action(Make_route_config(c), "update", true)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
}
