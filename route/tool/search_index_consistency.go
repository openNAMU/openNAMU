package tool

import (
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	encoding_json "encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	bleve_index "github.com/blevesearch/bleve_index_api"
)

const search_index_consistency_batch_size = 100
const search_index_consistency_interval = 7 * 24 * time.Hour
const search_index_consistency_state_version = 1

var err_search_index_consistency_state_invalid = errors.New("invalid search-index audit checkpoint")

type Search_index_consistency_result struct {
	Index_name        string
	Phase             string
	Processed         int
	Missing           int
	Stale             int
	Orphaned          int
	Repaired          int
	Cycle_completed   bool
	Rebuild_requested bool
	Skipped           bool
	Next_due          time.Time
	Total_missing     uint64
	Total_stale       uint64
	Total_orphaned    uint64
	Total_repaired    uint64
}

type search_index_consistency_state struct {
	State_version     int       `json:"state_version"`
	Phase             string    `json:"phase"`
	Database_identity string    `json:"database_identity"`
	Index_version     string    `json:"index_version"`
	Index_signature   string    `json:"index_signature"`
	Cycle_started     time.Time `json:"cycle_started"`
	Last_completed    time.Time `json:"last_completed"`
	SQL_cursor_id     string    `json:"sql_cursor_id,omitempty"`
	SQL_cursor_code   string    `json:"sql_cursor_code,omitempty"`
	SQL_cursor_valid  bool      `json:"sql_cursor_valid,omitempty"`
	Index_after       []string  `json:"index_after,omitempty"`
	Checked           uint64    `json:"checked"`
	Missing           uint64    `json:"missing"`
	Stale             uint64    `json:"stale"`
	Orphaned          uint64    `json:"orphaned"`
	Repair_count      uint64    `json:"repair_count"`
}

func Search_index_fingerprint_field_mapping() *mapping.FieldMapping {
	field_mapping := mapping.NewKeywordFieldMapping()
	field_mapping.Index = false
	field_mapping.Store = true
	field_mapping.IncludeTermVectors = false
	field_mapping.IncludeInAll = false
	field_mapping.DocValues = false
	return field_mapping
}

func Search_index_content_fingerprint(values ...string) string {
	hash := sha256.New()
	var length_buffer [8]byte
	for _, value := range values {
		binary.BigEndian.PutUint64(length_buffer[:], uint64(len(value)))
		_, _ = hash.Write(length_buffer[:])
		_, _ = hash.Write([]byte(value))
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func Search_index_document(doc_name string, data string) Search_document {
	title_search := Do_remove_spaces(doc_name)
	return Search_document{
		Title_search: title_search,
		Data:         data,
		Fingerprint:  Search_index_content_fingerprint(doc_name, title_search, data),
	}
}

func Bbs_search_index_prepare_document(document Bbs_search_document) Bbs_search_document {
	if document.Comment_code != "" {
		values := []string{"comment", document.Set_id, document.Set_code, document.Comment_code, document.Comment_data, document.Date}
		document.Fingerprint = Search_index_content_fingerprint(values...)
		return document
	}

	fingerprint_tags := append([]string(nil), document.Tags...)
	sort.Strings(fingerprint_tags)
	values := []string{"post", document.Set_id, document.Set_code, document.Title, document.Prefix, strconv.Itoa(len(fingerprint_tags))}
	values = append(values, fingerprint_tags...)
	values = append(values, document.Data, document.Date)
	document.Fingerprint = Search_index_content_fingerprint(values...)
	return document
}

func Search_index_consistency_audit_batch(db *sql.DB, state_dir string, now time.Time) (Search_index_consistency_result, error) {
	result := Search_index_consistency_result{Index_name: "document"}
	if now.IsZero() {
		now = time.Now()
	}
	if db == nil {
		return result, errors.New("document index audit requires a database")
	}
	Search_index_check_runtime()

	search_index_update_lock.Lock()
	var rebuild_reason string
	var audit_err error
	func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				rebuild_reason = "document index consistency audit panicked"
				audit_err = fmt.Errorf("document index audit panic: %v", recovered)
			}
		}()
		result, rebuild_reason, audit_err = Search_index_consistency_audit_document_locked(db, state_dir, now)
	}()
	search_index_update_lock.Unlock()
	if rebuild_reason != "" {
		Search_index_request_rebuild(rebuild_reason, audit_err)
		result.Rebuild_requested = true
	}
	return result, audit_err
}

