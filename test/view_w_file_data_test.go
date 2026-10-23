package test

import (
	"bytes"
	"database/sql"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

// Setup_file_data_test_db builds a sqlite DB whose file storage directory is an
// isolated temp directory, so tests can seed on-disk files directly.
func Setup_file_data_test_db(t *testing.T) (*sql.DB, string) {
	t.Helper()

	db_dir := filepath.Join(t.TempDir(), "db")
	img_dir := t.TempDir()
	os.MkdirAll(db_dir, 0755)
	os.Setenv("NAMU_DB", db_dir)
	os.Setenv("NAMU_DB_TYPE", "sqlite")
	defer os.Unsetenv("NAMU_DB")
	defer os.Unsetenv("NAMU_DB_TYPE")

	tool.DB_init()
	db := tool.DB_connect()

	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed %q: %v", q, err)
		}
	}
	exec("insert or replace into other (name, data, coverage) values ('image_where', ?, '')", img_dir)

	return db, img_dir
}

// fileStoragePath mirrors the path computation in View_w_file_data so the test
// can place the on-disk file where the production code expects it.
func fileStoragePath(file_name string) string {
	extension := strings.TrimPrefix(strings.ToLower(filepath.Ext(file_name)), ".")
	base_name := strings.TrimSuffix(file_name, filepath.Ext(file_name))
	return tool.File_name_to_dir(base_name, extension)
}

func pngBytes(t *testing.T, w, h int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}
	return buf.Bytes()
}

func jpegBytes(t *testing.T, w, h int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90}); err != nil {
		t.Fatalf("encode jpeg: %v", err)
	}
	return buf.Bytes()
}

