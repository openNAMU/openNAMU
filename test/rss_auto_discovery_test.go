package test

import (
	"strings"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

// Test_view_list_recent_change_rss verifies the recent-change menu adds an RSS
// link (/rss/recent) and that the rendered page head carries the RSS
// auto-discovery <link rel="alternate"> tag.
func Test_view_list_recent_change_rss(t *testing.T) {
	db := Setup_rss_db(t, false)
	defer tool.DB_close(db)

	seed_change(db, t, "1", "RSS Test Doc", "2026-10-19 12:00:00", "")

	config := tool.Config{IP: "192.168.0.1", Session: &emptySession{}}
	out := route.View_list_recent_change(config, "normal", "10", "1")

	// 1) menu link is present.
	menu_link := `<a href="/rss/recent">(RSS)</a>`
	if !strings.Contains(out, menu_link) {
		t.Errorf("recent-change menu missing RSS link %q", menu_link)
	}

	// 2) head auto-discovery <link> is present (added to the common generated
	// head, not in a per-skin template).
	head_link := `<link rel="alternate" type="application/rss+xml" title="Recent changes" href="/rss/recent">`
	if !strings.Contains(out, head_link) {
		t.Errorf("rendered page head missing RSS <link> %q", head_link)
	}
}
