package route

import (
	"database/sql"
	"math"
	"strings"

	"opennamu/route/tool"
)

const ua_page_size = 50
const ua_batch_size = 256

// Page_offset returns the raw row offset for a 1-based page number. It clamps
// the page so that (page_num-1)*ua_page_size never overflows; when it would,
// it returns math.MaxInt so the caller treats the page as past-the-end and
// yields an empty page instead of panicking.
func Page_offset(page_num int) int {
	if page_num < 1 {
		page_num = 1
	}
	if page_num-1 > math.MaxInt/ua_page_size {
		return math.MaxInt
	}
	return (page_num - 1) * ua_page_size
}

// Api_user_check_rows returns the user-check rows for the inspector, subject,
// mode and page. It enforces the check permission and:
//   - a name lookup masks a hideip subject's ip/ua (existing behavior);
//   - an IP reverse lookup excludes hideip accounts that were not explicitly
//     named on either side of the compare.
func Api_user_check_rows(
	db *sql.DB,
	inspector_ip string,
	name string,
	check_type string,
	page_num int,
	plus_name string,
) map[string]any {
	if !tool.Check_permission(db, "check", inspector_ip) {
		return map[string]any{"response": "require auth"}
	}

	name = strings.TrimSpace(name)
	plus_name = strings.TrimSpace(plus_name)
	if page_num < 1 {
		page_num = 1
	}
	offset := Page_offset(page_num)
	name_is_ip := tool.IP_or_user(name)
	plus_is_ip := tool.IP_or_user(plus_name)

	if check_type == "simple" {
		return Simple_rows(db, name, name_is_ip, offset)
	}
	return Normal_rows(db, name, plus_name, name_is_ip, plus_is_ip, offset)
}

// Hideip_names returns the set of user names that hold the hideip authority.
// It uses Check_permission so that group, rule, regex, cidr and reserved-name
// rules are honored. IP-like strings and empty names are skipped.
func Hideip_names(db *sql.DB, names []string) map[string]bool {
	hideip := make(map[string]bool)
	for _, n := range names {
		if n == "" || tool.IP_or_user(n) {
			continue
		}
		hideip[n] = tool.Check_permission(db, "hideip", n)
	}
	return hideip
}

// Simple_rows returns the "check" rows for one subject. When the subject is an
// IP it reverse-lookups names and excludes hideip accounts; when the subject
// is a name it returns the subject's distinct IPs, masked when the subject
// holds hideip.
func Simple_rows(db *sql.DB, name string, name_is_ip bool, offset int) map[string]any {
	maps := []map[string]string{}

	// A page past the end (offset == math.MaxInt as returned by Page_offset)
	// yields an empty page without querying the DB, so an out-of-range page
	// number cannot trigger a batch walk over the table.
	if offset == math.MaxInt {
		return Simple_response(nil)
	}

	if name_is_ip {
		var visible int
		for raw := 0; ; raw += ua_batch_size {
			batch := tool.Get_ua_simple_rows(db, name, true, raw, ua_batch_size)
			hideip := Hideip_names(db, batch)
			for _, other_name := range batch {
				if hideip[other_name] {
					continue
				}
				if visible < offset {
					visible++
					continue
				}
				maps = append(maps, map[string]string{"other_name": other_name, "masked": "false"})
				visible++
				if len(maps) == ua_page_size {
					return Simple_response(maps)
				}
			}
			if len(batch) < ua_batch_size {
				return Simple_response(maps)
			}
		}
	}

	masked := tool.Check_permission(db, "hideip", name)
	raw := tool.Get_ua_simple_rows(db, name, false, offset, ua_page_size)
	for _, other_name := range raw {
		if masked {
			maps = append(maps, map[string]string{"other_name": "", "masked": "true"})
		} else {
			maps = append(maps, map[string]string{"other_name": other_name, "masked": "false"})
		}
	}
	return Simple_response(maps)
}

