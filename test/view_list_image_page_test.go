package test

import (
	"fmt"
	"strings"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

// runImageList seeds N image file docs and renders the first image-list page.
func runImageList(t *testing.T, n int, page string) string {
	t.Helper()
	db := Setup_rebuild_db(t)
	defer tool.DB_close(db)
	for i := 1; i <= n; i++ {
		if _, err := db.Exec(`insert into data (title, data) values (?, '')`,
			"file:image"+fmt.Sprint(i)+".png"); err != nil {
			t.Fatalf("seed image %d: %v", i, err)
		}
	}
	config := tool.Config{IP: "admin_ip", Session: &emptySession{}}
	return route.View_list_image_page(config, page)
}

// listImageCount counts list image cells (only rendered images have alt="file:image...").
func listImageCount(s string) int {
	return strings.Count(s, `alt="file:image`)
}

// fileLinkCount counts file-document link cells (href="/w/file%3A...").
func fileLinkCount(s string) int {
	return strings.Count(s, `href="/w/file%3A`)
}

// Test_view_list_image_page_last_row_file_links verifies that the file-document
// links for the last row are rendered (previously the trailing sub_data was not
// flushed into render_data), while preserving the existing 4-per-row layout,
// the 50-per-page pagination, and the "no data" notice.
func Test_view_list_image_page_last_row_file_links(t *testing.T) {
	cases := []struct {
		name       string
		n          int
		wantImages int
		wantFiles  int
		wantNext   bool
		wantEmpty  bool
	}{
		{"no images: data-missing notice, no next page", 0, 0, 0, false, true},
		{"single image: image and file link shown", 1, 1, 1, false, false},
		{"four images: one row, all file links", 4, 4, 4, false, false},
		{"five images: last-row file link rendered", 5, 5, 5, false, false},
		{"fifty images: last row complete, no next page", 50, 50, 50, false, false},
		{"fifty-one: first page shows 50, next page kept", 51, 50, 50, true, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out := runImageList(t, c.n, "1")
			if got := listImageCount(out); got != c.wantImages {
				t.Fatalf("want %d list images, got %d", c.wantImages, got)
			}
			if got := fileLinkCount(out); got != c.wantFiles {
				t.Fatalf("want %d file links, got %d", c.wantFiles, got)
			}
			if want := strings.Contains(out, `href="/list/image/2"`); want != c.wantNext {
				t.Fatalf("want next page = %v, got %v", c.wantNext, want)
			}
			if c.wantEmpty && !strings.Contains(out, "정보 없음") {
				t.Fatalf("want data-missing notice, got none")
			}
		})
	}
}
