package route

import (
	"database/sql"
	"math"
	"strconv"
	"strings"

	"opennamu/route/tool"
)

const bbs_best_comment_enabled_setting = "bbs_best_comment_enabled"
const bbs_best_comment_min_score_setting = "bbs_best_comment_min_score"
const bbs_best_comment_default_min_score = 5
const bbs_best_comment_max_min_score = 1000000000
const bbs_best_comment_candidate_limit = 50
const bbs_best_comment_limit = 3

func Bbs_best_comment_enabled_normalize(value string) (string, bool) {
	value = strings.TrimSpace(value)
	if value == "" || value == "1" {
		return "1", true
	}
	if value == "0" {
		return "0", true
	}
	return "", false
}

func Bbs_best_comment_min_score_normalize(value string) (string, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return strconv.Itoa(bbs_best_comment_default_min_score), true
	}

	min_score, err := strconv.Atoi(value)
	if err != nil || min_score < 0 || min_score > bbs_best_comment_max_min_score {
		return "", false
	}
	return strconv.Itoa(min_score), true
}

func Bbs_best_comment_enabled(db *sql.DB, set_id string) bool {
	value := strings.TrimSpace(Bbs_set_value(db, set_id, bbs_best_comment_enabled_setting))
	return value != "0"
}

func Bbs_best_comment_min_score(db *sql.DB, set_id string) int {
	value, valid := Bbs_best_comment_min_score_normalize(Bbs_set_value(db, set_id, bbs_best_comment_min_score_setting))
	if !valid {
		return bbs_best_comment_default_min_score
	}
	min_score, err := strconv.Atoi(value)
	if err != nil {
		return bbs_best_comment_default_min_score
	}
	return min_score
}

type bbs_best_comment_candidate struct {
	comment_set_id   string
	comment_set_code string
	net_score        string
}

func Bbs_best_comment_candidate_path_sql(alias string) string {
	if tool.Get_DB_type() == "mysql" {
		return alias + ".set_id regexp '^-?[0-9]+(-[0-9]+)*$' and " + alias + ".set_code regexp '^[1-9][0-9]{0,17}$'"
	}

	return alias + ".set_id not glob '*[^0-9-]*' and " +
		alias + ".set_id not like '%--%' and " +
		alias + ".set_id not like '%-' and " +
		alias + ".set_code not glob '*[^0-9]*' and " +
		"length(" + alias + ".set_code) between 1 and 18 and " +
		"substr(" + alias + ".set_code, 1, 1) between '1' and '9'"
}

func Bbs_comment_vote_counter_sql(column string) string {
	trimmed := "trim(" + column + ")"
	if tool.Get_DB_type() == "mysql" {
		return "case when " + trimmed + " regexp '^[0-9]{1,18}$' then cast(" + trimmed + " as signed) else 0 end"
	}
	return "case when " + trimmed + " <> '' and " + trimmed + " not glob '*[^0-9]*' and length(" + trimmed + ") <= 18 then cast(" + trimmed + " as integer) else 0 end"
}

func Bbs_best_comment_comment_code_valid(comment_code string) bool {
	if !bbs_comment_code_regex.MatchString(comment_code) {
		return false
	}
	for _, part := range strings.Split(comment_code, "-") {
		value, err := strconv.ParseInt(part, 10, 64)
		if err != nil || value <= 0 || strconv.FormatInt(value, 10) != part {
			return false
		}
	}
	return true
}

func Bbs_comment_first_child_code_sql(alias string) string {
	suffix := "substr(" + alias + ".set_id, length(?) + 2)"
	return "substr(" + alias + ".set_id, length(?) + 2, case when instr(" + suffix + ", '-') = 0 then length(" + alias + ".set_id) else instr(" + suffix + ", '-') - 1 end)"
}

func Bbs_comment_node_name_sql(alias string) string {
	return "(" + alias + ".set_name = 'comment' or " + alias + ".set_name like 'comment%' or " + alias + ".set_name in ('blind', 'pinned', 'tabom_count', 'tabom_down_count'))"
}

func Bbs_comment_page_for_code(db *sql.DB, set_id string, set_code string, comment_code string) (string, bool) {
	if !Bbs_best_comment_comment_code_valid(comment_code) || strings.ContainsAny(set_id, "%_\\") || strings.ContainsAny(set_code, "-%_\\") {
		return "", false
	}

	root_id := set_id + "-" + set_code
	parts := strings.Split(comment_code, "-")
	page_rank := int64(len(parts)) // the target and its ancestors occupy traversal slots

	for part_index, target_part := range parts {
		parent_path := root_id
		if part_index > 0 {
			parent_path += "-" + strings.Join(parts[:part_index], "-")
		}

		var ambiguous_code string
		if tool.QueryRow_DB(
			db,
			"select comment_node.set_code from bbs_data comment_node where "+Bbs_comment_node_name_sql("comment_node")+" and comment_node.set_id = ? and comment_node.set_code + 0 = ? and comment_node.set_code <> ? limit 1",
			[]any{&ambiguous_code},
			parent_path,
			target_part,
			target_part,
		) {
			return "", false
		}

		comment_set_code := ""
		if !tool.QueryRow_DB(
			db,
			"select comment_node.set_code from bbs_data comment_node where "+Bbs_comment_node_name_sql("comment_node")+" and comment_node.set_id = ? and comment_node.set_code = ?",
			[]any{&comment_set_code},
			parent_path,
			target_part,
		) {
			return "", false
		}

		first_child_code := Bbs_comment_first_child_code_sql("before_comment")
		parent_first_child_code := Bbs_comment_first_child_code_sql("before_comment")
		query := "select count(*) from (" +
			"select before_comment.set_id, before_comment.set_code from bbs_data before_comment " +
			"where " + Bbs_comment_node_name_sql("before_comment") + " and ((before_comment.set_id = ? and before_comment.set_code + 0 < ?) or (" +
			"before_comment.set_id like ? and (" + first_child_code + ") + 0 < ? and exists (" +
			"select 1 from bbs_data before_parent where " + Bbs_comment_node_name_sql("before_parent") + " and before_parent.set_id = ? and before_parent.set_code = " + parent_first_child_code +
			"))) group by before_comment.set_id, before_comment.set_code) preceding_comments"
		query_values := []any{
			parent_path,
			target_part,
			parent_path + "-%",
			parent_path,
			parent_path,
			parent_path,
			target_part,
			parent_path,
			parent_path,
			parent_path,
			parent_path,
		}

		var preceding_count int64
		if err := db.QueryRow(tool.DB_change(query), query_values...).Scan(&preceding_count); err != nil {
			return "", false
		}
		if preceding_count < 0 || page_rank > math.MaxInt64-preceding_count {
			return "", false
		}
		page_rank += preceding_count
	}

	return strconv.FormatInt((page_rank-1)/bbs_comment_page_size+1, 10), true
}

