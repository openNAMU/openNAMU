package route

import (
	"bytes"
	"database/sql"
	stdjson "encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"opennamu/route/tool"
)

type indexnow_request struct {
	Host        string   `json:"host"`
	Key         string   `json:"key"`
	KeyLocation string   `json:"keyLocation"`
	URLList     []string `json:"urlList"`
}

func indexnow_key_valid(key string) bool {
	if key == "" {
		return false
	}

	for _, value := range key {
		if (value < 'a' || value > 'z') &&
			(value < 'A' || value > 'Z') &&
			(value < '0' || value > '9') && value != '-' && value != '_' {
			return false
		}
	}

	return true
}

func save_indexnow_key(key string) bool {
	key = strings.TrimSpace(key)
	if !indexnow_key_valid(key) {
		return false
	}

	if err := os.WriteFile(key+".txt", []byte(key), 0644); err != nil {
		log.Printf("indexnow key file save failed: %v", err)
		return false
	}

	return true
}

func sync_indexnow_key(db *sql.DB) {
	key := tool.Get_setting_value(db, "indexnow_key", "", "")
	if key != "" {
		save_indexnow_key(key)
	}
}

func send_indexnow(host string, key string, key_location string, page_url string) {
	request_data, err := stdjson.Marshal(indexnow_request{
		Host:        host,
		Key:         key,
		KeyLocation: key_location,
		URLList:     []string{page_url},
	})
	if err != nil {
		return
	}

	request, err := http.NewRequest(
		http.MethodPost,
		"https://api.indexnow.org/indexnow",
		bytes.NewReader(request_data),
	)
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := http.Client{Timeout: 2 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		log.Printf("indexnow request failed: %v", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusAccepted {
		log.Printf("indexnow request failed: status %d", response.StatusCode)
	}
}

func notify_indexnow(db *sql.DB, doc_name string) {
	key := strings.TrimSpace(tool.Get_setting_value(db, "indexnow_key", "", ""))
	if !indexnow_key_valid(key) || !save_indexnow_key(key) {
		return
	}

	domain := strings.TrimRight(tool.Get_domain(db, true), "/")
	parsed_domain, err := url.Parse(domain)
	if err != nil || parsed_domain.Host == "" {
		return
	}

	page_url := domain + "/w/" + tool.Url_parser(doc_name)
	key_location := domain + "/" + key + ".txt"
	go send_indexnow(parsed_domain.Host, key, key_location, page_url)
}
