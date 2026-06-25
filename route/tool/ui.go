package tool

import (
	"database/sql"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/flosch/pongo2/v6"
)

func Get_skin_route(skin_name string, route string) string {
    return filepath.Join("..", "views", skin_name, route)
}

func Get_template_set(skin_name string) map[string]string {
    set_file_path := Get_skin_route(skin_name, "set.json")
    if _, err := os.Stat(set_file_path); err == nil {
        data, err := os.ReadFile(set_file_path)
        if err != nil {
            panic(err)
        }

        set_json := map[string]string{}
        json.Unmarshal([]byte(data), &set_json)

        return set_json
    }

    return map[string]string{}
}

func Get_use_skin_name(db *sql.DB, ip string) string {
    skin_list := Get_skin_list("ringo", true)
    skin := skin_list[0]

    user_skin_name := ""
    if !IP_or_user(ip) {
        QueryRow_DB(
            db,
            "select data from user_set where name = 'skin' and id = ?",
            []any{ &user_skin_name },
            ip,
        )
    }

    if user_skin_name == "default" {
        user_skin_name = ""
    }

    if user_skin_name == "" {
        QueryRow_DB(
            db,
            "select data from other where name = 'skin'",
            []any{ &user_skin_name },
        )
    }

    if user_skin_name != "" && Arr_in_str(skin_list, user_skin_name) {
        skin = user_skin_name
    }

    return skin
}

func Get_template(db *sql.DB, config Config, name string, data string, other []any, menu [][]any, option map[string]string) string {
    skin_name := Get_use_skin_name(db, config.IP)
    
    template_set := Get_template_set(skin_name)
    for k, v := range template_set {
        data = strings.ReplaceAll(data, k, v)
    }

    menu_func := func(menu [][]any) any {
        if len(menu) == 0 {
            return 0
        } else {
            return menu
        }
    }
    menu_func_result := menu_func(menu)

    if len(other) < 1 {
        other = append(other, 0)
    }

    if len(other) < 2 {
        other = append(other, 0)
    }

    for k := range other {
        switch v := other[k].(type) {
            case nil:
                other[k] = 0
            case float64:
                other[k] = int(v)
            case int64:
                other[k] = int(v)
            case int:
                other[k] = v
            case string:
                if v == "" {
                    other[k] = 0
                } else {
                    other[k] = v
                }
            default:
                other[k] = 0
        }   
    }

    imp_1 := Get_wiki_set(db, config.IP, config.Cookies)
    imp_2 := Get_wiki_custom(db, config.IP, config.Session, config.Cookies)
    imp_3 := Get_wiki_css(other, config.Cookies)

    if len(imp_3) < 8 {
        imp_3 = append(imp_3, 0)
    }

    added_menu := []string{}
    switch imp_1[7].(type) {
        case []string:
            added_menu = imp_1[7].([]string)
        default:
            added_menu = []string{"", "", ""}
    }

    if len(added_menu) < 3 {
        for i := len(added_menu); i < 3; i++ {
            added_menu = append(added_menu, "")
        }
    }
    
    imp_1[7] = added_menu

    path := ""
    if option["path"] != "" {
        path = option["path"]
    }

    doc_length := ""
    if _, ok := option["length_doc"]; ok {
        doc_length = option["length_doc"]
    }

    context := pongo2.Context{
        "imp" : []any{
            name,
            imp_1,
            imp_2,
            imp_3,
        },
        "data" : `<div class="opennamu_main">` + data + `</div>`,
        "menu" : menu_func_result,

        "title" : name,
        
        "wiki_name" : imp_1[0],
        "license" : imp_1[1],
        "wiki_logo" : imp_1[4],
        "global_head" : imp_1[5],
        "add_menu" : imp_1[6],
        "template_var_1" : added_menu[0],
        "template_var_2" : added_menu[1],
        "template_var_3" : added_menu[2],

        "user_login" : imp_2[2],
        "user_head" : imp_2[3],
        "user_email" : imp_2[4],
        "user_name" : imp_2[5],
        "user_is_admin" : imp_2[6],
        "user_is_ban" : imp_2[7],
        "user_alarm_count" : imp_2[8],
        "user_auth" : imp_2[9],
        "user_ip" : imp_2[10],
        "user_discuss" : imp_2[11],
        "user_path" : path,
        "user_level" : imp_2[13],

        "sub_title" : imp_3[0],
        "last_edit" : imp_3[1],
        "main_head" : imp_3[3],
        "star_doc" : imp_3[4],
        "main_head_dark" : imp_3[5],
        "description_doc" : imp_3[6],
        "view_count" : imp_3[7],
        "length_doc" : doc_length,
    }

    tpl, err := pongo2.FromFile(Get_skin_route(skin_name, "index.html"))
    if err != nil {
        panic(err)
    }

    out, err := tpl.Execute(context)
    if err != nil {
        panic(err)
    }

    return out
}

