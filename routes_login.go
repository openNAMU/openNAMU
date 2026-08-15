package main

import (
	"net/http"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_login_routes(r *gin.Engine) {
	r.GET("/login/2fa", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_login_2fa(make_route_config(c), nil)))
	})
	r.POST("/login/2fa", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_login_2fa(make_route_config(c), c.Request.PostForm)))
	})
	r.GET("/login/2fa/email", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_login_2fa_email(make_route_config(c), nil)))
	})
	r.POST("/login/2fa/email", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_login_2fa_email(make_route_config(c), c.Request.PostForm)))
	})
	r.GET("/login/find", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_login_find(make_route_config(c), nil)))
	})
	r.POST("/login/find", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_login_find(make_route_config(c), c.Request.PostForm)))
	})
	r.GET("/login/find/key", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_login_find_key(make_route_config(c), nil)))
	})
	r.POST("/login/find/key", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_login_find_key(make_route_config(c), c.Request.PostForm)))
	})
	r.GET("/login/find/email", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_login_find_email(make_route_config(c), nil)))
	})
	r.POST("/login/find/email", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_login_find_email(make_route_config(c), c.Request.PostForm)))
	})
	r.GET("/login/find/email/check", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_login_find_email_check(make_route_config(c), nil)))
	})
	r.POST("/login/find/email/check", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_login_find_email_check(make_route_config(c), c.Request.PostForm)))
	})
}
