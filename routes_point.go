package main

import (
	"net/http"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_point_routes(r *gin.Engine) {
	r.GET("/point/give", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_point_give(make_route_config(c), nil)))
	})
	r.POST("/point/give", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_point_give(make_route_config(c), admin_post_values(c))))
	})
	r.POST("/api/point/give", func(c *gin.Context) {
		c.JSON(http.StatusOK, route.Api_point_give_post(make_route_config(c), c.PostForm("user_name"), c.PostForm("amount")))
	})
}
