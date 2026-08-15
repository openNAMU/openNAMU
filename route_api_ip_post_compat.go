package main

import (
	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func compat_ip_single_post(c *gin.Context) {
	compat_api_data(c, route.Api_func_ip_post(make_route_config(c), compat_ip_post_data(c)))
}

func compat_ip_menu_post(c *gin.Context) {
	compat_api_data(c, route.Api_func_ip_menu(make_route_config(c), compat_doc_name(c, "ip"), "user"))
}
