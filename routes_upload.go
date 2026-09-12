package main

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"opennamu/route"
	"opennamu/route/tool"

	"github.com/gin-gonic/gin"
)

func upload_post(c *gin.Context) {
	upload_db := tool.DB_connect()
	file_max_size := tool.Get_file_max_size(upload_db)
	tool.DB_close(upload_db)
	if file_max_size <= 0 {
		file_max_size = 2
	}

	max_file_bytes := int64(file_max_size) * 1000 * 1000
	max_request_bytes := max_file_bytes*10 + 16*1000*1000
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, max_request_bytes)

	form, err := c.MultipartForm()
	if err != nil || form == nil {
		status := http.StatusBadRequest
		if err != nil && strings.Contains(err.Error(), "request body too large") {
			status = http.StatusRequestEntityTooLarge
		}
		c.String(status, "invalid multipart form")
		return
	}

	files := form.File["f_data[]"]
	if len(files) == 0 {
		c.String(http.StatusBadRequest, "no file")
		return
	}

	posted_name := strings.TrimSpace(c.PostForm("f_name"))
	license := c.PostForm("f_lice_sel")
	license_text := c.PostForm("f_lice")
	replace := c.PostForm("f_replace")
	captcha := captcha_response(c)
	upload_files := []route.Upload_file_data{}

	count := 1
	for _, fh := range files {
		f, err := fh.Open()
		if err != nil {
			continue
		}

		original_name := filepath.Base(strings.TrimSpace(fh.Filename))
		ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(original_name)), ".")
		name := strings.TrimSuffix(original_name, filepath.Ext(original_name))
		if posted_name != "" {
			name = posted_name
			if len(files) > 1 {
				name += " " + strconv.Itoa(count)
			}
		}
		name = strings.TrimSpace(name)
		ext = strings.TrimSpace(ext)

		upload_files = append(upload_files, route.Upload_file_data{
			File_name:    name,
			File_reader:  f,
			File_ext:     ext,
			License:      license,
			License_text: license_text,
			Replace:      replace,
			Captcha:      captcha,
		})
		count += 1
	}

	route_data := route.View_edit_file_upload_post(make_route_config(c), upload_files)
	status := http.StatusOK
	if _, ok := tool.Get_redirect_target(route_data); !ok {
		status = http.StatusBadRequest
	}
	write_data(c, status, "text/html; charset=utf-8", []byte(route_data))
}
