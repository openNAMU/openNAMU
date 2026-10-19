package backlink

import (
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/dlclark/regexp2"
	"golang.org/x/net/idna"
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
	`\[(?<type>a|include|ex)\((?<link>[^,\])\r\n]+)(?:,[^\])\r\n]*)?\)\]`,
	regexp2.None,
)

// Markdown code-block/inline-code regexes (mirroring markup/markdown.go) are
// used to strip code regions before link extraction, so URLs inside code are
// never collected.
var backlink_fenced_code_regex = regexp2.MustCompile(
	"(?ms)^ {0,3}\\x60\\x60\\x60([^\\n]*)\\n(.*?)^ {0,3}\\x60\\x60\\x60[ \\t]*(?:\\n|$)",
	0,
)
var backlink_tilde_code_regex = regexp2.MustCompile(
	"(?ms)^ {0,3}~~~([^\\n]*)\\n(.*?)^ {0,3}~~~[ \\t]*(?:\\n|$)",
	0,
)
var backlink_inline_code_regex = regexp2.MustCompile("`([^`\\n]+)`", 0)

// backlink_autolink_regex extracts autolinked URLs from namumark text. It uses
// the same character class as the namumark renderer (markdown.go) so that a URL
// containing "s" is matched correctly. The V4 raw-string form passed the literal
// "\s" to the regex, which excluded "\" and "s" rather than whitespace.
var backlink_autolink_regex = regexp2.MustCompile("(?i)<(https?://[^<>\\s]+)>", 0)

