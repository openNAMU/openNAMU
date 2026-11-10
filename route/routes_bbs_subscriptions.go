package route

import (
	"net/http"

	"opennamu/route/tool"

	"github.com/gin-gonic/gin"
)

func Register_bbs_subscription_routes(r *gin.Engine) {
	r.GET("/bbs_bookmark_list", func(c *gin.Context) {
		data := View_user_watch_list(Make_route_config(c), "1", "bbs_bookmark")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/bbs_bookmark_list_page/:num", func(c *gin.Context) {
		data := View_user_watch_list(Make_route_config(c), c.Param("num"), "bbs_bookmark")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/bbs_board_watch_list", func(c *gin.Context) {
		data := View_user_watch_list(Make_route_config(c), "1", "bbs_board_watchlist")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/bbs_board_watch_list_page/:num", func(c *gin.Context) {
		data := View_user_watch_list(Make_route_config(c), c.Param("num"), "bbs_board_watchlist")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/bbs_tag_watch_list", func(c *gin.Context) {
		data := View_user_watch_list(Make_route_config(c), "1", "bbs_tag_watchlist")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/bbs_tag_watch_list_page/:num", func(c *gin.Context) {
		data := View_user_watch_list(Make_route_config(c), c.Param("num"), "bbs_tag_watchlist")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/bbs_bookmark/:set_id/:set_code", func(c *gin.Context) {
		data := View_bbs_bookmark(Make_route_config(c), c.Param("set_id"), c.Param("set_code"), nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.POST("/bbs_bookmark/:set_id/:set_code", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := View_bbs_bookmark(Make_route_config(c), c.Param("set_id"), c.Param("set_code"), c.Request.PostForm)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/bbs_board_watch/:set_id", func(c *gin.Context) {
		data := View_bbs_board_watch(Make_route_config(c), c.Param("set_id"), nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.POST("/bbs_board_watch/:set_id", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := View_bbs_board_watch(Make_route_config(c), c.Param("set_id"), c.Request.PostForm)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/bbs_tag_watch/:set_id", func(c *gin.Context) {
		data := View_bbs_tag_watch(Make_route_config(c), "", c.Param("set_id"), nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.POST("/bbs_tag_watch/:set_id", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := View_bbs_tag_watch(Make_route_config(c), "", c.Param("set_id"), c.Request.PostForm)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/bbs_tag_watch/:set_id/:tag_token", func(c *gin.Context) {
		data := View_bbs_tag_watch(Make_route_config(c), c.Param("set_id"), c.Param("tag_token"), nil)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.POST("/bbs_tag_watch/:set_id/:tag_token", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := View_bbs_tag_watch(Make_route_config(c), c.Param("set_id"), c.Param("tag_token"), c.Request.PostForm)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
}

func Register_bbs_subscription_api_routes(r *gin.Engine) {
	r.POST("/api/v2/bbs_bookmark/:set_id/:set_code", func(c *gin.Context) {
		c.JSON(http.StatusOK, Api_bbs_bookmark_post(Make_route_config(c), c.Param("set_id"), c.Param("set_code")))
	})

	r.POST("/api/v2/bbs_board_watch/:set_id", func(c *gin.Context) {
		c.JSON(http.StatusOK, Api_bbs_board_watch_post(Make_route_config(c), c.Param("set_id")))
	})

	r.POST("/api/v2/bbs_tag_watch/:set_id", func(c *gin.Context) {
		tag, err := tool.Get_base64_decode(c.Param("set_id"))
		if err != nil {
			c.JSON(http.StatusOK, map[string]any{"response": "error", "data": "tag"})
			return
		}
		c.JSON(http.StatusOK, Api_bbs_tag_watch_post(Make_route_config(c), "", tag))
	})

	r.POST("/api/v2/bbs_tag_watch/:set_id/:tag_token", func(c *gin.Context) {
		tag, err := tool.Get_base64_decode(c.Param("tag_token"))
		if err != nil {
			c.JSON(http.StatusOK, map[string]any{"response": "error", "data": "tag"})
			return
		}
		c.JSON(http.StatusOK, Api_bbs_tag_watch_post(Make_route_config(c), c.Param("set_id"), tag))
	})
}

func Register_bbs_post_create_route(r *gin.Engine) {
	r.POST("/bbs/edit/:set_id", func(c *gin.Context) {
		config := Make_route_config(c)
		captcha_response := tool.Captcha_response(c.PostForm("g-recaptcha"), c.PostForm("g-recaptcha-response"), c.PostForm("h-captcha-response"), c.PostForm("cf-turnstile-response"), c.PostForm("altcha"))
		data := View_bbs_edit_post(
			config,
			c.Param("set_id"),
			"",
			"",
			c.PostForm("title"),
			c.PostForm("content"),
			c.PostForm("prefix"),
			c.PostForm("tags"),
			captcha_response,
			"",
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
}
