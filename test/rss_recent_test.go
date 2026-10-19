package test

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

// parse_rss_feed is a struct used to unmarshal and assert on the generated RSS.
type parse_rss struct {
	Channel parse_channel `xml:"channel"`
}

type parse_channel struct {
	Title         string   `xml:"title"`
	Link          string   `xml:"link"`
	Description   string   `xml:"description"`
	Language      string   `xml:"language"`
	LastBuildDate string   `xml:"lastBuildDate"`
	Item          []parse_item `xml:"item"`
}

type parse_item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Guid        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

func parse_rss_xml(t *testing.T, body string) parse_rss {
	var feed parse_rss
	if err := xml.Unmarshal([]byte(body), &feed); err != nil {
		t.Fatalf("parse rss: %v (xml=%s)", err, body)
	}
	return feed
}

func rss_item_titles(feed parse_rss) []string {
	titles := []string{}
	for _, item := range feed.Channel.Item {
		titles = append(titles, item.Title)
	}
	return titles
}

func contains_str(list []string, target string) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}

// Setup_rss_db builds an in-memory sqlite wiki with the base schema and the
// permission groups + wiki settings this task relies on. withDomain seeds the
// wiki "domain" setting (wiki.example.com) when true; when false the feed must
// fall back to the request host.
func Setup_rss_db(t *testing.T, withDomain bool) *sql.DB {
	t.Helper()

	db_dir := filepath.Join(t.TempDir(), "db")
	if err := os.MkdirAll(db_dir, 0755); err != nil {
		t.Fatalf("mkdir db dir: %v", err)
	}

	os.Setenv("NAMU_DB", db_dir)
	os.Setenv("NAMU_DB_TYPE", "sqlite")
	defer os.Unsetenv("NAMU_DB")
	defer os.Unsetenv("NAMU_DB_TYPE")

	tool.DB_init()
	db := tool.DB_connect()

	exec := func(query string, args ...any) {
		if _, err := db.Exec(query, args...); err != nil {
			t.Fatalf("seed exec %v: %v", query, err)
		}
	}

	// permission groups (checked by acl resolution)
	exec("insert into alist (name, acl) values ('owner', 'owner')")
	exec("insert into alist (name, acl) values ('admin', 'admin')")
	exec("insert into alist (name, acl) values ('user', 'user')")
	exec("insert into alist (name, acl) values ('ip', 'ip')")

	// wiki settings
	exec(`insert into other (name, data, coverage) values ('name', 'TestWiki', '')`)
	exec(`insert into other (name, data, coverage) values ('language', 'ko-KR', '')`)
	exec(`insert into other (name, data, coverage) values ('http_select', 'http', '')`)
	if withDomain {
		exec(`insert into other (name, data, coverage) values ('domain', 'wiki.example.com', '')`)
	}

	return db
}

// seed_change inserts one recent change: a history row (type '' = normal edit)
// plus the matching rc row (type '' so it is picked up by the default
// recent-change view, which reads set_type ''). hide is the hide flag ("1" =
// hidden, "" = public).
func seed_change(db *sql.DB, t *testing.T, id, title, date, hide string) {
	_, err := db.Exec(
		`insert into history (id, title, data, date, ip, send, leng, hide, type) values (?, ?, ?, ?, ?, ?, ?, ?, '')`,
		id, title, "content "+id, date, "192.168.0.1", "edit", "+5", hide,
	)
	if err != nil {
		t.Fatalf("seed history %s: %v", id, err)
	}
	_, err = db.Exec(`insert into rc (id, title, date, type) values (?, ?, ?, '')`, id, title, date)
	if err != nil {
		t.Fatalf("seed rc %s: %v", id, err)
	}
}

// seed_admin creates the "rss_admin" user that the permission resolver maps to
// the "admin" acl group (via user_set acl=admin + acl_group admin->admin).
func seed_admin(db *sql.DB, t *testing.T) {
	// The "admin" group is already seeded in Setup_rss_db; this adds the
	// "rss_admin" user, mapped to the "admin" group via user_set (id, name='acl', data='admin').
	_, err := db.Exec(`insert into user_set (id, name, data) values ('rss_admin', 'acl', 'admin')`)
	if err != nil {
		t.Fatalf("seed admin: %v", err)
	}
}

