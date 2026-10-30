package tool

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const auto_maintenance_lock_stale = 30 * time.Minute
const auto_maintenance_lock_refresh = time.Minute
const auto_scheduler_check_interval = time.Minute
const auto_scheduler_retry_delay = 5 * time.Minute
const auto_daily_maintenance_interval = 24 * time.Hour

type auto_maintenance_lock struct {
	path string
	stop chan struct{}
	done chan struct{}
	once sync.Once
}

func (lock *auto_maintenance_lock) Release() {
	lock.once.Do(func() {
		close(lock.stop)
		<-lock.done
		if err := os.Remove(lock.path); err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Printf("[MAINTENANCE] lock removal failed: %v", err)
		}
	})
}

func Acquire_auto_maintenance_lock(lock_path string) (*auto_maintenance_lock, error) {
	now := time.Now()

	for attempt := 0; attempt < 3; attempt++ {
		file, err := os.OpenFile(lock_path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if err == nil {
			content := fmt.Sprintf("%d|%d", os.Getpid(), now.UnixNano())
			if _, err := file.WriteString(content); err != nil {
				_ = file.Close()
				_ = os.Remove(lock_path)
				return nil, fmt.Errorf("auto maintenance lock write failed: %w", err)
			}
			if err := file.Sync(); err != nil {
				_ = file.Close()
				_ = os.Remove(lock_path)
				return nil, fmt.Errorf("auto maintenance lock sync failed: %w", err)
			}
			if err := file.Close(); err != nil {
				_ = os.Remove(lock_path)
				return nil, fmt.Errorf("auto maintenance lock close failed: %w", err)
			}

			lock := &auto_maintenance_lock{
				path: lock_path,
				stop: make(chan struct{}),
				done: make(chan struct{}),
			}
			go lock.Refresh()
			return lock, nil
		}

		if !errors.Is(err, os.ErrExist) {
			return nil, fmt.Errorf("auto maintenance lock create failed: %w", err)
		}

		lock_info, stat_err := os.Stat(lock_path)
		if errors.Is(stat_err, os.ErrNotExist) {
			continue
		}
		if stat_err != nil {
			return nil, fmt.Errorf("auto maintenance lock stat failed: %w", stat_err)
		}
		if now.Sub(lock_info.ModTime()) < auto_maintenance_lock_stale {
			return nil, fmt.Errorf("auto maintenance lock is active")
		}

		current_info, stat_err := os.Stat(lock_path)
		if errors.Is(stat_err, os.ErrNotExist) {
			continue
		}
		if stat_err != nil {
			return nil, fmt.Errorf("auto maintenance lock recheck failed: %w", stat_err)
		}
		if !current_info.ModTime().Equal(lock_info.ModTime()) {
			continue
		}
		if err := os.Remove(lock_path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("auto maintenance stale lock remove failed: %w", err)
		}

		time.Sleep(100 * time.Millisecond)
	}

	return nil, fmt.Errorf("auto maintenance lock is active")
}

func (lock *auto_maintenance_lock) Refresh() {
	ticker := time.NewTicker(auto_maintenance_lock_refresh)
	defer ticker.Stop()
	defer close(lock.done)

	for {
		select {
		case <-lock.stop:
			return
		case now := <-ticker.C:
			if err := os.Chtimes(lock.path, now, now); err != nil {
				log.Printf("[MAINTENANCE] lock refresh failed: %v", err)
			}
		}
	}
}

func Start_auto_backup_and_maintenance() {
	db_set := Get_DB_set()

	state_dir, err := filepath.Abs(".")
	if err != nil {
		log.Printf("[MAINTENANCE] scheduler state path failed: %v", err)
		return
	}

	source_path := ""
	if db_set["db_type"] == "sqlite" || db_set["db_type"] == "turso" {
		source_path, err = filepath.Abs(filepath.Join("", db_set["db_name"]+".db"))
		if err == nil {
			if resolved_source_path, resolve_err := filepath.EvalSymlinks(source_path); resolve_err == nil {
				source_path = resolved_source_path
			}
			state_dir = filepath.Dir(source_path)
		}
	}

	if db_set["db_type"] == "sqlite" {
		if err != nil {
			log.Printf("[BACKUP] source db path failed: %v", err)
			source_path = ""
		} else if _, stat_err := os.Stat(source_path); stat_err != nil {
			log.Printf("[BACKUP] source db not found: %v", stat_err)
			source_path = ""
		}
	}

	lock_path := filepath.Join(state_dir, ".opennamu_auto_backup_lock")
	backup_last_run_path := filepath.Join(state_dir, ".opennamu_auto_backup_last_run")
	maintenance_last_run_path := filepath.Join(state_dir, ".opennamu_auto_maintenance_last_run")
	backup_retry_after_path := filepath.Join(state_dir, ".opennamu_auto_backup_retry_after")
	maintenance_retry_after_path := filepath.Join(state_dir, ".opennamu_auto_maintenance_retry_after")
	backup_enabled := db_set["db_type"] == "sqlite" && source_path != ""

	go func() {
		for {
			lock, lock_err := Acquire_auto_maintenance_lock(lock_path)
			if lock_err != nil {
				log.Printf("[MAINTENANCE] scheduler lock busy: %v", lock_err)
				time.Sleep(auto_scheduler_check_interval)
				continue
			}

			db := DB_connect()
			now := time.Now()
			backup_due := false
			if backup_enabled {
				back_up_value, setting_err := Maintenance_get_setting_value(db, "back_up")
				if setting_err != nil {
					log.Printf("[BACKUP] interval setting read failed: %v", setting_err)
				} else {
					interval_hours, parse_err := Backup_parse_interval_hours(back_up_value)
					if parse_err != nil {
						log.Printf("[BACKUP] invalid back_up setting: %v", parse_err)
					} else if interval_hours > 0 {
						interval := time.Duration(interval_hours * float64(time.Hour))
						backup_due = Auto_task_is_due(backup_last_run_path, backup_retry_after_path, interval, now)
					}
				}
			}

			maintenance_due := Auto_task_is_due(maintenance_last_run_path, maintenance_retry_after_path, auto_daily_maintenance_interval, now)
			if backup_due {
				run_err := Run_auto_scheduled_task(func() error {
					return Backup_run(db, source_path)
				})
				if run_err != nil {
					log.Printf("[BACKUP] failed: %v", run_err)
					if retry_err := Record_auto_task_failure(backup_retry_after_path, time.Now()); retry_err != nil {
						log.Printf("[BACKUP] retry state write failed: %v", retry_err)
					}
				} else if state_err := Record_auto_task_success(backup_last_run_path, backup_retry_after_path, time.Now()); state_err != nil {
					log.Printf("[BACKUP] success state write failed: %v", state_err)
					if retry_err := Record_auto_task_failure(backup_retry_after_path, time.Now()); retry_err != nil {
						log.Printf("[BACKUP] retry state write failed: %v", retry_err)
					}
				}
			}

			if maintenance_due {
				run_err := Run_auto_scheduled_task(func() error {
					return Maintenance_daily(db)
				})
				if run_err != nil {
					log.Printf("[MAINTENANCE] failed: %v", run_err)
					if retry_err := Record_auto_task_failure(maintenance_retry_after_path, time.Now()); retry_err != nil {
						log.Printf("[MAINTENANCE] retry state write failed: %v", retry_err)
					}
				} else if state_err := Record_auto_task_success(maintenance_last_run_path, maintenance_retry_after_path, time.Now()); state_err != nil {
					log.Printf("[MAINTENANCE] success state write failed: %v", state_err)
					if retry_err := Record_auto_task_failure(maintenance_retry_after_path, time.Now()); retry_err != nil {
						log.Printf("[MAINTENANCE] retry state write failed: %v", retry_err)
					}
				}
			}

			lock.Release()
			DB_close(db)
			time.Sleep(auto_scheduler_check_interval)
		}
	}()
}

func Auto_run_is_due(last_run_path string, interval time.Duration, now time.Time) bool {
	if interval <= 0 {
		return false
	}

	last_run, exists, err := Read_auto_run_time(last_run_path)
	if err != nil {
		log.Printf("[MAINTENANCE] last-run state read failed: %v", err)
		return true
	}
	if !exists {
		return true
	}

	return now.Sub(last_run) >= interval
}

func Auto_task_is_due(last_success_path string, retry_after_path string, interval time.Duration, now time.Time) bool {
	if !Auto_run_is_due(last_success_path, interval, now) {
		return false
	}

	retry_after, exists, err := Read_auto_run_time(retry_after_path)
	if err != nil {
		log.Printf("[MAINTENANCE] retry state read failed: %v", err)
		return true
	}
	return !exists || !now.Before(retry_after)
}

func Record_auto_task_failure(retry_after_path string, attempt_time time.Time) error {
	return Write_auto_run_time(retry_after_path, attempt_time.Add(auto_scheduler_retry_delay))
}

func Record_auto_task_success(last_success_path string, retry_after_path string, success_time time.Time) error {
	if err := Write_auto_run_time(last_success_path, success_time); err != nil {
		return err
	}
	if err := os.Remove(retry_after_path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("retry state remove failed: %w", err)
	}

	return nil
}

func Run_auto_scheduled_task(task func() error) (task_err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			if recovered_err, ok := recovered.(error); ok {
				task_err = fmt.Errorf("scheduled task panic: %w", recovered_err)
			} else {
				task_err = fmt.Errorf("scheduled task panic: %v", recovered)
			}
		}
	}()

	return task()
}

