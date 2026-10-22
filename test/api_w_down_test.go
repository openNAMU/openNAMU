package test

import (
	"database/sql"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

// Setup_api_w_down_db builds a schema and seeds groups, users, docs, and one
// view-restricted doc used by Test_api_w_down.
func Setup_api_w_down_db(t *testing.T) *sql.DB {
	t.Helper()

	db_dir := filepath.Join(t.TempDir(), "db")
	os.MkdirAll(db_dir, 0755)
	os.Setenv("NAMU_DB", db_dir)
	os.Setenv("NAMU_DB_TYPE", "sqlite")
	defer os.Unsetenv("NAMU_DB")
	defer os.Unsetenv("NAMU_DB_TYPE")

	tool.DB_init()
	db := tool.DB_connect()

	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("seed %q: %v", q, err)
		}
	}

	// owner is admin (full access); regular has only the view permission
	// (which implies site_view). anon_ip is intentionally unbound.
	exec("insert into alist (name, acl) values ('owner','owner')")
	exec("insert into alist (name, acl) values ('regular','view')")
	exec("insert into user_set (name, id, data) values ('acl','admin_ip','owner')")
	exec("insert into user_set (name, id, data) values ('acl','regular_ip','regular')")

	docs := []string{
		"Test",
		"Test/Child",
		"Test/Child/Deep",
		"Testing/Child",
		"Te_st",
		"Te_st/RealChild",
		"TeXst/FakeChild",
		"Te%st",
		"Te%st/RealChild",
		"TeAst/FakeChild",
		"Empty",
		"AclTest",
		"AclTest/PublicChild",
		"AclTest/RestrictedChild",
		"NoParent/Child",
	}
	for _, d := range docs {
		exec("insert into data (title, data) values (?, '')", d)
	}

	// Restrict AclTest/RestrictedChild to the owner group for view.
	exec("insert into acl (title, data, type) values ('AclTest/RestrictedChild','owner','view')")

	return db
}

func Test_api_w_down(t *testing.T) {
	db := Setup_api_w_down_db(t)
	defer tool.DB_close(db)

	admin := tool.Config{IP: "admin_ip", Session: &emptySession{}}
	regular := tool.Config{IP: "regular_ip", Session: &emptySession{}}
	anon := tool.Config{IP: "anon_ip", Session: &emptySession{}}

	setEq := func(name string, got map[string]any, want []string) {
		if got["response"] != "ok" {
			t.Fatalf("%s: response = %v, want ok", name, got["response"])
		}
		got_data := got["data"].([]string)
		if !reflect.DeepEqual(got_data, want) {
			t.Errorf("%s: data = %v, want %v", name, got_data, want)
		}
	}

	// 1. Only descendants returned; parent and similar-prefix docs excluded.
	setEq("descendants", route.Api_w_down(admin, "Test"), []string{
		"Test/Child",
		"Test/Child/Deep",
	})

	// 2. Parent doc need not exist; children still returned.
	setEq("no parent", route.Api_w_down(admin, "NoParent"), []string{
		"NoParent/Child",
	})

	// 3. Nested descendants of a child are returned.
	setEq("nested", route.Api_w_down(admin, "Test/Child"), []string{
		"Test/Child/Deep",
	})

	// 4/5. Percent and underscore treated literally; similar names excluded.
	setEq("escape underscore", route.Api_w_down(admin, "Te_st"), []string{
		"Te_st/RealChild",
	})
	setEq("escape percent", route.Api_w_down(admin, "Te%st"), []string{
		"Te%st/RealChild",
	})

	// 6. No children -> empty list.
	setEq("no children", route.Api_w_down(admin, "Empty"), []string{})

	// 7a. Admin can view the restricted child.
	setEq("admin restricted", route.Api_w_down(admin, "AclTest"), []string{
		"AclTest/PublicChild",
		"AclTest/RestrictedChild",
	})

	// 7b. Regular views the public child but not the restricted one.
	setEq("regular restricted", route.Api_w_down(regular, "AclTest"), []string{
		"AclTest/PublicChild",
	})

	// 7c. Unauthorized user -> require auth, empty data.
	out := route.Api_w_down(anon, "Test")
	if out["response"] != "require auth" {
		t.Fatalf("anonymous: response = %v, want require auth", out["response"])
	}
	if len(out["data"].([]string)) != 0 {
		t.Fatalf("anonymous: data should be empty, got %v", out["data"].([]string))
	}
}
