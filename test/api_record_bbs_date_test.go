package test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func record_test_seed_post_with_date(db *sql.DB, t *testing.T, set_id, set_code, user, date string) {
	record_test_seed_user_id(db, t, set_code, set_id, user)

	if date == "" {
		return
	}

	if _, err := db.Exec(
		`insert into bbs_data (set_name, set_code, set_id, set_data) values ('date', ?, ?, ?)`,
		set_code, set_id, date,
	); err != nil {
		t.Fatalf("seed date: %v", err)
	}
}

func bbs_date_backends() []string {
	backends := []string{"sqlite", "turso"}
	if os.Getenv("MYSQL_AUDIT_DSN") != "" {
		backends = append(backends, "mysql")
	}
	return backends
}

func assert_bbs_data(t *testing.T, res map[string]any, want [][]string) {
	if response := res["response"]; response != "ok" {
		t.Fatalf("response = %v, want ok", res)
	}

	got, ok := res["data"].([][]string)
	if !ok {
		t.Fatalf("data is not [][]string: %v", res)
	}
	if len(got) != len(want) {
		t.Fatalf("data len = %d, want %d; got=%v want=%v", len(got), len(want), got, want)
	}
	for i, row := range want {
		for j, value := range row {
			if i >= len(got) || j >= len(got[i]) || got[i][j] != value {
				t.Fatalf("data[%d][%d] = %q, want %q", i, j, got[i][j], value)
			}
		}
	}
}

