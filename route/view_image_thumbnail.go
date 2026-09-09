package route

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"opennamu/route/tool"

	"github.com/gin-gonic/gin"
)

const thumbnail_max_size = 1024
const thumbnail_max_pixels = int64(100000000)

func View_image_thumbnail(c *gin.Context) {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	config := tool.Config{
		IP:      tool.Get_IP(c),
		Cookies: tool.Get_Cookies(c),
		Session: tool.Get_session(c),
	}
	if !tool.Check_permission(db, "image_view", config.IP) {
		c.Data(http.StatusForbidden, "text/html; charset=utf-8", []byte(tool.Get_error_page(db, config, "auth")))
		return
	}

	size, err := strconv.Atoi(c.Param("size"))
	if err != nil || size <= 0 || size > thumbnail_max_size {
		c.String(http.StatusBadRequest, "")
		return
	}

	raw_path := strings.TrimPrefix(c.Param("name"), "/")
	raw_path = strings.ReplaceAll(raw_path, "\\", "/")
	raw_file_name := path.Base(raw_path)
	if raw_file_name == "." || raw_file_name == ".." || raw_file_name == "/" || strings.ContainsAny(raw_file_name, `/\`) {
		c.String(http.StatusBadRequest, "")
		return
	}

	file_name := view_image_file_cache_regex.ReplaceAllString(raw_file_name, "")
	if file_name == "" || file_name == "." || file_name == ".." {
		c.String(http.StatusBadRequest, "")
		return
	}

	extension := strings.ToLower(filepath.Ext(file_name))
	if extension != ".jpg" && extension != ".jpeg" && extension != ".png" && extension != ".gif" {
		redirect_image_original(c, raw_file_name)
		return
	}

	image_dir := tool.Get_file_main_dir(db)
	file_path := filepath.Join(image_dir, file_name)
	file_info, err := os.Stat(file_path)
	if err != nil || file_info.IsDir() {
		c.String(http.StatusNotFound, "")
		return
	}

	cache_key := strconv.FormatInt(file_info.ModTime().UnixNano(), 10)
	if cache_version := view_image_file_cache_regex.FindString(raw_file_name); cache_version != "" {
		cache_key = cache_version
	}
	cache_dir := filepath.Join(image_dir, ".thumbnail")
	cache_path := filepath.Join(cache_dir, fmt.Sprintf("%d_%s_%s.png", size, cache_key, file_name))
	if _, err := os.Stat(cache_path); err == nil {
		write_thumbnail_file(c, cache_path, view_image_file_cache_regex.MatchString(raw_file_name))
		return
	}

	file, err := os.Open(file_path)
	if err != nil {
		c.String(http.StatusInternalServerError, "read error")
		return
	}
	image_config, _, err := image.DecodeConfig(file)
	if err != nil {
		_ = file.Close()
		c.String(http.StatusUnsupportedMediaType, "")
		return
	}
	if image_config.Width <= 0 || image_config.Height <= 0 || int64(image_config.Width)*int64(image_config.Height) > thumbnail_max_pixels {
		_ = file.Close()
		c.String(http.StatusRequestEntityTooLarge, "")
		return
	}
	if image_config.Width <= size && image_config.Height <= size {
		_ = file.Close()
		redirect_image_original(c, raw_file_name)
		return
	}
	if _, err := file.Seek(0, 0); err != nil {
		_ = file.Close()
		c.String(http.StatusInternalServerError, "read error")
		return
	}
	source, _, err := image.Decode(file)
	_ = file.Close()
	if err != nil {
		c.String(http.StatusUnsupportedMediaType, "")
		return
	}

	width, height := thumbnail_size(image_config.Width, image_config.Height, size)
	thumbnail := resize_image(source, width, height)
	if err := os.MkdirAll(cache_dir, 0o755); err != nil {
		c.String(http.StatusInternalServerError, "cache error")
		return
	}

	temporary_file, err := os.CreateTemp(cache_dir, ".thumbnail-*.png")
	if err != nil {
		c.String(http.StatusInternalServerError, "cache error")
		return
	}
	temporary_path := temporary_file.Name()
	defer os.Remove(temporary_path)
	if err := png.Encode(temporary_file, thumbnail); err != nil {
		_ = temporary_file.Close()
		c.String(http.StatusInternalServerError, "thumbnail error")
		return
	}
	if err := temporary_file.Close(); err != nil {
		c.String(http.StatusInternalServerError, "cache error")
		return
	}
	if err := os.Rename(temporary_path, cache_path); err != nil {
		if _, stat_err := os.Stat(cache_path); stat_err != nil {
			c.String(http.StatusInternalServerError, "cache error")
			return
		}
	}

	write_thumbnail_file(c, cache_path, view_image_file_cache_regex.MatchString(raw_file_name))
}

func thumbnail_size(width int, height int, size int) (int, int) {
	if width >= height {
		result_height := int(float64(height) * float64(size) / float64(width))
		if result_height < 1 {
			result_height = 1
		}
		return size, result_height
	}

	result_width := int(float64(width) * float64(size) / float64(height))
	if result_width < 1 {
		result_width = 1
	}
	return result_width, size
}

func resize_image(source image.Image, width int, height int) image.Image {
	output := image.NewRGBA(image.Rect(0, 0, width, height))
	bound := source.Bounds()
	for y := 0; y < height; y++ {
		source_y := bound.Min.Y + y*bound.Dy()/height
		for x := 0; x < width; x++ {
			source_x := bound.Min.X + x*bound.Dx()/width
			output.Set(x, y, source.At(source_x, source_y))
		}
	}
	return output
}

func write_thumbnail_file(c *gin.Context, file_path string, immutable bool) {
	c.Header("Content-Type", "image/png")
	if immutable {
		c.Header("Cache-Control", "private, max-age=31536000, immutable")
	} else {
		c.Header("Cache-Control", "private, max-age=3600")
	}
	c.File(file_path)
}

func redirect_image_original(c *gin.Context, file_name string) {
	c.Redirect(http.StatusFound, "/image/"+tool.Url_parser(file_name))
}
