package main

import (
	"net/http"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_register_routes(r *gin.Engine) {
	r.GET("/register/email", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_register_email(make_route_config(c), nil)))
	})
	r.POST("/register/email", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_register_email(make_route_config(c), c.Request.PostForm)))
	})
	r.GET("/register/email/check", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_register_email_check(make_route_config(c), nil)))
	})
	r.POST("/register/email/check", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_register_email_check(make_route_config(c), c.Request.PostForm)))
	})
	r.GET("/register/submit", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_register_submit(make_route_config(c), nil)))
	})
	r.POST("/register/submit", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_register_submit(make_route_config(c), c.Request.PostForm)))
	})
}
