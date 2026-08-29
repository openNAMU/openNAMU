package route

import (
	"opennamu/route/tool"

	altcha "github.com/altcha-org/altcha-lib-go/v2"
)

func Api_captcha_challenge() (altcha.Challenge, error) {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return tool.Captcha_challenge(db)
}
