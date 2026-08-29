package route

import (
	"database/sql"
	"strings"

	"opennamu/route/tool"
)

func View_setting_external(config tool.Config) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_permission(db, "setting_external", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}

	return view_setting_external_data(db, config, setting_load_fields(db, setting_external_fields()))
}

func setting_external_fields() []setting_field {
	return []setting_field{
		{name: "recaptcha"},
		{name: "sec_re"},
		{name: "smtp_server"},
		{name: "smtp_port"},
		{name: "smtp_security"},
		{name: "smtp_email"},
		{name: "smtp_pass"},
		{name: "recaptcha_ver"},
		{name: "oauth_client_id"},
		{name: "email_have"},
	}
}

func view_setting_external_data(db *sql.DB, config tool.Config, values map[string]string) string {
	lang := func(name string) string {
		return tool.Get_language(db, name, true)
	}

	recaptcha_labels := map[string]string{
		"":              "reCAPTCHA v2",
		"v3":            "reCAPTCHA v3",
		"h":             "hCAPTCHA",
		"cf":            "Turnstile",
		"altcha_low":    "ALTCHA Low",
		"altcha_medium": "ALTCHA Medium",
		"altcha_high":   "ALTCHA High",
	}

	data := strings.Builder{}
	data.WriteString(`<form method="post">`)
	data.WriteString(`<h2>` + lang("captcha") + `</h2>`)
	data.WriteString(`<a href="https://www.google.com/recaptcha/">(` + lang("recaptcha") + `)</a> <a href="https://www.hcaptcha.com/">(` + lang("hcaptcha") + `)</a> <a href="https://altcha.org/">(ALTCHA)</a>` + setting_hr())
	data.WriteString(`<span>` + lang("public_key") + `</span>` + setting_hr())
	data.WriteString(setting_input("recaptcha", values["recaptcha"], "text") + setting_hr())
	data.WriteString(`<span>` + lang("secret_key") + `</span>` + setting_hr())
	data.WriteString(setting_input("sec_re", values["sec_re"], "text") + setting_hr())
	data.WriteString(`<span>` + lang("version") + `</span>` + setting_hr())
	data.WriteString(`<select name="recaptcha_ver">` + setting_options(values["recaptcha_ver"], []string{"", "v3", "h", "cf", "altcha_low", "altcha_medium", "altcha_high"}, recaptcha_labels) + `</select>` + setting_hr())

	data.WriteString(`<h2>` + lang("email_setting") + `</h2>`)
	data.WriteString(`<a href="/setting/phrase#s-6">(` + lang("text_setting") + `)</a>` + setting_hr())
	data.WriteString(`<label><input type="checkbox" name="email_have" ` + setting_checked(values["email_have"]) + `> ` + lang("email_required") + `</label>`)
	data.WriteString(`<h3>` + lang("smtp_setting") + `</h3>`)
	data.WriteString(`<a href="https://support.google.com/mail/answer/7126229">(Google)</a>` + setting_hr())
	data.WriteString(`<a href="/setting/email_test">(` + lang("test") + `)</a>` + setting_hr())
	data.WriteString(`<span>` + lang("smtp_server") + `</span>` + setting_hr())
	data.WriteString(setting_input("smtp_server", values["smtp_server"], "text") + setting_hr())
	data.WriteString(`<span>` + lang("smtp_port") + `</span>` + setting_hr())
	data.WriteString(setting_input("smtp_port", values["smtp_port"], "text") + setting_hr())
	data.WriteString(`<span>` + lang("smtp_security") + `</span>` + setting_hr())
	data.WriteString(`<select name="smtp_security">` + setting_options(values["smtp_security"], []string{"tls", "starttls", "plain"}, nil) + `</select>` + setting_hr())
	data.WriteString(`<span>` + lang("smtp_username") + `</span>` + setting_hr())
	data.WriteString(setting_input("smtp_email", values["smtp_email"], "text") + setting_hr())
	data.WriteString(`<span>` + lang("smtp_password") + `</span>` + setting_hr())
	data.WriteString(setting_input("smtp_pass", values["smtp_pass"], "password") + setting_hr())

	data.WriteString(`<h2>` + lang("oauth") + ` (` + lang("not_working") + `)</h2>`)
	data.WriteString(`<a href="https://developers.google.com/identity/protocols/oauth2">(Google)</a>` + setting_hr())
	data.WriteString(`<span>` + lang("oauth_client_id") + `</span>` + setting_hr())
	data.WriteString(setting_input("oauth_client_id", values["oauth_client_id"], "text") + setting_hr())
	data.WriteString(`<button id="opennamu_save_button" type="submit">` + lang("save") + `</button></form>`)

	return setting_page(db, config, lang("ext_api_req_set"), data.String(), "setting")
}
