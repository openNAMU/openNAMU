package main

import (
	"net/http"
	"net/url"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_admin_misc_routes(r *gin.Engine) {
	r.GET("/list/user/check", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_check(make_route_config(c), "", "normal", "1", "")))
	})
	r.POST("/list/user/check", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_check(make_route_config(c), c.PostForm("name"), "normal", "1", "")))
	})
	r.GET("/list/user/check/:name", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_check(make_route_config(c), c.Param("name"), "normal", "1", "")))
	})
	r.GET("/list/user/check/:name/:check_type", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_check(make_route_config(c), c.Param("name"), c.Param("check_type"), "1", "")))
	})
	r.GET("/list/user/check/:name/:check_type/:page", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_check(make_route_config(c), c.Param("name"), c.Param("check_type"), c.Param("page"), "")))
	})
	r.GET("/list/user/check/:name/:check_type/:page/:plus_name", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_check(make_route_config(c), c.Param("name"), c.Param("check_type"), c.Param("page"), c.Param("plus_name"))))
	})

	r.GET("/list/user/check/delete/*data", func(c *gin.Context) {
		admin_user_check_delete(c, nil)
	})
	r.POST("/list/user/check/delete/*data", func(c *gin.Context) {
		admin_user_check_delete(c, admin_post_values(c))
	})

	r.GET("/delete_multiple", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_delete_multiple(make_route_config(c), nil)))
	})
	r.POST("/delete_multiple", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_delete_multiple(make_route_config(c), admin_post_values(c))))
	})
}

func admin_user_check_delete(c *gin.Context, form_values url.Values) {
	path := strings.TrimPrefix(c.Param("data"), "/")
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		c.Redirect(http.StatusFound, "/manager")
		return
	}
	values := []string{}
	for _, part := range parts[:4] {
		value, err := url.PathUnescape(part)
		if err != nil {
			value = part
		}
		values = append(values, value)
	}
	write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_check_delete(make_route_config(c), values[0], values[1], values[2], values[3], form_values)))
}
