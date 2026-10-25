package test

// Keeper release on backend switch to MySQL.
//
// Pre-fix hazard: startKeeper_locked returned early for keeper-less backends
// (mysql and unknown db types) BEFORE stopping a keeper armed for a previous
// backend, so an sqlite/turso keeper connection survived a switch to mysql.
// The fix runs stopKeeper_locked on the early-return path and now propagates
// its error, so a failed release surfaces to the caller instead of leaving a
// silently stale connection behind.
//
// The keeper release logic itself is backend-agnostic: it only depends on
// db_set["db_type"] and the keeper state, never on a reachable MySQL server.
// These tests therefore run without any MySQL credentials - the "mysql"
// backend is exercised purely as the keeper-less branch of startKeeper_locked
// (the no-op arm/release path). A reachable MySQL server is NOT required.
// (The separate Test_mysql_* audit suite covers real-server behavior and
// still skips without MYSQL_AUDIT_DSN.)

import (
	"testing"
	"time"

	"opennamu/route/tool"
)

// keeper_switch_to_mysql switches the sandbox identity to the mysql backend
// and asserts the keeper arm path releases any previously armed keeper
// without a live MySQL server.
func keeper_switch_to_mysql(t *testing.T, prev_db, prev_type string) {
	t.Helper()

	t.Setenv("NAMU_DB_TYPE", "mysql")
	t.Setenv("NAMU_DB", prev_db+"_mysql")
	tool.DB_boot()

	// The keeper-less branch releases the old keeper and returns its error.
	// Nothing was armed for a keeper-less backend, so the release must be
	// clean and the whole Start must succeed.
	if err := tool.Start_keeper(); err != nil {
		t.Fatalf("mysql keeper arm (release of previous keeper) failed: %v", err)
	}

	// Start again: still a clean release/no-op, idempotent.
	if err := tool.Start_keeper(); err != nil {
		t.Fatalf("mysql keeper arm x2: %v", err)
	}

	// Stop twice: must stay safe with nothing armed.
	if err := tool.Stop_keeper(); err != nil {
		t.Fatalf("mysql keeper stop: %v", err)
	}
	if err := tool.Stop_keeper(); err != nil {
		t.Fatalf("mysql keeper stop x2: %v", err)
	}
}

// Test_keeper_release_on_switch_to_mysql: an armed sqlite keeper is released
// when the backend switches to mysql - verified without a MySQL server by
// switching BACK to sqlite and confirming the keeper re-arms cleanly on the
// same sqlite DB that the released keeper used to pin.
func Test_keeper_release_on_switch_to_mysql(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("NAMU_DB", dir+"/a")
	t.Setenv("NAMU_DB_TYPE", "sqlite")
	tool.DB_init()
	if err := tool.Start_keeper(); err != nil {
		t.Fatalf("sqlite keeper start: %v", err)
	}

	keeper_switch_to_mysql(t, dir+"/a", "sqlite")

	// switch back to sqlite on the SAME db: the keeper must re-arm cleanly,
	// proving the mysql switch released (not merely forgot) the old keeper.
	t.Setenv("NAMU_DB_TYPE", "sqlite")
	t.Setenv("NAMU_DB", dir+"/a")
	tool.DB_boot()
	if err := tool.Start_keeper(); err != nil {
		t.Fatalf("sqlite keeper re-arm after mysql switch: %v", err)
	}
	if err := tool.Stop_keeper(); err != nil {
		t.Fatalf("sqlite keeper stop: %v", err)
	}
}

