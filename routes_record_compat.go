package main

import (
	"net/http"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_record_compat_routes(r *gin.Engine) {
	r.GET("/record/bbs/:user_name", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_record_bbs_safe(make_route_config(c), c.Param("user_name"), "1")))
	})
	r.GET("/record/bbs/:user_name/:page", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_record_bbs_safe(make_route_config(c), c.Param("user_name"), c.Param("page"))))
	})
	r.GET("/record/bbs_comment/:user_name", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_record_bbs_comment_safe(make_route_config(c), c.Param("user_name"), "1")))
	})
	r.GET("/record/bbs_comment/:user_name/:page", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_record_bbs_comment_safe(make_route_config(c), c.Param("user_name"), c.Param("page"))))
	})
}
