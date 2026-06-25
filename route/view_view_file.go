package route

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

func View_view_file(c *gin.Context) {
    raw_path := strings.TrimPrefix(c.Param("name"), "/")
    if raw_path == "" {
        c.String(http.StatusOK, "")
        return
    }

    dir_name := path.Dir(raw_path)
    file_name := path.Base(raw_path)

    if file_name == "." || file_name == "/" {
        c.String(http.StatusOK, "")
        return
    }

    re_cache := regexp.MustCompile(`\.cache_v[0-9]+$`)
    file_name = re_cache.ReplaceAllString(file_name, "")

    re_dots := regexp.MustCompile(`\.{2,}`)
    dir_name = re_dots.ReplaceAllString(dir_name, "")
    dir_name = filepath.ToSlash(filepath.Clean(dir_name))
    
    if strings.HasPrefix(dir_name, "../") || strings.Contains(dir_name, "/../") {
        c.String(http.StatusBadRequest, "bad path")
        return
    }

    parts := strings.Split(file_name, ".")
    mime_type := "text/plain"
    if len(parts) >= 2 {
        ext := strings.ToLower(parts[len(parts) - 1])
        switch ext {
        case "jpeg", "jpg", "gif", "png", "webp", "ico":
            mime_type = "image/" + ext
        case "svg":
            mime_type = "image/svg+xml"
        case "js":
            mime_type = "text/javascript"
        case "txt":
            mime_type = "text/plain"
        default:
            mime_type = "text/" + ext
        }
    }

    final_path := filepath.Join("..", "views", dir_name, file_name)
    if _, err := os.Stat(final_path); err != nil {
        if os.IsNotExist(err) {
            c.String(http.StatusOK, "")
            return
        }

        c.String(http.StatusInternalServerError, "read error")
        return
    }

    if strings.HasPrefix(mime_type, "image/") && mime_type != "image/svg+xml" {
        c.Header("Content-Type", mime_type)
    } else {
        c.Header("Content-Type", mime_type+"; charset=utf-8")
    }

    c.File(final_path)
}