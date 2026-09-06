package route

import (
	"opennamu/route/tool"
)

func View_setting_backlink_reset(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "setting_backlink", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	return tool.Get_template(
		db,
		config,
		tool.Get_language(db, "reset_all_backlink", true),
		`<form method="post">
            <span>`+tool.Get_language(db, "backlink_reset_load", true)+`</span><hr class="main_hr">
            <select name="backlink_reset_load">
                <option value="fast" selected="selected">`+tool.Get_language(db, "backlink_reset_load_fast", true)+`</option>
                <option value="normal">`+tool.Get_language(db, "backlink_reset_load_normal", true)+`</option>
                <option value="slow">`+tool.Get_language(db, "backlink_reset_load_slow", true)+`</option>
            </select><hr class="main_hr">
            <button id="opennamu_save_button" type="submit">`+tool.Get_language(db, "reset_all_backlink", true)+`</button>
        </form>`,
		[]any{},
		[][]any{
			{"setting", tool.Get_language(db, "return", true)},
		},
		map[string]string{},
	)
}
