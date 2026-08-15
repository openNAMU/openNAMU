package main

import (
	"net/http"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_record_page_compat_routes(r *gin.Engine) {
	r.GET("/record/topic/:user_name/:page", func(c *gin.Context) {
		data := route.View_record_page(make_route_config(c), c.Param("user_name"), "topic", c.Param("page"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/record/:user_name/:set_type/*record_user_name", func(c *gin.Context) {
		user_name := strings.TrimPrefix(c.Param("record_user_name"), "/")
		data := route.View_record_page(make_route_config(c), user_name, c.Param("set_type"), c.Param("user_name"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
}