func Get_redirect(target string) string {
    attrURL := html.EscapeString(target)
    jsURL := strconv.Quote(target)

    return fmt.Sprintf(`<!doctype html>
<html lang="ko">
<head>
<meta charset="utf-8">
<title>Redirecting…</title>
<script>
location.replace(%s);
</script>
<noscript>
<meta http-equiv="refresh" content="0; url=%s">
</noscript>
</head>
<body>
<p>Redirecting… <a href="%s">continue</a></p>
</body>
</html>`, jsURL, attrURL, attrURL)
}

func Cache_v() string {
    return ".cache_v288"
}

func Get_wiki_css(data []any, cookies string) []any {
    for len(data) < 4 {
        data = append(data, "")
    }

    data_css := ""
    data_css_dark := ""

    data_css_ver := Cache_v()

    // Cache Control
    data_css += `<meta http-equiv="Cache-Control" content="max-age=31536000">`

    // External JS
    data_css += `<script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.38/dist/katex.min.js" integrity="sha384-H6s1ZrH2CKpFpqR680poRdStIRJGXty7fSkxAcIfxwl9iu6A4BOPtTk7vQ58Ovio" crossorigin="anonymous"></script>`
    data_css += `<script defer src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.8.0/highlight.min.js" integrity="sha512-rdhY3cbXURo13l/WU9VlaRyaIYeJ/KBakckXIvJNAQde8DgpOmE+eZf7ha4vdqVjTtwQt69bD2wH2LXob/LB7Q==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>`
    data_css += `<script defer src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.8.0/languages/x86asm.min.js" integrity="sha512-HeAchnWb+wLjUb2njWKqEXNTDlcd1QcyOVxb+Mc9X0bWY0U5yNHiY5hTRUt/0twG8NEZn60P3jttqBvla/i2gA==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>`
    data_css += `<script defer src="https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.48.0/min/vs/loader.min.js" integrity="sha512-ZG31AN9z/CQD1YDDAK4RUAvogwbJHv6bHrumrnMLzdCrVu4HeAqrUX7Jsal/cbUwXGfaMUNmQU04tQ8XXl5Znw==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>`
    data_css += `<script defer src="https://cdnjs.cloudflare.com/ajax/libs/highlightjs-line-numbers.js/2.8.0/highlightjs-line-numbers.min.js"></script>`
    data_css += `<script defer src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>`

    // Func JS
    data_css += `<script defer src="/views/main_css/js/func/func.js` + data_css_ver + `"></script>`
    data_css += `<script defer src="/views/main_css/js/func/insert_version.js` + data_css_ver + `"></script>`
    data_css += `<script defer src="/views/main_css/js/func/insert_user_info.js` + data_css_ver + `"></script>`
    data_css += `<script defer src="/views/main_css/js/func/insert_version_skin.js` + data_css_ver + `"></script>`
    data_css += `<script defer src="/views/main_css/js/func/insert_http_warning_text.js` + data_css_ver + `"></script>`
    data_css += `<script defer src="/views/main_css/js/func/ie_end_of_life.js` + data_css_ver + `"></script>`
    data_css += `<script defer src="/views/main_css/js/func/shortcut.js` + data_css_ver + `"></script>`
    data_css += `<script defer src="/views/main_css/js/func/editor.js` + data_css_ver + `"></script>`
    data_css += `<script defer src="/views/main_css/js/func/render.js` + data_css_ver + `"></script>`

    // Main CSS
    data_css += `<link rel="stylesheet" href="/views/main_css/css/main.css` + data_css_ver + `">`

    // External CSS
    data_css += `<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/katex@0.16.38/dist/katex.min.css" integrity="sha384-/L6i+LN3dyoaK2jYG5ZLh5u13cjdsPDcFOSNJeFBFa/KgVXR5kOfTdiN3ft1uMAq" crossorigin="anonymous">`
    data_css += `<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.8.0/styles/default.min.css" integrity="sha512-hasIneQUHlh06VNBe7f6ZcHmeRTLIaQWFd43YriJ0UND19bvYRauxthDg8E4eVNPm9bRUhr5JGeqH7FRFXQu5g==" crossorigin="anonymous" referrerpolicy="no-referrer" />`
    data_css += `<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.41.0/min/vs/editor/editor.main.min.css" integrity="sha512-MFDhxgOYIqLdcYTXw7en/n5BshKoduTitYmX8TkQ+iJOGjrWusRi8+KmfZOrgaDrCjZSotH2d1U1e/Z1KT6nWw==" crossorigin="anonymous" referrerpolicy="no-referrer" />`

    cookie_map := Get_cookie_header(cookies)
    if cookie_map["main_css_darkmode"] == "1" {
        // Main CSS
        data_css_dark += `<link rel="stylesheet" href="/views/main_css/css/sub/dark.css` + data_css_ver + `">`

        // External CSS
        data_css_dark += `<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.8.0/styles/dark.min.css" integrity="sha512-bfLTSZK4qMP/TWeS1XJAR/VDX0Uhe84nN5YmpKk5x8lMkV0D+LwbuxaJMYTPIV13FzEv4CUOhHoc+xZBDgG9QA==" crossorigin="anonymous" referrerpolicy="no-referrer" />`
    }

    end := 2
    if end > len(data) {
        end = len(data)
    }

    new_data := append([]any{}, data[:end]...)
    new_data = append(new_data, "", data_css)

    if len(data) >= 3 {
        new_data = append(new_data, data[2])
    }

    new_data = append(new_data, data_css_dark)

    if len(data) >= 3 {
        new_data = append(new_data, data[3:]...)
    }

    return new_data
}

