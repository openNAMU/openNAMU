package main

import (
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func Register_wiki_compat_routes(r *gin.Engine) {
	r.GET("/w_from/*doc_name", func(c *gin.Context) {
		data, status := route.View_w(c, Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "from")
		Write_data(c, status, "text/html; charset=utf-8", []byte(data))
	})
}
