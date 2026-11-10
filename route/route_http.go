package route

import (
	"net/http"

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

func Write_data(c *gin.Context, status int, content_type string, data []byte) {
	target, ok := tool.Get_redirect_target(string(data))
	if ok {
		c.Redirect(http.StatusFound, target)
		return
	}

	c.Data(status, content_type, data)
}

func Site_view_middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := tool.DB_connect()
		defer tool.DB_close(db)

		user_id := tool.Get_IP(c)
		if !tool.Get_auth_info(db, user_id)["site_view"] {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		tool.Update_online_user(user_id)

		c.Next()
	}
}
