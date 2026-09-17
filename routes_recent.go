package main

import (
	"net/http"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func Register_recent_routes(r *gin.Engine) {
	r.GET("/recent_block", func(c *gin.Context) {
		route_data := route.View_list_recent_block(Make_route_config(c), "1", "", "", "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/all", func(c *gin.Context) {
		route_data := route.View_list_recent_block(Make_route_config(c), "1", "", "", "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/all/:num", func(c *gin.Context) {
		route_data := route.View_list_recent_block(Make_route_config(c), c.Param("num"), "", "", "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/all/:num/*why", func(c *gin.Context) {
		route_data := route.View_list_recent_block(Make_route_config(c), c.Param("num"), "", strings.TrimPrefix(c.Param("why"), "/"), "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/user/:user_name", func(c *gin.Context) {
		route_data := route.View_list_recent_block(Make_route_config(c), "1", "user", "", c.Param("user_name"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/user/:user_name/:num", func(c *gin.Context) {
		route_data := route.View_list_recent_block(Make_route_config(c), c.Param("num"), "user", "", c.Param("user_name"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/admin/:user_name", func(c *gin.Context) {
		route_data := route.View_list_recent_block(Make_route_config(c), "1", "admin", "", c.Param("user_name"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/admin/:user_name/:num", func(c *gin.Context) {
		route_data := route.View_list_recent_block(Make_route_config(c), c.Param("num"), "admin", "", c.Param("user_name"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/regex", func(c *gin.Context) {
		route_data := route.View_list_recent_block(Make_route_config(c), "1", "regex", "", "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/regex/:num", func(c *gin.Context) {
		route_data := route.View_list_recent_block(Make_route_config(c), c.Param("num"), "regex", "", "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/cidr", func(c *gin.Context) {
		route_data := route.View_list_recent_block(Make_route_config(c), "1", "cidr", "", "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/cidr/:num", func(c *gin.Context) {
		route_data := route.View_list_recent_block(Make_route_config(c), c.Param("num"), "cidr", "", "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/private", func(c *gin.Context) {
		route_data := route.View_list_recent_block(Make_route_config(c), "1", "private", "", "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/private/:num", func(c *gin.Context) {
		route_data := route.View_list_recent_block(Make_route_config(c), c.Param("num"), "private", "", "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/ongoing", func(c *gin.Context) {
		route_data := route.View_list_recent_block(Make_route_config(c), "1", "ongoing", "", "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/ongoing/:num", func(c *gin.Context) {
		route_data := route.View_list_recent_block(Make_route_config(c), c.Param("num"), "ongoing", "", "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/login/:login_type", func(c *gin.Context) {
		route_data := route.View_list_recent_block(Make_route_config(c), "1", c.Param("login_type"), "", "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/login/:login_type/:num", func(c *gin.Context) {
		route_data := route.View_list_recent_block(Make_route_config(c), c.Param("num"), c.Param("login_type"), "", "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_change", func(c *gin.Context) {
		route_data := route.View_list_recent_change(Make_route_config(c), "", "50", "1")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_changes", func(c *gin.Context) {
		route_data := route.View_list_recent_change(Make_route_config(c), "", "50", "1")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_change/:num/:set_type", func(c *gin.Context) {
		route_data := route.View_list_recent_change(Make_route_config(c), c.Param("set_type"), "50", c.Param("num"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_discuss", func(c *gin.Context) {
		route_data := route.View_list_recent_discuss(Make_route_config(c), "50", "1", "normal")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_discuss/:num/:set_type", func(c *gin.Context) {
		route_data := route.View_list_recent_discuss(Make_route_config(c), "50", c.Param("num"), c.Param("set_type"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
}
