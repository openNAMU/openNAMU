package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register_bbs_vote_routes(r *gin.Engine) {
	r.POST("/bbs/w/:set_id/:set_code/tabom", func(c *gin.Context) {
		data := View_bbs_in_w_tabom_post(
			Make_route_config(c),
			c.Param("set_id"),
			c.Param("set_code"),
			c.PostForm("vote_type"),
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/bbs/w/:set_id/:set_code/comment_tabom/:comment_code/:vote_type", func(c *gin.Context) {
		data := View_bbs_in_w_comment_tabom(
			Make_route_config(c),
			c.Param("set_id"),
			c.Param("set_code"),
			c.Param("comment_code"),
			c.Param("vote_type"),
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.POST("/bbs/w/:set_id/:set_code/comment/:comment_code/tabom", func(c *gin.Context) {
		data := View_bbs_in_w_comment_tabom_post(
			Make_route_config(c),
			c.Param("set_id"),
			c.Param("set_code"),
			c.Param("comment_code"),
			c.PostForm("vote_type"),
		)
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
}
