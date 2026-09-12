package main

import (
	"io"
	"net/http"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

const document_backup_max_size = 64 * 1024 * 1024

func register_document_backup_routes(r *gin.Engine) {
	r.GET("/backup", func(c *gin.Context) {
		data := route.View_document_backup(make_route_config(c), nil)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.POST("/backup", func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, document_backup_max_size)
		file, _, err := c.Request.FormFile("backup")
		if err != nil {
			data := route.View_document_backup(make_route_config(c), map[string]any{"response": "error"})
			write_data(c, http.StatusBadRequest, "text/html; charset=utf-8", []byte(data))
			return
		}
		defer file.Close()

		raw_data, err := io.ReadAll(file)
		if err != nil || len(raw_data) > document_backup_max_size {
			data := route.View_document_backup(make_route_config(c), map[string]any{"response": "error"})
			write_data(c, http.StatusBadRequest, "text/html; charset=utf-8", []byte(data))
			return
		}

		result := route.Api_document_backup_import(make_route_config(c), raw_data)
		data := route.View_document_backup(make_route_config(c), result)
		write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	r.GET("/backup/export", func(c *gin.Context) {
		result := route.Api_document_backup_export(make_route_config(c))
		if result["response"] != "ok" {
			data := route.View_document_backup(make_route_config(c), result)
			write_data(c, http.StatusForbidden, "text/html; charset=utf-8", []byte(data))
			return
		}

		raw_data, ok := result["data"].([]byte)
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Header("Content-Disposition", `attachment; filename="opennamu-document-backup.json"`)
		c.Data(http.StatusOK, "application/json; charset=utf-8", raw_data)
	})
}
