package main

import (
	"net/http"
	"strconv"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_vote_routes(r *gin.Engine) {
	r.GET("/vote/add", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_vote_add(make_route_config(c), nil)))
	})
	r.POST("/vote/add", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_vote_add(make_route_config(c), c.Request.PostForm)))
	})
	r.GET("/vote/list", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_vote_list(make_route_config(c), "open", "1")))
	})
	r.GET("/vote/list/:type", func(c *gin.Context) {
		page := "1"
		type_data := "open"
		if _, err := strconv.Atoi(c.Param("type")); err == nil {
			page = c.Param("type")
		} else if c.Param("type") == "close" {
			type_data = "close"
		}
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_vote_list(make_route_config(c), type_data, page)))
	})
	r.GET("/vote/list/:type/:page", func(c *gin.Context) {
		type_data := "open"
		if c.Param("type") == "close" {
			type_data = "close"
		}
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_vote_list(make_route_config(c), type_data, c.Param("page"))))
	})
	r.GET("/vote/:id", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_vote_select(make_route_config(c), c.Param("id"), nil)))
	})
	r.POST("/vote/:id", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_vote_select(make_route_config(c), c.Param("id"), c.Request.PostForm)))
	})
	r.GET("/vote/end/:id", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_vote_end(make_route_config(c), c.Param("id"))))
	})
	r.GET("/vote/close/:id", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_vote_close(make_route_config(c), c.Param("id"))))
	})
}
