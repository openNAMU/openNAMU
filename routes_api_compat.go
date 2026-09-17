package main

import (
	"net/http"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func Compat_doc_name(c *gin.Context, key string) string {
	return strings.TrimPrefix(c.Param(key), "/")
}

func Compat_api_data(c *gin.Context, data map[string]any) {
	c.JSON(http.StatusOK, data)
}

func Register_api_compat_routes(r *gin.Engine) {
	r.GET("/api/raw_exist/*doc_name", func(c *gin.Context) {
		Compat_api_data(c, route.Api_w_raw(Make_route_config(c), Compat_doc_name(c, "doc_name"), "1", ""))
	})
	r.GET("/api/raw_rev/:rev/*doc_name", func(c *gin.Context) {
		Compat_api_data(c, route.Api_w_raw(Make_route_config(c), Compat_doc_name(c, "doc_name"), "", c.Param("rev")))
	})
	r.GET("/api/raw/*doc_name", func(c *gin.Context) {
		Compat_api_data(c, route.Api_w_raw(Make_route_config(c), Compat_doc_name(c, "doc_name"), "", ""))
	})
	r.GET("/api/random", func(c *gin.Context) {
		Compat_api_data(c, route.Api_w_random(Make_route_config(c)))
	})
	r.GET("/api/sha224/*data", func(c *gin.Context) {
		Compat_api_data(c, route.Api_func_sha224(Make_route_config(c), Compat_doc_name(c, "data")))
	})
	r.GET("/api/ip/*data", func(c *gin.Context) {
		Compat_api_data(c, route.Api_func_ip(Make_route_config(c), Compat_doc_name(c, "data")))
	})
	r.GET("/api/lang/*data", func(c *gin.Context) {
		Compat_api_data(c, route.Api_func_language(Make_route_config(c), Compat_doc_name(c, "data"), "", "on"))
	})
	r.GET("/api/lang_safe/*data", func(c *gin.Context) {
		Compat_api_data(c, route.Api_func_language(Make_route_config(c), Compat_doc_name(c, "data"), "on", "on"))
	})
	r.POST("/api/lang", func(c *gin.Context) {
		Compat_api_data(c, route.Api_func_language(Make_route_config(c), c.PostForm("data"), c.PostForm("safe"), "on"))
	})
	r.GET("/api/xref/:page/*doc_name", func(c *gin.Context) {
		Compat_api_data(c, route.Api_w_xref(Make_route_config(c), c.Param("page"), Compat_doc_name(c, "doc_name"), "1"))
	})
	r.GET("/api/xref_this/:page/*doc_name", func(c *gin.Context) {
		Compat_api_data(c, route.Api_w_xref(Make_route_config(c), c.Param("page"), Compat_doc_name(c, "doc_name"), "2"))
	})
	r.GET("/api/image/*name", func(c *gin.Context) {
		Compat_api_data(c, route.Api_image_exist(Compat_doc_name(c, "name")))
	})

	r.GET("/api/skin_info", Compat_skin_info)
	r.GET("/api/skin_info/:name", Compat_skin_info)
	r.GET("/api/bbs/w/:sub_code", func(c *gin.Context) {
		set_id, set_code := Parse_bbs_code(c.Param("sub_code"))
		Compat_api_data(c, route.Api_bbs_w(Make_route_config(c), set_id, set_code))
	})
	r.GET("/api/bbs/w/comment/:sub_code", func(c *gin.Context) {
		Compat_api_data(c, route.Api_bbs_w_comment(Make_route_config(c), "normal", c.Param("sub_code")))
	})
	r.GET("/api/bbs/w/comment_one/:sub_code", func(c *gin.Context) {
		Compat_api_data(c, route.Api_bbs_w_comment_one(Make_route_config(c), false, "normal", c.Param("sub_code")))
	})

	r.GET("/api/recent_changes", func(c *gin.Context) {
		Compat_api_data_cors(c, route.Api_list_recent_change(Make_route_config(c), "normal", "10", "1"))
	})
	r.GET("/api/recent_discuss/*data", Compat_recent_discuss)
	r.GET("/api/recent_change/:limit/:set_type/:num", func(c *gin.Context) {
		Compat_api_data_cors(c, route.Api_list_recent_change(Make_route_config(c), c.Param("set_type"), c.Param("limit"), c.Param("num")))
	})
	r.GET("/api/search/*keyword", func(c *gin.Context) {
		Compat_api_data(c, route.Api_func_search(Make_route_config(c), Compat_doc_name(c, "keyword"), "1", "title"))
	})
	r.GET("/api/search_page/:num/*keyword", func(c *gin.Context) {
		Compat_api_data(c, route.Api_func_search(Make_route_config(c), Compat_doc_name(c, "keyword"), c.Param("num"), "title"))
	})
	r.GET("/api/search_data/*keyword", func(c *gin.Context) {
		Compat_api_data(c, route.Api_func_search(Make_route_config(c), Compat_doc_name(c, "keyword"), "1", "data"))
	})
	r.GET("/api/search_data_page/:num/*keyword", func(c *gin.Context) {
		Compat_api_data(c, route.Api_func_search(Make_route_config(c), Compat_doc_name(c, "keyword"), c.Param("num"), "data"))
	})

	r.GET("/api/v2/recent_change/:set_type/:num", func(c *gin.Context) {
		Compat_api_data_cors(c, route.Api_list_recent_change(Make_route_config(c), c.Param("set_type"), "50", c.Param("num")))
	})
	r.GET("/api/v2/recent_discuss/:set_type/:num", func(c *gin.Context) {
		Compat_api_data_cors(c, route.Api_list_recent_discuss(Make_route_config(c), "50", c.Param("num"), c.Param("set_type")))
	})
	r.GET("/api/v2/recent_block/:set_type/:num", func(c *gin.Context) {
		Compat_api_data(c, route.Api_list_recent_block(Make_route_config(c), c.Param("num"), c.Param("set_type"), "", ""))
	})
	r.GET("/api/v2/recent_block_user/:set_type/:num/:user_name", func(c *gin.Context) {
		Compat_api_data(c, route.Api_list_recent_block(Make_route_config(c), c.Param("num"), c.Param("set_type"), "", c.Param("user_name")))
	})
	r.GET("/api/v2/list/document/old/:num", func(c *gin.Context) {
		Compat_api_data(c, route.Api_list_old_page(Make_route_config(c), c.Param("num"), "old"))
	})
	r.GET("/api/v2/list/document/new/:num", func(c *gin.Context) {
		Compat_api_data(c, route.Api_list_old_page(Make_route_config(c), c.Param("num"), "new"))
	})
	r.GET("/api/v2/list/document/:num", func(c *gin.Context) {
		Compat_api_data(c, route.Api_list_title_index(Make_route_config(c), c.Param("num")))
	})
	r.GET("/api/v2/history/:num/:set_type/*doc_name", func(c *gin.Context) {
		Compat_api_data(c, route.Api_list_history(Make_route_config(c), Compat_doc_name(c, "doc_name"), c.Param("set_type"), c.Param("num")))
	})
	r.GET("/api/v2/topic/:num/:set_type/*doc_name", func(c *gin.Context) {
		Compat_api_data(c, route.Api_topic_list(Make_route_config(c), c.Param("num"), Compat_doc_name(c, "doc_name"), c.Param("set_type")))
	})
	r.GET("/api/v2/doc_star_doc/:num/*doc_name", func(c *gin.Context) {
		Compat_api_data(c, route.Api_w_watch_list(Make_route_config(c), Compat_doc_name(c, "doc_name"), c.Param("num"), "star_doc"))
	})
	r.GET("/api/v2/doc_watch_list/:num/*doc_name", func(c *gin.Context) {
		Compat_api_data(c, route.Api_w_watch_list(Make_route_config(c), Compat_doc_name(c, "doc_name"), c.Param("num"), "watchlist"))
	})
	r.GET("/api/v2/user/rankup", func(c *gin.Context) {
		Compat_api_data(c, route.Api_user_rankup(Make_route_config(c), ""))
	})
	r.GET("/api/v2/user/rankup/:name", func(c *gin.Context) {
		Compat_api_data(c, route.Api_user_rankup(Make_route_config(c), c.Param("name")))
	})
	r.PATCH("/api/v2/user/rankup", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		Compat_api_data(c, route.Api_user_rankup_patch(Make_route_config(c), c.Request.PostForm))
	})
	r.GET("/api/v2/user_menu/*ip", func(c *gin.Context) {
		Compat_api_data(c, route.Api_func_ip_menu(Make_route_config(c), Compat_doc_name(c, "ip"), "user"))
	})
}

func Compat_skin_info(c *gin.Context) {
	if c.Param("name") == "all" {
		data, ok := route.Api_skin_info_all(Make_route_config(c))
		if ok {
			Write_data(c, http.StatusOK, "application/json; charset=utf-8", data)
			return
		}
	}
	data, ok := route.Api_skin_info(Make_route_config(c), c.Param("name"))
	if !ok {
		c.JSON(http.StatusNotFound, map[string]string{"response": "not found"})
		return
	}
	Write_data(c, http.StatusOK, "application/json; charset=utf-8", data)
}

func Compat_api_data_cors(c *gin.Context, data map[string]any) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type")
	c.Header("Access-Control-Allow-Methods", "GET")
	Compat_api_data(c, data)
}
