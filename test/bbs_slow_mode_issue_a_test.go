// Issue A — thread comment must NOT reset the topic author's board cooldown.
//
// Root cause: in a thread the topic's set_name='date' row is overwritten on
// every comment (see route/api_bbs_w_comment_post.go and
// route/api_thread_bbs_post.go: "update bbs_data set set_data = ? where
// set_name = 'date' ..."). The record-based board query attributes that
// (mutable) date to the topic author via user_id, so any commenter's write
// falsely extends the author's board cooldown.
//
// walltimeAt is provided by bbs_slow_mode_mariadb_test.go in this package
// (it takes a time.Time). All timestamps use the same local wall clock as
// tool.Get_time, so bss_parse_last_write interprets them consistently.
package test

import (
	"testing"
	"time"

	"opennamu/route"
	"opennamu/route/tool"
)

// Test_Bbs_slow_issue_a_thread_comment_does_not_reset_author reproduces the
// P0 attribution defect: userB's thread comment updates the topic's date and
// must NOT reset author userA's board cooldown.
func Test_Bbs_slow_issue_a_thread_comment_does_not_reset_author(t *testing.T) {
	db := new_slow_db(t)
	defer tool.DB_close(db)
	seedBbsNonExempt(db, t, "userA")
	seedBbsNonExempt(db, t, "userB")
	setSlowBoard(db, t, "-1", 60) // 60 s board-level cooldown on the thread

	// 1) User A created topic 1 two hours ago (date + user_id + an immutable
	//    initial comment authored by A).
	seedPostWrite(db, "userA", "-1", "1", walltimeAt(time.Now().Add(-2*time.Hour)))
	seedCommentWrite(db, "userA", "-1", "1", "1", walltimeAt(time.Now().Add(-2*time.Hour)))

	// A's real last write is 2h ago -> the 60 s cooldown is long expired.
	if rem, err := enforceOnce(db, "userA", "-1", "", 60, 0); rem != 0 {
		t.Fatalf("baseline: author cooldown should be expired before anyone comments; rem=%d err=%v", rem, err)
	}

	// 2) User B comments on A's topic. The thread write path inserts a comment
	//    record for B and ALSO updates the topic's date to now (side effect).
	seedCommentWrite(db, "userB", "-1", "1", "2", nowWall())
	if _, err := db.Exec(`update bbs_data set set_data = ? where set_name = 'date' and set_id = ? and set_code = ?`,
		nowWall(), "-1", "1"); err != nil {
		t.Fatalf("update topic date: %v", err)
	}

	// 3) User A tries another write. B's comment must NOT reset A's cooldown.
	rem, err := enforceOnce(db, "userA", "-1", "", 60, 0)
	if err != nil {
		t.Fatalf("enforceOnce error: %v", err)
	}
	if rem != 0 {
		t.Errorf("Issue A: userB's thread comment reset author userA's board cooldown; expected rem=0, got rem=%d", rem)
	}
}

// Test_Bbs_slow_issue_a_author_post_edit_does_not_refresh_creation_time verifies
// that mutating a post's display/sort date cannot replace its immutable creation
// timestamp in board cooldown attribution.
func Test_Bbs_slow_issue_a_author_post_edit_does_not_refresh_creation_time(t *testing.T) {
	db := new_slow_db(t)
	defer tool.DB_close(db)
	seedBbsNonExempt(db, t, "userA")
	setSlowBoard(db, t, "board", 60)

	created_at := walltimeAt(time.Now().Add(-2 * time.Hour))
	seedPostWrite(db, "userA", "board", "1", created_at)
	if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, tool.Bbs_post_creation_time_name, "1", "board", created_at); err != nil {
		t.Fatalf("insert creation timestamp: %v", err)
	}
	// An edit updates the mutable date but must not replace the creation time.
	if _, err := db.Exec(`update bbs_data set set_data = ? where set_name = 'date' and set_id = ? and set_code = ?`,
		nowWall(), "board", "1"); err != nil {
		t.Fatalf("update post date: %v", err)
	}

	rem, err := enforceOnce(db, "userA", "board", "", 60, 0)
	if err != nil {
		t.Fatalf("enforceOnce error: %v", err)
	}
	if rem != 0 {
		t.Errorf("post edit refreshed cooldown from mutable date; creation time is expired, got rem=%d", rem)
	}
}

