package main

import (
	"net/http"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func Register_list_routes(r *gin.Engine) {
	r.GET("/list/random", func(c *gin.Context) {
		route_data := route.View_list_random(Make_route_config(c))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/list/document/old", func(c *gin.Context) {
		route_data := route.View_list_old_page(Make_route_config(c), "1", "old")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/list/document/old/:num", func(c *gin.Context) {
		route_data := route.View_list_old_page(Make_route_config(c), c.Param("num"), "old")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/list/document/new", func(c *gin.Context) {
		route_data := route.View_list_old_page(Make_route_config(c), "1", "new")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/list/document/new/:num", func(c *gin.Context) {
		route_data := route.View_list_old_page(Make_route_config(c), c.Param("num"), "new")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/list/document/long", func(c *gin.Context) {
		route_data := route.View_list_long_page(Make_route_config(c), "1", "long")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/list/document/long/:num", func(c *gin.Context) {
		route_data := route.View_list_long_page(Make_route_config(c), c.Param("num"), "long")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/list/document/short", func(c *gin.Context) {
		route_data := route.View_list_long_page(Make_route_config(c), "1", "short")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/list/document/short/:num", func(c *gin.Context) {
		route_data := route.View_list_long_page(Make_route_config(c), c.Param("num"), "short")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/random", func(c *gin.Context) {
		route_data := route.View_w_random(Make_route_config(c))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/random/category/*category_name", func(c *gin.Context) {
		category_name := strings.TrimPrefix(c.Param("category_name"), "/")
		route_data := route.View_w_random_category(Make_route_config(c), category_name)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/other", func(c *gin.Context) {
		route_data := route.View_main_other(Make_route_config(c))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/easter_egg", func(c *gin.Context) {
		route_data := route.View_easter_egg(Make_route_config(c), nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
	r.POST("/easter_egg", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		route_data := route.View_easter_egg(Make_route_config(c), c.Request.PostForm)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
}
