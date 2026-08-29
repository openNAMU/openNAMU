package tool

import (
	"database/sql"

	altcha "github.com/altcha-org/altcha-lib-go/v2"
	"github.com/gin-contrib/sessions"
)

func Captcha_check(db *sql.DB, session sessions.Session, ip string, response string) bool {
	return captcha_check(db, session, ip, response)
}

func Captcha_response(response string, recaptcha_response string, hcaptcha_response string, turnstile_response string, altcha_response ...string) string {
	altcha_value := ""
	if len(altcha_response) > 0 {
		altcha_value = altcha_response[0]
	}
	return captcha_response(response, recaptcha_response, hcaptcha_response, turnstile_response, altcha_value)
}

func Captcha_challenge(db *sql.DB) (altcha.Challenge, error) {
	return captcha_challenge(db)
}