func Maintenance_get_setting_value(db *sql.DB, name string) (string, error) {
	rows, err := db.Query(DB_change(`SELECT data, coverage FROM other WHERE name = ?`), name)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	value := ""
	found := false
	for rows.Next() {
		setting_data := ""
		coverage := ""
		if err := rows.Scan(&setting_data, &coverage); err != nil {
			return "", err
		}
		if !found {
			value = setting_data
			found = true
		}
	}
	if err := rows.Err(); err != nil {
		return "", err
	}

	return value, nil
}

func Read_auto_run_time(path string) (time.Time, bool, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return time.Time{}, false, nil
	}
	if err != nil {
		return time.Time{}, false, err
	}

	unix_nano, err := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
	if err != nil {
		return time.Time{}, false, fmt.Errorf("invalid last-run timestamp: %w", err)
	}

	return time.Unix(0, unix_nano), true, nil
}

func Write_auto_run_time(path string, run_time time.Time) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), ".opennamu_last_run.tmp")
	if err != nil {
		return fmt.Errorf("last-run temporary file create failed: %w", err)
	}
	tmp_path := tmp.Name()
	defer os.Remove(tmp_path)

	if _, err := tmp.WriteString(strconv.FormatInt(run_time.UnixNano(), 10)); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("last-run timestamp write failed: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("last-run timestamp sync failed: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("last-run timestamp close failed: %w", err)
	}
	if err := os.Chmod(tmp_path, 0600); err != nil {
		return fmt.Errorf("last-run timestamp permissions failed: %w", err)
	}
	if err := os.Rename(tmp_path, path); err != nil {
		return fmt.Errorf("last-run timestamp publish failed: %w", err)
	}

	return nil
}

