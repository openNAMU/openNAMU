package main

import "github.com/gin-gonic/gin"

func Register_api_post_routes(r *gin.Engine) {
	r.POST("/api/v2/ip", Compat_ip_post)
	r.POST("/api/v2/ip/*ip", Compat_ip_single_post)
	r.POST("/api/v2/ip_menu/*ip", Compat_ip_menu_post)
}
