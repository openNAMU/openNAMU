package tool

import (
	"database/sql"
	"math"
	"strings"
)

type UA_multiple_data struct {
	IP    string
	Count int
	Date  string
	Names []string
}

const multiple_page_size = 50
const multiple_batch_size = 256

// Get_ua_multiple_rows returns multi-account groups that share an IP. A group is
// excluded when any of its members holds hideip, and it walks the matching
// groups in bounded batches (multiple_batch_size at a time) so the visible-row
// offset can be applied after Go-side permission filtering without loading the
// whole result set.
//
// A page past the end (offset == math.MaxInt as returned by Page_offset) yields
// an empty page without querying, so an out-of-range page number cannot
// trigger a full walk of the table.
func Get_ua_multiple_rows(db *sql.DB, search string, sort string, offset int) []UA_multiple_data {
	if offset == math.MaxInt {
		return nil
	}

	base_query := "select base.ip, count(distinct base.name), max(base.today) from ua_d as base where base.ip != '' and base.name != ''"
	search_values := []any{}
	if search != "" {
		base_query += " and (base.ip like ? or exists (select 1 from ua_d as matched where matched.ip = base.ip and matched.name like ?))"
		search_value := "%" + search + "%"
		search_values = append(search_values, search_value, search_value)
	}
	base_query += " group by base.ip having count(distinct base.name) >= 2"
	if sort == "count" {
		base_query += " order by count(distinct base.name) desc, max(base.today) desc, base.ip asc"
	} else {
		base_query += " order by max(base.today) desc, count(distinct base.name) desc, base.ip asc"
	}

	var collected []UA_multiple_data
	var visible_count int
	hideip_cache := make(map[string]*bool)

	for raw := 0; ; raw += multiple_batch_size {
		batch := Get_group_batch(db, base_query, search_values, raw, multiple_batch_size)
		if len(batch) == 0 {
			break
		}
		Attach_names(db, batch)

		// Populate the effective hideip result for every member read in this
		// batch.
		for _, group := range batch {
			for _, name := range group.Names {
				if _, done := hideip_cache[name]; done {
					continue
				}
				res := false
				if !IP_or_user(name) {
					res = Check_permission(db, "hideip", name)
				}
				hideip_cache[name] = &res
			}
		}

		// Fail-closed: a group is shown only if its full member set was read
		// (len(Names) == Count) AND none of its members holds hideip. A scan
		// failure or genuine DB error that leaves a member unread yields
		// len(Names) < Count, and the group is excluded rather than exposed
		// as if it were safe.
		for _, group := range batch {
			excluded := len(group.Names) != group.Count
			for _, name := range group.Names {
				if hideip_cache[name] != nil && *hideip_cache[name] {
					excluded = true
					break
				}
			}
			if excluded {
				continue
			}
			if visible_count < offset {
				visible_count++
				continue
			}
			collected = append(collected, group)
			visible_count++
			if len(collected) == multiple_page_size {
				return collected
			}
		}
	}
	return collected
}

// Get_group_batch returns one bounded batch of grouped IP rows. max(base.today)
// is nullable, so Date is read with sql.NullString to avoid a scan failure
// dropping a whole group.
func Get_group_batch(db *sql.DB, base_query string, search_values []any, offset, limit int) []UA_multiple_data {
	query := base_query + " limit ?, ?"
	args := append([]any{}, search_values...)
	args = append(args, offset, limit)

	rows := Query_DB(db, query, args...)
	defer rows.Close()

	data := []UA_multiple_data{}
	for rows.Next() {
		var g UA_multiple_data
		var date_s sql.NullString
		if err := rows.Scan(&g.IP, &g.Count, &date_s); err != nil {
			// Unexpected scan failure (NULLs are handled by NullString). Aborting
			// (panic) is fail-closed: the request becomes a 500 via the central
			// recover, so no partial group data is ever returned to a caller.
			panic(err)
		}
		g.Date = date_s.String
		data = append(data, g)
	}
	// A row-level error that was not surfaced by Scan (e.g. during iteration)
	// must not be silently turned into "no more rows".
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return data
}

// Attach_names fills each group's member names from the bounded IP batch. The
// ip, name and max(today) columns are read with sql.NullString so an expected
// NULL (notably today) does not drop a member; a dropped hideip member would
// otherwise let its group be exposed as safe.
func Attach_names(db *sql.DB, groups []UA_multiple_data) {
	if len(groups) == 0 {
		return
	}
	placeholders := make([]string, len(groups))
	args := make([]any, len(groups))
	for i, g := range groups {
		placeholders[i] = "?"
		args[i] = g.IP
	}
	name_query := "select ip, name, max(today) from ua_d where ip in (" + strings.Join(placeholders, ",") + ") and name != '' group by ip, name order by max(today) desc, name asc"
	name_rows := Query_DB(db, name_query, args...)
	defer name_rows.Close()

	by_ip := make(map[string]int)
	for i, g := range groups {
		by_ip[g.IP] = i
	}
	for name_rows.Next() {
		var ip_s sql.NullString
		var name_s sql.NullString
		var today_s sql.NullString
		if err := name_rows.Scan(&ip_s, &name_s, &today_s); err != nil {
			// Unexpected scan failure (NULLs are handled by NullString). If a
			// member row cannot be read, we must not treat the member as absent:
			// a missing hideip member would wrongly let its IP group be exposed
			// as safe. Aborting is fail-closed.
			panic(err)
		}
		if !name_s.Valid || name_s.String == "" {
			continue
		}
		idx, ok := by_ip[ip_s.String]
		if ok {
			groups[idx].Names = append(groups[idx].Names, name_s.String)
		}
	}
	// A row-level error that was not surfaced by Scan (e.g. during iteration)
	// must not be silently turned into "no more rows".
	if err := name_rows.Err(); err != nil {
		panic(err)
	}
}
