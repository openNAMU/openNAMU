package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_api_routes(r *gin.Engine) {
	r.GET("/api/version", func(c *gin.Context) {
		route_data := route.Api_version(make_route_config(c))
		c.JSON(http.StatusOK, route_data)
	})

	r.POST("/api/render", func(c *gin.Context) {
		route_data := route.Api_w_render(make_route_config(c), c.PostForm("name"), c.PostForm("data"), "api_view", c.PostForm("option"))
		c.JSON(http.StatusOK, route_data)
	})

	r.POST("/api/render/:tool", func(c *gin.Context) {
		render_type := "api_thread"
		switch c.Param("tool") {
		case "from":
			render_type = "api_from"
		case "include":
			render_type = "api_include"
		case "backlink":
			render_type = "backlink"
		}
		route_data := route.Api_w_render(make_route_config(c), c.PostForm("name"), c.PostForm("data"), render_type, c.PostForm("option"))
		c.JSON(http.StatusOK, route_data)
	})

	r.POST("/api/v2/vote/add", func(c *gin.Context) {
		data := c.PostForm("data")
		if data == "" && strings.HasPrefix(c.GetHeader("Content-Type"), "application/json") {
			body := map[string]any{}
			if c.ShouldBindJSON(&body) == nil {
				encoded, _ := json.Marshal(body)
				data = string(encoded)
			}
		}
		compat_api_data(c, route.Api_vote_add_post(make_route_config(c), data))
	})
	r.POST("/api/vote/add", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		compat_api_data(c, route.Api_vote_add_post(make_route_config(c), c.Request.PostForm.Encode()))
	})

	r.POST("/api/v2/set_reset/*doc_name", func(c *gin.Context) {
		route_data := route.Api_w_set_reset(
			make_route_config(c),
			strings.TrimPrefix(c.Param("doc_name"), "/"),
		)
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/list/markup", func(c *gin.Context) {
		route_data := route.Api_list_markup(make_route_config(c))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/list/acl/:type", func(c *gin.Context) {
		route_data := route.Api_list_acl(make_route_config(c), c.Param("type"))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/user_info/:user_name", func(c *gin.Context) {
		route_data := route.Api_user_info(make_route_config(c), c.Param("user_name"))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/recent_change", func(c *gin.Context) {
		route_data := route.Api_list_recent_change(make_route_config(c), "normal", "10", "1")
		compat_api_data_cors(c, route_data)
	})

	r.GET("/api/recent_change/:limit", func(c *gin.Context) {
		route_data := route.Api_list_recent_change(make_route_config(c), "normal", c.Param("limit"), "1")
		compat_api_data_cors(c, route_data)
	})

	r.GET("/api/recent_discuss", func(c *gin.Context) {
		route_data := route.Api_list_recent_discuss(make_route_config(c), "10", "1", "normal")
		compat_api_data_cors(c, route_data)
	})

	r.POST("/api/v2/lang", func(c *gin.Context) {
		data := c.PostForm("data")
		safe := c.PostForm("safe")
		legacy := c.PostForm("legacy")

		route_data := route.Api_func_language(make_route_config(c), data, safe, legacy)
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/ip/*ip", func(c *gin.Context) {
		route_data := route.Api_func_ip(make_route_config(c), compat_doc_name(c, "ip"))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/ip_menu/*ip", func(c *gin.Context) {
		route_data := route.Api_func_ip_menu(make_route_config(c), compat_doc_name(c, "ip"), "user")
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/user/setting/editor", func(c *gin.Context) {
		route_data := route.Api_user_setting_editor(make_route_config(c))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/page_view/*doc_name", func(c *gin.Context) {
		route_data := route.Api_w_page_view(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"))
		c.JSON(http.StatusOK, route_data)
	})

	r.POST("/api/v2/user/setting/editor", func(c *gin.Context) {
		route_data := route.Api_user_setting_editor_post(make_route_config(c), c.Request.FormValue("data"))
		c.JSON(http.StatusOK, route_data)
	})

	r.DELETE("/api/v2/user/setting/editor", func(c *gin.Context) {
		route_data := route.Api_user_setting_editor_delete(make_route_config(c), c.Request.FormValue("data"))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/page_view_post/*doc_name", func(c *gin.Context) {
		route_data := route.Api_w_page_view_post(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/bbs/w/page_view_post/:set_id/:set_code", func(c *gin.Context) {
		route_data := route.Api_bbs_w_page_view_post(make_route_config(c), c.Param("set_id"), c.Param("set_code"))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/bbs/w/page_view/:set_id/:set_code", func(c *gin.Context) {
		route_data := route.Api_bbs_w_page_view(make_route_config(c), c.Param("set_id"), c.Param("set_code"))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/bbs/main", func(c *gin.Context) {
		route_data := route.Api_bbs(make_route_config(c), "", "1", "")
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/bbs", func(c *gin.Context) {
		route_data := route.Api_bbs_list(make_route_config(c))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/bbs/in/:set_id/:page", func(c *gin.Context) {
		route_data := route.Api_bbs(make_route_config(c), c.Param("set_id"), c.Param("page"), "")
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/bbs/in/:set_id/view/:page", func(c *gin.Context) {
		route_data := route.Api_bbs(make_route_config(c), c.Param("set_id"), c.Param("page"), "view")
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/bbs/set/:set_id/:set_name", func(c *gin.Context) {
		route_data := route.Api_bbs_w_set(make_route_config(c), c.Param("set_id"), c.Param("set_name"))
		c.JSON(http.StatusOK, route_data)
	})

	r.PUT("/api/v2/bbs/set/:set_id/:set_name", func(c *gin.Context) {
		route_data := route.Api_bbs_w_set_put(
			make_route_config(c),
			c.Param("set_id"),
			c.Param("set_name"),
			c.Request.FormValue("data"),
			c.Request.FormValue("coverage"),
		)
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/bbs/w/tabom/:sub_code", func(c *gin.Context) {
		set_id, set_code := parse_bbs_code(c.Param("sub_code"))
		route_data := route.Api_bbs_w_tabom(make_route_config(c), set_id, set_code)
		c.JSON(http.StatusOK, route_data)
	})

	r.POST("/api/v2/bbs/w/tabom/:sub_code", func(c *gin.Context) {
		route_data := route.Api_bbs_w_tabom_post(make_route_config(c), c.Param("sub_code"))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/bbs/w/comment/:sub_code/:tool", func(c *gin.Context) {
		route_data := route.Api_bbs_w_comment(make_route_config(c), c.Param("tool"), c.Param("sub_code"))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/bbs/w/comment_one/:sub_code/:tool", func(c *gin.Context) {
		route_data := route.Api_bbs_w_comment_one(make_route_config(c), false, c.Param("tool"), c.Param("sub_code"))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/bbs/w/:sub_code", func(c *gin.Context) {
		set_id, set_code := parse_bbs_code(c.Param("sub_code"))
		route_data := route.Api_bbs_w(make_route_config(c), set_id, set_code)
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/thread/:topic_num/:tool/:end_num", func(c *gin.Context) {
		route_data := route.Api_topic(
			make_route_config(c),
			"",
			c.Param("topic_num"),
			c.Param("tool"),
			c.Param("end_num"),
		)
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/thread/:topic_num/:tool", func(c *gin.Context) {
		route_data := route.Api_topic(make_route_config(c), c.Param("tool"), c.Param("topic_num"), "", "")
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/thread/:topic_num", func(c *gin.Context) {
		route_data := route.Api_topic(make_route_config(c), "", c.Param("topic_num"), "", "")
		c.JSON(http.StatusOK, route_data)
	})

	// Draft adapter: legacy thread data with the BBS comment response shape.
	r.GET("/api/v2/bbs/thread/:topic_num", func(c *gin.Context) {
		route_data := route.Api_thread_bbs(
			make_route_config(c),
			"",
			c.Param("topic_num"),
			"",
			"",
		)
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/bbs/thread/:topic_num/:tool", func(c *gin.Context) {
		route_data := route.Api_thread_bbs(
			make_route_config(c),
			c.Param("tool"),
			c.Param("topic_num"),
			"",
			"",
		)
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/bbs/thread/:topic_num/:tool/:end_num", func(c *gin.Context) {
		route_data := route.Api_thread_bbs(
			make_route_config(c),
			"",
			c.Param("topic_num"),
			c.Param("tool"),
			c.Param("end_num"),
		)
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/raw/*doc_name", func(c *gin.Context) {
		route_data := route.Api_w_raw(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "", "")
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/raw_exist/*doc_name", func(c *gin.Context) {
		route_data := route.Api_w_raw(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "true", "")
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/raw_rev/:rev/*doc_name", func(c *gin.Context) {
		route_data := route.Api_w_raw(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "", c.Param("rev"))
		c.JSON(http.StatusOK, route_data)
	})
}
