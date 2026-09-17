package main

import (
	"strings"

	"github.com/gin-gonic/gin"
	"opennamu/route"
)

func Compat_recent_discuss(c *gin.Context) {
	parts := strings.Split(strings.TrimPrefix(c.Param("data"), "/"), "/")
	limit := "10"
	set_type := "normal"
	if len(parts) == 1 {
		limit = parts[0]
	} else if len(parts) > 1 {
		set_type = parts[0]
		limit = parts[len(parts)-1]
	}
	Compat_api_data_cors(c, route.Api_list_recent_discuss(Make_route_config(c), limit, "1", set_type))
}
