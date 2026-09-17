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

	return View_setting_external_data(db, config, Setting_load_fields(db, Setting_external_fields()))
}

func Setting_external_fields() []setting_field {
	return []setting_field{
		{name: "recaptcha"},
		{name: "sec_re"},
		{name: "altcha_sec_re"},
		{name: "smtp_server"},
		{name: "smtp_port"},
		{name: "smtp_security"},
		{name: "smtp_email"},
		{name: "smtp_pass"},
		{name: "recaptcha_ver", default_value: "altcha_high"},
		{name: "oauth_client_id"},
		{name: "ai_provider", default_value: "ollama"},
		{name: "ai_model"},
		{name: "openai_api_key"},
		{name: "google_api_key"},
		{name: "email_have"},
	}
}

func View_setting_external_data(db *sql.DB, config tool.Config, values map[string]string) string {
	lang := func(name string) string {
		return tool.Get_language(db, name, true)
	}

	recaptcha_labels := map[string]string{
		"v2":            "reCAPTCHA v2",
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
	data.WriteString(`<a href="https://www.google.com/recaptcha/">(` + lang("recaptcha") + `)</a> <a href="https://www.hcaptcha.com/">(` + lang("hcaptcha") + `)</a> <a href="https://altcha.org/">(ALTCHA)</a>` + Main_hr())
	data.WriteString(`<p>` + lang("altcha_info") + `</p>` + Main_hr())
	data.WriteString(`<span>` + lang("recaptcha") + ` ` + lang("public_key") + `</span>` + Main_hr())
	data.WriteString(Setting_input("recaptcha", values["recaptcha"], "text") + Main_hr())
	data.WriteString(`<span>` + lang("recaptcha") + ` ` + lang("secret_key") + `</span>` + Main_hr())
	data.WriteString(Setting_input("sec_re", values["sec_re"], "text") + Main_hr())
	data.WriteString(`<span>` + lang("altcha_secret_key") + `</span>` + Main_hr())
	data.WriteString(Setting_input("altcha_sec_re", values["altcha_sec_re"], "text") + Main_hr())
	data.WriteString(`<span>` + lang("version") + `</span>` + Main_hr())
	data.WriteString(`<select name="recaptcha_ver">` + Setting_options(values["recaptcha_ver"], []string{"v2", "v3", "h", "cf", "altcha_low", "altcha_medium", "altcha_high"}, recaptcha_labels) + `</select>` + Main_hr())

	data.WriteString(`<h2>` + lang("local_ai") + `</h2>`)
	data.WriteString(`<p>` + lang("ai_external_warning") + `</p>` + Main_hr())
	data.WriteString(`<span>` + lang("ai_provider") + `</span>` + Main_hr())
	data.WriteString(`<select name="ai_provider">` + Setting_options(values["ai_provider"], []string{"ollama", "openai", "google"}, map[string]string{
		"ollama": "Ollama",
		"openai": "OpenAI",
		"google": "Google Gemini",
	}) + `</select>` + Main_hr())
	ai_model := values["ai_model"]
	if ai_model == "" {
		ai_model = Ai_default_model(values["ai_provider"])
	}
	data.WriteString(`<span>` + lang("ai_model") + `</span>` + Main_hr())
	data.WriteString(Setting_input("ai_model", ai_model, "text") + Main_hr())
	data.WriteString(`<span>` + lang("openai_api_key") + `</span>` + Main_hr())
	data.WriteString(Setting_input("openai_api_key", "", "password") + Main_hr())
	data.WriteString(`<span>` + lang("google_api_key") + `</span>` + Main_hr())
	data.WriteString(Setting_input("google_api_key", "", "password") + Main_hr())

	data.WriteString(`<h2>` + lang("email_setting") + `</h2>`)
	data.WriteString(`<a href="/setting/phrase#s-6">(` + lang("text_setting") + `)</a>` + Main_hr())
	data.WriteString(`<label><input type="checkbox" name="email_have" ` + Setting_checked(values["email_have"]) + `> ` + lang("email_required") + `</label>` + Main_hr())
	data.WriteString(`<h3>` + lang("smtp_setting") + `</h3>`)
	data.WriteString(`<a href="https://support.google.com/mail/answer/7126229">(Google)</a>` + Main_hr())
	data.WriteString(`<a href="/setting/email_test">(` + lang("test") + `)</a>` + Main_hr())
	data.WriteString(`<span>` + lang("smtp_server") + `</span>` + Main_hr())
	data.WriteString(Setting_input("smtp_server", values["smtp_server"], "text") + Main_hr())
	data.WriteString(`<span>` + lang("smtp_port") + `</span>` + Main_hr())
	data.WriteString(Setting_input("smtp_port", values["smtp_port"], "text") + Main_hr())
	data.WriteString(`<span>` + lang("smtp_security") + `</span>` + Main_hr())
	data.WriteString(`<select name="smtp_security">` + Setting_options(values["smtp_security"], []string{"tls", "starttls", "plain"}, nil) + `</select>` + Main_hr())
	data.WriteString(`<span>` + lang("smtp_username") + `</span>` + Main_hr())
	data.WriteString(Setting_input("smtp_email", values["smtp_email"], "text") + Main_hr())
	data.WriteString(`<span>` + lang("smtp_password") + `</span>` + Main_hr())
	data.WriteString(Setting_input("smtp_pass", values["smtp_pass"], "password") + Main_hr())

	data.WriteString(`<h2>` + lang("oauth") + ` (` + lang("not_working") + `)</h2>`)
	data.WriteString(`<a href="https://developers.google.com/identity/protocols/oauth2">(Google)</a>` + Main_hr())
	data.WriteString(`<span>` + lang("oauth_client_id") + `</span>` + Main_hr())
	data.WriteString(Setting_input("oauth_client_id", values["oauth_client_id"], "text") + Main_hr())
	data.WriteString(`<button id="opennamu_save_button" type="submit">` + lang("save") + `</button></form>`)

	return Setting_page(db, config, lang("ext_api_req_set"), data.String(), "setting")
}