func Search_index_consistency_audit_document_locked(db *sql.DB, state_dir string, now time.Time) (Search_index_consistency_result, string, error) {
	result := Search_index_consistency_result{Index_name: "document"}
	search_index_lock.RLock()
	index := search_index
	ready := search_index_ready && index != nil
	search_index_lock.RUnlock()
	if !ready {
		result.Skipped = true
		return result, "", nil
	}

	generation := search_index_recovery.rebuild_generation.Load()
	signature, signature_err := Search_index_consistency_signature(search_index_directory, search_index_version_file, search_index_version, generation)
	if signature_err != nil {
		return result, "document index consistency audit could not identify the active index", signature_err
	}
	index_count, err := index.DocCount()
	if err != nil {
		return result, "document index consistency audit could not read the index", err
	}

	state_path := Search_index_consistency_state_path(state_dir, "document")
	database_identity := Search_index_consistency_database_identity()
	state, exists, err := Search_index_consistency_read_state(state_path)
	if err != nil && !errors.Is(err, err_search_index_consistency_state_invalid) {
		return result, "", err
	}
	if errors.Is(err, err_search_index_consistency_state_invalid) {
		log.Printf("[SEARCH] invalid document-index audit checkpoint; restarting the bounded audit: %v", err)
		exists = false
	}
	if !exists || state.State_version != search_index_consistency_state_version || state.Database_identity != database_identity || state.Index_version != search_index_version || state.Index_signature != signature || (state.Phase != "sql" && state.Phase != "index" && state.Phase != "complete") {
		state = Search_index_consistency_new_state("sql", database_identity, search_index_version, signature, now)
	}
	if state.Phase == "complete" {
		result.Phase = state.Phase
		result.Next_due = state.Last_completed.Add(search_index_consistency_interval)
		if now.Before(result.Next_due) {
			result.Skipped = true
			return result, "", nil
		}
		last_completed := state.Last_completed
		state = Search_index_consistency_new_state("sql", database_identity, search_index_version, signature, now)
		state.Last_completed = last_completed
	}

	result.Phase = state.Phase
	if state.Phase == "sql" {
		titles, err := Search_index_consistency_document_key_page(db, state.SQL_cursor_id)
		if err != nil {
			return result, "", fmt.Errorf("read document audit source page: %w", err)
		}
		for _, title := range titles {
			state.SQL_cursor_id = title
			state.Checked++
			result.Processed++
			data, exists, err := Search_index_consistency_document_source(db, title)
			if err != nil {
				return result, "", fmt.Errorf("read document audit source %q: %w", title, err)
			}
			if !exists {
				continue
			}
			expected := Search_index_document(title, data).Fingerprint
			actual, indexed, err := Search_index_consistency_index_fingerprint(index, title)
			if err != nil {
				return result, "document index consistency audit encountered an index read failure", err
			}
			if !indexed {
				state.Missing++
				result.Missing++
			} else if actual != expected {
				state.Stale++
				result.Stale++
			} else {
				continue
			}
			if err := Search_index_sync_locked(db, title); err != nil {
				return result, "document index consistency audit could not repair a source record", fmt.Errorf("repair document index entry %q: %w", title, err)
			}
			state.Repair_count++
			result.Repaired++
		}
		if len(titles) < search_index_consistency_batch_size {
			state.Phase = "index"
			state.SQL_cursor_id = ""
			state.Index_after = nil
		}
	} else if state.Phase == "index" {
		request := bleve.NewSearchRequest(bleve.NewMatchAllQuery())
		request.Size = search_index_consistency_batch_size
		request.SortBy([]string{"_id"})
		request.SearchAfter = state.Index_after
		search_result, err := index.Search(request)
		if err != nil {
			return result, "document index consistency audit encountered an index scan failure", err
		}
		for _, hit := range search_result.Hits {
			state.Index_after = hit.Sort
			if len(state.Index_after) == 0 {
				state.Index_after = []string{hit.ID}
			}
			state.Checked++
			result.Processed++
			data, exists, err := Search_index_consistency_document_source(db, hit.ID)
			if err != nil {
				return result, "", fmt.Errorf("read document audit source %q: %w", hit.ID, err)
			}
			actual, indexed, err := Search_index_consistency_index_fingerprint(index, hit.ID)
			if err != nil {
				return result, "document index consistency audit encountered an index read failure", err
			}
			if !exists {
				state.Orphaned++
				result.Orphaned++
			} else {
				expected := Search_index_document(hit.ID, data).Fingerprint
				if indexed && actual == expected {
					continue
				}
				state.Stale++
				result.Stale++
			}
			if err := Search_index_sync_locked(db, hit.ID); err != nil {
				return result, "document index consistency audit could not repair an index record", fmt.Errorf("repair document index entry %q: %w", hit.ID, err)
			}
			state.Repair_count++
			result.Repaired++
		}
		if len(search_result.Hits) < search_index_consistency_batch_size {
			state.Phase = "complete"
			state.Last_completed = now
			state.Index_after = nil
			result.Cycle_completed = true
		}
	}

	result.Phase = state.Phase
	if state.Phase == "complete" {
		result.Next_due = state.Last_completed.Add(search_index_consistency_interval)
	}
	result.Total_missing = state.Missing
	result.Total_stale = state.Stale
	result.Total_orphaned = state.Orphaned
	result.Total_repaired = state.Repair_count
	if err := Search_index_consistency_write_state(state_path, state); err != nil {
		return result, "", err
	}
	if state.Repair_count >= Search_index_consistency_rebuild_threshold(index_count) {
		return result, "document index consistency audit found widespread divergence", nil
	}
	return result, "", nil
}

