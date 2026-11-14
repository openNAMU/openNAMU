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
	"github.com/blevesearch/bleve/v2/mapping"
	bleve_query "github.com/blevesearch/bleve/v2/search/query"
)

type Bbs_search_document struct {
	Set_id       string   `json:"set_id"`
	Set_code     string   `json:"set_code"`
	Title        string   `json:"title"`
	Prefix       string   `json:"prefix"`
	Tags         []string `json:"tags"`
	Data         string   `json:"data"`
	Date         string   `json:"date"`
	Comment_code string   `json:"comment_code"`
	Comment_data string   `json:"comment_data"`
	Fingerprint  string   `json:"fingerprint"`
}

type bbs_search_index_change struct {
	Document Bbs_search_document
	Deleted  bool
}

var bbs_search_index_directory = filepath.Join("data", "bleve_bbs")
var bbs_search_index_version_file = filepath.Join("data", "bleve_bbs.version")

const bbs_search_index_version = "5"

var bbs_search_index bleve.Index
var bbs_search_index_lock sync.RWMutex
var bbs_search_index_ready bool
var bbs_search_index_managed bool
var bbs_search_index_update_lock sync.Mutex
var bbs_search_index_recovery search_index_recovery_state
var bbs_search_index_pending = map[string]bbs_search_index_change{}
var bbs_search_index_runtime_check_lock sync.Mutex
var bbs_search_index_runtime_check_time time.Time

func Bbs_search_index_keyword_mapping() *mapping.FieldMapping {
	field_mapping := mapping.NewKeywordFieldMapping()
	field_mapping.Store = false
	field_mapping.IncludeTermVectors = false
	field_mapping.IncludeInAll = false
	field_mapping.DocValues = false
	return field_mapping
}

func Bbs_search_index_date_mapping() *mapping.FieldMapping {
	field_mapping := Bbs_search_index_keyword_mapping()
	field_mapping.DocValues = true
	return field_mapping
}

func Bbs_search_index_mapping() *mapping.IndexMappingImpl {
	index_mapping := bleve.NewIndexMapping()
	index_mapping.StoreDynamic = false
	index_mapping.IndexDynamic = false
	index_mapping.DocValuesDynamic = false

	document_mapping := mapping.NewDocumentMapping()
	document_mapping.Dynamic = false
	document_mapping.AddFieldMappingsAt("set_id", Bbs_search_index_keyword_mapping())
	document_mapping.AddFieldMappingsAt("title", Bbs_search_index_keyword_mapping())
	document_mapping.AddFieldMappingsAt("prefix", Bbs_search_index_keyword_mapping())
	document_mapping.AddFieldMappingsAt("tags", Bbs_search_index_keyword_mapping())
	document_mapping.AddFieldMappingsAt("data", Search_index_field_mapping())
	document_mapping.AddFieldMappingsAt("comment_data", Search_index_field_mapping())
	document_mapping.AddFieldMappingsAt("date", Bbs_search_index_date_mapping())
	document_mapping.AddFieldMappingsAt("fingerprint", Search_index_fingerprint_field_mapping())
	index_mapping.DefaultMapping = document_mapping

	return index_mapping
}

func Search_bbs_index_start() {
	bbs_search_index_recovery.started.Store(true)
	Search_bbs_index_open()
}

func Search_bbs_index_open() {
	bbs_search_index_lock.RLock()
	ready := bbs_search_index_ready
	bbs_search_index_lock.RUnlock()
	if ready {
		return
	}
	bbs_search_index_recovery.schedule("BBS index", "startup or explicit open", Search_bbs_index_open_attempt)
}

func Search_bbs_index_mark_rebuild() error {
	return Search_bbs_index_request_rebuild("explicit rebuild requested", errors.New("full BBS index rebuild requested"))
}

func Search_bbs_index_request_rebuild(reason string, cause error) error {
	bbs_search_index_update_lock.Lock()
	defer bbs_search_index_update_lock.Unlock()
	bbs_search_index_recovery.started.Store(true)
	bbs_search_index_recovery.force_rebuild.Store(true)
	bbs_search_index_recovery.rebuild_generation.Add(1)
	marker_err := search_index_invalidate_version(bbs_search_index_version_file)
	if marker_err != nil {
		cause = errors.Join(cause, marker_err)
	}
	Search_bbs_index_mark_unavailable(reason, cause)
	return marker_err
}

func Search_bbs_index_version_valid() bool {
	valid, err := search_index_version_matches(bbs_search_index_version_file, bbs_search_index_version)
	return err == nil && valid
}

func Search_bbs_index_key(set_id string, set_code string) string {
	return set_id + "\x00" + set_code
}

