package tool

import (
	"crypto/sha256"
	"encoding/hex"
	"html/template"
	"net/url"
)

func Sha224(data string) string {
	hasher := sha256.New224()
	hasher.Write([]byte(data))
	hash_byte := hasher.Sum(nil)
	hash_str := hex.EncodeToString(hash_byte)

	return hash_str
}

func Url_parser(data string) string {
	return url.QueryEscape(data)
}

func HTML_escape(data string) string {
	return template.HTMLEscapeString(data)
}