func Search_bbs_index_consistency_audit_batch(db *sql.DB, state_dir string, now time.Time) (Search_index_consistency_result, error) {
	result := Search_index_consistency_result{Index_name: "bbs"}
	if now.IsZero() {
		now = time.Now()
	}
	if db == nil {
		return result, errors.New("BBS index audit requires a database")
	}
	Search_bbs_index_check_runtime()

	bbs_search_index_update_lock.Lock()
	var rebuild_reason string
	var audit_err error
	func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				rebuild_reason = "BBS index consistency audit panicked"
				audit_err = fmt.Errorf("BBS index audit panic: %v", recovered)
			}
		}()
		result, rebuild_reason, audit_err = Search_bbs_index_consistency_audit_locked(db, state_dir, now)
	}()
	bbs_search_index_update_lock.Unlock()
	if rebuild_reason != "" {
		marker_err := Search_bbs_index_request_rebuild(rebuild_reason, audit_err)
		result.Rebuild_requested = true
		if marker_err != nil {
			audit_err = errors.Join(audit_err, marker_err)
		}
	}
	return result, audit_err
}

func Search_bbs_index_consistency_audit_locked(db *sql.DB, state_dir string, now time.Time) (Search_index_consistency_result, string, error) {
	result := Search_index_consistency_result{Index_name: "bbs"}
	bbs_search_index_lock.RLock()
	index := bbs_search_index
	ready := bbs_search_index_ready && index != nil
	bbs_search_index_lock.RUnlock()
	if !ready {
		result.Skipped = true
		return result, "", nil
	}

	generation := bbs_search_index_recovery.rebuild_generation.Load()
	signature, signature_err := Search_index_consistency_signature(bbs_search_index_directory, bbs_search_index_version_file, bbs_search_index_version, generation)
	if signature_err != nil {
		return result, "BBS index consistency audit could not identify the active index", signature_err
	}
	index_count, err := index.DocCount()
	if err != nil {
		return result, "BBS index consistency audit could not read the index", err
	}

	state_path := Search_index_consistency_state_path(state_dir, "bbs")
	database_identity := Search_index_consistency_database_identity()
	state, exists, err := Search_index_consistency_read_state(state_path)
	if err != nil && !errors.Is(err, err_search_index_consistency_state_invalid) {
		return result, "", err
	}
	if errors.Is(err, err_search_index_consistency_state_invalid) {
		log.Printf("[SEARCH] invalid BBS-index audit checkpoint; restarting the bounded audit: %v", err)
		exists = false
	}
	if !exists || state.State_version != search_index_consistency_state_version || state.Database_identity != database_identity || state.Index_version != bbs_search_index_version || state.Index_signature != signature || (state.Phase != "posts" && state.Phase != "comments" && state.Phase != "index" && state.Phase != "complete") {
		state = Search_index_consistency_new_state("posts", database_identity, bbs_search_index_version, signature, now)
	}
	if state.Phase == "complete" {
		result.Phase = state.Phase
		result.Next_due = state.Last_completed.Add(search_index_consistency_interval)
		if now.Before(result.Next_due) {
			result.Skipped = true
			return result, "", nil
		}
		last_completed := state.Last_completed
		state = Search_index_consistency_new_state("posts", database_identity, bbs_search_index_version, signature, now)
		state.Last_completed = last_completed
	}

	result.Phase = state.Phase
	switch state.Phase {
	case "posts":
		keys, err := Search_index_consistency_bbs_key_page(db, "title", state.SQL_cursor_id, state.SQL_cursor_code, state.SQL_cursor_valid)
		if err != nil {
			return result, "", fmt.Errorf("read BBS post audit source page: %w", err)
		}
		for _, key := range keys {
			state.SQL_cursor_id = key[0]
			state.SQL_cursor_code = key[1]
			state.SQL_cursor_valid = true
			state.Checked++
			result.Processed++
			document, exists, err := Search_index_consistency_bbs_document_source(db, key[0], key[1])
			if err != nil {
				return result, "", fmt.Errorf("read BBS post audit source %q/%q: %w", key[0], key[1], err)
			}
			if !exists {
				continue
			}
			document = Bbs_search_index_prepare_document(document)
			index_key := Search_bbs_index_key(key[0], key[1])
			actual, indexed, err := Search_index_consistency_index_fingerprint(index, index_key)
			if err != nil {
				return result, "BBS index consistency audit encountered an index read failure", err
			}
			if !indexed {
				state.Missing++
				result.Missing++
			} else if actual != document.Fingerprint {
				state.Stale++
				result.Stale++
			} else {
				continue
			}
			if err := Search_bbs_index_sync_locked(db, key[0], key[1]); err != nil {
				return result, "BBS index consistency audit could not repair a source post", fmt.Errorf("repair BBS post %q/%q: %w", key[0], key[1], err)
			}
			state.Repair_count++
			result.Repaired++
		}
		if len(keys) < search_index_consistency_batch_size {
			state.Phase = "comments"
			state.SQL_cursor_id = ""
			state.SQL_cursor_code = ""
			state.SQL_cursor_valid = false
		}
	case "comments":
		keys, err := Search_index_consistency_bbs_key_page(db, "comment", state.SQL_cursor_id, state.SQL_cursor_code, state.SQL_cursor_valid)
		if err != nil {
			return result, "", fmt.Errorf("read BBS comment audit source page: %w", err)
		}
		for _, key := range keys {
			state.SQL_cursor_id = key[0]
			state.SQL_cursor_code = key[1]
			state.SQL_cursor_valid = true
			state.Checked++
			result.Processed++
			root_id, root_code, comment_code, valid := Search_bbs_index_comment_location(key[0], key[1])
			if !valid {
				continue
			}
			document, exists, err := Search_index_consistency_bbs_comment_source(db, root_id, root_code, comment_code)
			if err != nil {
				return result, "", fmt.Errorf("read BBS comment audit source %q/%q: %w", key[0], key[1], err)
			}
			if !exists {
				continue
			}
			document = Bbs_search_index_prepare_document(document)
			index_key := Search_bbs_index_comment_key(root_id, root_code, comment_code)
			actual, indexed, err := Search_index_consistency_index_fingerprint(index, index_key)
			if err != nil {
				return result, "BBS index consistency audit encountered an index read failure", err
			}
			if !indexed {
				state.Missing++
				result.Missing++
			} else if actual != document.Fingerprint {
				state.Stale++
				result.Stale++
			} else {
				continue
			}
			if err := Search_bbs_index_sync_comment_locked(db, root_id, root_code, comment_code); err != nil {
				return result, "BBS index consistency audit could not repair a source comment", fmt.Errorf("repair BBS comment %q/%q/%q: %w", root_id, root_code, comment_code, err)
			}
			state.Repair_count++
			result.Repaired++
		}
		if len(keys) < search_index_consistency_batch_size {
			state.Phase = "index"
			state.SQL_cursor_id = ""
			state.SQL_cursor_code = ""
			state.SQL_cursor_valid = false
			state.Index_after = nil
		}
	case "index":
		request := bleve.NewSearchRequest(bleve.NewMatchAllQuery())
		request.Size = search_index_consistency_batch_size
		request.SortBy([]string{"_id"})
		request.SearchAfter = state.Index_after
		search_result, err := index.Search(request)
		if err != nil {
			return result, "BBS index consistency audit encountered an index scan failure", err
		}
		for _, hit := range search_result.Hits {
			state.Index_after = hit.Sort
			if len(state.Index_after) == 0 {
				state.Index_after = []string{hit.ID}
			}
			state.Checked++
			result.Processed++
			comment_set_id, comment_set_code, comment_code, is_comment := Search_bbs_index_comment_key_data(hit.ID)
			if is_comment {
				document, exists, err := Search_index_consistency_bbs_comment_source(db, comment_set_id, comment_set_code, comment_code)
				if err != nil {
					return result, "", fmt.Errorf("read BBS comment audit source for %q: %w", hit.ID, err)
				}
				if !exists {
					state.Orphaned++
					result.Orphaned++
				} else {
					document = Bbs_search_index_prepare_document(document)
					actual, indexed, err := Search_index_consistency_index_fingerprint(index, hit.ID)
					if err != nil {
						return result, "BBS index consistency audit encountered an index read failure", err
					}
					if indexed && actual == document.Fingerprint {
						continue
					}
					state.Stale++
					result.Stale++
				}
				if err := Search_bbs_index_sync_comment_locked(db, comment_set_id, comment_set_code, comment_code); err != nil {
					return result, "BBS index consistency audit could not repair an index comment", fmt.Errorf("repair BBS index comment %q: %w", hit.ID, err)
				}
				state.Repair_count++
				result.Repaired++
				continue
			}

			set_id, set_code, valid_post_key := Search_bbs_index_key_data(hit.ID)
			valid_post_key = valid_post_key && strings.Count(hit.ID, "\x00") == 1
			if !valid_post_key {
				state.Orphaned++
				result.Orphaned++
				if err := Search_bbs_index_change_locked(hit.ID, bbs_search_index_change{Deleted: true}); err != nil {
					return result, "BBS index consistency audit could not remove an invalid index key", err
				}
				state.Repair_count++
				result.Repaired++
				continue
			}
			document, exists, err := Search_index_consistency_bbs_document_source(db, set_id, set_code)
			if err != nil {
				return result, "", fmt.Errorf("read BBS post audit source for %q: %w", hit.ID, err)
			}
			if !exists {
				state.Orphaned++
				result.Orphaned++
			} else {
				document = Bbs_search_index_prepare_document(document)
				actual, indexed, err := Search_index_consistency_index_fingerprint(index, hit.ID)
				if err != nil {
					return result, "BBS index consistency audit encountered an index read failure", err
				}
				if indexed && actual == document.Fingerprint {
					continue
				}
				state.Stale++
				result.Stale++
			}
			if err := Search_bbs_index_sync_locked(db, set_id, set_code); err != nil {
				return result, "BBS index consistency audit could not repair an index post", fmt.Errorf("repair BBS index post %q: %w", hit.ID, err)
			}
			state.Repair_count++
			result.Repaired++
		}
		if len(search_result.Hits) < search_index_consistency_batch_size {
			state.Phase = "complete"
			state.Last_completed = now
			state.Index_after = nil
			result.Cycle_completed = true
		}
	}

	result.Phase = state.Phase
	if state.Phase == "complete" {
		result.Next_due = state.Last_completed.Add(search_index_consistency_interval)
	}
	result.Total_missing = state.Missing
	result.Total_stale = state.Stale
	result.Total_orphaned = state.Orphaned
	result.Total_repaired = state.Repair_count
	if err := Search_index_consistency_write_state(state_path, state); err != nil {
		return result, "", err
	}
	if state.Repair_count >= Search_index_consistency_rebuild_threshold(index_count) {
		return result, "BBS index consistency audit found widespread divergence", nil
	}
	return result, "", nil
}

