package tool

import (
	"crypto/sha256"
	"crypto/sha512"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

const session_cookie_name = "opennamu_session"

func Session_middleware() gin.HandlerFunc {
	secret := os.Getenv("NAMU_SESSION_KEY")
	if secret == "" {
		db := DB_connect()
		QueryRow_DB(db, `select data from other where name = "session_key"`, []any{&secret})
		DB_close(db)
	}
	if secret == "" {
		secret = Get_random_key(128)
	}
	auth_key := sha512.Sum512([]byte(secret))
	encryption_key := sha256.Sum256([]byte("opennamu-session:" + secret))

	store := memstore.NewStore(auth_key[:], encryption_key[:])
	store.Options(Session_options(false))

	return sessions.Sessions(session_cookie_name, store)
}

func Get_session(c *gin.Context) sessions.Session {
	session := sessions.Default(c)
	session.Options(Session_options(Is_https(c)))

	return session
}

func Session_options(secure bool) sessions.Options {
	return sessions.Options{
		Path:     "/",
		MaxAge:   0,
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func Is_https(c *gin.Context) bool {
	return c.Request.TLS != nil || strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https")
}
