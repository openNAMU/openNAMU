package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register_calendar_routes(r *gin.Engine) {
	r.GET("/calendar", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(View_calendar_list(Make_route_config(c))))
	})

	r.GET("/calendar/add", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(View_calendar_add(Make_route_config(c), nil)))
	})

	r.POST("/calendar/add", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(View_calendar_add(Make_route_config(c), c.Request.PostForm)))
	})

	r.GET("/calendar/:id", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(View_calendar(Make_route_config(c), c.Param("id"), "")))
	})

	r.GET("/calendar/:id/edit", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(View_calendar_edit(Make_route_config(c), c.Param("id"), nil)))
	})

	r.POST("/calendar/:id/edit", func(c *gin.Context) {
		_ = c.Request.ParseForm()
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(View_calendar_edit(Make_route_config(c), c.Param("id"), c.Request.PostForm)))
	})

	r.GET("/calendar/:id/month/:month", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(View_calendar(Make_route_config(c), c.Param("id"), c.Param("month"))))
	})

	r.POST("/calendar/:id/event/add", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(View_calendar_event_add_post(
			Make_route_config(c),
			c.Param("id"),
			c.PostForm("date"),
			c.PostForm("title"),
			c.PostForm("memo"),
		)))
	})

	r.GET("/calendar/:id/event/:event", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(View_calendar_event(Make_route_config(c), c.Param("id"), c.Param("event"))))
	})

	r.POST("/calendar/:id/event/:event", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(View_calendar_event_edit_post(
			Make_route_config(c),
			c.Param("id"),
			c.Param("event"),
			c.PostForm("date"),
			c.PostForm("title"),
			c.PostForm("memo"),
		)))
	})

	r.GET("/calendar/:id/event/:event/delete", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(View_calendar_event_delete(Make_route_config(c), c.Param("id"), c.Param("event"), false)))
	})

	r.POST("/calendar/:id/event/:event/delete", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(View_calendar_event_delete(Make_route_config(c), c.Param("id"), c.Param("event"), true)))
	})

	r.GET("/calendar/:id/delete", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(View_calendar_delete(Make_route_config(c), c.Param("id"), false)))
	})

	r.POST("/calendar/:id/delete", func(c *gin.Context) {
		Write_data(c, http.StatusOK, "text/html; charset=utf-8", []byte(View_calendar_delete(Make_route_config(c), c.Param("id"), true)))
	})
}
