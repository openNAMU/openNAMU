package main

import (
	"net/http"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_auth_routes(r *gin.Engine) {
	r.GET("/login", func(c *gin.Context) {
		route_data := route.View_login_login(make_route_config(c))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/register", func(c *gin.Context) {
		route_data := route.View_login_register(make_route_config(c))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/register", func(c *gin.Context) {
		password := c.PostForm("password")
		if password == "" {
			password = c.PostForm("pw")
		}
		password_check := c.PostForm("password_check")
		if password_check == "" {
			password_check = c.PostForm("pw2")
		}
		route_data := route.View_login_register_post_full(
			make_route_config(c),
			c.PostForm("id"),
			password,
			password_check,
			captcha_response(c),
		)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/login", func(c *gin.Context) {
		password := c.PostForm("password")
		if password == "" {
			password = c.PostForm("pw")
		}
		route_data := route.View_login_login_post_full(
			make_route_config(c),
			c.PostForm("id"),
			password,
			captcha_response(c),
		)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/logout", func(c *gin.Context) {
		route_data := route.View_login_logout(make_route_config(c))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/logout", func(c *gin.Context) {
		route_data := route.View_login_logout_post(make_route_config(c))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/user", func(c *gin.Context) {
		route_data := route.View_user_safe(make_route_config(c), "")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/user/*user_name", func(c *gin.Context) {
		route_data := route.View_user_safe(make_route_config(c), strings.TrimPrefix(c.Param("user_name"), "/"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
}
