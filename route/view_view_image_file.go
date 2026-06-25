package route

import (
	"net/http"
	"opennamu/route/tool"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

func View_view_image_file(c *gin.Context) {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    c.Header("X-Content-Type-Options", "nosniff")

    raw_path := strings.TrimPrefix(c.Param("name"), "/")
    if raw_path == "" {
        c.String(http.StatusOK, "")
        return
    }

    raw_path = strings.ReplaceAll(raw_path, "\\", "/")
    file_name := path.Base(raw_path)

    if file_name == "." || file_name == "/" || strings.ContainsAny(file_name, `/\`) {
        c.String(http.StatusBadRequest, "")
        return
    }

    re_cache := regexp.MustCompile(`\.cache_v[0-9]+$`)
    file_name = re_cache.ReplaceAllString(file_name, "")

    parts := strings.Split(file_name, ".")
    mime_type := "application/octet-stream"
    if len(parts) >= 2 {
        ext := strings.ToLower(parts[len(parts) - 1])
        switch ext {
        // Images
        case "jpeg", "jpg":
            mime_type = "image/jpeg"
        case "png":
            mime_type = "image/png"
        case "gif":
            mime_type = "image/gif"
        case "webp":
            mime_type = "image/webp"
        case "bmp":
            mime_type = "image/bmp"
        case "tif", "tiff":
            mime_type = "image/tiff"
        case "ico":
            mime_type = "image/x-icon"
        case "svg":
            mime_type = "image/svg+xml"
        case "avif":
            mime_type = "image/avif"
        case "heic":
            mime_type = "image/heic"

        // Video
        case "mp4", "m4v":
            mime_type = "video/mp4"
        case "webm":
            mime_type = "video/webm"
        case "mkv":
            mime_type = "video/x-matroska"
        case "mov":
            mime_type = "video/quicktime"
        case "avi":
            mime_type = "video/x-msvideo"
        case "mpeg", "mpg":
            mime_type = "video/mpeg"
        case "ts":
            mime_type = "video/mp2t"
        case "3gp":
            mime_type = "video/3gpp"
        case "3g2":
            mime_type = "video/3gpp2"

        // Audio
        case "mp3":
            mime_type = "audio/mpeg"
        case "wav":
            mime_type = "audio/wav"
        case "flac":
            mime_type = "audio/flac"
        case "aac":
            mime_type = "audio/aac"
        case "m4a":
            mime_type = "audio/mp4"
        case "ogg", "oga":
            mime_type = "audio/ogg"
        case "opus":
            mime_type = "audio/opus"
        case "amr":
            mime_type = "audio/amr"
        case "weba":
            mime_type = "audio/webm"

        // Docs / text / code
        case "txt":
            mime_type = "text/plain; charset=utf-8"
        case "html", "htm":
            mime_type = "text/html; charset=utf-8"
        case "css":
            mime_type = "text/css; charset=utf-8"
        case "csv":
            mime_type = "text/csv; charset=utf-8"
        case "tsv":
            mime_type = "text/tab-separated-values; charset=utf-8"
        case "js":
            mime_type = "application/javascript; charset=utf-8"
        case "json":
            mime_type = "application/json"
        case "xml":
            mime_type = "application/xml"
        case "pdf":
            mime_type = "application/pdf"
        case "wasm":
            mime_type = "application/wasm"

        // Fonts
        case "woff":
            mime_type = "font/woff"
        case "woff2":
            mime_type = "font/woff2"
        case "ttf":
            mime_type = "font/ttf"
        case "otf":
            mime_type = "font/otf"
        case "eot":
            mime_type = "application/vnd.ms-fontobject"

        // Archives / binaries
        case "zip":
            mime_type = "application/zip"
        case "7z":
            mime_type = "application/x-7z-compressed"
        case "rar":
            mime_type = "application/vnd.rar"
        case "tar":
            mime_type = "application/x-tar"
        case "gz":
            mime_type = "application/gzip"
        case "bz2":
            mime_type = "application/x-bzip2"
        case "xz":
            mime_type = "application/x-xz"
        case "zst":
            mime_type = "application/zstd"
        }
    }

    final_path := filepath.Join(tool.Get_file_main_dir(db), file_name)
    if _, err := os.Stat(final_path); err != nil {
        if os.IsNotExist(err) {
            c.String(http.StatusNotFound, "")
            return
        }

        c.String(http.StatusInternalServerError, "read error")
        return
    }

    if strings.HasPrefix(mime_type, "text/") || mime_type == "application/javascript" || strings.HasPrefix(mime_type, "application/xml") {
        c.Header("Content-Type", mime_type+"; charset=utf-8")
    } else {
        c.Header("Content-Type", mime_type)
    }

    c.File(final_path)
}