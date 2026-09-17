package main

import (
	"net/http"
	"net/url"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func Register_alarm_routes(r *gin.Engine) {
	r.GET("/alarm", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm(Make_route_config(c), "", nil)))
	})
	r.POST("/alarm", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm(Make_route_config(c), "", c.Request.PostForm)))
	})
	r.GET("/alarm/page/:num", func(c *gin.Context) {
		values := url.Values{"num": {c.Param("num")}}
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm(Make_route_config(c), "", values)))
	})
	r.POST("/alarm/read", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm_read(Make_route_config(c), "")))
	})
	r.GET("/alarm/delete", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm_delete(Make_route_config(c), "")))
	})
	r.POST("/alarm/delete", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm(Make_route_config(c), "", url.Values{"all": {"1"}})))
	})
	r.GET("/alarm/delete/:id", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm_delete(Make_route_config(c), "")))
	})
	r.POST("/alarm/delete/:id", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm(Make_route_config(c), "", url.Values{"delete": {c.Param("id")}})))
	})

	r.GET("/alarm_user/:user_name", func(c *gin.Context) {
		name := c.Param("user_name")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm(Make_route_config(c), name, nil)))
	})
	r.POST("/alarm_user/:user_name", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		name := c.Param("user_name")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm(Make_route_config(c), name, c.Request.PostForm)))
	})
	r.GET("/alarm_user/:user_name/send", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm_send(Make_route_config(c), c.Param("user_name"), nil)))
	})
	r.POST("/alarm_user/:user_name/send", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm_send(Make_route_config(c), c.Param("user_name"), c.Request.PostForm)))
	})
	r.GET("/alarm_user/:user_name/page/:num", func(c *gin.Context) {
		values := url.Values{"num": {c.Param("num")}}
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm(Make_route_config(c), c.Param("user_name"), values)))
	})
	r.POST("/alarm_user/:user_name/read", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm_read(Make_route_config(c), c.Param("user_name"))))
	})
	r.GET("/alarm_user/:user_name/delete", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm_delete(Make_route_config(c), c.Param("user_name"))))
	})
	r.POST("/alarm_user/:user_name/delete", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm(Make_route_config(c), c.Param("user_name"), url.Values{"all": {"1"}})))
	})
	r.GET("/alarm_user/:user_name/delete/:id", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm_delete(Make_route_config(c), c.Param("user_name"))))
	})
	r.POST("/alarm_user/:user_name/delete/:id", func(c *gin.Context) {
		values := url.Values{"delete": {c.Param("id")}}
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_alarm(Make_route_config(c), c.Param("user_name"), values)))
	})
}
