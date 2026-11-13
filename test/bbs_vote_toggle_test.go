package test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"

	"opennamu/route"
	"opennamu/route/tool"
)

func Bbs_vote_test_seed_post(db *sql.DB, t *testing.T, set_id string, set_code string, post_user string) {
	t.Helper()
	for _, value := range []struct {
		name string
		data string
	}{
		{name: "title", data: "Vote test post"},
		{name: "data", data: "Vote test body"},
		{name: "date", data: "2026-01-01 00:00:00"},
		{name: "user_id", data: post_user},
		{name: "comment_count", data: "0"},
	} {
		if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, value.name, set_code, set_id, value.data); err != nil {
			t.Fatalf("seed vote post %s/%s %s: %v", set_id, set_code, value.name, err)
		}
	}
}

func Bbs_vote_test_seed_owner(db *sql.DB, t *testing.T, user string) {
	t.Helper()
	if user == "admin_ip" {
		return
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', ?, 'owner')`, user); err != nil {
		t.Fatalf("seed vote owner %q: %v", user, err)
	}
}

func Bbs_vote_test_config(user string) tool.Config {
	return tool.Config{IP: user, Session: &emptySession{}}
}

func Bbs_vote_test_call(t *testing.T, target string, config tool.Config, set_id string, set_code string, comment_code string, vote_type string) map[string]any {
	t.Helper()
	if target == "post" {
		return route.Api_bbs_w_tabom_post(config, set_id, set_code, vote_type)
	}
	return route.Api_bbs_w_comment_tabom_post(config, set_id, set_code, comment_code, vote_type)
}

func Bbs_vote_test_comment_data(t *testing.T, config tool.Config, set_id string, set_code string, comment_code string) map[string]string {
	t.Helper()
	result := route.Api_bbs_w_comment_one(config, false, "", set_id+"-"+set_code+"-"+comment_code)
	if result["response"] != "ok" {
		t.Fatalf("read comment vote state: %#v", result)
	}
	comments, ok := result["data"].([]map[string]string)
	if !ok || len(comments) != 1 {
		t.Fatalf("comment vote state data = %#v, want one comment", result["data"])
	}
	return comments[0]
}

func Bbs_vote_test_count(t *testing.T, db *sql.DB, set_name string, set_id string, set_code string) int {
	t.Helper()
	var count int
	if err := db.QueryRow(`select count(*) from bbs_data where set_name = ? and set_id = ? and set_code = ?`, set_name, set_id, set_code).Scan(&count); err != nil {
		t.Fatalf("count %s rows at %s/%s: %v", set_name, set_id, set_code, err)
	}
	return count
}

func Bbs_vote_test_counter_values(t *testing.T, db *sql.DB, set_name string, set_id string, set_code string) []string {
	t.Helper()
	rows, err := db.Query(`select set_data from bbs_data where set_name = ? and set_id = ? and set_code = ?`, set_name, set_id, set_code)
	if err != nil {
		t.Fatalf("read %s counters at %s/%s: %v", set_name, set_id, set_code, err)
	}
	defer rows.Close()

	values := []string{}
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			t.Fatalf("scan %s counter: %v", set_name, err)
		}
		values = append(values, value)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("read %s counters: %v", set_name, err)
	}
	return values
}

func Bbs_vote_test_assert_state(t *testing.T, db *sql.DB, set_id string, set_code string, state string, up_count int, down_count int) {
	t.Helper()
	up_votes := Bbs_vote_test_count(t, db, "tabom_list", set_id, set_code)
	down_votes := Bbs_vote_test_count(t, db, "tabom_down_list", set_id, set_code)
	want_up_votes := 0
	want_down_votes := 0
	if state == "up" {
		want_up_votes = 1
	}
	if state == "down" {
		want_down_votes = 1
	}
	if up_votes != want_up_votes || down_votes != want_down_votes {
		t.Fatalf("vote rows at %s/%s = up:%d down:%d, want state %s", set_id, set_code, up_votes, down_votes, state)
	}

	up_values := Bbs_vote_test_counter_values(t, db, "tabom_count", set_id, set_code)
	down_values := Bbs_vote_test_counter_values(t, db, "tabom_down_count", set_id, set_code)
	if up_count == 0 && len(up_values) == 0 {
		// Keeping a missing zero-valued legacy counter is valid.
	} else if len(up_values) == 0 || up_values[0] != fmt.Sprint(up_count) {
		t.Fatalf("up counter values at %s/%s = %v, want %d", set_id, set_code, up_values, up_count)
	}
	if down_count == 0 && len(down_values) == 0 {
		// Keeping a missing zero-valued legacy counter is valid.
	} else if len(down_values) == 0 || down_values[0] != fmt.Sprint(down_count) {
		t.Fatalf("down counter values at %s/%s = %v, want %d", set_id, set_code, down_values, down_count)
	}
	for _, value := range up_values {
		if value != fmt.Sprint(up_count) {
			t.Fatalf("duplicate up counters at %s/%s = %v, want %d", set_id, set_code, up_values, up_count)
		}
	}
	for _, value := range down_values {
		if value != fmt.Sprint(down_count) {
			t.Fatalf("duplicate down counters at %s/%s = %v, want %d", set_id, set_code, down_values, down_count)
		}
	}
}

