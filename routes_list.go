package main

import (
	"net/http"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_list_routes(r *gin.Engine) {
	r.GET("/list/random", func(c *gin.Context) {
		route_data := route.View_list_random(make_route_config(c))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/list/document/old", func(c *gin.Context) {
		route_data := route.View_list_old_page(make_route_config(c), "1", "old")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/list/document/old/:num", func(c *gin.Context) {
		route_data := route.View_list_old_page(make_route_config(c), c.Param("num"), "old")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/list/document/new", func(c *gin.Context) {
		route_data := route.View_list_old_page(make_route_config(c), "1", "new")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/list/document/new/:num", func(c *gin.Context) {
		route_data := route.View_list_old_page(make_route_config(c), c.Param("num"), "new")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/list/document/long", func(c *gin.Context) {
		route_data := route.View_list_long_page(make_route_config(c), "1", "long")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/list/document/long/:num", func(c *gin.Context) {
		route_data := route.View_list_long_page(make_route_config(c), c.Param("num"), "long")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/list/document/short", func(c *gin.Context) {
		route_data := route.View_list_long_page(make_route_config(c), "1", "short")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/list/document/short/:num", func(c *gin.Context) {
		route_data := route.View_list_long_page(make_route_config(c), c.Param("num"), "short")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/random", func(c *gin.Context) {
		route_data := route.View_w_random(make_route_config(c))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/other", func(c *gin.Context) {
		route_data := route.View_main_other(make_route_config(c))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/easter_egg", func(c *gin.Context) {
		route_data := route.View_easter_egg(make_route_config(c))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
}
