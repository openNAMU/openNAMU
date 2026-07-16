package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opennamu/route"
	"strings"
)

func register_wiki_routes(r *gin.Engine) {
	r.GET("/upload", func(c *gin.Context) {
		route_data := route.View_edit_file_upload(make_route_config(c))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/w/*doc_name", func(c *gin.Context) {
		route_data, status_code := route.View_w(c, make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"))
		c.Data(status_code, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/down/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_down(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/raw/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_raw(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/raw_rev/:rev/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_raw(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Param("rev"), "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/raw_acl/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_raw(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "", "document_acl")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/topic/*doc_name", func(c *gin.Context) {
		route_data := route.View_topic_list(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "", "1")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/topic_page/:num/*doc_name", func(c *gin.Context) {
		route_data := route.View_topic_list(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "", c.Param("num"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/topic_close/:num/*doc_name", func(c *gin.Context) {
		route_data := route.View_topic_list(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "close", c.Param("num"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/topic_agree/:num/*doc_name", func(c *gin.Context) {
		route_data := route.View_topic_list(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "agree", c.Param("num"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
}
