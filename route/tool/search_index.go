package tool

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/blevesearch/bleve/v2"
	_ "github.com/blevesearch/bleve/v2/analysis/lang/cjk"
	"github.com/blevesearch/bleve/v2/mapping"
)

type Search_document struct {
	Title_search string `json:"title_search"`
	Data         string `json:"data"`
}

type search_index_change struct {
	Data    string
	Deleted bool
}

const search_index_directory = "data/bleve"
const search_index_version_file = "data/bleve.version"
const search_index_version = "2"

var search_index bleve.Index
var search_index_lock sync.RWMutex
var search_index_ready bool
var search_index_start_once sync.Once
var search_index_pending = map[string]search_index_change{}

func search_index_field_mapping() *mapping.FieldMapping {
	field_mapping := mapping.NewTextFieldMapping()
	field_mapping.Analyzer = "cjk"
	field_mapping.Store = false
	field_mapping.IncludeTermVectors = false
	field_mapping.IncludeInAll = false
	field_mapping.DocValues = false
	return field_mapping
}

func search_index_mapping() *mapping.IndexMappingImpl {
	index_mapping := bleve.NewIndexMapping()
	index_mapping.DefaultAnalyzer = "cjk"
	index_mapping.StoreDynamic = false
	index_mapping.IndexDynamic = false
	index_mapping.DocValuesDynamic = false

	document_mapping := mapping.NewDocumentMapping()
	document_mapping.Dynamic = false
	document_mapping.AddFieldMappingsAt("title_search", search_index_field_mapping())
	document_mapping.AddFieldMappingsAt("data", search_index_field_mapping())
	index_mapping.DefaultMapping = document_mapping

	return index_mapping
}

func Search_index_start() {
	search_index_start_once.Do(func() {
		go search_index_open()
	})
}

func Search_index_ready() bool {
	search_index_lock.RLock()
	defer search_index_lock.RUnlock()
	return search_index_ready
}

func search_index_version_valid() bool {
	data, err := os.ReadFile(search_index_version_file)
	return err == nil && string(data) == search_index_version
}

func search_index_apply_change(index bleve.Index, doc_name string, change search_index_change) error {
	if change.Deleted {
		return index.Delete(doc_name)
	}
	return index.Index(doc_name, Search_document{Do_remove_spaces(doc_name), change.Data})
}

func search_index_set_ready(index bleve.Index) {
	search_index_lock.Lock()
	for doc_name, change := range search_index_pending {
		if err := search_index_apply_change(index, doc_name, change); err != nil {
			log.Printf("[SEARCH] pending document update failed: %v", err)
		}
	}
	search_index_pending = map[string]search_index_change{}
	search_index = index
	search_index_ready = true
	search_index_lock.Unlock()
}

func search_index_open() {
	if err := os.MkdirAll(filepath.Dir(search_index_directory), 0o755); err != nil {
		log.Printf("[SEARCH] index directory failed: %v", err)
		return
	}

	index, err := bleve.Open(search_index_directory)
	if err == nil && search_index_version_valid() {
		search_index_set_ready(index)
		log.Println("[SEARCH] index opened")
		return
	}
	if err == nil {
		_ = index.Close()
		log.Println("[SEARCH] index mapping changed; rebuilding")
	} else if _, stat_err := os.Stat(search_index_directory); stat_err == nil {
		log.Printf("[SEARCH] index open failed: %v", err)
	}

	_ = os.RemoveAll(search_index_directory)
	_ = os.Remove(search_index_version_file)
	search_index_rebuild()
}

func search_index_rebuild() {
	db := DB_connect()
	defer DB_close(db)

	temp_directory := search_index_directory + ".tmp"
	_ = os.RemoveAll(temp_directory)
	index, err := bleve.New(temp_directory, search_index_mapping())
	if err != nil {
		log.Printf("[SEARCH] index create failed: %v", err)
		return
	}

	rows := Query_DB(db, "select title, coalesce(data, '') from data")
	batch := index.NewBatch()
	batch_count := 0
	for rows.Next() {
		title := ""
		data := ""
		if rows.Scan(&title, &data) != nil || title == "" {
			continue
		}
		if err := batch.Index(title, Search_document{Do_remove_spaces(title), data}); err != nil {
			log.Printf("[SEARCH] document index failed: %v", err)
			continue
		}
		batch_count++
		if batch_count >= 500 {
			if err := index.Batch(batch); err != nil {
				rows.Close()
				index.Close()
				log.Printf("[SEARCH] batch index failed: %v", err)
				return
			}
			batch = index.NewBatch()
			batch_count = 0
		}
	}
	rows.Close()

	if batch_count > 0 {
		if err := index.Batch(batch); err != nil {
			index.Close()
			log.Printf("[SEARCH] batch index failed: %v", err)
			return
		}
	}
	if err := index.Close(); err != nil {
		log.Printf("[SEARCH] index close failed: %v", err)
		return
	}
	if err := os.Rename(temp_directory, search_index_directory); err != nil {
		log.Printf("[SEARCH] index install failed: %v", err)
		return
	}
	if err := os.WriteFile(search_index_version_file, []byte(search_index_version), 0o644); err != nil {
		_ = os.RemoveAll(search_index_directory)
		log.Printf("[SEARCH] index version failed: %v", err)
		return
	}

	index, err = bleve.Open(search_index_directory)
	if err != nil {
		log.Printf("[SEARCH] index reopen failed: %v", err)
		return
	}
	search_index_set_ready(index)
	log.Println("[SEARCH] index built")
}

func Search_index_update(doc_name string, data string) {
	if doc_name == "" {
		return
	}
	search_index_lock.Lock()
	defer search_index_lock.Unlock()
	if !search_index_ready {
		search_index_pending[doc_name] = search_index_change{Data: data}
		return
	}
	if err := search_index.Index(doc_name, Search_document{Do_remove_spaces(doc_name), data}); err != nil {
		log.Printf("[SEARCH] document update failed: %v", err)
	}
}

func Search_index_delete(doc_name string) {
	if doc_name == "" {
		return
	}
	search_index_lock.Lock()
	defer search_index_lock.Unlock()
	if !search_index_ready {
		search_index_pending[doc_name] = search_index_change{Deleted: true}
		return
	}
	if err := search_index.Delete(doc_name); err != nil {
		log.Printf("[SEARCH] document delete failed: %v", err)
	}
}

func Search_index_sync(db *sql.DB, doc_name string) {
	if doc_name == "" {
		return
	}
	data := ""
	if QueryRow_DB(db, "select coalesce(data, '') from data where title = ?", []any{&data}, doc_name) {
		Search_index_update(doc_name, data)
	} else {
		Search_index_delete(doc_name)
	}
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

	var request *bleve.SearchRequest
	if search_type == "title" {
		query := bleve.NewWildcardQuery("*" + strings.ToLower(keyword) + "*")
		query.SetField("title_search")
		request = bleve.NewSearchRequestOptions(query, limit, offset, false)
		request.SortBy([]string{"_id"})
	} else {
		query := bleve.NewMatchQuery(keyword)
		query.SetField("data")
		request = bleve.NewSearchRequestOptions(query, limit, offset, false)
		request.SortBy([]string{"_id"})
	}

	search_index_lock.RLock()
	defer search_index_lock.RUnlock()
	if !search_index_ready {
		return []string{}, false
	}
	result, err := search_index.Search(request)
	if err != nil {
		log.Printf("[SEARCH] search failed: %v", err)
		return []string{}, false
	}

	data_list := make([]string, 0, len(result.Hits))
	for _, hit := range result.Hits {
		data_list = append(data_list, hit.ID)
	}
	return data_list, true
}
