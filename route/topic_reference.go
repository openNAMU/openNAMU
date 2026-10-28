package route

import (
	"database/sql"
	"regexp"
	"strings"

	"opennamu/route/tool"
)

var topic_reference_regex = regexp.MustCompile(`#([0-9]+)(?:-([0-9]+))?`)
var topic_call_regex = regexp.MustCompile(`@([^ \n]+)`)
var topic_legacy_thread_regex = regexp.MustCompile(`(?:&lt;|<)topic_a(?:_thread)?(?:&gt;|>)#([0-9]+)(?:&lt;|<)/topic_a(?:_thread)?(?:&gt;|>)`)

func Replace_reference_or_call_tokens(data string, token_re *regexp.Regexp, callback func(value string) string) string {
	var out strings.Builder
	pos := 0

	for {
		idx := token_re.FindStringSubmatchIndex(data[pos:])
		if idx == nil {
			break
		}

		matchStart := pos + idx[0]
		matchEnd := pos + idx[1]

		prev_boundary := matchStart == 0 || data[matchStart-1] == ' ' || data[matchStart-1] == '\n'
		next_boundary := matchEnd == len(data) || data[matchEnd] == ' ' || data[matchEnd] == '\n'
		if !prev_boundary || !next_boundary {
			out.WriteString(data[pos:matchEnd])
			pos = matchEnd
			continue
		}

		out.WriteString(data[pos:matchStart])
		out.WriteString(callback(data[matchStart:matchEnd]))
		pos = matchEnd
	}

	out.WriteString(data[pos:])
	return out.String()
}

func Topic_reference_notify(db *sql.DB, config tool.Config, data string, num string, topic_num string, set_id string, name string, sub string, do_type string) {
	if data == "" {
		return
	}

	sent := map[string]bool{}
	send := func(target string) {
		if target == "" || tool.IP_or_user(target) || target == config.IP || sent[target] {
			return
		}

		sent[target] = true
		path := "/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(topic_num) + "#" + tool.Url_parser(num)
		tool.Send_alarm(
			db,
			config.IP,
			target,
			`<a href="`+path+`">`+tool.HTML_escape(name)+` - `+tool.HTML_escape(sub)+`#`+tool.HTML_escape(num)+`</a>`,
		)
	}

	Replace_reference_or_call_tokens(data, topic_reference_regex, func(value string) string {
		match := topic_reference_regex.FindStringSubmatch(value)
		if match == nil {
			return value
		}

		reference_id := match[1]
		reference_code := topic_num
		if match[2] != "" {
			reference_code = match[2]
		}

		target := ""
		if do_type == "thread" {
			if reference_id == "0" {
				tool.QueryRow_DB(db, "select set_data from bbs_data where set_name = 'user_id' and set_id = ? and set_code = ?", []any{&target}, thread_bbs_id, reference_code)
			} else {
				tool.QueryRow_DB(db, "select set_data from bbs_data where set_name = 'comment_user_id' and set_id = ? and set_code = ?", []any{&target}, thread_bbs_id+"-"+reference_code, reference_id)
			}
		} else if reference_id == "0" {
			tool.QueryRow_DB(db, "select set_data from bbs_data where set_name = 'user_id' and set_id = ? and set_code = ?", []any{&target}, set_id, topic_num)
		} else {
			comment_set_id := set_id + "-" + topic_num
			if match[2] != "" {
				comment_set_id = set_id + "-" + match[2]
			}
			tool.QueryRow_DB(db, "select set_data from bbs_data where set_name = 'comment_user_id' and set_id = ? and set_code = ?", []any{&target}, comment_set_id, reference_id)
		}

		send(target)
		return value
	})

	Replace_reference_or_call_tokens(data, topic_call_regex, func(value string) string {
		match := topic_call_regex.FindStringSubmatch(value)
		if match == nil {
			return value
		}

		target := ""
		tool.QueryRow_DB(db, "select ip from history where ip = ? limit 1", []any{&target}, match[1])
		if target == "" {
			tool.QueryRow_DB(db, "select set_data from bbs_data where set_name = 'comment_user_id' and set_data = ? and set_id like '-1-%' limit 1", []any{&target}, match[1])
		}
		send(target)
		return value
	})
}

func Render_topic_reference(data string, topic_num string, set_id string, set_code string, do_type string) string {
	data = topic_legacy_thread_regex.ReplaceAllStringFunc(data, func(value string) string {
		match := topic_legacy_thread_regex.FindStringSubmatch(value)
		return `<a href="#` + tool.Url_parser(match[1]) + `">#` + tool.HTML_escape(match[1]) + `</a>`
	})

	data = Replace_reference_or_call_tokens(data, topic_reference_regex, func(value string) string {
		match := topic_reference_regex.FindStringSubmatch(value)
		if match == nil {
			return value
		}

		label := "#" + match[1]
		reference_code := topic_num
		if match[2] != "" {
			label += "-" + match[2]
			if do_type == "thread" {
				reference_code = match[2]
			}
		}

		path := ""
		if do_type == "thread" && reference_code != "" {
			path = "/bbs/w/" + tool.Url_parser(thread_bbs_id) + "/" + tool.Url_parser(reference_code) + "#" + tool.Url_parser(match[1])
		} else if do_type != "thread" && set_id != "" && set_code != "" {
			target_code := set_code
			anchor := match[1]
			if match[2] != "" {
				target_code = match[2]
			}
			path = "/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(target_code) + "#" + tool.Url_parser(anchor)
		}
		if path == "" {
			return value
		}

		return `<a href="` + path + `">` + label + `</a>`
	})

	return Replace_reference_or_call_tokens(data, topic_call_regex, func(value string) string {
		match := topic_call_regex.FindStringSubmatch(value)
		if match == nil {
			return value
		}

		label := "@" + match[1]
		return `<a href="/w/user:` + tool.Url_parser(match[1]) + `">` + tool.HTML_escape(label) + `</a>`
	})
}
