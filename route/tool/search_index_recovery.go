package tool

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/blevesearch/bleve/v2"
)

const search_index_retry_initial = time.Second
const search_index_retry_max = 30 * time.Second
const search_index_runtime_check_interval = 5 * time.Second

type search_index_recovery_state struct {
	mu                 sync.Mutex
	started            atomic.Bool
	rebuilding         atomic.Bool
	force_rebuild      atomic.Bool
	rebuild_generation atomic.Uint64
	running            bool
	retry_count        int
	next_retry         time.Time
}

func search_index_close(index bleve.Index) (close_err error, panicked bool) {
	if index == nil {
		return nil, false
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			close_err = fmt.Errorf("close index panic: %v", recovered)
			panicked = true
		}
	}()
	return index.Close(), false
}

func search_index_attempt(attempt func() error) (attempt_err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			attempt_err = fmt.Errorf("index recovery panic: %v", recovered)
		}
	}()
	return attempt()
}

func (state *search_index_recovery_state) schedule(index_name string, reason string, attempt func() error) {
	state.mu.Lock()
	if state.running || time.Now().Before(state.next_retry) {
		state.mu.Unlock()
		return
	}
	state.running = true
	state.started.Store(true)
	state.rebuilding.Store(true)
	state.mu.Unlock()

	go func() {
		started_at := time.Now()
		log.Printf("[SEARCH] %s recovery started: %s", index_name, reason)
		err := search_index_attempt(attempt)
		if err == nil && state.force_rebuild.Load() {
			err = errors.New("a requested full index rebuild is still pending")
		}
		duration := time.Since(started_at)

		state.mu.Lock()
		state.running = false
		state.rebuilding.Store(false)
		if err == nil {
			state.retry_count = 0
			state.next_retry = time.Time{}
			state.mu.Unlock()
			log.Printf("[SEARCH] %s recovery completed in %s", index_name, duration.Round(time.Millisecond))
			return
		}

		state.retry_count++
		delay := search_index_retry_initial
		for retry := 1; retry < state.retry_count && delay < search_index_retry_max; retry++ {
			delay *= 2
		}
		if delay > search_index_retry_max {
			delay = search_index_retry_max
		}
		state.next_retry = time.Now().Add(delay)
		state.mu.Unlock()

		log.Printf("[SEARCH] %s recovery failed after %s: %v; retrying in %s", index_name, duration.Round(time.Millisecond), err, delay)
		time.AfterFunc(delay, func() {
			state.schedule(index_name, "scheduled retry", attempt)
		})
	}()
}

func search_index_version_matches(version_file string, expected_version string) (bool, error) {
	info, err := os.Lstat(version_file)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if info.IsDir() {
		return false, nil
	}
	data, err := os.ReadFile(version_file)
	if err != nil {
		return false, err
	}
	return string(data) == expected_version, nil
}

func search_index_cleanup_staging(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("remove abandoned index staging directory %q: %w", path, err)
	}
	return nil
}

func search_index_cleanup_version_staging(version_file string) error {
	staged_path := version_file + ".tmp"
	if err := os.Remove(staged_path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("remove abandoned staged version %q: %w", staged_path, err)
	}
	return nil
}

func search_index_invalidate_version(version_file string) error {
	if err := os.Remove(version_file); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("invalidate index version marker %q: %w", version_file, err)
	}
	return nil
}

func search_index_open(path string) (bleve.Index, error) {
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		if _, lstat_err := os.Lstat(path); lstat_err == nil {
			return nil, err
		} else if !errors.Is(lstat_err, os.ErrNotExist) {
			return nil, errors.Join(err, lstat_err)
		}
		return nil, bleve.ErrorIndexPathDoesNotExist
	}
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, bleve.ErrorIndexMetaCorrupt
	}

	meta_path := filepath.Join(path, "index_meta.json")
	if _, err := os.Stat(meta_path); errors.Is(err, os.ErrNotExist) {
		if _, lstat_err := os.Lstat(meta_path); lstat_err == nil {
			return nil, err
		} else if !errors.Is(lstat_err, os.ErrNotExist) {
			return nil, errors.Join(err, lstat_err)
		}
		return nil, bleve.ErrorIndexMetaMissing
	} else if err != nil {
		return nil, err
	}
	if _, err := os.ReadFile(meta_path); err != nil {
		return nil, err
	}
	return bleve.Open(path)
}

