package markup

import (
	"database/sql"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dlclark/regexp2"
	"golang.org/x/net/html"
	"opennamu/route/tool"
)

// namumark_compat_renderer is a temporary compatibility layer for syntax that
// the current small Namumark renderer intentionally does not understand yet.
// It owns all compatibility state so the layer can be removed without
// changing namumark.go.
type namumark_compat_renderer struct {
	db            *sql.DB
	doc_name      string
	render_type   string
	parameter     map[string]any
	data          string
	include_depth int
	collect_only  bool

	tokens         map[string]string
	slash_tokens   map[string]string
	token_prefix   string
	token_count    int
	syntax_count   int
	inter_data     map[string]string
	inter_raw      map[string]string
	inter_wrappers map[string]string
	literal_tokens map[string]bool
	inter_count    int

	backlinks      map[string]namumark_compat_backlink
	backlink_order []string
	link_count     int
	redirect       bool

	categories       []namumark_compat_category
	toc_items        []namumark_compat_toc_item
	footnotes        []namumark_compat_footnote
	footnote_pending []namumark_compat_footnote_group
	footnote_count   int
	footnote_map     map[string]int
	footnote_token   string
	footnote_prefix  string
}

type namumark_compat_backlink struct {
	target    string
	link_type string
	data      string
}

type namumark_compat_category struct {
	target string
	label  string
	blur   bool
	exists bool
}

type namumark_compat_footnote struct {
	name    string
	text    string
	numbers []string
	named   bool
}

type namumark_compat_footnote_group struct {
	index   int
	numbers []string
}

type namumark_compat_toc_item struct {
	number string
	text   string
}

var namumark_compat_list_regex = regexp.MustCompile(`^([ \t]*)(\*|(?:1|a|A|i|I)\.(?:#[0-9]+)?)[ \t]*(.*)$`)
var namumark_compat_middle_regex = regexp2.MustCompile(`\{\{\{([^{](?:(?!\{\{\{|\}\}\})[\s\S])*)?(?:\}|(OPENNAMU_COMPAT_TOKEN_[^<> \t\r\n]+X))\}\}`, 0)
var namumark_compat_if_regex = regexp.MustCompile(`^\s*([^\s]+)\s*(==|!=)\s*(.*?)\s*$`)
var namumark_compat_style_regex = regexp.MustCompile(`(?is)^(style|dark-style)\s*=\s*(?:\x22([^\x22]*)\x22|\x27([^\x27]*)\x27|&quot;(.*?)&quot;|&#x27;(.*?)&#x27;)(?:[ \t\n]+)`)
var namumark_compat_language_regex = regexp.MustCompile(`^[a-zA-Z0-9_+.-]+$`)

var namumark_compat_parameter_regex = regexp.MustCompile(`@([ㄱ-ㅣ가-힣a-zA-Z0-9_]+)(?:=([^@\n]+))?@`)
var namumark_compat_html_dimension_regex = regexp.MustCompile(`^[0-9]+(?:px|%)?$`)

var namumark_compat_html_tags = map[string]bool{
	"a": true, "b": true, "del": true, "div": true, "em": true,
	"i": true, "iframe": true, "s": true, "span": true, "strong": true,
	"sub": true, "sup": true, "bold": true,
}

var namumark_compat_html_style_properties = map[string]bool{
	"background":       true,
	"background-color": true,
	"background-image": true,
	"border":           true,
	"border-color":     true,
	"border-radius":    true,
	"border-right":     true,
	"border-top":       true,
	"color":            true,
	"display":          true,
	"float":            true,
	"font-size":        true,
	"font-family":      true,
	"font-style":       true,
	"font-weight":      true,
	"height":           true,
	"image-rendering":  true,
	"line-height":      true,
	"letter-spacing":   true,
	"margin":           true,
	"margin-bottom":    true,
	"margin-left":      true,
	"margin-right":     true,
	"margin-top":       true,
	"max-height":       true,
	"max-width":        true,
	"min-height":       true,
	"min-width":        true,
	"opacity":          true,
	"overflow":         true,
	"padding":          true,
	"padding-bottom":   true,
	"padding-left":     true,
	"padding-right":    true,
	"padding-top":      true,
	"text-align":       true,
	"text-decoration":  true,
	"text-shadow":      true,
	"vertical-align":   true,
	"white-space":      true,
	"width":            true,
	"word-break":       true,
}

func compat_html_safe_url(value string, iframe bool) string {
	value = strings.TrimSpace(value)
	if value == "" || strings.HasPrefix(value, "//") {
		return ""
	}
	for _, char := range value {
		if char < 0x20 || char == 0x7f {
			return ""
		}
	}

	parsed, err := url.Parse(value)
	if err != nil {
		return ""
	}
	scheme := strings.ToLower(parsed.Scheme)
	if iframe {
		if (scheme != "http" && scheme != "https") || parsed.Hostname() == "" || parsed.User != nil {
			return ""
		}
		switch strings.ToLower(parsed.Hostname()) {
		case "www.youtube.com", "www.google.com", "play-tv.kakao.com":
			return value
		default:
			return ""
		}
	}

	switch scheme {
	case "", "http", "https", "mailto", "tel":
		if (scheme == "http" || scheme == "https") && parsed.Hostname() == "" {
			return ""
		}
		return value
	default:
		return ""
	}
}

func compat_html_safe_style(value string) string {
	style_data := []string{}
	for _, declaration := range strings.Split(value, ";") {
		parts := strings.SplitN(declaration, ":", 2)
		if len(parts) != 2 {
			continue
		}

		property := strings.ToLower(strings.TrimSpace(parts[0]))
		if !namumark_compat_html_style_properties[property] {
			continue
		}
		property_value := strings.TrimSpace(parts[1])
		lower_value := strings.ToLower(property_value)
		if property_value == "" || strings.ContainsAny(property_value, "\r\n") {
			continue
		}
		unsafe_value := false
		for _, value := range []string{"url", "expression", "javascript", "vbscript", "binding", "behavior", "@import"} {
			if strings.Contains(lower_value, value) {
				unsafe_value = true
				break
			}
		}
		if unsafe_value {
			continue
		}
		style_data = append(style_data, declaration)
	}
	if len(style_data) == 0 {
		return ""
	}

	result := strings.Join(style_data, ";")
	if strings.HasSuffix(value, ";") {
		result += ";"
	}
	return result
}

func compat_html_dimension(value string) string {
	value = strings.TrimSpace(value)
	if !namumark_compat_html_dimension_regex.MatchString(value) {
		return ""
	}
	return value
}

func compat_sanitize_html_node(node *html.Node) string {
	if node == nil {
		return ""
	}
	if node.Type == html.TextNode {
		return compat_html_escape(node.Data)
	}
	if node.Type != html.ElementNode {
		data := ""
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			data += compat_sanitize_html_node(child)
		}
		return data
	}

	tag_name := strings.ToLower(node.Data)
	if !namumark_compat_html_tags[tag_name] {
		data := ""
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			data += compat_sanitize_html_node(child)
		}
		return data
	}

	attributes := ""
	if tag_name == "div" || tag_name == "span" {
		for _, attr := range node.Attr {
			if strings.EqualFold(attr.Key, "style") {
				if style := compat_html_safe_style(attr.Val); style != "" {
					attributes = ` style="` + compat_html_escape(style) + `"`
				}
				break
			}
		}
	} else if tag_name == "a" {
		href := ""
		for _, attr := range node.Attr {
			if strings.EqualFold(attr.Key, "href") {
				href = compat_html_safe_url(attr.Val, false)
				break
			}
		}
		attributes = ` class="opennamu_link_out" href="` + compat_html_escape(href) + `"`
	} else if tag_name == "iframe" {
		src := ""
		width := ""
		height := ""
		for _, attr := range node.Attr {
			switch strings.ToLower(attr.Key) {
			case "src":
				src = compat_html_safe_url(attr.Val, true)
			case "width":
				width = compat_html_dimension(attr.Val)
			case "height":
				height = compat_html_dimension(attr.Val)
			}
		}
		if src == "" {
			return ""
		}
		attributes = ` src="` + compat_html_escape(src) + `"`
		if width != "" {
			attributes += ` width="` + compat_html_escape(width) + `"`
		}
		if height != "" {
			attributes += ` height="` + compat_html_escape(height) + `"`
		}
		attributes += ` allowfullscreen frameborder="0"`
	}

	data := "<" + tag_name + attributes + ">"
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		data += compat_sanitize_html_node(child)
	}
	return data + "</" + tag_name + ">"
}

func compat_sanitize_html(data string) string {
	nodes, err := html.ParseFragment(strings.NewReader(data), nil)
	if err != nil {
		return compat_html_escape(data)
	}

	result := ""
	for _, node := range nodes {
		result += compat_sanitize_html_node(node)
	}
	return result
}

func compat_render_parameter_value(value any) string {
	if value == nil {
		return ""
	}
	if value_string, ok := value.(string); ok {
		return value_string
	}
	return fmt.Sprint(value)
}

func compat_render_parameter_data(data string, parameter map[string]any) string {
	if len(parameter) == 0 {
		return data
	}
	return namumark_compat_parameter_regex.ReplaceAllStringFunc(data, func(value string) string {
		match := namumark_compat_parameter_regex.FindStringSubmatch(value)
		parameter_value, ok := parameter[match[1]]
		if !ok {
			return match[2]
		}
		return compat_render_parameter_value(parameter_value)
	})
}

var namumark_compat_single_macro_regex = regexp.MustCompile(`(?is)\[([a-zA-Z가-힣]+)\]`)

var namumark_compat_heading_regex = regexp.MustCompile(`(?s)<h([1-6])><a href="#toc">([0-9.]+)\. </a>(.*?)</h[1-6]>`)
var namumark_compat_toc_regex = regexp.MustCompile(`(?s)<div class="opennamu_TOC" id="toc">.*?</div>`)

var namumark_compat_bold_regex = regexp.MustCompile(`(?s)<b>(.*?)</b>`)
var namumark_compat_strike_regex = regexp.MustCompile(`(?s)<s>(.*?)</s>`)

func compat_fix_heading_data(db *sql.DB, data string) string {
	toc_item := [][2]string{}
	data = namumark_compat_heading_regex.ReplaceAllStringFunc(data, func(raw string) string {
		match := namumark_compat_heading_regex.FindStringSubmatch(raw)
		if len(match) < 4 {
			return raw
		}
		toc_item = append(toc_item, [2]string{match[2], match[3]})
		return `<h` + match[1] + `><a href="#toc" id="s-` + compat_html_escape(match[2]) + `">` + match[2] + `. </a>` + match[3] + `</h` + match[1] + `>`
	})
	if len(toc_item) == 0 || !namumark_compat_toc_regex.MatchString(data) {
		return data
	}

	toc_title := "toc"
	if db != nil {
		toc_title = tool.Get_language(db, "toc", true)
	}
	toc_data := `<div class="opennamu_TOC" id="toc"><span class="opennamu_TOC_title">` + compat_html_escape(toc_title) + `</span><br>`
	for _, item := range toc_item {
		indent := strings.Count(item[0], ".")
		toc_data += `<br>` + strings.Repeat(`<span style="margin-left: 10px;"></span>`, indent) + `<span class="opennamu_TOC_list"><a href="#s-` + compat_html_escape(item[0]) + `">` + item[0] + `. </a>` + item[1] + `</span>`
	}
	toc_data += `</div>`
	return namumark_compat_toc_regex.ReplaceAllString(data, toc_data)
}

func new_namumark_compat_renderer(
	db *sql.DB,
	doc_name string,
	data string,
	render_type string,
	parameter map[string]any,
	include_depth int,
	collect_only bool,
) *namumark_compat_renderer {
	if parameter == nil {
		parameter = map[string]any{}
	}

	renderer := &namumark_compat_renderer{
		db:              db,
		doc_name:        doc_name,
		render_type:     render_type,
		parameter:       parameter,
		data:            compat_html_escape(strings.ReplaceAll(data, "\r", "")),
		include_depth:   include_depth,
		collect_only:    collect_only,
		tokens:          map[string]string{},
		slash_tokens:    map[string]string{},
		inter_data:      map[string]string{},
		inter_raw:       map[string]string{},
		inter_wrappers:  map[string]string{},
		literal_tokens:  map[string]bool{},
		backlinks:       map[string]namumark_compat_backlink{},
		footnote_map:    map[string]int{},
		footnote_prefix: "open_namu_fn_" + tool.Sha224(doc_name + data)[:12],
	}
	renderer.token_prefix = "OPENNAMU_COMPAT_TOKEN_" + fmt.Sprintf("%p", renderer) + "_"
	return renderer
}

func (class *namumark_compat_renderer) reserve(data string) string {
	class.token_count++
	token := class.token_prefix + strconv.Itoa(class.token_count) + "X"
	class.tokens[token] = data
	return token
}

func (class *namumark_compat_renderer) reserve_literal(data string) string {
	token := class.reserve(data)
	class.literal_tokens[token] = true
	return token
}

