package route

import (
	"database/sql"
	"regexp"
	"strings"

	"opennamu/route/tool"
)

var topic_reference_regex = regexp.MustCompile(`(^|[ \n])#([0-9]+)(?:-([0-9]+))?($|[ \n])`)
var topic_call_regex = regexp.MustCompile(`(^|[ \n])@([^ \n]+)($|[ \n])`)

func topic_reference_notify(db *sql.DB, config tool.Config, data string, num string, topic_num string, set_id string, name string, sub string, do_type string) {
	if data == "" {
		return
	}

	topic_reference_regex.ReplaceAllStringFunc(data, func(value string) string {
		match := topic_reference_regex.FindStringSubmatch(value)
		if len(match) < 5 {
			return value
		}

		reference_id := match[2]
		reference_code := topic_num
		if match[3] != "" {
			reference_code = match[3]
		}

		target := ""
		if do_type == "thread" {
			tool.QueryRow_DB(db, "select ip from topic where code = ? and id = ?", []any{&target}, reference_code, reference_id)
		} else if reference_id == "0" {
			tool.QueryRow_DB(db, "select set_data from bbs_data where set_name = 'user_id' and set_id = ? and set_code = ?", []any{&target}, set_id, topic_num)
		} else {
			comment_set_id := set_id + "-" + topic_num
			if match[3] != "" {
				comment_set_id = set_id + "-" + match[3]
			}
			tool.QueryRow_DB(db, "select set_data from bbs_data where set_name = 'comment_user_id' and set_id = ? and set_code = ?", []any{&target}, comment_set_id, reference_id)
		}

		if target != "" && !tool.IP_or_user(target) {
			path := "/thread/" + tool.Url_parser(topic_num) + "#" + tool.Url_parser(num)
			if do_type != "thread" {
				path = "/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(topic_num) + "#" + tool.Url_parser(num)
			}
			tool.Send_alarm(db, config.IP, target, `<a href="`+path+`">`+tool.HTML_escape(name)+` - `+tool.HTML_escape(sub)+`#`+tool.HTML_escape(num)+`</a>`)
		}
		return value
	})

	topic_call_regex.ReplaceAllStringFunc(data, func(value string) string {
		match := topic_call_regex.FindStringSubmatch(value)
		if len(match) < 4 {
			return value
		}

		target := ""
		tool.QueryRow_DB(db, "select ip from history where ip = ? limit 1", []any{&target}, match[2])
		if target == "" {
			tool.QueryRow_DB(db, "select ip from topic where ip = ? limit 1", []any{&target}, match[2])
		}
		if target != "" && !tool.IP_or_user(target) {
			path := "/thread/" + tool.Url_parser(topic_num) + "#" + tool.Url_parser(num)
			if do_type != "thread" {
				path = "/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(topic_num) + "#" + tool.Url_parser(num)
			}
			tool.Send_alarm(db, config.IP, target, `<a href="`+path+`">`+tool.HTML_escape(name)+` - `+tool.HTML_escape(sub)+`#`+tool.HTML_escape(num)+`</a>`)
		}
		return value
	})
}

func render_topic_reference(data string, topic_num string, set_id string, set_code string, do_type string) string {
	data = topic_reference_regex.ReplaceAllStringFunc(data, func(value string) string {
		match := topic_reference_regex.FindStringSubmatch(value)
		if len(match) < 5 {
			return value
		}

		label := "#" + match[2]
		reference_code := topic_num
		if match[3] != "" {
			label += "-" + match[3]
			if do_type == "thread" {
				reference_code = match[3]
			}
		}

		path := ""
		if do_type == "thread" && reference_code != "" {
			path = "/thread/" + tool.Url_parser(reference_code) + "#" + tool.Url_parser(match[2])
		} else if do_type != "thread" && set_id != "" && set_code != "" {
			target_code := set_code
			anchor := match[2]
			if match[3] != "" {
				target_code = match[3]
			}
			path = "/bbs/w/" + tool.Url_parser(set_id) + "/" + tool.Url_parser(target_code) + "#" + tool.Url_parser(anchor)
		}
		if path == "" {
			return value
		}

		return match[1] + `<a href="` + path + `">` + label + `</a>` + match[4]
	})

	return topic_call_regex.ReplaceAllStringFunc(data, func(value string) string {
		match := topic_call_regex.FindStringSubmatch(value)
		if len(match) < 4 {
			return value
		}

		label := "@" + match[2]
		return match[1] + `<a href="/w/user:` + tool.Url_parser(match[2]) + `">` + tool.HTML_escape(label) + `</a>` + match[3]
	})
}

func topic_thread_notify(db *sql.DB, config tool.Config, topic_num string, comment_num string, name string, sub string) {
	alarm := `<a href="/thread/` + tool.Url_parser(topic_num) + `#` + tool.Url_parser(comment_num) + `">` + tool.HTML_escape(name) + ` - ` + tool.HTML_escape(sub) + `#` + tool.HTML_escape(comment_num) + `</a>`

	if strings.HasPrefix(name, "user:") {
		target := strings.TrimPrefix(name, "user:")
		target_exists := false
		if tool.IP_or_user(target) {
			var value string
			target_exists = tool.QueryRow_DB(db, "select ip from history where ip = ? limit 1", []any{&value}, target)
			if !target_exists {
				target_exists = tool.QueryRow_DB(db, "select ip from topic where ip = ? limit 1", []any{&value}, target)
			}
		} else {
			var value string
			target_exists = tool.QueryRow_DB(db, "select id from user_set where id = ? limit 1", []any{&value}, target)
		}

		if target_exists {
			tool.Send_alarm(db, config.IP, target, alarm)
		}
	}

	first_user := ""
	if tool.QueryRow_DB(db, "select ip from topic where code = ? and id = '1'", []any{&first_user}, topic_num) && !tool.IP_or_user(first_user) {
		tool.Send_alarm(db, config.IP, first_user, alarm)
	}
}
