package main

import (
	"net/http"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func Register_include_routes(r *gin.Engine) {
	r.GET("/include/:payload", func(c *gin.Context) {
		data := route.View_include(Make_route_config(c), c.Param("payload"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
}
