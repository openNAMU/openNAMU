package main

import (
	"net/http"
	"net/url"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func admin_post_values(c *gin.Context) url.Values {
	_ = c.Request.ParseForm()
	return c.Request.PostForm
}

func register_admin_routes(r *gin.Engine) {
	register_admin_api_routes(r)
	r.GET("/auth/give", func(c *gin.Context) {
		admin_give(c, "", "normal", nil)
	})
	r.POST("/auth/give", func(c *gin.Context) {
		admin_give(c, "", "normal", admin_post_values(c))
	})
	r.GET("/auth/give_total", func(c *gin.Context) {
		admin_give_total(c, nil)
	})
	r.POST("/auth/give_total", func(c *gin.Context) {
		admin_give_total(c, admin_post_values(c))
	})
	r.GET("/auth/give/fix/:user_name", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_fix(make_route_config(c), c.Param("user_name"), nil)))
	})
	r.POST("/auth/give/fix/:user_name", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_fix(make_route_config(c), c.Param("user_name"), admin_post_values(c))))
	})
	r.GET("/auth/give/regex/*name", func(c *gin.Context) {
		admin_give(c, strings.TrimPrefix(c.Param("name"), "/"), "regex", nil)
	})
	r.POST("/auth/give/regex/*name", func(c *gin.Context) {
		admin_give(c, strings.TrimPrefix(c.Param("name"), "/"), "regex", admin_post_values(c))
	})
	r.GET("/auth/give/cidr/*name", func(c *gin.Context) {
		admin_give(c, strings.TrimPrefix(c.Param("name"), "/"), "cidr", nil)
	})
	r.POST("/auth/give/cidr/*name", func(c *gin.Context) {
		admin_give(c, strings.TrimPrefix(c.Param("name"), "/"), "cidr", admin_post_values(c))
	})
	r.GET("/auth/give/private/*name", func(c *gin.Context) {
		admin_give(c, strings.TrimPrefix(c.Param("name"), "/"), "private", nil)
	})
	r.POST("/auth/give/private/*name", func(c *gin.Context) {
		admin_give(c, strings.TrimPrefix(c.Param("name"), "/"), "private", admin_post_values(c))
	})
	r.GET("/auth/give/:user_name", func(c *gin.Context) {
		admin_give(c, c.Param("user_name"), "normal", nil)
	})
	r.POST("/auth/give/:user_name", func(c *gin.Context) {
		admin_give(c, c.Param("user_name"), "normal", admin_post_values(c))
	})

	r.GET("/auth/list", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_list(make_route_config(c))))
	})
	r.GET("/auth/list/add", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_group(make_route_config(c), "", nil)))
	})
	r.POST("/auth/list/add", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_group(make_route_config(c), "", admin_post_values(c))))
	})
	r.GET("/auth/list/add/*name", func(c *gin.Context) {
		name := strings.TrimPrefix(c.Param("name"), "/")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_group(make_route_config(c), name, nil)))
	})
	r.POST("/auth/list/add/*name", func(c *gin.Context) {
		name := strings.TrimPrefix(c.Param("name"), "/")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_group(make_route_config(c), name, admin_post_values(c))))
	})
	r.GET("/auth/list/delete/*name", func(c *gin.Context) {
		name := strings.TrimPrefix(c.Param("name"), "/")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_group_delete(make_route_config(c), name, nil)))
	})
	r.POST("/auth/list/delete/*name", func(c *gin.Context) {
		name := strings.TrimPrefix(c.Param("name"), "/")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_group_delete(make_route_config(c), name, admin_post_values(c))))
	})

	r.GET("/app_submit", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_app_submit(make_route_config(c), nil)))
	})
	r.POST("/app_submit", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_app_submit(make_route_config(c), admin_post_values(c))))
	})

	r.GET("/acl", func(c *gin.Context) { admin_acl(c, "", false, nil) })
	r.POST("/acl", func(c *gin.Context) { admin_acl(c, "", false, admin_post_values(c)) })
	r.GET("/acl_multiple", func(c *gin.Context) { admin_acl(c, "", true, nil) })
	r.POST("/acl_multiple", func(c *gin.Context) { admin_acl(c, "", true, admin_post_values(c)) })
	r.GET("/acl/*name", func(c *gin.Context) {
		admin_acl(c, strings.TrimPrefix(c.Param("name"), "/"), false, nil)
	})
	r.POST("/acl/*name", func(c *gin.Context) {
		admin_acl(c, strings.TrimPrefix(c.Param("name"), "/"), false, admin_post_values(c))
	})
}

func admin_give(c *gin.Context, name string, target_type string, values url.Values) {
	write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_give(make_route_config(c), "one", name, target_type, values)))
}

func admin_give_total(c *gin.Context, values url.Values) {
	write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_give(make_route_config(c), "total", "", "normal", values)))
}

func admin_acl(c *gin.Context, doc_name string, multiple bool, values url.Values) {
	write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_acl(make_route_config(c), doc_name, multiple, values)))
}