func Bbs_vote_test_assert_totals(t *testing.T, db *sql.DB, set_id string, set_code string, up_count int, down_count int) {
	t.Helper()
	if got := Bbs_vote_test_count(t, db, "tabom_list", set_id, set_code); got != up_count {
		t.Fatalf("total upvote rows at %s/%s = %d, want %d", set_id, set_code, got, up_count)
	}
	if got := Bbs_vote_test_count(t, db, "tabom_down_list", set_id, set_code); got != down_count {
		t.Fatalf("total downvote rows at %s/%s = %d, want %d", set_id, set_code, got, down_count)
	}
	up_values := Bbs_vote_test_counter_values(t, db, "tabom_count", set_id, set_code)
	down_values := Bbs_vote_test_counter_values(t, db, "tabom_down_count", set_id, set_code)
	for _, value := range up_values {
		if value != fmt.Sprint(up_count) {
			t.Fatalf("up counters at %s/%s = %v, want %d", set_id, set_code, up_values, up_count)
		}
	}
	for _, value := range down_values {
		if value != fmt.Sprint(down_count) {
			t.Fatalf("down counters at %s/%s = %v, want %d", set_id, set_code, down_values, down_count)
		}
	}
}

func Test_bbs_vote_toggle_transition_matrix(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso", "mysql"} {
		t.Run(backend, func(t *testing.T) {
			db := new_record_db(t, backend)
			defer tool.DB_close(db)
			Bbs_vote_test_seed_board(db, t, "80")
			Bbs_vote_test_seed_post(db, t, "80", "1", "admin_ip")
			bbs_comment_seed_parent(db, t, "80", "1", "1", "admin_ip")

			config := Bbs_vote_test_config("admin_ip")
			steps := []struct {
				vote_type string
				want      string
			}{
				{vote_type: "up", want: "up"},
				{vote_type: "up", want: "none"},
				{vote_type: "down", want: "down"},
				{vote_type: "down", want: "none"},
				{vote_type: "up", want: "up"},
				{vote_type: "down", want: "down"},
				{vote_type: "up", want: "up"},
			}

			for _, target := range []string{"post", "comment"} {
				for _, step := range steps {
					result := Bbs_vote_test_call(t, target, config, "80", "1", "1", step.vote_type)
					if result["response"] != "ok" {
						t.Fatalf("%s vote %s: %#v", target, step.vote_type, result)
					}
					data, ok := result["data"].(map[string]any)
					if !ok || data["vote"] != step.want {
						t.Fatalf("%s vote %s result = %#v, want state %s", target, step.vote_type, result["data"], step.want)
					}
					Bbs_vote_test_assert_state(t, db, func() string {
						if target == "post" {
							return "80"
						}
						return "80-1"
					}(), "1", step.want, data["up_count"].(int), data["down_count"].(int))

					if target == "post" {
						read_result := route.Api_bbs_w_tabom(config, "80", "1")
						if read_result["response"] != "ok" || read_result["my_vote"] != step.want {
							t.Fatalf("read post vote = %#v, want %s", read_result, step.want)
						}
					} else if got := Bbs_vote_test_comment_data(t, config, "80", "1", "1")["my_vote"]; got != step.want {
						t.Fatalf("read comment my_vote = %q, want %s", got, step.want)
					}
				}
			}

			invalid := route.Api_bbs_w_tabom_post(config, "80", "1", "UP")
			if invalid["response"] != "error" || invalid["data"] != "vote_type" {
				t.Fatalf("invalid post vote = %#v, want vote_type error", invalid)
			}
			invalid = route.Api_bbs_w_comment_tabom_post(config, "80", "1", "1", "")
			if invalid["response"] != "error" || invalid["data"] != "vote_type" {
				t.Fatalf("invalid comment vote = %#v, want vote_type error", invalid)
			}
		})
	}
}

func Bbs_vote_test_seed_board(db *sql.DB, t *testing.T, set_id string) {
	t.Helper()
	Seed_bbs_subscription_board(db, t, set_id, "normal")
}

