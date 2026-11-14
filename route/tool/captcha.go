package tool

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	stdjson "encoding/json"
	"errors"
	"fmt"
	"log"
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

func Captcha_setting_value(db *sql.DB, name string) string {
	for _, value := range Get_setting(db, name, "") {
		if len(value) > 1 && value[1] == "" {
			return value[0]
		}
	}
	return ""
}

func Captcha_setting(db *sql.DB) (string, string, string, string) {
	return Captcha_setting_value(db, "recaptcha"), Captcha_setting_value(db, "sec_re"), Captcha_setting_value(db, "altcha_sec_re"), Captcha_setting_value(db, "recaptcha_ver")
}

type Captcha_config struct {
	Provider          string
	public_key        string
	secret_key        string
	altcha_secret_key string
}

func Captcha_setting_value_read(db *sql.DB, name string) (string, bool, error) {
	if db == nil {
		return "", false, errors.New("database is required")
	}

	rows, err := db.Query(DB_change("select coalesce(data, '') from other where name = ? and coverage = ''"), name)
	if err != nil {
		return "", false, fmt.Errorf("read CAPTCHA setting %s: %w", name, err)
	}
	defer rows.Close()

	value := ""
	value_found := false
	row_found := false
	for rows.Next() {
		var current_value string
		if err := rows.Scan(&current_value); err != nil {
			return "", false, fmt.Errorf("read CAPTCHA setting %s: %w", name, err)
		}
		row_found = true
		if name == "recaptcha_ver" {
			current_value = strings.TrimSpace(current_value)
		}
		if strings.TrimSpace(current_value) == "" {
			continue
		}
		if value_found && value != current_value {
			return "", row_found, fmt.Errorf("conflicting CAPTCHA setting %s", name)
		}
		value = current_value
		value_found = true
	}
	if err := rows.Err(); err != nil {
		return "", row_found, fmt.Errorf("read CAPTCHA setting %s: %w", name, err)
	}

	return value, row_found, nil
}

func Captcha_setting_value_unique(db *sql.DB, name string) (string, error) {
	value, _, err := Captcha_setting_value_read(db, name)
	return value, err
}

func Captcha_settings_read(db *sql.DB) (map[string]string, map[string]bool, error) {
	if db == nil {
		return nil, nil, errors.New("database is required")
	}

	values := make(map[string]string, 4)
	value_found := make(map[string]bool, 4)
	found := make(map[string]bool, 4)
	rows, err := db.Query(DB_change("select name, coalesce(data, '') from other where coverage = '' and name in (?, ?, ?, ?)"), "recaptcha", "sec_re", "altcha_sec_re", "recaptcha_ver")
	if err != nil {
		return nil, nil, fmt.Errorf("read global CAPTCHA settings: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var current_value string
		if err := rows.Scan(&name, &current_value); err != nil {
			return nil, nil, fmt.Errorf("read global CAPTCHA settings: %w", err)
		}
		found[name] = true
		if name == "recaptcha_ver" {
			current_value = strings.TrimSpace(current_value)
		}
		if strings.TrimSpace(current_value) == "" {
			continue
		}
		if value_found[name] && values[name] != current_value {
			return nil, nil, fmt.Errorf("conflicting CAPTCHA setting %s", name)
		}
		values[name] = current_value
		value_found[name] = true
	}
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("read global CAPTCHA settings: %w", err)
	}

	for _, name := range []string{"recaptcha", "sec_re", "altcha_sec_re", "recaptcha_ver"} {
		if _, exists := values[name]; !exists {
			values[name] = ""
		}
	}
	return values, found, nil
}

func Captcha_provider_normalize(provider string) (string, bool) {
	provider = strings.TrimSpace(provider)
	switch provider {
	case "", "v2", "v3", "h", "cf", "altcha_low", "altcha_medium", "altcha_high":
		return provider, true
	default:
		return "", false
	}
}

