package route

import (
	"database/sql"
	"strings"

	"opennamu/route/tool"
)

func View_ollama(config tool.Config, question string, model string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "ai_use", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	question = strings.TrimSpace(question)
	if tool.Get_len(question) > 1000 {
		question = tool.Get_slice(question, 0, 1000)
	}
	model = strings.TrimSpace(model)
	if model == "" {
		model = "gemma4:e4b"
	}

	data_html := `<p>` + tool.Get_language(db, "ollama_requirement", true) + `</p><form method="post" action="/ai">` +
		`<textarea class="opennamu_textarea_100" name="question" placeholder="` + tool.Get_language(db, "ollama_question", true) + `">` + tool.HTML_escape(question) + `</textarea>` +
		`<input name="model" value="` + tool.HTML_escape(model) + `" placeholder="` + tool.Get_language(db, "ollama_model", true) + `">` +
		`<hr class="main_hr"><button type="submit">` + tool.Get_language(db, "go", true) + `</button></form>`

	if question != "" {
		context_data, source_list := ollama_document_context(db, config, question)
		prompt := "너는 위키 문서 검색을 돕는 AI다. 아래 참고 문서에 있는 내용만 근거로 답변하고, 근거가 없으면 모른다고 답변해라. 참고 문서:\n\n" + context_data + "\n질문: " + question
		answer, err := Api_ollama_stream(model, prompt)
		if err != nil {
			data_html += `<hr class="main_hr"><p>` + tool.Get_language(db, "ollama_error", true) + `</p>`
		} else {
			answer = strings.ReplaceAll(answer, "\r", "")
			answer_html := strings.ReplaceAll(tool.HTML_escape(answer), "\n", "<br>")
			data_html += `<hr class="main_hr"><h2>` + tool.Get_language(db, "ollama_answer", true) + `</h2><div>` + answer_html + `</div>`
		}

		if len(source_list) > 0 {
			data_html += `<h3>` + tool.Get_language(db, "ollama_source", true) + `</h3><ul>`
			for _, title := range source_list {
				data_html += `<li><a href="/w/` + tool.Url_parser(title) + `">` + tool.HTML_escape(title) + `</a></li>`
			}
			data_html += `</ul>`
		}
	}

	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "local_ai", true),
		data_html,
		[]any{},
		[][]any{{"other", tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}

func ollama_document_context(db *sql.DB, config tool.Config, question string) (string, []string) {
	search_data := Api_func_search(config, question, "1", "data")
	title_list, _ := search_data["data"].([]string)
	if len(title_list) > 3 {
		title_list = title_list[:3]
	}

	context_data := strings.Builder{}
	source_list := []string{}
	for _, title := range title_list {
		if !tool.Check_acl(db, title, "", "render", config.IP) {
			continue
		}
		raw_data, exists := tool.Get_data_content(db, title)
		if !exists || raw_data == "" {
			continue
		}
		if tool.Get_len(raw_data) > 3000 {
			raw_data = tool.Get_slice(raw_data, 0, 3000)
		}
		context_data.WriteString("[문서: ")
		context_data.WriteString(title)
		context_data.WriteString("]\n")
		context_data.WriteString(raw_data)
		context_data.WriteString("\n\n")
		source_list = append(source_list, title)
	}

	return context_data.String(), source_list
}
