package test

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_View_edit_file_delete_post_flow(t *testing.T) {
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

			admin_config := tool.Config{IP: "admin_ip", Session: &emptySession{}}
			admin_file_name := "file:delete-me.png"
			admin_file_path, valid := route.File_storage_path(file_dir, admin_file_name)
			if !valid {
				t.Fatal("valid file document name rejected")
			}
			if err := os.WriteFile(admin_file_path, []byte("image"), 0o600); err != nil {
				t.Fatal(err)
			}

			redirect_data := route.View_edit_file_delete(admin_config, admin_file_name, url.Values{})
			redirect_target, redirected := tool.Get_redirect_target(redirect_data)
			if !redirected || redirect_target != "/list/file" {
				t.Fatalf("delete result=%q, want redirect to /list/file", redirect_data)
			}
			if _, err := os.Stat(admin_file_path); !os.IsNotExist(err) {
				t.Fatalf("deleted file still exists (stat error=%v)", err)
			}

			with_document_name := "file:delete-with-document.png"
			with_document_path, valid := route.File_storage_path(file_dir, with_document_name)
			if !valid {
				t.Fatal("valid file document name rejected")
			}
			if err := os.WriteFile(with_document_path, []byte("image"), 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := db.Exec("insert into data (title, data) values (?, 'file content')", with_document_name); err != nil {
				t.Fatal(err)
			}
			with_document_output := route.View_edit_file_delete(admin_config, with_document_name, url.Values{"with_doc": {"1"}})
			with_document_target, redirected := tool.Get_redirect_target(with_document_output)
			want_document_target := "/w/" + tool.Url_parser(with_document_name)
			if !redirected || with_document_target != want_document_target {
				t.Fatalf("delete-with-document result=%q, want redirect to %q", with_document_output, want_document_target)
			}
			if _, err := os.Stat(with_document_path); !os.IsNotExist(err) {
				t.Fatalf("delete-with-document left the file in place (stat error=%v)", err)
			}
			var remaining_documents int
			if err := db.QueryRow("select count(*) from data where title = ?", with_document_name).Scan(&remaining_documents); err != nil {
				t.Fatal(err)
			}
			if remaining_documents != 0 {
				t.Fatalf("delete-with-document left %d document rows", remaining_documents)
			}

			unauthorized_file_name := "file:keep-me.png"
			unauthorized_file_path, valid := route.File_storage_path(file_dir, unauthorized_file_name)
			if !valid {
				t.Fatal("valid unauthorized file document name rejected")
			}
			if err := os.WriteFile(unauthorized_file_path, []byte("private"), 0o600); err != nil {
				t.Fatal(err)
			}

			unauthorized_output := route.View_edit_file_delete(tool.Config{IP: "guest_ip", Session: &emptySession{}}, unauthorized_file_name, url.Values{})
			if !strings.Contains(unauthorized_output, tool.Get_language(db, "authority_error", true)) {
				t.Fatalf("unauthorized delete did not render the auth error: %q", unauthorized_output)
			}
			if _, err := os.Stat(unauthorized_file_path); err != nil {
				t.Fatalf("unauthorized delete changed the file: %v", err)
			}
			unauthorized_invalid_output := route.View_edit_file_delete(tool.Config{IP: "guest_ip", Session: &emptySession{}}, "file:../invalid.png", url.Values{})
			if !strings.Contains(unauthorized_invalid_output, tool.Get_language(db, "file_exist_error", true)) {
				t.Fatalf("invalid file no longer takes precedence over auth in the View: %q", unauthorized_invalid_output)
			}

			invalid_response := route.Api_edit_file_delete_post(admin_config, "file:../invalid.png", url.Values{})
			if invalid_response["response"] != "error" || invalid_response["data"] != "invalid file" {
				t.Fatalf("invalid file API response=%v, want invalid file error", invalid_response)
			}
			invalid_output := route.View_edit_file_delete(admin_config, "file:../invalid.png", url.Values{})
			if !strings.Contains(invalid_output, tool.Get_language(db, "file_exist_error", true)) {
				t.Fatalf("invalid file did not retain its user-facing error: %q", invalid_output)
			}
		})
	}
}
