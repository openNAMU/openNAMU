package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"opennamu/route"
	"opennamu/route/tool"

	"net/http"

	"github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func error_handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}

				stackTrace := debug.Stack()

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"response" : "error",
					"error" : err.Error(),
					"stack" : string(stackTrace),
				})
			}
		}()

		c.Next()
	}
}

func pongo_init() {
	pongo2.RegisterFilter("md5_replace", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
		h := md5.Sum([]byte(in.String()))

		return pongo2.AsValue(hex.EncodeToString(h[:])), nil
	})

	pongo2.RegisterFilter("load_lang", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
		db := tool.DB_connect()
		defer tool.DB_close(db)

		return pongo2.AsValue(tool.Get_language(db, in.String(), false)), nil
	})

	pongo2.RegisterFilter("cut_100", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
		s := in.String()
		if len(s) > 100 {
			s = s[:100]
		}

		return pongo2.AsValue(s), nil
	})
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	port := "3001"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	var r *gin.Engine
	if len(os.Args) > 2 && os.Args[2] == "dev" {
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}

	if len(os.Args) <= 3 || os.Args[3] != "api" {
		tool.IN_mod_OUT_mod(true)
	}

	r.Use(error_handler())
	pongo_init()
	
	tool.Main_init()

	r.GET("/api/user_info/:user_name", func(c *gin.Context) {
		route_data := route.Api_user_info(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("user_name"))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/recent_change", func(c *gin.Context) {
		route_data := route.Api_list_recent_change(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "normal", "10", "1")
		in_data := route_data["data"].([][]string)
		c.JSON(http.StatusOK, in_data)
	})

	r.GET("/api/recent_change/:limit", func(c *gin.Context) {
		route_data := route.Api_list_recent_change(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "normal", c.Param("limit"), "1")
		in_data := route_data["data"].([][]string)
		c.JSON(http.StatusOK, in_data)
	})

	r.GET("/api/recent_discuss", func(c *gin.Context) {
		route_data := route.Api_list_recent_discuss(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "10", "1", "")
		in_data := route_data["data"].([][]string)
		c.JSON(http.StatusOK, in_data)
	})

	r.GET("/api/recent_discuss/:limit", func(c *gin.Context) {
		route_data := route.Api_list_recent_discuss(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("limit"), "1", "")
		in_data := route_data["data"].([][]string)
		c.JSON(http.StatusOK, in_data)
	})

	r.POST("/api/v2/lang", func(c *gin.Context) {
		data := c.PostForm("data")
		safe := c.PostForm("safe")
		legacy := c.PostForm("legacy")

		route_data := route.Api_func_language(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, data, safe, legacy)
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/ip/:ip", func(c *gin.Context) {
		route_data := route.Api_func_ip(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("ip"))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/ip_menu/:ip", func(c *gin.Context) {
		route_data := route.Api_func_ip_menu(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("ip"), "")
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/user/setting/editor", func(c *gin.Context) {
		route_data := route.Api_user_setting_editor(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		})
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/page_view/*doc_name", func(c *gin.Context) {
		route_data := route.Api_w_page_view(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("doc_name"), "/"))
		c.JSON(http.StatusOK, route_data)
	})

	r.POST("/api/v2/user/setting/editor", func(c *gin.Context) {
		route_data := route.Api_user_setting_editor_post(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Request.FormValue("data"))
		c.JSON(http.StatusOK, route_data)
	})

	r.DELETE("/api/v2/user/setting/editor", func(c *gin.Context) {
		route_data := route.Api_user_setting_editor_delete(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Request.FormValue("data"))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/page_view_post/*doc_name", func(c *gin.Context) {
		route_data := route.Api_w_page_view_post(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("doc_name"), "/"))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/bbs/w/page_view_post/:set_id/:set_code", func(c *gin.Context) {
		route_data := route.Api_bbs_w_page_view_post(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("set_id"), "/"), strings.TrimPrefix(c.Param("set_code"), "/"))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/bbs/w/page_view/:set_id/:set_code", func(c *gin.Context) {
		route_data := route.Api_bbs_w_page_view(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("set_id"), "/"), strings.TrimPrefix(c.Param("set_code"), "/"))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/bbs/main", func(c *gin.Context) {
		route_data := route.Api_bbs(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "", "1")
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/raw/*doc_name", func(c *gin.Context) {
		route_data := route.Api_w_raw(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("doc_name"), "/"), "", "")
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/api/v2/raw_exist/*doc_name", func(c *gin.Context) {
		route_data := route.Api_w_raw(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("doc_name"), "/"), "true", "")
		c.JSON(http.StatusOK, route_data)
	})


	r.GET("/api/v2/raw_rev/:rev/*doc_name", func(c *gin.Context) {
		route_data := route.Api_w_raw(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("doc_name"), "/"), "", c.Param("rev"))
		c.JSON(http.StatusOK, route_data)
	})

	r.GET("/watch_list", func(c *gin.Context) {
		route_data := route.View_user_watch_list(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "1", "watchlist")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/star_doc", func(c *gin.Context) {
		route_data := route.View_user_watch_list(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "1", "star_doc")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/doc_watch_list/:count/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_watch_list(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("doc_name"), "/"), strings.TrimPrefix(c.Param("count"), "/"), "watchlist")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/doc_star_doc/:count/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_watch_list(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("doc_name"), "/"), strings.TrimPrefix(c.Param("count"), "/"), "star_doc")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/star_doc_from/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_watch_list_add_post(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("doc_name"), "/"), "star_doc_from")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/star_doc_from/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_watch_list_add(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("doc_name"), "/"), "star_doc_from")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/star_doc/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_watch_list_add_post(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("doc_name"), "/"), "star_doc")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/star_doc/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_watch_list_add(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("doc_name"), "/"), "star_doc")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/watch_list_from/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_watch_list_add_post(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("doc_name"), "/"), "watchlist_from")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/watch_list_from/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_watch_list_add(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("doc_name"), "/"), "watchlist_from")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/watch_list/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_watch_list_add_post(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("doc_name"), "/"), "watchlist")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/watch_list/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_watch_list_add(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("doc_name"), "/"), "watchlist")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/record_bbs/:user_name", func(c *gin.Context) {
		route_data := route.View_record_bbs(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("user_name"), "1")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/record_bbs/:user_name/:page", func(c *gin.Context) {
		route_data := route.View_record_bbs(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("user_name"), c.Param("page"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/record_bbs_in/:set_id/:user_name", func(c *gin.Context) {
		route_data := route.View_record_bbs_in(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("user_name"), c.Param("set_id"), "1")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/record_bbs_in/:set_id/:user_name/:page", func(c *gin.Context) {
		route_data := route.View_record_bbs_in(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("user_name"), c.Param("set_id"), c.Param("page"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/record_bbs_comment/:user_name", func(c *gin.Context) {
		route_data := route.View_record_bbs_comment(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("user_name"), "1")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/record_bbs_comment/:user_name/:page", func(c *gin.Context) {
		route_data := route.View_record_bbs_comment(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("user_name"), c.Param("page"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/record_bbs_comment_in/:set_id/:user_name", func(c *gin.Context) {
		route_data := route.View_record_bbs_comment_in(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("user_name"), c.Param("set_id"), "1")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/record_bbs_comment_in/:set_id/:user_name/:page", func(c *gin.Context) {
		route_data := route.View_record_bbs_comment_in(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("user_name"), c.Param("set_id"), c.Param("page"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/list/random", func(c *gin.Context) {
		route_data := route.View_list_random(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		})
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/list/document/old", func(c *gin.Context) {
		route_data := route.View_list_old_page(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "1", "old")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/list/document/old/:num", func(c *gin.Context) {
		route_data := route.View_list_old_page(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("num"), "old")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/list/document/new", func(c *gin.Context) {
		route_data := route.View_list_old_page(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "1", "new")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/list/document/new/:num", func(c *gin.Context) {
		route_data := route.View_list_old_page(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("num"), "new")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/list/document/long", func(c *gin.Context) {
		route_data := route.View_list_long_page(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "1", "long")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/list/document/long/:num", func(c *gin.Context) {
		route_data := route.View_list_long_page(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("num"), "long")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/list/document/short", func(c *gin.Context) {
		route_data := route.View_list_long_page(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "1", "short")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/list/document/short/:num", func(c *gin.Context) {
		route_data := route.View_list_long_page(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("num"), "short")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/random", func(c *gin.Context) {
		route_data := route.View_w_random(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		})
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/other", func(c *gin.Context) {
		route_data := route.View_main_other(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		})
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/user", func(c *gin.Context) {
		route_data := route.View_user(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, tool.Get_IP(c))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/user/*user_name", func(c *gin.Context) {
		route_data := route.View_user(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("user_name"), "/"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/upload", func(c *gin.Context) {
		route_data := route.View_edit_file_upload(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		})
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/w/*doc_name", func(c *gin.Context) {
		route_data, status_code := route.View_w(c, tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("doc_name"), "/"))
		c.Data(status_code, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/down/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_down(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("doc_name"), "/"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/raw/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_raw(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("doc_name"), "/"), "", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/raw_rev/:rev/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_raw(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("doc_name"), "/"), c.Param("rev"), "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/raw_acl/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_raw(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("doc_name"), "/"), "", "document_acl")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/topic/*doc_name", func(c *gin.Context) {
		route_data := route.View_topic_list(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("doc_name"), "/"), "", "1")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/topic_page/:num/*doc_name", func(c *gin.Context) {
		route_data := route.View_topic_list(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("doc_name"), "/"), "", c.Param("num"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/topic_close/:num/*doc_name", func(c *gin.Context) {
		route_data := route.View_topic_list(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("doc_name"), "/"), "close", c.Param("num"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/topic_agree/:num/*doc_name", func(c *gin.Context) {
		route_data := route.View_topic_list(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("doc_name"), "/"), "agree", c.Param("num"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block", func(c *gin.Context) {
		route_data := route.View_list_recent_block(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "1", "", "", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/all", func(c *gin.Context) {
		route_data := route.View_list_recent_block(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "1", "", "", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/all/:num", func(c *gin.Context) {
		route_data := route.View_list_recent_block(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("num"), "", "", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/all/:num/*why", func(c *gin.Context) {
		route_data := route.View_list_recent_block(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("num"), "", strings.TrimPrefix(c.Param("why"), "/"), "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/user/:user_name", func(c *gin.Context) {
		route_data := route.View_list_recent_block(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "1", "user", "", c.Param("user_name"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/user/:user_name/:num", func(c *gin.Context) {
		route_data := route.View_list_recent_block(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("num"), "user", "", c.Param("user_name"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/admin/:user_name", func(c *gin.Context) {
		route_data := route.View_list_recent_block(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "1", "admin", "", c.Param("user_name"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/admin/:user_name/:num", func(c *gin.Context) {
		route_data := route.View_list_recent_block(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("num"), "admin", "", c.Param("user_name"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/regex", func(c *gin.Context) {
		route_data := route.View_list_recent_block(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "1", "regex", "", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/regex/:num", func(c *gin.Context) {
		route_data := route.View_list_recent_block(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("num"), "regex", "", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/cidr", func(c *gin.Context) {
		route_data := route.View_list_recent_block(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "1", "cidr", "", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/cidr/:num", func(c *gin.Context) {
		route_data := route.View_list_recent_block(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("num"), "cidr", "", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/private", func(c *gin.Context) {
		route_data := route.View_list_recent_block(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "1", "private", "", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/private/:num", func(c *gin.Context) {
		route_data := route.View_list_recent_block(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("num"), "private", "", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/ongoing", func(c *gin.Context) {
		route_data := route.View_list_recent_block(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "1", "ongoing", "", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_block/ongoing/:num", func(c *gin.Context) {
		route_data := route.View_list_recent_block(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("num"), "ongoing", "", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_change", func(c *gin.Context) {
		route_data := route.View_list_recent_change(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "", "50", "1")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_changes", func(c *gin.Context) {
		route_data := route.View_list_recent_change(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "", "50", "1")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_change/:num/:set_type", func(c *gin.Context) {
		route_data := route.View_list_recent_change(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("set_type"), "50", c.Param("num"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_discuss", func(c *gin.Context) {
		route_data := route.View_list_recent_discuss(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "50", "1", "")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/recent_discuss/:num/:set_type", func(c *gin.Context) {
		route_data := route.View_list_recent_discuss(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "50", c.Param("num"), c.Param("set_type"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/main", func(c *gin.Context) {
		route_data := route.View_bbs_main(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "1")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/make", func(c *gin.Context) {
		route_data := route.View_bbs_make(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		})
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/bbs/make", func(c *gin.Context) {
		bbs_name := c.PostForm("bbs_name")
		bbs_type := c.PostForm("bbs_type")

		route_data := route.View_bbs_make_post(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, bbs_name, bbs_type)
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/in/:set_id", func(c *gin.Context) {
		route_data := route.View_bbs_in(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("set_id"), "1")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/in/:set_id/:page_num", func(c *gin.Context) {
		route_data := route.View_bbs_in(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("set_id"), c.Param("page_num"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/bbs/w/:set_id/:set_code", func(c *gin.Context) {
		route_data := route.View_bbs_in_w(c, tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Param("set_id"), c.Param("set_code"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/history/*doc_name", func(c *gin.Context) {
		route_data := route.View_list_history(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("doc_name"), "/"), "", "1")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/history/*doc_name", func(c *gin.Context) {
		doc_name := strings.TrimPrefix(c.Param("doc_name"), "/")
		a := c.PostForm("a")
		b := c.PostForm("b")

		route_data := route.View_list_history_post(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, doc_name, a, b)
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/history_page/:num/:set_type/*doc_name", func(c *gin.Context) {
		route_data := route.View_list_history(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("doc_name"), "/"), strings.TrimPrefix(c.Param("set_type"), "/"), strings.TrimPrefix(c.Param("num"), "/"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/history_page/:num/:set_type/*doc_name", func(c *gin.Context) {
		doc_name := strings.TrimPrefix(c.Param("doc_name"), "/")
		a := c.PostForm("a")
		b := c.PostForm("b")

		route_data := route.View_list_history_post(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, doc_name, a, b)
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/edit/*doc_name", func(c *gin.Context) {
		route_data := route.View_edit(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("doc_name"), "/"), c.Query("load"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/edit/*doc_name", func(c *gin.Context) {
		doc_name := strings.TrimPrefix(c.Param("doc_name"), "/")
		data := c.PostForm("content")
		send := c.PostForm("send")
		agree := c.PostForm("copyright_agreement")

		route_data := route.View_edit_post(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, doc_name, data, send, agree)
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/upload", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil || form == nil {
			c.String(http.StatusBadRequest, "invalid multipart form")
			return
		}

		files := form.File["f_data[]"]
		if len(files) == 0 {
			c.String(http.StatusBadRequest, "no file")
			return
		}

		posted_name := strings.TrimSpace(c.PostForm("f_name"))
		other_set_arr := []map[string]string{}

		count := 1
		for _, fh := range files {
			f, err := fh.Open()
			if err != nil {
				continue
			}

			b, err := io.ReadAll(f)

			_ = f.Close()
			if err != nil {
				continue
			}

			name := posted_name

			name = strings.TrimSpace(name)
			ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(name)), ".")
			ext = strings.TrimSpace(ext)

			b64 := base64.StdEncoding.EncodeToString(b)

			other_set := map[string]string{
				"file_name": name,
				"file_ext":  ext,
				"file_data": b64,
			}

			other_set_arr = append(other_set_arr, other_set)
			count += 1
		}

		other_set_arr_str, _ := json.MarshalToString(other_set_arr)

		route_data := route.View_edit_file_upload_post(tool.Config{
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
			Other_set: other_set_arr_str,
		})
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/view/*name", route.View_view_file)
	r.GET("/views/*name", route.View_view_file)
	r.GET("/image/*name", route.View_view_image_file)

	r.POST("/goto", func(c *gin.Context) {
		route_data := route.View_main_search_post(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "", true, c.PostForm("search"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/goto/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search_post(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "", true, c.PostForm("search"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/search", func(c *gin.Context) {
		route_data := route.View_main_search_post(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "", false, c.PostForm("search"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/search/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("keyword"), "/"), "1", "title")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/search/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search_post(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "", false, c.PostForm("search"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/search_page/:num/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search_post(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "", false, c.PostForm("search"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/search_page/:num/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("keyword"), "/"), c.Param("num"), "title")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/search_data/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search_post(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "", false, c.PostForm("search"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/search_data_page/:num/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search_post(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, "", false, c.PostForm("search"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/search_data_page/:num/*keyword", func(c *gin.Context) {
		route_data := route.View_main_search(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, strings.TrimPrefix(c.Param("keyword"), "/"), c.Param("num"), "data")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.NoRoute(func(c *gin.Context) {
		route_data := route.View_main_404_page(tool.Config{
			Other_set: "",
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   "",
		}, c.Request.URL.Path)
		c.Data(http.StatusNotFound, "text/html; charset=utf-8", []byte(route_data))
	})

	log.Default().Println("Run in http://127.0.0.1:" + port)
	if err := r.Run("0.0.0.0:" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
