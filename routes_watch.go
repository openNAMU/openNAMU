package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opennamu/route"
	"strings"
)

func register_watch_routes(r *gin.Engine) {
	r.GET("/watch_list", func(c *gin.Context) {
		route_data := route.View_user_watch_list(make_route_config(c), "1", "watchlist")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/star_doc", func(c *gin.Context) {
		route_data := route.View_user_watch_list(make_route_config(c), "1", "star_doc")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/doc_watch_list/:count/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_watch_list(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), strings.TrimPrefix(c.Param("count"), "/"), "watchlist")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/doc_star_doc/:count/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_watch_list(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), strings.TrimPrefix(c.Param("count"), "/"), "star_doc")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/star_doc_from/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_watch_list_add_post(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "star_doc_from")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/star_doc_from/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_watch_list_add(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "star_doc_from")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/star_doc/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_watch_list_add_post(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "star_doc")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/star_doc/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_watch_list_add(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "star_doc")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/watch_list_from/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_watch_list_add_post(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "watchlist_from")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/watch_list_from/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_watch_list_add(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "watchlist_from")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/watch_list/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_watch_list_add_post(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "watchlist")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/watch_list/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_watch_list_add(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "watchlist")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/record_bbs/:user_name", func(c *gin.Context) {
		route_data := route.View_record_bbs(make_route_config(c), c.Param("user_name"), "1")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/record_bbs/:user_name/:page", func(c *gin.Context) {
		route_data := route.View_record_bbs(make_route_config(c), c.Param("user_name"), c.Param("page"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/record_bbs_in/:set_id/:user_name", func(c *gin.Context) {
		route_data := route.View_record_bbs_in(make_route_config(c), c.Param("user_name"), c.Param("set_id"), "1")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/record_bbs_in/:set_id/:user_name/:page", func(c *gin.Context) {
		route_data := route.View_record_bbs_in(make_route_config(c), c.Param("user_name"), c.Param("set_id"), c.Param("page"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/record_bbs_comment/:user_name", func(c *gin.Context) {
		route_data := route.View_record_bbs_comment(make_route_config(c), c.Param("user_name"), "1")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/record_bbs_comment/:user_name/:page", func(c *gin.Context) {
		route_data := route.View_record_bbs_comment(make_route_config(c), c.Param("user_name"), c.Param("page"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/record_bbs_comment_in/:set_id/:user_name", func(c *gin.Context) {
		route_data := route.View_record_bbs_comment_in(make_route_config(c), c.Param("user_name"), c.Param("set_id"), "1")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/record_bbs_comment_in/:set_id/:user_name/:page", func(c *gin.Context) {
		route_data := route.View_record_bbs_comment_in(make_route_config(c), c.Param("user_name"), c.Param("set_id"), c.Param("page"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})
}
