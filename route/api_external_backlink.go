package route

import (
	"database/sql"
	"sort"
	"strconv"
	"strings"

	"opennamu/route/tool"
	backlink_parser "opennamu/route/tool/markup/backlink"
)

// external_doc_in_batch bounds the placeholder count of the existing-doc IN
// query so a large doc set stays within SQLite/MySQL statement limits.
const external_doc_in_batch = 500

// Get_existing_doc_names returns the subset of names that exist as rows in the
// data table, checking existence with bounded IN batches.
func Get_existing_doc_names(db *sql.DB, names []string) map[string]bool {
	existing := make(map[string]bool, len(names))
	for start := 0; start < len(names); start += external_doc_in_batch {
		end := start + external_doc_in_batch
		if end > len(names) {
			end = len(names)
		}
		chunk := names[start:end]
		placeholders := make([]string, 0, len(chunk))
		args := make([]any, 0, len(chunk))
		for _, name := range chunk {
			placeholders = append(placeholders, "?")
			args = append(args, name)
		}
		rows := tool.Query_DB(
			db,
			"select title from data where title in ("+strings.Join(placeholders, ",")+")",
			args...,
		)
		defer rows.Close()
		for rows.Next() {
			title := ""
			if err := rows.Scan(&title); err != nil {
				continue
			}
			existing[title] = true
		}
	}
	return existing
}

func Viewable_external_doc_names(db *sql.DB, config tool.Config, normalized_host string) []string {
	candidates := tool.Query_DB(db, "select link from back where type = 'external' and title = ?", normalized_host)
	defer candidates.Close()

	doc_set := map[string]bool{}
	for candidates.Next() {
		link := ""
		if err := candidates.Scan(&link); err != nil {
			panic(err)
		}
		if link == "" {
			continue
		}
		doc_set[link] = true
	}

	doc_names := make([]string, 0, len(doc_set))
	for link := range doc_set {
		doc_names = append(doc_names, link)
	}
	if len(doc_names) == 0 {
		return doc_names
	}

	existing := Get_existing_doc_names(db, doc_names)

	result := make([]string, 0, len(doc_names))
	for _, name := range doc_names {
		if !existing[name] {
			continue
		}
		if !tool.Check_acl(db, name, "", "render", config.IP) {
			continue
		}
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}

func Api_external_docs_for_host(config tool.Config, page_str string, host string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)

	normalized_host := backlink_parser.Normalize_external_host(strings.TrimSpace(host))
	if normalized_host == "" {
		return_data["response"] = "ok"
		return_data["data"] = [][]string{}
		return return_data
	}

	doc_names := Viewable_external_doc_names(db, config, normalized_host)
	if len(doc_names) == 0 {
		return_data["response"] = "ok"
		return_data["data"] = [][]string{}
		return return_data
	}

	exact_rows := make([][]string, 0, len(doc_names))
	for _, name := range doc_names {
		exact_rows = append(exact_rows, []string{name, normalized_host})
	}

	return_data["response"] = "ok"
	return_data["data"] = Paginate_external_backlink(page_str, exact_rows)
	return return_data
}

func Api_external_hosts_of_doc(config tool.Config, doc_name string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)

	if !tool.Check_acl(db, doc_name, "", "render", config.IP) {
		return_data["response"] = "require auth"
		return_data["data"] = []string{}
		return return_data
	}

	candidates := tool.Query_DB(
		db,
		"select link, title from back where type = 'external' and "+tool.Exact_binary_eq("link"),
		doc_name,
	)
	defer candidates.Close()

	host_set := map[string]bool{}
	var hosts []string
	for candidates.Next() {
		link := ""
		title := ""
		if err := candidates.Scan(&link, &title); err != nil {
			panic(err)
		}
		if link != doc_name {
			continue
		}
		if title == "" {
			continue
		}
		if host_set[title] {
			continue
		}
		host_set[title] = true
		hosts = append(hosts, title)
	}

	sort.Strings(hosts)

	return_data["response"] = "ok"
	return_data["data"] = hosts
	return return_data
}

func Api_external_count_for_host(config tool.Config, host string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)

	normalized_host := backlink_parser.Normalize_external_host(strings.TrimSpace(host))
	if normalized_host == "" {
		return_data["response"] = "ok"
		return_data["count"] = 0
		return return_data
	}

	doc_names := Viewable_external_doc_names(db, config, normalized_host)
	return_data["response"] = "ok"
	return_data["count"] = len(doc_names)
	return return_data
}

func Api_external_all_hosts(config tool.Config, page_str string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !tool.Check_permission(db, "setting_backlink", config.IP) {
		return_data["response"] = "require auth"
		return_data["data"] = [][]string{}
		return return_data
	}

	// Count per host against the set of docs that actually exist in `data`,
	// case-sensitively. This excludes deleted docs and keeps case-sibling
	// names (Page / page) distinct, matching the per-host detailed view.
	candidates := tool.Query_DB(
		db,
		"select title, link from back where type = 'external' and title != '' and link != ''",
	)
	defer candidates.Close()

	host_links := map[string]map[string]bool{}
	all_links := map[string]bool{}
	for candidates.Next() {
		title := ""
		link := ""
		if err := candidates.Scan(&title, &link); err != nil {
			panic(err)
		}
		if title == "" || link == "" {
			continue
		}
		if host_links[title] == nil {
			host_links[title] = make(map[string]bool)
		}
		host_links[title][link] = true
		all_links[link] = true
	}

	all_link_names := make([]string, 0, len(all_links))
	for name := range all_links {
		all_link_names = append(all_link_names, name)
	}
	existing := Get_existing_doc_names(db, all_link_names)

	host_counts := map[string]int{}
	for host, links := range host_links {
		count := 0
		for link := range links {
			if existing[link] {
				count++
			}
		}
		host_counts[host] = count
	}

	hosts := make([]string, 0, len(host_counts))
	for host := range host_counts {
		hosts = append(hosts, host)
	}
	sort.Strings(hosts)

	exact_rows := make([][]string, 0, len(hosts))
	for _, host := range hosts {
		exact_rows = append(exact_rows, []string{host, strconv.Itoa(host_counts[host])})
	}

	return_data["response"] = "ok"
	return_data["data"] = Paginate_external_backlink(page_str, exact_rows)
	return return_data
}

func Paginate_external_backlink(num_str string, exact_rows [][]string) [][]string {
	page := tool.Str_to_int(num_str)
	num := 0
	if page*50 > 0 {
		num = page*50 - 50
	}
	if num >= len(exact_rows) {
		return [][]string{}
	}
	end := num + 50
	if end > len(exact_rows) {
		end = len(exact_rows)
	}
	return exact_rows[num:end]
}
