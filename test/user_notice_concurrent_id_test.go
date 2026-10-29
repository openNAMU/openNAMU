package test

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_user_notice_concurrent_id_allocation(t *testing.T) {
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

			ensure_other_ver_row(db, t)

			target := "notice_conc_target"
			const workers = 32

			var wg sync.WaitGroup
			var panic_count int64
			start := time.Now()

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

					tool.Send_alarm(db, "conc_from", target, fmt.Sprintf("concurrent %d", i))
				}(i)
			}

			wg.Wait()
			elapsed := time.Since(start)
			count := notice_count(db, target)
			dup_ids := duplicate_notice_ids(db, target)
			panics := atomic.LoadInt64(&panic_count)

			t.Logf(
				"backend=%s workers=%d count=%d panics=%d duplicate_ids=%v elapsed=%s",
				backend,
				workers,
				count,
				panics,
				dup_ids,
				elapsed,
			)

			if panics != 0 {
				t.Errorf("concurrent Send_alarm panicked %d times", panics)
			}
			if count != workers {
				t.Errorf("user_notice count for %s = %d, want %d", target, count, workers)
			}
			if len(dup_ids) > 0 {
				t.Errorf("duplicate user_notice id values found: %v", dup_ids)
			}

			var first_id string
			if tool.QueryRow_DB(
				db,
				"select id from user_notice where name = ? order by id + 0 asc limit 1",
				[]any{&first_id},
				target,
			) && first_id != "" {
				before_delete := count

				out := route.Api_alarm_delete_post(tool.Config{IP: target}, "", first_id, false)
				response := ""
				if v, ok := out["response"].(string); ok {
					response = v
				}
				after_delete := notice_count(db, target)

				t.Logf(
					"backend=%s delete_one id=%s before=%d after=%d response=%q",
					backend,
					first_id,
					before_delete,
					after_delete,
					response,
				)

				if response != "ok" {
					t.Errorf("Api_alarm_delete_post response = %q, want ok", response)
				}
				if after_delete != before_delete-1 {
					t.Errorf(
						"delete by id removed %d rows, want 1 (before=%d after=%d)",
						before_delete-after_delete,
						before_delete,
						after_delete,
					)
				}

				var remaining int
				_ = tool.QueryRow_DB(
					db,
					"select count(*) from user_notice where name = ? and id = ?",
					[]any{&remaining},
					target,
					first_id,
				)
				if remaining != 0 {
					t.Errorf("deleted id %s still exists with %d rows", first_id, remaining)
				}
			}
		})
	}
}

func ensure_other_ver_row(db *sql.DB, t *testing.T) {
	t.Helper()

	var count int
	exists := tool.QueryRow_DB(
		db,
		"select count(*) from other where name = 'ver'",
		[]any{&count},
	)
	if !exists || count == 0 {
		tool.Exec_DB(
			db,
			"insert into other (name, data, coverage) values ('ver', '20261028', '')",
		)
	}
}

func duplicate_notice_ids(db *sql.DB, name string) []string {
	rows, err := db.Query(
		"select id, count(*) from user_notice where name = ? group by id having count(*) > 1",
		name,
	)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		var count int
		if rows.Scan(&id, &count) == nil {
			ids = append(ids, id)
		}
	}
	return ids
}
