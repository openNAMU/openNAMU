package test

import (
	"errors"
	"os"
	"strings"
	"testing"

	"opennamu/route/tool"
)

// Make_sitemap used to leave sitemap_N.xml files from earlier runs in place when
// the document count shrank, so the stale parts (and the deleted documents they
// list) kept being served. Regenerating must remove parts the new index does not
// reference.
func Test_make_sitemap_removes_stale_files(t *testing.T) {
	t.Chdir(t.TempDir())
	db := Create_sitemap_db(t)
	set_test_setting(db, t, "domain", "wiki.example.org")

	// simulate a previous run that needed three parts
	for _, name := range []string{"sitemap_1.xml", "sitemap_2.xml"} {
		if err := os.WriteFile(name, []byte("<urlset><url><loc>stale</loc></url></urlset>"), 0644); err != nil {
			t.Fatalf("create stale %s: %v", name, err)
		}
	}

	if _, _, err := tool.Make_sitemap(db); err != nil {
		t.Fatalf("Make_sitemap: %v", err)
	}

	if _, err := os.Stat("sitemap_0.xml"); err != nil {
		t.Fatalf("sitemap_0.xml missing: %v", err)
	}
	for _, name := range []string{"sitemap_1.xml", "sitemap_2.xml"} {
		if _, err := os.Stat(name); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("stale %s was not removed: stat error=%v", name, err)
		}
	}

	index_data, err := os.ReadFile("sitemap.xml")
	if err != nil {
		t.Fatalf("read sitemap.xml: %v", err)
	}
	if strings.Contains(string(index_data), "sitemap_1.xml") {
		t.Errorf("sitemap index still references a removed part: %s", index_data)
	}

	// unrelated files that only look similar must be left alone
	if err := os.WriteFile("sitemap_abc.xml", []byte("keep"), 0644); err != nil {
		t.Fatalf("create sitemap_abc.xml: %v", err)
	}
	if _, _, err := tool.Make_sitemap(db); err != nil {
		t.Fatalf("Make_sitemap with a non-numeric part name: %v", err)
	}
	if _, err := os.Stat("sitemap_abc.xml"); err != nil {
		t.Fatalf("unrelated sitemap_abc.xml must be kept: %v", err)
	}
}

// A failing removal must be reported with the typed error the API understands
// instead of silently keeping stale parts.
func Test_make_sitemap_reports_remove_error(t *testing.T) {
	t.Chdir(t.TempDir())
	db := Create_sitemap_db(t)
	set_test_setting(db, t, "domain", "wiki.example.org")

	// a non-empty directory cannot be removed, which forces the removal failure
	if err := os.Mkdir("sitemap_1.xml", 0755); err != nil {
		t.Fatalf("create sitemap_1.xml directory: %v", err)
	}
	if err := os.WriteFile("sitemap_1.xml/keep", []byte("x"), 0644); err != nil {
		t.Fatalf("create blocking file: %v", err)
	}

	_, _, err := tool.Make_sitemap(db)
	if err == nil {
		t.Fatal("Make_sitemap must report a removal failure")
	}
	if !errors.Is(err, tool.Err_sitemap_file_remove) {
		t.Fatalf("removal failure must wrap Err_sitemap_file_remove, got %v", err)
	}
	if !strings.Contains(err.Error(), "sitemap_1.xml") {
		t.Errorf("removal failure should name the file, got %v", err)
	}
}
