package main

import (
	"net/http"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_thread_routes(r *gin.Engine) {
	r.GET("/thread/0/*doc_name", func(c *gin.Context) {
		data := route.View_thread_route(make_route_config(c), "0", strings.TrimPrefix(c.Param("doc_name"), "/"), nil)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/thread/0/*doc_name", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_thread_route(make_route_config(c), "0", strings.TrimPrefix(c.Param("doc_name"), "/"), c.Request.PostForm)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/thread/:topic_num/comment/:num/raw", func(c *gin.Context) {
		data := route.View_thread_raw(make_route_config(c), c.Param("topic_num"), c.Param("num"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/thread/:topic_num/comment/:num/tool", func(c *gin.Context) {
		data := route.View_thread_comment_tool(make_route_config(c), c.Param("topic_num"), c.Param("num"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/thread/:topic_num/comment/:num/notice", func(c *gin.Context) {
		data := route.View_thread_comment_notice(make_route_config(c), c.Param("topic_num"), c.Param("num"), nil)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/thread/:topic_num/comment/:num/notice", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_thread_comment_notice(make_route_config(c), c.Param("topic_num"), c.Param("num"), c.Request.PostForm)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/thread/:topic_num/comment/:num/blind", func(c *gin.Context) {
		data := route.View_thread_comment_blind(make_route_config(c), c.Param("topic_num"), c.Param("num"), nil)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/thread/:topic_num/comment/:num/blind", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_thread_comment_blind(make_route_config(c), c.Param("topic_num"), c.Param("num"), c.Request.PostForm)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/thread/:topic_num/comment/:num/delete", func(c *gin.Context) {
		data := route.View_thread_comment_delete(make_route_config(c), c.Param("topic_num"), c.Param("num"), nil)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/thread/:topic_num/comment/:num/delete", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_thread_comment_delete(make_route_config(c), c.Param("topic_num"), c.Param("num"), c.Request.PostForm)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/thread/:topic_num/tool", func(c *gin.Context) {
		data := route.View_thread_tool(make_route_config(c), c.Param("topic_num"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/thread/:topic_num/setting", func(c *gin.Context) {
		data := route.View_thread_setting(make_route_config(c), c.Param("topic_num"), nil)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/thread/:topic_num/setting", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_thread_setting(make_route_config(c), c.Param("topic_num"), c.Request.PostForm)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/thread/:topic_num/acl", func(c *gin.Context) {
		data := route.View_thread_acl(make_route_config(c), c.Param("topic_num"), nil)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/thread/:topic_num/acl", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_thread_acl(make_route_config(c), c.Param("topic_num"), c.Request.PostForm)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/thread/:topic_num/delete", func(c *gin.Context) {
		data := route.View_thread_delete(make_route_config(c), c.Param("topic_num"), nil)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/thread/:topic_num/delete", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_thread_delete(make_route_config(c), c.Param("topic_num"), c.Request.PostForm)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/thread/:topic_num/change", func(c *gin.Context) {
		data := route.View_thread_change(make_route_config(c), c.Param("topic_num"), nil)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/thread/:topic_num/change", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_thread_change(make_route_config(c), c.Param("topic_num"), c.Request.PostForm)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/thread/:topic_num", func(c *gin.Context) {
		data := route.View_thread_route(make_route_config(c), c.Param("topic_num"), "", nil)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/thread/:topic_num", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_thread_route(make_route_config(c), c.Param("topic_num"), "", c.Request.PostForm)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
}
