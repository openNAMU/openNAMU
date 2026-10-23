package test

import (
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

// Test_view_w_watch_list_add_post verifies the post-redirect target after a
// watchlist/star-list add/remove:
//   - name_from requests (watchlist_from / star_doc_from) redirect to the
//     original document, and the document name must be URL-encoded via
//     tool.Url_parser (the pre-fix bug left special characters unescaped).
//   - non-from requests still redirect to /watch_list and /star_doc.
//   - an anonymous user is redirected to /login.
//   - the user_set toggle is still written (API untouched by this fix).
func Test_view_w_watch_list_add_post(t *testing.T) {
	db := Setup_rebuild_db(t)
	defer tool.DB_close(db)

	user := tool.Config{IP: "user1", Session: &emptySession{}}
	anon := tool.Config{IP: "10.0.0.1", Session: &emptySession{}}

	// redirectTarget returns the target extracted from a Get_redirect() marker.
	redirectTarget := func(name, doType string, cfg tool.Config) (string, bool) {
		out := route.View_w_watch_list_add_post(cfg, name, doType)
		return tool.Get_redirect_target(out)
	}

	t.Run("star_doc_from encodes ? into %3F", func(t *testing.T) {
		target, ok := redirectTarget("What?why", "star_doc_from", user)
		if !ok {
			t.Fatalf("What?why star_doc_from: not a redirect")
		}
		if target != "/w/What%3Fwhy" {
			t.Fatalf("What?why star_doc_from: target %q, expected /w/What%%3Fwhy", target)
		}
	})

	t.Run("watchlist_from encodes # into %23", func(t *testing.T) {
		target, ok := redirectTarget("A#B", "watchlist_from", user)
		if !ok {
			t.Fatalf("A#B watchlist_from: not a redirect")
		}
		if target != "/w/A%23B" {
			t.Fatalf("A#B watchlist_from: target %q, expected /w/A%%23B", target)
		}
	})

	t.Run("various names match /w/+Url_parser(name)", func(t *testing.T) {
		cases := []string{"50%", "A B", "C++?", "한글", "NormalDoc", "123"}
		for _, n := range cases {
			target, ok := redirectTarget(n, "watchlist_from", user)
			if !ok {
				t.Fatalf("%s: not a redirect", n)
			}
			if target != "/w/"+tool.Url_parser(n) {
				t.Fatalf("%s: target %q, expected /w/%s", n, target, tool.Url_parser(n))
			}
		}
	})

	t.Run("plain doc name is unchanged", func(t *testing.T) {
		target, ok := redirectTarget("NormalDoc", "star_doc_from", user)
		if !ok {
			t.Fatalf("NormalDoc: not a redirect")
		}
		if target != "/w/NormalDoc" {
			t.Fatalf("NormalDoc: target %q, expected /w/NormalDoc", target)
		}
	})

	t.Run("non-from star_doc stays /star_doc", func(t *testing.T) {
		target, ok := redirectTarget("NormalDoc", "star_doc", user)
		if !ok {
			t.Fatalf("NormalDoc star_doc: not a redirect")
		}
		if target != "/star_doc" {
			t.Fatalf("NormalDoc star_doc: target %q, expected /star_doc", target)
		}
	})

	t.Run("non-from watchlist stays /watch_list", func(t *testing.T) {
		target, ok := redirectTarget("NormalDoc", "watchlist", user)
		if !ok {
			t.Fatalf("NormalDoc watchlist: not a redirect")
		}
		if target != "/watch_list" {
			t.Fatalf("NormalDoc watchlist: target %q, expected /watch_list", target)
		}
	})

	t.Run("anonymous user is redirected to /login", func(t *testing.T) {
		target, ok := redirectTarget("NormalDoc", "star_doc_from", anon)
		if !ok {
			t.Fatalf("anon star_doc_from: not a redirect")
		}
		if target != "/login" {
			t.Fatalf("anon: target %q, expected /login", target)
		}
	})

	t.Run("user_set toggle is still written by the API", func(t *testing.T) {
		_ = route.View_w_watch_list_add_post(user, "ToggleCheck.pdf", "watchlist_from")
		var cnt int
		if err := db.QueryRow("select count(*) from user_set where id=? and name=? and data=?",
			"user1", "watchlist", "ToggleCheck.pdf").Scan(&cnt); err != nil {
			t.Fatalf("count user_set: %v", err)
		}
		if cnt != 1 {
			t.Fatalf("expected 1 watchlist user_set row for ToggleCheck.pdf, got %d", cnt)
		}
	})
}