func Backup_parse_interval_hours(value string) (float64, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, nil
	}

	hours, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("back_up value is not a number: %s", value)
	}

	if math.IsNaN(hours) || math.IsInf(hours, 0) {
		return 0, fmt.Errorf("back_up value must be finite: %s", value)
	}
	if hours < 0 {
		return 0, fmt.Errorf("back_up value must be zero or positive: %s", value)
	}
	max_interval_hours := float64(time.Duration(1<<63-1)) / float64(time.Hour)
	if hours > max_interval_hours {
		return 0, fmt.Errorf("back_up value exceeds supported duration: %s", value)
	}

	return hours, nil
}

func Backup_parse_count(value string) (int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 3, nil
	}

	count, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("backup_count value is not an integer: %s", value)
	}

	if count < 0 {
		return 0, fmt.Errorf("backup_count value must be zero or positive: %s", value)
	}

	return count, nil
}

func Backup_run(db *sql.DB, source_path string) error {
	source_path, err := filepath.Abs(source_path)
	if err != nil {
		return err
	}

	if _, stat_err := os.Stat(source_path); stat_err != nil {
		return fmt.Errorf("source db not found: %w", stat_err)
	}

	interval_hours, err := Backup_parse_interval_hours(Get_setting_value(db, "back_up", "", ""))
	if err != nil {
		return fmt.Errorf("backup disabled: %w", err)
	}
	if interval_hours == 0 {
		return nil
	}

	backup_count, err := Backup_parse_count(Get_setting_value(db, "backup_count", "", ""))
	if err != nil {
		return fmt.Errorf("backup disabled: %w", err)
	}

	backup_where := strings.TrimSpace(Get_setting_value(db, "backup_where", "", source_path))
	if backup_where == "" {
		backup_where = source_path
	}

	if info, stat_err := os.Stat(backup_where); stat_err == nil && info.IsDir() {
		backup_where = filepath.Join(backup_where, filepath.Base(source_path))
	}

	backup_dir := filepath.Dir(backup_where)
	file_name := filepath.Base(backup_where)
	base_name := file_name
	if strings.EqualFold(filepath.Ext(file_name), ".db") {
		base_name = strings.TrimSuffix(file_name, filepath.Ext(file_name))
	}
	if base_name == "" {
		return fmt.Errorf("invalid backup_where value")
	}

	backup_dir_info, stat_err := os.Stat(backup_dir)
	if stat_err != nil {
		return fmt.Errorf("backup directory not found: %w", stat_err)
	}
	if !backup_dir_info.IsDir() {
		return fmt.Errorf("backup destination parent is not a directory")
	}
	if err := Backup_cleanup_temporary_files(backup_dir, base_name, time.Now()); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102150405")
	candidate, err := Unique_backup_name(backup_dir, base_name, timestamp)
	if err != nil {
		return fmt.Errorf("backup name conflict: %w", err)
	}

	tmp, err := os.CreateTemp(backup_dir, fmt.Sprintf("%s.backup.tmp", base_name))
	if err != nil {
		return fmt.Errorf("temporary backup file create failed: %w", err)
	}
	defer func() {
		_ = os.Remove(tmp.Name())
	}()

	if close_err := tmp.Close(); close_err != nil {
		return fmt.Errorf("temporary backup file close failed: %w", close_err)
	}

	ctx := context.Background()
	_, exec_err := db.ExecContext(ctx, `VACUUM INTO ?`, tmp.Name())
	if exec_err != nil {
		return fmt.Errorf("backup copy failed: %w", exec_err)
	}

	validate_db, err := sql.Open("sqlite", tmp.Name())
	if err != nil {
		return fmt.Errorf("backup validation failed: %w", err)
	}
	defer validate_db.Close()

	status := ""
	if err := validate_db.QueryRow("PRAGMA integrity_check").Scan(&status); err != nil {
		return fmt.Errorf("backup integrity check failed: %w", err)
	}
	if status != "ok" {
		return fmt.Errorf("backup integrity check failed: %s", status)
	}
	if err := os.Chmod(tmp.Name(), 0600); err != nil {
		return fmt.Errorf("backup permissions failed: %w", err)
	}
	backup_file, err := os.OpenFile(tmp.Name(), os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("backup sync open failed: %w", err)
	}
	if err := backup_file.Sync(); err != nil {
		_ = backup_file.Close()
		return fmt.Errorf("backup sync failed: %w", err)
	}
	if err := backup_file.Close(); err != nil {
		return fmt.Errorf("backup sync close failed: %w", err)
	}

	if err := os.Rename(tmp.Name(), candidate); err != nil {
		return fmt.Errorf("backup publish failed: %w", err)
	}

	if err := Backup_retention(backup_dir, base_name, backup_count); err != nil {
		return fmt.Errorf("backup retention failed: %w", err)
	}

	log.Printf("[BACKUP] ok: %s", candidate)
	return nil
}

