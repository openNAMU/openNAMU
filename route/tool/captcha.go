package tool

import (
	"crypto/sha256"
	"database/sql"
	stdjson "encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	altcha "github.com/altcha-org/altcha-lib-go/v2"
	"github.com/gin-contrib/sessions"
)

var altcha_used_challenges = struct {
	sync.Mutex
	data map[string]time.Time
}{data: make(map[string]time.Time)}

func captcha_setting_value(db *sql.DB, name string) string {
	for _, value := range Get_setting(db, name, "") {
		if len(value) > 1 && value[1] == "" {
			return value[0]
		}
	}
	return ""
}

func captcha_setting(db *sql.DB) (string, string, string) {
	return captcha_setting_value(db, "recaptcha"), captcha_setting_value(db, "sec_re"), captcha_setting_value(db, "recaptcha_ver")
}

func captcha_altcha_cost(rec_ver string) (int, bool) {
	switch rec_ver {
	case "altcha_low":
		return 1000, true
	case "altcha_medium":
		return 5000, true
	case "altcha_high":
		return 10000, true
	}
	return 0, false
}

func captcha_check(db *sql.DB, session sessions.Session, ip string, response string) bool {
	if Check_acl(db, "", "", "recaptcha", ip) {
		return true
	}

	pub_key, sec_key, rec_ver := captcha_setting(db)
	if altcha_cost, ok := captcha_altcha_cost(rec_ver); ok {
		return captcha_check_altcha(response, sec_key, altcha_cost)
	}

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

func captcha_check_altcha(response string, sec_key string, altcha_cost int) bool {
	if response == "" || sec_key == "" || len(response) > 1<<20 {
		return false
	}

	payload := altcha.Payload{}
	if err := stdjson.Unmarshal([]byte(response), &payload); err != nil {
		return false
	}

	result, err := altcha.VerifySolution(altcha.VerifySolutionOptions{
		Challenge:           payload.Challenge,
		Solution:            payload.Solution,
		DeriveKey:           altcha.DeriveKeyPBKDF2(),
		HMACSignatureSecret: sec_key,
	})
	if err != nil || !result.Verified {
		return false
	}
	if payload.Challenge.Parameters.Cost != altcha_cost {
		return false
	}

	now := time.Now()
	expires_at := now.Add(20 * time.Minute)
	if payload.Challenge.Parameters.ExpiresAt > 0 {
		expires_at = time.Unix(payload.Challenge.Parameters.ExpiresAt, 0)
	}
	if !expires_at.After(now) {
		return false
	}

	challenge_data, err := stdjson.Marshal(payload.Challenge)
	if err != nil {
		return false
	}
	hash := sha256.Sum256(challenge_data)
	key := string(hash[:])

	altcha_used_challenges.Lock()
	defer altcha_used_challenges.Unlock()
	for used_key, used_until := range altcha_used_challenges.data {
		if !used_until.After(now) {
			delete(altcha_used_challenges.data, used_key)
		}
	}
	if _, exists := altcha_used_challenges.data[key]; exists {
		return false
	}
	altcha_used_challenges.data[key] = expires_at
	return true
}

func captcha_challenge(db *sql.DB) (altcha.Challenge, error) {
	_, sec_key, rec_ver := captcha_setting(db)
	altcha_cost, ok := captcha_altcha_cost(rec_ver)
	if !ok || sec_key == "" {
		return altcha.Challenge{}, errors.New("altcha is not enabled")
	}

	expires_at := time.Now().Add(20 * time.Minute)
	return altcha.CreateChallenge(altcha.CreateChallengeOptions{
		Algorithm:           "PBKDF2/SHA-256",
		Cost:                altcha_cost,
		DeriveKey:           altcha.DeriveKeyPBKDF2(),
		ExpiresAt:           &expires_at,
		HMACSignatureSecret: sec_key,
		KeyLength:           32,
	})
}

func captcha_response(response string, recaptcha_response string, hcaptcha_response string, turnstile_response string, altcha_response string) string {
	if response != "" {
		return response
	}
	if recaptcha_response != "" {
		return recaptcha_response
	}
	if hcaptcha_response != "" {
		return hcaptcha_response
	}
	if turnstile_response != "" {
		return turnstile_response
	}
	return altcha_response
}