// Test_view_w_file_data verifies that file documents render media tags only for
// image/audio/video files, keep the /image/ download link and size, and use "-"
// (not Vector) as the resolution for non-media files.
func Test_view_w_file_data(t *testing.T) {
	db, imgDir := Setup_file_data_test_db(t)
	defer tool.DB_close(db)

	// createFileDoc seeds a "file:..." data row and writes the on-disk file
	// (when content is non-nil) under the hash-based storage name.
	createFileDoc := func(docName, body string, content []byte) {
		if _, err := db.Exec("insert into data (title, data) values (?, ?)", docName, body); err != nil {
			t.Fatalf("insert doc %q: %v", docName, err)
		}
		if content != nil {
			path := filepath.Join(imgDir, fileStoragePath(strings.TrimPrefix(docName, "file:")))
			if err := os.WriteFile(path, content, 0644); err != nil {
				t.Fatalf("write file for %q: %v", docName, err)
			}
		}
	}

	checkImage := func(docName string) {
		out := route.View_w_file_data(db, docName)
		if !strings.Contains(out, `<img src="/image/`) {
			t.Fatalf("%s: missing <img>, out: %q", docName, out)
		}
		if strings.Contains(out, "<audio") {
			t.Fatalf("%s: unexpected <audio>, out: %q", docName, out)
		}
		if strings.Contains(out, "<video") {
			t.Fatalf("%s: unexpected <video>, out: %q", docName, out)
		}
		if !strings.Contains(out, `href="/image/`) {
			t.Fatalf("%s: missing /image/ download link, out: %q", docName, out)
		}
		if !strings.Contains(out, "KB") {
			t.Fatalf("%s: missing file size (KB), out: %q", docName, out)
		}
	}

	checkAudio := func(docName string) {
		out := route.View_w_file_data(db, docName)
		if !strings.Contains(out, `<audio controls src="`) {
			t.Fatalf("%s: missing <audio>, out: %q", docName, out)
		}
		if strings.Contains(out, "<img") {
			t.Fatalf("%s: unexpected <img>, out: %q", docName, out)
		}
		if strings.Contains(out, "<video") {
			t.Fatalf("%s: unexpected <video>, out: %q", docName, out)
		}
		if !strings.Contains(out, `href="/image/`) {
			t.Fatalf("%s: missing /image/ download link, out: %q", docName, out)
		}
		if strings.Contains(out, "Vector") {
			t.Fatalf("%s: unexpected Vector for audio, out: %q", docName, out)
		}
	}

	checkVideo := func(docName string) {
		out := route.View_w_file_data(db, docName)
		if !strings.Contains(out, `<video controls preload="metadata" src="`) {
			t.Fatalf("%s: missing <video>, out: %q", docName, out)
		}
		if strings.Contains(out, "<img") {
			t.Fatalf("%s: unexpected <img>, out: %q", docName, out)
		}
		if strings.Contains(out, "<audio") {
			t.Fatalf("%s: unexpected <audio>, out: %q", docName, out)
		}
		if !strings.Contains(out, `href="/image/`) {
			t.Fatalf("%s: missing /image/ download link, out: %q", docName, out)
		}
		if strings.Contains(out, "Vector") {
			t.Fatalf("%s: unexpected Vector for video, out: %q", docName, out)
		}
	}

	// Non-media files must show no media tags, keep the download link/size,
	// and use "-" (not Vector) as the resolution.
	checkNonMedia := func(docName string) {
		out := route.View_w_file_data(db, docName)
		if strings.Contains(out, "<img") {
			t.Fatalf("%s: unexpected <img>, out: %q", docName, out)
		}
		if strings.Contains(out, "<audio") {
			t.Fatalf("%s: unexpected <audio>, out: %q", docName, out)
		}
		if strings.Contains(out, "<video") {
			t.Fatalf("%s: unexpected <video>, out: %q", docName, out)
		}
		if !strings.Contains(out, `href="/image/`) {
			t.Fatalf("%s: missing /image/ download link, out: %q", docName, out)
		}
		if strings.Contains(out, "Vector") {
			t.Fatalf("%s: unexpected Vector for a non-media file, out: %q", docName, out)
		}
		if !strings.Contains(out, "KB") {
			t.Fatalf("%s: missing file size (KB), out: %q", docName, out)
		}
	}

	// Image files: <img> + real (decodable) resolution.
	// 1. PNG
	createFileDoc("file:Photo.png", "png body", pngBytes(t, 100, 50))
	checkImage("file:Photo.png")
	if !strings.Contains(route.View_w_file_data(db, "file:Photo.png"), "100x50") {
		t.Fatalf("Photo.png: expected resolution 100x50")
	}
	if strings.Contains(route.View_w_file_data(db, "file:Photo.png"), "Vector") {
		t.Fatalf("Photo.png: unexpected Vector for a decodable image")
	}

	// 2. JPG
	createFileDoc("file:Photo.jpg", "jpg body", jpegBytes(t, 100, 50))
	checkImage("file:Photo.jpg")
	if !strings.Contains(route.View_w_file_data(db, "file:Photo.jpg"), "100x50") {
		t.Fatalf("Photo.jpg: expected resolution 100x50")
	}
	if strings.Contains(route.View_w_file_data(db, "file:Photo.jpg"), "Vector") {
		t.Fatalf("Photo.jpg: unexpected Vector for a decodable image")
	}

	// 3. WEBP: <img> is kept; Go stdlib cannot decode webp so resolution
	//    remains the vector fallback (existing behavior for undecodable images).
	createFileDoc("file:Photo.webp", "webp body", []byte("fake webp bytes"))
	checkImage("file:Photo.webp")
	if strings.Contains(route.View_w_file_data(db, "file:Photo.webp"), "100x50") {
		t.Fatalf("Photo.webp: unexpected pixel resolution for webp")
	}
	if !strings.Contains(route.View_w_file_data(db, "file:Photo.webp"), "Vector") {
		t.Fatalf("Photo.webp: expected Vector for an undecodable image")
	}

	// 4. Uppercase extension is normalized and still detected as an image.
	//    Doc uses .JPG; storage is hash.jpg; content is a real 100x50 JPG.
	createFileDoc("file:CapPhoto.JPG", "jpg body", jpegBytes(t, 100, 50))
	checkImage("file:CapPhoto.JPG")
	if !strings.Contains(route.View_w_file_data(db, "file:CapPhoto.JPG"), "100x50") {
		t.Fatalf("CapPhoto.JPG: expected resolution 100x50 (case normalization)")
	}

	// Audio files: <audio> only.
	// 5. MP3
	createFileDoc("file:Song.mp3", "audio body", []byte("fake mp3"))
	checkAudio("file:Song.mp3")

	// Video files: <video> only.
	// 6. MP4
	createFileDoc("file:Movie.mp4", "video body", []byte("fake mp4"))
	checkVideo("file:Movie.mp4")

	// Non-media files: no media tags, resolution "-", link/size kept.
	// 7-10. PDF, DOCX, ZIP, TXT
	createFileDoc("file:Report.pdf", "pdf body", []byte("%PDF-1.4 fake"))
	checkNonMedia("file:Report.pdf")
	createFileDoc("file:Notes.docx", "docx body", []byte("fake docx"))
	checkNonMedia("file:Notes.docx")
	createFileDoc("file:Archive.zip", "zip body", []byte("fake zip"))
	checkNonMedia("file:Archive.zip")
	createFileDoc("file:Readme.txt", "txt body", []byte("plain text body"))
	checkNonMedia("file:Readme.txt")

	// 11. Document whose file is missing from disk -> empty return (existing behavior).
	createFileDoc("file:Missing.pdf", "missing body", nil)
	if route.View_w_file_data(db, "file:Missing.pdf") != "" {
		t.Fatalf("Missing.pdf: expected empty (file not on disk)")
	}
}
