package main

import (
	"net/http"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_wiki_routes(r *gin.Engine) {
	r.GET("/xref/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_xref(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "1", "1")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/xref_page/:num/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_xref(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "1", c.Param("num"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/xref_this/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_xref(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "2", "1")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/xref_this_page/:num/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_xref(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "2", c.Param("num"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/xref_reset/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_xref_reset(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/xref_reset/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_xref_reset_post(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/upload", func(c *gin.Context) {
		route_data := route.View_edit_file_upload(make_route_config(c), c.Query("name"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/w/*doc_name", func(c *gin.Context) {
		route_data, status_code := route.View_w(c, make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "")
		write_data(c, status_code, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/down/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_down(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/raw/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_raw(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "", "")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/raw_rev/:rev/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_raw(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Param("rev"), "")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/raw_acl/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_raw(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "", "document_acl")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/topic/*doc_name", func(c *gin.Context) {
		route_data := route.View_topic_list(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Query("tool"), "1")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
	r.POST("/topic/*doc_name", func(c *gin.Context) {
		route_data := route.View_topic_list(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Query("tool"), "1")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/topic_page/:num/*doc_name", func(c *gin.Context) {
		route_data := route.View_topic_list(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "", c.Param("num"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
	r.POST("/topic_page/:num/*doc_name", func(c *gin.Context) {
		route_data := route.View_topic_list(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "", c.Param("num"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/topic_close/:num/*doc_name", func(c *gin.Context) {
		route_data := route.View_topic_list(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "close", c.Param("num"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/topic_agree/:num/*doc_name", func(c *gin.Context) {
		route_data := route.View_topic_list(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "agree", c.Param("num"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
}