func Backup_cleanup_temporary_files(backup_dir string, base_name string, now time.Time) error {
	pattern := regexp.MustCompile("^" + regexp.QuoteMeta(base_name) + `\.backup\.tmp[0-9]+$`)
	entries, err := os.ReadDir(backup_dir)
	if err != nil {
		return fmt.Errorf("backup temporary directory read failed: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || entry.Type()&os.ModeSymlink != 0 || !pattern.MatchString(entry.Name()) {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("backup temporary file stat failed: %w", err)
		}
		if now.Sub(info.ModTime()) <= auto_maintenance_lock_stale {
			continue
		}
		if err := os.Remove(filepath.Join(backup_dir, entry.Name())); err != nil {
			return fmt.Errorf("stale backup temporary file remove failed: %w", err)
		}
	}

	return nil
}

func Unique_backup_name(backup_dir string, base_name string, timestamp string) (string, error) {
	candidate := filepath.Join(backup_dir, base_name+"_"+timestamp+".db")
	_, stat_err := os.Stat(candidate)
	if errors.Is(stat_err, os.ErrNotExist) {
		return candidate, nil
	}
	if stat_err != nil {
		return "", stat_err
	}

	for n := 1; ; n++ {
		candidate := filepath.Join(backup_dir, fmt.Sprintf("%s_%s_%d.db", base_name, timestamp, n))
		_, stat_err = os.Stat(candidate)
		if errors.Is(stat_err, os.ErrNotExist) {
			return candidate, nil
		}
		if stat_err != nil {
			return "", stat_err
		}
	}
}