func Test_bbs_vote_legacy_rows_and_counter_reconciliation(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso", "mysql"} {
		t.Run(backend, func(t *testing.T) {
			db := new_record_db(t, backend)
			defer tool.DB_close(db)
			Bbs_vote_test_seed_board(db, t, "81")
			Bbs_vote_test_seed_post(db, t, "81", "1", "admin_ip")
			Bbs_vote_test_seed_post(db, t, "81", "2", "admin_ip")
			Bbs_vote_test_seed_post(db, t, "81", "3", "admin_ip")
			bbs_comment_seed_parent(db, t, "81", "1", "1", "admin_ip")

			seed_vote := func(list_name string, user string, set_id string, set_code string) {
				t.Helper()
				if _, err := db.Exec(`insert into bbs_data (set_name, set_data, set_id, set_code) values (?, ?, ?, ?)`, list_name, user, set_id, set_code); err != nil {
					t.Fatalf("seed legacy %s vote for %q at %s/%s: %v", list_name, user, set_id, set_code, err)
				}
			}
			seed_counter := func(name string, value string, set_id string, set_code string) {
				t.Helper()
				if _, err := db.Exec(`insert into bbs_data (set_name, set_data, set_id, set_code) values (?, ?, ?, ?)`, name, value, set_id, set_code); err != nil {
					t.Fatalf("seed legacy counter %s: %v", name, err)
				}
			}

			for _, vote := range []struct {
				list_name string
				user      string
				set_id    string
				set_code  string
			}{
				{"tabom_list", "admin_ip", "81", "1"},
				{"tabom_list", "admin_ip", "81", "1"},
				{"tabom_down_list", "admin_ip", "81", "1"},
				{"tabom_list", "other_voter", "81", "1"},
				{"tabom_list", "other_voter", "81", "1"},
				{"tabom_down_list", "down_voter", "81", "1"},
				{"tabom_down_list", "down_voter", "81", "1"},
				{"tabom_list", "admin_ip", "81-1", "1"},
				{"tabom_list", "admin_ip", "81-1", "1"},
				{"tabom_list", "other_voter", "81-1", "1"},
				{"tabom_list", "other_voter", "81-1", "1"},
				{"tabom_down_list", "down_voter", "81-1", "1"},
				{"tabom_down_list", "down_voter", "81-1", "1"},
				{"tabom_list", "unrelated_voter", "81", "2"},
				{"tabom_list", "other_voter_no_counter", "81", "3"},
			} {
				seed_vote(vote.list_name, vote.user, vote.set_id, vote.set_code)
			}
			for _, counter := range []struct {
				name     string
				value    string
				set_id   string
				set_code string
			}{
				{"tabom_count", "malformed", "81", "1"},
				{"tabom_count", "-99", "81", "1"},
				{"tabom_down_count", "800", "81", "1"},
				{"tabom_down_count", "NaN", "81", "1"},
				{"tabom_count", "777", "81-1", "1"},
				{"tabom_down_count", "888", "81-1", "1"},
				{"tabom_count", "41", "81", "2"},
			} {
				seed_counter(counter.name, counter.value, counter.set_id, counter.set_code)
			}

			post_result := route.Api_bbs_w_tabom_post(Bbs_vote_test_config("admin_ip"), "81", "1", "up")
			if post_result["response"] != "ok" {
				t.Fatalf("repair legacy post votes: %#v", post_result)
			}
			post_data := post_result["data"].(map[string]any)
			if post_data["vote"] != "up" || post_data["up_count"] != 2 || post_data["down_count"] != 1 {
				t.Fatalf("repaired post state = %#v, want up with counts 2/1", post_data)
			}
			if got := Bbs_vote_test_count(t, db, "tabom_list", "81", "1"); got != 3 {
				t.Fatalf("post up rows = %d, want current + two legacy other-user rows", got)
			}
			if got := Bbs_vote_test_count(t, db, "tabom_down_list", "81", "1"); got != 2 {
				t.Fatalf("post down rows = %d, want both legacy rows for unrelated down voter", got)
			}
			if got := Bbs_vote_test_counter_values(t, db, "tabom_count", "81", "1"); strings.Join(got, ",") != "2,2" {
				t.Fatalf("reconciled duplicate post up counters = %v, want [2 2]", got)
			}
			if got := Bbs_vote_test_counter_values(t, db, "tabom_down_count", "81", "1"); strings.Join(got, ",") != "1,1" {
				t.Fatalf("reconciled duplicate post down counters = %v, want [1 1]", got)
			}

			comment_result := route.Api_bbs_w_comment_tabom_post(Bbs_vote_test_config("admin_ip"), "81", "1", "1", "up")
			if comment_result["response"] != "ok" {
				t.Fatalf("repair legacy comment votes: %#v", comment_result)
			}
			comment_data := comment_result["data"].(map[string]any)
			if comment_data["vote"] != "none" || comment_data["up_count"] != 1 || comment_data["down_count"] != 1 {
				t.Fatalf("repaired comment state = %#v, want none with counts 1/1", comment_data)
			}
			if got := Bbs_vote_test_count(t, db, "tabom_list", "81-1", "1"); got != 2 {
				t.Fatalf("comment up rows = %d, want two duplicate rows for unrelated voter", got)
			}
			if got := Bbs_vote_test_count(t, db, "tabom_down_list", "81-1", "1"); got != 2 {
				t.Fatalf("comment down rows = %d, want two duplicate rows for unrelated voter", got)
			}
			if got := Bbs_vote_test_counter_values(t, db, "tabom_count", "81-1", "1"); strings.Join(got, ",") != "1" {
				t.Fatalf("reconciled comment up counter = %v, want [1]", got)
			}
			if got := Bbs_vote_test_counter_values(t, db, "tabom_down_count", "81-1", "1"); strings.Join(got, ",") != "1" {
				t.Fatalf("reconciled comment down counter = %v, want [1]", got)
			}

			if got := Bbs_vote_test_count(t, db, "tabom_list", "81", "2"); got != 1 {
				t.Fatalf("unrelated post up rows changed: %d, want 1", got)
			}
			if got := Bbs_vote_test_counter_values(t, db, "tabom_count", "81", "2"); strings.Join(got, ",") != "41" {
				t.Fatalf("unrelated post counter changed: %v, want [41]", got)
			}
			missing_counter_result := route.Api_bbs_w_tabom_post(Bbs_vote_test_config("admin_ip"), "81", "3", "up")
			if missing_counter_result["response"] != "ok" {
				t.Fatalf("repair missing post counter: %#v", missing_counter_result)
			}
			Bbs_vote_test_assert_totals(t, db, "81", "3", 2, 0)
			if got := Bbs_vote_test_counter_values(t, db, "tabom_count", "81", "3"); strings.Join(got, ",") != "2" {
				t.Fatalf("missing post counter recreated as %v, want [2]", got)
			}
			if got := Bbs_vote_test_counter_values(t, db, "tabom_down_count", "81", "3"); len(got) != 0 {
				t.Fatalf("zero-valued missing down counter was unnecessarily created: %v", got)
			}
		})
	}
}

