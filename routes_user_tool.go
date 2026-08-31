package main

import (
	"net/http"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_user_tool_routes(r *gin.Engine) {
	r.GET("/user_tool/*user_name", func(c *gin.Context) {
		user_name := strings.TrimPrefix(c.Param("user_name"), "/")
		data := route.View_user_tool(make_route_config(c), user_name)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
}
