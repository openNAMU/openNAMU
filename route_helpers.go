package main

import (
	"net/http"
	"strings"

	"opennamu/route/tool"

	"github.com/gin-gonic/gin"
)

func Make_route_config(c *gin.Context) tool.Config {
	return tool.Config{
		IP:        tool.Get_IP(c),
		Cookies:   tool.Get_Cookies(c),
		Session:   tool.Get_session(c),
		UserAgent: c.GetHeader("User-Agent"),
	}
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
	target, ok := tool.Get_redirect_target(string(data))
	if ok {
		c.Redirect(http.StatusFound, target)
		return
	}

	c.Data(status, content_type, data)
}
