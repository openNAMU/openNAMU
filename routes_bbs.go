package main

import (
	"net/http"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func Register_bbs_routes(r *gin.Engine) {
	r.GET("/vote", func(c *gin.Context) {
		route_data := route.View_vote_list(Make_route_config(c), "open", "1")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/main", func(c *gin.Context) {
		route_data := route.View_bbs_main(Make_route_config(c), "1")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/contributor", func(c *gin.Context) {
		route_data := route.View_bbs_contributor(Make_route_config(c))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/make", func(c *gin.Context) {
		route_data := route.View_bbs_make(Make_route_config(c))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/bbs/make", func(c *gin.Context) {
		bbs_name := c.PostForm("bbs_name")
		bbs_type := c.PostForm("bbs_type")

		route_data := route.View_bbs_make_post(Make_route_config(c), bbs_name, bbs_type)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/edit/:set_id", func(c *gin.Context) {
		route_data := route.View_bbs_edit(Make_route_config(c), c.Param("set_id"), "", "", "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/bbs/edit/:set_id", func(c *gin.Context) {
		route_data := route.View_bbs_edit_post(
			Make_route_config(c),
			c.Param("set_id"),
			"",
			"",
			c.PostForm("title"),
			c.PostForm("content"),
			c.PostForm("prefix"),
			c.PostForm("tags"),
			Captcha_response_internal(c),
			"",
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/edit/:set_id/document/:document", func(c *gin.Context) {
		route_data := route.View_bbs_edit(Make_route_config(c), c.Param("set_id"), "", "", c.Param("document"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/bbs/edit/:set_id/document/:document", func(c *gin.Context) {
		route_data := route.View_bbs_edit_post(
			Make_route_config(c),
			c.Param("set_id"),
			"",
			"",
			c.PostForm("title"),
			c.PostForm("content"),
			c.PostForm("prefix"),
			c.PostForm("tags"),
			Captcha_response_internal(c),
			c.Param("document"),
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/edit/:set_id/:set_code", func(c *gin.Context) {
		route_data := route.View_bbs_edit(Make_route_config(c), c.Param("set_id"), c.Param("set_code"), "", "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/bbs/edit/:set_id/:set_code", func(c *gin.Context) {
		route_data := route.View_bbs_edit_post(
			Make_route_config(c),
			c.Param("set_id"),
			c.Param("set_code"),
			"",
			c.PostForm("title"),
			c.PostForm("content"),
			c.PostForm("prefix"),
			c.PostForm("tags"),
			Captcha_response_internal(c),
			"",
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/edit/:set_id/:set_code/:comment_code", func(c *gin.Context) {
		route_data := route.View_bbs_edit(
			Make_route_config(c),
			c.Param("set_id"),
			c.Param("set_code"),
			c.Param("comment_code"),
			"",
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/bbs/edit/:set_id/:set_code/:comment_code", func(c *gin.Context) {
		route_data := route.View_bbs_edit_post(
			Make_route_config(c),
			c.Param("set_id"),
			c.Param("set_code"),
			c.Param("comment_code"),
			c.PostForm("title"),
			c.PostForm("content"),
			"",
			"",
			Captcha_response_internal(c),
			"",
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/in/:set_id", func(c *gin.Context) {
		route_data := route.View_bbs_in(Make_route_config(c), c.Param("set_id"), "1", "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/in/:set_id/:page_num", func(c *gin.Context) {
		route_data := route.View_bbs_in(Make_route_config(c), c.Param("set_id"), c.Param("page_num"), "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/in/:set_id/view/:page_num", func(c *gin.Context) {
		route_data := route.View_bbs_in(Make_route_config(c), c.Param("set_id"), c.Param("page_num"), "view")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/in/:set_id/activity/:page_num", func(c *gin.Context) {
		route_data := route.View_bbs_in(Make_route_config(c), c.Param("set_id"), c.Param("page_num"), "activity")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/in/:set_id/comment/:page_num", func(c *gin.Context) {
		route_data := route.View_bbs_in(Make_route_config(c), c.Param("set_id"), c.Param("page_num"), "comment")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/in/:set_id/tabom/:page_num", func(c *gin.Context) {
		route_data := route.View_bbs_in(Make_route_config(c), c.Param("set_id"), c.Param("page_num"), "tabom")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/in/:set_id/excellent/:page_num", func(c *gin.Context) {
		route_data := route.View_bbs_in(Make_route_config(c), c.Param("set_id"), c.Param("page_num"), "excellent")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/in/:set_id/filter/*filter_data", func(c *gin.Context) {
		route_data := route.View_bbs_in_filter(Make_route_config(c), c.Param("set_id"), c.Param("filter_data"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/bbs/in/:set_id/filter", func(c *gin.Context) {
		route_data := route.View_bbs_in_filter_post(
			c.Param("set_id"),
			c.PostForm("comment_min"),
			c.PostForm("commented"),
			c.PostForm("comment_user"),
			c.PostForm("tabom_min"),
			c.PostForm("mine"),
			c.PostForm("participate"),
			c.PostForm("tabom_user"),
			c.PostForm("author"),
			c.PostForm("prefix"),
			c.PostForm("tag"),
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/set/:set_id", func(c *gin.Context) {
		route_data := route.View_bbs_set(Make_route_config(c), c.Param("set_id"), nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/bbs/set/:set_id", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		route_data := route.View_bbs_set(Make_route_config(c), c.Param("set_id"), c.Request.PostForm)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/tool/:set_id/:set_code", func(c *gin.Context) {
		route_data := route.View_bbs_in_w_tool(Make_route_config(c), c.Param("set_id"), c.Param("set_code"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/bbs/tool/:set_id/:set_code", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		route_data := route.View_bbs_in_w_tool_post(
			Make_route_config(c),
			c.Param("set_id"),
			c.Param("set_code"),
			c.Request.PostForm,
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/tool/:set_id/:set_code/:comment_code", func(c *gin.Context) {
		route_data := route.View_bbs_in_w_comment_tool(
			Make_route_config(c),
			c.Param("set_id"),
			c.Param("set_code"),
			c.Param("comment_code"),
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/delete/:set_id", func(c *gin.Context) {
		route_data := route.View_bbs_delete(Make_route_config(c), c.Param("set_id"), "", "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/bbs/delete/:set_id", func(c *gin.Context) {
		route_data := route.View_bbs_delete_post(Make_route_config(c), c.Param("set_id"), "", "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/delete/:set_id/:set_code", func(c *gin.Context) {
		route_data := route.View_bbs_delete(Make_route_config(c), c.Param("set_id"), c.Param("set_code"), "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/bbs/delete/:set_id/:set_code", func(c *gin.Context) {
		route_data := route.View_bbs_delete_post(Make_route_config(c), c.Param("set_id"), c.Param("set_code"), "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/delete/:set_id/:set_code/:comment_code", func(c *gin.Context) {
		route_data := route.View_bbs_delete(
			Make_route_config(c),
			c.Param("set_id"),
			c.Param("set_code"),
			c.Param("comment_code"),
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/bbs/delete/:set_id/:set_code/:comment_code", func(c *gin.Context) {
		route_data := route.View_bbs_delete_post(
			Make_route_config(c),
			c.Param("set_id"),
			c.Param("set_code"),
			c.Param("comment_code"),
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/pinned/:set_id/:set_code", func(c *gin.Context) {
		route_data := route.View_bbs_pinned(Make_route_config(c), c.Param("set_id"), c.Param("set_code"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/bbs/pinned/:set_id/:set_code", func(c *gin.Context) {
		route_data := route.View_bbs_pinned_post(Make_route_config(c), c.Param("set_id"), c.Param("set_code"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/pinned/:set_id/:set_code/:comment_code", func(c *gin.Context) {
		route_data := route.View_bbs_comment_pinned(Make_route_config(c), c.Param("set_id"), c.Param("set_code"), c.Param("comment_code"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/bbs/pinned/:set_id/:set_code/:comment_code", func(c *gin.Context) {
		route_data := route.View_bbs_comment_pinned_post(Make_route_config(c), c.Param("set_id"), c.Param("set_code"), c.Param("comment_code"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
	r.GET("/bbs/blind/:set_id/:set_code", func(c *gin.Context) {
		route_data := route.View_bbs_in_w_blind(Make_route_config(c), c.Param("set_id"), c.Param("set_code"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/bbs/blind/:set_id/:set_code", func(c *gin.Context) {
		route_data := route.View_bbs_in_w_blind_post(Make_route_config(c), c.Param("set_id"), c.Param("set_code"), c.PostForm("blind"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
	r.GET("/bbs/blind/:set_id/:set_code/:comment_code", func(c *gin.Context) {
		route_data := route.View_bbs_in_w_comment_blind(Make_route_config(c), c.Param("set_id"), c.Param("set_code"), c.Param("comment_code"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/bbs/blind/:set_id/:set_code/:comment_code", func(c *gin.Context) {
		route_data := route.View_bbs_in_w_comment_blind_post(Make_route_config(c), c.Param("set_id"), c.Param("set_code"), c.Param("comment_code"), c.PostForm("blind"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/raw/:set_id/:set_code", func(c *gin.Context) {
		route_data := route.View_bbs_raw(Make_route_config(c), c.Param("set_id"), c.Param("set_code"), "")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/raw/:set_id/:set_code/:comment_code", func(c *gin.Context) {
		route_data := route.View_bbs_raw(Make_route_config(c), c.Param("set_id"), c.Param("set_code"), c.Param("comment_code"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/w/:set_id/:set_code", func(c *gin.Context) {
		route_data := route.View_bbs_in_w(c, Make_route_config(c), c.Param("set_id"), c.Param("set_code"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
	r.GET("/bbs/w/:set_id/:set_code/page/:page", func(c *gin.Context) {
		route_data := route.View_bbs_in_w(c, Make_route_config(c), c.Param("set_id"), c.Param("set_code"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
	r.GET("/bbs/w/:set_id/:set_code/page/:page/comment/:comment_select", func(c *gin.Context) {
		route_data := route.View_bbs_in_w(c, Make_route_config(c), c.Param("set_id"), c.Param("set_code"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/w/:set_id/:set_code/comment/:comment_select", func(c *gin.Context) {
		route_data := route.View_bbs_in_w(c, Make_route_config(c), c.Param("set_id"), c.Param("set_code"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/bbs/w/:set_id/:set_code", func(c *gin.Context) {
		route_data := route.View_bbs_in_w_post_secure(
			Make_route_config(c),
			c.Param("set_id"),
			c.Param("set_code"),
			c.PostForm("comment_select"),
			c.PostForm("content"),
			Captcha_response_internal(c),
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/bbs/w/:set_id/:set_code/tabom", func(c *gin.Context) {
		route_data := route.View_bbs_in_w_tabom_post(
			Make_route_config(c),
			c.Param("set_id"),
			c.Param("set_code"),
			c.PostForm("vote_type"),
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/w/:set_id/:set_code/comment_tabom/:comment_code/:vote_type", func(c *gin.Context) {
		route_data := route.View_bbs_in_w_comment_tabom(
			Make_route_config(c),
			c.Param("set_id"),
			c.Param("set_code"),
			c.Param("comment_code"),
			c.Param("vote_type"),
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/bbs/w/:set_id/:set_code/comment/:comment_code/tabom", func(c *gin.Context) {
		route_data := route.View_bbs_in_w_comment_tabom_post(
			Make_route_config(c),
			c.Param("set_id"),
			c.Param("set_code"),
			c.Param("comment_code"),
			c.PostForm("vote_type"),
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
	r.GET("/bbs/report/:set_id/:set_code", func(c *gin.Context) {
		route_data := route.View_bbs_report(
			Make_route_config(c),
			c.Param("set_id"),
			c.Param("set_code"),
			"",
			"",
			"",
			false,
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/bbs/report/:set_id/:set_code", func(c *gin.Context) {
		route_data := route.View_bbs_report(
			Make_route_config(c),
			c.Param("set_id"),
			c.Param("set_code"),
			"",
			c.PostForm("reason"),
			Captcha_response_internal(c),
			true,
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/report/:set_id/:set_code/comment/:comment_code", func(c *gin.Context) {
		route_data := route.View_bbs_report(
			Make_route_config(c),
			c.Param("set_id"),
			c.Param("set_code"),
			c.Param("comment_code"),
			"",
			"",
			false,
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/bbs/report/:set_id/:set_code/comment/:comment_code", func(c *gin.Context) {
		route_data := route.View_bbs_report(
			Make_route_config(c),
			c.Param("set_id"),
			c.Param("set_code"),
			c.Param("comment_code"),
			c.PostForm("reason"),
			Captcha_response_internal(c),
			true,
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
}
