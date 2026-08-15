package main

import (
	"net/http"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_list_compat_routes(r *gin.Engine) {
	r.GET("/list/document/acl", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_list_document_acl(make_route_config(c), "1")))
	})
	r.GET("/list/document/acl/:page", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_list_document_acl(make_route_config(c), c.Param("page"))))
	})
	r.GET("/list/document/need/:page", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_list_need_page(make_route_config(c), c.Param("page"))))
	})
	r.GET("/list/document/no_link/:page", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_list_no_link_page(make_route_config(c), c.Param("page"))))
	})
	r.GET("/list/file/:page", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_list_file_page(make_route_config(c), c.Param("page"), false)))
	})
	r.GET("/list/image/:page", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_list_file_page(make_route_config(c), c.Param("page"), true)))
	})
	r.POST("/list/admin/auth_use", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_list_admin_page(make_route_config(c), "1", true, c.PostForm("search"))))
	})
	r.GET("/list/admin/auth_use_page/:page", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_list_admin_page(make_route_config(c), c.Param("page"), true, "")))
	})
	r.POST("/list/admin/auth_use_page/:page", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_list_admin_page(make_route_config(c), c.Param("page"), true, c.PostForm("search"))))
	})
	r.GET("/list/admin/auth_use_page/:page/*search", func(c *gin.Context) {
		search := strings.TrimPrefix(c.Param("search"), "/")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_list_admin_page(make_route_config(c), c.Param("page"), true, search)))
	})
	r.POST("/list/admin/auth_use_page/:page/*search", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_list_admin_page(make_route_config(c), c.Param("page"), true, c.PostForm("search"))))
	})
	r.GET("/list/admin/:page", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_list_admin_page(make_route_config(c), c.Param("page"), false, "")))
	})
	r.GET("/list/user/:page", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_list_user_page(make_route_config(c), c.Param("page"))))
	})
}
