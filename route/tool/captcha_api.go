package tool

import (
	"database/sql"

	"github.com/gin-contrib/sessions"
)

func Captcha_check(db *sql.DB, session sessions.Session, ip string, response string) bool {
	return captcha_check(db, session, ip, response)
}

func Captcha_response(response string, recaptcha_response string, hcaptcha_response string, turnstile_response string) string {
	return captcha_response(response, recaptcha_response, hcaptcha_response, turnstile_response)
}
