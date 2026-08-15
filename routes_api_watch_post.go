package main

import (
	"net/http"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_api_watch_post_routes(r *gin.Engine) {
	r.POST("/api/v2/doc_star_doc/:num/*doc_name", func(c *gin.Context) {
		data := route.Api_w_watch_list_post(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "star_doc")
		c.JSON(http.StatusOK, data)
	})
	r.POST("/api/v2/doc_watch_list/:num/*doc_name", func(c *gin.Context) {
		data := route.Api_w_watch_list_post(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "watchlist")
		c.JSON(http.StatusOK, data)
	})
}