// Normal_rows returns the compare rows for two subjects. A row is excluded when
// it belongs to a hideip account and was reached only through an IP lookup (not
// an explicitly named side); hideip rows explicitly named on either side are
// shown with ip/ua masked.
func Normal_rows(db *sql.DB, name string, plus_name string, name_is_ip bool, plus_is_ip bool, offset int) map[string]any {
	maps := []map[string]string{}

	// A page past the end (offset == math.MaxInt as returned by Page_offset)
	// yields an empty page without querying the DB, so an out-of-range page
	// number cannot trigger a batch walk over the table.
	if offset == math.MaxInt {
		return Normal_response(db, nil, name, plus_name, name_is_ip, plus_is_ip, offset)
	}

	if !name_is_ip && !plus_is_ip {
		// No IP involved: every returned row is an explicit name match, so no
		// row is excluded. One bounded SQL batch suffices.
		raw := tool.Get_ua_rows(db, name, plus_name, false, false, offset, ua_page_size)
		hideip := Hideip_names(db, Row_names(raw))
		for _, row := range raw {
			maps = append(maps, Normal_row_map(row, hideip[row["name"]]))
		}
	} else {
		var visible int
		for raw := 0; ; raw += ua_batch_size {
			batch := tool.Get_ua_rows(db, name, plus_name, name_is_ip, plus_is_ip, raw, ua_batch_size)
			hideip := Hideip_names(db, Row_names(batch))
			for _, row := range batch {
				protect := hideip[row["name"]]
				explicit := (row["name"] == name && !name_is_ip) ||
					(row["name"] == plus_name && plus_name != "" && !plus_is_ip)
				if protect && !explicit {
					continue
				}
				if visible < offset {
					visible++
					continue
				}
				maps = append(maps, Normal_row_map(row, protect))
				visible++
				if len(maps) == ua_page_size {
					return Normal_response(db, maps, name, plus_name, name_is_ip, plus_is_ip, offset)
				}
			}
			if len(batch) < ua_batch_size {
				return Normal_response(db, maps, name, plus_name, name_is_ip, plus_is_ip, offset)
			}
		}
	}
	return Normal_response(db, maps, name, plus_name, name_is_ip, plus_is_ip, offset)
}

// Row_names collects the subject names of a compare-row batch.
func Row_names(rows []map[string]string) []string {
	names := make([]string, 0, len(rows))
	for _, row := range rows {
		names = append(names, row["name"])
	}
	return names
}

// Normal_row_map builds the output map for one compare row, masking ip/ua when
// the subject holds hideip.
func Normal_row_map(row map[string]string, protect bool) map[string]string {
	m := map[string]string{
		"name":   row["name"],
		"ip":     row["ip"],
		"ua":     row["ua"],
		"today":  row["today"],
		"masked": "false",
	}
	if protect {
		m["ip"] = ""
		m["ua"] = ""
		m["masked"] = "true"
	}
	return m
}

func Simple_response(maps []map[string]string) map[string]any {
	return map[string]any{
		"response":  "success",
		"simple":    true,
		"rows":      maps,
		"row_count": len(maps),
	}
}

func Normal_response(db *sql.DB, maps []map[string]string, name string, plus_name string, name_is_ip bool, plus_is_ip bool, offset int) map[string]any {
	page_num := offset/ua_page_size + 1

	same_ip_exist := false
	if plus_name != "" && page_num == 1 &&
		!name_is_ip && !plus_is_ip &&
		!tool.Check_permission(db, "hideip", name) && !tool.Check_permission(db, "hideip", plus_name) {
		all_count := len(tool.Get_ua_distinct_ips(db, name, plus_name, name_is_ip, plus_is_ip))
		name_count := len(tool.Get_ua_distinct_ips(db, name, "", name_is_ip, false))
		plus_count := len(tool.Get_ua_distinct_ips(db, plus_name, "", plus_is_ip, false))
		if name_count+plus_count != all_count {
			same_ip_exist = true
		}
	}

	qa_question := ""
	qa_answer := ""
	show_qa := false
	if !name_is_ip && !tool.Check_permission(db, "hideip", name) {
		qa_question, qa_answer = tool.Get_user_approval_data(db, name)
		if qa_question != "" && qa_answer != "" {
			show_qa = true
		}
	}

	return map[string]any{
		"response":      "success",
		"simple":        false,
		"rows":          maps,
		"row_count":     len(maps),
		"same_ip_exist": same_ip_exist,
		"show_qa":       show_qa,
		"qa_question":   qa_question,
		"qa_answer":     qa_answer,
	}
}
