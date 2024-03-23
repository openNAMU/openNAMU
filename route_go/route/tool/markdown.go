package tool

import (
	"bytes"
	"database/sql"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

func Markdown(db *sql.DB, db_set map[string]string, data map[string]string) map[string]interface{} {
	backlink := map[string]map[string]string{}
	link_count := 0

	raw_input := data["data"]

	r := regexp.MustCompile(`\[\]\(([^\(\)]+)\)`)
	raw_input = r.ReplaceAllStringFunc(raw_input, func(m string) string {
		match := r.FindStringSubmatch(m)

		return "[" + match[1] + "](" + match[1] + ")"
	})

	r = regexp.MustCompile(`\[([^\[\]]+)\]\(\)`)
	raw_input = r.ReplaceAllStringFunc(raw_input, func(m string) string {
		match := r.FindStringSubmatch(m)

		return "[" + match[1] + "](" + match[1] + ")"
	})

	input := []byte(raw_input)
	markdown := goldmark.New(
		goldmark.WithExtensions(
			extension.Strikethrough,
			extension.Table,
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)

	var buf bytes.Buffer
	if err := markdown.Convert(input, &buf); err != nil {
		panic(err)
	}

	string_data := buf.String()

	code_stack := []int{}
	code_stack_idx := 0
	code_stack_end := map[string]string{}

	r = regexp.MustCompile(`(<code>|<\/code>)`)
	for idx := r.FindStringIndex(string_data); len(idx) != 0; idx = r.FindStringIndex(string_data) {
		if string_data[idx[0]:idx[1]] == "<code>" {
			code_stack = []int{idx[0], idx[1]}
			string_data = strings.Replace(string_data, "<code>", "<0001>", 1)
		} else {
			string_data = strings.Replace(string_data, "<0001>", "<code>", -1)

			code_stack_idx_str := strconv.Itoa(code_stack_idx)
			code_stack_end["code_"+code_stack_idx_str] = string_data[code_stack[0]:idx[1]]
			code_stack_idx++

			string_data = string_data[:code_stack[0]] + "<code_" + code_stack_idx_str + ">" + string_data[idx[1]:]
		}
	}

	// p := bluemonday.UGCPolicy()
	// string_data := p.Sanitize(string_data)

	r = regexp.MustCompile(`\[([^\[\]]+)\]\(([^\(\)]*)\)`)
	string_data = r.ReplaceAllStringFunc(string_data, func(m string) string {
		match := r.FindStringSubmatch(m)

		link := match[2]
		if link == "" {
			link = match[1]
		}

		return "<a href=\"" + link + "\">" + match[1] + "</a>"
	})

	r = regexp.MustCompile(`<code_[0-9]+>`)
	string_data = r.ReplaceAllStringFunc(string_data, func(m string) string {
		m = strings.Replace(m, "<", "", 1)
		m = strings.Replace(m, ">", "", 1)

		return code_stack_end[m]
	})

	r = regexp.MustCompile(`<a href="([^"]+)"`)
	string_data = r.ReplaceAllStringFunc(string_data, func(m string) string {
		match := r.FindStringSubmatch(m)

		m1, _ := regexp.MatchString(`^https?:\/\/`, match[1])
		if m1 {
			return "<a href=\"" + match[1] + "\" class=\"opennamu_link_out\" target=\"_blank\""
		} else {
			link := ""
			link, _ = url.QueryUnescape(match[1])

			if _, ok := backlink[link]; !ok {
				backlink[link] = map[string]string{}
			}

			var exist string

			stmt, err := db.Prepare(DB_change(db_set, "select title from data where title = ?"))
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()

			err = stmt.QueryRow(link).Scan(&exist)
			if err != nil {
				if err == sql.ErrNoRows {
					exist = ""
				} else {
					log.Fatal(err)
				}
			}

			backlink[link][""] = ""
			link_count += 1

			class := ""
			if exist == "" {
				class = "opennamu_not_exist_link"
			}

			return "<a href=\"/w/" + match[1] + "\" class=\"" + class + "\""
		}
	})

	string_data = strings.Replace(string_data, "<ul>", "<ul class=\"opennamu_ul\">", -1)

	end_backlink := [][]string{}
	for k1, v1 := range backlink {
		for k2, v2 := range v1 {
			end_backlink = append(end_backlink, []string{
				data["doc_name"],
				k1,
				k2,
				v2,
			})
		}
	}

	end_data := make(map[string]interface{})
	end_data["data"] = string_data
	end_data["js_data"] = "opennamu_do_toc();"
	end_data["backlink"] = end_backlink
	end_data["link_count"] = link_count

	return end_data
}
