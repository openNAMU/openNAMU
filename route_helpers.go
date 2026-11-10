package main

import (
	"strings"

	"opennamu/route"
	"opennamu/route/tool"

	"github.com/gin-gonic/gin"
)

func Make_route_config(c *gin.Context) tool.Config {
	return route.Make_route_config(c)
}

func Captcha_response_internal(c *gin.Context) string {
	return tool.Captcha_response(c.PostForm("g-recaptcha"), c.PostForm("g-recaptcha-response"), c.PostForm("h-captcha-response"), c.PostForm("cf-turnstile-response"), c.PostForm("altcha"))
}

func Parse_bbs_code(value string) (string, string) {
	if strings.HasPrefix(value, "-1-") {
		return "-1", strings.TrimPrefix(value, "-1-")
	}
	parts := strings.SplitN(value, "-", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

func Write_data(c *gin.Context, status int, content_type string, data []byte) {
	route.Write_data(c, status, content_type, data)
}
