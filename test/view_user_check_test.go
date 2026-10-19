package test

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/flosch/pongo2/v6"
	"github.com/gin-contrib/sessions"

	"opennamu/route"
	"opennamu/route/tool"
)

// Pongo_init_test registers the template filters that main.go::Pongo_init
// registers, so route.View_user_check can render templates standalone.
func init() {
	pongo2.RegisterFilter("md5_replace", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
		h := md5.Sum([]byte(in.String()))
		return pongo2.AsValue(hex.EncodeToString(h[:])), nil
	})
	pongo2.RegisterFilter("load_lang", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
		db := tool.DB_connect()
		defer tool.DB_close(db)
		return pongo2.AsValue(tool.Get_language(db, in.String(), false)), nil
	})
	pongo2.RegisterFilter("cut_100", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
		data := in.String()
		data = tool.Get_slice(data, 0, 100)
		return pongo2.AsValue(data), nil
	})
}

// emptySession implements sessions.Session so route.View_user_check can run without
// a real request/session. It returns nil from Get, causing the template to
// fall back to DB lookups only.
type emptySession struct{}

func (e *emptySession) ID() string { return "" }

func (e *emptySession) Get(key interface{}) interface{} {
	return nil
}

func (e *emptySession) Set(key, value interface{}) {}

func (e *emptySession) Delete(key interface{}) {}

func (e *emptySession) Clear() {}

func (e *emptySession) AddFlash(value interface{}, vars ...string) {}

func (e *emptySession) Flashes(vars ...string) []interface{} { return nil }

func (e *emptySession) Options(o sessions.Options) {}

func (e *emptySession) Save() error { return nil }

func Setup_view_db(t *testing.T) *sql.DB {
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
		_, err := db.Exec(query, args...)
		if err != nil {
			t.Fatalf("seed: %v", err)
		}
	}

	exec(`delete from alist where name in ('inspector_group','victim_group','open_group')`)
	exec(`delete from user_set where name = 'acl' and id in ('inspector','victim','open')`)
	exec(`delete from ua_d where name in ('victim','open')`)

	exec(`insert into alist (name, acl) values ('inspector_group','check')`)
	// victim only needs hideip (masks own data when inspected); check is the
	// inspector's permission, not the subject's.
	exec(`insert into alist (name, acl) values ('victim_group','hideip')`)
	exec(`insert into alist (name, acl) values ('open_group','check')`)

	exec(`insert into user_set (id, name, data) values ('inspector','acl','inspector_group')`)
	exec(`insert into user_set (id, name, data) values ('victim','acl','victim_group')`)
	exec(`insert into user_set (id, name, data) values ('open','acl','open_group')`)

	exec(`insert into ua_d (name, ip, ua, today, sub) values ('victim','192.168.0.10','Mozilla/5.0',?,'0')`, "20261008")
	exec(`insert into ua_d (name, ip, ua, today, sub) values ('victim','192.168.0.11','Mozilla/5.1',?,'0')`, "20261008")
	exec(`insert into ua_d (name, ip, ua, today, sub) values ('open','192.168.0.20','Mozilla/5.2',?,'0')`, "20261008")

	return db
}

func Test_View_user_check_masked_subject_normal(t *testing.T) {
	db := Setup_view_db(t)
	defer tool.DB_close(db)

	config := tool.Config{IP: "inspector", Session: &emptySession{}}
	out := route.View_user_check(config, "victim", "normal", "1", "")

	if !strings.Contains(out, "victim") {
		t.Fatalf("did not render subject (likely auth error): %q", out)
	}
	if strings.Contains(out, "192.168.0.10") || strings.Contains(out, "192.168.0.11") {
		t.Fatalf("masked subject still exposed real IP")
	}
}

