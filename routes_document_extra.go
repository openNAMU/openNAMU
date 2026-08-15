package main

import (
	"net/http"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_document_extra_routes(r *gin.Engine) {
	r.GET("/edit_from/*doc_name", func(c *gin.Context) {
		data := route.View_edit(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Query("load"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/edit_from/*doc_name", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_edit_post(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.PostForm("content"), c.PostForm("send"), c.PostForm("copyright_agreement"), captcha_response(c), c.PostForm("ver"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/history_tool/:rev/*doc_name", func(c *gin.Context) {
		data := route.View_history_tool(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Param("rev"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/history_hidden/:rev/*doc_name", func(c *gin.Context) {
		data := route.View_history_hidden_safe(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Param("rev"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/history_send/:rev/*doc_name", func(c *gin.Context) {
		data := route.View_history_send_safe(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Param("rev"), nil)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/history_send/:rev/*doc_name", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_history_send_safe(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Param("rev"), c.Request.PostForm)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/history_delete/:rev/*doc_name", func(c *gin.Context) {
		data := route.View_history_delete(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Param("rev"), nil)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/history_delete/:rev/*doc_name", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_history_delete(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Param("rev"), c.Request.PostForm)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/history_reset/*doc_name", func(c *gin.Context) {
		data := route.View_history_reset(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), nil)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/history_reset/*doc_name", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_history_reset(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Request.PostForm)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/history_add/*doc_name", func(c *gin.Context) {
		data := route.View_history_add_safe(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), nil)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/history_add/*doc_name", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_history_add_safe(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Request.PostForm)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
}
