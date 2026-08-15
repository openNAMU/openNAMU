package main

import (
	"net/http"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func compat_doc_name(c *gin.Context, key string) string {
	return strings.TrimPrefix(c.Param(key), "/")
}

func compat_api_data(c *gin.Context, data map[string]any) {
	c.JSON(http.StatusOK, data)
}

func register_api_compat_routes(r *gin.Engine) {
	r.GET("/api/raw_exist/*doc_name", func(c *gin.Context) {
		compat_api_data(c, route.Api_w_raw(make_route_config(c), compat_doc_name(c, "doc_name"), "1", ""))
	})
	r.GET("/api/raw_rev/:rev/*doc_name", func(c *gin.Context) {
		compat_api_data(c, route.Api_w_raw(make_route_config(c), compat_doc_name(c, "doc_name"), "", c.Param("rev")))
	})
	r.GET("/api/raw/*doc_name", func(c *gin.Context) {
		compat_api_data(c, route.Api_w_raw(make_route_config(c), compat_doc_name(c, "doc_name"), "", ""))
	})
	r.GET("/api/random", func(c *gin.Context) {
		compat_api_data(c, route.Api_w_random(make_route_config(c)))
	})
	r.GET("/api/sha224/*data", func(c *gin.Context) {
		compat_api_data(c, route.Api_func_sha224(make_route_config(c), compat_doc_name(c, "data")))
	})
	r.GET("/api/ip/*data", func(c *gin.Context) {
		compat_api_data(c, route.Api_func_ip(make_route_config(c), compat_doc_name(c, "data")))
	})
	r.GET("/api/lang/*data", func(c *gin.Context) {
		compat_api_data(c, route.Api_func_language(make_route_config(c), compat_doc_name(c, "data"), "", "on"))
	})
	r.GET("/api/lang", func(c *gin.Context) {
		compat_api_data(c, route.Api_func_language(make_route_config(c), c.Query("data"), c.Query("safe"), "on"))
	})
	r.POST("/api/lang", func(c *gin.Context) {
		compat_api_data(c, route.Api_func_language(make_route_config(c), c.PostForm("data"), c.PostForm("safe"), "on"))
	})
	r.GET("/api/xref/:page/*doc_name", func(c *gin.Context) {
		compat_api_data(c, route.Api_w_xref(make_route_config(c), c.Param("page"), compat_doc_name(c, "doc_name"), "1"))
	})
	r.GET("/api/xref_this/:page/*doc_name", func(c *gin.Context) {
		compat_api_data(c, route.Api_w_xref(make_route_config(c), c.Param("page"), compat_doc_name(c, "doc_name"), "2"))
	})
	r.GET("/api/image/*name", func(c *gin.Context) {
		compat_api_data(c, route.Api_image_exist(compat_doc_name(c, "name")))
	})

	r.GET("/api/skin_info", compat_skin_info)
	r.GET("/api/skin_info/:name", compat_skin_info)
	r.GET("/api/bbs/w/:sub_code", func(c *gin.Context) {
		set_id, set_code := parse_bbs_code(c.Param("sub_code"))
		compat_api_data(c, route.Api_bbs_w(make_route_config(c), set_id, set_code))
	})
	r.GET("/api/bbs/w/comment/:sub_code", func(c *gin.Context) {
		compat_api_data(c, route.Api_bbs_w_comment(make_route_config(c), "normal", c.Param("sub_code")))
	})
	r.GET("/api/bbs/w/comment_one/:sub_code", func(c *gin.Context) {
		compat_api_data(c, route.Api_bbs_w_comment_one(make_route_config(c), false, "normal", c.Param("sub_code")))
	})

	r.GET("/api/recent_changes", func(c *gin.Context) {
		compat_api_data_cors(c, route.Api_list_recent_change(make_route_config(c), "normal", "10", "1"))
	})
	r.GET("/api/recent_discuss/*data", compat_recent_discuss)
	r.GET("/api/recent_change/:limit/:set_type/:num", func(c *gin.Context) {
		compat_api_data_cors(c, route.Api_list_recent_change(make_route_config(c), c.Param("set_type"), c.Param("limit"), c.Param("num")))
	})
	r.GET("/api/search/*keyword", func(c *gin.Context) {
		compat_api_data(c, route.Api_func_search(make_route_config(c), compat_doc_name(c, "keyword"), "1", compat_search_type(c)))
	})
	r.GET("/api/search_page/:num/*keyword", func(c *gin.Context) {
		compat_api_data(c, route.Api_func_search(make_route_config(c), compat_doc_name(c, "keyword"), c.Param("num"), compat_search_type(c)))
	})
	r.GET("/api/search_data/*keyword", func(c *gin.Context) {
		compat_api_data(c, route.Api_func_search(make_route_config(c), compat_doc_name(c, "keyword"), "1", "data"))
	})
	r.GET("/api/search_data_page/:num/*keyword", func(c *gin.Context) {
		compat_api_data(c, route.Api_func_search(make_route_config(c), compat_doc_name(c, "keyword"), c.Param("num"), "data"))
	})

	r.GET("/api/v2/recent_change/:set_type/:num", func(c *gin.Context) {
		compat_api_data_cors(c, route.Api_list_recent_change(make_route_config(c), c.Param("set_type"), "50", c.Param("num")))
	})
	r.GET("/api/v2/recent_discuss/:set_type/:num", func(c *gin.Context) {
		compat_api_data_cors(c, route.Api_list_recent_discuss(make_route_config(c), "50", c.Param("num"), c.Param("set_type")))
	})
	r.GET("/api/v2/recent_block/:set_type/:num", func(c *gin.Context) {
		compat_api_data(c, route.Api_list_recent_block(make_route_config(c), c.Param("num"), c.Param("set_type"), c.Query("why"), ""))
	})
	r.GET("/api/v2/recent_block_user/:set_type/:num/:user_name", func(c *gin.Context) {
		compat_api_data(c, route.Api_list_recent_block(make_route_config(c), c.Param("num"), c.Param("set_type"), c.Query("why"), c.Param("user_name")))
	})
	r.GET("/api/v2/list/document/old/:num", func(c *gin.Context) {
		compat_api_data(c, route.Api_list_old_page(make_route_config(c), c.Param("num"), "old"))
	})
	r.GET("/api/v2/list/document/new/:num", func(c *gin.Context) {
		compat_api_data(c, route.Api_list_old_page(make_route_config(c), c.Param("num"), "new"))
	})
	r.GET("/api/v2/list/document/:num", func(c *gin.Context) {
		compat_api_data(c, route.Api_list_title_index(make_route_config(c), c.Param("num")))
	})
	r.GET("/api/v2/history/:num/:set_type/*doc_name", func(c *gin.Context) {
		compat_api_data(c, route.Api_list_history(make_route_config(c), compat_doc_name(c, "doc_name"), c.Param("set_type"), c.Param("num")))
	})
	r.GET("/api/v2/topic/:num/:set_type/*doc_name", func(c *gin.Context) {
		compat_api_data(c, route.Api_topic_list(make_route_config(c), c.Param("num"), compat_doc_name(c, "doc_name"), c.Param("set_type")))
	})
	r.GET("/api/v2/doc_star_doc/:num/*doc_name", func(c *gin.Context) {
		compat_api_data(c, route.Api_w_watch_list(make_route_config(c), compat_doc_name(c, "doc_name"), c.Param("num"), "star_doc"))
	})
	r.GET("/api/v2/doc_watch_list/:num/*doc_name", func(c *gin.Context) {
		compat_api_data(c, route.Api_w_watch_list(make_route_config(c), compat_doc_name(c, "doc_name"), c.Param("num"), "watchlist"))
	})
	r.GET("/api/v2/user/rankup", func(c *gin.Context) {
		compat_api_data(c, route.Api_user_rankup(make_route_config(c), c.Query("name")))
	})
	r.PATCH("/api/v2/user/rankup", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		compat_api_data(c, route.Api_user_rankup_patch(make_route_config(c), c.Request.PostForm))
	})
	r.GET("/api/v2/user_menu/*ip", func(c *gin.Context) {
		compat_api_data(c, route.Api_func_ip_menu(make_route_config(c), compat_doc_name(c, "ip"), "user"))
	})
}

func compat_search_type(c *gin.Context) string {
	search_type := c.Query("type")
	if search_type == "" {
		search_type = "title"
	}
	return search_type
}

func compat_skin_info(c *gin.Context) {
	if c.Query("all") != "" {
		data, ok := route.Api_skin_info_all(make_route_config(c))
		if ok {
			write_data(c, http.StatusOK, "application/json; charset=utf-8", data)
			return
		}
	}
	data, ok := route.Api_skin_info(make_route_config(c), c.Param("name"))
	if !ok {
		c.JSON(http.StatusNotFound, map[string]string{"response": "not found"})
		return
	}
	write_data(c, http.StatusOK, "application/json; charset=utf-8", data)
}

func compat_api_data_cors(c *gin.Context, data map[string]any) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type")
	c.Header("Access-Control-Allow-Methods", "GET")
	compat_api_data(c, data)
}
