package route

import (
	"errors"
	"io/fs"
	"net/http"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"opennamu/route/tool"

	"github.com/gin-gonic/gin"
)

var view_file_cache_regex = regexp.MustCompile(`\.cache_v[0-9]+$`)

func View_view_file(c *gin.Context) {
	raw_path := strings.TrimPrefix(c.Param("name"), "/")
	if raw_path == "" {
		c.String(http.StatusOK, "")
		return
	}

	dir_name := path.Dir(raw_path)
	file_name := path.Base(raw_path)

	if file_name == "." || file_name == "/" || file_name == ".." {
		c.String(http.StatusOK, "")
		return
	}

	has_cache_version := view_file_cache_regex.MatchString(file_name)
	file_name = view_file_cache_regex.ReplaceAllString(file_name, "")

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
		ext := strings.ToLower(parts[len(parts)-1])
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

	final_path := path.Join(dir_name, file_name)
	file_data, err := tool.Read_view_file(final_path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			c.String(http.StatusOK, "")
			return
		}

		c.String(http.StatusInternalServerError, "read error")
		return
	}

	content_type := mime_type
	if strings.HasPrefix(mime_type, "image/") && mime_type != "image/svg+xml" {
		content_type = mime_type
	} else {
		content_type = mime_type + "; charset=utf-8"
	}
	if has_cache_version {
		c.Header("Cache-Control", "public, max-age=31536000, immutable")
	}

	c.Data(http.StatusOK, content_type, file_data)
}
