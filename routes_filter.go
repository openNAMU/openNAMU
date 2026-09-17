package main

import (
	"net/http"
	"strings"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func Register_filter_routes(r *gin.Engine) {
	r.GET("/filter/:kind", func(c *gin.Context) {
		data := route.View_filter(Make_route_config(c), c.Param("kind"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/filter/:kind/add", func(c *gin.Context) {
		data := route.View_filter_add(Make_route_config(c), c.Param("kind"), "", nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/filter/:kind/add", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := route.View_filter_add(Make_route_config(c), c.Param("kind"), "", c.Request.PostForm)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/filter/:kind/add/*name", func(c *gin.Context) {
		name := strings.TrimPrefix(c.Param("name"), "/")
		data := route.View_filter_add(Make_route_config(c), c.Param("kind"), name, nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/filter/:kind/add/*name", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		name := strings.TrimPrefix(c.Param("name"), "/")
		data := route.View_filter_add(Make_route_config(c), c.Param("kind"), name, c.Request.PostForm)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/filter/:kind/del/*name", func(c *gin.Context) {
		name := strings.TrimPrefix(c.Param("name"), "/")
		data := route.View_filter_delete(Make_route_config(c), c.Param("kind"), name, nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/filter/:kind/del/*name", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		name := strings.TrimPrefix(c.Param("name"), "/")
		data := route.View_filter_delete(Make_route_config(c), c.Param("kind"), name, c.Request.PostForm)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
}
