package main

import (
	"net/http"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func Register_wiki_routes(r *gin.Engine) {
	r.GET("/xref/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_xref(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "1", "1")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/xref_page/:num/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_xref(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "1", c.Param("num"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/xref_this/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_xref(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "2", "1")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/xref_this_page/:num/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_xref(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "2", c.Param("num"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/xref_reset/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_xref_reset(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/xref_reset/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_xref_reset_post(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/upload", func(c *gin.Context) {
		route_data := route.View_edit_file_upload(Make_route_config(c), "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
	r.GET("/upload/*name", func(c *gin.Context) {
		file_name := strings.TrimPrefix(c.Param("name"), "/")
		route_data := route.View_edit_file_upload(Make_route_config(c), file_name)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/w/*doc_name", func(c *gin.Context) {
		route_data, status_code := route.View_w(c, Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "")
		Write_data(c, status_code, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/down/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_down(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/raw/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_raw(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "", "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/raw_rev/:rev/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_raw(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Param("rev"), "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/raw_acl/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_raw(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "", "document_acl")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/topic/*doc_name", func(c *gin.Context) {
		route_data := route.View_topic_list(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "", "1")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
	r.POST("/topic/*doc_name", func(c *gin.Context) {
		route_data := route.View_topic_list(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "", "1")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/topic_page/:num/*doc_name", func(c *gin.Context) {
		route_data := route.View_topic_list(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "", c.Param("num"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
	r.POST("/topic_page/:num/*doc_name", func(c *gin.Context) {
		route_data := route.View_topic_list(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "", c.Param("num"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/topic_close/:num/*doc_name", func(c *gin.Context) {
		route_data := route.View_topic_list(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "close", c.Param("num"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/topic_agree/:num/*doc_name", func(c *gin.Context) {
		route_data := route.View_topic_list(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "agree", c.Param("num"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
}
