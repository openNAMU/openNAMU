package route

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"opennamu/route/tool"
)

// Related-document candidate policy lives in the API layer (§5): it is candidate filtering, not a
// general Bleve search concern. Prefix constants are defined here in one place alongside their only use.
const (
	related_namespace_prefix_user     = "user:"
	related_namespace_prefix_file     = "file:"
	related_namespace_prefix_category = "category:"
	related_document_enabled_setting  = "related_document_enabled"
)

func Related_document_enabled(db *sql.DB) bool {
	return tool.Get_setting_value_exact(db, related_document_enabled_setting, "", "1") != "0"
}

func Related_is_namespace(name string) bool {
	return strings.HasPrefix(name, related_namespace_prefix_user) ||
		strings.HasPrefix(name, related_namespace_prefix_file) ||
		strings.HasPrefix(name, related_namespace_prefix_category)
}

// Query_related_redirects reports which of the given candidate ids are redirect sources. It builds a single
// `link in (?, ?, ...)` query over exactly len(ids) placeholders and runs at most one statement, so it never
// loads the whole back table per page (§2). Empty id lists short-circuit before any SQL is built, avoiding a
// degenerate IN clause on both SQLite and MySQL. Candidate ids are internal document identifiers coming from
// Bleve hits (already stored in the DB), not user input, so positional binding keeps them injection-safe (§3).
func Query_related_redirects(db tool.DB_runner, candidate_ids []string) (map[string]bool, error) {
	redirect_set := map[string]bool{}
	if len(candidate_ids) == 0 {
		return redirect_set, nil
	}

	placeholder := ""
	for range candidate_ids {
		if placeholder != "" {
			placeholder += ", "
		}
		placeholder += "?"
	}

	id_args := make([]any, len(candidate_ids))
	for i, id := range candidate_ids {
		id_args[i] = id
	}

	query := "select link from back where type = 'redirect' and link in (" + placeholder + ")"
	stmt, err := db.Prepare(tool.DB_change(query))
	if err != nil {
		return nil, fmt.Errorf("prepare redirect lookup: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(id_args...)
	if err != nil {
		return nil, fmt.Errorf("query redirect lookup: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var link string
		if err := rows.Scan(&link); err != nil {
			return nil, fmt.Errorf("scan redirect lookup: %w", err)
		}
		redirect_set[link] = true
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate redirect lookup: %w", err)
	}

	return redirect_set, nil
}

func Api_w_related(config tool.Config, doc_name string, raw_data string) (return_data map[string]any) {
	defer func() {
		if recovered := recover(); recovered != nil {
			return_data = Related_failure_result("related document processing", fmt.Errorf("%v", recovered))
		}
	}()

	if doc_name == "" {
		return map[string]any{"response": "error", "data": "invalid document"}
	}

	db := tool.DB_connect() // opened once; reused by the redirect-source check below and by later ACL/redirect queries (§6)
	defer tool.DB_close(db)
	if !Related_document_enabled(db) {
		return Empty_related_result()
	}

	// Enforce the source ACL here so direct API callers do not rely on View_w or Api_w_raw to authorize the document.
	if !tool.Check_acl(db, doc_name, "", "render", config.IP) {
		return map[string]any{"response": "require auth"}
	}

	// These are successful no-result cases, not processing failures.
	if raw_data == "" || Related_is_namespace(doc_name) {
		return Empty_related_result()
	}

	// §1 explicit redirect-source exclusion: a document that is (or looks like) a /w_from redirect source must not get
	// related search even though status==200 && raw_data != "" there. Uses only the `link` column to avoid Get_back_redirect's
	// title/target conflation (§1). View_w already redirects plain `/w/<source>` before this code, but /w_from/ reaches it.
	redirect_source_set, err := Query_related_redirects(db, []string{doc_name})
	if err != nil {
		return Related_failure_result("redirect source lookup", err)
	}
	if redirect_source_set[doc_name] {
		return Empty_related_result()
	}

	candidates, err := tool.Search_index_related(doc_name, raw_data) // readiness + lock handled internally (§4)
	if err != nil {
		return Related_failure_result("Bleve search", err)
	}

	related_ids := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		related_ids = append(related_ids, candidate.ID)
	}

	redirect_set, err := Query_related_redirects(db, related_ids) // §6 redirect-source exclusion (§2 batch query)
	if err != nil {
		return Related_failure_result("redirect candidate lookup", err)
	}

	result := Related_candidate_names(candidates, redirect_set, func(candidate_name string) bool {
		return tool.Check_acl(db, candidate_name, "", "render", config.IP)
	})

	return_data = make(map[string]any)
	return_data["response"] = "ok"
	return_data["data"] = result
	return return_data
}

// Empty_related_result : consistent success shape for excluded sources — same response/data keys as the happy path,
// so every caller validates uniformly. §16 never returns an unkeyed empty map.
func Empty_related_result() map[string]any {
	return_data := make(map[string]any)
	return_data["response"] = "ok"
	return_data["data"] = []string{}
	return return_data
}

func Related_failure_result(operation string, err error) map[string]any {
	log.Printf("[RELATED] %s failed: %v", operation, err)
	return map[string]any{"response": "error"}
}

func Related_candidate_names(candidates []tool.Search_related_hit, redirect_set map[string]bool, can_render func(string) bool) []string {
	valid_candidates := make([]tool.Search_related_hit, 0, len(candidates))
	best_score := 0.0

	for _, candidate := range candidates {
		if Related_is_namespace(candidate.ID) || redirect_set[candidate.ID] || !can_render(candidate.ID) {
			continue
		}
		valid_candidates = append(valid_candidates, candidate)
		if candidate.Score > best_score {
			best_score = candidate.Score
		}
	}

	result := make([]string, 0, len(valid_candidates))
	if best_score > 0 {
		threshold := tool.Search_related_relative_threshold * best_score
		for _, candidate := range valid_candidates {
			if candidate.Score >= threshold {
				result = append(result, candidate.ID)
				if len(result) == tool.Search_related_result_limit {
					break
				}
			}
		}
	}

	return result
}
