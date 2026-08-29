package main

import (
	"net/http"
	"sync"
	"time"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

var captcha_challenge_last = struct {
	sync.Mutex
	data map[string]time.Time
}{data: make(map[string]time.Time)}

func captcha_challenge_allowed(ip string) bool {
	if ip == "" {
		ip = "unknown"
	}

	now := time.Now()
	captcha_challenge_last.Lock()
	defer captcha_challenge_last.Unlock()

	for key, last := range captcha_challenge_last.data {
		if now.Sub(last) >= time.Minute {
			delete(captcha_challenge_last.data, key)
		}
	}

	if last, ok := captcha_challenge_last.data[ip]; ok && now.Sub(last) < time.Second {
		return false
	}
	captcha_challenge_last.data[ip] = now
	return true
}

func register_captcha_routes(r *gin.Engine) {
	r.GET("/api/altcha/challenge", func(c *gin.Context) {
		if !captcha_challenge_allowed(c.ClientIP()) {
			c.Header("Retry-After", "1")
			c.JSON(http.StatusTooManyRequests, map[string]string{"error": "too many requests"})
			return
		}
		c.Header("Cache-Control", "no-store")
		challenge, err := route.Api_captcha_challenge()
		if err != nil {
			c.JSON(http.StatusNotFound, map[string]string{"error": "captcha disabled"})
			return
		}
		c.JSON(http.StatusOK, challenge)
	})
}