func Backup_retention(backup_dir string, base_name string, count int) error {
	if count <= 0 {
		return nil
	}

	base_regex := regexp.QuoteMeta(base_name)
	pattern := fmt.Sprintf("^%s_[0-9]{14}(_[0-9]+)?\\.db$", base_regex)
	backup_regex := regexp.MustCompile(pattern)

	entries, err := os.ReadDir(backup_dir)
	if err != nil {
		return fmt.Errorf("backup directory read failed: %w", err)
	}

	type backup_file_info struct {
		path  string
		mtime time.Time
	}

	files := []backup_file_info{}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		file_name := entry.Name()
		if !backup_regex.MatchString(file_name) {
			continue
		}

		file_path := filepath.Join(backup_dir, file_name)
		info, err := os.Stat(file_path)
		if err != nil {
			continue
		}

		files = append(files, backup_file_info{path: file_path, mtime: info.ModTime()})
	}

	sort.Slice(files, func(i, j int) bool {
		if files[i].mtime != files[j].mtime {
			return files[i].mtime.Before(files[j].mtime)
		}
		return files[i].path < files[j].path
	})

	for len(files) > count {
		err := os.Remove(files[0].path)
		if err != nil {
			return fmt.Errorf("backup retention remove failed: %w", err)
		}
		files = files[1:]
	}

	return nil
}

