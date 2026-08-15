package main

import (
	"strings"

	"github.com/gin-gonic/gin"
	"opennamu/route"
)

func compat_recent_discuss(c *gin.Context) {
	parts := strings.Split(strings.TrimPrefix(c.Param("data"), "/"), "/")
	limit := "10"
	set_type := "normal"
	if len(parts) == 1 {
		limit = parts[0]
	} else if len(parts) > 1 {
		set_type = parts[0]
		limit = parts[len(parts)-1]
	}
	compat_api_data_cors(c, route.Api_list_recent_discuss(make_route_config(c), limit, "1", set_type))
}