// Test_Bbs_slow_issue_a_thread_comment_api_does_not_reset_author exercises the
// real thread write API: user B's comment updates the topic's mutable date but
// must not reset user A's board cooldown.
func Test_Bbs_slow_issue_a_moderator_post_edit_does_not_refresh_owner(t *testing.T) {
	db := apiSlowDB(t)
	defer tool.DB_close(db)
	const board, post, author = "21", "1", "userA"
	bbs_comment_seed_board(db, t, board, post, "old title", author)
	seedBbsNonExempt(db, t, author)
	bbsSlowBoardSeed(db, t, board, 60)

	created_at := walltimeAt(time.Now().Add(-2 * time.Hour))
	for _, value := range []struct {
		name string
		data string
	}{
		{"data", "old body"},
		{tool.Bbs_post_creation_time_name, created_at},
	} {
		if _, err := db.Exec(`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, ?, ?, ?)`, value.name, post, board, value.data); err != nil {
			t.Fatalf("seed post metadata %s: %v", value.name, err)
		}
	}
	if _, err := db.Exec(`update bbs_data set set_data = ? where set_name = 'date' and set_id = ? and set_code = ?`, created_at, board, post); err != nil {
		t.Fatalf("seed post date: %v", err)
	}

	moderator := route.Api_bbs_w_edit_post(tool.Config{IP: "admin_ip", Session: &emptySession{}}, board, post, "", "moderated title", "moderated body", "", "", "")
	if moderator["response"] != "ok" {
		t.Fatalf("moderator edit response=%v, want ok", moderator)
	}
	new_post := route.Api_bbs_w_edit_post(tool.Config{IP: author, Session: &emptySession{}}, board, "", "", "new title", "new body", "", "", "")
	if new_post["response"] != "ok" {
		t.Fatalf("moderator edit incorrectly reset owner's board cooldown: response=%v data=%v", new_post["response"], new_post["data"])
	}
	new_code, ok := new_post["data"].(string)
	if !ok || new_code == "" {
		t.Fatalf("new post code = %v, want non-empty string", new_post["data"])
	}
	if _, exists := tool.Get_bbs_data_value(db, board, new_code, tool.Bbs_post_creation_time_name); !exists {
		t.Fatalf("new top-level post did not persist immutable creation metadata")
	}
	read := route.Api_bbs_w(tool.Config{IP: "admin_ip", Session: &emptySession{}}, board, new_code)
	read_data, ok := read["data"].(map[string]string)
	if read["response"] != "ok" || !ok {
		t.Fatalf("post read response=%v data type=%T", read["response"], read["data"])
	}
	if _, visible := read_data[tool.Bbs_post_creation_time_name]; visible {
		t.Fatalf("internal creation metadata leaked through Api_bbs_w: %q", tool.Bbs_post_creation_time_name)
	}
	listing := route.Api_bbs(tool.Config{IP: "admin_ip", Session: &emptySession{}}, board, "1", "date")
	listed_posts, ok := listing["data"].([]map[string]string)
	if listing["response"] != "ok" || !ok {
		t.Fatalf("board listing response=%v data type=%T", listing["response"], listing["data"])
	}
	for _, listed_post := range listed_posts {
		if listed_post["set_code"] == new_code {
			if _, visible := listed_post[tool.Bbs_post_creation_time_name]; visible {
				t.Fatalf("internal creation metadata leaked through Api_bbs: %q", tool.Bbs_post_creation_time_name)
			}
			break
		}
	}
}

func Test_Bbs_slow_issue_a_thread_comment_api_does_not_reset_author(t *testing.T) {
	db := apiSlowDB(t)
	defer tool.DB_close(db)
	seedBbsNonExempt(db, t, "userA")
	seedBbsNonExempt(db, t, "userB")
	bbsSlowBoardSeed(db, t, "-1", 60)
	cfg := func(user string) tool.Config {
		return tool.Config{IP: user, Session: &emptySession{}}
	}

	created := route.Api_thread_bbs_post(cfg("userA"), "0", "", "topic body", "Issue A topic", "")
	if created["response"] != "ok" {
		t.Fatalf("topic creation response=%v, want ok", created)
	}
	topic_num, ok := created["topic_num"].(string)
	if !ok || topic_num == "" {
		t.Fatalf("topic_num = %v, want non-empty string", created["topic_num"])
	}
	old_wall := walltimeAt(time.Now().Add(-2 * time.Hour))
	if _, err := db.Exec(`update bbs_data set set_data = ? where set_name = 'date' and set_id = '-1' and set_code = ?`, old_wall, topic_num); err != nil {
		t.Fatalf("age topic date: %v", err)
	}
	if _, err := db.Exec(`update bbs_data set set_data = ? where set_name = 'comment_date' and set_id = ? and set_code = '1'`, old_wall, "-1-"+topic_num); err != nil {
		t.Fatalf("age initial topic comment: %v", err)
	}

	comment_b := route.Api_thread_bbs_post(cfg("userB"), topic_num, "", "comment by B", "", "")
	if comment_b["response"] != "ok" {
		t.Fatalf("user B comment response=%v, want ok", comment_b)
	}
	comment_a := route.Api_thread_bbs_post(cfg("userA"), topic_num, "", "comment by A", "", "")
	if comment_a["response"] != "ok" {
		t.Fatalf("user B's comment incorrectly reset user A's cooldown: response=%v data=%v", comment_a["response"], comment_a["data"])
	}
}
