package main

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	"opennamu/route"
	"opennamu/route/tool"
	"path/filepath"
	"strings"
)

func register_history_edit_routes(r *gin.Engine) {
	r.GET("/history/*doc_name", func(c *gin.Context) {
		route_data := route.View_list_history(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), "", "1")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/history/*doc_name", func(c *gin.Context) {
		doc_name := strings.TrimPrefix(c.Param("doc_name"), "/")
		a := c.PostForm("a")
		b := c.PostForm("b")

		route_data := route.View_list_history_post(make_route_config(c), doc_name, a, b)
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/history_page/:num/:set_type/*doc_name", func(c *gin.Context) {
		route_data := route.View_list_history(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), strings.TrimPrefix(c.Param("set_type"), "/"), strings.TrimPrefix(c.Param("num"), "/"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/history_page/:num/:set_type/*doc_name", func(c *gin.Context) {
		doc_name := strings.TrimPrefix(c.Param("doc_name"), "/")
		a := c.PostForm("a")
		b := c.PostForm("b")

		route_data := route.View_list_history_post(make_route_config(c), doc_name, a, b)
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/diff/:before_rev/:after_rev/*doc_name", func(c *gin.Context) {
		route_data := route.View_w_diff(
			make_route_config(c),
			strings.TrimPrefix(c.Param("doc_name"), "/"),
			c.Param("before_rev"),
			c.Param("after_rev"),
		)
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/edit/*doc_name", func(c *gin.Context) {
		route_data := route.View_edit(make_route_config(c), strings.TrimPrefix(c.Param("doc_name"), "/"), c.Query("load"))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.POST("/edit/*doc_name", func(c *gin.Context) {
		doc_name := strings.TrimPrefix(c.Param("doc_name"), "/")
		data := c.PostForm("content")
		send := c.PostForm("send")
		agree := c.PostForm("copyright_agreement")

		route_data := route.View_edit_post(make_route_config(c), doc_name, data, send, agree)
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

		other_set_arr_str, _ := jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(other_set_arr)

		route_data := route.View_edit_file_upload_post(tool.Config{
			IP:        tool.Get_IP(c),
			Cookies:   tool.Get_Cookies(c),
			Session:   tool.Get_session(c),
			Other_set: other_set_arr_str,
		})
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(route_data))
	})

	r.GET("/view/*name", route.View_view_file)
	r.GET("/views/*name", route.View_view_file)
	r.GET("/image/*name", route.View_view_image_file)
}
