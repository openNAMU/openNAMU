package backlink

import (
	"strings"

	"github.com/dlclark/regexp2"
)

var namumark_link_regex = regexp2.MustCompile(
	`\[\[(?<link>[^\[\]|]+)(?:\|[^\[\]]*)?\]\]`,
	regexp2.None,
)

var markdown_link_regex = regexp2.MustCompile(
	`\[[^\]\r\n]*\]\((?<link>[^)\r\n]+)\)`,
	regexp2.None,
)

var macromark_link_regex = regexp2.MustCompile(
	`\[(?<type>a|include)\((?<link>[^,\])\r\n]+)(?:,[^\])\r\n]*)?\)\]`,
	regexp2.None,
)

func Get_backlink(raw_data string, markup string) (map[string][]string, int, bool) {
	var link_regex *regexp2.Regexp

	switch markup {
	case "", "namumark", "namumark_beta":
		link_regex = namumark_link_regex
	case "markdown":
		link_regex = markdown_link_regex
	case "macromark":
		link_regex = macromark_link_regex
	case "raw":
		return map[string][]string{}, 0, true
	default:
		return map[string][]string{}, 0, false
	}

	backlink_list := map[string][]string{}
	link_count := 0

	match, err := link_regex.FindStringMatch(raw_data)
	for err == nil && match != nil {
		link := normalize_link(match.GroupByName("link").String())
		if link == "" {
			match, err = link_regex.FindNextMatch(match)
			continue
		}

		link_type := ""
		if markup == "macromark" && match.GroupByName("type").String() == "include" {
			link_type = "include"
		}

		link_count++
		add_backlink(backlink_list, link, link_type)

		match, err = link_regex.FindNextMatch(match)
	}

	return backlink_list, link_count, true
}

func add_backlink(
	backlink_list map[string][]string,
	link string,
	link_type string,
) {
	for _, saved_type := range backlink_list[link] {
		if saved_type == link_type {
			return
		}
	}

	backlink_list[link] = append(backlink_list[link], link_type)
}

func normalize_link(raw_link string) string {
	link := strings.TrimSpace(raw_link)
	link = strings.TrimPrefix(link, "<")
	link = strings.TrimSuffix(link, ">")
	link = strings.TrimSpace(link)

	if hash_index := strings.Index(link, "#"); hash_index >= 0 {
		link = strings.TrimSpace(link[:hash_index])
	}

	link_lower := strings.ToLower(link)
	if link == "" ||
		strings.HasPrefix(link, "/") ||
		strings.HasPrefix(link, "//") ||
		strings.HasPrefix(link_lower, "http://") ||
		strings.HasPrefix(link_lower, "https://") ||
		strings.HasPrefix(link_lower, "mailto:") {
		return ""
	}

	return link
}
