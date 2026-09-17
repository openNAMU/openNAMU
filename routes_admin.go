package main

import (
	"net/http"
	"net/url"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func Admin_post_values(c *gin.Context) url.Values {
	_ = c.Request.ParseForm()
	return c.Request.PostForm
}

func Register_admin_routes(r *gin.Engine) {
	Register_admin_api_routes(r)
	r.GET("/auth/give", func(c *gin.Context) {
		Admin_give(c, "", "normal", nil)
	})
	r.POST("/auth/give", func(c *gin.Context) {
		Admin_give(c, "", "normal", Admin_post_values(c))
	})
	r.GET("/auth/give_total", func(c *gin.Context) {
		Admin_give_total(c, nil)
	})
	r.POST("/auth/give_total", func(c *gin.Context) {
		Admin_give_total(c, Admin_post_values(c))
	})
	r.GET("/auth/give/fix/:user_name", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_fix(Make_route_config(c), c.Param("user_name"), nil)))
	})
	r.POST("/auth/give/fix/:user_name", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_fix(Make_route_config(c), c.Param("user_name"), Admin_post_values(c))))
	})
	r.GET("/auth/give/regex/*name", func(c *gin.Context) {
		Admin_give(c, strings.TrimPrefix(c.Param("name"), "/"), "regex", nil)
	})
	r.POST("/auth/give/regex/*name", func(c *gin.Context) {
		Admin_give(c, strings.TrimPrefix(c.Param("name"), "/"), "regex", Admin_post_values(c))
	})
	r.GET("/auth/give/cidr/*name", func(c *gin.Context) {
		Admin_give(c, strings.TrimPrefix(c.Param("name"), "/"), "cidr", nil)
	})
	r.POST("/auth/give/cidr/*name", func(c *gin.Context) {
		Admin_give(c, strings.TrimPrefix(c.Param("name"), "/"), "cidr", Admin_post_values(c))
	})
	r.GET("/auth/give/private/*name", func(c *gin.Context) {
		Admin_give(c, strings.TrimPrefix(c.Param("name"), "/"), "private", nil)
	})
	r.POST("/auth/give/private/*name", func(c *gin.Context) {
		Admin_give(c, strings.TrimPrefix(c.Param("name"), "/"), "private", Admin_post_values(c))
	})
	r.GET("/auth/give/:user_name", func(c *gin.Context) {
		Admin_give(c, c.Param("user_name"), "normal", nil)
	})
	r.POST("/auth/give/:user_name", func(c *gin.Context) {
		Admin_give(c, c.Param("user_name"), "normal", Admin_post_values(c))
	})
	r.GET("/auth/invite", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_invite(Make_route_config(c), nil)))
	})
	r.POST("/auth/invite", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_invite(Make_route_config(c), Admin_post_values(c))))
	})

	r.GET("/auth/list", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_list(Make_route_config(c))))
	})
	r.GET("/auth/list/add", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_group(Make_route_config(c), "", nil)))
	})
	r.POST("/auth/list/add", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_group(Make_route_config(c), "", Admin_post_values(c))))
	})
	r.GET("/auth/list/add/*name", func(c *gin.Context) {
		name := strings.TrimPrefix(c.Param("name"), "/")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_group(Make_route_config(c), name, nil)))
	})
	r.POST("/auth/list/add/*name", func(c *gin.Context) {
		name := strings.TrimPrefix(c.Param("name"), "/")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_group(Make_route_config(c), name, Admin_post_values(c))))
	})
	r.GET("/auth/list/delete/*name", func(c *gin.Context) {
		name := strings.TrimPrefix(c.Param("name"), "/")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_group_delete(Make_route_config(c), name, nil)))
	})
	r.POST("/auth/list/delete/*name", func(c *gin.Context) {
		name := strings.TrimPrefix(c.Param("name"), "/")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_group_delete(Make_route_config(c), name, Admin_post_values(c))))
	})

	r.GET("/app_submit", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_app_submit(Make_route_config(c), nil)))
	})
	r.POST("/app_submit", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_app_submit(Make_route_config(c), Admin_post_values(c))))
	})

	r.GET("/acl", func(c *gin.Context) { Admin_acl(c, "", false, nil) })
	r.POST("/acl", func(c *gin.Context) { Admin_acl(c, "", false, Admin_post_values(c)) })
	r.GET("/acl_multiple", func(c *gin.Context) { Admin_acl(c, "", true, nil) })
	r.POST("/acl_multiple", func(c *gin.Context) { Admin_acl(c, "", true, Admin_post_values(c)) })
	r.GET("/acl/*name", func(c *gin.Context) {
		Admin_acl(c, strings.TrimPrefix(c.Param("name"), "/"), false, nil)
	})
	r.POST("/acl/*name", func(c *gin.Context) {
		Admin_acl(c, strings.TrimPrefix(c.Param("name"), "/"), false, Admin_post_values(c))
	})
}

func Admin_give(c *gin.Context, name string, target_type string, values url.Values) {
	Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_give(Make_route_config(c), "one", name, target_type, values)))
}

func Admin_give_total(c *gin.Context, values url.Values) {
	Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_auth_give(Make_route_config(c), "total", "", "normal", values)))
}

func Admin_acl(c *gin.Context, doc_name string, multiple bool, values url.Values) {
	Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_acl(Make_route_config(c), doc_name, multiple, values)))
}
