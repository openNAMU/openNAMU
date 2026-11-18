package route

import (
	"database/sql"
	"strconv"
	"time"

	"opennamu/route/tool"
	"opennamu/route/tool/markup"
)

// Setting_markup_trusted_options lists the markups selectable on trusted
// phrase surfaces (setting_phrase administrators, document owners). These
// surfaces may store raw html, unlike document body or BBS markup.
func Setting_markup_trusted_options() []string {
	return append(markup.List_markup(), "html")
}

// Setting_markup_normalize normalizes a trusted phrase markup. Unknown or
// retired values (custom, js, ...) stay fail-closed as raw escaped text.
func Setting_markup_normalize(markup_name string) string {
	if markup_name == "" || !tool.Arr_in_str(Setting_markup_trusted_options(), markup_name) {
		return "raw"
	}
	return markup_name
}

// Setting_markup_value_valid checks untrusted markup writes such as document
// body or BBS markup, where raw html stays rejected.
func Setting_markup_value_valid(markup_name string) bool {
	return markup_name == "" || tool.Arr_in_str(Setting_markup_options(), markup_name)
}

// Setting_markup_trusted_value_valid checks trusted phrase markup writes,
// which may additionally select raw html.
func Setting_markup_trusted_value_valid(markup_name string) bool {
	return markup_name == "" || tool.Arr_in_str(Setting_markup_trusted_options(), markup_name)
}

// Setting_markup_storage_value canonicalizes an untrusted markup for storage.
func Setting_markup_storage_value(markup_name string) string {
	if markup_name == "" {
		return ""
	}
	markup_name = tool.Normalize_document_markup(markup_name)
	if markup_name == "raw" {
		return "plain"
	}
	return markup_name
}

// Setting_markup_trusted_storage_value canonicalizes a trusted phrase markup
// for storage. html stays selectable; anything invalid falls back to plain.
func Setting_markup_trusted_storage_value(markup_name string) string {
	if markup_name == "" {
		return ""
	}
	markup_name = Setting_markup_normalize(markup_name)
	if markup_name == "raw" {
		return "plain"
	}
	return markup_name
}

func Setting_markup_value(db *sql.DB, setting_name string) string {
	return Setting_markup_normalize(
		tool.Get_setting_value(db, setting_name+"_markup", "", ""),
	)
}

// Setting_markup_options lists the markups selectable for untrusted content
// (document body, BBS). Trusted phrase surfaces use
// Setting_markup_trusted_options instead.
func Setting_markup_options() []string {
	return markup.List_markup()
}

// Setting_markup_select_ui builds the markup select for trusted phrase
// surfaces, so it includes the raw html option.
func Setting_markup_select_ui(name string, current string, disabled string) string {
	current = Setting_markup_normalize(current)
	data := `<span class="__ON_SELECT_DIV__"><select class="__ON_SELECT__" name="` + tool.HTML_escape(name) + `"` + disabled + `>`
	for _, markup_option := range Setting_markup_trusted_options() {
		selected := ""
		if markup_option == current {
			selected = ` selected`
		}
		data += `<option value="` + tool.HTML_escape(markup_option) + `"` + selected + `>` + tool.HTML_escape(markup_option) + `</option>`
	}
	data += `</select></span>`
	return data
}

func Setting_render_markup(db *sql.DB, data string, markup_name string) string {
	return Setting_render_markup_with_requester(db, data, markup_name, "")
}

// Setting_render_markup_with_requester renders trusted phrase content
// (global and document top/bottom phrases). html passes through unescaped
// because only setting_phrase administrators and owners can store it.
func Setting_render_markup_with_requester(db *sql.DB, data string, markup_name string, requester_identity string) string {
	markup_name = Setting_markup_normalize(markup_name)
	if markup_name == "html" {
		return data
	}
	render_data := markup.Get_render_direct_with_requester(
		db,
		"",
		data,
		markup_name,
		strconv.FormatInt(time.Now().UnixNano(), 10),
		"setting",
		nil,
		requester_identity,
	)

	return render_data["data"]
}

func Document_editor_top_render(db *sql.DB, doc_name string) string {
	return Document_editor_top_render_with_requester(db, doc_name, "")
}

func Document_editor_top_render_with_requester(db *sql.DB, doc_name string, requester_identity string) string {
	top_data := tool.Get_document_setting_value_exact(db, doc_name, "document_editor_top", "")
	top_markup := ""
	if top_data != "" {
		top_markup = Setting_markup_normalize(tool.Get_document_setting_value_exact(db, doc_name, "document_editor_top_markup", ""))
	} else {
		top_data = tool.Get_document_setting_value_exact(db, doc_name, "document_top", "")
		if top_data == "" {
			return ""
		}
		top_markup = Setting_markup_normalize(tool.Get_document_setting_value_exact(db, doc_name, "document_top_markup", ""))
	}

	return Setting_render_markup_with_requester(db, top_data, top_markup, requester_identity)
}
