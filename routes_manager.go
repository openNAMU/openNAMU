package main

import (
	"net/http"
	"strconv"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_manager_routes(r *gin.Engine) {
	r.GET("/manager", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_main_manager(make_route_config(c))))
	})
	r.POST("/manager", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_main_manager(make_route_config(c))))
	})

	r.GET("/manager/:num", func(c *gin.Context) {
		manager_redirect(c, "", false)
	})
	r.POST("/manager/:num", func(c *gin.Context) {
		manager_redirect(c, c.PostForm("name"), true)
	})
	r.GET("/manager/:num/*add_2", func(c *gin.Context) {
		manager_redirect(c, "", false)
	})
	r.POST("/manager/:num/*add_2", func(c *gin.Context) {
		manager_redirect(c, c.PostForm("name"), true)
	})
}

func manager_redirect(c *gin.Context, name string, post bool) {
	num, err := strconv.Atoi(c.Param("num"))
	if err != nil {
		c.Redirect(http.StatusFound, "/manager")
		return
	}

	data := route.View_manager_redirect(
		make_route_config(c),
		num,
		strings.TrimPrefix(c.Param("add_2"), "/"),
		name,
		post,
	)
	write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
}
