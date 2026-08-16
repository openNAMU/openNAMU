package main

import (
	"net/http"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_filter_routes(r *gin.Engine) {
	r.GET("/filter/:kind", func(c *gin.Context) {
		data := route.View_filter(make_route_config(c), c.Param("kind"))
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/filter/:kind/add", func(c *gin.Context) {
		data := route.View_filter_add(make_route_config(c), c.Param("kind"), "", nil)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/filter/:kind/add", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_filter_add(make_route_config(c), c.Param("kind"), "", c.Request.PostForm)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/filter/:kind/add/*name", func(c *gin.Context) {
		name := strings.TrimPrefix(c.Param("name"), "/")
		data := route.View_filter_add(make_route_config(c), c.Param("kind"), name, nil)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/filter/:kind/add/*name", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		name := strings.TrimPrefix(c.Param("name"), "/")
		data := route.View_filter_add(make_route_config(c), c.Param("kind"), name, c.Request.PostForm)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/filter/:kind/del/*name", func(c *gin.Context) {
		name := strings.TrimPrefix(c.Param("name"), "/")
		data := route.View_filter_delete(make_route_config(c), c.Param("kind"), name, nil)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/filter/:kind/del/*name", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		name := strings.TrimPrefix(c.Param("name"), "/")
		data := route.View_filter_delete(make_route_config(c), c.Param("kind"), name, c.Request.PostForm)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
}