func Backlink_redirect_target(raw_data string) string {

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
		return Normalize_link(target)
	}

	return ""
}
func Backlink_redirect_line(line string) bool {
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
	case "raw", "plain":
		return map[string][]string{}, 0, true
	default:
		return map[string][]string{}, 0, false
	}

	backlink_list := map[string][]string{}
	link_count := 0

	// Markdown: strip code blocks/inline code and collect autolinks (<https://…>)
	// before the bracketed-link pass. Images (V1) stay excluded by the existing
	// "image" group check.
	if markup == "markdown" {
		raw_data = Markdown_exclude_code(raw_data)
		Markdown_extract_autolinks(raw_data, backlink_list)
	}

	if markup == "" || markup == "namumark" || markup == "namumark_beta" {
		if redirect_target := Backlink_redirect_target(raw_data); redirect_target != "" {
			Add_backlink(backlink_list, redirect_target, "redirect")
			lines := strings.Split(raw_data, "\n")
			for index, line := range lines {
				if Backlink_redirect_line(line) {
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

		raw_link := match.GroupByName("link").String()

		// External extraction applies only to markdown bracketed links and the
		// macromark "ex" macro. Namumark external links are collected by the
		// renderer's Process_links, so they stay out of this pass.
		is_external_macro := false
		is_macromark_ex := false
		if markup == "markdown" {
			is_external_macro = true
		} else if markup == "macromark" && match.GroupByName("type").String() == "ex" {
			is_external_macro = true
			is_macromark_ex = true
		}
		if is_external_macro {
			ext_host := Normalize_external_host(Clean_link(raw_link))
			if ext_host != "" {
				Add_backlink(backlink_list, ext_host, "external")
				match, err = link_regex.FindNextMatch(match)
				continue
			}
			if is_macromark_ex {
				// The "ex" macro is exclusively external. An invalid target is
				// not an internal link, so skip it to avoid internal backlink
				// pollution (V4 fell through and treated it as internal).
				match, err = link_regex.FindNextMatch(match)
				continue
			}
			// markdown: an invalid external target may still be an internal
			// link, so fall through to the internal handling below.
		}

		var link string
		if markup == "markdown" {
			link = Normalize_markdown_link(raw_link)
		} else {
			link = Normalize_link(raw_link)
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
		Add_backlink(backlink_list, link, link_type)

		match, err = link_regex.FindNextMatch(match)
	}

	return backlink_list, link_count, true
}

func Add_backlink(
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

// Clean_link strips markdown-style URL wrappers and surrounding whitespace so a
// single raw target can be normalized for either external or internal links.
func Clean_link(raw_link string) string {
	link := strings.TrimSpace(raw_link)
	if strings.HasPrefix(link, "<") && strings.Contains(link, ">") {
		link = link[1:strings.Index(link, ">")]
	} else if index := strings.IndexAny(link, " \t\r\n"); index >= 0 {
		link = link[:index]
	}
	return strings.TrimSpace(link)
}

func Normalize_markdown_link(raw_link string) string {
	link := Clean_link(raw_link)
	if strings.HasPrefix(link, "./") || strings.HasPrefix(link, "../") {
		return ""
	}
	return Normalize_link(link)
}

func Normalize_link(raw_link string) string {
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

// Normalize_external_host converts a raw external URL into the external backlink
// title: the lower-cased, scheme- and path-less hostname. IDN host names are
// converted to ASCII via golang.org/x/net/idna (Punycode).
//
// V5 rules:
//
//   - only http/https schemes are accepted;
//   - the host is lower-cased and the trailing DNS dot is stripped;
//   - scheme, path, query, fragment, and port are not stored;
//   - the port, when present, must be a valid value in 1..65535;
//   - user info is rejected;
//   - the stored value is the full hostname (including subdomains);
//   - IP addresses (IPv4/IPv6) and single-label names (for example localhost)
//     are rejected;
//   - no external HTTP/DNS lookup is performed.
//
// It returns "" for anything that is not a valid external hostname.
func Normalize_external_host(raw_link string) string {
	link := strings.TrimSpace(raw_link)
	host := ""
	if strings.Contains(link, "://") {
		parsed, err := url.Parse(link)
		if err != nil {
			return ""
		}
		if parsed.Scheme != "http" && parsed.Scheme != "https" {
			return ""
		}
		if parsed.User != nil {
			return ""
		}
		if port := parsed.Port(); port != "" {
			if n, parse_err := strconv.Atoi(port); parse_err != nil || n < 1 || n > 65535 {
				return ""
			}
		}
		host = parsed.Hostname()
	} else {
		// Bare hostname input (query path). The stored value is already
		// scheme-free, so re-normalize and validate it as a plain name.
		host = link
	}
	host = strings.TrimSuffix(strings.ToLower(host), ".")
	if host == "" {
		return ""
	}
	// Reject IP addresses (IPv4, IPv6, ::1, etc.).
	if net.ParseIP(host) != nil {
		return ""
	}
	// Normalize IDN to ASCII Punycode. idna.ToASCII does not lowercase, so the
	// input is lower-cased above; ASCII input is returned verbatim.
	ascii_host, err := idna.ToASCII(host)
	if err != nil {
		return ""
	}
	if !Is_valid_dns_host(ascii_host) {
		return ""
	}
	return ascii_host
}

// Is_valid_dns_host reports whether host (ASCII, as returned by idna) is a valid
// fully qualified DNS hostname: at least one dot, no leading/trailing dot, no
// consecutive dots, every label uses [a-z0-9-] with [a-z0-9] at each end,
// each label at most 63 bytes, and the whole name at most 253 bytes.
func Is_valid_dns_host(host string) bool {
	if host == "" {
		return false
	}
	// Whole-name length limit (byte count, ASCII/Punycode input).
	if len(host) > 253 {
		return false
	}
	// Must be a fully qualified name: at least one dot. This also rejects
	// single-label names (for example localhost).
	if strings.Index(host, ".") < 0 {
		return false
	}
	// No leading or trailing dot.
	if host[0] == '.' || host[len(host)-1] == '.' {
		return false
	}
	// No consecutive dots (empty label).
	if strings.Contains(host, "..") {
		return false
	}
	// Each label must be non-empty, at most 63 bytes, and use [a-z0-9-] with
	// [a-z0-9] at each end.
	for _, part := range strings.Split(host, ".") {
		if part == "" {
			return false
		}
		if len(part) > 63 {
			return false
		}
		if !Is_label_end(part[0]) || !Is_label_end(part[len(part)-1]) {
			return false
		}
		for i := 1; i < len(part)-1; i++ {
			if !Is_label_char(part[i]) {
				return false
			}
		}
	}
	return true
}

// Is_label_char reports whether ch is a byte valid inside a DNS label.
func Is_label_char(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '-'
}

// Is_label_end reports whether ch is a byte valid at the start or end of a DNS label.
func Is_label_end(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9')
}

// Markdown_exclude_code returns data with fenced, tilde, and inline code
// regions replaced by an inert placeholder so their URLs are not collected.
func Markdown_exclude_code(data string) string {
	placeholder := "\x00CODE\x00"
	fn := func(_ regexp2.Match) string {
		return placeholder
	}
	data, _ = backlink_fenced_code_regex.ReplaceFunc(data, fn, -1, -1)
	data, _ = backlink_tilde_code_regex.ReplaceFunc(data, fn, -1, -1)
	data, _ = backlink_inline_code_regex.ReplaceFunc(data, fn, -1, -1)
	return data
}

// Markdown_extract_autolinks collects raw autolinks (<https://…>) as external
// backlinks. They are matched before the bracketed-link pass so they do not
// double-count.
func Markdown_extract_autolinks(data string, backlink_list map[string][]string) {
	match, err := backlink_autolink_regex.FindStringMatch(data)
	for err == nil && match != nil {
		raw_link := match.GroupByNumber(1).String()
		if ext_host := Normalize_external_host(raw_link); ext_host != "" {
			Add_backlink(backlink_list, ext_host, "external")
		}
		match, err = backlink_autolink_regex.FindNextMatch(match)
	}
}
