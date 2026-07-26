package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opennamu/route"
	"strings"
)

func register_api_routes(r *gin.Engine) {
	r.GET("/api/version", func(c *gin.Context) {
		route_data := route.Api_version(make_route_config(c))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/user_info/:user_name", func(c *gin.Context) {
		route_data := route.Api_user_info(make_route_config(c), c.Param("user_name"))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/recent_change", func(c *gin.Context) {
		route_data := route.Api_list_recent_change(make_route_config(c), "normal", "10", "1")
		in_data := route_data["data"].([][]string)
		c.JSON(http.StatusOK, in_data)
	})

	r.GET("/api/recent_change/:limit", func(c *gin.Context) {
		route_data := route.Api_list_recent_change(make_route_config(c), "normal", c.Param("limit"), "1")
		in_data := route_data["data"].([][]string)
		c.JSON(http.StatusOK, in_data)
	})

	r.GET("/api/recent_discuss", func(c *gin.Context) {
		route_data := route.Api_list_recent_discuss(make_route_config(c), "10", "1", "")
		in_data := route_data["data"].([][]string)
		c.JSON(http.StatusOK, in_data)
	})

	r.GET("/api/recent_discuss/:limit", func(c *gin.Context) {
		route_data := route.Api_list_recent_discuss(make_route_config(c), c.Param("limit"), "1", "")
		in_data := route_data["data"].([][]string)
		c.JSON(http.StatusOK, in_data)
	})

	r.POST("/api/v2/lang", func(c *gin.Context) {
		data := c.PostForm("data")
		safe := c.PostForm("safe")
		legacy := c.PostForm("legacy")

		route_data := route.Api_func_language(make_route_config(c), data, safe, legacy)
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/ip/:ip", func(c *gin.Context) {
		route_data := route.Api_func_ip(make_route_config(c), c.Param("ip"))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/ip_menu/:ip", func(c *gin.Context) {
		route_data := route.Api_func_ip_menu(make_route_config(c), c.Param("ip"), "")
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/user/setting/editor", func(c *gin.Context) {
		route_data := route.Api_user_setting_editor(make_route_config(c))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/page_view/*doc_name", func(c *gin.Context) {
		route_data := route.Api_w_page_view(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"))
		c.JSON(http.StatusOK, route_data)
	})

	r.POST("/api/v2/user/setting/editor", func(c *gin.Context) {
		route_data := route.Api_user_setting_editor_post(make_route_config(c), c.Request.FormValue("data"))
		c.JSON(http.StatusOK, route_data)
	})

	r.DELETE("/api/v2/user/setting/editor", func(c *gin.Context) {
		route_data := route.Api_user_setting_editor_delete(make_route_config(c), c.Request.FormValue("data"))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/page_view_post/*doc_name", func(c *gin.Context) {
		route_data := route.Api_w_page_view_post(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/bbs/w/page_view_post/:set_id/:set_code", func(c *gin.Context) {
		route_data := route.Api_bbs_w_page_view_post(make_route_config(c), strings.TrimPrefix(c.Param("set_id"), "/"), strings.TrimPrefix(c.Param("set_code"), "/"))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/bbs/w/page_view/:set_id/:set_code", func(c *gin.Context) {
		route_data := route.Api_bbs_w_page_view(make_route_config(c), strings.TrimPrefix(c.Param("set_id"), "/"), strings.TrimPrefix(c.Param("set_code"), "/"))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/bbs/main", func(c *gin.Context) {
		route_data := route.Api_bbs(make_route_config(c), "", "1")
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/raw/*doc_name", func(c *gin.Context) {
		route_data := route.Api_w_raw(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "", "")
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/raw_exist/*doc_name", func(c *gin.Context) {
		route_data := route.Api_w_raw(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "true", "")
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/raw_rev/:rev/*doc_name", func(c *gin.Context) {
		route_data := route.Api_w_raw(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "", c.Param("rev"))
		c.JSON(http.StatusOK, route_data)
	})
}