// Test_rss_recent_parsing verifies the feed is valid XML with the expected
// channel fields and that item titles, links, guids, and pubDate are correct,
// including Korean and special characters that must survive XML escaping.
func Test_rss_recent_parsing(t *testing.T) {
	db := Setup_rss_db(t, true)
	defer tool.DB_close(db)

	seed_change(db, t, "1", "한국 문서", "2026-10-19 12:34:56", "")
	seed_change(db, t, "2", "문서 & 예외 <name>", "2026-10-19 12:33:56", "")
	seed_change(db, t, "3", `Quote "Special" Doc`, "2026-10-19 12:32:56", "")

	config := tool.Config{IP: "192.168.0.1"}
	xml_out := route.View_rss_recent(config, "ignored.local")

	feed := parse_rss_xml(t, xml_out)

	if feed.Channel.Title != "TestWiki" {
		t.Errorf("channel title = %q, want TestWiki", feed.Channel.Title)
	}
	if feed.Channel.Language != "ko-KR" {
		t.Errorf("channel language = %q, want ko-KR", feed.Channel.Language)
	}
	if len(feed.Channel.Item) != 3 {
		t.Fatalf("items = %d, want 3", len(feed.Channel.Item))
	}

	titles := rss_item_titles(feed)
	for _, want := range []string{"한국 문서", "문서 & 예외 <name>", `Quote "Special" Doc`} {
		if !contains_str(titles, want) {
			t.Errorf("missing title %q (got %v)", want, titles)
		}
	}

	for i, item := range feed.Channel.Item {
		if item.Link == "" {
			t.Fatalf("item %d link empty", i)
		}
		if !strings.HasPrefix(item.Link, "http://wiki.example.com/w/") {
			t.Errorf("item %d link not absolute wiki.example.com URL: %q", i, item.Link)
		}
		if item.Guid == "" {
			t.Fatalf("item %d guid empty", i)
		}
		if !strings.Contains(item.Guid, "|") {
			t.Errorf("item %d guid should be doc|revision: %q", i, item.Guid)
		}
		if item.PubDate == "" {
			t.Fatalf("item %d pubDate empty", i)
		}
		// pubDate must be a valid RFC 822 date.
		if !strings.HasPrefix(item.PubDate, "Mon, ") && !strings.HasPrefix(item.PubDate, "Tue, ") &&
			!strings.HasPrefix(item.PubDate, "Wed, ") && !strings.HasPrefix(item.PubDate, "Thu, ") &&
			!strings.HasPrefix(item.PubDate, "Fri, ") && !strings.HasPrefix(item.PubDate, "Sat, ") &&
			!strings.HasPrefix(item.PubDate, "Sun, ") {
			t.Errorf("item %d pubDate not RFC822 weekday-prefixed: %q", i, item.PubDate)
		}
	}

	// channel link is the recent change page.
	if feed.Channel.Link != "http://wiki.example.com/recent_change" {
		t.Errorf("channel link = %q, want http://wiki.example.com/recent_change", feed.Channel.Link)
	}
}

// Test_rss_recent_max_20 verifies the feed never exceeds 20 items.
func Test_rss_recent_max_20(t *testing.T) {
	db := Setup_rss_db(t, true)
	defer tool.DB_close(db)

	for i := 1; i <= 22; i++ {
		id := fmt.Sprintf("%d", i)
		title := "doc " + id
		date := fmt.Sprintf("2026-10-19 12:%02d:00", i)
		seed_change(db, t, id, title, date, "")
	}

	config := tool.Config{IP: "192.168.0.1"}
	xml_out := route.View_rss_recent(config, "ignored.local")
	feed := parse_rss_xml(t, xml_out)

	if len(feed.Channel.Item) != 20 {
		t.Fatalf("items = %d, want 20 (max)", len(feed.Channel.Item))
	}
}

