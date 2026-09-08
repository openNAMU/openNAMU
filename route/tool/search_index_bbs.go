package tool

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	bleve_query "github.com/blevesearch/bleve/v2/search/query"
)

type Bbs_search_document struct {
	Set_id string   `json:"set_id"`
	Title  string   `json:"title"`
	Prefix string   `json:"prefix"`
	Tags   []string `json:"tags"`
	Date   string   `json:"date"`
}

type bbs_search_index_change struct {
	Document Bbs_search_document
	Deleted  bool
}

const bbs_search_index_directory = "data/bleve_bbs"
const bbs_search_index_version_file = "data/bleve_bbs.version"
const bbs_search_index_version = "1"

var bbs_search_index bleve.Index
var bbs_search_index_lock sync.RWMutex
var bbs_search_index_ready bool
var bbs_search_index_start_once sync.Once
var bbs_search_index_pending = map[string]bbs_search_index_change{}

func bbs_search_index_keyword_mapping() *mapping.FieldMapping {
	field_mapping := mapping.NewKeywordFieldMapping()
	field_mapping.Store = false
	field_mapping.IncludeTermVectors = false
	field_mapping.IncludeInAll = false
	field_mapping.DocValues = false
	return field_mapping
}

func bbs_search_index_date_mapping() *mapping.FieldMapping {
	field_mapping := bbs_search_index_keyword_mapping()
	field_mapping.DocValues = true
	return field_mapping
}

func bbs_search_index_mapping() *mapping.IndexMappingImpl {
	index_mapping := bleve.NewIndexMapping()
	index_mapping.StoreDynamic = false
	index_mapping.IndexDynamic = false
	index_mapping.DocValuesDynamic = false

	document_mapping := mapping.NewDocumentMapping()
	document_mapping.Dynamic = false
	document_mapping.AddFieldMappingsAt("set_id", bbs_search_index_keyword_mapping())
	document_mapping.AddFieldMappingsAt("title", bbs_search_index_keyword_mapping())
	document_mapping.AddFieldMappingsAt("prefix", bbs_search_index_keyword_mapping())
	document_mapping.AddFieldMappingsAt("tags", bbs_search_index_keyword_mapping())
	document_mapping.AddFieldMappingsAt("date", bbs_search_index_date_mapping())
	index_mapping.DefaultMapping = document_mapping

	return index_mapping
}

func Search_bbs_index_start() {
	bbs_search_index_start_once.Do(func() {
		go search_bbs_index_open()
	})
}

func search_bbs_index_version_valid() bool {
	data, err := os.ReadFile(bbs_search_index_version_file)
	return err == nil && string(data) == bbs_search_index_version
}

func search_bbs_index_key(set_id string, set_code string) string {
	return set_id + "\x00" + set_code
}

func Search_bbs_index_key_data(key string) (string, string, bool) {
	data := strings.SplitN(key, "\x00", 2)
	if len(data) != 2 || data[0] == "" || data[1] == "" {
		return "", "", false
	}
	return data[0], data[1], true
}

func search_bbs_index_apply_change(index bleve.Index, key string, change bbs_search_index_change) error {
	if change.Deleted {
		return index.Delete(key)
	}
	return index.Index(key, change.Document)
}

func search_bbs_index_set_ready(index bleve.Index) {
	bbs_search_index_lock.Lock()
	for key, change := range bbs_search_index_pending {
		if err := search_bbs_index_apply_change(index, key, change); err != nil {
			log.Printf("[SEARCH] pending bbs update failed: %v", err)
		}
	}
	bbs_search_index_pending = map[string]bbs_search_index_change{}
	bbs_search_index = index
	bbs_search_index_ready = true
	bbs_search_index_lock.Unlock()
}

func search_bbs_index_open() {
	if err := os.MkdirAll(filepath.Dir(bbs_search_index_directory), 0o755); err != nil {
		log.Printf("[SEARCH] bbs index directory failed: %v", err)
		return
	}

	index, err := bleve.Open(bbs_search_index_directory)
	if err == nil && search_bbs_index_version_valid() {
		search_bbs_index_set_ready(index)
		log.Println("[SEARCH] bbs index opened")
		return
	}
	if err == nil {
		_ = index.Close()
		log.Println("[SEARCH] bbs index mapping changed; rebuilding")
	} else if _, stat_err := os.Stat(bbs_search_index_directory); stat_err == nil {
		log.Printf("[SEARCH] bbs index open failed: %v", err)
	}

	_ = os.RemoveAll(bbs_search_index_directory)
	_ = os.Remove(bbs_search_index_version_file)
	search_bbs_index_rebuild()
}

