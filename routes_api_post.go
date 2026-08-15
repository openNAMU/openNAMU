package main

import "github.com/gin-gonic/gin"

func register_api_post_routes(r *gin.Engine) {
	r.POST("/api/v2/ip", compat_ip_post)
	r.POST("/api/v2/ip/*ip", compat_ip_single_post)
	r.POST("/api/v2/ip_menu/*ip", compat_ip_menu_post)
}