func Test_bbs_vote_updates_existing_recommendation_scores(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		t.Run(backend, func(t *testing.T) {
			db := new_record_db(t, backend)
			defer tool.DB_close(db)
			Bbs_global_excellent_seed_board(t, db, "87", "normal", "0")
			Bbs_global_excellent_seed_threshold(t, db, "0")
			Bbs_global_excellent_seed_post(t, db, "87", "1", "Scored post", "admin_ip", 0, 0, false, false)
			Bbs_vote_test_seed_post_comment_lock(db, t, "87", "1")
			bbs_comment_seed_parent(db, t, "87", "1", "1", "admin_ip")
			bbs_best_comment_set_value(db, t, "87", "bbs_best_comment_enabled", "1")
			bbs_best_comment_set_value(db, t, "87", "bbs_best_comment_min_score", "1")

			config := Bbs_vote_test_config("admin_ip")
			if got := Bbs_global_excellent_has_post(t, route.Api_bbs(config, "", "1", "excellent"), "87", "1"); got {
				t.Fatal("zero-score post appeared in excellent recommendations")
			}
			best_before := route.Api_bbs_w_comment_best(config, "87", "1")["data"].([]map[string]string)
			if len(best_before) != 0 {
				t.Fatalf("zero-score comment appeared in BEST: %#v", best_before)
			}

			if result := route.Api_bbs_w_tabom_post(config, "87", "1", "up"); result["response"] != "ok" {
				t.Fatalf("post upvote: %#v", result)
			}
			if !Bbs_global_excellent_has_post(t, route.Api_bbs(config, "", "1", "excellent"), "87", "1") {
				t.Fatal("upvoted post did not enter existing excellent ranking")
			}
			if result := route.Api_bbs_w_comment_tabom_post(config, "87", "1", "1", "up"); result["response"] != "ok" {
				t.Fatalf("comment upvote: %#v", result)
			}
			best_after_up := route.Api_bbs_w_comment_best(config, "87", "1")["data"].([]map[string]string)
			if len(best_after_up) != 1 || best_after_up[0]["best_score"] != "1" {
				t.Fatalf("upvoted comment BEST result = %#v, want score 1", best_after_up)
			}

			if result := route.Api_bbs_w_tabom_post(config, "87", "1", "down"); result["response"] != "ok" {
				t.Fatalf("post vote switch: %#v", result)
			}
			if Bbs_global_excellent_has_post(t, route.Api_bbs(config, "", "1", "excellent"), "87", "1") {
				t.Fatal("downvoted post remained in excellent ranking")
			}
			if result := route.Api_bbs_w_comment_tabom_post(config, "87", "1", "1", "down"); result["response"] != "ok" {
				t.Fatalf("comment vote switch: %#v", result)
			}
			best_after_down := route.Api_bbs_w_comment_best(config, "87", "1")["data"].([]map[string]string)
			if len(best_after_down) != 0 {
				t.Fatalf("downvoted comment remained in BEST: %#v", best_after_down)
			}
		})
	}
}

