package main

import (
	"net/http"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_auth_give_list_routes(r *gin.Engine) {
	r.GET("/auth/give_list", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_give_list(make_route_config(c), "1", false)))
	})
	r.GET("/auth/give_list/all", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_give_list(make_route_config(c), "1", true)))
	})
	r.GET("/auth/give_list/all/:num", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_give_list(make_route_config(c), c.Param("num"), true)))
	})
	r.GET("/auth/give_list/:num", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_give_list(make_route_config(c), c.Param("num"), false)))
	})
}
