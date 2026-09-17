package main

import (
	"net/http"
	"strings"

	"opennamu/route"
	"opennamu/route/tool"

	"github.com/gin-gonic/gin"
)

func Register_search_routes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, route.Get_frontpage_url())
	})

	r.POST("/goto", func(c *gin.Context) {
		route_data := route.View_main_search_post(Make_route_config(c), "", true, c.PostForm("search"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/goto/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search_post(Make_route_config(c), "", true, c.PostForm("search"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/search", func(c *gin.Context) {
		route_data := route.View_main_search_post(Make_route_config(c), "", false, c.PostForm("search"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/search/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search(Make_route_config(c), strings.TrimPrefix(c.Param("keyword"), "/"), "1", "title")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/search/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search_post(Make_route_config(c), "", false, c.PostForm("search"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/search_page/:num/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search_post(Make_route_config(c), "", false, c.PostForm("search"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/search_page/:num/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search(Make_route_config(c), strings.TrimPrefix(c.Param("keyword"), "/"), c.Param("num"), "title")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/search_data/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search_post(Make_route_config(c), "", false, c.PostForm("search"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/search_data_page/:num/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search_post(Make_route_config(c), "", false, c.PostForm("search"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/search_data_page/:num/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search(Make_route_config(c), strings.TrimPrefix(c.Param("keyword"), "/"), c.Param("num"), "data")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/favicon.ico", func(c *gin.Context) {
		data, err := tool.Read_view_file("main_css/file/favicon.ico")
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		Write_data(c, http.StatusOK, "image/x-icon", data)
	})

	r.NoRoute(func(c *gin.Context) {
		if data, content_type, ok := route.Read_root_file(c.Request.URL.Path); ok {
			Write_data(c, http.StatusOK, content_type, data)
			return
		}
		route_data := route.View_main_404_page(Make_route_config(c), c.Request.URL.Path)
		Write_data(c, http.StatusNotFound, "text/html; charset=utf-8", []byte(route_data))
	})
}
