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

	r.GET("/ext_backlink/*host", func(c *gin.Context) {
		host_param := strings.TrimPrefix(c.Param("host"), "/")
		route_data := route.View_external_docs_for_host(Make_route_config(c), host_param, "1")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/ext_backlink_page/:num/*host", func(c *gin.Context) {
		host_param := strings.TrimPrefix(c.Param("host"), "/")
		route_data := route.View_external_docs_for_host(Make_route_config(c), host_param, c.Param("num"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/ext_backlink_urls/*doc_name", func(c *gin.Context) {
		doc_name := strings.TrimPrefix(c.Param("doc_name"), "/")
		route_data := route.View_external_hosts_of_doc(Make_route_config(c), doc_name)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/ext_backlink_count/*host", func(c *gin.Context) {
		host_param := strings.TrimPrefix(c.Param("host"), "/")
		route_data := route.View_external_count_for_host(Make_route_config(c), host_param)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/ext_backlink_hosts", func(c *gin.Context) {
		route_data := route.View_external_all_hosts(Make_route_config(c), "1")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
	r.GET("/ext_backlink_hosts/:page", func(c *gin.Context) {
		route_data := route.View_external_all_hosts(Make_route_config(c), c.Param("page"))
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
