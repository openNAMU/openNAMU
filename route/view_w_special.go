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
	target, anchor, redirect_exists := tool.Get_back_redirect_data(db, doc_name)
	if redirect_exists && target != "" {
		_, target_exists := tool.Get_data_title(db, target)
		if !target_exists {
			return "", ""
		}
	}
	return target, anchor
}

func view_w_child_exists(db *sql.DB, doc_name string) bool {
	return tool.Get_data_title_like(db, doc_name+"/%")
}

func view_w_user_data(db *sql.DB, doc_name string) string {
	user_name := strings.TrimPrefix(doc_name, "user:")
	if slash_index := strings.Index(user_name, "/"); slash_index >= 0 {
		user_name = user_name[:slash_index]
	}

	phrase := ""
	if tool.Get_user_document(db, user_name) && !tool.Check_permission(db, "treat_as_admin", user_name) {
		phrase_name := "phrase_user_page_owner"
		if tool.Check_permission(db, "owner", user_name) {
			phrase_name = "phrase_user_page_admin"
		}
		phrase = tool.Get_other_data(db, phrase_name)
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
	manual_documents := map[string]bool{}
	manual_rows := tool.Get_category_manual_document_rows(db, doc_name)
	for manual_rows.Next() {
		name := ""
		if manual_rows.Scan(&name) == nil {
			manual_documents[name] = true
		}
	}
	manual_rows.Close()
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
		if show_category_view {
			view, _ = tool.Get_category_meta(db, doc_name, name, "cat_view")
		}
		_, blur := tool.Get_category_meta(db, doc_name, name, "cat_blur")
		return view, blur
	}

	rows := tool.Get_category_rows(db, doc_name)
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
		rows = tool.Get_category_back_rows(db, doc_name)
		for rows.Next() {
			name, view := "", ""
			if rows.Scan(&name, &view) == nil {
				add_category(name, view)
			}
		}
		rows.Close()
	}

	if len(category_list) == 0 {
		if tool.Check_permission(db, "edit", config.IP) {
			return category_manual_add_form(db, doc_name, "", doc_name) + `<hr class="main_hr">`
		}
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
		manual_form := ""
		if manual_documents[entry.name] && tool.Check_acl(db, entry.name, "", "document_edit", config.IP) {
			manual_form = category_manual_delete_form(db, doc_name, entry.name, doc_name)
		}
		if strings.HasPrefix(entry.name, "category:") {
			category_sub += item + manual_form + `</li>`
			category_sub_count++
		} else {
			category_doc += item + ` <a class="opennamu_link_inter" href="/xref/` + tool.Url_parser(entry.name) + `">(` + tool.Get_language(db, "backlink", true) + `)</a>` + manual_form + `</li>`
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
	random_link := ""
	if category_doc_count > 0 {
		random_link = `(<a href="/random/category/` + tool.Url_parser(doc_name) + `">` + tool.Get_language(db, "random_category", false) + `</a>)`
	}
	if tool.Check_permission(db, "edit", config.IP) {
		data += category_manual_add_form(db, doc_name, "", doc_name)
	}
	return data + random_link + `<hr class="main_hr">`
}

func category_manual_delete_form(db *sql.DB, category_name string, doc_name string, return_name string) string {
	return `<form method="post" action="/category/delete" style="display:inline"><input type="hidden" name="category" value="` + tool.HTML_escape(category_name) + `"><input type="hidden" name="document" value="` + tool.HTML_escape(doc_name) + `"><input type="hidden" name="return" value="` + tool.HTML_escape(return_name) + `"><button type="submit">` + tool.Get_language(db, "delete", true) + `</button></form>`
}

func category_manual_add_form(db *sql.DB, category_name string, doc_name string, return_name string) string {
	data := `<form method="post" action="/category/add">`
	if category_name == "" {
		data += `<input type="text" name="category" placeholder="` + tool.HTML_escape(tool.Get_language(db, "category", true)) + `">`
	} else {
		data += `<input type="hidden" name="category" value="` + tool.HTML_escape(category_name) + `">`
	}
	if doc_name == "" {
		data += `<input type="text" name="document" placeholder="` + tool.HTML_escape(tool.Get_language(db, "document", true)) + `">`
	} else {
		data += `<input type="hidden" name="document" value="` + tool.HTML_escape(doc_name) + `">`
	}
	data += `<input type="hidden" name="return" value="` + tool.HTML_escape(return_name) + `"><button type="submit">` + tool.Get_language(db, "add", true) + `</button></form>`
	return data
}

func view_w_manual_category_data(db *sql.DB, config tool.Config, doc_name string) string {
	can_edit := tool.Check_acl(db, doc_name, "", "document_edit", config.IP)
	rows := tool.Get_category_manual_rows(db, doc_name)
	defer rows.Close()

	data := `<div class="opennamu_category" id="cate_manual">` + tool.Get_language(db, "category_manual", true) + " : "
	count := 0
	for rows.Next() {
		category_name, view := "", ""
		if rows.Scan(&category_name, &view) != nil {
			continue
		}
		if count > 0 {
			data += " | "
		}
		label := strings.TrimPrefix(category_name, "category:")
		if view == "" {
			view = label
		}
		data += `<a href="/w/` + tool.Url_parser(category_name) + `">` + tool.HTML_escape(view) + `</a>`
		if can_edit {
			data += category_manual_delete_form(db, doc_name, category_name, doc_name)
		}
		count++
	}
	data += `</div>`
	if !can_edit && count == 0 {
		return ""
	}
	if can_edit {
		data += category_manual_add_form(db, "", doc_name, doc_name)
	}
	return `<hr class="main_hr">` + data
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

	cache_revision := tool.Get_history_last_revision(db, doc_name)
	if cache_revision == "" {
		cache_revision = "1"
	}
	image_url := "/image/" + tool.Url_parser(storage_name) + ".cache_v" + tool.Url_parser(cache_revision)

	is_audio := tool.Is_audio_extension(extension)
	is_video := tool.Is_video_extension(extension)
	resolution := "Vector"
	if !is_audio && !is_video {
		if file, open_err := os.Open(file_path); open_err == nil {
			image_config, _, decode_err := image.DecodeConfig(file)
			_ = file.Close()
			if decode_err == nil {
				resolution = strconv.Itoa(image_config.Width) + "x" + strconv.Itoa(image_config.Height)
			}
		}
	} else {
		resolution = "-"
	}

	media := `<img src="` + image_url + `">`
	if is_audio {
		media = `<audio controls src="` + image_url + `"></audio>`
	}
	if is_video {
		media = `<video controls preload="metadata" src="` + image_url + `"></video>`
	}

	return media + `<h2>` + tool.Get_language(db, "data", true) + `</h2><table><tr><td>` + tool.Get_language(db, "url", true) + `</td><td><a href="/image/` + tool.Url_parser(storage_name) + `">` + tool.Get_language(db, "link", true) + `</a></td></tr><tr><td>` + tool.Get_language(db, "volume", true) + `</td><td>` + strconv.FormatFloat(float64(file_info.Size())/1000, 'f', 1, 64) + `KB</td></tr><tr><td>` + tool.Get_language(db, "resolution", true) + `</td><td>` + resolution + `</td></tr></table><h2>` + tool.Get_language(db, "content", true) + `</h2>`
}
