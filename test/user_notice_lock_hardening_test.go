package test

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_user_notice_lock_hardening(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso", "mysql"} {
		if backend == "mysql" && os.Getenv("MYSQL_AUDIT_DSN") == "" {
			t.Run(backend, func(t *testing.T) {
				t.Skip("MySQL audit: set MYSQL_AUDIT_DSN to run against a test-only MariaDB")
			})
			continue
		}

		t.Run(backend, func(t *testing.T) {
			db := new_record_db(t, backend)
			defer tool.DB_close(db)

			t.Run("missing_ver", func(t *testing.T) {
				set_other_ver_rows(db, t, nil)

				target := "missing_ver_target"
				latencies, panics := send_round(db, t, target, 32)
				assert_round(db, t, target, 32, latencies, panics, "missing_ver")

				if rows := other_ver_rows(db); len(rows) != 0 {
					t.Errorf("unexpected other.ver rows after missing_ver: %v", rows)
				}
			})

			t.Run("duplicate_ver", func(t *testing.T) {
				before := []string{"dup_a", "dup_b"}
				set_other_ver_rows(db, t, before)

				target := "duplicate_ver_target"
				latencies, panics := send_round(db, t, target, 32)
				assert_round(db, t, target, 32, latencies, panics, "duplicate_ver")

				after := other_ver_rows(db)
				if !same_multiset(before, after) {
					t.Errorf("other.ver changed from %v to %v", before, after)
				}
			})

			t.Run("concurrent_ver_update", func(t *testing.T) {
				set_other_ver_rows(db, t, []string{"init_before"})

				target := "concurrent_init_target"
				workers := 32
				init_workers := 8

				latencies := make([]time.Duration, workers)
				var panic_count int64
				var wg sync.WaitGroup

				for i := 0; i < init_workers; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						defer func() {
							if r := recover(); r != nil {
								_ = r
								atomic.AddInt64(&panic_count, 1)
							}
						}()

						tool.Exec_DB(
							db,
							"update other set data = 'init_after' where name = 'ver'",
						)
					}()
				}

				for i := 0; i < workers; i++ {
					wg.Add(1)
					go func(i int) {
						defer wg.Done()
						defer func() {
							if r := recover(); r != nil {
								_ = r
								atomic.AddInt64(&panic_count, 1)
							}
						}()

						start := time.Now()
						tool.Send_alarm(db, "h_from", target, fmt.Sprintf("init %d", i))
						latencies[i] = time.Since(start)
					}(i)
				}

				wg.Wait()
				panics := atomic.LoadInt64(&panic_count)
				assert_round(db, t, target, workers, latencies, panics, "concurrent_init")

				allowed := map[string]bool{
					"init_before": true,
					"init_after":  true,
				}
				for _, value := range other_ver_rows(db) {
					if !allowed[value] {
						t.Errorf("unexpected other.ver data %q", value)
					}
				}
			})

			t.Run("workers_32_single", func(t *testing.T) {
				set_other_ver_rows(db, t, []string{"harden_ver"})

				target := "single32_target"
				latencies, panics := send_round(db, t, target, 32)
				assert_round(db, t, target, 32, latencies, panics, "workers_32_single")
			})

			t.Run("workers_64_single", func(t *testing.T) {
				set_other_ver_rows(db, t, []string{"harden_ver"})

				target := "single64_target"
				latencies, panics := send_round(db, t, target, 64)
				assert_round(db, t, target, 64, latencies, panics, "workers_64_single")
			})

			t.Run("workers_32_multi_recipient", func(t *testing.T) {
				set_other_ver_rows(db, t, []string{"harden_ver"})

				targets := []string{
					"multi32_target_0",
					"multi32_target_1",
					"multi32_target_2",
					"multi32_target_3",
				}

				latencies, panics := send_multi(db, t, targets, 8)
				assert_multi(db, t, targets, 8, latencies, panics, "workers_32_multi_recipient")
			})

			t.Run("workers_64_multi_recipient", func(t *testing.T) {
				set_other_ver_rows(db, t, []string{"harden_ver"})

				targets := []string{
					"multi64_target_0",
					"multi64_target_1",
					"multi64_target_2",
					"multi64_target_3",
				}

				latencies, panics := send_multi(db, t, targets, 16)
				assert_multi(db, t, targets, 16, latencies, panics, "workers_64_multi_recipient")
			})

			t.Run("multiple_connections", func(t *testing.T) {
				set_other_ver_rows(db, t, []string{"harden_ver"})

				targets := []string{
					"conn_target_0",
					"conn_target_1",
					"conn_target_2",
					"conn_target_3",
				}

				conns := make([]*sql.DB, 4)
				for i := range conns {
					conns[i] = tool.DB_connect()
					defer tool.DB_close(conns[i])
				}

				latencies, panics := send_multi_connections(conns, t, targets, 8)
				assert_multi(db, t, targets, 8, latencies, panics, "multiple_connections")
			})

			t.Run("repeated_rounds", func(t *testing.T) {
				set_other_ver_rows(db, t, []string{"harden_ver"})

				target := "repeated_target"
				const (
					rounds        = 3
					workersPerRound = 32
				)

				var all_latencies []time.Duration
				for round := 0; round < rounds; round++ {
					latencies, panics := send_round(db, t, target, workersPerRound)

					if panics != 0 {
						t.Errorf("round %d: Send_alarm panicked %d times", round, panics)
						continue
					}

					expected := (round + 1) * workersPerRound
					if count := notice_count(db, target); count != expected {
						t.Errorf(
							"round %d count = %d, want %d",
							round,
							count,
							expected,
						)
					}

					if dup := duplicate_notice_ids(db, target); len(dup) != 0 {
						t.Errorf("round %d duplicate ids: %v", round, dup)
					}

					t.Logf(
						"repeated_rounds round=%d p50=%s p95=%s worst=%s",
						round,
						percentile(latencies, 50),
						percentile(latencies, 95),
						percentile(latencies, 100),
					)

					all_latencies = append(all_latencies, latencies...)
				}

				if count := notice_count(db, target); count != rounds*workersPerRound {
					t.Errorf("final repeated count = %d, want %d", count, rounds*workersPerRound)
				}
				if dup := duplicate_notice_ids(db, target); len(dup) != 0 {
					t.Errorf("final repeated duplicate ids: %v", dup)
				}
				t.Logf(
					"repeated_rounds final p50=%s p95=%s worst=%s",
					percentile(all_latencies, 50),
					percentile(all_latencies, 95),
					percentile(all_latencies, 100),
				)
			})

			t.Run("inside_existing_tx", func(t *testing.T) {
				set_other_ver_rows(db, t, []string{"harden_ver"})

				target := "inside_tx_target"
				workers := 16

				latencies, panics := send_inside_existing_tx(db, t, target, workers)
				assert_round(db, t, target, workers, latencies, panics, "inside_existing_tx")
			})

			t.Run("delete_by_id", func(t *testing.T) {
				set_other_ver_rows(db, t, []string{"harden_ver"})

				target := "delete_target"
				latencies, panics := send_round(db, t, target, 16)
				assert_round(db, t, target, 16, latencies, panics, "delete_by_id")

				var id string
				if !tool.QueryRow_DB(
					db,
					"select id from user_notice where name = ? order by id + 0 asc limit 1",
					[]any{&id},
					target,
				) || id == "" {
					t.Fatalf("no user_notice id found for %s", target)
				}

				before := notice_count(db, target)
				out := route.Api_alarm_delete_post(tool.Config{IP: target}, "", id, false)

				response := ""
				if value, ok := out["response"].(string); ok {
					response = value
				}

				after := notice_count(db, target)

				t.Logf(
					"delete_by_id id=%s before=%d after=%d response=%q",
					id,
					before,
					after,
					response,
				)

				if response != "ok" {
					t.Errorf("Api_alarm_delete_post response = %q, want ok", response)
				}

				if after != before-1 {
					t.Errorf(
						"delete by id removed %d rows, want 1 (before=%d after=%d)",
						before-after,
						before,
						after,
					)
				}

				var remaining int
				_ = tool.QueryRow_DB(
					db,
					"select count(*) from user_notice where name = ? and id = ?",
					[]any{&remaining},
					target,
					id,
				)
				if remaining != 0 {
					t.Errorf("deleted id %s still exists with %d rows", id, remaining)
				}

				if dup := duplicate_notice_ids(db, target); len(dup) != 0 {
					t.Errorf("duplicate ids after delete: %v", dup)
				}
			})
		})
	}
}

