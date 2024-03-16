package tool

import (
	"bytes"
	"database/sql"
	"regexp"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

func Markdown(db *sql.DB, db_set map[string]string, data map[string]string) map[string]interface{} {
	input := []byte(data["data"])
	backlink := map[string]map[string]string{}
	link_count := 0

	markdown := goldmark.New(
		goldmark.WithExtensions(extension.Strikethrough, extension.Table),
		goldmark.WithRendererOptions(html.WithHardWraps()),
	)

	var buf bytes.Buffer
	if err := markdown.Convert(input, &buf); err != nil {
		panic(err)
	}

	string_data := buf.String()

	r := regexp.MustCompile(`\[([^\[\]]+)\]\(([^\(\)]*)\)`)
	string_data = r.ReplaceAllStringFunc(string_data, func(m string) string {
		match := r.FindStringSubmatch(m)

		return "<a href=\"" + match[2] + "\">" + match[1] + "</a>"
	})

	// p := bluemonday.UGCPolicy()
	// result := p.Sanitize(string_data)

	r = regexp.MustCompile(`<a href="([^"]+)"`)
	result := r.ReplaceAllStringFunc(string_data, func(m string) string {
		match := r.FindStringSubmatch(m)

		m1, _ := regexp.MatchString(`^https?:\/\/`, match[1])
		if m1 {
			return "<a href=\"" + match[1] + "\" class=\"opennamu_link_out\" target=\"_blank\""
		} else {
			if _, ok := backlink[match[1]]; !ok {
				backlink[match[1]] = map[string]string{}
			}

			var exist string

			stmt, err := db.Prepare(DB_change(db_set, "select title from data where title = ?"))
			if err != nil {
				exist = ""
			}
			defer stmt.Close()

			err = stmt.QueryRow(match[1]).Scan(&exist)
			if err != nil {
				exist = ""
			}

			backlink[match[1]][""] = ""
			link_count += 1

			class := ""
			if exist == "" {
				class = "opennamu_not_exist_link"
			}

			return "<a href=\"/w/" + match[1] + "\" class=\"" + class + "\""
		}
	})

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
	end_data["data"] = result
	end_data["js_data"] = ""
	end_data["backlink"] = end_backlink
	end_data["link_count"] = link_count

	return end_data
}
