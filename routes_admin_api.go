package main

import (
	"net/http"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_admin_api_routes(r *gin.Engine) {
	r.GET("/api/v2/list/auth", func(c *gin.Context) {
		c.JSON(http.StatusOK, route.Api_list_auth(make_route_config(c)))
	})
	r.GET("/api/v2/list/auth_give/:num", func(c *gin.Context) {
		c.JSON(http.StatusOK, route.Api_list_auth_give(make_route_config(c), c.Param("num"), false))
	})
	r.GET("/api/v2/list/auth_give/all/:num", func(c *gin.Context) {
		c.JSON(http.StatusOK, route.Api_list_auth_give(make_route_config(c), c.Param("num"), true))
	})
	r.GET("/api/v2/auth", func(c *gin.Context) {
		c.JSON(http.StatusOK, route.Api_func_auth(make_route_config(c), ""))
	})
	r.GET("/api/v2/auth/*user_name", func(c *gin.Context) {
		user_name := strings.TrimPrefix(c.Param("user_name"), "/")
		c.JSON(http.StatusOK, route.Api_func_auth(make_route_config(c), user_name))
	})
	r.PATCH("/api/v2/auth/give", func(c *gin.Context) {
		data := route.Api_give_auth_patch(
			make_route_config(c),
			c.PostForm("auth"),
			c.PostForm("change_auth"),
			c.PostForm("user_name"),
			c.PostForm("end_date"),
			c.PostForm("target_type"),
			c.PostForm("why"),
			c.PostForm("action") == "release",
		)
		c.JSON(http.StatusOK, data)
	})
}
