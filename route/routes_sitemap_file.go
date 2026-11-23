package route

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func Register_sitemap_file_routes(r *gin.Engine) {
	r.GET("/sitemap.xml", func(c *gin.Context) {
		data, ok := Read_sitemap_file("sitemap.xml")
		if !ok {
			c.Status(http.StatusNotFound)
			return
		}
		Write_data(c, http.StatusOK, "application/xml; charset=utf-8", []byte(data))
	})

	// gin v1.12 interprets a param name up to the end of the segment, so the
	// pattern "/sitemap_:num.xml" makes the param name "num.xml" and
	// c.Param("num") always stays empty. The pattern must end with ":num" and
	// the ".xml" suffix is validated in the handler instead.
	r.GET("/sitemap_:num", func(c *gin.Context) {
		num_data := c.Param("num")
		if !strings.HasSuffix(num_data, ".xml") {
			c.Status(http.StatusNotFound)
			return
		}

		num, err := strconv.Atoi(strings.TrimSuffix(num_data, ".xml"))
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		name := "sitemap_" + strconv.Itoa(num) + ".xml"
		data, ok := Read_sitemap_file(name)
		if !ok {
			c.Status(http.StatusNotFound)
			return
		}
		Write_data(c, http.StatusOK, "application/xml; charset=utf-8", []byte(data))
	})
}
