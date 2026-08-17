package main

import (
	"net/http"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_list_extra_routes(r *gin.Engine) {
	r.GET("/list/document/all", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_list_document_all(make_route_config(c), "1")))
	})
	r.GET("/list/document/all/:page", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_list_document_all(make_route_config(c), c.Param("page"))))
	})
	r.GET("/list/document/need", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_list_need_page(make_route_config(c), "1")))
	})
	r.GET("/list/document/no_link", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_list_no_link_page(make_route_config(c), "1")))
	})
	r.GET("/list/file", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_list_file_page(make_route_config(c), "1")))
	})
	r.GET("/list/image", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_list_image_page(make_route_config(c), "1")))
	})
	r.GET("/list/user", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_list_user_page(make_route_config(c), "1")))
	})
	r.GET("/list/admin", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_list_admin_page(make_route_config(c), "1")))
	})
	r.GET("/list/admin/auth_use", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_list_admin_auth_use_page(make_route_config(c), "1", "")))
	})
	r.GET("/list/user/check_submit/:user_name", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_list_user_check_submit(make_route_config(c), c.Param("user_name"))))
	})

	r.GET("/record", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_record_page(make_route_config(c), "", "edit", "1")))
	})
	r.GET("/record/:user_name", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_record_page(make_route_config(c), c.Param("user_name"), "edit", "1")))
	})
	r.GET("/record/topic/:user_name", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_record_page(make_route_config(c), c.Param("user_name"), "topic", "1")))
	})
	r.GET("/record/reset", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_record_simple(make_route_config(c), "", "edit")))
	})
	r.GET("/record/reset/:user_name", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_record_reset(make_route_config(c), c.Param("user_name"), nil)))
	})
	r.POST("/record/reset/:user_name", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_record_reset(make_route_config(c), c.Param("user_name"), c.Request.PostForm)))
	})
	r.GET("/count", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_record_count(make_route_config(c), "")))
	})
	r.GET("/count/:user_name", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_record_count(make_route_config(c), strings.TrimPrefix(c.Param("user_name"), "/"))))
	})
}