func Get_list_ui(left string, right string, bottom string, class_name string) string {
    data_html := ""

    data_html += `<span class="` + class_name + `">`
    data_html += `<div class="opennamu_recent_change">`
    data_html += left

    data_html += `<div style="float: right;">`
    data_html += right
    data_html += `</div>`

    data_html += `<div style="clear: both;"></div>`

    if bottom != "" {
        data_html += "<hr>"
        data_html += bottom
    }

    data_html += "</div>"
    data_html += `<hr class="main_hr">`
    data_html += "</span>"

    return data_html
}

// Get_error_page : auth, slow edit limit, edit filter (content), edit filter (send), send require, checkbox check require, overflow max length
func Get_error_page(db *sql.DB, config Config, error_name string) string {
    data := ""
    switch error_name {
    case "auth":
        data = Get_language(db, "authority_error", true)
    case "slow edit limit":
        data = Get_language(db, "fast_edit_error", true)
    case "edit filter (content)":
        data = Get_language(db, "edit_filter_error", true) + " (content)"
    case "edit filter (send)":
        data = Get_language(db, "edit_filter_error", true) + " (send)"
    case "send require":
        data = Get_language(db, "error_edit_send_request", true)
    case "checkbox check require":
        data = Get_language(db, "copyright_disagreed", true)
    case "overflow max length":
        data = Get_language(db, "document_content_max_length_error", true)
    default:
        data = Get_language(db, "inter_error", true)
    }

    return Get_template(
        db,
        config,
        Get_language(db, "error", true),
        `<h2>` + Get_language(db, "error", true) + `</h2>` +
        `<ul>` +
            `<li>` + data + `</li>` +
        `</ul>`,
        []any{},
        [][]any{},
        map[string]string{},
    )
}

func Get_page_control(db *sql.DB, page int, count int, max_count int, url string) string {
    data_html := "<hr class=\"main_hr\">"

    if page > 1 {
        prev_page := page - 1
        before_url := strings.ReplaceAll(url, "{}", strconv.Itoa(prev_page))

        data_html += `<a href="` + before_url + `">(` + Get_language(db, "previous", true) + `)</a> `
    }
    
    if count == max_count {
        prev_page := page + 1
        after_url := strings.ReplaceAll(url, "{}", strconv.Itoa(prev_page))

        data_html += `<a href="` + after_url + `">(` + Get_language(db, "next", true) + `)</a> `
    }

    return data_html
}