func Search_index_consistency_state_path(state_dir string, index_name string) string {
	return filepath.Join(state_dir, ".opennamu_"+index_name+"_search_index_audit.json")
}

func Search_index_consistency_new_state(phase string, database_identity string, index_version string, signature string, now time.Time) search_index_consistency_state {
	return search_index_consistency_state{
		State_version:     search_index_consistency_state_version,
		Phase:             phase,
		Database_identity: database_identity,
		Index_version:     index_version,
		Index_signature:   signature,
		Cycle_started:     now,
	}
}

func Search_index_consistency_signature(index_directory string, version_file string, expected_version string, generation uint64) (string, error) {
	version_valid, err := search_index_version_matches(version_file, expected_version)
	if err != nil {
		return "", err
	}
	if !version_valid {
		return "", fmt.Errorf("index version marker %q is not %q", version_file, expected_version)
	}
	meta_path := filepath.Join(index_directory, "index_meta.json")
	meta_info, err := os.Stat(meta_path)
	if err != nil {
		return "", fmt.Errorf("stat index metadata %q: %w", meta_path, err)
	}
	if !meta_info.Mode().IsRegular() {
		return "", fmt.Errorf("index metadata %q is not a regular file", meta_path)
	}
	absolute_directory, err := filepath.Abs(index_directory)
	if err != nil {
		return "", err
	}
	return Search_index_content_fingerprint(
		absolute_directory,
		expected_version,
		strconv.FormatUint(generation, 10),
		strconv.FormatInt(meta_info.Size(), 10),
		strconv.FormatInt(meta_info.ModTime().UnixNano(), 10),
	), nil
}