func Captcha_config_validate(provider string, public_key string, secret_key string, altcha_secret_key string) error {
	provider, supported := Captcha_provider_normalize(provider)
	if !supported || provider == "" {
		return errors.New("unsupported CAPTCHA provider")
	}

	if _, altcha_enabled := Captcha_altcha_cost(provider); altcha_enabled {
		if strings.TrimSpace(altcha_secret_key) == "" {
			return errors.New("ALTCHA secret key is required")
		}
		return nil
	}

	if strings.TrimSpace(public_key) == "" || strings.TrimSpace(secret_key) == "" {
		return errors.New("CAPTCHA provider keys are required")
	}
	return nil
}

func Captcha_config_resolve(db *sql.DB) (Captcha_config, error) {
	values, _, err := Captcha_settings_read(db)
	if err != nil {
		return Captcha_config{}, err
	}

	config := Captcha_config{
		public_key:        values["recaptcha"],
		secret_key:        values["sec_re"],
		altcha_secret_key: values["altcha_sec_re"],
	}
	provider, supported := Captcha_provider_normalize(values["recaptcha_ver"])
	if !supported {
		return config, errors.New("unsupported CAPTCHA provider")
	}
	if provider == "" {
		if strings.TrimSpace(config.public_key) != "" || strings.TrimSpace(config.secret_key) != "" {
			provider = "v2"
		} else if strings.TrimSpace(config.altcha_secret_key) != "" {
			provider = "altcha_high"
		}
	}
	config.Provider = provider
	if provider == "" {
		return config, errors.New("no CAPTCHA provider is configured")
	}
	if err := Captcha_config_validate(provider, config.public_key, config.secret_key, config.altcha_secret_key); err != nil {
		return config, err
	}
	return config, nil
}

func Captcha_altcha_cost(rec_ver string) (int, bool) {
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

func Captcha_check_internal(db *sql.DB, session sessions.Session, ip string, response string) bool {
	if Check_acl(db, "", "", "recaptcha", ip) {
		return true
	}

	config, err := Captcha_config_resolve(db)
	if err != nil {
		log.Printf("[CAPTCHA] invalid configuration; review the selected provider and its required keys in External settings")
		return false
	}
	if config.Provider == "" {
		return true
	}
	if altcha_cost, ok := Captcha_altcha_cost(config.Provider); ok {
		return Captcha_check_altcha(response, config.altcha_secret_key, altcha_cost)
	}
	if response == "" {
		return false
	}

	verify_url := "https://www.google.com/recaptcha/api/siteverify"
	switch config.Provider {
	case "v2", "v3":
	case "h":
		verify_url = "https://hcaptcha.com/siteverify"
	case "cf":
		verify_url = "https://challenges.cloudflare.com/turnstile/v0/siteverify"
	default:
		return false
	}

	form := url.Values{}
	form.Set("secret", config.secret_key)
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

func Captcha_check_altcha(response string, sec_key string, altcha_cost int) bool {
	if response == "" || sec_key == "" || len(response) > 1<<20 {
		return false
	}

	payload_data, err := base64.StdEncoding.DecodeString(response)
	if err != nil {
		return false
	}

	payload := altcha.Payload{}
	if err := stdjson.Unmarshal(payload_data, &payload); err != nil {
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

func Captcha_challenge_internal(db *sql.DB) (altcha.Challenge, error) {
	config, err := Captcha_config_resolve(db)
	if err != nil {
		return altcha.Challenge{}, errors.New("invalid CAPTCHA configuration")
	}
	altcha_cost, ok := Captcha_altcha_cost(config.Provider)
	if !ok {
		return altcha.Challenge{}, errors.New("ALTCHA is not enabled")
	}

	expires_at := time.Now().Add(20 * time.Minute)
	return altcha.CreateChallenge(altcha.CreateChallengeOptions{
		Algorithm:           "PBKDF2/SHA-256",
		Cost:                altcha_cost,
		DeriveKey:           altcha.DeriveKeyPBKDF2(),
		ExpiresAt:           &expires_at,
		HMACSignatureSecret: config.altcha_secret_key,
		KeyLength:           32,
	})
}

func Captcha_response_internal(response string, recaptcha_response string, hcaptcha_response string, turnstile_response string, altcha_response string) string {
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