func Get_editor_ui(db *sql.DB, config Config, data string, do_type string, add_on string, doc_name string) string {
    monaco_editor_top := ""
    help_text := ""
    document_top := ""

    switch do_type {
    case "edit":
        QueryRow_DB(
            db,
            `select data from other where name = "edit_help"`,
            []any{ &help_text },
        )
        
        QueryRow_DB(
            db,
            `select set_data from data_set where doc_name = ? and set_name = 'document_top'`,
            []any{ &document_top },
            doc_name,
        )
    case "bbs":
        QueryRow_DB(
            db,
            `select data from other where name = "bbs_help"`,
            []any{ &help_text },
        )
    case "bbs_comment":
        QueryRow_DB(
            db,
            `select data from other where name = "bbs_comment_help"`,
            []any{ &help_text },
        )
    default:
        QueryRow_DB(
            db,
            `select data from other where name = "topic_text"`,
            []any{ &help_text },
        )
    }

    if help_text == "" {
        help_text = Get_language(db, "default_edit_help", true)
    } else {
        help_text = HTML_escape(help_text)
    }

    editor_type := "edit"
    if do_type == "bbs_comment" || do_type == "thread" {
        editor_type = "thread"
    }

    monaco_editor_top += `
        <a href="javascript:opennamu_do_editor_temp_save();">(` + Get_language(db, "load_temp_save", true) + `)</a> <a href="javascript:opennamu_do_editor_temp_save_load();">(` + Get_language(db, "load_temp_save_load", true) + `)</a>
        <hr class="main_hr">
    `

    dark_mode := false

    cookie_map := Get_cookie_header(config.Cookies)
    if cookie_map["main_css_darkmode"] == "1" {
        dark_mode = true
    }

    monaco_theme := ""
    if dark_mode {
        monaco_theme = "vs-dark"
    }

    monaco_on := false
    if Get_main_skin_set(db, config, "main_css_monaco") == "use" {
        monaco_on = true
    }

    editor_display := []string{}
    for for_a := 0; for_a < 3; for_a++ {
        editor_display = append(editor_display, `style="display: none;"`)
    }

    select_A := ""
    select_B := ""

    if monaco_on {
        editor_display[1] = ""
        select_B = "selected"
    } else {
        editor_display[0] = ""
        select_A = "selected"
    }

    monaco_editor_top += `
        <span class="__ON_SELECT_DIV__">
            <select class="__ON_SELECT__" onclick="do_sync_monaco_and_textarea();" id="opennamu_select_editor" onchange="opennamu_edit_turn_off_monaco();">
                <option value="default" ` + select_A + `>` + Get_language(db, "default", true) + `</option>
                <option value="monaco" ` + select_B + `>` + Get_language(db, "monaco_editor", true) + `</option>
            </select>
        </span>
    `

    if editor_type == "edit" {
        monaco_editor_top += Get_markup_select_ui(db, config, doc_name, "", `id="opennamu_editor_markup" onclick="opennamu_do_sync_monaco_markup();"`, "")
    } else {
        monaco_editor_top += Get_markup_select_ui(db, config, doc_name, "", `id="opennamu_editor_markup" onclick="opennamu_do_sync_monaco_markup();"`, "disabled")
    }

    textarea_size := "opennamu_textarea_500"
    if editor_type != "edit" {
        textarea_size = "opennamu_textarea_100"
    }

    out_field := Get_captcha_ui(db, config) + Get_IP_warning_ui(db, config) + add_on

    return `
        <textarea class="__ON_TEXTAREA__" style="display: none;" id="opennamu_edit_origin" name="doc_data_org">` + HTML_escape(data) + `</textarea>
        <div>
            ` + monaco_editor_top + `
            <hr class="main_hr">
            ` + Get_editor_button_ui(db) + `
            <div id="opennamu_editor_user_button"></div>
        </div>

        ` + document_top + `

        <div id="opennamu_monaco_editor" class="` + textarea_size + `" ` + editor_display[1] + `></div>
        <textarea id="opennamu_edit_textarea" class="` + textarea_size + ` __ON_TEXTAREA__" ` + editor_display[0] + ` name="content" placeholder="` + help_text + `">` + HTML_escape(data) + `</textarea>
        <hr class="main_hr">
        ` + out_field + `

        <script>
            window.addEventListener('DOMContentLoaded', function() {
                do_stop_exit();
                do_paste_image();
                do_monaco_init("` + monaco_theme + `");
                opennnamu_do_user_editor();
            });
        </script>

        <button class="__ON_BUTTON__" id="opennamu_save_button" type="submit" onclick="do_stop_exit_release();">` + Get_language(db, "send", true) + `</button>
        <button class="__ON_BUTTON__" id="opennamu_preview_button" type="button" onclick="opennamu_do_editor_preview();">` + Get_language(db, "preview", true) + `</button>
        <hr class="main_hr">

        <div id="opennamu_preview_area"></div>
    `
}