func Search_index_consistency_database_identity() string {
	db_set := Get_DB_set()
	identity := []string{db_set["db_type"], db_set["db_name"]}
	if db_set["db_type"] == "mysql" {
		mysql_set := Get_DB_set_MySQL(db_set)
		identity = append(identity, mysql_set["db_mysql_host"], mysql_set["db_mysql_port"])
	}
	return Search_index_content_fingerprint(identity...)
}

func Search_index_consistency_read_state(path string) (search_index_consistency_state, bool, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return search_index_consistency_state{}, false, nil
	}
	if err != nil {
		return search_index_consistency_state{}, false, err
	}
	state := search_index_consistency_state{}
	if err := encoding_json.Unmarshal(data, &state); err != nil {
		return search_index_consistency_state{}, false, fmt.Errorf("%w: %v", err_search_index_consistency_state_invalid, err)
	}
	valid_phase := state.Phase == "sql" || state.Phase == "posts" || state.Phase == "comments" || state.Phase == "index" || state.Phase == "complete"
	invalid_counters := state.Missing > state.Checked || state.Stale > state.Checked-state.Missing || state.Orphaned > state.Checked-state.Missing-state.Stale || state.Repair_count > state.Checked
	if state.State_version != search_index_consistency_state_version || !valid_phase || state.Database_identity == "" || state.Index_version == "" || state.Index_signature == "" || state.Cycle_started.IsZero() || (state.Phase == "complete" && state.Last_completed.IsZero()) || invalid_counters || len(state.Index_after) > 1 {
		return search_index_consistency_state{}, false, fmt.Errorf("%w: unsupported or inconsistent checkpoint data", err_search_index_consistency_state_invalid)
	}
	return state, true, nil
}

