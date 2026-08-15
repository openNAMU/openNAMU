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
	`(?<image>!)?\[[^\]\r\n]*\]\((?<link>(?:[^()\r\n]+|\([^()\r\n]*\))+?)\)`,
	regexp2.None,
)

var macromark_link_regex = regexp2.MustCompile(
	`\[(?<type>a|include)\((?<link>[^,\])\r\n]+)(?:,[^\])\r\n]*)?\)\]`,
	regexp2.None,
)

func backlink_redirect_target(raw_data string) string {

	for _, line := range strings.Split(raw_data, "\n") {
		line = strings.TrimSpace(line)
		line_lower := strings.ToLower(line)
		prefix := ""
		if strings.HasPrefix(line_lower, "#redirect ") {
			prefix = "#redirect "
		} else if strings.HasPrefix(line_lower, "#넘겨주기 ") {
			prefix = "#넘겨주기 "
		}
		if prefix == "" {
			continue
		}

		target := strings.TrimSpace(line[len(prefix):])
		if strings.HasPrefix(target, "[[") && strings.HasSuffix(target, "]]") {
			target = strings.TrimSuffix(strings.TrimPrefix(target, "[["), "]]")
			if pipe_index := strings.Index(target, "|"); pipe_index >= 0 {
				target = target[:pipe_index]
			}
		}
		return normalize_link(target)
	}

	return ""
}
func backlink_redirect_line(line string) bool {
	line = strings.ToLower(strings.TrimSpace(line))
	return strings.HasPrefix(line, "#redirect ") || strings.HasPrefix(line, "#넘겨주기 ")
}

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
	if markup == "" || markup == "namumark" || markup == "namumark_beta" {
		if redirect_target := backlink_redirect_target(raw_data); redirect_target != "" {
			add_backlink(backlink_list, redirect_target, "redirect")
			lines := strings.Split(raw_data, "\n")
			for index, line := range lines {
				if backlink_redirect_line(line) {
					lines[index] = ""
				}
			}
			raw_data = strings.Join(lines, "\n")
		}
	}

	match, err := link_regex.FindStringMatch(raw_data)
	for err == nil && match != nil {
		if markup == "markdown" && match.GroupByName("image").String() != "" {
			match, err = link_regex.FindNextMatch(match)
			continue
		}

		link := match.GroupByName("link").String()
		if markup == "markdown" {
			link = normalize_markdown_link(link)
		} else {
			link = normalize_link(link)
		}
		if link == "" {
			match, err = link_regex.FindNextMatch(match)
			continue
		}

		link_type := ""
		if strings.HasPrefix(strings.ToLower(link), "category:") {
			link_type = "cat"
		}
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

func normalize_markdown_link(raw_link string) string {
	link := strings.TrimSpace(raw_link)
	if strings.HasPrefix(link, "<") && strings.Contains(link, ">") {
		link = link[1:strings.Index(link, ">")]
	} else if index := strings.IndexAny(link, " \t\r\n"); index >= 0 {
		link = link[:index]
	}

	if strings.HasPrefix(link, "./") || strings.HasPrefix(link, "../") {
		return ""
	}

	return normalize_link(link)
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

	if strings.HasPrefix(link, "분류:") {
		link = "category:" + strings.TrimPrefix(link, "분류:")
	}

	return link
}
