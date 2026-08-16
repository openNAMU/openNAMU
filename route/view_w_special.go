package route

import (
	"database/sql"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"opennamu/route/tool"
)

func view_w_redirect_target(db *sql.DB, doc_name string) (string, string) {
	target := ""
	anchor := ""
	tool.QueryRow_DB(
		db,
		"select title, data from back where link = ? and type = 'redirect' limit 1",
		[]any{&target, &anchor},
		doc_name,
	)
	if target != "" {
		target_exists := ""
		if !tool.QueryRow_DB(db, "select title from data where title = ?", []any{&target_exists}, target) {
			return "", ""
		}
	}

	return target, anchor
}

func view_w_child_exists(db *sql.DB, doc_name string) bool {
	child := ""
	return tool.QueryRow_DB(
		db,
		"select title from data where title like ? limit 1",
		[]any{&child},
		doc_name+"/%",
	)
}

func view_w_user_data(db *sql.DB, doc_name string) string {
	user_name := strings.TrimPrefix(doc_name, "user:")
	if slash_index := strings.Index(user_name, "/"); slash_index >= 0 {
		user_name = user_name[:slash_index]
	}

	phrase := ""
	if tool.Get_user_document(db, user_name) && !tool.Check_acl(db, "", "", "all_admin_auth", user_name) {
		phrase_name := "phrase_user_page_owner"
		if tool.Check_acl(db, "", "", "owner_auth", user_name) {
			phrase_name = "phrase_user_page_admin"
		}
		tool.QueryRow_DB(db, "select data from other where name = ?", []any{&phrase}, phrase_name)
	}
	if phrase != "" {
		phrase += "<br>"
	}
	return phrase + `<div id="opennamu_get_user_info">` + tool.HTML_escape(user_name) + `</div><hr class="main_hr">`
}

func view_w_category_data(db *sql.DB, config tool.Config, doc_name string) string {
	type category_entry struct {
		name string
		view string
	}

	category_list := []category_entry{}
	category_seen := map[string]bool{}
	category_blur := map[string]bool{}
	add_category := func(name string, view string) {
		if name == "" || category_seen[name] {
			return
		}
		if view == "" {
			view = name
		}
		category_seen[name] = true
		category_list = append(category_list, category_entry{name, view})
	}
	show_category_view := tool.Get_main_skin_set(db, config, "main_css_category_change_title") != "off"
	get_category_meta := func(name string) (string, bool) {
		view := ""
		blur_value := ""
		if show_category_view {
			tool.QueryRow_DB(db, "select data from back where title = ? and link = ? and type = 'cat_view' limit 1", []any{&view}, doc_name, name)
		}
		blur := tool.QueryRow_DB(db, "select data from back where title = ? and link = ? and type = 'cat_blur' limit 1", []any{&blur_value}, doc_name, name)
		return view, blur
	}

	rows := tool.Query_DB(db, "select link, data from back where title = ? and type = 'cat' order by link", doc_name)
	for rows.Next() {
		name, view := "", ""
		if rows.Scan(&name, &view) == nil {
			category_view, blur := get_category_meta(name)
			if category_view != "" {
				view = category_view
			}
			if blur {
				category_blur[name] = true
			}
			add_category(name, view)
		}
	}
	rows.Close()

	if len(category_list) == 0 {
		rows = tool.Query_DB(db, "select title, data from back where link = ? and (type = 'cat' or type = '') order by title", doc_name)
		for rows.Next() {
			name, view := "", ""
			if rows.Scan(&name, &view) == nil {
				add_category(name, view)
			}
		}
		rows.Close()
	}

	if len(category_list) == 0 {
		return ""
	}

	category_sub := ""
	category_doc := ""
	category_sub_count := 0
	category_doc_count := 0
	for _, entry := range category_list {
		class_name := ""
		if category_blur[entry.name] {
			class_name = ` class="opennamu_category_blur"`
		}
		item := `<li><a` + class_name + ` href="/w/` + tool.Url_parser(entry.name) + `">` + tool.HTML_escape(entry.view) + `</a>`
		if strings.HasPrefix(entry.name, "category:") {
			category_sub += item + `</li>`
			category_sub_count++
		} else {
			category_doc += item + ` <a class="opennamu_link_inter" href="/xref/` + tool.Url_parser(entry.name) + `">(` + tool.Get_language(db, "backlink", true) + `)</a></li>`
			category_doc_count++
		}
	}

	data := ""
	if category_sub != "" {
		data += `<h2>` + tool.Get_language(db, "under_category", true) + `</h2><ul><li>` + tool.Get_language(db, "all", true) + " : " + strconv.Itoa(category_sub_count) + `</li>` + category_sub + `</ul>`
	}
	if category_doc != "" {
		data += `<h2>` + tool.Get_language(db, "category_title", true) + `</h2><ul><li>` + tool.Get_language(db, "all", true) + " : " + strconv.Itoa(category_doc_count) + `</li>` + category_doc + `</ul>`
	}
	return data + `<hr class="main_hr">`
}

func view_w_file_data(db *sql.DB, doc_name string) string {
	file_name := strings.TrimPrefix(doc_name, "file:")
	extension := strings.TrimPrefix(strings.ToLower(filepath.Ext(file_name)), ".")
	if extension == "" {
		return ""
	}

	base_name := strings.TrimSuffix(file_name, filepath.Ext(file_name))
	storage_name := tool.File_name_to_dir(base_name, extension)
	file_path := filepath.Join(tool.Get_file_main_dir(db), storage_name)
	file_info, err := os.Stat(file_path)
	if err != nil || file_info.IsDir() {
		return ""
	}

	cache_revision := "1"
	tool.QueryRow_DB(db, "select id from history where title = ? order by date desc limit 1", []any{&cache_revision}, doc_name)
	image_url := "/image/" + tool.Url_parser(storage_name) + ".cache_v" + tool.Url_parser(cache_revision)

	resolution := "Vector"
	if file, open_err := os.Open(file_path); open_err == nil {
		image_config, _, decode_err := image.DecodeConfig(file)
		_ = file.Close()
		if decode_err == nil {
			resolution = strconv.Itoa(image_config.Width) + "x" + strconv.Itoa(image_config.Height)
		}
	}

	return `<img src="` + image_url + `"><h2>` + tool.Get_language(db, "data", true) + `</h2><table><tr><td>` + tool.Get_language(db, "url", true) + `</td><td><a href="/image/` + tool.Url_parser(storage_name) + `">` + tool.Get_language(db, "link", true) + `</a></td></tr><tr><td>` + tool.Get_language(db, "volume", true) + `</td><td>` + strconv.FormatFloat(float64(file_info.Size())/1000, 'f', 1, 64) + `KB</td></tr><tr><td>` + tool.Get_language(db, "resolution", true) + `</td><td>` + resolution + `</td></tr></table><h2>` + tool.Get_language(db, "content", true) + `</h2>`
}