func Search_index_consistency_write_state(path string, state search_index_consistency_state) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create search-index audit state directory: %w", err)
	}
	data, err := encoding_json.Marshal(state)
	if err != nil {
		return fmt.Errorf("encode search-index audit state: %w", err)
	}
	temporary_file, err := os.CreateTemp(filepath.Dir(path), ".opennamu_search_index_audit.tmp-*")
	if err != nil {
		return fmt.Errorf("create search-index audit checkpoint: %w", err)
	}
	temporary_path := temporary_file.Name()
	defer os.Remove(temporary_path)
	if _, err := temporary_file.Write(data); err != nil {
		_ = temporary_file.Close()
		return fmt.Errorf("write search-index audit checkpoint: %w", err)
	}
	if err := temporary_file.Sync(); err != nil {
		_ = temporary_file.Close()
		return fmt.Errorf("sync search-index audit checkpoint: %w", err)
	}
	if err := temporary_file.Close(); err != nil {
		return fmt.Errorf("close search-index audit checkpoint: %w", err)
	}
	if err := os.Chmod(temporary_path, 0o600); err != nil {
		return fmt.Errorf("set search-index audit checkpoint permissions: %w", err)
	}
	if err := os.Rename(temporary_path, path); err != nil {
		return fmt.Errorf("publish search-index audit checkpoint: %w", err)
	}
	return nil
}