func Maintenance_daily(db *sql.DB) error {
	today := time.Now().Format("2006-01-02")
	var maintenance_errors []error

	maintenance_tasks := []struct {
		name string
		run  func() error
	}{
		{name: "vote close", run: func() error { return Maintenance_vote_close(db, today) }},
		{name: "rb expire", run: func() error { return Maintenance_rb_expire(db) }},
		{name: "auth_date expire", run: func() error { return Maintenance_auth_date_expire(db, today) }},
		{name: "document acl expire", run: func() error { return Maintenance_document_acl_expire(db, today) }},
		{name: "ua expire", run: func() error { return Maintenance_ua_expire(db) }},
		{name: "auth history expire", run: func() error { return Maintenance_auth_history_expire(db) }},
		{name: "sitemap", run: func() error { return Maintenance_sitemap_auto(db) }},
		{name: "user title", run: func() error { return Maintenance_user_title(db) }},
	}

	for _, maintenance_task := range maintenance_tasks {
		if err := Run_auto_scheduled_task(maintenance_task.run); err != nil {
			maintenance_errors = append(maintenance_errors, fmt.Errorf("%s: %w", maintenance_task.name, err))
		}
	}

	return errors.Join(maintenance_errors...)
}

func Maintenance_vote_close(db *sql.DB, today string) error {
	return DB_transaction(db, func(tx *sql.Tx) error {
		rows, err := tx.Query(`SELECT id, type FROM vote WHERE (type = 'open' OR type = 'n_open') AND user = ''`)
		if err != nil {
			return err
		}

		type vote_close_info struct {
			id  string
			typ string
		}

		vote_list := []vote_close_info{}

		for rows.Next() {
			var id string
			var typ string
			if err := rows.Scan(&id, &typ); err != nil {
				return err
			}
			vote_list = append(vote_list, vote_close_info{id: id, typ: typ})
		}
		if err := rows.Err(); err != nil {
			_ = rows.Close()
			return err
		}
		if err := rows.Close(); err != nil {
			return err
		}

		for _, vote := range vote_list {
			end_date := ""
			err := tx.QueryRow(`SELECT data FROM vote WHERE id = ? AND name = 'end_date' AND type = 'option'`, vote.id).Scan(&end_date)
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			if err != nil {
				return err
			}
			if end_date == "" {
				continue
			}

			if !Maintenance_date_expired(end_date, today) {
				continue
			}

			new_type := "n_close"
			if vote.typ == "open" {
				new_type = "close"
			}

			_, err = tx.Exec(`UPDATE vote SET type = ? WHERE user = '' AND id = ? AND type = ?`, new_type, vote.id, vote.typ)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func Maintenance_rb_expire(db *sql.DB) error {
	return DB_transaction(db, func(tx *sql.Tx) error {
		rows, err := tx.Query(`SELECT COALESCE(block, ''), COALESCE(end, ''), COALESCE(today, ''), COALESCE(blocker, ''), COALESCE(why, ''), COALESCE(band, ''), COALESCE(login, '') FROM rb WHERE end != '' AND ongoing = '1'`)
		if err != nil {
			return err
		}

		type expired_ban_info struct {
			block   string
			end     string
			today   string
			blocker string
			why     string
			band    string
			login   string
		}
		expired_ban_list := []expired_ban_info{}

		for rows.Next() {
			ban := expired_ban_info{}
			if err := rows.Scan(&ban.block, &ban.end, &ban.today, &ban.blocker, &ban.why, &ban.band, &ban.login); err != nil {
				_ = rows.Close()
				return err
			}
			expired_ban_list = append(expired_ban_list, ban)
		}
		if err := rows.Err(); err != nil {
			_ = rows.Close()
			return err
		}
		if err := rows.Close(); err != nil {
			return err
		}

		for _, ban := range expired_ban_list {
			if Auth_end_active(ban.end) {
				continue
			}
			_, err := tx.Exec(`UPDATE rb SET ongoing = '' WHERE COALESCE(block, '') = ? AND COALESCE(end, '') = ? AND COALESCE(today, '') = ? AND COALESCE(blocker, '') = ? AND COALESCE(why, '') = ? AND COALESCE(band, '') = ? AND COALESCE(login, '') = ? AND ongoing = '1'`, ban.block, ban.end, ban.today, ban.blocker, ban.why, ban.band, ban.login)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func Maintenance_auth_date_expire(db *sql.DB, today string) error {
	return DB_transaction(db, func(tx *sql.Tx) error {
		rows, err := tx.Query(`SELECT id, data FROM user_set WHERE name = 'auth_date'`)
		if err != nil {
			return err
		}

		type auth_date_info struct {
			id   string
			data string
		}
		auth_date_list := []auth_date_info{}

		for rows.Next() {
			id := ""
			data := ""
			if err := rows.Scan(&id, &data); err != nil {
				_ = rows.Close()
				return err
			}
			auth_date_list = append(auth_date_list, auth_date_info{id: id, data: data})
		}
		if err := rows.Err(); err != nil {
			_ = rows.Close()
			return err
		}
		if err := rows.Close(); err != nil {
			return err
		}

		for _, auth_date := range auth_date_list {
			id := auth_date.id
			data := auth_date.data
			if !Maintenance_date_expired(data, today) {
				continue
			}

			acl := ""
			acl_err := tx.QueryRow(`SELECT data FROM user_set WHERE id = ? AND name = 'acl'`, id).Scan(&acl)
			if acl_err != nil && !errors.Is(acl_err, sql.ErrNoRows) {
				return acl_err
			}
			acl_end := ""
			err := tx.QueryRow(`SELECT data FROM user_set WHERE id = ? AND name = 'acl_end'`, id).Scan(&acl_end)
			if err == nil {
				if Auth_group_end_active(acl, acl_end) {
					continue
				}
			} else if !errors.Is(err, sql.ErrNoRows) {
				return err
			}

			_, err = tx.Exec(`UPDATE user_set SET data = 'user' WHERE id = ? AND name = 'acl'`, id)
			if err != nil {
				return err
			}

			_, err = tx.Exec(`DELETE FROM user_set WHERE id = ? AND name = 'auth_date'`, id)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func Maintenance_date_expired(value string, today string) bool {
	date_part := strings.SplitN(value, " ", 2)[0]
	value_date, value_err := time.Parse("2006-01-02", date_part)
	today_date, today_err := time.Parse("2006-01-02", today)
	if today_err != nil {
		return false
	}
	if value_err != nil {
		return true
	}

	return value_date.Before(today_date)
}

func Maintenance_document_acl_expire(db *sql.DB, today string) error {
	return DB_transaction(db, func(tx *sql.Tx) error {
		rows, err := tx.Query(`SELECT doc_name, doc_rev, set_data FROM data_set WHERE set_name = 'acl_date'`)
		if err != nil {
			return err
		}

		type document_acl_expiry_info struct {
			doc_name string
			doc_rev  string
			set_data string
		}
		document_acl_expiry_list := []document_acl_expiry_info{}

		for rows.Next() {
			doc_name := ""
			doc_rev := ""
			set_data := ""
			if err := rows.Scan(&doc_name, &doc_rev, &set_data); err != nil {
				_ = rows.Close()
				return err
			}
			document_acl_expiry_list = append(document_acl_expiry_list, document_acl_expiry_info{
				doc_name: doc_name,
				doc_rev:  doc_rev,
				set_data: set_data,
			})
		}
		if err := rows.Err(); err != nil {
			_ = rows.Close()
			return err
		}
		if err := rows.Close(); err != nil {
			return err
		}

		for _, document_acl_expiry := range document_acl_expiry_list {
			doc_name := document_acl_expiry.doc_name
			doc_rev := document_acl_expiry.doc_rev
			set_data := document_acl_expiry.set_data
			if !Maintenance_date_expired(set_data, today) {
				continue
			}

			_, err = tx.Exec(`DELETE FROM acl WHERE title = ? AND type = ?`, doc_name, doc_rev)
			if err != nil {
				return err
			}

			_, err = tx.Exec(`DELETE FROM data_set WHERE doc_name = ? AND doc_rev = ? AND set_name = 'acl_date'`, doc_name, doc_rev)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func Maintenance_expiration_cutoff(days_raw string, now time.Time) (time.Time, bool, error) {
	days_raw = strings.TrimSpace(days_raw)
	if days_raw == "" {
		return time.Time{}, false, nil
	}

	days, err := strconv.Atoi(days_raw)
	if err != nil {
		return time.Time{}, false, fmt.Errorf("expiration days must be a non-negative integer: %s", days_raw)
	}
	if days < 0 {
		return time.Time{}, false, fmt.Errorf("expiration days must not be negative: %s", days_raw)
	}
	max_days := int64((1<<63 - 1) / int64(24*time.Hour))
	if int64(days) > max_days {
		return time.Time{}, false, fmt.Errorf("expiration days exceed supported duration: %s", days_raw)
	}

	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	cutoff := midnight.Add(-time.Duration(days) * 24 * time.Hour)
	return cutoff, true, nil
}

func Maintenance_ua_expire(db *sql.DB) error {
	days_raw, err := Maintenance_get_setting_value(db, "ua_expiration_date")
	if err != nil {
		return err
	}
	cutoff, enabled, err := Maintenance_expiration_cutoff(days_raw, time.Now())
	if err != nil {
		return err
	}
	if !enabled {
		return nil
	}
	cutoff_time := cutoff.Format("2006-01-02 15:04:05")

	return DB_transaction(db, func(tx *sql.Tx) error {
		_, err := tx.Exec(`DELETE FROM ua_d WHERE today < ?`, cutoff_time)
		return err
	})
}

func Maintenance_auth_history_expire(db *sql.DB) error {
	days_raw, err := Maintenance_get_setting_value(db, "auth_history_expiration_date")
	if err != nil {
		return err
	}
	cutoff, enabled, err := Maintenance_expiration_cutoff(days_raw, time.Now())
	if err != nil {
		return err
	}
	if !enabled {
		return nil
	}
	cutoff_time := cutoff.Format("2006-01-02 15:04:05")

	return DB_transaction(db, func(tx *sql.Tx) error {
		_, err := tx.Exec(`DELETE FROM re_admin WHERE time < ?`, cutoff_time)
		return err
	})
}

func Maintenance_sitemap_auto(db *sql.DB) error {
	value, err := Maintenance_get_setting_value(db, "sitemap_auto_make")
	if err != nil {
		return err
	}
	value = strings.TrimSpace(value)
	if value == "" || value == "0" {
		return nil
	}

	_, _, err = Make_sitemap(db)
	if err != nil {
		return err
	}

	return nil
}

func Maintenance_user_title(db *sql.DB) error {
	rows, err := db.Query(`SELECT id FROM user_set WHERE name = 'user_title' AND data = '✅'`)
	if err != nil {
		return err
	}

	title_ids := []string{}
	for rows.Next() {
		id := ""
		if err := rows.Scan(&id); err != nil {
			_ = rows.Close()
			return err
		}
		title_ids = append(title_ids, id)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return err
	}
	if err := rows.Close(); err != nil {
		return err
	}

	if len(title_ids) == 0 {
		return nil
	}

	admin_title_ids := []string{}
	for _, id := range title_ids {
		if Check_permission(db, "all_admin_auth", id) {
			admin_title_ids = append(admin_title_ids, id)
		}
	}

	return DB_transaction(db, func(tx *sql.Tx) error {
		for _, id := range admin_title_ids {
			_, err := tx.Exec(`UPDATE user_set SET data = '☑️' WHERE name = 'user_title' AND data = '✅' AND id = ?`, id)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
