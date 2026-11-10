package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register_bbs_post_tool_routes(r *gin.Engine) {
	r.GET("/bbs/tool/:set_id/:set_code", func(c *gin.Context) {
		data := View_bbs_in_w_tool(Make_route_config(c), c.Param("set_id"), c.Param("set_code"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.POST("/bbs/tool/:set_id/:set_code", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		data := View_bbs_in_w_tool_post(
			Make_route_config(c),
			c.Param("set_id"),
			c.Param("set_code"),
			c.Request.PostForm,
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
}