func send_round(db *sql.DB, t *testing.T, target string, workers int) ([]time.Duration, int64) {
	latencies := make([]time.Duration, workers)
	var panic_count int64
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					_ = r
					atomic.AddInt64(&panic_count, 1)
				}
			}()

			start := time.Now()
			tool.Send_alarm(db, "h_from", target, fmt.Sprintf("hardening_%d", i))
			latencies[i] = time.Since(start)
		}(i)
	}

	wg.Wait()
	return latencies, atomic.LoadInt64(&panic_count)
}

func send_multi(db *sql.DB, t *testing.T, targets []string, workersPerTarget int) ([]time.Duration, int64) {
	total := len(targets) * workersPerTarget
	latencies := make([]time.Duration, total)
	var panic_count int64
	var wg sync.WaitGroup

	for targetIdx, target := range targets {
		for i := 0; i < workersPerTarget; i++ {
			idx := targetIdx*workersPerTarget + i
			wg.Add(1)
			go func(targetIdx int, target string, idx int) {
				defer wg.Done()
				defer func() {
					if r := recover(); r != nil {
						_ = r
						atomic.AddInt64(&panic_count, 1)
					}
				}()

				start := time.Now()
				tool.Send_alarm(db, "h_from", target, fmt.Sprintf("multi_%d_%d", targetIdx, idx))
				latencies[idx] = time.Since(start)
			}(targetIdx, target, idx)
		}
	}

	wg.Wait()
	return latencies, atomic.LoadInt64(&panic_count)
}

