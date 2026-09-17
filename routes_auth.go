package main

import (
	"net/http"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func Register_auth_routes(r *gin.Engine) {
	r.GET("/login", func(c *gin.Context) {
		route_data := route.View_login_login(Make_route_config(c))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/register", func(c *gin.Context) {
		route_data := route.View_login_register(Make_route_config(c))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
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
			Make_route_config(c),
			c.PostForm("id"),
			password,
			password_check,
			Captcha_response_internal(c),
			c.PostForm("invite"),
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/login", func(c *gin.Context) {
		password := c.PostForm("password")
		if password == "" {
			password = c.PostForm("pw")
		}
		route_data := route.View_login_login_post_full(
			Make_route_config(c),
			c.PostForm("id"),
			password,
			Captcha_response_internal(c),
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/logout", func(c *gin.Context) {
		route_data := route.View_login_logout(Make_route_config(c))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/logout", func(c *gin.Context) {
		route_data := route.View_login_logout_post(Make_route_config(c))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/user", func(c *gin.Context) {
		route_data := route.View_user_safe(Make_route_config(c), "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/user/*user_name", func(c *gin.Context) {
		route_data := route.View_user_safe(Make_route_config(c), strings.TrimPrefix(c.Param("user_name"), "/"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
}