func Search_bbs_index_comment_key(set_id string, set_code string, comment_code string) string {
	return "comment\x00" + set_id + "\x00" + set_code + "\x00" + comment_code
}

func Search_bbs_index_key_data(key string) (string, string, bool) {
	data := strings.SplitN(key, "\x00", 2)
	if len(data) != 2 || data[0] == "" || data[1] == "" {
		return "", "", false
	}
	return data[0], data[1], true
}

func Search_bbs_index_comment_key_data(key string) (string, string, string, bool) {
	data := strings.Split(key, "\x00")
	if len(data) != 4 || data[0] != "comment" || data[1] == "" || data[2] == "" || data[3] == "" {
		return "", "", "", false
	}
	return data[1], data[2], data[3], true
}

func Search_bbs_index_comment_location(set_id string, set_code string) (string, string, string, bool) {
	parts := strings.Split(set_id, "-")
	root_id := ""
	root_code := ""
	comment_parts := []string{}

	if strings.HasPrefix(set_id, "-1-") {
		parts = strings.Split(strings.TrimPrefix(set_id, "-1-"), "-")
		if len(parts) < 1 || parts[0] == "" {
			return "", "", "", false
		}
		root_id = "-1"
		root_code = parts[0]
		comment_parts = parts[1:]
	} else {
		if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
			return "", "", "", false
		}
		root_id = parts[0]
		root_code = parts[1]
		comment_parts = parts[2:]
	}

	comment_code := set_code
	if len(comment_parts) > 0 {
		comment_code = strings.Join(append(comment_parts, set_code), "-")
	}
	if comment_code == "" {
		return "", "", "", false
	}

	return root_id, root_code, comment_code, true
}

func Search_bbs_index_apply_change(index bleve.Index, key string, change bbs_search_index_change) error {
	if change.Deleted {
		return index.Delete(key)
	}
	return index.Index(key, Bbs_search_index_prepare_document(change.Document))
}

func Search_bbs_index_set_ready(index bleve.Index) {
	if index == nil {
		log.Printf("[SEARCH] cannot activate a nil BBS index")
		return
	}
	bbs_search_index_update_lock.Lock()
	bbs_search_index_lock.Lock()
	for key, change := range bbs_search_index_pending {
		if err := Search_bbs_index_apply_change(index, key, change); err != nil {
			bbs_search_index_ready = false
			bbs_search_index_lock.Unlock()
			bbs_search_index_update_lock.Unlock()
			log.Printf("[SEARCH] pending BBS update failed for %q: %v", key, err)
			if bbs_search_index_recovery.started.Load() {
				bbs_search_index_recovery.schedule("BBS index", "pending update replay failed", Search_bbs_index_open_attempt)
			}
			return
		}
	}
	if bbs_search_index != nil && bbs_search_index != index {
		close_err, close_panicked := search_index_close(bbs_search_index)
		if close_err != nil && !errors.Is(close_err, bleve.ErrorIndexClosed) {
			if !close_panicked {
				bbs_search_index_ready = false
				bbs_search_index_lock.Unlock()
				bbs_search_index_update_lock.Unlock()
				log.Printf("[SEARCH] close previous BBS index failed: %v", close_err)
				return
			}
			log.Printf("[SEARCH] previous BBS index was already closed: %v", close_err)
		}
	}
	bbs_search_index = index
	bbs_search_index_ready = true
	bbs_search_index_managed = false
	bbs_search_index_pending = map[string]bbs_search_index_change{}
	bbs_search_index_runtime_check_lock.Lock()
	bbs_search_index_runtime_check_time = time.Now()
	bbs_search_index_runtime_check_lock.Unlock()
	bbs_search_index_lock.Unlock()
	bbs_search_index_recovery.rebuilding.Store(false)
	bbs_search_index_update_lock.Unlock()
}

