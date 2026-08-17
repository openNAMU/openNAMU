package route

import (
	"opennamu/route/tool"
	"os"
	"strings"
)

func View_robots_txt(config tool.Config, request_host string) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	robot_default := setting_value(db, "robot_default", "", "")
	robot := setting_value(db, "robot", "", "")
	if robot_default == "" && robot != "" {
		return robot
	}

	domain := tool.Get_domain(db, true)
	if domain == "" || domain == "http://" || domain == "https://" {
		domain = "http://" + request_host
	}

	data := strings.Builder{}
	data.WriteString("User-agent: *\n")
	data.WriteString("Disallow: /\n")
	data.WriteString("Allow: /$\n")
	data.WriteString("Allow: /w/\n")
	data.WriteString("Allow: /bbs/w/\n")
	data.WriteString("Allow: /sitemap.xml$\n")
	data.WriteString("Allow: /sitemap_*.xml$")
	if _, err := os.Stat("sitemap.xml"); err == nil {
		data.WriteString("\nSitemap: " + domain + "/sitemap.xml")
	}

	return data.String()
}
