package main

import (
	"net/http"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_bbs_search_routes(r *gin.Engine) {
	r.GET("/bbs/search/:set_id", func(c *gin.Context) {
		route_data := route.View_bbs_search(make_route_config(c), c.Param("set_id"), "")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/bbs/search/:set_id", func(c *gin.Context) {
		route_data := route.View_bbs_search(make_route_config(c), c.Param("set_id"), c.PostForm("keyword"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
}
