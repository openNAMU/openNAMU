package tool

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/blevesearch/bleve/v2"
	_ "github.com/blevesearch/bleve/v2/analysis/lang/cjk"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/search/query"
)

type Search_document struct {
	Title_search string `json:"title_search"`
	Data         string `json:"data"`
	Fingerprint  string `json:"fingerprint"`
}

type search_index_change struct {
	Data    string
	Deleted bool
}

var search_index_directory = filepath.Join("data", "bleve")
var search_index_version_file = filepath.Join("data", "bleve.version")

const search_index_version = "3"

var search_index bleve.Index
var search_index_lock sync.RWMutex
var search_index_ready bool
var search_index_managed bool
var search_index_update_lock sync.Mutex
var search_index_recovery search_index_recovery_state
var search_index_pending = map[string]search_index_change{}
var search_index_runtime_check_lock sync.Mutex
var search_index_runtime_check_time time.Time

func Search_index_field_mapping() *mapping.FieldMapping {
	field_mapping := mapping.NewTextFieldMapping()
	field_mapping.Analyzer = "cjk"
	field_mapping.Store = false
	field_mapping.IncludeTermVectors = false
	field_mapping.IncludeInAll = false
	field_mapping.DocValues = false
	return field_mapping
}

func Search_index_mapping() *mapping.IndexMappingImpl {
	index_mapping := bleve.NewIndexMapping()
	index_mapping.DefaultAnalyzer = "cjk"
	index_mapping.StoreDynamic = false
	index_mapping.IndexDynamic = false
	index_mapping.DocValuesDynamic = false

	document_mapping := mapping.NewDocumentMapping()
	document_mapping.Dynamic = false
	document_mapping.AddFieldMappingsAt("title_search", Search_index_field_mapping())
	document_mapping.AddFieldMappingsAt("data", Search_index_field_mapping())
	document_mapping.AddFieldMappingsAt("fingerprint", Search_index_fingerprint_field_mapping())
	index_mapping.DefaultMapping = document_mapping

	return index_mapping
}

func Search_index_start() {
	search_index_recovery.started.Store(true)
	Search_index_open()
}

func Search_index_open() {
	search_index_lock.RLock()
	ready := search_index_ready
	search_index_lock.RUnlock()
	if ready {
		return
	}
	search_index_recovery.schedule("document index", "startup or explicit open", Search_index_open_attempt)
}

func Search_index_ready() bool {
	Search_index_check_runtime()
	search_index_lock.RLock()
	defer search_index_lock.RUnlock()
	return search_index_ready
}

func Search_index_version_valid() bool {
	valid, err := search_index_version_matches(search_index_version_file, search_index_version)
	return err == nil && valid
}

func Search_index_apply_change(index bleve.Index, doc_name string, change search_index_change) error {
	if change.Deleted {
		return index.Delete(doc_name)
	}
	return index.Index(doc_name, Search_index_document(doc_name, change.Data))
}

func Search_index_set_ready(index bleve.Index) {
	if index == nil {
		log.Printf("[SEARCH] cannot activate a nil document index")
		return
	}
	search_index_update_lock.Lock()
	search_index_lock.Lock()
	for doc_name, change := range search_index_pending {
		if err := Search_index_apply_change(index, doc_name, change); err != nil {
			search_index_ready = false
			search_index_lock.Unlock()
			search_index_update_lock.Unlock()
			log.Printf("[SEARCH] pending document update failed for %q: %v", doc_name, err)
			if search_index_recovery.started.Load() {
				search_index_recovery.schedule("document index", "pending update replay failed", Search_index_open_attempt)
			}
			return
		}
	}
	if search_index != nil && search_index != index {
		close_err, close_panicked := search_index_close(search_index)
		if close_err != nil && !errors.Is(close_err, bleve.ErrorIndexClosed) {
			if !close_panicked {
				search_index_ready = false
				search_index_lock.Unlock()
				search_index_update_lock.Unlock()
				log.Printf("[SEARCH] close previous document index failed: %v", close_err)
				return
			}
			log.Printf("[SEARCH] previous document index was already closed: %v", close_err)
		}
	}
	search_index = index
	search_index_ready = true
	search_index_managed = false
	search_index_pending = map[string]search_index_change{}
	search_index_runtime_check_lock.Lock()
	search_index_runtime_check_time = time.Now()
	search_index_runtime_check_lock.Unlock()
	search_index_lock.Unlock()
	search_index_recovery.rebuilding.Store(false)
	search_index_update_lock.Unlock()
}