func search_index_validate(index bleve.Index) error {
	if index == nil {
		return errors.New("index handle is nil")
	}
	if _, err := index.DocCount(); err != nil {
		return err
	}
	request := bleve.NewSearchRequestOptions(bleve.NewMatchAllQuery(), 1, 0, false)
	_, err := index.Search(request)
	return err
}

func search_index_open_valid(path string) (bleve.Index, error) {
	index, err := search_index_open(path)
	if err != nil {
		return nil, err
	}
	if err := search_index_validate(index); err != nil {
		close_err := index.Close()
		return nil, errors.Join(err, close_err)
	}
	return index, nil
}

func search_index_missing_or_corrupt(err error) bool {
	return errors.Is(err, bleve.ErrorIndexPathDoesNotExist) ||
		errors.Is(err, bleve.ErrorIndexMetaMissing) ||
		errors.Is(err, bleve.ErrorIndexMetaCorrupt) ||
		errors.Is(err, bleve.ErrorIndexReadInconsistency) ||
		errors.Is(err, bleve.ErrorUnknownIndexType)
}

func search_index_stage_version(version_file string, version string) (string, error) {
	staged_path := version_file + ".tmp"
	if err := os.Remove(staged_path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("remove stale staged version %q: %w", staged_path, err)
	}

	file, err := os.OpenFile(staged_path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return "", fmt.Errorf("create staged version %q: %w", staged_path, err)
	}
	if _, err := file.WriteString(version); err != nil {
		close_err := file.Close()
		cleanup_err := os.Remove(staged_path)
		return "", errors.Join(fmt.Errorf("write staged version %q: %w", staged_path, err), close_err, cleanup_err)
	}
	if err := file.Sync(); err != nil {
		close_err := file.Close()
		cleanup_err := os.Remove(staged_path)
		return "", errors.Join(fmt.Errorf("sync staged version %q: %w", staged_path, err), close_err, cleanup_err)
	}
	if err := file.Close(); err != nil {
		cleanup_err := os.Remove(staged_path)
		return "", errors.Join(fmt.Errorf("close staged version %q: %w", staged_path, err), cleanup_err)
	}
	return staged_path, nil
}

func search_index_restore_previous(index_directory string, version_file string) error {
	if _, err := os.Lstat(index_directory); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	previous_directories, err := filepath.Glob(index_directory + ".previous-*")
	if err != nil {
		return err
	}
	sort.Sort(sort.Reverse(sort.StringSlice(previous_directories)))
	for _, previous_directory := range previous_directories {
		if strings.HasSuffix(previous_directory, ".version") {
			continue
		}
		index, err := search_index_open_valid(previous_directory)
		if err != nil {
			if search_index_missing_or_corrupt(err) {
				continue
			}
			return fmt.Errorf("inspect interrupted index backup %q: %w", previous_directory, err)
		}
		if err := index.Close(); err != nil {
			return fmt.Errorf("close interrupted index backup %q: %w", previous_directory, err)
		}
		if err := os.Rename(previous_directory, index_directory); err != nil {
			return fmt.Errorf("restore interrupted index backup %q: %w", previous_directory, err)
		}

		previous_version := previous_directory + ".version"
		if _, err := os.Lstat(version_file); errors.Is(err, os.ErrNotExist) {
			if err := os.Rename(previous_version, version_file); err != nil && !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("restore interrupted index version %q: %w", previous_version, err)
			}
		} else if err != nil {
			return fmt.Errorf("inspect index version %q: %w", version_file, err)
		}
		log.Printf("[SEARCH] restored interrupted index installation from %q", previous_directory)
		return nil
	}
	return nil
}

