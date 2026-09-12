package route

import (
	"net/url"
	"strings"

	"opennamu/route/tool"
)

func View_auth_invite(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "user_manage", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	token := ""
	if values != nil {
		if values.Get("make") != "" {
			result := Api_auth_invite_post(config, "make", "")
			if result["response"] != "ok" {
				return tool.Get_error_page(db, config, "error")
			}
			token, _ = result["token"].(string)
		} else if revoke := strings.TrimSpace(values.Get("revoke")); revoke != "" {
			result := Api_auth_invite_post(config, "revoke", revoke)
			if result["response"] != "ok" {
				return tool.Get_error_page(db, config, "error")
			}
			return tool.Get_redirect("/auth/invite")
		}
	}

	lang := func(name string) string {
		return tool.Get_language(db, name, true)
	}

	data := `<p>` + lang("invite_help") + `</p>` + main_hr()
	if token != "" {
		data += `<p>` + lang("invite_created") + `</p><code>` + tool.HTML_escape(token) + `</code>` + main_hr()
	}
	data += `<form method="post"><button name="make" value="1" type="submit">` + lang("invite_make") + `</button></form>` + main_hr()

	for _, record := range tool.Get_invite_list(db) {
		data += `<div class="opennamu_list_1">`
		data += lang("invite_issued_by") + ` : ` + tool.HTML_escape(record.Issuer)
		data += ` | ` + lang("invite_created_at") + ` : ` + tool.HTML_escape(record.Date)
		data += ` | ` + lang("invite_expire") + ` : ` + tool.HTML_escape(record.End)
		if !tool.Invite_valid(db, record.Hash) {
			data += ` (` + lang("invite_expired") + `)`
		}
		data += `<form method="post"><input type="hidden" name="revoke" value="` + tool.HTML_escape(record.Hash) + `"><button type="submit">` + lang("invite_revoke") + `</button></form></div>` + main_hr()
	}

	return tool.Get_template(
		db,
		config,
		lang("invite_manage"),
		data,
		[]any{},
		[][]any{{"manager", lang("return")}},
		map[string]string{},
	)
}
