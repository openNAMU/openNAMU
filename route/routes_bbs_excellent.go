package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register_bbs_excellent_routes(r *gin.Engine) {
	r.GET("/bbs/excellent", func(c *gin.Context) {
		data := View_bbs_excellent(Make_route_config(c), "1")
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})

	r.GET("/bbs/excellent/page/:page", func(c *gin.Context) {
		data := View_bbs_excellent(Make_route_config(c), c.Param("page"))
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(data))
	})
}

func Register_bbs_excellent_api_routes(r *gin.Engine) {
	r.GET("/api/v2/bbs/excellent/:page", func(c *gin.Context) {
		c.JSON(http.StatusOK, Api_bbs(Make_route_config(c), "", c.Param("page"), "excellent"))
	})

	r.POST("/api/v2/bbs/excellent_exclusion/:set_id/:set_code/:state", func(c *gin.Context) {
		excluded := false
		switch c.Param("state") {
		case "exclude":
			excluded = true
		case "restore":
		default:
			c.JSON(http.StatusOK, map[string]any{"response": "error", "data": "state"})
			return
		}
		c.JSON(http.StatusOK, Api_bbs_excellent_exclusion_set(Make_route_config(c), c.Param("set_id"), c.Param("set_code"), excluded))
	})
}