func Get_markup_select_ui(db *sql.DB, config Config, doc_name string, markup string, add_on string, disable string) string {
    default_markup := ""
    QueryRow_DB(
        db,
        `select data from other where name = "markup"`, 
        []any{ &default_markup },
    )

    markup_load := markup
    if markup == "" {
        QueryRow_DB(
            db,
            `select set_data from data_set where doc_name = ? and set_name = 'document_markup'`,
            []any{ &markup_load },
            doc_name,
        )
    }

    markup_list := []string{ "normal" }
    markup_list = append(markup_list, Get_init_set_list("markup")["markup"]["list"].([]string)...)

    markup_html := ""
    for _, v := range markup_list {
        selected := ""
        if markup_load == v {
            selected = "selected"
        }

        value := v
        if v == "normal" {
            value = default_markup
        }

        markup_html += `<option value="` + value + `" ` + selected + `>` + v + `</option>`
    }

    markup_html = `
        <span class="__ON_SELECT_DIV__">
            <select class="__ON_SELECT__" name="document_markup" ` + disable + ` ` + add_on + `>` + markup_html + `</select>
        </span>
    `

    return markup_html
}

func Get_captcha_ui(db *sql.DB, config Config) string {
    data := ""

    if !Check_acl(db, "", "", "recaptcha", config.IP) {
        pub_key := ""
        QueryRow_DB(
            db,
            `select data from other where name = "recaptcha"`, 
            []any{ &pub_key },
        )

        sec_key := ""
        QueryRow_DB(
            db,
            `select data from other where name = "sec_re"`, 
            []any{ &sec_key },
        )

        if pub_key != "" && sec_key != "" {
            rec_ver := ""
            QueryRow_DB(
                db,
                `select data from other where name = "recaptcha_ver"`, 
                []any{ &rec_ver },
            )

            switch rec_ver {
            case "":
                data += `
                    <script defer src="https://www.google.com/recaptcha/api.js"></script>
                    <div class="g-recaptcha" data-sitekey="` + pub_key + `"></div>
                    <hr class="main_hr">
                `
            case "v3":
                data += `
                    <script defer src="https://www.google.com/recaptcha/api.js?render=` + pub_key + `"></script>
                    <input class="__ON_INPUT__" type="hidden" id="g-recaptcha" name="g-recaptcha">
                    <script type="text/javascript">
                        document.addEventListener('DOMContentLoaded', function () {
                            grecaptcha.ready(function() {
                                grecaptcha.execute('` + pub_key + `', {action: 'homepage'}).then(function(token) {
                                    document.getElementById('g-recaptcha').value = token;
                                });
                            });
                        });
                    </script>
                `
            case "cf":
                data += `
                    <script defer src="https://challenges.cloudflare.com/turnstile/v0/api.js?compat=recaptcha"></script>
                    <div class="g-recaptcha" data-sitekey="` + pub_key + `"></div>
                    <hr class="main_hr">
                `
            default:
                data += `
                    <script defer src="https://js.hcaptcha.com/1/api.js"></script>
                    <div class="h-captcha" data-sitekey="` + pub_key + `"></div>
                    <hr class="main_hr">
                `
            }
        }
    }

    return data
}

