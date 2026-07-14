package tool

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

const session_cookie_name = "opennamu_session"

func Session_middleware() gin.HandlerFunc {
	auth_key := make([]byte, 64)
	encryption_key := make([]byte, 32)

	if _, err := rand.Read(auth_key); err != nil {
		panic(fmt.Errorf("session authentication key generation failed: %w", err))
	}
	if _, err := rand.Read(encryption_key); err != nil {
		panic(fmt.Errorf("session encryption key generation failed: %w", err))
	}

	store := memstore.NewStore(auth_key, encryption_key)
	store.Options(session_options(false))

	return sessions.Sessions(session_cookie_name, store)
}

func Get_session(c *gin.Context) sessions.Session {
	session := sessions.Default(c)
	session.Options(session_options(is_https(c)))

	return session
}

func session_options(secure bool) sessions.Options {
	return sessions.Options{
		Path:     "/",
		MaxAge:   0,
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func is_https(c *gin.Context) bool {
	return c.Request.TLS != nil || strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https")
}
