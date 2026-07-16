package main

import (
	"net/http"
	"strings"

	"opennamu/route"
	"opennamu/route/tool"

	"github.com/gin-gonic/gin"
)

func register_auth_routes(r *gin.Engine) {
	r.GET("/login", func(c *gin.Context) {
		route_data := route.View_login_login(make_route_config(c))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/register", func(c *gin.Context) {
		route_data := route.View_login_register(make_route_config(c))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/register", func(c *gin.Context) {
		route_data := route.View_login_register_post(
			make_route_config(c),
			c.PostForm("id"),
			c.PostForm("password"),
			c.PostForm("password_check"),
		)
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/login", func(c *gin.Context) {
		route_data := route.View_login_login_post(
			make_route_config(c),
			c.PostForm("id"),
			c.PostForm("password"),
		)
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/logout", func(c *gin.Context) {
		route_data := route.View_login_logout(make_route_config(c))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/logout", func(c *gin.Context) {
		route_data := route.View_login_logout_post(make_route_config(c))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/user", func(c *gin.Context) {
		route_data := route.View_user(make_route_config(c), tool.Get_IP(c))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/user/*user_name", func(c *gin.Context) {
		route_data := route.View_user(make_route_config(c), strings.TrimPrefix(c.Param("user_name"), "/"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
}
