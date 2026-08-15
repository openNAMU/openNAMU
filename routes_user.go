package main

import (
	"net/http"

	"opennamu/route"
	"opennamu/route/tool"

	"github.com/gin-gonic/gin"
)

func set_skin_language_cookie(c *gin.Context) {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	language := "ko-KR"
	tool.QueryRow_DB(db, "select data from other where name = 'language'", []any{&language})
	if language == "" {
		language = "ko-KR"
	}
	user_language := language
	tool.QueryRow_DB(db, "select data from user_set where name = 'lang' and id = ?", []any{&user_language}, tool.Get_IP(c))
	if user_language == "" {
		user_language = language
	}

	http.SetCookie(c.Writer, &http.Cookie{Name: "language", Value: language, Path: "/"})
	http.SetCookie(c.Writer, &http.Cookie{Name: "user_language", Value: user_language, Path: "/"})
}

func register_user_routes(r *gin.Engine) {
	r.GET("/challenge", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_challenge(make_route_config(c), nil)))
	})
	r.POST("/challenge", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_challenge(make_route_config(c), c.Request.PostForm)))
	})

	r.GET("/change", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_setting(make_route_config(c), nil)))
	})
	r.POST("/change", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_setting(make_route_config(c), c.Request.PostForm)))
	})

	r.GET("/change/pw", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_password(make_route_config(c), nil)))
	})
	r.POST("/change/pw", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_password(make_route_config(c), c.Request.PostForm)))
	})
	r.GET("/change/key", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_key(make_route_config(c), nil)))
	})
	r.POST("/change/key", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_key(make_route_config(c), c.Request.PostForm)))
	})

	r.GET("/change/email", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_email(make_route_config(c), nil)))
	})
	r.POST("/change/email", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_email(make_route_config(c), c.Request.PostForm)))
	})
	r.GET("/change/head", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_head_skin(make_route_config(c), "", nil)))
	})
	r.POST("/change/head", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_head_skin(make_route_config(c), "", c.Request.PostForm)))
	})
	r.GET("/change/top_menu", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_top_menu(make_route_config(c), nil)))
	})
	r.POST("/change/top_menu", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_top_menu(make_route_config(c), c.Request.PostForm)))
	})
	r.GET("/change/user_name", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_name_for(make_route_config(c), "", nil)))
	})
	r.POST("/change/user_name", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_name_for(make_route_config(c), "", c.Request.PostForm)))
	})

	r.GET("/change/skin_set", func(c *gin.Context) {
		set_skin_language_cookie(c)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_skin_set(make_route_config(c))))
	})
	r.POST("/change/skin_set", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		set_skin_language_cookie(c)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_user_skin_set(make_route_config(c))))
	})
}
