package route

import (
	"encoding/xml"
	"strings"
	"time"

	"database/sql"

	"opennamu/route/tool"
)

type rss struct {
	Version string   `xml:"version,attr"`
	Channel rss_channel `xml:"channel"`
}

type rss_channel struct {
	Title         string   `xml:"title"`
	Link          string   `xml:"link"`
	Description   string   `xml:"description"`
	Language      string   `xml:"language"`
	LastBuildDate string   `xml:"lastBuildDate"`
	Item          []rss_item `xml:"item"`
}

type rss_guid struct {
	IsPermaLink bool   `xml:"isPermaLink,attr"`
	Value       string `xml:",chardata"`
}

type rss_item struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Guid        rss_guid `xml:"guid"`
	PubDate     string   `xml:"pubDate,omitempty"`
	Description string   `xml:"description"`
}

func Get_rss_pub_date(date string) string {
	if date == "" {
		return ""
	}

	t, err := time.Parse("2006-01-02 15:04:05", date)
	if err != nil {
		return ""
	}

	return t.Format("Mon, 02 Jan 2006 15:04:05 -0700")
}

func Get_rss_base_url(db *sql.DB, host string) string {
	full := tool.Get_domain(db, true)

	host_part := full
	if i := strings.Index(full, "://"); i >= 0 {
		host_part = full[i+3:]
	}
	if j := strings.Index(host_part, "/"); j >= 0 {
		host_part = host_part[:j]
	}

	if host_part == "" {
		if host != "" {
			return "http://" + host
		}
		return "http://"
	}

	return full
}

func View_rss_recent(config tool.Config, host string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	base_url := Get_rss_base_url(db, host)

	wiki_name := "Wiki"
	tool.QueryRow_DB(db, "select data from other where name = 'name'", []any{&wiki_name})

	language := ""
	tool.QueryRow_DB(db, "select data from other where name = 'language'", []any{&language})

	api_data := Api_list_recent_change(config, "", "50", "1")

	var items []rss_item
	if api_data["response"] == "ok" {
		data_list, _ := api_data["data"].([][]string)
		for _, row := range data_list {
			if len(items) >= 20 {
				break
			}

			if row[6] != "" {
				continue
			}

			revision := row[0]
			doc_name := row[1]

			type_label := ""
			if row[8] != "" {
				type_label = tool.Get_language(db, row[8], true)
			}

			description := "r" + revision
			if type_label != "" {
				description += " (" + type_label + ")"
			}

			items = append(items, rss_item{
				Title: doc_name,
				Link: base_url + "/w/" + tool.Url_parser(doc_name),
				Guid: rss_guid{
					IsPermaLink: false,
					Value: doc_name + "|" + revision,
				},
				PubDate:     Get_rss_pub_date(row[2]),
				Description: description,
			})
		}
	}

	feed := rss{
		Version: "2.0",
		Channel: rss_channel{
			Title:         wiki_name,
			Link:          base_url + "/recent_change",
			Description:   tool.Get_language(db, "recent_change", true),
			Language:      language,
			LastBuildDate: time.Now().Format("Mon, 02 Jan 2006 15:04:05 -0700"),
			Item:          items,
		},
	}

	body, err := xml.Marshal(&feed)
	if err != nil {
		return ""
	}

	return `<?xml version="1.0" encoding="UTF-8"?>
` + string(body)
}