func (class *namumark_compat_renderer) reserve_inter(data string) string {
	class.inter_count++
	token := "OPENNAMU_COMPAT_INTER_TOKEN_" + strconv.Itoa(class.inter_count) + "X"
	class.inter_data[token] = data
	return token
}

func (class *namumark_compat_renderer) reserve_slash(data string) string {
	token := class.reserve(data)
	class.slash_tokens[token] = data
	return token
}

func (class *namumark_compat_renderer) restore_slash(data string) string {
	for token, token_data := range class.slash_tokens {
		data = strings.ReplaceAll(data, token, tool.HTML_unescape(token_data))
	}
	return data
}

func (class *namumark_compat_renderer) middle_literal(data string) string {
	for token, token_data := range class.slash_tokens {
		data = strings.ReplaceAll(data, token, "\\"+token_data)
	}
	return data
}

func (class *namumark_compat_renderer) restore_inter_literal(data string) string {
	for count := 0; count <= len(class.inter_raw)+1; count++ {
		changed := false
		for token, wrapper := range class.inter_wrappers {
			if strings.Contains(data, wrapper) {
				data = strings.ReplaceAll(data, wrapper, class.inter_raw[token])
				changed = true
			}
		}
		for token, token_data := range class.inter_raw {
			if _, has_wrapper := class.inter_wrappers[token]; has_wrapper {
				continue
			}
			if strings.Contains(data, token) {
				data = strings.ReplaceAll(data, token, token_data)
				changed = true
			}
		}
		if !changed {
			break
		}
	}
	return data
}

func (class *namumark_compat_renderer) restore(data string) string {
	for count := 0; count <= len(class.tokens)+1; count++ {
		changed := false
		for token, token_data := range class.tokens {
			if strings.Contains(data, token) {
				data = strings.ReplaceAll(data, token, token_data)
				changed = true
			}
		}
		if !changed {
			break
		}
	}
	return data
}

func (class *namumark_compat_renderer) get_render_setting(name string) string {
	if setting, ok := class.parameter["__opennamu_skin_set"].(map[string]string); ok {
		if value := setting[name]; value != "" {
			return value
		}
	}
	if setting, ok := class.parameter["__opennamu_skin_set"].(map[string]any); ok {
		if value, ok := setting[name].(string); ok && value != "" {
			return value
		}
	}
	return "default"
}

func (class *namumark_compat_renderer) add_backlink(target string, link_type string, data string) {
	if target == "" {
		return
	}

	key := target + "\x00" + link_type
	if old_data, ok := class.backlinks[key]; ok {
		if old_data.data == "" && data != "" {
			old_data.data = data
			class.backlinks[key] = old_data
		}
		return
	}

	class.backlinks[key] = namumark_compat_backlink{target, link_type, data}
	class.backlink_order = append(class.backlink_order, key)
}

func (class *namumark_compat_renderer) find_document(name string) (string, bool) {
	actual := ""
	query := "select title from data where title = ?"
	case_insensitive := ""
	tool.QueryRow_DB(class.db, "select data from other where name = 'link_case_insensitive'", []any{&case_insensitive})
	if case_insensitive != "" {
		query += " collate nocase"
	}

	if tool.QueryRow_DB(class.db, query, []any{&actual}, name) {
		return actual, true
	}
	return name, false
}

func (class *namumark_compat_renderer) compat_split_anchor(target string) (string, string) {
	hash_index := strings.LastIndex(target, "#")
	if hash_index < 0 || hash_index == len(target)-1 {
		return target, ""
	}

	anchor := strings.TrimSpace(target[hash_index+1:])
	if anchor == "" {
		return target[:hash_index], ""
	}
	return target[:hash_index], "#" + class.compat_url_parser(anchor)
}

func (class *namumark_compat_renderer) normalize_target(target string) string {
	target = strings.TrimSpace(target)
	target = class.restore_slash(target)

	base_name := class.doc_name
	relative_path := false

	for strings.HasPrefix(target, "../") {
		relative_path = true
		target = strings.TrimPrefix(target, "../")
		if slash_index := strings.LastIndex(base_name, "/"); slash_index >= 0 {
			base_name = base_name[:slash_index]
		}
	}

	if relative_path && base_name != "" {
		target = base_name + "/" + target
	}

	if strings.HasPrefix(target, "/") {
		target = strings.TrimPrefix(target, "/")
		if base_name != "" {
			target = base_name + "/" + target
		}
	}

	target = normalize_namumark_link(target)
	return target
}

func (class *namumark_compat_renderer) add_document_link(target string, link_type string, data string) string {
	target = class.normalize_target(target)
	if target == "" {
		return ""
	}

	actual, exists := class.find_document(target)
	if exists {
		target = actual
	}
	class.add_backlink(target, link_type, data)
	if !exists {
		class.add_backlink(target, "no", "")
	}
	return target
}

func compat_is_external_link(target string) bool {
	target_lower := strings.ToLower(strings.TrimSpace(target))
	return strings.HasPrefix(target_lower, "http://") ||
		strings.HasPrefix(target_lower, "https://")
}

func compat_is_interwiki_link(target string) bool {
	parts := strings.SplitN(target, ":", 3)
	return len(parts) == 3 && (strings.EqualFold(parts[0], "inter") || strings.EqualFold(parts[0], "인터"))
}

func (class *namumark_compat_renderer) get_outer_link_data(target string, label string) (string, string, bool) {
	parsed, err := url.Parse(target)
	if err != nil || parsed.Host == "" {
		return "", "", false
	}

	domain := strings.ToLower(parsed.Host)
	link_name := ""
	icon := ""
	if !tool.QueryRow_DB(class.db, "select html, plus_t from html_filter where kind = 'outer_link' and plus = ?", []any{&link_name, &icon}, domain) {
		return "", "", false
	}

	if icon == "" {
		return compat_escape_value(link_name + ":"), "opennamu_link_inter", true
	}
	if strings.ContainsAny(icon, "<>") {
		return icon, "opennamu_link_inter", true
	}
	if strings.Contains(class.restore(label), `"`+icon+`"`) {
		return "", "opennamu_link_inter", true
	}
	return `<img src="` + compat_html_escape(icon) + `">`, "opennamu_link_inter", true
}

func (class *namumark_compat_renderer) get_interwiki_url(target string) (string, bool) {
	parts := strings.SplitN(target, ":", 3)
	if len(parts) != 3 || !compat_is_interwiki_link(target) {
		return "", false
	}

	page, anchor := class.compat_split_anchor(parts[2])
	plus := ""
	icon := ""
	if !tool.QueryRow_DB(class.db, "select plus, plus_t from html_filter where kind = 'inter_wiki' and html = ?", []any{&plus, &icon}, parts[1]) {
		return "", false
	}

	link := plus + class.compat_url_parser(page) + anchor
	inter_sub_mode := ""
	tool.QueryRow_DB(class.db, "select plus_t from html_filter where kind = 'inter_wiki_sub' and html = ?", []any{&inter_sub_mode}, parts[1])
	if inter_sub_mode == "under_bar" {
		link = strings.ReplaceAll(link, "%20", "_")
	}
	return link, true
}

func (class *namumark_compat_renderer) process_redirect(data string) string {
	lines := strings.Split(data, "\n")
	for index, line := range lines {
		if index > 0 {
			break
		}
		trimmed := strings.TrimSpace(line)
		lower := strings.ToLower(trimmed)
		prefix := ""
		if strings.HasPrefix(lower, "#redirect ") {
			prefix = "#redirect "
		} else if strings.HasPrefix(lower, "#넘겨주기 ") {
			prefix = "#넘겨주기 "
		}
		if prefix == "" {
			continue
		}

		target := tool.HTML_unescape(strings.TrimSpace(trimmed[len(prefix):]))
		if strings.HasPrefix(target, "[[") && strings.HasSuffix(target, "]]") {
			target = strings.TrimSuffix(strings.TrimPrefix(target, "[["), "]]")
			if pipe_index := strings.Index(target, "|"); pipe_index >= 0 {
				target = target[:pipe_index]
			}
		}

		main_target, anchor := class.compat_split_anchor(target)
		if compat_is_external_link(main_target) {
			class.redirect = true
			if class.collect_only || class.include_depth > 0 {
				return ""
			}
			lines[index] = class.reserve(`<a href="` + compat_html_escape(main_target+anchor) + `">(GO)</a>`)
			return strings.Join(lines, "\n")
		}
		if compat_is_interwiki_link(main_target) {
			class.redirect = true
			if class.collect_only || class.include_depth > 0 {
				return ""
			}
			link, ok := class.get_interwiki_url(main_target + anchor)
			if !ok {
				lines[index] = ""
				return strings.Join(lines, "\n")
			}
			lines[index] = class.reserve(`<a href="` + compat_html_escape(link) + `">(GO)</a>`)
			return strings.Join(lines, "\n")
		}

		main_target = class.add_document_link(main_target, "redirect", anchor)
		class.redirect = true
		if class.collect_only || class.include_depth > 0 {
			return ""
		}

		link := "/w_from/" + class.compat_url_parser(main_target) + anchor
		lines[index] = class.reserve(`<a href="` + compat_html_escape(link) + `">(GO)</a>`)
		return strings.Join(lines, "\n")
	}

	return data
}

func (class *namumark_compat_renderer) process_category(target string, label string) string {
	target = strings.TrimSpace(target)
	target = class.normalize_target(target)
	target = strings.TrimPrefix(strings.TrimPrefix(target, "category:"), "분류:")

	blur := false
	if strings.HasSuffix(strings.ToLower(target), "#blur") {
		target = target[:len(target)-5]
		blur = true
	}
	if label == "" {
		label = target
	}

	category_target := "category:" + target
	actual, exists := class.find_document(category_target)
	if exists {
		category_target = actual
	}
	class.add_backlink(category_target, "cat", "")
	if label != target {
		class.add_backlink(category_target, "cat_view", label)
	}
	if blur {
		class.add_backlink(category_target, "cat_blur", "")
	}
	if !exists {
		class.add_backlink(category_target, "no", "")
	}

	for _, category := range class.categories {
		if category.target == category_target {
			return class.reserve("")
		}
	}
	class.categories = append(class.categories, namumark_compat_category{category_target, label, blur, exists})
	return class.reserve("")
}

func compat_file_option_value(value string) string {
	value = tool.HTML_unescape(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, ";", "")
	value = strings.ReplaceAll(value, "{", "")
	value = strings.ReplaceAll(value, "}", "")
	value = strings.ReplaceAll(value, "\"", "")
	value = strings.ReplaceAll(value, "'", "")
	return value
}

func compat_file_px(value string) string {
	value = compat_file_option_value(value)
	if regexp.MustCompile(`^[0-9]+$`).MatchString(value) {
		return value + "px"
	}
	return value
}

func compat_file_options(raw string, default_alt string) (string, string, string) {
	raw = tool.HTML_unescape(raw)
	alt := default_alt
	style := ""
	theme := ""
	for _, item := range strings.Split(raw, "&") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		parts := strings.SplitN(item, "=", 2)
		if len(parts) != 2 {
			if alt == default_alt {
				alt = item
			}
			continue
		}

		name := strings.ToLower(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(parts[1])
		switch name {
		case "alt", "title":
			alt = value
		case "width":
			style += "width:" + compat_file_px(value) + ";"
		case "height":
			style += "height:" + compat_file_px(value) + ";"
		case "align":
			if value == "left" || value == "right" {
				style += "float:" + value + ";"
			} else if value == "center" {
				style += "display:block;margin-left:auto;margin-right:auto;"
			}
		case "bgcolor":
			style += "background:" + compat_file_option_value(value) + ";"
		case "theme":
			value = strings.ToLower(strings.TrimSpace(value))
			if value == "dark" || value == "light" {
				theme = value
			}
		case "border-radius":
			style += "border-radius:" + compat_file_px(value) + ";"
		case "rendering":
			if value == "pixelated" {
				style += "image-rendering:pixelated;"
			}
		}
	}
	return alt, style, theme
}

func (class *namumark_compat_renderer) file_theme_visible(theme string) bool {
	if theme == "" {
		return true
	}
	darkmode := class.get_render_setting("main_css_darkmode")
	if darkmode == "" || darkmode == "default" {
		darkmode = "0"
	}
	return (theme == "dark" && darkmode == "1") || (theme == "light" && darkmode != "1")
}