func Bbs_vote_test_seed_post_comment_lock(db *sql.DB, t *testing.T, set_id string, set_code string) {
	t.Helper()
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values ('comment_count', ?, ?, '0')`, set_code, set_id); err != nil {
		t.Fatalf("seed post comment lock %s/%s: %v", set_id, set_code, err)
	}
}

func Test_bbs_vote_authorization_and_target_visibility(t *testing.T) {
	db := new_record_db(t, "sqlite")
	defer tool.DB_close(db)
	Bbs_vote_test_seed_board(db, t, "82")
	Bbs_vote_test_seed_post(db, t, "82", "1", "admin_ip")
	bbs_comment_seed_parent(db, t, "82", "1", "1", "admin_ip")

	if _, err := db.Exec(`insert into alist (name, acl) values ('vote_member_group', 'view')`); err != nil {
		t.Fatalf("seed view permission: %v", err)
	}
	if _, err := db.Exec(`insert into alist (name, acl) values ('vote_member_group', 'bbs_use')`); err != nil {
		t.Fatalf("seed board permission: %v", err)
	}
	if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', 'vote_member', 'vote_member_group')`); err != nil {
		t.Fatalf("seed limited viewer: %v", err)
	}
	member := Bbs_vote_test_config("vote_member")
	admin := Bbs_vote_test_config("admin_ip")

	if _, err := db.Exec(`insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_only_my_data_view_acl', '', '82', 'owner')`); err != nil {
		t.Fatalf("set private post ACL: %v", err)
	}
	private_result := route.Api_bbs_w_tabom_post(member, "82", "1", "up")
	if private_result["response"] != "require auth" {
		t.Fatalf("private post vote = %#v, want require auth", private_result)
	}
	if count := Bbs_vote_test_count(t, db, "tabom_list", "82", "1"); count != 0 {
		t.Fatalf("private post vote created %d rows", count)
	}
	if _, err := db.Exec(`update bbs_set set set_data = 'normal' where set_name = 'bbs_only_my_data_view_acl' and set_id = '82'`); err != nil {
		t.Fatalf("restore post ACL: %v", err)
	}

	if _, err := db.Exec(`insert into bbs_data (set_name, set_data, set_id, set_code) values ('blind', 'O', '82', '1')`); err != nil {
		t.Fatalf("blind post: %v", err)
	}
	blind_post_result := route.Api_bbs_w_tabom_post(member, "82", "1", "up")
	if blind_post_result["response"] != "require auth" {
		t.Fatalf("blind post vote = %#v, want require auth", blind_post_result)
	}
	if _, err := db.Exec(`delete from bbs_data where set_name = 'blind' and set_id = '82' and set_code = '1'`); err != nil {
		t.Fatalf("unblind post: %v", err)
	}

	if _, err := db.Exec(`insert into bbs_data (set_name, set_data, set_id, set_code) values ('blind', 'O', '82-1', '1')`); err != nil {
		t.Fatalf("blind comment: %v", err)
	}
	blind_comment_result := route.Api_bbs_w_comment_tabom_post(member, "82", "1", "1", "up")
	if blind_comment_result["response"] != "require auth" {
		t.Fatalf("blind comment vote = %#v, want require auth", blind_comment_result)
	}
	if count := Bbs_vote_test_count(t, db, "tabom_list", "82-1", "1"); count != 0 {
		t.Fatalf("blind comment vote created %d rows", count)
	}
	blind_comment_html := route.View_bbs_in_w_comment(db, member, "82", "1", "", 1)
	if strings.Contains(blind_comment_html, `/comment/1/tabom`) {
		t.Fatal("blind comment vote controls were rendered to a non-manager")
	}

	invalid_post := route.Api_bbs_w_tabom_post(admin, "82", "1", "neutral")
	invalid_comment := route.Api_bbs_w_comment_tabom_post(admin, "82", "1", "1", "neutral")
	if invalid_post["response"] != "error" || invalid_post["data"] != "vote_type" || invalid_comment["response"] != "error" || invalid_comment["data"] != "vote_type" {
		t.Fatalf("invalid vote responses = post %#v comment %#v", invalid_post, invalid_comment)
	}

	if _, err := db.Exec(`delete from bbs_data where set_name = 'blind' and set_id = '82-1' and set_code = '1'`); err != nil {
		t.Fatalf("unblind comment: %v", err)
	}
	if _, err := db.Exec(`update bbs_data set set_data = '' where set_name = 'comment' and set_id = '82-1' and set_code = '1'`); err != nil {
		t.Fatalf("delete comment body: %v", err)
	}
	deleted_result := route.Api_bbs_w_comment_tabom_post(admin, "82", "1", "1", "up")
	if deleted_result["response"] != "not exist" || deleted_result["data"] != "comment" {
		t.Fatalf("deleted comment vote = %#v, want not exist", deleted_result)
	}
	if count := Bbs_vote_test_count(t, db, "tabom_list", "82-1", "1"); count != 0 {
		t.Fatalf("deleted comment vote created %d rows", count)
	}

	if _, err := db.Exec(`update bbs_set set set_data = 'owner' where set_name = 'bbs_view_acl' and set_id = '82'`); err != nil {
		t.Fatalf("make board private: %v", err)
	}
	private_board := route.Api_bbs_w_comment_tabom_post(member, "82", "1", "1", "down")
	if private_board["response"] != "require auth" {
		t.Fatalf("private board comment vote = %#v, want require auth", private_board)
	}
}