// Test_rss_recent_hidden_excluded_normal verifies hidden changes are not shown
// to a normal (non-admin) viewer. For a non-admin, the API returns a hidden
// row as a placeholder (empty title); the view must still skip it.
func Test_rss_recent_hidden_excluded_normal(t *testing.T) {
	db := Setup_rss_db(t, true)
	defer tool.DB_close(db)

	seed_change(db, t, "1", "Visible Doc", "2026-10-19 12:00:00", "")
	seed_change(db, t, "2", "Hidden Doc", "2026-10-19 11:00:00", "1")

	config := tool.Config{IP: "192.168.0.1"}
	xml_out := route.View_rss_recent(config, "ignored.local")
	feed := parse_rss_xml(t, xml_out)

	titles := rss_item_titles(feed)
	if !contains_str(titles, "Visible Doc") {
		t.Errorf("visible doc missing, got %v", titles)
	}
	if contains_str(titles, "Hidden Doc") {
		t.Fatalf("hidden doc leaked to non-admin: %v", titles)
	}
	// raw-xml check: the hidden title must not appear anywhere.
	if strings.Contains(xml_out, "Hidden Doc") {
		t.Fatalf("raw xml contains hidden doc title")
	}
}

// Test_rss_recent_hidden_excluded_admin is the stronger case: an admin gets
// the full hidden row (non-empty title, hide='1') from the API. The view must
// filter it out using the hide flag, not an empty-title heuristic.
func Test_rss_recent_hidden_excluded_admin(t *testing.T) {
	db := Setup_rss_db(t, true)
	defer tool.DB_close(db)

	seed_admin(db, t)

	seed_change(db, t, "1", "Visible Doc", "2026-10-19 12:00:00", "")
	seed_change(db, t, "2", "Hidden Doc", "2026-10-19 11:00:00", "1")

	// "rss_admin" resolves to the "admin" acl group, which grants the
	// hidden-records permission, so the API returns the hidden row with full data.
	config := tool.Config{IP: "rss_admin"}
	xml_out := route.View_rss_recent(config, "ignored.local")
	feed := parse_rss_xml(t, xml_out)

	titles := rss_item_titles(feed)
	if !contains_str(titles, "Visible Doc") {
		t.Errorf("visible doc missing for admin, got %v", titles)
	}
	if contains_str(titles, "Hidden Doc") {
		t.Fatalf("hidden doc leaked to admin despite hide filter: %v", titles)
	}
	if strings.Contains(xml_out, "Hidden Doc") {
		t.Fatalf("raw xml contains hidden doc title for admin")
	}
}

// Test_rss_recent_empty verifies a feed with zero changes is still valid XML
// and has no items.
func Test_rss_recent_empty(t *testing.T) {
	db := Setup_rss_db(t, true)
	defer tool.DB_close(db)

	config := tool.Config{IP: "192.168.0.1"}
	xml_out := route.View_rss_recent(config, "ignored.local")

	feed := parse_rss_xml(t, xml_out)

	if feed.Channel.Title != "TestWiki" {
		t.Errorf("channel title = %q, want TestWiki", feed.Channel.Title)
	}
	if len(feed.Channel.Item) != 0 {
		t.Errorf("items = %d, want 0 for empty feed", len(feed.Channel.Item))
	}
}

// Test_rss_recent_domain_fallback verifies that when the wiki domain setting is
// unset, the feed uses the request Host (as an absolute URL) instead of a
// proxy header.
func Test_rss_recent_domain_fallback(t *testing.T) {
	db := Setup_rss_db(t, false)
	defer tool.DB_close(db)

	seed_change(db, t, "1", "Fallback Doc", "2026-10-19 12:00:00", "")

	config := tool.Config{IP: "192.168.0.1"}
	xml_out := route.View_rss_recent(config, "fallback.example.com")
	feed := parse_rss_xml(t, xml_out)

	if feed.Channel.Link != "http://fallback.example.com/recent_change" {
		t.Errorf("channel link = %q, want http://fallback.example.com/recent_change", feed.Channel.Link)
	}
	if len(feed.Channel.Item) != 1 {
		t.Fatalf("items = %d, want 1", len(feed.Channel.Item))
	}
	want := "http://fallback.example.com/w/" + tool.Url_parser("Fallback Doc")
	if feed.Channel.Item[0].Link != want {
		t.Errorf("item link = %q, want %q", feed.Channel.Item[0].Link, want)
	}
}