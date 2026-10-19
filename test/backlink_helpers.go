package test

import (
	"database/sql"
	"testing"
)

// hasBacklink reports whether the back table holds a row for the given link
// whose title matches the given title. On SQLite the test DB's BINARY collation
// makes the comparison byte-exact.
func hasBacklink(db *sql.DB, link, title string, t *testing.T) bool {
	var c int
	if err := db.QueryRow(
		"select count(*) from back where link = ? and title = ?",
		link, title,
	).Scan(&c); err != nil {
		t.Fatalf("hasBacklink: %v", err)
	}
	return c > 0
}

// backlink_titles returns the titles of backlink rows for the given link.
func backlink_titles(db *sql.DB, link string, t *testing.T) []string {
	rows, err := db.Query("select title from back where link = ?", link)
	if err != nil {
		t.Fatalf("backlink_titles: %v", err)
	}
	defer rows.Close()

	var titles []string
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			t.Fatalf("backlink_titles scan: %v", err)
		}
		titles = append(titles, title)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("backlink_titles: %v", err)
	}
	return titles
}
