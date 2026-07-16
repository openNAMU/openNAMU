package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opennamu/route"
)

func register_bbs_routes(r *gin.Engine) {
	r.GET("/vote", func(c *gin.Context) {
		route_data := route.View_vote_list(make_route_config(c), "open", "1")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/main", func(c *gin.Context) {
		route_data := route.View_bbs_main(make_route_config(c), "1")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/make", func(c *gin.Context) {
		route_data := route.View_bbs_make(make_route_config(c))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/bbs/make", func(c *gin.Context) {
		bbs_name := c.PostForm("bbs_name")
		bbs_type := c.PostForm("bbs_type")

		route_data := route.View_bbs_make_post(make_route_config(c), bbs_name, bbs_type)
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/in/:set_id", func(c *gin.Context) {
		route_data := route.View_bbs_in(make_route_config(c), c.Param("set_id"), "1")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/in/:set_id/:page_num", func(c *gin.Context) {
		route_data := route.View_bbs_in(make_route_config(c), c.Param("set_id"), c.Param("page_num"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/w/:set_id/:set_code", func(c *gin.Context) {
		route_data := route.View_bbs_in_w(c, make_route_config(c), c.Param("set_id"), c.Param("set_code"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
}
