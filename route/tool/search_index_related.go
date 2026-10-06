package tool

import (
	"errors"
	"strings"

	"github.com/blevesearch/bleve/v2"
)

// Related-document search parameters (§5·§10·§17; adjust in one place).
const (
	search_related_candidate_limit      = 20   // §10 Bleve candidate limit (TOP N request size)
	Search_related_result_limit         = 5    // §16 exported: route slices result to this max
	search_related_title_boost          = 3.0  // §5 title field weight
	search_related_body_boost           = 1.0  // §5 body field weight
	search_related_query_body_max_runes = 6000 // §8 UTF-8-safe query-body prefix
	Search_related_relative_threshold   = 0.35 // §17 exported: route threshold = best_score * this
)

// Search_related_hit : related-document search candidate. ID is the real document identifier, Score is Bleve relevance score.
type Search_related_hit struct {
	ID    string
	Score float64
}

// Search_index_related : content-based related-document search that reuses the existing shared index.
//
// It owns the readiness check and the RLock (lock ownership stays in tool per §4), so route never needs
// to know about a raw bleve.Index or tool.Search_index(). Missing input returns no results; an unavailable index
// and search errors are returned to the caller as processing failures.
func Search_index_related(doc_name string, raw_data string) ([]Search_related_hit, error) {
	if doc_name == "" {
		return []Search_related_hit{}, nil
	}
	if !Search_index_ready() {
		return nil, errors.New("search index is not ready")
	}

	search_index_lock.RLock()
	defer search_index_lock.RUnlock()

	return Search_index_related_with_index(search_index, doc_name, raw_data)
}

// Search_index_related_with_index : pure-Bleve relevance core. It takes the index explicitly so unit tests
// can inject a temporary index (§10), and it does no DB/ACL work (those live in route). Self document is
// excluded by ID only — normalization is intentionally not applied here (§11·§31).
func Search_index_related_with_index(index bleve.Index, doc_name string, raw_data string) ([]Search_related_hit, error) {
	if index == nil || doc_name == "" {
		return []Search_related_hit{}, nil
	}

	query_body := strings.ReplaceAll(strings.ReplaceAll(raw_data, "\r", ""), "\n", " ")
	query_body = Get_slice(query_body, 0, search_related_query_body_max_runes)

	title_query := bleve.NewMatchQuery(Do_remove_spaces(doc_name))
	title_query.SetField("title_search")
	title_query.SetBoost(search_related_title_boost)

	body_query := bleve.NewMatchQuery(query_body)
	body_query.SetField("data")
	body_query.SetBoost(search_related_body_boost)

	query := bleve.NewDisjunctionQuery(title_query, body_query)
	request := bleve.NewSearchRequestOptions(query, search_related_candidate_limit, 0, false)

	result, err := index.Search(request)
	if err != nil {
		return nil, err
	}

	related_hits := make([]Search_related_hit, 0, len(result.Hits))
	for _, hit := range result.Hits {
		if hit.ID == doc_name {
			continue
		}
		related_hits = append(related_hits, Search_related_hit{ID: hit.ID, Score: hit.Score})
	}

	return related_hits, nil
}