func (class *namumark_compat_renderer) process_file(target string, label string) string {
	lower_target := strings.ToLower(target)
	external := strings.HasPrefix(lower_target, "out:") || strings.HasPrefix(lower_target, "외부:")
	file_name := target
	if external {
		file_name = target[strings.Index(target, ":")+1:]
	}
	if !external {
		if colon_index := strings.Index(file_name, ":"); colon_index >= 0 {
			prefix := strings.ToLower(file_name[:colon_index])
			if prefix == "file" || prefix == "파일" {
				file_name = file_name[colon_index+1:]
			}
		}
	}
	file_name = class.restore_slash(tool.HTML_unescape(strings.TrimSpace(file_name)))
	if file_name == "" {
		return class.reserve("")
	}

	_, style, theme := compat_file_options(label, file_name)
	alt := strings.TrimSpace(target)
	if external {
		if !class.file_theme_visible(theme) {
			return class.reserve("")
		}
		if !strings.HasPrefix(strings.ToLower(file_name), "http://") && !strings.HasPrefix(strings.ToLower(file_name), "https://") {
			return class.reserve(compat_html_escape(alt))
		}
		image := `<img style="` + compat_html_escape(style) + `" alt="` + compat_html_escape(alt) + `" src="` + compat_html_escape(file_name) + `">`
		return class.reserve(`<a title="` + compat_html_escape(alt) + `" href="` + compat_html_escape(file_name) + `">` + image + `</a>`)
	}

	file_target := "file:" + file_name
	extension := "jpg"
	base_name := file_name
	if dot_index := strings.LastIndex(file_name, "."); dot_index > 0 && dot_index < len(file_name)-1 {
		base_name = file_name[:dot_index]
		extension = strings.ToLower(file_name[dot_index+1:])
	}
	actual, exists := class.find_document(file_target)
	if exists {
		file_target = actual
	}
	class.add_backlink(file_target, "file", "")
	if !exists {
		class.add_backlink(file_target, "no", "")
		return class.reserve(`<a class="opennamu_not_exist_link" title="` + compat_html_escape(alt) + `" href="/upload?name=` + class.compat_url_parser(base_name) + `">(` + compat_html_escape(alt) + `)</a>`)
	}
	if !class.file_theme_visible(theme) {
		return class.reserve("")
	}

	storage_name := tool.File_name_to_dir(base_name, extension)
	rev := "1"
	tool.QueryRow_DB(class.db, "select id from history where title = ? order by date desc limit 1", []any{&rev}, file_target)
	image_url := "/image/" + class.compat_url_parser(storage_name) + ".cache_v" + class.compat_url_parser(rev)
	file_url := "/w/file:" + class.compat_url_parser(file_name)
	image := `<img style="` + compat_html_escape(style) + `" alt="` + compat_html_escape(alt) + `" src="` + compat_html_escape(image_url) + `">`
	return class.reserve(`<a title="` + compat_html_escape(alt) + `" href="` + compat_html_escape(file_url) + `">` + image + `</a>`)
}

func (class *namumark_compat_renderer) process_interwiki(target string, label string) string {
	parts := strings.SplitN(target, ":", 3)
	if len(parts) != 3 {
		return ""
	}
	name := parts[1]
	page, anchor := class.compat_split_anchor(parts[2])
	plus := ""
	icon := ""
	if !tool.QueryRow_DB(class.db, "select plus, plus_t from html_filter where kind = 'inter_wiki' and html = ?", []any{&plus, &icon}, name) {
		return ""
	}

	if label == "" {
		label = page
	}
	link := plus + class.compat_url_parser(page) + anchor
	inter_sub_mode := ""
	tool.QueryRow_DB(class.db, "select plus_t from html_filter where kind = 'inter_wiki_sub' and html = ?", []any{&inter_sub_mode}, name)
	if inter_sub_mode == "under_bar" {
		link = strings.ReplaceAll(link, "%20", "_")
	}
	if icon == "" {
		icon = name + ":"
	}
	icon_data := compat_escape_value(icon)
	if strings.ContainsAny(icon, "<>") {
		icon_data = icon
	}
	return class.reserve(`<a class="opennamu_link_inter" title="` + compat_html_escape(name+":"+page) + `" href="` + compat_html_escape(link) + `">` + icon_data + compat_escape_value(label) + `</a>`)
}

func (class *namumark_compat_renderer) process_links(data string) string {
	for {
		previous := data
		data = compat_replace_regex2(data, `(?i)\[\[((?:(?!\[\[|\]\]|\||<|>).)+)(?:\|((?:(?!\[\[|\]\]|\|).)+))?\]\]`, func(match regexp2.Match) string {
			raw := match.String()
			body := strings.TrimSpace(tool.HTML_unescape(match.GroupByNumber(1).String()))
			label := strings.TrimSpace(tool.HTML_unescape(match.GroupByNumber(2).String()))

			lower_body := strings.ToLower(body)
			switch {
			case strings.HasPrefix(lower_body, "file:"), strings.HasPrefix(lower_body, "파일:"),
				strings.HasPrefix(lower_body, "out:"), strings.HasPrefix(lower_body, "외부:"):
				return class.process_file(body, label)
			case strings.HasPrefix(lower_body, "category:"), strings.HasPrefix(lower_body, "분류:"):
				return class.process_category(body, label)
			case strings.HasPrefix(lower_body, "inter:"), strings.HasPrefix(lower_body, "인터:"):
				if result := class.process_interwiki(body, label); result != "" {
					return result
				}
				return raw
			}

			main_target, anchor := class.compat_split_anchor(body)
			if strings.HasPrefix(body, "#") {
				if label == "" {
					label = body
				}
				return class.reserve(`<a class=" " title="`+compat_html_escape(body)+`" href="`+compat_html_escape("#"+class.compat_url_parser(strings.TrimPrefix(body, "#")))+`">`) + compat_html_escape(label) + class.reserve(`</a>`)
			}
			if compat_is_external_link(main_target) {
				main_target = class.restore_slash(main_target)
				if label == "" {
					label = body
				}
				link_target := main_target + anchor
				icon, link_class, matched := class.get_outer_link_data(main_target, label)
				if !matched {
					link_class = "opennamu_link_out"
				}
				return class.reserve(`<a class="`+link_class+`" target="_blank" title="`+compat_html_escape(link_target)+`" href="`+compat_html_escape(link_target)+`">`) + class.reserve(icon) + compat_escape_value(label) + class.reserve(`</a>`)
			}

			if main_target == "" {
				return raw
			}
			class.link_count++

			normalized_target := class.normalize_target(main_target)
			actual_target, exists := class.find_document(normalized_target)
			if exists {
				normalized_target = actual_target
			}
			class.add_backlink(normalized_target, "", "")
			if !exists {
				class.add_backlink(normalized_target, "no", "")
			}
			if label == "" {
				label = body
			}
			label_lower := strings.ToLower(label)
			if strings.HasPrefix(label_lower, ":분류:") ||
				strings.HasPrefix(label_lower, ":category:") ||
				strings.HasPrefix(label_lower, ":파일:") ||
				strings.HasPrefix(label_lower, ":file:") {
				label = strings.TrimPrefix(label, ":")
			}
			link_exist := ""
			link_same := ""
			if !exists {
				link_exist = "opennamu_not_exist_link"
			}
			if normalized_target == class.doc_name {
				link_same = "opennamu_same_link"
			}
			link_class := ` class="` + link_exist + " " + link_same + `"`
			return class.reserve(`<a`+link_class+` title="`+compat_html_escape(normalized_target+anchor)+`" href="/w/`+class.compat_url_parser(normalized_target)+anchor+`">`) + compat_html_escape(label) + class.reserve(`</a>`)
		})
		if data == previous {
			break
		}
	}
	return data
}

func compat_split_macro_args(data string) []string {
	data = strings.ReplaceAll(data, ",,", "\x00")
	parts := strings.Split(data, ",")
	for index := range parts {
		parts[index] = strings.ReplaceAll(parts[index], "\x00", ",")
		parts[index] = strings.TrimSpace(parts[index])
	}
	return parts
}

func (class *namumark_compat_renderer) merge_child(child *namumark_compat_renderer) {
	for _, key := range child.backlink_order {
		entry := child.backlinks[key]
		class.add_backlink(entry.target, entry.link_type, entry.data)
	}
	class.link_count += child.link_count
}

