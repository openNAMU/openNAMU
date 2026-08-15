package main

import (
	"net/http"
	"net/url"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func user_compat_values(c *gin.Context) url.Values {
	_ = c.Request.ParseForm()
	return c.Request.PostForm
}

func register_user_compat_routes(r *gin.Engine) {
	r.GET("/edit_filter/*name", func(c *gin.Context) {
		data := route.View_user_edit_filter(make_route_config(c), strings.TrimPrefix(c.Param("name"), "/"), nil)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/edit_filter/*name", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_user_edit_filter(make_route_config(c), strings.TrimPrefix(c.Param("name"), "/"), c.Request.PostForm)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/change/key/delete", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_key_delete(make_route_config(c))))
	})
	r.POST("/change/key/delete", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_key_delete(make_route_config(c))))
	})
	r.GET("/change/email/delete", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_field_delete(make_route_config(c), "email")))
	})
	r.POST("/change/email/delete", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_field_delete(make_route_config(c), "email")))
	})
	r.GET("/change/email/check", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_email_check(make_route_config(c), nil)))
	})
	r.POST("/change/email/check", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_email_check(make_route_config(c), user_compat_values(c))))
	})
	r.GET("/change/head_reset", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_head_reset(make_route_config(c), nil)))
	})
	r.POST("/change/head_reset", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_head_reset(make_route_config(c), user_compat_values(c))))
	})
	r.GET("/change/skin_set/main", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_skin_main(make_route_config(c), nil)))
	})
	r.POST("/change/skin_set/main", func(c *gin.Context) {
		value := c.PostForm("main_css_darkmode")
		if value != "0" && value != "1" && value != "default" {
			value = "default"
		}
		http.SetCookie(c.Writer, &http.Cookie{Name: "main_css_darkmode", Value: value, Path: "/"})
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_skin_main(make_route_config(c), user_compat_values(c))))
	})
	r.GET("/skin_set", func(c *gin.Context) {
		set_skin_language_cookie(c)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_skin_set(make_route_config(c))))
	})
	r.POST("/skin_set", func(c *gin.Context) {
		set_skin_language_cookie(c)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_skin_set(make_route_config(c))))
	})
	r.GET("/change/head/:skin_name", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_head_skin(make_route_config(c), c.Param("skin_name"), nil)))
	})
	r.POST("/change/head/:skin_name", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_head_skin(make_route_config(c), c.Param("skin_name"), user_compat_values(c))))
	})
	r.GET("/change/user_name/:user_name", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_name_for(make_route_config(c), c.Param("user_name"), nil)))
	})
	r.POST("/change/user_name/:user_name", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_name_for(make_route_config(c), c.Param("user_name"), user_compat_values(c))))
	})
}
