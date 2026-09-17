package main

import (
	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func Compat_ip_single_post(c *gin.Context) {
	Compat_api_data(c, route.Api_func_ip_post(Make_route_config(c), Compat_ip_post_data(c)))
}

func Compat_ip_menu_post(c *gin.Context) {
	Compat_api_data(c, route.Api_func_ip_menu(Make_route_config(c), Compat_doc_name(c, "ip"), "user"))
}