func Bbs_best_comment_candidates(db *sql.DB, set_id string, set_code string, min_score int, can_view_blind bool) []bbs_best_comment_candidate {
	root_id := set_id + "-" + set_code
	up_count := Bbs_comment_vote_counter_sql("up_vote.set_data")
	down_count := Bbs_comment_vote_counter_sql("down_vote.set_data")
	net_score := "(coalesce(max(" + up_count + "), 0) - coalesce(max(" + down_count + "), 0))"
	blind_filter := ""
	if !can_view_blind {
		blind_filter = " and not exists (select 1 from bbs_data blind_data where blind_data.set_name = 'blind' and blind_data.set_data = 'O' and blind_data.set_id = comment.set_id and blind_data.set_code = comment.set_code)"
	}

	query := "select comment.set_id, comment.set_code, " + net_score + " as net_score " +
		"from bbs_data comment " +
		"left join bbs_data up_vote on up_vote.set_name = 'tabom_count' and up_vote.set_id = comment.set_id and up_vote.set_code = comment.set_code " +
		"left join bbs_data down_vote on down_vote.set_name = 'tabom_down_count' and down_vote.set_id = comment.set_id and down_vote.set_code = comment.set_code " +
		"left join bbs_data comment_date on comment_date.set_name = 'comment_date' and comment_date.set_id = comment.set_id and comment_date.set_code = comment.set_code " +
		"where comment.set_name = 'comment' and trim(comment.set_data) <> '' and (comment.set_id = ? or comment.set_id like ?) and " +
		Bbs_best_comment_candidate_path_sql("comment") + blind_filter +
		" and not exists (select 1 from bbs_data pinned_data where pinned_data.set_name = 'pinned' and pinned_data.set_id = comment.set_id and pinned_data.set_code = comment.set_code) " +
		"group by comment.set_id, comment.set_code having " + net_score + " >= ? " +
		"order by net_score desc, min(comment_date.set_data) asc, comment.set_id asc, comment.set_code asc limit ?"

	rows, err := db.Query(
		tool.DB_change(query),
		root_id,
		root_id+"-%",
		min_score,
		bbs_best_comment_candidate_limit,
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	candidates := make([]bbs_best_comment_candidate, 0, bbs_best_comment_candidate_limit)
	for rows.Next() {
		candidate := bbs_best_comment_candidate{}
		if err := rows.Scan(&candidate.comment_set_id, &candidate.comment_set_code, &candidate.net_score); err != nil {
			panic(err)
		}
		candidates = append(candidates, candidate)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return candidates
}

func Api_bbs_w_comment_best(config tool.Config, set_id string, set_code string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := map[string]any{"response": "ok", "data": []map[string]string{}}
	if set_id == "" || set_code == "" || strings.ContainsAny(set_id, "%_\\") || strings.ContainsAny(set_code, "-%_\\") {
		return_data["response"] = "not exist"
		return return_data
	}
	if !tool.Check_acl(db, set_id, "", "bbs_view", config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}
	if _, allowed := Bbs_post_view_auth(db, set_id, set_code, config.IP); !allowed {
		return_data["response"] = "require auth"
		return return_data
	}
	if set_id == report_bbs_id || !Bbs_best_comment_enabled(db, set_id) {
		return return_data
	}

	can_view_blind := tool.Check_permission(db, "bbs_comment_manage", config.IP)
	candidates := Bbs_best_comment_candidates(db, set_id, set_code, Bbs_best_comment_min_score(db, set_id), can_view_blind)
	best_comments := make([]map[string]string, 0, bbs_best_comment_limit)
	for _, candidate := range candidates {
		post_set_id, post_set_code, comment_code, valid := Bbs_comment_storage_location(candidate.comment_set_id, candidate.comment_set_code)
		if !valid || post_set_id != set_id || post_set_code != set_code || !Bbs_best_comment_comment_code_valid(comment_code) {
			continue
		}

		comment_page, page_exists := Bbs_comment_page_for_code(db, set_id, set_code, comment_code)
		if !page_exists {
			continue
		}

		comment_list := Bbs_comment_one_data(db, config, "", candidate.comment_set_id+"-"+candidate.comment_set_code)
		if len(comment_list) == 0 || strings.TrimSpace(comment_list[0]["comment"]) == "" || (comment_list[0]["blind"] == "O" && !can_view_blind) {
			continue
		}
		comment_list[0]["best_score"] = candidate.net_score
		comment_list[0]["best_page"] = comment_page
		best_comments = append(best_comments, comment_list[0])
		if len(best_comments) == bbs_best_comment_limit {
			break
		}
	}

	return_data["data"] = best_comments
	return return_data
}
