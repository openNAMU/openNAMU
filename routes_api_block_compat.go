package main

import (
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_api_block_compat_routes(r *gin.Engine) {
	r.GET("/api/v2/recent_block/:set_type/:num/*why", func(c *gin.Context) {
		why := strings.TrimPrefix(c.Param("why"), "/")
		compat_api_data(c, route.Api_list_recent_block(make_route_config(c), c.Param("num"), c.Param("set_type"), why, ""))
	})
	r.GET("/api/v2/recent_block_user/:set_type/:num/:user_name/*why", func(c *gin.Context) {
		why := strings.TrimPrefix(c.Param("why"), "/")
		compat_api_data(c, route.Api_list_recent_block(make_route_config(c), c.Param("num"), c.Param("set_type"), why, c.Param("user_name")))
	})
}
