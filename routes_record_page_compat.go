package main

import (
	"net/http"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func Register_record_page_compat_routes(r *gin.Engine) {
	r.GET("/record/topic/:user_name/:page", func(c *gin.Context) {
		data := route.View_record_page(Make_route_config(c), c.Param("user_name"), "topic", c.Param("page"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/record/:user_name/:set_type/*record_user_name", func(c *gin.Context) {
		user_name := strings.TrimPrefix(c.Param("record_user_name"), "/")
		if c.Param("user_name") == "topic" {
			data := route.View_record_page(Make_route_config(c), c.Param("set_type"), "topic", user_name)
			Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
			return
		}
		data := route.View_record_page(Make_route_config(c), user_name, c.Param("set_type"), c.Param("user_name"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
}