func Test_View_user_check_unmasked_subject_normal(t *testing.T) {
	db := Setup_view_db(t)
	defer tool.DB_close(db)

	config := tool.Config{IP: "inspector", Session: &emptySession{}}
	out := route.View_user_check(config, "open", "normal", "1", "")

	if !strings.Contains(out, "open") {
		t.Fatalf("did not render subject (likely auth error): %q", out)
	}
	if !strings.Contains(out, "192.168.0.20") {
		t.Fatalf("unmasked subject did not expose IP")
	}
}

func Test_View_user_check_masked_subject_simple(t *testing.T) {
	db := Setup_view_db(t)
	defer tool.DB_close(db)

	config := tool.Config{IP: "inspector", Session: &emptySession{}}
	out := route.View_user_check(config, "victim", "simple", "1", "")

	if !strings.Contains(out, "victim") {
		t.Fatalf("did not render subject (likely auth error): %q", out)
	}
	if strings.Contains(out, "192.168.0.10") || strings.Contains(out, "192.168.0.11") {
		t.Fatalf("masked subject exposed IP in simple mode")
	}
}

func Test_View_user_check_unmasked_subject_simple(t *testing.T) {
	db := Setup_view_db(t)
	defer tool.DB_close(db)

	config := tool.Config{IP: "inspector", Session: &emptySession{}}
	out := route.View_user_check(config, "open", "simple", "1", "")

	if !strings.Contains(out, "open") {
		t.Fatalf("did not render subject (likely auth error): %q", out)
	}
	if !strings.Contains(out, "192.168.0.20") {
		t.Fatalf("unmasked subject did not expose IP in simple mode")
	}
}

func Test_View_user_check_requires_auth(t *testing.T) {
	db := Setup_view_db(t)
	defer tool.DB_close(db)

	config := tool.Config{IP: "anonymous", Session: &emptySession{}}
	out := route.View_user_check(config, "victim", "normal", "1", "")

	if strings.Contains(out, "192.168.0.11") {
		t.Fatalf("unpermitted inspector was able to see IP: %q", out)
	}
}

// Test_View_user_check_simple_ip_excludes_hideip verifies that an IP
// reverse-lookup in simple mode does not leak a hideip account's name,
// while still listing non-hideip accounts that share the IP.
func Test_View_user_check_simple_ip_excludes_hideip(t *testing.T) {
	db := Setup_view_db(t)
	defer tool.DB_close(db)

	// Make IP 192.168.0.10 shared by the hideip user "victim" and the
	// non-hideip user "open".
	if _, err := db.Exec(`insert into ua_d (name, ip, ua, today, sub) values ('open','192.168.0.10','Mozilla/5.3',?,'0')`, "20261008"); err != nil {
		t.Fatalf("seed shared ip: %v", err)
	}

	config := tool.Config{IP: "inspector", Session: &emptySession{}}
	out := route.View_user_check(config, "192.168.0.10", "simple", "1", "")

	if strings.Contains(out, "victim") {
		t.Fatalf("hideip account leaked through IP reverse lookup")
	}
	if !strings.Contains(out, "open") {
		t.Fatalf("non-hideip account not shown after IP reverse lookup")
	}
}

// Test_View_user_check_normal_ip_excludes_hideip verifies that an IP
// reverse-lookup in normal mode excludes hideip accounts entirely.
func Test_View_user_check_normal_ip_excludes_hideip(t *testing.T) {
	db := Setup_view_db(t)
	defer tool.DB_close(db)

	if _, err := db.Exec(`insert into ua_d (name, ip, ua, today, sub) values ('open','192.168.0.10','Mozilla/5.3',?,'0')`, "20261008"); err != nil {
		t.Fatalf("seed shared ip: %v", err)
	}

	config := tool.Config{IP: "inspector", Session: &emptySession{}}
	out := route.View_user_check(config, "192.168.0.10", "normal", "1", "")

	if strings.Contains(out, "Mozilla/5.0") {
		t.Fatalf("hideip account's UA leaked through IP reverse lookup")
	}
	if !strings.Contains(out, "Mozilla/5.3") {
		t.Fatalf("non-hideip account's UA not shown after IP reverse lookup")
	}
}
