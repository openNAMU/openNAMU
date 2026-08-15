package main

import (
	"net/http"
	"net/url"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_alarm_routes(r *gin.Engine) {
	r.GET("/alarm", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm(make_route_config(c), "", c.Request.URL.Query())))
	})
	r.POST("/alarm", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm(make_route_config(c), "", c.Request.PostForm)))
	})
	r.GET("/alarm/delete", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm(make_route_config(c), "", url.Values{"all": {"1"}})))
	})
	r.POST("/alarm/delete", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm(make_route_config(c), "", url.Values{"all": {"1"}})))
	})
	r.GET("/alarm/delete/:id", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm(make_route_config(c), "", url.Values{"delete": {c.Param("id")}})))
	})
	r.POST("/alarm/delete/:id", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm(make_route_config(c), "", url.Values{"delete": {c.Param("id")}})))
	})
	r.GET("/alarm/:user_name", func(c *gin.Context) {
		name := c.Param("user_name")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm(make_route_config(c), name, c.Request.URL.Query())))
	})
	r.POST("/alarm/:user_name", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		name := c.Param("user_name")
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm(make_route_config(c), name, c.Request.PostForm)))
	})
	r.GET("/alarm/:user_name/delete", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm(make_route_config(c), c.Param("user_name"), url.Values{"all": {"1"}})))
	})
	r.POST("/alarm/:user_name/delete", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm(make_route_config(c), c.Param("user_name"), url.Values{"all": {"1"}})))
	})
	r.GET("/alarm/:user_name/delete/:id", func(c *gin.Context) {
		values := url.Values{"delete": {c.Param("id")}}
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm(make_route_config(c), c.Param("user_name"), values)))
	})
	r.POST("/alarm/:user_name/delete/:id", func(c *gin.Context) {
		values := url.Values{"delete": {c.Param("id")}}
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm(make_route_config(c), c.Param("user_name"), values)))
	})
}