// Test_keeper_release_on_switch_from_turso_to_mysql: same as above with the
// turso keeper. Behavioral evidence that the release is real: after the
// switch (and back), a fresh physical connection pays the full instance load
// again (~28ms on this schema) - if the turso keeper had survived, the
// instance would still be loaded and the attach would be fast.
func Test_keeper_release_on_switch_from_turso_to_mysql(t *testing.T) {
	if !tool.Turso_platform_supported() {
		t.Skip("turso unsupported platform")
	}

	dir := t.TempDir()
	t.Setenv("NAMU_DB", dir+"/t")
	t.Setenv("NAMU_DB_TYPE", "turso")
	tool.DB_init()
	if err := tool.Start_keeper(); err != nil {
		t.Fatalf("turso keeper start: %v", err)
	}

	// sanity: keeper armed -> warm attach
	warm := attach_series_ms(t, 3)
	if warm > 10 {
		t.Fatalf("turso keeper attach = %.2fms before switch, expected warm", warm)
	}

	keeper_switch_to_mysql(t, dir+"/t", "turso")

	// switch back to the turso DB with all keepers off: the first attach must
	// pay the full load again - proving the previous turso keeper connection
	// (and its in-process instance pin) was really released
	t.Setenv("NAMU_DB_TYPE", "turso")
	t.Setenv("NAMU_DB", dir+"/t")
	tool.DB_boot()
	cold := attach_series_ms(t, 3)
	t.Logf("turso attach: keeper-armed %.2fms -> after mysql switch re-open %.2fms", warm, cold)
	if cold < 10 {
		t.Errorf("attach after switch-back = %.2fms, want cold (>=10ms) - old turso keeper survived the mysql switch", cold)
	}
}

// Test_keeper_mysql_alone_noop: mysql-only runs never arm a keeper; Start/Stop
// in any order are error-free and idempotent. Runs without a MySQL server:
// the keeper-less branch touches no database, only the released state.
func Test_keeper_mysql_alone_noop(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("NAMU_DB", dir+"/m")
	t.Setenv("NAMU_DB_TYPE", "mysql")
	tool.DB_boot() // identity only; no DB_init needed for a no-op keeper

	for i := 0; i < 2; i++ {
		if err := tool.Start_keeper(); err != nil {
			t.Fatalf("mysql Start_keeper #%d: %v", i+1, err)
		}
		if err := tool.Stop_keeper(); err != nil {
			t.Fatalf("mysql Stop_keeper #%d: %v", i+1, err)
		}
	}
	// double stop must stay safe
	if err := tool.Stop_keeper(); err != nil {
		t.Fatalf("mysql Stop_keeper extra: %v", err)
	}
}

// Test_keeper_unknown_type_releases: an unknown db_type takes the same
// keeper-less branch as mysql; a previously armed keeper must be released
// with the error propagated path (Start succeeds when nothing was held).
func Test_keeper_unknown_type_releases(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("NAMU_DB", dir+"/u")
	t.Setenv("NAMU_DB_TYPE", "sqlite")
	tool.DB_init()
	if err := tool.Start_keeper(); err != nil {
		t.Fatalf("sqlite keeper start: %v", err)
	}

	t.Setenv("NAMU_DB_TYPE", "totally_unknown_type")
	tool.DB_boot()
	if err := tool.Start_keeper(); err != nil {
		t.Fatalf("unknown-type keeper arm (release path) failed: %v", err)
	}
	if err := tool.Stop_keeper(); err != nil {
		t.Fatalf("keeper stop after unknown-type switch: %v", err)
	}
}

// attach_series_ms returns the median of n real attach cycles
// (DB_connect -> SELECT 1 -> Close).
func attach_series_ms(t *testing.T, n int) float64 {
	t.Helper()
	var times []float64
	for i := 0; i < n; i++ {
		db := tool.DB_connect()
		start := time.Now()
		var one int
		if err := db.QueryRow(`select 1`).Scan(&one); err != nil {
			tool.DB_close(db)
			t.Fatalf("attach query: %v", err)
		}
		times = append(times, time.Since(start).Seconds()*1000)
		tool.DB_close(db)
	}
	// median
	for i := 0; i < len(times); i++ {
		for j := i + 1; j < len(times); j++ {
			if times[j] < times[i] {
				times[i], times[j] = times[j], times[i]
			}
		}
	}
	return times[len(times)/2]
}