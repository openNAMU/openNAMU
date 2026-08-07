package main

import (
	"net/http"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_recent_routes(r *gin.Engine) {
	r.GET("/recent_block", func(c *gin.Context) {
		route_data := route.View_list_recent_block(make_route_config(c), "1", "", "", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/all", func(c *gin.Context) {
		route_data := route.View_list_recent_block(make_route_config(c), "1", "", "", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/all/:num", func(c *gin.Context) {
		route_data := route.View_list_recent_block(make_route_config(c), c.Param("num"), "", "", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/all/:num/*why", func(c *gin.Context) {
		route_data := route.View_list_recent_block(make_route_config(c), c.Param("num"), "", strings.TrimPrefix(c.Param("why"), "/"), "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/user/:user_name", func(c *gin.Context) {
		route_data := route.View_list_recent_block(make_route_config(c), "1", "user", "", c.Param("user_name"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/user/:user_name/:num", func(c *gin.Context) {
		route_data := route.View_list_recent_block(make_route_config(c), c.Param("num"), "user", "", c.Param("user_name"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/admin/:user_name", func(c *gin.Context) {
		route_data := route.View_list_recent_block(make_route_config(c), "1", "admin", "", c.Param("user_name"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/admin/:user_name/:num", func(c *gin.Context) {
		route_data := route.View_list_recent_block(make_route_config(c), c.Param("num"), "admin", "", c.Param("user_name"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/regex", func(c *gin.Context) {
		route_data := route.View_list_recent_block(make_route_config(c), "1", "regex", "", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/regex/:num", func(c *gin.Context) {
		route_data := route.View_list_recent_block(make_route_config(c), c.Param("num"), "regex", "", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/cidr", func(c *gin.Context) {
		route_data := route.View_list_recent_block(make_route_config(c), "1", "cidr", "", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/cidr/:num", func(c *gin.Context) {
		route_data := route.View_list_recent_block(make_route_config(c), c.Param("num"), "cidr", "", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/private", func(c *gin.Context) {
		route_data := route.View_list_recent_block(make_route_config(c), "1", "private", "", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/private/:num", func(c *gin.Context) {
		route_data := route.View_list_recent_block(make_route_config(c), c.Param("num"), "private", "", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/ongoing", func(c *gin.Context) {
		route_data := route.View_list_recent_block(make_route_config(c), "1", "ongoing", "", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/ongoing/:num", func(c *gin.Context) {
		route_data := route.View_list_recent_block(make_route_config(c), c.Param("num"), "ongoing", "", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_change", func(c *gin.Context) {
		route_data := route.View_list_recent_change(make_route_config(c), "", "50", "1")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_changes", func(c *gin.Context) {
		route_data := route.View_list_recent_change(make_route_config(c), "", "50", "1")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_change/:num/:set_type", func(c *gin.Context) {
		route_data := route.View_list_recent_change(make_route_config(c), c.Param("set_type"), "50", c.Param("num"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_discuss", func(c *gin.Context) {
		route_data := route.View_list_recent_discuss(make_route_config(c), "50", "1", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_discuss/:num/:set_type", func(c *gin.Context) {
		route_data := route.View_list_recent_discuss(make_route_config(c), "50", c.Param("num"), c.Param("set_type"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
}