func (class *namumark_compat_renderer) process_includes(data string) string {
	for {
		previous := data
		data = compat_replace_regex2(data, `(?i)\[include\(((?:(?!\[include\(|\)\]|</div>).)+)\)\](\n?)`, func(match regexp2.Match) string {
			suffix := match.GroupByNumber(2).String()
			if class.render_type == "include" {
				return ""
			}

			include_name := ""
			include_parameter := map[string]any{}
			for key, value := range class.parameter {
				include_parameter[key] = value
			}
			for _, item := range compat_split_macro_args(match.GroupByNumber(1).String()) {
				if item == "" {
					continue
				}
				item = tool.HTML_unescape(item)
				parts := strings.SplitN(item, "=", 2)
				if len(parts) == 1 {
					include_name = item
					continue
				}
				include_parameter[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
			include_name_org := strings.TrimSpace(include_name)
			include_name = class.normalize_target(include_name)
			include_name_url := include_name
			if strings.Contains(include_name, "OPENNAMU_COMPAT_TOKEN_") {
				include_name_url = "<" + include_name + ">"
				include_name_org = regexp.MustCompile(`<[^<>]*>`).ReplaceAllString(class.restore(include_name_org), "")
			}
			if include_name == "" {
				return suffix
			}

			actual, exists := class.find_document(include_name)
			if exists {
				include_name = actual
			}
			class.add_backlink(include_name, "include", "")
			if !exists {
				class.add_backlink(include_name, "no", "")
				return class.reserve(`<a class="opennamu_not_exist_link" href="/w/`+class.compat_url_parser(include_name_url)+`">(`+compat_escape_value(include_name_org)+`)</a>`) + suffix
			}
			if class.collect_only || class.include_depth >= 8 {
				return suffix
			}

			include_data := ""
			tool.QueryRow_DB(class.db, "select data from data where title = ?", []any{&include_data}, include_name)
			child := new_namumark_compat_renderer(class.db, class.doc_name, include_data, "include", include_parameter, class.include_depth+1, false)
			child.prepare()
			child_output := `<div class="opennamu_render_complete">` + child.render_output() + `</div>`
			include_link := ""
			if class.get_render_setting("main_css_include_link") == "use" {
				include_link = `<a href="/w/` + class.compat_url_parser(include_name_url) + `">(` + compat_escape_value(include_name_org) + `)</a><br>`
			}
			child_output = include_link + child_output
			class.merge_child(child)
			return class.reserve(child_output) + suffix
		})
		if data == previous {
			break
		}
	}
	return data
}

func (class *namumark_compat_renderer) get_footnote_label(footnote namumark_compat_footnote, number string, view_number bool) string {
	label := footnote.name
	if class.get_render_setting("main_css_footnote_number") == "only_number" {
		label = footnote.numbers[0]
	}
	if view_number && footnote.named && class.get_render_setting("main_css_view_real_footnote_num") == "on" {
		label += " (" + number + ")"
	}
	return label
}
func (class *namumark_compat_renderer) add_footnote_pending(index int, number string) {
	for pending_index := range class.footnote_pending {
		if class.footnote_pending[pending_index].index == index {
			class.footnote_pending[pending_index].numbers = append(class.footnote_pending[pending_index].numbers, number)
			return
		}
	}
	class.footnote_pending = append(class.footnote_pending, namumark_compat_footnote_group{index: index, numbers: []string{number}})
}
func (class *namumark_compat_renderer) process_footnotes(data string) string {
	footnote_regex := `(?i)(?:\[\*((?:(?!\[\*|\]| ).)+)?(?: ((?:(?!\[\*|\]).)+))?\]|\[(?:각주|footnote)\])`
	data = compat_replace_regex2(data, footnote_regex, func(match regexp2.Match) string {
		raw := match.String()
		if strings.EqualFold(raw, "[각주]") || strings.EqualFold(raw, "[footnote]") {
			class.footnote_count++
			if class.collect_only {
				return ""
			}
			footnote_html := class.make_footnotes(class.footnote_pending)
			class.footnote_pending = nil
			if footnote_html == "" {
				return ""
			}
			return class.reserve(footnote_html)
		}

		name := strings.TrimSpace(match.GroupByNumber(1).String())
		class.footnote_count++
		text := match.GroupByNumber(2).String()
		text = class.process_macros(text)
		named := name != ""
		if name == "" {
			name = strconv.Itoa(class.footnote_count)
		}
		index, exists := class.footnote_map[name]
		if !exists {
			class.footnotes = append(class.footnotes, namumark_compat_footnote{name: name, text: text, named: named})
			index = len(class.footnotes) - 1
			class.footnote_map[name] = index
		}
		number := strconv.Itoa(class.footnote_count)
		class.footnotes[index].numbers = append(class.footnotes[index].numbers, number)
		class.add_footnote_pending(index, number)

		if class.collect_only {
			return ""
		}
		first_number := class.footnotes[index].numbers[0]
		footnote_label := class.get_footnote_label(class.footnotes[index], number, true)
		fn := class.footnote_prefix + "fn_" + first_number
		rfn := class.footnote_prefix + "rfn_" + number
		title_data := class.render_inline(class.footnotes[index].text)
		title_data = class.restore(title_data)
		title := compat_escape_value(regexp.MustCompile(`<[^<>]*>`).ReplaceAllString(title_data, ""))
		footnote_label = compat_escape_value(footnote_label)
		footnote_set := class.get_render_setting("main_css_footnote_set")
		if footnote_set == "spread" || footnote_set == "popup" {
			return class.reserve(`<sup><a title="` + title + `" id="` + rfn + `" href="javascript:void(0);">(` + footnote_label + `)</a></sup><span class="opennamu_spead_footnote" id="` + rfn + `_load" style="display: none;"></span>`)
		}
		if footnote_set == "popover" {
			return class.reserve(`<span id="` + rfn + `_over"><sup><a title="` + title + `" id="` + rfn + `" href="javascript:void(0);">(` + footnote_label + `)</a></sup><span class="opennamu_popup_footnote" id="` + rfn + `_load" style="display: none;"></span></span>`)
		}
		return class.reserve(`<sup><a title="` + title + `" id="` + rfn + `" href="#` + fn + `">(` + footnote_label + `)</a></sup>`)
	})

	if !class.collect_only && len(class.footnote_pending) > 0 {
		class.footnote_token = class.reserve(class.make_footnotes(class.footnote_pending))
		class.footnote_pending = nil
	}
	return data
}

func (class *namumark_compat_renderer) make_footnotes(groups []namumark_compat_footnote_group) string {
	if len(groups) == 0 {
		return ""
	}

	data := `<div class="opennamu_footnote">`
	for group_index, group := range groups {
		if group_index > 0 {
			data += `<hr class="main_hr">`
		}
		footnote := class.footnotes[group.index]
		first_number := group.numbers[0]
		if len(group.numbers) > 1 {
			data += `(` + compat_escape_value(footnote.name) + `) `
			for _, number := range group.numbers {
				data += `<sup><a id="` + class.footnote_prefix + `fn_` + number + `" href="#` + class.footnote_prefix + `rfn_` + number + `">(` + compat_html_escape(number) + `)</a></sup> `
			}
		} else {
			data += `<a id="` + class.footnote_prefix + `fn_` + first_number + `" href="#` + class.footnote_prefix + `rfn_` + first_number + `">(` + compat_escape_value(footnote.name) + `) </a> `
		}
		data += `<footnote_title id="` + class.footnote_prefix + `fn_` + first_number + `_title">` + class.render_inline(footnote.text) + `</footnote_title>`
	}
	return data + `</div>`
}

func (class *namumark_compat_renderer) process_math(data string) string {
	render := func(expression string) string {
		if class.collect_only {
			return ""
		}
		expression = strings.TrimSpace(strings.ReplaceAll(expression, "\n", " "))
		expression = tool.HTML_unescape(expression)
		return class.reserve(`<code class="opennamu_math" data-math="` + compat_html_escape(expression) + `">` + compat_html_escape(expression) + `</code>`)
	}
	data = compat_replace_regex2(data, `(?i)\[math\(((?:(?!\[math\(|\)\]).|\n)+)\)\]`, func(match regexp2.Match) string {
		return render(match.GroupByNumber(1).String())
	})
	return compat_replace_regex2(data, `(?i)&lt;math&gt;((?:(?!&lt;math&gt;|&lt;\/math&gt;).)+)&lt;\/math&gt;`, func(match regexp2.Match) string {
		return render(match.GroupByNumber(1).String())
	})
}

func compat_middle_style(value string) string {
	value = strings.TrimSpace(value)
	value = tool.HTML_unescape(value)
	lower_value := strings.ToLower(value)
	for _, blocked := range []string{"url(", "expression", "javascript:", "@import", "position", "<", ">"} {
		if strings.Contains(lower_value, blocked) {
			return ""
		}
	}
	value = strings.NewReplacer("\"", "", "'", "", "\\", "").Replace(value)
	return value
}

func compat_dark_mode_value(values []string, darkmode string) string {
	if len(values) == 0 {
		return ""
	}
	if darkmode == "1" && len(values) > 1 {
		return values[1]
	}
	return values[0]
}

func compat_middle_color_style(name string, body string, darkmode string) (string, string, bool) {
	if !strings.HasPrefix(name, "#") && !strings.HasPrefix(name, "@") {
		return "", body, false
	}
	background := strings.HasPrefix(name, "@")
	color_data := strings.TrimPrefix(strings.TrimPrefix(name, "@"), "#")
	color_parts := strings.Split(color_data, ",")
	for index := range color_parts {
		color_parts[index] = strings.TrimPrefix(color_parts[index], "@")
	}
	if len(color_parts) == 0 || color_parts[0] == "" {
		return "", body, false
	}
	color_regex := regexp.MustCompile(`^(?:#?[0-9a-fA-F]{3,8}|[\p{L}][\p{L}\p{N}_-]*)$`)
	hex_color_regex := regexp.MustCompile(`^[0-9a-fA-F]{3,8}$`)
	for index, color := range color_parts {
		if !color_regex.MatchString(color) {
			return "", body, false
		}
		if hex_color_regex.MatchString(color) {
			color_parts[index] = "#" + color
		}
	}
	body = strings.TrimLeft(body, " \t\n")
	if body == "" {
		body = name
	}
	property := "color"
	if background {
		property = "background-color"
	}
	return property + ":" + compat_dark_mode_value(color_parts, darkmode), body, true
}
func compat_middle_size_style(name string, body string) (string, string, bool) {
	sizes := map[string]string{
		"+5": "200", "+4": "180", "+3": "160", "+2": "140", "+1": "120",
		"-1": "90", "-2": "80", "-3": "70", "-4": "60", "-5": "50",
	}
	size, ok := sizes[name]
	if !ok {
		return "", body, false
	}
	return "font-size:" + size + "%", strings.TrimLeft(body, " \t\n"), true
}
func (class *namumark_compat_renderer) process_middle_block(middle_data string) string {
	middle_data = strings.TrimPrefix(middle_data, "\n")
	middle_name := middle_data
	body := ""
	if space_index := strings.IndexAny(middle_data, " \t\n"); space_index >= 0 {
		middle_name = middle_data[:space_index]
		body = middle_data[space_index:]
	}
	middle_name = strings.ToLower(strings.TrimSpace(middle_name))
	body = strings.TrimPrefix(body, " ")

	if strings.HasPrefix(middle_name, "opennamu_compat_token_") {
		return class.reserve_literal(class.middle_literal(strings.Trim(middle_data, "\n")))
	}

	switch {
	case middle_name == "#!if":
		condition_line := body
		condition_body := ""
		if newline_index := strings.IndexByte(body, '\n'); newline_index >= 0 {
			condition_line = body[:newline_index]
			condition_body = body[newline_index+1:]
		}
		condition_match := namumark_compat_if_regex.FindStringSubmatch(condition_line)
		if len(condition_match) < 4 {
			return ""
		}
		parameter_value, parameter_exists := class.parameter[condition_match[1]]
		expected := strings.TrimSpace(condition_match[3])
		expected = strings.Trim(expected, "\"'")
		matches := false
		if strings.EqualFold(expected, "null") {
			matches = !parameter_exists || parameter_value == nil
		} else if parameter_exists && parameter_value != nil {
			matches = compat_render_parameter_value(parameter_value) == expected
		}
		if condition_match[2] == "!=" {
			matches = !matches
		}
		if matches {
			return condition_body
		}
		return ""

	case middle_name == "#!html":
		if class.collect_only {
			return ""
		}
		for token := range class.literal_tokens {
			if strings.Contains(body, token) {
				class.tokens[token] = compat_legacy_html_literal(class.tokens[token])
			}
		}
		body = class.restore_inter_literal(body)
		body = tool.HTML_unescape(body)
		body = strings.ReplaceAll(body, "&amp;nbsp;", "&nbsp;")
		return class.reserve(compat_sanitize_html(strings.Trim(body, "\n")))

	case middle_name == "#!syntax":
		language := "python"
		code := body
		if newline_index := strings.IndexByte(body, '\n'); newline_index >= 0 {
			language = strings.TrimSpace(body[:newline_index])
			code = body[newline_index+1:]
		}
		if language == "" {
			language = "python"
		}
		if language == "asm" || language == "assembly" {
			language = "x86arm"
		}
		if !namumark_compat_language_regex.MatchString(language) {
			language = "text"
		}
		if class.collect_only {
			return ""
		}
		code = tool.HTML_unescape(code)
		syntax_id := "opennamu_syntax_" + strconv.Itoa(class.syntax_count)
		class.syntax_count++
		return class.reserve(`<pre id="syntax"><code class="` + compat_html_escape(language) + `" id="` + syntax_id + `">` + compat_html_escape(code) + `</code></pre>`)

	case middle_name == "#!wiki":
		wiki_body := body
		style := ""
		darkmode := class.get_render_setting("main_css_darkmode")
		for {
			style_match := namumark_compat_style_regex.FindStringSubmatch(wiki_body)
			if len(style_match) < 4 {
				break
			}
			if strings.EqualFold(style_match[1], "style") || darkmode == "1" {
				style_value := ""
				for style_index := 2; style_index < len(style_match); style_index++ {
					if style_match[style_index] != "" {
						style_value = compat_middle_style(style_match[style_index])
						break
					}
				}
				style += style_value
			}
			wiki_body = wiki_body[len(style_match[0]):]
		}
		wiki_body = strings.Trim(wiki_body, "\n")
		if class.collect_only {
			return wiki_body
		}
		inter_token := class.reserve_inter(wiki_body)
		class.inter_raw[inter_token] = class.middle_literal("{{{" + middle_data + "}}}")
		open_token := class.reserve("<div style=\"\">")
		close_token := class.reserve("</div>")
		if style != "" {
			open_token = class.reserve(`<div style="` + compat_html_escape(style) + `">`)
		}
		class.inter_wrappers[inter_token] = open_token + inter_token + close_token
		return open_token + inter_token + close_token

	case middle_name == "#!folding":
		title := "test"
		folding_body := body
		if newline_index := strings.IndexByte(body, '\n'); newline_index >= 0 {
			title = strings.TrimSpace(body[:newline_index])
			folding_body = strings.Trim(body[newline_index+1:], "\n")
		}
		if title == "" {
			title = "test"
		}
		if class.collect_only {
			return folding_body
		}
		inter_token := class.reserve_inter(folding_body)
		class.inter_raw[inter_token] = class.middle_literal("{{{" + middle_data + "}}}")
		open_token := class.reserve("<details><summary>" + compat_escape_value(title) + "</summary><div class=\"opennamu_folding\">")
		close_token := class.reserve("</div></details>")
		class.inter_wrappers[inter_token] = open_token + inter_token + close_token
		return open_token + inter_token + close_token

	case middle_name == "#!dark" || middle_name == "#!white":
		darkmode := class.get_render_setting("main_css_darkmode")
		if darkmode == "" || darkmode == "default" {
			darkmode = "0"
		}
		if (middle_name == "#!dark" && darkmode == "1") || (middle_name == "#!white" && darkmode != "1") {
			return body
		}
		return ""
	}
	if style, style_body, ok := compat_middle_size_style(middle_name, body); ok {
		if class.collect_only {
			return style_body
		}
		return class.reserve(`<span style="`+style+`">`) + style_body + class.reserve("</span>")
	}
	if style, color_body, ok := compat_middle_color_style(middle_name, body, class.get_render_setting("main_css_darkmode")); ok {
		if class.collect_only {
			return color_body
		}
		return class.reserve(`<span style="`+style+`">`) + color_body + class.reserve("</span>")
	}
	return class.reserve_literal(class.middle_literal(strings.Trim(middle_data, "\n")))
}

func (class *namumark_compat_renderer) process_legacy_html_suffix(data string) string {
	open_index := strings.Index(data, "{{{")
	if open_index < 0 {
		return data
	}
	close_offset := strings.Index(data[open_index+3:], "}}}")
	if close_offset < 0 {
		return data
	}
	close_index := open_index + 3 + close_offset
	middle_data := data[open_index+3 : close_index]
	middle_name := strings.TrimSpace(middle_data)
	if space_index := strings.IndexAny(middle_name, " \t\n"); space_index >= 0 {
		middle_name = middle_name[:space_index]
	}
	known_middle := strings.HasPrefix(middle_name, "#!") || strings.HasPrefix(middle_name, "#")
	if !known_middle || !strings.HasPrefix(data[close_index+3:], "}}}") {
		return data
	}
	literal_token := class.reserve_literal(class.middle_literal("{{{" + middle_data + "}}}"))
	return data[:open_index] + literal_token + data[close_index+6:]
}

func (class *namumark_compat_renderer) process_middle(data string) string {
	find_slash_close := func(value string) int {
		for token, token_data := range class.slash_tokens {
			if token_data == "}" && strings.HasPrefix(value, token) {
				return len(token)
			}
		}
		return 0
	}

	search_index := 0
	for {
		open_index := strings.Index(data[search_index:], "{{{")
		if open_index < 0 {
			break
		}
		open_index += search_index

		depth := 0
		close_index := -1
		for index := open_index; index < len(data); {
			if strings.HasPrefix(data[index:], "{{{") {
				depth++
				index += 3
				continue
			}
			if slash_length := find_slash_close(data[index:]); slash_length > 0 && strings.HasPrefix(data[index+slash_length:], "}}") {
				depth--
				if depth == 0 {
					close_index = index + slash_length + 2
					break
				}
				index += slash_length + 2
				continue
			}
			if strings.HasPrefix(data[index:], "}}}") {
				depth--
				if depth == 0 {
					close_index = index + 3
					break
				}
				index += 3
				continue
			}
			index++
		}
		if close_index < 0 {
			break
		}

		if strings.HasPrefix(data[open_index:], "{{{{{{") {
			search_index = close_index
			continue
		}

		middle_data := data[open_index+3 : close_index-3]
		middle_start := strings.TrimSpace(middle_data)
		if strings.HasPrefix(middle_start, "#!html") && strings.Contains(middle_data, "{{{") {
			inner_index := strings.Index(middle_data, "{{{")
			inner_open := open_index + 3 + inner_index
			inner_depth := 0
			inner_close := -1
			for index := inner_open; index < len(data); {
				if strings.HasPrefix(data[index:], "{{{") {
					inner_depth++
					index += 3
					continue
				}
				if slash_length := find_slash_close(data[index:]); slash_length > 0 && strings.HasPrefix(data[index+slash_length:], "}}") {
					inner_depth--
					if inner_depth == 0 {
						inner_close = index + slash_length + 2
						break
					}
					index += slash_length + 2
					continue
				}
				if strings.HasPrefix(data[index:], "}}}") {
					inner_depth--
					if inner_depth == 0 {
						inner_close = index + 3
						break
					}
					index += 3
					continue
				}
				index++
			}
			if inner_close > inner_open {
				literal_token := class.reserve_literal(data[inner_open:inner_close])
				data = data[:inner_open] + literal_token + data[inner_close:]
				search_index = open_index
				continue
			}
		}

		middle_name := middle_start
		if space_index := strings.IndexAny(middle_name, " \t\n"); space_index >= 0 {
			middle_name = middle_name[:space_index]
		}
		known_middle := strings.HasPrefix(middle_name, "#!") || strings.HasPrefix(middle_name, "#")
		for _, name := range []string{"+1", "+2", "+3", "+4", "+5", "-1", "-2", "-3", "-4", "-5"} {
			if middle_name == name {
				known_middle = true
				break
			}
		}
		if !strings.Contains(middle_data, "{{{") || known_middle {
			search_index = open_index + 3
			continue
		}
		literal_token := class.reserve_literal(class.middle_literal(strings.Trim(middle_data, "\n")))
		data = data[:open_index] + literal_token + data[close_index:]
		search_index = open_index + len(literal_token)
	}

	search_index = 0
	for {
		open_index := strings.Index(data[search_index:], "{{{{{{")
		if open_index < 0 {
			break
		}
		open_index += search_index

		close_start := -1
		close_end := -1
		close_marker_start := -1
		close_marker_end := -1
		depth := 2
		for index := open_index + 6; index < len(data); {
			if close_marker_start >= 0 && depth == 1 {
				look_index := index
				for look_index < len(data) && (data[look_index] == ' ' || data[look_index] == '\t' || data[look_index] == '\r' || data[look_index] == '\n') {
					look_index++
				}
				if look_index < len(data) && data[look_index] != '}' {
					close_start = close_marker_start
					close_end = close_marker_end
					break
				}
			}
			if data[index] == '{' {
				run_end := index
				for run_end < len(data) && data[run_end] == '{' {
					run_end++
				}
				if run_end-index >= 3 {
					depth += (run_end - index) / 3
				}
				index = run_end
				continue
			}
			if data[index] == '}' {
				run_end := index
				for run_end < len(data) && data[run_end] == '}' {
					run_end++
				}
				close_count := (run_end - index) / 3
				if close_count >= depth {
					if close_marker_start >= 0 {
						close_start = close_marker_start
						close_end = run_end
					} else {
						close_start = index + (depth-2)*3
						close_end = close_start + 6
					}
					break
				}
				if close_count > 0 {
					if close_marker_start < 0 && depth-close_count == 1 {
						close_marker_start = index
						close_marker_end = run_end
					}
					depth -= close_count
				}
				index = run_end
				continue
			}
			index++
		}
		if close_start < 0 {
			search_index = open_index + 6
			continue
		}

		middle_data := data[open_index+6 : close_start]
		middle_token := class.reserve_literal(class.middle_literal("{{{" + middle_data + "}}}"))
		remainder := data[close_end:]
		if strings.HasPrefix(strings.TrimSpace(middle_data), "#!html") {
			remainder = class.process_legacy_html_suffix(remainder)
		}
		data = data[:open_index] + middle_token + remainder
		search_index = open_index + len(middle_token)
	}

	middle_count := strings.Count(data, "{{{") * 10
	if middle_count < 20 {
		middle_count = 20
	}
	for count := 0; count < middle_count; count++ {
		previous := data
		result, err := namumark_compat_middle_regex.ReplaceFunc(data, func(match regexp2.Match) string {
			middle_slash := match.GroupByNumber(2).String()
			middle_data := match.GroupByNumber(1).String()
			if middle_slash != "" {
				if class.slash_tokens[middle_slash] != "}" {
					return match.String()
				}
				middle_name := strings.TrimSpace(middle_data)
				if space_index := strings.IndexAny(middle_name, " \t\n"); space_index >= 0 {
					middle_name = middle_name[:space_index]
				}
				known_middle := strings.HasPrefix(middle_name, "#!") || strings.HasPrefix(middle_name, "#")
				for _, name := range []string{"+1", "+2", "+3", "+4", "+5", "-1", "-2", "-3", "-4", "-5"} {
					if middle_name == name {
						known_middle = true
						break
					}
				}
				if known_middle {
					return class.reserve_literal(class.middle_literal("{{{" + middle_data + "}}}"))
				}
				middle_data += "\\"
			}
			return class.process_middle_block(middle_data)
		}, -1, 1)
		if err != nil || result == previous {
			break
		}
		data = result
	}
	if len(class.inter_data) > 0 {
		data = class.resolve_inter_data(data)
	}
	return data
}
func (class *namumark_compat_renderer) render_inter(data string) string {
	child := new_namumark_compat_renderer(class.db, class.doc_name, data, "inter", class.parameter, class.include_depth, false)
	child.data = strings.ReplaceAll(data, "\r", "")
	child.prepare()
	result := child.render_output()
	result = strings.ReplaceAll(result, "||", class.reserve("||"))
	class.merge_child(child)
	return result
}

func (class *namumark_compat_renderer) resolve_inter_data(data string) string {
	resolved := map[string]string{}
	var resolve func(string) string
	resolve = func(token string) string {
		if value, ok := resolved[token]; ok {
			return value
		}
		body, ok := class.inter_data[token]
		if !ok {
			return token
		}
		for nested_token := range class.inter_data {
			if strings.Contains(body, nested_token) {
				body = strings.ReplaceAll(body, nested_token, resolve(nested_token))
			}
		}
		resolved[token] = class.render_inter(body)
		return resolved[token]
	}
	replace_inter := func(value string) string {
		for token := range class.inter_data {
			if strings.Contains(value, token) {
				value = strings.ReplaceAll(value, token, resolve(token))
			}
		}
		return value
	}
	data = replace_inter(data)
	for token := range class.literal_tokens {
		if token_data, ok := class.tokens[token]; ok {
			class.tokens[token] = class.restore_inter_literal(token_data)
		}
	}
	return data
}

func (class *namumark_compat_renderer) render_inline(data string) string {
	return class.render_text(data)
}

func (class *namumark_compat_renderer) render_table(lines []string) string {
	table_style := ""
	table_class := ""
	div_style := ""
	caption := ""
	type table_cell struct {
		data    string
		colspan string
	}
	rows := [][]table_cell{}

	style_value := func(value string) string {
		return compat_middle_style(strings.TrimSpace(value))
	}
	add_style := func(target *string, property string, value string) {
		value = style_value(value)
		if value != "" {
			*target += property + ":" + value + ";"
		}
	}
	color_value := func(value string) string {
		value = style_value(value)
		return compat_dark_mode_value(strings.Split(value, ","), class.get_render_setting("main_css_darkmode"))
	}
	add_color_style := func(target *string, property string, value string) {
		value = color_value(value)
		if value != "" {
			*target += property + ":" + value + ";"
		}
	}
	known_parameter := map[string]bool{
		"tablebgcolor": true, "tablewidth": true, "tableheight": true,
		"tablealign": true, "tableclass": true, "tabletextalign": true,
		"tablecolor": true, "tablebordercolor": true,
		"rowbgcolor": true, "rowtextalign": true, "rowcolor": true,
		"colcolor": true, "colbgcolor": true, "coltextalign": true,
		"bgcolor": true, "color": true, "width": true, "height": true,
		"keepall": true, "rowkeepall": true, "colkeepall": true, "nopad": true,
	}
	html_tag_name := map[string]bool{"a": true, "b": true, "br": true, "code": true, "del": true, "div": true, "em": true, "i": true, "mark": true, "p": true, "pre": true, "s": true, "small": true, "span": true, "strong": true, "sub": true, "sup": true, "u": true}
	parse_parameter := func(parameter string, table *string, row *string, cell *string, col *string, div *string, colspan *string, rowspan *string) bool {
		recognized := false
		for _, item := range strings.Split(parameter, "<") {
			item = strings.TrimSuffix(strings.TrimSpace(item), ">")
			if item == "" {
				continue
			}
			parts := strings.SplitN(item, "=", 2)
			name := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(parts[0]), " ", ""))
			if known_parameter[name] {
				recognized = true
			}
			value := ""
			if len(parts) == 2 {
				value = strings.TrimSpace(parts[1])
			}
			switch name {
			case "tablebgcolor":
				add_color_style(table, "background", value)
			case "tablewidth":
				add_style(div, "width", compat_file_px(value))
				*table += "width:100%;"
			case "tableheight":
				add_style(table, "height", compat_file_px(value))
			case "tablealign":
				if value == "right" {
					*div += "float:right;"
				} else if value == "center" {
					*div += "margin:auto;"
					*table += "margin:auto;"
				}
			case "tableclass":
				table_class = value
			case "tabletextalign":
				add_style(table, "text-align", value)
			case "tablecolor":
				add_color_style(table, "color", value)
			case "tablebordercolor":
				value = color_value(value)
				if value != "" {
					*table += "border:2px solid " + value + ";"
				}
			case "rowbgcolor":
				add_color_style(row, "background", value)
			case "rowtextalign":
				add_style(row, "text-align", value)
			case "rowcolor":
				add_color_style(row, "color", value)
			case "colcolor":
				add_color_style(col, "color", value)
			case "colbgcolor":
				add_color_style(col, "background", value)
			case "coltextalign":
				add_style(col, "text-align", value)
			case "bgcolor":
				add_color_style(cell, "background", value)
			case "color":
				add_color_style(cell, "color", value)
			case "width":
				add_style(cell, "width", compat_file_px(value))
			case "height":
				add_style(cell, "height", compat_file_px(value))
			case "keepall":
				*cell += "word-break:keep-all !important;"
			case "rowkeepall":
				*row += "word-break:keep-all !important;"
			case "colkeepall":
				*col += "word-break:keep-all !important;"
			case "nopad":
				*cell += "padding:0 !important;"
			default:
				if len(parts) == 1 && strings.HasPrefix(name, "-") {
					recognized = true
					*colspan = strings.TrimPrefix(name, "-")
				} else if len(parts) == 1 && (strings.HasPrefix(name, "|") || strings.HasPrefix(name, "^|") || strings.HasPrefix(name, "v|")) {
					recognized = true
					*rowspan = strings.TrimLeft(name, "^v|")
					if strings.HasPrefix(name, "^|") {
						*cell += "vertical-align:top;"
					} else if strings.HasPrefix(name, "v|") {
						*cell += "vertical-align:bottom;"
					}
				} else if len(parts) == 1 && (name == "(" || name == ":" || name == ")") {
					recognized = true
					align := map[string]string{"(": "left", ":": "center", ")": "right"}[name]
					*cell += "text-align:" + align + " !important;"
				} else if len(parts) == 1 && !html_tag_name[name] && len(name) > 2 && regexp.MustCompile(`^(?:#?[0-9a-fA-F]{3,8}|[a-zA-Z][a-zA-Z0-9_-]*)$`).MatchString(name) {
					recognized = true
					add_color_style(cell, "background", name)
				}
			}
		}
		return recognized
	}

	for _, line := range lines {
		line = strings.TrimLeft(strings.TrimSuffix(line, "\r"), " \t")
		if strings.HasPrefix(line, "|") && !strings.HasPrefix(line, "||") {
			caption_end := strings.Index(line[1:], "|")
			if caption_end >= 0 {
				caption_end++
				if strings.Contains(line[caption_end+1:], "||") {
					if caption == "" {
						caption = strings.TrimSpace(line[1:caption_end])
					}
					line = "||" + line[caption_end+1:]
				}
			}
		}
		caption_line := strings.TrimPrefix(strings.TrimSuffix(line, "||"), "||")
		if strings.HasPrefix(caption_line, "|") && !strings.HasPrefix(caption_line, "||") && strings.HasSuffix(caption_line, "|") && !strings.HasSuffix(caption_line, "||") && caption == "" {
			caption = strings.Trim(caption_line, "|")
			rows = append(rows, []table_cell{})
			continue
		}

		row := []table_cell{}
		line_index := 0
		for line_index < len(line) && strings.HasPrefix(line[line_index:], "||") {
			colspan := 0
			for strings.HasPrefix(line[line_index:], "||") {
				colspan++
				line_index += 2
			}
			if line_index >= len(line) {
				break
			}
			data_end := strings.Index(line[line_index:], "||")
			if data_end < 0 {
				data_end = len(line)
			} else {
				data_end += line_index
			}
			row = append(row, table_cell{line[line_index:data_end], strconv.Itoa(colspan)})
			if data_end >= len(line) {
				break
			}
			line_index = data_end
		}
		rows = append(rows, row)
	}

	rendered_rows := []string{}
	column_styles := map[int]string{}
	rowspan_slots := map[int]int{}
	for _, row := range rows {
		row_style := ""
		if len(row) == 0 {
			rendered_rows = append(rendered_rows, `<tr style=""></tr>`)
			continue
		}
		column_index := 0
		rendered_cells := ""
		for _, table_cell_data := range row {
			for rowspan_slots[column_index] > 0 {
				rowspan_slots[column_index]--
				column_index++
			}
			cell_style := ""
			col_style := ""
			colspan := table_cell_data.colspan
			rowspan := ""
			cell_data := table_cell_data.data
			cell_data = strings.TrimLeft(cell_data, "\n")
			for strings.HasPrefix(cell_data, "<") {
				end_index := strings.Index(cell_data, ">")
				if end_index < 0 {
					break
				}
				if !parse_parameter(cell_data[1:end_index], &table_style, &row_style, &cell_style, &col_style, &div_style, &colspan, &rowspan) {
					break
				}
				cell_data = cell_data[end_index+1:]
			}
			for strings.HasPrefix(cell_data, "&lt;") {
				end_index := strings.Index(cell_data, "&gt;")
				if end_index < 0 {
					break
				}
				if !parse_parameter(tool.HTML_unescape(cell_data[:end_index+4]), &table_style, &row_style, &cell_style, &col_style, &div_style, &colspan, &rowspan) {
					break
				}
				cell_data = cell_data[end_index+4:]
			}
			leading_space := strings.HasPrefix(cell_data, " ")
			if leading_space {
				cell_data = cell_data[1:]
			}
			auto_cell_data := cell_data
			if leading_space {
				auto_cell_data = " " + auto_cell_data
			}
			if strings.HasSuffix(cell_data, " ") {
				cell_data = strings.TrimSuffix(cell_data, " ")
			}
			if !strings.Contains(cell_style, "text-align:") {
				switch {
				case strings.HasPrefix(auto_cell_data, " ") && strings.HasSuffix(auto_cell_data, " "):
					cell_style += "text-align: center;"
				case strings.HasPrefix(auto_cell_data, " "):
					cell_style += "text-align: right;"
				case strings.HasSuffix(auto_cell_data, " "):
					cell_style += "text-align: left;"
				default:
					cell_style += "text-align: left;"
				}
			}
			column_styles[column_index] += col_style
			cell_style_data := column_styles[column_index] + cell_style
			cell_data_rendered := cell_data
			trailing_break_count := len(cell_data) - len(strings.TrimRight(cell_data, "\n"))
			if strings.Contains(cell_data, "\n") {
				cell_data_rendered = class.process_blocks(strings.Trim(cell_data, "\n"))
			}
			cell_rendered := strings.ReplaceAll(class.render_inline(cell_data_rendered), "\n", "<br>")
			cell_rendered += strings.Repeat("<br>", trailing_break_count)
			rendered_cells += `<td colspan="` + compat_html_escape(colspan) + `" rowspan="` + compat_html_escape(rowspan) + `" style="` + compat_html_escape(cell_style_data) + `">` + cell_rendered + `</td>`
			column_span := tool.Str_to_int(colspan)
			if column_span < 1 {
				column_span = 1
			}
			row_span := tool.Str_to_int(rowspan)
			if row_span > 1 {
				for column_offset := 0; column_offset < column_span; column_offset++ {
					rowspan_slots[column_index+column_offset] = row_span - 1
				}
			}
			column_index += column_span
		}
		rendered_rows = append(rendered_rows, `<tr style="`+compat_html_escape(row_style)+`">`+rendered_cells+`</tr>`)
	}

	data := `<div class="table_safe" style="` + compat_html_escape(div_style) + `"><table class="` + compat_html_escape(table_class) + `" style="` + compat_html_escape(table_style) + `">`
	if caption != "" {
		data += `<caption>` + class.render_inline(caption) + `</caption>`
	}
	return data + strings.Join(rendered_rows, "") + `</table></div>`
}

