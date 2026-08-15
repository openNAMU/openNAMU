package main

import (
	"net/http"
	"strings"

	"opennamu/route"
	"opennamu/route/tool"

	"github.com/gin-gonic/gin"
)

func register_search_routes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, route.Get_frontpage_url())
	})

	r.POST("/goto", func(c *gin.Context) {
		route_data := route.View_main_search_post(make_route_config(c), "", true, c.PostForm("search"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/goto/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search_post(make_route_config(c), "", true, c.PostForm("search"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/search", func(c *gin.Context) {
		route_data := route.View_main_search_post(make_route_config(c), "", false, c.PostForm("search"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/search/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search(make_route_config(c), strings.TrimPrefix(c.Param("keyword"), "/"), "1", "title")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/search/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search_post(make_route_config(c), "", false, c.PostForm("search"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/search_page/:num/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search_post(make_route_config(c), "", false, c.PostForm("search"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/search_page/:num/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search(make_route_config(c), strings.TrimPrefix(c.Param("keyword"), "/"), c.Param("num"), "title")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/search_data/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search_post(make_route_config(c), "", false, c.PostForm("search"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/search_data_page/:num/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search_post(make_route_config(c), "", false, c.PostForm("search"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/search_data_page/:num/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search(make_route_config(c), strings.TrimPrefix(c.Param("keyword"), "/"), c.Param("num"), "data")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/favicon.ico", func(c *gin.Context) {
		data, err := tool.Read_view_file("main_css/file/favicon.ico")
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		write_data(c, http.StatusOK, "image/x-icon", data)
	})

	r.NoRoute(func(c *gin.Context) {
		if data, content_type, ok := route.Read_root_file(c.Request.URL.Path); ok {
			write_data(c, http.StatusOK, content_type, data)
			return
		}
		route_data := route.View_main_404_page(make_route_config(c), c.Request.URL.Path)
		write_data(c, http.StatusNotFound, "text/html; charset=utf-8", []byte(route_data))
	})
}
