package test

import (
	"bytes"
	"database/sql"
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"opennamu/route"
	"opennamu/route/tool"
)

func Make_security_png(t *testing.T, width int, height int) []byte {
	t.Helper()
	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			canvas.Set(x, y, color.RGBA{R: uint8(x), G: uint8(y), B: 127, A: 255})
		}
	}

	var encoded bytes.Buffer
	if err := png.Encode(&encoded, canvas); err != nil {
		t.Fatalf("encode PNG fixture: %v", err)
	}
	return encoded.Bytes()
}

func Seed_security_file(t *testing.T, db *sql.DB, image_dir string, document_name string, file_data []byte) string {
	t.Helper()
	base_name, extension, valid := route.File_rename_parts(document_name)
	if !valid {
		t.Fatalf("invalid file document name %q", document_name)
	}
	storage_name := tool.File_name_to_dir(base_name, extension)
	if _, err := db.Exec(`insert into data (title, data, type) values (?, '', '')`, document_name); err != nil {
		t.Fatalf("seed file document %q: %v", document_name, err)
	}
	if err := os.WriteFile(filepath.Join(image_dir, storage_name), file_data, 0600); err != nil {
		t.Fatalf("write stored attachment %q: %v", storage_name, err)
	}
	return storage_name
}

func Test_SEC03_media_handlers_enforce_file_document_acl(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		t.Run(backend, func(t *testing.T) {
			db := audit_setup(t, backend)
			audit_seed_auth(db, t)

			image_dir := filepath.Join(t.TempDir(), "images")
			if err := os.MkdirAll(image_dir, 0700); err != nil {
				t.Fatalf("create image fixture directory: %v", err)
			}
			loadSetting(t, db, "image_where", image_dir)
			loadSetting(t, db, "load_ip_select", "default")
			loadSetting(t, db, "trusted_proxies", "127.0.0.1")
			tool.Get_ip_select(db)
			tool.Get_trusted_proxies(db)

			const (
				restricted_document = "file:restricted-audit.png"
				public_document     = "file:public-audit.png"
				missing_document    = "file:missing-audit.png"
				owner_identity      = "203.0.113.42"
			)
			private_bytes := Make_security_png(t, 32, 16)
			restricted_storage := Seed_security_file(t, db, image_dir, restricted_document, private_bytes)
			public_bytes := Make_security_png(t, 24, 12)
			public_storage := Seed_security_file(t, db, image_dir, public_document, public_bytes)
			if _, err := db.Exec(`insert into data (title, data, type) values (?, '', '')`, missing_document); err != nil {
				t.Fatalf("seed missing file document: %v", err)
			}
			if _, err := db.Exec(`insert into acl (title, data, type) values (?, 'owner', 'view')`, restricted_document); err != nil {
				t.Fatalf("restrict file document: %v", err)
			}
			if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', ?, 'owner')`, owner_identity); err != nil {
				t.Fatalf("grant owner identity: %v", err)
			}

			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set(sessions.DefaultKey, &emptySession{})
				c.Next()
			})
			router.GET("/image/*name", route.View_view_image_file)
			router.GET("/thumbnail/:size/*name", route.View_image_thumbnail)
			server := httptest.NewServer(router)
			defer server.Close()

			type http_result struct {
				status  int
				headers http.Header
				body    []byte
			}
			request := func(path string, identity string) http_result {
				t.Helper()
				http_request, err := http.NewRequest(http.MethodGet, server.URL+path, nil)
				if err != nil {
					t.Fatalf("build request %q: %v", path, err)
				}
				if identity != "" {
					http_request.Header.Set("X-Real-IP", identity)
				}
				response, err := http.DefaultClient.Do(http_request)
				if err != nil {
					t.Fatalf("request %q: %v", path, err)
				}
				defer response.Body.Close()
				body, err := io.ReadAll(response.Body)
				if err != nil {
					t.Fatalf("read response %q: %v", path, err)
				}
				return http_result{status: response.StatusCode, headers: response.Header, body: body}
			}

			// Anonymous HTTP requests cannot fetch the restricted original or thumbnail,
			// even when the storage hash is already known.
			for _, path := range []string{
				"/image/" + restricted_storage,
				"/thumbnail/8/" + restricted_storage,
			} {
				response := request(path, "")
				if response.status != http.StatusForbidden || bytes.Contains(response.body, private_bytes) {
					t.Fatalf("anonymous request %q = (%d, %d bytes), want empty 403", path, response.status, len(response.body))
				}
			}

			// Alternate encodings and the immutable cache suffix still resolve to the
			// same file document and cannot bypass its ACL.
			encoded_name := strings.Replace(restricted_storage, ".", "%2e", 1)
			for _, path := range []string{
				"/image/" + encoded_name,
				"/image/" + restricted_storage + ".cache_v9",
				"/thumbnail/8/" + restricted_storage + ".cache_v9",
				"/image/%2e%2e%2f" + restricted_storage,
			} {
				response := request(path, "")
				if response.status == http.StatusOK || bytes.Contains(response.body, private_bytes) {
					t.Fatalf("alternate request %q bypassed file ACL: (%d, %d bytes)", path, response.status, len(response.body))
				}
			}

			// A valid trusted identity still receives both the source bytes and a
			// generated thumbnail through the actual Gin HTTP handlers.
			original := request("/image/"+restricted_storage, owner_identity)
			if original.status != http.StatusOK || !bytes.Equal(original.body, private_bytes) {
				t.Fatalf("owner original response = (%d, %d bytes), want 200 and fixture bytes", original.status, len(original.body))
			}
			thumbnail_url := "/thumbnail/8/" + restricted_storage + ".cache_v9"
			thumbnail := request(thumbnail_url, owner_identity)
			if thumbnail.status != http.StatusOK {
				t.Fatalf("owner thumbnail status = %d, want 200: %q", thumbnail.status, thumbnail.body)
			}
			decoded_thumbnail, err := png.Decode(bytes.NewReader(thumbnail.body))
			if err != nil {
				t.Fatalf("decode authorized thumbnail: %v", err)
			}
			if decoded_thumbnail.Bounds().Dx() > 8 || decoded_thumbnail.Bounds().Dy() > 8 {
				t.Fatalf("authorized thumbnail dimensions = %v, want at most 8x8", decoded_thumbnail.Bounds())
			}
			if response := request(thumbnail_url, ""); response.status != http.StatusForbidden || bytes.Contains(response.body, private_bytes) {
				t.Fatalf("anonymous cached-thumbnail response = (%d, %d bytes), want empty 403", response.status, len(response.body))
			}

			// An ordinary public attachment remains available on both routes.
			public_original := request("/image/"+public_storage, "")
			if public_original.status != http.StatusOK || !bytes.Equal(public_original.body, public_bytes) {
				t.Fatalf("public original response = (%d, %d bytes), want 200 and fixture bytes", public_original.status, len(public_original.body))
			}
			public_thumbnail := request("/thumbnail/8/"+public_storage, "")
			if public_thumbnail.status != http.StatusOK {
				t.Fatalf("public thumbnail status = %d, want 200", public_thumbnail.status)
			}
			if _, err := db.Exec(`insert into alist (name, acl) values ('no_image_view', 'site_view')`); err != nil {
				t.Fatalf("seed site-view-only group: %v", err)
			}
			grant(t, db, "203.0.113.44", "no_image_view")
			if response := request("/image/"+public_storage, "203.0.113.44"); response.status != http.StatusForbidden {
				t.Fatalf("user without global image_view permission status = %d, want 403", response.status)
			}

			// Updating the file-document ACL immediately changes both media routes.
			if _, err := db.Exec(`update acl set data = 'ip' where title = ? and type = 'view'`, restricted_document); err != nil {
				t.Fatalf("change file ACL to IP: %v", err)
			}
			for _, path := range []string{"/image/" + restricted_storage, "/thumbnail/8/" + restricted_storage} {
				if response := request(path, ""); response.status != http.StatusOK {
					t.Fatalf("public-by-ACL request %q status = %d, want 200", path, response.status)
				}
			}
			if _, err := db.Exec(`update acl set data = 'owner' where title = ? and type = 'view'`, restricted_document); err != nil {
				t.Fatalf("restore owner file ACL: %v", err)
			}
			if response := request("/image/"+restricted_storage, ""); response.status != http.StatusForbidden {
				t.Fatalf("restored owner ACL status = %d, want 403", response.status)
			}

			// A denied document view must not emit a usable media URL or attachment
			// contents, even though the file document itself exists.
			page_context, _ := gin.CreateTestContext(httptest.NewRecorder())
			page_context.Request = httptest.NewRequest(http.MethodGet, "/w/"+restricted_document, nil)
			page, page_status := route.View_w(page_context, tool.Config{IP: "127.0.0.1", Session: &emptySession{}}, restricted_document, "")
			if page_status != http.StatusUnauthorized {
				t.Fatalf("restricted file page status = %d, want 401", page_status)
			}
			if strings.Contains(page, "/image/"+restricted_storage) || bytes.Contains([]byte(page), private_bytes) {
				t.Fatalf("denied file page leaked a usable URL or private file bytes")
			}

			// Orphaned storage and missing files fail safely; invalid names are rejected.
			orphan_name := tool.File_name_to_dir("orphan-audit", "png")
			if err := os.WriteFile(filepath.Join(image_dir, orphan_name), private_bytes, 0600); err != nil {
				t.Fatalf("write orphan fixture: %v", err)
			}
			if response := request("/image/"+orphan_name, ""); response.status != http.StatusNotFound {
				t.Fatalf("orphan image status = %d, want 404", response.status)
			}
			missing_storage := tool.File_name_to_dir("missing-audit", "png")
			if response := request("/image/"+missing_storage, ""); response.status != http.StatusNotFound {
				t.Fatalf("missing image status = %d, want 404", response.status)
			}
			if response := request("/image/..cache_v1", ""); response.status != http.StatusBadRequest {
				t.Fatalf("invalid image name status = %d, want 400", response.status)
			}

			// Distinct file documents that normalize to the same storage name are
			// ambiguous and must never select the less restrictive document.
			ambiguous_storage := Seed_security_file(t, db, image_dir, "file:ambiguous-audit.png", private_bytes)
			if other_storage := Seed_security_file(t, db, image_dir, "file:ambiguous-audit.PNG", private_bytes); other_storage != ambiguous_storage {
				t.Fatalf("case-normalized ambiguous storage names differ: %q != %q", ambiguous_storage, other_storage)
			}
			if _, err := db.Exec(`insert into acl (title, data, type) values (?, 'owner', 'view')`, "file:ambiguous-audit.PNG"); err != nil {
				t.Fatalf("restrict ambiguous file document: %v", err)
			}
			if response := request("/image/"+ambiguous_storage, ""); response.status != http.StatusNotFound {
				t.Fatalf("ambiguous file mapping status = %d, want fail-closed 404", response.status)
			}
		})
	}
}
