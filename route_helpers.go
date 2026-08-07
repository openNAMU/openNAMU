package main

import (
	"opennamu/route/tool"

	"github.com/gin-gonic/gin"
)

func make_route_config(c *gin.Context) tool.Config {
	return tool.Config{
		IP:      tool.Get_IP(c),
		Cookies: tool.Get_Cookies(c),
		Session: tool.Get_session(c),
	}
}
