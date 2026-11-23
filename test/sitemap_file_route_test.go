package test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func Test_sitemap_file_routes(t *testing.T) {
	prev_mode := gin.Mode()
	gin.SetMode(gin.TestMode)
	t.Cleanup(func() { gin.SetMode(prev_mode) })

	// Read_sitemap_file reads sitemap files from the working directory.
	t.Chdir(t.TempDir())

	index_data := `<?xml version="1.0" encoding="UTF-8"?><sitemapindex></sitemapindex>`
	sitemap_zero_data := `<?xml version="1.0" encoding="UTF-8"?><urlset>zero</urlset>`
	sitemap_one_data := `<?xml version="1.0" encoding="UTF-8"?><urlset>one</urlset>`

	for name, data := range map[string]string{
		"sitemap.xml":   index_data,
		"sitemap_0.xml": sitemap_zero_data,
		"sitemap_1.xml": sitemap_one_data,
	} {
		if err := os.WriteFile(name, []byte(data), 0644); err != nil {
			t.Fatal(err)
		}
	}

	r := gin.New()
	route.Register_sitemap_file_routes(r)

	do := func(path string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		r.ServeHTTP(w, req)
		return w
	}

	// the sitemap index must serve the file content
	w := do("/sitemap.xml")
	if w.Code != http.StatusOK || w.Body.String() != index_data {
		t.Errorf("GET /sitemap.xml = %d %q, want 200 %q", w.Code, w.Body.String(), index_data)
	}
	if got := w.Header().Get("Content-Type"); got != "application/xml; charset=utf-8" {
		t.Errorf("GET /sitemap.xml content-type = %q", got)
	}

	// regression: sub sitemap files must be reachable; the old pattern
	// "/sitemap_:num.xml" left c.Param("num") empty in gin v1.12 and every
	// sub sitemap request returned 404
	for path, want := range map[string]string{
		"/sitemap_0.xml": sitemap_zero_data,
		"/sitemap_1.xml": sitemap_one_data,
	} {
		w := do(path)
		if w.Code != http.StatusOK || w.Body.String() != want {
			t.Errorf("GET %s = %d %q, want 200 %q", path, w.Code, w.Body.String(), want)
		}
	}

	// missing files and invalid paths must stay 404
	for _, path := range []string{
		"/sitemap_2.xml",
		"/sitemap_abc.xml",
		"/sitemap_0.json",
		"/sitemap_0",
		"/sitemap_.xml",
		"/sitemap_-1.xml",
	} {
		if w := do(path); w.Code != http.StatusNotFound {
			t.Errorf("GET %s = %d, want 404", path, w.Code)
		}
	}
}
