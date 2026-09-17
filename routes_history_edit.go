package main

import (
	"net/http"
	"strings"

	"opennamu/route"
	"opennamu/route/tool"

	"github.com/gin-gonic/gin"
)

func Register_history_edit_routes(r *gin.Engine) {
	r.GET("/move_all", func(c *gin.Context) {
		data := route.View_edit_move_all(Make_route_config(c), nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/move_all", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_edit_move_all(Make_route_config(c), c.Request.PostForm)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/render/:rev/*doc_name", func(c *gin.Context) {
		route_data := route.View_render(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Param("rev"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/history/*doc_name", func(c *gin.Context) {
		route_data := route.View_list_history(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "", "1")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/history/*doc_name", func(c *gin.Context) {
		doc_name := strings.TrimPrefix(c.Param("doc_name"), "/")
		a := c.PostForm("a")
		b := c.PostForm("b")

		route_data := route.View_list_history_post(Make_route_config(c), doc_name, a, b)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/history_page/:num/:set_type/*doc_name", func(c *gin.Context) {
		route_data := route.View_list_history(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Param("set_type"), c.Param("num"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/history_page/:num/:set_type/*doc_name", func(c *gin.Context) {
		doc_name := strings.TrimPrefix(c.Param("doc_name"), "/")
		a := c.PostForm("a")
		b := c.PostForm("b")

		route_data := route.View_list_history_post(Make_route_config(c), doc_name, a, b)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/diff/:before_rev/:after_rev/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_diff(
			Make_route_config(c),
			strings.TrimPrefix(c.Param("doc_name"), "/"),
			c.Param("before_rev"),
			c.Param("after_rev"),
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/edit/*doc_name", func(c *gin.Context) {
		route_data := route.View_edit(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/edit_load/:load/*doc_name", func(c *gin.Context) {
		load_data, err := tool.Get_base64_decode(c.Param("load"))
		if err != nil {
			c.Redirect(http.StatusFound, "/manager")
			return
		}
		route_data := route.View_edit(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), load_data)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/edit/*doc_name", func(c *gin.Context) {
		doc_name := strings.TrimPrefix(c.Param("doc_name"), "/")
		data := c.PostForm("content")
		preview := c.PostForm("preview")
		if preview == "normal" || preview == "dark" {
			route_data := route.View_edit_preview(Make_route_config(c), doc_name, data, preview, c.PostForm("send"))
			Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
			return
		}
		send := c.PostForm("send")
		agree := c.PostForm("copyright_agreement")

		route_data := route.View_edit_post(Make_route_config(c), doc_name, data, send, agree, Captcha_response_internal(c), c.PostForm("ver"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/edit_preview/:mode/*doc_name", func(c *gin.Context) {
		doc_name := strings.TrimPrefix(c.Param("doc_name"), "/")
		route_data := route.View_edit_preview(Make_route_config(c), doc_name, c.PostForm("content"), c.Param("mode"), c.PostForm("send"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/move/*doc_name", func(c *gin.Context) {
		route_data := route.View_edit_move(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/move/*doc_name", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		route_data := route.View_edit_move(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Request.PostForm)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/delete/*doc_name", func(c *gin.Context) {
		doc_name := strings.TrimPrefix(c.Param("doc_name"), "/")
		route_data := route.View_edit_delete(Make_route_config(c), doc_name)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/delete/*doc_name", func(c *gin.Context) {
		doc_name := strings.TrimPrefix(c.Param("doc_name"), "/")
		route_data := route.View_edit_delete_post(
			Make_route_config(c),
			doc_name,
			c.PostForm("send"),
			c.PostForm("copyright_agreement"),
			Captcha_response_internal(c),
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/delete_file/*doc_name", func(c *gin.Context) {
		route_data := route.View_edit_file_delete(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/delete_file/*doc_name", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		route_data := route.View_edit_file_delete(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Request.PostForm)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/rename_file/*doc_name", func(c *gin.Context) {
		route_data := route.View_edit_file_rename(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/rename_file/*doc_name", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		route_data := route.View_edit_file_rename(Make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Request.PostForm)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/revert/:rev/*doc_name", func(c *gin.Context) {
		route_data := route.View_edit_revert(
			Make_route_config(c),
			strings.TrimPrefix(c.Param("doc_name"), "/"),
			c.Param("rev"),
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/revert/:rev/*doc_name", func(c *gin.Context) {
		route_data := route.View_edit_revert_post(
			Make_route_config(c),
			strings.TrimPrefix(c.Param("doc_name"), "/"),
			c.Param("rev"),
			c.PostForm("send"),
			c.PostForm("copyright_agreement"),
			Captcha_response_internal(c),
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/upload", Upload_post)
	r.POST("/upload/*name", Upload_post)

	r.GET("/view/*name", route.View_view_file)
	r.GET("/views/*name", route.View_view_file)
	r.GET("/image/*name", route.View_view_image_file)
}
