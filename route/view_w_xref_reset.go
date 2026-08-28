package route

import (
	"opennamu/route/tool"
)

func View_w_xref_reset(config tool.Config, doc_name string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, doc_name, "", "document_edit", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	api_data := Api_w_raw(config, doc_name, "", "")
	response := api_data["response"].(string)
	if response == "require auth" {
		return tool.Get_error_page(db, config, "auth")
	} else if response != "ok" {
		return tool.Get_redirect("/w/" + tool.Url_parser(doc_name))
	}

	return tool.Get_template(
		db,
		config,
		doc_name,
		`<form method="post">
            <button id="opennamu_save_button" type="submit">`+tool.Get_language(db, "reset_backlink", true)+`</button>
        </form>`,
		[]any{"(" + tool.Get_language(db, "reset_backlink", true) + ")"},
		[][]any{
			{"xref/" + tool.Url_parser(doc_name), tool.Get_language(db, "return", true)},
		},
		map[string]string{},
	)
}
