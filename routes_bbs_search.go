package main

import (
	"net/http"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_bbs_search_routes(r *gin.Engine) {
	r.GET("/bbs/search", func(c *gin.Context) {
		route_data := route.View_bbs_search(make_route_config(c), "", "", "1")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/bbs/search", func(c *gin.Context) {
		route_data := route.View_bbs_search(make_route_config(c), "", c.PostForm("keyword"), "1")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/search/:set_id", func(c *gin.Context) {
		route_data := route.View_bbs_search(make_route_config(c), c.Param("set_id"), "", "1")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/bbs/search/:set_id", func(c *gin.Context) {
		route_data := route.View_bbs_search(make_route_config(c), c.Param("set_id"), c.PostForm("keyword"), "1")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/search_page/:page/*keyword", func(c *gin.Context) {
		route_data := route.View_bbs_search(
			make_route_config(c),
			"",
			strings.TrimPrefix(c.Param("keyword"), "/"),
			c.Param("page"),
		)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/search_board_page/:set_id/:page/*keyword", func(c *gin.Context) {
		route_data := route.View_bbs_search(
			make_route_config(c),
			c.Param("set_id"),
			strings.TrimPrefix(c.Param("keyword"), "/"),
			c.Param("page"),
		)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
}
