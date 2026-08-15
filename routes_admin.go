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
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_give(make_route_config(c), "many", "", nil)))
	})
	r.POST("/auth/give", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_give(make_route_config(c), "many", "", admin_post_values(c))))
	})
	r.GET("/auth/give_total", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_give(make_route_config(c), "total", "", nil)))
	})
	r.POST("/auth/give_total", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_give(make_route_config(c), "total", "", admin_post_values(c))))
	})
	r.GET("/auth/give/fix/:user_name", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_fix(make_route_config(c), c.Param("user_name"), nil)))
	})
	r.POST("/auth/give/fix/:user_name", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_fix(make_route_config(c), c.Param("user_name"), admin_post_values(c))))
	})
	r.GET("/auth/give/:user_name", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_give(make_route_config(c), "one", c.Param("user_name"), nil)))
	})
	r.POST("/auth/give/:user_name", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_give(make_route_config(c), "one", c.Param("user_name"), admin_post_values(c))))
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

	r.GET("/auth/ban", func(c *gin.Context) { admin_ban(c, "", "", false, nil) })
	r.POST("/auth/ban", func(c *gin.Context) { admin_ban(c, "", "", false, admin_post_values(c)) })
	r.GET("/auth/ban/*name", func(c *gin.Context) {
		name := strings.TrimPrefix(c.Param("name"), "/")
		multiple := name == "multiple"
		if multiple {
			name = ""
		}
		admin_ban(c, name, "", multiple, nil)
	})
	r.POST("/auth/ban/*name", func(c *gin.Context) {
		name := strings.TrimPrefix(c.Param("name"), "/")
		multiple := name == "multiple"
		if multiple {
			name = ""
		}
		admin_ban(c, name, "", multiple, admin_post_values(c))
	})
	r.GET("/auth/ban_cidr/*name", func(c *gin.Context) {
		admin_ban(c, strings.TrimPrefix(c.Param("name"), "/"), "cidr", false, nil)
	})
	r.POST("/auth/ban_cidr/*name", func(c *gin.Context) {
		admin_ban(c, strings.TrimPrefix(c.Param("name"), "/"), "cidr", false, admin_post_values(c))
	})
	r.GET("/auth/ban_regex/*name", func(c *gin.Context) {
		admin_ban(c, strings.TrimPrefix(c.Param("name"), "/"), "regex", false, nil)
	})
	r.POST("/auth/ban_regex/*name", func(c *gin.Context) {
		admin_ban(c, strings.TrimPrefix(c.Param("name"), "/"), "regex", false, admin_post_values(c))
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

func admin_ban(c *gin.Context, name string, ban_type string, multiple bool, values url.Values) {
	write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_ban(make_route_config(c), name, ban_type, multiple, values)))
}

func admin_acl(c *gin.Context, doc_name string, multiple bool, values url.Values) {
	write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_acl(make_route_config(c), doc_name, multiple, values)))
}
