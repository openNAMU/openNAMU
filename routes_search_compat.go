package main

import (
	"net/http"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_search_compat_routes(r *gin.Engine) {
	r.GET("/goto/*keyword", func(c *gin.Context) {
		keyword := strings.TrimPrefix(c.Param("keyword"), "/")
		data := route.View_main_search_post(make_route_config(c), "", true, keyword)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/search_data/*keyword", func(c *gin.Context) {
		keyword := strings.TrimPrefix(c.Param("keyword"), "/")
		data := route.View_main_search(make_route_config(c), keyword, "1", "data")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
}