func Search_index_consistency_document_key_page(db *sql.DB, cursor string) ([]string, error) {
	query := "select distinct title collate binary from data where title collate binary > ? collate binary and title != '' order by title collate binary limit ?"
	values := []any{cursor, search_index_consistency_batch_size}
	if cursor == "" {
		query = "select distinct title collate binary from data where title != '' order by title collate binary limit ?"
		values = []any{search_index_consistency_batch_size}
	}
	if Get_DB_type() == "mysql" {
		query = "select distinct binary title from data where binary title > binary ? and binary title != binary '' order by binary title limit ?"
		values = []any{cursor, search_index_consistency_batch_size}
		if cursor == "" {
			query = "select distinct binary title from data where binary title != binary '' order by binary title limit ?"
			values = []any{search_index_consistency_batch_size}
		}
	}
	rows, err := db.Query(DB_change(query), values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	titles := make([]string, 0, search_index_consistency_batch_size)
	for rows.Next() {
		title := ""
		if err := rows.Scan(&title); err != nil {
			return nil, err
		}
		titles = append(titles, title)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return titles, rows.Close()
}

func Search_index_consistency_document_source(db *sql.DB, title string) (string, bool, error) {
	query := "select coalesce(data, '') from data where title collate binary = ? collate binary limit 1"
	if Get_DB_type() == "mysql" {
		query = "select coalesce(data, '') from data where binary title = binary ? limit 1"
	}
	data := ""
	err := db.QueryRow(DB_change(query), title).Scan(&data)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return data, true, nil
}

func Search_index_consistency_bbs_document_source(db *sql.DB, set_id string, set_code string) (Bbs_search_document, bool, error) {
	document := Bbs_search_document{Set_id: strings.ToLower(set_id)}
	set_id_condition := "set_id collate binary = ? collate binary"
	set_code_condition := "set_code collate binary = ? collate binary"
	if Get_DB_type() == "mysql" {
		set_id_condition = "binary set_id = binary ?"
		set_code_condition = "binary set_code = binary ?"
	}
	query := "select set_name, set_data from bbs_data where " + set_id_condition + " and " + set_code_condition + " and set_name in ('title', 'prefix', 'tag', 'data', 'date')"
	rows, err := db.Query(DB_change(query), set_id, set_code)
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

func Search_index_consistency_bbs_comment_source(db *sql.DB, set_id string, set_code string, comment_code string) (Bbs_search_document, bool, error) {
	comment_parts := strings.Split(comment_code, "-")
	if len(comment_parts) == 0 || comment_parts[len(comment_parts)-1] == "" {
		return Bbs_search_document{}, false, nil
	}
	comment_set_id := set_id + "-" + set_code
	if len(comment_parts) > 1 {
		comment_set_id += "-" + strings.Join(comment_parts[:len(comment_parts)-1], "-")
	}
	comment_set_code := comment_parts[len(comment_parts)-1]
	set_id_condition := "set_id collate binary = ? collate binary"
	set_code_condition := "set_code collate binary = ? collate binary"
	if Get_DB_type() == "mysql" {
		set_id_condition = "binary set_id = binary ?"
		set_code_condition = "binary set_code = binary ?"
	}
	query := "select set_name, set_data from bbs_data where " + set_id_condition + " and " + set_code_condition + " and set_name in ('comment', 'comment_date')"
	document := Bbs_search_document{Set_id: set_id, Set_code: set_code, Comment_code: comment_code}
	rows, err := db.Query(DB_change(query), comment_set_id, comment_set_code)
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

func Search_index_consistency_index_fingerprint(index bleve.Index, document_id string) (string, bool, error) {
	document, err := index.Document(document_id)
	if err != nil {
		return "", false, err
	}
	if document == nil {
		return "", false, nil
	}
	fingerprint := ""
	field_count := 0
	document.VisitFields(func(field bleve_index.Field) {
		if field.Name() == "fingerprint" {
			fingerprint = string(field.Value())
			field_count++
		}
	})
	if field_count != 1 {
		return "", true, nil
	}
	return fingerprint, true, nil
}

func Search_index_consistency_bbs_key_page(db *sql.DB, set_name string, cursor_id string, cursor_code string, cursor_valid bool) ([][2]string, error) {
	query := "select distinct set_id collate binary, set_code collate binary from bbs_data where set_name = ? and (set_id collate binary > ? collate binary or (set_id collate binary = ? collate binary and set_code collate binary > ? collate binary)) order by set_id collate binary, set_code collate binary limit ?"
	values := []any{set_name, cursor_id, cursor_id, cursor_code, search_index_consistency_batch_size}
	if !cursor_valid {
		query = "select distinct set_id collate binary, set_code collate binary from bbs_data where set_name = ? order by set_id collate binary, set_code collate binary limit ?"
		values = []any{set_name, search_index_consistency_batch_size}
	}
	if Get_DB_type() == "mysql" {
		query = "select distinct binary set_id, binary set_code from bbs_data where set_name = ? and (binary set_id > binary ? or (binary set_id = binary ? and binary set_code > binary ?)) order by binary set_id, binary set_code limit ?"
		values = []any{set_name, cursor_id, cursor_id, cursor_code, search_index_consistency_batch_size}
		if !cursor_valid {
			query = "select distinct binary set_id, binary set_code from bbs_data where set_name = ? order by binary set_id, binary set_code limit ?"
			values = []any{set_name, search_index_consistency_batch_size}
		}
	}
	rows, err := db.Query(DB_change(query), values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	keys := make([][2]string, 0, search_index_consistency_batch_size)
	for rows.Next() {
		key := [2]string{}
		if err := rows.Scan(&key[0], &key[1]); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return keys, rows.Close()
}

func Search_index_consistency_rebuild_threshold(index_count uint64) uint64 {
	threshold := index_count / 20
	if threshold < 10 {
		threshold = 10
	}
	if threshold > search_index_consistency_batch_size {
		threshold = search_index_consistency_batch_size
	}
	return threshold
}

func Search_index_sync_locked(db *sql.DB, doc_name string) error {
	data, exists, err := Search_index_consistency_document_source(db, doc_name)
	if err != nil {
		return err
	}
	change := search_index_change{Data: data, Deleted: !exists}
	return Search_index_change_locked(doc_name, change)
}

func Search_bbs_index_sync_locked(db *sql.DB, set_id string, set_code string) error {
	document, exists, err := Search_index_consistency_bbs_document_source(db, set_id, set_code)
	if err != nil {
		return err
	}
	key := Search_bbs_index_key(set_id, set_code)
	return Search_bbs_index_change_locked(key, bbs_search_index_change{Document: document, Deleted: !exists})
}

func Search_bbs_index_sync_comment_locked(db *sql.DB, set_id string, set_code string, comment_code string) error {
	document, exists, err := Search_index_consistency_bbs_comment_source(db, set_id, set_code, comment_code)
	if err != nil {
		return err
	}
	key := Search_bbs_index_comment_key(set_id, set_code, comment_code)
	return Search_bbs_index_change_locked(key, bbs_search_index_change{Document: document, Deleted: !exists})
}