func Search_index_open_attempt() error {
	if err := search_index_restore_previous(search_index_directory, search_index_version_file); err != nil {
		return fmt.Errorf("restore interrupted document index installation: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(search_index_directory), 0o755); err != nil {
		return fmt.Errorf("create document index parent: %w", err)
	}
	if search_index_recovery.force_rebuild.Load() {
		return Search_index_rebuild_from_database()
	}

	index, err := search_index_open(search_index_directory)
	if err == nil {
		version_valid, version_err := search_index_version_matches(search_index_version_file, search_index_version)
		if version_err != nil {
			close_err := index.Close()
			return errors.Join(fmt.Errorf("read document index version: %w", version_err), close_err)
		}
		if version_valid {
			if err := search_index_validate(index); err != nil {
				close_err := index.Close()
				if !search_index_missing_or_corrupt(err) {
					return errors.Join(fmt.Errorf("validate document index: %w", err), close_err)
				}
				log.Printf("[SEARCH] document index is corrupt: %v", err)
			} else {
				if err := Search_index_activate(index, true); err != nil {
					close_err := index.Close()
					return errors.Join(fmt.Errorf("activate document index: %w", err), close_err)
				}
				staging_directory := filepath.Join(filepath.Dir(search_index_directory), filepath.Base(search_index_directory)+".tmp")
				if err := search_index_cleanup_staging(staging_directory); err != nil {
					return err
				}
				if err := search_index_cleanup_version_staging(search_index_version_file); err != nil {
					return err
				}
				search_index_cleanup_previous(search_index_directory)
				log.Printf("[SEARCH] document index opened")
				return nil
			}
		} else {
			if err := index.Close(); err != nil && !errors.Is(err, bleve.ErrorIndexClosed) {
				return fmt.Errorf("close document index with stale version: %w", err)
			}
			log.Printf("[SEARCH] document index version is missing or invalid; rebuilding")
		}
	} else if !search_index_missing_or_corrupt(err) {
		return fmt.Errorf("open document index: %w", err)
	} else {
		log.Printf("[SEARCH] document index is missing or has corrupt metadata: %v", err)
	}

	return Search_index_rebuild_from_database()
}

func Search_index_rebuild() {
	Search_index_request_rebuild("explicit rebuild requested", errors.New("full document index rebuild requested"))
}

func Search_index_request_rebuild(reason string, cause error) {
	search_index_update_lock.Lock()
	defer search_index_update_lock.Unlock()
	search_index_recovery.started.Store(true)
	search_index_recovery.force_rebuild.Store(true)
	search_index_recovery.rebuild_generation.Add(1)
	if err := search_index_invalidate_version(search_index_version_file); err != nil {
		cause = errors.Join(cause, err)
	}
	Search_index_mark_unavailable(reason, cause)
}

func Search_index_rebuild_from_database() (return_err error) {
	rebuild_generation := search_index_recovery.rebuild_generation.Load()
	defer func() {
		if recovered := recover(); recovered != nil {
			return_err = fmt.Errorf("document index rebuild panic: %v", recovered)
		}
	}()

	temporary_directory := filepath.Join(filepath.Dir(search_index_directory), filepath.Base(search_index_directory)+".tmp")
	if err := search_index_cleanup_staging(temporary_directory); err != nil {
		return err
	}
	index, err := bleve.New(temporary_directory, Search_index_mapping())
	if err != nil {
		return fmt.Errorf("create document index staging directory: %w", err)
	}
	index_closed := false
	defer func() {
		if !index_closed {
			if err := index.Close(); err != nil {
				log.Printf("[SEARCH] close failed document staging index: %v", err)
			}
		}
		if cleanup_err := os.RemoveAll(temporary_directory); cleanup_err != nil {
			log.Printf("[SEARCH] remove failed document staging directory %q: %v", temporary_directory, cleanup_err)
		}
	}()

	db := DB_connect()
	defer DB_close(db)
	rows, err := db.Query(DB_change("select title, coalesce(data, '') from data"))
	if err != nil {
		return fmt.Errorf("read document source rows: %w", err)
	}
	batch := index.NewBatch()
	batch_count := 0
	for rows.Next() {
		title := ""
		data := ""
		if err := rows.Scan(&title, &data); err != nil {
			return errors.Join(fmt.Errorf("scan document source row: %w", err), rows.Close())
		}
		if title == "" {
			continue
		}
		if err := batch.Index(title, Search_index_document(title, data)); err != nil {
			return errors.Join(fmt.Errorf("prepare document index entry %q: %w", title, err), rows.Close())
		}
		batch_count++
		if batch_count >= 500 {
			if err := index.Batch(batch); err != nil {
				return errors.Join(fmt.Errorf("write document index batch: %w", err), rows.Close())
			}
			batch = index.NewBatch()
			batch_count = 0
		}
	}
	if err := rows.Err(); err != nil {
		return errors.Join(fmt.Errorf("iterate document source rows: %w", err), rows.Close())
	}
	if err := rows.Close(); err != nil {
		return fmt.Errorf("close document source rows: %w", err)
	}
	if batch_count > 0 {
		if err := index.Batch(batch); err != nil {
			return fmt.Errorf("write final document index batch: %w", err)
		}
	}
	if err := index.Close(); err != nil {
		index_closed = true
		return fmt.Errorf("close document staging index: %w", err)
	}
	index_closed = true

	candidate, err := search_index_open_valid(temporary_directory)
	if err != nil {
		return fmt.Errorf("verify document staging index: %w", err)
	}
	defer func() {
		if candidate != nil {
			if err := candidate.Close(); err != nil && !errors.Is(err, bleve.ErrorIndexClosed) {
				log.Printf("[SEARCH] close document staging verification index failed: %v", err)
			}
		}
	}()

	search_index_update_lock.Lock()
	defer search_index_update_lock.Unlock()
	if rebuild_generation != search_index_recovery.rebuild_generation.Load() {
		return errors.New("a newer document index rebuild was requested during the rebuild")
	}
	search_index_lock.RLock()
	pending_changes := make(map[string]search_index_change, len(search_index_pending))
	for doc_name, change := range search_index_pending {
		pending_changes[doc_name] = change
	}
	search_index_lock.RUnlock()
	for doc_name, change := range pending_changes {
		if err := Search_index_apply_change(candidate, doc_name, change); err != nil {
			return fmt.Errorf("replay pending document change for %q: %w", doc_name, err)
		}
	}
	if err := search_index_validate(candidate); err != nil {
		return fmt.Errorf("validate document index after pending replay: %w", err)
	}
	if err := candidate.Close(); err != nil {
		candidate = nil
		return fmt.Errorf("close verified document staging index: %w", err)
	}
	candidate = nil

	search_index_lock.Lock()
	defer search_index_lock.Unlock()
	if search_index != nil {
		close_err, close_panicked := search_index_close(search_index)
		if close_err != nil && !errors.Is(close_err, bleve.ErrorIndexClosed) && !close_panicked {
			search_index_ready = false
			return fmt.Errorf("close active document index before replacement: %w", close_err)
		}
		if close_panicked {
			log.Printf("[SEARCH] active document index close panicked; continuing replacement: %v", close_err)
		}
		search_index = nil
	}

	installed_index, err := search_index_install(search_index_directory, search_index_version_file, temporary_directory, search_index_version)
	if err != nil {
		restored_index, restore_err := Search_index_restore_active()
		if restore_err != nil {
			return errors.Join(fmt.Errorf("install rebuilt document index: %w", err), restore_err)
		}
		if restored_index != nil {
			search_index = restored_index
			search_index_ready = !search_index_recovery.force_rebuild.Load()
			search_index_managed = true
			search_index_pending = map[string]search_index_change{}
		}
		return fmt.Errorf("install rebuilt document index: %w", err)
	}
	search_index = installed_index
	search_index_ready = true
	search_index_managed = true
	search_index_pending = map[string]search_index_change{}
	search_index_runtime_check_lock.Lock()
	search_index_runtime_check_time = time.Now()
	search_index_runtime_check_lock.Unlock()
	search_index_recovery.force_rebuild.Store(false)
	search_index_recovery.rebuilding.Store(false)
	log.Printf("[SEARCH] document index built")
	return nil
}

func Search_index_restore_active() (bleve.Index, error) {
	index, err := search_index_open_valid(search_index_directory)
	if err != nil {
		return nil, err
	}
	version_valid, version_err := search_index_version_matches(search_index_version_file, search_index_version)
	if version_err != nil {
		close_err := index.Close()
		return nil, errors.Join(version_err, close_err)
	}
	if !version_valid {
		close_err := index.Close()
		return nil, errors.Join(errors.New("restored document index has an invalid version marker"), close_err)
	}
	for doc_name, change := range search_index_pending {
		if err := Search_index_apply_change(index, doc_name, change); err != nil {
			close_err := index.Close()
			return nil, errors.Join(fmt.Errorf("replay pending document change for %q on restored index: %w", doc_name, err), close_err)
		}
	}
	if err := search_index_validate(index); err != nil {
		close_err := index.Close()
		return nil, errors.Join(err, close_err)
	}
	search_index_pending = map[string]search_index_change{}
	return index, nil
}

func Search_index_activate(index bleve.Index, managed bool) error {
	if search_index_recovery.force_rebuild.Load() {
		return errors.New("a full document index rebuild is pending")
	}
	search_index_update_lock.Lock()
	defer search_index_update_lock.Unlock()
	search_index_lock.Lock()
	defer search_index_lock.Unlock()

	for doc_name, change := range search_index_pending {
		if err := Search_index_apply_change(index, doc_name, change); err != nil {
			return fmt.Errorf("replay pending document change for %q: %w", doc_name, err)
		}
	}
	if err := search_index_validate(index); err != nil {
		return fmt.Errorf("validate active document index: %w", err)
	}
	if search_index != nil && search_index != index {
		close_err, close_panicked := search_index_close(search_index)
		if close_err != nil && !errors.Is(close_err, bleve.ErrorIndexClosed) && !close_panicked {
			return fmt.Errorf("close previous document index: %w", close_err)
		}
		if close_panicked {
			log.Printf("[SEARCH] previous document index close panicked: %v", close_err)
		}
	}
	search_index = index
	search_index_ready = true
	search_index_managed = managed
	search_index_pending = map[string]search_index_change{}
	search_index_runtime_check_lock.Lock()
	search_index_runtime_check_time = time.Now()
	search_index_runtime_check_lock.Unlock()
	search_index_recovery.rebuilding.Store(false)
	return nil
}

func Search_index_change_locked(doc_name string, change search_index_change) error {
	search_index_lock.Lock()
	defer search_index_lock.Unlock()

	if search_index_recovery.rebuilding.Load() || !search_index_ready || search_index == nil {
		search_index_pending[doc_name] = change
	}
	if !search_index_ready || search_index == nil {
		return nil
	}
	if err := Search_index_apply_change(search_index, doc_name, change); err != nil {
		search_index_pending[doc_name] = change
		search_index_ready = false
		return err
	}
	return nil
}

func Search_index_update(doc_name string, data string) {
	if doc_name == "" {
		return
	}
	search_index_update_lock.Lock()
	err := Search_index_change_locked(doc_name, search_index_change{Data: data})
	search_index_update_lock.Unlock()
	if err != nil {
		Search_index_mark_unavailable("document write failed", err)
	}
}

func Search_index_delete(doc_name string) {
	if doc_name == "" {
		return
	}
	search_index_update_lock.Lock()
	err := Search_index_change_locked(doc_name, search_index_change{Deleted: true})
	search_index_update_lock.Unlock()
	if err != nil {
		Search_index_mark_unavailable("document delete failed", err)
	}
}

func Search_index_sync(db *sql.DB, doc_name string) {
	if db == nil || doc_name == "" {
		return
	}
	search_index_update_lock.Lock()
	var data string
	err := db.QueryRow(DB_change("select coalesce(data, '') from data where title = ?"), doc_name).Scan(&data)
	var change search_index_change
	if errors.Is(err, sql.ErrNoRows) {
		change.Deleted = true
	} else if err != nil {
		search_index_update_lock.Unlock()
		Search_index_request_rebuild("document source read failed", err)
		return
	} else {
		change.Data = data
	}
	apply_err := Search_index_change_locked(doc_name, change)
	search_index_update_lock.Unlock()
	if apply_err != nil {
		Search_index_mark_unavailable("document synchronization failed", apply_err)
	}
}

func Search_index_mark_unavailable(reason string, err error) {
	search_index_lock.Lock()
	search_index_ready = false
	search_index_lock.Unlock()
	log.Printf("[SEARCH] document index unavailable (%s): %v", reason, err)
	if search_index_recovery.started.Load() {
		search_index_recovery.schedule("document index", reason, Search_index_open_attempt)
	}
}

func Search_index_check_runtime() {
	search_index_lock.RLock()
	managed := search_index_managed && search_index_ready
	search_index_lock.RUnlock()
	if !managed {
		return
	}

	search_index_runtime_check_lock.Lock()
	if time.Since(search_index_runtime_check_time) < search_index_runtime_check_interval {
		search_index_runtime_check_lock.Unlock()
		return
	}
	search_index_runtime_check_time = time.Now()
	search_index_runtime_check_lock.Unlock()

	info, err := os.Stat(search_index_directory)
	if err == nil && info.IsDir() {
		return
	}
	if err == nil {
		err = fmt.Errorf("index path %q is not a directory", search_index_directory)
	}
	Search_index_mark_unavailable("runtime filesystem check failed", err)
}

func Search_index_search(keyword string, search_type string, offset int, limit int) ([]string, bool) {
	if keyword == "" || strings.ContainsAny(keyword, "*?") {
		return []string{}, false
	}
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 50
	}

	var search_query query.FieldableQuery
	search_field := "data"
	if search_type == "title" {
		search_query = bleve.NewWildcardQuery("*" + strings.ToLower(keyword) + "*")
		search_field = "title_search"
	} else {
		search_query = bleve.NewMatchQuery(keyword)
	}
	search_query.SetField(search_field)
	request := bleve.NewSearchRequestOptions(search_query, limit, offset, false)
	request.SortBy([]string{"_id"})

	Search_index_check_runtime()
	search_index_lock.RLock()
	if !search_index_ready || search_index == nil {
		search_index_lock.RUnlock()
		return []string{}, false
	}
	result, err := search_index.Search(request)
	search_index_lock.RUnlock()
	if err != nil {
		Search_index_mark_unavailable("document read failed", err)
		return []string{}, false
	}

	data_list := make([]string, 0, len(result.Hits))
	for _, hit := range result.Hits {
		data_list = append(data_list, hit.ID)
	}
	return data_list, true
}
