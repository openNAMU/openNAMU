package test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_View_edit_file_upload_unifies_reader_inputs(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		t.Run(backend, func(t *testing.T) {
			db := audit_setup(t, backend)
			audit_seed_auth(db, t)

			file_dir := filepath.Join(t.TempDir(), "images")
			if err := os.MkdirAll(file_dir, 0o755); err != nil {
				t.Fatal(err)
			}
			if _, err := db.Exec("insert into other (name, data, coverage) values ('image_where', ?, '')", file_dir); err != nil {
				t.Fatal(err)
			}
			if _, err := db.Exec("insert into html_filter (html, kind, plus, plus_t) values ('png', 'extension', '', '')"); err != nil {
				t.Fatal(err)
			}

			config := tool.Config{IP: "admin_ip", Session: &emptySession{}}
			for _, test_case := range []struct {
				name           string
				file_name      string
				file_data      []byte
				want_file_data []byte
				file_reader    io.ReadCloser
			}{
				{
					name:           "byte input",
					file_name:      "bytes-upload",
					file_data:      []byte("byte-backed content"),
					want_file_data: []byte("byte-backed content"),
				},
				{
					name:           "reader input",
					file_name:      "reader-upload",
					want_file_data: []byte("stream-backed content"),
					file_reader:    io.NopCloser(bytes.NewReader([]byte("stream-backed content"))),
				},
			} {
				t.Run(backend+" "+test_case.name, func(t *testing.T) {
					output := route.View_edit_file_upload_post(config, []route.Upload_file_data{{
						File_name:   test_case.file_name,
						File_data:   test_case.file_data,
						File_reader: test_case.file_reader,
						File_ext:    "png",
						License:     "direct_input",
					}})
					redirect_target, redirected := tool.Get_redirect_target(output)
					want_doc := "file:" + test_case.file_name + ".png"
					if !redirected || redirect_target != "/w/"+tool.Url_parser(want_doc) {
						t.Fatalf("upload output=%q, want redirect to %q", output, "/w/"+tool.Url_parser(want_doc))
					}

					var stored_data string
					if err := db.QueryRow("select data from data where title = ?", want_doc).Scan(&stored_data); err != nil {
						t.Fatalf("read uploaded file document: %v", err)
					}
					file_path := filepath.Join(file_dir, tool.File_name_to_dir(test_case.file_name, "png"))
					file_data, err := os.ReadFile(file_path)
					if err != nil {
						t.Fatalf("read stored upload: %v", err)
					}
					if !bytes.Equal(file_data, test_case.want_file_data) || stored_data == "" {
						t.Fatalf("upload did not preserve bytes and document (file=%q, want=%q, document=%q)", file_data, test_case.want_file_data, stored_data)
					}
				})
			}
		})
	}
}