func Test_api_record_bbs_latest_post_date(t *testing.T) {
	for _, backend := range bbs_date_backends() {
		t.Run(backend, func(t *testing.T) {
			t.Run("other_user_date_not_mixed", func(t *testing.T) {
				db := new_record_db(t, backend)
				defer tool.DB_close(db)

				seed_record_access_auth(db, t)
				record_test_seed_post_with_date(db, t, "5", "1", "alice", "2020-01-01")
				record_test_seed_post_with_date(db, t, "5", "2", "bob", "2025-01-01")

				res := route.Api_record_bbs(tool.Config{IP: "alice", Session: &emptySession{}}, "alice", "1")
				assert_bbs_data(t, res, [][]string{{"5", "2020-01-01"}})
			})

			t.Run("own_latest_date", func(t *testing.T) {
				db := new_record_db(t, backend)
				defer tool.DB_close(db)

				seed_record_access_auth(db, t)
				record_test_seed_post_with_date(db, t, "5", "1", "alice", "2020-01-01")
				record_test_seed_post_with_date(db, t, "5", "2", "alice", "2024-01-01")
				record_test_seed_post_with_date(db, t, "5", "3", "bob", "2025-01-01")

				res := route.Api_record_bbs(tool.Config{IP: "alice", Session: &emptySession{}}, "alice", "1")
				assert_bbs_data(t, res, [][]string{{"5", "2024-01-01"}})
			})

			t.Run("legacy_missing_date", func(t *testing.T) {
				db := new_record_db(t, backend)
				defer tool.DB_close(db)

				seed_record_access_auth(db, t)
				record_test_seed_post_with_date(db, t, "5", "1", "alice", "")
				record_test_seed_post_with_date(db, t, "5", "2", "bob", "2025-01-01")

				res := route.Api_record_bbs(tool.Config{IP: "alice", Session: &emptySession{}}, "alice", "1")
				assert_bbs_data(t, res, [][]string{{"5", ""}})
			})

			t.Run("special_and_numeric_sort", func(t *testing.T) {
				db := new_record_db(t, backend)
				defer tool.DB_close(db)

				seed_record_access_auth(db, t)
				record_test_seed_post_with_date(db, t, "2", "1", "alice", "2020-01-01")
				record_test_seed_post_with_date(db, t, "9", "1", "alice", "2020-01-09")
				record_test_seed_post_with_date(db, t, "10", "1", "alice", "2020-01-10")
				record_test_seed_post_with_date(db, t, "11", "1", "alice", "2020-01-11")
				record_test_seed_post_with_date(db, t, "0", "1", "alice", "2020-02-01")
				record_test_seed_post_with_date(db, t, "-1", "1", "alice", "2020-03-01")

				res := route.Api_record_bbs(tool.Config{IP: "alice", Session: &emptySession{}}, "alice", "1")
				assert_bbs_data(t, res, [][]string{
					{"11", "2020-01-11"},
					{"10", "2020-01-10"},
					{"9", "2020-01-09"},
					{"2", "2020-01-01"},
					{"0", "2020-02-01"},
					{"-1", "2020-03-01"},
				})
			})

			t.Run("page_boundary_51", func(t *testing.T) {
				db := new_record_db(t, backend)
				defer tool.DB_close(db)

				seed_record_access_auth(db, t)
				for i := 1; i <= 51; i++ {
					record_test_seed_post_with_date(db, t, fmt.Sprintf("%d", i), "1", "alice", "2026-01-01")
				}

				page1 := route.Api_record_bbs(tool.Config{IP: "alice", Session: &emptySession{}}, "alice", "1")
				page2 := route.Api_record_bbs(tool.Config{IP: "alice", Session: &emptySession{}}, "alice", "2")

				data1 := page1["data"].([][]string)
				if len(data1) != 50 {
					t.Fatalf("page1 len = %d, want 50", len(data1))
				}
				if data1[0][0] != "51" || data1[len(data1)-1][0] != "2" {
					t.Fatalf("page1 boundary wrong: first=%q last=%q", data1[0][0], data1[len(data1)-1][0])
				}
				for _, row := range data1 {
					if row[0] == "1" {
						t.Fatalf("page1 incorrectly contains board 1")
					}
				}

				data2 := page2["data"].([][]string)
				if len(data2) != 1 {
					t.Fatalf("page2 len = %d, want 1", len(data2))
				}
				if data2[0][0] != "1" || data2[0][1] != "2026-01-01" {
					t.Fatalf("page2 wrong: %v", data2[0])
				}
			})

			t.Run("exact_50_no_page2", func(t *testing.T) {
				db := new_record_db(t, backend)
				defer tool.DB_close(db)

				seed_record_access_auth(db, t)
				for i := 1; i <= 50; i++ {
					record_test_seed_post_with_date(db, t, fmt.Sprintf("%d", i), "1", "alice", "2026-01-01")
				}

				page1 := route.Api_record_bbs(tool.Config{IP: "alice", Session: &emptySession{}}, "alice", "1")
				page2 := route.Api_record_bbs(tool.Config{IP: "alice", Session: &emptySession{}}, "alice", "2")

				data1 := page1["data"].([][]string)
				if len(data1) != 50 {
					t.Fatalf("page1 len = %d, want 50", len(data1))
				}
				if data1[0][0] != "50" || data1[len(data1)-1][0] != "1" {
					t.Fatalf("page1 boundary wrong: first=%q last=%q", data1[0][0], data1[len(data1)-1][0])
				}

				data2 := page2["data"].([][]string)
				if len(data2) != 0 {
					t.Fatalf("page2 len = %d, want 0; got=%v", len(data2), data2)
				}
			})

			t.Run("permission_regression", func(t *testing.T) {
				db := new_record_db(t, backend)
				defer tool.DB_close(db)

				seed_record_access_auth(db, t)
				record_test_seed_post_with_date(db, t, "5", "1", "alice", "2026-01-01")

				self := route.Api_record_bbs(tool.Config{IP: "alice", Session: &emptySession{}}, "alice", "1")
				hidel := route.Api_record_bbs(tool.Config{IP: "hidel_user", Session: &emptySession{}}, "alice", "1")
				regular := route.Api_record_bbs(tool.Config{IP: "regular", Session: &emptySession{}}, "alice", "1")

				if response := self["response"]; response != "ok" {
					t.Fatalf("self response = %v, want ok", response)
				}
				if response := hidel["response"]; response != "ok" {
					t.Fatalf("hidel response = %v, want ok", response)
				}
				if response := regular["response"]; response != "require auth" {
					t.Fatalf("regular response = %v, want require auth", response)
				}
			})
		})
	}
}