func (class *namumark_compat_renderer) render_list(lines []string) string {
	list_numbers := map[string][]int{}
	data := `<ul>`
	to_alpha := func(number int) string {
		result := ""
		for number > 0 {
			number--
			result = string(rune('a'+number%26)) + result
			number /= 26
		}
		return result
	}
	to_roman := func(number int) string {
		values := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
		symbols := []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
		result := ""
		for index, value := range values {
			for number >= value {
				result += symbols[index]
				number -= value
			}
		}
		return result
	}
	for _, line := range lines {
		match := namumark_compat_list_regex.FindStringSubmatch(line)
		if len(match) < 3 {
			continue
		}
		indent := match[1]
		marker := match[2]
		level := 1
		item_class := "opennamu_list_none"
		prefix := ""
		if strings.HasPrefix(marker, "*") {
			level = len(indent)
			if level < 1 {
				level = 1
			}
			item_class = "opennamu_list_5"
			if level <= 4 {
				item_class = "opennamu_list_" + strconv.Itoa(level)
			}
		} else {
			level = len(indent)
			if level < 1 {
				level = 1
			}
			kind := string(marker[0])
			start := 0
			if hash_index := strings.Index(marker, "#"); hash_index >= 0 {
				start = tool.Str_to_int(marker[hash_index+1:])
			}
			if start < 1 {
				start = 1
			}
			for old_kind := range list_numbers {
				if old_kind != kind {
					delete(list_numbers, old_kind)
				}
			}
			if list_numbers[kind] == nil {
				list_numbers[kind] = []int{}
			}
			if len(list_numbers[kind]) < level {
				for len(list_numbers[kind]) < level {
					list_numbers[kind] = append(list_numbers[kind], 1)
				}
			} else {
				list_numbers[kind][level-1]++
				list_numbers[kind] = list_numbers[kind][:level]
			}
			if start != 1 {
				list_numbers[kind][level-1] = start
			}
			number := list_numbers[kind][level-1]
			if kind == "1" && class.get_render_setting("main_css_list_view_change") == "on" {
				number_data := []string{}
				for _, number_value := range list_numbers[kind] {
					if number_value != 0 {
						number_data = append(number_data, strconv.Itoa(number_value))
					}
				}
				prefix = strings.Join(number_data, "-")
			} else {
				switch kind {
				case "a":
					prefix = to_alpha(number)
				case "A":
					prefix = strings.ToUpper(to_alpha(number))
				case "i":
					prefix = strings.ToLower(to_roman(number))
				case "I":
					prefix = to_roman(number)
				default:
					prefix = strconv.Itoa(number)
				}
			}
			prefix += ". "
		}
		data += `<li style="margin-left: ` + strconv.Itoa((level-1)*20) + `px;" class="` + item_class + `">` + prefix + class.render_inline(match[3]) + `</li>`
	}
	return data + `</ul>`
}

