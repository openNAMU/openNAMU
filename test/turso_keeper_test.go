package test

// B) Keeper extension to the Turso local embedded MVCC backend.
//
// Start_keeper/Stop_keeper previously handled sqlite only; turso now holds one
// warmed idle *sql.Conn so the local DB instance stays loaded in-process and
// requests after startup or idle periods attach cheaply instead of paying the
// full per-connection instance load (~27ms on this build's schema).
//
// Behavioral evidence used here (deterministic margins, not tight timing):
//   - right after Start_keeper(turso), a fresh DB_connect + SELECT 1 is fast
//     because the instance is already loaded (warm attach);
//   - the keeper key includes the DB type, so a sqlite -> turso switch on the
//     SAME path must rebind (a no-op keeper would leave the turso instance
//     cold and the first attach would be slow).

import (
	"testing"
	"time"

	"opennamu/route/tool"
)

// keeper_attach_ms measures a real physical connection cycle (not sql.Open,
// which is lazy): open pool -> SELECT 1 -> close.
func keeper_attach_ms(t *testing.T) float64 {
	t.Helper()
	db := tool.DB_connect()
	defer tool.DB_close(db)

	start := time.Now()
	var one int
	if err := db.QueryRow(`select 1`).Scan(&one); err != nil {
		t.Fatalf("attach query: %v", err)
	}
	return time.Since(start).Seconds() * 1000
}

// Test_keeper_turso_start_stop_idempotent: Start x2, Stop x2, re-arm; each
// phase must be error-free and the warm-attach effect must appear after Start.
func Test_keeper_turso_start_stop_idempotent(t *testing.T) {
	if !tool.Turso_platform_supported() {
		t.Skip("turso unsupported platform")
	}
	audit_setup(t, "turso")
	tool.Stop_keeper() // ensure OFF baseline (audit_setup may have started it)

	if err := tool.Start_keeper(); err != nil {
		t.Fatalf("keeper start: %v", err)
	}
	if err := tool.Start_keeper(); err != nil { // idempotent
		t.Fatalf("keeper start x2: %v", err)
	}

	warm := keeper_attach_ms(t)
	if warm > 10 {
		t.Errorf("attach after Start_keeper = %.2fms, want warm (<10ms) - held connection not warmed", warm)
	}
	t.Logf("turso keeper ON attach=%.2fms", warm)

	if err := tool.Stop_keeper(); err != nil {
		t.Fatalf("keeper stop: %v", err)
	}
	if err := tool.Stop_keeper(); err != nil { // idempotent
		t.Fatalf("keeper stop x2: %v", err)
	}

	// re-arm after stop
	if err := tool.Start_keeper(); err != nil {
		t.Fatalf("keeper re-start: %v", err)
	}
	if warm2 := keeper_attach_ms(t); warm2 > 10 {
		t.Errorf("attach after re-arm = %.2fms, want warm", warm2)
	}
	tool.Stop_keeper()
}

// Test_keeper_turso_rebind_on_type_change: switching sqlite -> turso on the
// SAME DB path must rebind the keeper (key includes the DB type). With a
// type-blind key the keeper would keep the sqlite connection and the turso
// instance would stay cold (~27ms attach on this schema).
func Test_keeper_turso_rebind_on_type_change(t *testing.T) {
	if !tool.Turso_platform_supported() {
		t.Skip("turso unsupported platform")
	}
	dir := t.TempDir()

	// phase 1: sqlite keeper on <dir>/shared.db
	t.Setenv("NAMU_DB", dir+"/shared")
	t.Setenv("NAMU_DB_TYPE", "sqlite")
	tool.DB_init()
	tool.Stop_keeper()
	if err := tool.Start_keeper(); err != nil {
		t.Fatalf("sqlite keeper start: %v", err)
	}

	// phase 2: same path, turso backend
	t.Setenv("NAMU_DB_TYPE", "turso")
	tool.DB_init()
	if err := tool.Start_keeper(); err != nil {
		t.Fatalf("turso keeper start (rebind): %v", err)
	}
	if warm := keeper_attach_ms(t); warm > 10 {
		t.Errorf("attach after type-change rebind = %.2fms, want warm (<10ms) - keeper likely no-op'd on same path", warm)
	}
	tool.Stop_keeper()
}

// Test_keeper_turso_off_baseline: without the keeper and with all pools
// closed, the first attach pays the full instance load. Documents the OFF
// baseline that the ON numbers are compared against (see the probe in tmp/).
func Test_keeper_turso_off_baseline(t *testing.T) {
	if !tool.Turso_platform_supported() {
		t.Skip("turso unsupported platform")
	}
	audit_setup(t, "turso")
	tool.Stop_keeper()

	cold := keeper_attach_ms(t)
	t.Logf("turso keeper OFF cold attach=%.2fms", cold)
	if cold < 10 {
		t.Logf("note: cold attach was fast (%.2fms) - instance may still be warm from audit_setup; informational only", cold)
	}
}


// keeper regression sanity for sqlite: Start x2 + churn must stay error-free
// (full behavior is covered by the existing sqlite_busy keeper tests).
func Test_keeper_sqlite_sanity(t *testing.T) {
	audit_setup(t, "sqlite")

	if err := tool.Start_keeper(); err != nil {
		t.Fatalf("sqlite keeper start: %v", err)
	}
	if err := tool.Start_keeper(); err != nil {
		t.Fatalf("sqlite keeper start x2: %v", err)
	}

	db := tool.DB_connect()
	var one int
	if err := db.QueryRow(`select 1`).Scan(&one); err != nil {
		t.Fatalf("query with keeper on: %v", err)
	}
	tool.DB_close(db)
}

func database_sql_pkg_marker() error { return nil }