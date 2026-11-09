package test

import (
	"encoding/json"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_Api_add_user_preserves_creation_and_invite_flow(t *testing.T) {
	db := Setup_rebuild_db(t)
	defer tool.DB_close(db)

	ordinary_result := route.Api_add_user(tool.Config{}, "new_user", "ordinary-password", "user@example.test", "")
	if ordinary_result["response"] != "ok" {
		t.Fatalf("ordinary registration response = %v, want ok", ordinary_result)
	}
	if !tool.Password_check(db, "new_user", "ordinary-password") {
		t.Fatal("ordinary registration did not store a verifiable password")
	}
	var email string
	if !tool.QueryRow_DB(db, "select data from user_set where id = ? and name = 'email'", []any{&email}, "new_user") {
		t.Fatal("ordinary registration did not store the email")
	}
	if email != "user@example.test" {
		t.Fatalf("stored email = %q, want user@example.test", email)
	}

	hashed_password := tool.Password_encode(db, "prehashed-password", "sha256")
	hash_result := route.Api_add_user_hash(tool.Config{}, "hashed_user", hashed_password, "", "sha256")
	if hash_result["response"] != "ok" {
		t.Fatalf("pre-hashed registration response = %v, want ok", hash_result)
	}
	if !tool.Password_check(db, "hashed_user", "prehashed-password") {
		t.Fatal("pre-hashed registration did not preserve the supplied encoding")
	}

	if _, err := db.Exec("insert into other (name, data, coverage) values ('invite_required', 'on', '')"); err != nil {
		t.Fatalf("enable invite requirement: %v", err)
	}
	invite_hash := tool.Invite_hash("invite-token")
	invite_data, err := json.Marshal(tool.Invite_record{Issuer: "admin", Date: tool.Get_time()})
	if err != nil {
		t.Fatalf("encode invite: %v", err)
	}
	if _, err := db.Exec("insert into user_set (id, name, data) values (?, 'invite', ?)", invite_hash, string(invite_data)); err != nil {
		t.Fatalf("insert invite: %v", err)
	}

	invalid_result := route.Api_add_user_invite(tool.Config{}, "invalid_invite_user", "password", "", "", "invalid-invite-hash")
	if invalid_result["response"] != "error" || invalid_result["data"] != "invite error" {
		t.Fatalf("invalid invite response = %v, want invite error", invalid_result)
	}
	var invalid_user_count int
	if err := db.QueryRow("select count(*) from user_set where id = ?", "invalid_invite_user").Scan(&invalid_user_count); err != nil {
		t.Fatalf("count invalid-invite user rows: %v", err)
	}
	if invalid_user_count != 0 {
		t.Fatalf("invalid invite created %d user rows", invalid_user_count)
	}

	invite_result := route.Api_add_user_invite(tool.Config{}, "invited_user", "invited-password", "", "", invite_hash)
	if invite_result["response"] != "ok" {
		t.Fatalf("valid invite response = %v, want ok", invite_result)
	}
	if !tool.Password_check(db, "invited_user", "invited-password") {
		t.Fatal("invited registration did not store a verifiable password")
	}
	var invite_count int
	if err := db.QueryRow("select count(*) from user_set where id = ? and name = 'invite'", invite_hash).Scan(&invite_count); err != nil {
		t.Fatalf("count invite rows: %v", err)
	}
	if invite_count != 0 {
		t.Fatalf("consumed invite row remains: count=%d", invite_count)
	}
}
