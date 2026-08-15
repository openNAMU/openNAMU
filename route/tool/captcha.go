package tool

import (
	"database/sql"
	stdjson "encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
)

func captcha_check(db *sql.DB, session sessions.Session, ip string, response string) bool {
	if Check_acl(db, "", "", "recaptcha", ip) {
		return true
	}

	pub_key := ""
	sec_key := ""
	rec_ver := ""
	QueryRow_DB(db, `select data from other where name = "recaptcha"`, []any{&pub_key})
	QueryRow_DB(db, `select data from other where name = "sec_re"`, []any{&sec_key})
	QueryRow_DB(db, `select data from other where name = "recaptcha_ver"`, []any{&rec_ver})

	if pub_key == "" || sec_key == "" {
		return true
	}
	if response == "" {
		return false
	}

	verify_url := "https://www.google.com/recaptcha/api/siteverify"
	if rec_ver == "cf" {
		verify_url = "https://challenges.cloudflare.com/turnstile/v0/siteverify"
	} else if rec_ver == "h" {
		verify_url = "https://hcaptcha.com/siteverify"
	}

	form := url.Values{}
	form.Set("secret", sec_key)
	form.Set("response", response)
	req, err := http.NewRequest(http.MethodPost, verify_url, strings.NewReader(form.Encode()))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 8 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return false
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return false
	}

	result := struct {
		Success bool `json:"success"`
	}{}
	if err := stdjson.NewDecoder(res.Body).Decode(&result); err != nil {
		return false
	}

	return result.Success
}

func captcha_response(response string, recaptcha_response string, hcaptcha_response string, turnstile_response string) string {
	if response != "" {
		return response
	}
	if recaptcha_response != "" {
		return recaptcha_response
	}
	if hcaptcha_response != "" {
		return hcaptcha_response
	}
	return turnstile_response
}