func (class *namumark_compat_renderer) compat_media_macro(name string, data string) string {
	code := ""
	data = tool.HTML_unescape(data)
	width := "640px"
	height := "360px"
	start := ""
	end := ""
	if name == "instagram" || name == "tiktok" {
		width, height = "360px", "480px"
	} else if name == "facebook" {
		width, height = "500px", "616px"
	} else if name == "twitter" {
		width, height = "480px", "480px"
	}

	for _, item := range compat_split_macro_args(data) {
		parts := strings.SplitN(item, "=", 2)
		if len(parts) == 2 {
			switch strings.ToLower(strings.TrimSpace(parts[0])) {
			case "width":
				width = compat_middle_style(compat_file_px(parts[1]))
			case "height":
				height = compat_middle_style(compat_file_px(parts[1]))
			case "start":
				start = parts[1]
			case "end":
				end = parts[1]
			case "https://www.youtube.com/watch?v":
				if name == "youtube" {
					code = parts[1]
				}
			}
		} else if code == "" {
			code = strings.TrimSpace(item)
		}
	}
	if code == "" || width == "" || height == "" {
		return ""
	}

	valid_code := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	switch name {
	case "youtube":
		if strings.HasPrefix(code, "https://youtu.be/") {
			code = strings.TrimPrefix(code, "https://youtu.be/")
		} else if strings.Contains(code, "youtube.com/watch?v=") {
			code = strings.SplitN(strings.SplitN(code, "youtube.com/watch?v=", 2)[1], "&", 2)[0]
		} else if strings.Contains(code, "youtube.com/embed/") {
			code = strings.TrimPrefix(strings.SplitN(code, "?", 2)[0], "https://www.youtube.com/embed/")
		}
		if !valid_code.MatchString(code) {
			return ""
		}
		src := "https://www.youtube.com/embed/" + code
		query := []string{}
		if regexp.MustCompile(`^[0-9]+$`).MatchString(start) {
			query = append(query, "start="+start)
		}
		if regexp.MustCompile(`^[0-9]+$`).MatchString(end) {
			query = append(query, "end="+end)
		}
		if len(query) > 0 {
			src += "?" + strings.Join(query, "&")
		}
		return compat_media_iframe(src, width, height, "YouTube")
	case "instagram":
		code = strings.TrimSuffix(strings.TrimPrefix(code, "https://www.instagram.com/p/"), "/")
		if !valid_code.MatchString(code) {
			return ""
		}
		return compat_media_iframe("https://www.instagram.com/p/"+code+"/embed/", width, height, "Instagram")
	case "facebook":
		if !compat_is_external_link(code) {
			return ""
		}
		src := "https://www.facebook.com/plugins/post.php?href=" + url.QueryEscape(code) + "&width=" + url.QueryEscape(width) + "&height=" + url.QueryEscape(height)
		return compat_media_iframe(src, width, height, "Facebook")
	case "tiktok":
		code = strings.TrimPrefix(code, "https://www.tiktok.com/@")
		if strings.Contains(code, "/video/") {
			code = strings.SplitN(code, "/video/", 2)[1]
		}
		if !valid_code.MatchString(code) {
			return ""
		}
		return compat_media_iframe("https://www.tiktok.com/embed/v2/"+code, width, height, "TikTok")
	case "kakaotv":
		code = strings.TrimPrefix(code, "https://tv.kakao.com/v/")
		if !valid_code.MatchString(code) {
			return ""
		}
		return compat_media_iframe("https://tv.kakao.com/embed/player/cliplink/"+code+"?service=kakao_tv", width, height, "KakaoTV")
	case "navertv":
		code = strings.TrimPrefix(code, "https://tv.naver.com/v/")
		if !valid_code.MatchString(code) {
			return ""
		}
		return compat_media_iframe("https://tv.naver.com/embed/"+code, width, height, "NaverTV")
	case "nicovideo":
		code = strings.TrimPrefix(code, "https://www.nicovideo.jp/watch/")
		if !valid_code.MatchString(code) {
			return ""
		}
		return compat_media_iframe("https://embed.nicovideo.jp/watch/"+code, width, height, "Niconico")
	case "vimeo":
		code = strings.TrimPrefix(code, "https://vimeo.com/")
		if !valid_code.MatchString(code) {
			return ""
		}
		return compat_media_iframe("https://player.vimeo.com/video/"+code, width, height, "Vimeo")
	case "twitter":
		media_url, err := url.Parse(code)
		if err != nil || (strings.ToLower(media_url.Scheme) != "http" && strings.ToLower(media_url.Scheme) != "https") {
			return ""
		}
		media_host := strings.ToLower(media_url.Hostname())
		if media_host != "twitter.com" && media_host != "www.twitter.com" && media_host != "x.com" && media_host != "www.x.com" {
			return ""
		}
		src := "https://twitframe.com/show?url=" + url.QueryEscape(code)
		if class.get_render_setting("main_css_darkmode") == "1" {
			src += "&theme=dark"
		}
		return compat_media_iframe(src, width, height, "Twitter")
	}
	return ""
}

func compat_media_iframe(src string, width string, height string, title string) string {
	return `<iframe title="` + compat_html_escape(title) + `" style="width:` + compat_html_escape(width) + `;height:` + compat_html_escape(height) + `;" src="` + compat_html_escape(src) + `" frameborder="0" allowfullscreen loading="lazy"></iframe>`
}

func (class *namumark_compat_renderer) render_quote(lines []string) string {
	quote_lines := []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "&gt;")
		quote_lines = append(quote_lines, strings.TrimSpace(line))
	}
	quote_data := class.process_blocks(strings.Join(quote_lines, "\n"))
	quote_data = class.render_inline(quote_data)
	quote_data = strings.ReplaceAll(quote_data, "\n", "<br>")
	return `<hr class="mini_hr"><blockquote><div>` + quote_data + `</div></blockquote><hr class="mini_hr">`
}