func Search_bbs_index_open_attempt() error {
	if err := search_index_restore_previous(bbs_search_index_directory, bbs_search_index_version_file); err != nil {
		return fmt.Errorf("restore interrupted BBS index installation: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(bbs_search_index_directory), 0o755); err != nil {
		return fmt.Errorf("create BBS index parent: %w", err)
	}
	if bbs_search_index_recovery.force_rebuild.Load() {
		return Search_bbs_index_rebuild_from_database()
	}

	index, err := search_index_open(bbs_search_index_directory)
	if err == nil {
		version_valid, version_err := search_index_version_matches(bbs_search_index_version_file, bbs_search_index_version)
		if version_err != nil {
			close_err := index.Close()
			return errors.Join(fmt.Errorf("read BBS index version: %w", version_err), close_err)
		}
		if version_valid {
			if err := search_index_validate(index); err != nil {
				close_err := index.Close()
				if !search_index_missing_or_corrupt(err) {
					return errors.Join(fmt.Errorf("validate BBS index: %w", err), close_err)
				}
				log.Printf("[SEARCH] BBS index is corrupt: %v", err)
			} else {
				if err := Search_bbs_index_activate(index, true); err != nil {
					close_err := index.Close()
					return errors.Join(fmt.Errorf("activate BBS index: %w", err), close_err)
				}
				staging_directory := filepath.Join(filepath.Dir(bbs_search_index_directory), filepath.Base(bbs_search_index_directory)+".tmp")
				if err := search_index_cleanup_staging(staging_directory); err != nil {
					return err
				}
				if err := search_index_cleanup_version_staging(bbs_search_index_version_file); err != nil {
					return err
				}
				search_index_cleanup_previous(bbs_search_index_directory)
				log.Printf("[SEARCH] BBS index opened")
				return nil
			}
		} else {
			if err := index.Close(); err != nil && !errors.Is(err, bleve.ErrorIndexClosed) {
				return fmt.Errorf("close BBS index with stale version: %w", err)
			}
			log.Printf("[SEARCH] BBS index version is missing or invalid; rebuilding")
		}
	} else if !search_index_missing_or_corrupt(err) {
		return fmt.Errorf("open BBS index: %w", err)
	} else {
		log.Printf("[SEARCH] BBS index is missing or has corrupt metadata: %v", err)
	}

	return Search_bbs_index_rebuild_from_database()
}

func Search_bbs_index_rebuild() {
	if err := Search_bbs_index_mark_rebuild(); err != nil {
		log.Printf("[SEARCH] request BBS index rebuild failed: %v", err)
	}
}

func Search_bbs_index_rebuild_from_database() (return_err error) {
	rebuild_generation := bbs_search_index_recovery.rebuild_generation.Load()
	defer func() {
		if recovered := recover(); recovered != nil {
			return_err = fmt.Errorf("BBS index rebuild panic: %v", recovered)
		}
	}()

	temporary_directory := filepath.Join(filepath.Dir(bbs_search_index_directory), filepath.Base(bbs_search_index_directory)+".tmp")
	if err := search_index_cleanup_staging(temporary_directory); err != nil {
		return err
	}
	index, err := bleve.New(temporary_directory, Bbs_search_index_mapping())
	if err != nil {
		return fmt.Errorf("create BBS index staging directory: %w", err)
	}
	index_closed := false
	defer func() {
		if !index_closed {
			if err := index.Close(); err != nil {
				log.Printf("[SEARCH] close failed BBS staging index: %v", err)
			}
		}
		if cleanup_err := os.RemoveAll(temporary_directory); cleanup_err != nil {
			log.Printf("[SEARCH] remove failed BBS staging directory %q: %v", temporary_directory, cleanup_err)
		}
	}()

	db := DB_connect()
	defer DB_close(db)
	rows, err := db.Query(DB_change("select set_id, set_code, set_name, set_data from bbs_data where set_name in ('title', 'prefix', 'tag', 'data', 'date') order by set_id, set_code"))
	if err != nil {
		return fmt.Errorf("read BBS source rows: %w", err)
	}
	current_key := ""
	document := Bbs_search_document{}
	batch := index.NewBatch()
	batch_count := 0
	flush := func() error {
		if current_key == "" || document.Title == "" {
			return nil
		}
		if err := batch.Index(current_key, Bbs_search_index_prepare_document(document)); err != nil {
			return err
		}
		batch_count++
		if batch_count < 500 {
			return nil
		}
		if err := index.Batch(batch); err != nil {
			return err
		}
		batch = index.NewBatch()
		batch_count = 0
		return nil
	}
	for rows.Next() {
		set_id := ""
		set_code := ""
		set_name := ""
		set_data := ""
		if err := rows.Scan(&set_id, &set_code, &set_name, &set_data); err != nil {
			return errors.Join(fmt.Errorf("scan BBS source row: %w", err), rows.Close())
		}
		key := Search_bbs_index_key(set_id, set_code)
		if current_key != "" && current_key != key {
			if err := flush(); err != nil {
				return errors.Join(fmt.Errorf("prepare BBS index entry %q: %w", current_key, err), rows.Close())
			}
			document = Bbs_search_document{}
		}
		current_key = key
		document.Set_id = strings.ToLower(set_id)
		switch set_name {
		case "title":
			document.Title = strings.ToLower(set_data)
		case "prefix":
			document.Prefix = strings.ToLower(set_data)
		case "tag":
			document.Tags = append(document.Tags, strings.ToLower(set_data))
		case "data":
			document.Data = set_data
		case "date":
			document.Date = set_data
		}
	}
	if err := rows.Err(); err != nil {
		return errors.Join(fmt.Errorf("iterate BBS source rows: %w", err), rows.Close())
	}
	if err := rows.Close(); err != nil {
		return fmt.Errorf("close BBS source rows: %w", err)
	}
	if err := flush(); err != nil {
		return fmt.Errorf("prepare final BBS index entry: %w", err)
	}
	if batch_count > 0 {
		if err := index.Batch(batch); err != nil {
			return fmt.Errorf("write final BBS index batch: %w", err)
		}
	}
	if err := Search_bbs_index_rebuild_comments(db, index); err != nil {
		return fmt.Errorf("rebuild nested BBS comment documents: %w", err)
	}
	if err := index.Close(); err != nil {
		index_closed = true
		return fmt.Errorf("close BBS staging index: %w", err)
	}
	index_closed = true

	candidate, err := search_index_open_valid(temporary_directory)
	if err != nil {
		return fmt.Errorf("verify BBS staging index: %w", err)
	}
	defer func() {
		if candidate != nil {
			if err := candidate.Close(); err != nil && !errors.Is(err, bleve.ErrorIndexClosed) {
				log.Printf("[SEARCH] close BBS staging verification index failed: %v", err)
			}
		}
	}()

	bbs_search_index_update_lock.Lock()
	defer bbs_search_index_update_lock.Unlock()
	if rebuild_generation != bbs_search_index_recovery.rebuild_generation.Load() {
		return errors.New("a newer BBS index rebuild was requested during the rebuild")
	}
	bbs_search_index_lock.RLock()
	pending_changes := make(map[string]bbs_search_index_change, len(bbs_search_index_pending))
	for key, change := range bbs_search_index_pending {
		pending_changes[key] = change
	}
	bbs_search_index_lock.RUnlock()
	for key, change := range pending_changes {
		if err := Search_bbs_index_apply_change(candidate, key, change); err != nil {
			return fmt.Errorf("replay pending BBS change for %q: %w", key, err)
		}
	}
	if err := search_index_validate(candidate); err != nil {
		return fmt.Errorf("validate BBS index after pending replay: %w", err)
	}
	if err := candidate.Close(); err != nil {
		candidate = nil
		return fmt.Errorf("close verified BBS staging index: %w", err)
	}
	candidate = nil

	bbs_search_index_lock.Lock()
	defer bbs_search_index_lock.Unlock()
	if bbs_search_index != nil {
		close_err, close_panicked := search_index_close(bbs_search_index)
		if close_err != nil && !errors.Is(close_err, bleve.ErrorIndexClosed) && !close_panicked {
			bbs_search_index_ready = false
			return fmt.Errorf("close active BBS index before replacement: %w", close_err)
		}
		if close_panicked {
			log.Printf("[SEARCH] active BBS index close panicked; continuing replacement: %v", close_err)
		}
		bbs_search_index = nil
	}
	installed_index, err := search_index_install(bbs_search_index_directory, bbs_search_index_version_file, temporary_directory, bbs_search_index_version)
	if err != nil {
		restored_index, restore_err := Search_bbs_index_restore_active()
		if restore_err != nil {
			return errors.Join(fmt.Errorf("install rebuilt BBS index: %w", err), restore_err)
		}
		if restored_index != nil {
			bbs_search_index = restored_index
			bbs_search_index_ready = !bbs_search_index_recovery.force_rebuild.Load()
			bbs_search_index_managed = true
			bbs_search_index_pending = map[string]bbs_search_index_change{}
		}
		return fmt.Errorf("install rebuilt BBS index: %w", err)
	}
	bbs_search_index = installed_index
	bbs_search_index_ready = true
	bbs_search_index_managed = true
	bbs_search_index_pending = map[string]bbs_search_index_change{}
	bbs_search_index_runtime_check_lock.Lock()
	bbs_search_index_runtime_check_time = time.Now()
	bbs_search_index_runtime_check_lock.Unlock()
	bbs_search_index_recovery.force_rebuild.Store(false)
	bbs_search_index_recovery.rebuilding.Store(false)
	log.Printf("[SEARCH] BBS index built")
	return nil
}

func Search_bbs_index_restore_active() (bleve.Index, error) {
	index, err := search_index_open_valid(bbs_search_index_directory)
	if err != nil {
		return nil, err
	}
	version_valid, version_err := search_index_version_matches(bbs_search_index_version_file, bbs_search_index_version)
	if version_err != nil {
		close_err := index.Close()
		return nil, errors.Join(version_err, close_err)
	}
	if !version_valid {
		close_err := index.Close()
		return nil, errors.Join(errors.New("restored BBS index has an invalid version marker"), close_err)
	}
	for key, change := range bbs_search_index_pending {
		if err := Search_bbs_index_apply_change(index, key, change); err != nil {
			close_err := index.Close()
			return nil, errors.Join(fmt.Errorf("replay pending BBS change for %q on restored index: %w", key, err), close_err)
		}
	}
	if err := search_index_validate(index); err != nil {
		close_err := index.Close()
		return nil, errors.Join(err, close_err)
	}
	bbs_search_index_pending = map[string]bbs_search_index_change{}
	return index, nil
}

func Search_bbs_index_activate(index bleve.Index, managed bool) error {
	if bbs_search_index_recovery.force_rebuild.Load() {
		return errors.New("a full BBS index rebuild is pending")
	}
	bbs_search_index_update_lock.Lock()
	defer bbs_search_index_update_lock.Unlock()
	bbs_search_index_lock.Lock()
	defer bbs_search_index_lock.Unlock()

	for key, change := range bbs_search_index_pending {
		if err := Search_bbs_index_apply_change(index, key, change); err != nil {
			return fmt.Errorf("replay pending BBS change for %q: %w", key, err)
		}
	}
	if err := search_index_validate(index); err != nil {
		return fmt.Errorf("validate active BBS index: %w", err)
	}
	if bbs_search_index != nil && bbs_search_index != index {
		close_err, close_panicked := search_index_close(bbs_search_index)
		if close_err != nil && !errors.Is(close_err, bleve.ErrorIndexClosed) && !close_panicked {
			return fmt.Errorf("close previous BBS index: %w", close_err)
		}
		if close_panicked {
			log.Printf("[SEARCH] previous BBS index close panicked: %v", close_err)
		}
	}
	bbs_search_index = index
	bbs_search_index_ready = true
	bbs_search_index_managed = managed
	bbs_search_index_pending = map[string]bbs_search_index_change{}
	bbs_search_index_runtime_check_lock.Lock()
	bbs_search_index_runtime_check_time = time.Now()
	bbs_search_index_runtime_check_lock.Unlock()
	bbs_search_index_recovery.rebuilding.Store(false)
	return nil
}

func Search_bbs_index_change_locked(key string, change bbs_search_index_change) error {
	bbs_search_index_lock.Lock()
	defer bbs_search_index_lock.Unlock()

	if bbs_search_index_recovery.rebuilding.Load() || !bbs_search_index_ready || bbs_search_index == nil {
		bbs_search_index_pending[key] = change
	}
	if !bbs_search_index_ready || bbs_search_index == nil {
		return nil
	}
	if err := Search_bbs_index_apply_change(bbs_search_index, key, change); err != nil {
		bbs_search_index_pending[key] = change
		bbs_search_index_ready = false
		return err
	}
	return nil
}

func Search_bbs_index_mark_unavailable(reason string, err error) {
	bbs_search_index_lock.Lock()
	bbs_search_index_ready = false
	bbs_search_index_lock.Unlock()
	log.Printf("[SEARCH] BBS index unavailable (%s): %v", reason, err)
	if bbs_search_index_recovery.started.Load() {
		bbs_search_index_recovery.schedule("BBS index", reason, Search_bbs_index_open_attempt)
	}
}

func Search_bbs_index_ready() bool {
	Search_bbs_index_check_runtime()
	bbs_search_index_lock.RLock()
	defer bbs_search_index_lock.RUnlock()
	return bbs_search_index_ready
}

func Search_bbs_index_check_runtime() {
	bbs_search_index_lock.RLock()
	managed := bbs_search_index_managed && bbs_search_index_ready
	bbs_search_index_lock.RUnlock()
	if !managed {
		return
	}

	bbs_search_index_runtime_check_lock.Lock()
	if time.Since(bbs_search_index_runtime_check_time) < search_index_runtime_check_interval {
		bbs_search_index_runtime_check_lock.Unlock()
		return
	}
	bbs_search_index_runtime_check_time = time.Now()
	bbs_search_index_runtime_check_lock.Unlock()

	info, err := os.Stat(bbs_search_index_directory)
	if err == nil && info.IsDir() {
		return
	}
	if err == nil {
		err = fmt.Errorf("index path %q is not a directory", bbs_search_index_directory)
	}
	Search_bbs_index_mark_unavailable("runtime filesystem check failed", err)
}

func Search_bbs_index_rebuild_comments(db *sql.DB, index bleve.Index) error {
	rows, err := db.Query(DB_change("select set_id, set_code, set_name, set_data from bbs_data where set_name in ('comment', 'comment_date') order by set_id, set_code, set_name"))
	if err != nil {
		return fmt.Errorf("query BBS comment source rows: %w", err)
	}

	current_key := ""
	document := Bbs_search_document{}
	batch := index.NewBatch()
	batch_count := 0

	flush := func() error {
		if current_key == "" || document.Comment_data == "" {
			return nil
		}
		key := Search_bbs_index_comment_key(document.Set_id, document.Set_code, document.Comment_code)
		if err := batch.Index(key, Bbs_search_index_prepare_document(document)); err != nil {
			return err
		}
		batch_count++
		if batch_count < 500 {
			return nil
		}
		if err := index.Batch(batch); err != nil {
			return err
		}
		batch = index.NewBatch()
		batch_count = 0
		return nil
	}

	for rows.Next() {
		comment_set_id := ""
		comment_set_code := ""
		set_name := ""
		set_data := ""
		if err := rows.Scan(&comment_set_id, &comment_set_code, &set_name, &set_data); err != nil {
			return errors.Join(fmt.Errorf("scan BBS comment source row: %w", err), rows.Close())
		}

		post_set_id, post_set_code, comment_code, valid := Search_bbs_index_comment_location(comment_set_id, comment_set_code)
		if !valid {
			continue
		}

		key := comment_set_id + "\x00" + comment_set_code
		if current_key != "" && current_key != key {
			if err := flush(); err != nil {
				return err
			}
			document = Bbs_search_document{}
		}
		current_key = key

		if document.Set_id == "" {
			document.Set_id = post_set_id
			document.Set_code = post_set_code
			document.Comment_code = comment_code
		}

		switch set_name {
		case "comment":
			document.Comment_data = set_data
		case "comment_date":
			document.Date = set_data
		}
	}

	if err := rows.Err(); err != nil {
		return errors.Join(fmt.Errorf("iterate BBS comment source rows: %w", err), rows.Close())
	}
	if err := rows.Close(); err != nil {
		return fmt.Errorf("close BBS comment source rows: %w", err)
	}
	if err := flush(); err != nil {
		return err
	}
	if batch_count > 0 {
		if err := index.Batch(batch); err != nil {
			return err
		}
	}
	return nil
}

func search_bbs_index_comment_document(db *sql.DB, set_id string, set_code string, comment_code string) (Bbs_search_document, bool, error) {
	parts := strings.Split(comment_code, "-")
	if len(parts) == 0 || parts[len(parts)-1] == "" {
		return Bbs_search_document{}, false, nil
	}

	comment_set_id := set_id + "-" + set_code
	if len(parts) > 1 {
		comment_set_id += "-" + strings.Join(parts[:len(parts)-1], "-")
	}
	comment_set_code := parts[len(parts)-1]

	document := Bbs_search_document{
		Set_id:       set_id,
		Set_code:     set_code,
		Comment_code: comment_code,
	}
	rows, err := db.Query(DB_change("select set_name, set_data from bbs_data where set_id = ? and set_code = ? and set_name in ('comment', 'comment_date')"), comment_set_id, comment_set_code)
	if err != nil {
		return Bbs_search_document{}, false, err
	}
	for rows.Next() {
		set_name := ""
		set_data := ""
		if err := rows.Scan(&set_name, &set_data); err != nil {
			return Bbs_search_document{}, false, errors.Join(err, rows.Close())
		}
		if set_name == "comment" {
			document.Comment_data = set_data
		} else if set_name == "comment_date" {
			document.Date = set_data
		}
	}
	if err := rows.Err(); err != nil {
		return Bbs_search_document{}, false, errors.Join(err, rows.Close())
	}
	if err := rows.Close(); err != nil {
		return Bbs_search_document{}, false, err
	}
	return document, document.Comment_data != "", nil
}

func Search_bbs_index_comment_document(db *sql.DB, set_id string, set_code string, comment_code string) (Bbs_search_document, bool) {
	document, exists, err := search_bbs_index_comment_document(db, set_id, set_code, comment_code)
	if err != nil {
		panic(err)
	}
	return document, exists
}

func Search_bbs_index_update_comment(db *sql.DB, set_id string, set_code string, comment_code string) {
	if db == nil || set_id == "" || set_code == "" || comment_code == "" {
		return
	}

	bbs_search_index_update_lock.Lock()
	document, exists, err := search_bbs_index_comment_document(db, set_id, set_code, comment_code)
	if err != nil {
		bbs_search_index_update_lock.Unlock()
		Search_bbs_index_request_rebuild("BBS comment source read failed", err)
		return
	}
	key := Search_bbs_index_comment_key(set_id, set_code, comment_code)
	change := bbs_search_index_change{Document: document, Deleted: !exists}
	apply_err := Search_bbs_index_change_locked(key, change)
	bbs_search_index_update_lock.Unlock()
	if apply_err != nil {
		Search_bbs_index_mark_unavailable("BBS comment write failed", apply_err)
	}
}

func Search_bbs_index_delete_comment(set_id string, set_code string, comment_code string) {
	if set_id == "" || set_code == "" || comment_code == "" {
		return
	}

	key := Search_bbs_index_comment_key(set_id, set_code, comment_code)
	bbs_search_index_update_lock.Lock()
	err := Search_bbs_index_change_locked(key, bbs_search_index_change{Deleted: true})
	bbs_search_index_update_lock.Unlock()
	if err != nil {
		Search_bbs_index_mark_unavailable("BBS comment delete failed", err)
	}
}

func Search_bbs_index_delete_comments(db *sql.DB, set_id string, set_code string) {
	comment_set_id := set_id + "-" + set_code
	rows := Query_DB(
		db,
		"select distinct set_id, set_code from bbs_data where set_name = 'comment' and (set_id = ? or set_id like ?)",
		comment_set_id,
		comment_set_id+"-%",
	)
	defer rows.Close()

	for rows.Next() {
		storage_set_id := ""
		storage_set_code := ""
		if rows.Scan(&storage_set_id, &storage_set_code) != nil {
			continue
		}
		comment_set_id, comment_set_code, comment_code, valid := Search_bbs_index_comment_location(storage_set_id, storage_set_code)
		if !valid {
			continue
		}
		Search_bbs_index_delete_comment(comment_set_id, comment_set_code, comment_code)
	}
}

func search_bbs_index_document(db *sql.DB, set_id string, set_code string) (Bbs_search_document, bool, error) {
	document := Bbs_search_document{Set_id: strings.ToLower(set_id)}
	rows, err := db.Query(DB_change("select set_name, set_data from bbs_data where set_id = ? and set_code = ? and set_name in ('title', 'prefix', 'tag', 'data', 'date')"), set_id, set_code)
	if err != nil {
		return Bbs_search_document{}, false, err
	}
	for rows.Next() {
		set_name := ""
		set_data := ""
		if err := rows.Scan(&set_name, &set_data); err != nil {
			return Bbs_search_document{}, false, errors.Join(err, rows.Close())
		}
		switch set_name {
		case "title":
			document.Title = strings.ToLower(set_data)
		case "prefix":
			document.Prefix = strings.ToLower(set_data)
		case "tag":
			document.Tags = append(document.Tags, strings.ToLower(set_data))
		case "data":
			document.Data = set_data
		case "date":
			document.Date = set_data
		}
	}
	if err := rows.Err(); err != nil {
		return Bbs_search_document{}, false, errors.Join(err, rows.Close())
	}
	if err := rows.Close(); err != nil {
		return Bbs_search_document{}, false, err
	}
	return document, document.Title != "", nil
}

func Search_bbs_index_document(db *sql.DB, set_id string, set_code string) (Bbs_search_document, bool) {
	document, exists, err := search_bbs_index_document(db, set_id, set_code)
	if err != nil {
		panic(err)
	}
	return document, exists
}

func Search_bbs_index_update(db *sql.DB, set_id string, set_code string) {
	if db == nil || set_id == "" || set_code == "" {
		return
	}
	bbs_search_index_update_lock.Lock()
	document, exists, err := search_bbs_index_document(db, set_id, set_code)
	if err != nil {
		bbs_search_index_update_lock.Unlock()
		Search_bbs_index_request_rebuild("BBS source read failed", err)
		return
	}
	key := Search_bbs_index_key(set_id, set_code)
	change := bbs_search_index_change{Document: document, Deleted: !exists}
	apply_err := Search_bbs_index_change_locked(key, change)
	bbs_search_index_update_lock.Unlock()
	if apply_err != nil {
		Search_bbs_index_mark_unavailable("BBS write failed", apply_err)
	}
}

func Search_bbs_index_delete(set_id string, set_code string) {
	if set_id == "" || set_code == "" {
		return
	}
	key := Search_bbs_index_key(set_id, set_code)
	bbs_search_index_update_lock.Lock()
	err := Search_bbs_index_change_locked(key, bbs_search_index_change{Deleted: true})
	bbs_search_index_update_lock.Unlock()
	if err != nil {
		Search_bbs_index_mark_unavailable("BBS delete failed", err)
	}
}

func Search_bbs_index_delete_set(db *sql.DB, set_id string) {
	if db == nil {
		return
	}
	rows, err := db.Query(DB_change("select set_code from bbs_data where set_name = 'title' and set_id = ?"), set_id)
	if err != nil {
		Search_bbs_index_request_rebuild("BBS delete-set source read failed", err)
		return
	}
	set_codes := []string{}
	for rows.Next() {
		set_code := ""
		if err := rows.Scan(&set_code); err != nil {
			close_err := rows.Close()
			Search_bbs_index_request_rebuild("BBS delete-set source scan failed", errors.Join(err, close_err))
			return
		}
		set_codes = append(set_codes, set_code)
	}
	if err := rows.Err(); err != nil {
		close_err := rows.Close()
		Search_bbs_index_request_rebuild("BBS delete-set source iteration failed", errors.Join(err, close_err))
		return
	}
	if err := rows.Close(); err != nil {
		Search_bbs_index_request_rebuild("BBS delete-set source close failed", err)
		return
	}
	for _, set_code := range set_codes {
		Search_bbs_index_delete(set_id, set_code)
	}
}

func Search_bbs_index_search_internal(keyword string, set_id string, search_type string, offset int, limit int) ([]string, bool) {
	if keyword == "" || strings.ContainsAny(keyword, "*?%_") {
		return []string{}, false
	}
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 50
	}

	var search_query bleve_query.Query
	if search_type == "data" {
		query := bleve.NewMatchQuery(keyword)
		query.SetField("data")
		search_query = query
	} else {
		queries := []bleve_query.Query{}
		for _, field := range []string{"title", "prefix", "tags"} {
			query := bleve.NewWildcardQuery("*" + strings.ToLower(keyword) + "*")
			query.SetField(field)
			queries = append(queries, query)
		}
		search_query = bleve.NewDisjunctionQuery(queries...)
	}
	if set_id != "" {
		set_query := bleve.NewTermQuery(strings.ToLower(set_id))
		set_query.SetField("set_id")
		search_query = bleve.NewConjunctionQuery(search_query, set_query)
	}

	request := bleve.NewSearchRequestOptions(search_query, limit, offset, false)
	request.SortBy([]string{"-date", "_id"})

	Search_bbs_index_check_runtime()
	bbs_search_index_lock.RLock()
	if !bbs_search_index_ready || bbs_search_index == nil {
		bbs_search_index_lock.RUnlock()
		return []string{}, false
	}
	result, err := bbs_search_index.Search(request)
	bbs_search_index_lock.RUnlock()
	if err != nil {
		Search_bbs_index_mark_unavailable("BBS read failed", err)
		return []string{}, false
	}

	data_list := make([]string, 0, len(result.Hits))
	for _, hit := range result.Hits {
		data_list = append(data_list, hit.ID)
	}
	return data_list, true
}

