package main

import (
	"net/http"
	"strconv"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func register_setting_routes(r *gin.Engine) {
	r.GET("/setting", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting(make_route_config(c))))
	})

	r.GET("/setting/main", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_main(make_route_config(c))))
	})
	r.POST("/setting/main", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_main_post(make_route_config(c), setting_form(c))))
	})

	r.GET("/setting/rankup", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_rankup(make_route_config(c))))
	})
	r.POST("/setting/rankup", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_rankup_post(make_route_config(c), setting_form(c))))
	})

	r.GET("/setting/main/logo", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_main_logo(make_route_config(c))))
	})
	r.POST("/setting/main/logo", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_main_logo_post(make_route_config(c), setting_form(c))))
	})

	r.GET("/setting/phrase", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_phrase(make_route_config(c))))
	})
	r.POST("/setting/phrase", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_phrase_post(make_route_config(c), setting_form(c))))
	})

	r.GET("/setting/top_menu", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_top_menu(make_route_config(c))))
	})
	r.POST("/setting/top_menu", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_top_menu_post(make_route_config(c), setting_form(c))))
	})

	r.GET("/setting/head", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_head(make_route_config(c), "head", "")))
	})
	r.GET("/setting/head/:skin_name", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_head(make_route_config(c), "head", c.Param("skin_name"))))
	})
	r.POST("/setting/head", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_head_post(make_route_config(c), "head", "", c.PostForm("content"))))
	})
	r.POST("/setting/head/:skin_name", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_head_post(make_route_config(c), "head", c.Param("skin_name"), c.PostForm("content"))))
	})

	r.GET("/setting/body/top", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_head(make_route_config(c), "body/top", "")))
	})
	r.POST("/setting/body/top", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_head_post(make_route_config(c), "body/top", "", c.PostForm("content"))))
	})
	r.POST("/setting_preview/body/top", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_head_preview(make_route_config(c), "body/top", c.PostForm("content"))))
	})

	r.GET("/setting/body/bottom", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_head(make_route_config(c), "body/bottom", "")))
	})
	r.POST("/setting/body/bottom", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_head_post(make_route_config(c), "body/bottom", "", c.PostForm("content"))))
	})
	r.POST("/setting_preview/body/bottom", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_head_preview(make_route_config(c), "body/bottom", c.PostForm("content"))))
	})

	r.GET("/setting/robot", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_robot(make_route_config(c))))
	})
	r.POST("/setting/robot", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_robot_post(make_route_config(c), setting_form(c))))
	})

	r.GET("/setting/external", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_external(make_route_config(c))))
	})
	r.POST("/setting/external", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_external_post(make_route_config(c), setting_form(c))))
	})

	r.GET("/setting/sitemap", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_sitemap(make_route_config(c))))
	})
	r.POST("/setting/sitemap", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_sitemap_post(make_route_config(c))))
	})

	r.GET("/setting/sitemap_set", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_sitemap_set(make_route_config(c))))
	})
	r.POST("/setting/sitemap_set", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_sitemap_set_post(make_route_config(c), setting_form(c))))
	})

	r.GET("/setting/skin_set", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_skin_set(make_route_config(c))))
	})
	r.POST("/setting/skin_set", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_skin_set_post(make_route_config(c), setting_form(c))))
	})

	r.GET("/setting/404_page", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_404_page(make_route_config(c))))
	})
	r.POST("/setting/404_page", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_404_page_post(make_route_config(c), setting_form(c))))
	})
	r.GET("/setting/email_test", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_email_test(make_route_config(c))))
	})
	r.POST("/setting/email_test", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_email_test_post(make_route_config(c), setting_form(c))))
	})

	r.GET("/setting/backlink_reset", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_backlink_reset(make_route_config(c))))
	})
	r.POST("/setting/backlink_reset", func(c *gin.Context) {
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(route.View_setting_backlink_reset_post(make_route_config(c))))
	})

	r.GET("/api/v2/setting/:set_name", func(c *gin.Context) {
		data := route.Api_setting(make_route_config(c), c.Param("set_name"), "")
		c.JSON(http.StatusOK, data)
	})
	r.GET("/api/v2/setting/:set_name/:coverage", func(c *gin.Context) {
		data := route.Api_setting(make_route_config(c), c.Param("set_name"), c.Param("coverage"))
		c.JSON(http.StatusOK, data)
	})
	r.PUT("/api/v2/setting/:set_name", func(c *gin.Context) {
		data := route.Api_setting_put(make_route_config(c), c.Param("set_name"), c.PostForm("data"), c.PostForm("coverage"))
		c.JSON(http.StatusOK, data)
	})
	r.DELETE("/api/v2/setting/:set_name", func(c *gin.Context) {
		data := route.Api_setting_delete(make_route_config(c), c.Param("set_name"))
		c.JSON(http.StatusOK, data)
	})

	r.GET("/robots.txt", func(c *gin.Context) {
		data := route.View_robots_txt(make_route_config(c), c.Request.Host)
		write_data(c, http.StatusOK, "text/plain; charset=utf-8", []byte(data))
	})
	r.GET("/sitemap.xml", func(c *gin.Context) {
		data, ok := route.Read_sitemap_file("sitemap.xml")
		if !ok {
			c.Status(http.StatusNotFound)
			return
		}
		write_data(c, http.StatusOK, "application/xml; charset=utf-8", []byte(data))
	})
	r.GET("/sitemap_:num.xml", func(c *gin.Context) {
		num, err := strconv.Atoi(c.Param("num"))
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		name := "sitemap_" + strconv.Itoa(num) + ".xml"
		data, ok := route.Read_sitemap_file(name)
		if !ok {
			c.Status(http.StatusNotFound)
			return
		}
		write_data(c, http.StatusOK, "application/xml; charset=utf-8", []byte(data))
	})
}

func setting_form(c *gin.Context) map[string]string {
	_ = c.Request.ParseForm()
	form := map[string]string{}
	for name, values := range c.Request.PostForm {
		if len(values) > 0 {
			form[name] = values[0]
		}
	}

	return form
}
