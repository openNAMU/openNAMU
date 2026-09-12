package route

import (
	"strconv"

	"opennamu/route/tool"
)

func View_document_backup(config tool.Config, result map[string]any) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "admin", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	data := `<p>` + tool.Get_language(db, "document_backup_help", true) + `</p>
<a href="/backup/export">` + tool.Get_language(db, "document_backup_export", true) + `</a><hr class="main_hr">
<form method="post" enctype="multipart/form-data">
<input type="file" name="backup" accept=".json,application/json"><hr class="main_hr">
<button type="submit">` + tool.Get_language(db, "document_backup_import", true) + `</button>
</form>`

	if result != nil {
		if response, ok := result["response"].(string); ok && response == "ok" {
			imported, _ := result["imported"].(int)
			skipped, _ := result["skipped"].(int)
			data = `<p>` + tool.Get_language(db, "document_backup_imported", true) + `: ` + strconv.Itoa(imported) + `<br>` +
				tool.Get_language(db, "document_backup_skipped", true) + `: ` + strconv.Itoa(skipped) + `</p><hr class="main_hr">` + data
		} else if response == "error" {
			data = `<p>` + tool.Get_language(db, "document_backup_invalid", true) + `</p><hr class="main_hr">` + data
		}
	}

	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "document_backup", true),
		data,
		[]any{},
		[][]any{{"manager", tool.Get_language(db, "return", true)}},
		map[string]string{},
	)
}
