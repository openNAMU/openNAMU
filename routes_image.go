package main

import (
	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_image_routes(r *gin.Engine) {
	r.GET("/thumbnail/:size/*name", route.View_image_thumbnail)
}
