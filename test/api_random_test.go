package test

import (
	"database/sql"
	"strings"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

// seed_general_doc inserts a non-redirect document into data.
func seed_general_doc(db *sql.DB, t *testing.T, title string) {
	if _, err := db.Exec(`insert into data (title, data, type) values (?, '', '')`, title); err != nil {
		t.Fatalf("seed data %q: %v", title, err)
	}
}

// seed_redirect_doc inserts a redirect document: a data row plus a back row
// of type 'redirect' whose link points at the document, matching
// Get_document_is_redirect / Get_back_title.
func seed_redirect_doc(db *sql.DB, t *testing.T, doc, target string) {
	seed_general_doc(db, t, doc)
	if _, err := db.Exec(`insert into back (title, link, type, data) values (?, ?, 'redirect', '')`, target, doc); err != nil {
		t.Fatalf("seed back redirect %q: %v", doc, err)
	}
}

func Test_api_w_random(t *testing.T) {
	t.Run("general_docs_selected", func(t *testing.T) {
		db := Setup_rss_db(t, false)
		defer tool.DB_close(db)

		seed_general_doc(db, t, "Alpha")
		seed_general_doc(db, t, "Beta")

		res := route.Api_w_random(tool.Config{})
		if res["response"] != "ok" {
			t.Errorf("response = %v, want ok", res["response"])
		}
		got, _ := res["data"].(string)
		if got != "Alpha" && got != "Beta" {
			t.Errorf("data = %q, want Alpha or Beta", got)
		}
	})

	// A redirect document must never be selected. With one general doc the
	// result is deterministic: always the general doc, never the redirect.
	t.Run("redirects_excluded", func(t *testing.T) {
		db := Setup_rss_db(t, false)
		defer tool.DB_close(db)

		seed_general_doc(db, t, "Alpha")
		seed_redirect_doc(db, t, "RedirA", "Alpha")
		seed_redirect_doc(db, t, "RedirB", "Alpha")

		for i := 0; i < 20; i++ {
			res := route.Api_w_random(tool.Config{})
			if res["response"] != "ok" {
				t.Fatalf("iter %d response = %v, want ok", i, res["response"])
			}
			got, _ := res["data"].(string)
			if got != "Alpha" {
				t.Fatalf("iter %d data = %q, want Alpha (redirects excluded)", i, got)
			}
		}
	})

	// No selectable documents -> clear failure state, not an arbitrary name.
	t.Run("empty_not_exist", func(t *testing.T) {
		db := Setup_rss_db(t, false)
		defer tool.DB_close(db)

		res := route.Api_w_random(tool.Config{})
		if res["response"] != "not exist" {
			t.Errorf("response = %v, want not exist", res["response"])
		}
		if got := res["data"]; got != "" {
			t.Errorf("data = %v, want empty (no arbitrary name like Test)", got)
		}
	})

	// All docs are redirects -> nothing selectable -> not exist.
	t.Run("only_redirects_not_exist", func(t *testing.T) {
		db := Setup_rss_db(t, false)
		defer tool.DB_close(db)

		seed_redirect_doc(db, t, "Redir1", "Target1")

		res := route.Api_w_random(tool.Config{})
		if res["response"] != "not exist" {
			t.Errorf("response = %v, want not exist (only redirects present)", res["response"])
		}
	})
}

func Test_api_list_random(t *testing.T) {
	t.Run("general_docs_listed", func(t *testing.T) {
		db := Setup_rss_db(t, false)
		defer tool.DB_close(db)

		seed_general_doc(db, t, "Alpha")
		seed_general_doc(db, t, "Beta")

		res := route.Api_list_random(tool.Config{}, 10)
		if res["response"] != "ok" {
			t.Fatalf("response = %v, want ok", res["response"])
		}
		got := res["data"].([]string)
		if len(got) != 2 {
			t.Fatalf("list = %v, want 2 docs", got)
		}
		for _, v := range got {
			if v != "Alpha" && v != "Beta" {
				t.Fatalf("list contains unexpected %q", v)
			}
		}
	})

	// Redirects must be excluded from the list.
	t.Run("redirects_excluded", func(t *testing.T) {
		db := Setup_rss_db(t, false)
		defer tool.DB_close(db)

		seed_general_doc(db, t, "Alpha")
		seed_redirect_doc(db, t, "RedirA", "Alpha")

		res := route.Api_list_random(tool.Config{}, 10)
		got := res["data"].([]string)
		for _, v := range got {
			if v == "RedirA" {
				t.Fatalf("list contains redirect %q", v)
			}
		}
		if len(got) != 1 || got[0] != "Alpha" {
			t.Fatalf("list = %v, want [Alpha]", got)
		}
	})

	// No docs -> empty list.
	t.Run("empty_list", func(t *testing.T) {
		db := Setup_rss_db(t, false)
		defer tool.DB_close(db)

		res := route.Api_list_random(tool.Config{}, 10)
		if res["response"] != "ok" {
			t.Fatalf("response = %v, want ok", res["response"])
		}
		if got := res["data"].([]string); len(got) != 0 {
			t.Fatalf("list = %v, want empty", got)
		}
	})
}

// View_w_random reuses the existing error-page path for the no-document case
// instead of redirecting to an arbitrary name.
func Test_view_w_random(t *testing.T) {
	t.Run("found_redirects_to_doc", func(t *testing.T) {
		db := Setup_rss_db(t, false)
		defer tool.DB_close(db)

		seed_general_doc(db, t, "OnlyDoc")

		out := route.View_w_random(tool.Config{IP: "192.168.0.1", Session: &emptySession{}})
		target, is_redirect := tool.Get_redirect_target(out)
		if !is_redirect {
			t.Fatalf("expected redirect, got: %s", out)
		}
		if !strings.Contains(target, "OnlyDoc") {
			t.Errorf("redirect target = %q, want to contain OnlyDoc", target)
		}
	})

	// No selectable docs -> error page (not a redirect to "Test").
	t.Run("empty_error_page", func(t *testing.T) {
		db := Setup_rss_db(t, false)
		defer tool.DB_close(db)

		out := route.View_w_random(tool.Config{IP: "192.168.0.1", Session: &emptySession{}})
		if _, is_redirect := tool.Get_redirect_target(out); is_redirect {
			t.Fatalf("expected error page, got redirect: %s", out)
		}
		if out == "" {
			t.Fatal("expected non-empty error page")
		}
	})
}