func search_bbs_index_rebuild() {
	db := DB_connect()
	defer DB_close(db)

	temp_directory := bbs_search_index_directory + ".tmp"
	_ = os.RemoveAll(temp_directory)
	index, err := bleve.New(temp_directory, bbs_search_index_mapping())
	if err != nil {
		log.Printf("[SEARCH] bbs index create failed: %v", err)
		return
	}

	rows := Query_DB(
		db,
		"select set_id, set_code, set_name, set_data from bbs_data where set_name in ('title', 'prefix', 'tag', 'date') order by set_id, set_code",
	)
	defer rows.Close()

	current_key := ""
	document := Bbs_search_document{}
	batch := index.NewBatch()
	batch_count := 0

	flush := func() bool {
		if current_key == "" || document.Title == "" {
			return true
		}
		if err := batch.Index(current_key, document); err != nil {
			log.Printf("[SEARCH] bbs document index failed: %v", err)
			return false
		}
		batch_count++
		if batch_count < 500 {
			return true
		}
		if err := index.Batch(batch); err != nil {
			log.Printf("[SEARCH] bbs batch index failed: %v", err)
			return false
		}
		batch = index.NewBatch()
		batch_count = 0
		return true
	}

	for rows.Next() {
		set_id := ""
		set_code := ""
		set_name := ""
		set_data := ""
		if rows.Scan(&set_id, &set_code, &set_name, &set_data) != nil {
			continue
		}

		key := search_bbs_index_key(set_id, set_code)
		if current_key != "" && current_key != key {
			if !flush() {
				_ = index.Close()
				return
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
		case "date":
			document.Date = set_data
		}
	}

	if !flush() {
		_ = index.Close()
		return
	}
	if batch_count > 0 {
		if err := index.Batch(batch); err != nil {
			_ = index.Close()
			log.Printf("[SEARCH] bbs batch index failed: %v", err)
			return
		}
	}
	if err := index.Close(); err != nil {
		log.Printf("[SEARCH] bbs index close failed: %v", err)
		return
	}
	if err := os.Rename(temp_directory, bbs_search_index_directory); err != nil {
		log.Printf("[SEARCH] bbs index install failed: %v", err)
		return
	}
	if err := os.WriteFile(bbs_search_index_version_file, []byte(bbs_search_index_version), 0o644); err != nil {
		_ = os.RemoveAll(bbs_search_index_directory)
		log.Printf("[SEARCH] bbs index version failed: %v", err)
		return
	}

	index, err = bleve.Open(bbs_search_index_directory)
	if err != nil {
		log.Printf("[SEARCH] bbs index reopen failed: %v", err)
		return
	}
	search_bbs_index_set_ready(index)
	log.Println("[SEARCH] bbs index built")
}

func search_bbs_index_document(db *sql.DB, set_id string, set_code string) (Bbs_search_document, bool) {
	document := Bbs_search_document{Set_id: strings.ToLower(set_id)}
	rows := Query_DB(
		db,
		"select set_name, set_data from bbs_data where set_id = ? and set_code = ? and set_name in ('title', 'prefix', 'tag', 'date')",
		set_id,
		set_code,
	)
	defer rows.Close()

	for rows.Next() {
		set_name := ""
		set_data := ""
		if rows.Scan(&set_name, &set_data) != nil {
			continue
		}
		switch set_name {
		case "title":
			document.Title = strings.ToLower(set_data)
		case "prefix":
			document.Prefix = strings.ToLower(set_data)
		case "tag":
			document.Tags = append(document.Tags, strings.ToLower(set_data))
		case "date":
			document.Date = set_data
		}
	}

	return document, document.Title != ""
}

func Search_bbs_index_update(db *sql.DB, set_id string, set_code string) {
	if set_id == "" || set_code == "" {
		return
	}
	document, exists := search_bbs_index_document(db, set_id, set_code)
	key := search_bbs_index_key(set_id, set_code)

	bbs_search_index_lock.Lock()
	defer bbs_search_index_lock.Unlock()
	if !exists {
		if !bbs_search_index_ready {
			bbs_search_index_pending[key] = bbs_search_index_change{Deleted: true}
			return
		}
		if err := bbs_search_index.Delete(key); err != nil {
			log.Printf("[SEARCH] bbs document delete failed: %v", err)
		}
		return
	}
	change := bbs_search_index_change{Document: document}
	if !bbs_search_index_ready {
		bbs_search_index_pending[key] = change
		return
	}
	if err := bbs_search_index.Index(key, document); err != nil {
		log.Printf("[SEARCH] bbs document update failed: %v", err)
	}
}

func Search_bbs_index_delete(set_id string, set_code string) {
	if set_id == "" || set_code == "" {
		return
	}
	key := search_bbs_index_key(set_id, set_code)
	bbs_search_index_lock.Lock()
	defer bbs_search_index_lock.Unlock()
	if !bbs_search_index_ready {
		bbs_search_index_pending[key] = bbs_search_index_change{Deleted: true}
		return
	}
	if err := bbs_search_index.Delete(key); err != nil {
		log.Printf("[SEARCH] bbs document delete failed: %v", err)
	}
}

func Search_bbs_index_delete_set(db *sql.DB, set_id string) {
	rows := Query_DB(db, "select set_code from bbs_data where set_name = 'title' and set_id = ?", set_id)
	defer rows.Close()
	for rows.Next() {
		set_code := ""
		if rows.Scan(&set_code) == nil {
			Search_bbs_index_delete(set_id, set_code)
		}
	}
}

func Search_bbs_index_search(keyword string, set_id string, offset int, limit int) ([]string, bool) {
	if keyword == "" || strings.ContainsAny(keyword, "*?%_") {
		return []string{}, false
	}
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 50
	}

	queries := []bleve_query.Query{}
	for _, field := range []string{"title", "prefix", "tags"} {
		query := bleve.NewWildcardQuery("*" + strings.ToLower(keyword) + "*")
		query.SetField(field)
		queries = append(queries, query)
	}
	var search_query bleve_query.Query = bleve.NewDisjunctionQuery(queries...)
	if set_id != "" {
		set_query := bleve.NewTermQuery(strings.ToLower(set_id))
		set_query.SetField("set_id")
		search_query = bleve.NewConjunctionQuery(search_query, set_query)
	}

	request := bleve.NewSearchRequestOptions(search_query, limit, offset, false)
	request.SortBy([]string{"-date", "_id"})

	bbs_search_index_lock.RLock()
	defer bbs_search_index_lock.RUnlock()
	if !bbs_search_index_ready {
		return []string{}, false
	}
	result, err := bbs_search_index.Search(request)
	if err != nil {
		log.Printf("[SEARCH] bbs search failed: %v", err)
		return []string{}, false
	}

	data_list := make([]string, 0, len(result.Hits))
	for _, hit := range result.Hits {
		data_list = append(data_list, hit.ID)
	}
	return data_list, true
}
