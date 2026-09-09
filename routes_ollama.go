package main

import (
	"net/http"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_ollama_routes(r *gin.Engine) {
	r.GET("/ai", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_ollama(make_route_config(c), "", "")))
	})
	r.POST("/ai", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_ollama(make_route_config(c), c.PostForm("question"), c.PostForm("model"))))
	})
	r.POST("/api/ai/stream", func(c *gin.Context) {
		c.JSON(http.StatusOK, route.Api_ollama_stream_post(make_route_config(c), c.PostForm("question"), c.PostForm("model")))
	})
}
