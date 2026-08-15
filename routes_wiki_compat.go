package main

import (
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_wiki_compat_routes(r *gin.Engine) {
	r.GET("/w_from/*doc_name", func(c *gin.Context) {
		data, status := route.View_w(c, make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "from")
		write_data(c, status, "text/html; charset=utf-8", []byte(data))
	})
}