func search_index_install(index_directory string, version_file string, temporary_directory string, version string) (bleve.Index, error) {
	staged_version, err := search_index_stage_version(version_file, version)
	if err != nil {
		return nil, err
	}

	if version_info, err := os.Lstat(version_file); err == nil && version_info.IsDir() {
		cleanup_err := os.Remove(staged_version)
		return nil, errors.Join(fmt.Errorf("index version marker %q is a directory", version_file), cleanup_err)
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		cleanup_err := os.Remove(staged_version)
		return nil, errors.Join(fmt.Errorf("inspect index version marker %q: %w", version_file, err), cleanup_err)
	}

	previous_directory := fmt.Sprintf("%s.previous-%d", index_directory, time.Now().UnixNano())
	previous_version := previous_directory + ".version"
	active_moved := false
	version_moved := false
	temporary_installed := false
	staged_version_installed := false

	rollback := func(cause error) error {
		rollback_errors := []error{cause}
		if temporary_installed {
			if err := os.RemoveAll(index_directory); err != nil {
				rollback_errors = append(rollback_errors, fmt.Errorf("remove failed replacement %q: %w", index_directory, err))
			}
		}
		if staged_version_installed {
			if err := os.Remove(version_file); err != nil && !errors.Is(err, os.ErrNotExist) {
				rollback_errors = append(rollback_errors, fmt.Errorf("remove failed replacement version %q: %w", version_file, err))
			}
		}
		if version_moved {
			if err := os.Rename(previous_version, version_file); err != nil {
				rollback_errors = append(rollback_errors, fmt.Errorf("restore previous version %q: %w", previous_version, err))
			}
		}
		if active_moved {
			if err := os.Rename(previous_directory, index_directory); err != nil {
				rollback_errors = append(rollback_errors, fmt.Errorf("restore previous index %q: %w", previous_directory, err))
			}
		}
		if err := os.Remove(staged_version); err != nil && !errors.Is(err, os.ErrNotExist) {
			rollback_errors = append(rollback_errors, fmt.Errorf("remove staged version %q: %w", staged_version, err))
		}
		return errors.Join(rollback_errors...)
	}

	if _, err := os.Lstat(index_directory); err == nil {
		if err := os.Rename(index_directory, previous_directory); err != nil {
			return nil, rollback(fmt.Errorf("move active index to backup: %w", err))
		}
		active_moved = true
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, rollback(fmt.Errorf("inspect active index %q: %w", index_directory, err))
	}

	if _, err := os.Lstat(version_file); err == nil {
		if err := os.Rename(version_file, previous_version); err != nil {
			return nil, rollback(fmt.Errorf("move active index version to backup: %w", err))
		}
		version_moved = true
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, rollback(fmt.Errorf("inspect active index version %q: %w", version_file, err))
	}

	if err := os.Rename(temporary_directory, index_directory); err != nil {
		return nil, rollback(fmt.Errorf("install replacement index: %w", err))
	}
	temporary_installed = true
	if err := os.Rename(staged_version, version_file); err != nil {
		return nil, rollback(fmt.Errorf("install replacement index version: %w", err))
	}
	staged_version_installed = true

	index, err := search_index_open_valid(index_directory)
	if err != nil {
		return nil, rollback(fmt.Errorf("reopen replacement index: %w", err))
	}

	if active_moved {
		if err := os.RemoveAll(previous_directory); err != nil {
			log.Printf("[SEARCH] remove previous index backup %q failed: %v", previous_directory, err)
		}
	}
	if version_moved {
		if err := os.Remove(previous_version); err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Printf("[SEARCH] remove previous index version %q failed: %v", previous_version, err)
		}
	}
	return index, nil
}

func search_index_cleanup_previous(index_directory string) {
	previous_directories, err := filepath.Glob(index_directory + ".previous-*")
	if err != nil {
		log.Printf("[SEARCH] list previous index backups failed: %v", err)
		return
	}
	for _, previous_directory := range previous_directories {
		if strings.HasSuffix(previous_directory, ".version") {
			continue
		}
		if err := os.RemoveAll(previous_directory); err != nil {
			log.Printf("[SEARCH] remove stale index backup %q failed: %v", previous_directory, err)
		}
		if err := os.Remove(previous_directory + ".version"); err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Printf("[SEARCH] remove stale index version backup %q failed: %v", previous_directory+".version", err)
		}
	}
}