func send_multi_connections(
	conns []*sql.DB,
	t *testing.T,
	targets []string,
	workersPerTarget int,
) ([]time.Duration, int64) {
	total := len(targets) * workersPerTarget
	latencies := make([]time.Duration, total)
	var panic_count int64
	var wg sync.WaitGroup

	for targetIdx, target := range targets {
		for i := 0; i < workersPerTarget; i++ {
			idx := targetIdx*workersPerTarget + i
			conn := conns[idx%len(conns)]
			wg.Add(1)
			go func(targetIdx int, target string, idx int, conn *sql.DB) {
				defer wg.Done()
				defer func() {
					if r := recover(); r != nil {
						_ = r
						atomic.AddInt64(&panic_count, 1)
					}
				}()

				start := time.Now()
				tool.Send_alarm(conn, "h_from", target, fmt.Sprintf("conn_%d_%d", targetIdx, idx))
				latencies[idx] = time.Since(start)
			}(targetIdx, target, idx, conn)
		}
	}

	wg.Wait()
	return latencies, atomic.LoadInt64(&panic_count)
}

func send_inside_existing_tx(db *sql.DB, t *testing.T, target string, workers int) ([]time.Duration, int64) {
	latencies := make([]time.Duration, workers)
	var panic_count int64
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					_ = r
					atomic.AddInt64(&panic_count, 1)
				}
			}()

			tx, err := db.Begin()
			if err != nil {
				_ = err
				atomic.AddInt64(&panic_count, 1)
				return
			}
			defer tx.Rollback()

			start := time.Now()
			tool.Send_alarm(tx, "tx_from", target, fmt.Sprintf("tx_%d", i))
			latencies[i] = time.Since(start)

			if err := tx.Commit(); err != nil {
				_ = err
				atomic.AddInt64(&panic_count, 1)
			}
		}(i)
	}

	wg.Wait()
	return latencies, atomic.LoadInt64(&panic_count)
}

func assert_round(
	db *sql.DB,
	t *testing.T,
	target string,
	want_count int,
	latencies []time.Duration,
	panics int64,
	label string,
) {
	if panics != 0 {
		t.Errorf("%s: Send_alarm panicked %d times", label, panics)
	}

	count := notice_count(db, target)
	if count != want_count {
		t.Errorf("%s: user_notice count for %s = %d, want %d", label, target, count, want_count)
	}

	dup := duplicate_notice_ids(db, target)
	if len(dup) != 0 {
		t.Errorf("%s: duplicate user_notice ids: %v", label, dup)
	}

	p50 := percentile(latencies, 50)
	p95 := percentile(latencies, 95)
	worst := percentile(latencies, 100)

	t.Logf(
		"%s count=%d panics=%d p50=%s p95=%s worst=%s",
		label,
		count,
		panics,
		p50,
		p95,
		worst,
	)
}

func assert_multi(
	db *sql.DB,
	t *testing.T,
	targets []string,
	want_per_target int,
	latencies []time.Duration,
	panics int64,
	label string,
) {
	if panics != 0 {
		t.Errorf("%s: Send_alarm panicked %d times", label, panics)
	}

	for _, target := range targets {
		count := notice_count(db, target)
		if count != want_per_target {
			t.Errorf(
				"%s: user_notice count for %s = %d, want %d",
				label,
				target,
				count,
				want_per_target,
			)
		}

		dup := duplicate_notice_ids(db, target)
		if len(dup) != 0 {
			t.Errorf("%s: duplicate user_notice ids for %s: %v", label, target, dup)
		}
	}

	p50 := percentile(latencies, 50)
	p95 := percentile(latencies, 95)
	worst := percentile(latencies, 100)

	t.Logf(
		"%s total=%d per_target=%d panics=%d p50=%s p95=%s worst=%s",
		label,
		len(latencies),
		want_per_target,
		panics,
		p50,
		p95,
		worst,
	)
}

func set_other_ver_rows(db *sql.DB, t *testing.T, values []string) {
	tool.Exec_DB(db, "delete from other where name = 'ver'")

	for _, value := range values {
		tool.Exec_DB(
			db,
			"insert into other (name, data, coverage) values ('ver', ?, '')",
			value,
		)
	}
}

func other_ver_rows(db *sql.DB) []string {
	rows, err := db.Query("select data from other where name = 'ver'")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var values []string
	for rows.Next() {
		var value string
		if rows.Scan(&value) == nil {
			values = append(values, value)
		}
	}
	return values
}

func same_multiset(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	a_copy := make([]string, len(a))
	b_copy := make([]string, len(b))
	copy(a_copy, a)
	copy(b_copy, b)

	sort.Strings(a_copy)
	sort.Strings(b_copy)

	for i := range a_copy {
		if a_copy[i] != b_copy[i] {
			return false
		}
	}
	return true
}

func percentile(durations []time.Duration, p int) time.Duration {
	n := len(durations)
	if n == 0 {
		return 0
	}

	sorted := make([]time.Duration, n)
	copy(sorted, durations)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	idx := int(float64(p/100) * float64(n-1))
	if idx < 0 {
		idx = 0
	}
	if idx >= n {
		idx = n - 1
	}

	return sorted[idx]
}