func (class *namumark_compat_renderer) toc_html() string {
	if len(class.toc_items) == 0 {
		return ""
	}
	toc_title := "toc"
	if class.db != nil {
		toc_title = tool.Get_language(class.db, "toc", true)
	}
	data := `<div class="opennamu_TOC" id="toc"><span class="opennamu_TOC_title">` + compat_html_escape(toc_title) + `</span><br>`
	for _, item := range class.toc_items {
		indent := strings.Count(item.number, ".")
		data += `<br>` + strings.Repeat(`<span style="margin-left: 10px;"></span>`, indent) + `<span class="opennamu_TOC_list"><a href="#s-` + compat_html_escape(item.number) + `">` + item.number + `. </a>` + item.text + `</span>`
	}
	return data + `</div>`
}

func (class *namumark_compat_renderer) process_headings(data string) string {
	lines := strings.Split(data, "\n")
	heading_stack := [6]int{}
	toc_requested := strings.Contains(data, "[toc]")
	first_heading_index := -1
	for index, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		level := 0
		for level < len(line) && level < 6 && line[level] == '=' {
			level++
		}
		if level == 0 || level == len(line) {
			continue
		}

		heading_data := strings.TrimSpace(line[level:])
		heading_folding := strings.HasPrefix(heading_data, "#")
		if heading_folding {
			heading_data = strings.TrimSpace(heading_data[1:])
		}
		trailing := 0
		for trailing < len(heading_data) && trailing < 6 && heading_data[len(heading_data)-trailing-1] == '=' {
			trailing++
		}
		if trailing > 0 {
			heading_data = strings.TrimSpace(heading_data[:len(heading_data)-trailing])
		}
		trailing_folding := strings.HasSuffix(heading_data, "#")
		if heading_folding != trailing_folding {
			continue
		}
		if trailing_folding {
			heading_data = strings.TrimSpace(strings.TrimSuffix(heading_data, "#"))
		}
		if heading_data == "" {
			continue
		}

		if first_heading_index < 0 {
			first_heading_index = index
		}
		heading_stack[level-1]++
		for stack_index := level; stack_index < len(heading_stack); stack_index++ {
			heading_stack[stack_index] = 0
		}
		number_list := []string{}
		for stack_index := 0; stack_index < level; stack_index++ {
			if heading_stack[stack_index] > 0 {
				number_list = append(number_list, strconv.Itoa(heading_stack[stack_index]))
			}
		}
		number := strings.Join(number_list, ".")
		rendered_heading := class.render_inline(heading_data)
		class.toc_items = append(class.toc_items, namumark_compat_toc_item{number, rendered_heading})
		heading_id := regexp.MustCompile(`<[^<>]*>`).ReplaceAllString(class.restore(rendered_heading), "")
		heading_id = tool.HTML_unescape(heading_id)
		lines[index] = class.reserve(`<h` + strconv.Itoa(level) + ` id="` + compat_html_escape(heading_id) + `"><a href="#toc" id="s-` + compat_html_escape(number) + `">` + number + `.</a> ` + rendered_heading + `</h` + strconv.Itoa(level) + `>`)
	}

	toc_set := class.get_render_setting("main_css_toc_set")
	data = strings.Join(lines, "\n")
	if toc_requested {
		if toc_set == "off" {
			data = strings.ReplaceAll(data, "[toc]", "")
		} else {
			data = strings.ReplaceAll(data, "[toc]", class.reserve(class.toc_html()))
		}
	} else if class.is_view() && toc_set != "off" && toc_set != "half_off" && first_heading_index >= 0 {
		toc_data := class.reserve(class.toc_html())
		toc_lines := []string{toc_data}
		previous_is_list := first_heading_index > 0 && namumark_compat_list_regex.MatchString(strings.TrimSpace(lines[first_heading_index-1]))
		if first_heading_index == 0 || previous_is_list {
			toc_lines = []string{"", toc_data}
		}
		lines = append(lines[:first_heading_index], append(toc_lines, lines[first_heading_index:]...)...)
		data = strings.Join(lines, "\n")
	}
	return data
}
func compat_remove_comments(data string) string {
	lines := strings.Split(data, "\n")
	result := []string{lines[0]}
	for index := 1; index < len(lines); index++ {
		if strings.HasPrefix(lines[index], "##") && len(lines[index]) > 2 {
			continue
		}
		result = append(result, lines[index])
	}
	return strings.Join(result, "\n")
}

func (class *namumark_compat_renderer) process_blocks(data string) string {
	lines := strings.Split(data, "\n")
	var result strings.Builder
	has_result := false
	last_was_block := false
	last_was_break := false
	append_result := func(value string, is_block bool) {
		if has_result && !last_was_block && !is_block && !last_was_break {
			result.WriteByte('\n')
		}
		result.WriteString(value)
		has_result = true
		last_was_block = is_block
		last_was_break = strings.HasSuffix(value, "\n")
		if token_data, token_exists := class.tokens[strings.TrimSpace(value)]; token_exists {
			last_was_break = last_was_break || strings.HasSuffix(strings.TrimSpace(token_data), "<br>") || strings.HasSuffix(token_data, "\n")
		}
	}
	for index := 0; index < len(lines); index++ {
		line := lines[index]
		trimmed := strings.TrimSpace(line)
		token_data, token_exists := class.tokens[trimmed]
		if trimmed == "" && last_was_break {
			continue
		}
		if token_exists && strings.HasPrefix(token_data, "<h") {
			append_result(line, true)
			continue
		}
		if token_exists && strings.HasPrefix(token_data, `<div class="opennamu_TOC"`) {
			append_result(line, false)
			continue
		}

		if strings.HasPrefix(line, "##") && len(line) > 2 {
			continue
		}
		if regexp.MustCompile(`^-{4,9}$`).MatchString(trimmed) {
			append_result(class.reserve("<hr>"), true)
			continue
		}

		table_caption := false
		if strings.HasPrefix(trimmed, "|") && !strings.HasPrefix(trimmed, "||") {
			caption_end := strings.Index(trimmed[1:], "|")
			if caption_end >= 0 {
				caption_end++
				table_caption = strings.Contains(trimmed[caption_end+1:], "||")
			}
		}
		if strings.HasPrefix(trimmed, "||") || table_caption {
			table_lines := []string{line}
			table_row_open := !strings.HasSuffix(trimmed, "||")
			if trimmed == "||" && index+1 < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[index+1]), "||") {
				table_row_open = true
			}
			for index+1 < len(lines) {
				next := strings.TrimSpace(lines[index+1])
				if !table_row_open && !strings.HasPrefix(next, "||") {
					break
				}
				index++
				if table_row_open {
					table_lines[len(table_lines)-1] += "\n" + lines[index]
				} else {
					table_lines = append(table_lines, lines[index])
				}
				table_row_open = !strings.HasSuffix(next, "||")
			}
			table_data := class.render_table(table_lines)
			table_suffix := ""
			if index+1 < len(lines) {
				next_line := strings.TrimSpace(lines[index+1])
				is_hr := len(next_line) >= 4 && len(next_line) <= 9 && strings.Trim(next_line, "-") == ""
				if next_line != "" && !is_hr && !namumark_compat_list_regex.MatchString(next_line) {
					table_suffix = "\n"
				}
			}
			if has_result && !last_was_block {
				result.WriteByte(10)
			}
			append_result(class.reserve(table_data)+table_suffix, true)
			continue
		}

		if namumark_compat_list_regex.MatchString(line) {
			first_match := namumark_compat_list_regex.FindStringSubmatch(line)
			ordered := len(first_match) >= 3 && !strings.HasPrefix(first_match[2], "*")
			list_lines := []string{line}
			for index+1 < len(lines) {
				next_match := namumark_compat_list_regex.FindStringSubmatch(lines[index+1])
				if len(next_match) < 3 {
					break
				}
				next_ordered := !strings.HasPrefix(next_match[2], "*")
				if next_ordered != ordered {
					break
				}
				index++
				list_lines = append(list_lines, lines[index])
			}
			if ordered && len(list_lines) < 2 {
				append_result(line, false)
				continue
			}
			append_result(class.reserve(class.render_list(list_lines)), true)
			continue
		}

		if strings.HasPrefix(trimmed, "&gt;") {
			quote_lines := []string{line}
			for index+1 < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[index+1]), "&gt;") {
				index++
				quote_lines = append(quote_lines, lines[index])
			}
			append_result(class.reserve(class.render_quote(quote_lines)), true)
			continue
		}

		append_result(line, false)
	}
	return result.String()
}

func compat_date_now() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.UTC)
}

func compat_parse_date(data string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", strings.TrimSpace(data), time.UTC)
}

func compat_date_days(now time.Time, parsed time.Time) int {
	duration := now.Sub(parsed)
	days := int(duration / (24 * time.Hour))
	if duration < 0 && duration%(24*time.Hour) != 0 {
		days--
	}
	return days
}

func compat_date_component(days int, month bool) int {
	sign := 1
	if days <= 0 {
		sign = -1
		days = -days
	}
	date_value := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, days)
	if month {
		return sign * ((date_value.Year()-1)*12 + int(date_value.Month()) - 1)
	}
	return sign * (date_value.Year() - 1)
}

func compat_macro_date(data string, macro_name string) string {
	parsed, err := compat_parse_date(data)
	if err != nil {
		return "invalid date"
	}
	now := compat_date_now()
	days := compat_date_days(now, parsed)
	switch macro_name {
	case "age":
		if parsed.After(now) {
			return "invalid date"
		}
		return strconv.Itoa(days / 365)
	case "dday":
		date_value := days
		if date_value > 0 {
			return "+" + strconv.Itoa(date_value)
		}
		if date_value == 0 {
			return "-0"
		}
		return strconv.Itoa(date_value)
	case "dmonth":
		date_value := compat_date_component(days, true)
		if date_value > 0 {
			return "+" + strconv.Itoa(date_value)
		}
		if date_value == 0 {
			return "-0"
		}
		return strconv.Itoa(date_value)
	case "dyear":
		date_value := compat_date_component(days, false)
		if date_value > 0 {
			return "+" + strconv.Itoa(date_value)
		}
		if date_value == 0 {
			return "-0"
		}
		return strconv.Itoa(date_value)
	}
	return "invalid date"
}

func (class *namumark_compat_renderer) macro_lastedit(data string) string {
	parts := compat_split_macro_args(data)
	target := ""
	full := false
	for _, part := range parts {
		option := strings.SplitN(part, "=", 2)
		if len(option) == 2 {
			if strings.EqualFold(strings.TrimSpace(option[0]), "view") && strings.EqualFold(strings.TrimSpace(option[1]), "full") {
				full = true
			}
			continue
		}
		if target == "" {
			target = strings.TrimSpace(part)
		}
	}
	if class.db == nil || target == "" {
		return "0"
	}
	target = class.normalize_target(tool.HTML_unescape(target))
	date_data := ""
	tool.QueryRow_DB(class.db, "select set_data from data_set where doc_name = ? and set_name = 'last_edit'", []any{&date_data}, target)
	if date_data == "" {
		return "0"
	}
	if !full && strings.Contains(date_data, " ") {
		return strings.SplitN(date_data, " ", 2)[0]
	}
	return date_data
}

