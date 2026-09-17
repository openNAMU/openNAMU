package main

import (
	"net/http"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func Register_ollama_routes(r *gin.Engine) {
	r.GET("/ai", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_ollama(Make_route_config(c), "", "")))
	})
	r.POST("/ai", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_ollama(Make_route_config(c), c.PostForm("question"), c.PostForm("model"))))
	})
	r.POST("/api/ai/stream", func(c *gin.Context) {
		c.JSON(http.StatusOK, route.Api_ollama_stream_post(Make_route_config(c), c.PostForm("question"), c.PostForm("model")))
	})
}