func Test_bbs_vote_database_errors_roll_back_transitions(t *testing.T) {
	for _, target := range []string{"post", "comment"} {
		t.Run(target, func(t *testing.T) {
			db := new_record_db(t, "sqlite")
			defer tool.DB_close(db)
			Bbs_vote_test_seed_board(db, t, "83")
			Bbs_vote_test_seed_post(db, t, "83", "1", "admin_ip")
			bbs_comment_seed_parent(db, t, "83", "1", "1", "admin_ip")
			config := Bbs_vote_test_config("admin_ip")

			if result := Bbs_vote_test_call(t, target, config, "83", "1", "1", "up"); result["response"] != "ok" {
				t.Fatalf("initial vote: %#v", result)
			}
			set_id := "83"
			if target == "comment" {
				set_id = "83-1"
			}
			trigger_sql := fmt.Sprintf(`create trigger reject_down_vote before insert on bbs_data when NEW.set_name = 'tabom_down_list' and NEW.set_id = '%s' and NEW.set_code = '1' begin select raise(abort, 'rejected test vote'); end`, set_id)
			if _, err := db.Exec(trigger_sql); err != nil {
				t.Fatalf("install rollback trigger: %v", err)
			}

			failed := Bbs_vote_test_call(t, target, config, "83", "1", "1", "down")
			if failed["response"] != "error" || failed["data"] != "database" {
				t.Fatalf("failed transition response = %#v, want database error", failed)
			}
			Bbs_vote_test_assert_state(t, db, set_id, "1", "up", 1, 0)
		})
	}
}

func Test_bbs_vote_revalidates_delete_and_blind_state_under_lock(t *testing.T) {
	for _, scenario := range []struct {
		name          string
		set_id        string
		target        string
		mutation_sql  string
		want_response string
		want_data     string
	}{
		{
			name:          "post deleted after precheck",
			set_id:        "84",
			target:        "post",
			mutation_sql:  `delete from bbs_data where set_name = 'title' and set_id = '84' and set_code = '1'`,
			want_response: "not exist",
			want_data:     "post",
		},
		{
			name:          "comment deleted after precheck",
			set_id:        "84-1",
			target:        "comment",
			mutation_sql:  `delete from bbs_data where set_name = 'comment' and set_id = '84-1' and set_code = '1'`,
			want_response: "not exist",
			want_data:     "comment",
		},
		{
			name:          "comment becomes blind after precheck",
			set_id:        "84-1",
			target:        "comment",
			mutation_sql:  `insert into bbs_data (set_name, set_data, set_id, set_code) values ('blind', 'O', '84-1', '1')`,
			want_response: "require auth",
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			db := new_record_db(t, "sqlite")
			defer tool.DB_close(db)
			Bbs_vote_test_seed_board(db, t, "84")
			Bbs_vote_test_seed_post(db, t, "84", "1", "admin_ip")
			bbs_comment_seed_parent(db, t, "84", "1", "1", "admin_ip")
			if scenario.name == "comment becomes blind after precheck" {
				if _, err := db.Exec(`insert into alist (name, acl) values ('vote_member_group', 'view')`); err != nil {
					t.Fatalf("seed view permission: %v", err)
				}
				if _, err := db.Exec(`insert into alist (name, acl) values ('vote_member_group', 'bbs_use')`); err != nil {
					t.Fatalf("seed BBS permission: %v", err)
				}
				if _, err := db.Exec(`insert into user_set (name, id, data) values ('acl', 'vote_member', 'vote_member_group')`); err != nil {
					t.Fatalf("seed vote member: %v", err)
				}
			}

			// The existing local transaction lock changes comment_count before
			// invoking the callback. The trigger simulates a concurrent delete or
			// blind operation that committed after the API's early check but before
			// the locked target revalidation.
			trigger_sql := `create trigger mutate_vote_target after update of set_data on bbs_data when NEW.set_name = 'comment_count' and NEW.set_id = '84' and NEW.set_code = '1' and NEW.set_data = '0' || char(1) begin ` + scenario.mutation_sql + `; end`
			if _, err := db.Exec(trigger_sql); err != nil {
				t.Fatalf("install target mutation trigger: %v", err)
			}
			config := Bbs_vote_test_config("admin_ip")
			if scenario.name == "comment becomes blind after precheck" {
				config = Bbs_vote_test_config("vote_member")
			}
			result := Bbs_vote_test_call(t, scenario.target, config, "84", "1", "1", "up")
			if result["response"] != scenario.want_response || (scenario.want_data != "" && result["data"] != scenario.want_data) {
				t.Fatalf("locked target validation = %#v, want %s/%s", result, scenario.want_response, scenario.want_data)
			}
			if got := Bbs_vote_test_count(t, db, "tabom_list", scenario.set_id, "1"); got != 0 {
				t.Fatalf("invalidated target received %d upvote rows", got)
			}
		})
	}
}