func (class *namumark_compat_renderer) process_macro_double(name string, data string, raw string) string {
	name = strings.ToLower(name)
	switch name {
	case "youtube", "nicovideo", "navertv", "kakaotv", "vimeo", "instagram", "twitter", "tiktok", "facebook":
		if class.collect_only {
			return ""
		}
		if media_data := class.compat_media_macro(name, data); media_data != "" {
			return class.reserve(media_data)
		}
		return raw
	case "toc":
		return "[toc()]"
	case "pagecount":
		return "0"
	case "joke":
		if class.collect_only {
			return data
		}
		return class.reserve(`<span class="opennamu_joke">`) + data + class.reserve("</span>")
	case "anchor":
		anchor := tool.HTML_unescape(strings.TrimSpace(compat_split_macro_args(data)[0]))
		if anchor == "" {
			return ""
		}
		if class.collect_only {
			return ""
		}
		return class.reserve(`<span id="`+compat_html_escape(anchor)+`">`) + class.reserve("</span>")
	case "age", "dday", "dmonth", "dyear":
		return compat_macro_date(data, name)
	case "timeif":
		main_data := ""
		before_data := ""
		after_data := ""
		for _, item := range compat_split_macro_args(data) {
			parts := strings.SplitN(item, "=", 2)
			if len(parts) == 2 {
				switch strings.ToLower(strings.TrimSpace(parts[0])) {
				case "before":
					before_data = parts[1]
				case "after":
					after_data = parts[1]
				}
			} else if main_data == "" {
				main_data = strings.TrimSpace(item)
			}
		}
		parsed, err := compat_parse_date(main_data)
		if err != nil {
			return "invalid date"
		}
		if parsed.After(compat_date_now()) {
			return before_data
		}
		return after_data
	case "ruby":
		main_text := ""
		ruby_text := ""
		color := ""
		for _, item := range compat_split_macro_args(data) {
			parts := strings.SplitN(item, "=", 2)
			if len(parts) == 2 {
				switch strings.ToLower(strings.TrimSpace(parts[0])) {
				case "ruby":
					ruby_text = parts[1]
				case "color":
					color = compat_middle_style(parts[1])
					if color != "" {
						color += ";"
					}
				}
			} else {
				main_text = item
			}
		}
		if class.collect_only {
			return main_text + ruby_text
		}
		ruby_data := class.reserve("<ruby>") + main_text + class.reserve("<rp>(</rp><rt>")
		if color != "" {
			ruby_data += class.reserve(`<span style="color:`+compat_html_escape(tool.HTML_unescape(color))+`">`) + ruby_text + class.reserve("</span>")
		} else {
			ruby_data += ruby_text
		}
		return ruby_data + class.reserve("</rt><rp>)</rp></ruby>")
	case "username":
		if strings.TrimSpace(data) == "" {
			return raw
		}
		username := ""
		render := true
		load_name := false
		for _, item := range compat_split_macro_args(data) {
			parts := strings.SplitN(item, "=", 2)
			if len(parts) == 2 {
				switch strings.ToLower(strings.TrimSpace(parts[0])) {
				case "render":
					render = strings.TrimSpace(parts[1]) != "0"
				case "load_name":
					load_name = strings.TrimSpace(parts[1]) == "1"
				}
			} else if username == "" {
				username = item
			}
		}
		if load_name {
			username = compat_render_parameter_value(class.parameter["ip"])
		}
		if class.collect_only {
			return username
		}
		class_name := ""
		if render {
			class_name = "opennamu_render_ip"
		}
		return class.reserve(`<span class="`+class_name+`">`) + compat_escape_value(username) + class.reserve("</span>")
	case "lastedit":
		return class.macro_lastedit(data)
	default:
		return raw
	}
}

func (class *namumark_compat_renderer) process_macros(data string) string {
	data = compat_replace_regex2(data, `(?is)(?<!\[)\[(?!\*)([^[(\]]+)\(((?:(?!\[\[|\)\]).)+)\)\]`, func(match regexp2.Match) string {
		return class.process_macro_double(match.GroupByNumber(1).String(), match.GroupByNumber(2).String(), match.String())
	})
	return namumark_compat_single_macro_regex.ReplaceAllStringFunc(data, func(raw string) string {
		match := namumark_compat_single_macro_regex.FindStringSubmatch(raw)
		if len(match) < 2 {
			return raw
		}
		switch strings.ToLower(match[1]) {
		case "br":
			if class.collect_only {
				return ""
			}
			return class.reserve("<br>")
		case "clearfix":
			if class.collect_only {
				return ""
			}
			return class.reserve(`<div style="clear: both;"></div>`)
		case "date", "datetime":
			return tool.Get_time()
		case "pagecount":
			count := "0"
			if class.db != nil {
				tool.QueryRow_DB(class.db, "select data from other where name = 'count_all_title'", []any{&count})
			}
			return count
		case "toc", "목차", "tableofcontents":
			return "[toc]"
		default:
			return raw
		}
	})
}

func (class *namumark_compat_renderer) process_slash(data string) string {
	return compat_replace_regex2(data, `\\(&lt;|&gt;|&#x27;|&quot;|&amp;|.)`, func(match regexp2.Match) string {
		return class.reserve_slash(match.GroupByNumber(1).String())
	})
}

func (class *namumark_compat_renderer) prepare() {
	class.data = compat_render_parameter_data(class.data, class.parameter)
	class.data = compat_remove_comments(class.data)
	class.data = class.process_slash(class.data)
	class.data = class.process_redirect(class.data)
	class.data = class.process_middle(class.data)
	class.data = class.process_macros(class.data)
	class.data = class.process_links(class.data)
	class.data = class.process_includes(class.data)
	class.data = class.process_math(class.data)
	if class.render_type != "include" && class.render_type != "inter" {
		class.data = class.process_footnotes(class.data)
	}
	if !class.collect_only {
		if class.render_type != "include" && class.render_type != "inter" {
			class.data = class.process_headings(class.data)
		}
		class.data = class.render_text(class.data)
		class.data = class.process_blocks(class.data)
		if class.footnote_token != "" {
			class.data += class.footnote_token
		}
	}
}

func (class *namumark_compat_renderer) category_html() string {
	if len(class.categories) == 0 {
		return ""
	}
	data := `<div class="opennamu_category" id="cate">` + tool.Get_language(class.db, "category", true) + " : "
	for index, category := range class.categories {
		if index > 0 {
			data += " | "
		}
		category_classes := []string{}
		if category.blur {
			category_classes = append(category_classes, "opennamu_category_blur")
		}
		if !category.exists {
			category_classes = append(category_classes, "opennamu_not_exist_link")
		}
		class_name := ""
		if len(category_classes) > 0 {
			class_name = ` class="` + strings.Join(category_classes, " ") + `"`
		}
		label := compat_html_escape(category.label)
		category_url := class.compat_url_parser(category.target)
		category_parts := strings.SplitN(category.target, ":", 2)
		if len(category_parts) == 2 && strings.EqualFold(category_parts[0], "category") {
			category_url = "category:" + class.compat_url_parser(category_parts[1])
		}
		data += `<a` + class_name + ` title="` + label + `" href="/w/` + category_url + `">` + label + `</a>`
	}
	return data + "</div>"
}

func (class *namumark_compat_renderer) is_view() bool {
	return class.include_depth == 0 && (class.render_type == "normal" || class.render_type == "view" || class.render_type == "from")
}

func (class *namumark_compat_renderer) apply_text_setting(data string) string {
	switch class.get_render_setting("main_css_bold") {
	case "delete":
		data = namumark_compat_bold_regex.ReplaceAllString(data, "")
	case "change":
		data = namumark_compat_bold_regex.ReplaceAllString(data, "$1")
	}
	switch class.get_render_setting("main_css_strike") {
	case "delete":
		data = namumark_compat_strike_regex.ReplaceAllString(data, "")
	case "change":
		data = namumark_compat_strike_regex.ReplaceAllString(data, "$1")
	}
	return data
}

func compat_remove_nested_links(data string) string {
	link_depth := 0
	link_regex := regexp.MustCompile(`(?i)(<a(?: [^<>]*)?>|</a>)`)
	return link_regex.ReplaceAllStringFunc(data, func(value string) string {
		if strings.HasPrefix(strings.ToLower(value), "</a>") {
			if link_depth == 0 {
				return ""
			}
			link_depth--
			if link_depth > 0 {
				return ""
			}
			return value
		}
		if link_depth > 0 {
			link_depth++
			return ""
		}
		link_depth = 1
		return value
	})
}

func (class *namumark_compat_renderer) render_output() string {
	data := class.data
	data = strings.ReplaceAll(data, "\n", "<br>")
	data = class.restore(data)
	data = compat_remove_nested_links(data)
	if class.is_view() && len(class.categories) > 0 {
		data = strings.TrimSuffix(data, "<br>")
	}
	if class.is_view() {
		if category_data := class.category_html(); category_data != "" {
			category_separator := "<hr class=\"main_hr\">"
			if class.get_render_setting("main_css_category_set") == "bottom" {
				category_separator = "<hr>"
			}
			category_data = category_separator + category_data
			footnote_index := strings.Index(data, "<div class=\"opennamu_footnote\">")
			if footnote_index >= 0 {
				data = strings.TrimSuffix(data[:footnote_index], "<br>") + category_data + data[footnote_index:]
			} else {
				data += category_data
			}
		}
	}
	return data
}

func (class *namumark_compat_renderer) result() map[string]any {
	class.prepare()
	if class.collect_only {
		return map[string]any{
			"data":       "",
			"js_data":    "",
			"backlinks":  class.backlink_entries(),
			"link_count": class.link_count,
			"redirect":   class.redirect,
		}
	}
	return map[string]any{
		"data":       class.render_output(),
		"js_data":    "",
		"backlinks":  class.backlink_entries(),
		"link_count": class.link_count,
		"redirect":   class.redirect,
	}
}

func (class *namumark_compat_renderer) backlink_entries() []namumark_compat_backlink {
	result := []namumark_compat_backlink{}
	for _, key := range class.backlink_order {
		result = append(result, class.backlinks[key])
	}
	return result
}

func render_namumark_compat(
	db *sql.DB,
	doc_name string,
	data string,
	render_type string,
	parameter map[string]any,
	include_depth int,
	collect_only bool,
) map[string]any {
	class := new_namumark_compat_renderer(db, doc_name, data, render_type, parameter, include_depth, collect_only)
	return class.result()
}

func render_namumark_compat_backlink(db *sql.DB, doc_name string, data string) map[string]string {
	result := render_namumark_compat(db, doc_name, data, "backlink", nil, 0, true)
	entries, _ := result["backlinks"].([]namumark_compat_backlink)
	link_count, _ := result["link_count"].(int)
	redirect, _ := result["redirect"].(bool)

	tool.Exec_DB(db, "delete from back where link = ?", doc_name)
	tool.Exec_DB(db, "delete from back where title = ? and type = 'no'", doc_name)
	tool.Exec_DB(db, "delete from data_set where doc_name = ? and set_name in ('link_count', 'doc_type')", doc_name)
	for _, entry := range entries {
		tool.Exec_DB(db, "insert into back (link, title, type, data) values (?, ?, ?, ?)", doc_name, entry.target, entry.link_type, entry.data)
	}
	tool.Exec_DB(db, "insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', 'link_count', ?)", doc_name, link_count)

	doc_type := ""
	if strings.HasPrefix(doc_name, "user:") {
		doc_type = "user"
	} else if strings.HasPrefix(doc_name, "file:") {
		doc_type = "file"
	} else if strings.HasPrefix(doc_name, "category:") {
		doc_type = "category"
	} else if redirect {
		doc_type = "redirect"
	}
	tool.Exec_DB(db, "insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', 'doc_type', ?)", doc_name, doc_type)

	return map[string]string{"data": `<div class="opennamu_render_complete"></div>`, "js_data": ""}
}
func compat_url_parser(data string) string {
	data = strings.ReplaceAll(data, "OPENNAMU_COMPAT_TOKEN_", "OPENNAMU_COMPAT_T0KEN_")
	if strings.HasPrefix(data, ".") {
		data = "\\" + data
	}
	return strings.ReplaceAll(url.QueryEscape(data), "+", "%20")
}
func (class *namumark_compat_renderer) compat_url_parser(data string) string {
	return compat_url_parser(class.restore_slash(data))
}

func compat_legacy_html_literal(data string) string {
	data = strings.ReplaceAll(data, "&amp;amp;nbsp;", "&amp;nbsp;")
	return data
}

func compat_html_escape(data string) string {
	data = strings.ReplaceAll(data, "&", "&amp;")
	data = strings.ReplaceAll(data, "<", "&lt;")
	data = strings.ReplaceAll(data, ">", "&gt;")
	data = strings.ReplaceAll(data, `"`, "&quot;")
	return strings.ReplaceAll(data, "'", "&#x27;")
}
func compat_escape_value(data string) string {
	return compat_html_escape(tool.HTML_unescape(data))
}

func compat_replace_regex2(data string, pattern string, fn func(regexp2.Match) string) string {
	regex := regexp2.MustCompile(pattern, 0)
	result, err := regex.ReplaceFunc(data, fn, -1, -1)
	if err != nil {
		return data
	}
	return result
}

func (class *namumark_compat_renderer) render_text(data string) string {
	text_data := []struct {
		pattern string
		open    string
		close   string
		setting string
	}{
		{`&#x27;&#x27;&#x27;((?:(?!&#x27;&#x27;&#x27;).)+)&#x27;&#x27;&#x27;`, "<b>", "</b>", "main_css_bold"},
		{`&#x27;&#x27;((?:(?!&#x27;&#x27;).)+)&#x27;&#x27;`, "<i>", "</i>", ""},
		{`__((?:(?!__).)+)__`, "<u>", "</u>", ""},
		{`\^\^\^((?:(?!\^\^\^).)+)\^\^\^`, "<sup>", "</sup>", ""},
		{`\^\^((?:(?!\^\^).)+)\^\^`, "<sup>", "</sup>", ""},
		{`,,,((?:(?!,,,).)+),,,`, "<sub>", "</sub>", ""},
		{`,,((?:(?!,,).)+),,`, "<sub>", "</sub>", ""},
		{`--((?:(?!--).)+)--`, "<s>", "</s>", "main_css_strike"},
		{`~~((?:(?!~~).)+)~~`, "<s>", "</s>", "main_css_strike"},
	}

	for _, item := range text_data {
		setting := class.get_render_setting(item.setting)
		data = compat_replace_regex2(data, item.pattern, func(match regexp2.Match) string {
			inside := match.GroupByNumber(1).String()
			if item.setting != "" && setting == "delete" {
				return ""
			}
			if item.setting != "" && setting == "change" {
				return inside
			}
			return item.open + inside + item.close
		})
	}
	return data
}
