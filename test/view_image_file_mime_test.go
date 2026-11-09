package test

import (
	"mime"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"opennamu/route"
	"opennamu/route/tool"
)

func Test_View_view_image_file_mime_and_response_contract(t *testing.T) {
	db := apiSlowDB(t)
	defer tool.DB_close(db)

	project_root, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("find project root: %v", err)
	}
	temporary_root := filepath.Join(project_root, "tmp")
	if err := os.MkdirAll(temporary_root, 0755); err != nil {
		t.Fatalf("create project tmp directory: %v", err)
	}
	image_dir, err := os.MkdirTemp(temporary_root, "phase9_image_mime_")
	if err != nil {
		t.Fatalf("create temporary image directory: %v", err)
	}
	defer os.RemoveAll(image_dir)
	if _, err := db.Exec(`insert into other (name, data, coverage) values ('image_where', ?, '')`, image_dir); err != nil {
		t.Fatalf("configure image directory: %v", err)
	}

	mime_cases := []struct {
		extension string
		mime_type string
		inline    bool
	}{
		{extension: "jpeg", mime_type: "image/jpeg", inline: true},
		{extension: "jpg", mime_type: "image/jpeg", inline: true},
		{extension: "png", mime_type: "image/png", inline: true},
		{extension: "gif", mime_type: "image/gif", inline: true},
		{extension: "webp", mime_type: "image/webp", inline: true},
		{extension: "bmp", mime_type: "image/bmp", inline: true},
		{extension: "tif", mime_type: "image/tiff", inline: true},
		{extension: "tiff", mime_type: "image/tiff", inline: true},
		{extension: "ico", mime_type: "image/x-icon", inline: true},
		{extension: "svg", mime_type: "image/svg+xml"},
		{extension: "avif", mime_type: "image/avif", inline: true},
		{extension: "heic", mime_type: "image/heic", inline: true},
		{extension: "mp4", mime_type: "video/mp4", inline: true},
		{extension: "m4v", mime_type: "video/mp4", inline: true},
		{extension: "webm", mime_type: "video/webm", inline: true},
		{extension: "mkv", mime_type: "video/x-matroska", inline: true},
		{extension: "mov", mime_type: "video/quicktime", inline: true},
		{extension: "avi", mime_type: "video/x-msvideo", inline: true},
		{extension: "mpeg", mime_type: "video/mpeg", inline: true},
		{extension: "mpg", mime_type: "video/mpeg", inline: true},
		{extension: "ts", mime_type: "video/mp2t", inline: true},
		{extension: "3gp", mime_type: "video/3gpp", inline: true},
		{extension: "3g2", mime_type: "video/3gpp2", inline: true},
		{extension: "mp3", mime_type: "audio/mpeg", inline: true},
		{extension: "wav", mime_type: "audio/wav", inline: true},
		{extension: "flac", mime_type: "audio/flac", inline: true},
		{extension: "aac", mime_type: "audio/aac", inline: true},
		{extension: "m4a", mime_type: "audio/mp4", inline: true},
		{extension: "ogg", mime_type: "audio/ogg", inline: true},
		{extension: "oga", mime_type: "audio/ogg", inline: true},
		{extension: "opus", mime_type: "audio/opus", inline: true},
		{extension: "amr", mime_type: "audio/amr", inline: true},
		{extension: "weba", mime_type: "audio/webm", inline: true},
		{extension: "txt", mime_type: "text/plain; charset=utf-8"},
		{extension: "html", mime_type: "text/html; charset=utf-8"},
		{extension: "htm", mime_type: "text/html; charset=utf-8"},
		{extension: "css", mime_type: "text/css; charset=utf-8"},
		{extension: "csv", mime_type: "text/csv; charset=utf-8"},
		{extension: "tsv", mime_type: "text/tab-separated-values; charset=utf-8"},
		{extension: "js", mime_type: "application/javascript; charset=utf-8"},
		{extension: "json", mime_type: "application/json"},
		{extension: "xml", mime_type: "application/xml"},
		{extension: "pdf", mime_type: "application/pdf", inline: true},
		{extension: "wasm", mime_type: "application/wasm"},
		{extension: "woff", mime_type: "font/woff"},
		{extension: "woff2", mime_type: "font/woff2"},
		{extension: "ttf", mime_type: "font/ttf"},
		{extension: "otf", mime_type: "font/otf"},
		{extension: "eot", mime_type: "application/vnd.ms-fontobject"},
		{extension: "zip", mime_type: "application/zip"},
		{extension: "7z", mime_type: "application/x-7z-compressed"},
		{extension: "rar", mime_type: "application/vnd.rar"},
		{extension: "tar", mime_type: "application/x-tar"},
		{extension: "gz", mime_type: "application/gzip"},
		{extension: "bz2", mime_type: "application/x-bzip2"},
		{extension: "xz", mime_type: "application/x-xz"},
		{extension: "zst", mime_type: "application/zstd"},
		{extension: "unknown", mime_type: "application/octet-stream"},
		{extension: "UPPER.PNG", mime_type: "image/png", inline: true},
	}

	router := gin.New()
	router.Use(tool.Session_middleware())
	router.GET("/image/*name", route.View_view_image_file)
	make_request := func(name string, remote_address string) *httptest.ResponseRecorder {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/image/"+name, nil)
		request.RemoteAddr = remote_address
		router.ServeHTTP(recorder, request)
		return recorder
	}

	for _, test_case := range mime_cases {
		t.Run(test_case.extension, func(t *testing.T) {
			file_name := "phase9." + test_case.extension
			if test_case.extension == "UPPER.PNG" {
				file_name = "phase9.UPPER.PNG"
			}
			file_body := []byte("file body:" + test_case.extension)
			if err := os.WriteFile(filepath.Join(image_dir, file_name), file_body, 0600); err != nil {
				t.Fatalf("write fixture %s: %v", file_name, err)
			}

			response := make_request(file_name, "192.0.2.88:1234")
			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d: %s", response.Code, http.StatusOK, response.Body.String())
			}
			want_mime := test_case.mime_type
			want_disposition := ""
			if !test_case.inline {
				if test_case.mime_type != "image/svg+xml" {
					want_mime = "application/octet-stream"
				}
				want_disposition = mime.FormatMediaType("attachment", map[string]string{"filename": file_name})
			}
			if got := response.Header().Get("Content-Type"); got != want_mime {
				t.Fatalf("Content-Type = %q, want %q", got, want_mime)
			}
			if got := response.Header().Get("Content-Disposition"); got != want_disposition {
				t.Fatalf("Content-Disposition = %q, want %q", got, want_disposition)
			}
			if got := response.Header().Get("X-Content-Type-Options"); got != "nosniff" {
				t.Fatalf("X-Content-Type-Options = %q, want nosniff", got)
			}
			if response.Body.String() != string(file_body) {
				t.Fatalf("response body = %q, want %q", response.Body.String(), file_body)
			}
		})
	}

	t.Run("extensionless file uses attachment fallback", func(t *testing.T) {
		file_body := []byte("extensionless file")
		if err := os.WriteFile(filepath.Join(image_dir, "extensionless"), file_body, 0600); err != nil {
			t.Fatalf("write extensionless fixture: %v", err)
		}
		response := make_request("extensionless", "192.0.2.88:1234")
		if response.Code != http.StatusOK || response.Body.String() != string(file_body) {
			t.Fatalf("extensionless response = (%d, %q), want file body", response.Code, response.Body.String())
		}
		if got := response.Header().Get("Content-Type"); got != "application/octet-stream" {
			t.Fatalf("extensionless Content-Type = %q", got)
		}
		if got := response.Header().Get("Content-Disposition"); got != mime.FormatMediaType("attachment", map[string]string{"filename": "extensionless"}) {
			t.Fatalf("extensionless Content-Disposition = %q", got)
		}
	})

	t.Run("cache version serves base file with immutable cache header", func(t *testing.T) {
		file_body := []byte("cached image")
		if err := os.WriteFile(filepath.Join(image_dir, "cached.png"), file_body, 0600); err != nil {
			t.Fatalf("write cache fixture: %v", err)
		}
		response := make_request("cached.png.cache_v42", "192.0.2.88:1234")
		if response.Code != http.StatusOK || response.Body.String() != string(file_body) {
			t.Fatalf("cache response = (%d, %q), want file body", response.Code, response.Body.String())
		}
		if got := response.Header().Get("Cache-Control"); got != "private, max-age=31536000, immutable" {
			t.Fatalf("Cache-Control = %q", got)
		}
	})

	t.Run("missing directory and traversal are not served", func(t *testing.T) {
		if err := os.Mkdir(filepath.Join(image_dir, "directory.png"), 0700); err != nil {
			t.Fatalf("create directory fixture: %v", err)
		}
		outside_file := filepath.Join(filepath.Dir(image_dir), "outside.png")
		if err := os.WriteFile(outside_file, []byte("outside"), 0600); err != nil {
			t.Fatalf("write outside fixture: %v", err)
		}
		defer os.Remove(outside_file)

		for _, name := range []string{"missing.png", "directory.png", "../outside.png"} {
			response := make_request(name, "192.0.2.88:1234")
			if response.Code != http.StatusNotFound {
				t.Fatalf("%q status = %d, want 404", name, response.Code)
			}
		}
	})

	t.Run("invalid normalized name and unauthorized peer are rejected", func(t *testing.T) {
		invalid_name := make_request("..cache_v1", "192.0.2.88:1234")
		if invalid_name.Code != http.StatusBadRequest {
			t.Fatalf("invalid normalized name status = %d, want 400", invalid_name.Code)
		}
		unauthorized := make_request("phase9.png", "not-an-ip")
		if unauthorized.Code != http.StatusForbidden {
			t.Fatalf("invalid peer status = %d, want 403", unauthorized.Code)
		}
		if got := unauthorized.Header().Get("X-Content-Type-Options"); got != "nosniff" {
			t.Fatalf("forbidden X-Content-Type-Options = %q, want nosniff", got)
		}
	})

	t.Run("empty path preserves empty response", func(t *testing.T) {
		response := make_request("", "192.0.2.88:1234")
		if response.Code != http.StatusOK || strings.TrimSpace(response.Body.String()) != "" {
			t.Fatalf("empty path response = (%d, %q), want empty 200", response.Code, response.Body.String())
		}
	})
}