func Get_IP_warning_ui(db *sql.DB, config Config) string {
    text_data := ""

    if IP_or_user(config.IP) {
        text_db := ""

        QueryRow_DB(
            db,
            `select data from other where name = "no_login_warning"`,
            []any{ &text_db },
        )

        if text_db == "" {
            text_db = Get_language(db, "no_login_warning", true)
        }
        
        text_data = `<span>` + text_db + `</span><hr class="main_hr">`
    }

    return text_data
}

func Get_editor_button_ui(db *sql.DB) string {
    data_html := ""

    rows := Query_DB(
        db,
        `select html, plus from html_filter where kind = 'edit_top'`,
    )
    defer rows.Close()

    for rows.Next() {
        var html string
        var plus string

        err := rows.Scan(&html, &plus)
        if err != nil {
            panic(err)
        }

        data_html += `<a href="javascript:do_insert_data('` + JS_escape(plus) + `');">` + HTML_escape(html) + `</a>`
    }

    if data_html != "" {
        data_html += " "
    }

    data_html += `<a href="/filter/edit_top">(` + Get_language(db, "add", true) + `)</a><hr class="main_hr">`

    return data_html
}

func Get_thread_ui(user_name string, date string, data string, code string, color string, blind string, add_style string, topic_num string) string {
    color_b := ""
    class_b := ""

    if blind == "O" {
        if data == "" {
            color_b = "opennamu_comment_blind"
        } else {
            color_b = "opennamu_comment_blind_admin"
        }

        class_b = "opennamu_comment_blind_js opennamu_list_hidden"
    } else {
        color_b = "opennamu_comment_blind_not"
    }

    admin_check_box := ""
    if topic_num != "" {
        admin_check_box = `<input type="checkbox" class="opennamu_blind_button" id="opennamu_blind_` + topic_num + `_` + code + `">`
    }

    return `
        <span class="` + class_b + `">
            <table class="opennamu_comment" style="` + add_style + `">
                <tr>
                    <td class="opennamu_comment_color_` + color + `">
                        ` + admin_check_box + `
                        <a href="#thread_shortcut" id="` + code + `">#` + code + `</a>
                        ` + user_name + `
                        <span style="float: right;">` + date + `</span>
                    </td>
                </tr>
                <tr>
                    <td class="` + color_b + ` opennamu_comment_data_main" id="thread_` + code + `">
                        <div class="opennamu_comment_scroll" id="opennamu_thread_render_` + code + `">` + HTML_escape(data) + `</div>
                    </td>
                    <script>
                        window.addEventListener('DOMContentLoaded', function() {
                            opennamu_do_render_with_dom("opennamu_thread_render_` + code + `", "opennamu_thread_render_` + code + `", "", "thread");
                        })
                    </script>
                </tr>
            </table>
            <hr class="main_hr">
        </span>
    `
}

func Get_edit_check_box_ui(db *sql.DB) string {
    cccb_text := ""

    QueryRow_DB(
        db,
        `select data from other where name = "copyright_checkbox_text"`,
        []any{ &cccb_text },
    )

    result := ""
    if cccb_text != "" {
        result = `
            <label class="__ON_CHECKLABEL__"><input class="__ON_CHECKBOX__" type="checkbox" name="copyright_agreement" value="yes" checked> ` + cccb_text + `</label>
            <hr class="main_hr">
        `
    }

    return result
}

func Get_edit_bottom_text_ui(db *sql.DB, do_type string) string {
    b_text := ""

    QueryRow_DB(
        db,
        `select data from other where name = "edit_bottom_text"`,
        []any{ &b_text },
    )

    db_data := ""

    switch do_type {
    case "edit":
        QueryRow_DB(
            db,
            `select data from other where name = "edit_only_bottom_text"`,
            []any{ &db_data },
        )
    case "move":
        QueryRow_DB(
            db,
            `select data from other where name = "move_bottom_text"`,
            []any{ &db_data },
        )
    case "delete":
        QueryRow_DB(
            db,
            `select data from other where name = "delete_bottom_text"`,
            []any{ &db_data },
        )
    default:
        QueryRow_DB(
            db,
            `select data from other where name = "revert_bottom_text"`,
            []any{ &db_data },
        )
    }
    
    result := ""
    if db_data != "" {
        result = db_data + `<hr class="main_hr">`
    } else if b_text != "" {
        result = b_text + `<hr class="main_hr">`
    }

    return result
}