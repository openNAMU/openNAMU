package test

// Regression for the vote list ordering fix (route/api_vote_list.go).
// The vote.id column is TEXT, so without an explicit numeric ORDER BY the
// list order is engine-dependent and, on real data, wrong (id "10" sorts
// before "2" lexicographically). The fix applies `order by id + 0 desc`
// (the existing numeric-text convention used across the codebase, e.g.
// route/tool/db_read.go) so both open and close lists are sorted by numeric
// ID descending and pagination (limit offset, 51 -> has_next) is stable.

import (
	"fmt"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

// vote_list_api_call invokes the real Api_vote_list through the exact
// production signature.
func vote_list_api_call(t *testing.T, num string, set_type string) map[string]any {
	t.Helper()
	return route.Api_vote_list(tool.Config{IP: "admin_ip", Session: &emptySession{}}, set_type, num)
}

func platform_turso_ok() bool { return tool.Turso_platform_supported() }

// Test_vote_list_ordering_regression: seeds votes with ids in scrambled order
// including the 2/10 lexicographic trap and the 50-row page boundary, then
// asserts the API returns numeric-descending order with no duplicates or gaps
// across pages.
func Test_vote_list_ordering_regression(t *testing.T) {
	db := audit_setup(t, "sqlite")
	audit_seed_auth(db, t)

	// seed ids 1..60 in a scrambled insertion order to make engine ordering
	// irrelevant: insert 2, 10, 1, 11, 3, 60, 4, ... mixing types
	type seed_item struct {
		id   int
		typ  string
	}
	order := []int{2, 10, 1, 11, 3, 60, 50, 51, 49, 4, 5, 6, 7, 8, 9, 12, 13, 14, 15}
	for i := 16; i <= 60; i++ {
		if i != 49 && i != 50 && i != 51 && i != 60 {
			order = append(order, i)
		}
	}
	for _, id := range order {
		vote_type := "open"
		if id%3 == 0 {
			vote_type = "n_open"
		} else if id%7 == 0 {
			vote_type = "close"
		} else if id%11 == 0 {
			vote_type = "n_close"
		}
		if _, err := db.Exec(
			`insert into vote (name, id, subject, data, user, type, acl) values (?, ?, ?, ?, '', ?, '')`,
			fmt.Sprintf("vote_%d", id), fmt.Sprintf("%d", id), "subject", "data", vote_type,
		); err != nil {
			t.Fatalf("seed vote %d: %v", id, err)
		}
	}
	// also add option rows (type = 'option') which must never appear in lists
	if _, err := db.Exec(
		`insert into vote (name, id, subject, data, user, type, acl) values ('', '1', '', 'opt', '', 'option', '')`,
	); err != nil {
		t.Fatalf("seed option: %v", err)
	}

	// expected sorted ids per list type (numeric descending)
	collect := func(set_type string) []string {
		out := []string{}
		for _, id := range order {
			vote_type := "open"
			if id%3 == 0 {
				vote_type = "n_open"
			} else if id%7 == 0 {
				vote_type = "close"
			} else if id%11 == 0 {
				vote_type = "n_close"
			}
			in_list := false
			if set_type == "open" {
				in_list = vote_type == "open" || vote_type == "n_open"
			} else {
				in_list = vote_type == "close" || vote_type == "n_close"
			}
			if in_list {
				out = append(out, fmt.Sprintf("%d", id))
			}
		}
		// sort numeric desc
		for i := 0; i < len(out); i++ {
			for j := i + 1; j < len(out); j++ {
				a, b := 0, 0
				fmt.Sscanf(out[i], "%d", &a)
				fmt.Sscanf(out[j], "%d", &b)
				if b > a {
					out[i], out[j] = out[j], out[i]
				}
			}
		}
		return out
	}

	// walk pages via the real API and verify: ordering, no duplicates,
	// has_next correctness at the 50 boundary
	seen := map[string]bool{}
	page := 1
	total_open := 0
	for {
		result := vote_list_api_call(t, fmt.Sprintf("%d", page), "open")
		list, _ := result["data"].([][]string)
		has_next, _ := result["has_next"].(bool)

		// ordering check within the page
		for i := 1; i < len(list); i++ {
			var a, b int
			fmt.Sscanf(list[i-1][1], "%d", &a)
			fmt.Sscanf(list[i][1], "%d", &b)
			if a <= b {
				t.Fatalf("page %d: ordering broken at index %d: %d then %d (ids 2 vs 10 lexicographic trap?)",
					page, i, a, b)
			}
		}
		for _, row := range list {
			if seen[row[1]] {
				t.Fatalf("page %d: duplicate id %s across pages", page, row[1])
			}
			seen[row[1]] = true
			total_open++
		}

		if has_next {
			if len(list) != 50 {
				t.Fatalf("page %d: has_next with %d rows, want exactly 50", page, len(list))
			}
			page++
			if page > 10 {
				t.Fatalf("pagination did not terminate")
			}
			continue
		}
		// last page: must be <= 50 rows
		if len(list) > 50 {
			t.Fatalf("page %d: last page has %d rows (limit leak)", page, len(list))
		}
		break
	}

	want := collect("open")
	if total_open != len(want) {
		t.Fatalf("open list total=%d, want %d (missing/dropped rows)", total_open, len(want))
	}
	for _, id := range want {
		if !seen[id] {
			t.Fatalf("open list missing id %s", id)
		}
	}

	// explicit 2 vs 10 ordering check via a page-1 query
	first := vote_list_api_call(t, "1", "open")
	first_list, _ := first["data"].([][]string)
	pos2, pos10 := -1, -1
	for i, row := range first_list {
		if row[1] == "2" {
			pos2 = i
		}
		if row[1] == "10" {
			pos10 = i
		}
	}
	if pos2 < 0 || pos10 < 0 || pos10 >= pos2 {
		t.Fatalf("id 2 must come BEFORE id 10 (numeric desc): pos2=%d pos10=%d (lexicographic bug if 10 first)",
			pos2, pos10)
	}
	// close list ordering sanity
	result := vote_list_api_call(t, "1", "close")
	close_list, _ := result["data"].([][]string)
	for i := 1; i < len(close_list); i++ {
		var a, b int
		fmt.Sscanf(close_list[i-1][1], "%d", &a)
		fmt.Sscanf(close_list[i][1], "%d", &b)
		if a <= b {
			t.Fatalf("close list ordering broken at %d: %d then %d", i, a, b)
		}
	}
}

// Test_vote_list_ordering_turso: same ordering assertions on the turso backend.
func Test_vote_list_ordering_turso(t *testing.T) {
	if !platform_turso_ok() {
		t.Skip("turso unsupported platform")
	}
	db := audit_setup(t, "turso")
	audit_seed_auth(db, t)

	for _, id := range []int{2, 10, 1, 3, 9, 60} {
		vote_type := "open"
		if _, err := db.Exec(
			`insert into vote (name, id, subject, data, user, type, acl) values (?, ?, ?, ?, '', ?, '')`,
			fmt.Sprintf("vote_%d", id), fmt.Sprintf("%d", id), "subject", "data", vote_type,
		); err != nil {
			t.Fatalf("seed vote %d: %v", id, err)
		}
	}

	result := vote_list_api_call(t, "1", "open")
	list, _ := result["data"].([][]string)
	if len(list) < 2 {
		t.Fatalf("turso vote list rows=%d", len(list))
	}
	// ids seeded: 60, 9, 3, 2, 1 open (10 is open too) -> numeric desc: 60,10,9,3,2,1
	want := []string{"60", "10", "9", "3", "2", "1"}
	for i, id := range want {
		if i >= len(list) || list[i][1] != id {
			t.Fatalf("turso ordering broken at %d: got %v, want first ids %v", i, list, want)
		}
	}
}