func Test_bbs_vote_concurrent_toggles_are_serialized(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		t.Run(backend, func(t *testing.T) {
			db := new_record_db(t, backend)
			defer tool.DB_close(db)
			Bbs_vote_test_seed_board(db, t, "85")
			Bbs_vote_test_seed_post(db, t, "85", "1", "admin_ip")
			bbs_comment_seed_parent(db, t, "85", "1", "1", "admin_ip")

			for _, target := range []string{"post", "comment"} {
				const toggle_count = 10
				var group sync.WaitGroup
				errors := make(chan map[string]any, toggle_count)
				for index := 0; index < toggle_count; index++ {
					group.Add(1)
					go func() {
						defer group.Done()
						result := Bbs_vote_test_call(t, target, Bbs_vote_test_config("admin_ip"), "85", "1", "1", "up")
						if result["response"] != "ok" {
							errors <- result
						}
					}()
				}
				group.Wait()
				close(errors)
				for result := range errors {
					t.Errorf("concurrent %s toggle failed: %#v", target, result)
				}
				set_id := "85"
				if target == "comment" {
					set_id = "85-1"
				}
				Bbs_vote_test_assert_state(t, db, set_id, "1", "none", 0, 0)
			}

			const distinct_voters = 6
			for index := 0; index < distinct_voters; index++ {
				Bbs_vote_test_seed_owner(db, t, fmt.Sprintf("vote_concurrent_%d", index))
			}
			var group sync.WaitGroup
			errors := make(chan map[string]any, distinct_voters*2)
			for index := 0; index < distinct_voters; index++ {
				user := fmt.Sprintf("vote_concurrent_%d", index)
				for _, target := range []string{"post", "comment"} {
					group.Add(1)
					go func(target string, user string) {
						defer group.Done()
						result := Bbs_vote_test_call(t, target, Bbs_vote_test_config(user), "85", "1", "1", "up")
						if result["response"] != "ok" {
							errors <- result
						}
					}(target, user)
				}
			}
			group.Wait()
			close(errors)
			for result := range errors {
				t.Errorf("concurrent distinct vote failed: %#v", result)
			}
			Bbs_vote_test_assert_totals(t, db, "85", "1", distinct_voters, 0)
			Bbs_vote_test_assert_totals(t, db, "85-1", "1", distinct_voters, 0)
		})
	}
}

