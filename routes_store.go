package main

import (
	"net/http"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func Register_store_routes(r *gin.Engine) {
	r.GET("/store", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_store(Make_route_config(c))))
	})
	r.GET("/api/store", func(c *gin.Context) {
		c.JSON(http.StatusOK, route.Api_store_list(Make_route_config(c)))
	})
}