func Search_bbs_index_search(keyword string, set_id string, offset int, limit int) ([]string, bool) {
	return Search_bbs_index_search_internal(keyword, set_id, "title", offset, limit)
}

func Search_bbs_index_search_data(keyword string, set_id string, offset int, limit int) ([]string, bool) {
	return Search_bbs_index_search_internal(keyword, set_id, "data", offset, limit)
}

func Search_bbs_index_search_comment(keyword string, set_id string, offset int, limit int) ([]string, bool) {
	if keyword == "" || strings.ContainsAny(keyword, "*?%_") {
		return []string{}, false
	}
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 50
	}

	comment_query := bleve.NewMatchQuery(keyword)
	comment_query.SetField("comment_data")
	search_query := bleve_query.Query(comment_query)
	if set_id != "" {
		set_query := bleve.NewTermQuery(strings.ToLower(set_id))
		set_query.SetField("set_id")
		search_query = bleve.NewConjunctionQuery(search_query, set_query)
	}

	request := bleve.NewSearchRequestOptions(search_query, limit, offset, false)
	request.SortBy([]string{"-date", "_id"})

	Search_bbs_index_check_runtime()
	bbs_search_index_lock.RLock()
	if !bbs_search_index_ready || bbs_search_index == nil {
		bbs_search_index_lock.RUnlock()
		return []string{}, false
	}
	result, err := bbs_search_index.Search(request)
	bbs_search_index_lock.RUnlock()
	if err != nil {
		Search_bbs_index_mark_unavailable("BBS comment read failed", err)
		return []string{}, false
	}

	data_list := make([]string, 0, len(result.Hits))
	for _, hit := range result.Hits {
		data_list = append(data_list, hit.ID)
	}
	return data_list, true
}
