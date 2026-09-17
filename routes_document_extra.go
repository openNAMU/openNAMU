package main

import (
	"net/http"
	"strings"

	"opennamu/route"
	"opennamu/route/tool"

	"github.com/gin-gonic/gin"
)

func Register_document_extra_routes(r *gin.Engine) {
	r.GET("/edit_from/*doc_name", func(c *gin.Context) {
		data := route.View_edit(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/edit_from_load/:load/*doc_name", func(c *gin.Context) {
		load_data, err := tool.Get_base64_decode(c.Param("load"))
		if err != nil {
			c.Redirect(http.StatusFound, "/manager")
			return
		}
		data := route.View_edit(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), load_data)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/edit_from/*doc_name", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_edit_post(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.PostForm("content"), c.PostForm("send"), c.PostForm("copyright_agreement"), Captcha_response_internal(c), c.PostForm("ver"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/history_tool/:rev/*doc_name", func(c *gin.Context) {
		data := route.View_history_tool(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Param("rev"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/history_hidden/:rev/*doc_name", func(c *gin.Context) {
		data := route.View_history_hidden_safe(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Param("rev"), nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/history_hidden/:rev/*doc_name", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_history_hidden_safe(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Param("rev"), c.Request.PostForm)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/history_send/:rev/*doc_name", func(c *gin.Context) {
		data := route.View_history_send_safe(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Param("rev"), nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/history_send/:rev/*doc_name", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_history_send_safe(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Param("rev"), c.Request.PostForm)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/history_delete/:rev/*doc_name", func(c *gin.Context) {
		data := route.View_history_delete(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Param("rev"), nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/history_delete/:rev/*doc_name", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_history_delete(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Param("rev"), c.Request.PostForm)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/history_reset/*doc_name", func(c *gin.Context) {
		data := route.View_history_reset(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/history_reset/*doc_name", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_history_reset(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Request.PostForm)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/history_add/*doc_name", func(c *gin.Context) {
		data := route.View_history_add_safe(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/history_add/*doc_name", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_history_add_safe(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Request.PostForm)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
}
