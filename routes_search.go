package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opennamu/route"
	"strings"
)

func register_search_routes(r *gin.Engine) {
	r.POST("/goto", func(c *gin.Context) {
		route_data := route.View_main_search_post(make_route_config(c), "", true, c.PostForm("search"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/goto/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search_post(make_route_config(c), "", true, c.PostForm("search"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/search", func(c *gin.Context) {
		route_data := route.View_main_search_post(make_route_config(c), "", false, c.PostForm("search"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/search/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search(make_route_config(c), strings.TrimPrefix(c.Param("keyword"), "/"), "1", "title")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/search/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search_post(make_route_config(c), "", false, c.PostForm("search"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/search_page/:num/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search_post(make_route_config(c), "", false, c.PostForm("search"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/search_page/:num/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search(make_route_config(c), strings.TrimPrefix(c.Param("keyword"), "/"), c.Param("num"), "title")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/search_data/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search_post(make_route_config(c), "", false, c.PostForm("search"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/search_data_page/:num/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search_post(make_route_config(c), "", false, c.PostForm("search"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/search_data_page/:num/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search(make_route_config(c), strings.TrimPrefix(c.Param("keyword"), "/"), c.Param("num"), "data")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.NoRoute(func(c *gin.Context) {
		route_data := route.View_main_404_page(make_route_config(c), c.Request.URL.Path)
		c.Data(http.StatusNotFound, "text/html; charset=utf-8", []byte(route_data))
	})
}