func Test_bbs_vote_registered_http_routes_and_comment_page_redirect(t *testing.T) {
	db := new_record_db(t, "sqlite")
	defer tool.DB_close(db)
	Bbs_vote_test_seed_board(db, t, "86")
	Bbs_vote_test_seed_post(db, t, "86", "1", "admin_ip")
	for code := 1; code <= 51; code++ {
		bbs_comment_seed_parent(db, t, "86", "1", fmt.Sprint(code), "admin_ip")
	}

	previous_mode := gin.Mode()
	gin.SetMode(gin.TestMode)
	t.Cleanup(func() { gin.SetMode(previous_mode) })
	t.Setenv("NAMU_SESSION_KEY", "bbs_vote_http_test_session_key")
	session_middleware := tool.Session_middleware()
	admin_cookie := Bbs_subscription_http_session_cookie(t, "admin_ip", session_middleware)

	http_router := gin.New()
	http_router.Use(session_middleware)
	http_router.Use(route.Site_view_middleware())
	route.Register_bbs_vote_routes(http_router)
	http_router.GET("/bbs/w/:set_id/:set_code", func(c *gin.Context) {
		data := route.View_bbs_in_w(c, route.Make_route_config(c), c.Param("set_id"), c.Param("set_code"))
		route.Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	http_router.GET("/bbs/w/:set_id/:set_code/page/:page", func(c *gin.Context) {
		data := route.View_bbs_in_w(c, route.Make_route_config(c), c.Param("set_id"), c.Param("set_code"))
		route.Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	http_router.GET("/bbs/w/:set_id/:set_code/page/:page/comment/:comment_select", func(c *gin.Context) {
		data := route.View_bbs_in_w(c, route.Make_route_config(c), c.Param("set_id"), c.Param("set_code"))
		route.Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	http_router.GET("/bbs/w/:set_id/:set_code/comment/:comment_select", func(c *gin.Context) {
		data := route.View_bbs_in_w(c, route.Make_route_config(c), c.Param("set_id"), c.Param("set_code"))
		route.Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
	server := httptest.NewServer(http_router)
	t.Cleanup(server.Close)

	status, page := Bbs_subscription_http_request(t, server.URL, http.MethodGet, "/bbs/w/86/1/page/2", admin_cookie, nil)
	if status != http.StatusOK || !strings.Contains(string(page), `action="/bbs/w/86/1/comment/51/tabom"`) || !strings.Contains(string(page), `name="vote_type" value="up"`) {
		t.Fatalf("page-two comment vote controls = %d %s", status, page)
	}

	status, confirmation := Bbs_subscription_http_request(t, server.URL, http.MethodGet, "/bbs/w/86/1/comment_tabom/51/up", admin_cookie, nil)
	if status != http.StatusOK || !strings.Contains(string(confirmation), `value="up"`) {
		t.Fatalf("legacy GET confirmation route = %d %s", status, confirmation)
	}
	if count := Bbs_vote_test_count(t, db, "tabom_list", "86-1", "51"); count != 0 {
		t.Fatalf("legacy GET confirmation changed state: %d upvote rows", count)
	}

	form := url.Values{"vote_type": {"up"}}
	status, body, location := Bbs_subscription_http_request_with_location(t, server.URL, http.MethodPost, "/bbs/w/86/1/comment/51/tabom", admin_cookie, form)
	if status != http.StatusFound || location != "/bbs/w/86/1/page/2/comment/51#51" {
		t.Fatalf("comment page-two POST redirect = %d location=%q body=%s", status, location, body)
	}
	if count := Bbs_vote_test_count(t, db, "tabom_list", "86-1", "51"); count != 1 {
		t.Fatalf("comment POST created %d upvote rows, want one", count)
	}

	status, selected_page := Bbs_subscription_http_request(t, server.URL, http.MethodGet, "/bbs/w/86/1/page/2/comment/51", admin_cookie, nil)
	if status != http.StatusOK || !strings.Contains(string(selected_page), `aria-pressed="true"`) || !strings.Contains(string(selected_page), tool.Get_language(db, "cancel_upvote", true)) {
		t.Fatalf("selected comment vote state = %d %s", status, selected_page)
	}

	status, post_page := Bbs_subscription_http_request(t, server.URL, http.MethodGet, "/bbs/w/86/1", admin_cookie, nil)
	if status != http.StatusOK || !strings.Contains(string(post_page), `action="/bbs/w/86/1/tabom"`) {
		t.Fatalf("post vote controls = %d %s", status, post_page)
	}
	status, _, location = Bbs_subscription_http_request_with_location(t, server.URL, http.MethodPost, "/bbs/w/86/1/tabom", admin_cookie, form)
	if status != http.StatusFound || location != "/bbs/w/86/1" {
		t.Fatalf("post upvote redirect = %d location=%q", status, location)
	}
	status, post_page = Bbs_subscription_http_request(t, server.URL, http.MethodGet, "/bbs/w/86/1", admin_cookie, nil)
	if status != http.StatusOK || !strings.Contains(string(post_page), tool.Get_language(db, "cancel_upvote", true)) {
		t.Fatalf("selected post vote state = %d %s", status, post_page)
	}

	status, _ = Bbs_subscription_http_request(t, server.URL, http.MethodGet, "/bbs/w/86/1/tabom", admin_cookie, nil)
	if status == http.StatusOK || Bbs_vote_test_count(t, db, "tabom_list", "86", "1") != 1 {
		t.Fatalf("GET on POST-only post vote route changed state or was accepted: status=%d", status)
	}
}
