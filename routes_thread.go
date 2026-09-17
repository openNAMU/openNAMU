package main

import (
	"net/http"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func Register_thread_routes(r *gin.Engine) {
	r.GET("/thread/0/*doc_name", func(c *gin.Context) {
		data := route.View_thread_route(Make_route_config(c), "0", strings.TrimPrefix(c.Param("doc_name"), "/"), "1", nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/thread/0/*doc_name", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_thread_route(Make_route_config(c), "0", strings.TrimPrefix(c.Param("doc_name"), "/"), "1", c.Request.PostForm)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/thread/:topic_num/comment/:num/raw", func(c *gin.Context) {
		data := route.View_thread_raw(Make_route_config(c), c.Param("topic_num"), c.Param("num"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/thread/:topic_num/comment/:num/tool", func(c *gin.Context) {
		data := route.View_thread_comment_tool(Make_route_config(c), c.Param("topic_num"), c.Param("num"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/thread/:topic_num/comment/:num/notice", func(c *gin.Context) {
		data := route.View_thread_comment_notice(Make_route_config(c), c.Param("topic_num"), c.Param("num"), nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/thread/:topic_num/comment/:num/notice", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_thread_comment_notice(Make_route_config(c), c.Param("topic_num"), c.Param("num"), c.Request.PostForm)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/thread/:topic_num/comment/:num/blind", func(c *gin.Context) {
		data := route.View_thread_comment_blind(Make_route_config(c), c.Param("topic_num"), c.Param("num"), nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/thread/:topic_num/comment/:num/blind", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_thread_comment_blind(Make_route_config(c), c.Param("topic_num"), c.Param("num"), c.Request.PostForm)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/thread/:topic_num/comment/:num/delete", func(c *gin.Context) {
		data := route.View_thread_comment_delete(Make_route_config(c), c.Param("topic_num"), c.Param("num"), nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/thread/:topic_num/comment/:num/delete", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_thread_comment_delete(Make_route_config(c), c.Param("topic_num"), c.Param("num"), c.Request.PostForm)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/thread/:topic_num/tool", func(c *gin.Context) {
		data := route.View_thread_tool(Make_route_config(c), c.Param("topic_num"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/thread/:topic_num/setting", func(c *gin.Context) {
		data := route.View_thread_setting(Make_route_config(c), c.Param("topic_num"), nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/thread/:topic_num/setting", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_thread_setting(Make_route_config(c), c.Param("topic_num"), c.Request.PostForm)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/thread/:topic_num/acl", func(c *gin.Context) {
		data := route.View_thread_acl(Make_route_config(c), c.Param("topic_num"), nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/thread/:topic_num/acl", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_thread_acl(Make_route_config(c), c.Param("topic_num"), c.Request.PostForm)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/thread/:topic_num/delete", func(c *gin.Context) {
		data := route.View_thread_delete(Make_route_config(c), c.Param("topic_num"), nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/thread/:topic_num/delete", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_thread_delete(Make_route_config(c), c.Param("topic_num"), c.Request.PostForm)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/thread/:topic_num/change", func(c *gin.Context) {
		data := route.View_thread_change(Make_route_config(c), c.Param("topic_num"), nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/thread/:topic_num/change", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_thread_change(Make_route_config(c), c.Param("topic_num"), c.Request.PostForm)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/thread/:topic_num", func(c *gin.Context) {
		data := route.View_thread_route(Make_route_config(c), c.Param("topic_num"), "", "1", nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/thread/:topic_num", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_thread_route(Make_route_config(c), c.Param("topic_num"), "", "1", c.Request.PostForm)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/thread/:topic_num/page/:page", func(c *gin.Context) {
		data := route.View_thread_route(Make_route_config(c), c.Param("topic_num"), "", c.Param("page"), nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
}
